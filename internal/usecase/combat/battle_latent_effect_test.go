// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// === 受け入れ基準 7,8: 潜在効果のフィルタリング ===

// TestApplySkillEffect_LatentEffect_PerfectTriggersLatent は
// パーフェクト時に潜在効果が発動判定されることを確認します。
func TestApplySkillEffect_LatentEffect_PerfectTriggersLatent(t *testing.T) {
	engine, state, agent, skill := setupLatentEffectTest(t, true)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
		IsPerfect:      true,
	}

	initialHP := state.Enemy.HP
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	if state.Enemy.HP >= initialHP {
		t.Error("パーフェクト時に潜在効果（IsLatent=true）が発動してダメージを与えるべき")
	}
}

// TestApplySkillEffect_LatentEffect_NonPerfectSkipsLatent は
// 非パーフェクト時に潜在効果がスキップされることを確認します。
func TestApplySkillEffect_LatentEffect_NonPerfectSkipsLatent(t *testing.T) {
	engine, state, agent, skill := setupLatentEffectTest(t, true)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       0.9,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.9,
		IsPerfect:      false,
	}

	initialHP := state.Enemy.HP
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	if state.Enemy.HP != initialHP {
		t.Errorf("非パーフェクト時に潜在効果（IsLatent=true）はスキップされるべき: got HP=%d, want HP=%d", state.Enemy.HP, initialHP)
	}
}

// TestApplySkillEffect_NormalEffect_AlwaysTriggered は
// 通常効果（IsLatent=false）がパーフェクト状態に関係なく発動することを確認します。
func TestApplySkillEffect_NormalEffect_AlwaysTriggered(t *testing.T) {
	engine, state, agent, skill := setupLatentEffectTest(t, false)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       0.8,
		SpeedFactor:    1.0,
		AccuracyFactor: 0.8,
		IsPerfect:      false,
	}

	initialHP := state.Enemy.HP
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	if state.Enemy.HP >= initialHP {
		t.Error("通常効果（IsLatent=false）は非パーフェクト時でもダメージを与えるべき")
	}
}

// === 受け入れ基準 10: 複数潜在効果の独立判定 ===

// TestApplySkillEffect_MultipleLatentEffects_IndependentTrigger は
// 複数の潜在効果が独立して判定されることを確認します。
func TestApplySkillEffect_MultipleLatentEffects_IndependentTrigger(t *testing.T) {
	// 2つの潜在効果: 確率1.0（必ず発動）と確率0.0（絶対に発動しない）
	// 各効果が独立して判定されることを検証
	skillType := domain.SkillType{
		ID:   "multi_latent",
		Name: "複数潜在効果スキル",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				IsLatent:    true,
				Icon:        "⚔️",
			},
			{
				Target:      domain.TargetSelf,
				HPFormula:   &domain.HPFormula{Base: 30, StatCoef: 0, StatRef: "WIL"},
				Probability: 0.0,
				IsLatent:    true,
				Icon:        "💚",
			},
		},
	}

	engine, state, agent := createLatentTestSetup(t, skillType)
	skill := domain.NewSkillFromType(skillType, nil)
	agent.Skills = []*domain.SkillModel{skill}

	// HP状態を調整
	state.Player.TakeDamage(50)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
		IsPerfect:      true,
	}

	initialEnemyHP := state.Enemy.HP
	initialPlayerHP := state.Player.HP

	engine.ApplySkillEffect(state, agent, skill, typingResult)

	// 敵にダメージが入った（確率1.0の潜在効果が発動）
	if state.Enemy.HP >= initialEnemyHP {
		t.Error("確率1.0の潜在効果が敵にダメージを与えるべき")
	}
	// プレイヤーは回復しない（確率0.0の潜在効果は発動しない）
	if state.Player.HP != initialPlayerHP {
		t.Errorf("確率0.0の潜在効果は発動しないべき: got HP=%d, want HP=%d", state.Player.HP, initialPlayerHP)
	}
}

// === 受け入れ基準 11: Echo/DoubleCastでパーフェクト状態引き継ぎ ===

// TestApplySkillEffectWithEcho_LatentEffect_PerfectInherited は
// Echo/DoubleCastの追加発動で潜在効果が引き継がれることを確認します。
func TestApplySkillEffectWithEcho_LatentEffect_PerfectInherited(t *testing.T) {
	engine, state, agent, skill := setupLatentEffectTest(t, true)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
		IsPerfect:      true,
	}

	initialHP := state.Enemy.HP
	// 2回発動（エコー）
	totalEffect := engine.ApplySkillEffectWithEcho(state, agent, skill, typingResult, 2)

	if totalEffect == 0 {
		t.Error("エコースキルで潜在効果が2回発動してダメージを与えるべき")
	}
	if state.Enemy.HP >= initialHP {
		t.Error("エコースキルの追加発動でも潜在効果が発動するべき")
	}
}

// TestApplySkillEffectWithCombo_LatentEffect_PerfectInherited は
// DoubleCast（ApplySkillEffectWithCombo）の追加発動でパーフェクト状態が引き継がれることを確認します。
func TestApplySkillEffectWithCombo_LatentEffect_PerfectInherited(t *testing.T) {
	engine, state, agent, skill := setupLatentEffectTest(t, true)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
		IsPerfect:      true,
	}

	initialHP := state.Enemy.HP
	// DoubleCast経路: ApplySkillEffectWithComboを2回呼び出し
	effect1 := engine.ApplySkillEffectWithCombo(state, agent, skill, typingResult, 0)
	effect2 := engine.ApplySkillEffectWithCombo(state, agent, skill, typingResult, 0)

	if effect1 == 0 || effect2 == 0 {
		t.Errorf("DoubleCastの両方の発動で潜在効果がダメージを与えるべき: effect1=%d, effect2=%d", effect1, effect2)
	}
	if state.Enemy.HP >= initialHP {
		t.Error("DoubleCastの追加発動でも潜在効果が発動するべき")
	}
}

// === ヘルパー関数 ===

// setupLatentEffectTest は潜在効果テスト用のセットアップを行います。
func setupLatentEffectTest(t *testing.T, isLatent bool) (*BattleEngine, *BattleState, *domain.AgentModel, *domain.SkillModel) {
	t.Helper()

	skillType := domain.SkillType{
		ID:   "test_latent",
		Name: "テスト潜在効果スキル",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 100, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				IsLatent:    isLatent,
				Icon:        "⚔️",
			},
		},
	}

	engine, state, agent := createLatentTestSetup(t, skillType)
	skill := domain.NewSkillFromType(skillType, nil)
	agent.Skills = []*domain.SkillModel{skill}

	return engine, state, agent, skill
}

// createLatentTestSetup はテスト用のエンジン、状態、エージェントを作成します。
func createLatentTestSetup(t *testing.T, skillType domain.SkillType) (*BattleEngine, *BattleState, *domain.AgentModel) {
	t.Helper()

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_none", Name: "なし"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	skill := domain.NewSkillFromType(skillType, nil)
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	enemyTypes := []domain.EnemyType{
		{
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 10000,
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	return engine, state, agent
}
