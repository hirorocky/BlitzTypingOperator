// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strconv"
	"strings"

	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.2: バトル選択画面 ====================

// BattleSelectState はバトル選択画面の状態を表します。
type BattleSelectState int

const (
	// StateInput は入力状態です。
	StateInput BattleSelectState = iota
	// StateConfirm は確認状態です。
	StateConfirm
)

// StartBattleMsg はバトル開始を要求するメッセージです。
type StartBattleMsg struct {
	Level int
}

// BattleSelectScreen はバトル選択画面を表します。

type BattleSelectScreen struct {
	input             *components.InputField
	maxLevelReached   int
	maxChallengeLevel int
	agentProvider     AgentProvider // 装備エージェントを取得するプロバイダー
	state             BattleSelectState
	selectedLevel     int
	error             string
	styles            *styles.GameStyles
	width             int
	height            int
}

// NewBattleSelectScreen は新しいBattleSelectScreenを作成します。

func NewBattleSelectScreen(maxLevelReached int, agentProvider AgentProvider) *BattleSelectScreen {
	input := components.NewInputField("レベル番号を入力 (例: 1)")
	input.InputMode = components.InputModeNumeric
	input.MinValue = 1
	input.MaxValue = maxLevelReached + 1
	input.MaxLength = 3

	return &BattleSelectScreen{
		input:             input,
		maxLevelReached:   maxLevelReached,
		maxChallengeLevel: maxLevelReached + 1,
		agentProvider:     agentProvider,
		state:             StateInput,
		styles:            styles.NewGameStyles(),
		width:             140,
		height:            40,
	}
}

// Init は画面の初期化を行います。
func (s *BattleSelectScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *BattleSelectScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *BattleSelectScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch s.state {
	case StateInput:
		return s.handleInputState(msg)
	case StateConfirm:
		return s.handleConfirmState(msg)
	}
	return s, nil
}

// handleInputState は入力状態でのキー処理を行います。
func (s *BattleSelectScreen) handleInputState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":

		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "enter":
		// 入力を検証
		valid, errMsg := s.validateInput()
		if !valid {
			s.error = errMsg
			return s, nil
		}
		// 確認画面へ遷移
		level, _ := strconv.Atoi(s.input.Value)
		s.selectedLevel = level
		s.state = StateConfirm
		s.error = ""
		return s, nil
	case "backspace":
		s.input.HandleBackspace()
		s.error = ""
	default:
		if len(msg.Runes) == 1 {
			s.input.HandleInput(msg.Runes[0])
			s.error = ""
		}
	}
	return s, nil
}

// handleConfirmState は確認状態でのキー処理を行います。
func (s *BattleSelectScreen) handleConfirmState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		// 入力画面に戻る
		s.state = StateInput
		return s, nil
	case "enter", "y":

		equippedAgents := s.agentProvider.GetEquippedAgents()
		if len(equippedAgents) == 0 {
			s.error = "エージェントが装備されていません。\nエージェント管理でエージェントを装備してください。"
			return s, nil
		}

		return s, func() tea.Msg {
			return StartBattleMsg{Level: s.selectedLevel}
		}
	}
	return s, nil
}

// validateInput は入力を検証します。

func (s *BattleSelectScreen) validateInput() (bool, string) {
	if s.input.Value == "" {
		return false, "レベル番号を入力してください"
	}

	level, err := strconv.Atoi(s.input.Value)
	if err != nil {
		return false, "有効な数値を入力してください"
	}

	if level < 1 {
		return false, "レベルは1以上を入力してください"
	}

	if level > s.maxChallengeLevel {
		return false, fmt.Sprintf("挑戦可能な最大レベルはLv.%dです", s.maxChallengeLevel)
	}

	return true, ""
}

// View は画面をレンダリングします。
func (s *BattleSelectScreen) View() string {
	switch s.state {
	case StateInput:
		return s.renderInputState()
	case StateConfirm:
		return s.renderConfirmState()
	}
	return ""
}

// renderInputState は入力状態の画面をレンダリングします。
func (s *BattleSelectScreen) renderInputState() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("バトル選択"))
	builder.WriteString("\n\n")

	// レベル情報

	infoStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	levelInfo := fmt.Sprintf("到達最高レベル: Lv.%d  |  挑戦可能: Lv.1 〜 Lv.%d",
		s.maxLevelReached, s.maxChallengeLevel)
	builder.WriteString(infoStyle.Render(levelInfo))
	builder.WriteString("\n\n")

	// 入力フィールド
	inputBox := s.input.Render(30)
	centeredInput := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(inputBox)
	builder.WriteString(centeredInput)
	builder.WriteString("\n\n")

	// エラーメッセージ
	if s.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorDamage).
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(errorStyle.Render(s.error))
		builder.WriteString("\n\n")
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(hintStyle.Render("Enter: 確認  Esc: 戻る"))

	return builder.String()
}

// renderConfirmState は確認状態の画面をレンダリングします。

func (s *BattleSelectScreen) renderConfirmState() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("バトル確認"))
	builder.WriteString("\n\n")

	// 確認内容
	contentStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(s.width)

	confirmContent := fmt.Sprintf("Lv.%d に挑戦しますか？", s.selectedLevel)
	builder.WriteString(contentStyle.Render(confirmContent))
	builder.WriteString("\n\n")

	// 予想敵情報
	infoPanel := components.NewInfoPanel("予想敵情報")
	infoPanel.AddItem("レベル", fmt.Sprintf("Lv.%d", s.selectedLevel))
	infoPanel.AddItem("予想HP", fmt.Sprintf("約 %d", s.selectedLevel*100))
	infoPanel.AddItem("予想攻撃力", fmt.Sprintf("約 %d", 10+s.selectedLevel*2))

	infoPanelRendered := infoPanel.Render(40)
	centeredInfo := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(infoPanelRendered)
	builder.WriteString(centeredInfo)
	builder.WriteString("\n\n")

	// 装備中エージェント情報（AgentProviderから最新の状態を取得）
	equippedAgents := s.agentProvider.GetEquippedAgents()
	agentPanel := components.NewInfoPanel("装備中エージェント")
	if len(equippedAgents) == 0 {
		agentPanel.AddItem("状態", "未装備")
	} else {
		for i, agent := range equippedAgents {
			agentPanel.AddItem(fmt.Sprintf("スロット%d", i+1),
				fmt.Sprintf("%s (Lv.%d)", agent.GetCoreTypeName(), agent.Level))
		}
	}

	agentPanelRendered := agentPanel.Render(40)
	centeredAgent := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(agentPanelRendered)
	builder.WriteString(centeredAgent)
	builder.WriteString("\n\n")

	// エラーメッセージ
	if s.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorDamage).
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(errorStyle.Render(s.error))
		builder.WriteString("\n\n")
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(hintStyle.Render("Enter/y: バトル開始  Esc/n: 戻る"))

	return builder.String()
}

// SetMaxLevelReached は到達最高レベルを設定します。
func (s *BattleSelectScreen) SetMaxLevelReached(level int) {
	s.maxLevelReached = level
	s.maxChallengeLevel = level + 1
	s.input.MaxValue = s.maxChallengeLevel
}
