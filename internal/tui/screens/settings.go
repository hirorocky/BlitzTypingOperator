// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.7: 設定画面 ====================

// SettingsScreen は設定画面を表します。

type SettingsScreen struct {
	settings      *SettingsData
	selectedIndex int
	editing       bool
	editKey       string
	styles        *styles.GameStyles
	width         int
	height        int
}

// KeybindItem はキーバインド項目を表します。
type KeybindItem struct {
	ID    string
	Label string
	Key   string
}

// NewSettingsScreen は新しいSettingsScreenを作成します。
func NewSettingsScreen(settings *SettingsData) *SettingsScreen {
	return &SettingsScreen{
		settings:      settings,
		selectedIndex: 0,
		editing:       false,
		styles:        styles.NewGameStyles(),
		width:         140,
		height:        40,
	}
}

// Init は画面の初期化を行います。
func (s *SettingsScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return s, nil
}

// handleKeyMsg はキーボード入力を処理します。
func (s *SettingsScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 編集モード中
	if s.editing {
		return s.handleEditingMode(msg)
	}

	// 通常モード
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "up", "k":
		s.moveUp()
	case "down", "j":
		s.moveDown()
	case "enter":
		s.startEditing()
	}

	return s, nil
}

// handleEditingMode は編集モードの入力を処理します。
func (s *SettingsScreen) handleEditingMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// 編集キャンセル
		s.editing = false
		s.editKey = ""
		return s, nil
	default:
		// 新しいキーを設定
		newKey := msg.String()
		if newKey != "" && newKey != "enter" {
			s.applyKeybindChange(newKey)
			s.editing = false
			s.editKey = ""
		}
	}
	return s, nil
}

// moveUp は選択を上に移動します。
func (s *SettingsScreen) moveUp() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// moveDown は選択を下に移動します。
func (s *SettingsScreen) moveDown() {
	maxIndex := len(s.getKeybindItems())
	if s.selectedIndex < maxIndex-1 {
		s.selectedIndex++
	}
}

// startEditing は編集モードを開始します。
func (s *SettingsScreen) startEditing() {
	items := s.getKeybindItems()
	if s.selectedIndex < len(items) {
		s.editing = true
		s.editKey = items[s.selectedIndex].ID
	}
}

// applyKeybindChange はキーバインドの変更を適用します。

func (s *SettingsScreen) applyKeybindChange(newKey string) {
	if s.editKey != "" {
		s.settings.Keybinds[s.editKey] = newKey
	}
}

// getKeybindItems はキーバインド項目のリストを返します。
func (s *SettingsScreen) getKeybindItems() []KeybindItem {
	return []KeybindItem{
		{ID: "select", Label: "決定", Key: s.settings.Keybinds["select"]},
		{ID: "cancel", Label: "キャンセル", Key: s.settings.Keybinds["cancel"]},
		{ID: "move_up", Label: "上移動", Key: s.settings.Keybinds["move_up"]},
		{ID: "move_down", Label: "下移動", Key: s.settings.Keybinds["move_down"]},
		{ID: "move_left", Label: "左移動", Key: s.settings.Keybinds["move_left"]},
		{ID: "move_right", Label: "右移動", Key: s.settings.Keybinds["move_right"]},
	}
}

// View は画面をレンダリングします。
func (s *SettingsScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("設定"))
	builder.WriteString("\n\n")

	// キーバインド設定
	builder.WriteString(s.renderKeybindSettings())
	builder.WriteString("\n\n")

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hints string
	if s.editing {
		hints = "新しいキーを押してください  Esc: キャンセル"
	} else {
		hints = "↑/↓: 選択  Enter: 変更  Esc: 戻る"
	}
	builder.WriteString(hintStyle.Render(hints))

	return builder.String()
}

// renderKeybindSettings はキーバインド設定をレンダリングします。

func (s *SettingsScreen) renderKeybindSettings() string {
	var builder strings.Builder

	sectionTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary).
		Render("キーバインド設定")

	builder.WriteString(lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(sectionTitle))
	builder.WriteString("\n\n")

	// キーバインドリスト
	items := s.getKeybindItems()
	var itemLines []string

	for i, item := range items {
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			prefix = "> "
			if s.editing {
				style = style.Bold(true).Foreground(styles.ColorWarning)
			} else {
				style = style.Bold(true).
					Foreground(styles.ColorSelectedFg).
					Background(styles.ColorSelectedBg)
			}
		}

		keyDisplay := item.Key
		if i == s.selectedIndex && s.editing {
			keyDisplay = "_"
		}

		line := fmt.Sprintf("%s%-12s : %s", prefix, item.Label, keyDisplay)
		itemLines = append(itemLines, style.Render(line))
	}

	listContent := strings.Join(itemLines, "\n")
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Width(40).
		Render(listContent)

	builder.WriteString(lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box))

	return builder.String()
}

// ==================== Screenインターフェース実装 ====================

// SetSize は画面サイズを設定します。
// Screenインターフェースの実装です。
func (s *SettingsScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// GetTitle は画面のタイトルを返します。
// Screenインターフェースの実装です。
func (s *SettingsScreen) GetTitle() string {
	return "設定"
}

// GetSize は現在の画面サイズを返します。
func (s *SettingsScreen) GetSize() (width, height int) {
	return s.width, s.height
}
