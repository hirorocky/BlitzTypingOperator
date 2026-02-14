// Package combat はバトルエンジンを提供します。
// ボルテージシステム統合テスト

package combat

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// TestBattleEngine_VoltageInitialization はバトル初期化時のボルテージ初期化をテストします。
func TestBattleEngine_VoltageInitialization(t *testing.T) {
	// 敵タイプを定義
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 100,
			VoltageRisePer10s:      20.0,
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	// テスト用エージェント
	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 1.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ボルテージが100%で初期化されていることを確認
	expectedVoltage := 100.0
	if state.Enemy.GetVoltage() != expectedVoltage {
		t.Errorf("expected voltage %.1f, got %.1f", expectedVoltage, state.Enemy.GetVoltage())
	}
}

// TestBattleEngine_UpdateEffects_VoltageRise はUpdateEffectsでボルテージが上昇することをテストします。
func TestBattleEngine_UpdateEffects_VoltageRise(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 100,
			VoltageRisePer10s:      20.0, // 10秒で20ポイント上昇
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 1.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 5秒経過を更新（10ポイント上昇するはず）
	engine.UpdateEffects(state, 5.0)

	expectedVoltage := 110.0
	if state.Enemy.GetVoltage() != expectedVoltage {
		t.Errorf("expected voltage %.1f, got %.1f", expectedVoltage, state.Enemy.GetVoltage())
	}
}

// TestBattleEngine_UpdateEffects_VoltageContinuesOnPhaseTransition はフェーズ遷移時もボルテージが継続することをテストします。
func TestBattleEngine_UpdateEffects_VoltageContinuesOnPhaseTransition(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 100,
			VoltageRisePer10s:      20.0,
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 1.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ボルテージを150%に設定
	state.Enemy.SetVoltage(150.0)

	// 敵をダメージで強化フェーズに遷移させる（HP50%以下に）
	state.Enemy.TakeDamage(state.Enemy.MaxHP / 2)

	// フェーズ遷移を確認
	phaseChanged := engine.CheckPhaseTransition(state)
	if !phaseChanged {
		t.Error("expected phase transition")
	}

	// フェーズ遷移後もボルテージが150%のまま
	expectedVoltage := 150.0
	if state.Enemy.GetVoltage() != expectedVoltage {
		t.Errorf("expected voltage %.1f after phase transition, got %.1f", expectedVoltage, state.Enemy.GetVoltage())
	}
}

// TestBattleEngine_UpdateEffects_VoltageSmallIncrement は0.1秒単位のボルテージ上昇をテストします。
func TestBattleEngine_UpdateEffects_VoltageSmallIncrement(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 100,
			VoltageRisePer10s:      10.0, // 10秒で10ポイント = 1秒で1ポイント = 0.1秒で0.1ポイント
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 1.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 0.1秒経過を更新（BattleTickInterval相当）
	engine.UpdateEffects(state, 0.1)

	expectedVoltage := 100.1
	if state.Enemy.GetVoltage() != expectedVoltage {
		t.Errorf("expected voltage %.1f, got %.1f", expectedVoltage, state.Enemy.GetVoltage())
	}
}

// TestBattleEngine_UpdateEffects_VoltageZeroRise は上昇率0の場合をテストします。
func TestBattleEngine_UpdateEffects_VoltageZeroRise(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 100,
			VoltageRisePer10s:      0.0, // 上昇しない設定
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 1.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 10秒経過
	engine.UpdateEffects(state, 10.0)

	// ボルテージが100%のまま
	expectedVoltage := 100.0
	if state.Enemy.GetVoltage() != expectedVoltage {
		t.Errorf("expected voltage %.1f, got %.1f", expectedVoltage, state.Enemy.GetVoltage())
	}
}

// TestBattleEngine_VoltageManager_NotNil はVoltageManagerが設定されていることをテストします。
func TestBattleEngine_VoltageManager_NotNil(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "test_enemy",
			Name:   "テスト敵",
			BaseHP: 100,
		},
	}

	engine := NewBattleEngine(enemyTypes)

	if engine.voltageManager == nil {
		t.Error("expected voltageManager to be initialized")
	}
}

// TestBattleEngine_VoltageDamageMultiplier_100Percent はボルテージ100%で敵攻撃が等倍ダメージをテストします。
func TestBattleEngine_VoltageDamageMultiplier_100Percent(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 1000,
			VoltageRisePer10s:      0.0,
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second, DamageBase: 10, DamagePerLevel: 2},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 10.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ボルテージ100%（等倍）
	initialPlayerHP := state.Player.HP
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// 敵攻撃力 = DamageBase(10) + level(1) * DamagePerLevel(2) = 12
	// ボルテージ100% -> 乗算 x1.0 -> 最終ダメージ = 12
	expectedDamage := 12
	if damage != expectedDamage {
		t.Errorf("expected damage %d, got %d", expectedDamage, damage)
	}

	actualDamage := initialPlayerHP - state.Player.HP
	if actualDamage != expectedDamage {
		t.Errorf("expected HP reduction %d, got %d", expectedDamage, actualDamage)
	}
}

// TestBattleEngine_VoltageDamageMultiplier_150Percent はボルテージ150%で敵攻撃が1.5倍ダメージをテストします。
func TestBattleEngine_VoltageDamageMultiplier_150Percent(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 1000,
			VoltageRisePer10s:      0.0,
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second, DamageBase: 10, DamagePerLevel: 2},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 10.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ボルテージを150%に設定
	state.Enemy.SetVoltage(150.0)

	initialPlayerHP := state.Player.HP
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// 敵攻撃力 = DamageBase(10) + level(1) * DamagePerLevel(2) = 12
	// ボルテージ150% -> 乗算 x1.5 -> 最終ダメージ = 18
	expectedDamage := 18
	if damage != expectedDamage {
		t.Errorf("expected damage %d, got %d", expectedDamage, damage)
	}

	actualDamage := initialPlayerHP - state.Player.HP
	if actualDamage != expectedDamage {
		t.Errorf("expected HP reduction %d, got %d", expectedDamage, actualDamage)
	}
}

// TestBattleEngine_VoltageDamageMultiplier_200Percent はボルテージ200%で敵攻撃が2倍ダメージをテストします。
func TestBattleEngine_VoltageDamageMultiplier_200Percent(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                     "test_enemy",
			Name:                   "テスト敵",
			BaseHP:                 1000,
			VoltageRisePer10s:      0.0,
			NormalActionPatternIDs: []string{},
			ResolvedNormalActions: []domain.EnemyAction{
				{ID: "attack_1", Name: "攻撃", ActionType: domain.EnemyActionAttack, ChargeTime: 2 * time.Second, DamageBase: 10, DamagePerLevel: 2},
			},
		},
	}

	engine := NewBattleEngine(enemyTypes)

	core := &domain.CoreModel{
		ID:    "test_core",
		Name:  "テストコア",
		Stats: domain.Stats{STR: 10, INT: 10, WIL: 10, LUK: 10},
	}
	skill := newTestDamageSkill("m1", "ダメージスキル", []string{"physical"}, 10.0, "STR", "テスト")
	agent := domain.NewAgent("agent_001", core, []*domain.SkillModel{skill})
	agents := []*domain.AgentModel{agent}

	state, err := engine.initializeBattleForTest(1, agents)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ボルテージを200%に設定
	state.Enemy.SetVoltage(200.0)

	initialPlayerHP := state.Player.HP
	damage := engine.ProcessEnemyAttackDamage(state, "physical")

	// 敵攻撃力 = DamageBase(10) + level(1) * DamagePerLevel(2) = 12
	// ボルテージ200% -> 乗算 x2.0 -> 最終ダメージ = 24
	expectedDamage := 24
	if damage != expectedDamage {
		t.Errorf("expected damage %d, got %d", expectedDamage, damage)
	}

	actualDamage := initialPlayerHP - state.Player.HP
	if actualDamage != expectedDamage {
		t.Errorf("expected HP reduction %d, got %d", expectedDamage, actualDamage)
	}
}
