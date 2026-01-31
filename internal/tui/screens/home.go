// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/ascii"
	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.1: ホーム画面 ====================

// AgentProvider は装備エージェントを提供するインターフェースです。
// HomeScreenやBattleSelectScreenがAgentManagerから最新の装備状態を取得するために使用します。
type AgentProvider interface {
	GetEquippedAgents() []*domain.AgentModel
}

// SlotReadinessProvider はエージェントスロットの準備状態を提供するインターフェースです。
// v3.0.0 AgentSlotManager で実装されます。
type SlotReadinessProvider interface {
	// GetReadySlotCount はバトルに使用可能なスロット数を返します。
	GetReadySlotCount() int
}

// DefeatedEnemyProvider は撃破済み敵情報を提供するインターフェースです。
// BattleSelectScreenがレベル選択の制約を判断するために使用します。
type DefeatedEnemyProvider interface {
	// GetDefeatedEnemies は敵タイプIDをキー、撃破最高レベルを値とするマップを返します。
	GetDefeatedEnemies() map[string]int
	// IsEnemyDefeated は指定した敵タイプが撃破済みかどうかを返します。
	IsEnemyDefeated(enemyTypeID string) bool
	// GetDefeatedLevel は指定した敵タイプの撃破最高レベルを返します。
	GetDefeatedLevel(enemyTypeID string) int
	// GetMaxDefeatedLevel は全敵種類を通じた最高撃破レベル（到達Lv）を返します。
	GetMaxDefeatedLevel() int
	// GetMaxLevelReached は到達最高レベル（敵のデフォルトレベルで更新）を返します。
	GetMaxLevelReached() int
}

// EnemyProgressProvider はランクベースの敵進行情報を提供するインターフェースです。
// BattleSelectScreenRankBasedが敵表示とレベル選択の制約を判断するために使用します。
type EnemyProgressProvider interface {
	// GetCurrentRank は現在解放済みランクを返します。
	GetCurrentRank() int
	// GetCurrentRankEnemies は現在ランクの敵リストを返します。
	GetCurrentRankEnemies() []domain.EnemyType
	// IsDefeated は指定した敵が撃破済みかどうかを返します。
	IsDefeated(enemyTypeID string) bool
	// GetMaxDefeatedLevel は指定した敵の撃破済み最大レベルを返します。
	GetMaxDefeatedLevel(enemyTypeID string) int
	// GetSelectableLevelRange は敵の選択可能レベル範囲を取得します。
	GetSelectableLevelRange(enemyTypeID string) (min, max int)
}

// HomeScreen はホーム画面を表します。

// UI-Improvement Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6
type HomeScreen struct {
	menu            *components.Menu
	maxLevelReached int
	agentProvider   AgentProvider         // 装備エージェントを取得するプロバイダー
	slotProvider    SlotReadinessProvider // v3.0.0 スロット準備状態プロバイダー
	styles          *styles.GameStyles
	width           int
	height          int
	statusMessage   string // セーブ/ロード結果などのステータスメッセージ
	// UI改善: ASCIIアートレンダラー
	logoRenderer   ascii.ASCIILogoRenderer
	numberRenderer ascii.ASCIINumberRenderer
	// セーブ確認ダイアログ
	confirmDialog  *components.ConfirmDialog
	showingConfirm bool
}

// ChangeSceneMsg はシーン遷移を要求するメッセージです。
type ChangeSceneMsg struct {
	Scene string
}

// SaveRequestMsg はセーブ実行を要求するメッセージです。
type SaveRequestMsg struct{}

// NewHomeScreen は新しいHomeScreenを作成します。

// UI-Improvement Requirement 1.6, 5.3: 装備エージェントが空の場合、バトル選択メニューを無効化
func NewHomeScreen(maxLevelReached int, agentProvider AgentProvider) *HomeScreen {
	// 装備エージェントがあるかチェック
	hasEquippedAgents := false
	if agentProvider != nil {
		equippedAgents := agentProvider.GetEquippedAgents()
		hasEquippedAgents = len(equippedAgents) > 0
	}

	// UI-Improvement Requirement 1.6, 5.3: 装備がない場合はバトル選択を無効化
	items := []components.MenuItem{
		{Label: "エージェントカスタマイズ", Value: "agent_customization"},
		{Label: "インベントリ", Value: "inventory"},
		{Label: "バトル選択", Value: "battle_select", Disabled: !hasEquippedAgents},
		{Label: "図鑑", Value: "encyclopedia"},
		{Label: "統計/実績", Value: "stats_achievements"},
		{Label: "セーブ", Value: "save"},
		{Label: "設定", Value: "settings"},
	}

	return &HomeScreen{
		menu:            components.NewMenuWithTitle("メインメニュー", items),
		maxLevelReached: maxLevelReached,
		agentProvider:   agentProvider,
		styles:          styles.NewGameStyles(),
		width:           140,
		height:          40,
		// UI改善: ASCIIアートレンダラーを初期化
		logoRenderer:   ascii.NewASCIILogo(),
		numberRenderer: ascii.NewASCIINumbers(),
		// セーブ確認ダイアログを初期化
		confirmDialog: components.NewConfirmDialog(
			"セーブ確認",
			"ゲームの進行状況をセーブしますか？",
		),
	}
}

// Init は画面の初期化を行います。
func (s *HomeScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *HomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s *HomeScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// セーブ確認ダイアログ表示中
	if s.showingConfirm {
		result := s.confirmDialog.HandleKey(msg.String())
		switch result {
		case components.ConfirmResultYes:
			s.showingConfirm = false
			return s, func() tea.Msg { return SaveRequestMsg{} }
		case components.ConfirmResultNo, components.ConfirmResultCancelled:
			s.showingConfirm = false
		}
		return s, nil
	}

	// 通常のメニュー操作
	switch msg.String() {
	case "up", "k":
		s.menu.MoveUp()
	case "down", "j":
		s.menu.MoveDown()
	case "enter":
		selected := s.menu.GetSelected()
		return s, s.handleMenuSelection(selected.Value)
	case "q", "ctrl+c":
		return s, tea.Quit
	}

	return s, nil
}

// handleMenuSelection はメニュー選択を処理します。

func (s *HomeScreen) handleMenuSelection(value string) tea.Cmd {
	// セーブ選択時は確認ダイアログを表示
	if value == "save" {
		s.showingConfirm = true
		s.confirmDialog.Show()
		return nil
	}

	return func() tea.Msg {
		return ChangeSceneMsg{Scene: value}
	}
}

// View は画面をレンダリングします。
// UI-Improvement Requirement 1.1: ASCIIアートロゴを表示
func (s *HomeScreen) View() string {
	var builder strings.Builder

	// UI改善: ASCIIアートロゴを表示

	logo := s.logoRenderer.Render(true) // カラーモード
	logoLines := strings.Split(logo, "\n")

	// ロゴを中央揃えで表示
	for _, line := range logoLines {
		lineWidth := len([]rune(line))
		padding := (s.width - lineWidth) / 2
		if padding < 0 {
			padding = 0
		}
		builder.WriteString(strings.Repeat(" ", padding))
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	// サブタイトル
	subtitleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	subtitle := subtitleStyle.Render("Terminal Typing Battle Game")
	builder.WriteString(subtitle)
	builder.WriteString("\n\n")

	// メインコンテンツ（メニューと進行状況を横並び）
	menuContent := s.menu.Render()
	statusContent := s.renderStatusPanel()

	// レイアウト調整
	menuBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Width(40).
		Render(menuContent)

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1, 2).
		Width(50).
		Render(statusContent)

	// 横に並べる
	content := lipgloss.JoinHorizontal(lipgloss.Top, menuBox, "  ", statusBox)

	// 中央揃え
	centeredContent := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)

	builder.WriteString(centeredContent)
	builder.WriteString("\n\n")

	// ステータスメッセージ（セーブ/ロード結果など）
	if s.statusMessage != "" {
		statusStyle := lipgloss.NewStyle().
			Foreground(styles.ColorHeal).
			Align(lipgloss.Center).
			Width(s.width)

		status := statusStyle.Render("💾 " + s.statusMessage)
		builder.WriteString(status)
		builder.WriteString("\n\n")
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	hint := hintStyle.Render("↑/k: 上  ↓/j: 下  Enter: 選択  q: 終了")
	builder.WriteString(hint)

	// セーブ確認ダイアログのオーバーレイ
	if s.showingConfirm {
		baseView := builder.String()
		dialog := s.confirmDialog.Render(s.width, s.height)
		return baseView + "\n\n" + dialog
	}

	return builder.String()
}

// renderStatusPanel は進行状況パネルをレンダリングします。

// UI-Improvement Requirement 1.4: 到達レベルをASCII数字アートで表示
func (s *HomeScreen) renderStatusPanel() string {
	var builder strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary)

	labelStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle)

	valueStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSecondary).
		Bold(true)

	builder.WriteString(titleStyle.Render("進行状況"))
	builder.WriteString("\n\n")

	// UI改善: 到達最高レベルをASCII数字アートで表示

	builder.WriteString(labelStyle.Render("到達最高レベル:"))
	builder.WriteString("\n")
	if s.maxLevelReached == 0 {
		builder.WriteString(labelStyle.Render("  まだなし"))
	} else {
		// ASCII数字でレベルを表示
		levelArt := s.numberRenderer.RenderNumber(s.maxLevelReached, styles.ColorPrimary)
		builder.WriteString(levelArt)
	}
	builder.WriteString("\n")

	// 挑戦可能最大レベル
	builder.WriteString(labelStyle.Render("挑戦可能レベル: "))
	nextLevel := s.maxLevelReached + 1
	builder.WriteString(valueStyle.Render(fmt.Sprintf("Lv.%d まで", nextLevel)))
	builder.WriteString("\n\n")

	// 装備中エージェント
	builder.WriteString(titleStyle.Render("装備中エージェント"))
	builder.WriteString("\n\n")

	// AgentProviderから最新の装備状態を取得
	var equippedAgents []*domain.AgentModel
	if s.agentProvider != nil {
		equippedAgents = s.agentProvider.GetEquippedAgents()
	}

	// UI-Improvement Requirement 1.6: 装備なし時の誘導メッセージ
	if len(equippedAgents) == 0 {
		// 警告アイコンとメッセージ
		warningStyle := lipgloss.NewStyle().
			Foreground(styles.ColorWarning).
			Bold(true)
		builder.WriteString(warningStyle.Render("⚠ エージェントが装備されていません"))
		builder.WriteString("\n\n")

		// 誘導メッセージ
		guideStyle := lipgloss.NewStyle().
			Foreground(styles.ColorInfo)
		builder.WriteString(guideStyle.Render("→ エージェント管理で"))
		builder.WriteString("\n")
		builder.WriteString(guideStyle.Render("  合成・装備を行ってください"))
	} else {
		for i, agent := range equippedAgents {
			slotLabel := fmt.Sprintf("スロット%d: ", i+1)
			builder.WriteString(labelStyle.Render(slotLabel))
			agentInfo := agent.GetCoreTypeName()
			builder.WriteString(valueStyle.Render(agentInfo))
			builder.WriteString("\n")
		}

		// 空きスロットを表示
		for i := len(equippedAgents); i < 3; i++ {
			slotLabel := fmt.Sprintf("スロット%d: ", i+1)
			builder.WriteString(labelStyle.Render(slotLabel))
			builder.WriteString(labelStyle.Render("(空)"))
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// SetMaxLevelReached は到達最高レベルを設定します。
func (s *HomeScreen) SetMaxLevelReached(level int) {
	s.maxLevelReached = level
}

// SetSlotProvider はスロット準備状態プロバイダーを設定します。
// v3.0.0 AgentSlotManager を使用してバトル選択の有効/無効を判断します。
func (s *HomeScreen) SetSlotProvider(provider SlotReadinessProvider) {
	s.slotProvider = provider
}

// RefreshMenuState はメニューの有効/無効状態を最新の装備状態に基づいて更新します。
func (s *HomeScreen) RefreshMenuState() {
	hasEquippedAgents := false

	// v3.0.0: slotProviderがある場合は優先使用（AgentSlotManager経由）
	if s.slotProvider != nil {
		hasEquippedAgents = s.slotProvider.GetReadySlotCount() > 0
	} else if s.agentProvider != nil {
		// 後方互換: 旧来のagentProvider経由
		equippedAgents := s.agentProvider.GetEquippedAgents()
		hasEquippedAgents = len(equippedAgents) > 0
	}

	// "バトル選択"メニュー項目の無効状態を更新
	for i := range s.menu.Items {
		if s.menu.Items[i].Value == "battle_select" {
			s.menu.Items[i].Disabled = !hasEquippedAgents
			break
		}
	}
}

// SetStatusMessage はステータスメッセージを設定します。
func (s *HomeScreen) SetStatusMessage(msg string) {
	s.statusMessage = msg
}

// ClearStatusMessage はステータスメッセージをクリアします。
func (s *HomeScreen) ClearStatusMessage() {
	s.statusMessage = ""
}

// ==================== Screenインターフェース実装 ====================

// SetSize は画面サイズを設定します。
// Screenインターフェースの実装です。
func (s *HomeScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}

// GetTitle は画面のタイトルを返します。
// Screenインターフェースの実装です。
func (s *HomeScreen) GetTitle() string {
	return "ホーム"
}

// GetSize は現在の画面サイズを返します。
func (s *HomeScreen) GetSize() (width, height int) {
	return s.width, s.height
}
