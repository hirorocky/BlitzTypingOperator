// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// ==================== Task 1.1: Module → Skill ユビキタス言語変更テスト ====================

// TestSkillType_フィールドの確認 はSkillType構造体のフィールドが正しく設定されることを確認します。
func TestSkillType_フィールドの確認(t *testing.T) {
	skillType := SkillType{
		ID:              "fireball_lv1",
		Name:            "ファイアボール",
		Icon:            "🔥",
		Tags:            []string{"magic_low"},
		Description:     "炎の魔法で敵に魔法ダメージを与える",
		CooldownSeconds: 10.0,
		Difficulty:      1,
		MinDropLevel:    1,
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "INT"},
				Probability: 1.0,
				Icon:        "🔥",
			},
		},
	}

	if skillType.ID != "fireball_lv1" {
		t.Errorf("IDが期待値と異なります: got %s, want fireball_lv1", skillType.ID)
	}
	if skillType.Name != "ファイアボール" {
		t.Errorf("Nameが期待値と異なります: got %s, want ファイアボール", skillType.Name)
	}
	if len(skillType.Tags) != 1 || skillType.Tags[0] != "magic_low" {
		t.Errorf("Tagsが期待値と異なります: got %v, want [magic_low]", skillType.Tags)
	}
	if skillType.Description != "炎の魔法で敵に魔法ダメージを与える" {
		t.Errorf("Descriptionが期待値と異なります: got %s", skillType.Description)
	}
	if len(skillType.Effects) != 1 {
		t.Errorf("Effectsの長さが期待値と異なります: got %d, want 1", len(skillType.Effects))
	}
}

// TestSkillType_HasTag はHasTagメソッドがタグの存在を正しく判定することを確認します。
func TestSkillType_HasTag(t *testing.T) {
	skillType := SkillType{
		ID:   "test_skill",
		Tags: []string{"physical_low", "fire"},
	}

	if !skillType.HasTag("physical_low") {
		t.Error("physical_lowタグが存在するはずですがfalseが返されました")
	}
	if !skillType.HasTag("fire") {
		t.Error("fireタグが存在するはずですがfalseが返されました")
	}
	if skillType.HasTag("magic_low") {
		t.Error("magic_lowタグは存在しないはずですがtrueが返されました")
	}
}

// TestSkillType_HasTag_空タグリスト はTagsが空の場合に常にfalseを返すことを確認します。
func TestSkillType_HasTag_空タグリスト(t *testing.T) {
	skillType := SkillType{
		ID:   "test_skill",
		Tags: []string{},
	}

	if skillType.HasTag("physical_low") {
		t.Error("Tagsが空の場合、falseを返すべきです")
	}
}

// TestSkillType_IsCompatibleWithCore はスキルがコアに装備可能かを判定するメソッドをテストします。
func TestSkillType_IsCompatibleWithCore(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		AllowedTags: []string{"physical_low", "magic_low"},
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
	}
	core := NewCoreWithTypeID("test", coreType, PassiveSkill{})

	// 互換性のあるスキル
	compatibleSkill := SkillType{
		ID:   "physical_attack_lv1",
		Tags: []string{"physical_low"},
	}

	// 互換性のないスキル
	incompatibleSkill := SkillType{
		ID:   "heal_lv2",
		Tags: []string{"heal_mid"},
	}

	if !compatibleSkill.IsCompatibleWithCore(core) {
		t.Error("physical_lowタグを持つスキルはコアと互換性があるはずです")
	}

	if incompatibleSkill.IsCompatibleWithCore(core) {
		t.Error("heal_midタグを持つスキルはコアと互換性がないはずです")
	}
}

// TestSkillEffect_IsDamageEffect はダメージ効果の判定をテストします。
func TestSkillEffect_IsDamageEffect(t *testing.T) {
	damageEffect := SkillEffect{
		Target:    TargetEnemy,
		HPFormula: &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
	}
	if !damageEffect.IsDamageEffect() {
		t.Error("敵対象のHPFormula効果はダメージ効果であるべきです")
	}

	healEffect := SkillEffect{
		Target:    TargetSelf,
		HPFormula: &HPFormula{Base: 0, StatCoef: 0.8, StatRef: "INT"},
	}
	if healEffect.IsDamageEffect() {
		t.Error("自身対象のHPFormula効果はダメージ効果ではないべきです")
	}
}

// TestSkillEffect_IsHealEffect は回復効果の判定をテストします。
func TestSkillEffect_IsHealEffect(t *testing.T) {
	healEffect := SkillEffect{
		Target:    TargetSelf,
		HPFormula: &HPFormula{Base: 0, StatCoef: 0.8, StatRef: "INT"},
	}
	if !healEffect.IsHealEffect() {
		t.Error("自身対象のHPFormula効果は回復効果であるべきです")
	}

	damageEffect := SkillEffect{
		Target:    TargetEnemy,
		HPFormula: &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
	}
	if damageEffect.IsHealEffect() {
		t.Error("敵対象のHPFormula効果は回復効果ではないべきです")
	}
}

// TestSkillEffect_IsBuffEffect はバフ効果の判定をテストします。
func TestSkillEffect_IsBuffEffect(t *testing.T) {
	buffEffect := SkillEffect{
		Target: TargetSelf,
		ColumnSpec: &EffectColumnSpec{
			Column:   ColDamageBonus,
			Value:    10.0,
			Duration: 10.0,
		},
	}
	if !buffEffect.IsBuffEffect() {
		t.Error("自身対象のColumnSpec効果はバフ効果であるべきです")
	}

	debuffEffect := SkillEffect{
		Target: TargetEnemy,
		ColumnSpec: &EffectColumnSpec{
			Column:   ColDamageCut,
			Value:    -10.0,
			Duration: 8.0,
		},
	}
	if debuffEffect.IsBuffEffect() {
		t.Error("敵対象のColumnSpec効果はバフ効果ではないべきです")
	}
}

// TestSkillEffect_IsDebuffEffect はデバフ効果の判定をテストします。
func TestSkillEffect_IsDebuffEffect(t *testing.T) {
	debuffEffect := SkillEffect{
		Target: TargetEnemy,
		ColumnSpec: &EffectColumnSpec{
			Column:   ColDamageCut,
			Value:    -10.0,
			Duration: 8.0,
		},
	}
	if !debuffEffect.IsDebuffEffect() {
		t.Error("敵対象のColumnSpec効果はデバフ効果であるべきです")
	}

	buffEffect := SkillEffect{
		Target: TargetSelf,
		ColumnSpec: &EffectColumnSpec{
			Column:   ColDamageBonus,
			Value:    10.0,
			Duration: 10.0,
		},
	}
	if buffEffect.IsDebuffEffect() {
		t.Error("自身対象のColumnSpec効果はデバフ効果ではないべきです")
	}
}
