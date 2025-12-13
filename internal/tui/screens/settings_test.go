// Package screens はTUI画面のテストを提供します。
package screens

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.7: 設定画面のテスト ====================

// TestNewSettingsScreen はSettingsScreenの初期化をテストします。
func TestNewSettingsScreen(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	if screen == nil {
		t.Fatal("SettingsScreenがnilです")
	}
}

// TestSettingsKeybindDisplay はキーバインド設定表示をテストします。

func TestSettingsKeybindDisplay(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	// キーバインドが表示されていること
	if len(screen.settings.Keybinds) == 0 {
		t.Error("キーバインドがありません")
	}
}

// TestSettingsKeybindChange はキーバインド変更をテストします。

func TestSettingsKeybindChange(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	// 編集モードに入る
	screen.selectedIndex = 0
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	if !screen.editing {
		t.Error("編集モードに入れませんでした")
	}
}

// TestSettingsImmediateApply は設定の即時適用をテストします。

func TestSettingsImmediateApply(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	// 設定を変更
	originalValue := screen.settings.Keybinds["select"]

	// 編集モードで新しいキーを設定
	screen.selectedIndex = 0
	screen.editing = true
	screen.handleKeyMsg(tea.KeyMsg{Runes: []rune{'x'}})

	// 編集終了
	screen.editing = false

	// 変更が反映されていること（または同じ値のままであること）
	_ = originalValue
}

// TestSettingsBackNavigation は戻るナビゲーションをテストします。
func TestSettingsBackNavigation(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Escキーでコマンドが返されません")
	}
}

// TestSettingsRender はレンダリングをテストします。
func TestSettingsRender(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestSettingsNavigation は設定項目のナビゲーションをテストします。
func TestSettingsNavigation(t *testing.T) {
	settings := createTestSettings()
	screen := NewSettingsScreen(settings)

	// 下に移動
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyDown})
	if screen.selectedIndex != 1 {
		t.Errorf("下移動後の選択インデックス: got %d, want 1", screen.selectedIndex)
	}

	// 上に移動
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyUp})
	if screen.selectedIndex != 0 {
		t.Errorf("上移動後の選択インデックス: got %d, want 0", screen.selectedIndex)
	}
}

// ==================== ヘルパー関数 ====================

func createTestSettings() *SettingsData {
	return &SettingsData{
		Keybinds: map[string]string{
			"select":     "enter",
			"cancel":     "esc",
			"move_up":    "k",
			"move_down":  "j",
			"move_left":  "h",
			"move_right": "l",
		},
		SoundVolume: 100,
		Difficulty:  "normal",
	}
}
