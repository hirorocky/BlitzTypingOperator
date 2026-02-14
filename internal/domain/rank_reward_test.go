package domain

import "testing"

// 受け入れ基準1: ランクアップ時にマスタデータで定義された報酬アイテムを表現できる

// TestRankReward_HasItems はRankRewardがアイテムリストを保持できることをテストします。
func TestRankReward_HasItems(t *testing.T) {
	reward := RankReward{
		Rank: 2,
		Items: []RankRewardItem{
			{Category: "core", TypeID: "magic_balance"},
			{Category: "skill", TypeID: "heal_lv1"},
			{Category: "chain_effect", TypeID: "chain_heal"},
		},
	}

	if reward.Rank != 2 {
		t.Errorf("Rankが期待と異なる: got %d, want 2", reward.Rank)
	}
	if len(reward.Items) != 3 {
		t.Errorf("Items数が期待と異なる: got %d, want 3", len(reward.Items))
	}
}

// TestRankReward_EmptyItems は空の報酬リスト（報酬なしランクアップ）をテストします。
func TestRankReward_EmptyItems(t *testing.T) {
	reward := RankReward{
		Rank:  3,
		Items: []RankRewardItem{},
	}

	if len(reward.Items) != 0 {
		t.Errorf("Items数が期待と異なる: got %d, want 0", len(reward.Items))
	}
}

// TestRankRewardItem_Categories は各カテゴリ値が正しく保持されることをテストします。
func TestRankRewardItem_Categories(t *testing.T) {
	tests := []struct {
		name     string
		category string
		typeID   string
	}{
		{"コア報酬", "core", "attack_balance"},
		{"スキル報酬", "skill", "heal_lv1"},
		{"チェイン効果報酬", "chain_effect", "damage_bonus"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := RankRewardItem{Category: tt.category, TypeID: tt.typeID}
			if item.Category != tt.category {
				t.Errorf("Categoryが期待と異なる: got %s, want %s", item.Category, tt.category)
			}
			if item.TypeID != tt.typeID {
				t.Errorf("TypeIDが期待と異なる: got %s, want %s", item.TypeID, tt.typeID)
			}
		})
	}
}
