// Package screens はTUIゲームの画面を提供します。
// このファイルはAgentCustomizationScreen（エージェントカスタマイズ画面）を定義します。

package screens

import (
	"fmt"
	"sort"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== モード定義 ====================

// CustomizationMode はカスタマイズ画面のモードを表します。
type CustomizationMode int

const (
	// CustomizationModeSlotSelect はスロット選択モードです。
	CustomizationModeSlotSelect CustomizationMode = iota
	// CustomizationModeCoreSelect はコア選択モードです。
	CustomizationModeCoreSelect
	// CustomizationModeLevelSelect はレベル選択モードです。
	CustomizationModeLevelSelect
	// CustomizationModeSkillSlotSelect はスキルスロット選択モードです。
	CustomizationModeSkillSlotSelect
	// CustomizationModeSkillSelect はスキル選択モードです。
	CustomizationModeSkillSelect
	// CustomizationModeChainSelect はチェイン効果選択モードです。
	CustomizationModeChainSelect
)

// ==================== 表示用アイテム定義 ====================

// CoreSelectItem はコア選択リストの表示アイテムです。
type CoreSelectItem struct {
	TypeID   string
	TypeName string
	MaxLevel int
}

// SkillSelectItem はスキル選択リストの表示アイテムです。
type SkillSelectItem struct {
	TypeID          string
	TypeName        string
	Icon            string
	IsCompatible    bool
	ChainCount      int
	ChainVariations []string
}

// ==================== AgentCustomizationScreen ====================

// AgentCustomizationScreen はエージェントカスタマイズ画面を表します。
// 3つのエージェントスロットに対してコアとスキルを自由に付け替えできます。
type AgentCustomizationScreen struct {
	// インベントリマネージャーへの参照
	invManager *inventory.InventoryManager

	// エージェントスロットマネージャーへの参照
	slotManager *slot.AgentSlotManager

	// マスタデータ
	coreTypes  map[string]domain.CoreType
	skillTypes map[string]domain.SkillType

	// 現在のモード
	currentMode CustomizationMode

	// スロット選択
	selectedSlotIndex int

	// スキルスロット選択
	selectedSkillSlotIndex int

	// コア選択
	coreList           []CoreSelectItem
	selectedCoreIndex  int
	selectedCoreTypeID string

	// レベル選択
	selectedLevelIndex int
	maxSelectableLevel int

	// スキル選択
	compatibleSkillList []SkillSelectItem
	selectedSkillIndex  int
	selectedSkillTypeID string

	// チェイン効果選択
	chainVariationList []string
	selectedChainIndex int

	// スタイル
	styles *styles.GameStyles
	width  int
	height int

	// ステータス/エラーメッセージ
	statusMessage string
	errorMessage  string
}

// NewAgentCustomizationScreen は新しいAgentCustomizationScreenを作成します。
func NewAgentCustomizationScreen(
	invManager *inventory.InventoryManager,
	slotManager *slot.AgentSlotManager,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
) *AgentCustomizationScreen {
	screen := &AgentCustomizationScreen{
		invManager:  invManager,
		slotManager: slotManager,
		coreTypes:   coreTypes,
		skillTypes:  skillTypes,
		currentMode: CustomizationModeSlotSelect,
		styles:      styles.NewGameStyles(),
		width:       140,
		height:      40,
	}

	return screen
}

// Init は画面の初期化を行います。
func (s *AgentCustomizationScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *AgentCustomizationScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *AgentCustomizationScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch s.currentMode {
	case CustomizationModeSlotSelect:
		return s.handleSlotSelectKey(msg)
	case CustomizationModeCoreSelect:
		return s.handleCoreSelectKey(msg)
	case CustomizationModeLevelSelect:
		return s.handleLevelSelectKey(msg)
	case CustomizationModeSkillSlotSelect:
		return s.handleSkillSlotSelectKey(msg)
	case CustomizationModeSkillSelect:
		return s.handleSkillSelectKey(msg)
	case CustomizationModeChainSelect:
		return s.handleChainSelectKey(msg)
	}

	return s, nil
}

// ==================== スロット選択モード ====================

func (s *AgentCustomizationScreen) handleSlotSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "up", "k":
		if s.selectedSlotIndex > 0 {
			s.selectedSlotIndex--
		}
	case "down", "j":
		if s.selectedSlotIndex < slot.MaxAgentSlotCount-1 {
			s.selectedSlotIndex++
		}
	case "enter":
		s.enterCoreSelectMode()
	case "s":
		// スキル設定モードへ（コアが設定されている場合のみ）
		if s.slotManager.IsSlotReady(s.selectedSlotIndex) {
			s.enterSkillSlotSelectMode()
		} else {
			s.errorMessage = "先にコアを設定してください"
		}
	case "delete", "d":
		// スロットのコアをクリア
		s.clearCurrentSlotCore()
	}

	return s, nil
}

// enterCoreSelectMode はコア選択モードに遷移します。
func (s *AgentCustomizationScreen) enterCoreSelectMode() {
	s.currentMode = CustomizationModeCoreSelect
	s.selectedCoreIndex = 0
	s.updateCoreList()
	s.errorMessage = ""
	s.statusMessage = ""
}

// updateCoreList はコア選択リストを更新します。
func (s *AgentCustomizationScreen) updateCoreList() {
	s.coreList = []CoreSelectItem{}

	if s.invManager == nil {
		return
	}

	// 保有コアを取得
	ownedCores := s.invManager.Cores().GetOwnedCores()

	// TypeIDでソート
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

		item := CoreSelectItem{
			TypeID:   typeID,
			TypeName: typeName,
			MaxLevel: maxLevel,
		}
		s.coreList = append(s.coreList, item)
	}
}

// clearCurrentSlotCore は現在選択中のスロットのコアをクリアします。
func (s *AgentCustomizationScreen) clearCurrentSlotCore() {
	if err := s.slotManager.ClearCore(s.selectedSlotIndex); err != nil {
		s.errorMessage = fmt.Sprintf("クリアに失敗: %v", err)
	} else {
		s.statusMessage = fmt.Sprintf("スロット%dをクリアしました", s.selectedSlotIndex+1)
		s.errorMessage = ""
	}
}

// ==================== コア選択モード ====================

func (s *AgentCustomizationScreen) handleCoreSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = CustomizationModeSlotSelect
	case "up", "k":
		if s.selectedCoreIndex > 0 {
			s.selectedCoreIndex--
		}
	case "down", "j":
		if s.selectedCoreIndex < len(s.coreList)-1 {
			s.selectedCoreIndex++
		}
	case "enter":
		if s.selectedCoreIndex < len(s.coreList) {
			s.selectedCoreTypeID = s.coreList[s.selectedCoreIndex].TypeID
			s.maxSelectableLevel = s.coreList[s.selectedCoreIndex].MaxLevel
			s.selectedLevelIndex = s.maxSelectableLevel - 1 // デフォルトで最大レベルを選択
			s.currentMode = CustomizationModeLevelSelect
		}
	}

	return s, nil
}

// ==================== レベル選択モード ====================

func (s *AgentCustomizationScreen) handleLevelSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = CustomizationModeCoreSelect
	case "up", "k":
		if s.selectedLevelIndex < s.maxSelectableLevel-1 {
			s.selectedLevelIndex++
		}
	case "down", "j":
		if s.selectedLevelIndex > 0 {
			s.selectedLevelIndex--
		}
	case "enter":
		// コアをスロットに設定
		level := s.selectedLevelIndex + 1
		if err := s.slotManager.SetCore(s.selectedSlotIndex, s.selectedCoreTypeID, level); err != nil {
			s.errorMessage = fmt.Sprintf("コア設定に失敗: %v", err)
		} else {
			s.statusMessage = fmt.Sprintf("スロット%dにコアを設定しました", s.selectedSlotIndex+1)
			s.errorMessage = ""
			s.currentMode = CustomizationModeSlotSelect
		}
	}

	return s, nil
}

// ==================== スキルスロット選択モード ====================

func (s *AgentCustomizationScreen) enterSkillSlotSelectMode() {
	s.currentMode = CustomizationModeSkillSlotSelect
	s.selectedSkillSlotIndex = 0
	s.errorMessage = ""
	s.statusMessage = ""
}

func (s *AgentCustomizationScreen) handleSkillSlotSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = CustomizationModeSlotSelect
	case "up", "k":
		if s.selectedSkillSlotIndex > 0 {
			s.selectedSkillSlotIndex--
		}
	case "down", "j":
		if s.selectedSkillSlotIndex < domain.MaxSkillSlotCount-1 {
			s.selectedSkillSlotIndex++
		}
	case "enter":
		s.enterSkillSelectMode()
	case "backspace", "delete", "d":
		// スキルスロットをクリア
		s.clearCurrentSkillSlot()
	}

	return s, nil
}

// clearCurrentSkillSlot は現在選択中のスキルスロットをクリアします。
func (s *AgentCustomizationScreen) clearCurrentSkillSlot() {
	if err := s.slotManager.ClearSkill(s.selectedSlotIndex, s.selectedSkillSlotIndex); err != nil {
		s.errorMessage = fmt.Sprintf("スキルクリアに失敗: %v", err)
	} else {
		s.statusMessage = fmt.Sprintf("スキルスロット%dをクリアしました", s.selectedSkillSlotIndex+1)
		s.errorMessage = ""
	}
}

// ==================== スキル選択モード ====================

func (s *AgentCustomizationScreen) enterSkillSelectMode() {
	s.currentMode = CustomizationModeSkillSelect
	s.selectedSkillIndex = 0
	s.updateSkillList()
	s.errorMessage = ""
	s.statusMessage = ""
}

// updateSkillList は互換スキルリストを更新します。
func (s *AgentCustomizationScreen) updateSkillList() {
	s.compatibleSkillList = []SkillSelectItem{}

	if s.invManager == nil || s.slotManager == nil {
		return
	}

	// 互換スキルのTypeIDリストを取得
	compatibleSkillIDs := s.slotManager.GetCompatibleSkills(s.selectedSlotIndex)

	// 保有スキル情報を取得
	ownedSkills := s.invManager.Skills().GetOwnedSkills()

	// TypeIDでソート
	sort.Strings(compatibleSkillIDs)

	// リストを構築
	for _, typeID := range compatibleSkillIDs {
		ownership, exists := ownedSkills[typeID]
		if !exists {
			continue
		}

		// スキルタイプ情報を取得
		typeName := typeID
		icon := "?"
		if skillType, ok := s.skillTypes[typeID]; ok {
			typeName = skillType.Name
			icon = skillType.Icon
		}

		// チェイン効果バリエーションを取得
		chainVariations := ownership.GetChainVariations()

		item := SkillSelectItem{
			TypeID:          typeID,
			TypeName:        typeName,
			Icon:            icon,
			IsCompatible:    true,
			ChainCount:      len(chainVariations),
			ChainVariations: chainVariations,
		}
		s.compatibleSkillList = append(s.compatibleSkillList, item)
	}
}

func (s *AgentCustomizationScreen) handleSkillSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = CustomizationModeSkillSlotSelect
	case "up", "k":
		if s.selectedSkillIndex > 0 {
			s.selectedSkillIndex--
		}
	case "down", "j":
		if s.selectedSkillIndex < len(s.compatibleSkillList)-1 {
			s.selectedSkillIndex++
		}
	case "enter":
		if s.selectedSkillIndex < len(s.compatibleSkillList) {
			skill := s.compatibleSkillList[s.selectedSkillIndex]
			s.selectedSkillTypeID = skill.TypeID

			// チェイン効果バリエーションがある場合は選択モードへ
			if len(skill.ChainVariations) > 0 {
				s.chainVariationList = skill.ChainVariations
				s.selectedChainIndex = 0
				s.currentMode = CustomizationModeChainSelect
			} else {
				// チェイン効果なしで直接設定
				s.setSkillToSlot("")
			}
		}
	}

	return s, nil
}

// ==================== チェイン効果選択モード ====================

func (s *AgentCustomizationScreen) handleChainSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = CustomizationModeSkillSelect
	case "up", "k":
		if s.selectedChainIndex > 0 {
			s.selectedChainIndex--
		}
	case "down", "j":
		// +1 は「なし」オプション用
		if s.selectedChainIndex < len(s.chainVariationList) {
			s.selectedChainIndex++
		}
	case "enter":
		var chainEffectID string
		if s.selectedChainIndex < len(s.chainVariationList) {
			chainEffectID = s.chainVariationList[s.selectedChainIndex]
		} else {
			chainEffectID = "" // 「なし」を選択
		}
		s.setSkillToSlot(chainEffectID)
	}

	return s, nil
}

// setSkillToSlot はスキルをスロットに設定します。
func (s *AgentCustomizationScreen) setSkillToSlot(chainEffectID string) {
	if err := s.slotManager.SetSkill(s.selectedSlotIndex, s.selectedSkillSlotIndex, s.selectedSkillTypeID, chainEffectID); err != nil {
		s.errorMessage = fmt.Sprintf("スキル設定に失敗: %v", err)
	} else {
		s.statusMessage = "スキルを設定しました"
		s.errorMessage = ""
		s.currentMode = CustomizationModeSkillSlotSelect
	}
}

// ==================== View ====================

// View は画面をレンダリングします。
func (s *AgentCustomizationScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("エージェントカスタマイズ"))
	builder.WriteString("\n\n")

	// メインコンテンツ
	builder.WriteString(s.renderMainContent())
	builder.WriteString("\n\n")

	// ステータス/エラーメッセージ
	if s.errorMessage != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")). // 赤色
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(errorStyle.Render(s.errorMessage))
		builder.WriteString("\n")
	} else if s.statusMessage != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSecondary).
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(statusStyle.Render(s.statusMessage))
		builder.WriteString("\n")
	}

	// ヒント
	builder.WriteString(s.renderHints())

	return builder.String()
}

// renderMainContent はメインコンテンツをレンダリングします。
func (s *AgentCustomizationScreen) renderMainContent() string {
	// 左側: スロット一覧
	leftContent := s.renderSlotList()

	// 右側: 選択パネル（モードによって変化）
	rightContent := s.renderSelectionPanel()

	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(50).
		Render(leftContent)

	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(55).
		Render(rightContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, "  ", rightBox)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderSlotList はスロット一覧をレンダリングします。
func (s *AgentCustomizationScreen) renderSlotList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("エージェントスロット"))
	builder.WriteString("\n\n")

	slots := s.slotManager.GetSlots()

	for i := 0; i < slot.MaxAgentSlotCount; i++ {
		agentSlot := slots[i]

		style := lipgloss.NewStyle()
		prefix := "  "

		// スロット選択モードで選択中の場合
		if s.currentMode == CustomizationModeSlotSelect && i == s.selectedSlotIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		// スロット内容を表示
		var slotContent string
		if agentSlot == nil || agentSlot.IsEmpty() {
			slotContent = fmt.Sprintf("スロット %d: (空)", i+1)
		} else {
			// コア名を取得
			coreName := agentSlot.CoreTypeID
			if coreType, ok := s.coreTypes[agentSlot.CoreTypeID]; ok {
				coreName = coreType.Name
			}
			skillCount := agentSlot.GetSkillCount()
			slotContent = fmt.Sprintf("スロット %d: %s Lv.%d (%d/%dスキル)",
				i+1, coreName, agentSlot.CoreLevel, skillCount, domain.MaxSkillSlotCount)
		}

		builder.WriteString(style.Render(prefix + slotContent))
		builder.WriteString("\n")

		// スキルスロット選択モードでこのスロットが選択されている場合、スキル詳細を表示
		if s.currentMode != CustomizationModeSlotSelect && i == s.selectedSlotIndex && agentSlot != nil && !agentSlot.IsEmpty() {
			builder.WriteString(s.renderSkillSlots(agentSlot))
		}
	}

	return builder.String()
}

// renderSkillSlots はスキルスロットをレンダリングします。
func (s *AgentCustomizationScreen) renderSkillSlots(agentSlot *domain.AgentSlot) string {
	var builder strings.Builder

	for j := 0; j < domain.MaxSkillSlotCount; j++ {
		skillConfig := agentSlot.GetSkill(j)

		style := lipgloss.NewStyle().Padding(0, 0, 0, 4)
		prefix := "  "

		// スキルスロット選択モードで選択中の場合
		if s.currentMode == CustomizationModeSkillSlotSelect && j == s.selectedSkillSlotIndex {
			style = style.Bold(true).Foreground(styles.ColorPrimary)
			prefix = "> "
		}

		var skillContent string
		if skillConfig == nil || skillConfig.IsEmpty() {
			skillContent = fmt.Sprintf("スキル%d: (空)", j+1)
		} else {
			skillName := skillConfig.TypeID
			icon := "?"
			if skillType, ok := s.skillTypes[skillConfig.TypeID]; ok {
				skillName = skillType.Name
				icon = skillType.Icon
			}
			skillContent = fmt.Sprintf("スキル%d: %s %s", j+1, icon, skillName)
			if skillConfig.ChainEffectID != "" {
				skillContent += fmt.Sprintf(" [%s]", skillConfig.ChainEffectID)
			}
		}

		builder.WriteString(style.Render(prefix + skillContent))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderSelectionPanel は右側の選択パネルをレンダリングします。
func (s *AgentCustomizationScreen) renderSelectionPanel() string {
	switch s.currentMode {
	case CustomizationModeSlotSelect:
		return s.renderSlotDetail()
	case CustomizationModeCoreSelect:
		return s.renderCoreSelectList()
	case CustomizationModeLevelSelect:
		return s.renderLevelSelectList()
	case CustomizationModeSkillSlotSelect:
		return s.renderSkillSlotDetail()
	case CustomizationModeSkillSelect:
		return s.renderSkillSelectList()
	case CustomizationModeChainSelect:
		return s.renderChainSelectList()
	}
	return ""
}

// renderSlotDetail はスロットの詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderSlotDetail() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("スロット詳細"))
	builder.WriteString("\n\n")

	slot := s.slotManager.GetSlot(s.selectedSlotIndex)
	if slot == nil || slot.IsEmpty() {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("スロットが空です"))
		builder.WriteString("\n\n")
		builder.WriteString("Enterでコアを選択してください")
		return builder.String()
	}

	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary)

	// コア情報
	coreName := slot.CoreTypeID
	if coreType, ok := s.coreTypes[slot.CoreTypeID]; ok {
		coreName = coreType.Name
	}

	builder.WriteString(labelStyle.Render("コア: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%s Lv.%d", coreName, slot.CoreLevel)))
	builder.WriteString("\n\n")

	// スキル情報
	builder.WriteString(labelStyle.Render("スキル:"))
	builder.WriteString("\n")
	for j := 0; j < domain.MaxSkillSlotCount; j++ {
		skillConfig := slot.GetSkill(j)
		if skillConfig == nil || skillConfig.IsEmpty() {
			builder.WriteString(fmt.Sprintf("  %d: (空)\n", j+1))
		} else {
			skillName := skillConfig.TypeID
			icon := "?"
			if skillType, ok := s.skillTypes[skillConfig.TypeID]; ok {
				skillName = skillType.Name
				icon = skillType.Icon
			}
			builder.WriteString(fmt.Sprintf("  %d: %s %s\n", j+1, icon, skillName))
		}
	}

	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("Sキーでスキルを設定"))

	return builder.String()
}

// renderCoreSelectList はコア選択リストをレンダリングします。
func (s *AgentCustomizationScreen) renderCoreSelectList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("コア選択"))
	builder.WriteString("\n\n")

	if len(s.coreList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("保有コアがありません"))
		return builder.String()
	}

	for i, core := range s.coreList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedCoreIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		item := fmt.Sprintf("%s (最大Lv.%d)", core.TypeName, core.MaxLevel)
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderLevelSelectList はレベル選択リストをレンダリングします。
func (s *AgentCustomizationScreen) renderLevelSelectList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("レベル選択"))
	builder.WriteString("\n\n")

	// コア名を表示
	coreName := s.selectedCoreTypeID
	if coreType, ok := s.coreTypes[s.selectedCoreTypeID]; ok {
		coreName = coreType.Name
	}
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	builder.WriteString(labelStyle.Render(fmt.Sprintf("コア: %s", coreName)))
	builder.WriteString("\n\n")

	// レベルリスト（高い順に表示）
	for i := s.maxSelectableLevel; i >= 1; i-- {
		idx := s.maxSelectableLevel - i
		style := lipgloss.NewStyle()
		prefix := "  "

		if idx == s.selectedLevelIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		item := fmt.Sprintf("Lv.%d", s.maxSelectableLevel-idx)
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderSkillSlotDetail はスキルスロットの詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderSkillSlotDetail() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("スキルスロット選択"))
	builder.WriteString("\n\n")

	slot := s.slotManager.GetSlot(s.selectedSlotIndex)
	if slot == nil {
		return "スロットが見つかりません"
	}

	builder.WriteString("Enterでスキルを選択、Backspaceでクリア")

	return builder.String()
}

// renderSkillSelectList はスキル選択リストをレンダリングします。
func (s *AgentCustomizationScreen) renderSkillSelectList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("スキル選択"))
	builder.WriteString("\n\n")

	if len(s.compatibleSkillList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("互換スキルがありません"))
		return builder.String()
	}

	for i, skill := range s.compatibleSkillList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedSkillIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		item := fmt.Sprintf("%s %s", skill.Icon, skill.TypeName)
		if skill.ChainCount > 0 {
			item += fmt.Sprintf(" (%d種のチェイン効果)", skill.ChainCount)
		}
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderChainSelectList はチェイン効果選択リストをレンダリングします。
func (s *AgentCustomizationScreen) renderChainSelectList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("チェイン効果選択"))
	builder.WriteString("\n\n")

	// スキル名を表示
	skillName := s.selectedSkillTypeID
	if skillType, ok := s.skillTypes[s.selectedSkillTypeID]; ok {
		skillName = skillType.Name
	}
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	builder.WriteString(labelStyle.Render(fmt.Sprintf("スキル: %s", skillName)))
	builder.WriteString("\n\n")

	// チェイン効果リスト
	for i, chainID := range s.chainVariationList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedChainIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		builder.WriteString(style.Render(prefix + chainID))
		builder.WriteString("\n")
	}

	// 「なし」オプション
	style := lipgloss.NewStyle()
	prefix := "  "
	if s.selectedChainIndex == len(s.chainVariationList) {
		style = style.Bold(true).
			Foreground(styles.ColorSelectedFg).
			Background(styles.ColorSelectedBg)
		prefix = "> "
	}
	builder.WriteString(style.Render(prefix + "(なし)"))
	builder.WriteString("\n")

	return builder.String()
}

// renderHints はヒントをレンダリングします。
func (s *AgentCustomizationScreen) renderHints() string {
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hints string
	switch s.currentMode {
	case CustomizationModeSlotSelect:
		hints = "↑/↓: スロット選択  Enter: コア設定  S: スキル設定  D: クリア  Esc: ホーム"
	case CustomizationModeCoreSelect:
		hints = "↑/↓: コア選択  Enter: 決定  Esc: 戻る"
	case CustomizationModeLevelSelect:
		hints = "↑/↓: レベル選択  Enter: 決定  Esc: 戻る"
	case CustomizationModeSkillSlotSelect:
		hints = "↑/↓: スキルスロット選択  Enter: スキル設定  Backspace: クリア  Esc: 戻る"
	case CustomizationModeSkillSelect:
		hints = "↑/↓: スキル選択  Enter: 決定  Esc: 戻る"
	case CustomizationModeChainSelect:
		hints = "↑/↓: チェイン効果選択  Enter: 決定  Esc: 戻る"
	}

	return hintStyle.Render(hints)
}

// ==================== Screenインターフェース実装 ====================

// SetSize は画面サイズを設定します。
func (s *AgentCustomizationScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// GetTitle は画面のタイトルを返します。
func (s *AgentCustomizationScreen) GetTitle() string {
	return "エージェントカスタマイズ"
}

// GetSize は現在の画面サイズを返します。
func (s *AgentCustomizationScreen) GetSize() (width, height int) {
	return s.width, s.height
}

// RefreshData はデータを最新状態に更新します。
func (s *AgentCustomizationScreen) RefreshData() {
	s.updateCoreList()
	s.updateSkillList()
}
