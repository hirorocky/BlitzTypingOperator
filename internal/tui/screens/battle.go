// Package screens はTUIゲームの画面を提供します。
// battle.go はバトル画面のModel構造体とInit/Updateメソッドを担当します。
// UIレンダリングはbattle_view.go、ゲームロジックはbattle_logic.goに分離されています。
package screens

import (
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/ascii"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/typing"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.3: バトル画面 ====================

// tickInterval はバトル画面の更新間隔です。
// config.BattleTickIntervalを参照しています。
var tickInterval = config.BattleTickInterval

// ==================== メッセージ型 ====================

// BattleTickMsg はバトル画面の定期更新メッセージです。
type BattleTickMsg struct{}

// BattleResultMsg はバトル結果メッセージです。
type BattleResultMsg struct {
	Victory bool
	Level   int
	Stats   *combat.BattleStatistics // バトル統計
	EnemyID string                   // 敵図鑑更新用
}

// ==================== モジュールスロット ====================

// ModuleSlot はモジュールスロットを表します。
type ModuleSlot struct {
	Module            *domain.ModuleModel
	Agent             *domain.AgentModel
	AgentIndex        int
	ModuleIndex       int
	CooldownRemaining float64
	CooldownTotal     float64
}

// IsReady はモジュールが使用可能かを返します。
func (s *ModuleSlot) IsReady() bool {
	return s.CooldownRemaining <= 0
}

// ==================== BattleScreen構造体 ====================

// BattleScreen はバトル画面を表します。

// UI-Improvement Requirements: 3.1, 3.2, 3.9
type BattleScreen struct {
	// 戦闘参加者
	enemy          *domain.EnemyModel
	player         *domain.PlayerModel
	equippedAgents []*domain.AgentModel

	// モジュールスロット
	moduleSlots  []ModuleSlot
	selectedSlot int

	// エージェント選択状態（UI改善: 3エリアレイアウト用）
	selectedAgentIdx int

	// タイピング状態
	isTyping          bool
	typingText        string
	typingIndex       int
	typingMistakes    []int
	typingStartTime   time.Time
	typingTimeLimit   time.Duration
	selectedModuleIdx int

	// タイピングシステム
	challengeGenerator *typing.ChallengeGenerator
	evaluator          *typing.Evaluator
	typingState        *typing.ChallengeState

	// バトルエンジン
	battleEngine *combat.BattleEngine
	battleState  *combat.BattleState

	// 敵攻撃
	nextEnemyAttack time.Time

	// ゲーム終了状態
	gameOver      bool
	victory       bool
	showingResult bool

	// UI
	styles          *styles.GameStyles
	winLoseRenderer ascii.WinLoseRenderer
	width           int
	height          int
	message         string

	// アニメーション（UI改善: フローティングダメージ、HPバーアニメーション）
	floatingDamageManager *styles.FloatingDamageManager
	playerHPBar           *styles.AnimatedHPBar
	enemyHPBar            *styles.AnimatedHPBar
}

// ==================== コンストラクタ ====================

// NewBattleScreen は新しいBattleScreenを作成します。
func NewBattleScreen(enemy *domain.EnemyModel, player *domain.PlayerModel, agents []*domain.AgentModel) *BattleScreen {
	// デフォルト辞書を作成
	dictionary := createDefaultDictionary()

	// 敵タイプリストを作成（BattleEngine用）
	enemyTypes := []domain.EnemyType{enemy.Type}

	gs := styles.NewGameStyles()
	screen := &BattleScreen{
		enemy:              enemy,
		player:             player,
		equippedAgents:     agents,
		moduleSlots:        make([]ModuleSlot, 0),
		selectedSlot:       0,
		selectedAgentIdx:   0,
		isTyping:           false,
		challengeGenerator: typing.NewChallengeGenerator(dictionary),
		evaluator:          typing.NewEvaluator(),
		battleEngine:       combat.NewBattleEngine(enemyTypes),
		styles:             gs,
		winLoseRenderer:    ascii.NewWinLoseRenderer(gs),
		width:              140,
		height:             40,
		// UI改善: アニメーション初期化
		floatingDamageManager: styles.NewFloatingDamageManager(),
		playerHPBar:           styles.NewAnimatedHPBar(player.MaxHP),
		enemyHPBar:            styles.NewAnimatedHPBar(enemy.MaxHP),
	}

	// バトル状態を初期化
	screen.battleState = &combat.BattleState{
		Enemy:          enemy,
		Player:         player,
		EquippedAgents: agents,
		Level:          enemy.Level,
		Stats: &combat.BattleStatistics{
			StartTime: time.Now(),
		},
		NextAttackTime: time.Now().Add(enemy.AttackInterval),
	}

	screen.battleState.NextAction = screen.battleEngine.DetermineNextAction(screen.battleState)

	// モジュールスロットを初期化

	for agentIdx, agent := range agents {
		for modIdx, module := range agent.Modules {
			screen.moduleSlots = append(screen.moduleSlots, ModuleSlot{
				Module:            module,
				Agent:             agent,
				AgentIndex:        agentIdx,
				ModuleIndex:       modIdx,
				CooldownRemaining: 0,
				CooldownTotal:     config.DefaultModuleCooldown,
			})
		}
	}

	return screen
}

// createDefaultDictionary はデフォルトのタイピング辞書を作成します。
func createDefaultDictionary() *typing.Dictionary {
	return &typing.Dictionary{
		// Easy: 3-6文字の単語
		Easy: []string{
			"cat", "dog", "run", "jump", "fire", "ice",
			"hit", "cut", "heal", "buff", "fast", "slow",
			"axe", "bow", "sun", "moon", "star", "wind",
			"red", "blue", "gold", "dark", "life", "mana",
		},
		// Medium: 7-11文字の単語
		Medium: []string{
			"warrior", "monster", "defense", "attack",
			"healing", "protect", "thunder", "blizzard",
			"fireball", "critical", "accuracy", "strength",
			"powerful", "ultimate", "blessing", "cursed",
		},
		// Hard: 12-20文字の単語
		Hard: []string{
			"thunderstorm", "annihilation", "resurrection",
			"extraordinary", "invulnerable", "battleground",
			"concentration", "determination", "acceleration",
			"purification", "hallucination", "obliteration",
		},
	}
}

// ==================== tea.Modelインターフェース実装 ====================

// Init は画面の初期化を行います。
func (s *BattleScreen) Init() tea.Cmd {
	s.nextEnemyAttack = time.Now().Add(s.enemy.AttackInterval)
	return s.tick()
}

// tick は次のtickコマンドを返します。
func (s *BattleScreen) tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return BattleTickMsg{}
	})
}

// Update はメッセージを処理します。
func (s *BattleScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case BattleTickMsg:
		return s.handleTick()

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return s, nil
}

// ==================== メッセージハンドラ ====================

// handleTick は定期更新を処理します。
func (s *BattleScreen) handleTick() (tea.Model, tea.Cmd) {
	// ゲーム終了済みなら何もしない（最優先）
	if s.gameOver {
		return s, nil
	}

	// 結果表示中はtickを継続するが、ゲーム進行はしない
	if s.showingResult {
		return s, s.tick()
	}

	deltaSeconds := tickInterval.Seconds()

	// 勝敗判定（結果表示状態に入る）
	if s.checkGameOver() {
		s.showingResult = true
		// HP表示を実際のHPに即座に合わせる
		if s.enemy.HP <= 0 {
			s.enemyHPBar.SetTarget(0)
			s.enemyHPBar.ForceComplete()
		}
		if s.player.HP <= 0 {
			s.playerHPBar.SetTarget(0)
			s.playerHPBar.ForceComplete()
		}
		return s, s.tick()
	}

	// クールダウンを更新
	s.UpdateCooldowns(deltaSeconds)

	// UI改善: アニメーション更新
	deltaMS := int(tickInterval.Milliseconds())
	s.floatingDamageManager.Update(deltaMS)
	s.playerHPBar.Update(deltaMS)
	s.enemyHPBar.Update(deltaMS)

	// タイピング中の時間切れチェック
	if s.isTyping {
		elapsed := time.Since(s.typingStartTime)
		if elapsed >= s.typingTimeLimit {
			s.CancelTyping()
			s.message = "タイムアップ！"
		}
	}

	// 敵攻撃チェック
	if time.Now().After(s.nextEnemyAttack) {
		s.processEnemyAttack()

		// 攻撃後の敗北判定（結果表示状態に入る）
		if s.checkGameOver() {
			s.showingResult = true
			// HP表示を実際のHPに即座に合わせる
			if s.enemy.HP <= 0 {
				s.enemyHPBar.SetTarget(0)
				s.enemyHPBar.ForceComplete()
			}
			if s.player.HP <= 0 {
				s.playerHPBar.SetTarget(0)
				s.playerHPBar.ForceComplete()
			}
			return s, s.tick()
		}
	}

	// バフ・デバフの持続時間を更新
	s.updateEffectDurations(deltaSeconds)

	// 次のtickを返す
	return s, s.tick()
}

// handleKeyMsg はキーボード入力を処理します。
func (s *BattleScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 結果表示中はEnterでのみ遷移
	if s.showingResult {
		return s.handleResultInput(msg)
	}

	if s.isTyping {
		return s.handleTypingInput(msg)
	}

	return s.handleModuleSelection(msg)
}

// handleResultInput は結果表示中のキー入力を処理します。
func (s *BattleScreen) handleResultInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		// Enterで結果を確定してホームに戻る
		return s, s.createGameOverCmd()
	}
	// Enter以外のキーは無視
	return s, nil
}

// handleModuleSelection はモジュール選択時のキー処理を行います。
// UI-Improvement: 左右キーでエージェント切替、上下キーでモジュール選択
func (s *BattleScreen) handleModuleSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		// 前のエージェントに切り替え
		s.selectedAgentIdx--
		if s.selectedAgentIdx < 0 {
			s.selectedAgentIdx = len(s.equippedAgents) - 1
		}
		// そのエージェントの最初のモジュールを選択
		s.selectFirstModuleOfAgent(s.selectedAgentIdx)
	case "right", "l":
		// 次のエージェントに切り替え
		s.selectedAgentIdx++
		if s.selectedAgentIdx >= len(s.equippedAgents) {
			s.selectedAgentIdx = 0
		}
		// そのエージェントの最初のモジュールを選択
		s.selectFirstModuleOfAgent(s.selectedAgentIdx)
	case "up", "k":
		// 現在のエージェント内で前のモジュールに移動
		s.moveToPrevModuleInAgent()
	case "down", "j":
		// 現在のエージェント内で次のモジュールに移動
		s.moveToNextModuleInAgent()
	case "enter":
		if len(s.moduleSlots) > 0 && s.moduleSlots[s.selectedSlot].IsReady() {
			// モジュール選択 → タイピングチャレンジ開始
			s.selectedModuleIdx = s.selectedSlot
			module := s.moduleSlots[s.selectedSlot].Module

			// モジュールレベルに応じた難易度でチャレンジを生成
			difficulty := typing.GetDifficultyForModuleLevel(module.Level)
			timeLimit := typing.GetDefaultTimeLimit(difficulty)
			challenge := s.challengeGenerator.Generate(difficulty, timeLimit)

			if challenge != nil {
				s.StartTypingChallenge(challenge.Text, challenge.TimeLimit)
			}
		}
	case "esc":
		// バトルを中断してホームに戻る（デバッグ用）
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	}

	return s, nil
}

// handleTypingInput はタイピング中のキー処理を行います。
func (s *BattleScreen) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// タイピングをキャンセル
		s.CancelTyping()
		return s, nil
	default:
		if len(msg.Runes) == 1 {
			s.ProcessTypingInput(msg.Runes[0])
		}
	}

	return s, nil
}
