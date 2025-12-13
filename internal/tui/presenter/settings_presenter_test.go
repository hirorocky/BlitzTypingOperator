package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/usecase/session"
)

// TestCreateSettingsData は設定データ作成をテストします。
func TestCreateSettingsData(t *testing.T) {
	gs := session.NewGameState()

	data := CreateSettingsData(gs)

	if data == nil {
		t.Fatal("CreateSettingsData returned nil")
	}

	// デフォルト設定の確認
	if data.SoundVolume != 100 {
		t.Errorf("SoundVolume expected 100, got %d", data.SoundVolume)
	}
	if data.Difficulty != "normal" {
		t.Errorf("Difficulty expected 'normal', got %s", data.Difficulty)
	}
	if data.Keybinds == nil {
		t.Error("Keybinds is nil")
	}
}

// TestCreateSettingsData_ModifiedSettings は変更後の設定データ作成をテストします。
func TestCreateSettingsData_ModifiedSettings(t *testing.T) {
	gs := session.NewGameState()

	// 設定を変更
	gs.Settings().SetSoundVolume(50)
	gs.Settings().SetDifficulty(session.DifficultyHard)
	gs.Settings().SetKeybind("custom_action", "x")

	data := CreateSettingsData(gs)

	if data.SoundVolume != 50 {
		t.Errorf("SoundVolume expected 50, got %d", data.SoundVolume)
	}
	if data.Difficulty != "hard" {
		t.Errorf("Difficulty expected 'hard', got %s", data.Difficulty)
	}
	if data.Keybinds["custom_action"] != "x" {
		t.Errorf("Keybind expected 'x', got %s", data.Keybinds["custom_action"])
	}
}
