// Package integration_test は統合テストを提供します。

package integration_test

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/typing"
)

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
	core := domain.NewCore("core_1", "テストコア", 5, coreType, passiveSkill)

	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "物理打撃Lv1", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "ファイアボールLv1", domain.MagicAttack, 1, []string{"magic_low"}, 10.0, "MAG", ""),
		domain.NewModule("m3", "ヒールLv1", domain.Heal, 1, []string{"heal_low"}, 8.0, "MAG", ""),
		domain.NewModule("m4", "バフLv1", domain.Buff, 1, []string{"buff_low"}, 5.0, "SPD", ""),
	}

	return []*domain.AgentModel{
		domain.NewAgent("agent_1", core, modules),
	}
}

// createTestEnemyTypes はテスト用の敵タイプを作成します。
func createTestEnemyTypes() []domain.EnemyType {
	return []domain.EnemyType{
		{
			ID:                 "goblin",
			Name:               "ゴブリン",
			BaseHP:             50,
			BaseAttackPower:    5,
			BaseAttackInterval: 3 * time.Second,
			AttackType:         "physical",
		},
	}
}

func TestBattleFlow_Initialize(t *testing.T) {

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, err := engine.InitializeBattle(1, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)
	initialHP := state.Player.HP

	// 敵攻撃を処理
	damage := engine.ProcessEnemyAttack(state)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)
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
	damage := engine.ApplyModuleEffect(state, agent, module, typingResult)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

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
	healAmount := engine.ApplyModuleEffect(state, agent, module, typingResult)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

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

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

	// プレイヤーに防御バフを付与
	duration := 10.0
	state.Player.EffectTable.AddRow(domain.EffectRow{
		ID:         "def_buff",
		SourceType: domain.SourceBuff,
		Name:       "防御UP",
		Duration:   &duration,
		Modifiers:  domain.StatModifiers{DamageReduction: 0.5}, // 50%軽減
	})

	// 敵攻撃
	damage := engine.ProcessEnemyAttack(state)
	expectedMaxDamage := state.Enemy.AttackPower // バフなしの場合

	// ダメージが軽減されている
	if damage >= expectedMaxDamage {
		t.Errorf("ダメージが軽減されるべきです: got %d, expected < %d", damage, expectedMaxDamage)
	}
}

func TestBattleFlow_AccuracyPenalty(t *testing.T) {

	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	// バトル初期化（stateは使用しないが、エンジン初期化のために呼び出す）
	_, _ = engine.InitializeBattle(1, agents)

	agent := agents[0]
	module := agent.Modules[0]

	// 高い正確性
	highAccuracyResult := &typing.TypingResult{
		Completed:      true,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.95,
	}
	highDamage := engine.CalculateModuleEffect(agent, module, highAccuracyResult)

	// 低い正確性（50%未満）
	lowAccuracyResult := &typing.TypingResult{
		Completed:      true,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.4,
	}
	lowDamage := engine.CalculateModuleEffect(agent, module, lowAccuracyResult)

	// 低い正確性の方が効果が低い（半減ペナルティ適用）
	if lowDamage >= highDamage {
		t.Errorf("低い正確性の効果は高い正確性より低いべきです: low=%d, high=%d", lowDamage, highDamage)
	}
}

func TestBattleFlow_Statistics(t *testing.T) {
	// バトル統計の記録
	engine := combat.NewBattleEngine(createTestEnemyTypes())
	agents := createTestAgents()

	state, _ := engine.InitializeBattle(1, agents)

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
