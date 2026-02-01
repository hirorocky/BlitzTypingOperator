// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// ===== Phase 1: コンボカウンターとps_combo_masterのテスト =====

// TestBattleEngine_ComboMaster_StackedDamage はコンボによるダメージ増加をテストします。
func TestBattleEngine_ComboMaster_StackedDamage(t *testing.T) {
	// Arrange: ps_combo_masterの定義
	// ミスなし連続でダメージ+10%累積（最大+50%）
	comboMasterDef := domain.PassiveSkill{
		ID:          "ps_combo_master",
		Name:        "コンボマスター",
		TriggerType: domain.PassiveTriggerStack,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionNoMissStreak,
			Value: 1, // 1回以上のコンボで発動
		},
		EffectType:     domain.PassiveEffectMultiplier,
		EffectValue:    1.0, // ベース倍率（スタック0時）
		MaxStacks:      5,
		StackIncrement: 0.1, // スタックごとに+10%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_combo_master": comboMasterDef,
	}

	// エージェントを作成
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_combo_master", Name: "コンボマスター"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	moduleType := domain.SkillType{
		ID:   "test_attack",
		Name: "テスト攻撃",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
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
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 50000,
			Rank:   1,
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// コンボ0でのベースラインダメージ
	baselineResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}
	baselineDamage := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 0)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// コンボ3でのダメージ（+30% = 1.3倍）
	combo3Damage := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 3)

	// Assert: コンボ3で約1.3倍
	expectedRatio := 1.3
	actualRatio := float64(combo3Damage) / float64(baselineDamage)

	if math.Abs(actualRatio-expectedRatio) > 0.1 {
		t.Errorf("ps_combo_master: コンボ3でのダメージ倍率が期待値と異なります。got=%.2f, want=%.2f", actualRatio, expectedRatio)
	}
}

// TestBattleEngine_ComboMaster_MaxStacks はコンボ上限（5スタック）をテストします。
func TestBattleEngine_ComboMaster_MaxStacks(t *testing.T) {
	// Arrange
	comboMasterDef := domain.PassiveSkill{
		ID:          "ps_combo_master",
		Name:        "コンボマスター",
		TriggerType: domain.PassiveTriggerStack,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionNoMissStreak,
			Value: 1,
		},
		EffectType:     domain.PassiveEffectMultiplier,
		EffectValue:    1.0,
		MaxStacks:      5,
		StackIncrement: 0.1,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_combo_master": comboMasterDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_combo_master", Name: "コンボマスター"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	moduleType := domain.SkillType{
		ID:   "test_attack",
		Name: "テスト攻撃",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
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
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 50000,
			Rank:   1,
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	baselineResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// コンボ0でのベースライン
	baselineDamage := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 0)
	state.Enemy.HP = state.Enemy.MaxHP

	// コンボ5でのダメージ（最大+50% = 1.5倍）
	combo5Damage := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 5)
	state.Enemy.HP = state.Enemy.MaxHP

	// コンボ7でのダメージ（5で頭打ち、1.5倍のまま）
	combo7Damage := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 7)

	// Assert: コンボ5で1.5倍
	expectedRatio5 := 1.5
	actualRatio5 := float64(combo5Damage) / float64(baselineDamage)
	if math.Abs(actualRatio5-expectedRatio5) > 0.1 {
		t.Errorf("ps_combo_master: コンボ5でのダメージ倍率が期待値と異なります。got=%.2f, want=%.2f", actualRatio5, expectedRatio5)
	}

	// Assert: コンボ7でも1.5倍（キャップ）
	actualRatio7 := float64(combo7Damage) / float64(baselineDamage)
	if math.Abs(actualRatio7-expectedRatio5) > 0.1 {
		t.Errorf("ps_combo_master: コンボ7でのダメージ倍率が1.5でキャップされるべき。got=%.2f, want=%.2f", actualRatio7, expectedRatio5)
	}
}

// TestBattleEngine_ComboMaster_ZeroCombo はコンボ0で通常ダメージをテストします。
func TestBattleEngine_ComboMaster_ZeroCombo(t *testing.T) {
	// Arrange
	comboMasterDef := domain.PassiveSkill{
		ID:          "ps_combo_master",
		Name:        "コンボマスター",
		TriggerType: domain.PassiveTriggerStack,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionNoMissStreak,
			Value: 1,
		},
		EffectType:     domain.PassiveEffectMultiplier,
		EffectValue:    1.0,
		MaxStacks:      5,
		StackIncrement: 0.1,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_combo_master": comboMasterDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_combo_master", Name: "コンボマスター"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	moduleType := domain.SkillType{
		ID:   "test_attack",
		Name: "テスト攻撃",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
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
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 50000,
			Rank:   1,
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	baselineResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// コンボ0で2回実行
	damage1 := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 0)
	state.Enemy.HP = state.Enemy.MaxHP
	damage2 := engine.ApplySkillEffectWithCombo(state, agent, module, baselineResult, 0)

	// Assert: コンボ0では一貫したダメージ（倍率なし）
	ratio := float64(damage2) / float64(damage1)
	if math.Abs(ratio-1.0) > 0.05 {
		t.Errorf("ps_combo_master: コンボ0でダメージが変動してはいけません。ratio=%.2f", ratio)
	}
}
