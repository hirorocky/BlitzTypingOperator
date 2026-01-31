// Package domain はゲームのドメインモデルを定義します。
// このテストファイルはCoreInventoryのテストを含みます。

package domain

import (
	"testing"
)

// ==================== CoreInventory テスト ====================
// 新仕様: TypeIDベースのフラグ管理（レベル概念を削除）

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

	updated := inv.AddCore("attack_balance")

	if !updated {
		t.Error("新規コア追加でtrueが返されるべきです")
	}

	if !inv.HasCore("attack_balance") {
		t.Error("コアが保存されていません")
	}
}

// TestCoreInventory_AddCore_重複追加 は既存TypeIDの重複追加でfalseが返されることを確認します。
func TestCoreInventory_AddCore_重複追加(t *testing.T) {
	inv := NewCoreInventory()

	// 初回追加
	inv.AddCore("attack_balance")

	// 同じTypeIDで再追加
	updated := inv.AddCore("attack_balance")

	if updated {
		t.Error("既存TypeIDの重複追加でfalseが返されるべきです")
	}

	// 保有状態は変わらない
	if !inv.HasCore("attack_balance") {
		t.Error("コアが保存されていません")
	}
}

// TestCoreInventory_GetOwnedCores は保有している全TypeIDを取得できることを確認します。
func TestCoreInventory_GetOwnedCores(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance")
	inv.AddCore("healer")
	inv.AddCore("paladin")

	owned := inv.GetOwnedCores()

	if len(owned) != 3 {
		t.Errorf("保有コア数が正しくない: got %d, want 3", len(owned))
	}

	// 返されたスライスに各TypeIDが含まれていることを確認
	found := map[string]bool{}
	for _, typeID := range owned {
		found[typeID] = true
	}

	if !found["attack_balance"] {
		t.Error("attack_balanceが保有リストに含まれていません")
	}
	if !found["healer"] {
		t.Error("healerが保有リストに含まれていません")
	}
	if !found["paladin"] {
		t.Error("paladinが保有リストに含まれていません")
	}
}

// TestCoreInventory_HasCore_保有 は保有しているTypeIDに対してtrueが返されることを確認します。
func TestCoreInventory_HasCore_保有(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance")

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

// TestCoreInventory_AddCore_空TypeID は空のTypeIDが拒否されることを確認します。
func TestCoreInventory_AddCore_空TypeID(t *testing.T) {
	inv := NewCoreInventory()

	updated := inv.AddCore("")
	if updated {
		t.Error("空TypeIDのコア追加はfalseが返されるべきです")
	}

	owned := inv.GetOwnedCores()
	if len(owned) != 0 {
		t.Error("空TypeIDのコアは保存されるべきではありません")
	}
}

// TestCoreInventory_GetOwnedCores_コピー返却 は返されるスライスが内部状態に影響しないことを確認します。
func TestCoreInventory_GetOwnedCores_コピー返却(t *testing.T) {
	inv := NewCoreInventory()

	inv.AddCore("attack_balance")

	owned := inv.GetOwnedCores()
	// 返されたスライスを変更
	if len(owned) > 0 {
		owned[0] = "modified"
	}

	// 内部状態は変更されていないこと
	if !inv.HasCore("attack_balance") {
		t.Error("GetOwnedCoresの返り値の変更が内部状態に影響しています")
	}
}

// TestCoreInventory_複数TypeIDの管理 は複数のTypeIDが独立して管理されることを確認します。
func TestCoreInventory_複数TypeIDの管理(t *testing.T) {
	inv := NewCoreInventory()

	// 複数のTypeIDを追加
	inv.AddCore("attack_balance")
	inv.AddCore("healer")

	// それぞれ独立して管理されている
	if !inv.HasCore("attack_balance") {
		t.Error("attack_balanceが保存されていません")
	}
	if !inv.HasCore("healer") {
		t.Error("healerが保存されていません")
	}

	// 保有数を確認
	owned := inv.GetOwnedCores()
	if len(owned) != 2 {
		t.Errorf("保有コア数が正しくない: got %d, want 2", len(owned))
	}
}
