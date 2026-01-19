// Package savedata はセーブデータの永続化を担当します。
// このファイルはユニークインベントリ関連のセーブデータのテストを提供します。

package savedata

import (
	"encoding/json"
	"testing"
)

// ==================== タスク7.1: 新セーブデータ構造体のテスト ====================

// TestCoreInventorySave_JSONSerialization はCoreInventorySaveのJSON化をテストします。
func TestCoreInventorySave_JSONSerialization(t *testing.T) {
	save := CoreInventorySave{
		Cores: map[string]int{
			"all_rounder":    5,
			"attack_balance": 3,
		},
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded CoreInventorySave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 検証
	if len(loaded.Cores) != 2 {
		t.Errorf("Cores count: got %d, want 2", len(loaded.Cores))
	}
	if loaded.Cores["all_rounder"] != 5 {
		t.Errorf("all_rounder level: got %d, want 5", loaded.Cores["all_rounder"])
	}
	if loaded.Cores["attack_balance"] != 3 {
		t.Errorf("attack_balance level: got %d, want 3", loaded.Cores["attack_balance"])
	}
}

// TestSkillInventorySave_JSONSerialization はSkillInventorySaveのJSON化をテストします。
func TestSkillInventorySave_JSONSerialization(t *testing.T) {
	save := SkillInventorySave{
		Skills: map[string][]string{
			"physical_lv1": {"damage_bonus", "life_steal"},
			"heal_lv1":     {},
		},
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded SkillInventorySave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 検証
	if len(loaded.Skills) != 2 {
		t.Errorf("Skills count: got %d, want 2", len(loaded.Skills))
	}
	if len(loaded.Skills["physical_lv1"]) != 2 {
		t.Errorf("physical_lv1 chain effects: got %d, want 2", len(loaded.Skills["physical_lv1"]))
	}
	if len(loaded.Skills["heal_lv1"]) != 0 {
		t.Errorf("heal_lv1 chain effects: got %d, want 0", len(loaded.Skills["heal_lv1"]))
	}
}

// TestSkillSlotSaveCfg_JSONSerialization はSkillSlotSaveCfgのJSON化をテストします。
func TestSkillSlotSaveCfg_JSONSerialization(t *testing.T) {
	// チェイン効果ありのスキル
	withChain := SkillSlotSaveCfg{
		TypeID:        "physical_lv1",
		ChainEffectID: "damage_bonus",
	}

	data, err := json.Marshal(withChain)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	var loaded SkillSlotSaveCfg
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	if loaded.TypeID != "physical_lv1" {
		t.Errorf("TypeID: got %s, want physical_lv1", loaded.TypeID)
	}
	if loaded.ChainEffectID != "damage_bonus" {
		t.Errorf("ChainEffectID: got %s, want damage_bonus", loaded.ChainEffectID)
	}

	// 空のスロット（omitemptyで省略される）
	empty := SkillSlotSaveCfg{}
	emptyData, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("空スロットのJSON化に失敗: %v", err)
	}

	// 空のJSONは "{}" になるはず
	if string(emptyData) != "{}" {
		t.Errorf("空スロットのJSON: got %s, want {}", string(emptyData))
	}
}

// TestAgentSlotSave_JSONSerialization はAgentSlotSaveのJSON化をテストします。
func TestAgentSlotSave_JSONSerialization(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "all_rounder",
		CoreLevel:  5,
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1", ChainEffectID: "damage_bonus"},
			{TypeID: "heal_lv1", ChainEffectID: ""},
			{TypeID: "", ChainEffectID: ""},
			{TypeID: "", ChainEffectID: ""},
		},
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded AgentSlotSave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 検証
	if loaded.CoreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", loaded.CoreTypeID)
	}
	if loaded.CoreLevel != 5 {
		t.Errorf("CoreLevel: got %d, want 5", loaded.CoreLevel)
	}
	if loaded.Skills[0].TypeID != "physical_lv1" {
		t.Errorf("Skills[0].TypeID: got %s, want physical_lv1", loaded.Skills[0].TypeID)
	}
	if loaded.Skills[0].ChainEffectID != "damage_bonus" {
		t.Errorf("Skills[0].ChainEffectID: got %s, want damage_bonus", loaded.Skills[0].ChainEffectID)
	}
	if loaded.Skills[1].TypeID != "heal_lv1" {
		t.Errorf("Skills[1].TypeID: got %s, want heal_lv1", loaded.Skills[1].TypeID)
	}
}

// TestAgentSlotSave_EmptySlot は空のAgentSlotSaveのJSON化をテストします。
func TestAgentSlotSave_EmptySlot(t *testing.T) {
	// 空のスロット
	save := AgentSlotSave{}

	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	var loaded AgentSlotSave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 空スロットの検証
	if loaded.CoreTypeID != "" {
		t.Errorf("空スロットのCoreTypeID: got %s, want empty", loaded.CoreTypeID)
	}
	if loaded.CoreLevel != 0 {
		t.Errorf("空スロットのCoreLevel: got %d, want 0", loaded.CoreLevel)
	}
}

// TestInventorySaveDataV3_JSONSerialization は新しいInventorySaveDataのJSON化をテストします。
func TestInventorySaveDataV3_JSONSerialization(t *testing.T) {
	save := &InventorySaveData{
		// 新フィールド
		UniqueCores: &CoreInventorySave{
			Cores: map[string]int{"all_rounder": 5},
		},
		UniqueSkills: &SkillInventorySave{
			Skills: map[string][]string{"physical_lv1": {"damage_bonus"}},
		},
		// レガシーフィールド（後方互換性）
		CoreInstances:   []CoreInstanceSave{},
		ModuleInstances: []ModuleInstanceSave{},
		AgentInstances:  []AgentInstanceSave{},
		MaxCoreSlots:    100,
		MaxModuleSlots:  200,
		MaxAgentSlots:   20,
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded InventorySaveData
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 新フィールドの検証
	if loaded.UniqueCores == nil {
		t.Fatal("UniqueCoresがnilです")
	}
	if loaded.UniqueCores.Cores["all_rounder"] != 5 {
		t.Errorf("UniqueCores[all_rounder]: got %d, want 5", loaded.UniqueCores.Cores["all_rounder"])
	}
	if loaded.UniqueSkills == nil {
		t.Fatal("UniqueSkillsがnilです")
	}
	if len(loaded.UniqueSkills.Skills["physical_lv1"]) != 1 {
		t.Errorf("UniqueSkills[physical_lv1] chain count: got %d, want 1", len(loaded.UniqueSkills.Skills["physical_lv1"]))
	}
}

// TestPlayerSaveDataV3_JSONSerialization は新しいPlayerSaveDataのJSON化をテストします。
func TestPlayerSaveDataV3_JSONSerialization(t *testing.T) {
	save := &PlayerSaveData{
		// レガシーフィールド
		EquippedAgentIDs: [3]string{},
		// 新フィールド
		AgentSlots: [3]AgentSlotSave{
			{
				CoreTypeID: "all_rounder",
				CoreLevel:  5,
				Skills: [4]SkillSlotSaveCfg{
					{TypeID: "physical_lv1", ChainEffectID: "damage_bonus"},
					{},
					{},
					{},
				},
			},
			{},
			{},
		},
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded PlayerSaveData
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 新フィールドの検証
	if loaded.AgentSlots[0].CoreTypeID != "all_rounder" {
		t.Errorf("AgentSlots[0].CoreTypeID: got %s, want all_rounder", loaded.AgentSlots[0].CoreTypeID)
	}
	if loaded.AgentSlots[0].CoreLevel != 5 {
		t.Errorf("AgentSlots[0].CoreLevel: got %d, want 5", loaded.AgentSlots[0].CoreLevel)
	}
	if loaded.AgentSlots[0].Skills[0].TypeID != "physical_lv1" {
		t.Errorf("AgentSlots[0].Skills[0].TypeID: got %s, want physical_lv1", loaded.AgentSlots[0].Skills[0].TypeID)
	}

	// 2つ目と3つ目のスロットは空であること
	if loaded.AgentSlots[1].CoreTypeID != "" {
		t.Errorf("AgentSlots[1] should be empty, got CoreTypeID: %s", loaded.AgentSlots[1].CoreTypeID)
	}
	if loaded.AgentSlots[2].CoreTypeID != "" {
		t.Errorf("AgentSlots[2] should be empty, got CoreTypeID: %s", loaded.AgentSlots[2].CoreTypeID)
	}
}

// TestNewSaveDataV3 は新しいバージョンのNewSaveDataをテストします。
func TestNewSaveDataV3(t *testing.T) {
	saveData := NewSaveData()

	// バージョンが3.0.0であること
	if saveData.Version != "3.0.0" {
		t.Errorf("Version: got %s, want 3.0.0", saveData.Version)
	}

	// 新フィールドが初期化されていること
	if saveData.Inventory.UniqueCores == nil {
		t.Error("UniqueCoresがnilです")
	}
	if saveData.Inventory.UniqueSkills == nil {
		t.Error("UniqueSkillsがnilです")
	}
	if len(saveData.Inventory.UniqueCores.Cores) != 0 {
		t.Errorf("UniqueCores should be empty, got %d", len(saveData.Inventory.UniqueCores.Cores))
	}
	if len(saveData.Inventory.UniqueSkills.Skills) != 0 {
		t.Errorf("UniqueSkills should be empty, got %d", len(saveData.Inventory.UniqueSkills.Skills))
	}

	// PlayerのAgentSlotsが初期化されていること
	for i, slot := range saveData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			t.Errorf("AgentSlots[%d].CoreTypeID should be empty, got %s", i, slot.CoreTypeID)
		}
	}
}

// ==================== タスク7.2: セーブ/ロード機能のテスト ====================

// TestSaveAndLoadV3_FullData は新スキーマでのセーブとロードをテストします。
func TestSaveAndLoadV3_FullData(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// v3.0.0形式のセーブデータを作成
	saveData := NewSaveData()

	// ユニークコアを追加
	saveData.Inventory.UniqueCores.Cores["all_rounder"] = 5
	saveData.Inventory.UniqueCores.Cores["attack_balance"] = 3

	// ユニークスキルを追加
	saveData.Inventory.UniqueSkills.Skills["physical_lv1"] = []string{"damage_bonus", "life_steal"}
	saveData.Inventory.UniqueSkills.Skills["heal_lv1"] = []string{}

	// エージェントスロットを設定
	saveData.Player.AgentSlots[0] = AgentSlotSave{
		CoreTypeID: "all_rounder",
		CoreLevel:  5,
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1", ChainEffectID: "damage_bonus"},
			{TypeID: "heal_lv1"},
			{},
			{},
		},
	}

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// バージョン検証
	if loadedData.Version != "3.0.0" {
		t.Errorf("Version: got %s, want 3.0.0", loadedData.Version)
	}

	// ユニークコアの検証
	if loadedData.Inventory.UniqueCores == nil {
		t.Fatal("UniqueCoresがnilです")
	}
	if loadedData.Inventory.UniqueCores.Cores["all_rounder"] != 5 {
		t.Errorf("all_rounder level: got %d, want 5", loadedData.Inventory.UniqueCores.Cores["all_rounder"])
	}
	if loadedData.Inventory.UniqueCores.Cores["attack_balance"] != 3 {
		t.Errorf("attack_balance level: got %d, want 3", loadedData.Inventory.UniqueCores.Cores["attack_balance"])
	}

	// ユニークスキルの検証
	if loadedData.Inventory.UniqueSkills == nil {
		t.Fatal("UniqueSkillsがnilです")
	}
	if len(loadedData.Inventory.UniqueSkills.Skills["physical_lv1"]) != 2 {
		t.Errorf("physical_lv1 chain count: got %d, want 2", len(loadedData.Inventory.UniqueSkills.Skills["physical_lv1"]))
	}

	// エージェントスロットの検証
	slot0 := loadedData.Player.AgentSlots[0]
	if slot0.CoreTypeID != "all_rounder" {
		t.Errorf("AgentSlots[0].CoreTypeID: got %s, want all_rounder", slot0.CoreTypeID)
	}
	if slot0.CoreLevel != 5 {
		t.Errorf("AgentSlots[0].CoreLevel: got %d, want 5", slot0.CoreLevel)
	}
	if slot0.Skills[0].TypeID != "physical_lv1" {
		t.Errorf("AgentSlots[0].Skills[0].TypeID: got %s, want physical_lv1", slot0.Skills[0].TypeID)
	}
	if slot0.Skills[0].ChainEffectID != "damage_bonus" {
		t.Errorf("AgentSlots[0].Skills[0].ChainEffectID: got %s, want damage_bonus", slot0.Skills[0].ChainEffectID)
	}
}

// TestConvertCoreInventoryToSave はCoreInventoryからセーブ形式への変換をテストします。
func TestConvertCoreInventoryToSave(t *testing.T) {
	// ドメインモデルを作成
	cores := map[string]int{
		"all_rounder":    5,
		"attack_balance": 3,
	}

	// セーブ形式に変換
	save := ConvertCoreInventoryToSave(cores)

	// 検証
	if len(save.Cores) != 2 {
		t.Errorf("Cores count: got %d, want 2", len(save.Cores))
	}
	if save.Cores["all_rounder"] != 5 {
		t.Errorf("all_rounder: got %d, want 5", save.Cores["all_rounder"])
	}
	if save.Cores["attack_balance"] != 3 {
		t.Errorf("attack_balance: got %d, want 3", save.Cores["attack_balance"])
	}
}

// TestConvertSaveToCoreInventory はセーブ形式からCoreInventoryへの変換をテストします。
func TestConvertSaveToCoreInventory(t *testing.T) {
	save := &CoreInventorySave{
		Cores: map[string]int{
			"all_rounder":    5,
			"attack_balance": 3,
		},
	}

	// マスタデータに存在するTypeIDのセット
	validTypeIDs := map[string]bool{
		"all_rounder": true,
		// attack_balanceはマスタに存在しない
	}

	// ドメイン形式に変換（存在しないTypeIDは無視）
	cores := ConvertSaveToCoreInventory(save, validTypeIDs)

	// 検証: マスタに存在するもののみ復元
	if len(cores) != 1 {
		t.Errorf("Cores count: got %d, want 1", len(cores))
	}
	if cores["all_rounder"] != 5 {
		t.Errorf("all_rounder: got %d, want 5", cores["all_rounder"])
	}
	if _, exists := cores["attack_balance"]; exists {
		t.Error("attack_balance should be ignored (not in master data)")
	}
}

// TestConvertSkillInventoryToSave はSkillInventoryからセーブ形式への変換をテストします。
func TestConvertSkillInventoryToSave(t *testing.T) {
	// ドメインモデルを作成
	skills := map[string][]string{
		"physical_lv1": {"damage_bonus", "life_steal"},
		"heal_lv1":     {},
	}

	// セーブ形式に変換
	save := ConvertSkillInventoryToSave(skills)

	// 検証
	if len(save.Skills) != 2 {
		t.Errorf("Skills count: got %d, want 2", len(save.Skills))
	}
	if len(save.Skills["physical_lv1"]) != 2 {
		t.Errorf("physical_lv1 chains: got %d, want 2", len(save.Skills["physical_lv1"]))
	}
}

// TestConvertSaveToSkillInventory はセーブ形式からSkillInventoryへの変換をテストします。
func TestConvertSaveToSkillInventory(t *testing.T) {
	save := &SkillInventorySave{
		Skills: map[string][]string{
			"physical_lv1": {"damage_bonus", "life_steal"},
			"heal_lv1":     {},
			"unknown":      {"effect1"},
		},
	}

	// マスタデータに存在するTypeIDのセット
	validTypeIDs := map[string]bool{
		"physical_lv1": true,
		"heal_lv1":     true,
		// unknownはマスタに存在しない
	}

	// ドメイン形式に変換（存在しないTypeIDは無視）
	skills := ConvertSaveToSkillInventory(save, validTypeIDs)

	// 検証: マスタに存在するもののみ復元
	if len(skills) != 2 {
		t.Errorf("Skills count: got %d, want 2", len(skills))
	}
	if _, exists := skills["unknown"]; exists {
		t.Error("unknown should be ignored (not in master data)")
	}
}

// TestConvertAgentSlotsToSave はAgentSlotからセーブ形式への変換をテストします。
func TestConvertAgentSlotsToSave(t *testing.T) {
	// ドメインモデル形式（CoreTypeID, CoreLevel, Skills配列）
	slots := [3]struct {
		CoreTypeID string
		CoreLevel  int
		Skills     [4]struct {
			TypeID        string
			ChainEffectID string
		}
	}{
		{
			CoreTypeID: "all_rounder",
			CoreLevel:  5,
			Skills: [4]struct {
				TypeID        string
				ChainEffectID string
			}{
				{TypeID: "physical_lv1", ChainEffectID: "damage_bonus"},
				{TypeID: "heal_lv1"},
				{},
				{},
			},
		},
		{},
		{},
	}

	// セーブ形式に変換
	var saves [3]AgentSlotSave
	for i, slot := range slots {
		saves[i] = ConvertAgentSlotToSave(slot.CoreTypeID, slot.CoreLevel, slot.Skills)
	}

	// 検証
	if saves[0].CoreTypeID != "all_rounder" {
		t.Errorf("Slot 0 CoreTypeID: got %s, want all_rounder", saves[0].CoreTypeID)
	}
	if saves[0].CoreLevel != 5 {
		t.Errorf("Slot 0 CoreLevel: got %d, want 5", saves[0].CoreLevel)
	}
	if saves[0].Skills[0].TypeID != "physical_lv1" {
		t.Errorf("Slot 0 Skills[0].TypeID: got %s, want physical_lv1", saves[0].Skills[0].TypeID)
	}
	if saves[1].CoreTypeID != "" {
		t.Errorf("Slot 1 should be empty, got CoreTypeID: %s", saves[1].CoreTypeID)
	}
}

// TestConvertSaveToAgentSlot はセーブ形式からAgentSlotへの変換をテストします。
func TestConvertSaveToAgentSlot(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "all_rounder",
		CoreLevel:  5,
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1", ChainEffectID: "damage_bonus"},
			{TypeID: "unknown"},
			{},
			{},
		},
	}

	// マスタデータに存在するTypeIDのセット
	validCoreTypeIDs := map[string]bool{
		"all_rounder": true,
	}
	validSkillTypeIDs := map[string]bool{
		"physical_lv1": true,
		// unknownはマスタに存在しない
	}

	// ドメイン形式に変換（存在しないTypeIDは無視）
	coreTypeID, coreLevel, skills := ConvertSaveToAgentSlot(save, validCoreTypeIDs, validSkillTypeIDs)

	// 検証
	if coreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", coreTypeID)
	}
	if coreLevel != 5 {
		t.Errorf("CoreLevel: got %d, want 5", coreLevel)
	}
	if skills[0].TypeID != "physical_lv1" {
		t.Errorf("Skills[0].TypeID: got %s, want physical_lv1", skills[0].TypeID)
	}
	if skills[1].TypeID != "" {
		t.Errorf("Skills[1].TypeID should be empty (unknown not in master), got %s", skills[1].TypeID)
	}
}

// TestConvertSaveToAgentSlot_InvalidCore はマスタに存在しないコアが無視されることをテストします。
func TestConvertSaveToAgentSlot_InvalidCore(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "invalid_core",
		CoreLevel:  5,
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{},
			{},
			{},
		},
	}

	validCoreTypeIDs := map[string]bool{
		"all_rounder": true,
		// invalid_coreは存在しない
	}
	validSkillTypeIDs := map[string]bool{
		"physical_lv1": true,
	}

	// 変換（コアが無効な場合は全体が空になる）
	coreTypeID, coreLevel, skills := ConvertSaveToAgentSlot(save, validCoreTypeIDs, validSkillTypeIDs)

	// コアが無効なのでスロット全体が空
	if coreTypeID != "" {
		t.Errorf("CoreTypeID should be empty for invalid core, got %s", coreTypeID)
	}
	if coreLevel != 0 {
		t.Errorf("CoreLevel should be 0 for invalid core, got %d", coreLevel)
	}
	for i, skill := range skills {
		if skill.TypeID != "" {
			t.Errorf("Skills[%d].TypeID should be empty for invalid core, got %s", i, skill.TypeID)
		}
	}
}
