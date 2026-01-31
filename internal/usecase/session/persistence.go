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
// v4.0.0形式: EnemyProgress + ユニークインベントリ + エージェントスロットシステム
// 注意: 旧形式(AgentInstances, CoreInstances, ModuleInstances)は後方互換性のため空で保存
func (g *GameState) ToSaveData() *savedata.SaveData {
	saveData := savedata.NewSaveData()

	// v4.0.0: 最高到達レベルはEnemyProgress経由で取得
	saveData.Statistics.MaxLevelReached = g.GetMaxLevelReached()

	// v3.0.0: 旧形式は空で保存（後方互換性）
	saveData.Inventory.CoreInstances = []savedata.CoreInstanceSave{}
	saveData.Inventory.ModuleInstances = []savedata.ModuleInstanceSave{}
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{}

	// v3.0.0: 新システムでは容量制限なし（後方互換性のため固定値を設定）
	saveData.Inventory.MaxCoreSlots = 100
	saveData.Inventory.MaxModuleSlots = 200

	// v3.0.0: EquippedAgentIDsは空（AgentSlotsを使用）
	saveData.Player.EquippedAgentIDs = [3]string{}

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

	// v4.0.0: 撃破済み敵情報はEnemyProgress経由で取得（後方互換用）
	saveData.Statistics.DefeatedEnemies = g.GetDefeatedEnemies()

	return saveData
}

// GameStateFromSaveData はセーブデータからGameStateを生成します。
// v3.0.0形式のセーブデータからオブジェクトを再構築します。
// sourcesにはマスタデータから変換されたドメイン型データを渡す必要があります。
// 注意: エージェントスロットの復元はRootModel側で行います（AgentSlotManagerへの依存を避けるため）
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

	// インベントリマネージャーを作成
	invManager := NewInventoryManager()

	// v3.0.0: 旧形式のコア/モジュールはスキップ（後方互換性のためデータは読み飛ばす）
	// 新形式のUniqueCores/UniqueSkillsはRootModel側で別途ロード

	// プレイヤーを作成
	player := domain.NewPlayer()

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

	// v4.0.0: エンカウント敵リスト、撃破済み敵情報を取得
	var encounteredEnemies []string
	var defeatedEnemies map[string]int
	if data.Statistics != nil {
		encounteredEnemies = data.Statistics.EncounteredEnemies
		defeatedEnemies = data.Statistics.DefeatedEnemies
	}

	gs := &GameState{
		player:             player,
		inventory:          invManager,
		statistics:         statsMgr,
		achievements:       achievementMgr,
		settings:           settings,
		rewardCalculator:   rewardCalc,
		tempStorage:        &rewarding.TempStorage{},
		enemyGenerator:     enemyGen,
		encounteredEnemies: encounteredEnemies,
		enemyProgress:      domain.NewEnemyProgress(),
	}

	// v4.0.0: 撃破済み敵情報をEnemyProgressに復元
	if defeatedEnemies != nil {
		gs.SetDefeatedEnemies(defeatedEnemies)
	}

	return gs
}
