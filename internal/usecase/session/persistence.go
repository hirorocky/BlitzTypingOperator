// Package game_state はゲーム全体の状態管理を提供するユースケースです。
// このファイルはセーブ/ロードの変換ロジックを担当します。
package session

import (
	"log/slog"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/usecase/achievement"
	"hirorocky/type-battle/internal/usecase/rewarding"
	"hirorocky/type-battle/internal/usecase/spawning"
	"hirorocky/type-battle/internal/usecase/synthesize"
)

// DomainDataSources はセーブデータ復元時に使用するドメイン型データソースです。
type DomainDataSources struct {
	CoreTypes              []domain.CoreType
	ModuleTypes            []rewarding.ModuleDropInfo
	EnemyTypes             []domain.EnemyType
	PassiveSkills          map[string]domain.PassiveSkill
	ChainEffectDefinitions []rewarding.ChainEffectDefinition
}

// ToSaveData はGameStateをセーブデータに変換します。
// v1.0.0形式: TypeIDとLevel（コア）、TypeIDとChainEffect（モジュール）を保存します。
func (g *GameState) ToSaveData() *savedata.SaveData {
	saveData := savedata.NewSaveData()

	// 最高到達レベル
	saveData.Statistics.MaxLevelReached = g.MaxLevelReached

	// コアをv1.0.0形式で保存（IDなし）
	coreInstances := make([]savedata.CoreInstanceSave, 0)
	for _, core := range g.inventory.GetCores() {
		coreInstances = append(coreInstances, savedata.CoreInstanceSave{
			CoreTypeID: core.TypeID,
			Level:      core.Level,
		})
	}
	saveData.Inventory.CoreInstances = coreInstances

	// モジュールをModuleInstancesとして保存（チェイン効果対応）
	moduleInstances := make([]savedata.ModuleInstanceSave, 0)
	for _, module := range g.inventory.GetModules() {
		modSave := savedata.ModuleInstanceSave{
			TypeID: module.TypeID,
		}
		if module.ChainEffect != nil {
			modSave.ChainEffect = &savedata.ChainEffectSave{
				Type:  string(module.ChainEffect.Type),
				Value: module.ChainEffect.Value,
			}
		}
		moduleInstances = append(moduleInstances, modSave)
	}
	saveData.Inventory.ModuleInstances = moduleInstances

	// エージェントを保存（コア情報を直接埋め込み、スキルはオブジェクト配列）
	agentInstances := make([]savedata.AgentInstanceSave, 0)
	for _, ag := range g.agentManager.GetAgents() {
		skills := make([]savedata.SkillInstanceSave, len(ag.Modules))
		for i, m := range ag.Modules {
			skills[i] = savedata.SkillInstanceSave{
				TypeID: m.TypeID,
			}
			if m.ChainEffect != nil {
				skills[i].ChainEffect = &savedata.ChainEffectSave{
					Type:  string(m.ChainEffect.Type),
					Value: m.ChainEffect.Value,
				}
			}
		}
		agentInstances = append(agentInstances, savedata.AgentInstanceSave{
			ID: ag.ID,
			Core: savedata.CoreInstanceSave{
				CoreTypeID: ag.Core.TypeID,
				Level:      ag.Core.Level,
			},
			Skills: skills,
		})
	}
	saveData.Inventory.AgentInstances = agentInstances

	saveData.Inventory.MaxCoreSlots = g.inventory.Cores().MaxSlots()
	saveData.Inventory.MaxModuleSlots = g.inventory.Modules().MaxSlots()

	// 装備中のエージェントIDをスロット番号順に取得
	var equippedIDs [synthesize.MaxEquipmentSlots]string
	for slot := 0; slot < synthesize.MaxEquipmentSlots; slot++ {
		if equippedAgent := g.agentManager.GetEquippedAgentAt(slot); equippedAgent != nil {
			equippedIDs[slot] = equippedAgent.ID
		}
		// nilの場合は空文字列のまま
	}
	saveData.Player.EquippedAgentIDs = equippedIDs

	// 統計
	stats := g.statistics
	saveData.Statistics.TotalBattles = stats.Battle().TotalBattles
	saveData.Statistics.Victories = stats.Battle().Wins
	saveData.Statistics.Defeats = stats.Battle().Losses
	saveData.Statistics.HighestWPM = float64(stats.Typing().MaxWPM)
	saveData.Statistics.AverageWPM = stats.GetAverageWPM()
	saveData.Statistics.PerfectAccuracyCount = stats.Typing().PerfectAccuracyCount
	saveData.Statistics.TotalCharactersTyped = stats.Typing().TotalCharacters
	saveData.Statistics.EncounteredEnemies = g.encounteredEnemies

	// 実績（ドメイン型を経由してセーブデータ型に変換）
	saveData.Achievements = savedata.AchievementStateToSaveData(g.achievements.GetUnlockedIDs())

	// 設定
	saveData.Settings.KeyBindings = g.settings.Keybinds()

	// 撃破済み敵情報を保存
	saveData.Statistics.DefeatedEnemies = g.GetDefeatedEnemies()

	return saveData
}

// GameStateFromSaveData はセーブデータからGameStateを生成します。
// v1.0.0形式のセーブデータからオブジェクトを再構築します。
// sourcesにはマスタデータから変換されたドメイン型データを渡す必要があります。
func GameStateFromSaveData(data *savedata.SaveData, sources *DomainDataSources) *GameState {
	if sources == nil {
		slog.Error("マスタデータソースが必要です")
		return nil
	}

	// マスタデータを取得（ドメイン型）
	coreTypes := sources.CoreTypes
	moduleTypes := sources.ModuleTypes
	passiveSkills := sources.PassiveSkills
	enemyTypes := sources.EnemyTypes
	chainEffectDefs := sources.ChainEffectDefinitions

	// インベントリマネージャーを作成
	invManager := NewInventoryManager()

	// セーブデータからコアを再構築（v1.0.0形式: TypeIDとLevelのみ）
	if data.Inventory != nil {
		for _, coreSave := range data.Inventory.CoreInstances {
			// コア特性を検索（ドメイン型）
			coreType := findCoreType(coreTypes, coreSave.CoreTypeID)
			passiveSkill := findPassiveSkill(passiveSkills, coreType.PassiveSkillID)

			// コアを再構築（v1.0.0形式: TypeIDベース）
			core := domain.NewCoreWithTypeID(
				coreSave.CoreTypeID,
				coreSave.Level,
				coreType,
				passiveSkill,
			)
			if err := invManager.AddCore(core); err != nil {
				slog.Error("コア追加に失敗",
					slog.String("core_type_id", core.TypeID),
					slog.String("core_type", core.Type.ID),
					slog.Any("error", err),
				)
			}
		}

		// モジュールを再構築（v1.0.0形式: ModuleInstances）
		for _, modSave := range data.Inventory.ModuleInstances {
			moduleDropInfo := findModuleDropInfo(moduleTypes, modSave.TypeID)
			if moduleDropInfo != nil {
				// チェイン効果を復元
				var chainEffect *domain.ChainEffect
				if modSave.ChainEffect != nil {
					effectType := domain.ChainEffectType(modSave.ChainEffect.Type)
					chainEffectDef := findChainEffectDefinition(chainEffectDefs, effectType)
					if chainEffectDef != nil {
						ce := domain.NewChainEffectWithTemplate(
							effectType,
							modSave.ChainEffect.Value,
							chainEffectDef.Description,
							chainEffectDef.ShortDescription,
						)
						chainEffect = &ce
					}
				}
				module := moduleDropInfo.ToDomainWithChainEffect(chainEffect)
				if err := invManager.AddModule(module); err != nil {
					slog.Error("モジュール追加に失敗",
						slog.String("module_type_id", module.TypeID),
						slog.String("module_name", module.Name()),
						slog.Any("error", err),
					)
				}
			}
		}

		// 後方互換性: 旧形式ModuleCountsからの復元
		for moduleID, count := range data.Inventory.ModuleCounts {
			moduleDropInfo := findModuleDropInfo(moduleTypes, moduleID)
			if moduleDropInfo != nil {
				for i := 0; i < count; i++ {
					module := moduleDropInfo.ToDomain()
					if err := invManager.AddModule(module); err != nil {
						slog.Error("モジュール追加に失敗（旧形式）",
							slog.String("module_type_id", module.TypeID),
							slog.String("module_name", module.Name()),
							slog.Any("error", err),
						)
					}
				}
			}
		}
	}

	// エージェントマネージャーを作成
	agentMgr := synthesize.NewAgentManager(
		invManager.Cores(),
		invManager.Modules(),
	)

	// セーブデータからエージェントを再構築（コア情報は各エージェントに埋め込まれている）
	if data.Inventory != nil {
		for _, agentSave := range data.Inventory.AgentInstances {
			// エージェント内のコア情報からコアを再構築（v1.0.0形式）
			coreType := findCoreType(coreTypes, agentSave.Core.CoreTypeID)
			passiveSkill := findPassiveSkill(passiveSkills, coreType.PassiveSkillID)
			core := domain.NewCoreWithTypeID(
				agentSave.Core.CoreTypeID,
				agentSave.Core.Level,
				coreType,
				passiveSkill,
			)

			// スキルを再構築（オブジェクト配列形式）
			modules := make([]*domain.SkillModel, 0, len(agentSave.Skills))
			for _, skillSave := range agentSave.Skills {
				moduleDropInfo := findModuleDropInfo(moduleTypes, skillSave.TypeID)
				if moduleDropInfo != nil {
					// チェイン効果を復元
					var chainEffect *domain.ChainEffect
					if skillSave.ChainEffect != nil {
						effectType := domain.ChainEffectType(skillSave.ChainEffect.Type)
						chainEffectDef := findChainEffectDefinition(chainEffectDefs, effectType)
						if chainEffectDef != nil {
							ce := domain.NewChainEffectWithTemplate(
								effectType,
								skillSave.ChainEffect.Value,
								chainEffectDef.Description,
								chainEffectDef.ShortDescription,
							)
							chainEffect = &ce
						}
					}
					modules = append(modules, moduleDropInfo.ToDomainWithChainEffect(chainEffect))
				}
			}

			// エージェントを再構築
			agentModel := domain.NewAgent(agentSave.ID, core, modules)
			if err := agentMgr.AddAgent(agentModel); err != nil {
				slog.Error("エージェント追加に失敗",
					slog.String("agent_id", agentModel.ID),
					slog.Any("error", err),
				)
			}
		}
	}

	// 装備エージェントを復元（スロット番号を保持して復元）
	player := domain.NewPlayer()
	if data.Player != nil {
		for slot, agentID := range data.Player.EquippedAgentIDs {
			if agentID != "" {
				if err := agentMgr.EquipAgent(slot, agentID, player); err != nil {
					slog.Error("エージェント装備に失敗",
						slog.Int("slot", slot),
						slog.String("agent_id", agentID),
						slog.Any("error", err),
					)
				}
			}
		}
	}

	// 実績マネージャーを作成（セーブデータ型からドメイン型に変換してロード）
	achievementMgr := achievement.NewAchievementManager()
	if data.Achievements != nil {
		unlockedIDs := savedata.SaveDataToAchievementState(data.Achievements)
		achievementMgr.LoadFromUnlockedIDs(unlockedIDs)
	}

	// 統計マネージャーを作成して復元
	statsMgr := NewStatisticsManager()
	if data.Statistics != nil {
		statsSaveData := &StatisticsSaveData{
			TotalBattles:         data.Statistics.TotalBattles,
			Victories:            data.Statistics.Victories,
			Defeats:              data.Statistics.Defeats,
			MaxLevelReached:      data.Statistics.MaxLevelReached,
			HighestWPM:           data.Statistics.HighestWPM,
			AverageWPM:           data.Statistics.AverageWPM,
			PerfectAccuracyCount: data.Statistics.PerfectAccuracyCount,
			TotalCharactersTyped: data.Statistics.TotalCharactersTyped,
		}
		statsMgr.LoadFromSaveData(statsSaveData)
	}

	// 設定を復元
	settings := NewSettings()
	if data.Settings != nil && data.Settings.KeyBindings != nil {
		for action, key := range data.Settings.KeyBindings {
			settings.SetKeybind(action, key)
		}
	}

	// RewardCalculatorを作成
	rewardCalc := rewarding.NewRewardCalculator(coreTypes, moduleTypes, passiveSkills)

	// チェイン効果プールを設定
	if len(sources.ChainEffectDefinitions) > 0 {
		chainEffectPool := rewarding.NewChainEffectPool(sources.ChainEffectDefinitions)
		rewardCalc.SetChainEffectPool(chainEffectPool)
	}

	// EnemyGeneratorを作成
	enemyGen := spawning.NewEnemyGenerator(enemyTypes)

	// 最高到達レベル、エンカウント敵リスト、撃破済み敵情報を取得
	maxLevelReached := 0
	var encounteredEnemies []string
	var defeatedEnemies map[string]int
	if data.Statistics != nil {
		maxLevelReached = data.Statistics.MaxLevelReached
		encounteredEnemies = data.Statistics.EncounteredEnemies
		defeatedEnemies = data.Statistics.DefeatedEnemies
	}

	gs := &GameState{
		MaxLevelReached:    maxLevelReached,
		player:             player,
		inventory:          invManager,
		agentManager:       agentMgr,
		statistics:         statsMgr,
		achievements:       achievementMgr,
		settings:           settings,
		rewardCalculator:   rewardCalc,
		tempStorage:        &rewarding.TempStorage{},
		enemyGenerator:     enemyGen,
		encounteredEnemies: encounteredEnemies,
		defeatedEnemies:    make(map[string]int),
	}

	// 撃破済み敵情報を復元
	if defeatedEnemies != nil {
		gs.SetDefeatedEnemies(defeatedEnemies)
	}

	return gs
}

// findCoreType はコア特性リストから指定IDのコア特性を検索します。
func findCoreType(coreTypes []domain.CoreType, coreTypeID string) domain.CoreType {
	for _, ct := range coreTypes {
		if ct.ID == coreTypeID {
			return ct
		}
	}
	// 見つからない場合は最初のコア特性を返す（フォールバック）
	if len(coreTypes) > 0 {
		return coreTypes[0]
	}
	return domain.CoreType{}
}

// findPassiveSkill はパッシブスキルマップから指定コア特性に対応するスキルを検索します。
func findPassiveSkill(passiveSkills map[string]domain.PassiveSkill, coreTypeID string) domain.PassiveSkill {
	// まずコア特性IDに対応するパッシブスキルを検索
	if skill, ok := passiveSkills[coreTypeID]; ok {
		return skill
	}
	// 見つからない場合は空のパッシブスキルを返す
	return domain.PassiveSkill{}
}

// findModuleDropInfo はモジュール定義リストから指定IDのモジュール定義を検索します。
func findModuleDropInfo(moduleDefs []rewarding.ModuleDropInfo, moduleID string) *rewarding.ModuleDropInfo {
	for i := range moduleDefs {
		if moduleDefs[i].ID == moduleID {
			return &moduleDefs[i]
		}
	}
	return nil
}

// findChainEffectDefinition はチェイン効果定義を検索します。
func findChainEffectDefinition(defs []rewarding.ChainEffectDefinition, effectType domain.ChainEffectType) *rewarding.ChainEffectDefinition {
	for i := range defs {
		if defs[i].EffectType == effectType {
			return &defs[i]
		}
	}
	return nil
}
