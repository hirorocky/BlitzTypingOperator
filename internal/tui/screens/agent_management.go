// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"log/slog"
	"strings"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.4: エージェント管理画面 ====================

// AgentManagementTab はタブの種類を表します。
type AgentManagementTab int

const (
	// TabCoreList はコア一覧タブです。
	TabCoreList AgentManagementTab = iota
	// TabModuleList はモジュール一覧タブです。
	TabModuleList
	// TabSynthesis は合成タブです。
	TabSynthesis
	// TabEquip は装備タブです。
	TabEquip
)

// InventoryProvider はインベントリデータを提供するインターフェースです。
type InventoryProvider interface {
	GetCores() []*domain.CoreModel
	GetModules() []*domain.ModuleModel
	GetAgents() []*domain.AgentModel
	GetEquippedAgents() []*domain.AgentModel
	AddAgent(agent *domain.AgentModel) error
	RemoveCore(id string) error
	RemoveModule(id string) error
	EquipAgent(slot int, agent *domain.AgentModel) error
	UnequipAgent(slot int) error
}

// SynthesisState は合成状態を表します。
type SynthesisState struct {
	selectedCore    *domain.CoreModel
	selectedModules []*domain.ModuleModel
	step            int // 0: コア選択, 1: モジュール選択, 2: 確認
}

// AgentManagementScreen はエージェント管理画面を表します。

// UI-Improvement Requirements: 2.5, 2.6, 2.7, 2.8, 2.9
type AgentManagementScreen struct {
	inventory      InventoryProvider
	currentTab     AgentManagementTab
	selectedIndex  int
	coreList       []*domain.CoreModel
	moduleList     []*domain.ModuleModel
	agentList      []*domain.AgentModel
	equipSlots     []*domain.AgentModel
	synthesisState SynthesisState
	styles         *styles.GameStyles
	width          int
	height         int
	// UI改善: 確認ダイアログ
	confirmDialog    *components.ConfirmDialog
	pendingDeleteIdx int // 削除待ちのエージェントインデックス
	// 装備タブ用: 選択中のスロットインデックス (0-2)
	selectedEquipSlot int
	// エラー/ステータスメッセージ
	errorMessage  string
	statusMessage string
}

// NewAgentManagementScreen は新しいAgentManagementScreenを作成します。
func NewAgentManagementScreen(inventory InventoryProvider) *AgentManagementScreen {
	screen := &AgentManagementScreen{
		inventory:     inventory,
		currentTab:    TabCoreList,
		selectedIndex: 0,
		equipSlots:    make([]*domain.AgentModel, config.MaxAgentEquipSlots),
		synthesisState: SynthesisState{
			selectedModules: []*domain.ModuleModel{},
		},
		styles: styles.NewGameStyles(),
		width:  140,
		height: 40,
	}
	screen.updateCurrentList()
	return screen
}

// Init は画面の初期化を行います。
func (s *AgentManagementScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *AgentManagementScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
// UI-Improvement Requirement 2.9: 確認ダイアログ対応
func (s *AgentManagementScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 確認ダイアログが表示中の場合はダイアログのキー処理を優先
	if s.confirmDialog != nil && s.confirmDialog.Visible {
		result := s.confirmDialog.HandleKey(msg.String())
		switch result {
		case components.ConfirmResultYes:
			// 削除を実行
			s.executeDelete()
			return s, nil
		case components.ConfirmResultNo, components.ConfirmResultCancelled:
			// キャンセル
			return s, nil
		}
		return s, nil
	}

	// 装備タブの場合は専用の処理
	if s.currentTab == TabEquip {
		return s.handleEquipKeyMsg(msg)
	}

	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "backspace":
		return s.handleBackspace()
	case "left", "h":
		s.prevTab()
	case "right", "l":
		s.nextTab()
	case "up", "k":
		s.moveUp()
	case "down", "j":
		s.moveDown()
	case "enter":
		return s.handleEnter()
	case "d":
		return s.handleDelete()
	}

	return s, nil
}

// handleEquipKeyMsg は装備タブ専用のキー処理を行います。
// Tab: スロット切り替え、上下: エージェント選択、Enter: 装備、Backspace: 取り外し
func (s *AgentManagementScreen) handleEquipKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "tab":
		// スロット切り替え (0 → 1 → 2 → 0...)
		s.selectedEquipSlot = (s.selectedEquipSlot + 1) % 3
	case "left", "h":
		s.prevTab()
	case "right", "l":
		s.nextTab()
	case "up", "k":
		// エージェント一覧で上に移動
		if s.selectedIndex > 0 {
			s.selectedIndex--
		}
	case "down", "j":
		// エージェント一覧で下に移動
		if s.selectedIndex < len(s.agentList)-1 {
			s.selectedIndex++
		}
	case "enter":
		// 選択中のエージェントを選択中のスロットに装備
		if s.selectedIndex < len(s.agentList) {
			agent := s.agentList[s.selectedIndex]
			if err := s.inventory.EquipAgent(s.selectedEquipSlot, agent); err != nil {
				slog.Error("エージェント装備に失敗",
					slog.Int("slot", s.selectedEquipSlot),
					slog.String("agent_id", agent.ID),
					slog.Any("error", err),
				)
				s.errorMessage = fmt.Sprintf("装備に失敗しました: %v", err)
				s.statusMessage = ""
			} else {
				s.statusMessage = fmt.Sprintf("'%s'をスロット%dに装備しました", agent.GetCoreTypeName(), s.selectedEquipSlot+1)
				s.errorMessage = ""
			}
			s.updateCurrentList()
		}
	case "backspace":
		// 選択中のスロットからエージェントを取り外し
		if s.equipSlots[s.selectedEquipSlot] != nil {
			agentName := s.equipSlots[s.selectedEquipSlot].GetCoreTypeName()
			if err := s.inventory.UnequipAgent(s.selectedEquipSlot); err != nil {
				slog.Error("エージェント装備解除に失敗",
					slog.Int("slot", s.selectedEquipSlot),
					slog.Any("error", err),
				)
				s.errorMessage = fmt.Sprintf("装備解除に失敗しました: %v", err)
				s.statusMessage = ""
			} else {
				s.statusMessage = fmt.Sprintf("'%s'をスロット%dから取り外しました", agentName, s.selectedEquipSlot+1)
				s.errorMessage = ""
			}
			s.updateCurrentList()
		}
	case "d":
		return s.handleDelete()
	}

	return s, nil
}

// prevTab は前のタブに移動します。
func (s *AgentManagementScreen) prevTab() {
	if s.currentTab > TabCoreList {
		s.currentTab--
		s.selectedIndex = 0
		s.updateCurrentList()
	}
}

// nextTab は次のタブに移動します。
func (s *AgentManagementScreen) nextTab() {
	if s.currentTab < TabEquip {
		s.currentTab++
		s.selectedIndex = 0
		s.updateCurrentList()
	}
}

// moveUp は選択を上に移動します。
func (s *AgentManagementScreen) moveUp() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// moveDown は選択を下に移動します。
func (s *AgentManagementScreen) moveDown() {
	maxIndex := s.getMaxIndex()
	if s.selectedIndex < maxIndex-1 {
		s.selectedIndex++
	}
}

// getMaxIndex は現在のタブの最大インデックスを返します。
func (s *AgentManagementScreen) getMaxIndex() int {
	switch s.currentTab {
	case TabCoreList:
		return len(s.coreList)
	case TabModuleList:
		return len(s.moduleList)
	case TabSynthesis:
		if s.synthesisState.step == 0 {
			return len(s.coreList)
		}
		return len(s.moduleList)
	case TabEquip:
		return len(s.agentList) + len(s.equipSlots)
	}
	return 0
}

// updateCurrentList は現在のリストを更新します。
func (s *AgentManagementScreen) updateCurrentList() {
	s.coreList = s.inventory.GetCores()
	s.moduleList = s.inventory.GetModules()
	s.agentList = s.inventory.GetAgents()

	equipped := s.inventory.GetEquippedAgents()
	s.equipSlots = make([]*domain.AgentModel, config.MaxAgentEquipSlots)
	for i := 0; i < config.MaxAgentEquipSlots && i < len(equipped); i++ {
		s.equipSlots[i] = equipped[i]
	}
}

// handleBackspace は合成中のBackspace処理を行います。
func (s *AgentManagementScreen) handleBackspace() (tea.Model, tea.Cmd) {
	if s.currentTab != TabSynthesis {
		return s, nil
	}

	switch s.synthesisState.step {
	case 0: // コア選択中 - 何もしない
		return s, nil
	case 1: // モジュール選択中
		if len(s.synthesisState.selectedModules) > 0 {
			// 最後に選択したモジュールを取り消し
			s.synthesisState.selectedModules = s.synthesisState.selectedModules[:len(s.synthesisState.selectedModules)-1]
		} else {
			// モジュール未選択ならコア選択に戻る
			s.synthesisState.selectedCore = nil
			s.synthesisState.step = 0
			s.selectedIndex = 0
		}
	case 2: // 確認画面
		// モジュール選択に戻る
		s.synthesisState.step = 1
	}
	return s, nil
}

// handleEnter はEnterキーの処理を行います。
func (s *AgentManagementScreen) handleEnter() (tea.Model, tea.Cmd) {
	switch s.currentTab {
	case TabSynthesis:
		return s.handleSynthesisEnter()
	case TabEquip:
		return s.handleEquipEnter()
	}
	return s, nil
}

// handleSynthesisEnter は合成タブでのEnter処理を行います。
func (s *AgentManagementScreen) handleSynthesisEnter() (tea.Model, tea.Cmd) {
	switch s.synthesisState.step {
	case 0: // コア選択
		if s.selectedIndex < len(s.coreList) {
			s.synthesisState.selectedCore = s.coreList[s.selectedIndex]
			s.synthesisState.step = 1
			s.selectedIndex = 0
		}
	case 1: // モジュール選択
		if s.selectedIndex < len(s.moduleList) {
			module := s.moduleList[s.selectedIndex]
			// タグ互換性チェック + 重複チェック
			if s.synthesisState.selectedCore != nil &&
				s.isModuleCompatible(module) &&
				!s.isModuleAlreadySelected(module) {
				if len(s.synthesisState.selectedModules) < 4 {
					s.synthesisState.selectedModules = append(s.synthesisState.selectedModules, module)
				}
				if len(s.synthesisState.selectedModules) == 4 {
					s.synthesisState.step = 2
				}
			}
		}
	case 2: // 確認・合成実行
		if s.canSynthesize() {
			s.executeSynthesis()
			s.resetSynthesisState()
		}
	}
	return s, nil
}

// handleEquipEnter は装備タブでのEnter処理を行います。
func (s *AgentManagementScreen) handleEquipEnter() (tea.Model, tea.Cmd) {
	// 装備スロット（0-2）または所持エージェント（3以降）
	if s.selectedIndex < 3 {
		// 装備解除
		if s.equipSlots[s.selectedIndex] != nil {
			if err := s.inventory.UnequipAgent(s.selectedIndex); err != nil {
				slog.Error("エージェント装備解除に失敗",
					slog.Int("slot", s.selectedIndex),
					slog.Any("error", err),
				)
			}
			s.updateCurrentList()
		}
	} else {
		// エージェント装備
		agentIndex := s.selectedIndex - 3
		if agentIndex < len(s.agentList) {
			// 空きスロットを探す
			for i := 0; i < 3; i++ {
				if s.equipSlots[i] == nil {
					if err := s.inventory.EquipAgent(i, s.agentList[agentIndex]); err != nil {
						slog.Error("エージェント装備に失敗",
							slog.Int("slot", i),
							slog.String("agent_id", s.agentList[agentIndex].ID),
							slog.Any("error", err),
						)
					}
					s.updateCurrentList()
					break
				}
			}
		}
	}
	return s, nil
}

// handleDelete は削除処理を行います。
// UI-Improvement Requirement 2.9: エージェント削除時に確認ダイアログを表示
func (s *AgentManagementScreen) handleDelete() (tea.Model, tea.Cmd) {
	switch s.currentTab {
	case TabCoreList:
		if s.selectedIndex < len(s.coreList) {
			core := s.coreList[s.selectedIndex]
			s.confirmDialog = components.NewConfirmDialog(
				"コアの削除",
				fmt.Sprintf("「%s Lv.%d」を削除しますか？", core.Name, core.Level),
			)
			s.confirmDialog.Show()
			s.pendingDeleteIdx = s.selectedIndex
		}
	case TabModuleList:
		if s.selectedIndex < len(s.moduleList) {
			module := s.moduleList[s.selectedIndex]
			s.confirmDialog = components.NewConfirmDialog(
				"モジュールの削除",
				fmt.Sprintf("「%s」を削除しますか？", module.Name),
			)
			s.confirmDialog.Show()
			s.pendingDeleteIdx = s.selectedIndex
		}
	case TabEquip:
		// エージェント削除（装備タブで所持エージェントを選択中の場合）
		if s.selectedIndex >= 3 && s.selectedIndex-3 < len(s.agentList) {
			agent := s.agentList[s.selectedIndex-3]
			s.confirmDialog = components.NewConfirmDialog(
				"エージェントの削除",
				fmt.Sprintf("「%s Lv.%d」を削除しますか？", agent.GetCoreTypeName(), agent.Level),
			)
			s.confirmDialog.Show()
			s.pendingDeleteIdx = s.selectedIndex
		}
	}
	return s, nil
}

// executeDelete は確認後の削除を実行します。
func (s *AgentManagementScreen) executeDelete() {
	switch s.currentTab {
	case TabCoreList:
		if s.pendingDeleteIdx < len(s.coreList) {
			if err := s.inventory.RemoveCore(s.coreList[s.pendingDeleteIdx].ID); err != nil {
				slog.Error("コア削除に失敗",
					slog.String("core_id", s.coreList[s.pendingDeleteIdx].ID),
					slog.Any("error", err),
				)
			}
			s.updateCurrentList()
		}
	case TabModuleList:
		if s.pendingDeleteIdx < len(s.moduleList) {
			if err := s.inventory.RemoveModule(s.moduleList[s.pendingDeleteIdx].ID); err != nil {
				slog.Error("モジュール削除に失敗",
					slog.String("module_id", s.moduleList[s.pendingDeleteIdx].ID),
					slog.Any("error", err),
				)
			}
			s.updateCurrentList()
		}
	case TabEquip:
		// エージェント削除は未実装（インベントリにRemoveAgentメソッドが必要）
		// 今後の拡張として実装
	}
	s.pendingDeleteIdx = -1
}

// isModuleAlreadySelected は指定されたモジュールが既に選択済みかをチェックします。
func (s *AgentManagementScreen) isModuleAlreadySelected(module *domain.ModuleModel) bool {
	for _, selected := range s.synthesisState.selectedModules {
		if selected.ID == module.ID {
			return true
		}
	}
	return false
}

// isModuleCompatible はモジュールがコアと互換性があるかチェックします。
func (s *AgentManagementScreen) isModuleCompatible(module *domain.ModuleModel) bool {
	if s.synthesisState.selectedCore == nil {
		return false
	}
	allowedTags := s.synthesisState.selectedCore.Type.AllowedTags
	for _, moduleTag := range module.Tags {
		for _, allowedTag := range allowedTags {
			if moduleTag == allowedTag {
				return true
			}
		}
	}
	return false
}

// canSynthesize は合成可能かどうかを返します。
func (s *AgentManagementScreen) canSynthesize() bool {
	return s.synthesisState.selectedCore != nil &&
		len(s.synthesisState.selectedModules) == 4
}

// executeSynthesis は合成を実行します。
func (s *AgentManagementScreen) executeSynthesis() {
	if !s.canSynthesize() {
		return
	}

	// エージェント作成
	agentID := fmt.Sprintf("agent_%d", len(s.agentList)+1)
	agent := domain.NewAgent(agentID, s.synthesisState.selectedCore, s.synthesisState.selectedModules)

	// インベントリに追加
	if err := s.inventory.AddAgent(agent); err != nil {
		slog.Error("エージェント追加に失敗",
			slog.String("agent_id", agent.ID),
			slog.Any("error", err),
		)
		s.errorMessage = fmt.Sprintf("エージェントの合成に失敗しました: %v", err)
		return // 素材を消費しないように処理を中断
	}

	// エラーメッセージをクリアしてステータスメッセージを設定
	s.errorMessage = ""
	s.statusMessage = "エージェントを合成しました"

	// 使用した素材を削除
	if err := s.inventory.RemoveCore(s.synthesisState.selectedCore.ID); err != nil {
		slog.Error("合成素材のコア削除に失敗",
			slog.String("core_id", s.synthesisState.selectedCore.ID),
			slog.Any("error", err),
		)
	}
	for _, m := range s.synthesisState.selectedModules {
		if err := s.inventory.RemoveModule(m.ID); err != nil {
			slog.Error("合成素材のモジュール削除に失敗",
				slog.String("module_id", m.ID),
				slog.Any("error", err),
			)
		}
	}

	s.updateCurrentList()
}

// resetSynthesisState は合成状態をリセットします。
func (s *AgentManagementScreen) resetSynthesisState() {
	s.synthesisState = SynthesisState{
		selectedModules: []*domain.ModuleModel{},
	}
	s.selectedIndex = 0
}

// getSelectedCoreDetail は選択中のコアの詳細を返します。
func (s *AgentManagementScreen) getSelectedCoreDetail() *domain.CoreModel {
	if s.selectedIndex < len(s.coreList) {
		return s.coreList[s.selectedIndex]
	}
	return nil
}

// getSelectedModuleDetail は選択中のモジュールの詳細を返します。
func (s *AgentManagementScreen) getSelectedModuleDetail() *domain.ModuleModel {
	if s.selectedIndex < len(s.moduleList) {
		return s.moduleList[s.selectedIndex]
	}
	return nil
}

// View は画面をレンダリングします。
// UI-Improvement Requirement 2.9: 確認ダイアログのオーバーレイ表示
func (s *AgentManagementScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("エージェント管理"))
	builder.WriteString("\n\n")

	// タブバー
	builder.WriteString(s.renderTabBar())
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
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hints string
	if s.confirmDialog != nil && s.confirmDialog.Visible {
		hints = "←/→: 選択切替  Enter: 決定  Esc: キャンセル"
	} else if s.currentTab == TabEquip {
		hints = "←/→: タブ切替  Tab: スロット切替  ↑/↓: エージェント選択  Enter: 装備  Backspace: 取り外し  Esc: ホーム"
	} else {
		hints = "←/→: タブ切替  ↑/↓: 選択  Enter: 決定  Backspace: 戻る  d: 削除  Esc: ホーム"
	}
	builder.WriteString(hintStyle.Render(hints))

	// 確認ダイアログのオーバーレイ
	if s.confirmDialog != nil && s.confirmDialog.Visible {
		builder.WriteString("\n\n")
		dialog := s.confirmDialog.Render(s.width, s.height)
		builder.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(dialog))
	}

	return builder.String()
}

// renderTabBar はタブバーをレンダリングします。
func (s *AgentManagementScreen) renderTabBar() string {
	tabs := []string{"コア一覧", "モジュール一覧", "合成", "装備"}

	var tabItems []string
	for i, tab := range tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		if AgentManagementTab(i) == s.currentTab {
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
func (s *AgentManagementScreen) renderMainContent() string {
	switch s.currentTab {
	case TabCoreList:
		return s.renderCoreList()
	case TabModuleList:
		return s.renderModuleList()
	case TabSynthesis:
		return s.renderSynthesis()
	case TabEquip:
		return s.renderEquip()
	}
	return ""
}

// renderCoreList はコア一覧をレンダリングします。
func (s *AgentManagementScreen) renderCoreList() string {
	var builder strings.Builder

	if len(s.coreList) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("コアがありません")
	}

	// リストとプレビューを横に並べる
	listContent := s.renderCoreListItems()
	previewContent := s.renderCorePreview()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(50).
		Render(listContent)

	previewBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(50).
		Render(previewContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, listBox, "  ", previewBox)
	builder.WriteString(lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content))

	return builder.String()
}

// renderCoreListItems はコアリストの項目をレンダリングします。
func (s *AgentManagementScreen) renderCoreListItems() string {
	var items []string
	for i, core := range s.coreList {
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}
		item := fmt.Sprintf("%s Lv.%d", core.Type.Name, core.Level)
		items = append(items, style.Render(prefix+item))
	}
	return strings.Join(items, "\n")
}

// renderCorePreview はコアのプレビューをレンダリングします。
func (s *AgentManagementScreen) renderCorePreview() string {
	core := s.getSelectedCoreDetail()
	if core == nil {
		return "コアを選択してください"
	}

	panel := components.NewInfoPanel(core.Name)
	panel.AddItem("レベル", fmt.Sprintf("Lv.%d", core.Level))
	panel.AddItem("特性", core.Type.Name)
	panel.AddItem("STR", fmt.Sprintf("%d", core.Stats.STR))
	panel.AddItem("MAG", fmt.Sprintf("%d", core.Stats.MAG))
	panel.AddItem("SPD", fmt.Sprintf("%d", core.Stats.SPD))
	panel.AddItem("LUK", fmt.Sprintf("%d", core.Stats.LUK))

	return panel.Render(45)
}

// renderModuleList はモジュール一覧をレンダリングします。
func (s *AgentManagementScreen) renderModuleList() string {
	if len(s.moduleList) == 0 {
		return lipgloss.NewStyle().
			Width(s.width).
			Align(lipgloss.Center).
			Foreground(styles.ColorSubtle).
			Render("モジュールがありません")
	}

	listContent := s.renderModuleListItems()
	previewContent := s.renderModulePreview()

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(50).
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

// renderModuleListItems はモジュールリストの項目をレンダリングします。
func (s *AgentManagementScreen) renderModuleListItems() string {
	var items []string
	for i, module := range s.moduleList {
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}
		item := fmt.Sprintf("%s [%s] Lv.%d", module.Name, module.Category.String(), module.Level)
		items = append(items, style.Render(prefix+item))
	}
	return strings.Join(items, "\n")
}

// renderModulePreview はモジュールのプレビューをレンダリングします。
func (s *AgentManagementScreen) renderModulePreview() string {
	module := s.getSelectedModuleDetail()
	if module == nil {
		return "モジュールを選択してください"
	}

	panel := components.NewInfoPanel(module.Name)
	panel.AddItem("カテゴリ", module.Category.String())
	panel.AddItem("レベル", fmt.Sprintf("Lv.%d", module.Level))
	panel.AddItem("基礎効果", fmt.Sprintf("%.0f", module.BaseEffect))
	panel.AddItem("参照ステータス", module.StatRef)
	panel.AddItem("説明", module.Description)

	return panel.Render(45)
}

// renderSynthesis は合成画面をレンダリングします。
// design.md準拠: 左右2カラムレイアウト
func (s *AgentManagementScreen) renderSynthesis() string {
	// 左側パネル: 選択可能なパーツ
	leftContent := s.renderSynthesisLeftPanel()

	// 右側パネル: 選択済みパーツ + 完成予測ステータス
	rightContent := s.renderSynthesisRightPanel()

	// 左右のパネルを結合
	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(50).
		Height(20).
		Render(leftContent)

	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(55).
		Height(20).
		Render(rightContent)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftBox, "  ", rightBox)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderSynthesisLeftPanel は合成画面の左側パネルをレンダリングします。
func (s *AgentManagementScreen) renderSynthesisLeftPanel() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("選択可能なパーツ"))
	builder.WriteString("\n\n")

	// ステップ表示
	steps := []string{"コア選択", "モジュール選択", "確認"}
	stepStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(stepStyle.Render(fmt.Sprintf("ステップ %d: %s", s.synthesisState.step+1, steps[s.synthesisState.step])))
	builder.WriteString("\n\n")

	// リスト
	switch s.synthesisState.step {
	case 0:
		builder.WriteString(s.renderSynthesisCoreListItems())
	case 1:
		builder.WriteString(s.renderSynthesisModuleListItems())
	case 2:
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Enterキーで合成を実行"))
	}

	builder.WriteString("\n\n")

	// 区切り線
	builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("────────────────────────────────────────"))
	builder.WriteString("\n")

	// 選択中パーツの詳細
	builder.WriteString(titleStyle.Render("選択中パーツの詳細"))
	builder.WriteString("\n\n")
	builder.WriteString(s.renderSelectedPartDetail())

	return builder.String()
}

// renderSynthesisRightPanel は合成画面の右側パネルをレンダリングします。
func (s *AgentManagementScreen) renderSynthesisRightPanel() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("選択済みパーツ"))
	builder.WriteString("\n\n")

	// コア
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary)

	builder.WriteString(labelStyle.Render("コア: "))
	if s.synthesisState.selectedCore != nil {
		builder.WriteString(valueStyle.Render(fmt.Sprintf("%s Lv.%d", s.synthesisState.selectedCore.Name, s.synthesisState.selectedCore.Level)))
	} else {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("(未選択)"))
	}
	builder.WriteString("\n")

	// モジュール
	for i := 0; i < 4; i++ {
		builder.WriteString(labelStyle.Render(fmt.Sprintf("モジュール%d: ", i+1)))
		if i < len(s.synthesisState.selectedModules) {
			builder.WriteString(valueStyle.Render(s.synthesisState.selectedModules[i].Name))
		} else {
			builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("(未選択)"))
		}
		builder.WriteString("\n")
	}

	builder.WriteString("\n")

	// 区切り線
	builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("──────────────────────────────────────────────"))
	builder.WriteString("\n")

	// 完成予測ステータス
	builder.WriteString(titleStyle.Render("完成予測ステータス"))
	builder.WriteString("\n\n")
	builder.WriteString(s.renderSynthesisPreview())

	return builder.String()
}

// renderSelectedPartDetail は選択中パーツの詳細をレンダリングします。
func (s *AgentManagementScreen) renderSelectedPartDetail() string {
	switch s.synthesisState.step {
	case 0:
		// コア選択中
		core := s.getSelectedCoreDetail()
		if core == nil {
			return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("コアを選択してください")
		}
		return s.formatCoreDetail(core)
	case 1:
		// モジュール選択中
		module := s.getSelectedModuleDetail()
		if module == nil {
			return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("モジュールを選択してください")
		}
		return s.formatModuleDetail(module)
	case 2:
		// 確認画面
		return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("合成準備完了")
	}
	return ""
}

// formatCoreDetail はコアの詳細をフォーマットします。
func (s *AgentManagementScreen) formatCoreDetail(core *domain.CoreModel) string {
	var builder strings.Builder
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)

	builder.WriteString(nameStyle.Render(fmt.Sprintf("%s Lv.%d", core.Name, core.Level)))
	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("特性: "))
	builder.WriteString(core.Type.Name)
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("STR: %-3d  MAG: %-3d  SPD: %-3d  LUK: %-3d",
		core.Stats.STR, core.Stats.MAG, core.Stats.SPD, core.Stats.LUK))

	return builder.String()
}

// formatModuleDetail はモジュールの詳細をフォーマットします。
func (s *AgentManagementScreen) formatModuleDetail(module *domain.ModuleModel) string {
	var builder strings.Builder
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)

	builder.WriteString(nameStyle.Render(module.Name))
	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("カテゴリ: "))
	builder.WriteString(s.getModuleIcon(module.Category) + " " + module.Category.String())
	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("基礎効果: "))
	builder.WriteString(fmt.Sprintf("%.0f", module.BaseEffect))
	builder.WriteString("\n")
	builder.WriteString(labelStyle.Render("説明: "))
	builder.WriteString(module.Description)

	return builder.String()
}

// renderSynthesisPreview は合成後の予測ステータスをレンダリングします。
func (s *AgentManagementScreen) renderSynthesisPreview() string {
	if s.synthesisState.selectedCore == nil {
		return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("コアを選択すると\n予測ステータスが表示されます")
	}

	var builder strings.Builder
	core := s.synthesisState.selectedCore
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	// 名前
	builder.WriteString(labelStyle.Render("名前: "))
	builder.WriteString(valueStyle.Render(core.Type.Name))
	builder.WriteString("\n")

	// レベル
	builder.WriteString(labelStyle.Render("Lv: "))
	builder.WriteString(valueStyle.Render(fmt.Sprintf("%d", core.Level)))
	builder.WriteString("\n")

	// ステータス
	stats := core.Stats
	builder.WriteString(fmt.Sprintf("STR: %-3d  MAG: %-3d\n", stats.STR, stats.MAG))
	builder.WriteString(fmt.Sprintf("SPD: %-3d  LUK: %-3d\n", stats.SPD, stats.LUK))

	// モジュール情報
	if len(s.synthesisState.selectedModules) > 0 {
		builder.WriteString("\n")
		builder.WriteString(labelStyle.Render("攻撃: "))
		var attacks []string
		for _, m := range s.synthesisState.selectedModules {
			icon := s.getModuleIcon(m.Category)
			attacks = append(attacks, fmt.Sprintf("%s(%s+%.0f)", m.Name, icon, m.BaseEffect))
		}
		builder.WriteString(strings.Join(attacks, ", "))
	}

	return builder.String()
}

// renderSynthesisCoreListItems はコアリストの項目をレンダリングします（合成用）。
func (s *AgentManagementScreen) renderSynthesisCoreListItems() string {
	if len(s.coreList) == 0 {
		return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("コアがありません")
	}

	var items []string
	for i, core := range s.coreList {
		style := lipgloss.NewStyle()
		prefix := "  "
		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}
		item := fmt.Sprintf("%s Lv.%d (%s)", core.Name, core.Level, core.Type.Name)
		items = append(items, style.Render(prefix+item))
	}
	return strings.Join(items, "\n")
}

// renderSynthesisModuleListItems はモジュールリストの項目をレンダリングします（合成用）。
func (s *AgentManagementScreen) renderSynthesisModuleListItems() string {
	if len(s.moduleList) == 0 {
		return lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("モジュールがありません")
	}

	var items []string
	for i, module := range s.moduleList {
		compatible := s.isModuleCompatible(module)
		alreadySelected := s.isModuleAlreadySelected(module)
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		} else if !compatible || alreadySelected {
			style = style.Foreground(styles.ColorSubtle)
		}

		icon := s.getModuleIcon(module.Category)
		item := fmt.Sprintf("%s %s", icon, module.Name)
		if !compatible {
			item += " (互換性なし)"
		} else if alreadySelected {
			item += " (選択済み)"
		}
		items = append(items, style.Render(prefix+item))
	}
	return strings.Join(items, "\n")
}

// renderEquip は装備画面をレンダリングします。
// UI-Improvement Requirement 2.5, 2.6, 2.7: 上下分割レイアウト
func (s *AgentManagementScreen) renderEquip() string {
	var builder strings.Builder

	// 上部エリア: 所持エージェント一覧（左）と選択中エージェント詳細（右）
	topSection := s.renderEquipTopSection()
	topBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1, 2).
		Width(s.width - 10)

	topTitle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Render("──────────────────────────  所持エージェント  ──────────────────────────")

	builder.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(
		lipgloss.JoinVertical(lipgloss.Center, topTitle, topBox.Render(topSection)),
	))
	builder.WriteString("\n\n")

	// 下部エリア: 装備中エージェント（3スロット横並び）
	bottomSection := s.renderEquipBottomSection()
	bottomBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Width(s.width - 10)

	bottomTitle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Render("────────────────────────  装備中エージェント  ────────────────────────")

	builder.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(
		lipgloss.JoinVertical(lipgloss.Center, bottomTitle, bottomBox.Render(bottomSection)),
	))

	return builder.String()
}

// renderEquipTopSection は装備タブの上部セクション（エージェント一覧と詳細）をレンダリングします。
func (s *AgentManagementScreen) renderEquipTopSection() string {
	// 左側: エージェント一覧
	leftContent := s.renderEquipAgentList()
	leftBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(0, 1).
		Width(35).
		Height(12)

	// 右側: 選択中エージェントの詳細
	rightContent := s.renderEquipAgentDetail()
	rightBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(0, 1).
		Width(55).
		Height(12)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftBox.Render(leftContent),
		"  ",
		rightBox.Render(rightContent),
	)
}

// renderEquipAgentList は装備タブのエージェント一覧をレンダリングします。
func (s *AgentManagementScreen) renderEquipAgentList() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("エージェント一覧"))
	builder.WriteString("\n\n")

	if len(s.agentList) == 0 {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("(なし)"))
		return builder.String()
	}

	for i, agent := range s.agentList {
		style := lipgloss.NewStyle()
		prefix := "  "

		if i == s.selectedIndex {
			style = style.Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
			prefix = "> "
		}

		// 装備中かチェックしてマークを付ける
		equipMark := ""
		for slotIdx, equipped := range s.equipSlots {
			if equipped != nil && equipped.ID == agent.ID {
				equipMark = fmt.Sprintf(" [E%d]", slotIdx+1)
				break
			}
		}

		item := fmt.Sprintf("%s Lv.%d%s", agent.GetCoreTypeName(), agent.Level, equipMark)
		builder.WriteString(style.Render(prefix + item))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderEquipAgentDetail は選択中エージェントの詳細をレンダリングします。
func (s *AgentManagementScreen) renderEquipAgentDetail() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	builder.WriteString(titleStyle.Render("選択中エージェント詳細"))
	builder.WriteString("\n\n")

	// 選択中のエージェントを取得（一覧から）
	var selectedAgent *domain.AgentModel
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.agentList) {
		selectedAgent = s.agentList[s.selectedIndex]
	}

	if selectedAgent == nil {
		builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("エージェントを選択してください"))
		return builder.String()
	}

	// エージェント名とレベル
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	builder.WriteString(nameStyle.Render(fmt.Sprintf("%s Lv.%d", selectedAgent.GetCoreTypeName(), selectedAgent.Level)))
	builder.WriteString("\n")

	// コアタイプ
	labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
	builder.WriteString(labelStyle.Render("コアタイプ: "))
	builder.WriteString(valueStyle.Render(selectedAgent.Core.Type.Name))
	builder.WriteString("\n")

	// 区切り線
	builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("────────────────────────────────────"))
	builder.WriteString("\n")

	// ステータス（コアから取得）
	stats := selectedAgent.Core.Stats
	builder.WriteString(fmt.Sprintf("STR: %-4d  MAG: %-4d  SPD: %-4d  LUK: %-4d\n",
		stats.STR, stats.MAG, stats.SPD, stats.LUK))

	// 区切り線
	builder.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("────────────────────────────────────"))
	builder.WriteString("\n")

	// モジュール
	builder.WriteString(labelStyle.Render("モジュール:"))
	builder.WriteString("\n")
	for _, module := range selectedAgent.Modules {
		icon := s.getModuleIcon(module.Category)
		builder.WriteString(fmt.Sprintf("  %s %s\n", icon, module.Name))
	}

	return builder.String()
}

// renderEquipBottomSection は装備タブの下部セクション（装備スロット）をレンダリングします。
func (s *AgentManagementScreen) renderEquipBottomSection() string {
	var cards []string
	cardWidth := 30

	for i := 0; i < 3; i++ {
		var cardContent strings.Builder
		isSelected := i == s.selectedEquipSlot

		// スロットタイトル
		slotTitle := fmt.Sprintf("スロット%d", i+1)
		if isSelected {
			slotTitle = "*" + slotTitle + "*"
		}
		cardContent.WriteString(lipgloss.NewStyle().Bold(true).Render(slotTitle))
		cardContent.WriteString("\n\n")

		if s.equipSlots[i] != nil {
			agent := s.equipSlots[i]
			cardContent.WriteString(fmt.Sprintf("%s Lv.%d\n", agent.GetCoreTypeName(), agent.Level))

			// モジュールアイコン
			var icons []string
			for _, module := range agent.Modules {
				icons = append(icons, s.getModuleIcon(module.Category))
			}
			cardContent.WriteString(strings.Join(icons, ""))
		} else {
			cardContent.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("(空)\n\n"))
			cardContent.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("Enterで装備"))
		}

		// カードスタイル
		borderColor := styles.ColorSubtle
		if isSelected {
			borderColor = styles.ColorPrimary
		}

		cardStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(cardWidth).
			Height(6)

		cards = append(cards, cardStyle.Render(cardContent.String()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, cards[0], "  ", cards[1], "  ", cards[2])
}

// getModuleIcon はモジュールカテゴリのアイコンを返します。
func (s *AgentManagementScreen) getModuleIcon(category domain.ModuleCategory) string {
	switch category {
	case domain.PhysicalAttack:
		return "⚔"
	case domain.MagicAttack:
		return "✦"
	case domain.Heal:
		return "♥"
	case domain.Buff:
		return "▲"
	case domain.Debuff:
		return "▼"
	default:
		return "•"
	}
}
