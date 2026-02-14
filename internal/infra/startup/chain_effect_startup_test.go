// Package startup は初回起動時の初期化処理のテストを提供します。

package startup

import (
	"testing"
)

// 受け入れ基準17: 初期エージェントが異なるコアを持つ

func TestInitializeNewGame_DifferentCoresForEachSlot(t *testing.T) {
	initializer := NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()

	// 3スロットが異なるコアを持つことを確認
	coreIDs := make(map[string]bool)
	for i, slot := range saveData.Player.AgentSlots {
		if slot.CoreTypeID == "" {
			t.Errorf("スロット%dのコアが空です", i)
			continue
		}
		if coreIDs[slot.CoreTypeID] {
			t.Errorf("スロット%dのコア %s が他のスロットと重複しています", i, slot.CoreTypeID)
		}
		coreIDs[slot.CoreTypeID] = true
	}

	if len(coreIDs) != 3 {
		t.Errorf("3種類の異なるコアが必要: got %d種類", len(coreIDs))
	}
}

func TestInitializeNewGame_CoreInventoryContainsAllSlotCores(t *testing.T) {
	initializer := NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()

	// コアインベントリに3種類のコアが含まれること
	coreInvSet := make(map[string]bool)
	for _, c := range saveData.Inventory.UniqueCores.Cores {
		coreInvSet[c] = true
	}

	for i, slot := range saveData.Player.AgentSlots {
		if !coreInvSet[slot.CoreTypeID] {
			t.Errorf("スロット%dのコア %s がインベントリに含まれていません", i, slot.CoreTypeID)
		}
	}
}

// 受け入れ基準18: 初期コアと初期スキルの互換性

func TestInitializeNewGame_CoreSkillCompatibility(t *testing.T) {
	initializer := NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()

	// テスト用コアのAllowedTagsマップ
	coreAllowedTags := map[string]map[string]bool{
		"all_rounder":     {"physical_low": true, "magic_low": true, "heal_low": true, "buff_low": true, "debuff_low": true},
		"healer":          {"heal_low": true, "heal_mid": true, "heal_high": true},
		"quick_recoverer": {"heal_low": true, "heal_mid": true, "buff_low": true},
	}

	// テスト用スキルタグマップ
	skillTags := map[string][]string{
		"physical_strike_lv1": {"physical_low"},
		"heal_lv1":            {"heal_low"},
		"attack_buff_lv1":     {"buff_low"},
		"str_buff_lv1":        {"buff_low"},
	}

	for i, slot := range saveData.Player.AgentSlots {
		allowedSet, ok := coreAllowedTags[slot.CoreTypeID]
		if !ok {
			// テストデータに登録されていないコアの場合はスキップ
			continue
		}

		for j, skillSlot := range slot.Skills {
			if skillSlot.TypeID == "" {
				continue
			}

			tags, exists := skillTags[skillSlot.TypeID]
			if !exists {
				continue
			}

			compatible := false
			for _, tag := range tags {
				if allowedSet[tag] {
					compatible = true
					break
				}
			}

			if !compatible {
				t.Errorf("スロット%dのスキル%d (%s) がコア %s と非互換です",
					i, j, skillSlot.TypeID, slot.CoreTypeID)
			}
		}
	}
}
