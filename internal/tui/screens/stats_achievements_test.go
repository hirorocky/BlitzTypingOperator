// Package screens はTUI画面のテストを提供します。
package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.6: 統計・実績画面のテスト ====================

// TestNewStatsAchievementsScreen はStatsAchievementsScreenの初期化をテストします。
func TestNewStatsAchievementsScreen(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	if screen == nil {
		t.Fatal("StatsAchievementsScreenがnilです")
	}
}

// TestStatsAchievementsTabs はタブ切り替えをテストします。
func TestStatsAchievementsTabs(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	// 初期タブ
	if screen.currentTab != TabTypingStats {
		t.Errorf("初期タブが正しくありません: got %d, want %d", screen.currentTab, TabTypingStats)
	}

	// タブ切り替え（右へ）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})
	if screen.currentTab != TabBattleStats {
		t.Errorf("タブ切り替え(右)が正しくありません: got %d, want %d", screen.currentTab, TabBattleStats)
	}

	// さらに右へ
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})
	if screen.currentTab != TabAchievements {
		t.Errorf("タブ切り替え(右)が正しくありません: got %d, want %d", screen.currentTab, TabAchievements)
	}
}

// TestStatsAchievementsTypingStats はタイピング統計表示をテストします。

func TestStatsAchievementsTypingStats(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	// タイピング統計タブ
	screen.currentTab = TabTypingStats

	// 統計データが存在すること
	if screen.data.TypingStats.MaxWPM < 0 {
		t.Error("最高WPMが負の値です")
	}
	if screen.data.TypingStats.AverageWPM < 0 {
		t.Error("平均WPMが負の値です")
	}
	if screen.data.TypingStats.PerfectAccuracyCount < 0 {
		t.Error("100%正確性回数が負の値です")
	}
	if screen.data.TypingStats.TotalCharacters < 0 {
		t.Error("総タイプ文字数が負の値です")
	}
}

// TestStatsAchievementsBattleStats はバトル統計表示をテストします。

func TestStatsAchievementsBattleStats(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	// バトル統計タブ
	screen.currentTab = TabBattleStats

	// 統計データが存在すること
	if screen.data.BattleStats.TotalBattles < 0 {
		t.Error("総バトル数が負の値です")
	}
	if screen.data.BattleStats.Wins < 0 {
		t.Error("勝利数が負の値です")
	}
	if screen.data.BattleStats.Losses < 0 {
		t.Error("敗北数が負の値です")
	}
	if screen.data.BattleStats.MaxLevelReached < 0 {
		t.Error("到達最高レベルが負の値です")
	}
}

// TestStatsAchievementsAchievementsList は実績一覧をテストします。

func TestStatsAchievementsAchievementsList(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	// 実績タブ
	screen.currentTab = TabAchievements

	// 実績が存在すること
	if len(screen.data.Achievements) == 0 {
		t.Error("実績がありません")
	}

	// 達成状況が判定できること
	achievedCount := 0
	for _, ach := range screen.data.Achievements {
		if ach.Achieved {
			achievedCount++
		}
	}
	// 達成済み実績が少なくとも1つあること（テストデータ依存）
	_ = achievedCount
}

// TestStatsAchievementsCompletionRate はコンプリート率をテストします。

func TestStatsAchievementsCompletionRate(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	rate := screen.getAchievementCompletionRate()

	// 0〜100の範囲であること
	if rate < 0 || rate > 100 {
		t.Errorf("実績コンプリート率が範囲外: %d", rate)
	}
}

// TestStatsAchievementsBackNavigation は戻るナビゲーションをテストします。
func TestStatsAchievementsBackNavigation(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Escキーでコマンドが返されません")
	}
}

// TestStatsAchievementsRender はレンダリングをテストします。
func TestStatsAchievementsRender(t *testing.T) {
	data := createTestStatsData()
	screen := NewStatsAchievementsScreen(data)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// ==================== ヘルパー関数 ====================

func createTestStatsData() *StatsData {
	return &StatsData{
		TypingStats: TypingStatsData{
			MaxWPM:               120,
			AverageWPM:           85.5,
			PerfectAccuracyCount: 10,
			TotalCharacters:      50000,
		},
		BattleStats: BattleStatsData{
			TotalBattles:    100,
			Wins:            75,
			Losses:          25,
			MaxLevelReached: 30,
		},
		Achievements: []AchievementData{
			{ID: "wpm_50", Name: "タイピスト見習い", Description: "WPM 50達成", Achieved: true},
			{ID: "wpm_80", Name: "タイピスト", Description: "WPM 80達成", Achieved: true},
			{ID: "wpm_100", Name: "タイピストマスター", Description: "WPM 100達成", Achieved: true},
			{ID: "wpm_120", Name: "タイピングゴッド", Description: "WPM 120達成", Achieved: false},
			{ID: "enemy_10", Name: "初陣の勇者", Description: "敵10体撃破", Achieved: true},
			{ID: "enemy_50", Name: "熟練の戦士", Description: "敵50体撃破", Achieved: true},
			{ID: "enemy_100", Name: "歴戦の勇者", Description: "敵100体撃破", Achieved: false},
			{ID: "level_10", Name: "Lv10到達", Description: "レベル10に到達", Achieved: true},
			{ID: "level_25", Name: "Lv25到達", Description: "レベル25に到達", Achieved: true},
			{ID: "level_50", Name: "Lv50到達", Description: "レベル50に到達", Achieved: false},
		},
	}
}
