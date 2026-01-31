// Package screens はTUIゲームの画面を提供します。
// このファイルはInventoryScreen（インベントリ画面）を定義します。

package screens

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rivo/uniseg"
)

// InventoryTab はインベントリ画面のタブを表します。
type InventoryTab int

const (
	// TabCoreInventory はコア一覧タブです。
	TabCoreInventory InventoryTab = iota
	// TabSkillInventory はスキル一覧タブです。
	TabSkillInventory
)

// UniqueInventoryProvider はユニークインベントリデータを提供するインターフェースです。
type UniqueInventoryProvider interface {
	// GetOwnedCoreTypes は保有コアTypeID一覧を返します。
	GetOwnedCoreTypes() []string
	// GetOwnedSkillTypes は保有スキルTypeID一覧を返します。
	GetOwnedSkillTypes() []string
	// Cores はコアインベントリを返します。
	Cores() *domain.CoreInventory
	// Skills はスキルインベントリを返します。
	Skills() *domain.SkillInventory
}

// InventoryScreen はインベントリ画面を表します。
// コアとスキルの一覧表示、装備状況の確認ができます。
type InventoryScreen struct {
	// インベントリマネージャーへの参照
	invManager *inventory.InventoryManager

	// エージェントスロットマネージャーへの参照（装備状況表示用）
	slotManager *slot.AgentSlotManager

	// マスタデータ
	coreTypes  map[string]domain.CoreType
	skillTypes map[string]domain.SkillType

	// 現在のタブ
	currentTab InventoryTab

	// コア一覧関連
	coreList          []CoreInventoryItem
	selectedCoreIndex int

	// スキル一覧関連
	skillList          []SkillInventoryItem
	selectedSkillIndex int

	// スキル詳細表示モード（チェイン効果バリエーション表示）
	showingSkillDetail bool

	// スタイル
	styles *styles.GameStyles
	width  int
	height int
}

// CoreInventoryItem はコア一覧の表示アイテムです。
type CoreInventoryItem struct {
	TypeID        string
	TypeName      string
	MaxLevel      int
	IsEquipped    bool  // いずれかのスロットに装備中か
	EquippedSlots []int // 装備中のスロット番号リスト
}

// SkillInventoryItem はスキル一覧の表示アイテムです。
type SkillInventoryItem struct {
	TypeID          string
	TypeName        string
	Icon            string
	ChainCount      int      // 利用可能なチェイン効果バリエーション数
	ChainVariations []string // チェイン効果IDリスト
	IsEquipped      bool     // いずれかのスロットに装備中か
	EquippedSlots   []int    // 装備中のスロット番号リスト
}

// NewInventoryScreen は新しいInventoryScreenを作成します。
func NewInventoryScreen(
	invManager *inventory.InventoryManager,
	slotManager *slot.AgentSlotManager,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
) *InventoryScreen {
	screen := &InventoryScreen{
		invManager:  invManager,
		slotManager: slotManager,
		coreTypes:   coreTypes,
		skillTypes:  skillTypes,
		currentTab:  TabCoreInventory,
		styles:      styles.NewGameStyles(),
		width:       140,
		height:      40,
	}

	// リストを初期化
	screen.updateCoreList()
	screen.updateSkillList()

	return screen
}

// Init は画面の初期化を行います。
func (s *InventoryScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *InventoryScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *InventoryScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// スキル詳細表示中の場合
	if s.showingSkillDetail {
		switch msg.String() {
		case "esc", "backspace":
			s.showingSkillDetail = false
		}
		return s, nil
	}

	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "left", "h":
		s.prevTab()
	case "right", "l":
		s.nextTab()
	case "up", "k":
		s.moveUp()
	case "down", "j":
		s.moveDown()
	case "enter":
		// スキルタブでEnterを押すと詳細表示
		if s.currentTab == TabSkillInventory && len(s.skillList) > 0 {
			s.showingSkillDetail = true
		}
	}

	return s, nil
}

// prevTab は前のタブに移動します。
func (s *InventoryScreen) prevTab() {
	if s.currentTab > TabCoreInventory {
		s.currentTab--
		s.showingSkillDetail = false
	}
}

// nextTab は次のタブに移動します。
func (s *InventoryScreen) nextTab() {
	if s.currentTab < TabSkillInventory {
		s.currentTab++
		s.showingSkillDetail = false
	}
}

// moveUp は選択を上に移動します。
func (s *InventoryScreen) moveUp() {
	switch s.currentTab {
	case TabCoreInventory:
		if s.selectedCoreIndex > 0 {
			s.selectedCoreIndex--
		}
	case TabSkillInventory:
		if s.selectedSkillIndex > 0 {
			s.selectedSkillIndex--
		}
	}
}

// moveDown は選択を下に移動します。
func (s *InventoryScreen) moveDown() {
	switch s.currentTab {
	case TabCoreInventory:
		if s.selectedCoreIndex < len(s.coreList)-1 {
			s.selectedCoreIndex++
		}
	case TabSkillInventory:
		if s.selectedSkillIndex < len(s.skillList)-1 {
			s.selectedSkillIndex++
		}
	}
}

// updateCoreList はコア一覧を更新します。
func (s *InventoryScreen) updateCoreList() {
	s.coreList = []CoreInventoryItem{}

	if s.invManager == nil {
		return
	}

	// 保有コアを取得
	ownedCores := s.invManager.Cores().GetOwnedCores()

	// 装備状況を取得
	equippedCores := s.getEquippedCoreTypeIDs()

	// TypeIDでソートするためにキーを取得
	typeIDs := make([]string, 0, len(ownedCores))
	for typeID := range ownedCores {
		typeIDs = append(typeIDs, typeID)
	}
	sort.Strings(typeIDs)

	// リストを構築
	for _, typeID := range typeIDs {
		maxLevel := ownedCores[typeID]

		// コアタイプ名を取得
		typeName := typeID
		if coreType, ok := s.coreTypes[typeID]; ok {
			typeName = coreType.Name
		}

		// 装備状況を確認
		equippedSlots := equippedCores[typeID]

		item := CoreInventoryItem{
			TypeID:        typeID,
			TypeName:      typeName,
			MaxLevel:      maxLevel,
			IsEquipped:    len(equippedSlots) > 0,
			EquippedSlots: equippedSlots,
		}
		s.coreList = append(s.coreList, item)
	}
}

// updateSkillList はスキル一覧を更新します。
func (s *InventoryScreen) updateSkillList() {
	s.skillList = []SkillInventoryItem{}

	if s.invManager == nil {
		return
	}

	// 保有スキルを取得
	ownedSkills := s.invManager.Skills().GetOwnedSkills()

	// 装備状況を取得
	equippedSkills := s.getEquippedSkillTypeIDs()

	// TypeIDでソートするためにキーを取得
	typeIDs := make([]string, 0, len(ownedSkills))
	for typeID := range ownedSkills {
		typeIDs = append(typeIDs, typeID)
	}
	sort.Strings(typeIDs)

	// リストを構築
	for _, typeID := range typeIDs {
		ownership := ownedSkills[typeID]

		// スキルタイプ情報を取得
		typeName := typeID
		icon := "?"
		if skillType, ok := s.skillTypes[typeID]; ok {
			typeName = skillType.Name
			icon = skillType.Icon
		}

		// チェイン効果バリエーションを取得
		chainVariations := ownership.GetChainVariations()

		// 装備状況を確認
		equippedSlots := equippedSkills[typeID]

		item := SkillInventoryItem{
			TypeID:          typeID,
			TypeName:        typeName,
			Icon:            icon,
			ChainCount:      len(chainVariations),
			ChainVariations: chainVariations,
			IsEquipped:      len(equippedSlots) > 0,
			EquippedSlots:   equippedSlots,
		}
		s.skillList = append(s.skillList, item)
	}
}

// getEquippedCoreTypeIDs は装備中のコアTypeIDとスロット番号のマップを返します。
func (s *InventoryScreen) getEquippedCoreTypeIDs() map[string][]int {
	result := make(map[string][]int)

	if s.slotManager == nil {
		return result
	}

	slots := s.slotManager.GetSlots()
	for i, agentSlot := range slots {
		if agentSlot != nil && !agentSlot.IsEmpty() {
			typeID := agentSlot.CoreTypeID
			result[typeID] = append(result[typeID], i)
		}
	}

	return result
}

// getEquippedSkillTypeIDs は装備中のスキルTypeIDとスロット番号のマップを返します。
func (s *InventoryScreen) getEquippedSkillTypeIDs() map[string][]int {
	result := make(map[string][]int)

	if s.slotManager == nil {
		return result
	}

	slots := s.slotManager.GetSlots()
	for i, agentSlot := range slots {
		if agentSlot == nil || agentSlot.IsEmpty() {
			continue
		}

		// 各スキルスロットを確認
		for j := range domain.MaxSkillSlotCount {
			skillConfig := agentSlot.GetSkill(j)
			if skillConfig != nil && !skillConfig.IsEmpty() {
				typeID := skillConfig.TypeID
				// 同じスロットで重複しないようにチェック
				if !slices.Contains(result[typeID], i) {
					result[typeID] = append(result[typeID], i)
				}
			}
		}
	}

	return result
}

// View は画面をレンダリングします。
func (s *InventoryScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("インベントリ"))
	builder.WriteString("\n\n")

	// タブバー
	builder.WriteString(s.renderTabBar())
	builder.WriteString("\n\n")

	// メインコンテンツ
	builder.WriteString(s.renderMainContent())
	builder.WriteString("\n\n")

	// ヒント
	builder.WriteString(s.renderHints())

	return builder.String()
}

// renderTabBar はタブバーをレンダリングします。
func (s *InventoryScreen) renderTabBar() string {
	tabs := []string{"コア一覧", "スキル一覧"}

	var tabItems []string
	for i, tab := range tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if InventoryTab(i) == s.currentTab {
			style = style.
				Bold(true).
				Foreground(styles.ColorPrimary).
				Background(lipgloss.Color("236"))
		} else {
			style = style.Foreground(styles.ColorSubtle)
		}
		tabItems = append(tabItems, style.Render(tab))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Center, tabItems...)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(tabBar)
}

// renderMainContent はメインコンテンツをレンダリングします。
func (s *InventoryScreen) renderMainContent() string {
	switch s.currentTab {
	case TabCoreInventory:
		return s.renderCoreInventory()
	case TabSkillInventory:
		return s.renderSkillInventory()
	}
	return ""
}

// renderCoreInventory はコア一覧をレンダリングします。
func (s *InventoryScreen) renderCoreInventory() string {
	if len(s.coreList) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("コアがありません")
	}

	// リストとプレビューを横に並べる
	listContent := s.renderCoreListItems()
	previewContent := s.renderCorePreviewContent()

	// Width()を使わない - コンテンツが既に固定幅にパディングされているため
	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
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

// renderCoreListItems はコアリストの項目をレンダリングします。
func (s *InventoryScreen) renderCoreListItems() string {
	var items []string

	// 内部コンテンツの固定表示幅
	// lipglossのWidth()に依存せず、全行を同じ幅にパディングしてボーダーを正しく描画
	const contentWidth = 50

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	titleText := padToDisplayWidth("保有コア一覧", contentWidth)
	items = append(items, titleStyle.Render(titleText))
	items = append(items, padToDisplayWidth("", contentWidth))

	for i, core := range s.coreList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedCoreIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		// 装備マーカー
		equipMark := ""
		if core.IsEquipped {
			slots := make([]string, len(core.EquippedSlots))
			for j, slotNum := range core.EquippedSlots {
				slots[j] = fmt.Sprintf("%d", slotNum+1)
			}
			equipMark = fmt.Sprintf(" [E%s]", strings.Join(slots, ","))
		}

		item := fmt.Sprintf("%s Lv.%d%s", core.TypeName, core.MaxLevel, equipMark)

		// プレフィックスを含めた全体の内容を作成し、表示幅で切り詰め後パディング
		fullContent := prefix + item
		truncated := truncateToDisplayWidth(fullContent, contentWidth)
		padded := padToDisplayWidth(truncated, contentWidth)
		items = append(items, style.Render(padded))
	}

	return strings.Join(items, "\n")
}

// renderCorePreviewContent はコアのプレビューをレンダリングします。
func (s *InventoryScreen) renderCorePreviewContent() string {
	if s.selectedCoreIndex >= len(s.coreList) {
		return "コアを選択してください"
	}

	core := s.coreList[s.selectedCoreIndex]

	var builder strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	builder.WriteString(titleStyle.Render(core.TypeName))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("TypeID: "))
	builder.WriteString(valueStyle.Render(core.TypeID))
	builder.WriteString("\n")

	builder.WriteString(labelStyle.Render("取得済み最大レベル: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("Lv.%d", core.MaxLevel)))
	builder.WriteString("\n\n")

	// コアタイプの詳細情報
	if coreType, ok := s.coreTypes[core.TypeID]; ok {
		builder.WriteString(labelStyle.Render("許可タグ: "))
		builder.WriteString(valueStyle.Render(strings.Join(coreType.AllowedTags, ", ")))
		builder.WriteString("\n")

		// ステータス重み
		builder.WriteString(labelStyle.Render("ステータス重み:"))
		builder.WriteString("\n")
		for stat, weight := range coreType.StatWeights {
			fmt.Fprintf(&builder, "  %s: %.1f\n", stat, weight)
		}
	}

	// 装備状況
	builder.WriteString("\n")
	if core.IsEquipped {
		slots := make([]string, len(core.EquippedSlots))
		for j, slotNum := range core.EquippedSlots {
			slots[j] = fmt.Sprintf("スロット%d", slotNum+1)
		}
		equipStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)
		builder.WriteString(equipStyle.Render("装備中: " + strings.Join(slots, ", ")))
	} else {
		builder.WriteString(labelStyle.Render("未装備"))
	}

	return builder.String()
}

// renderSkillInventory はスキル一覧をレンダリングします。
func (s *InventoryScreen) renderSkillInventory() string {
	if len(s.skillList) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("スキルがありません")
	}

	// リストとプレビューを横に並べる
	listContent := s.renderSkillListItems()
	previewContent := s.renderSkillPreviewContent()

	// Width()を使わない - コンテンツが既に固定幅にパディングされているため
	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
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

// renderSkillListItems はスキルリストの項目をレンダリングします。
func (s *InventoryScreen) renderSkillListItems() string {
	var items []string

	// 内部コンテンツの固定表示幅
	// lipglossのWidth()に依存せず、全行を同じ幅にパディングしてボーダーを正しく描画
	const contentWidth = 50

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	titleText := padToDisplayWidth("保有スキル一覧", contentWidth)
	items = append(items, titleStyle.Render(titleText))
	items = append(items, padToDisplayWidth("", contentWidth))

	for i, skill := range s.skillList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedSkillIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		// 装備マーカー
		equipMark := ""
		if skill.IsEquipped {
			slots := make([]string, len(skill.EquippedSlots))
			for j, slotNum := range skill.EquippedSlots {
				slots[j] = fmt.Sprintf("%d", slotNum+1)
			}
			equipMark = fmt.Sprintf(" [E%s]", strings.Join(slots, ","))
		}

		// チェイン効果バリエーション数
		chainInfo := ""
		if skill.ChainCount > 0 {
			chainInfo = fmt.Sprintf("(%d種)", skill.ChainCount)
		}

		// アイコン + スキル名 + チェイン情報 + 装備マーカーを構成
		icon := normalizeInventoryIcon(skill.Icon)
		var itemContent string
		if chainInfo != "" {
			itemContent = fmt.Sprintf("%s %s %s%s", icon, skill.TypeName, chainInfo, equipMark)
		} else {
			itemContent = fmt.Sprintf("%s %s%s", icon, skill.TypeName, equipMark)
		}

		// プレフィックスを含めた全体の内容を作成し、表示幅で切り詰め後パディング
		fullContent := prefix + itemContent
		truncated := truncateToDisplayWidth(fullContent, contentWidth)
		padded := padToDisplayWidth(truncated, contentWidth)
		items = append(items, style.Render(padded))
	}

	return strings.Join(items, "\n")
}

// renderSkillPreviewContent はスキルのプレビューをレンダリングします。
func (s *InventoryScreen) renderSkillPreviewContent() string {
	if s.selectedSkillIndex >= len(s.skillList) {
		return "スキルを選択してください"
	}

	skill := s.skillList[s.selectedSkillIndex]

	var builder strings.Builder
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
	chainStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)

	builder.WriteString(titleStyle.Render(normalizeInventoryIcon(skill.Icon) + " " + skill.TypeName))
	builder.WriteString("\n\n")

	builder.WriteString(labelStyle.Render("TypeID: "))
	builder.WriteString(valueStyle.Render(skill.TypeID))
	builder.WriteString("\n")

	// スキルタイプの詳細情報
	if skillType, ok := s.skillTypes[skill.TypeID]; ok {
		builder.WriteString(labelStyle.Render("タグ: "))
		builder.WriteString(valueStyle.Render(strings.Join(skillType.Tags, ", ")))
		builder.WriteString("\n")

		if skillType.Description != "" {
			builder.WriteString(labelStyle.Render("説明: "))
			builder.WriteString(valueStyle.Render(skillType.Description))
			builder.WriteString("\n")
		}
	}

	// チェイン効果バリエーション
	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("チェイン効果バリエーション:"))
	builder.WriteString("\n")

	if len(skill.ChainVariations) == 0 {
		builder.WriteString(labelStyle.Render("  (なし)"))
		builder.WriteString("\n")
	} else {
		// 詳細表示モードかどうかで表示を切り替え
		if s.showingSkillDetail {
			for _, chainID := range skill.ChainVariations {
				builder.WriteString(chainStyle.Render(fmt.Sprintf("  - %s", chainID)))
				builder.WriteString("\n")
			}
		} else {
			builder.WriteString(chainStyle.Render(fmt.Sprintf("  %d種類の効果を保有", len(skill.ChainVariations))))
			builder.WriteString("\n")
			builder.WriteString(labelStyle.Render("  (Enterで詳細表示)"))
			builder.WriteString("\n")
		}
	}

	// 装備状況
	builder.WriteString("\n")
	if skill.IsEquipped {
		slots := make([]string, len(skill.EquippedSlots))
		for j, slotNum := range skill.EquippedSlots {
			slots[j] = fmt.Sprintf("スロット%d", slotNum+1)
		}
		equipStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)
		builder.WriteString(equipStyle.Render("装備中: " + strings.Join(slots, ", ")))
	} else {
		builder.WriteString(labelStyle.Render("未装備"))
	}

	return builder.String()
}

// renderHints はヒントをレンダリングします。
func (s *InventoryScreen) renderHints() string {
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hints string
	if s.showingSkillDetail {
		hints = "Esc/Backspace: 戻る"
	} else if s.currentTab == TabSkillInventory && len(s.skillList) > 0 {
		hints = "←/→: タブ切替  ↑/↓: 選択  Enter: チェイン効果詳細  Esc: ホーム"
	} else {
		hints = "←/→: タブ切替  ↑/↓: 選択  Esc: ホーム"
	}

	return hintStyle.Render(hints)
}

// RefreshData はデータを最新状態に更新します。
func (s *InventoryScreen) RefreshData() {
	s.updateCoreList()
	s.updateSkillList()
}

// ==================== Screenインターフェース実装 ====================

// SetSize は画面サイズを設定します。
func (s *InventoryScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// GetTitle は画面のタイトルを返します。
func (s *InventoryScreen) GetTitle() string {
	return "インベントリ"
}

// GetSize は現在の画面サイズを返します。
func (s *InventoryScreen) GetSize() (width, height int) {
	return s.width, s.height
}

// ==================== 表示幅ヘルパー ====================

// normalizeInventoryIcon はボックスの幅ズレを避けるためにアイコン表現を正規化します。
// VS16を含む絵文字は端末によって1幅表示になるため、テキスト表示に寄せて幅を合わせます。
func normalizeInventoryIcon(icon string) string {
	const variationSelectorEmoji = "\ufe0f"
	const variationSelectorText = "\ufe0e"
	if !strings.Contains(icon, variationSelectorEmoji) {
		return icon
	}
	return strings.ReplaceAll(icon, variationSelectorEmoji, variationSelectorText)
}

// stringDisplayWidth はlipglossで文字列の表示幅を計算します。
// lipglossのレンダリングと一致する幅計算を行います。
func stringDisplayWidth(s string) int {
	return lipgloss.Width(s)
}

// truncateToDisplayWidth は文字列を指定した表示幅に切り詰めます。
// グラフェムクラスター単位で処理し、lipgloss.Width()で幅計算を行います。
func truncateToDisplayWidth(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	var result strings.Builder
	gr := uniseg.NewGraphemes(s)
	currentWidth := 0

	for gr.Next() {
		cluster := gr.Str()
		w := lipgloss.Width(cluster)
		if currentWidth+w > maxWidth {
			break
		}
		result.WriteString(cluster)
		currentWidth += w
	}

	return result.String()
}

// padToDisplayWidth は文字列を指定した表示幅になるようにスペースでパディングします。
func padToDisplayWidth(s string, targetWidth int) string {
	currentWidth := lipgloss.Width(s)
	if currentWidth >= targetWidth {
		return s
	}
	padding := targetWidth - currentWidth
	return s + strings.Repeat(" ", padding)
}
