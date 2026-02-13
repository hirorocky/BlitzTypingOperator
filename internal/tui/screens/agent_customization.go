// Package screens はTUIゲームの画面を提供します。
// このファイルはAgentCustomizationScreen（エージェントカスタマイズ画面）を定義します。

package screens

import (
	"fmt"
	"sort"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/components"
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
	// ModeSkillSelect はモーダル: スキル選択モードです。
	ModeSkillSelect
	// ModeChainSelect はモーダル: チェイン効果選択モードです。
	ModeChainSelect
)

// ==================== 表示用アイテム定義 ====================

// CoreSelectItem はコア選択リストの表示アイテムです。
type CoreSelectItem struct {
	TypeID              string
	TypeName            string
	IsEquippedElsewhere bool
}

// SkillSelectItem はスキル選択リストの表示アイテムです。
type SkillSelectItem struct {
	TypeID       string
	TypeName     string
	Icon         string
	IsCompatible bool
}

// ChainEffectSelectItem はチェイン効果選択リストの表示アイテムです。
type ChainEffectSelectItem struct {
	TypeID              string
	ShortDescription    string
	IsEquippedElsewhere bool
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

	// カード内フォーカス位置
	// 0=コア、奇数=スキル(1,3,5,7)、偶数(2,4,6,8)=チェイン効果
	// スキルインデックス: (focusPosition-1)/2、チェインインデックス: (focusPosition-2)/2
	focusPosition int

	// コア選択
	coreList          []CoreSelectItem
	selectedCoreIndex int

	// スキル選択
	compatibleSkillList []SkillSelectItem
	selectedSkillIndex  int
	selectedSkillTypeID string

	// チェイン効果選択
	chainEffectList    []ChainEffectSelectItem
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
			maxPos := 2 * domain.MaxSkillSlotCount // 8
			if s.focusPosition < maxPos {
				s.focusPosition++
			}
		}
	case "enter":
		// モーダルを開く
		if s.focusPosition == 0 {
			// コア選択
			s.enterCoreSelectMode()
		} else if s.isChainFocusPosition() {
			// チェイン効果選択（独立フロー）
			s.enterChainSelectMode()
		} else {
			// スキル選択
			currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
			if currentSlot != nil && !currentSlot.IsEmpty() {
				s.enterSkillSelectMode(s.getSelectedSkillSlotIndex())
			} else {
				s.errorMessage = "先にコアを設定してください"
			}
		}
	case "delete", "backspace", "d":
		// 選択中のコア/スキル/チェイン効果を外す
		if s.focusPosition == 0 {
			s.clearCurrentSlotCore()
		} else if s.isChainFocusPosition() {
			s.clearChainEffectSlot()
		} else {
			currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
			if currentSlot != nil && !currentSlot.IsEmpty() {
				s.clearSkillSlot(s.getSelectedSkillSlotIndex())
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
	sort.Strings(ownedCores)

	// リストを構築
	for _, typeID := range ownedCores {
		// コアタイプ名を取得
		typeName := typeID
		if coreType, ok := s.coreTypes[typeID]; ok {
			typeName = coreType.Name
		}

		item := CoreSelectItem{
			TypeID:              typeID,
			TypeName:            typeName,
			IsEquippedElsewhere: s.isCoreEquippedElsewhere(typeID),
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
			// コアを直接スロットに設定
			typeID := s.coreList[s.selectedCoreIndex].TypeID
			if err := s.slotManager.SetCore(s.selectedSlotIndex, typeID); err != nil {
				s.errorMessage = fmt.Sprintf("コア設定に失敗: %v", err)
			} else {
				s.statusMessage = fmt.Sprintf("スロット%dにコアを設定しました", s.selectedSlotIndex+1)
				s.errorMessage = ""
				s.currentMode = ModeCardSelect
			}
		}
	}

	return s, nil
}

// ==================== スキル選択モード ====================

// enterSkillSelectMode はスキル選択モードに遷移します。
func (s *AgentCustomizationScreen) enterSkillSelectMode(_ int) {
	s.currentMode = ModeSkillSelect
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
		_, exists := ownedSkills[typeID]
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

		item := SkillSelectItem{
			TypeID:       typeID,
			TypeName:     typeName,
			Icon:         icon,
			IsCompatible: true,
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
			s.setSkillToSlot()
		}
	}

	return s, nil
}

// isChainFocusPosition はフォーカスがチェイン効果位置にあるかを返します。
// チェイン効果位置: 2, 4, 6, 8（偶数かつ2以上）
func (s *AgentCustomizationScreen) isChainFocusPosition() bool {
	return s.focusPosition >= 2 && s.focusPosition%2 == 0
}

// getSelectedSkillSlotIndex は現在のfocusPositionに対応するスキルスロットインデックスを返します。
// スキル位置(1,3,5,7): (pos-1)/2、チェイン位置(2,4,6,8): (pos-2)/2
func (s *AgentCustomizationScreen) getSelectedSkillSlotIndex() int {
	if s.focusPosition <= 0 {
		return 0
	}
	if s.isChainFocusPosition() {
		return (s.focusPosition - 2) / 2
	}
	return (s.focusPosition - 1) / 2
}

// isCoreEquippedElsewhere は指定コアが他のスロットで装備されているかを返します。
func (s *AgentCustomizationScreen) isCoreEquippedElsewhere(typeID string) bool {
	for i := 0; i < slot.MaxAgentSlotCount; i++ {
		if i == s.selectedSlotIndex {
			continue
		}
		agentSlot := s.slotManager.GetSlot(i)
		if agentSlot != nil && agentSlot.CoreTypeID == typeID {
			return true
		}
	}
	return false
}

// isChainEffectEquippedElsewhere は指定チェイン効果が他の位置で装備されているかを返します。
func (s *AgentCustomizationScreen) isChainEffectEquippedElsewhere(typeID string) bool {
	chainSlotIdx := s.getSelectedSkillSlotIndex()
	for i := 0; i < slot.MaxAgentSlotCount; i++ {
		agentSlot := s.slotManager.GetSlot(i)
		if agentSlot == nil || agentSlot.IsEmpty() {
			continue
		}
		for j := 0; j < domain.MaxSkillSlotCount; j++ {
			if i == s.selectedSlotIndex && j == chainSlotIdx {
				continue
			}
			chainCfg := agentSlot.GetChainEffect(j)
			if chainCfg != nil && !chainCfg.IsEmpty() && chainCfg.TypeID == typeID {
				return true
			}
		}
	}
	return false
}

// ==================== チェイン効果選択モード ====================

// enterChainSelectMode はチェイン効果選択モードに遷移します。
func (s *AgentCustomizationScreen) enterChainSelectMode() {
	currentSlot := s.slotManager.GetSlot(s.selectedSlotIndex)
	if currentSlot == nil || currentSlot.IsEmpty() {
		s.errorMessage = "先にコアを設定してください"
		return
	}

	chainSlotIdx := s.getSelectedSkillSlotIndex()
	skillConfig := currentSlot.GetSkill(chainSlotIdx)
	if skillConfig == nil || skillConfig.IsEmpty() {
		s.errorMessage = "先にスキルを設定してください"
		return
	}

	s.currentMode = ModeChainSelect
	s.selectedChainIndex = 0
	s.updateChainEffectList()
	s.errorMessage = ""
	s.statusMessage = ""
}

// updateChainEffectList は保有チェイン効果リストを更新します。
func (s *AgentCustomizationScreen) updateChainEffectList() {
	s.chainEffectList = []ChainEffectSelectItem{}

	if s.invManager == nil || s.invManager.ChainEffects() == nil {
		return
	}

	ownedChainEffects := s.invManager.ChainEffects().GetOwnedChainEffects()
	sort.Strings(ownedChainEffects)

	for _, typeID := range ownedChainEffects {
		shortDesc := typeID
		if ce, ok := s.chainEffects[typeID]; ok {
			shortDesc = ce.ShortDescription
		}

		item := ChainEffectSelectItem{
			TypeID:              typeID,
			ShortDescription:    shortDesc,
			IsEquippedElsewhere: s.isChainEffectEquippedElsewhere(typeID),
		}
		s.chainEffectList = append(s.chainEffectList, item)
	}
}

func (s *AgentCustomizationScreen) handleChainSelectKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		s.currentMode = ModeCardSelect
	case "up", "k":
		if s.selectedChainIndex > 0 {
			s.selectedChainIndex--
		}
	case "down", "j":
		// +1 は「なし」オプション用
		if s.selectedChainIndex < len(s.chainEffectList) {
			s.selectedChainIndex++
		}
	case "enter":
		s.setChainEffectToSlot()
	}

	return s, nil
}

// setChainEffectToSlot はチェイン効果をスロットに設定します。
func (s *AgentCustomizationScreen) setChainEffectToSlot() {
	chainSlotIdx := s.getSelectedSkillSlotIndex()

	if s.selectedChainIndex < len(s.chainEffectList) {
		// チェイン効果を設定
		typeID := s.chainEffectList[s.selectedChainIndex].TypeID
		if err := s.slotManager.SetChainEffect(s.selectedSlotIndex, chainSlotIdx, typeID); err != nil {
			s.errorMessage = fmt.Sprintf("チェイン効果設定に失敗: %v", err)
		} else {
			s.statusMessage = "チェイン効果を設定しました"
			s.errorMessage = ""
			s.currentMode = ModeCardSelect
		}
	} else {
		// 「なし」を選択: チェイン効果をクリア
		if err := s.slotManager.ClearChainEffect(s.selectedSlotIndex, chainSlotIdx); err != nil {
			s.errorMessage = fmt.Sprintf("チェイン効果クリアに失敗: %v", err)
		} else {
			s.statusMessage = "チェイン効果をクリアしました"
			s.errorMessage = ""
			s.currentMode = ModeCardSelect
		}
	}
}

// clearChainEffectSlot はフォーカス中のチェイン効果スロットをクリアします。
func (s *AgentCustomizationScreen) clearChainEffectSlot() {
	chainSlotIdx := s.getSelectedSkillSlotIndex()
	if err := s.slotManager.ClearChainEffect(s.selectedSlotIndex, chainSlotIdx); err != nil {
		s.errorMessage = fmt.Sprintf("チェイン効果クリアに失敗: %v", err)
	} else {
		s.statusMessage = fmt.Sprintf("チェイン効果スロット%dをクリアしました", chainSlotIdx+1)
		s.errorMessage = ""
	}
}

// setSkillToSlot はスキルをスロットに設定します。
func (s *AgentCustomizationScreen) setSkillToSlot() {
	skillSlotIndex := s.getSelectedSkillSlotIndex()
	if err := s.slotManager.SetSkill(s.selectedSlotIndex, skillSlotIndex, s.selectedSkillTypeID); err != nil {
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
		cardContent.WriteString(coreStyle.Render(fmt.Sprintf("%s%s", corePrefix, coreName)))
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

		// スキルスロット表示（各スキルとチェイン効果を交互に表示）
		for j := 0; j < domain.MaxSkillSlotCount; j++ {
			skillConfig := agentSlot.GetSkill(j)
			skillFocusPos := 2*j + 1 // 1, 3, 5, 7
			chainFocusPos := 2*j + 2 // 2, 4, 6, 8

			// スキル行
			skillStyle := lipgloss.NewStyle()
			prefix := "  "

			if isSelected && s.focusPosition == skillFocusPos {
				skillStyle = skillStyle.Bold(true).
					Foreground(styles.ColorSelectedFg).
					Background(styles.ColorSelectedBg)
				prefix = "> "
			}

			var skillContent string
			if skillConfig == nil || skillConfig.IsEmpty() {
				skillContent = fmt.Sprintf("スキル%d: (空)", j+1)
				if !isSelected || s.focusPosition != skillFocusPos {
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

			// チェイン効果行
			chainEffectCfg := agentSlot.GetChainEffect(j)
			chainStyle := lipgloss.NewStyle()
			chainPrefix := "    "

			if isSelected && s.focusPosition == chainFocusPos {
				chainStyle = chainStyle.Bold(true).
					Foreground(styles.ColorSelectedFg).
					Background(styles.ColorSelectedBg)
				chainPrefix = "  > "
			}

			if chainEffectCfg != nil && !chainEffectCfg.IsEmpty() {
				if !isSelected || s.focusPosition != chainFocusPos {
					chainStyle = chainStyle.Foreground(styles.ColorBuff)
				}
				chainText := chainEffectCfg.TypeID
				if ce, ok := s.chainEffects[chainEffectCfg.TypeID]; ok {
					chainText = ce.ShortDescription
				}
				cardContent.WriteString(chainStyle.Render(fmt.Sprintf("%s[%s]", chainPrefix, chainText)))
			} else if isSelected && s.focusPosition == chainFocusPos {
				// フォーカス中の空チェイン効果スロット
				cardContent.WriteString(chainStyle.Render(chainPrefix + "[チェイン: 空]"))
			} else {
				cardContent.WriteString("")
			}
			cardContent.WriteString("\n")
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

		displayName := core.TypeName
		if core.IsEquippedElsewhere {
			displayName += " (装備中)"
		}
		builder.WriteString(style.Render(prefix + displayName))
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

	// ステータス表示
	if coreType, ok := s.coreTypes[core.TypeID]; ok {
		statsStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
		builder.WriteString(statsStyle.Render(components.FormatStats(coreType)))
		builder.WriteString("\n\n")

		// パッシブスキル情報（PassiveSkillIDがある場合のみ）
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

		// チェイン効果情報はチェイン効果タブで確認
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

	if len(s.chainEffectList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("保有チェイン効果がありません"))
		return builder.String()
	}

	// 全アイテム数（チェイン効果 + 「なし」オプション）
	totalItems := len(s.chainEffectList) + 1

	// 表示可能な行数
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
		if i < len(s.chainEffectList) {
			ce := s.chainEffectList[i]
			itemText = ce.ShortDescription
			if ce.IsEquippedElsewhere {
				itemText += " (装備中)"
			}
		} else {
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

	if s.selectedChainIndex < len(s.chainEffectList) {
		ce := s.chainEffectList[s.selectedChainIndex]

		// チェイン効果名
		chainStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff).Bold(true)
		builder.WriteString(chainStyle.Render(ce.ShortDescription))
		builder.WriteString("\n\n")

		// チェイン効果の説明
		if ceData, ok := s.chainEffects[ce.TypeID]; ok && ceData.Description != "" {
			descStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
			builder.WriteString(descStyle.Render(ceData.Description))
		}

		// 装備中の警告
		if ce.IsEquippedElsewhere {
			builder.WriteString("\n\n")
			warnStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
			builder.WriteString(warnStyle.Render("他のスロットで装備中"))
		}
	} else {
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
