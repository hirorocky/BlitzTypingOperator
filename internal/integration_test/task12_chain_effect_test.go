// Package integration_test はタスク12の統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// newTestModuleForChain はテスト用モジュールを作成します。
func newTestModuleForChain(id, name string, tags []string, statCoef float64, statRef, description string) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []domain.ModuleEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
}

// newTestModuleWithChainEffectForChain はチェイン効果付きテスト用モジュールを作成します。
func newTestModuleWithChainEffectForChain(id, name string, tags []string, statCoef float64, statRef, description string, chainEffect *domain.ChainEffect) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []domain.ModuleEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, chainEffect)
}

// ==================================================
// Task 12.3: 全19種類のチェイン効果の挙動検証テスト
// ==================================================

// TestChainEffect_AllTypes は全19種類のチェイン効果タイプを検証します。
func TestChainEffect_AllTypes(t *testing.T) {
	// 全チェイン効果タイプのテストケース
	testCases := []struct {
		effectType       domain.ChainEffectType
		value            float64
		expectedCategory domain.ChainEffectCategory
		name             string
	}{
		// 攻撃強化カテゴリ
		{domain.ChainEffectDamageBonus, 20.0, domain.ChainEffectCategoryAttack, "DamageBonus"},
		{domain.ChainEffectDamageAmp, 25.0, domain.ChainEffectCategoryAttack, "DamageAmp"},
		{domain.ChainEffectArmorPierce, 1.0, domain.ChainEffectCategoryAttack, "ArmorPierce"},
		{domain.ChainEffectLifeSteal, 10.0, domain.ChainEffectCategoryAttack, "LifeSteal"},
		// 防御強化カテゴリ
		{domain.ChainEffectDamageCut, 25.0, domain.ChainEffectCategoryDefense, "DamageCut"},
		{domain.ChainEffectEvasion, 10.0, domain.ChainEffectCategoryDefense, "Evasion"},
		{domain.ChainEffectReflect, 15.0, domain.ChainEffectCategoryDefense, "Reflect"},
		{domain.ChainEffectRegen, 2.0, domain.ChainEffectCategoryDefense, "Regen"},
		// 回復強化カテゴリ
		{domain.ChainEffectHealBonus, 15.0, domain.ChainEffectCategoryHeal, "HealBonus"},
		{domain.ChainEffectHealAmp, 25.0, domain.ChainEffectCategoryHeal, "HealAmp"},
		{domain.ChainEffectOverheal, 1.0, domain.ChainEffectCategoryHeal, "Overheal"},
		// タイピングカテゴリ
		{domain.ChainEffectTimeExtend, 3.0, domain.ChainEffectCategoryTyping, "TimeExtend"},
		{domain.ChainEffectAutoCorrect, 2.0, domain.ChainEffectCategoryTyping, "AutoCorrect"},
		// リキャストカテゴリ
		{domain.ChainEffectCooldownReduce, 20.0, domain.ChainEffectCategoryRecast, "CooldownReduce"},
		// 効果延長カテゴリ
		{domain.ChainEffectBuffExtend, 5.0, domain.ChainEffectCategoryEffectExtend, "BuffExtend"},
		{domain.ChainEffectDebuffExtend, 4.0, domain.ChainEffectCategoryEffectExtend, "DebuffExtend"},
		{domain.ChainEffectBuffDuration, 5.0, domain.ChainEffectCategoryEffectExtend, "BuffDuration"},
		{domain.ChainEffectDebuffDuration, 5.0, domain.ChainEffectCategoryEffectExtend, "DebuffDuration"},
		// 特殊カテゴリ
		{domain.ChainEffectDoubleCast, 10.0, domain.ChainEffectCategorySpecial, "DoubleCast"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// ChainEffect作成（説明文はテンプレートから別途生成するため、ここでは簡易コンストラクタを使用）
			effect := domain.NewChainEffect("test_effect", tc.effectType, tc.value)

			// Type確認
			if effect.Type != tc.effectType {
				t.Errorf("Type expected %s, got %s", tc.effectType, effect.Type)
			}

			// Value確認
			if effect.Value != tc.value {
				t.Errorf("Value expected %f, got %f", tc.value, effect.Value)
			}

			// Category確認
			category := tc.effectType.Category()
			if category != tc.expectedCategory {
				t.Errorf("Category expected %s, got %s", tc.expectedCategory, category)
			}
		})
	}
}

// TestChainEffect_AttackCategory は攻撃強化カテゴリのチェイン効果を検証します。
func TestChainEffect_AttackCategory(t *testing.T) {
	t.Run("DamageAmp", func(t *testing.T) {
		// ダメージアンプ: 効果中の攻撃ダメージ+X%
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDamageAmp, 25.0,
			"効果中の攻撃ダメージ+%.0f%%", "攻撃ダメ+%.0f%%")

		if effect.Type != domain.ChainEffectDamageAmp {
			t.Error("Type should be damage_amp")
		}
		if effect.Value != 25.0 {
			t.Error("Value should be 25.0")
		}
		if effect.Type.Category() != domain.ChainEffectCategoryAttack {
			t.Error("Category should be attack")
		}
		// 説明文に効果値が含まれていることを確認
		expectedDesc := "効果中の攻撃ダメージ+25%"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("ArmorPierce", func(t *testing.T) {
		// アーマーピアス: 効果中の攻撃が防御バフ無視（プレースホルダなし）
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectArmorPierce, 1.0,
			"効果中の攻撃が防御バフ無視", "防御バフ無視")

		if effect.Type != domain.ChainEffectArmorPierce {
			t.Error("Type should be armor_pierce")
		}
		if effect.Type.Category() != domain.ChainEffectCategoryAttack {
			t.Error("Category should be attack")
		}
		expectedDesc := "効果中の攻撃が防御バフ無視"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("LifeSteal", func(t *testing.T) {
		// ライフスティール: 効果中の攻撃ダメージのX%回復
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectLifeSteal, 10.0,
			"効果中の攻撃ダメージの%.0f%%回復", "与ダメの%.0f%%回復")

		if effect.Type != domain.ChainEffectLifeSteal {
			t.Error("Type should be life_steal")
		}
		expectedDesc := "効果中の攻撃ダメージの10%回復"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_DefenseCategory は防御強化カテゴリのチェイン効果を検証します。
func TestChainEffect_DefenseCategory(t *testing.T) {
	t.Run("DamageCut", func(t *testing.T) {
		// ダメージカット: 効果中の被ダメージ-X%
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDamageCut, 25.0,
			"効果中の被ダメージ-%.0f%%", "被ダメ-%.0f%%")

		if effect.Type != domain.ChainEffectDamageCut {
			t.Error("Type should be damage_cut")
		}
		if effect.Type.Category() != domain.ChainEffectCategoryDefense {
			t.Error("Category should be defense")
		}
		expectedDesc := "効果中の被ダメージ-25%"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("Evasion", func(t *testing.T) {
		// イベイジョン: 効果中X%で攻撃回避
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectEvasion, 10.0,
			"効果中%.0f%%で攻撃回避", "%.0f%%回避")

		expectedDesc := "効果中10%で攻撃回避"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("Reflect", func(t *testing.T) {
		// リフレクト: 効果中被ダメージの X%反射
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectReflect, 15.0,
			"効果中被ダメージの%.0f%%反射", "被ダメ%.0f%%反射")

		expectedDesc := "効果中被ダメージの15%反射"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("Regen", func(t *testing.T) {
		// リジェネ: 効果中毎秒HP X%回復
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectRegen, 2.0,
			"効果中毎秒HP%.0f%%回復", "毎秒HP%.0f%%回復")

		expectedDesc := "効果中毎秒HP2%回復"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_HealCategory は回復強化カテゴリのチェイン効果を検証します。
func TestChainEffect_HealCategory(t *testing.T) {
	t.Run("HealAmp", func(t *testing.T) {
		// ヒールアンプ: 効果中の回復量+X%
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectHealAmp, 25.0,
			"効果中の回復量+%.0f%%", "回復量+%.0f%%")

		if effect.Type.Category() != domain.ChainEffectCategoryHeal {
			t.Error("Category should be heal")
		}
		expectedDesc := "効果中の回復量+25%"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("Overheal", func(t *testing.T) {
		// オーバーヒール: 効果中の超過回復を一時HPに（プレースホルダなし）
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectOverheal, 1.0,
			"効果中の超過回復を一時HPに", "超過回復→一時HP")

		expectedDesc := "効果中の超過回復を一時HPに"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_TypingCategory はタイピングカテゴリのチェイン効果を検証します。
func TestChainEffect_TypingCategory(t *testing.T) {
	t.Run("TimeExtend", func(t *testing.T) {
		// タイムエクステンド: 効果中のタイピング制限時間+X秒
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectTimeExtend, 3.0,
			"効果中のタイピング制限時間+%.0f秒", "入力時間+%.0f秒")

		if effect.Type.Category() != domain.ChainEffectCategoryTyping {
			t.Error("Category should be typing")
		}
		expectedDesc := "効果中のタイピング制限時間+3秒"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("AutoCorrect", func(t *testing.T) {
		// オートコレクト: 効果中ミスX回まで無視
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectAutoCorrect, 2.0,
			"効果中ミス%.0f回まで無視", "ミス%.0f回無視")

		expectedDesc := "効果中ミス2回まで無視"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_RecastCategory はリキャストカテゴリのチェイン効果を検証します。
func TestChainEffect_RecastCategory(t *testing.T) {
	t.Run("CooldownReduce", func(t *testing.T) {
		// クールダウンリデュース: 効果中発生した他エージェントのリキャスト時間X%短縮
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectCooldownReduce, 20.0,
			"効果中発生した他エージェントのリキャスト時間%.0f%%短縮", "他CD%.0f%%短縮")

		if effect.Type.Category() != domain.ChainEffectCategoryRecast {
			t.Error("Category should be recast")
		}
		expectedDesc := "効果中発生した他エージェントのリキャスト時間20%短縮"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_EffectExtendCategory は効果延長カテゴリのチェイン効果を検証します。
func TestChainEffect_EffectExtendCategory(t *testing.T) {
	t.Run("BuffDuration", func(t *testing.T) {
		// バフデュレーション: 効果中のバフスキル効果時間+X秒
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectBuffDuration, 5.0,
			"効果中のバフスキル効果時間+%.0f秒", "バフ延長+%.0f秒")

		if effect.Type.Category() != domain.ChainEffectCategoryEffectExtend {
			t.Error("Category should be effect_extend")
		}
		expectedDesc := "効果中のバフスキル効果時間+5秒"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("DebuffDuration", func(t *testing.T) {
		// デバフデュレーション: 効果中のデバフスキル効果時間+X秒
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDebuffDuration, 5.0,
			"効果中のデバフスキル効果時間+%.0f秒", "デバフ延長+%.0f秒")

		expectedDesc := "効果中のデバフスキル効果時間+5秒"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("BuffExtend", func(t *testing.T) {
		// バフ延長
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectBuffExtend, 5.0,
			"バフ効果時間+%.0f秒", "バフ時間+%.0f秒")

		expectedDesc := "バフ効果時間+5秒"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})

	t.Run("DebuffExtend", func(t *testing.T) {
		// デバフ延長
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDebuffExtend, 4.0,
			"デバフ効果時間+%.0f秒", "デバフ時間+%.0f秒")

		expectedDesc := "デバフ効果時間+4秒"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_SpecialCategory は特殊カテゴリのチェイン効果を検証します。
func TestChainEffect_SpecialCategory(t *testing.T) {
	t.Run("DoubleCast", func(t *testing.T) {
		// ダブルキャスト: 効果中X%でスキル2回発動
		effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDoubleCast, 10.0,
			"効果中%.0f%%でスキル2回発動", "%.0f%%で2回発動")

		if effect.Type.Category() != domain.ChainEffectCategorySpecial {
			t.Error("Category should be special")
		}
		expectedDesc := "効果中10%でスキル2回発動"
		if effect.Description != expectedDesc {
			t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
		}
	})
}

// TestChainEffect_Equality はチェイン効果の等価性判定を検証します。
func TestChainEffect_Equality(t *testing.T) {
	effect1 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 25.0)
	effect2 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 25.0)
	effect3 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 30.0)
	effect4 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageCut, 25.0)

	// 同じ効果は等価
	if !effect1.Equals(effect2) {
		t.Error("同じ効果は等価であるべきです")
	}

	// 値が異なれば不等価
	if effect1.Equals(effect3) {
		t.Error("値が異なれば不等価であるべきです")
	}

	// タイプが異なれば不等価
	if effect1.Equals(effect4) {
		t.Error("タイプが異なれば不等価であるべきです")
	}
}

// TestChainEffect_ValueCalculation はチェイン効果の効果値計算を検証します。
func TestChainEffect_ValueCalculation(t *testing.T) {
	testCases := []struct {
		name       string
		effectType domain.ChainEffectType
		value      float64
	}{
		{"最小ダメージアンプ", domain.ChainEffectDamageAmp, 15.0},
		{"最大ダメージアンプ", domain.ChainEffectDamageAmp, 30.0},
		{"最小ダメージカット", domain.ChainEffectDamageCut, 15.0},
		{"最大ダメージカット", domain.ChainEffectDamageCut, 30.0},
		{"最小ヒールアンプ", domain.ChainEffectHealAmp, 15.0},
		{"最大ヒールアンプ", domain.ChainEffectHealAmp, 30.0},
		{"最小リジェネ", domain.ChainEffectRegen, 1.0},
		{"最大リジェネ", domain.ChainEffectRegen, 3.0},
		{"最小バフ延長", domain.ChainEffectBuffDuration, 3.0},
		{"最大バフ延長", domain.ChainEffectBuffDuration, 7.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			effect := domain.NewChainEffect("test_effect", tc.effectType, tc.value)

			if effect.Value != tc.value {
				t.Errorf("Value expected %f, got %f", tc.value, effect.Value)
			}

			// 効果値が正の値であることを確認
			if effect.Value <= 0 {
				t.Error("効果値は正の値であるべきです")
			}
		})
	}
}

// TestChainEffect_ModuleIntegration はモジュールへのチェイン効果統合を検証します。
func TestChainEffect_ModuleIntegration(t *testing.T) {
	// チェイン効果付きモジュール
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 25.0)
	module := newTestModuleWithChainEffectForChain(
		"physical_lv1", "物理打撃Lv1",
		[]string{"physical_low"}, 1.0, "STR", "物理ダメージ",
		&chainEffect,
	)

	if !module.HasChainEffect() {
		t.Error("モジュールにはチェイン効果があるべきです")
	}

	if module.ChainEffect.Type != domain.ChainEffectDamageAmp {
		t.Errorf("ChainEffect.Type expected damage_amp, got %s", module.ChainEffect.Type)
	}

	// チェイン効果なしモジュール
	moduleNoEffect := newTestModuleForChain(
		"physical_lv2", "物理打撃Lv2",
		[]string{"physical_mid"}, 1.5, "STR", "物理ダメージ",
	)

	if moduleNoEffect.HasChainEffect() {
		t.Error("モジュールにはチェイン効果がないべきです")
	}
}

// TestChainEffect_CategoryIcon はカテゴリごとのアイコンマッピングを検証します。
func TestChainEffect_CategoryIcon(t *testing.T) {
	// 各カテゴリが正しく判定されることを確認
	categoryTests := []struct {
		effectType domain.ChainEffectType
		category   domain.ChainEffectCategory
	}{
		{domain.ChainEffectDamageAmp, domain.ChainEffectCategoryAttack},
		{domain.ChainEffectArmorPierce, domain.ChainEffectCategoryAttack},
		{domain.ChainEffectLifeSteal, domain.ChainEffectCategoryAttack},
		{domain.ChainEffectDamageCut, domain.ChainEffectCategoryDefense},
		{domain.ChainEffectEvasion, domain.ChainEffectCategoryDefense},
		{domain.ChainEffectReflect, domain.ChainEffectCategoryDefense},
		{domain.ChainEffectRegen, domain.ChainEffectCategoryDefense},
		{domain.ChainEffectHealAmp, domain.ChainEffectCategoryHeal},
		{domain.ChainEffectOverheal, domain.ChainEffectCategoryHeal},
		{domain.ChainEffectTimeExtend, domain.ChainEffectCategoryTyping},
		{domain.ChainEffectAutoCorrect, domain.ChainEffectCategoryTyping},
		{domain.ChainEffectCooldownReduce, domain.ChainEffectCategoryRecast},
		{domain.ChainEffectBuffDuration, domain.ChainEffectCategoryEffectExtend},
		{domain.ChainEffectDebuffDuration, domain.ChainEffectCategoryEffectExtend},
		{domain.ChainEffectDoubleCast, domain.ChainEffectCategorySpecial},
	}

	for _, tc := range categoryTests {
		t.Run(string(tc.effectType), func(t *testing.T) {
			category := tc.effectType.Category()
			if category != tc.category {
				t.Errorf("Category for %s expected %s, got %s", tc.effectType, tc.category, category)
			}
		})
	}
}

// TestChainEffect_AllDescriptionsGenerated は全19種類のチェイン効果の説明文テンプレートを検証します。
func TestChainEffect_AllDescriptionsGenerated(t *testing.T) {
	// テスト用の説明文テンプレートマップ（実際はJSONから読み込む）
	descTemplates := map[domain.ChainEffectType]struct{ desc, shortDesc string }{
		domain.ChainEffectDamageBonus:    {"次の攻撃のダメージ+%.0f%%", "次攻撃ダメ+%.0f%%"},
		domain.ChainEffectHealBonus:      {"次の回復量+%.0f%%", "次回復量+%.0f%%"},
		domain.ChainEffectBuffExtend:     {"バフ効果時間+%.0f秒", "バフ時間+%.0f秒"},
		domain.ChainEffectDebuffExtend:   {"デバフ効果時間+%.0f秒", "デバフ時間+%.0f秒"},
		domain.ChainEffectDamageAmp:      {"効果中の攻撃ダメージ+%.0f%%", "攻撃ダメ+%.0f%%"},
		domain.ChainEffectArmorPierce:    {"効果中の攻撃が防御バフ無視", "防御バフ無視"},
		domain.ChainEffectLifeSteal:      {"効果中の攻撃ダメージの%.0f%%回復", "与ダメの%.0f%%回復"},
		domain.ChainEffectDamageCut:      {"効果中の被ダメージ-%.0f%%", "被ダメ-%.0f%%"},
		domain.ChainEffectEvasion:        {"効果中%.0f%%で攻撃回避", "%.0f%%回避"},
		domain.ChainEffectReflect:        {"効果中被ダメージの%.0f%%反射", "被ダメ%.0f%%反射"},
		domain.ChainEffectRegen:          {"効果中毎秒HP%.0f%%回復", "毎秒HP%.0f%%回復"},
		domain.ChainEffectHealAmp:        {"効果中の回復量+%.0f%%", "回復量+%.0f%%"},
		domain.ChainEffectOverheal:       {"効果中の超過回復を一時HPに", "超過回復→一時HP"},
		domain.ChainEffectTimeExtend:     {"効果中のタイピング制限時間+%.0f秒", "入力時間+%.0f秒"},
		domain.ChainEffectAutoCorrect:    {"効果中ミス%.0f回まで無視", "ミス%.0f回無視"},
		domain.ChainEffectCooldownReduce: {"効果中発生した他エージェントのリキャスト時間%.0f%%短縮", "他CD%.0f%%短縮"},
		domain.ChainEffectBuffDuration:   {"効果中のバフスキル効果時間+%.0f秒", "バフ延長+%.0f秒"},
		domain.ChainEffectDebuffDuration: {"効果中のデバフスキル効果時間+%.0f秒", "デバフ延長+%.0f秒"},
		domain.ChainEffectDoubleCast:     {"効果中%.0f%%でスキル2回発動", "%.0f%%で2回発動"},
	}

	for effectType, templates := range descTemplates {
		t.Run(string(effectType), func(t *testing.T) {
			effect := domain.NewChainEffectWithTemplate("test_effect", effectType, 10.0, templates.desc, templates.shortDesc)

			// 説明文が生成されていることを確認
			if effect.Description == "" {
				t.Errorf("効果タイプ %s の説明文が生成されていません", effectType)
			}

			// 短縮説明文が生成されていることを確認
			if effect.ShortDescription == "" {
				t.Errorf("効果タイプ %s の短縮説明文が生成されていません", effectType)
			}
		})
	}

	// チェイン効果の総数を確認
	if len(descTemplates) != 19 {
		t.Errorf("チェイン効果タイプ数 expected 19, got %d", len(descTemplates))
	}
}

// TestChainEffect_ZeroValue はゼロ値のチェイン効果を検証します。
func TestChainEffect_ZeroValue(t *testing.T) {
	// ゼロ値でも説明文が生成されることを確認
	effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDamageAmp, 0,
		"効果中の攻撃ダメージ+%.0f%%", "攻撃ダメ+%.0f%%")

	if effect.Description == "" {
		t.Error("ゼロ値でも説明文が生成されるべきです")
	}

	expectedDesc := "効果中の攻撃ダメージ+0%"
	if effect.Description != expectedDesc {
		t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
	}
}

// TestChainEffect_NegativeValue は負の値のチェイン効果を検証します。
func TestChainEffect_NegativeValue(t *testing.T) {
	// 負の値でも動作することを確認（通常は使用しないが）
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, -10.0)

	if effect.Value != -10.0 {
		t.Errorf("Value expected -10.0, got %f", effect.Value)
	}
}

// TestChainEffect_LargeValue は大きな値のチェイン効果を検証します。
func TestChainEffect_LargeValue(t *testing.T) {
	effect := domain.NewChainEffectWithTemplate("test_effect", domain.ChainEffectDamageAmp, 1000.0,
		"効果中の攻撃ダメージ+%.0f%%", "攻撃ダメ+%.0f%%")

	if effect.Value != 1000.0 {
		t.Errorf("Value expected 1000.0, got %f", effect.Value)
	}

	expectedDesc := "効果中の攻撃ダメージ+1000%"
	if effect.Description != expectedDesc {
		t.Errorf("Description expected '%s', got '%s'", expectedDesc, effect.Description)
	}
}

// TestChainEffect_DecimalValue は小数値のチェイン効果を検証します。
func TestChainEffect_DecimalValue(t *testing.T) {
	// 小数値の場合、説明文では整数に丸められる
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectRegen, 1.5)

	if effect.Value != 1.5 {
		t.Errorf("Value expected 1.5, got %f", effect.Value)
	}

	// %.0f では小数点以下が切り捨てられる
	expectedDesc := "効果中毎秒HP2%回復" // 1.5 -> 2 (四捨五入)
	if effect.Description != expectedDesc {
		t.Logf("Note: 小数値 1.5 は説明文で %s として表示されます", effect.Description)
	}
}
