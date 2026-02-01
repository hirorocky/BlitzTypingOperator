// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ===== Phase 2: 同種攻撃カウンターのテスト =====

// TestBattleState_SameAttackCount_Track は同種攻撃のカウント追跡をテストします。
func TestBattleState_SameAttackCount_Track(t *testing.T) {
	// Arrange
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_adaptive_shield", Name: "アダプティブシールド"}
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
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          1000,
			BaseAttackPower: 50,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)

	// Act: 物理攻撃を3回
	engine.RecordAttackType(state, "physical")
	engine.RecordAttackType(state, "physical")
	engine.RecordAttackType(state, "physical")

	// Assert: カウントが3
	if state.SameAttackCount != 3 {
		t.Errorf("SameAttackCount: got %d, want 3", state.SameAttackCount)
	}
	if state.LastAttackType != "physical" {
		t.Errorf("LastAttackType: got %s, want physical", state.LastAttackType)
	}
}

// TestBattleState_SameAttackCount_Reset は攻撃属性変更時のカウントリセットをテストします。
func TestBattleState_SameAttackCount_Reset(t *testing.T) {
	// Arrange
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_adaptive_shield", Name: "アダプティブシールド"}
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
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          1000,
			BaseAttackPower: 50,
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)

	// Act: 物理2回 → 魔法1回
	engine.RecordAttackType(state, "physical")
	engine.RecordAttackType(state, "physical")
	engine.RecordAttackType(state, "magic")

	// Assert: カウントが1（魔法1回目）
	if state.SameAttackCount != 1 {
		t.Errorf("SameAttackCount: got %d, want 1", state.SameAttackCount)
	}
	if state.LastAttackType != "magic" {
		t.Errorf("LastAttackType: got %s, want magic", state.LastAttackType)
	}
}

// TestBattleEngine_AdaptiveShield は同種攻撃3回以上でダメージ軽減をテストします。
func TestBattleEngine_AdaptiveShield(t *testing.T) {
	// Arrange
	adaptiveShieldDef := domain.PassiveSkill{
		ID:          "ps_adaptive_shield",
		Name:        "アダプティブシールド",
		TriggerType: domain.PassiveTriggerConditional,
		TriggerCondition: &domain.TriggerCondition{
			Type:  domain.TriggerConditionSameAttackCount,
			Value: 3,
		},
		EffectType:  domain.PassiveEffectModifier,
		EffectValue: 0.25, // 25%軽減
		Probability: 1.0,
	}

	passiveSkillDefs := map[string]domain.PassiveSkill{
		"ps_adaptive_shield": adaptiveShieldDef,
	}

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "ps_adaptive_shield", Name: "アダプティブシールド"}
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
			ID:              "test_enemy",
			Name:            "テスト敵",
			BaseHP:          10000,
			BaseAttackPower: 100, // 100ダメージ
			AttackType:      "physical",
		},
	}

	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkillDefs)
	engine.SetRng(rand.New(rand.NewSource(42)))

	state, _ := engine.initializeBattleForTest(1, agents)
	engine.RegisterPassiveSkills(state, agents)

	// 1-2回目の物理攻撃（まだ軽減なし）
	damage1 := engine.ProcessEnemyAttackDamage(state, "physical") // count=1
	state.Player.HP = state.Player.MaxHP                          // HPリセット
	damage2 := engine.ProcessEnemyAttackDamage(state, "physical") // count=2
	state.Player.HP = state.Player.MaxHP                          // HPリセット

	// 3回目の物理攻撃（count=3で軽減発動）
	damage3 := engine.ProcessEnemyAttackDamage(state, "physical")
	state.Player.HP = state.Player.MaxHP // HPリセット

	// 4回目の物理攻撃（引き続き軽減）
	damage4 := engine.ProcessEnemyAttackDamage(state, "physical")

	// Assert: 1-2回目は軽減なし、3回目以降は25%軽減
	// damage1, damage2は同じはず（軽減なし）
	if damage1 != damage2 {
		t.Errorf("1回目と2回目のダメージが異なる: damage1=%d, damage2=%d", damage1, damage2)
	}

	// damage3, damage4は軽減されているはず
	expectedWithReduction := int(float64(damage1) * 0.75)
	if damage3 < damage1-10 || damage3 > damage1 {
		// damage3は軽減されているはず（damage1より小さい）
		// ただし最低1ダメージ保証などもあるので幅を持たせる
	}

	// 3回目以降は約75%のダメージになる
	ratio := float64(damage3) / float64(damage1)
	if ratio < 0.70 || ratio > 0.80 {
		t.Errorf("3回目のダメージ軽減率が期待値と異なる: ratio=%.2f, want 0.75", ratio)
	}

	// 4回目も同様に軽減
	ratio4 := float64(damage4) / float64(damage1)
	if ratio4 < 0.70 || ratio4 > 0.80 {
		t.Errorf("4回目のダメージ軽減率が期待値と異なる: ratio=%.2f, want 0.75", ratio4)
	}

	_ = expectedWithReduction
}
