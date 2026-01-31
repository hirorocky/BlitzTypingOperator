// Package integration_test はタスク11の移行を検証する統合テストを提供します。
// 旧システム（synthesize.AgentManager）から新システム（slot.AgentSlotManager）への
// 移行が正しく行われていることを確認します。
package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"
)

// TestTask11_NewSystemCanReplaceOldSystem は新システムが旧システムの機能を
// 代替できることを確認します。
func TestTask11_NewSystemCanReplaceOldSystem(t *testing.T) {
	// テスト用のマスタデータを準備
	coreTypes := map[string]domain.CoreType{
		"warrior": {
			ID:             "warrior",
			Name:           "ウォリアー",
			AllowedTags:    []string{"attack", "physical"},
			PassiveSkillID: "warrior_passive",
			MinDropLevel:   1,
		},
		"mage": {
			ID:             "mage",
			Name:           "メイジ",
			AllowedTags:    []string{"attack", "magic"},
			PassiveSkillID: "mage_passive",
			MinDropLevel:   1,
		},
	}

	skillTypes := map[string]domain.SkillType{
		"slash": {
			ID:          "slash",
			Name:        "斬撃",
			Icon:        "🗡️",
			Tags:        []string{"attack", "physical"},
			Description: "物理攻撃",
		},
		"fireball": {
			ID:          "fireball",
			Name:        "ファイアボール",
			Icon:        "🔥",
			Tags:        []string{"attack", "magic"},
			Description: "魔法攻撃",
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{
		"warrior_passive": {
			ID:   "warrior_passive",
			Name: "ウォリアーの力",
		},
		"mage_passive": {
			ID:   "mage_passive",
			Name: "メイジの知恵",
		},
	}

	t.Run("AgentSlotManagerでコア設定とスキル設定ができる", func(t *testing.T) {
		// 新システムのインベントリを作成
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior", 10)
		invMgr.AddCore("mage", 5)

		// スキルを追加
		invMgr.AddSkill("slash", "")
		invMgr.AddSkill("fireball", "chain_damage")

		// AgentSlotManagerを作成
		slotMgr := slot.NewAgentSlotManager(
			invMgr.Cores(),
			invMgr.Skills(),
			coreTypes,
			skillTypes,
			passiveSkills,
			nil, // chainEffects
		)

		// コアを設定（旧システムの「合成」に相当）
		err := slotMgr.SetCore(0, "warrior", 5)
		if err != nil {
			t.Fatalf("SetCoreに失敗: %v", err)
		}

		// スキルを設定
		err = slotMgr.SetSkill(0, 0, "slash", "")
		if err != nil {
			t.Fatalf("SetSkillに失敗: %v", err)
		}

		// スロットが準備完了であることを確認
		if !slotMgr.IsSlotReady(0) {
			t.Error("スロット0が準備完了でない")
		}
	})

	t.Run("BuildAgentsForBattleでバトル用AgentModelを構築できる", func(t *testing.T) {
		// 新システムのインベントリを作成
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior", 10)

		// スキルを追加
		invMgr.AddSkill("slash", "")

		// AgentSlotManagerを作成
		slotMgr := slot.NewAgentSlotManager(
			invMgr.Cores(),
			invMgr.Skills(),
			coreTypes,
			skillTypes,
			passiveSkills,
			nil, // chainEffects
		)

		// コアとスキルを設定
		_ = slotMgr.SetCore(0, "warrior", 5)
		_ = slotMgr.SetSkill(0, 0, "slash", "")

		// バトル用エージェントを構築（旧システムのGetEquippedAgentsに相当）
		agents := slotMgr.BuildAgentsForBattle()

		if len(agents) != 1 {
			t.Errorf("エージェント数が不正: got %d, want 1", len(agents))
		}

		if agents[0].Core == nil {
			t.Error("エージェントのコアがnil")
		}

		if agents[0].Core.TypeID != "warrior" {
			t.Errorf("コアTypeIDが不正: got %s, want warrior", agents[0].Core.TypeID)
		}

		if agents[0].Core.Level != 5 {
			t.Errorf("コアレベルが不正: got %d, want 5", agents[0].Core.Level)
		}
	})

	t.Run("複数スロットにエージェントを設定できる", func(t *testing.T) {
		// 新システムのインベントリを作成
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior", 10)
		invMgr.AddCore("mage", 5)

		// スキルを追加
		invMgr.AddSkill("slash", "")
		invMgr.AddSkill("fireball", "")

		// AgentSlotManagerを作成
		slotMgr := slot.NewAgentSlotManager(
			invMgr.Cores(),
			invMgr.Skills(),
			coreTypes,
			skillTypes,
			passiveSkills,
			nil, // chainEffects
		)

		// スロット0にウォリアーを設定
		_ = slotMgr.SetCore(0, "warrior", 5)
		_ = slotMgr.SetSkill(0, 0, "slash", "")

		// スロット1にメイジを設定
		_ = slotMgr.SetCore(1, "mage", 3)
		_ = slotMgr.SetSkill(1, 0, "fireball", "")

		// バトル用エージェントを構築
		agents := slotMgr.BuildAgentsForBattle()

		if len(agents) != 2 {
			t.Errorf("エージェント数が不正: got %d, want 2", len(agents))
		}

		// 準備完了スロット数を確認
		if slotMgr.GetReadySlotCount() != 2 {
			t.Errorf("準備完了スロット数が不正: got %d, want 2", slotMgr.GetReadySlotCount())
		}
	})

	t.Run("空スロットはバトルに含まれない", func(t *testing.T) {
		// 新システムのインベントリを作成
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior", 10)

		// スキルを追加
		invMgr.AddSkill("slash", "")

		// AgentSlotManagerを作成
		slotMgr := slot.NewAgentSlotManager(
			invMgr.Cores(),
			invMgr.Skills(),
			coreTypes,
			skillTypes,
			passiveSkills,
			nil, // chainEffects
		)

		// スロット1にのみ設定（スロット0は空）
		_ = slotMgr.SetCore(1, "warrior", 5)
		_ = slotMgr.SetSkill(1, 0, "slash", "")

		// バトル用エージェントを構築
		agents := slotMgr.BuildAgentsForBattle()

		// 空スロットは含まれないので1つだけ
		if len(agents) != 1 {
			t.Errorf("エージェント数が不正: got %d, want 1", len(agents))
		}
	})
}

// TestTask11_InventoryManagerUniqueness は新システムのユニーク管理を確認します。
func TestTask11_InventoryManagerUniqueness(t *testing.T) {
	t.Run("同一TypeIDの高レベルコアで最大レベルが更新される", func(t *testing.T) {
		invMgr := inventory.NewInventoryManager()

		// レベル5のコアを追加
		invMgr.AddCore("warrior", 5)

		// レベル10のコアを追加
		updated := invMgr.AddCore("warrior", 10)

		if !updated {
			t.Error("高レベルコアで更新されるべき")
		}

		// 最大レベルが10であることを確認
		maxLevel := invMgr.Cores().GetMaxLevel("warrior")
		if maxLevel != 10 {
			t.Errorf("最大レベルが不正: got %d, want 10", maxLevel)
		}
	})

	t.Run("同一TypeIDの低レベルコアでは更新されない", func(t *testing.T) {
		invMgr := inventory.NewInventoryManager()

		// レベル10のコアを追加
		invMgr.AddCore("warrior", 10)

		// レベル5のコアを追加（更新されないはず）
		updated := invMgr.AddCore("warrior", 5)

		if updated {
			t.Error("低レベルコアでは更新されないべき")
		}

		// 最大レベルが10のままであることを確認
		maxLevel := invMgr.Cores().GetMaxLevel("warrior")
		if maxLevel != 10 {
			t.Errorf("最大レベルが不正: got %d, want 10", maxLevel)
		}
	})

	t.Run("スキルのチェイン効果バリエーションが蓄積される", func(t *testing.T) {
		invMgr := inventory.NewInventoryManager()

		// チェイン効果なしのスキルを追加
		invMgr.AddSkill("slash", "")

		// チェイン効果付きのスキルを追加
		invMgr.AddSkill("slash", "chain_damage")
		invMgr.AddSkill("slash", "chain_heal")

		// チェイン効果バリエーションを確認
		variations := invMgr.Skills().GetChainVariations("slash")

		// 空文字列は含まれないので、chain_damageとchain_healの2つ
		if len(variations) != 2 {
			t.Errorf("チェイン効果バリエーション数が不正: got %d, want 2", len(variations))
		}
	})
}
