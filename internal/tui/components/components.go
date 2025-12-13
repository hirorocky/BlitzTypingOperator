// Package components はTUI共通コンポーネントを提供します。
// メニュー、入力フィールド、情報パネルなどの再利用可能なUIコンポーネントを含みます。
package components

import (
	"fmt"
	"strconv"
	"strings"

	"hirorocky/type-battle/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// ==================== メニューコンポーネント ====================

// MenuItem はメニューアイテムを表します。
type MenuItem struct {
	Label    string
	Value    string
	Disabled bool
}

// Menu はメニューコンポーネントを表します。
type Menu struct {
	Items         []MenuItem
	SelectedIndex int
	Title         string
	styles        *styles.GameStyles
}

// NewMenu は新しいMenuを作成します。
func NewMenu(items []MenuItem) *Menu {
	return &Menu{
		Items:         items,
		SelectedIndex: 0,
		styles:        styles.NewGameStyles(),
	}
}

// NewMenuWithTitle はタイトル付きMenuを作成します。
func NewMenuWithTitle(title string, items []MenuItem) *Menu {
	m := NewMenu(items)
	m.Title = title
	return m
}

// MoveUp は選択を上に移動します。
func (m *Menu) MoveUp() {
	m.SelectedIndex--
	if m.SelectedIndex < 0 {
		m.SelectedIndex = len(m.Items) - 1
	}
	// 無効なアイテムをスキップ
	for m.Items[m.SelectedIndex].Disabled && m.SelectedIndex > 0 {
		m.SelectedIndex--
	}
}

// MoveDown は選択を下に移動します。
func (m *Menu) MoveDown() {
	m.SelectedIndex++
	if m.SelectedIndex >= len(m.Items) {
		m.SelectedIndex = 0
	}
	// 無効なアイテムをスキップ
	for m.Items[m.SelectedIndex].Disabled && m.SelectedIndex < len(m.Items)-1 {
		m.SelectedIndex++
	}
}

// GetSelected は現在選択されているアイテムを返します。
func (m *Menu) GetSelected() MenuItem {
	if m.SelectedIndex < 0 || m.SelectedIndex >= len(m.Items) {
		return MenuItem{}
	}
	return m.Items[m.SelectedIndex]
}

// Render はメニューをレンダリングします。
func (m *Menu) Render() string {
	var builder strings.Builder

	if m.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorPrimary).
			MarginBottom(1)
		builder.WriteString(titleStyle.Render(m.Title))
		builder.WriteString("\n\n")
	}

	for i, item := range m.Items {
		var style lipgloss.Style
		prefix := "  "

		if i == m.SelectedIndex {
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if item.Disabled {
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSubtle)
		} else {
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSecondary)
		}

		builder.WriteString(style.Render(prefix + item.Label))
		builder.WriteString("\n")
	}

	return builder.String()
}

// ==================== 入力フィールドコンポーネント ====================

// InputMode は入力モードを表します。
type InputMode int

const (
	// InputModeText はテキスト入力モードです。
	InputModeText InputMode = iota
	// InputModeNumeric は数値入力モードです。
	InputModeNumeric
)

// InputField は入力フィールドコンポーネントを表します。
type InputField struct {
	Value       string
	Placeholder string
	InputMode   InputMode
	MinValue    int
	MaxValue    int
	MaxLength   int
	Focused     bool
	Error       string
	styles      *styles.GameStyles
}

// NewInputField は新しいInputFieldを作成します。
func NewInputField(placeholder string) *InputField {
	return &InputField{
		Value:       "",
		Placeholder: placeholder,
		InputMode:   InputModeText,
		MinValue:    0,
		MaxValue:    0,
		MaxLength:   100,
		Focused:     true,
		styles:      styles.NewGameStyles(),
	}
}

// HandleInput は文字入力を処理します。
func (f *InputField) HandleInput(r rune) {
	// 最大長チェック
	if f.MaxLength > 0 && len(f.Value) >= f.MaxLength {
		return
	}

	// 入力モードに応じたフィルタリング
	if f.InputMode == InputModeNumeric {
		if r < '0' || r > '9' {
			return
		}
	}

	f.Value += string(r)
	f.Error = "" // 入力時にエラーをクリア
}

// HandleBackspace はバックスペースを処理します。
func (f *InputField) HandleBackspace() {
	if len(f.Value) > 0 {
		f.Value = f.Value[:len(f.Value)-1]
		f.Error = ""
	}
}

// Clear は入力をクリアします。
func (f *InputField) Clear() {
	f.Value = ""
	f.Error = ""
}

// GetIntValue は数値として値を取得します。
func (f *InputField) GetIntValue() (int, error) {
	if f.Value == "" {
		return 0, fmt.Errorf("値が入力されていません")
	}
	return strconv.Atoi(f.Value)
}

// Validate は入力を検証します。
func (f *InputField) Validate() (bool, string) {
	if f.Value == "" {
		return false, "値を入力してください"
	}

	if f.InputMode == InputModeNumeric {
		val, err := strconv.Atoi(f.Value)
		if err != nil {
			return false, "有効な数値を入力してください"
		}

		if f.MinValue > 0 && val < f.MinValue {
			return false, fmt.Sprintf("%d以上の値を入力してください", f.MinValue)
		}

		if f.MaxValue > 0 && val > f.MaxValue {
			return false, fmt.Sprintf("%d以下の値を入力してください", f.MaxValue)
		}
	}

	return true, ""
}

// Render は入力フィールドをレンダリングします。
func (f *InputField) Render(width int) string {
	var displayValue string
	if f.Value == "" && !f.Focused {
		displayValue = f.Placeholder
	} else {
		displayValue = f.Value
	}

	if f.Focused {
		displayValue += "_"
	}

	var style lipgloss.Style
	if f.Error != "" {
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorDamage).
			Padding(0, 1).
			Width(width)
	} else if f.Focused {
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorPrimary).
			Padding(0, 1).
			Width(width)
	} else {
		style = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.ColorSubtle).
			Padding(0, 1).
			Width(width)
	}

	result := style.Render(displayValue)

	if f.Error != "" {
		errorStyle := lipgloss.NewStyle().Foreground(styles.ColorDamage)
		result += "\n" + errorStyle.Render(f.Error)
	}

	return result
}

// ==================== 情報パネルコンポーネント ====================

// InfoItem は情報パネルのアイテムを表します。
type InfoItem struct {
	Label string
	Value string
}

// InfoPanel は情報パネルコンポーネントを表します。
type InfoPanel struct {
	Title  string
	Items  []InfoItem
	styles *styles.GameStyles
}

// NewInfoPanel は新しいInfoPanelを作成します。
func NewInfoPanel(title string) *InfoPanel {
	return &InfoPanel{
		Title:  title,
		Items:  make([]InfoItem, 0),
		styles: styles.NewGameStyles(),
	}
}

// AddItem はアイテムを追加します。
func (p *InfoPanel) AddItem(label, value string) {
	p.Items = append(p.Items, InfoItem{
		Label: label,
		Value: value,
	})
}

// ClearItems はアイテムをクリアします。
func (p *InfoPanel) ClearItems() {
	p.Items = make([]InfoItem, 0)
}

// Render は情報パネルをレンダリングします。
func (p *InfoPanel) Render(width int) string {
	var builder strings.Builder

	for _, item := range p.Items {
		labelStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSubtle)
		valueStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSecondary).
			Bold(true)

		builder.WriteString(labelStyle.Render(item.Label + ": "))
		builder.WriteString(valueStyle.Render(item.Value))
		builder.WriteString("\n")
	}

	content := strings.TrimSuffix(builder.String(), "\n")
	return p.styles.RenderBoxWithTitle(p.Title, content, width)
}

// ==================== リストコンポーネント ====================

// ListItem はリストアイテムを表します。
type ListItem struct {
	ID          string
	Title       string
	Description string
	Selected    bool
	Disabled    bool
}

// List はリストコンポーネントを表します。
type List struct {
	Items         []ListItem
	SelectedIndex int
	Title         string
	MaxVisible    int
	ScrollOffset  int
	styles        *styles.GameStyles
}

// NewList は新しいListを作成します。
func NewList(title string, maxVisible int) *List {
	return &List{
		Items:         make([]ListItem, 0),
		SelectedIndex: 0,
		Title:         title,
		MaxVisible:    maxVisible,
		ScrollOffset:  0,
		styles:        styles.NewGameStyles(),
	}
}

// AddItem はアイテムを追加します。
func (l *List) AddItem(item ListItem) {
	l.Items = append(l.Items, item)
}

// ClearItems はアイテムをクリアします。
func (l *List) ClearItems() {
	l.Items = make([]ListItem, 0)
	l.SelectedIndex = 0
	l.ScrollOffset = 0
}

// MoveUp は選択を上に移動します。
func (l *List) MoveUp() {
	if len(l.Items) == 0 {
		return
	}
	l.SelectedIndex--
	if l.SelectedIndex < 0 {
		l.SelectedIndex = len(l.Items) - 1
	}
	l.updateScroll()
}

// MoveDown は選択を下に移動します。
func (l *List) MoveDown() {
	if len(l.Items) == 0 {
		return
	}
	l.SelectedIndex++
	if l.SelectedIndex >= len(l.Items) {
		l.SelectedIndex = 0
	}
	l.updateScroll()
}

// updateScroll はスクロール位置を更新します。
func (l *List) updateScroll() {
	if l.SelectedIndex < l.ScrollOffset {
		l.ScrollOffset = l.SelectedIndex
	}
	if l.SelectedIndex >= l.ScrollOffset+l.MaxVisible {
		l.ScrollOffset = l.SelectedIndex - l.MaxVisible + 1
	}
}

// GetSelected は選択されているアイテムを返します。
func (l *List) GetSelected() *ListItem {
	if l.SelectedIndex < 0 || l.SelectedIndex >= len(l.Items) {
		return nil
	}
	return &l.Items[l.SelectedIndex]
}

// Render はリストをレンダリングします。
func (l *List) Render(width int) string {
	var builder strings.Builder

	if l.Title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorPrimary).
			MarginBottom(1)
		builder.WriteString(titleStyle.Render(l.Title))
		builder.WriteString("\n\n")
	}

	if len(l.Items) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSubtle).
			Italic(true)
		builder.WriteString(emptyStyle.Render("(アイテムがありません)"))
		return builder.String()
	}

	// 表示範囲を計算
	endIndex := l.ScrollOffset + l.MaxVisible
	if endIndex > len(l.Items) {
		endIndex = len(l.Items)
	}

	for i := l.ScrollOffset; i < endIndex; i++ {
		item := l.Items[i]
		var style lipgloss.Style
		prefix := "  "

		if i == l.SelectedIndex {
			style = lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if item.Disabled {
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSubtle)
		} else {
			style = lipgloss.NewStyle().
				Foreground(styles.ColorSecondary)
		}

		builder.WriteString(style.Render(prefix + item.Title))
		if item.Description != "" {
			descStyle := lipgloss.NewStyle().
				Foreground(styles.ColorSubtle)
			builder.WriteString(descStyle.Render(" - " + item.Description))
		}
		builder.WriteString("\n")
	}

	// スクロールインジケーター
	if len(l.Items) > l.MaxVisible {
		indicator := fmt.Sprintf("(%d/%d)", l.SelectedIndex+1, len(l.Items))
		indicatorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSubtle)
		builder.WriteString("\n")
		builder.WriteString(indicatorStyle.Render(indicator))
	}

	return builder.String()
}

// ==================== AgentCardコンポーネント ====================

// AgentCardStyle はカードのスタイルバリエーションです。

type AgentCardStyle int

const (
	// AgentCardCompact はコンパクト（横並び用）スタイルです。
	AgentCardCompact AgentCardStyle = iota
	// AgentCardDetailed は詳細（単体表示用）スタイルです。
	AgentCardDetailed
)

// AgentCard はエージェント情報カードを表します。

type AgentCard struct {
	// AgentName はエージェント名です（空の場合は空スロット）
	AgentName string
	// AgentLevel はエージェントのレベルです
	AgentLevel int
	// CoreTypeName はコア特性の名前です
	CoreTypeName string
	// ModuleIcons はモジュールのアイコンリストです
	ModuleIcons []string
	// Style はカードのスタイルです
	Style AgentCardStyle
	// Selected は選択状態かどうかです
	Selected bool
	// ShowHP はHP表示を行うかどうかです
	ShowHP bool
	// CurrentHP は現在のHP値です（ShowHP=true時に使用）
	CurrentHP int
	// MaxHP は最大HP値です
	MaxHP int
	// gameStyles はゲームスタイルです
	gameStyles *styles.GameStyles
}

// NewAgentCard は新しいAgentCardを作成します。
// agentがnilの場合は空スロット表示用のカードを作成します。
func NewAgentCard(agent interface{}, style AgentCardStyle) *AgentCard {
	return &AgentCard{
		Style:      style,
		gameStyles: styles.NewGameStyles(),
	}
}

// SetSelected は選択状態を設定します。
func (c *AgentCard) SetSelected(selected bool) {
	c.Selected = selected
}

// SetHP はHP表示を設定します。
func (c *AgentCard) SetHP(current, max int) {
	c.ShowHP = true
	c.CurrentHP = current
	c.MaxHP = max
}

// Render はカードをレンダリングします。
func (c *AgentCard) Render(width int) string {
	// 空スロットの場合
	if c.AgentName == "" {
		return c.renderEmptySlot(width)
	}

	switch c.Style {
	case AgentCardDetailed:
		return c.renderDetailed(width)
	default:
		return c.renderCompact(width)
	}
}

// renderEmptySlot は空スロットを描画します。
func (c *AgentCard) renderEmptySlot(width int) string {
	var content strings.Builder
	content.WriteString("(空)\n")
	content.WriteString("\n")
	content.WriteString("Enterで装備")

	borderColor := styles.ColorSubtle
	if c.Selected {
		borderColor = styles.ColorPrimary
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Foreground(styles.ColorSubtle).
		Width(width-2).
		Padding(0, 1).
		Align(lipgloss.Center)

	return style.Render(content.String())
}

// renderCompact はコンパクトスタイルを描画します。
func (c *AgentCard) renderCompact(width int) string {
	var content strings.Builder

	// エージェント名とレベル
	nameStyle := lipgloss.NewStyle().Bold(true)
	content.WriteString(nameStyle.Render(fmt.Sprintf("%s Lv.%d", c.AgentName, c.AgentLevel)))
	content.WriteString("\n")

	// モジュールアイコン
	if len(c.ModuleIcons) > 0 {
		content.WriteString(strings.Join(c.ModuleIcons, ""))
		content.WriteString("\n")
	}

	// HP表示（オプション）
	if c.ShowHP {
		hpBar := c.gameStyles.RenderHPBar(c.CurrentHP, c.MaxHP, width-6)
		content.WriteString(hpBar)
	}

	borderColor := styles.ColorSubtle
	if c.Selected {
		borderColor = styles.ColorPrimary
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width-2).
		Padding(0, 1)

	return style.Render(content.String())
}

// renderDetailed は詳細スタイルを描画します。
func (c *AgentCard) renderDetailed(width int) string {
	var content strings.Builder

	// エージェント名とレベル
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	content.WriteString(nameStyle.Render(fmt.Sprintf("%s Lv.%d", c.AgentName, c.AgentLevel)))
	content.WriteString("\n")

	// コアタイプ
	if c.CoreTypeName != "" {
		typeStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
		content.WriteString(typeStyle.Render(fmt.Sprintf("コアタイプ: %s", c.CoreTypeName)))
		content.WriteString("\n")
	}

	// 区切り線
	divider := strings.Repeat("─", width-6)
	dividerStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	content.WriteString(dividerStyle.Render(divider))
	content.WriteString("\n")

	// モジュールアイコン
	if len(c.ModuleIcons) > 0 {
		content.WriteString("モジュール: ")
		content.WriteString(strings.Join(c.ModuleIcons, " "))
		content.WriteString("\n")
	}

	// HP表示（オプション）
	if c.ShowHP {
		content.WriteString("\n")
		hpBar := c.gameStyles.RenderHPBarWithValue(c.CurrentHP, c.MaxHP, width-10)
		content.WriteString(hpBar)
	}

	borderColor := styles.ColorSubtle
	if c.Selected {
		borderColor = styles.ColorPrimary
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width-2).
		Padding(0, 1)

	return style.Render(content.String())
}

// ==================== ConfirmDialogコンポーネント ====================

// ConfirmResult はダイアログの結果を表します。

type ConfirmResult int

const (
	// ConfirmResultNone は結果がまだないことを表します
	ConfirmResultNone ConfirmResult = iota
	// ConfirmResultYes は「はい」が選択されたことを表します
	ConfirmResultYes
	// ConfirmResultNo は「いいえ」が選択されたことを表します
	ConfirmResultNo
	// ConfirmResultCancelled はキャンセルされたことを表します
	ConfirmResultCancelled
)

// ConfirmDialog は確認ダイアログを表します。

type ConfirmDialog struct {
	// Title はダイアログのタイトルです
	Title string
	// Message はダイアログのメッセージです
	Message string
	// Visible は表示中かどうかです
	Visible bool
	// SelectedYes は「はい」が選択されているかどうかです（falseなら「いいえ」）
	SelectedYes bool
	// Result はダイアログの結果です
	Result ConfirmResult
	// gameStyles はゲームスタイルです
	gameStyles *styles.GameStyles
}

// NewConfirmDialog は新しいConfirmDialogを作成します。
func NewConfirmDialog(title, message string) *ConfirmDialog {
	return &ConfirmDialog{
		Title:       title,
		Message:     message,
		Visible:     false,
		SelectedYes: false, // 初期選択は「いいえ」（安全側）
		Result:      ConfirmResultNone,
		gameStyles:  styles.NewGameStyles(),
	}
}

// Show はダイアログを表示します。
func (d *ConfirmDialog) Show() {
	d.Visible = true
	d.SelectedYes = false // 表示時に「いいえ」にリセット
	d.Result = ConfirmResultNone
}

// Hide はダイアログを非表示にします。
func (d *ConfirmDialog) Hide() {
	d.Visible = false
}

// HandleKey はキー入力を処理します。
// 左右キーで選択、Enterで決定、Escでキャンセル
func (d *ConfirmDialog) HandleKey(key string) ConfirmResult {
	if !d.Visible {
		return ConfirmResultNone
	}

	switch key {
	case "left", "h":
		d.SelectedYes = true
		return ConfirmResultNone
	case "right", "l":
		d.SelectedYes = false
		return ConfirmResultNone
	case "enter":
		if d.SelectedYes {
			d.Result = ConfirmResultYes
		} else {
			d.Result = ConfirmResultNo
		}
		d.Hide()
		return d.Result
	case "esc", "escape":
		d.Result = ConfirmResultCancelled
		d.Hide()
		return d.Result
	}

	return ConfirmResultNone
}

// Render はダイアログをレンダリングします。
// screenWidth, screenHeightは画面サイズ（中央配置に使用）
func (d *ConfirmDialog) Render(screenWidth, screenHeight int) string {
	if !d.Visible {
		return ""
	}

	// ダイアログの幅を計算
	dialogWidth := 50
	if len(d.Message) > dialogWidth-4 {
		dialogWidth = len(d.Message) + 8
	}
	if dialogWidth > screenWidth-4 {
		dialogWidth = screenWidth - 4
	}

	var content strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary)
	content.WriteString(titleStyle.Render(d.Title))
	content.WriteString("\n\n")

	// メッセージ
	messageStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSecondary)
	content.WriteString(messageStyle.Render(d.Message))
	content.WriteString("\n")
	content.WriteString(messageStyle.Render("この操作は取り消せません。"))
	content.WriteString("\n\n")

	// ボタン
	yesStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Bold(d.SelectedYes)
	noStyle := lipgloss.NewStyle().
		Padding(0, 2).
		Bold(!d.SelectedYes)

	if d.SelectedYes {
		yesStyle = yesStyle.
			Background(styles.ColorPrimary).
			Foreground(lipgloss.Color("#000000"))
		noStyle = noStyle.
			Foreground(styles.ColorSubtle)
	} else {
		yesStyle = yesStyle.
			Foreground(styles.ColorSubtle)
		noStyle = noStyle.
			Background(styles.ColorPrimary).
			Foreground(lipgloss.Color("#000000"))
	}

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		yesStyle.Render("はい"),
		"        ",
		noStyle.Render("いいえ"),
	)
	content.WriteString(buttons)

	// ダイアログボックス
	dialogStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Width(dialogWidth).
		Align(lipgloss.Center)

	return dialogStyle.Render(content.String())
}
