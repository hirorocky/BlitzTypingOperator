// Package integration_test は統合テストを提供します。
// このファイルはタスク12.2: スロット→バトル連携の統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"
)

// ==================== タスク 12.2: スロット→バトル連携テスト ====================

// createBattleTestCoreTypes はバトルテスト用のCoreTypeマスタデータを作成します。
func createBattleTestCoreTypes() map[string]domain.CoreType {
	return map[string]domain.CoreType{
		"warrior": {
			ID:             "warrior",
			Name:           "ウォーリアー",
			AllowedTags:    []string{"physical_low", "physical_mid"},
			StatWeights:    map[string]float64{"STR": 1.5, "INT": 0.5, "WIL": 1.0, "LUK": 0.5},
			PassiveSkillID: "berserker_rage",
		},
		"mage": {
			ID:             "mage",
			Name:           "メイジ",
			AllowedTags:    []string{"magic_low", "magic_mid"},
			StatWeights:    map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 1.0, "LUK": 0.5},
			PassiveSkillID: "spell_mastery",
		},
		"support": {
			ID:             "support",
			Name:           "サポート",
			AllowedTags:    []string{"heal_low", "buff_low"},
			StatWeights:    map[string]float64{"STR": 0.5, "INT": 1.0, "WIL": 1.5, "LUK": 0.5},
			PassiveSkillID: "healing_light",
		},
	}
}

// createBattleTestSkillTypes はバトルテスト用のSkillTypeマスタデータを作成します。
func createBattleTestSkillTypes() map[string]domain.SkillType {
	return map[string]domain.SkillType{
		"slash": {
			ID:          "slash",
			Name:        "スラッシュ",
			Tags:        []string{"physical_low"},
			Description: "物理攻撃",
			Effects: []domain.SkillEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 10, StatCoef: 1.0, StatRef: "STR"},
					Probability: 1.0,
				},
			},
		},
		"heavy_strike": {
			ID:          "heavy_strike",
			Name:        "ヘビーストライク",
			Tags:        []string{"physical_mid"},
			Description: "強力な物理攻撃",
			Effects: []domain.SkillEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 20, StatCoef: 1.5, StatRef: "STR"},
					Probability: 1.0,
				},
			},
		},
		"fire": {
			ID:          "fire",
			Name:        "ファイア",
			Tags:        []string{"magic_low"},
			Description: "魔法攻撃",
			Effects: []domain.SkillEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 15, StatCoef: 1.0, StatRef: "INT"},
					Probability: 1.0,
				},
			},
		},
		"cure": {
			ID:          "cure",
			Name:        "キュア",
			Tags:        []string{"heal_low"},
			Description: "回復魔法",
			Effects: []domain.SkillEffect{
				{
					Target:      domain.TargetSelf,
					HPFormula:   &domain.HPFormula{Base: 20, StatCoef: 0.8, StatRef: "INT"},
					Probability: 1.0,
				},
			},
		},
	}
}

// createBattleTestPassiveSkills はバトルテスト用のPassiveSkillマスタデータを作成します。
func createBattleTestPassiveSkills() map[string]domain.PassiveSkill {
	return map[string]domain.PassiveSkill{
		"berserker_rage": {
			ID:          "berserker_rage",
			Name:        "バーサーカーの怒り",
			Description: "攻撃力上昇",
		},
		"spell_mastery": {
			ID:          "spell_mastery",
			Name:        "スペルマスタリー",
			Description: "魔法威力上昇",
		},
		"healing_light": {
			ID:          "healing_light",
			Name:        "癒やしの光",
			Description: "回復効果上昇",
		},
	}
}

// ==================== BuildAgentsForBattleの生成AgentModel検証 ====================

// TestSlotBattle_BuildAgentsForBattle_SingleAgent は単一エージェントのBuildAgentsForBattle動作をテストします。
func TestSlotBattle_BuildAgentsForBattle_SingleAgent(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("warrior")
	invManager.AddSkill("slash", "")
	invManager.AddSkill("heavy_strike", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// スロット0にコアとスキルを設定
	slotManager.SetCore(0, "warrior")
	slotManager.SetSkill(0, 0, "slash", "")
	slotManager.SetSkill(0, 1, "heavy_strike", "")

	// バトル用エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 生成されたエージェント数を確認
	if len(agents) != 1 {
		t.Fatalf("エージェント数が不正: got %d, want 1", len(agents))
	}

	// エージェントの内容を確認
	agent := agents[0]
	if agent == nil {
		t.Fatal("エージェントがnilです")
	}

	// コアTypeIDが設定されている
	if agent.Core == nil {
		t.Fatal("コアがnilです")
	}
	if agent.Core.TypeID != "warrior" {
		t.Errorf("コアTypeIDが不正: got %s, want warrior", agent.Core.TypeID)
	}

	// スキル数が一致
	if len(agent.Modules) != 2 {
		t.Errorf("スキル数が不正: got %d, want 2", len(agent.Modules))
	}
}

// TestSlotBattle_BuildAgentsForBattle_MultipleAgents は複数エージェントのBuildAgentsForBattle動作をテストします。
func TestSlotBattle_BuildAgentsForBattle_MultipleAgents(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// 複数のコアとスキルを追加
	invManager.AddCore("warrior")
	invManager.AddCore("mage")
	invManager.AddCore("support")
	invManager.AddSkill("slash", "")
	invManager.AddSkill("fire", "")
	invManager.AddSkill("cure", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 全スロットにコアとスキルを設定
	slotManager.SetCore(0, "warrior")
	slotManager.SetSkill(0, 0, "slash", "")

	slotManager.SetCore(1, "mage")
	slotManager.SetSkill(1, 0, "fire", "")

	slotManager.SetCore(2, "support")
	slotManager.SetSkill(2, 0, "cure", "")

	// バトル用エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 生成されたエージェント数を確認
	if len(agents) != 3 {
		t.Fatalf("エージェント数が不正: got %d, want 3", len(agents))
	}

	// 各エージェントのコアTypeIDを確認
	expectedCoreTypeIDs := []string{"warrior", "mage", "support"}
	for i, agent := range agents {
		if agent.Core.TypeID != expectedCoreTypeIDs[i] {
			t.Errorf("エージェント%dのコアTypeIDが不正: got %s, want %s", i, agent.Core.TypeID, expectedCoreTypeIDs[i])
		}
	}
}

// TestSlotBattle_BuildAgentsForBattle_EmptySlotsExcluded は空スロットがバトルから除外されることをテストします。
func TestSlotBattle_BuildAgentsForBattle_EmptySlotsExcluded(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアを追加
	invManager.AddCore("warrior")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// スロット1のみにコアを設定（スロット0,2は空）
	slotManager.SetCore(1, "warrior")

	// バトル用エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 空スロットは除外され、1体のみ
	if len(agents) != 1 {
		t.Errorf("エージェント数が不正: got %d, want 1", len(agents))
	}
}

// TestSlotBattle_BuildAgentsForBattle_NoAgents は全スロット空の場合のテストです。
func TestSlotBattle_BuildAgentsForBattle_NoAgents(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// AgentSlotManagerを作成（何も設定しない）
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// バトル用エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 空の場合は空スライス
	if len(agents) != 0 {
		t.Errorf("エージェント数が不正: got %d, want 0", len(agents))
	}
}

// TestSlotBattle_BuildAgentsForBattle_CoreOnlyNoSkill はスキルなしコアのみの場合のテストです。
func TestSlotBattle_BuildAgentsForBattle_CoreOnlyNoSkill(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアのみ追加
	invManager.AddCore("warrior")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// コアのみ設定（スキルなし）
	slotManager.SetCore(0, "warrior")

	// バトル用エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// コアがあればエージェントは生成される
	if len(agents) != 1 {
		t.Fatalf("エージェント数が不正: got %d, want 1", len(agents))
	}

	// スキルは空
	if len(agents[0].Modules) != 0 {
		t.Errorf("スキル数が不正: got %d, want 0", len(agents[0].Modules))
	}
}

// ==================== スロット構成変更後のバトル反映テスト ====================

// TestSlotBattle_SlotUpdateReflectsOnBattle はスロット変更がバトル用エージェントに反映されることをテストします。
func TestSlotBattle_SlotUpdateReflectsOnBattle(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("warrior")
	invManager.AddCore("mage")
	invManager.AddSkill("slash", "")
	invManager.AddSkill("fire", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定：ウォーリアー
	slotManager.SetCore(0, "warrior")
	slotManager.SetSkill(0, 0, "slash", "")

	// バトル用エージェントを構築（初回）
	agents1 := slotManager.BuildAgentsForBattle()
	if agents1[0].Core.TypeID != "warrior" {
		t.Errorf("初期コアTypeIDが不正: got %s, want warrior", agents1[0].Core.TypeID)
	}

	// スロット設定を変更：メイジに
	slotManager.SetCore(0, "mage")
	slotManager.SetSkill(0, 0, "fire", "")

	// バトル用エージェントを再構築
	agents2 := slotManager.BuildAgentsForBattle()

	// 変更が反映されている
	if agents2[0].Core.TypeID != "mage" {
		t.Errorf("変更後コアTypeIDが不正: got %s, want mage", agents2[0].Core.TypeID)
	}
}

// TestSlotBattle_CoreReplacementReflectsOnBattle はコア再設定がバトルに反映されることをテストします。
func TestSlotBattle_CoreReplacementReflectsOnBattle(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアを追加
	invManager.AddCore("warrior")
	invManager.AddCore("mage")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定
	slotManager.SetCore(0, "warrior")

	// バトル用エージェントを構築
	agents1 := slotManager.BuildAgentsForBattle()
	if agents1[0].Core.TypeID != "warrior" {
		t.Errorf("初期コアTypeIDが不正: got %s, want warrior", agents1[0].Core.TypeID)
	}

	// コアを変更
	slotManager.SetCore(0, "mage")

	// バトル用エージェントを再構築
	agents2 := slotManager.BuildAgentsForBattle()

	// コアが変更されている
	if agents2[0].Core.TypeID != "mage" {
		t.Errorf("変更後コアTypeIDが不正: got %s, want mage", agents2[0].Core.TypeID)
	}
}

// TestSlotBattle_SkillChangeReflectsOnBattle はスキル変更がバトルに反映されることをテストします。
func TestSlotBattle_SkillChangeReflectsOnBattle(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("warrior")
	invManager.AddSkill("slash", "")
	invManager.AddSkill("heavy_strike", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定：slashのみ
	slotManager.SetCore(0, "warrior")
	slotManager.SetSkill(0, 0, "slash", "")

	// バトル用エージェントを構築
	agents1 := slotManager.BuildAgentsForBattle()
	if len(agents1[0].Modules) != 1 {
		t.Fatalf("初期スキル数が不正: got %d, want 1", len(agents1[0].Modules))
	}

	// スキルを追加
	slotManager.SetSkill(0, 1, "heavy_strike", "")

	// バトル用エージェントを再構築
	agents2 := slotManager.BuildAgentsForBattle()

	// スキルが追加されている
	if len(agents2[0].Modules) != 2 {
		t.Errorf("変更後スキル数が不正: got %d, want 2", len(agents2[0].Modules))
	}
}

// TestSlotBattle_ClearSlotReflectsOnBattle はスロットクリアがバトルに反映されることをテストします。
func TestSlotBattle_ClearSlotReflectsOnBattle(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアを追加
	invManager.AddCore("warrior")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定
	slotManager.SetCore(0, "warrior")
	slotManager.SetCore(1, "warrior")

	// バトル用エージェントを構築
	agents1 := slotManager.BuildAgentsForBattle()
	if len(agents1) != 2 {
		t.Fatalf("初期エージェント数が不正: got %d, want 2", len(agents1))
	}

	// スロット1をクリア
	slotManager.ClearCore(1)

	// バトル用エージェントを再構築
	agents2 := slotManager.BuildAgentsForBattle()

	// エージェント数が減っている
	if len(agents2) != 1 {
		t.Errorf("変更後エージェント数が不正: got %d, want 1", len(agents2))
	}
}

// ==================== バトル中のスロット変更禁止テスト ====================

// TestSlotBattle_LockDuringBattle はバトル中のスロット変更が禁止されることをテストします。
func TestSlotBattle_LockDuringBattle(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("warrior")
	invManager.AddCore("mage")
	invManager.AddSkill("slash", "")
	invManager.AddSkill("fire", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定
	slotManager.SetCore(0, "warrior")
	slotManager.SetSkill(0, 0, "slash", "")

	// バトル開始：ロック
	slotManager.Lock()

	// ロック中の状態確認
	if !slotManager.IsLocked() {
		t.Error("IsLockedがtrueを返すべきです")
	}

	// ロック中のコア設定はエラー
	err := slotManager.SetCore(1, "mage")
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のSetCoreはErrSlotLockedを返すべき: got %v", err)
	}

	// ロック中のスキル設定はエラー
	err = slotManager.SetSkill(0, 1, "fire", "")
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のSetSkillはErrSlotLockedを返すべき: got %v", err)
	}

	// ロック中のコアクリアはエラー
	err = slotManager.ClearCore(0)
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のClearCoreはErrSlotLockedを返すべき: got %v", err)
	}

	// ロック中のスキルクリアはエラー
	err = slotManager.ClearSkill(0, 0)
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のClearSkillはErrSlotLockedを返すべき: got %v", err)
	}

	// バトル終了：アンロック
	slotManager.Unlock()

	// アンロック後の状態確認
	if slotManager.IsLocked() {
		t.Error("IsLockedがfalseを返すべきです")
	}

	// アンロック後はコア設定可能
	err = slotManager.SetCore(1, "mage")
	if err != nil {
		t.Errorf("アンロック後のSetCoreは成功すべき: %v", err)
	}
}

// TestSlotBattle_BuildAgentsForBattleWhileLocked はロック中でもBuildAgentsForBattleは動作することをテストします。
func TestSlotBattle_BuildAgentsForBattleWhileLocked(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createBattleTestCoreTypes()
	skillTypes := createBattleTestSkillTypes()
	passiveSkills := createBattleTestPassiveSkills()

	// コアを追加
	invManager.AddCore("warrior")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 初期設定
	slotManager.SetCore(0, "warrior")

	// バトル開始：ロック
	slotManager.Lock()

	// ロック中でもBuildAgentsForBattleは動作する（読み取り専用）
	agents := slotManager.BuildAgentsForBattle()
	if len(agents) != 1 {
		t.Errorf("ロック中でもBuildAgentsForBattleは動作すべき: got %d agents", len(agents))
	}
}
