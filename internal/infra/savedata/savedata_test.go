// Package savedata はセーブデータの永続化を担当します。
package savedata

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewSaveData は新規SaveDataの作成をテストします。
func TestNewSaveData(t *testing.T) {
	saveData := NewSaveData()

	if saveData.Version == "" {
		t.Error("Versionが空です")
	}
	if saveData.Timestamp.IsZero() {
		t.Error("Timestampが設定されていません")
	}
	if saveData.Player == nil {
		t.Error("Playerがnilです")
	}
	if saveData.Inventory == nil {
		t.Error("Inventoryがnilです")
	}
	if saveData.Statistics == nil {
		t.Error("Statisticsがnilです")
	}
	if saveData.Achievements == nil {
		t.Error("Achievementsがnilです")
	}
	if saveData.Settings == nil {
		t.Error("Settingsがnilです")
	}
}

// TestSaveAndLoadGame はセーブとロードの基本動作をテストします。
func TestSaveAndLoadGame(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// テスト用のセーブデータを作成
	saveData := NewSaveData()
	saveData.Statistics.TotalBattles = 10
	saveData.Statistics.Victories = 8
	saveData.Statistics.MaxLevelReached = 5

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if loadedData.Statistics.TotalBattles != 10 {
		t.Errorf("TotalBattles: got %d, want 10", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.Victories != 8 {
		t.Errorf("Victories: got %d, want 8", loadedData.Statistics.Victories)
	}
	if loadedData.Statistics.MaxLevelReached != 5 {
		t.Errorf("MaxLevelReached: got %d, want 5", loadedData.Statistics.MaxLevelReached)
	}
}

// TestAtomicWrite は原子的書き込み（一時ファイル→リネーム）をテストします。

func TestAtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()

	// セーブ実行
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// 一時ファイルが残っていないことを確認
	tmpFile := filepath.Join(tmpDir, "save.json.tmp")
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Error("一時ファイルが残っています")
	}

	// セーブファイルが存在することを確認
	saveFile := filepath.Join(tmpDir, "save.json")
	if _, err := os.Stat(saveFile); os.IsNotExist(err) {
		t.Error("セーブファイルが作成されていません")
	}
}

// TestBackupRotation はバックアップローテーション（直近3世代）をテストします。

func TestBackupRotation(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 4回セーブして、バックアップローテーションを確認
	for i := 0; i < 4; i++ {
		saveData := NewSaveData()
		saveData.Statistics.TotalBattles = i + 1
		if err := io.SaveGame(saveData); err != nil {
			t.Fatalf("セーブ%dに失敗: %v", i+1, err)
		}
	}

	// バックアップファイルの存在確認
	bak1 := filepath.Join(tmpDir, "save.json.bak1")
	bak2 := filepath.Join(tmpDir, "save.json.bak2")
	bak3 := filepath.Join(tmpDir, "save.json.bak3")

	if _, err := os.Stat(bak1); os.IsNotExist(err) {
		t.Error("save.json.bak1が存在しません")
	}
	if _, err := os.Stat(bak2); os.IsNotExist(err) {
		t.Error("save.json.bak2が存在しません")
	}
	if _, err := os.Stat(bak3); os.IsNotExist(err) {
		t.Error("save.json.bak3が存在しません")
	}
}

// TestLoadFromBackup は破損時のバックアップ復元をテストします。

func TestLoadFromBackup(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 1回目のセーブ（これがバックアップになる）
	saveData1 := NewSaveData()
	saveData1.Statistics.TotalBattles = 100
	if err := io.SaveGame(saveData1); err != nil {
		t.Fatalf("セーブ1に失敗: %v", err)
	}

	// 2回目のセーブ（これによりバックアップが作成される）
	saveData2 := NewSaveData()
	saveData2.Statistics.TotalBattles = 200
	if err := io.SaveGame(saveData2); err != nil {
		t.Fatalf("セーブ2に失敗: %v", err)
	}

	// メインのセーブファイルを破損させる
	saveFile := filepath.Join(tmpDir, "save.json")
	if err := os.WriteFile(saveFile, []byte("corrupted data"), 0644); err != nil {
		t.Fatalf("ファイル破損に失敗: %v", err)
	}

	// ロード（バックアップから復元されるはず）
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("バックアップからのロードに失敗: %v", err)
	}

	// バックアップのデータが読み込まれていることを確認
	// バックアップは2回目セーブ時に作成されるので、1回目のデータ(100)が入っている
	if loadedData.Statistics.TotalBattles != 100 {
		t.Errorf("TotalBattles: got %d, want 100", loadedData.Statistics.TotalBattles)
	}
}

// TestVersionCheck はセーブデータのバージョンチェックをテストします。

func TestVersionCheck(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// バージョンが設定されていることを確認
	if loadedData.Version == "" {
		t.Error("Versionが空です")
	}
}

// TestLoadGameFileNotFound はセーブファイルが存在しない場合のエラーをテストします。
func TestLoadGameFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	_, err := io.LoadGame()
	if err == nil {
		t.Error("ファイルが存在しない場合はエラーが返されるべき")
	}
}

// TestSaveDataWithInventory はインベントリを含むセーブデータをテストします。
// v1.0.0形式のセーブデータ構造をテスト
func TestSaveDataWithInventory(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成（v1.0.0形式: IDなし）
	saveData := NewSaveData()
	saveData.Inventory.CoreInstances = append(saveData.Inventory.CoreInstances, CoreInstanceSave{
		CoreTypeID: "test_core",
		Level:      5,
	})
	// v1.0.0ではModuleInstancesを使用
	saveData.Inventory.ModuleInstances = append(saveData.Inventory.ModuleInstances, ModuleInstanceSave{
		TypeID:      "module_001",
		ChainEffect: nil,
	})

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if len(loadedData.Inventory.CoreInstances) != 1 {
		t.Errorf("CoreInstances: got %d, want 1", len(loadedData.Inventory.CoreInstances))
	}
	if len(loadedData.Inventory.ModuleInstances) != 1 {
		t.Errorf("ModuleInstances: got %d, want 1", len(loadedData.Inventory.ModuleInstances))
	}
	if loadedData.Inventory.CoreInstances[0].CoreTypeID != "test_core" {
		t.Errorf("Core CoreTypeID: got %s, want test_core", loadedData.Inventory.CoreInstances[0].CoreTypeID)
	}
	if loadedData.Inventory.CoreInstances[0].Level != 5 {
		t.Errorf("Core Level: got %d, want 5", loadedData.Inventory.CoreInstances[0].Level)
	}
}

// TestSaveDataWithAgents はエージェントを含むセーブデータをテストします。
// v1.0.0形式のセーブデータ構造をテスト
func TestSaveDataWithAgents(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成（v1.0.0形式: Core.IDなし）
	saveData := NewSaveData()
	// エージェントインスタンスを追加（コア情報を埋め込み）
	saveData.Inventory.AgentInstances = append(saveData.Inventory.AgentInstances, AgentInstanceSave{
		ID: "agent_001",
		Core: CoreInstanceSave{
			CoreTypeID: "test_core",
			Level:      5,
		},
		Skills: []SkillInstanceSave{
			{TypeID: "mod_1"},
			{TypeID: "mod_2"},
			{TypeID: "mod_3"},
			{TypeID: "mod_4"},
		},
	})

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if len(loadedData.Inventory.AgentInstances) != 1 {
		t.Errorf("AgentInstances: got %d, want 1", len(loadedData.Inventory.AgentInstances))
	}
	if loadedData.Inventory.AgentInstances[0].ID != "agent_001" {
		t.Errorf("Agent ID: got %s, want agent_001", loadedData.Inventory.AgentInstances[0].ID)
	}
	if loadedData.Inventory.AgentInstances[0].Core.CoreTypeID != "test_core" {
		t.Errorf("Agent Core.CoreTypeID: got %s, want test_core", loadedData.Inventory.AgentInstances[0].Core.CoreTypeID)
	}
	if loadedData.Inventory.AgentInstances[0].Core.Level != 5 {
		t.Errorf("Agent Core.Level: got %d, want 5", loadedData.Inventory.AgentInstances[0].Core.Level)
	}
	if len(loadedData.Inventory.AgentInstances[0].Skills) != 4 {
		t.Errorf("Agent Skills count: got %d, want 4", len(loadedData.Inventory.AgentInstances[0].Skills))
	}
}

// TestSaveDataTimestamp はタイムスタンプが更新されることをテストします。

func TestSaveDataTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 1回目のセーブ
	saveData1 := NewSaveData()
	time1 := saveData1.Timestamp
	if err := io.SaveGame(saveData1); err != nil {
		t.Fatalf("セーブ1に失敗: %v", err)
	}

	// 少し待機
	time.Sleep(10 * time.Millisecond)

	// 2回目のセーブ
	saveData2 := NewSaveData()
	if err := io.SaveGame(saveData2); err != nil {
		t.Fatalf("セーブ2に失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// タイムスタンプが更新されていることを確認
	if !loadedData.Timestamp.After(time1) {
		t.Error("Timestampが更新されていません")
	}
}

// TestResetSaveData はセーブデータのリセットをテストします。

func TestResetSaveData(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成
	saveData := NewSaveData()
	saveData.Statistics.TotalBattles = 100
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// リセット
	if err := io.ResetSaveData(); err != nil {
		t.Fatalf("リセットに失敗: %v", err)
	}

	// セーブファイルが存在しないことを確認
	_, err := io.LoadGame()
	if err == nil {
		t.Error("リセット後にセーブファイルが存在しています")
	}
}

// TestValidateSaveData はセーブデータのバリデーションをテストします。
func TestValidateSaveData(t *testing.T) {
	// 正常なデータ
	validData := NewSaveData()
	if err := ValidateSaveData(validData); err != nil {
		t.Errorf("正常なデータでエラー: %v", err)
	}

	// バージョンなし
	invalidData := NewSaveData()
	invalidData.Version = ""
	if err := ValidateSaveData(invalidData); err == nil {
		t.Error("Versionが空でもエラーにならない")
	}
}

// ==================== タスク3: 永続化層リファクタリングのテスト ====================

// TestCoreInstanceSaveWithoutID はIDフィールドを削除したCoreInstanceSaveをテストします。
// CoreInstanceSaveはcore_type_idとlevelのみを保持する。
func TestCoreInstanceSaveWithoutID(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成（新形式: IDなし）
	saveData := NewSaveData()
	saveData.Inventory.CoreInstances = append(saveData.Inventory.CoreInstances, CoreInstanceSave{
		CoreTypeID: "all_rounder",
		Level:      5,
	})

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if len(loadedData.Inventory.CoreInstances) != 1 {
		t.Fatalf("CoreInstances: got %d, want 1", len(loadedData.Inventory.CoreInstances))
	}
	core := loadedData.Inventory.CoreInstances[0]
	if core.CoreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", core.CoreTypeID)
	}
	if core.Level != 5 {
		t.Errorf("Level: got %d, want 5", core.Level)
	}
}

// TestModuleInstanceSaveWithChainEffect はチェイン効果付きModuleInstanceSaveをテストします。
func TestModuleInstanceSaveWithChainEffect(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成（新形式: ModuleInstances）
	saveData := NewSaveData()
	saveData.Inventory.ModuleInstances = append(saveData.Inventory.ModuleInstances, ModuleInstanceSave{
		TypeID: "physical_lv1",
		ChainEffect: &ChainEffectSave{
			Type:  "damage_bonus",
			Value: 15.0,
		},
	})
	// チェイン効果なしのモジュールも追加
	saveData.Inventory.ModuleInstances = append(saveData.Inventory.ModuleInstances, ModuleInstanceSave{
		TypeID:      "heal_lv1",
		ChainEffect: nil,
	})

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if len(loadedData.Inventory.ModuleInstances) != 2 {
		t.Fatalf("ModuleInstances: got %d, want 2", len(loadedData.Inventory.ModuleInstances))
	}

	// チェイン効果ありのモジュール
	mod1 := loadedData.Inventory.ModuleInstances[0]
	if mod1.TypeID != "physical_lv1" {
		t.Errorf("TypeID: got %s, want physical_lv1", mod1.TypeID)
	}
	if mod1.ChainEffect == nil {
		t.Fatal("ChainEffectがnilです")
	}
	if mod1.ChainEffect.Type != "damage_bonus" {
		t.Errorf("ChainEffect.Type: got %s, want damage_bonus", mod1.ChainEffect.Type)
	}
	if mod1.ChainEffect.Value != 15.0 {
		t.Errorf("ChainEffect.Value: got %f, want 15.0", mod1.ChainEffect.Value)
	}

	// チェイン効果なしのモジュール
	mod2 := loadedData.Inventory.ModuleInstances[1]
	if mod2.TypeID != "heal_lv1" {
		t.Errorf("TypeID: got %s, want heal_lv1", mod2.TypeID)
	}
	if mod2.ChainEffect != nil {
		t.Error("ChainEffectがnilであるべき")
	}
}

// TestAgentInstanceSaveWithChainEffects はチェイン効果付きAgentInstanceSaveをテストします。
func TestAgentInstanceSaveWithChainEffects(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成
	saveData := NewSaveData()
	saveData.Inventory.AgentInstances = append(saveData.Inventory.AgentInstances, AgentInstanceSave{
		ID: "agent_001",
		Core: CoreInstanceSave{
			CoreTypeID: "attack_balance",
			Level:      3,
		},
		Skills: []SkillInstanceSave{
			{TypeID: "physical_lv1", ChainEffect: &ChainEffectSave{Type: "damage_bonus", Value: 15.0}},
			{TypeID: "heal_lv1"}, // チェイン効果なし
			{TypeID: "buff_lv1", ChainEffect: &ChainEffectSave{Type: "buff_extend", Value: 2.0}},
			{TypeID: "debuff_lv1"}, // チェイン効果なし
		},
	})

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証
	if len(loadedData.Inventory.AgentInstances) != 1 {
		t.Fatalf("AgentInstances: got %d, want 1", len(loadedData.Inventory.AgentInstances))
	}

	agent := loadedData.Inventory.AgentInstances[0]
	if agent.ID != "agent_001" {
		t.Errorf("Agent ID: got %s, want agent_001", agent.ID)
	}
	if agent.Core.CoreTypeID != "attack_balance" {
		t.Errorf("Core.CoreTypeID: got %s, want attack_balance", agent.Core.CoreTypeID)
	}
	if agent.Core.Level != 3 {
		t.Errorf("Core.Level: got %d, want 3", agent.Core.Level)
	}
	if len(agent.Skills) != 4 {
		t.Fatalf("Skills count: got %d, want 4", len(agent.Skills))
	}

	// チェイン効果の検証
	if agent.Skills[0].ChainEffect == nil {
		t.Fatal("Skills[0].ChainEffectがnilです")
	}
	if agent.Skills[0].ChainEffect.Type != "damage_bonus" {
		t.Errorf("Skills[0].ChainEffect.Type: got %s, want damage_bonus", agent.Skills[0].ChainEffect.Type)
	}
	if agent.Skills[1].ChainEffect != nil {
		t.Error("Skills[1].ChainEffectはnilであるべき")
	}
	if agent.Skills[2].ChainEffect == nil {
		t.Fatal("Skills[2].ChainEffectがnilです")
	}
	if agent.Skills[2].ChainEffect.Type != "buff_extend" {
		t.Errorf("Skills[2].ChainEffect.Type: got %s, want buff_extend", agent.Skills[2].ChainEffect.Type)
	}
	if agent.Skills[3].ChainEffect != nil {
		t.Error("Skills[3].ChainEffectはnilであるべき")
	}
}

// TestSaveDataVersionV1 はv1.0.0形式のセーブデータをテストします。
func TestSaveDataVersionV1(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// セーブデータを作成
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

	// バージョンを検証
	if loadedData.Version != CurrentSaveDataVersion {
		t.Errorf("Version: got %s, want %s", loadedData.Version, CurrentSaveDataVersion)
	}
}

// TestModuleInstancesReplacesModuleCounts はModuleCountsがModuleInstancesに置き換わることをテストします。
func TestModuleInstancesReplacesModuleCounts(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 新形式のセーブデータを作成
	saveData := NewSaveData()
	// ModuleCountsは空のまま
	// ModuleInstancesに追加
	saveData.Inventory.ModuleInstances = append(saveData.Inventory.ModuleInstances,
		ModuleInstanceSave{TypeID: "physical_lv1", ChainEffect: &ChainEffectSave{Type: "damage_amp", Value: 20.0}},
		ModuleInstanceSave{TypeID: "physical_lv1", ChainEffect: &ChainEffectSave{Type: "life_steal", Value: 10.0}},
		ModuleInstanceSave{TypeID: "heal_lv1", ChainEffect: nil},
	)

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 検証: 同一TypeIDでも異なるChainEffectで別インスタンスとして保持される
	if len(loadedData.Inventory.ModuleInstances) != 3 {
		t.Errorf("ModuleInstances count: got %d, want 3", len(loadedData.Inventory.ModuleInstances))
	}

	// 同じTypeIDでも異なるChainEffectを持つことを確認
	physicalCount := 0
	for _, m := range loadedData.Inventory.ModuleInstances {
		if m.TypeID == "physical_lv1" {
			physicalCount++
		}
	}
	if physicalCount != 2 {
		t.Errorf("physical_lv1 count: got %d, want 2", physicalCount)
	}
}

// ==================== v4.0.0: 敵進行システム対応テスト ====================

// TestSaveDataVersionV4 はv4.0.0形式のセーブデータをテストします。
func TestSaveDataVersionV4(t *testing.T) {
	saveData := NewSaveData()

	// バージョンが4.0.0であることを確認
	if saveData.Version != "4.0.0" {
		t.Errorf("Version: got %s, want 4.0.0", saveData.Version)
	}
}

// TestPlayerSaveData_MaxHP はPlayerSaveDataにMaxHPフィールドがあることをテストします。
func TestPlayerSaveData_MaxHP(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()
	saveData.Player.MaxHP = 1500

	// セーブ
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// MaxHPが正しく保存・読み込みされることを確認
	if loadedData.Player.MaxHP != 1500 {
		t.Errorf("MaxHP: got %d, want 1500", loadedData.Player.MaxHP)
	}
}

// TestEnemyProgressSave は敵進行データの保存と読み込みをテストします。
func TestEnemyProgressSave(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()
	saveData.EnemyProgress = &EnemyProgressSave{
		CurrentRank: 2,
		DefeatRecords: map[string]DefeatRecordSave{
			"slime": {Defeated: true, MaxDefeatedLevel: 10},
			"bat":   {Defeated: true, MaxDefeatedLevel: 5},
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

	// EnemyProgressが正しく保存・読み込みされることを確認
	if loadedData.EnemyProgress == nil {
		t.Fatal("EnemyProgressがnilです")
	}
	if loadedData.EnemyProgress.CurrentRank != 2 {
		t.Errorf("CurrentRank: got %d, want 2", loadedData.EnemyProgress.CurrentRank)
	}
	if len(loadedData.EnemyProgress.DefeatRecords) != 2 {
		t.Errorf("DefeatRecords count: got %d, want 2", len(loadedData.EnemyProgress.DefeatRecords))
	}
	slimeRecord, exists := loadedData.EnemyProgress.DefeatRecords["slime"]
	if !exists {
		t.Fatal("slimeの記録が存在しません")
	}
	if !slimeRecord.Defeated {
		t.Error("slimeが撃破済みでない")
	}
	if slimeRecord.MaxDefeatedLevel != 10 {
		t.Errorf("slime MaxDefeatedLevel: got %d, want 10", slimeRecord.MaxDefeatedLevel)
	}
}

// TestAgentSlotSave_NoLevel はAgentSlotSaveからCoreLevelフィールドが削除されていることをテストします。
func TestAgentSlotSave_NoLevel(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()
	saveData.Player.AgentSlots[0] = AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]SkillSlotSaveCfg{
			{TypeID: "skill1"},
			{TypeID: "skill2"},
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

	// CoreTypeIDが保存されていることを確認
	if loadedData.Player.AgentSlots[0].CoreTypeID != "all_rounder" {
		t.Errorf("CoreTypeID: got %s, want all_rounder", loadedData.Player.AgentSlots[0].CoreTypeID)
	}
}

// TestCoreInventorySave_TypeIDOnly はCoreInventorySaveがTypeIDリスト形式であることをテストします。
func TestCoreInventorySave_TypeIDOnly(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	saveData := NewSaveData()
	saveData.Inventory.UniqueCores = &CoreInventorySave{
		Cores: []string{"all_rounder", "paladin", "healer"},
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

	// UniqueCoresが正しく保存・読み込みされることを確認
	if loadedData.Inventory.UniqueCores == nil {
		t.Fatal("UniqueCoresがnilです")
	}
	if len(loadedData.Inventory.UniqueCores.Cores) != 3 {
		t.Errorf("Cores count: got %d, want 3", len(loadedData.Inventory.UniqueCores.Cores))
	}
}

// TestNewSaveData_V4Defaults はNewSaveDataがv4.0.0の正しいデフォルト値を設定することをテストします。
func TestNewSaveData_V4Defaults(t *testing.T) {
	saveData := NewSaveData()

	// バージョン確認
	if saveData.Version != "4.0.0" {
		t.Errorf("Version: got %s, want 4.0.0", saveData.Version)
	}

	// PlayerのMaxHPデフォルト確認
	if saveData.Player.MaxHP != 1000 {
		t.Errorf("Player.MaxHP: got %d, want 1000", saveData.Player.MaxHP)
	}

	// EnemyProgressデフォルト確認
	if saveData.EnemyProgress == nil {
		t.Fatal("EnemyProgressがnilです")
	}
	if saveData.EnemyProgress.CurrentRank != 1 {
		t.Errorf("EnemyProgress.CurrentRank: got %d, want 1", saveData.EnemyProgress.CurrentRank)
	}
}

// ==================== v4.0.0: バージョンチェックテスト ====================

// TestValidateSaveVersion_ValidVersion は有効なバージョンの検証をテストします。
func TestValidateSaveVersion_ValidVersion(t *testing.T) {
	// v4.0.0は有効
	err := ValidateSaveVersion("4.0.0")
	if err != nil {
		t.Errorf("v4.0.0 should be valid, got error: %v", err)
	}
}

// TestValidateSaveVersion_OldVersion は旧バージョンがエラーになることをテストします。
func TestValidateSaveVersion_OldVersion(t *testing.T) {
	testCases := []struct {
		version string
	}{
		{"3.0.0"},
		{"2.0.0"},
		{"1.0.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.version, func(t *testing.T) {
			err := ValidateSaveVersion(tc.version)
			if err == nil {
				t.Errorf("Version %s should return error", tc.version)
			}
		})
	}
}

// TestValidateSaveVersion_EmptyVersion は空のバージョンがエラーになることをテストします。
func TestValidateSaveVersion_EmptyVersion(t *testing.T) {
	err := ValidateSaveVersion("")
	if err == nil {
		t.Error("Empty version should return error")
	}
}

// TestLoadGame_VersionCheck はLoadGame時のバージョンチェックをテストします。
func TestLoadGame_VersionCheck(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// 旧バージョンのセーブデータを作成（手動でJSONを書き込み）
	oldVersionSave := `{
		"version": "3.0.0",
		"timestamp": "2025-01-01T00:00:00Z",
		"player": {
			"equipped_agent_ids": ["", "", ""],
			"agent_slots": [{}, {}, {}]
		},
		"inventory": {
			"core_instances": [],
			"module_instances": [],
			"agent_instances": [],
			"max_core_slots": 100,
			"max_module_slots": 200,
			"max_agent_slots": 20
		},
		"statistics": {
			"total_battles": 0,
			"victories": 0,
			"defeats": 0,
			"max_level_reached": 0,
			"highest_wpm": 0,
			"average_wpm": 0,
			"perfect_accuracy_count": 0,
			"total_characters_typed": 0
		},
		"achievements": {
			"unlocked": [],
			"progress": {}
		},
		"settings": {
			"key_bindings": {}
		}
	}`

	savePath := filepath.Join(tmpDir, SaveFileName)
	if err := os.WriteFile(savePath, []byte(oldVersionSave), 0644); err != nil {
		t.Fatalf("セーブファイル作成に失敗: %v", err)
	}

	// LoadGameがエラーを返すことを確認
	_, err := io.LoadGame()
	if err == nil {
		t.Error("旧バージョンのセーブデータをロードするとエラーが返されるべき")
	}
}

// TestLoadGame_CurrentVersion は現在バージョンが正常にロードされることをテストします。
func TestLoadGame_CurrentVersion(t *testing.T) {
	tmpDir := t.TempDir()
	io := NewSaveDataIO(tmpDir, false)

	// v4.0.0のセーブデータを作成
	saveData := NewSaveData()
	if err := io.SaveGame(saveData); err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// LoadGameが成功することを確認
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// バージョンが正しいことを確認
	if loadedData.Version != CurrentSaveDataVersion {
		t.Errorf("Version: got %s, want %s", loadedData.Version, CurrentSaveDataVersion)
	}
}
