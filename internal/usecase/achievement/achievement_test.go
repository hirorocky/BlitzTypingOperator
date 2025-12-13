// Package achievement は実績システムを担当します。

package achievement

import (
	"testing"
)

// ==================================================
// Task 12.1: タイピング実績テスト
// ==================================================

func TestTypingAchievements_WPMMilestones(t *testing.T) {

	manager := NewAchievementManager()

	tests := []struct {
		wpm          float64
		expectedID   string
		shouldUnlock bool
	}{
		{30, AchievementWPM50, false},
		{50, AchievementWPM50, true},
		{79, AchievementWPM80, false},
		{80, AchievementWPM80, true},
		{99, AchievementWPM100, false},
		{100, AchievementWPM100, true},
		{119, AchievementWPM120, false},
		{120, AchievementWPM120, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			manager = NewAchievementManager()                             // リセット
			notifications := manager.CheckTypingAchievements(tt.wpm, 100) // 100%正確性
			unlocked := manager.IsUnlocked(tt.expectedID)

			if tt.shouldUnlock && !unlocked {
				t.Errorf("WPM %.0f で実績 %s が解除されるべきです", tt.wpm, tt.expectedID)
			}
			if tt.shouldUnlock && len(notifications) == 0 {
				t.Errorf("実績解除通知が返されるべきです")
			}
		})
	}
}

func TestTypingAchievements_PerfectAccuracy(t *testing.T) {

	manager := NewAchievementManager()

	// 99%正確性では解除されない
	_ = manager.CheckTypingAchievements(50, 99)
	if manager.IsUnlocked(AchievementPerfectAccuracy) {
		t.Error("99%正確性では実績が解除されるべきではありません")
	}

	// 100%正確性で解除
	manager = NewAchievementManager()
	notifications := manager.CheckTypingAchievements(50, 100)
	if !manager.IsUnlocked(AchievementPerfectAccuracy) {
		t.Error("100%正確性で実績が解除されるべきです")
	}
	if len(notifications) == 0 {
		t.Error("実績解除通知が返されるべきです")
	}
}

func TestAchievementNotification(t *testing.T) {

	manager := NewAchievementManager()

	notifications := manager.CheckTypingAchievements(120, 100)
	// 120 WPMでは WPM50, WPM80, WPM100, WPM120, PerfectAccuracy が解除される
	if len(notifications) < 2 {
		t.Errorf("複数の通知が返されるべきです: got %d", len(notifications))
	}

	// 通知には実績名が含まれる
	for _, n := range notifications {
		if n.AchievementID == "" {
			t.Error("通知には実績IDが含まれるべきです")
		}
		if n.Name == "" {
			t.Error("通知には実績名が含まれるべきです")
		}
	}
}

// ==================================================
// Task 12.2: バトル実績テスト
// ==================================================

func TestBattleAchievements_EnemyDefeatedMilestones(t *testing.T) {

	manager := NewAchievementManager()

	tests := []struct {
		totalDefeated int
		expectedID    string
		shouldUnlock  bool
	}{
		{9, AchievementDefeat10, false},
		{10, AchievementDefeat10, true},
		{49, AchievementDefeat50, false},
		{50, AchievementDefeat50, true},
		{99, AchievementDefeat100, false},
		{100, AchievementDefeat100, true},
		{499, AchievementDefeat500, false},
		{500, AchievementDefeat500, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			manager = NewAchievementManager()
			manager.CheckBattleAchievements(tt.totalDefeated, 0, false)
			unlocked := manager.IsUnlocked(tt.expectedID)

			if tt.shouldUnlock && !unlocked {
				t.Errorf("撃破数 %d で実績 %s が解除されるべきです", tt.totalDefeated, tt.expectedID)
			}
		})
	}
}

func TestBattleAchievements_LevelMilestones(t *testing.T) {

	manager := NewAchievementManager()

	tests := []struct {
		maxLevel     int
		expectedID   string
		shouldUnlock bool
	}{
		{9, AchievementLevel10, false},
		{10, AchievementLevel10, true},
		{24, AchievementLevel25, false},
		{25, AchievementLevel25, true},
		{49, AchievementLevel50, false},
		{50, AchievementLevel50, true},
		{99, AchievementLevel100, false},
		{100, AchievementLevel100, true},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			manager = NewAchievementManager()
			manager.CheckBattleAchievements(0, tt.maxLevel, false)
			unlocked := manager.IsUnlocked(tt.expectedID)

			if tt.shouldUnlock && !unlocked {
				t.Errorf("レベル %d で実績 %s が解除されるべきです", tt.maxLevel, tt.expectedID)
			}
		})
	}
}

func TestBattleAchievements_NoDamageClear(t *testing.T) {

	manager := NewAchievementManager()

	// ダメージを受けた場合は解除されない
	_ = manager.CheckBattleAchievements(1, 1, false)
	if manager.IsUnlocked(AchievementNoDamage) {
		t.Error("ダメージを受けた場合は実績が解除されるべきではありません")
	}

	// ノーダメージで解除
	manager = NewAchievementManager()
	notifications := manager.CheckBattleAchievements(1, 1, true)
	if !manager.IsUnlocked(AchievementNoDamage) {
		t.Error("ノーダメージクリアで実績が解除されるべきです")
	}
	if len(notifications) == 0 {
		t.Error("実績解除通知が返されるべきです")
	}
}

// ==================================================
// 実績の重複解除防止テスト
// ==================================================

func TestAchievement_NoDuplicateUnlock(t *testing.T) {
	manager := NewAchievementManager()

	// 1回目の解除
	notifications1 := manager.CheckTypingAchievements(50, 100)
	// 2回目は解除されない
	notifications2 := manager.CheckTypingAchievements(50, 100)

	// 2回目は通知がない（既に解除済み）
	has50WPM := false
	for _, n := range notifications2 {
		if n.AchievementID == AchievementWPM50 {
			has50WPM = true
		}
	}
	if has50WPM {
		t.Error("既に解除済みの実績は再度通知されるべきではありません")
	}
	if len(notifications1) == 0 {
		t.Error("初回は通知されるべきです")
	}
}

// ==================================================
// 実績一覧取得テスト
// ==================================================

func TestAchievement_GetAllAchievements(t *testing.T) {
	manager := NewAchievementManager()

	achievements := manager.GetAllAchievements()
	if len(achievements) == 0 {
		t.Error("実績一覧が返されるべきです")
	}

	// 各実績には必須フィールドがある
	for _, a := range achievements {
		if a.ID == "" {
			t.Error("実績IDが必要です")
		}
		if a.Name == "" {
			t.Error("実績名が必要です")
		}
		if a.Description == "" {
			t.Error("実績の説明が必要です")
		}
	}
}

func TestAchievement_CompletionRate(t *testing.T) {
	manager := NewAchievementManager()

	// 初期状態は0%
	if rate := manager.GetCompletionRate(); rate != 0.0 {
		t.Errorf("初期状態の達成率は0%%であるべきです: got %.1f%%", rate*100)
	}

	// いくつか解除
	manager.CheckTypingAchievements(50, 100) // WPM50 + PerfectAccuracy
	rate := manager.GetCompletionRate()
	if rate <= 0.0 {
		t.Error("実績解除後、達成率は0より大きいべきです")
	}
}

// ==================================================
// セーブ/ロードテスト
// ==================================================

func TestAchievement_SaveLoad(t *testing.T) {
	manager := NewAchievementManager()

	// いくつか解除
	manager.CheckTypingAchievements(80, 100)
	manager.CheckBattleAchievements(50, 25, true)

	// 解除済み実績IDを取得
	unlockedIDs := manager.GetUnlockedIDs()
	if len(unlockedIDs) == 0 {
		t.Error("解除済み実績がセーブデータに含まれるべきです")
	}

	// 新しいマネージャーにロード
	newManager := NewAchievementManager()
	newManager.LoadFromUnlockedIDs(unlockedIDs)

	// 解除状態が復元されている
	if !newManager.IsUnlocked(AchievementWPM50) {
		t.Error("WPM50実績がロードされるべきです")
	}
	if !newManager.IsUnlocked(AchievementWPM80) {
		t.Error("WPM80実績がロードされるべきです")
	}
	if !newManager.IsUnlocked(AchievementDefeat50) {
		t.Error("Defeat50実績がロードされるべきです")
	}
}
