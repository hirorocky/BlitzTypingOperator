// Package presenter はUI用のデータ変換を提供します。
// usecase層のデータをtui/screens用のViewModelに変換します。
package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/usecase/session"
)

// TestCreateStatsData は統計データ作成をテストします。
func TestCreateStatsData(t *testing.T) {
	gs := session.NewGameState()

	// タイピング結果を記録
	gs.RecordTypingResult(60, 95.0, 100, 95, 5)

	// バトル勝利を記録
	gs.RecordBattleVictory(1)

	data := CreateStatsData(gs)

	if data == nil {
		t.Fatal("CreateStatsData returned nil")
	}

	// タイピング統計
	if data.TypingStats.MaxWPM != 60 {
		t.Errorf("MaxWPM expected 60, got %d", data.TypingStats.MaxWPM)
	}

	// バトル統計
	if data.BattleStats.Wins != 1 {
		t.Errorf("Wins expected 1, got %d", data.BattleStats.Wins)
	}
	if data.BattleStats.MaxLevelReached != 1 {
		t.Errorf("MaxLevelReached expected 1, got %d", data.BattleStats.MaxLevelReached)
	}
}

// TestCreateStatsData_Empty は空の状態での統計データ作成をテストします。
func TestCreateStatsData_Empty(t *testing.T) {
	gs := session.NewGameState()

	data := CreateStatsData(gs)

	if data == nil {
		t.Fatal("CreateStatsData returned nil")
	}

	// 初期状態
	if data.TypingStats.MaxWPM != 0 {
		t.Errorf("MaxWPM expected 0, got %d", data.TypingStats.MaxWPM)
	}
	if data.BattleStats.TotalBattles != 0 {
		t.Errorf("TotalBattles expected 0, got %d", data.BattleStats.TotalBattles)
	}
}
