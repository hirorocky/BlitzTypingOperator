package domain_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// 受け入れ基準1: ChainEffectInventoryのユニーク管理
func TestChainEffectInventory_AddAndHas(t *testing.T) {
	inv := domain.NewChainEffectInventory()

	added := inv.AddChainEffect("damage_bonus")
	if !added {
		t.Error("新規追加でtrueが返るべき")
	}

	if !inv.HasChainEffect("damage_bonus") {
		t.Error("追加したチェイン効果が保有されているべき")
	}
}

func TestChainEffectInventory_MultipleTypes(t *testing.T) {
	inv := domain.NewChainEffectInventory()

	inv.AddChainEffect("damage_bonus")
	inv.AddChainEffect("armor_pierce")
	inv.AddChainEffect("heal_boost")

	owned := inv.GetOwnedChainEffects()
	if len(owned) != 3 {
		t.Errorf("保有数が不正: got %d, want 3", len(owned))
	}

	// ソート済みであること
	if owned[0] != "armor_pierce" || owned[1] != "damage_bonus" || owned[2] != "heal_boost" {
		t.Errorf("ソート順が不正: %v", owned)
	}
}

// 受け入れ基準2: 既保有TypeIDの追加は無視
func TestChainEffectInventory_DuplicateAddReturnsFalse(t *testing.T) {
	inv := domain.NewChainEffectInventory()

	inv.AddChainEffect("damage_bonus")
	added := inv.AddChainEffect("damage_bonus")

	if added {
		t.Error("重複追加でfalseが返るべき")
	}

	owned := inv.GetOwnedChainEffects()
	if len(owned) != 1 {
		t.Errorf("重複追加後の保有数が不正: got %d, want 1", len(owned))
	}
}

func TestChainEffectInventory_EmptyTypeIDReturnsFalse(t *testing.T) {
	inv := domain.NewChainEffectInventory()

	added := inv.AddChainEffect("")
	if added {
		t.Error("空文字列でfalseが返るべき")
	}

	owned := inv.GetOwnedChainEffects()
	if len(owned) != 0 {
		t.Errorf("空文字列追加後の保有数が不正: got %d, want 0", len(owned))
	}
}

func TestChainEffectInventory_NewIsEmpty(t *testing.T) {
	inv := domain.NewChainEffectInventory()

	owned := inv.GetOwnedChainEffects()
	if len(owned) != 0 {
		t.Errorf("新規作成時の保有数が不正: got %d, want 0", len(owned))
	}

	if inv.HasChainEffect("anything") {
		t.Error("新規作成時にHasChainEffectがfalseを返すべき")
	}
}
