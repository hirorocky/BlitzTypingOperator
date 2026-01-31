// Package progress は敵進行管理のユースケースを提供します。
// 敵の撃破状況、ランク進行、プレイヤーHP成長を一元管理します。
package progress

import (
	"hirorocky/type-battle/internal/domain"
)

// VictoryResult は撃破記録の結果を表す構造体です。
type VictoryResult struct {
	// HPGain は獲得最大HP（PlayerModelに適用済み）です。
	HPGain int

	// RankUnlocked は新ランク解放フラグです。
	RankUnlocked bool

	// NewMaxLevel は新しい選択可能最大レベルです（変化なしなら0）。
	NewMaxLevel int
}

// EnemyProgressManager は敵進行とHP成長を一元管理するユースケースです。
// EnemyProgress（ドメイン）とPlayerModel（HP成長の適用先）を統合的に操作します。
type EnemyProgressManager struct {
	// progress は敵撃破状況を管理するドメインモデルです。
	progress *domain.EnemyProgress

	// player はHP成長の適用先（Single Source of Truth）です。
	player *domain.PlayerModel

	// enemyTypes は敵タイプのマスタデータです。
	enemyTypes map[string]domain.EnemyType
}

// NewEnemyProgressManager は新しいEnemyProgressManagerを作成します。
// progress: 敵撃破状況のドメインモデル
// player: HP成長を適用するPlayerModel
// enemyTypes: 敵タイプのマスタデータ（キー: 敵タイプID）
func NewEnemyProgressManager(
	progress *domain.EnemyProgress,
	player *domain.PlayerModel,
	enemyTypes map[string]domain.EnemyType,
) *EnemyProgressManager {
	return &EnemyProgressManager{
		progress:   progress,
		player:     player,
		enemyTypes: enemyTypes,
	}
}

// GetCurrentRankEnemies は現在ランクの敵リストを返します。
// 現在解放済みランクに属する全ての敵タイプを返します。
func (m *EnemyProgressManager) GetCurrentRankEnemies() []domain.EnemyType {
	currentRank := m.progress.CurrentRank
	enemies := make([]domain.EnemyType, 0)

	for _, et := range m.enemyTypes {
		if et.Rank == currentRank {
			enemies = append(enemies, et)
		}
	}

	return enemies
}

// GetSelectableLevelRange は敵の選択可能レベル範囲を取得します。
// 未撃破: デフォルトレベルのみ
// 撃破済み: 1 から (最大撃破レベル + 5) まで
func (m *EnemyProgressManager) GetSelectableLevelRange(enemyTypeID string) (min, max int) {
	et, exists := m.enemyTypes[enemyTypeID]
	if !exists {
		return 1, 1
	}

	return m.progress.GetSelectableLevelRange(enemyTypeID, et.DefaultLevel)
}

// RecordVictory は撃破を記録しHP成長を適用します。
// EnemyProgress.RecordDefeat()を呼び出し、hpGainをPlayerModel.IncreaseMaxHP()に適用します。
// ランク内全敵撃破時は次ランクを解放します。
func (m *EnemyProgressManager) RecordVictory(enemyTypeID string, level int) VictoryResult {
	// 撃破前の状態を記録
	wasDefeated := m.progress.IsDefeated(enemyTypeID)
	oldMaxLevel := m.progress.GetMaxDefeatedLevel(enemyTypeID)

	// 撃破を記録
	hpGain, _ := m.progress.RecordDefeat(enemyTypeID, level)

	// HP成長をPlayerModelに適用
	if hpGain > 0 {
		m.player.IncreaseMaxHP(hpGain)
	}

	// 新しい選択可能最大レベルを計算
	newMaxLevel := 0
	if !wasDefeated || level > oldMaxLevel {
		// 初撃破または新記録の場合
		newMaxLevel = m.progress.GetMaxDefeatedLevel(enemyTypeID) + domain.SelectableLevelBonus
	}

	// ランク解放チェック
	rankUnlocked := false
	if m.CheckRankComplete() {
		m.progress.CurrentRank++
		rankUnlocked = true
	}

	return VictoryResult{
		HPGain:       hpGain,
		RankUnlocked: rankUnlocked,
		NewMaxLevel:  newMaxLevel,
	}
}

// CheckRankComplete はランク内全敵撃破チェックを行います。
// 現在ランクの全敵が撃破済みならtrueを返します。
func (m *EnemyProgressManager) CheckRankComplete() bool {
	currentRank := m.progress.CurrentRank

	for _, et := range m.enemyTypes {
		if et.Rank == currentRank {
			if !m.progress.IsDefeated(et.ID) {
				return false
			}
		}
	}

	// 現在ランクに敵が1体でもいることを確認
	hasEnemiesInRank := false
	for _, et := range m.enemyTypes {
		if et.Rank == currentRank {
			hasEnemiesInRank = true
			break
		}
	}

	return hasEnemiesInRank
}

// GetEnemyByID は敵タイプIDから敵タイプを取得します。
// 存在しない場合は空のEnemyTypeとfalseを返します。
func (m *EnemyProgressManager) GetEnemyByID(enemyTypeID string) (domain.EnemyType, bool) {
	et, exists := m.enemyTypes[enemyTypeID]
	return et, exists
}

// IsEnemyInCurrentRank は敵が現在ランク内かどうかを判定します。
func (m *EnemyProgressManager) IsEnemyInCurrentRank(enemyTypeID string) bool {
	et, exists := m.enemyTypes[enemyTypeID]
	if !exists {
		return false
	}
	return et.Rank == m.progress.CurrentRank
}

// GetProgress は管理中のEnemyProgressを返します。
// セーブ時など、直接ドメインモデルにアクセスする場合に使用します。
func (m *EnemyProgressManager) GetProgress() *domain.EnemyProgress {
	return m.progress
}

// GetPlayer は管理中のPlayerModelを返します。
func (m *EnemyProgressManager) GetPlayer() *domain.PlayerModel {
	return m.player
}

// IsDefeated は指定した敵が撃破済みかどうかを返します。
func (m *EnemyProgressManager) IsDefeated(enemyTypeID string) bool {
	return m.progress.IsDefeated(enemyTypeID)
}

// GetMaxDefeatedLevel は指定した敵の撃破済み最大レベルを返します。
func (m *EnemyProgressManager) GetMaxDefeatedLevel(enemyTypeID string) int {
	return m.progress.GetMaxDefeatedLevel(enemyTypeID)
}

// GetCurrentRank は現在解放済みランクを返します。
func (m *EnemyProgressManager) GetCurrentRank() int {
	return m.progress.CurrentRank
}
