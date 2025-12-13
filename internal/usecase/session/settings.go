package session

// Settings はゲームの設定を管理する構造体です。
type Settings struct {
	// keybinds はキーバインド設定です。
	keybinds map[string]string

	// soundVolume は音量設定（0-100）です。
	soundVolume int

	// difficulty は難易度設定です。
	difficulty Difficulty
}

// Difficulty は難易度を表す型です。
type Difficulty string

const (
	// DifficultyEasy は簡単モードです。
	DifficultyEasy Difficulty = "easy"

	// DifficultyNormal は通常モードです。
	DifficultyNormal Difficulty = "normal"

	// DifficultyHard は難しいモードです。
	DifficultyHard Difficulty = "hard"
)

// DefaultKeybinds はデフォルトのキーバインド設定です。
var DefaultKeybinds = map[string]string{
	"select":     "enter",
	"cancel":     "esc",
	"move_up":    "k",
	"move_down":  "j",
	"move_left":  "h",
	"move_right": "l",
}

// NewSettings はデフォルト値で新しいSettingsを作成します。
func NewSettings() *Settings {
	keybinds := make(map[string]string)
	for k, v := range DefaultKeybinds {
		keybinds[k] = v
	}

	return &Settings{
		keybinds:    keybinds,
		soundVolume: 100,
		difficulty:  DifficultyNormal,
	}
}

// Keybinds はキーバインド設定を返します。
func (s *Settings) Keybinds() map[string]string {
	return s.keybinds
}

// SetKeybind はキーバインドを設定します。
func (s *Settings) SetKeybind(action, key string) {
	s.keybinds[action] = key
}

// GetKeybind は指定されたアクションのキーバインドを返します。
func (s *Settings) GetKeybind(action string) string {
	if key, ok := s.keybinds[action]; ok {
		return key
	}
	return ""
}

// SoundVolume は音量を返します。
func (s *Settings) SoundVolume() int {
	return s.soundVolume
}

// SetSoundVolume は音量を設定します。
func (s *Settings) SetSoundVolume(volume int) {
	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}
	s.soundVolume = volume
}

// Difficulty は難易度を返します。
func (s *Settings) Difficulty() Difficulty {
	return s.difficulty
}

// SetDifficulty は難易度を設定します。
func (s *Settings) SetDifficulty(difficulty Difficulty) {
	s.difficulty = difficulty
}

// ToScreensSettingsData は画面用のSettingsDataに変換します。
func (s *Settings) ToScreensSettingsData() map[string]interface{} {
	return map[string]interface{}{
		"keybinds":     s.keybinds,
		"sound_volume": s.soundVolume,
		"difficulty":   string(s.difficulty),
	}
}
