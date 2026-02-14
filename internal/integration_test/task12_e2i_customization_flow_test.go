// Package integration_test は統合テストを提供します。
// このファイルはタスク12.3: E2Iカスタマイズフローの統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"
)

// ==================== タスク 12.3: E2Iカスタマイズフローテスト ====================

// createE2ITestCoreTypes はE2Iテスト用のCoreTypeマスタデータを作成します。
func createE2ITestCoreTypes() map[string]domain.CoreType {
	return map[string]domain.CoreType{
		"balanced": {
			ID:             "balanced",
			Name:           "バランス",
			AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
			StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
			PassiveSkillID: "adaptability",
		},
		"berserker": {
			ID:             "berserker",
			Name:           "バーサーカー",
			AllowedTags:    []string{"physical_low", "physical_mid", "physical_high"},
			StatWeights:    map[string]float64{"STR": 2.0, "INT": 0.3, "WIL": 0.8, "LUK": 0.3},
			PassiveSkillID: "rage",
		},
		"archmage": {
			ID:             "archmage",
			Name:           "アークメイジ",
			AllowedTags:    []string{"magic_low", "magic_mid", "magic_high"},
			StatWeights:    map[string]float64{"STR": 0.3, "INT": 2.0, "WIL": 1.0, "LUK": 0.3},
			PassiveSkillID: "arcane_mastery",
		},
		"priest": {
			ID:             "priest",
			Name:           "プリースト",
			AllowedTags:    []string{"heal_low", "heal_mid", "buff_low"},
			StatWeights:    map[string]float64{"STR": 0.3, "INT": 1.2, "WIL": 1.5, "LUK": 0.5},
			PassiveSkillID: "divine_blessing",
		},
	}
}

// createE2ITestSkillTypes はE2Iテスト用のSkillTypeマスタデータを作成します。
func createE2ITestSkillTypes() map[string]domain.SkillType {
	return map[string]domain.SkillType{
		"basic_attack": {
			ID:          "basic_attack",
			Name:        "基本攻撃",
			Tags:        []string{"physical_low"},
			Description: "基本的な物理攻撃",
		},
		"power_attack": {
			ID:          "power_attack",
			Name:        "パワーアタック",
			Tags:        []string{"physical_mid"},
			Description: "強力な物理攻撃",
		},
		"berserk_strike": {
			ID:          "berserk_strike",
			Name:        "バーサークストライク",
			Tags:        []string{"physical_high"},
			Description: "最強の物理攻撃",
		},
		"fire_bolt": {
			ID:          "fire_bolt",
			Name:        "ファイアボルト",
			Tags:        []string{"magic_low"},
			Description: "基本的な火魔法",
		},
		"inferno": {
			ID:          "inferno",
			Name:        "インフェルノ",
			Tags:        []string{"magic_mid"},
			Description: "強力な火魔法",
		},
		"meteor": {
			ID:          "meteor",
			Name:        "メテオ",
			Tags:        []string{"magic_high"},
			Description: "最強の火魔法",
		},
		"heal": {
			ID:          "heal",
			Name:        "ヒール",
			Tags:        []string{"heal_low"},
			Description: "基本的な回復魔法",
		},
		"cure_all": {
			ID:          "cure_all",
			Name:        "キュアオール",
			Tags:        []string{"heal_mid"},
			Description: "強力な回復魔法",
		},
		"attack_up": {
			ID:          "attack_up",
			Name:        "攻撃UP",
			Tags:        []string{"buff_low"},
			Description: "攻撃力上昇バフ",
		},
		"poison": {
			ID:          "poison",
			Name:        "ポイズン",
			Tags:        []string{"debuff_low"},
			Description: "毒状態を付与",
		},
	}
}

// createE2ITestPassiveSkills はE2Iテスト用のPassiveSkillマスタデータを作成します。
func createE2ITestPassiveSkills() map[string]domain.PassiveSkill {
	return map[string]domain.PassiveSkill{
		"adaptability": {
			ID:          "adaptability",
			Name:        "適応力",
			Description: "全ステータス微増",
		},
		"rage": {
			ID:          "rage",
			Name:        "激怒",
			Description: "HP減少時攻撃力上昇",
		},
		"arcane_mastery": {
			ID:          "arcane_mastery",
			Name:        "秘術の極み",
			Description: "魔法威力上昇",
		},
		"divine_blessing": {
			ID:          "divine_blessing",
			Name:        "神の祝福",
			Description: "回復効果上昇",
		},
	}
}

// ==================== 一連フローテスト ====================

// TestE2I_FullCustomizationFlow はスロット選択→コア設定→レベル選択→スキル設定→チェイン効果選択の一連フローをテストします。
func TestE2I_FullCustomizationFlow(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// インベントリにアイテムを追加（バトル報酬シミュレーション）
	invManager.AddCore("berserker")
	invManager.AddSkill("basic_attack")   // チェイン効果付き
	invManager.AddSkill("basic_attack")   // 同じスキル、別のチェイン効果
	invManager.AddSkill("power_attack")   // チェイン効果なし
	invManager.AddSkill("berserk_strike") // 最強スキル

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// ステップ1: スロット選択（スロット0）
	slotIndex := 0

	// ステップ2: コア設定
	err := slotManager.SetCore(slotIndex, "berserker")
	if err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}

	// ステップ3: コア設定確認
	slot0 := slotManager.GetSlot(slotIndex)
	if slot0.CoreTypeID != "berserker" {
		t.Errorf("CoreTypeIDが不正: got %s, want berserker", slot0.CoreTypeID)
	}

	// ステップ4: 互換スキル一覧を取得
	compatibleSkills := slotManager.GetCompatibleSkills(slotIndex)
	// バーサーカーは physical_low, physical_mid, physical_high を許可
	// basic_attack(physical_low), power_attack(physical_mid), berserk_strike(physical_high)
	if len(compatibleSkills) != 3 {
		t.Errorf("互換スキル数が不正: got %d, want 3", len(compatibleSkills))
	}

	// ステップ5: スキル設定（チェイン効果付き）
	err = slotManager.SetSkill(slotIndex, 0, "basic_attack")
	if err != nil {
		t.Fatalf("スキル設定(チェイン効果付き)に失敗: %v", err)
	}

	// ステップ6: 同じスキルを別スロットに設定しようとするとエラーになることを確認
	err = slotManager.SetSkill(slotIndex, 1, "basic_attack")
	if err == nil {
		t.Error("同じスキルIDを複数スロットに設定できてしまった")
	}
	if err != slot.ErrSkillAlreadyEquipped {
		t.Errorf("期待するエラー: %v, 実際: %v", slot.ErrSkillAlreadyEquipped, err)
	}

	// ステップ7: スキル設定（チェイン効果なし）
	err = slotManager.SetSkill(slotIndex, 1, "power_attack")
	if err != nil {
		t.Fatalf("スキル設定(チェインなし)に失敗: %v", err)
	}

	// ステップ8: スキル設定（最強スキル）
	err = slotManager.SetSkill(slotIndex, 2, "berserk_strike")
	if err != nil {
		t.Fatalf("スキル設定(最強)に失敗: %v", err)
	}

	// ステップ9: 最終確認（3つのスキルが設定されているはず）
	if slot0.GetSkillCount() != 3 {
		t.Errorf("スキル数が不正: got %d, want 3", slot0.GetSkillCount())
	}

	// ステップ10: バトル用エージェント構築
	agents := slotManager.BuildAgentsForBattle()
	if len(agents) != 1 {
		t.Fatalf("エージェント数が不正: got %d, want 1", len(agents))
	}

	agent := agents[0]
	if agent.Core.TypeID != "berserker" {
		t.Errorf("エージェントコアTypeIDが不正: got %s, want berserker", agent.Core.TypeID)
	}
	if len(agent.Skills) != 3 {
		t.Errorf("エージェントスキル数が不正: got %d, want 3", len(agent.Skills))
	}
}

// TestE2I_MultipleAgentCustomizationFlow は3体のエージェントをカスタマイズするフローをテストします。
func TestE2I_MultipleAgentCustomizationFlow(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// 3種類のコアと各種スキルを追加
	invManager.AddCore("berserker")
	invManager.AddCore("archmage")
	invManager.AddCore("priest")

	invManager.AddSkill("basic_attack")
	invManager.AddSkill("power_attack")
	invManager.AddSkill("fire_bolt")
	invManager.AddSkill("inferno")
	invManager.AddSkill("heal")
	invManager.AddSkill("cure_all")
	invManager.AddSkill("attack_up")

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// スロット0: バーサーカー（物理アタッカー）
	slotManager.SetCore(0, "berserker")
	slotManager.SetSkill(0, 0, "basic_attack")
	slotManager.SetSkill(0, 1, "power_attack")

	// スロット1: アークメイジ（魔法アタッカー）
	slotManager.SetCore(1, "archmage")
	slotManager.SetSkill(1, 0, "fire_bolt")
	slotManager.SetSkill(1, 1, "inferno")

	// スロット2: プリースト（サポート）
	slotManager.SetCore(2, "priest")
	slotManager.SetSkill(2, 0, "heal")
	slotManager.SetSkill(2, 1, "cure_all")
	slotManager.SetSkill(2, 2, "attack_up")

	// 全スロット確認
	if slotManager.GetReadySlotCount() != 3 {
		t.Errorf("準備完了スロット数が不正: got %d, want 3", slotManager.GetReadySlotCount())
	}

	// バトル用エージェント構築
	agents := slotManager.BuildAgentsForBattle()
	if len(agents) != 3 {
		t.Fatalf("エージェント数が不正: got %d, want 3", len(agents))
	}

	// 各エージェントの構成確認
	expectedConfigs := []struct {
		coreTypeID string
		skillCount int
	}{
		{"berserker", 2},
		{"archmage", 2},
		{"priest", 3},
	}

	for i, expected := range expectedConfigs {
		agent := agents[i]
		if agent.Core.TypeID != expected.coreTypeID {
			t.Errorf("エージェント%d CoreTypeIDが不正: got %s, want %s", i, agent.Core.TypeID, expected.coreTypeID)
		}
		if len(agent.Skills) != expected.skillCount {
			t.Errorf("エージェント%d スキル数が不正: got %d, want %d", i, len(agent.Skills), expected.skillCount)
		}
	}
}

// ==================== 互換性フィルタリング動作確認 ====================

// TestE2I_CompatibilityFiltering_AllTagsMatch は全タグが一致するケースをテストします。
func TestE2I_CompatibilityFiltering_AllTagsMatch(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// バランスコア（ほぼすべてのタグを許可）
	invManager.AddCore("balanced")

	// 各種スキルを追加
	invManager.AddSkill("basic_attack")   // physical_low
	invManager.AddSkill("fire_bolt")      // magic_low
	invManager.AddSkill("heal")           // heal_low
	invManager.AddSkill("attack_up")      // buff_low
	invManager.AddSkill("poison")         // debuff_low
	invManager.AddSkill("power_attack")   // physical_mid - バランスは許可しない
	invManager.AddSkill("berserk_strike") // physical_high - バランスは許可しない

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// バランスコアを設定
	slotManager.SetCore(0, "balanced")

	// 互換スキル一覧を取得
	compatibleSkills := slotManager.GetCompatibleSkills(0)

	// バランスは physical_low, magic_low, heal_low, buff_low, debuff_low を許可
	// 7つのスキルのうち、5つが互換性あり
	if len(compatibleSkills) != 5 {
		t.Errorf("バランスコアの互換スキル数が不正: got %d, want 5", len(compatibleSkills))
	}
}

// TestE2I_CompatibilityFiltering_NoTagsMatch はタグが一致しないケースをテストします。
func TestE2I_CompatibilityFiltering_NoTagsMatch(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// プリーストコア（heal_low, heal_mid, buff_low のみ許可）
	invManager.AddCore("priest")

	// 物理スキルのみ追加（プリーストとは互換性なし）
	invManager.AddSkill("basic_attack")   // physical_low
	invManager.AddSkill("power_attack")   // physical_mid
	invManager.AddSkill("berserk_strike") // physical_high

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// プリーストコアを設定
	slotManager.SetCore(0, "priest")

	// 互換スキル一覧を取得
	compatibleSkills := slotManager.GetCompatibleSkills(0)

	// 物理スキルはプリーストと互換性がない
	if len(compatibleSkills) != 0 {
		t.Errorf("プリーストコアに物理スキルは互換性がないべき: got %d, want 0", len(compatibleSkills))
	}
}

// TestE2I_CompatibilityFiltering_CoreChangeRemovesIncompatible はコア変更時に互換性のないスキルが削除されることをテストします。
func TestE2I_CompatibilityFiltering_CoreChangeRemovesIncompatible(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// バランスとバーサーカーを追加
	invManager.AddCore("balanced")
	invManager.AddCore("berserker")

	// 各種スキルを追加
	invManager.AddSkill("basic_attack") // physical_low - 両方互換
	invManager.AddSkill("fire_bolt")    // magic_low - バランスのみ互換
	invManager.AddSkill("heal")         // heal_low - バランスのみ互換

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// バランスコアを設定
	slotManager.SetCore(0, "balanced")

	// スキルを設定
	slotManager.SetSkill(0, 0, "basic_attack")
	slotManager.SetSkill(0, 1, "fire_bolt")
	slotManager.SetSkill(0, 2, "heal")

	// スキル数を確認
	if slotManager.GetSlot(0).GetSkillCount() != 3 {
		t.Errorf("初期スキル数が不正: got %d, want 3", slotManager.GetSlot(0).GetSkillCount())
	}

	// バーサーカーに変更（fire_boltとhealは互換性がない）
	slotManager.SetCore(0, "berserker")

	// basic_attackのみ残る
	if slotManager.GetSlot(0).GetSkillCount() != 1 {
		t.Errorf("コア変更後のスキル数が不正: got %d, want 1", slotManager.GetSlot(0).GetSkillCount())
	}

	// 残っているのがbasic_attackであることを確認
	skill := slotManager.GetSlot(0).GetSkill(0)
	if skill.TypeID != "basic_attack" {
		t.Errorf("残っているスキルが不正: got %s, want basic_attack", skill.TypeID)
	}
}

// TestE2I_CompatibilityFiltering_ValidateMethod はValidateSkillCompatibilityメソッドの動作をテストします。
func TestE2I_CompatibilityFiltering_ValidateMethod(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// アークメイジコアを追加（magic系のみ許可）
	invManager.AddCore("archmage")

	// スキルを追加
	invManager.AddSkill("fire_bolt")    // magic_low - 互換
	invManager.AddSkill("basic_attack") // physical_low - 非互換

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypes,
		skillTypes,
		passiveSkills,
		nil, // chainEffects
	)

	// アークメイジコアを設定
	slotManager.SetCore(0, "archmage")

	// 互換性検証
	if !slotManager.ValidateSkillCompatibility(0, "fire_bolt") {
		t.Error("fire_boltはアークメイジと互換性があるべき")
	}
	if slotManager.ValidateSkillCompatibility(0, "basic_attack") {
		t.Error("basic_attackはアークメイジと互換性がないべき")
	}

	// 空スロットに対する検証
	if slotManager.ValidateSkillCompatibility(1, "fire_bolt") {
		t.Error("空スロットに対してはfalseを返すべき")
	}
}

// TestE2I_SkillOverwrite はスキル上書きの動作をテストします。
func TestE2I_SkillOverwrite(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// コアを追加
	invManager.AddCore("berserker")

	// スキルを追加
	invManager.AddSkill("basic_attack")
	invManager.AddSkill("power_attack")

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
	slotManager.SetCore(0, "berserker")

	// 最初にbasic_attackで設定
	err := slotManager.SetSkill(0, 0, "basic_attack")
	if err != nil {
		t.Errorf("basic_attackでの設定に失敗: %v", err)
	}

	// 同じスロット位置でスキルを上書きできることを確認
	err = slotManager.SetSkill(0, 0, "power_attack")
	if err != nil {
		t.Errorf("power_attackでの上書きに失敗: %v", err)
	}

	// スキルがpower_attackに変更されていることを確認
	skillConfig := slotManager.GetSlot(0).GetSkill(0)
	if skillConfig.TypeID != "power_attack" {
		t.Errorf("スキルが更新されていない: got %s, want power_attack", skillConfig.TypeID)
	}
}

// TestE2I_CoreSettingWorks はコア設定が動作することをテストします。
func TestE2I_CoreSettingWorks(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// コアを追加
	invManager.AddCore("berserker")

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
	err := slotManager.SetCore(0, "berserker")
	if err != nil {
		t.Errorf("コア設定に失敗: %v", err)
	}
	if slotManager.GetSlot(0).CoreTypeID != "berserker" {
		t.Errorf("CoreTypeIDが不正: got %s, want berserker", slotManager.GetSlot(0).CoreTypeID)
	}
}

// TestE2I_SlotClearAndReconfigure はスロットクリア後の再設定フローをテストします。
func TestE2I_SlotClearAndReconfigure(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("berserker")
	invManager.AddCore("archmage")
	invManager.AddSkill("basic_attack")
	invManager.AddSkill("fire_bolt")

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
	slotManager.SetCore(0, "berserker")
	slotManager.SetSkill(0, 0, "basic_attack")

	// 設定確認
	if slotManager.GetSlot(0).IsEmpty() {
		t.Error("スロットは空であるべきではない")
	}

	// スロットをクリア
	err := slotManager.ClearCore(0)
	if err != nil {
		t.Errorf("クリアに失敗: %v", err)
	}

	// クリア確認
	if !slotManager.GetSlot(0).IsEmpty() {
		t.Error("スロットは空であるべき")
	}

	// 再設定（別のコア）
	slotManager.SetCore(0, "archmage")
	slotManager.SetSkill(0, 0, "fire_bolt")

	// 再設定確認
	if slotManager.GetSlot(0).CoreTypeID != "archmage" {
		t.Errorf("CoreTypeIDが不正: got %s, want archmage", slotManager.GetSlot(0).CoreTypeID)
	}
}

// TestE2I_SameSkillMultipleTimes_ShouldFail は同じスキルを同一エージェント内で複数設定できないことをテストします。
func TestE2I_SameSkillMultipleTimes_ShouldFail(t *testing.T) {
	// テスト環境のセットアップ
	invManager := inventory.NewInventoryManager()
	coreTypes := createE2ITestCoreTypes()
	skillTypes := createE2ITestSkillTypes()
	passiveSkills := createE2ITestPassiveSkills()

	// コアとスキルを追加
	invManager.AddCore("berserker")
	invManager.AddSkill("basic_attack")

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
	slotManager.SetCore(0, "berserker")

	// 最初のスロットにスキルを設定
	err := slotManager.SetSkill(0, 0, "basic_attack")
	if err != nil {
		t.Fatalf("スロット0へのbasic_attack設定に失敗: %v", err)
	}

	// 同じスキルを別のスロットに設定しようとするとエラーになるべき
	err = slotManager.SetSkill(0, 1, "basic_attack")
	if err == nil {
		t.Error("同じスキルを別スロットに設定できてしまった")
	}
	if err != slot.ErrSkillAlreadyEquipped {
		t.Errorf("期待するエラー: %v, 実際: %v", slot.ErrSkillAlreadyEquipped, err)
	}

	// 1つだけが設定されていることを確認
	if slotManager.GetSlot(0).GetSkillCount() != 1 {
		t.Errorf("スキル数が不正: got %d, want 1", slotManager.GetSlot(0).GetSkillCount())
	}

	// バトル用エージェント構築
	agents := slotManager.BuildAgentsForBattle()
	if len(agents[0].Skills) != 1 {
		t.Errorf("エージェントスキル数が不正: got %d, want 1", len(agents[0].Skills))
	}
}
