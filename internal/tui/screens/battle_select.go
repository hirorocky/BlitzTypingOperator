// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strconv"
	"strings"

	"hirorocky/type-battle/internal/domain"
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
	Level       int
	EnemyTypeID string // カルーセル方式で選択した敵タイプID（空の場合はランダム）
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
			agentPanel.AddItem(fmt.Sprintf("スロット%d", i+1), agent.GetCoreTypeName())
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

// ==================== カルーセル方式のバトル選択画面 ====================

// EnemyTypeProvider は敵タイプリストを提供するインターフェースです。
type EnemyTypeProvider interface {
	GetEnemyTypes() []domain.EnemyType
}

// BattleSelectScreenCarousel はカルーセル方式のバトル選択画面を表します。
type BattleSelectScreenCarousel struct {
	agentProvider    AgentProvider
	defeatedProvider DefeatedEnemyProvider
	enemyTypes       []domain.EnemyType

	// 敵種類選択用
	selectedTypeIdx int

	// レベル選択用
	selectedLevel      int
	minSelectableLevel int // 敵タイプのデフォルトレベル
	maxSelectableLevel int // 撃破済み最高レベル+1（未撃破ならデフォルトレベル）

	error  string
	styles *styles.GameStyles
	width  int
	height int
}

// NewBattleSelectScreenCarousel は新しいカルーセル方式のBattleSelectScreenを作成します。
func NewBattleSelectScreenCarousel(
	agentProvider AgentProvider,
	defeatedProvider DefeatedEnemyProvider,
	enemyTypeProvider EnemyTypeProvider,
) *BattleSelectScreenCarousel {
	allEnemyTypes := enemyTypeProvider.GetEnemyTypes()

	// 到達最高レベル（敵のデフォルトレベルで更新）を取得
	maxLevelReached := defeatedProvider.GetMaxLevelReached()

	filteredEnemyTypes := make([]domain.EnemyType, 0)

	// 1. 撃破済み敵を全て追加
	for _, et := range allEnemyTypes {
		if defeatedProvider.IsEnemyDefeated(et.ID) {
			filteredEnemyTypes = append(filteredEnemyTypes, et)
		}
	}

	// 2. 未撃破敵: MaxLevelReached+1 以上のデフォルトLvを持つ敵の中で最小レベルの1体を追加
	var nextUndefeated *domain.EnemyType
	minNextLevel := 101 // 最大レベル+1

	for i := range allEnemyTypes {
		et := &allEnemyTypes[i]
		if defeatedProvider.IsEnemyDefeated(et.ID) {
			continue
		}
		defaultLevel := et.DefaultLevel
		if defaultLevel < 1 {
			defaultLevel = 1
		}
		// MaxLevelReached+1 以上かつ最小のデフォルトレベルを持つ敵を選択
		if defaultLevel >= maxLevelReached+1 && defaultLevel < minNextLevel {
			minNextLevel = defaultLevel
			nextUndefeated = et
		}
	}

	if nextUndefeated != nil {
		filteredEnemyTypes = append(filteredEnemyTypes, *nextUndefeated)
	}

	s := &BattleSelectScreenCarousel{
		agentProvider:    agentProvider,
		defeatedProvider: defeatedProvider,
		enemyTypes:       filteredEnemyTypes,
		selectedTypeIdx:  0,
		styles:           styles.NewGameStyles(),
		width:            140,
		height:           40,
	}

	// 初期選択敵タイプのレベル範囲を設定
	if len(filteredEnemyTypes) > 0 {
		s.updateLevelRange()
	}

	return s
}

// updateLevelRange は現在選択中の敵タイプに応じてレベル範囲を更新します。
func (s *BattleSelectScreenCarousel) updateLevelRange() {
	if len(s.enemyTypes) == 0 {
		return
	}

	enemyType := s.enemyTypes[s.selectedTypeIdx]
	defaultLevel := enemyType.DefaultLevel
	if defaultLevel < 1 {
		defaultLevel = 1
	}

	s.minSelectableLevel = defaultLevel

	// 撃破済みの場合は到達最高レベル（MaxLevelReached）まで選択可能
	if s.defeatedProvider.IsEnemyDefeated(enemyType.ID) {
		maxLevelReached := s.defeatedProvider.GetMaxLevelReached()
		s.maxSelectableLevel = maxLevelReached
		if s.maxSelectableLevel > 100 {
			s.maxSelectableLevel = 100
		}
		// minがmaxより大きくならないように調整
		if s.maxSelectableLevel < s.minSelectableLevel {
			s.maxSelectableLevel = s.minSelectableLevel
		}
	} else {
		// 未撃破の場合はデフォルトレベルのみ
		s.maxSelectableLevel = defaultLevel
	}

	// 選択レベルをデフォルトレベルにリセット
	s.selectedLevel = defaultLevel
}

// Init は画面の初期化を行います。
func (s *BattleSelectScreenCarousel) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *BattleSelectScreenCarousel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *BattleSelectScreenCarousel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}

	case tea.KeyLeft:
		// 左キーで前の敵タイプへ（ループ）
		if len(s.enemyTypes) > 0 {
			s.selectedTypeIdx--
			if s.selectedTypeIdx < 0 {
				s.selectedTypeIdx = len(s.enemyTypes) - 1
			}
			s.updateLevelRange()
		}
		return s, nil

	case tea.KeyRight:
		// 右キーで次の敵タイプへ（ループ）
		if len(s.enemyTypes) > 0 {
			s.selectedTypeIdx++
			if s.selectedTypeIdx >= len(s.enemyTypes) {
				s.selectedTypeIdx = 0
			}
			s.updateLevelRange()
		}
		return s, nil

	case tea.KeyUp:
		// 上キーでレベル上昇（撃破済みの場合のみ有効）
		if s.selectedLevel < s.maxSelectableLevel {
			s.selectedLevel++
		}
		return s, nil

	case tea.KeyDown:
		// 下キーでレベル下降
		if s.selectedLevel > s.minSelectableLevel {
			s.selectedLevel--
		}
		return s, nil

	case tea.KeyEnter:
		// バトル開始
		equippedAgents := s.agentProvider.GetEquippedAgents()
		if len(equippedAgents) == 0 {
			s.error = "エージェントが装備されていません。\nエージェント管理でエージェントを装備してください。"
			return s, nil
		}

		if len(s.enemyTypes) == 0 {
			s.error = "敵タイプが読み込まれていません。"
			return s, nil
		}

		selectedEnemy := s.enemyTypes[s.selectedTypeIdx]
		return s, func() tea.Msg {
			return StartBattleMsg{
				Level:       s.selectedLevel,
				EnemyTypeID: selectedEnemy.ID,
			}
		}
	}

	return s, nil
}

// View は画面をレンダリングします。
func (s *BattleSelectScreenCarousel) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("バトル選択"))
	builder.WriteString("\n\n")

	if len(s.enemyTypes) == 0 {
		builder.WriteString("敵タイプが読み込まれていません")
		return builder.String()
	}

	// 敵選択カルーセル
	s.renderEnemyCarousel(&builder)

	// 敵情報パネル
	s.renderEnemyInfoPanel(&builder)

	// レベル選択
	s.renderLevelSelector(&builder)

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

	builder.WriteString(hintStyle.Render("←→: 敵選択  ↑↓: レベル選択  Enter: バトル開始  Esc: 戻る"))

	return builder.String()
}

// renderEnemyCarousel は敵選択カルーセルをレンダリングします。
func (s *BattleSelectScreenCarousel) renderEnemyCarousel(builder *strings.Builder) {
	// カルーセル表示：< [敵名] >
	selectedEnemy := s.enemyTypes[s.selectedTypeIdx]

	carouselStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	carousel := fmt.Sprintf("◀  %s  ▶", selectedEnemy.Name)
	builder.WriteString(carouselStyle.Render(carousel))
	builder.WriteString("\n")

	// 敵インデックス表示
	indexStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	indexInfo := fmt.Sprintf("(%d / %d)", s.selectedTypeIdx+1, len(s.enemyTypes))
	builder.WriteString(indexStyle.Render(indexInfo))
	builder.WriteString("\n\n")
}

// renderEnemyInfoPanel は敵情報パネルをレンダリングします。
func (s *BattleSelectScreenCarousel) renderEnemyInfoPanel(builder *strings.Builder) {
	selectedEnemy := s.enemyTypes[s.selectedTypeIdx]

	infoPanel := components.NewInfoPanel("敵情報")
	infoPanel.AddItem("名前", selectedEnemy.Name)
	infoPanel.AddItem("攻撃属性", s.formatAttackType(selectedEnemy.AttackType))
	infoPanel.AddItem("基礎HP", fmt.Sprintf("%d", selectedEnemy.BaseHP))
	infoPanel.AddItem("デフォルトLv", fmt.Sprintf("%d", selectedEnemy.DefaultLevel))

	// パッシブスキル情報（descriptionを表示）
	if selectedEnemy.NormalPassive != nil {
		infoPanel.AddItem("通常パッシブ", "★"+selectedEnemy.NormalPassive.Description)
	}
	if selectedEnemy.EnhancedPassive != nil {
		infoPanel.AddItem("強化パッシブ", "★"+selectedEnemy.EnhancedPassive.Description)
	}

	// 撃破状態
	if s.defeatedProvider.IsEnemyDefeated(selectedEnemy.ID) {
		defeatedLevel := s.defeatedProvider.GetDefeatedLevel(selectedEnemy.ID)
		infoPanel.AddItem("撃破済み", fmt.Sprintf("最高Lv.%d", defeatedLevel))
	} else {
		infoPanel.AddItem("撃破状態", "未撃破")
	}

	infoPanelRendered := infoPanel.Render(50)
	centeredInfo := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(infoPanelRendered)
	builder.WriteString(centeredInfo)
	builder.WriteString("\n\n")
}

// formatAttackType は攻撃タイプを日本語に変換します。
func (s *BattleSelectScreenCarousel) formatAttackType(attackType string) string {
	switch attackType {
	case "physical":
		return "物理"
	case "magic":
		return "魔法"
	default:
		return attackType
	}
}

// renderLevelSelector はレベル選択をレンダリングします。
func (s *BattleSelectScreenCarousel) renderLevelSelector(builder *strings.Builder) {
	levelStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Width(s.width)

	var levelDisplay string
	if s.minSelectableLevel == s.maxSelectableLevel {
		// 未撃破：レベル固定
		levelDisplay = fmt.Sprintf("挑戦レベル: Lv.%d (固定)", s.selectedLevel)
	} else {
		// 撃破済み：レベル選択可能
		levelDisplay = fmt.Sprintf("挑戦レベル: ▲ Lv.%d ▼ (Lv.%d 〜 Lv.%d)",
			s.selectedLevel, s.minSelectableLevel, s.maxSelectableLevel)
	}

	builder.WriteString(levelStyle.Render(levelDisplay))
	builder.WriteString("\n\n")
}
