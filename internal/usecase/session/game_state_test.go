// Package game_state はゲーム全体の状態管理を提供するユースケースです。
package session

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestNewGameState は新しいGameStateの作成をテストします。
func TestNewGameState(t *testing.T) {
	gs := NewGameState()

	if gs == nil {
		t.Fatal("NewGameState() returned nil")
	}

	// 初期状態の確認
	if gs.MaxLevelReached != 0 {
		t.Errorf("MaxLevelReached expected 0, got %d", gs.MaxLevelReached)
	}

	if gs.Player() == nil {
		t.Error("Player() returned nil")
	}

	if gs.Inventory() == nil {
		t.Error("Inventory() returned nil")
	}

	if gs.AgentManager() == nil {
		t.Error("AgentManager() returned nil")
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
	gs := NewGameState()

	// レベル1で勝利
	gs.RecordBattleVictory(1)

	if gs.MaxLevelReached != 1 {
		t.Errorf("MaxLevelReached expected 1, got %d", gs.MaxLevelReached)
	}

	stats := gs.Statistics()
	if stats.Battle().Wins != 1 {
		t.Errorf("Wins expected 1, got %d", stats.Battle().Wins)
	}

	// レベル3で勝利（MaxLevelReachedが更新される）
	gs.RecordBattleVictory(3)
	if gs.MaxLevelReached != 3 {
		t.Errorf("MaxLevelReached expected 3, got %d", gs.MaxLevelReached)
	}

	// レベル2で勝利（MaxLevelReachedは更新されない）
	gs.RecordBattleVictory(2)
	if gs.MaxLevelReached != 3 {
		t.Errorf("MaxLevelReached expected 3, got %d", gs.MaxLevelReached)
	}
}

// TestRecordBattleDefeat はバトル敗北の記録をテストします。
func TestRecordBattleDefeat(t *testing.T) {
	gs := NewGameState()

	gs.RecordBattleDefeat(1)

	stats := gs.Statistics()
	if stats.Battle().Losses != 1 {
		t.Errorf("Losses expected 1, got %d", stats.Battle().Losses)
	}
}

// TestRecordTypingResult はタイピング結果の記録をテストします。
func TestRecordTypingResult(t *testing.T) {
	gs := NewGameState()

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
	gs := NewGameState()

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
	gs := NewGameState()

	// プレイヤーの準備
	gs.PreparePlayerForBattle()

	player := gs.Player()
	// HPが設定されていることを確認（最小でもBaseHP）
	if player.MaxHP < domain.BaseHP {
		t.Errorf("MaxHP should be at least BaseHP (%d), got %d", domain.BaseHP, player.MaxHP)
	}
	if player.HP != player.MaxHP {
		t.Errorf("HP should equal MaxHP after preparation")
	}
}
