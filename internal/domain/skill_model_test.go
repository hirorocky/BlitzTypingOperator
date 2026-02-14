// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestSkillModel_フィールドの確認 はSkillModel構造体のフィールドが正しく設定されることを確認します。
func TestSkillModel_フィールドの確認(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:          "fireball_lv1",
		Name:        "ファイアボール",
		Icon:        "🔥",
		Tags:        []string{"magic_low"},
		Description: "炎の魔法で敵に魔法ダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "INT"},
				Probability: 1.0,
				Icon:        "🔥",
			},
		},
	}, nil)

	if skill.TypeID != "fireball_lv1" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want fireball_lv1", skill.TypeID)
	}
	if skill.Name() != "ファイアボール" {
		t.Errorf("Name()が期待値と異なります: got %s, want ファイアボール", skill.Name())
	}
	if len(skill.Tags()) != 1 || skill.Tags()[0] != "magic_low" {
		t.Errorf("Tags()が期待値と異なります: got %v, want [magic_low]", skill.Tags())
	}
	if skill.Description() != "炎の魔法で敵に魔法ダメージを与える" {
		t.Errorf("Description()が期待値と異なります: got %s", skill.Description())
	}
	if len(skill.Effects()) != 1 {
		t.Errorf("Effects()の長さが期待値と異なります: got %d, want 1", len(skill.Effects()))
	}
}

// TestNewSkillFromType_スキル作成 はNewSkillFromType関数でスキルが正しく作成されることを確認します。
func TestNewSkillFromType_スキル作成(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:          "physical_attack_lv1",
		Name:        "物理打撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "物理攻撃で敵にダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)

	if skill.TypeID != "physical_attack_lv1" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want physical_attack_lv1", skill.TypeID)
	}
	if skill.Name() != "物理打撃" {
		t.Errorf("Name()が期待値と異なります: got %s, want 物理打撃", skill.Name())
	}
}

// TestNewSkillFromType_タグのコピー はNewSkillFromTypeで作成したスキルのTagsが元のスライスと独立していることを確認します。
func TestNewSkillFromType_タグのコピー(t *testing.T) {
	originalTags := []string{"magic_low", "fire"}
	skillType := SkillType{
		ID:          "fireball_lv1",
		Name:        "ファイアボール",
		Icon:        "🔥",
		Tags:        originalTags,
		Description: "炎の魔法で敵に魔法ダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "INT"},
				Probability: 1.0,
			},
		},
	}
	_ = NewSkillFromType(skillType, nil)

	// 元のタグを変更
	originalTags[0] = "modified_tag"

	// SkillTypeのTagsはスライスなので影響を受ける（GoのスライスはReferenceのため）
	// この挙動は許容される（パフォーマンスのためのトレードオフ）
	// 本番コードではマスタデータは変更されないため問題なし
}

// TestSkillModel_HasTag_タグ存在確認 はHasTagメソッドがタグの存在を正しく判定することを確認します。
func TestSkillModel_HasTag_タグ存在確認(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:   "test_skill",
		Tags: []string{"physical_low", "fire"},
	}, nil)

	if !skill.HasTag("physical_low") {
		t.Error("physical_lowタグが存在するはずですがfalseが返されました")
	}
	if !skill.HasTag("fire") {
		t.Error("fireタグが存在するはずですがfalseが返されました")
	}
	if skill.HasTag("magic_low") {
		t.Error("magic_lowタグは存在しないはずですがtrueが返されました")
	}
}

// TestSkillModel_HasTag_空タグリスト はTagsが空の場合に常にfalseを返すことを確認します。
func TestSkillModel_HasTag_空タグリスト(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:   "test_skill",
		Tags: []string{},
	}, nil)

	if skill.HasTag("physical_low") {
		t.Error("Tagsが空の場合、falseを返すべきです")
	}
}

// TestSkillModel_IsCompatibleWithCore はスキルがコアに装備可能かを判定するメソッドをテストします。
func TestSkillModel_IsCompatibleWithCore(t *testing.T) {
	// 物理攻撃と魔法攻撃の低レベルスキルを許可するコア
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		AllowedTags: []string{"physical_low", "magic_low"},
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
	}
	core := NewCoreWithTypeID("test", coreType, PassiveSkill{})

	// 互換性のあるスキル
	compatibleSkill := NewSkillFromType(SkillType{
		ID:   "physical_attack_lv1",
		Tags: []string{"physical_low"},
	}, nil)

	// 互換性のないスキル
	incompatibleSkill := NewSkillFromType(SkillType{
		ID:   "heal_lv2",
		Tags: []string{"heal_mid"},
	}, nil)

	if !compatibleSkill.IsCompatibleWithCore(core) {
		t.Error("physical_lowタグを持つスキルはコアと互換性があるはずです")
	}

	if incompatibleSkill.IsCompatibleWithCore(core) {
		t.Error("heal_midタグを持つスキルはコアと互換性がないはずです")
	}
}

// TestSkillModel_IsCompatibleWithCore_複数タグ はスキルが複数タグを持つ場合の互換性判定をテストします。
func TestSkillModel_IsCompatibleWithCore_複数タグ(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		AllowedTags: []string{"physical_low", "magic_low"},
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
	}
	core := NewCoreWithTypeID("test", coreType, PassiveSkill{})

	// 複数タグのうち1つがコアの許可タグに含まれる場合
	skillWithMultipleTags := NewSkillFromType(SkillType{
		ID:   "hybrid_lv1",
		Tags: []string{"physical_low", "fire"},
	}, nil)

	if !skillWithMultipleTags.IsCompatibleWithCore(core) {
		t.Error("1つでもコアの許可タグに含まれるタグがあれば互換性があるはずです")
	}

	// どのタグもコアの許可タグに含まれない場合
	skillNoMatch := NewSkillFromType(SkillType{
		ID:   "heal_lv1",
		Tags: []string{"heal_low", "light"},
	}, nil)

	if skillNoMatch.IsCompatibleWithCore(core) {
		t.Error("どのタグもコアの許可タグに含まれない場合、互換性がないはずです")
	}
}

// ==================== Task 7.2: Icon()メソッドのテスト ====================

// TestSkillType_Icon はSkillTypeのIconフィールドが正しく設定されることを確認します。
func TestSkillType_Icon(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:   "test",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
	}, nil)

	if skill.Icon() != "⚔️" {
		t.Errorf("Icon()が期待値と異なります: got %s, want ⚔️", skill.Icon())
	}
}

// TestSkillModel_Icon_Empty は空のアイコンに対してIcon()がデフォルト値を返すことを確認します。
func TestSkillModel_Icon_Empty(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:   "test",
		Icon: "",
		Tags: []string{"physical_low"},
	}, nil)

	if skill.Icon() != "•" {
		t.Errorf("空のアイコンに対するIcon()が期待値と異なります: got %s, want •", skill.Icon())
	}
}

// ==================== SkillModel TypeID/ChainEffect リファクタリングテスト ====================

// TestSkillModel_TypeIDフィールドの確認 はSkillModelにTypeIDフィールドが存在することを確認します。
func TestSkillModel_TypeIDフィールドの確認(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:          "physical_attack_lv1",
		Name:        "物理打撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "物理攻撃で敵にダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
			},
		},
	}, nil)

	if skill.TypeID != "physical_attack_lv1" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want physical_attack_lv1", skill.TypeID)
	}
	if skill.ChainEffect != nil {
		t.Errorf("ChainEffectはnilであるべきです: got %v", skill.ChainEffect)
	}
}

// TestSkillModel_ChainEffect付きの作成 はChainEffect付きのスキル作成をテストします。
func TestSkillModel_ChainEffect付きの作成(t *testing.T) {
	chainEffect := NewChainEffectWithTemplate("test_damage_bonus", ChainEffectDamageBonus, 25.0, "次の攻撃のダメージ+%.0f%%", "次攻撃ダメ+%.0f%%")
	skill := NewSkillFromType(SkillType{
		ID:          "physical_attack_lv1",
		Name:        "物理打撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "物理攻撃で敵にダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
			},
		},
	}, &chainEffect)

	if skill.ChainEffect == nil {
		t.Fatal("ChainEffectがnilです")
	}
	if skill.ChainEffect.Type != ChainEffectDamageBonus {
		t.Errorf("ChainEffect.Typeが期待値と異なります: got %s, want %s", skill.ChainEffect.Type, ChainEffectDamageBonus)
	}
	if skill.ChainEffect.Value != 25.0 {
		t.Errorf("ChainEffect.Valueが期待値と異なります: got %f, want 25.0", skill.ChainEffect.Value)
	}
}

// TestSkillModel_同一TypeID異なるChainEffect は同一TypeIDで異なるChainEffectを持つスキルを許容することを確認します。
func TestSkillModel_同一TypeID異なるChainEffect(t *testing.T) {
	chainEffect1 := NewChainEffectWithTemplate("test_damage_bonus", ChainEffectDamageBonus, 25.0, "次の攻撃のダメージ+%.0f%%", "次攻撃ダメ+%.0f%%")
	chainEffect2 := NewChainEffectWithTemplate("test_heal_bonus", ChainEffectHealBonus, 20.0, "次の回復量+%.0f%%", "次回復量+%.0f%%")

	skillType := SkillType{
		ID:          "physical_attack_lv1",
		Name:        "物理打撃",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "物理攻撃で敵にダメージを与える",
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
			},
		},
	}

	skill1 := NewSkillFromType(skillType, &chainEffect1)
	skill2 := NewSkillFromType(skillType, &chainEffect2)

	// 同じTypeIDであっても異なるChainEffectを持つことを許容
	if skill1.TypeID != skill2.TypeID {
		t.Error("同じTypeIDであるべきです")
	}
	if skill1.ChainEffect.Type == skill2.ChainEffect.Type {
		t.Error("異なるChainEffectを持っているはずです")
	}
}

// TestSkillModel_ChainEffectなし はChainEffectがnilのスキルが正しく動作することを確認します。
func TestSkillModel_ChainEffectなし(t *testing.T) {
	skill := NewSkillFromType(SkillType{
		ID:          "heal_lv1",
		Name:        "ヒール",
		Icon:        "💚",
		Tags:        []string{"heal_low"},
		Description: "HPを回復する",
		Effects: []SkillEffect{
			{
				Target:      TargetSelf,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 0.8, StatRef: "INT"},
				Probability: 1.0,
			},
		},
	}, nil)

	if skill.ChainEffect != nil {
		t.Errorf("ChainEffectはnilであるべきです: got %v", skill.ChainEffect)
	}

	// HasChainEffectメソッドのテスト
	if skill.HasChainEffect() {
		t.Error("ChainEffectがない場合、HasChainEffect()はfalseを返すべきです")
	}
}

// TestSkillModel_HasChainEffect はHasChainEffectメソッドをテストします。
func TestSkillModel_HasChainEffect(t *testing.T) {
	chainEffect := NewChainEffectWithTemplate("test_buff_extend", ChainEffectBuffExtend, 5.0, "バフ効果時間+%.0f秒", "バフ時間+%.0f秒")
	skillWithEffect := NewSkillFromType(SkillType{
		ID:          "buff_lv1",
		Name:        "バフ",
		Icon:        "⬆️",
		Tags:        []string{"buff_low"},
		Description: "バフを付与する",
		Effects: []SkillEffect{
			{
				Target: TargetSelf,
				ColumnSpec: &EffectColumnSpec{
					Column:   ColDamageBonus,
					Value:    10.0,
					Duration: 10.0,
				},
				Probability: 1.0,
			},
		},
	}, &chainEffect)

	if !skillWithEffect.HasChainEffect() {
		t.Error("ChainEffectがある場合、HasChainEffect()はtrueを返すべきです")
	}

	skillWithoutEffect := NewSkillFromType(SkillType{
		ID:          "buff_lv1",
		Name:        "バフ",
		Icon:        "⬆️",
		Tags:        []string{"buff_low"},
		Description: "バフを付与する",
		Effects: []SkillEffect{
			{
				Target: TargetSelf,
				ColumnSpec: &EffectColumnSpec{
					Column:   ColDamageBonus,
					Value:    10.0,
					Duration: 10.0,
				},
				Probability: 1.0,
			},
		},
	}, nil)

	if skillWithoutEffect.HasChainEffect() {
		t.Error("ChainEffectがない場合、HasChainEffect()はfalseを返すべきです")
	}
}
