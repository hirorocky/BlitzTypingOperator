// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// ===== Phase 1: タイピング完了時イベント発火のテスト =====

// TestBattleEngine_TypingDone_PerfectRhythm は正確性100%でps_perfect_rhythmが発動することをテストします。
func TestBattleEngine_TypingDone_PerfectRhythm(t *testing.T) {
	// Arrange: パッシブスキル定義を作成
	perfectRhythmDef := domain.PassiveSkill{
		ID:          "ps_perfect_rhythm",
		Name:        "パーフェクトリズム",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_perfect_rhythm": perfectRhythmDef,
	}

	// エージェントを作成（パッシブスキル付き）
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_perfect_rhythm", Name: "パーフェクトリズム"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	skillType := domain.SkillType{
		ID:          "test_attack",
		Name:        "テスト攻撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト用攻撃",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}
	skill := domain.NewSkillFromType(skillType, nil)
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	// 敵タイプを作成
	enemyTypes := []domain.EnemyType{
		{
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 10000,
		},
	}

	// BattleEngineを作成
	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)

	// バトルを初期化
	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("InitializeBattle failed: %v", err)
	}
	engine.RegisterPassiveSkills(state, agents)

	// ベースラインダメージを計算（パッシブなしの条件：Score=95）
	baselineResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     95, // パッシブ発動しない
	}
	baselineResult_ := engine.ApplySkillEffect(state, agent, skill, baselineResult)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// パーフェクトタイピング結果（パッシブ発動する）
	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
		IsPerfect: true,
	}
	result := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ダメージが約1.5倍になっていることを確認
	expectedRatio := 1.5
	actualRatio := float64(result.TotalDamage) / float64(baselineResult_.TotalDamage)

	if math.Abs(actualRatio-expectedRatio) > 0.1 {
		t.Errorf("ps_perfect_rhythm: ダメージ倍率が期待値と異なります。got=%.2f, want=%.2f (±0.1)", actualRatio, expectedRatio)
	}
}

// TestBattleEngine_TypingDone_PerfectRhythm_NotTriggered は正確性100%未満でps_perfect_rhythmが発動しないことをテストします。
func TestBattleEngine_TypingDone_PerfectRhythm_NotTriggered(t *testing.T) {
	// Arrange
	perfectRhythmDef := domain.PassiveSkill{
		ID:          "ps_perfect_rhythm",
		Name:        "パーフェクトリズム",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_perfect_rhythm": perfectRhythmDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_perfect_rhythm", Name: "パーフェクトリズム"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	skillType := domain.SkillType{
		ID:          "test_attack",
		Name:        "テスト攻撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト用攻撃",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}
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
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Score=95のタイピング結果（100未満）
	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     95, // 95%
	}
	result1 := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// 同じ条件でもう一度（パッシブなしの一貫性確認）
	result2 := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ダメージが1.5倍にならない（ほぼ同じダメージ）
	ratio := float64(result2.TotalDamage) / float64(result1.TotalDamage)
	if math.Abs(ratio-1.0) > 0.05 {
		t.Errorf("ps_perfect_rhythm: 正確性100%%未満で倍率が変わってはいけません。ratio=%.2f", ratio)
	}
}

// TestBattleEngine_TypingDone_SpeedBreak はWPM80以上でps_speed_breakが発動することをテストします。
func TestBattleEngine_TypingDone_SpeedBreak(t *testing.T) {
	// Arrange
	speedBreakDef := domain.PassiveSkill{
		ID:          "ps_speed_break",
		Name:        "スピードブレイク",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionWPMAbove,
			Value: 80,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.25,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_speed_break": speedBreakDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_speed_break", Name: "スピードブレイク"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	skillType := domain.SkillType{
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
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// ベースラインダメージを計算（WPM70、パッシブ発動しない）
	baselineResult := &typing.TypingResult{
		Completed: true,
		WPM:       70.0, // 80未満
		Score:     100,
	}
	baselineResult_ := engine.ApplySkillEffect(state, agent, skill, baselineResult)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// WPM=85のタイピング結果（パッシブ発動）
	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       85.0, // 80以上
		Score:     100,
	}
	result := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ダメージが約1.25倍
	expectedRatio := 1.25
	actualRatio := float64(result.TotalDamage) / float64(baselineResult_.TotalDamage)

	if math.Abs(actualRatio-expectedRatio) > 0.1 {
		t.Errorf("ps_speed_break: ダメージ倍率が期待値と異なります。got=%.2f, want=%.2f (±0.1)", actualRatio, expectedRatio)
	}
}

// TestBattleEngine_TypingDone_SpeedBreak_NotTriggered はWPM80未満でps_speed_breakが発動しないことをテストします。
func TestBattleEngine_TypingDone_SpeedBreak_NotTriggered(t *testing.T) {
	// Arrange
	speedBreakDef := domain.PassiveSkill{
		ID:          "ps_speed_break",
		Name:        "スピードブレイク",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionWPMAbove,
			Value: 80,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.25,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_speed_break": speedBreakDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_speed_break", Name: "スピードブレイク"}
	core := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill)

	skillType := domain.SkillType{
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
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// WPM=70のタイピング結果（80未満）
	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       70.0, // 80未満
		Score:     100,
	}
	result1 := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// 同じ条件でもう一度
	result2 := engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ダメージが1.25倍にならない
	ratio := float64(result2.TotalDamage) / float64(result1.TotalDamage)
	if math.Abs(ratio-1.0) > 0.05 {
		t.Errorf("ps_speed_break: WPM80未満で倍率が変わってはいけません。ratio=%.2f", ratio)
	}
}

// TestBattleEngine_TypingDone_Combined は両方の条件を満たした場合、両方のパッシブが発動することをテストします。
func TestBattleEngine_TypingDone_Combined(t *testing.T) {
	// Arrange: 両方のパッシブスキルを定義
	perfectRhythmDef := domain.PassiveSkill{
		ID:          "ps_perfect_rhythm",
		Name:        "パーフェクトリズム",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	speedBreakDef := domain.PassiveSkill{
		ID:          "ps_speed_break",
		Name:        "スピードブレイク",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionWPMAbove,
			Value: 80,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 1.25,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_perfect_rhythm": perfectRhythmDef,
		"ps_speed_break":    speedBreakDef,
	}

	// 2つのエージェントを作成（それぞれ異なるパッシブスキル）
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}

	passiveSkill1 := domain.PassiveSkill{ID: "ps_perfect_rhythm", Name: "パーフェクトリズム"}
	core1 := domain.NewCoreWithTypeID("core_001", coreType, passiveSkill1)

	passiveSkill2 := domain.PassiveSkill{ID: "ps_speed_break", Name: "スピードブレイク"}
	core2 := domain.NewCoreWithTypeID("core_002", coreType, passiveSkill2)

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent1 := domain.NewAgent("agent_001", core1, []*domain.SkillModel{skill})
	agent2 := domain.NewAgent("agent_002", core2, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent1, agent2}

	enemyTypes := []domain.EnemyType{
		{
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 10000,
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// ベースラインダメージ（パッシブなし: WPM=70, Score=95）
	baselineResult := &typing.TypingResult{
		Completed: true,
		WPM:       70.0,
		Score:     95,
	}
	baselineResult_ := engine.ApplySkillEffect(state, agent1, skill, baselineResult)

	// 敵HPをリセット
	state.Enemy.HP = state.Enemy.MaxHP

	// 両方の条件を満たすタイピング結果（パーフェクト, WPM=85）
	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       85.0, // 80以上
		Score:     100,
		IsPerfect: true,
	}
	result := engine.ApplySkillEffect(state, agent1, skill, typingResult)

	// Assert: 両方の効果が適用される（1.5 * 1.25 = 1.875倍）
	expectedRatio := 1.5 * 1.25 // 1.875
	actualRatio := float64(result.TotalDamage) / float64(baselineResult_.TotalDamage)

	if math.Abs(actualRatio-expectedRatio) > 0.15 {
		t.Errorf("両パッシブ発動: ダメージ倍率が期待値と異なります。got=%.2f, want=%.2f (±0.15)", actualRatio, expectedRatio)
	}
}
