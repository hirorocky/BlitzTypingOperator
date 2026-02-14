// Package rewarding はドロップ・報酬システムのチェイン効果分離テストを提供します。

package rewarding

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// 受け入れ基準16: 報酬でのスキル/チェイン効果分離追加

// TestAddRewardsToInventory_ChainEffectSeparation はチェイン効果付きスキルドロップ時に
// スキルとチェイン効果が別々のインベントリに追加されることをテストします。
func TestAddRewardsToInventory_ChainEffectSeparation(t *testing.T) {
	// Arrange: チェイン効果付きスキルを作成
	effect := domain.NewChainEffect("damage_bonus", domain.ChainEffectDamageBonus, 25)
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "physical_lv1",
		Name: "物理攻撃Lv1",
		Icon: "⚔️",
		Tags: []string{"physical"},
	}, &effect)

	result := &RewardResult{
		IsVictory:     true,
		DroppedSkills: []*domain.SkillModel{skill},
		DroppedCores:  make([]*domain.CoreModel, 0),
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert: スキルがSkillInventoryに追加される
	if !skillInv.HasSkill("physical_lv1") {
		t.Error("スキルがSkillInventoryに追加されるべき")
	}

	// Assert: チェイン効果がChainEffectInventoryに追加される
	if !chainEffectInv.HasChainEffect("damage_bonus") {
		t.Error("チェイン効果がChainEffectInventoryに追加されるべき")
	}
}

// TestAddRewardsToInventory_SkillWithoutChainEffect はチェイン効果なしスキルの場合
// スキルのみがインベントリに追加されることをテストします。
func TestAddRewardsToInventory_SkillWithoutChainEffect(t *testing.T) {
	// Arrange: チェイン効果なしスキル
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "heal_lv1",
		Name: "応急手当Lv1",
		Icon: "💚",
		Tags: []string{"heal"},
	}, nil)

	result := &RewardResult{
		IsVictory:     true,
		DroppedSkills: []*domain.SkillModel{skill},
		DroppedCores:  make([]*domain.CoreModel, 0),
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert: スキルは追加される
	if !skillInv.HasSkill("heal_lv1") {
		t.Error("スキルがSkillInventoryに追加されるべき")
	}

	// Assert: チェイン効果インベントリは空のまま
	if len(chainEffectInv.GetOwnedChainEffects()) != 0 {
		t.Error("チェイン効果なしスキルではChainEffectInventoryに追加されるべきではない")
	}
}

// TestAddRewardsToInventory_DuplicateChainEffectNotDuplicated は同じチェイン効果の重複追加が
// インベントリに重複しないことをテストします。
func TestAddRewardsToInventory_DuplicateChainEffectNotDuplicated(t *testing.T) {
	// Arrange: 同じチェイン効果IDの2つのスキル
	effect1 := domain.NewChainEffect("damage_bonus", domain.ChainEffectDamageBonus, 20)
	effect2 := domain.NewChainEffect("damage_bonus", domain.ChainEffectDamageBonus, 30)

	skill1 := domain.NewSkillFromType(domain.SkillType{
		ID:   "physical_lv1",
		Name: "物理攻撃Lv1",
	}, &effect1)
	skill2 := domain.NewSkillFromType(domain.SkillType{
		ID:   "physical_lv2",
		Name: "物理攻撃Lv2",
	}, &effect2)

	result := &RewardResult{
		IsVictory:     true,
		DroppedSkills: []*domain.SkillModel{skill1, skill2},
		DroppedCores:  make([]*domain.CoreModel, 0),
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert: チェイン効果はユニーク管理なので1つのみ
	owned := chainEffectInv.GetOwnedChainEffects()
	if len(owned) != 1 {
		t.Errorf("同一チェイン効果はユニーク管理で1つのみであるべき: got %d", len(owned))
	}
}

// TestAddRewardsToInventory_CoreAndSkillMixed はコアとスキルの混合ドロップを
// テストします。
func TestAddRewardsToInventory_CoreAndSkillMixed(t *testing.T) {
	// Arrange
	effect := domain.NewChainEffect("heal_amp", domain.ChainEffectHealAmp, 15)
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "heal_lv1",
		Name: "応急手当Lv1",
	}, &effect)
	core := domain.NewCoreWithTypeID("all_rounder", domain.CoreType{ID: "all_rounder"}, domain.PassiveSkill{})

	result := &RewardResult{
		IsVictory:     true,
		DroppedSkills: []*domain.SkillModel{skill},
		DroppedCores:  []*domain.CoreModel{core},
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert: 全て正しく追加される
	if !coreInv.HasCore("all_rounder") {
		t.Error("コアがCoreInventoryに追加されるべき")
	}
	if !skillInv.HasSkill("heal_lv1") {
		t.Error("スキルがSkillInventoryに追加されるべき")
	}
	if !chainEffectInv.HasChainEffect("heal_amp") {
		t.Error("チェイン効果がChainEffectInventoryに追加されるべき")
	}
}
