// Package savedata は実績データの永続化アダプターをテストします。
package savedata

import (
	"testing"
)

// TestAchievementAdapter_ToSaveData は実績データの変換をテストします。
func TestAchievementAdapter_ToSaveData(t *testing.T) {
	// 解除済み実績のリスト
	unlocked := []string{"wpm_50", "wpm_80", "defeat_10"}

	saveData := AchievementStateToSaveData(unlocked)

	if saveData == nil {
		t.Fatal("AchievementStateToSaveData returned nil")
	}

	if len(saveData.Unlocked) != 3 {
		t.Errorf("Expected 3 unlocked achievements, got %d", len(saveData.Unlocked))
	}

	// 各実績が含まれていることを確認
	unlockedMap := make(map[string]bool)
	for _, id := range saveData.Unlocked {
		unlockedMap[id] = true
	}

	if !unlockedMap["wpm_50"] {
		t.Error("wpm_50 should be in unlocked")
	}
	if !unlockedMap["wpm_80"] {
		t.Error("wpm_80 should be in unlocked")
	}
	if !unlockedMap["defeat_10"] {
		t.Error("defeat_10 should be in unlocked")
	}
}

// TestAchievementAdapter_FromSaveData はセーブデータからの復元をテストします。
func TestAchievementAdapter_FromSaveData(t *testing.T) {
	saveData := &AchievementsSaveData{
		Unlocked: []string{"wpm_100", "level_25"},
		Progress: make(map[string]int),
	}

	unlocked := SaveDataToAchievementState(saveData)

	if len(unlocked) != 2 {
		t.Errorf("Expected 2 unlocked achievements, got %d", len(unlocked))
	}

	// 各実績が含まれていることを確認
	unlockedMap := make(map[string]bool)
	for _, id := range unlocked {
		unlockedMap[id] = true
	}

	if !unlockedMap["wpm_100"] {
		t.Error("wpm_100 should be in unlocked")
	}
	if !unlockedMap["level_25"] {
		t.Error("level_25 should be in unlocked")
	}
}

// TestAchievementAdapter_NilSaveData はnilセーブデータの処理をテストします。
func TestAchievementAdapter_NilSaveData(t *testing.T) {
	unlocked := SaveDataToAchievementState(nil)

	if unlocked == nil {
		t.Error("SaveDataToAchievementState should return empty slice, not nil")
	}

	if len(unlocked) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(unlocked))
	}
}

// TestAchievementAdapter_EmptyUnlocked は空の解除リストをテストします。
func TestAchievementAdapter_EmptyUnlocked(t *testing.T) {
	saveData := AchievementStateToSaveData([]string{})

	if saveData == nil {
		t.Fatal("AchievementStateToSaveData returned nil")
	}

	if len(saveData.Unlocked) != 0 {
		t.Errorf("Expected empty unlocked list, got %d", len(saveData.Unlocked))
	}

	if saveData.Progress == nil {
		t.Error("Progress should not be nil")
	}
}
