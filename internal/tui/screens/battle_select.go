// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strconv"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/ascii"
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

	// 未撃破時のデフォルトレベルは1
	const defaultLevel = 1
	for i := range allEnemyTypes {
		et := &allEnemyTypes[i]
		if defeatedProvider.IsEnemyDefeated(et.ID) {
			continue
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
	// 未撃破時のデフォルトレベルは1
	const defaultLevel = 1

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
	infoPanel.AddItem("基礎HP", fmt.Sprintf("%d", selectedEnemy.BaseHP))

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

// ==================== ランクベースのバトル選択画面 ====================

// BattleSelectRankState はランクベースバトル選択画面の状態を表します。
type BattleSelectRankState int

const (
	// StateEnemySelect は敵選択状態です。
	StateEnemySelect BattleSelectRankState = iota
	// StateLevelModal はレベル選択モーダル表示状態です。
	StateLevelModal
)

// BattleSelectScreenRankBased はランクベースのバトル選択画面を表します。
// 現在のランクに属する敵のみを表示し、撃破状況に応じてレベル選択を制限します。
type BattleSelectScreenRankBased struct {
	agentProvider    AgentProvider
	progressProvider EnemyProgressProvider
	enemyTypes       []domain.EnemyType

	// 画面状態
	state BattleSelectRankState

	// ランク選択用
	selectedRank    int // 画面上で選択中のランク（ゲーム進行度とは独立）
	maxUnlockedRank int // 解放済み最大ランク

	// 敵種類選択用
	selectedTypeIdx int

	// レベル選択用（メイン画面の状態）
	selectedLevel      int
	minSelectableLevel int
	maxSelectableLevel int

	// モーダル内でのレベル選択用（キャンセル時に元に戻すため分離）
	modalLevel int

	error  string
	styles *styles.GameStyles
	width  int
	height int

	// ASCIIアートレンダラー（ランク表示用）
	shadowRenderer ascii.ShadowNumberRenderer
}

// NewBattleSelectScreenRankBased は新しいランクベースのBattleSelectScreenを作成します。
func NewBattleSelectScreenRankBased(
	agentProvider AgentProvider,
	progressProvider EnemyProgressProvider,
) *BattleSelectScreenRankBased {
	// 初期表示は最大解放ランク
	maxUnlockedRank := progressProvider.GetMaxUnlockedRank()
	currentRankEnemies := progressProvider.GetEnemiesByRank(maxUnlockedRank)

	s := &BattleSelectScreenRankBased{
		agentProvider:    agentProvider,
		progressProvider: progressProvider,
		enemyTypes:       currentRankEnemies,
		state:            StateEnemySelect,
		selectedRank:     maxUnlockedRank,
		maxUnlockedRank:  maxUnlockedRank,
		selectedTypeIdx:  0,
		styles:           styles.NewGameStyles(),
		width:            140,
		height:           40,
		shadowRenderer:   ascii.NewShadowNumbers(),
	}

	// 初期選択敵タイプのレベル範囲を設定
	if len(currentRankEnemies) > 0 {
		s.updateLevelRange()
	}

	return s
}

// updateLevelRange は現在選択中の敵タイプに応じてレベル範囲を更新します。
func (s *BattleSelectScreenRankBased) updateLevelRange() {
	if len(s.enemyTypes) == 0 {
		return
	}

	enemyType := s.enemyTypes[s.selectedTypeIdx]
	// 未撃破時のデフォルトレベルは1
	const defaultLevel = 1

	min, max := s.progressProvider.GetSelectableLevelRange(enemyType.ID)
	s.minSelectableLevel = min
	s.maxSelectableLevel = max

	// 選択レベルをデフォルトレベルにリセット（min以上max以下に収める）
	s.selectedLevel = defaultLevel
	if s.selectedLevel < s.minSelectableLevel {
		s.selectedLevel = s.minSelectableLevel
	}
	if s.selectedLevel > s.maxSelectableLevel {
		s.selectedLevel = s.maxSelectableLevel
	}
}

// Init は画面の初期化を行います。
func (s *BattleSelectScreenRankBased) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *BattleSelectScreenRankBased) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *BattleSelectScreenRankBased) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch s.state {
	case StateEnemySelect:
		return s.handleEnemySelectState(msg)
	case StateLevelModal:
		return s.handleLevelModalState(msg)
	}
	return s, nil
}

// handleEnemySelectState は敵選択状態でのキー処理を行います。
func (s *BattleSelectScreenRankBased) handleEnemySelectState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
			s.error = "" // エラーをクリア
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
			s.error = "" // エラーをクリア
		}
		return s, nil

	case tea.KeyUp:
		// ランク上昇（上限: maxUnlockedRank）
		if s.selectedRank < s.maxUnlockedRank {
			s.changeRank(s.selectedRank + 1)
			s.error = "" // エラーをクリア
		}
		return s, nil

	case tea.KeyDown:
		// ランク下降（下限: 1）
		if s.selectedRank > 1 {
			s.changeRank(s.selectedRank - 1)
			s.error = "" // エラーをクリア
		}
		return s, nil

	case tea.KeyEnter:
		// レベル選択モーダルを開く
		s.openLevelModal()
		return s, nil
	}

	// j/kキーにも対応
	switch msg.String() {
	case "k":
		if s.selectedRank < s.maxUnlockedRank {
			s.changeRank(s.selectedRank + 1)
			s.error = ""
		}
		return s, nil
	case "j":
		if s.selectedRank > 1 {
			s.changeRank(s.selectedRank - 1)
			s.error = ""
		}
		return s, nil
	case "h":
		if len(s.enemyTypes) > 0 {
			s.selectedTypeIdx--
			if s.selectedTypeIdx < 0 {
				s.selectedTypeIdx = len(s.enemyTypes) - 1
			}
			s.updateLevelRange()
			s.error = ""
		}
		return s, nil
	case "l":
		if len(s.enemyTypes) > 0 {
			s.selectedTypeIdx++
			if s.selectedTypeIdx >= len(s.enemyTypes) {
				s.selectedTypeIdx = 0
			}
			s.updateLevelRange()
			s.error = ""
		}
		return s, nil
	}

	return s, nil
}

// handleLevelModalState はレベル選択モーダル状態でのキー処理を行います。
func (s *BattleSelectScreenRankBased) handleLevelModalState(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// モーダルを閉じる（ホームには戻らない）
		s.closeLevelModal()
		return s, nil

	case tea.KeyUp:
		// レベル上昇
		if s.modalLevel < s.maxSelectableLevel {
			s.modalLevel++
		}
		return s, nil

	case tea.KeyDown:
		// レベル下降
		if s.modalLevel > s.minSelectableLevel {
			s.modalLevel--
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
				Level:       s.modalLevel,
				EnemyTypeID: selectedEnemy.ID,
			}
		}
	}

	// j/kキーにも対応
	switch msg.String() {
	case "k":
		if s.modalLevel < s.maxSelectableLevel {
			s.modalLevel++
		}
		return s, nil
	case "j":
		if s.modalLevel > s.minSelectableLevel {
			s.modalLevel--
		}
		return s, nil
	}

	return s, nil
}

// changeRank は選択ランクを変更し、敵リストを更新します。
func (s *BattleSelectScreenRankBased) changeRank(newRank int) {
	if newRank < 1 || newRank > s.maxUnlockedRank {
		return
	}
	s.selectedRank = newRank
	s.enemyTypes = s.progressProvider.GetEnemiesByRank(newRank)

	// 敵リストが空でない場合は選択インデックスをリセット
	if len(s.enemyTypes) > 0 {
		s.selectedTypeIdx = 0
		s.updateLevelRange()
	} else {
		s.selectedTypeIdx = 0
		s.selectedLevel = 1
		s.minSelectableLevel = 1
		s.maxSelectableLevel = 1
	}
}

// openLevelModal はレベル選択モーダルを開きます。
func (s *BattleSelectScreenRankBased) openLevelModal() {
	if len(s.enemyTypes) == 0 {
		s.error = "敵タイプが読み込まれていません。"
		return
	}

	// 装備チェックはモーダルを開く前に行う
	equippedAgents := s.agentProvider.GetEquippedAgents()
	if len(equippedAgents) == 0 {
		s.error = "エージェントが装備されていません。\nエージェント管理でエージェントを装備してください。"
		return
	}

	s.state = StateLevelModal
	s.error = "" // エラーをクリア
	// 初期値を最高選択可能レベルに設定
	s.modalLevel = s.maxSelectableLevel
}

// closeLevelModal はレベル選択モーダルを閉じます。
func (s *BattleSelectScreenRankBased) closeLevelModal() {
	s.state = StateEnemySelect
	s.error = "" // モーダル閉鎖時にエラーをクリア
	// modalLevelはキャンセル時に破棄される（selectedLevelには反映しない）
}

// View は画面をレンダリングします。
func (s *BattleSelectScreenRankBased) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("バトル選択"))
	builder.WriteString("\n\n")

	if len(s.enemyTypes) == 0 && s.state == StateEnemySelect {
		// 敵がいない場合のメッセージ
		emptyMsg := fmt.Sprintf("ランク %d には敵がいません", s.selectedRank)
		builder.WriteString(lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(s.width).
			Render(emptyMsg))
		builder.WriteString("\n\n")

		// エラーメッセージ（あれば表示）
		if s.error != "" {
			errorStyle := lipgloss.NewStyle().
				Foreground(styles.ColorDamage).
				Align(lipgloss.Center).
				Width(s.width)
			builder.WriteString(errorStyle.Render(s.error))
			builder.WriteString("\n\n")
		}

		// ランク移動のヒント
		hintStyle := lipgloss.NewStyle().
			Foreground(styles.ColorSubtle).
			Align(lipgloss.Center).
			Width(s.width)
		hint := "↑↓: ランク選択  Esc: 戻る"
		builder.WriteString(hintStyle.Render(hint))
		return builder.String()
	}

	// ランクと進行状況（上部にランク位置インジケーター含む）
	s.renderRankProgress(&builder)

	// 敵選択カルーセル（左部に敵位置インジケーター含む）
	s.renderEnemyCarousel(&builder)

	// 敵情報パネル（2ペイン）
	s.renderEnemyInfoPanel(&builder)

	// エラーメッセージ
	if s.error != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.ColorDamage).
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(errorStyle.Render(s.error))
		builder.WriteString("\n\n")
	}

	// モーダル表示
	if s.state == StateLevelModal {
		s.renderLevelModal(&builder)
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hint string
	switch s.state {
	case StateEnemySelect:
		hint = "←→: 敵選択  ↑↓: ランク選択  Enter: レベル選択  Esc: 戻る"
	case StateLevelModal:
		hint = "↑↓: レベル選択  Enter: バトル開始  Esc: 戻る"
	}
	builder.WriteString(hintStyle.Render(hint))

	return builder.String()
}

// renderRankProgress はランクと進行状況をレンダリングします。
// 新デザイン: 左に影付きランク表示（上下に▲▼）、右に敵選択インジケーター + 撃破数
func (s *BattleSelectScreenRankBased) renderRankProgress(builder *strings.Builder) {
	// 選択中ランクの撃破状況を計算
	defeatedCount := 0
	defeatedFlags := make([]bool, len(s.enemyTypes))
	for i, et := range s.enemyTypes {
		if s.progressProvider.IsDefeated(et.ID) {
			defeatedFlags[i] = true
			defeatedCount++
		}
	}
	totalCount := len(s.enemyTypes)

	// 左側: 影付きランク表示（選択中ランクを表示）
	rankArt := s.shadowRenderer.RenderShadowNumber(s.selectedRank)

	// ランク上下ナビゲーション
	navStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSecondary).
		Align(lipgloss.Center)

	// 上矢印（ランクが上げられる場合のみ強調）
	var upArrow string
	if s.selectedRank < s.maxUnlockedRank {
		upArrow = navStyle.Render("▲")
	} else {
		upArrow = lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("▲")
	}

	// 下矢印（ランクが下げられる場合のみ強調）
	var downArrow string
	if s.selectedRank > 1 {
		downArrow = navStyle.Render("▼")
	} else {
		downArrow = lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("▼")
	}

	// ランク表示を上下矢印で囲む
	rankWithNav := lipgloss.JoinVertical(
		lipgloss.Center,
		upArrow,
		rankArt,
		downArrow,
	)

	// 右側: 敵選択インジケーター + 撃破数
	enemyIndicator := components.RenderEnemySelectionIndicator(len(s.enemyTypes), s.selectedTypeIdx, defeatedFlags)
	defeatedInfo := fmt.Sprintf("撃破: %d/%d", defeatedCount, totalCount)

	rightContent := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(enemyIndicator + "\n" + defeatedInfo)

	// 左右を横に結合
	content := lipgloss.JoinHorizontal(
		lipgloss.Center,
		rankWithNav,
		lipgloss.NewStyle().Width(20).Render(""), // スペーサー
		rightContent,
	)

	centeredContent := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)

	builder.WriteString(centeredContent)
	builder.WriteString("\n\n")
}

// renderEnemyCarousel は敵選択カルーセルをレンダリングします。
func (s *BattleSelectScreenRankBased) renderEnemyCarousel(builder *strings.Builder) {
	// 敵選択インジケーターは上部のrenderRankProgressに移動したため、ここでは何もしない
}

// renderEnemyInfoPanel は敵情報パネルをレンダリングします。
// 新デザイン: 1段目に基本情報ボックス、2段目に通常時/強化時の2パネル
func (s *BattleSelectScreenRankBased) renderEnemyInfoPanel(builder *strings.Builder) {
	selectedEnemy := s.enemyTypes[s.selectedTypeIdx]

	// 1段目: 基本情報ボックス
	basicInfoBox := s.renderBasicInfoBox(selectedEnemy)
	centeredBasicInfo := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(basicInfoBox)
	builder.WriteString(centeredBasicInfo)
	builder.WriteString("\n\n")

	// 2段目: 通常時/強化時の2パネル
	normalPane := s.renderNormalModePane(selectedEnemy)
	enhancedPane := s.renderEnhancedModePane(selectedEnemy)

	content := lipgloss.JoinHorizontal(lipgloss.Top, normalPane, "  ", enhancedPane)

	centeredContent := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
	builder.WriteString(centeredContent)
	builder.WriteString("\n\n")
}

// renderBasicInfoBox は敵の基本情報ボックスをレンダリングします。
// 1行目: 敵名、撃破状態、基礎HP
// 2行目: ドロップアイテム
func (s *BattleSelectScreenRankBased) renderBasicInfoBox(enemy domain.EnemyType) string {
	var lines []string

	// 1行目: 敵名、撃破状態、基礎HP
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary)
	var defeatedStatus string
	if s.progressProvider.IsDefeated(enemy.ID) {
		defeatedLevel := s.progressProvider.GetMaxDefeatedLevel(enemy.ID)
		defeatedStatus = lipgloss.NewStyle().Foreground(styles.ColorBuff).Render(
			fmt.Sprintf("Lv.%d撃破済み", defeatedLevel))
	} else {
		defeatedStatus = lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("未撃破")
	}

	line1 := fmt.Sprintf("  %s    撃破状態: %s    基礎HP: %d  ",
		nameStyle.Render(enemy.Name),
		defeatedStatus,
		enemy.BaseHP)
	lines = append(lines, line1)

	// 2行目: ドロップアイテム
	dropItemName := s.progressProvider.GetDropItemName(enemy.DropItemCategory, enemy.DropItemTypeID)
	line2 := fmt.Sprintf("  ドロップ: %s  ", dropItemName)
	lines = append(lines, line2)

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(0, 1).
		Width(90).
		Render(content)
}

// renderNormalModePane は通常時のパネルをレンダリングします。
func (s *BattleSelectScreenRankBased) renderNormalModePane(enemy domain.EnemyType) string {
	var lines []string

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary)
	lines = append(lines, titleStyle.Render("── 通常時 ──"))
	lines = append(lines, "")

	// パッシブ
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("パッシブ:"))
	if enemy.NormalPassive != nil {
		lines = append(lines, "  "+enemy.NormalPassive.Description)
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("  なし"))
	}
	lines = append(lines, "")

	// 行動パターン
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("行動パターン:"))
	if len(enemy.ResolvedNormalActions) > 0 {
		for _, action := range enemy.ResolvedNormalActions {
			actionLine := s.formatActionInfo(action)
			lines = append(lines, "  "+actionLine)
		}
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("  データなし"))
	}

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(42).
		Render(content)
}

// renderEnhancedModePane は強化時のパネルをレンダリングします。
func (s *BattleSelectScreenRankBased) renderEnhancedModePane(enemy domain.EnemyType) string {
	var lines []string

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorDamage)
	lines = append(lines, titleStyle.Render("── 強化時 ──"))
	lines = append(lines, "")

	// パッシブ
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("パッシブ:"))
	if enemy.EnhancedPassive != nil {
		lines = append(lines, "  "+enemy.EnhancedPassive.Description)
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("  なし"))
	}
	lines = append(lines, "")

	// 行動パターン
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("行動パターン:"))
	if len(enemy.ResolvedEnhancedActions) > 0 {
		for _, action := range enemy.ResolvedEnhancedActions {
			actionLine := s.formatActionInfo(action)
			lines = append(lines, "  "+actionLine)
		}
	} else {
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("  通常行動を継続"))
	}

	content := strings.Join(lines, "\n")

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSubtle).
		Padding(1).
		Width(42).
		Render(content)
}

// formatActionInfo は敵行動の情報を整形して返します。
func (s *BattleSelectScreenRankBased) formatActionInfo(action domain.EnemyAction) string {
	line := fmt.Sprintf("・%s", action.Name)

	// 行動タイプに応じて追加情報を表示
	switch action.ActionType {
	case domain.EnemyActionAttack:
		if action.DamageBase > 0 {
			line += fmt.Sprintf(" (基礎威力: %.0f)", action.DamageBase)
		}
	case domain.EnemyActionBuff:
		line += " [バフ]"
	case domain.EnemyActionDebuff:
		line += " [デバフ]"
	case domain.EnemyActionDefense:
		line += " [防御]"
	}

	return line
}

// renderLevelModal はレベル選択モーダルをレンダリングします。
func (s *BattleSelectScreenRankBased) renderLevelModal(builder *strings.Builder) {
	if len(s.enemyTypes) == 0 {
		return
	}

	selectedEnemy := s.enemyTypes[s.selectedTypeIdx]

	// モーダルの内容を構築
	var lines []string

	// 敵名
	titleLine := fmt.Sprintf("%s に挑戦", selectedEnemy.Name)
	lines = append(lines, lipgloss.NewStyle().Bold(true).Foreground(styles.ColorPrimary).Render(titleLine))
	lines = append(lines, "")

	// レベル選択
	levelNavStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
	var levelDisplay string
	if s.minSelectableLevel == s.maxSelectableLevel {
		// 未撃破：レベル固定
		levelDisplay = fmt.Sprintf("Lv.%d (固定)", s.modalLevel)
	} else {
		// 撃破済み：レベル選択可能
		var upArrow, downArrow string
		if s.modalLevel >= s.maxSelectableLevel {
			upArrow = lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("▲")
		} else {
			upArrow = levelNavStyle.Render("▲")
		}
		if s.modalLevel <= s.minSelectableLevel {
			downArrow = lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("▼")
		} else {
			downArrow = levelNavStyle.Render("▼")
		}
		levelDisplay = fmt.Sprintf("%s Lv.%d %s", upArrow, s.modalLevel, downArrow)
	}
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(levelDisplay))

	// 選択可能範囲
	if s.minSelectableLevel != s.maxSelectableLevel {
		rangeInfo := fmt.Sprintf("(Lv.%d 〜 Lv.%d)", s.minSelectableLevel, s.maxSelectableLevel)
		lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render(rangeInfo))
	}
	lines = append(lines, "")

	// 操作ヒント
	hintLine := "Enter: 開始  Esc: 戻る"
	lines = append(lines, lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render(hintLine))

	content := strings.Join(lines, "\n")

	// モーダルボックス
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSecondary).
		Padding(1, 3).
		Align(lipgloss.Center)

	modalBox := modalStyle.Render(content)

	// 画面中央に配置
	centeredModal := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(modalBox)

	builder.WriteString(centeredModal)
	builder.WriteString("\n\n")
}
