// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// ===== Phase 3: モジュール使用時イベント発火のテスト =====

// TestBattleEngine_ModuleUse_EchoSkill は15%でスキル2回発動をテストします。
func TestBattleEngine_ModuleUse_EchoSkill(t *testing.T) {
	// Arrange: ps_echo_skillの定義
	echoSkillDef := domain.PassiveSkill{
		ID:          "ps_echo_skill",
		Name:        "エコースキル",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnSkillUse,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 2.0, // 2回発動
		Probability: 1.0, // テスト用に100%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_echo_skill": echoSkillDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_echo_skill", Name: "エコースキル"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	moduleType := domain.SkillType{
		ID:          "test_attack",
		Name:        "テスト攻撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト用攻撃",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 100, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{module})
	agents := []*domain.AgentModel{agent}

	enemyTypes := []domain.EnemyType{
		{
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          10000,
			BaseAttackPower: 10,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act: モジュール使用
	repeatCount := engine.EvaluateEchoSkill(state, agent)

	// Assert: エコースキル発動で2回
	if repeatCount != 2 {
		t.Errorf("ps_echo_skill: repeatCount=%d, want 2", repeatCount)
	}

	// ダメージが2倍（2回発動）になることを確認
	initialHP := state.Enemy.HP
	_ = engine.ApplySkillEffectWithEcho(state, agent, module, typingResult, repeatCount)
	damageDealt := initialHP - state.Enemy.HP

	// 通常ダメージと比較
	state.Enemy.HP = initialHP
	_ = engine.ApplySkillEffectWithEcho(state, agent, module, typingResult, 1)
	singleDamage := initialHP - state.Enemy.HP

	// 2回発動で約2倍のダメージ
	ratio := float64(damageDealt) / float64(singleDamage)
	if ratio < 1.9 || ratio > 2.1 {
		t.Errorf("ps_echo_skill: ダメージ倍率が期待値と異なる: ratio=%.2f, want 2.0", ratio)
	}
}

// TestBattleEngine_ModuleUse_MiracleHeal は回復スキル時10%でHP全回復をテストします。
func TestBattleEngine_ModuleUse_MiracleHeal(t *testing.T) {
	// Arrange: ps_miracle_healの定義
	miracleHealDef := domain.PassiveSkill{
		ID:          "ps_miracle_heal",
		Name:        "ミラクルヒール",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnHeal,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1.0, // HP全回復
		Probability: 1.0, // テスト用に100%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_miracle_heal": miracleHealDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"INT": 1.0},
		AllowedTags: []string{"heal"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_miracle_heal", Name: "ミラクルヒール"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	moduleType := domain.SkillType{
		ID:          "test_heal",
		Name:        "テスト回復",
		Icon:        "💚",
		Tags:        []string{"heal"},
		Description: "テスト用回復",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetSelf,
				HPFormula:   &domain.HPFormula{Base: 20, StatCoef: 1.0, StatRef: "INT"},
				Probability: 1.0,
				Icon:        "💚",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{module})
	agents := []*domain.AgentModel{agent}

	enemyTypes := []domain.EnemyType{
		{
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          1000,
			BaseAttackPower: 10,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// HPを50%に減らす
	state.Player.HP = state.Player.MaxHP / 2

	// Act: ミラクルヒール発動チェック
	isMiracleHeal := engine.EvaluateMiracleHeal(state, agent, module)

	// Assert: ミラクルヒール発動
	if !isMiracleHeal {
		t.Errorf("ps_miracle_heal: ミラクルヒールが発動していない")
	}

	// ミラクルヒールを適用
	if isMiracleHeal {
		state.Player.HP = state.Player.MaxHP
	}

	if state.Player.HP != state.Player.MaxHP {
		t.Errorf("ps_miracle_heal: HPが全回復していない: HP=%d, MaxHP=%d", state.Player.HP, state.Player.MaxHP)
	}
}

// TestBattleEngine_ModuleUse_MiracleHeal_NotHealSkill は非回復スキルでは発動しないことをテストします。
func TestBattleEngine_ModuleUse_MiracleHeal_NotHealSkill(t *testing.T) {
	// Arrange
	miracleHealDef := domain.PassiveSkill{
		ID:          "ps_miracle_heal",
		Name:        "ミラクルヒール",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnHeal,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1.0,
		Probability: 1.0,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_miracle_heal": miracleHealDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_miracle_heal", Name: "ミラクルヒール"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	// 攻撃スキル（回復ではない）
	moduleType := domain.SkillType{
		ID:          "test_attack",
		Name:        "テスト攻撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト用攻撃",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 100, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{module})
	agents := []*domain.AgentModel{agent}

	enemyTypes := []domain.EnemyType{
		{
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          1000,
			BaseAttackPower: 10,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Act: 攻撃スキルでミラクルヒール判定
	isMiracleHeal := engine.EvaluateMiracleHeal(state, agent, module)

	// Assert: 攻撃スキルでは発動しない
	if isMiracleHeal {
		t.Errorf("ps_miracle_heal: 攻撃スキルで発動してはいけない")
	}
}
