// Package inventory はインベントリ管理のユースケースを提供します。
// このファイルはInventoryManagerのテストを定義します。

package inventory

import (
	"testing"
)

// TestNewInventoryManager はInventoryManagerの生成をテストします。
func TestNewInventoryManager(t *testing.T) {
	manager := NewInventoryManager()

	if manager == nil {
		t.Fatal("NewInventoryManager() returned nil")
	}

	if manager.Cores() == nil {
		t.Error("Cores() returned nil")
	}

	if manager.Skills() == nil {
		t.Error("Skills() returned nil")
	}
}

// TestInventoryManager_AddCore_NewCore はコア新規追加のテストです。
func TestInventoryManager_AddCore_NewCore(t *testing.T) {
	manager := NewInventoryManager()

	// 新規コア追加
	updated := manager.AddCore("core_balance")

	if !updated {
		t.Error("AddCore() should return true for new core")
	}

	// 追加されたことを確認
	if !manager.Cores().HasCore("core_balance") {
		t.Error("Core should exist after AddCore")
	}
}

// TestInventoryManager_AddCore_Duplicate は重複追加のテストです。
func TestInventoryManager_AddCore_Duplicate(t *testing.T) {
	manager := NewInventoryManager()

	// 最初のコア追加
	manager.AddCore("core_balance")

	// 同じコアを再度追加（更新されないはず）
	updated := manager.AddCore("core_balance")

	if updated {
		t.Error("AddCore() should return false for duplicate core")
	}

	// 保有していることを確認
	if !manager.Cores().HasCore("core_balance") {
		t.Error("Core should still exist after duplicate AddCore")
	}
}

// TestInventoryManager_AddSkill_NewSkill はスキル新規追加のテストです。
func TestInventoryManager_AddSkill_NewSkill(t *testing.T) {
	manager := NewInventoryManager()

	// 新規スキル追加（チェイン効果あり）
	manager.AddSkill("skill_slash", "chain_fire")

	// 追加されたことを確認
	if !manager.Skills().HasSkill("skill_slash") {
		t.Error("Skill should exist after AddSkill")
	}

	// チェイン効果が登録されていることを確認
	if !manager.Skills().HasChainVariation("skill_slash", "chain_fire") {
		t.Error("Chain variation should exist after AddSkill")
	}
}

// TestInventoryManager_AddSkill_WithoutChainEffect はチェイン効果なしスキルのテストです。
func TestInventoryManager_AddSkill_WithoutChainEffect(t *testing.T) {
	manager := NewInventoryManager()

	// チェイン効果なしでスキル追加
	manager.AddSkill("skill_heal", "")

	// 追加されたことを確認
	if !manager.Skills().HasSkill("skill_heal") {
		t.Error("Skill should exist after AddSkill")
	}

	// チェイン効果は空のはず
	variations := manager.Skills().GetChainVariations("skill_heal")
	if len(variations) != 0 {
		t.Errorf("GetChainVariations() should return empty slice, got %v", variations)
	}
}

// TestInventoryManager_AddSkill_AdditionalChainEffect は追加チェイン効果のテストです。
func TestInventoryManager_AddSkill_AdditionalChainEffect(t *testing.T) {
	manager := NewInventoryManager()

	// 最初のスキル追加
	manager.AddSkill("skill_slash", "chain_fire")

	// 同じスキルに別のチェイン効果を追加
	manager.AddSkill("skill_slash", "chain_ice")

	// 両方のチェイン効果が登録されていることを確認
	if !manager.Skills().HasChainVariation("skill_slash", "chain_fire") {
		t.Error("First chain variation should exist")
	}

	if !manager.Skills().HasChainVariation("skill_slash", "chain_ice") {
		t.Error("Second chain variation should exist")
	}
}

// TestInventoryManager_GetOwnedCoreTypes は保有コアTypeID一覧のテストです。
func TestInventoryManager_GetOwnedCoreTypes(t *testing.T) {
	manager := NewInventoryManager()

	// 複数のコアを追加
	manager.AddCore("core_balance")
	manager.AddCore("core_attack")
	manager.AddCore("core_defense")

	// TypeID一覧を取得
	typeIDs := manager.GetOwnedCoreTypes()

	// 3つのTypeIDが含まれているはず
	if len(typeIDs) != 3 {
		t.Errorf("GetOwnedCoreTypes() returned %d items, want 3", len(typeIDs))
	}

	// 全てのTypeIDが含まれているか確認
	typeIDSet := make(map[string]bool)
	for _, id := range typeIDs {
		typeIDSet[id] = true
	}

	expectedIDs := []string{"core_balance", "core_attack", "core_defense"}
	for _, expected := range expectedIDs {
		if !typeIDSet[expected] {
			t.Errorf("GetOwnedCoreTypes() should contain %s", expected)
		}
	}
}

// TestInventoryManager_GetOwnedSkillTypes は保有スキルTypeID一覧のテストです。
func TestInventoryManager_GetOwnedSkillTypes(t *testing.T) {
	manager := NewInventoryManager()

	// 複数のスキルを追加
	manager.AddSkill("skill_slash", "chain_fire")
	manager.AddSkill("skill_heal", "")
	manager.AddSkill("skill_guard", "chain_shield")

	// TypeID一覧を取得
	typeIDs := manager.GetOwnedSkillTypes()

	// 3つのTypeIDが含まれているはず
	if len(typeIDs) != 3 {
		t.Errorf("GetOwnedSkillTypes() returned %d items, want 3", len(typeIDs))
	}

	// 全てのTypeIDが含まれているか確認
	typeIDSet := make(map[string]bool)
	for _, id := range typeIDs {
		typeIDSet[id] = true
	}

	expectedIDs := []string{"skill_slash", "skill_heal", "skill_guard"}
	for _, expected := range expectedIDs {
		if !typeIDSet[expected] {
			t.Errorf("GetOwnedSkillTypes() should contain %s", expected)
		}
	}
}

// TestInventoryManager_GetOwnedCoreTypes_Empty は空の状態でのテストです。
func TestInventoryManager_GetOwnedCoreTypes_Empty(t *testing.T) {
	manager := NewInventoryManager()

	typeIDs := manager.GetOwnedCoreTypes()

	if len(typeIDs) != 0 {
		t.Errorf("GetOwnedCoreTypes() should return empty slice for new manager, got %v", typeIDs)
	}
}

// TestInventoryManager_GetOwnedSkillTypes_Empty は空の状態でのテストです。
func TestInventoryManager_GetOwnedSkillTypes_Empty(t *testing.T) {
	manager := NewInventoryManager()

	typeIDs := manager.GetOwnedSkillTypes()

	if len(typeIDs) != 0 {
		t.Errorf("GetOwnedSkillTypes() should return empty slice for new manager, got %v", typeIDs)
	}
}

// TestInventoryManager_LoadFromSaveData はセーブデータからの復元をテストします。
func TestInventoryManager_LoadFromSaveData(t *testing.T) {
	// コアインベントリとスキルインベントリを事前に準備
	manager := NewInventoryManagerWithInventories(nil, nil)

	if manager == nil {
		t.Fatal("NewInventoryManagerWithInventories() returned nil")
	}

	// nilの場合は新規インベントリが作成されるはず
	if manager.Cores() == nil {
		t.Error("Cores() should not be nil even with nil input")
	}

	if manager.Skills() == nil {
		t.Error("Skills() should not be nil even with nil input")
	}
}
