// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestChainEffectType_定数の確認 はChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"ダメージボーナス", ChainEffectDamageBonus, "damage_bonus"},
		{"回復ボーナス", ChainEffectHealBonus, "heal_bonus"},
		{"バフ延長", ChainEffectBuffExtend, "buff_extend"},
		{"デバフ延長", ChainEffectDebuffExtend, "debuff_extend"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffect_フィールドの確認 はChainEffect構造体のフィールドが正しく設定されることを確認します。
func TestChainEffect_フィールドの確認(t *testing.T) {
	effect := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}

	if effect.Type != ChainEffectDamageBonus {
		t.Errorf("Typeが期待値と異なります: got %s, want %s", effect.Type, ChainEffectDamageBonus)
	}
	if effect.Value != 25.0 {
		t.Errorf("Valueが期待値と異なります: got %f, want 25.0", effect.Value)
	}
	if effect.Description != "攻撃ダメージ+25%" {
		t.Errorf("Descriptionが期待値と異なります: got %s, want 攻撃ダメージ+25%%", effect.Description)
	}
}

// TestChainEffect_Equals_等価なチェイン効果 は同一のチェイン効果が等価と判定されることを確認します。
func TestChainEffect_Equals_等価なチェイン効果(t *testing.T) {
	effect1 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}
	effect2 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}

	if !effect1.Equals(effect2) {
		t.Error("同一のチェイン効果は等価であるべきです")
	}
}

// TestChainEffect_Equals_異なるType は異なるTypeのチェイン効果が非等価と判定されることを確認します。
func TestChainEffect_Equals_異なるType(t *testing.T) {
	effect1 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}
	effect2 := ChainEffect{
		Type:        ChainEffectHealBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}

	if effect1.Equals(effect2) {
		t.Error("異なるTypeのチェイン効果は非等価であるべきです")
	}
}

// TestChainEffect_Equals_異なるValue は異なるValueのチェイン効果が非等価と判定されることを確認します。
func TestChainEffect_Equals_異なるValue(t *testing.T) {
	effect1 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}
	effect2 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       30.0,
		Description: "攻撃ダメージ+25%",
	}

	if effect1.Equals(effect2) {
		t.Error("異なるValueのチェイン効果は非等価であるべきです")
	}
}

// TestChainEffect_Equals_異なるDescription は異なるDescriptionのチェイン効果が非等価と判定されることを確認します。
func TestChainEffect_Equals_異なるDescription(t *testing.T) {
	effect1 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+25%",
	}
	effect2 := ChainEffect{
		Type:        ChainEffectDamageBonus,
		Value:       25.0,
		Description: "攻撃ダメージ+30%",
	}

	if effect1.Equals(effect2) {
		t.Error("異なるDescriptionのチェイン効果は非等価であるべきです")
	}
}

// TestNewChainEffectWithTemplate はNewChainEffectWithTemplate関数でチェイン効果が正しく作成されることを確認します。
func TestNewChainEffectWithTemplate(t *testing.T) {
	effect := NewChainEffectWithTemplate(
		"test_damage_bonus",
		ChainEffectDamageBonus,
		25.0,
		"次の攻撃のダメージ+%.0f%%",
		"次攻撃ダメ+%.0f%%",
	)

	if effect.Type != ChainEffectDamageBonus {
		t.Errorf("Typeが期待値と異なります: got %s, want %s", effect.Type, ChainEffectDamageBonus)
	}
	if effect.Value != 25.0 {
		t.Errorf("Valueが期待値と異なります: got %f, want 25.0", effect.Value)
	}
	// Descriptionがテンプレートから生成されていることを確認
	expectedDesc := "次の攻撃のダメージ+25%"
	if effect.Description != expectedDesc {
		t.Errorf("Descriptionが期待値と異なります: got %s, want %s", effect.Description, expectedDesc)
	}
	// ShortDescriptionも確認
	expectedShortDesc := "次攻撃ダメ+25%"
	if effect.ShortDescription != expectedShortDesc {
		t.Errorf("ShortDescriptionが期待値と異なります: got %s, want %s", effect.ShortDescription, expectedShortDesc)
	}
}

// TestChainEffect_イミュータブル性 はChainEffectがイミュータブルであることを確認します。
func TestChainEffect_イミュータブル性(t *testing.T) {
	effect := NewChainEffectWithTemplate(
		"test_heal_bonus",
		ChainEffectHealBonus,
		15.0,
		"次の回復量+%.0f%%",
		"次回復量+%.0f%%",
	)
	originalValue := effect.Value

	// 値のコピーを作成
	effectCopy := effect
	effectCopy.Value = 999.0

	// 元のeffectは変更されていないはず
	if effect.Value != originalValue {
		t.Errorf("元のChainEffectが変更されています: got %f, want %f", effect.Value, originalValue)
	}
}

// TestChainEffect_Equalsのゼロ値 はゼロ値同士が等価と判定されることを確認します。
func TestChainEffect_Equalsのゼロ値(t *testing.T) {
	effect1 := ChainEffect{}
	effect2 := ChainEffect{}

	if !effect1.Equals(effect2) {
		t.Error("ゼロ値同士は等価であるべきです")
	}
}

// TestChainEffectType_攻撃強化カテゴリ定数の確認 は攻撃強化カテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_攻撃強化カテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"ダメージアンプ", ChainEffectDamageAmp, "damage_amp"},
		{"アーマーピアス", ChainEffectArmorPierce, "armor_pierce"},
		{"ライフスティール", ChainEffectLifeSteal, "life_steal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestNewChainEffectWithTemplate_攻撃強化カテゴリ は攻撃強化カテゴリの効果が正しく作成されることを確認します。
func TestNewChainEffectWithTemplate_攻撃強化カテゴリ(t *testing.T) {
	tests := []struct {
		effectType        ChainEffectType
		value             float64
		descTemplate      string
		shortDescTemplate string
		expectedDesc      string
	}{
		{ChainEffectDamageAmp, 25.0, "効果中の攻撃ダメージ+%.0f%%", "攻撃ダメ+%.0f%%", "効果中の攻撃ダメージ+25%"},
		{ChainEffectArmorPierce, 1.0, "効果中の攻撃が防御バフ無視", "防御バフ無視", "効果中の攻撃が防御バフ無視"},
		{ChainEffectLifeSteal, 10.0, "効果中の攻撃ダメージの%.0f%%回復", "与ダメの%.0f%%回復", "効果中の攻撃ダメージの10%回復"},
	}

	for _, tt := range tests {
		t.Run(string(tt.effectType), func(t *testing.T) {
			effect := NewChainEffectWithTemplate("test_effect", tt.effectType, tt.value, tt.descTemplate, tt.shortDescTemplate)
			if effect.Type != tt.effectType {
				t.Errorf("Typeが期待値と異なります: got %s, want %s", effect.Type, tt.effectType)
			}
			if effect.Value != tt.value {
				t.Errorf("Valueが期待値と異なります: got %f, want %f", effect.Value, tt.value)
			}
			if effect.Description != tt.expectedDesc {
				t.Errorf("Descriptionが期待値と異なります: got %s, want %s", effect.Description, tt.expectedDesc)
			}
		})
	}
}

// TestChainEffectType_防御強化カテゴリ定数の確認 は防御強化カテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_防御強化カテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"ダメージカット", ChainEffectDamageCut, "damage_cut"},
		{"イベイジョン", ChainEffectEvasion, "evasion"},
		{"リフレクト", ChainEffectReflect, "reflect"},
		{"リジェネ", ChainEffectRegen, "regen"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_回復強化カテゴリ定数の確認 は回復強化カテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_回復強化カテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"ヒールアンプ", ChainEffectHealAmp, "heal_amp"},
		{"オーバーヒール", ChainEffectOverheal, "overheal"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_タイピングカテゴリ定数の確認 はタイピングカテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_タイピングカテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"タイムエクステンド", ChainEffectTimeExtend, "time_extend"},
		{"オートコレクト", ChainEffectAutoCorrect, "auto_correct"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_リキャストカテゴリ定数の確認 はリキャストカテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_リキャストカテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"クールダウンリデュース", ChainEffectCooldownReduce, "cooldown_reduce"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_効果延長カテゴリ定数の確認 は効果延長カテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_効果延長カテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"バフデュレーション", ChainEffectBuffDuration, "buff_duration"},
		{"デバフデュレーション", ChainEffectDebuffDuration, "debuff_duration"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_特殊カテゴリ定数の確認 は特殊カテゴリのChainEffectType定数が正しく定義されていることを確認します。
func TestChainEffectType_特殊カテゴリ定数の確認(t *testing.T) {
	tests := []struct {
		name       string
		effectType ChainEffectType
		expected   string
	}{
		{"ダブルキャスト", ChainEffectDoubleCast, "double_cast"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.effectType) != tt.expected {
				t.Errorf("ChainEffectTypeが期待値と異なります: got %s, want %s", tt.effectType, tt.expected)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_攻撃強化 は攻撃強化カテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_攻撃強化(t *testing.T) {
	attackTypes := []ChainEffectType{
		ChainEffectDamageBonus,
		ChainEffectDamageAmp,
		ChainEffectArmorPierce,
		ChainEffectLifeSteal,
	}

	for _, effectType := range attackTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryAttack {
				t.Errorf("%s はattackカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_防御強化 は防御強化カテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_防御強化(t *testing.T) {
	defenseTypes := []ChainEffectType{
		ChainEffectDamageCut,
		ChainEffectEvasion,
		ChainEffectReflect,
		ChainEffectRegen,
	}

	for _, effectType := range defenseTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryDefense {
				t.Errorf("%s はdefenseカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_回復強化 は回復強化カテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_回復強化(t *testing.T) {
	healTypes := []ChainEffectType{
		ChainEffectHealBonus,
		ChainEffectHealAmp,
		ChainEffectOverheal,
	}

	for _, effectType := range healTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryHeal {
				t.Errorf("%s はhealカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_タイピング はタイピングカテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_タイピング(t *testing.T) {
	typingTypes := []ChainEffectType{
		ChainEffectTimeExtend,
		ChainEffectAutoCorrect,
	}

	for _, effectType := range typingTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryTyping {
				t.Errorf("%s はtypingカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_リキャスト はリキャストカテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_リキャスト(t *testing.T) {
	recastTypes := []ChainEffectType{
		ChainEffectCooldownReduce,
	}

	for _, effectType := range recastTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryRecast {
				t.Errorf("%s はrecastカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_効果延長 は効果延長カテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_効果延長(t *testing.T) {
	effectExtendTypes := []ChainEffectType{
		ChainEffectBuffExtend,
		ChainEffectDebuffExtend,
		ChainEffectBuffDuration,
		ChainEffectDebuffDuration,
	}

	for _, effectType := range effectExtendTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategoryEffectExtend {
				t.Errorf("%s はeffect_extendカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectType_カテゴリ判定_特殊 は特殊カテゴリの判定が正しいことを確認します。
func TestChainEffectType_カテゴリ判定_特殊(t *testing.T) {
	specialTypes := []ChainEffectType{
		ChainEffectDoubleCast,
	}

	for _, effectType := range specialTypes {
		t.Run(string(effectType), func(t *testing.T) {
			if effectType.Category() != ChainEffectCategorySpecial {
				t.Errorf("%s はspecialカテゴリであるべきです", effectType)
			}
		})
	}
}

// TestChainEffectCategory_定数確認 はChainEffectCategory定数が正しく定義されていることを確認します。
func TestChainEffectCategory_定数確認(t *testing.T) {
	tests := []struct {
		name     string
		category ChainEffectCategory
		expected string
	}{
		{"攻撃強化", ChainEffectCategoryAttack, "attack"},
		{"防御強化", ChainEffectCategoryDefense, "defense"},
		{"回復強化", ChainEffectCategoryHeal, "heal"},
		{"タイピング", ChainEffectCategoryTyping, "typing"},
		{"リキャスト", ChainEffectCategoryRecast, "recast"},
		{"効果延長", ChainEffectCategoryEffectExtend, "effect_extend"},
		{"特殊", ChainEffectCategorySpecial, "special"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.category) != tt.expected {
				t.Errorf("ChainEffectCategoryが期待値と異なります: got %s, want %s", tt.category, tt.expected)
			}
		})
	}
}
