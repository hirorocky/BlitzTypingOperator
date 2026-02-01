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
	SkillTypes             []rewarding.ModuleDropInfo
	EnemyTypes             []domain.EnemyType
	PassiveSkills          map[string]domain.PassiveSkill
	ChainEffectDefinitions []rewarding.ChainEffectDefinition
}

// ToSaveData はGameStateをセーブデータに変換します。
func (g *GameState) ToSaveData() *savedata.SaveData {
	saveData := savedata.NewSaveData()

	saveData.Statistics.MaxLevelReached = g.GetMaxLevelReached()

	// プレイヤーのMaxHPを保存
	saveData.Player.MaxHP = g.player.MaxHP

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

	// 敵進行データを保存
	saveData.EnemyProgress.CurrentRank = g.enemyProgress.CurrentRank
	for id, record := range g.enemyProgress.DefeatRecords {
		saveData.EnemyProgress.DefeatRecords[id] = savedata.DefeatRecordSave{
			Defeated:         record.Defeated,
			MaxDefeatedLevel: record.MaxDefeatedLevel,
		}
	}

	return saveData
}

// GameStateFromSaveData はセーブデータからGameStateを生成します。
// sourcesにはマスタデータから変換されたドメイン型データを渡す必要があります。
// 注意: エージェントスロットの復元はRootModel側で行います（AgentSlotManagerへの依存を避けるため）
func GameStateFromSaveData(data *savedata.SaveData, sources *DomainDataSources) *GameState {
	if sources == nil {
		slog.Error("マスタデータソースが必要です")
		return nil
	}

	// マスタデータを取得（ドメイン型）
	coreTypes := sources.CoreTypes
	moduleTypes := sources.SkillTypes
	passiveSkills := sources.PassiveSkills
	enemyTypes := sources.EnemyTypes

	// インベントリマネージャーを作成
	invManager := NewInventoryManager()

	// UniqueCores/UniqueSkillsはRootModel側で別途ロード

	// プレイヤーを作成（MaxHPをセーブデータから復元）
	maxHP := domain.InitialMaxHP
	if data.Player != nil && data.Player.MaxHP > 0 {
		maxHP = data.Player.MaxHP
	}
	player := domain.NewPlayerWithMaxHP(maxHP)

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

	var encounteredEnemies []string
	var defeatedEnemies map[string]int
	if data.Statistics != nil {
		encounteredEnemies = data.Statistics.EncounteredEnemies
		defeatedEnemies = data.Statistics.DefeatedEnemies
	}

	// EnemyProgressを復元
	enemyProgress := domain.NewEnemyProgress()
	if data.EnemyProgress != nil && len(data.EnemyProgress.DefeatRecords) > 0 {
		// EnemyProgressセクションから復元
		enemyProgress.CurrentRank = data.EnemyProgress.CurrentRank
		if enemyProgress.CurrentRank < 1 {
			enemyProgress.CurrentRank = 1
		}
		for id, record := range data.EnemyProgress.DefeatRecords {
			enemyProgress.DefeatRecords[id] = domain.EnemyDefeatRecord{
				Defeated:         record.Defeated,
				MaxDefeatedLevel: record.MaxDefeatedLevel,
			}
		}
	} else if len(defeatedEnemies) > 0 {
		// EnemyProgressセクションがない場合、Statisticsから撃破記録を復元
		for id, level := range defeatedEnemies {
			enemyProgress.DefeatRecords[id] = domain.EnemyDefeatRecord{
				Defeated:         true,
				MaxDefeatedLevel: level,
			}
		}
		// ランク情報がない場合はランク1から開始（次回勝利時に正しく更新される）
		enemyProgress.CurrentRank = 1
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
		enemyProgress:      enemyProgress,
	}

	return gs
}
