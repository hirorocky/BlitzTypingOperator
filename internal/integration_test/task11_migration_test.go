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
		invMgr.AddCore("warrior")
		invMgr.AddCore("mage")

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
		err := slotMgr.SetCore(0, "warrior")
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
		invMgr.AddCore("warrior")

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
		_ = slotMgr.SetCore(0, "warrior")
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
	})

	t.Run("複数スロットにエージェントを設定できる", func(t *testing.T) {
		// 新システムのインベントリを作成
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior")
		invMgr.AddCore("mage")

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
		_ = slotMgr.SetCore(0, "warrior")
		_ = slotMgr.SetSkill(0, 0, "slash", "")

		// スロット1にメイジを設定
		_ = slotMgr.SetCore(1, "mage")
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
		invMgr.AddCore("warrior")

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
		_ = slotMgr.SetCore(1, "warrior")
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
	t.Run("同一TypeIDのコアは重複追加されない", func(t *testing.T) {
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior")

		// 同じコアを再度追加（重複なので更新されない）
		updated := invMgr.AddCore("warrior")

		if updated {
			t.Error("重複コアでは更新されないべき")
		}

		// コアが保有されていることを確認
		if !invMgr.Cores().HasCore("warrior") {
			t.Error("コアが保有されているべき")
		}
	})

	t.Run("異なるTypeIDのコアは追加される", func(t *testing.T) {
		invMgr := inventory.NewInventoryManager()

		// コアを追加
		invMgr.AddCore("warrior")

		// 別のコアを追加
		updated := invMgr.AddCore("mage")

		if !updated {
			t.Error("異なるコアでは更新されるべき")
		}

		// 両方のコアが保有されていることを確認
		if !invMgr.Cores().HasCore("warrior") || !invMgr.Cores().HasCore("mage") {
			t.Error("両方のコアが保有されているべき")
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
