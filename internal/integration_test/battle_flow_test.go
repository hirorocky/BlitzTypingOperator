// Package integration_test は統合テストを提供します。

package integration_test

import (
	"fmt"
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/typing"
)

// newTestDamageModuleBattle はテスト用のダメージモジュールを作成するヘルパー関数です。
func newTestDamageModuleBattle(id, name string, tags []string, statCoef float64, statRef, description string) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
}

// newTestHealModuleBattle はテスト用の回復モジュールを作成するヘルパー関数です。
func newTestHealModuleBattle(id, name string, tags []string, statCoef float64, statRef, description string) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          id,
		Name:        name,
		Icon:        "💚",
		Tags:        tags,
		Description: description,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetSelf,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "💚",
			},
		},
	}, nil)
}

// newTestBuffModuleBattle はテスト用のバフモジュールを作成するヘルパー関数です。
func newTestBuffModuleBattle(id, name string, tags []string, value float64, _, description string) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          id,
		Name:        name,
		Icon:        "⬆️",
		Tags:        tags,
		Description: description,
		Effects: []domain.SkillEffect{
			{
				Target: domain.TargetSelf,
				ColumnSpec: &domain.EffectColumnSpec{
					Column:   domain.ColDamageBonus,
					Value:    value,
					Duration: 10.0,
				},
				Probability: 1.0,
				Icon:        "⬆️",
			},
		},
	}, nil)
}

// ==================================================
// Task 15.2: バトルフロー統合テスト
// ==================================================

// createTestAgents はテスト用のエージェントを作成します。
func createTestAgents() []*domain.AgentModel {
	coreType := domain.CoreType{
		ID:   "all_rounder",
		Name: "オールラウンダー",
		StatWeights: map[string]float64{
			"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0,
		},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト", Description: ""}
	core := domain.NewCoreWithTypeID("core_1", coreType, passiveSkill)

	modules := []*domain.SkillModel{
		newTestDamageModuleBattle("m1", "物理打撃Lv1", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModuleBattle("m2", "ファイアボールLv1", []string{"magic_low"}, 1.0, "MAG", ""),
		newTestHealModuleBattle("m3", "ヒールLv1", []string{"heal_low"}, 0.8, "MAG", ""),
		newTestBuffModuleBattle("m4", "バフLv1", []string{"buff_low"}, 5.0, "SPD", ""),
	}

	return []*domain.AgentModel{
		domain.NewAgent("agent_1", core, modules),
	}
}

// createTestEnemyTypes はテスト用の敵タイプを作成します。
func createTestEnemyTypes() []domain.EnemyType {
	return []domain.EnemyType{
		{
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 50,
			Rank:   1,
		},
	}
}

// initBattleForTest はテスト用のバトル初期化ヘルパーです。
// BattleStateを直接構築し、敵の生成と行動準備を行います。
func initBattleForTest(engine *combat.BattleEngine, level int, agents []*domain.AgentModel, enemyTypes []domain.EnemyType) *combat.BattleState {
	// 敵を生成
	enemyType := enemyTypes[0]
	hp := enemyType.BaseHP * level
	attackPower := 5 + (level * 2) // 固定のベース攻撃力を使用
	enemy := domain.NewEnemy(
		"test-enemy-id",
		enemyType.Name+" Lv."+fmt.Sprintf("%d", level),
		level,
		hp,
		attackPower,
		enemyType,
	)

	// プレイヤーを初期化
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	player.PrepareForBattle()

	// バトル状態を作成
	state := &combat.BattleState{
		Enemy:          enemy,
		Player:         player,
		EquippedAgents: agents,
		Level:          level,
		Stats: &combat.BattleStatistics{
			StartTime: time.Now(),
		},
	}

	// 最初の行動を準備してチャージ開始
	state.Enemy.PrepareNextAction()
	engine.StartEnemyCharging(state, time.Now())

	return state
}

func TestBattleFlow_Initialize(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// 敵が生成されている
	if state.Enemy == nil {
		t.Error("敵が生成されるべきです")
	}

	// プレイヤーが初期化されている
	if state.Player == nil {
		t.Error("プレイヤーが初期化されるべきです")
	}

	// プレイヤーHPが最大値
	if state.Player.HP != state.Player.MaxHP {
		t.Error("プレイヤーHPは最大値であるべきです")
	}
}

func TestBattleFlow_EnemyAttack(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)
	initialHP := state.Player.HP

	// 敵攻撃を処理
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// ダメージが与えられた
	if damage <= 0 {
		t.Error("ダメージは0より大きいべきです")
	}

	// プレイヤーHPが減少
	if state.Player.HP >= initialHP {
		t.Error("プレイヤーHPが減少するべきです")
	}
}

func TestBattleFlow_ModuleUse_Attack(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)
	initialEnemyHP := state.Enemy.HP

	// タイピング結果
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	// 物理攻撃モジュールを使用
	agent := agents[0]
	module := agent.Modules[0] // 物理打撃
	damage := engine.ApplySkillEffect(state, agent, module, typingResult)

	// ダメージが与えられた
	if damage <= 0 {
		t.Errorf("ダメージは0より大きいべきです: got %d", damage)
	}

	// 敵HPが減少
	if state.Enemy.HP >= initialEnemyHP {
		t.Error("敵HPが減少するべきです")
	}
}

func TestBattleFlow_ModuleUse_Heal(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// プレイヤーにダメージを与える
	state.Player.TakeDamage(30)
	damagedHP := state.Player.HP

	// タイピング結果
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	// 回復モジュールを使用
	agent := agents[0]
	module := agent.Modules[2] // ヒール
	healAmount := engine.ApplySkillEffect(state, agent, module, typingResult)

	// 回復量が正の値
	if healAmount <= 0 {
		t.Errorf("回復量は0より大きいべきです: got %d", healAmount)
	}

	// プレイヤーHPが増加
	if state.Player.HP <= damagedHP {
		t.Error("プレイヤーHPが増加するべきです")
	}
}

func TestBattleFlow_VictoryCondition(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// 敵HPを0にする
	state.Enemy.HP = 0

	// 勝敗判定
	ended, result := engine.CheckBattleEnd(state)

	if !ended {
		t.Error("バトルが終了するべきです")
	}
	if !result.IsVictory {
		t.Error("プレイヤーの勝利であるべきです")
	}
}

func TestBattleFlow_DefeatCondition(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// プレイヤーHPを0にする
	state.Player.HP = 0

	// 勝敗判定
	ended, result := engine.CheckBattleEnd(state)

	if !ended {
		t.Error("バトルが終了するべきです")
	}
	if result.IsVictory {
		t.Error("プレイヤーの敗北であるべきです")
	}
}

func TestBattleFlow_PhaseTransition(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// 初期フェーズは通常
	if state.Enemy.Phase != domain.PhaseNormal {
		t.Error("初期フェーズは通常であるべきです")
	}

	// HP50%以下に設定
	state.Enemy.HP = state.Enemy.MaxHP / 2

	// フェーズ変化チェック
	transitioned := engine.CheckPhaseTransition(state)

	if !transitioned {
		t.Error("フェーズ変化が発生するべきです")
	}
	if state.Enemy.Phase != domain.PhaseEnhanced {
		t.Error("フェーズが強化になるべきです")
	}
}

func TestBattleFlow_TypingChallenge(t *testing.T) {
	// タイピングチャレンジ→完了→効果計算の流れ
	eval := typing.NewEvaluator()

	// チャレンジ作成
	challenge := &typing.Challenge{
		Text:       "hello",
		TimeLimit:  10 * time.Second,
		Difficulty: typing.DifficultyEasy,
	}

	// チャレンジ開始
	state := eval.StartChallenge(challenge)
	if state == nil {
		t.Fatal("チャレンジ状態が作成されるべきです")
	}

	// 文字入力をシミュレート
	for _, char := range "hello" {
		eval.ProcessInput(state, char)
	}

	// 完了判定
	if !eval.IsCompleted(state) {
		t.Error("チャレンジが完了しているべきです")
	}

	// 結果取得
	result := eval.CompleteChallenge(state)
	if !result.Completed {
		t.Error("結果のCompletedがtrueであるべきです")
	}
	if result.WPM <= 0 {
		t.Error("WPMが計算されるべきです")
	}
}

func TestBattleFlow_BuffDebuffInteraction(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// プレイヤーに防御バフを付与
	state.Player.EffectTable.AddBuff("防御UP", "防御UP", 10.0, map[domain.EffectColumn]float64{
		domain.ColDamageCut: 0.5, // 50%軽減
	})

	// 敵攻撃
	damage := engine.ProcessEnemyAttackDamage(state, "physical")
	expectedMaxDamage := state.Enemy.AttackPower // バフなしの場合

	// ダメージが軽減されている
	if damage >= expectedMaxDamage {
		t.Errorf("ダメージが軽減されるべきです: got %d, expected < %d", damage, expectedMaxDamage)
	}
}

func TestBattleFlow_AccuracyPenalty(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	// バトル初期化
	_ = initBattleForTest(engine, 1, agents, enemyTypes)

	agent := agents[0]
	module := agent.Modules[0]

	// 高い正確性
	highAccuracyResult := &typing.TypingResult{
		Completed:      true,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.95,
	}
	highDamage := engine.CalculateSkillEffectWithPassive(agent, module, highAccuracyResult)

	// 低い正確性（50%未満）
	lowAccuracyResult := &typing.TypingResult{
		Completed:      true,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.4,
	}
	lowDamage := engine.CalculateSkillEffectWithPassive(agent, module, lowAccuracyResult)

	// 低い正確性の方が効果が低い（半減ペナルティ適用）
	if lowDamage >= highDamage {
		t.Errorf("低い正確性の効果は高い正確性より低いべきです: low=%d, high=%d", lowDamage, highDamage)
	}
}

func TestBattleFlow_Statistics(t *testing.T) {
	// バトル統計の記録
	enemyTypes := createTestEnemyTypes()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 1, agents, enemyTypes)

	// タイピング結果を記録
	result := &typing.TypingResult{
		Completed: true,
		WPM:       100,
		Accuracy:  0.95,
	}
	engine.RecordTypingResult(state, result)

	if state.Stats.TotalTypingCount != 1 {
		t.Error("タイピング回数がカウントされるべきです")
	}
	if state.Stats.TotalWPM != 100 {
		t.Error("WPMが記録されるべきです")
	}
}

// ==================== Task 7.3: バトルシステム統合テスト ====================

// createTestEnemyTypesWithPatterns は行動パターンとパッシブを持つ敵タイプを作成します。
func createTestEnemyTypesWithPatterns() []domain.EnemyType {
	normalActions := []domain.EnemyAction{
		{
			ID:             "slash",
			Name:           "斬撃",
			ActionType:     domain.EnemyActionAttack,
			AttackType:     "physical",
			DamageBase:     10.0,
			DamagePerLevel: 2.0,
			ChargeTime:     500 * time.Millisecond,
		},
		{
			ID:            "power_up",
			Name:          "パワーアップ",
			ActionType:    domain.EnemyActionBuff,
			EffectColumn:  domain.ColDamageMultiplier,
			EffectValue:   1.5,
			Duration:      10.0,
			TimedEffectID: "power_up",
			ChargeTime:    300 * time.Millisecond,
		},
	}
	enhancedActions := []domain.EnemyAction{
		{
			ID:             "rage_strike",
			Name:           "怒りの一撃",
			ActionType:     domain.EnemyActionAttack,
			AttackType:     "physical",
			DamageBase:     20.0,
			DamagePerLevel: 4.0,
			ChargeTime:     1 * time.Second,
		},
	}

	normalPassive := &domain.EnemyPassiveSkill{
		ID:          "defensive_stance",
		Name:        "防御姿勢",
		Description: "ダメージ10%軽減",
		Effects: map[domain.EffectColumn]float64{
			domain.ColDamageCut: 0.1,
		},
	}
	enhancedPassive := &domain.EnemyPassiveSkill{
		ID:          "berserker_rage",
		Name:        "バーサーカーの怒り",
		Description: "攻撃力50%上昇",
		Effects: map[domain.EffectColumn]float64{
			domain.ColDamageMultiplier: 1.5,
		},
	}

	return []domain.EnemyType{
		{
			ID:                      "boss_goblin",
			Name:                    "ゴブリンリーダー",
			BaseHP:                  100,
			Rank:                    1,
			ResolvedNormalActions:   normalActions,
			ResolvedEnhancedActions: enhancedActions,
			NormalPassive:           normalPassive,
			EnhancedPassive:         enhancedPassive,
			DropItemCategory:        "core",
			DropItemTypeID:          "attack_balance",
		},
	}
}

// TestBattleFlow_LevelSelection_PatternBased はレベル選択→パターンベースバトル→報酬の一連フローをテストします。
func TestBattleFlow_LevelSelection_PatternBased(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithPatterns()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	// レベル10でバトル開始
	level := 10
	state := initBattleForTest(engine, level, agents, enemyTypes)

	// 敵タイプを手動で設定（テスト用）
	state.Enemy = domain.NewEnemy("test", "ゴブリンリーダー Lv.10", level, 1000, 30, enemyTypes[0])

	// 敵が正しく生成されている
	if state.Enemy == nil {
		t.Fatal("敵が生成されるべきです")
	}
	if state.Enemy.Level != level {
		t.Errorf("敵レベルが不正: got %d, want %d", state.Enemy.Level, level)
	}
	if state.Enemy.Type.ID != "boss_goblin" {
		t.Errorf("敵タイプIDが不正: got %s, want boss_goblin", state.Enemy.Type.ID)
	}

	// 行動パターンが設定されている
	if len(state.Enemy.Type.ResolvedNormalActions) == 0 {
		t.Error("通常行動パターンが設定されるべきです")
	}

	// 敵パッシブスキルを登録
	engine.RegisterEnemyPassive(state)

	// パッシブスキルが適用されている
	passives := state.Enemy.EffectTable.FindBySourceType(domain.SourcePassive)
	if len(passives) == 0 {
		t.Error("敵のパッシブスキルが登録されるべきです")
	}
	if state.Enemy.ActivePassiveID != "defensive_stance" {
		t.Errorf("ActivePassiveIDが不正: got %s", state.Enemy.ActivePassiveID)
	}

	// 行動決定がパターンベース
	nextAction := engine.DeterminePatternBasedAction(state)
	if nextAction.SourceAction == nil {
		t.Error("パターンベースの行動が決定されるべきです")
	}
	if nextAction.SourceAction.ID != "slash" {
		t.Errorf("最初の行動は斬撃であるべき: got %s", nextAction.SourceAction.ID)
	}
}

// TestBattleFlow_PhaseTransition_PassiveSwitch はフェーズ遷移とパッシブ切り替えの連携をテストします。
func TestBattleFlow_PhaseTransition_PassiveSwitch(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithPatterns()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 10, agents, enemyTypes)
	state.Enemy = domain.NewEnemy("test", "ゴブリンリーダー Lv.10", 10, 1000, 30, enemyTypes[0])
	engine.RegisterEnemyPassive(state)

	// 初期状態：通常フェーズ、通常パッシブ
	if state.Enemy.Phase != domain.PhaseNormal {
		t.Error("初期フェーズは通常であるべきです")
	}
	if state.Enemy.ActivePassiveID != "defensive_stance" {
		t.Errorf("初期パッシブが不正: got %s", state.Enemy.ActivePassiveID)
	}

	// 通常パッシブ効果の確認
	ctx := domain.NewEffectContext(state.Player.HP, state.Player.MaxHP, state.Enemy.HP, state.Enemy.MaxHP)
	effects := state.Enemy.EffectTable.Aggregate(ctx)
	if effects.DamageCut != 0.1 {
		t.Errorf("通常パッシブ効果が不正: DamageCut got %f, want 0.1", effects.DamageCut)
	}

	// HP50%以下に設定してフェーズ遷移
	state.Enemy.HP = state.Enemy.MaxHP / 2

	if !engine.CheckPhaseTransition(state) {
		t.Error("フェーズ遷移が発生するべきです")
	}

	// パッシブ切り替え
	state.Enemy.ResetActionIndex() // 行動パターンリセット
	engine.SwitchEnemyPassive(state)

	// 強化フェーズ、強化パッシブに変更
	if state.Enemy.Phase != domain.PhaseEnhanced {
		t.Error("フェーズが強化になるべきです")
	}
	if state.Enemy.ActivePassiveID != "berserker_rage" {
		t.Errorf("強化パッシブが適用されるべき: got %s", state.Enemy.ActivePassiveID)
	}

	// 強化パッシブ効果の確認
	effects = state.Enemy.EffectTable.Aggregate(ctx)
	if effects.DamageMultiplier != 1.5 {
		t.Errorf("強化パッシブ効果が不正: DamageMultiplier got %f, want 1.5", effects.DamageMultiplier)
	}
	if effects.DamageCut != 0.0 {
		t.Errorf("通常パッシブが無効化されるべき: DamageCut got %f, want 0.0", effects.DamageCut)
	}

	// 行動インデックスがリセットされている
	if state.Enemy.ActionIndex != 0 {
		t.Errorf("行動インデックスがリセットされるべき: got %d", state.Enemy.ActionIndex)
	}

	// 強化パターンの行動を取得
	nextAction := state.Enemy.GetCurrentAction()
	if nextAction.ID != "rage_strike" {
		t.Errorf("強化パターンの行動が取得されるべき: got %s", nextAction.ID)
	}
}

// TestBattleFlow_PatternLoopAndProgress は行動パターンのループと進行をテストします。
func TestBattleFlow_PatternLoopAndProgress(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithPatterns()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 5, agents, enemyTypes)
	state.Enemy = domain.NewEnemy("test", "ゴブリンリーダー Lv.5", 5, 500, 20, enemyTypes[0])

	// 通常パターンは2つ（斬撃→パワーアップ）
	if len(state.Enemy.Type.ResolvedNormalActions) != 2 {
		t.Fatalf("通常パターンは2つあるべき: got %d", len(state.Enemy.Type.ResolvedNormalActions))
	}

	// 行動1: 斬撃
	action1 := state.Enemy.GetCurrentAction()
	if action1.ID != "slash" {
		t.Errorf("1回目は斬撃であるべき: got %s", action1.ID)
	}
	state.Enemy.AdvanceActionIndex()

	// 行動2: パワーアップ
	action2 := state.Enemy.GetCurrentAction()
	if action2.ID != "power_up" {
		t.Errorf("2回目はパワーアップであるべき: got %s", action2.ID)
	}
	state.Enemy.AdvanceActionIndex()

	// 行動3: 斬撃（ループ）
	action3 := state.Enemy.GetCurrentAction()
	if action3.ID != "slash" {
		t.Errorf("3回目は斬撃（ループ）であるべき: got %s", action3.ID)
	}
	if state.Enemy.ActionIndex != 0 {
		t.Errorf("ActionIndexは0にループするべき: got %d", state.Enemy.ActionIndex)
	}
}

// TestBattleFlow_EnemyDefeatRecordIntegration は敵撃破記録の一連のフローをテストします。
func TestBattleFlow_EnemyDefeatRecordIntegration(t *testing.T) {
	// 撃破記録マップ
	defeatRecords := make(map[string]int)

	enemyTypes := createTestEnemyTypesWithPatterns()
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	// 初回バトル（レベル5）
	level := 5
	state := initBattleForTest(engine, level, agents, enemyTypes)
	state.Enemy = domain.NewEnemy("test", "ゴブリンリーダー Lv.5", level, 500, 20, enemyTypes[0])

	// バトル終了（勝利）
	state.Enemy.HP = 0
	ended, result := engine.CheckBattleEnd(state)

	if !ended || !result.IsVictory {
		t.Fatal("勝利であるべきです")
	}

	// 撃破記録を更新
	enemyTypeID := state.Enemy.Type.ID
	if currentMax, ok := defeatRecords[enemyTypeID]; !ok || level > currentMax {
		defeatRecords[enemyTypeID] = level
	}

	if defeatRecords["boss_goblin"] != 5 {
		t.Errorf("撃破記録が不正: got %d, want 5", defeatRecords["boss_goblin"])
	}

	// 2回目のバトル（より高いレベルで挑戦可能）
	level2 := 10
	state2 := initBattleForTest(engine, level2, agents, enemyTypes)
	state2.Enemy = domain.NewEnemy("test", "ゴブリンリーダー Lv.10", level2, 1000, 30, enemyTypes[0])
	state2.Enemy.HP = 0
	engine.CheckBattleEnd(state2)

	// 撃破記録を更新
	if currentMax, ok := defeatRecords[enemyTypeID]; !ok || level2 > currentMax {
		defeatRecords[enemyTypeID] = level2
	}

	if defeatRecords["boss_goblin"] != 10 {
		t.Errorf("撃破記録が更新されるべき: got %d, want 10", defeatRecords["boss_goblin"])
	}

	// 次回選択可能レベルの範囲確認
	maxSelectableLevel := defeatRecords["boss_goblin"] + 1
	if maxSelectableLevel != 11 {
		t.Errorf("次回選択可能な最大レベルは11であるべき: got %d", maxSelectableLevel)
	}
}

// TestBattleFlow_ChargingAndDefenseIntegration はチャージングとディフェンスの統合テストです。
func TestBattleFlow_ChargingAndDefenseIntegration(t *testing.T) {
	defenseAction := domain.EnemyAction{
		ID:          "guard",
		Name:        "ガード",
		ActionType:  domain.EnemyActionDefense,
		DefenseType: domain.DefensePhysicalCut,
		EffectValue: 0.5,
		Duration:    5.0,
	}

	enemyTypes := []domain.EnemyType{
		{
			ID:     "defender_goblin",
			Name:   "ディフェンダーゴブリン",
			BaseHP: 80,
			Rank:   1,
			ResolvedNormalActions: []domain.EnemyAction{
				{
					ID:             "attack",
					Name:           "攻撃",
					ActionType:     domain.EnemyActionAttack,
					AttackType:     "physical",
					DamageBase:     10.0,
					DamagePerLevel: 2.0,
					ChargeTime:     1 * time.Second,
				},
				defenseAction,
			},
		},
	}

	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 10, agents, enemyTypes)
	state.Enemy = domain.NewEnemy("test", "ディフェンダーゴブリン Lv.10", 10, 800, 8, enemyTypes[0])

	// 攻撃行動でチャージ開始
	now := time.Now()
	attackAction := state.Enemy.GetCurrentAction()
	state.Enemy.StartCharging(attackAction, now)

	if state.Enemy.WaitMode != domain.WaitModeCharging {
		t.Error("チャージモードになるべきです")
	}

	// チャージ進捗確認
	progress := state.Enemy.GetChargeProgress(now)
	if progress != 0 {
		t.Errorf("開始直後の進捗は0であるべき: got %f", progress)
	}

	// チャージ完了後、行動実行
	executedAction := state.Enemy.ExecuteChargedAction()
	if executedAction == nil || executedAction.ID != "attack" {
		t.Error("攻撃行動が実行されるべきです")
	}

	// 次の行動はディフェンス
	nextAction := state.Enemy.GetCurrentAction()
	if nextAction.ID != "guard" {
		t.Errorf("次の行動はガードであるべき: got %s", nextAction.ID)
	}

	// ディフェンス発動
	state.Enemy.StartDefense(
		defenseAction.DefenseType,
		defenseAction.EffectValue,
		time.Duration(defenseAction.Duration*float64(time.Second)),
		now,
	)

	if state.Enemy.WaitMode != domain.WaitModeDefending {
		t.Error("ディフェンスモードになるべきです")
	}

	// ディフェンス中のダメージ軽減確認
	baseDamage := 100
	reducedDamage := engine.ApplyDefenseReduction(state, baseDamage, "physical")
	if reducedDamage != 50 {
		t.Errorf("ダメージが50%%軽減されるべき: got %d", reducedDamage)
	}

	// 魔法ダメージは軽減されない
	magicDamage := engine.ApplyDefenseReduction(state, baseDamage, "magic")
	if magicDamage != baseDamage {
		t.Errorf("魔法ダメージは軽減されないべき: got %d", magicDamage)
	}
}

// TestBattleFlow_BuffDebuffPattern はバフ・デバフ行動パターンの統合テストです。
func TestBattleFlow_BuffDebuffPattern(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "buff_goblin",
			Name:   "バフゴブリン",
			BaseHP: 60,
			Rank:   1,
			ResolvedNormalActions: []domain.EnemyAction{
				{
					ID:            "self_buff",
					Name:          "自己強化",
					ActionType:    domain.EnemyActionBuff,
					EffectColumn:  domain.ColDamageMultiplier,
					EffectValue:   2.0,
					Duration:      10.0,
					TimedEffectID: "self_buff",
				},
				{
					ID:            "slow_curse",
					Name:          "スロウの呪い",
					ActionType:    domain.EnemyActionDebuff,
					EffectColumn:  domain.ColCooldownReduce,
					EffectValue:   -0.3,
					Duration:      8.0,
					TimedEffectID: "slow_curse",
				},
			},
		},
	}

	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	state := initBattleForTest(engine, 5, agents, enemyTypes)
	state.Enemy = domain.NewEnemy("test", "バフゴブリン Lv.5", 5, 300, 6, enemyTypes[0])

	// 自己バフ行動を実行
	buffAction := state.Enemy.GetCurrentAction()
	if buffAction.ID != "self_buff" {
		t.Errorf("最初の行動は自己強化であるべき: got %s", buffAction.ID)
	}
	engine.ApplyPatternBuff(state, buffAction)

	// バフが敵のEffectTableに登録されている
	buffs := state.Enemy.EffectTable.FindBySourceType(domain.SourceBuff)
	if len(buffs) == 0 {
		t.Error("バフがEffectTableに登録されるべきです")
	}

	// バフ効果の確認
	ctx := domain.NewEffectContext(state.Player.HP, state.Player.MaxHP, state.Enemy.HP, state.Enemy.MaxHP)
	effects := state.Enemy.EffectTable.Aggregate(ctx)
	if effects.DamageMultiplier != 2.0 {
		t.Errorf("DamageMultiplierが適用されるべき: got %f", effects.DamageMultiplier)
	}

	// 次の行動（デバフ）
	state.Enemy.AdvanceActionIndex()
	debuffAction := state.Enemy.GetCurrentAction()
	if debuffAction.ID != "slow_curse" {
		t.Errorf("次の行動はスロウの呪いであるべき: got %s", debuffAction.ID)
	}
	engine.ApplyPatternDebuff(state, debuffAction)

	// デバフがプレイヤーのEffectTableに登録されている
	debuffs := state.Player.EffectTable.FindBySourceType(domain.SourceDebuff)
	if len(debuffs) == 0 {
		t.Error("デバフがプレイヤーのEffectTableに登録されるべきです")
	}
}
