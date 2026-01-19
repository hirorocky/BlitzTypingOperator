// Package domain はゲームのドメインモデルを定義します。
// このテストファイルはCoreInventoryのテストを含みます。

package domain

import (
	"testing"
)

// ==================== CoreInventory テスト ====================

// TestCoreInventory_NewCoreInventory は新しいCoreInventoryが正しく作成されることを確認します。
func TestCoreInventory_NewCoreInventory(t *testing.T) {
	inv := NewCoreInventory()

	if inv == nil {
		t.Fatal("NewCoreInventoryがnilを返しました")
	}

	// 初期状態では空であること
	owned := inv.GetOwnedCores()
	if len(owned) != 0 {
		t.Errorf("初期状態で保有コア数が0でない: got %d, want 0", len(owned))
	}
}

// TestCoreInventory_AddCore_新規追加 は新規TypeIDのコア追加で保存されることを確認します。
func TestCoreInventory_AddCore_新規追加(t *testing.T) {
	inv := NewCoreInventory()

	updated := inv.AddCore("attack_balance", 5)

	if !updated {
		t.Error("新規コア追加でtrueが返されるべきです")
	}

	maxLevel := inv.GetMaxLevel("attack_balance")
	if maxLevel != 5 {
		t.Errorf("最大レベルが正しくない: got %d, want 5", maxLevel)
	}
}

// TestCoreInventory_AddCore_レベル更新 は既存より高いレベルのコア取得で更新されることを確認します。
func TestCoreInventory_AddCore_レベル更新(t *testing.T) {
	inv := NewCoreInventory()

	// 初回追加
	inv.AddCore("attack_balance", 5)

	// より高いレベルで更新
	updated := inv.AddCore("attack_balance", 10)

	if !updated {
		t.Error("より高いレベルのコア追加でtrueが返されるべきです")
	}

	maxLevel := inv.GetMaxLevel("attack_balance")
	if maxLevel != 10 {
		t.Errorf("最大レベルが更新されていない: got %d, want 10", maxLevel)
	}
}

// TestCoreInventory_AddCore_レベル据え置き は既存以下のレベルのコア取得で変更されないことを確認します。
func TestCoreInventory_AddCore_レベル据え置き(t *testing.T) {
	inv := NewCoreInventory()

	// 初回追加
	inv.AddCore("attack_balance", 10)

	// 同じレベルで追加
	updated := inv.AddCore("attack_balance", 10)
	if updated {
		t.Error("同じレベルのコア追加でfalseが返されるべきです")
	}

	// より低いレベルで追加
	updated = inv.AddCore("attack_balance", 5)
	if updated {
		t.Error("より低いレベルのコア追加でfalseが返されるべきです")
	}

	// レベルが変わっていないこと
	maxLevel := inv.GetMaxLevel("attack_balance")
	if maxLevel != 10 {
		t.Errorf("最大レベルが変更されている: got %d, want 10", maxLevel)
	}
}

// TestCoreInventory_GetMaxLevel_未保有 は未保有TypeIDに対して0が返されることを確認します。
func TestCoreInventory_GetMaxLevel_未保有(t *testing.T) {
	inv := NewCoreInventory()

	maxLevel := inv.GetMaxLevel("nonexistent")
	if maxLevel != 0 {
		t.Errorf("未保有TypeIDの最大レベルは0であるべきです: got %d", maxLevel)
	}
}

// TestCoreInventory_GetOwnedCores は保有している全コア情報を取得できることを確認します。
func TestCoreInventory_GetOwnedCores(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance", 5)
	inv.AddCore("healer", 10)
	inv.AddCore("paladin", 3)

	owned := inv.GetOwnedCores()

	if len(owned) != 3 {
		t.Errorf("保有コア数が正しくない: got %d, want 3", len(owned))
	}

	if owned["attack_balance"] != 5 {
		t.Errorf("attack_balanceの最大レベルが正しくない: got %d, want 5", owned["attack_balance"])
	}
	if owned["healer"] != 10 {
		t.Errorf("healerの最大レベルが正しくない: got %d, want 10", owned["healer"])
	}
	if owned["paladin"] != 3 {
		t.Errorf("paladinの最大レベルが正しくない: got %d, want 3", owned["paladin"])
	}
}

// TestCoreInventory_HasCore_保有 は保有しているTypeIDに対してtrueが返されることを確認します。
func TestCoreInventory_HasCore_保有(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance", 5)

	if !inv.HasCore("attack_balance") {
		t.Error("保有しているコアに対してtrueが返されるべきです")
	}
}

// TestCoreInventory_HasCore_未保有 は未保有TypeIDに対してfalseが返されることを確認します。
func TestCoreInventory_HasCore_未保有(t *testing.T) {
	inv := NewCoreInventory()

	if inv.HasCore("nonexistent") {
		t.Error("未保有のコアに対してfalseが返されるべきです")
	}
}

// TestCoreInventory_AddCore_レベル検証 はレベルが1以上の正整数であることを検証します。
func TestCoreInventory_AddCore_レベル検証(t *testing.T) {
	inv := NewCoreInventory()

	// レベル0は無効
	updated := inv.AddCore("test", 0)
	if updated {
		t.Error("レベル0のコア追加はfalseが返されるべきです")
	}
	if inv.HasCore("test") {
		t.Error("レベル0のコアは保存されるべきではありません")
	}

	// レベル-1は無効
	updated = inv.AddCore("test", -1)
	if updated {
		t.Error("負のレベルのコア追加はfalseが返されるべきです")
	}
	if inv.HasCore("test") {
		t.Error("負のレベルのコアは保存されるべきではありません")
	}
}

// TestCoreInventory_AddCore_空TypeID は空のTypeIDが拒否されることを確認します。
func TestCoreInventory_AddCore_空TypeID(t *testing.T) {
	inv := NewCoreInventory()

	updated := inv.AddCore("", 5)
	if updated {
		t.Error("空TypeIDのコア追加はfalseが返されるべきです")
	}

	owned := inv.GetOwnedCores()
	if len(owned) != 0 {
		t.Error("空TypeIDのコアは保存されるべきではありません")
	}
}

// TestCoreInventory_GetOwnedCores_コピー返却 は返されるマップが内部状態のコピーであることを確認します。
func TestCoreInventory_GetOwnedCores_コピー返却(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance", 5)

	owned := inv.GetOwnedCores()
	// 返されたマップを変更
	owned["attack_balance"] = 100
	owned["new_core"] = 50

	// 内部状態は変更されていないこと
	if inv.GetMaxLevel("attack_balance") != 5 {
		t.Error("GetOwnedCoresの返り値の変更が内部状態に影響しています")
	}
	if inv.HasCore("new_core") {
		t.Error("GetOwnedCoresの返り値への追加が内部状態に影響しています")
	}
}

// TestCoreInventory_複数TypeIDの管理 は複数のTypeIDが独立して管理されることを確認します。
func TestCoreInventory_複数TypeIDの管理(t *testing.T) {
	inv := NewCoreInventory()

	// 複数のTypeIDを追加
	inv.AddCore("attack_balance", 5)
	inv.AddCore("healer", 10)

	// それぞれ独立して管理されている
	if inv.GetMaxLevel("attack_balance") != 5 {
		t.Errorf("attack_balanceの最大レベルが正しくない: got %d, want 5", inv.GetMaxLevel("attack_balance"))
	}
	if inv.GetMaxLevel("healer") != 10 {
		t.Errorf("healerの最大レベルが正しくない: got %d, want 10", inv.GetMaxLevel("healer"))
	}

	// 一方を更新しても他方に影響しない
	inv.AddCore("attack_balance", 15)
	if inv.GetMaxLevel("attack_balance") != 15 {
		t.Errorf("attack_balanceの最大レベルが更新されていない: got %d, want 15", inv.GetMaxLevel("attack_balance"))
	}
	if inv.GetMaxLevel("healer") != 10 {
		t.Errorf("healerの最大レベルが変わっている: got %d, want 10", inv.GetMaxLevel("healer"))
	}
}
