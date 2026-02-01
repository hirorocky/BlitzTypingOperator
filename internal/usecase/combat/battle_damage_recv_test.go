// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ===== Phase 2: 被ダメージ時イベント発火のテスト =====

// TestBattleEngine_DamageRecv_LastStand はHP25%以下で被ダメージを1に固定するテストです。
func TestBattleEngine_DamageRecv_LastStand(t *testing.T) {
	// Arrange: ps_last_standの定義
	// HP25%以下で30%の確率で被ダメージ1
	lastStandDef := domain.PassiveSkill{
		ID:          "ps_last_stand",
		Name:        "ラストスタンド",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionHPBelowPercent,
			Value: 25,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1,   // ダメージを1に固定
		Probability: 1.0, // テスト用に100%で発動
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_last_stand": lastStandDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_last_stand", Name: "ラストスタンド"}
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
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 1.0, StatRef: "STR"},
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
			BaseAttackPower: 50, // 大きなダメージ
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	// テスト用に乱数シードを固定
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// プレイヤーHPを25%以下に設定
	state.Player.HP = 20 // 20/100 = 20%

	// Act: 敵の攻撃を処理
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// Assert: ダメージが1に固定されている
	if damage != 1 {
		t.Errorf("ps_last_stand: HP25%%以下でダメージが1になるべき。got=%d, want=1", damage)
	}
}

// TestBattleEngine_DamageRecv_LastStand_HPAbove25 はHP25%以上では通常ダメージのテストです。
func TestBattleEngine_DamageRecv_LastStand_HPAbove25(t *testing.T) {
	// Arrange
	lastStandDef := domain.PassiveSkill{
		ID:          "ps_last_stand",
		Name:        "ラストスタンド",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionHPBelowPercent,
			Value: 25,
		},
		EffectType:  domain.PassiveEffectSpecial,
		EffectValue: 1,
		Probability: 1.0,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_last_stand": lastStandDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_last_stand", Name: "ラストスタンド"}
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
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 1.0, StatRef: "STR"},
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
			BaseAttackPower: 50,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// プレイヤーHPを50%に設定（25%以上）
	state.Player.HP = 500 // 500/1000 = 50%

	// Act
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// Assert: 通常ダメージ（1より大きい）
	if damage <= 1 {
		t.Errorf("ps_last_stand: HP25%%以上では通常ダメージになるべき。got=%d", damage)
	}
}

// TestBattleEngine_DamageRecv_CounterCharge は被ダメージ時に次攻撃2倍バフのテストです。
func TestBattleEngine_DamageRecv_CounterCharge(t *testing.T) {
	// Arrange: ps_counter_chargeの定義
	// 被ダメージ時20%で次の攻撃2倍
	counterChargeDef := domain.PassiveSkill{
		ID:          "ps_counter_charge",
		Name:        "カウンターチャージ",
		TriggerType: domain.PassiveTriggerProbability,
		TriggerCondition: &domain.TriggerCondition{
			Type: domain.TriggerConditionOnDamageReceived,
		},
		EffectType:  domain.PassiveEffectMultiplier,
		EffectValue: 2.0, // 次攻撃2倍
		Probability: 1.0, // テスト用に100%で発動
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_counter_charge": counterChargeDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_counter_charge", Name: "カウンターチャージ"}
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
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 1.0, StatRef: "STR"},
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

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// Act: 敵の攻撃を処理
	engine.ProcessEnemyAttackDamage(state, "physical")

	// Assert: プレイヤーのEffectTableに「次攻撃2倍」バフが追加されている
	ctx := domain.NewEffectContext(state.Player.HP, state.Player.MaxHP, state.Enemy.HP, state.Enemy.MaxHP)
	effects := state.Player.EffectTable.Aggregate(ctx)

	if effects.DamageMultiplier < 1.9 {
		t.Errorf("ps_counter_charge: 被ダメージ後にDamageMultiplierが2.0になるべき。got=%.2f", effects.DamageMultiplier)
	}
}
