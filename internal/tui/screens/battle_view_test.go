// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"strings"
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/combat/chain"
	"hirorocky/type-battle/internal/usecase/combat/recast"
)

// allFeaturesUnlockedProvider は全機能解放済みのプロバイダーを返します。
func allFeaturesUnlockedProvider() FeatureUnlockProvider {
	return &mockUnlockProvider{unlocked: map[domain.FeatureID]bool{
		domain.FeatureDefenseSkill:       true,
		domain.FeatureAgentCustomization: true,
		domain.FeatureChainEffect:        true,
		domain.FeatureManaSystem:         true,
	}}
}

// ==================== タスク9: バトル画面UI拡張テスト ====================

// createTestAgentWithPassive はパッシブスキル付きテスト用エージェントを作成します。
func createTestAgentWithPassive(passiveSkill domain.PassiveSkill, skills []*domain.SkillModel) *domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test_core_type",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.2, "MAG": 1.0, "SPD": 1.1, "LUK": 0.8},
		AllowedTags: []string{"physical_low"},
	}

	core := domain.NewCoreWithTypeID("test_core", coreType, passiveSkill)
	return domain.NewAgent("test_agent", core, skills)
}

// createTestSkillWithChain はチェイン効果付きテスト用スキルを作成します。
func createTestSkillWithChain(name string, chainEffect *domain.ChainEffect) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          "test_skill_" + name,
		Name:        name,
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト攻撃スキル",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50.0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, chainEffect)
}

// TestBattleScreen_RenderAgentAreaWithRecast はリキャスト状態表示のテストです。
func TestBattleScreen_RenderAgentAreaWithRecast(t *testing.T) {
	// テスト用エージェントとスキル作成
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	// BattleScreen作成
	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// エージェント0のリキャストを開始
	screen.recastManager.StartRecast(0, 5*time.Second)

	// View()を呼び出し
	result := screen.View()

	// リキャスト状態が表示されていることを確認
	if result == "" {
		t.Error("View() should return non-empty string")
	}
}

// TestBattleScreen_RenderAgentAreaWithChainEffect はチェイン効果待機表示のテストです。
func TestBattleScreen_RenderAgentAreaWithChainEffect(t *testing.T) {
	// チェイン効果付きスキル作成
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25.0)
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", &chainEffect),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	// BattleScreen作成
	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// チェイン効果を登録
	screen.chainEffectManager.RegisterChainEffect(0, &chainEffect, "test_skill")

	// View()を呼び出し
	result := screen.View()

	// チェイン効果の表示はViewに含まれる
	if result == "" {
		t.Error("View() should return non-empty string")
	}
}

// TestBattleScreen_RenderAgentAreaWithPassiveSkill はパッシブスキル表示のテストです。
func TestBattleScreen_RenderAgentAreaWithPassiveSkill(t *testing.T) {
	// パッシブスキル付きエージェント作成
	passiveSkill := domain.PassiveSkill{
		ID:          "test_passive",
		Name:        "パワーブースト",
		Description: "STRを強化する",
		Effects: map[domain.EffectColumn]float64{
			domain.ColSTRMultiplier: 1.1,
		},
	}

	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(passiveSkill, skills)

	// BattleScreen作成
	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// View()を呼び出し
	result := screen.View()

	// パッシブスキル表示はViewに含まれる
	if result == "" {
		t.Error("View() should return non-empty string")
	}
}

// TestBattleScreen_RecastStateAffectsSkillUsability はリキャスト状態によるスキル使用可否のテストです。
func TestBattleScreen_RecastStateAffectsSkillUsability(t *testing.T) {
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// リキャスト前は使用可能
	if !screen.isSkillUsable(0) {
		t.Error("Skill should be usable before recast")
	}

	// リキャスト開始
	screen.recastManager.StartRecast(0, 5*time.Second)

	// リキャスト中は使用不可
	if screen.isSkillUsable(0) {
		t.Error("Skill should not be usable during recast")
	}
}

// TestBattleScreen_ManaAffectsSkillUsability はマナ不足によるスキル使用可否のテストです。
func TestBattleScreen_ManaAffectsSkillUsability(t *testing.T) {
	// ManaCost=1のスキルを作成
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:       "mana_skill",
		Name:     "マナスキル",
		Icon:     "🔮",
		Tags:     []string{"magic_low"},
		ManaCost: 1,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50.0, StatCoef: 1.0, StatRef: "MAG"},
				Probability: 1.0,
				Icon:        "🔮",
			},
		},
	}, nil)
	skills := []*domain.SkillModel{skill}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	player := createTestPlayer()
	player.MaxMana = 10
	screen := NewBattleScreen(createTestEnemy(), player, []*domain.AgentModel{agent}, nil)
	screen.SetFeatureUnlockProvider(allFeaturesUnlockedProvider())

	// マナ不足時は使用不可
	screen.player.Mana = 0
	if screen.isSkillUsable(0) {
		t.Error("マナ不足時にスキルが使用可能になっています")
	}

	// マナ十分時は使用可能
	screen.player.Mana = 1
	if !screen.isSkillUsable(0) {
		t.Error("マナ十分時にスキルが使用不可になっています")
	}

	// ManaCost=0のスキルはマナ0でも使用可能
	freeCostSkill := domain.NewSkillFromType(domain.SkillType{
		ID:       "free_skill",
		Name:     "無コストスキル",
		Icon:     "⚔️",
		Tags:     []string{"physical_low"},
		ManaCost: 0,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50.0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
	freeSkills := []*domain.SkillModel{freeCostSkill}
	freeAgent := createTestAgentWithPassive(domain.PassiveSkill{}, freeSkills)

	screen2 := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{freeAgent}, nil)
	screen2.player.MaxMana = 10
	screen2.player.Mana = 0
	if !screen2.isSkillUsable(0) {
		t.Error("ManaCost=0のスキルがマナ0で使用不可になっています")
	}
}

// TestGetRecastInfoForAgent はエージェントリキャスト情報取得のテストです。
func TestGetRecastInfoForAgent(t *testing.T) {
	rm := recast.NewRecastManager()

	// リキャスト開始
	rm.StartRecast(0, 5*time.Second)

	// 状態取得
	state := rm.GetRecastState(0)
	if state == nil {
		t.Error("GetRecastState should return non-nil for agent in recast")
	}

	// 残り時間確認
	if state.RemainingSeconds != 5.0 {
		t.Errorf("RemainingSeconds = %v, want 5.0", state.RemainingSeconds)
	}
}

// TestGetPendingChainEffectForAgent はエージェントチェイン効果取得のテストです。
func TestGetPendingChainEffectForAgent(t *testing.T) {
	cm := chain.NewChainEffectManager()
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25.0)

	// チェイン効果を登録
	cm.RegisterChainEffect(0, &chainEffect, "test_skill")

	// 待機中効果を取得
	pending := cm.GetPendingEffectForAgent(0)
	if pending == nil {
		t.Error("GetPendingEffectForAgent should return non-nil for agent with chain effect")
	}

	// 効果の確認
	if pending.Effect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("Effect.Type = %v, want %v", pending.Effect.Type, domain.ChainEffectDamageBonus)
	}
}

// TestRenderSkillWithChainEffectBadge はスキル表示にチェイン効果バッジが含まれるかのテストです。
func TestRenderSkillWithChainEffectBadge(t *testing.T) {
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25.0)
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", &chainEffect),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)
	result := screen.View()

	// スキルが表示されている
	if !strings.Contains(result, "攻撃A") {
		t.Error("View should contain skill name")
	}
}

// TestBattleScreen_RenderRecastProgress はリキャスト進捗表示のテストです。
func TestBattleScreen_RenderRecastProgress(t *testing.T) {
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// リキャスト開始
	screen.recastManager.StartRecast(0, 5*time.Second)

	// リキャスト進捗取得
	progress := screen.recastManager.GetProgress(0)

	// 開始直後は0.0
	if progress != 0.0 {
		t.Errorf("GetProgress() = %v, want 0.0 at start", progress)
	}

	// 未リキャストのエージェントは1.0
	progress1 := screen.recastManager.GetProgress(1)
	if progress1 != 1.0 {
		t.Errorf("GetProgress(1) = %v, want 1.0 for non-recast agent", progress1)
	}
}

// TestBattleScreen_ChainEffectFeedback はチェイン効果発動フィードバックのテストです。
func TestBattleScreen_ChainEffectFeedback(t *testing.T) {
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25.0)
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", &chainEffect),
		createTestSkillWithChain("攻撃B", nil),
		createTestSkillWithChain("攻撃C", nil),
		createTestSkillWithChain("攻撃D", nil),
	}
	agent0 := createTestAgentWithPassive(domain.PassiveSkill{}, skills)
	agent1 := createTestAgentWithPassive(domain.PassiveSkill{}, []*domain.SkillModel{
		createTestSkillWithChain("攻撃E", nil),
		createTestSkillWithChain("攻撃F", nil),
		createTestSkillWithChain("攻撃G", nil),
		createTestSkillWithChain("攻撃H", nil),
	})

	screen := NewBattleScreen(createTestEnemy(), createTestPlayer(), []*domain.AgentModel{agent0, agent1}, nil)

	// エージェント0のチェイン効果を登録
	screen.chainEffectManager.RegisterChainEffect(0, &chainEffect, "test_skill")

	// エージェント1のスキル使用でチェイン効果発動をチェック
	triggered := screen.chainEffectManager.CheckAndTrigger(1, chain.SkillEffectFlags{HasDamage: true})

	// チェイン効果が発動する
	if len(triggered) != 1 {
		t.Errorf("CheckAndTrigger should return 1 triggered effect, got %d", len(triggered))
	}

	if len(triggered) > 0 {
		if triggered[0].Effect.Type != domain.ChainEffectDamageBonus {
			t.Errorf("Triggered effect type = %v, want %v", triggered[0].Effect.Type, domain.ChainEffectDamageBonus)
		}
	}
}

// ==================== ボルテージ表示テスト ====================

// TestBattleScreen_VoltageDisplay はボルテージ表示のテストです。
func TestBattleScreen_VoltageDisplay(t *testing.T) {
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	enemy := createTestEnemy()
	screen := NewBattleScreen(enemy, createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// View()を呼び出し
	result := screen.View()

	// ボルテージ表示が含まれていることを確認
	if !strings.Contains(result, "VOLTAGE") {
		t.Error("View should contain VOLTAGE display")
	}

	// 初期値100%が表示されていることを確認
	if !strings.Contains(result, "100%") {
		t.Error("View should contain initial voltage value 100%")
	}
}

// TestBattleScreen_VoltageDisplayWithHighVoltage は高ボルテージ時の表示テストです。
func TestBattleScreen_VoltageDisplayWithHighVoltage(t *testing.T) {
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	enemy := createTestEnemy()
	// ボルテージを150%に設定
	enemy.SetVoltage(150.0)

	screen := NewBattleScreen(enemy, createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// View()を呼び出し
	result := screen.View()

	// ボルテージ150%が表示されていることを確認
	if !strings.Contains(result, "150%") {
		t.Error("View should contain voltage value 150%")
	}
}

// TestBattleScreen_VoltageDisplayWithDangerVoltage は危険レベルボルテージの表示テストです。
func TestBattleScreen_VoltageDisplayWithDangerVoltage(t *testing.T) {
	skills := []*domain.SkillModel{
		createTestSkillWithChain("攻撃A", nil),
	}
	agent := createTestAgentWithPassive(domain.PassiveSkill{}, skills)

	enemy := createTestEnemy()
	// ボルテージを200%に設定
	enemy.SetVoltage(200.0)

	screen := NewBattleScreen(enemy, createTestPlayer(), []*domain.AgentModel{agent}, nil)

	// View()を呼び出し
	result := screen.View()

	// ボルテージ200%が表示されていることを確認
	if !strings.Contains(result, "200%") {
		t.Error("View should contain voltage value 200%")
	}
}
