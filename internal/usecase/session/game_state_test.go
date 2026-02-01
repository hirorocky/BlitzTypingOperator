// Package game_state はゲーム全体の状態管理を提供するユースケースです。
package session

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestNewGameState は新しいGameStateの作成をテストします。
func TestNewGameState(t *testing.T) {
	gs := NewGameStateForTest()

	if gs == nil {
		t.Fatal("NewGameStateForTest() returned nil")
	}

	// 初期状態の確認
	if gs.GetMaxLevelReached() != 0 {
		t.Errorf("GetMaxLevelReached() expected 0, got %d", gs.GetMaxLevelReached())
	}

	if gs.Player() == nil {
		t.Error("Player() returned nil")
	}

	if gs.Inventory() == nil {
		t.Error("Inventory() returned nil")
	}

	if gs.Statistics() == nil {
		t.Error("Statistics() returned nil")
	}

	if gs.Settings() == nil {
		t.Error("Settings() returned nil")
	}
}

// TestRecordBattleVictory はバトル勝利の記録をテストします。
func TestRecordBattleVictory(t *testing.T) {
	gs := NewGameStateForTest()

	// 統計のみをテスト
	gs.RecordBattleVictory(1, 1)

	stats := gs.Statistics()
	if stats.Battle().Wins != 1 {
		t.Errorf("Wins expected 1, got %d", stats.Battle().Wins)
	}

	// 追加の勝利を記録
	gs.RecordBattleVictory(5, 3)
	if stats.Battle().Wins != 2 {
		t.Errorf("Wins expected 2, got %d", stats.Battle().Wins)
	}

	gs.RecordBattleVictory(4, 2)
	if stats.Battle().Wins != 3 {
		t.Errorf("Wins expected 3, got %d", stats.Battle().Wins)
	}
}

// TestRecordBattleDefeat はバトル敗北の記録をテストします。
func TestRecordBattleDefeat(t *testing.T) {
	gs := NewGameStateForTest()

	gs.RecordBattleDefeat(1)

	stats := gs.Statistics()
	if stats.Battle().Losses != 1 {
		t.Errorf("Losses expected 1, got %d", stats.Battle().Losses)
	}
}

// TestRecordTypingResult はタイピング結果の記録をテストします。
func TestRecordTypingResult(t *testing.T) {
	gs := NewGameStateForTest()

	gs.RecordTypingResult(60, 95.0, 100, 95, 5)

	stats := gs.Statistics()
	if stats.Typing().MaxWPM != 60 {
		t.Errorf("MaxWPM expected 60, got %d", stats.Typing().MaxWPM)
	}

	// より高いWPMで記録
	gs.RecordTypingResult(80, 98.0, 100, 98, 2)
	if stats.Typing().MaxWPM != 80 {
		t.Errorf("MaxWPM expected 80, got %d", stats.Typing().MaxWPM)
	}
}

// TestAddEncounteredEnemy は敵エンカウントの記録をテストします。
func TestAddEncounteredEnemy(t *testing.T) {
	gs := NewGameStateForTest()

	// 敵を追加
	gs.AddEncounteredEnemy("enemy_001")
	enemies := gs.GetEncounteredEnemies()
	if len(enemies) != 1 {
		t.Errorf("Expected 1 enemy, got %d", len(enemies))
	}

	// 同じ敵を追加（重複防止）
	gs.AddEncounteredEnemy("enemy_001")
	enemies = gs.GetEncounteredEnemies()
	if len(enemies) != 1 {
		t.Errorf("Expected 1 enemy (no duplicates), got %d", len(enemies))
	}

	// 別の敵を追加
	gs.AddEncounteredEnemy("enemy_002")
	enemies = gs.GetEncounteredEnemies()
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies, got %d", len(enemies))
	}

	// 空のIDは無視
	gs.AddEncounteredEnemy("")
	enemies = gs.GetEncounteredEnemies()
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies (empty ID ignored), got %d", len(enemies))
	}
}

// TestPreparePlayerForBattle はバトル準備をテストします。
func TestPreparePlayerForBattle(t *testing.T) {
	gs := NewGameStateForTest()

	// MaxHPが0で空のエージェントリストの場合、MaxHPは変更されない
	gs.Player().InitializeHP(domain.InitialMaxHP)
	gs.PreparePlayerForBattle([]*domain.AgentModel{})

	player := gs.Player()
	// HPが初期最大HP（1000）であることを確認
	if player.MaxHP != domain.InitialMaxHP {
		t.Errorf("MaxHP should be %d, got %d", domain.InitialMaxHP, player.MaxHP)
	}
	if player.HP != player.MaxHP {
		t.Errorf("HP should equal MaxHP after preparation")
	}
}

// ==================== EnemyProgress 統合テスト（Task 7.1） ====================

// TestGameState_EnemyProgress はGameStateがEnemyProgressを正しく統合していることをテストします。
func TestGameState_EnemyProgress(t *testing.T) {
	gs := NewGameStateForTest()

	// 初期状態の確認
	progress := gs.EnemyProgress()
	if progress == nil {
		t.Fatal("EnemyProgress() returned nil")
	}

	// 初期ランクは1
	if progress.CurrentRank != 1 {
		t.Errorf("Initial CurrentRank expected 1, got %d", progress.CurrentRank)
	}

	// 初期状態では敵は未撃破
	if gs.IsEnemyDefeated("enemy_001") {
		t.Error("Enemy should not be defeated initially")
	}
}

// TestGameState_DefeatedEnemyProvider_Interface はDefeatedEnemyProviderインターフェースの実装をテストします。
func TestGameState_DefeatedEnemyProvider_Interface(t *testing.T) {
	gs := NewGameStateForTest()

	// IsEnemyDefeated: 未撃破の敵
	if gs.IsEnemyDefeated("enemy_001") {
		t.Error("Enemy should not be defeated initially")
	}

	// GetDefeatedLevel: 未撃破の敵は0
	if level := gs.GetDefeatedLevel("enemy_001"); level != 0 {
		t.Errorf("GetDefeatedLevel for undefeated enemy expected 0, got %d", level)
	}

	// GetMaxDefeatedLevel: 初期状態は0
	if level := gs.GetMaxDefeatedLevel(); level != 0 {
		t.Errorf("GetMaxDefeatedLevel expected 0 initially, got %d", level)
	}

	// GetMaxLevelReached: 初期状態は0
	if level := gs.GetMaxLevelReached(); level != 0 {
		t.Errorf("GetMaxLevelReached expected 0 initially, got %d", level)
	}

	// GetDefeatedEnemies: 初期状態は空マップ
	defeatedEnemies := gs.GetDefeatedEnemies()
	if len(defeatedEnemies) != 0 {
		t.Errorf("GetDefeatedEnemies expected empty map, got %d entries", len(defeatedEnemies))
	}
}

// TestGameState_RecordEnemyDefeat_WithEnemyProgress は撃破記録がEnemyProgressを正しく更新することをテストします。
func TestGameState_RecordEnemyDefeat_WithEnemyProgress(t *testing.T) {
	gs := NewGameStateForTest()

	// 敵を撃破
	gs.RecordEnemyDefeat("enemy_001", 5)

	// 撃破状態を確認
	if !gs.IsEnemyDefeated("enemy_001") {
		t.Error("Enemy should be defeated after RecordEnemyDefeat")
	}

	// 撃破レベルを確認
	if level := gs.GetDefeatedLevel("enemy_001"); level != 5 {
		t.Errorf("GetDefeatedLevel expected 5, got %d", level)
	}

	// GetMaxDefeatedLevel を確認
	if level := gs.GetMaxDefeatedLevel(); level != 5 {
		t.Errorf("GetMaxDefeatedLevel expected 5, got %d", level)
	}

	// GetDefeatedEnemies を確認
	defeatedEnemies := gs.GetDefeatedEnemies()
	if len(defeatedEnemies) != 1 {
		t.Errorf("GetDefeatedEnemies expected 1 entry, got %d", len(defeatedEnemies))
	}
	if defeatedEnemies["enemy_001"] != 5 {
		t.Errorf("GetDefeatedEnemies[enemy_001] expected 5, got %d", defeatedEnemies["enemy_001"])
	}

	// より高いレベルで撃破
	gs.RecordEnemyDefeat("enemy_001", 10)
	if level := gs.GetDefeatedLevel("enemy_001"); level != 10 {
		t.Errorf("GetDefeatedLevel after higher level defeat expected 10, got %d", level)
	}

	// より低いレベルで撃破しても更新されない
	gs.RecordEnemyDefeat("enemy_001", 7)
	if level := gs.GetDefeatedLevel("enemy_001"); level != 10 {
		t.Errorf("GetDefeatedLevel should not decrease, expected 10, got %d", level)
	}
}

// TestGameState_SetEnemyProgress はSetEnemyProgressが正しく動作することをテストします。
func TestGameState_SetEnemyProgress(t *testing.T) {
	gs := NewGameStateForTest()

	// 新しいEnemyProgressを作成して設定
	newProgress := domain.NewEnemyProgress()
	newProgress.CurrentRank = 3
	newProgress.DefeatRecords["enemy_001"] = domain.EnemyDefeatRecord{
		Defeated:         true,
		MaxDefeatedLevel: 15,
	}

	gs.SetEnemyProgress(newProgress)

	// 設定が反映されていることを確認
	progress := gs.EnemyProgress()
	if progress.CurrentRank != 3 {
		t.Errorf("CurrentRank expected 3, got %d", progress.CurrentRank)
	}
	if !gs.IsEnemyDefeated("enemy_001") {
		t.Error("Enemy should be defeated after SetEnemyProgress")
	}
	if level := gs.GetDefeatedLevel("enemy_001"); level != 15 {
		t.Errorf("GetDefeatedLevel expected 15, got %d", level)
	}
}
