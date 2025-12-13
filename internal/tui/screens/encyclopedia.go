// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.5: 図鑑画面 ====================

// EncyclopediaCategory は図鑑カテゴリを表します。
type EncyclopediaCategory int

const (
	// CategoryCore はコア図鑑です。
	CategoryCore EncyclopediaCategory = iota
	// CategoryModule はモジュール図鑑です。
	CategoryModule
	// CategoryEnemy は敵図鑑です。
	CategoryEnemy
)

// EncyclopediaScreen は図鑑画面を表します。

type EncyclopediaScreen struct {
	data            *EncyclopediaData
	currentCategory EncyclopediaCategory
	selectedIndex   int
	styles          *styles.GameStyles
	width           int
	height          int
}

// NewEncyclopediaScreen は新しいEncyclopediaScreenを作成します。
func NewEncyclopediaScreen(data *EncyclopediaData) *EncyclopediaScreen {
	return &EncyclopediaScreen{
		data:            data,
		currentCategory: CategoryCore,
		selectedIndex:   0,
		styles:          styles.NewGameStyles(),
		width:           140,
		height:          40,
	}
}

// Init は画面の初期化を行います。
func (s *EncyclopediaScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *EncyclopediaScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *EncyclopediaScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "left", "h":
		s.prevCategory()
	case "right", "l":
		s.nextCategory()
	case "up", "k":
		s.moveUp()
	case "down", "j":
		s.moveDown()
	}

	return s, nil
}

// prevCategory は前のカテゴリに移動します。
func (s *EncyclopediaScreen) prevCategory() {
	if s.currentCategory > CategoryCore {
		s.currentCategory--
		s.selectedIndex = 0
	}
}

// nextCategory は次のカテゴリに移動します。
func (s *EncyclopediaScreen) nextCategory() {
	if s.currentCategory < CategoryEnemy {
		s.currentCategory++
		s.selectedIndex = 0
	}
}

// moveUp は選択を上に移動します。
func (s *EncyclopediaScreen) moveUp() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// moveDown は選択を下に移動します。
func (s *EncyclopediaScreen) moveDown() {
	maxIndex := s.getMaxIndex()
	if s.selectedIndex < maxIndex-1 {
		s.selectedIndex++
	}
}

// getMaxIndex は現在のカテゴリの最大インデックスを返します。
func (s *EncyclopediaScreen) getMaxIndex() int {
	switch s.currentCategory {
	case CategoryCore:
		return len(s.data.AllCoreTypes)
	case CategoryModule:
		return len(s.data.AllModuleTypes)
	case CategoryEnemy:
		return len(s.data.AllEnemyTypes)
	}
	return 0
}

// isCoreTypeAcquired はコア特性が獲得済みかを返します。
func (s *EncyclopediaScreen) isCoreTypeAcquired(id string) bool {
	for _, acquired := range s.data.AcquiredCoreTypes {
		if acquired == id {
			return true
		}
	}
	return false
}

// isModuleTypeAcquired はモジュールタイプが獲得済みかを返します。
func (s *EncyclopediaScreen) isModuleTypeAcquired(id string) bool {
	for _, acquired := range s.data.AcquiredModuleTypes {
		if acquired == id {
			return true
		}
	}
	return false
}

// isEnemyEncountered は敵が遭遇済みかを返します。
func (s *EncyclopediaScreen) isEnemyEncountered(id string) bool {
	for _, encountered := range s.data.EncounteredEnemies {
		if encountered == id {
			return true
		}
	}
	return false
}

// getCoreDisplayName はコアの表示名を返します（未獲得は???）。
func (s *EncyclopediaScreen) getCoreDisplayName(ct domain.CoreType) string {
	if s.isCoreTypeAcquired(ct.ID) {
		return ct.Name
	}
	return "???"
}

// getModuleDisplayName はモジュールの表示名を返します（未獲得は???）。
func (s *EncyclopediaScreen) getModuleDisplayName(mt ModuleTypeInfo) string {
	if s.isModuleTypeAcquired(mt.ID) {
		return mt.Name
	}
	return "???"
}

// getEnemyDisplayName は敵の表示名を返します（未遭遇は???）。
func (s *EncyclopediaScreen) getEnemyDisplayName(et domain.EnemyType) string {
	if s.isEnemyEncountered(et.ID) {
		return et.Name
	}
	return "???"
}

// getCoreCompletionRate はコア図鑑のコンプリート率を返します。
func (s *EncyclopediaScreen) getCoreCompletionRate() int {
	if len(s.data.AllCoreTypes) == 0 {
		return 0
	}
	return len(s.data.AcquiredCoreTypes) * 100 / len(s.data.AllCoreTypes)
}

// getModuleCompletionRate はモジュール図鑑のコンプリート率を返します。
func (s *EncyclopediaScreen) getModuleCompletionRate() int {
	if len(s.data.AllModuleTypes) == 0 {
		return 0
	}
	return len(s.data.AcquiredModuleTypes) * 100 / len(s.data.AllModuleTypes)
}

// getEnemyCompletionRate は敵図鑑のコンプリート率を返します。
func (s *EncyclopediaScreen) getEnemyCompletionRate() int {
	if len(s.data.AllEnemyTypes) == 0 {
		return 0
	}
	return len(s.data.EncounteredEnemies) * 100 / len(s.data.AllEnemyTypes)
}

// View は画面をレンダリングします。
func (s *EncyclopediaScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("図鑑"))
	builder.WriteString("\n\n")

	// カテゴリタブ
	builder.WriteString(s.renderCategoryTabs())
	builder.WriteString("\n\n")

	// メインコンテンツ
	builder.WriteString(s.renderMainContent())
	builder.WriteString("\n\n")

	// コンプリート率
	builder.WriteString(s.renderCompletionRate())
	builder.WriteString("\n\n")

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	hints := "←/→: カテゴリ切替  ↑/↓: 選択  Esc: 戻る"
	builder.WriteString(hintStyle.Render(hints))

	return builder.String()
}

// renderCategoryTabs はカテゴリタブをレンダリングします。
func (s *EncyclopediaScreen) renderCategoryTabs() string {
	categories := []string{"コア図鑑", "モジュール図鑑", "敵図鑑"}

	var tabItems []string
	for i, cat := range categories {
		style := lipgloss.NewStyle().Padding(0, 2)
		prefix := "  "
		if EncyclopediaCategory(i) == s.currentCategory {
			prefix = "> "
			style = style.
				Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
		} else {
			style = style.Foreground(styles.ColorSubtle)
		}
		tabItems = append(tabItems, style.Render(prefix+cat))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Center, tabItems...)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(tabBar)
}

// renderMainContent はメインコンテンツをレンダリングします。
func (s *EncyclopediaScreen) renderMainContent() string {
	switch s.currentCategory {
	case CategoryCore:
		return s.renderCoreEncyclopedia()
	case CategoryModule:
		return s.renderModuleEncyclopedia()
	case CategoryEnemy:
		return s.renderEnemyEncyclopedia()
	}
	return ""
}

// renderCoreEncyclopedia はコア図鑑をレンダリングします。
func (s *EncyclopediaScreen) renderCoreEncyclopedia() string {
	if len(s.data.AllCoreTypes) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("コア特性がありません")
	}

	// リストとプレビューを横に並べる
	listContent := s.renderCoreList()
	previewContent := s.renderCorePreview()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(40).
		Render(listContent)

	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(50).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listBox, "  ", previewBox)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderCoreList はコア図鑑のリストをレンダリングします。
func (s *EncyclopediaScreen) renderCoreList() string {
	var items []string
	for i, ct := range s.data.AllCoreTypes {
		acquired := s.isCoreTypeAcquired(ct.ID)
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if !acquired {
			style = style.Foreground(styles.ColorSubtle)
		}

		displayName := s.getCoreDisplayName(ct)
		status := ""
		if acquired {
			status = " [獲得済み]"
		}
		items = append(items, style.Render(prefix+displayName+status))
	}
	return strings.Join(items, "\n")
}

// renderCorePreview はコア図鑑のプレビューをレンダリングします。
func (s *EncyclopediaScreen) renderCorePreview() string {
	if s.selectedIndex >= len(s.data.AllCoreTypes) {
		return "選択してください"
	}

	ct := s.data.AllCoreTypes[s.selectedIndex]
	acquired := s.isCoreTypeAcquired(ct.ID)

	if !acquired {
		return "このコア特性はまだ獲得していません\n\n敵を倒してドロップを獲得しましょう"
	}

	panel := components.NewInfoPanel(ct.Name)
	panel.AddItem("ID", ct.ID)
	panel.AddItem("STR重み", fmt.Sprintf("%.1f", ct.StatWeights["STR"]))
	panel.AddItem("MAG重み", fmt.Sprintf("%.1f", ct.StatWeights["MAG"]))
	panel.AddItem("SPD重み", fmt.Sprintf("%.1f", ct.StatWeights["SPD"]))
	panel.AddItem("LUK重み", fmt.Sprintf("%.1f", ct.StatWeights["LUK"]))

	return panel.Render(45)
}

// renderModuleEncyclopedia はモジュール図鑑をレンダリングします。
func (s *EncyclopediaScreen) renderModuleEncyclopedia() string {
	if len(s.data.AllModuleTypes) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("モジュールタイプがありません")
	}

	listContent := s.renderModuleList()
	previewContent := s.renderModulePreviewEncyclopedia()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(40).
		Render(listContent)

	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(50).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listBox, "  ", previewBox)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderModuleList はモジュール図鑑のリストをレンダリングします。
func (s *EncyclopediaScreen) renderModuleList() string {
	var items []string
	for i, mt := range s.data.AllModuleTypes {
		acquired := s.isModuleTypeAcquired(mt.ID)
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if !acquired {
			style = style.Foreground(styles.ColorSubtle)
		}

		displayName := s.getModuleDisplayName(mt)
		status := ""
		if acquired {
			status = " [獲得済み]"
		}
		items = append(items, style.Render(prefix+displayName+status))
	}
	return strings.Join(items, "\n")
}

// renderModulePreviewEncyclopedia はモジュール図鑑のプレビューをレンダリングします。
func (s *EncyclopediaScreen) renderModulePreviewEncyclopedia() string {
	if s.selectedIndex >= len(s.data.AllModuleTypes) {
		return "選択してください"
	}

	mt := s.data.AllModuleTypes[s.selectedIndex]
	acquired := s.isModuleTypeAcquired(mt.ID)

	if !acquired {
		return "このモジュールはまだ獲得していません\n\n敵を倒してドロップを獲得しましょう"
	}

	panel := components.NewInfoPanel(mt.Name)
	panel.AddItem("カテゴリ", mt.Category.String())
	panel.AddItem("レベル", fmt.Sprintf("Lv.%d", mt.Level))
	if mt.Description != "" {
		panel.AddItem("説明", mt.Description)
	}

	return panel.Render(45)
}

// renderEnemyEncyclopedia は敵図鑑をレンダリングします。
func (s *EncyclopediaScreen) renderEnemyEncyclopedia() string {
	if len(s.data.AllEnemyTypes) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("敵タイプがありません")
	}

	listContent := s.renderEnemyList()
	previewContent := s.renderEnemyPreview()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(40).
		Render(listContent)

	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(50).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listBox, "  ", previewBox)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderEnemyList は敵図鑑のリストをレンダリングします。
func (s *EncyclopediaScreen) renderEnemyList() string {
	var items []string
	for i, et := range s.data.AllEnemyTypes {
		encountered := s.isEnemyEncountered(et.ID)
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if !encountered {
			style = style.Foreground(styles.ColorSubtle)
		}

		displayName := s.getEnemyDisplayName(et)
		status := ""
		if encountered {
			status = " [遭遇済み]"
		}
		items = append(items, style.Render(prefix+displayName+status))
	}
	return strings.Join(items, "\n")
}

// renderEnemyPreview は敵図鑑のプレビューをレンダリングします。
func (s *EncyclopediaScreen) renderEnemyPreview() string {
	if s.selectedIndex >= len(s.data.AllEnemyTypes) {
		return "選択してください"
	}

	et := s.data.AllEnemyTypes[s.selectedIndex]
	encountered := s.isEnemyEncountered(et.ID)

	if !encountered {
		return "この敵はまだ遭遇していません\n\nバトルで遭遇しましょう"
	}

	panel := components.NewInfoPanel(et.Name)
	panel.AddItem("ID", et.ID)
	panel.AddItem("基礎HP", fmt.Sprintf("%d", et.BaseHP))
	panel.AddItem("基礎攻撃力", fmt.Sprintf("%d", et.BaseAttackPower))

	return panel.Render(45)
}

// renderCompletionRate はコンプリート率をレンダリングします。
func (s *EncyclopediaScreen) renderCompletionRate() string {
	coreRate := s.getCoreCompletionRate()
	moduleRate := s.getModuleCompletionRate()
	enemyRate := s.getEnemyCompletionRate()

	rateStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSecondary)

	content := fmt.Sprintf(
		"コンプリート率  コア: %d%%  モジュール: %d%%  敵: %d%%",
		coreRate, moduleRate, enemyRate,
	)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(rateStyle.Render(content))
}

// ==================== Screenインターフェース実装 ====================

// SetSize は画面サイズを設定します。
// Screenインターフェースの実装です。
func (s *EncyclopediaScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// GetTitle は画面のタイトルを返します。
// Screenインターフェースの実装です。
func (s *EncyclopediaScreen) GetTitle() string {
	return "図鑑"
}

// GetSize は現在の画面サイズを返します。
func (s *EncyclopediaScreen) GetSize() (width, height int) {
	return s.width, s.height
}
