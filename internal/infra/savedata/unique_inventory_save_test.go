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
