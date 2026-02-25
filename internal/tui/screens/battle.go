// Package screens はTUIゲームの画面を提供します。
// battle.go はバトル画面のModel構造体とInit/Updateメソッドを担当します。
// UIレンダリングはbattle_view.go、ゲームロジックはbattle_logic.goに分離されています。
package screens

import (
	"fmt"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/ascii"
	"hirorocky/type-battle/internal/tui/challenges"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/combat/chain"
	"hirorocky/type-battle/internal/usecase/combat/recast"
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
	Victory   bool
	Level     int
	Stats     *combat.BattleStatistics // バトル統計
	EnemyID   string                   // 敵図鑑更新用
	EnemyType *domain.EnemyType        // 確定ドロップ設定参照用
}

// ==================== スキルスロット ====================

// SkillSlot はスキルスロットを表します。
type SkillSlot struct {
	Skill             *domain.SkillModel
	Agent             *domain.AgentModel
	AgentIndex        int
	SkillIndex        int
	CooldownRemaining float64
	CooldownTotal     float64
}

// IsReady はスキルが使用可能かを返します。
func (s *SkillSlot) IsReady() bool {
	return s.CooldownRemaining <= 0
}

// SetFeatureUnlockProvider は機能解放状態プロバイダーを設定します。
func (s *BattleScreen) SetFeatureUnlockProvider(provider FeatureUnlockProvider) {
	s.unlockProvider = provider
}

// isDefenseSkillUnlocked はディフェンススキルが解放済みかを返します。
func (s *BattleScreen) isDefenseSkillUnlocked() bool {
	return s.unlockProvider != nil && s.unlockProvider.IsUnlocked(domain.FeatureDefenseSkill)
}

// isManaSystemUnlocked はマナシステムが解放済みかを返します。
func (s *BattleScreen) isManaSystemUnlocked() bool {
	return s.unlockProvider != nil && s.unlockProvider.IsUnlocked(domain.FeatureManaSystem)
}

// isLatentEffectUnlocked は潜在効果が解放済みかを返します。
func (s *BattleScreen) isLatentEffectUnlocked() bool {
	return s.unlockProvider != nil && s.unlockProvider.IsUnlocked(domain.FeatureLatentEffect)
}

// SetPassiveSkills はパッシブスキル定義を設定します。
// これにより、RegisterPassiveSkills で条件付きパッシブスキルが EffectTable に登録されます。
func (s *BattleScreen) SetPassiveSkills(skills map[string]domain.PassiveSkill) {
	if s.battleEngine != nil {
		s.battleEngine.SetPassiveSkills(skills)
	}
}

// ==================== BattleScreen構造体 ====================

// BattleScreen はバトル画面を表します。

// UI-Improvement Requirements: 3.1, 3.2, 3.9
type BattleScreen struct {
	// 戦闘参加者
	enemy          *domain.EnemyModel
	player         *domain.PlayerModel
	equippedAgents []*domain.AgentModel

	// スキルスロット
	skillSlots   []SkillSlot
	selectedSlot int

	// エージェント選択状態（UI改善: 3エリアレイアウト用）
	selectedAgentIdx int

	// チャレンジシステム
	activeChallenge  challenges.ChallengeModel
	selectedSkillIdx int
	dictionary       []string // チャレンジ生成用辞書

	// バトルエンジン
	battleEngine *combat.BattleEngine
	battleState  *combat.BattleState

	// リキャスト・チェイン効果管理
	recastManager      *recast.RecastManager
	chainEffectManager *chain.ChainEffectManager

	// パッシブスキル関連
	comboCount            int // ミスなし連続タイピング回数
	firstStrikeAgentIndex int // ps_first_strike発動エージェント（-1は無効）

	// ゲーム終了状態
	gameOver      bool
	victory       bool
	showingResult bool

	// パーフェクト演出状態
	showingPerfect  bool
	perfectTimer    int
	perfectRenderer ascii.PerfectRenderer

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

	// 機能解放状態プロバイダー（ゲート判定用）
	unlockProvider FeatureUnlockProvider
}

// ==================== コンストラクタ ====================

// NewBattleScreen は新しいBattleScreenを作成します。
// dictionaryがnilの場合はデフォルト辞書を使用します。
func NewBattleScreen(enemy *domain.EnemyModel, player *domain.PlayerModel, agents []*domain.AgentModel, dictionary *typing.Dictionary) *BattleScreen {
	// 辞書がnilの場合はデフォルト辞書を使用
	if dictionary == nil {
		dictionary = createDefaultDictionary()
	}

	// 辞書を統合してフラットな[]stringにする
	allWords := make([]string, 0)
	allWords = append(allWords, dictionary.Easy...)
	allWords = append(allWords, dictionary.Medium...)
	allWords = append(allWords, dictionary.Hard...)

	// 敵タイプリストを作成（BattleEngine用）
	enemyTypes := []domain.EnemyType{enemy.Type}

	gs := styles.NewGameStyles()
	screen := &BattleScreen{
		enemy:            enemy,
		player:           player,
		equippedAgents:   agents,
		skillSlots:       make([]SkillSlot, 0),
		selectedSlot:     0,
		selectedAgentIdx: 0,
		dictionary:       allWords,
		battleEngine:     combat.NewBattleEngine(enemyTypes),
		// リキャスト・チェイン効果管理を初期化
		recastManager:      recast.NewRecastManager(),
		chainEffectManager: chain.NewChainEffectManager(),
		// パッシブスキル関連
		firstStrikeAgentIndex: -1, // 無効値で初期化
		styles:                gs,
		winLoseRenderer:       ascii.NewWinLoseRenderer(gs),
		perfectRenderer:       ascii.NewPerfectRenderer(gs),
		width:                 140,
		height:                40,
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
	}

	// 最初の行動を準備してチャージ開始
	screen.battleState.Enemy.PrepareNextAction()
	screen.battleEngine.StartEnemyCharging(screen.battleState, time.Now())

	// 敵のパッシブスキルを登録
	screen.battleEngine.RegisterEnemyPassive(screen.battleState)

	// スキルスロットを初期化

	for agentIdx, agent := range agents {
		for skillIdx, skill := range agent.Skills {
			screen.skillSlots = append(screen.skillSlots, SkillSlot{
				Skill:             skill,
				Agent:             agent,
				AgentIndex:        agentIdx,
				SkillIndex:        skillIdx,
				CooldownRemaining: 0,
				CooldownTotal:     config.DefaultSkillCooldown,
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
	// チャージシステムはInitBattleで初期化済み（PrepareNextAction + StartCharging）

	// ps_first_strike: バトル開始時に最初のスキル即発動を評価
	s.evaluateFirstStrike()

	return s.tick()
}

// evaluateFirstStrike はps_first_strikeの発動を評価します。
func (s *BattleScreen) evaluateFirstStrike() {
	if s.battleEngine == nil || s.battleState == nil {
		return
	}

	for agentIdx, agent := range s.battleState.EquippedAgents {
		if s.battleEngine.EvaluateFirstStrike(s.battleState, agent) {
			s.firstStrikeAgentIndex = agentIdx
			s.message = fmt.Sprintf("[ファーストストライク！ %sが即発動可能！]", agent.Core.Name)
			return
		}
	}
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

	default:
		// チャレンジ固有のtickメッセージ（standardTickMsg, symbolStormTickMsg等）を転送
		if s.activeChallenge != nil {
			updated, cmd := s.activeChallenge.Update(msg)
			s.activeChallenge = updated

			// チャレンジ完了チェック
			if result := s.activeChallenge.Result(); result != nil {
				s.handleChallengeComplete(result)
			}

			return s, cmd
		}
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

	// パーフェクト演出タイマー（約0.2秒 = 2 tick後に消去）
	if s.showingPerfect {
		s.perfectTimer++
		if s.perfectTimer >= 2 {
			s.showingPerfect = false
		}
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

	// リキャストを更新（チェイン効果の期限切れも処理）
	s.UpdateRecasts(deltaSeconds)

	// UI改善: アニメーション更新
	deltaMS := int(tickInterval.Milliseconds())
	s.floatingDamageManager.Update(deltaMS)
	s.playerHPBar.Update(deltaMS)
	s.enemyHPBar.Update(deltaMS)

	// チャレンジ完了チェック（ChallengeModel.Result()で判定）
	if s.activeChallenge != nil {
		if result := s.activeChallenge.Result(); result != nil {
			s.handleChallengeComplete(result)
		}
	}

	// ディフェンス終了チェック
	now := time.Now()
	if s.enemy.WaitMode == domain.WaitModeDefending && !s.enemy.IsDefenseActive(now) {
		s.enemy.EndDefense()
		// 次の行動を準備してチャージ開始
		s.enemy.PrepareNextAction()
		if action := s.enemy.GetNextAction(); action != nil {
			s.battleEngine.StartEnemyCharging(s.battleState, now)
		}
		s.message = fmt.Sprintf("%sのディフェンスが終了した", s.enemy.Name)
	}

	// 敵攻撃チェック（チャージ完了判定）
	if s.enemy.IsChargeComplete(now) {
		// ディフェンスチャレンジ中の場合、防御率を適用して攻撃解決後にチャレンジを自動終了
		if dp, ok := s.activeChallenge.(challenges.DefenseProvider); ok {
			s.processEnemyAttackWithDefense(dp)
		} else {
			s.processEnemyAttack()
		}

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

	// バフ・デバフの持続時間とボルテージを更新
	if s.battleEngine != nil && s.battleState != nil {
		s.battleEngine.UpdateEffects(s.battleState, deltaSeconds)
	}

	// 次のtickを返す
	return s, s.tick()
}

// handleKeyMsg はキーボード入力を処理します。
func (s *BattleScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// 結果表示中はEnterでのみ遷移
	if s.showingResult {
		return s.handleResultInput(msg)
	}

	if s.activeChallenge != nil {
		return s.handleTypingInput(msg)
	}

	return s.handleSkillSelection(msg)
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

// handleSkillSelection はスキル選択時のキー処理を行います。
// UI-Improvement: 左右キーでエージェント切替、上下キーでスキル選択
func (s *BattleScreen) handleSkillSelection(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "left", "h":
		// 前のエージェントに切り替え
		s.selectedAgentIdx--
		if s.selectedAgentIdx < 0 {
			s.selectedAgentIdx = len(s.equippedAgents) - 1
		}
		// そのエージェントの最初のスキルを選択
		s.selectFirstSkillOfAgent(s.selectedAgentIdx)
	case "right", "l":
		// 次のエージェントに切り替え
		s.selectedAgentIdx++
		if s.selectedAgentIdx >= len(s.equippedAgents) {
			s.selectedAgentIdx = 0
		}
		// そのエージェントの最初のスキルを選択
		s.selectFirstSkillOfAgent(s.selectedAgentIdx)
	case "up", "k":
		// 現在のエージェント内で前のスキルに移動
		s.moveToPrevSkillInAgent()
	case "down", "j":
		// 現在のエージェント内で次のスキルに移動
		s.moveToNextSkillInAgent()
	case "enter":
		// スキル使用可能チェック（クールダウンとリキャスト両方）
		if len(s.skillSlots) > 0 && s.isSkillUsable(s.selectedSlot) {
			s.selectedSkillIdx = s.selectedSlot
			skill := s.skillSlots[s.selectedSlot].Skill

			// ChallengeInput を構築してチャレンジを開始
			cmd := s.startChallenge(skill)

			// スキル選択直後にクールダウンとリキャストを開始
			slot := s.skillSlots[s.selectedSkillIdx]
			s.StartCooldown(s.selectedSkillIdx, slot.CooldownTotal)
			s.startAgentRecast(slot.AgentIndex, skill)

			if cmd != nil {
				return s, cmd
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

// handleTypingInput はチャレンジ中のキー処理を行います。
// ChallengeModel.Update() に委譲します。
func (s *BattleScreen) handleTypingInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if s.activeChallenge == nil {
		return s, nil
	}

	updated, cmd := s.activeChallenge.Update(msg)
	s.activeChallenge = updated

	// チャレンジ完了チェック
	if result := s.activeChallenge.Result(); result != nil {
		s.handleChallengeComplete(result)
	}

	return s, cmd
}
