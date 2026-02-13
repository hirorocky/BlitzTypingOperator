// Package chain はチェイン効果管理機能を提供します。
package chain

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestNewChainEffectManager はChainEffectManagerの生成をテストします。
func TestNewChainEffectManager(t *testing.T) {
	cem := NewChainEffectManager()
	if cem == nil {
		t.Fatal("ChainEffectManagerがnilです")
	}

	// 初期状態では待機中効果がない
	pending := cem.GetPendingEffects()
	if len(pending) != 0 {
		t.Errorf("初期状態の待機中効果数: got %d, want 0", len(pending))
	}
}

// TestRegisterChainEffect はチェイン効果登録をテストします。
func TestRegisterChainEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// チェイン効果を作成
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)

	// 登録
	cem.RegisterChainEffect(0, &effect, "slash_lv1")

	// 待機中効果を確認
	pending := cem.GetPendingEffects()
	if len(pending) != 1 {
		t.Fatalf("待機中効果数: got %d, want 1", len(pending))
	}

	pe := pending[0]
	if pe.AgentIndex != 0 {
		t.Errorf("AgentIndex: got %d, want 0", pe.AgentIndex)
	}
	if pe.Effect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("Effect.Type: got %v, want %v", pe.Effect.Type, domain.ChainEffectDamageBonus)
	}
	if pe.Effect.Value != 20.0 {
		t.Errorf("Effect.Value: got %f, want 20.0", pe.Effect.Value)
	}
	if pe.SourceModule != "slash_lv1" {
		t.Errorf("SourceModule: got %s, want slash_lv1", pe.SourceModule)
	}
}

// TestCheckAndTrigger は他エージェントモジュール使用時の発動をテストします。
func TestCheckAndTrigger(t *testing.T) {
	cem := NewChainEffectManager()

	// エージェント0のチェイン効果を登録（説明文テンプレート付き）
	effect := domain.NewChainEffectWithTemplate("test_damage_bonus", domain.ChainEffectDamageBonus, 25.0,
		"次の攻撃のダメージ+%.0f%%", "次攻撃ダメ+%.0f%%")
	cem.RegisterChainEffect(0, &effect, "slash_lv1")

	// エージェント1がモジュールを使用（チェイン効果が発動）
	triggered := cem.CheckAndTrigger(1, SkillEffectFlags{HasDamage: true})

	// 発動した効果を確認
	if len(triggered) != 1 {
		t.Fatalf("発動した効果数: got %d, want 1", len(triggered))
	}

	te := triggered[0]
	if te.Effect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("Effect.Type: got %v, want %v", te.Effect.Type, domain.ChainEffectDamageBonus)
	}
	if te.EffectValue != 25.0 {
		t.Errorf("EffectValue: got %f, want 25.0", te.EffectValue)
	}
	if te.Message == "" {
		t.Error("Messageが空です")
	}

	// 発動後は待機中効果から削除されている
	pending := cem.GetPendingEffects()
	if len(pending) != 0 {
		t.Errorf("発動後の待機中効果数: got %d, want 0", len(pending))
	}
}

// TestCheckAndTriggerSameAgent は同一エージェントでは発動しないことをテストします。
func TestCheckAndTriggerSameAgent(t *testing.T) {
	cem := NewChainEffectManager()

	// エージェント0のチェイン効果を登録
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25.0)
	cem.RegisterChainEffect(0, &effect, "slash_lv1")

	// 同じエージェント0がモジュールを使用（発動しない）
	triggered := cem.CheckAndTrigger(0, SkillEffectFlags{HasDamage: true})

	// 発動しない
	if len(triggered) != 0 {
		t.Errorf("同一エージェントで発動した効果数: got %d, want 0", len(triggered))
	}

	// 待機中効果は残っている
	pending := cem.GetPendingEffects()
	if len(pending) != 1 {
		t.Errorf("待機中効果数: got %d, want 1", len(pending))
	}
}

// TestClearEffectForAgent は指定エージェントの待機中チェイン効果を削除することをテストします。
func TestClearEffectForAgent(t *testing.T) {
	cem := NewChainEffectManager()

	// 複数エージェントのチェイン効果を登録
	effect1 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)
	effect2 := domain.NewChainEffect("test_effect", domain.ChainEffectHealBonus, 30.0)
	cem.RegisterChainEffect(0, &effect1, "slash_lv1")
	cem.RegisterChainEffect(1, &effect2, "heal_lv1")

	// エージェント0の効果を削除
	cem.ClearEffectForAgent(0)

	// エージェント0の効果が削除されている
	if cem.HasPendingEffect(0) {
		t.Error("エージェント0の待機中効果が削除されていません")
	}

	// エージェント1の効果は残っている
	pending := cem.GetPendingEffects()
	if len(pending) != 1 {
		t.Fatalf("残り待機中効果数: got %d, want 1", len(pending))
	}
	if pending[0].AgentIndex != 1 {
		t.Errorf("残りAgentIndex: got %d, want 1", pending[0].AgentIndex)
	}
}

// TestClearEffectForAgentNoEffect は効果がないエージェントに対するClearEffectForAgentが安全に動作することをテストします。
func TestClearEffectForAgentNoEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// 存在しないエージェントに対してクリア（エラーにならない）
	cem.ClearEffectForAgent(0)

	// 他の効果に影響しない
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)
	cem.RegisterChainEffect(1, &effect, "slash_lv1")
	cem.ClearEffectForAgent(0)

	if !cem.HasPendingEffect(1) {
		t.Error("エージェント1の待機中効果が消えてしまいました")
	}
}

// TestMultipleEffectsFromSameAgent は同一エージェントの複数効果をテストします。
func TestMultipleEffectsFromSameAgent(t *testing.T) {
	cem := NewChainEffectManager()

	// 同一エージェントから複数効果を登録（通常は1つだが、上書きされる）
	effect1 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)
	effect2 := domain.NewChainEffect("test_effect", domain.ChainEffectHealBonus, 30.0)
	cem.RegisterChainEffect(0, &effect1, "slash_lv1")
	cem.RegisterChainEffect(0, &effect2, "heal_lv1") // 上書き

	// 最新の効果のみ保持
	pending := cem.GetPendingEffects()
	if len(pending) != 1 {
		t.Fatalf("待機中効果数: got %d, want 1", len(pending))
	}
	if pending[0].Effect.Type != domain.ChainEffectHealBonus {
		t.Errorf("効果タイプ: got %v, want %v", pending[0].Effect.Type, domain.ChainEffectHealBonus)
	}
}

// TestClearAll は全効果クリアをテストします。
func TestClearAll(t *testing.T) {
	cem := NewChainEffectManager()

	// 複数効果を登録
	effect1 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)
	effect2 := domain.NewChainEffect("test_effect", domain.ChainEffectHealBonus, 30.0)
	cem.RegisterChainEffect(0, &effect1, "slash_lv1")
	cem.RegisterChainEffect(1, &effect2, "heal_lv1")

	// 全クリア
	cem.ClearAll()

	pending := cem.GetPendingEffects()
	if len(pending) != 0 {
		t.Errorf("クリア後の待機中効果数: got %d, want 0", len(pending))
	}
}

// TestNilEffect はnil効果の登録をテストします。
func TestNilEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// nil効果の登録（無視される）
	cem.RegisterChainEffect(0, nil, "slash_lv1")

	pending := cem.GetPendingEffects()
	if len(pending) != 0 {
		t.Errorf("nil効果後の待機中効果数: got %d, want 0", len(pending))
	}
}

// TestHasPendingEffect は待機中効果の存在確認をテストします。
func TestHasPendingEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// 初期状態
	if cem.HasPendingEffect(0) {
		t.Error("初期状態でエージェント0に待機中効果があります")
	}

	// 効果を登録
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 20.0)
	cem.RegisterChainEffect(0, &effect, "slash_lv1")

	// 登録後
	if !cem.HasPendingEffect(0) {
		t.Error("登録後にエージェント0の待機中効果がありません")
	}
	if cem.HasPendingEffect(1) {
		t.Error("登録していないエージェント1に待機中効果があります")
	}
}

// TestDamageAmpEffect はダメージアンプ効果をテストします。
func TestDamageAmpEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// ダメージアンプ効果を登録
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 15.0)
	cem.RegisterChainEffect(0, &effect, "amp_skill")

	// 他エージェントが攻撃（発動）
	triggered := cem.CheckAndTrigger(1, SkillEffectFlags{HasDamage: true})

	if len(triggered) != 1 {
		t.Fatalf("発動した効果数: got %d, want 1", len(triggered))
	}
	if triggered[0].Effect.Type != domain.ChainEffectDamageAmp {
		t.Errorf("Effect.Type: got %v, want %v", triggered[0].Effect.Type, domain.ChainEffectDamageAmp)
	}
}

// TestHealBonusEffect は回復ボーナス効果をテストします。
func TestHealBonusEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// 回復ボーナス効果を登録
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectHealBonus, 30.0)
	cem.RegisterChainEffect(0, &effect, "heal_amp")

	// 他エージェントが回復（発動）
	triggered := cem.CheckAndTrigger(1, SkillEffectFlags{HasHeal: true})

	if len(triggered) != 1 {
		t.Fatalf("発動した効果数: got %d, want 1", len(triggered))
	}
	if triggered[0].Effect.Type != domain.ChainEffectHealBonus {
		t.Errorf("Effect.Type: got %v, want %v", triggered[0].Effect.Type, domain.ChainEffectHealBonus)
	}
}

// TestBuffExtendEffect はバフ延長効果をテストします。
func TestBuffExtendEffect(t *testing.T) {
	cem := NewChainEffectManager()

	// バフ延長効果を登録
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectBuffExtend, 2.0)
	cem.RegisterChainEffect(0, &effect, "buff_extend")

	// 他エージェントがバフ（発動）
	triggered := cem.CheckAndTrigger(1, SkillEffectFlags{HasBuff: true})

	if len(triggered) != 1 {
		t.Fatalf("発動した効果数: got %d, want 1", len(triggered))
	}
	if triggered[0].Effect.Type != domain.ChainEffectBuffExtend {
		t.Errorf("Effect.Type: got %v, want %v", triggered[0].Effect.Type, domain.ChainEffectBuffExtend)
	}
}

// TestEffectCategoryMatching は効果カテゴリとモジュール効果フラグのマッチングをテストします。
func TestEffectCategoryMatching(t *testing.T) {
	tests := []struct {
		name          string
		effectType    domain.ChainEffectType
		moduleFlags   SkillEffectFlags
		shouldTrigger bool
	}{
		// 攻撃強化効果は攻撃モジュールで発動
		{"DamageBonus-Damage", domain.ChainEffectDamageBonus, SkillEffectFlags{HasDamage: true}, true},
		{"DamageBonus-Heal", domain.ChainEffectDamageBonus, SkillEffectFlags{HasHeal: true}, false},

		// 回復強化効果は回復モジュールで発動
		{"HealBonus-Heal", domain.ChainEffectHealBonus, SkillEffectFlags{HasHeal: true}, true},
		{"HealBonus-Damage", domain.ChainEffectHealBonus, SkillEffectFlags{HasDamage: true}, false},

		// バフ延長効果はバフモジュールで発動
		{"BuffExtend-Buff", domain.ChainEffectBuffExtend, SkillEffectFlags{HasBuff: true}, true},
		{"BuffExtend-Debuff", domain.ChainEffectBuffExtend, SkillEffectFlags{HasDebuff: true}, false},

		// デバフ延長効果はデバフモジュールで発動
		{"DebuffExtend-Debuff", domain.ChainEffectDebuffExtend, SkillEffectFlags{HasDebuff: true}, true},
		{"DebuffExtend-Buff", domain.ChainEffectDebuffExtend, SkillEffectFlags{HasBuff: true}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cem := NewChainEffectManager()
			effect := domain.NewChainEffect("test_effect", tt.effectType, 10.0)
			cem.RegisterChainEffect(0, &effect, "test_module")

			triggered := cem.CheckAndTrigger(1, tt.moduleFlags)

			if tt.shouldTrigger && len(triggered) == 0 {
				t.Error("発動すべきなのに発動しませんでした")
			}
			if !tt.shouldTrigger && len(triggered) != 0 {
				t.Error("発動すべきでないのに発動しました")
			}
		})
	}
}
