// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ===== Phase 4: バトル開始・タイムアウト・ミスイベントのテスト =====

// TestBattleEngine_BattleStart_FirstStrike はバトル開始時の即発動をテストします。
func TestBattleEngine_BattleStart_FirstStrike(t *testing.T) {
	// Arrange: ps_first_strikeの定義
	firstStrikeDef := domain.PassiveSkill{
		ID:          "ps_first_strike",
		Name:        "ファーストストライク",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnBattleStart,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1.0, // 最初のスキル即発動
		Probability: 1.0, // テスト用に100%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_first_strike": firstStrikeDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_first_strike", Name: "ファーストストライク"}
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
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Act: バトル開始時のファーストストライク判定
	isFirstStrike := engine.EvaluateFirstStrike(state, agent)

	// Assert: ファーストストライク発動
	if !isFirstStrike {
		t.Errorf("ps_first_strike: ファーストストライクが発動していない")
	}
}

// TestBattleEngine_BattleStart_FirstStrike_NotEquipped は装備していない場合発動しないことをテストします。
func TestBattleEngine_BattleStart_FirstStrike_NotEquipped(t *testing.T) {
	// Arrange: パッシブなしのエージェント
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_other", Name: "その他"}
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
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)

	// Act
	isFirstStrike := engine.EvaluateFirstStrike(state, agent)

	// Assert: 発動しない
	if isFirstStrike {
		t.Errorf("ps_first_strike: パッシブなしで発動してはいけない")
	}
}

// TestBattleEngine_TypoRecovery はミス時の時間延長をテストします。
func TestBattleEngine_TypoRecovery(t *testing.T) {
	// Arrange: ps_typo_recoveryの定義
	typoRecoveryDef := domain.PassiveSkill{
		ID:          "ps_typo_recovery",
		Name:        "タイポリカバリー",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnTypingMiss,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1.0, // +1秒
		Probability: 1.0, // テスト用に100%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_typo_recovery": typoRecoveryDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_typo_recovery", Name: "タイポリカバリー"}
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
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Act: タイポリカバリー発動チェック
	timeExtension := engine.EvaluateTypoRecovery(state, agent)

	// Assert: +1秒の延長
	if timeExtension != 1.0 {
		t.Errorf("ps_typo_recovery: timeExtension=%.1f, want 1.0", timeExtension)
	}
}

// TestBattleEngine_SecondChance はタイムアウト時の再挑戦をテストします。
func TestBattleEngine_SecondChance(t *testing.T) {
	// Arrange: ps_second_chanceの定義
	secondChanceDef := domain.PassiveSkill{
		ID:          "ps_second_chance",
		Name:        "セカンドチャンス",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnTimeout,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1.0, // 再挑戦
		Probability: 1.0, // テスト用に100%
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_second_chance": secondChanceDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_second_chance", Name: "セカンドチャンス"}
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
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Act: セカンドチャンス発動チェック
	isSecondChance := engine.EvaluateSecondChance(state, agent)

	// Assert: 発動
	if !isSecondChance {
		t.Errorf("ps_second_chance: セカンドチャンスが発動していない")
	}
}
