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
	// ModeCardSelect はカード選択モード（メイン画面）です。
	ModeCardSelect CustomizationMode = iota
	// ModeCoreSelect はモーダル: コア選択モードです。
	ModeCoreSelect
	// ModeLevelSelect はモーダル: レベル選択モードです。
	ModeLevelSelect
	// ModeSkillSelect はモーダル: スキル選択モードです。
	ModeSkillSelect
	// ModeChainSelect はモーダル: チェイン効果選択モードです。
	ModeChainSelect
)

// 旧モード定義のエイリアス（後方互換性のため）
const (
	CustomizationModeSlotSelect      = ModeCardSelect
	CustomizationModeCoreSelect      = ModeCoreSelect
	CustomizationModeLevelSelect     = ModeLevelSelect
	CustomizationModeSkillSlotSelect = ModeCardSelect // スキルスロット選択はカード選択に統合
	CustomizationModeSkillSelect     = ModeSkillSelect
	CustomizationModeChainSelect     = ModeChainSelect
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
	coreTypes     map[string]domain.CoreType
	skillTypes    map[string]domain.SkillType
	passiveSkills map[string]domain.PassiveSkill
	chainEffects  map[string]domain.ChainEffect

	// 現在のモード
	currentMode CustomizationMode

	// スロット選択
	selectedSlotIndex int

	// カード内フォーカス位置（0=コア、1-4=スキルスロット）
	focusPosition int

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
	passiveSkills map[string]domain.PassiveSkill,
	chainEffects map[string]domain.ChainEffect,
) *AgentCustomizationScreen {
	screen := &AgentCustomizationScreen{
		invManager:    invManager,
		slotManager:   slotManager,
		coreTypes:     coreTypes,
		skillTypes:    skillTypes,
		passiveSkills: passiveSkills,
		chainEffects:  chainEffects,
		currentMode:   ModeCardSelect,
		styles:        styles.NewGameStyles(),
		width:         140,
		height:        40,
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
	case ModeCardSelect:
		return s.handleCardSelectKey(msg)
	case ModeCoreSelect:
		return s.handleCoreSelectKey(msg)
	case ModeLevelSelect:
		return s.handleLevelSelectKey(msg)
	case ModeSkillSelect:
		return s.handleSkillSelectKey(msg)
	case ModeChainSelect:
		return s.handleChainSelectKey(msg)
	}

	return s, nil
}

// ==================== カード選択モード ====================

func (s *AgentCustomizationScreen) handleCardSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "left", "h":
		// スロット切替（左）
		if s.selectedSlotIndex > 0 {
			s.selectedSlotIndex--
			s.focusPosition = 0 // スロット変更時はコアにフォーカス
		}
	case "right", "l":
		// スロット切替（右）
		if s.selectedSlotIndex < slot.MaxAgentSlotCount-1 {
			s.selectedSlotIndex++
			s.focusPosition = 0 // スロット変更時はコアにフォーカス
		}
	case "up", "k":
		// カード内フォーカス移動（上）
		if s.focusPosition > 0 {
			s.focusPosition--
		}
	case "down", "j":
		// カード内フォーカス移動（下）
		// コアが設定されていない場合はコア位置から動かさない
		currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
		if currentSlot != nil && !currentSlot.IsEmpty() {
			if s.focusPosition < domain.MaxSkillSlotCount {
				s.focusPosition++
			}
		}
	case "enter":
		// モーダルを開く
		if s.focusPosition == 0 {
			// コア選択
			s.enterCoreSelectMode()
		} else {
			// スキル選択
			currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
			if currentSlot != nil && !currentSlot.IsEmpty() {
				s.enterSkillSelectMode(s.focusPosition - 1)
			} else {
				s.errorMessage = "先にコアを設定してください"
			}
		}
	case "delete", "backspace", "d":
		// 選択中のコア/スキルを外す
		if s.focusPosition == 0 {
			s.clearCurrentSlotCore()
		} else {
			currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
			if currentSlot != nil && !currentSlot.IsEmpty() {
				s.clearSkillSlot(s.focusPosition - 1)
			}
		}
	}

	return s, nil
}

// enterCoreSelectMode はコア選択モードに遷移します。
func (s *AgentCustomizationScreen) enterCoreSelectMode() {
	s.currentMode = ModeCoreSelect
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
		s.focusPosition = 0 // コア位置に戻す
	}
}

// clearSkillSlot は指定スキルスロットをクリアします。
func (s *AgentCustomizationScreen) clearSkillSlot(skillSlotIndex int) {
	if err := s.slotManager.ClearSkill(s.selectedSlotIndex, skillSlotIndex); err != nil {
		s.errorMessage = fmt.Sprintf("スキルクリアに失敗: %v", err)
	} else {
		s.statusMessage = fmt.Sprintf("スキルスロット%dをクリアしました", skillSlotIndex+1)
		s.errorMessage = ""
	}
}

// ==================== コア選択モード ====================

func (s *AgentCustomizationScreen) handleCoreSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = ModeCardSelect
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
			s.currentMode = ModeLevelSelect
		}
	}

	return s, nil
}

// ==================== レベル選択モード ====================

func (s *AgentCustomizationScreen) handleLevelSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = ModeCoreSelect
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
			s.currentMode = ModeCardSelect
		}
	}

	return s, nil
}

// ==================== スキル選択モード ====================

// enterSkillSelectMode はスキル選択モードに遷移します。
func (s *AgentCustomizationScreen) enterSkillSelectMode(skillSlotIndex int) {
	s.currentMode = ModeSkillSelect
	s.selectedSkillIndex = 0
	// focusPosition - 1 がスキルスロットインデックス
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
		s.currentMode = ModeCardSelect
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
				s.currentMode = ModeChainSelect
			} else {
				// チェイン効果なしで直接設定
				s.setSkillToSlot("")
			}
		}
	}

	return s, nil
}

// getSelectedSkillSlotIndex は現在選択中のスキルスロットインデックスを返します。
func (s *AgentCustomizationScreen) getSelectedSkillSlotIndex() int {
	if s.focusPosition > 0 {
		return s.focusPosition - 1
	}
	return 0
}

// 後方互換性のため
var _ = (*AgentCustomizationScreen).selectedSkillSlotIndex

func (s *AgentCustomizationScreen) selectedSkillSlotIndex() int {
	return s.getSelectedSkillSlotIndex()
}

// ==================== チェイン効果選択モード ====================

func (s *AgentCustomizationScreen) handleChainSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = ModeSkillSelect
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
	skillSlotIndex := s.getSelectedSkillSlotIndex()
	if err := s.slotManager.SetSkill(s.selectedSlotIndex, skillSlotIndex, s.selectedSkillTypeID, chainEffectID); err != nil {
		s.errorMessage = fmt.Sprintf("スキル設定に失敗: %v", err)
	} else {
		s.statusMessage = "スキルを設定しました"
		s.errorMessage = ""
		s.currentMode = ModeCardSelect
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

	// 3カード横並びエリア
	builder.WriteString(s.renderCardArea())
	builder.WriteString("\n\n")

	// モーダル（ModeCardSelect以外の場合に表示）
	if s.currentMode != ModeCardSelect {
		builder.WriteString(s.renderModal())
		builder.WriteString("\n")
	}

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

// renderCardArea は3カード横並びエリアをレンダリングします。
func (s *AgentCustomizationScreen) renderCardArea() string {
	// カード幅を計算（battle_view.goと同じ計算式）
	cardWidth := (s.width - 18) / 3
	if cardWidth < 30 {
		cardWidth = 30
	}

	var cards []string

	for i := 0; i < slot.MaxAgentSlotCount; i++ {
		isSelected := i == s.selectedSlotIndex
		card := s.renderAgentCard(i, isSelected, cardWidth)
		cards = append(cards, card)
	}

	// カードを横に並べる
	agentCards := lipgloss.JoinHorizontal(lipgloss.Top, cards[0], " ", cards[1], " ", cards[2])

	// 中央揃え
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(agentCards)
}

// renderAgentCard は1つのエージェントカードをレンダリングします。
func (s *AgentCustomizationScreen) renderAgentCard(slotIndex int, isSelected bool, cardWidth int) string {
	var cardContent strings.Builder

	agentSlot := s.slotManager.GetSlot(slotIndex)

	if agentSlot == nil || agentSlot.IsEmpty() {
		// 空スロット
		emptyStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
		cardContent.WriteString(emptyStyle.Render("(空)"))
		cardContent.WriteString("\n\n")
		cardContent.WriteString(emptyStyle.Render("Enterでコアを選択"))
	} else {
		// コア名とレベル
		coreName := agentSlot.CoreTypeID
		if coreType, ok := s.coreTypes[agentSlot.CoreTypeID]; ok {
			coreName = coreType.Name
		}

		// コア行のスタイル
		coreStyle := lipgloss.NewStyle()
		corePrefix := "  "
		if isSelected && s.focusPosition == 0 {
			coreStyle = coreStyle.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			corePrefix = "> "
		} else {
			coreStyle = coreStyle.Bold(true)
		}
		cardContent.WriteString(coreStyle.Render(fmt.Sprintf("%s%s Lv.%d", corePrefix, coreName, agentSlot.CoreLevel)))
		cardContent.WriteString("\n")

		// パッシブスキル表示（PassiveSkillIDがある場合のみ）
		if coreType, ok := s.coreTypes[agentSlot.CoreTypeID]; ok {
			if coreType.PassiveSkillID != "" {
				passiveStyle := lipgloss.NewStyle().
					Foreground(styles.ColorBuff).
					Bold(true)
				// パッシブスキルの短い説明を表示
				passiveText := coreType.PassiveSkillID
				if ps, ok := s.passiveSkills[coreType.PassiveSkillID]; ok {
					passiveText = ps.ShortDescription
				}
				cardContent.WriteString(passiveStyle.Render(fmt.Sprintf("  ★%s", passiveText)))
				cardContent.WriteString("\n")
			}
		}
		cardContent.WriteString("\n")

		// スキルスロット表示
		for j := 0; j < domain.MaxSkillSlotCount; j++ {
			skillConfig := agentSlot.GetSkill(j)

			skillStyle := lipgloss.NewStyle()
			prefix := "  "

			// フォーカス位置判定（j+1 がfocusPosition）
			if isSelected && s.focusPosition == j+1 {
				skillStyle = skillStyle.Bold(true).
					Foreground(styles.ColorSelectedFg).
					Background(styles.ColorSelectedBg)
				prefix = "> "
			}

			var skillContent string
			if skillConfig == nil || skillConfig.IsEmpty() {
				skillContent = fmt.Sprintf("スキル%d: (空)", j+1)
				if !isSelected || s.focusPosition != j+1 {
					skillStyle = skillStyle.Foreground(styles.ColorSubtle)
				}
			} else {
				skillName := skillConfig.TypeID
				icon := "?"
				if skillType, ok := s.skillTypes[skillConfig.TypeID]; ok {
					skillName = skillType.Name
					icon = skillType.Icon
				}
				skillContent = fmt.Sprintf("%s %s", icon, skillName)
			}

			cardContent.WriteString(skillStyle.Render(prefix + skillContent))
			cardContent.WriteString("\n")

			// チェイン効果表示（2行目）
			if skillConfig != nil && !skillConfig.IsEmpty() && skillConfig.ChainEffectID != "" {
				// チェイン効果の短い説明を表示
				chainStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)
				chainText := skillConfig.ChainEffectID
				if ce, ok := s.chainEffects[skillConfig.ChainEffectID]; ok {
					chainText = ce.ShortDescription
				}
				cardContent.WriteString("    ")
				cardContent.WriteString(chainStyle.Render(fmt.Sprintf("[%s]", chainText)))
				cardContent.WriteString("\n")
			} else {
				// チェイン効果がない場合は空行
				cardContent.WriteString("\n")
			}
		}
	}

	// カードボックス
	borderColor := styles.ColorSubtle
	if isSelected {
		borderColor = styles.ColorPrimary
	}

	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(0, 1).
		Width(cardWidth).
		Height(14) // 固定高さ

	return cardStyle.Render(cardContent.String())
}

// renderModal はモーダルをレンダリングします。
func (s *AgentCustomizationScreen) renderModal() string {
	// モーダルコンテンツ
	var content string
	var title string

	switch s.currentMode {
	case ModeCoreSelect:
		title = "コア選択"
		content = s.renderModalCoreSelect()
	case ModeLevelSelect:
		title = "レベル選択"
		content = s.renderModalLevelSelect()
	case ModeSkillSelect:
		title = "スキル選択"
		content = s.renderModalSkillSelect()
	case ModeChainSelect:
		title = "チェイン効果選択"
		content = s.renderModalChainSelect()
	}

	// タイトル行
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary).
		Align(lipgloss.Center).
		Width(80)
	titleRendered := titleStyle.Render(title)

	// モーダルボックス（二重線枠）
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Width(84)

	modalContent := titleRendered + "\n\n" + content

	// ヒント行
	var hintText string
	switch s.currentMode {
	case ModeCoreSelect, ModeSkillSelect:
		hintText = "↑/↓: 選択  Enter: 決定  Esc: キャンセル"
	case ModeLevelSelect:
		hintText = "↑/↓: レベル選択  Enter: 決定  Esc: 戻る"
	case ModeChainSelect:
		hintText = "↑/↓: 効果選択  Enter: 決定  Esc: 戻る"
	}
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(80)
	modalContent += "\n" + hintStyle.Render(hintText)

	modal := modalStyle.Render(modalContent)

	// 中央揃え
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(modal)
}

// truncateLines はコンテンツの行数を制限します。
func truncateLines(content string, maxLines int) string {
	// 末尾の改行を削除
	content = strings.TrimSuffix(content, "\n")
	lines := strings.Split(content, "\n")

	// 行数を制限
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	// 足りない行を埋める
	for len(lines) < maxLines {
		lines = append(lines, "")
	}
	return strings.Join(lines, "\n")
}

// renderTwoPanelLayout は左右2パネルのレイアウトをレンダリングします。
// 枠線問題を避けるため、内容のみを横に並べてから外枠を付けます。
func renderTwoPanelLayout(leftContent, rightContent string, panelWidth int) string {
	// 各パネルを固定幅でスタイル設定
	leftStyle := lipgloss.NewStyle().
		Width(panelWidth).
		Padding(0, 1)
	rightStyle := lipgloss.NewStyle().
		Width(panelWidth).
		Padding(0, 1).
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.ColorSubtle)

	// コンテンツをレンダリング
	left := leftStyle.Render(leftContent)
	right := rightStyle.Render(rightContent)

	// 横に並べる
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

// renderModalCoreSelect はモーダル内のコア選択をレンダリングします。
func (s *AgentCustomizationScreen) renderModalCoreSelect() string {
	leftContent := truncateLines(s.renderCoreList(), 10)
	rightContent := truncateLines(s.renderCoreDetail(), 10)
	return renderTwoPanelLayout(leftContent, rightContent, 38)
}

// renderCoreList はコアリストをレンダリングします。
func (s *AgentCustomizationScreen) renderCoreList() string {
	var builder strings.Builder

	if len(s.coreList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("保有コアがありません"))
		return builder.String()
	}

	// 表示可能な行数（Height - パディング）
	maxVisibleItems := 10

	// スクロール位置を計算
	startIdx := 0
	if s.selectedCoreIndex >= maxVisibleItems {
		startIdx = s.selectedCoreIndex - maxVisibleItems + 1
	}

	endIdx := startIdx + maxVisibleItems
	if endIdx > len(s.coreList) {
		endIdx = len(s.coreList)
	}

	for i := startIdx; i < endIdx; i++ {
		core := s.coreList[i]
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedCoreIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		item := fmt.Sprintf("%s (Lv.%d)", core.TypeName, core.MaxLevel)
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderCoreDetail は選択中コアの詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderCoreDetail() string {
	var builder strings.Builder

	if len(s.coreList) == 0 || s.selectedCoreIndex >= len(s.coreList) {
		return ""
	}

	core := s.coreList[s.selectedCoreIndex]

	// コア名
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(nameStyle.Render(core.TypeName))
	builder.WriteString("\n\n")

	// 最大レベル
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
	builder.WriteString(labelStyle.Render("最大レベル: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("Lv.%d", core.MaxLevel)))
	builder.WriteString("\n\n")

	// パッシブスキル情報（PassiveSkillIDがある場合のみ）
	if coreType, ok := s.coreTypes[core.TypeID]; ok {
		if coreType.PassiveSkillID != "" {
			passiveStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff).Bold(true)
			// パッシブスキルの説明を表示
			passiveText := coreType.PassiveSkillID
			if ps, ok := s.passiveSkills[coreType.PassiveSkillID]; ok {
				passiveText = ps.Description
			}
			builder.WriteString(passiveStyle.Render("★ " + passiveText))
		}
	}

	return builder.String()
}

// renderModalLevelSelect はモーダル内のレベル選択をレンダリングします。
func (s *AgentCustomizationScreen) renderModalLevelSelect() string {
	leftContent := truncateLines(s.renderLevelList(), 10)
	rightContent := truncateLines(s.renderLevelDetail(), 10)
	return renderTwoPanelLayout(leftContent, rightContent, 38)
}

// renderLevelList はレベル選択UIをレンダリングします（数値増減方式）。
func (s *AgentCustomizationScreen) renderLevelList() string {
	var builder strings.Builder

	// 現在選択中のレベル
	selectedLevel := s.selectedLevelIndex + 1

	// 数値増減UI
	builder.WriteString("\n")

	// 上矢印（増加可能な場合）
	arrowStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	activeArrowStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary).Bold(true)

	if selectedLevel < s.maxSelectableLevel {
		builder.WriteString(activeArrowStyle.Render("        ▲"))
	} else {
		builder.WriteString(arrowStyle.Render("        ▲"))
	}
	builder.WriteString("\n\n")

	// 現在のレベル（大きく表示）
	levelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Background(styles.ColorSelectedBg).
		Padding(0, 2)
	builder.WriteString(levelStyle.Render(fmt.Sprintf("Lv.%d", selectedLevel)))
	builder.WriteString("\n\n")

	// 下矢印（減少可能な場合）
	if selectedLevel > 1 {
		builder.WriteString(activeArrowStyle.Render("        ▼"))
	} else {
		builder.WriteString(arrowStyle.Render("        ▼"))
	}
	builder.WriteString("\n\n")

	// 範囲表示
	rangeStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	builder.WriteString(rangeStyle.Render(fmt.Sprintf("  (1 〜 %d)", s.maxSelectableLevel)))

	return builder.String()
}

// renderLevelDetail はレベル詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderLevelDetail() string {
	var builder strings.Builder

	// コア名
	coreName := s.selectedCoreTypeID
	if coreType, ok := s.coreTypes[s.selectedCoreTypeID]; ok {
		coreName = coreType.Name
	}

	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(nameStyle.Render(coreName))
	builder.WriteString("\n\n")

	// 選択中のレベル
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary).Bold(true)
	selectedLevel := s.selectedLevelIndex + 1
	builder.WriteString(labelStyle.Render("選択レベル: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("Lv.%d", selectedLevel)))

	return builder.String()
}

// renderModalSkillSelect はモーダル内のスキル選択をレンダリングします。
func (s *AgentCustomizationScreen) renderModalSkillSelect() string {
	leftContent := truncateLines(s.renderSkillList(), 10)
	rightContent := truncateLines(s.renderSkillDetail(), 10)
	return renderTwoPanelLayout(leftContent, rightContent, 38)
}

// renderSkillList はスキルリストをレンダリングします。
func (s *AgentCustomizationScreen) renderSkillList() string {
	var builder strings.Builder

	if len(s.compatibleSkillList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("互換スキルがありません"))
		return builder.String()
	}

	// 表示可能な行数
	maxVisibleItems := 10

	// スクロール位置を計算
	startIdx := 0
	if s.selectedSkillIndex >= maxVisibleItems {
		startIdx = s.selectedSkillIndex - maxVisibleItems + 1
	}

	endIdx := startIdx + maxVisibleItems
	if endIdx > len(s.compatibleSkillList) {
		endIdx = len(s.compatibleSkillList)
	}

	for i := startIdx; i < endIdx; i++ {
		skill := s.compatibleSkillList[i]
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedSkillIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		item := fmt.Sprintf("%s %s", skill.Icon, skill.TypeName)
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderSkillDetail は選択中スキルの詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderSkillDetail() string {
	var builder strings.Builder

	if len(s.compatibleSkillList) == 0 || s.selectedSkillIndex >= len(s.compatibleSkillList) {
		return ""
	}

	skill := s.compatibleSkillList[s.selectedSkillIndex]

	// スキル名
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(nameStyle.Render(fmt.Sprintf("%s %s", skill.Icon, skill.TypeName)))
	builder.WriteString("\n\n")

	// スキルタイプ情報
	if skillType, ok := s.skillTypes[skill.TypeID]; ok {
		labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
		// 説明
		builder.WriteString(labelStyle.Render(skillType.Description))
		builder.WriteString("\n\n")

		// チェイン効果情報
		if len(skill.ChainVariations) > 0 {
			chainStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)
			builder.WriteString(chainStyle.Render(fmt.Sprintf("🔗 %d種のチェイン効果あり", len(skill.ChainVariations))))
		}
	}

	return builder.String()
}

// renderModalChainSelect はモーダル内のチェイン効果選択をレンダリングします。
func (s *AgentCustomizationScreen) renderModalChainSelect() string {
	leftContent := truncateLines(s.renderChainList(), 10)
	rightContent := truncateLines(s.renderChainDetail(), 10)
	return renderTwoPanelLayout(leftContent, rightContent, 38)
}

// renderChainList はチェイン効果リストをレンダリングします（スクロール対応）。
func (s *AgentCustomizationScreen) renderChainList() string {
	var builder strings.Builder

	// 全アイテム数（チェイン効果 + 「なし」オプション）
	totalItems := len(s.chainVariationList) + 1

	// 表示可能な行数（Height - パディング）
	maxVisibleItems := 10

	// スクロール位置を計算
	startIdx := 0
	if s.selectedChainIndex >= maxVisibleItems {
		startIdx = s.selectedChainIndex - maxVisibleItems + 1
	}

	endIdx := startIdx + maxVisibleItems
	if endIdx > totalItems {
		endIdx = totalItems
	}

	// チェイン効果リスト
	for i := startIdx; i < endIdx; i++ {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedChainIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		var itemText string
		if i < len(s.chainVariationList) {
			// チェイン効果の短い説明を表示
			chainID := s.chainVariationList[i]
			itemText = chainID
			if ce, ok := s.chainEffects[chainID]; ok {
				itemText = ce.ShortDescription
			}
		} else {
			// 「なし」オプション
			itemText = "(なし)"
		}
		builder.WriteString(style.Render(prefix + itemText))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderChainDetail は選択中チェイン効果の詳細をレンダリングします。
func (s *AgentCustomizationScreen) renderChainDetail() string {
	var builder strings.Builder

	// スキル名
	skillName := s.selectedSkillTypeID
	if skillType, ok := s.skillTypes[s.selectedSkillTypeID]; ok {
		skillName = skillType.Name
	}
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(nameStyle.Render(skillName))
	builder.WriteString("\n\n")

	// 選択中のチェイン効果
	if s.selectedChainIndex < len(s.chainVariationList) {
		chainID := s.chainVariationList[s.selectedChainIndex]

		// チェイン効果の説明を表示
		chainStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff).Bold(true)
		chainName := chainID
		chainDesc := ""
		if ce, ok := s.chainEffects[chainID]; ok {
			chainName = ce.ShortDescription
			chainDesc = ce.Description
		}
		builder.WriteString(chainStyle.Render(chainName))
		builder.WriteString("\n\n")

		// チェイン効果の長い説明
		if chainDesc != "" {
			descStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
			builder.WriteString(descStyle.Render(chainDesc))
		}
	} else {
		// 「なし」選択中
		labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
		builder.WriteString(labelStyle.Render("チェイン効果なしで装備します"))
	}

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
	case ModeCardSelect:
		hints = "←/→: スロット切替  ↑/↓: 項目選択  Enter: 編集  Delete: 外す  Esc: 戻る"
	case ModeCoreSelect:
		hints = "↑/↓: コア選択  Enter: 決定  Esc: キャンセル"
	case ModeLevelSelect:
		hints = "↑/↓: レベル選択  Enter: 決定  Esc: 戻る"
	case ModeSkillSelect:
		hints = "↑/↓: スキル選択  Enter: 決定  Esc: キャンセル"
	case ModeChainSelect:
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
