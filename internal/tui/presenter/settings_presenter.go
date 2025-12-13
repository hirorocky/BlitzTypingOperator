package presenter

import (
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/usecase/session"
)

// CreateSettingsData はGameStateから設定データを生成します。
func CreateSettingsData(gs *session.GameState) *screens.SettingsData {
	settings := gs.Settings()
	return &screens.SettingsData{
		Keybinds:    settings.Keybinds(),
		SoundVolume: settings.SoundVolume(),
		Difficulty:  string(settings.Difficulty()),
	}
}
