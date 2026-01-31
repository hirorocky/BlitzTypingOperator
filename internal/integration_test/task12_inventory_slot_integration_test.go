// Package integration_test は統合テストを提供します。
// このファイルはタスク12.1: インベントリ→スロット連携の統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"
)

// ==================== タスク 12.1: インベントリ→スロット連携テスト ====================

// createTestCoreTypes はテスト用のCoreTypeマスタデータを作成します。
func createTestCoreTypes() map[string]domain.CoreType {
	return map[string]domain.CoreType{
		"all_rounder": {
			ID:          "all_rounder",
			Name:        "オールラウンダー",
			AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
			StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
		"attacker": {
			ID:          "attacker",
			Name:        "アタッカー",
			AllowedTags: []string{"physical_low", "physical_mid"},
			StatWeights: map[string]float64{"STR": 1.5, "INT": 0.5, "WIL": 0.5, "LUK": 0.5},
		},
		"healer": {
			ID:          "healer",
			Name:        "ヒーラー",
			AllowedTags: []string{"heal_low", "heal_mid"},
			StatWeights: map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 1.0, "LUK": 0.5},
		},
	}
}

// createTestSkillTypes はテスト用のSkillTypeマスタデータを作成します。
func createTestSkillTypes() map[string]domain.SkillType {
	return map[string]domain.SkillType{
		"physical_strike": {
			ID:          "physical_strike",
			Name:        "物理打撃",
			Tags:        []string{"physical_low"},
			Description: "物理ダメージを与える",
		},
		"fireball": {
			ID:          "fireball",
			Name:        "ファイアボール",
			Tags:        []string{"magic_low"},
			Description: "魔法ダメージを与える",
		},
		"heal": {
			ID:          "heal",
			Name:        "ヒール",
			Tags:        []string{"heal_low"},
			Description: "HPを回復する",
		},
		"power_slash": {
			ID:          "power_slash",
			Name:        "パワースラッシュ",
			Tags:        []string{"physical_mid"},
			Description: "強力な物理ダメージを与える",
		},
	}
}

// createTestPassiveSkills はテスト用のPassiveSkillマスタデータを作成します。
func createTestPassiveSkills() map[string]domain.PassiveSkill {
	return map[string]domain.PassiveSkill{
		"adaptability": {
			ID:          "adaptability",
			Name:        "適応力",
			Description: "全ステータス微増",
		},
	}
}

// ==================== コア追加後のスロット設定フローテスト ====================

// TestInventorySlot_AddCoreAndSetSlot はコア追加後にスロットに設定できることをテストします。
func TestInventorySlot_AddCoreAndSetSlot(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアを追加
	updated := invManager.AddCore("all_rounder", 5)
	if !updated {
		t.Fatal("コアの追加に失敗しました")
	}

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// スロットにコアを設定
	err := slotManager.SetCore(0, "all_rounder", 5)
	if err != nil {
		t.Errorf("コア設定に失敗: %v", err)
	}

	// スロットの状態を確認
	slot0 := slotManager.GetSlot(0)
	if slot0.CoreTypeID != "all_rounder" {
		t.Errorf("CoreTypeIDが不正: got %s, want all_rounder", slot0.CoreTypeID)
	}
	if slot0.CoreLevel != 5 {
		t.Errorf("CoreLevelが不正: got %d, want 5", slot0.CoreLevel)
	}
}

// TestInventorySlot_AddCoreAtLowerLevel はコアを最大レベル以下で設定できることをテストします。
func TestInventorySlot_AddCoreAtLowerLevel(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアをレベル10で追加
	invManager.AddCore("all_rounder", 10)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// レベル5で設定（最大レベル以下）
	err := slotManager.SetCore(0, "all_rounder", 5)
	if err != nil {
		t.Errorf("レベル5でのコア設定に失敗: %v", err)
	}

	// レベル1で設定
	err = slotManager.SetCore(1, "all_rounder", 1)
	if err != nil {
		t.Errorf("レベル1でのコア設定に失敗: %v", err)
	}

	// レベル10で設定
	err = slotManager.SetCore(2, "all_rounder", 10)
	if err != nil {
		t.Errorf("レベル10でのコア設定に失敗: %v", err)
	}

	// スロットの状態を確認
	if slotManager.GetSlot(0).CoreLevel != 5 {
		t.Errorf("スロット0のレベルが不正: got %d, want 5", slotManager.GetSlot(0).CoreLevel)
	}
	if slotManager.GetSlot(1).CoreLevel != 1 {
		t.Errorf("スロット1のレベルが不正: got %d, want 1", slotManager.GetSlot(1).CoreLevel)
	}
	if slotManager.GetSlot(2).CoreLevel != 10 {
		t.Errorf("スロット2のレベルが不正: got %d, want 10", slotManager.GetSlot(2).CoreLevel)
	}
}

// TestInventorySlot_CoreNotOwnedError は保有していないコアを設定できないことをテストします。
func TestInventorySlot_CoreNotOwnedError(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアを追加しない状態でAgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 保有していないコアを設定しようとする
	err := slotManager.SetCore(0, "all_rounder", 1)
	if err != slot.ErrCoreNotOwned {
		t.Errorf("ErrCoreNotOwnedが返されるべき: got %v", err)
	}
}

// TestInventorySlot_LevelOutOfRangeError は最大レベルを超えるレベルでエラーになることをテストします。
func TestInventorySlot_LevelOutOfRangeError(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアをレベル5で追加
	invManager.AddCore("all_rounder", 5)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 最大レベルを超えるレベルで設定しようとする
	err := slotManager.SetCore(0, "all_rounder", 10)
	if err != slot.ErrLevelOutOfRange {
		t.Errorf("ErrLevelOutOfRangeが返されるべき: got %v", err)
	}
}

// TestInventorySlot_CoreUpdateReflectsOnNewSlotManager はインベントリの更新が新しいSlotManagerに反映されることをテストします。
func TestInventorySlot_CoreUpdateReflectsOnNewSlotManager(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// 初期コア追加
	invManager.AddCore("all_rounder", 3)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// レベル3で設定できる
	err := slotManager.SetCore(0, "all_rounder", 3)
	if err != nil {
		t.Errorf("レベル3でのコア設定に失敗: %v", err)
	}

	// インベントリにより高いレベルのコアを追加
	invManager.AddCore("all_rounder", 7)

	// 同じSlotManagerからより高いレベルで設定できるようになる
	// (インベントリへの参照を保持しているため)
	err = slotManager.SetCore(1, "all_rounder", 7)
	if err != nil {
		t.Errorf("レベル7でのコア設定に失敗: %v", err)
	}
}

// ==================== スキル追加後のスロット設定フローテスト ====================

// TestInventorySlot_AddSkillAndSetSlot はスキル追加後にスロットに設定できることをテストします。
func TestInventorySlot_AddSkillAndSetSlot(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddSkill("physical_strike", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// コアを設定
	err := slotManager.SetCore(0, "all_rounder", 5)
	if err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}

	// スキルを設定
	err = slotManager.SetSkill(0, 0, "physical_strike", "")
	if err != nil {
		t.Errorf("スキル設定に失敗: %v", err)
	}

	// スロットの状態を確認
	skillConfig := slotManager.GetSlot(0).GetSkill(0)
	if skillConfig.TypeID != "physical_strike" {
		t.Errorf("スキルTypeIDが不正: got %s, want physical_strike", skillConfig.TypeID)
	}
}

// TestInventorySlot_AddSkillWithChainEffect はチェイン効果付きスキルを設定できることをテストします。
func TestInventorySlot_AddSkillWithChainEffect(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアとスキル（チェイン効果付き）を追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddSkill("fireball", "chain_fire_damage")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// コアを設定
	slotManager.SetCore(0, "all_rounder", 5)

	// チェイン効果付きでスキルを設定
	err := slotManager.SetSkill(0, 0, "fireball", "chain_fire_damage")
	if err != nil {
		t.Errorf("チェイン効果付きスキル設定に失敗: %v", err)
	}

	// スロットの状態を確認
	skillConfig := slotManager.GetSlot(0).GetSkill(0)
	if skillConfig.TypeID != "fireball" {
		t.Errorf("スキルTypeIDが不正: got %s, want fireball", skillConfig.TypeID)
	}
	if skillConfig.ChainEffectID != "chain_fire_damage" {
		t.Errorf("ChainEffectIDが不正: got %s, want chain_fire_damage", skillConfig.ChainEffectID)
	}
}

// TestInventorySlot_SkillNotOwnedError は保有していないスキルを設定できないことをテストします。
func TestInventorySlot_SkillNotOwnedError(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアのみ追加（スキルは追加しない）
	invManager.AddCore("all_rounder", 5)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// コアを設定
	slotManager.SetCore(0, "all_rounder", 5)

	// 保有していないスキルを設定しようとする
	err := slotManager.SetSkill(0, 0, "physical_strike", "")
	if err != slot.ErrSkillNotOwned {
		t.Errorf("ErrSkillNotOwnedが返されるべき: got %v", err)
	}
}

// TestInventorySlot_ChainVariationNotOwnedError は保有していないチェイン効果でエラーになることをテストします。
func TestInventorySlot_ChainVariationNotOwnedError(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアとスキル（チェイン効果なし）を追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddSkill("fireball", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// コアを設定
	slotManager.SetCore(0, "all_rounder", 5)

	// 保有していないチェイン効果でスキルを設定しようとする
	err := slotManager.SetSkill(0, 0, "fireball", "chain_not_owned")
	if err != slot.ErrChainVariationNotOwned {
		t.Errorf("ErrChainVariationNotOwnedが返されるべき: got %v", err)
	}
}

// TestInventorySlot_SkillIncompatibleError はコアと互換性のないスキルでエラーになることをテストします。
func TestInventorySlot_SkillIncompatibleError(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// ヒーラーコアと物理スキルを追加
	invManager.AddCore("healer", 5)
	invManager.AddSkill("power_slash", "") // physical_mid タグ

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// ヒーラーコアを設定
	slotManager.SetCore(0, "healer", 5)

	// 互換性のないスキルを設定しようとする
	err := slotManager.SetSkill(0, 0, "power_slash", "")
	if err != slot.ErrSkillIncompatible {
		t.Errorf("ErrSkillIncompatibleが返されるべき: got %v", err)
	}
}

// TestInventorySlot_CompatibleSkillsFilter はコア変更時に互換スキルがフィルタリングされることをテストします。
func TestInventorySlot_CompatibleSkillsFilter(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// 複数のコアとスキルを追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddCore("attacker", 5)
	invManager.AddSkill("physical_strike", "")
	invManager.AddSkill("fireball", "")
	invManager.AddSkill("heal", "")
	invManager.AddSkill("power_slash", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// オールラウンダーコアを設定
	slotManager.SetCore(0, "all_rounder", 5)

	// 互換スキル一覧を取得
	compatibleSkills := slotManager.GetCompatibleSkills(0)
	if len(compatibleSkills) != 3 { // physical_strike, fireball, heal
		t.Errorf("オールラウンダーの互換スキル数が不正: got %d, want 3", len(compatibleSkills))
	}

	// アタッカーコアを設定
	slotManager.SetCore(1, "attacker", 5)

	// 互換スキル一覧を取得
	compatibleSkills = slotManager.GetCompatibleSkills(1)
	if len(compatibleSkills) != 2 { // physical_strike, power_slash
		t.Errorf("アタッカーの互換スキル数が不正: got %d, want 2", len(compatibleSkills))
	}
}

// TestInventorySlot_SkillAutoRemoveOnCoreChange はコア変更時に互換性のないスキルが自動削除されることをテストします。
func TestInventorySlot_SkillAutoRemoveOnCoreChange(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// オールラウンダーとヒーラーを追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddCore("healer", 5)
	invManager.AddSkill("physical_strike", "")
	invManager.AddSkill("heal", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// オールラウンダーコアを設定
	slotManager.SetCore(0, "all_rounder", 5)

	// スキルを設定
	slotManager.SetSkill(0, 0, "physical_strike", "")
	slotManager.SetSkill(0, 1, "heal", "")

	// スキル設定を確認
	if slotManager.GetSlot(0).GetSkillCount() != 2 {
		t.Errorf("スキル数が不正: got %d, want 2", slotManager.GetSlot(0).GetSkillCount())
	}

	// ヒーラーコアに変更（physical_strikeは互換性がない）
	slotManager.SetCore(0, "healer", 5)

	// physical_strikeが自動削除され、healのみ残る
	skillCount := slotManager.GetSlot(0).GetSkillCount()
	if skillCount != 1 {
		t.Errorf("コア変更後のスキル数が不正: got %d, want 1", skillCount)
	}

	// 残っているのがhealであることを確認
	foundHeal := false
	for i := range domain.MaxSkillSlotCount {
		skill := slotManager.GetSlot(0).GetSkill(i)
		if !skill.IsEmpty() && skill.TypeID == "heal" {
			foundHeal = true
		}
	}
	if !foundHeal {
		t.Error("healスキルが残っているべきです")
	}
}

// TestInventorySlot_MultipleSlotsSameCore は同じコアを複数スロットに設定できることをテストします。
func TestInventorySlot_MultipleSlotsSameCore(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアを追加
	invManager.AddCore("all_rounder", 10)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 全スロットに同じコアを異なるレベルで設定
	err := slotManager.SetCore(0, "all_rounder", 10)
	if err != nil {
		t.Errorf("スロット0への設定に失敗: %v", err)
	}
	err = slotManager.SetCore(1, "all_rounder", 5)
	if err != nil {
		t.Errorf("スロット1への設定に失敗: %v", err)
	}
	err = slotManager.SetCore(2, "all_rounder", 1)
	if err != nil {
		t.Errorf("スロット2への設定に失敗: %v", err)
	}

	// 全スロットが設定されていることを確認
	if slotManager.GetReadySlotCount() != 3 {
		t.Errorf("バトル準備可能スロット数が不正: got %d, want 3", slotManager.GetReadySlotCount())
	}
}

// TestInventorySlot_MultipleSlotsSameSkill_ShouldFail は同じスキルを複数スロットに設定できないことをテストします。
func TestInventorySlot_MultipleSlotsSameSkill_ShouldFail(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createTestCoreTypes()
	skillTypes := createTestSkillTypes()
	passiveSkills := createTestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("all_rounder", 5)
	invManager.AddSkill("physical_strike", "")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// 全スロットにコアを設定
	for i := range 3 {
		slotManager.SetCore(i, "all_rounder", 5)
	}

	// 最初のスロットにスキルを設定
	err := slotManager.SetSkill(0, 0, "physical_strike", "")
	if err != nil {
		t.Fatalf("スロット0へのスキル設定に失敗: %v", err)
	}

	// 同じスキルを他のスロットに設定しようとするとエラーになるはず
	err = slotManager.SetSkill(1, 0, "physical_strike", "")
	if err == nil {
		t.Error("同じスキルを別スロットに設定できてしまった")
	}
	if err != slot.ErrSkillAlreadyEquipped {
		t.Errorf("期待するエラー: %v, 実際: %v", slot.ErrSkillAlreadyEquipped, err)
	}

	// スロット0のみにスキルが設定されていることを確認
	if slotManager.GetSlot(0).GetSkillCount() != 1 {
		t.Errorf("スロット0のスキル数が不正: got %d, want 1", slotManager.GetSlot(0).GetSkillCount())
	}
	if slotManager.GetSlot(1).GetSkillCount() != 0 {
		t.Errorf("スロット1のスキル数が不正: got %d, want 0", slotManager.GetSlot(1).GetSkillCount())
	}
}
