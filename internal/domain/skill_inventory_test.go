// Package domain はゲームのドメインモデルを定義します。
// このテストファイルはSkillInventoryのテストを含みます。

package domain

import (
	"sort"
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

	// 初期状態ではチェイン効果バリエーションが空
	if len(ownership.ChainVariations) != 0 {
		t.Errorf("初期状態でChainVariationsが空でない: got %d, want 0", len(ownership.ChainVariations))
	}
}

// TestSkillOwnership_AddChainVariation はチェイン効果バリエーションが正しく追加されることを確認します。
func TestSkillOwnership_AddChainVariation(t *testing.T) {
	ownership := NewSkillOwnership()

	// チェイン効果を追加
	ownership.AddChainVariation("chain_fire")

	if !ownership.HasChainVariation("chain_fire") {
		t.Error("追加したチェイン効果が見つかりません")
	}
}

// TestSkillOwnership_AddChainVariation_重複追加 は同じチェイン効果の重複追加が無視されることを確認します。
func TestSkillOwnership_AddChainVariation_重複追加(t *testing.T) {
	ownership := NewSkillOwnership()

	ownership.AddChainVariation("chain_fire")
	ownership.AddChainVariation("chain_fire")

	// 重複しても1つのみ
	variations := ownership.GetChainVariations()
	if len(variations) != 1 {
		t.Errorf("重複追加でバリエーション数が増えている: got %d, want 1", len(variations))
	}
}

// TestSkillOwnership_AddChainVariation_空文字列 は空のチェイン効果IDが追加されないことを確認します。
func TestSkillOwnership_AddChainVariation_空文字列(t *testing.T) {
	ownership := NewSkillOwnership()

	ownership.AddChainVariation("")

	variations := ownership.GetChainVariations()
	if len(variations) != 0 {
		t.Errorf("空文字列のチェイン効果が追加されている: got %d, want 0", len(variations))
	}
}

// TestSkillOwnership_GetChainVariations は保有しているチェイン効果の一覧が取得できることを確認します。
func TestSkillOwnership_GetChainVariations(t *testing.T) {
	ownership := NewSkillOwnership()

	ownership.AddChainVariation("chain_fire")
	ownership.AddChainVariation("chain_ice")
	ownership.AddChainVariation("chain_thunder")

	variations := ownership.GetChainVariations()
	if len(variations) != 3 {
		t.Errorf("チェイン効果バリエーション数が正しくない: got %d, want 3", len(variations))
	}

	// ソートして比較
	sort.Strings(variations)
	expected := []string{"chain_fire", "chain_ice", "chain_thunder"}
	sort.Strings(expected)

	for i, v := range expected {
		if variations[i] != v {
			t.Errorf("チェイン効果が正しくない: got %s, want %s", variations[i], v)
		}
	}
}

// TestSkillOwnership_HasChainVariation_保有 は保有しているチェイン効果に対してtrueが返されることを確認します。
func TestSkillOwnership_HasChainVariation_保有(t *testing.T) {
	ownership := NewSkillOwnership()
	ownership.AddChainVariation("chain_fire")

	if !ownership.HasChainVariation("chain_fire") {
		t.Error("保有しているチェイン効果に対してtrueが返されるべきです")
	}
}

// TestSkillOwnership_HasChainVariation_未保有 は未保有のチェイン効果に対してfalseが返されることを確認します。
func TestSkillOwnership_HasChainVariation_未保有(t *testing.T) {
	ownership := NewSkillOwnership()

	if ownership.HasChainVariation("nonexistent") {
		t.Error("未保有のチェイン効果に対してfalseが返されるべきです")
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

	inv.AddSkill("fire_attack", "")

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

// TestSkillInventory_AddSkill_チェイン効果付き はチェイン効果付きスキル取得でバリエーションが追加されることを確認します。
func TestSkillInventory_AddSkill_チェイン効果付き(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")

	ownership := inv.GetOwnedSkills()["fire_attack"]
	if ownership == nil {
		t.Fatal("追加したスキルがGetOwnedSkillsで取得できません")
	}

	if !ownership.HasChainVariation("chain_burn") {
		t.Error("追加したチェイン効果が見つかりません")
	}
}

// TestSkillInventory_AddSkill_再取得でチェイン効果追加 は同一スキル再取得で新しいチェイン効果が追加されることを確認します。
func TestSkillInventory_AddSkill_再取得でチェイン効果追加(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")
	inv.AddSkill("fire_attack", "chain_ember")

	ownership := inv.GetOwnedSkills()["fire_attack"]
	if ownership == nil {
		t.Fatal("追加したスキルがGetOwnedSkillsで取得できません")
	}

	if !ownership.HasChainVariation("chain_burn") {
		t.Error("最初のチェイン効果が見つかりません")
	}
	if !ownership.HasChainVariation("chain_ember") {
		t.Error("再取得で追加したチェイン効果が見つかりません")
	}

	variations := inv.GetChainVariations("fire_attack")
	if len(variations) != 2 {
		t.Errorf("チェイン効果バリエーション数が正しくない: got %d, want 2", len(variations))
	}
}

// TestSkillInventory_AddSkill_チェイン効果なし はチェイン効果なしスキルの取得で保有状態のみ更新されることを確認します。
func TestSkillInventory_AddSkill_チェイン効果なし(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("basic_attack", "")

	if !inv.HasSkill("basic_attack") {
		t.Error("スキル追加後にHasSkillがtrueを返すべきです")
	}

	variations := inv.GetChainVariations("basic_attack")
	if len(variations) != 0 {
		t.Errorf("チェイン効果なしスキルでバリエーションが空でない: got %d, want 0", len(variations))
	}
}

// TestSkillInventory_AddSkill_空TypeID は空のTypeIDが拒否されることを確認します。
func TestSkillInventory_AddSkill_空TypeID(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("", "chain_test")

	owned := inv.GetOwnedSkills()
	if len(owned) != 0 {
		t.Error("空TypeIDのスキルは保存されるべきではありません")
	}
}

// TestSkillInventory_GetOwnedSkills は保有している全スキル情報を取得できることを確認します。
func TestSkillInventory_GetOwnedSkills(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")
	inv.AddSkill("ice_attack", "")
	inv.AddSkill("thunder_attack", "chain_shock")

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

// TestSkillInventory_GetChainVariations_保有スキル は保有スキルのチェイン効果一覧が取得できることを確認します。
func TestSkillInventory_GetChainVariations_保有スキル(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")
	inv.AddSkill("fire_attack", "chain_ember")

	variations := inv.GetChainVariations("fire_attack")

	if len(variations) != 2 {
		t.Errorf("チェイン効果バリエーション数が正しくない: got %d, want 2", len(variations))
	}
}

// TestSkillInventory_GetChainVariations_未保有スキル は未保有スキルに対して空のリストが返されることを確認します。
func TestSkillInventory_GetChainVariations_未保有スキル(t *testing.T) {
	inv := NewSkillInventory()

	variations := inv.GetChainVariations("nonexistent")

	if len(variations) != 0 {
		t.Errorf("未保有スキルのチェイン効果は空であるべきです: got %d", len(variations))
	}
}

// TestSkillInventory_HasSkill_保有 は保有しているスキルに対してtrueが返されることを確認します。
func TestSkillInventory_HasSkill_保有(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "")

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

// TestSkillInventory_HasChainVariation_保有 は保有しているチェイン効果バリエーションに対してtrueが返されることを確認します。
func TestSkillInventory_HasChainVariation_保有(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")

	if !inv.HasChainVariation("fire_attack", "chain_burn") {
		t.Error("保有しているチェイン効果に対してtrueが返されるべきです")
	}
}

// TestSkillInventory_HasChainVariation_未保有スキル は未保有スキルのチェイン効果に対してfalseが返されることを確認します。
func TestSkillInventory_HasChainVariation_未保有スキル(t *testing.T) {
	inv := NewSkillInventory()

	if inv.HasChainVariation("nonexistent", "chain_burn") {
		t.Error("未保有スキルのチェイン効果に対してfalseが返されるべきです")
	}
}

// TestSkillInventory_HasChainVariation_未保有チェイン効果 は未保有のチェイン効果に対してfalseが返されることを確認します。
func TestSkillInventory_HasChainVariation_未保有チェイン効果(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")

	if inv.HasChainVariation("fire_attack", "chain_nonexistent") {
		t.Error("未保有のチェイン効果に対してfalseが返されるべきです")
	}
}

// TestSkillInventory_GetOwnedSkills_コピー返却 は返されるマップが内部状態のコピーであることを確認します。
func TestSkillInventory_GetOwnedSkills_コピー返却(t *testing.T) {
	inv := NewSkillInventory()

	inv.AddSkill("fire_attack", "chain_burn")

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
	inv.AddSkill("fire_attack", "chain_burn")
	inv.AddSkill("ice_attack", "chain_freeze")

	// それぞれ独立して管理されている
	if !inv.HasChainVariation("fire_attack", "chain_burn") {
		t.Error("fire_attackのchain_burnが見つかりません")
	}
	if !inv.HasChainVariation("ice_attack", "chain_freeze") {
		t.Error("ice_attackのchain_freezeが見つかりません")
	}

	// fire_attackにchain_freezeがないこと
	if inv.HasChainVariation("fire_attack", "chain_freeze") {
		t.Error("fire_attackにchain_freezeが誤って設定されています")
	}

	// ice_attackにchain_burnがないこと
	if inv.HasChainVariation("ice_attack", "chain_burn") {
		t.Error("ice_attackにchain_burnが誤って設定されています")
	}
}
