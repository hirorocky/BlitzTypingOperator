// Package domain はゲームのドメインモデルを定義します。
// このテストファイルはSkillInventoryのテストを含みます。

package domain

import (
	"testing"
)

// ==================== SkillOwnership テスト ====================

// TestSkillOwnership_NewSkillOwnership は新しいSkillOwnershipが正しく作成されることを確認します。
func TestSkillOwnership_NewSkillOwnership(t *testing.T) {
	ownership := NewSkillOwnership()

	if ownership == nil {
		t.Fatal("NewSkillOwnershipがnilを返しました")
	}

	// 初期状態では保有フラグがfalse
	if ownership.Owned {
		t.Error("初期状態でOwnedがtrueになっています")
	}
}

// ==================== SkillInventory テスト ====================

// TestSkillInventory_NewSkillInventory は新しいSkillInventoryが正しく作成されることを確認します。
func TestSkillInventory_NewSkillInventory(t *testing.T) {
	inv := NewSkillInventory()

	if inv == nil {
		t.Fatal("NewSkillInventoryがnilを返しました")
	}

	// 初期状態では空であること
	owned := inv.GetOwnedSkills()
	if len(owned) != 0 {
		t.Errorf("初期状態で保有スキル数が0でない: got %d, want 0", len(owned))
	}
}

// TestSkillInventory_AddSkill_初回取得 は初回スキル取得で保有状態がtrueになることを確認します。
func TestSkillInventory_AddSkill_初回取得(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack")

	if !inv.HasSkill("fire_attack") {
		t.Error("スキル追加後にHasSkillがtrueを返すべきです")
	}

	ownership := inv.GetOwnedSkills()["fire_attack"]
	if ownership == nil {
		t.Fatal("追加したスキルがGetOwnedSkillsで取得できません")
	}
	if !ownership.Owned {
		t.Error("追加したスキルのOwnedがtrueであるべきです")
	}
}

// TestSkillInventory_AddSkill_空TypeID は空のTypeIDが拒否されることを確認します。
func TestSkillInventory_AddSkill_空TypeID(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("")

	owned := inv.GetOwnedSkills()
	if len(owned) != 0 {
		t.Error("空TypeIDのスキルは保存されるべきではありません")
	}
}

// TestSkillInventory_GetOwnedSkills は保有している全スキル情報を取得できることを確認します。
func TestSkillInventory_GetOwnedSkills(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack")
	inv.AddSkill("ice_attack")
	inv.AddSkill("thunder_attack")

	owned := inv.GetOwnedSkills()

	if len(owned) != 3 {
		t.Errorf("保有スキル数が正しくない: got %d, want 3", len(owned))
	}

	if owned["fire_attack"] == nil {
		t.Error("fire_attackが保有リストに含まれていません")
	}
	if owned["ice_attack"] == nil {
		t.Error("ice_attackが保有リストに含まれていません")
	}
	if owned["thunder_attack"] == nil {
		t.Error("thunder_attackが保有リストに含まれていません")
	}
}

// TestSkillInventory_HasSkill_保有 は保有しているスキルに対してtrueが返されることを確認します。
func TestSkillInventory_HasSkill_保有(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack")

	if !inv.HasSkill("fire_attack") {
		t.Error("保有しているスキルに対してtrueが返されるべきです")
	}
}

// TestSkillInventory_HasSkill_未保有 は未保有スキルに対してfalseが返されることを確認します。
func TestSkillInventory_HasSkill_未保有(t *testing.T) {
	inv := NewSkillInventory()

	if inv.HasSkill("nonexistent") {
		t.Error("未保有のスキルに対してfalseが返されるべきです")
	}
}

// TestSkillInventory_GetOwnedSkills_コピー返却 は返されるマップが内部状態のコピーであることを確認します。
func TestSkillInventory_GetOwnedSkills_コピー返却(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack")

	owned := inv.GetOwnedSkills()
	// 返されたマップを変更
	owned["fire_attack"].Owned = false
	delete(owned, "fire_attack")
	owned["new_skill"] = NewSkillOwnership()

	// 内部状態は変更されていないこと
	if !inv.HasSkill("fire_attack") {
		t.Error("GetOwnedSkillsの返り値の削除が内部状態に影響しています")
	}
	if inv.HasSkill("new_skill") {
		t.Error("GetOwnedSkillsの返り値への追加が内部状態に影響しています")
	}
}

// TestSkillInventory_複数スキルの管理 は複数のスキルが独立して管理されることを確認します。
func TestSkillInventory_複数スキルの管理(t *testing.T) {
	inv := NewSkillInventory()

	// 複数のスキルを追加
	inv.AddSkill("fire_attack")
	inv.AddSkill("ice_attack")

	// それぞれ独立して管理されている
	if !inv.HasSkill("fire_attack") {
		t.Error("fire_attackが見つかりません")
	}
	if !inv.HasSkill("ice_attack") {
		t.Error("ice_attackが見つかりません")
	}
}
