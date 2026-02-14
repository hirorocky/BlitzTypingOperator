// Package game_state はゲーム全体の状態管理を提供するユースケースです。
// プレイヤー情報、インベントリ、統計、実績、設定などを一元管理します。
package session

import (
	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/achievement"
	"hirorocky/type-battle/internal/usecase/rewarding"
	"hirorocky/type-battle/internal/usecase/spawning"
)

// GameState はゲーム全体の状態を保持する構造体です。
// プレイヤー情報、インベントリ、統計、実績、設定などを含みます。
// セーブ/ロード時にはこの構造体がJSON形式で永続化されます。
type GameState struct {
	// player はプレイヤーの状態です。
	player *domain.PlayerModel

	// inventory はゲーム全体のインベントリマネージャーです。
	inventory *InventoryManager

	// statistics は統計情報を管理します。
	statistics *StatisticsManager

	// achievements は実績管理を担当します。
	achievements *achievement.AchievementManager

	// settings はゲーム設定です。
	settings *Settings

	// rewardCalculator は報酬計算を担当します。
	rewardCalculator *rewarding.RewardCalculator

	// tempStorage は一時保管を担当します（インベントリ満杯時用）。
	tempStorage *rewarding.TempStorage

	// enemyGenerator は敵生成を担当します。
	enemyGenerator *spawning.EnemyGenerator

	// encounteredEnemies はエンカウントした敵のIDリストです（敵図鑑用）。
	encounteredEnemies []string

	// enemyProgress は敵撃破状況を管理します。
	enemyProgress *domain.EnemyProgress
}

// NewGameState はマスタデータを使用して新しいGameStateを作成します。
// 初回起動時やセーブデータが存在しない場合に使用されます。
func NewGameState(
	coreTypes []domain.CoreType,
	moduleTypes []rewarding.ModuleDropInfo,
	passiveSkills map[string]domain.PassiveSkill,
) *GameState {
	// インベントリマネージャーを作成
	invManager := NewInventoryManager()

	// 実績マネージャーを作成
	achievementMgr := achievement.NewAchievementManager()

	// RewardCalculatorを作成
	rewardCalc := rewarding.NewRewardCalculator(coreTypes, moduleTypes, passiveSkills)

	// EnemyGeneratorを作成（デフォルト敵タイプを使用）
	enemyGen := spawning.NewEnemyGenerator(nil)

	return &GameState{
		player:           domain.NewPlayerWithMaxHP(0),
		inventory:        invManager,
		statistics:       NewStatisticsManager(),
		achievements:     achievementMgr,
		settings:         NewSettings(),
		rewardCalculator: rewardCalc,
		tempStorage:      &rewarding.TempStorage{},
		enemyGenerator:   enemyGen,
		enemyProgress:    domain.NewEnemyProgress(),
	}
}

// Player はプレイヤーの状態を返します。
func (g *GameState) Player() *domain.PlayerModel {
	return g.player
}

// Inventory はインベントリマネージャーを返します。
func (g *GameState) Inventory() *InventoryManager {
	return g.inventory
}

// Statistics は統計マネージャーを返します。
func (g *GameState) Statistics() *StatisticsManager {
	return g.statistics
}

// Achievements は実績マネージャーを返します。
func (g *GameState) Achievements() *achievement.AchievementManager {
	return g.achievements
}

// Settings は設定を返します。
func (g *GameState) Settings() *Settings {
	return g.settings
}

// EnemyGenerator は敵生成器を返します。
func (g *GameState) EnemyGenerator() *spawning.EnemyGenerator {
	return g.enemyGenerator
}

// UpdateEnemyGenerator は敵生成器を敵タイプで更新します。
func (g *GameState) UpdateEnemyGenerator(enemyTypes []domain.EnemyType) {
	if len(enemyTypes) > 0 {
		g.enemyGenerator = spawning.NewEnemyGenerator(enemyTypes)
	}
}

// UpdateRewardCalculator は報酬計算器を更新します。
func (g *GameState) UpdateRewardCalculator(coreTypes []domain.CoreType, moduleTypes []rewarding.ModuleDropInfo, passiveSkills map[string]domain.PassiveSkill) {
	if len(coreTypes) > 0 || len(moduleTypes) > 0 {
		g.rewardCalculator = rewarding.NewRewardCalculator(coreTypes, moduleTypes, passiveSkills)
	}
}

// RecordBattleVictory はバトル勝利を記録します。
// selectedLevel は選択したレベル（統計記録用）、defaultLevel は敵のデフォルトレベル（MaxLevelReached更新用）です。
func (g *GameState) RecordBattleVictory(selectedLevel int, _ int) {
	g.statistics.RecordBattleResult(true, selectedLevel)
	g.checkAchievements()
}

// GetMaxLevelReached は到達最高レベルを返します。
func (g *GameState) GetMaxLevelReached() int {
	return g.GetMaxDefeatedLevel()
}

// RecordBattleDefeat はバトル敗北を記録します。
func (g *GameState) RecordBattleDefeat(level int) {
	g.statistics.RecordBattleResult(false, level)
}

// RecordTypingResult はタイピング結果を記録します。
func (g *GameState) RecordTypingResult(wpm int, accuracy float64, characters int, correct int, missed int) {
	g.statistics.RecordTypingResult(wpm, accuracy, characters, correct, missed)

	// 実績チェック
	g.checkAchievements()
}

// checkAchievements は実績の達成状況をチェックします。
func (g *GameState) checkAchievements() {
	stats := g.statistics

	// タイピング実績をチェック
	g.achievements.CheckTypingAchievements(
		float64(stats.Typing().MaxWPM),
		g.statistics.GetAccuracyRate(),
	)

	g.achievements.CheckBattleAchievements(
		stats.Battle().TotalEnemiesDefeated,
		g.GetMaxLevelReached(),
		false,
	)
}

// CheckBattleAchievementsWithNoDamage はノーダメージ判定付きでバトル実績をチェックします。
func (g *GameState) CheckBattleAchievementsWithNoDamage(noDamage bool) {
	stats := g.statistics

	// タイピング実績をチェック
	g.achievements.CheckTypingAchievements(
		float64(stats.Typing().MaxWPM),
		g.statistics.GetAccuracyRate(),
	)

	// バトル実績をチェック（ノーダメージ判定付き）
	g.achievements.CheckBattleAchievements(
		stats.Battle().TotalEnemiesDefeated,
		g.GetMaxLevelReached(),
		noDamage,
	)
}

// AddEncounteredEnemy は敵をエンカウント済みとして記録します（敵図鑑用）。
func (g *GameState) AddEncounteredEnemy(enemyID string) {
	// 空のIDは無視
	if enemyID == "" {
		return
	}
	// 重複チェック
	for _, id := range g.encounteredEnemies {
		if id == enemyID {
			return
		}
	}
	g.encounteredEnemies = append(g.encounteredEnemies, enemyID)
}

// GetEncounteredEnemies はエンカウント済みの敵IDリストを返します。
func (g *GameState) GetEncounteredEnemies() []string {
	return g.encounteredEnemies
}

// PreparePlayerForBattle はプレイヤーをバトル用に準備します。
// agents: バトルに参加するエージェントのリスト（AgentSlotManager.BuildAgentsForBattle()から取得）
// 注意: MaxHPは事前に設定されている必要があります（セーブデータロード時 or NewPlayerWithMaxHP）
func (g *GameState) PreparePlayerForBattle(agents []*domain.AgentModel) {
	if g.player.MaxHP == 0 && len(agents) > 0 {
		g.player.InitializeHP(domain.InitialMaxHP)
	}
	g.player.MaxMana = config.DefaultMaxMana
	g.player.PrepareForBattle()
}

// RewardCalculator は報酬計算を返します。
func (g *GameState) RewardCalculator() *rewarding.RewardCalculator {
	return g.rewardCalculator
}

// TempStorage は一時保管を返します。
func (g *GameState) TempStorage() *rewarding.TempStorage {
	return g.tempStorage
}

// AddRewardsToInventory は報酬をインベントリに追加します。
func (g *GameState) AddRewardsToInventory(result *rewarding.RewardResult) *rewarding.InventoryWarning {
	return rewarding.AddRewardsToInventory(
		result,
		g.inventory.Cores(),
		g.inventory.Skills(),
		g.inventory.ChainEffects(),
	)
}

// ========== 撃破済み敵情報の管理 ==========

// EnemyProgress は敵撃破状況のドメインモデルを返します。
func (g *GameState) EnemyProgress() *domain.EnemyProgress {
	return g.enemyProgress
}

// SetEnemyProgress は敵撃破状況を設定します（セーブデータロード用）。
func (g *GameState) SetEnemyProgress(progress *domain.EnemyProgress) {
	if progress == nil {
		g.enemyProgress = domain.NewEnemyProgress()
		return
	}
	g.enemyProgress = progress
}

// RecordEnemyDefeat は敵の撃破を記録します。
// 既に記録されている敵の場合、より高いレベルで撃破した場合のみ更新します。
func (g *GameState) RecordEnemyDefeat(enemyTypeID string, level int) {
	if enemyTypeID == "" {
		return
	}

	if g.enemyProgress == nil {
		g.enemyProgress = domain.NewEnemyProgress()
	}

	// EnemyProgress.RecordDefeatを使用（HP増加は無視）
	g.enemyProgress.RecordDefeat(enemyTypeID, level)
}

// GetDefeatedEnemies は撃破済み敵のマップを返します。
// キーは敵タイプID、値は撃破した最高レベルです。
func (g *GameState) GetDefeatedEnemies() map[string]int {
	if g.enemyProgress == nil {
		return make(map[string]int)
	}
	// EnemyProgressのDefeatRecordsから変換
	result := make(map[string]int)
	for id, record := range g.enemyProgress.DefeatRecords {
		if record.Defeated {
			result[id] = record.MaxDefeatedLevel
		}
	}
	return result
}

// IsEnemyDefeated は指定した敵タイプが一度でも撃破されているかどうかを返します。
func (g *GameState) IsEnemyDefeated(enemyTypeID string) bool {
	if g.enemyProgress == nil {
		return false
	}
	return g.enemyProgress.IsDefeated(enemyTypeID)
}

// GetDefeatedLevel は指定した敵タイプの撃破最高レベルを返します。
// 未撃破の場合は0を返します。
func (g *GameState) GetDefeatedLevel(enemyTypeID string) int {
	if g.enemyProgress == nil {
		return 0
	}
	return g.enemyProgress.GetMaxDefeatedLevel(enemyTypeID)
}

// SetDefeatedEnemies は撃破済み敵情報を設定します（セーブデータロード用）。
func (g *GameState) SetDefeatedEnemies(defeated map[string]int) {
	if g.enemyProgress == nil {
		g.enemyProgress = domain.NewEnemyProgress()
	}

	if defeated == nil {
		return
	}

	// map[string]int形式からEnemyProgress形式に変換
	for id, level := range defeated {
		g.enemyProgress.RecordDefeat(id, level)
	}
}

// GetMaxDefeatedLevel は全敵種類を通じた最高撃破レベル（到達Lv）を返します。
// 一度も敵を撃破していない場合は0を返します。
func (g *GameState) GetMaxDefeatedLevel() int {
	if g.enemyProgress == nil {
		return 0
	}
	maxLevel := 0
	for _, record := range g.enemyProgress.DefeatRecords {
		if record.Defeated && record.MaxDefeatedLevel > maxLevel {
			maxLevel = record.MaxDefeatedLevel
		}
	}
	return maxLevel
}
