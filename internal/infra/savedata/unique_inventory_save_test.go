// Package savedata はセーブデータの永続化を担当します。
// このファイルはユニークインベントリ関連のセーブデータのテストを提供します。

package savedata

import (
	"encoding/json"
	"testing"
)

// ==================== タスク7.1: 新セーブデータ構造体のテスト ====================

// TestCoreInventorySave_JSONSerialization はCoreInventorySaveのJSON化をテストします。
// コアはTypeIDリスト形式で管理します。
func TestCoreInventorySave_JSONSerialization(t *testing.T) {
	save := CoreInventorySave{
		Cores: []string{"all_rounder", "attack_balance"},
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
	// TypeIDリストに含まれていることを確認
	hasAllRounder := false
	hasAttackBalance := false
	for _, c := range loaded.Cores {
		if c == "all_rounder" {
			hasAllRounder = true
		}
		if c == "attack_balance" {
			hasAttackBalance = true
		}
	}
	if !hasAllRounder {
		t.Error("all_rounder should be in Cores")
	}
	if !hasAttackBalance {
		t.Error("attack_balance should be in Cores")
	}
}

// TestSkillInventorySave_JSONSerialization はSkillInventorySaveのJSON化をテストします。
func TestSkillInventorySave_JSONSerialization(t *testing.T) {
	save := SkillInventorySave{
		Skills: []string{"physical_lv1", "heal_lv1"},
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
	if loaded.Skills[0] != "physical_lv1" {
		t.Errorf("Skills[0]: got %s, want physical_lv1", loaded.Skills[0])
	}
	if loaded.Skills[1] != "heal_lv1" {
		t.Errorf("Skills[1]: got %s, want heal_lv1", loaded.Skills[1])
	}
}

// TestSkillSlotSaveCfg_JSONSerialization はSkillSlotSaveCfgのJSON化をテストします。
func TestSkillSlotSaveCfg_JSONSerialization(t *testing.T) {
	// スキルスロット
	withSkill := SkillSlotSaveCfg{
		TypeID: "physical_lv1",
	}

	data, err := json.Marshal(withSkill)
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

// TestChainEffectInventorySave_JSONSerialization はChainEffectInventorySaveのJSON化をテストします。
func TestChainEffectInventorySave_JSONSerialization(t *testing.T) {
	save := ChainEffectInventorySave{
		ChainEffects: []string{"damage_bonus", "heal_bonus", "life_steal"},
	}

	// JSON化
	data, err := json.Marshal(save)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	// 復元
	var loaded ChainEffectInventorySave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 検証
	if len(loaded.ChainEffects) != 3 {
		t.Errorf("ChainEffects count: got %d, want 3", len(loaded.ChainEffects))
	}
	if loaded.ChainEffects[0] != "damage_bonus" {
		t.Errorf("ChainEffects[0]: got %s, want damage_bonus", loaded.ChainEffects[0])
	}
	if loaded.ChainEffects[1] != "heal_bonus" {
		t.Errorf("ChainEffects[1]: got %s, want heal_bonus", loaded.ChainEffects[1])
	}
	if loaded.ChainEffects[2] != "life_steal" {
		t.Errorf("ChainEffects[2]: got %s, want life_steal", loaded.ChainEffects[2])
	}
}

// TestChainEffectSlotSaveCfg_JSONSerialization はChainEffectSlotSaveCfgのJSON化をテストします。
func TestChainEffectSlotSaveCfg_JSONSerialization(t *testing.T) {
	// チェイン効果ありのスロット
	withEffect := ChainEffectSlotSaveCfg{
		TypeID: "damage_bonus",
	}

	data, err := json.Marshal(withEffect)
	if err != nil {
		t.Fatalf("JSON化に失敗: %v", err)
	}

	var loaded ChainEffectSlotSaveCfg
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	if loaded.TypeID != "damage_bonus" {
		t.Errorf("TypeID: got %s, want damage_bonus", loaded.TypeID)
	}

	// 空のスロット（omitemptyで省略される）
	empty := ChainEffectSlotSaveCfg{}
	emptyData, err := json.Marshal(empty)
	if err != nil {
		t.Fatalf("空スロットのJSON化に失敗: %v", err)
	}

	if string(emptyData) != "{}" {
		t.Errorf("空スロットのJSON: got %s, want {}", string(emptyData))
	}
}

// TestAgentSlotSave_JSONSerialization はAgentSlotSaveのJSON化をテストします。
func TestAgentSlotSave_JSONSerialization(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{TypeID: "heal_lv1"},
			{},
			{},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "damage_bonus"},
			{},
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
	var loaded AgentSlotSave
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("復元に失敗: %v", err)
	}

	// 検証
	if loaded.CoreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", loaded.CoreTypeID)
	}
	if loaded.Skills[0].TypeID != "physical_lv1" {
		t.Errorf("Skills[0].TypeID: got %s, want physical_lv1", loaded.Skills[0].TypeID)
	}
	if loaded.Skills[1].TypeID != "heal_lv1" {
		t.Errorf("Skills[1].TypeID: got %s, want heal_lv1", loaded.Skills[1].TypeID)
	}
	if loaded.ChainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("ChainEffects[0].TypeID: got %s, want damage_bonus", loaded.ChainEffects[0].TypeID)
	}
	if loaded.ChainEffects[1].TypeID != "" {
		t.Errorf("ChainEffects[1].TypeID: got %s, want empty", loaded.ChainEffects[1].TypeID)
	}
}

// TestAgentSlotSave_EmptySlot は空のAgentSlotSaveのJSON化をテストします。
// CoreLevelフィールドを削除。
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
}

// TestInventorySaveData_JSONSerialization はInventorySaveDataのJSON化をテストします。
func TestInventorySaveData_JSONSerialization(t *testing.T) {
	save := &InventorySaveData{
		UniqueCores: &CoreInventorySave{
			Cores: []string{"all_rounder"},
		},
		UniqueSkills: &SkillInventorySave{
			Skills: []string{"physical_lv1"},
		},
		UniqueChainEffects: &ChainEffectInventorySave{
			ChainEffects: []string{"damage_bonus"},
		},
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

	// UniqueCoresの検証
	if loaded.UniqueCores == nil {
		t.Fatal("UniqueCoresがnilです")
	}
	if len(loaded.UniqueCores.Cores) != 1 {
		t.Errorf("UniqueCores count: got %d, want 1", len(loaded.UniqueCores.Cores))
	}
	if loaded.UniqueCores.Cores[0] != "all_rounder" {
		t.Errorf("UniqueCores[0]: got %s, want all_rounder", loaded.UniqueCores.Cores[0])
	}

	// UniqueSkillsの検証
	if loaded.UniqueSkills == nil {
		t.Fatal("UniqueSkillsがnilです")
	}
	if len(loaded.UniqueSkills.Skills) != 1 {
		t.Errorf("UniqueSkills count: got %d, want 1", len(loaded.UniqueSkills.Skills))
	}
	if loaded.UniqueSkills.Skills[0] != "physical_lv1" {
		t.Errorf("UniqueSkills[0]: got %s, want physical_lv1", loaded.UniqueSkills.Skills[0])
	}

	// UniqueChainEffectsの検証
	if loaded.UniqueChainEffects == nil {
		t.Fatal("UniqueChainEffectsがnilです")
	}
	if len(loaded.UniqueChainEffects.ChainEffects) != 1 {
		t.Errorf("UniqueChainEffects count: got %d, want 1", len(loaded.UniqueChainEffects.ChainEffects))
	}
	if loaded.UniqueChainEffects.ChainEffects[0] != "damage_bonus" {
		t.Errorf("UniqueChainEffects[0]: got %s, want damage_bonus", loaded.UniqueChainEffects.ChainEffects[0])
	}
}

// TestPlayerSaveData_JSONSerialization はPlayerSaveDataのJSON化をテストします。
func TestPlayerSaveData_JSONSerialization(t *testing.T) {
	save := &PlayerSaveData{
		MaxHP: 1000,
		AgentSlots: [3]AgentSlotSave{
			{
				CoreTypeID: "all_rounder",
				Skills: [4]SkillSlotSaveCfg{
					{TypeID: "physical_lv1"},
					{},
					{},
					{},
				},
				ChainEffects: [4]ChainEffectSlotSaveCfg{
					{TypeID: "damage_bonus"},
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

	// MaxHPの検証
	if loaded.MaxHP != 1000 {
		t.Errorf("MaxHP: got %d, want 1000", loaded.MaxHP)
	}
	// AgentSlotsの検証
	if loaded.AgentSlots[0].CoreTypeID != "all_rounder" {
		t.Errorf("AgentSlots[0].CoreTypeID: got %s, want all_rounder", loaded.AgentSlots[0].CoreTypeID)
	}
	if loaded.AgentSlots[0].Skills[0].TypeID != "physical_lv1" {
		t.Errorf("AgentSlots[0].Skills[0].TypeID: got %s, want physical_lv1", loaded.AgentSlots[0].Skills[0].TypeID)
	}
	if loaded.AgentSlots[0].ChainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("AgentSlots[0].ChainEffects[0].TypeID: got %s, want damage_bonus", loaded.AgentSlots[0].ChainEffects[0].TypeID)
	}

	// 2つ目と3つ目のスロットは空であること
	if loaded.AgentSlots[1].CoreTypeID != "" {
		t.Errorf("AgentSlots[1] should be empty, got CoreTypeID: %s", loaded.AgentSlots[1].CoreTypeID)
	}
	if loaded.AgentSlots[2].CoreTypeID != "" {
		t.Errorf("AgentSlots[2] should be empty, got CoreTypeID: %s", loaded.AgentSlots[2].CoreTypeID)
	}
}

// TestNewSaveData_FullInitialization はNewSaveDataの完全な初期化をテストします。
func TestNewSaveData_FullInitialization(t *testing.T) {
	saveData := NewSaveData()

	// バージョンが0.0.1であること
	if saveData.Version != "0.0.1" {
		t.Errorf("Version: got %s, want 0.0.1", saveData.Version)
	}

	// フィールドが初期化されていること
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

	// UniqueChainEffectsが初期化されていること
	if saveData.Inventory.UniqueChainEffects == nil {
		t.Error("UniqueChainEffectsがnilです")
	}
	if len(saveData.Inventory.UniqueChainEffects.ChainEffects) != 0 {
		t.Errorf("UniqueChainEffects should be empty, got %d", len(saveData.Inventory.UniqueChainEffects.ChainEffects))
	}

	// PlayerのAgentSlotsが初期化されていること
	for i, slot := range saveData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			t.Errorf("AgentSlots[%d].CoreTypeID should be empty, got %s", i, slot.CoreTypeID)
		}
	}

	// MaxHPの初期値が1000であること
	if saveData.Player.MaxHP != 1000 {
		t.Errorf("Player.MaxHP: got %d, want 1000", saveData.Player.MaxHP)
	}

	// EnemyProgressが初期化されていること
	if saveData.EnemyProgress == nil {
		t.Error("EnemyProgressがnilです")
	}
	if saveData.EnemyProgress.CurrentRank != 1 {
		t.Errorf("EnemyProgress.CurrentRank: got %d, want 1", saveData.EnemyProgress.CurrentRank)
	}
	if len(saveData.EnemyProgress.DefeatRecords) != 0 {
		t.Errorf("EnemyProgress.DefeatRecords should be empty, got %d", len(saveData.EnemyProgress.DefeatRecords))
	}
}

// ==================== タスク7.2: セーブ/ロード機能のテスト ====================

// TestSaveAndLoadV4_FullData は新スキーマでのセーブとロードをテストします。
// TestSaveAndLoad_FullData はセーブデータのフル形式でのセーブ・ロードをテストします。
func TestSaveAndLoad_FullData(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成
	saveData := NewSaveData()

	// ユニークコアを追加（TypeIDリスト形式）
	saveData.Inventory.UniqueCores.Cores = append(saveData.Inventory.UniqueCores.Cores, "all_rounder", "attack_balance")

	// ユニークスキルを追加
	saveData.Inventory.UniqueSkills.Skills = append(saveData.Inventory.UniqueSkills.Skills,
		"physical_lv1", "heal_lv1")

	// ユニークチェイン効果を追加
	saveData.Inventory.UniqueChainEffects.ChainEffects = append(saveData.Inventory.UniqueChainEffects.ChainEffects,
		"damage_bonus", "life_steal")

	// エージェントスロットを設定
	saveData.Player.AgentSlots[0] = AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{TypeID: "heal_lv1"},
			{},
			{},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "damage_bonus"},
			{},
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
	if loadedData.Version != "0.0.1" {
		t.Errorf("Version: got %s, want 0.0.1", loadedData.Version)
	}

	// MaxHP検証
	if loadedData.Player.MaxHP != 1000 {
		t.Errorf("Player.MaxHP: got %d, want 1000", loadedData.Player.MaxHP)
	}

	// ユニークコアの検証（TypeIDリスト形式）
	if loadedData.Inventory.UniqueCores == nil {
		t.Fatal("UniqueCoresがnilです")
	}
	if len(loadedData.Inventory.UniqueCores.Cores) != 2 {
		t.Errorf("UniqueCores count: got %d, want 2", len(loadedData.Inventory.UniqueCores.Cores))
	}

	// ユニークスキルの検証
	if loadedData.Inventory.UniqueSkills == nil {
		t.Fatal("UniqueSkillsがnilです")
	}
	if len(loadedData.Inventory.UniqueSkills.Skills) != 2 {
		t.Errorf("UniqueSkills count: got %d, want 2", len(loadedData.Inventory.UniqueSkills.Skills))
	}

	// ユニークチェイン効果の検証
	if loadedData.Inventory.UniqueChainEffects == nil {
		t.Fatal("UniqueChainEffectsがnilです")
	}
	if len(loadedData.Inventory.UniqueChainEffects.ChainEffects) != 2 {
		t.Errorf("UniqueChainEffects count: got %d, want 2", len(loadedData.Inventory.UniqueChainEffects.ChainEffects))
	}

	// エージェントスロットの検証
	slot0 := loadedData.Player.AgentSlots[0]
	if slot0.CoreTypeID != "all_rounder" {
		t.Errorf("AgentSlots[0].CoreTypeID: got %s, want all_rounder", slot0.CoreTypeID)
	}
	if slot0.Skills[0].TypeID != "physical_lv1" {
		t.Errorf("AgentSlots[0].Skills[0].TypeID: got %s, want physical_lv1", slot0.Skills[0].TypeID)
	}
	if slot0.ChainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("AgentSlots[0].ChainEffects[0].TypeID: got %s, want damage_bonus", slot0.ChainEffects[0].TypeID)
	}
}

// TestConvertCoreInventoryToSave はCoreInventoryからセーブ形式への変換をテストします。
// コアはTypeIDリスト形式で管理します。
func TestConvertCoreInventoryToSave(t *testing.T) {
	// ドメインモデルを作成（TypeIDリスト）
	cores := []string{"all_rounder", "attack_balance"}

	// セーブ形式に変換
	save := ConvertCoreInventoryToSave(cores)

	// 検証
	if len(save.Cores) != 2 {
		t.Errorf("Cores count: got %d, want 2", len(save.Cores))
	}
	// TypeIDリストに含まれていることを確認
	hasAllRounder := false
	hasAttackBalance := false
	for _, c := range save.Cores {
		if c == "all_rounder" {
			hasAllRounder = true
		}
		if c == "attack_balance" {
			hasAttackBalance = true
		}
	}
	if !hasAllRounder {
		t.Error("all_rounder should be in Cores")
	}
	if !hasAttackBalance {
		t.Error("attack_balance should be in Cores")
	}
}

// TestConvertSaveToCoreInventory はセーブ形式からCoreInventoryへの変換をテストします。
// コアはTypeIDリスト形式で管理します。
func TestConvertSaveToCoreInventory(t *testing.T) {
	save := &CoreInventorySave{
		Cores: []string{"all_rounder", "attack_balance"},
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
	// all_rounderのみが返されることを確認
	hasAllRounder := false
	hasAttackBalance := false
	for _, c := range cores {
		if c == "all_rounder" {
			hasAllRounder = true
		}
		if c == "attack_balance" {
			hasAttackBalance = true
		}
	}
	if !hasAllRounder {
		t.Error("all_rounder should be in cores")
	}
	if hasAttackBalance {
		t.Error("attack_balance should be ignored (not in master data)")
	}
}

// TestConvertSkillInventoryToSave はSkillInventoryからセーブ形式への変換をテストします。
func TestConvertSkillInventoryToSave(t *testing.T) {
	// ドメインモデルを作成（TypeIDリスト形式）
	skills := []string{"physical_lv1", "heal_lv1"}

	// セーブ形式に変換
	save := ConvertSkillInventoryToSave(skills)

	// 検証
	if len(save.Skills) != 2 {
		t.Errorf("Skills count: got %d, want 2", len(save.Skills))
	}
	if save.Skills[0] != "physical_lv1" {
		t.Errorf("Skills[0]: got %s, want physical_lv1", save.Skills[0])
	}
	if save.Skills[1] != "heal_lv1" {
		t.Errorf("Skills[1]: got %s, want heal_lv1", save.Skills[1])
	}
}

// TestConvertSaveToSkillInventory はセーブ形式からSkillInventoryへの変換をテストします。
func TestConvertSaveToSkillInventory(t *testing.T) {
	save := &SkillInventorySave{
		Skills: []string{"physical_lv1", "heal_lv1", "unknown"},
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
	// unknownが含まれていないことを確認
	for _, s := range skills {
		if s == "unknown" {
			t.Error("unknown should be ignored (not in master data)")
		}
	}
}

// TestConvertChainEffectInventoryToSave はChainEffectInventoryからセーブ形式への変換をテストします。
func TestConvertChainEffectInventoryToSave(t *testing.T) {
	// ドメインモデルを作成（TypeIDリスト形式）
	chainEffects := []string{"damage_bonus", "heal_bonus"}

	// セーブ形式に変換
	save := ConvertChainEffectInventoryToSave(chainEffects)

	// 検証
	if len(save.ChainEffects) != 2 {
		t.Errorf("ChainEffects count: got %d, want 2", len(save.ChainEffects))
	}
	if save.ChainEffects[0] != "damage_bonus" {
		t.Errorf("ChainEffects[0]: got %s, want damage_bonus", save.ChainEffects[0])
	}
	if save.ChainEffects[1] != "heal_bonus" {
		t.Errorf("ChainEffects[1]: got %s, want heal_bonus", save.ChainEffects[1])
	}
}

// TestConvertSaveToChainEffectInventory はセーブ形式からChainEffectInventoryへの変換をテストします。
func TestConvertSaveToChainEffectInventory(t *testing.T) {
	save := &ChainEffectInventorySave{
		ChainEffects: []string{"damage_bonus", "heal_bonus", "unknown"},
	}

	// マスタデータに存在するTypeIDのセット
	validTypeIDs := map[string]bool{
		"damage_bonus": true,
		"heal_bonus":   true,
		// unknownはマスタに存在しない
	}

	// ドメイン形式に変換（存在しないTypeIDは無視）
	chainEffects := ConvertSaveToChainEffectInventory(save, validTypeIDs)

	// 検証: マスタに存在するもののみ復元
	if len(chainEffects) != 2 {
		t.Errorf("ChainEffects count: got %d, want 2", len(chainEffects))
	}
	// unknownが含まれていないことを確認
	for _, ce := range chainEffects {
		if ce == "unknown" {
			t.Error("unknown should be ignored (not in master data)")
		}
	}
}

// TestConvertAgentSlotsToSave はAgentSlotからセーブ形式への変換をテストします。
func TestConvertAgentSlotsToSave(t *testing.T) {
	// ドメインモデル形式（CoreTypeID, Skills配列, ChainEffects配列）
	slots := [3]struct {
		CoreTypeID   string
		Skills       [4]struct{ TypeID string }
		ChainEffects [4]struct{ TypeID string }
	}{
		{
			CoreTypeID: "all_rounder",
			Skills: [4]struct{ TypeID string }{
				{TypeID: "physical_lv1"},
				{TypeID: "heal_lv1"},
				{},
				{},
			},
			ChainEffects: [4]struct{ TypeID string }{
				{TypeID: "damage_bonus"},
				{},
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
		saves[i] = ConvertAgentSlotToSave(slot.CoreTypeID, slot.Skills, slot.ChainEffects)
	}

	// 検証
	if saves[0].CoreTypeID != "all_rounder" {
		t.Errorf("Slot 0 CoreTypeID: got %s, want all_rounder", saves[0].CoreTypeID)
	}
	if saves[0].Skills[0].TypeID != "physical_lv1" {
		t.Errorf("Slot 0 Skills[0].TypeID: got %s, want physical_lv1", saves[0].Skills[0].TypeID)
	}
	if saves[0].ChainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("Slot 0 ChainEffects[0].TypeID: got %s, want damage_bonus", saves[0].ChainEffects[0].TypeID)
	}
	if saves[1].CoreTypeID != "" {
		t.Errorf("Slot 1 should be empty, got CoreTypeID: %s", saves[1].CoreTypeID)
	}
}

// TestConvertSaveToAgentSlot はセーブ形式からAgentSlotへの変換をテストします。
func TestConvertSaveToAgentSlot(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{TypeID: "unknown"},
			{},
			{},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "damage_bonus"},
			{TypeID: "invalid_chain"},
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
	validChainEffectTypeIDs := map[string]bool{
		"damage_bonus": true,
		// invalid_chainはマスタに存在しない
	}

	// ドメイン形式に変換（存在しないTypeIDは無視）
	coreTypeID, skills, chainEffects := ConvertSaveToAgentSlot(save, validCoreTypeIDs, validSkillTypeIDs, validChainEffectTypeIDs)

	// 検証
	if coreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", coreTypeID)
	}
	if skills[0].TypeID != "physical_lv1" {
		t.Errorf("Skills[0].TypeID: got %s, want physical_lv1", skills[0].TypeID)
	}
	if skills[1].TypeID != "" {
		t.Errorf("Skills[1].TypeID should be empty (unknown not in master), got %s", skills[1].TypeID)
	}
	if chainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("ChainEffects[0].TypeID: got %s, want damage_bonus", chainEffects[0].TypeID)
	}
	if chainEffects[1].TypeID != "" {
		t.Errorf("ChainEffects[1].TypeID should be empty (invalid_chain not in master), got %s", chainEffects[1].TypeID)
	}
}

// TestConvertSaveToAgentSlot_InvalidCore はマスタに存在しないコアが無視されることをテストします。
func TestConvertSaveToAgentSlot_InvalidCore(t *testing.T) {
	save := AgentSlotSave{
		CoreTypeID: "invalid_core",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{},
			{},
			{},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "damage_bonus"},
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
	validChainEffectTypeIDs := map[string]bool{
		"damage_bonus": true,
	}

	// 変換（コアが無効な場合は全体が空になる）
	coreTypeID, skills, chainEffects := ConvertSaveToAgentSlot(save, validCoreTypeIDs, validSkillTypeIDs, validChainEffectTypeIDs)

	// コアが無効なのでスロット全体が空
	if coreTypeID != "" {
		t.Errorf("CoreTypeID should be empty for invalid core, got %s", coreTypeID)
	}
	for i, skill := range skills {
		if skill.TypeID != "" {
			t.Errorf("Skills[%d].TypeID should be empty for invalid core, got %s", i, skill.TypeID)
		}
	}
	for i, ce := range chainEffects {
		if ce.TypeID != "" {
			t.Errorf("ChainEffects[%d].TypeID should be empty for invalid core, got %s", i, ce.TypeID)
		}
	}
}

// ==================== タスク7.3: セーブ/ロードの統合テスト ====================

// TestSaveLoadIntegration_FullCycle はセーブデータの保存・復元の統合テストです。
func TestSaveLoadIntegration_FullCycle(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 1. セーブデータを作成
	saveData := NewSaveData()

	// ユニークコアを設定（TypeIDリスト形式）
	saveData.Inventory.UniqueCores.Cores = append(saveData.Inventory.UniqueCores.Cores,
		"all_rounder", "attack_balance", "defense_balance")

	// ユニークスキルを設定（TypeIDリスト形式）
	saveData.Inventory.UniqueSkills.Skills = append(saveData.Inventory.UniqueSkills.Skills,
		"physical_lv1", "heal_lv1", "buff_lv1")

	// ユニークチェイン効果を設定
	saveData.Inventory.UniqueChainEffects.ChainEffects = append(saveData.Inventory.UniqueChainEffects.ChainEffects,
		"damage_bonus", "heal_bonus", "life_steal")

	// エージェントスロット0を設定（フル装備）
	saveData.Player.AgentSlots[0] = AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{TypeID: "heal_lv1"},
			{TypeID: "buff_lv1"},
			{TypeID: "physical_lv1"},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "damage_bonus"},
			{TypeID: "heal_bonus"},
			{},
			{TypeID: "life_steal"},
		},
	}

	// エージェントスロット1を設定（一部装備）
	saveData.Player.AgentSlots[1] = AgentSlotSave{
		CoreTypeID: "attack_balance",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "physical_lv1"},
			{},
			{},
			{},
		},
	}

	// エージェントスロット2は空のまま

	// 統計も設定
	saveData.Statistics.TotalBattles = 100
	saveData.Statistics.Victories = 80

	// 2. セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// 3. ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 4. 検証

	// バージョン
	if loadedData.Version != "0.0.1" {
		t.Errorf("Version: got %s, want 0.0.1", loadedData.Version)
	}

	// MaxHP
	if loadedData.Player.MaxHP != 1000 {
		t.Errorf("Player.MaxHP: got %d, want 1000", loadedData.Player.MaxHP)
	}

	// ユニークコア（TypeIDリスト形式）
	if len(loadedData.Inventory.UniqueCores.Cores) != 3 {
		t.Errorf("UniqueCores count: got %d, want 3", len(loadedData.Inventory.UniqueCores.Cores))
	}
	expectedCores := map[string]bool{"all_rounder": true, "attack_balance": true, "defense_balance": true}
	for _, c := range loadedData.Inventory.UniqueCores.Cores {
		if !expectedCores[c] {
			t.Errorf("Unexpected core in UniqueCores: %s", c)
		}
	}

	// ユニークスキル
	if len(loadedData.Inventory.UniqueSkills.Skills) != 3 {
		t.Errorf("UniqueSkills count: got %d, want 3", len(loadedData.Inventory.UniqueSkills.Skills))
	}

	// ユニークチェイン効果
	if len(loadedData.Inventory.UniqueChainEffects.ChainEffects) != 3 {
		t.Errorf("UniqueChainEffects count: got %d, want 3", len(loadedData.Inventory.UniqueChainEffects.ChainEffects))
	}

	// エージェントスロット0
	slot0 := loadedData.Player.AgentSlots[0]
	if slot0.CoreTypeID != "all_rounder" {
		t.Errorf("Slot0 CoreTypeID: got %s, want all_rounder", slot0.CoreTypeID)
	}
	if slot0.Skills[0].TypeID != "physical_lv1" {
		t.Errorf("Slot0 Skill0 TypeID: got %s, want physical_lv1", slot0.Skills[0].TypeID)
	}
	if slot0.ChainEffects[0].TypeID != "damage_bonus" {
		t.Errorf("Slot0 ChainEffect0 TypeID: got %s, want damage_bonus", slot0.ChainEffects[0].TypeID)
	}
	if slot0.Skills[3].TypeID != "physical_lv1" {
		t.Errorf("Slot0 Skill3 TypeID: got %s, want physical_lv1", slot0.Skills[3].TypeID)
	}
	if slot0.ChainEffects[3].TypeID != "life_steal" {
		t.Errorf("Slot0 ChainEffect3 TypeID: got %s, want life_steal", slot0.ChainEffects[3].TypeID)
	}

	// エージェントスロット1
	slot1 := loadedData.Player.AgentSlots[1]
	if slot1.CoreTypeID != "attack_balance" {
		t.Errorf("Slot1 CoreTypeID: got %s, want attack_balance", slot1.CoreTypeID)
	}

	// エージェントスロット2（空）
	slot2 := loadedData.Player.AgentSlots[2]
	if slot2.CoreTypeID != "" {
		t.Errorf("Slot2 should be empty, got CoreTypeID: %s", slot2.CoreTypeID)
	}

	// 統計
	if loadedData.Statistics.TotalBattles != 100 {
		t.Errorf("TotalBattles: got %d, want 100", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.Victories != 80 {
		t.Errorf("Victories: got %d, want 80", loadedData.Statistics.Victories)
	}
}

// TestSaveLoadIntegration_InvalidTypeIDsIgnored はマスタに存在しないTypeIDが無視されることをテストします。
func TestSaveLoadIntegration_InvalidTypeIDsIgnored(t *testing.T) {
	// セーブデータにマスタに存在しないTypeIDを含める
	saveData := NewSaveData()

	// 一部有効、一部無効なコア（TypeIDリスト形式）
	saveData.Inventory.UniqueCores.Cores = append(saveData.Inventory.UniqueCores.Cores,
		"valid_core", "invalid_core")

	// 一部有効、一部無効なスキル（TypeIDリスト形式）
	saveData.Inventory.UniqueSkills.Skills = append(saveData.Inventory.UniqueSkills.Skills,
		"valid_skill", "invalid_skill")

	// 一部有効、一部無効なチェイン効果
	saveData.Inventory.UniqueChainEffects.ChainEffects = append(saveData.Inventory.UniqueChainEffects.ChainEffects,
		"valid_chain", "invalid_chain")

	// エージェントスロット: 有効なコアだが一部無効なスキル・チェイン効果
	saveData.Player.AgentSlots[0] = AgentSlotSave{
		CoreTypeID: "valid_core",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "valid_skill"},
			{TypeID: "invalid_skill"},
			{},
			{},
		},
		ChainEffects: [4]ChainEffectSlotSaveCfg{
			{TypeID: "valid_chain"},
			{TypeID: "invalid_chain"},
			{},
			{},
		},
	}

	// エージェントスロット: 無効なコア
	saveData.Player.AgentSlots[1] = AgentSlotSave{
		CoreTypeID: "invalid_core",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "valid_skill"},
			{},
			{},
			{},
		},
	}

	// マスタデータセット（実際にはapp層で用意される）
	validCoreTypeIDs := map[string]bool{
		"valid_core": true,
	}
	validSkillTypeIDs := map[string]bool{
		"valid_skill": true,
	}
	validChainEffectTypeIDs := map[string]bool{
		"valid_chain": true,
	}

	// コアインベントリの変換テスト（TypeIDリスト形式）
	cores := ConvertSaveToCoreInventory(saveData.Inventory.UniqueCores, validCoreTypeIDs)
	if len(cores) != 1 {
		t.Errorf("Cores after filter: got %d, want 1", len(cores))
	}
	hasValidCore := false
	hasInvalidCore := false
	for _, c := range cores {
		if c == "valid_core" {
			hasValidCore = true
		}
		if c == "invalid_core" {
			hasInvalidCore = true
		}
	}
	if !hasValidCore {
		t.Error("valid_core should be in cores")
	}
	if hasInvalidCore {
		t.Error("invalid_core should be filtered out")
	}

	// スキルインベントリの変換テスト
	skills := ConvertSaveToSkillInventory(saveData.Inventory.UniqueSkills, validSkillTypeIDs)
	if len(skills) != 1 {
		t.Errorf("Skills after filter: got %d, want 1", len(skills))
	}
	for _, s := range skills {
		if s == "invalid_skill" {
			t.Error("invalid_skill should be filtered out")
		}
	}

	// チェイン効果インベントリの変換テスト
	chainEffects := ConvertSaveToChainEffectInventory(saveData.Inventory.UniqueChainEffects, validChainEffectTypeIDs)
	if len(chainEffects) != 1 {
		t.Errorf("ChainEffects after filter: got %d, want 1", len(chainEffects))
	}
	for _, ce := range chainEffects {
		if ce == "invalid_chain" {
			t.Error("invalid_chain should be filtered out")
		}
	}

	// エージェントスロット0の変換テスト（コア有効、スキル・チェイン効果一部無効）
	coreTypeID, skillSlots, ceSlots := ConvertSaveToAgentSlot(
		saveData.Player.AgentSlots[0],
		validCoreTypeIDs,
		validSkillTypeIDs,
		validChainEffectTypeIDs,
	)
	if coreTypeID != "valid_core" {
		t.Errorf("Slot0 CoreTypeID: got %s, want valid_core", coreTypeID)
	}
	if skillSlots[0].TypeID != "valid_skill" {
		t.Errorf("Slot0 Skill0: got %s, want valid_skill", skillSlots[0].TypeID)
	}
	if skillSlots[1].TypeID != "" {
		t.Errorf("Slot0 Skill1 should be empty (invalid skill filtered), got %s", skillSlots[1].TypeID)
	}
	if ceSlots[0].TypeID != "valid_chain" {
		t.Errorf("Slot0 ChainEffect0: got %s, want valid_chain", ceSlots[0].TypeID)
	}
	if ceSlots[1].TypeID != "" {
		t.Errorf("Slot0 ChainEffect1 should be empty (invalid chain filtered), got %s", ceSlots[1].TypeID)
	}

	// エージェントスロット1の変換テスト（コア無効→スロット全体空）
	coreTypeID, skillSlots, ceSlots = ConvertSaveToAgentSlot(
		saveData.Player.AgentSlots[1],
		validCoreTypeIDs,
		validSkillTypeIDs,
		validChainEffectTypeIDs,
	)
	if coreTypeID != "" {
		t.Errorf("Slot1 CoreTypeID should be empty (invalid core), got %s", coreTypeID)
	}
	for i, skill := range skillSlots {
		if skill.TypeID != "" {
			t.Errorf("Slot1 Skill%d should be empty when core is invalid, got %s", i, skill.TypeID)
		}
	}
	for i, ce := range ceSlots {
		if ce.TypeID != "" {
			t.Errorf("Slot1 ChainEffect%d should be empty when core is invalid, got %s", i, ce.TypeID)
		}
	}
}

// TestSaveLoadIntegration_EmptyData は空データのセーブ・ロードをテストします。
func TestSaveLoadIntegration_EmptyData(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 空のセーブデータ
	saveData := NewSaveData()

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 空のデータが正しく復元されることを検証
	if loadedData.Inventory.UniqueCores == nil {
		t.Fatal("UniqueCores should not be nil")
	}
	if len(loadedData.Inventory.UniqueCores.Cores) != 0 {
		t.Errorf("UniqueCores should be empty, got %d", len(loadedData.Inventory.UniqueCores.Cores))
	}
	if loadedData.Inventory.UniqueSkills == nil {
		t.Fatal("UniqueSkills should not be nil")
	}
	if len(loadedData.Inventory.UniqueSkills.Skills) != 0 {
		t.Errorf("UniqueSkills should be empty, got %d", len(loadedData.Inventory.UniqueSkills.Skills))
	}
	if loadedData.Inventory.UniqueChainEffects == nil {
		t.Fatal("UniqueChainEffects should not be nil")
	}
	if len(loadedData.Inventory.UniqueChainEffects.ChainEffects) != 0 {
		t.Errorf("UniqueChainEffects should be empty, got %d", len(loadedData.Inventory.UniqueChainEffects.ChainEffects))
	}
	for i, slot := range loadedData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			t.Errorf("AgentSlots[%d] should be empty, got CoreTypeID: %s", i, slot.CoreTypeID)
		}
	}
}

// TestSaveLoadIntegration_NilConversion はnilセーブデータの変換をテストします。
func TestSaveLoadIntegration_NilConversion(t *testing.T) {
	validTypeIDs := map[string]bool{"valid": true}

	// nil CoreInventorySave
	cores := ConvertSaveToCoreInventory(nil, validTypeIDs)
	if len(cores) != 0 {
		t.Errorf("nil CoreInventorySave should return empty slice, got %d", len(cores))
	}

	// nil SkillInventorySave
	skills := ConvertSaveToSkillInventory(nil, validTypeIDs)
	if len(skills) != 0 {
		t.Errorf("nil SkillInventorySave should return empty slice, got %d", len(skills))
	}

	// nil ChainEffectInventorySave
	chainEffects := ConvertSaveToChainEffectInventory(nil, validTypeIDs)
	if len(chainEffects) != 0 {
		t.Errorf("nil ChainEffectInventorySave should return empty slice, got %d", len(chainEffects))
	}
}

// ==================== EnemyProgress 変換テスト ====================

// TestConvertEnemyProgressToSave はEnemyProgressからセーブ形式への変換をテストします。
func TestConvertEnemyProgressToSave(t *testing.T) {
	// ドメインモデルのデータを作成
	currentRank := 2
	defeatRecords := map[string]struct {
		Defeated         bool
		MaxDefeatedLevel int
	}{
		"slime":  {Defeated: true, MaxDefeatedLevel: 10},
		"bat":    {Defeated: true, MaxDefeatedLevel: 5},
		"goblin": {Defeated: false, MaxDefeatedLevel: 0},
	}

	// セーブ形式に変換
	save := ConvertEnemyProgressToSave(currentRank, defeatRecords)

	// 検証
	if save == nil {
		t.Fatal("save should not be nil")
	}
	if save.CurrentRank != 2 {
		t.Errorf("CurrentRank: got %d, want 2", save.CurrentRank)
	}
	if len(save.DefeatRecords) != 3 {
		t.Errorf("DefeatRecords count: got %d, want 3", len(save.DefeatRecords))
	}
	slimeRecord := save.DefeatRecords["slime"]
	if !slimeRecord.Defeated {
		t.Error("slime should be defeated")
	}
	if slimeRecord.MaxDefeatedLevel != 10 {
		t.Errorf("slime MaxDefeatedLevel: got %d, want 10", slimeRecord.MaxDefeatedLevel)
	}
}

// TestConvertSaveToEnemyProgress はセーブ形式からEnemyProgressへの変換をテストします。
func TestConvertSaveToEnemyProgress(t *testing.T) {
	save := &EnemyProgressSave{
		CurrentRank: 3,
		DefeatRecords: map[string]DefeatRecordSave{
			"slime": {Defeated: true, MaxDefeatedLevel: 15},
			"bat":   {Defeated: true, MaxDefeatedLevel: 8},
		},
	}

	// ドメイン形式に変換（DefeatRecordInput型で返される）
	currentRank, defeatRecords := ConvertSaveToEnemyProgress(save)

	// 検証
	if currentRank != 3 {
		t.Errorf("CurrentRank: got %d, want 3", currentRank)
	}
	if len(defeatRecords) != 2 {
		t.Errorf("DefeatRecords count: got %d, want 2", len(defeatRecords))
	}
	slimeRecord, exists := defeatRecords["slime"]
	if !exists {
		t.Fatal("slime should exist")
	}
	if !slimeRecord.Defeated {
		t.Error("slime should be defeated")
	}
	if slimeRecord.MaxDefeatedLevel != 15 {
		t.Errorf("slime MaxDefeatedLevel: got %d, want 15", slimeRecord.MaxDefeatedLevel)
	}
}

// TestConvertSaveToEnemyProgress_NilSave はnilセーブデータの変換をテストします。
func TestConvertSaveToEnemyProgress_NilSave(t *testing.T) {
	// nil EnemyProgressSave
	currentRank, defeatRecords := ConvertSaveToEnemyProgress(nil)

	// 検証：デフォルト値が返される
	if currentRank != 1 {
		t.Errorf("nil save should return rank 1, got %d", currentRank)
	}
	if len(defeatRecords) != 0 {
		t.Errorf("nil save should return empty records, got %d", len(defeatRecords))
	}
}

// TestConvertSaveToEnemyProgress_EmptyRecords は空の記録の変換をテストします。
func TestConvertSaveToEnemyProgress_EmptyRecords(t *testing.T) {
	save := &EnemyProgressSave{
		CurrentRank:   2,
		DefeatRecords: map[string]DefeatRecordSave{},
	}

	currentRank, defeatRecords := ConvertSaveToEnemyProgress(save)

	if currentRank != 2 {
		t.Errorf("CurrentRank: got %d, want 2", currentRank)
	}
	if len(defeatRecords) != 0 {
		t.Errorf("DefeatRecords should be empty, got %d", len(defeatRecords))
	}
}
