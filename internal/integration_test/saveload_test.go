// Package integration_test は統合テストを提供します。

package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/infra/startup"
)

// ==================================================
// Task 15.3: セーブ/ロードフロー統合テスト
// ==================================================

func TestSaveLoadFlow_WriteAndRead(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// 初期データを作成
	initializer := startup.NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()

	// 統計データを追加
	saveData.Statistics.TotalBattles = 10
	saveData.Statistics.Victories = 7
	saveData.Statistics.MaxLevelReached = 5
	saveData.Statistics.HighestWPM = 120.5

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 整合性確認
	if loadedData.Statistics.TotalBattles != 10 {
		t.Errorf("TotalBattles expected 10, got %d", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.Victories != 7 {
		t.Errorf("Victories expected 7, got %d", loadedData.Statistics.Victories)
	}
	if loadedData.Statistics.MaxLevelReached != 5 {
		t.Errorf("MaxLevelReached expected 5, got %d", loadedData.Statistics.MaxLevelReached)
	}
	if loadedData.Statistics.HighestWPM != 120.5 {
		t.Errorf("HighestWPM expected 120.5, got %f", loadedData.Statistics.HighestWPM)
	}
}

func TestSaveLoadFlow_InventoryPersistence(t *testing.T) {
	// インベントリの永続化テスト（ID化された構造）
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// テスト用データを作成
	saveData := savedata.NewSaveData()

	// コアインスタンスを追加（ID化された構造）
	saveData.Inventory.CoreInstances = append(saveData.Inventory.CoreInstances, savedata.CoreInstanceSave{
		ID:         "core_1",
		CoreTypeID: "test_type",
		Level:      5,
	})

	// モジュールカウントを追加（ID化された構造）
	saveData.Inventory.ModuleCounts["module_1"] = 1

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// コア確認（ID化された構造）
	if len(loadedData.Inventory.CoreInstances) != 1 {
		t.Fatalf("コアインスタンス数 expected 1, got %d", len(loadedData.Inventory.CoreInstances))
	}
	loadedCore := loadedData.Inventory.CoreInstances[0]
	if loadedCore.ID != "core_1" {
		t.Errorf("Core ID expected 'core_1', got '%s'", loadedCore.ID)
	}
	if loadedCore.Level != 5 {
		t.Errorf("Core Level expected 5, got %d", loadedCore.Level)
	}

	// モジュール確認（ID化された構造）
	if loadedData.Inventory.ModuleCounts["module_1"] != 1 {
		t.Errorf("モジュールカウント expected 1, got %d", loadedData.Inventory.ModuleCounts["module_1"])
	}
}

func TestSaveLoadFlow_CorruptedData_BackupRestore(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// 正常なデータをセーブ
	initializer := startup.NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()
	saveData.Statistics.MaxLevelReached = 10

	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// バックアップを確認
	bak1Path := filepath.Join(tempDir, "save.json.bak1")
	_, _ = os.Stat(bak1Path)
	// 初回セーブ時はバックアップがないかもしれない

	// もう一度セーブしてバックアップを作成
	saveData.Statistics.MaxLevelReached = 20
	err = io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("2回目のセーブに失敗: %v", err)
	}

	// メインファイルを破損させる
	savePath := filepath.Join(tempDir, "save.json")
	err = os.WriteFile(savePath, []byte("invalid json"), 0644)
	if err != nil {
		t.Fatalf("ファイル破損シミュレートに失敗: %v", err)
	}

	// ロードを試みる（バックアップから復元されるべき）
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("バックアップからの復元に失敗: %v", err)
	}

	// バックアップのデータが復元されている
	if loadedData.Statistics.MaxLevelReached != 10 {
		t.Errorf("バックアップのMaxLevelReached expected 10, got %d", loadedData.Statistics.MaxLevelReached)
	}
}

func TestSaveLoadFlow_BackupRotation(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	initializer := startup.NewNewGameInitializer(createTestExternalData())

	// 複数回セーブしてバックアップローテーションを確認
	for i := 1; i <= 5; i++ {
		saveData := initializer.InitializeNewGame()
		saveData.Statistics.MaxLevelReached = i
		err := io.SaveGame(saveData)
		if err != nil {
			t.Fatalf("セーブ%dに失敗: %v", i, err)
		}
	}

	// バックアップファイルの存在確認
	bak1 := filepath.Join(tempDir, "save.json.bak1")
	bak2 := filepath.Join(tempDir, "save.json.bak2")
	bak3 := filepath.Join(tempDir, "save.json.bak3")

	if _, err := os.Stat(bak1); os.IsNotExist(err) {
		t.Error("bak1が存在するべきです")
	}
	if _, err := os.Stat(bak2); os.IsNotExist(err) {
		t.Error("bak2が存在するべきです")
	}
	if _, err := os.Stat(bak3); os.IsNotExist(err) {
		t.Error("bak3が存在するべきです")
	}

	// bak4は存在しない（最大3世代）
	bak4 := filepath.Join(tempDir, "save.json.bak4")
	if _, err := os.Stat(bak4); !os.IsNotExist(err) {
		t.Error("bak4は存在しないべきです")
	}
}

func TestSaveLoadFlow_NewGameWhenNoSave(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// セーブファイルが存在しないことを確認
	if io.Exists() {
		t.Error("セーブファイルは存在しないべきです")
	}

	// ロード試行（失敗するはず）
	_, err := io.LoadGame()
	if err == nil {
		t.Error("セーブファイル不在時はエラーが返されるべきです")
	}

	// 新規ゲームを開始
	initializer := startup.NewNewGameInitializer(createTestExternalData())
	newGameData := initializer.InitializeNewGame()

	// セーブ
	err = io.SaveGame(newGameData)
	if err != nil {
		t.Fatalf("新規ゲームのセーブに失敗: %v", err)
	}

	// ロード成功
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 初期状態の確認
	if loadedData.Statistics.MaxLevelReached != 0 {
		t.Error("新規ゲームのMaxLevelReachedは0であるべきです")
	}
}

func TestSaveLoadFlow_DataValidation(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// バージョンが空のデータを作成
	savePath := filepath.Join(tempDir, "save.json")
	invalidData := `{"version": "", "player": null}`
	err := os.MkdirAll(tempDir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(savePath, []byte(invalidData), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// ロード試行（検証エラー）
	_, err = io.LoadGame()
	if err == nil {
		t.Error("不正なデータでは検証エラーが返されるべきです")
	}
}

func TestSaveLoadFlow_ResetSaveData(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)

	// データをセーブ
	initializer := startup.NewNewGameInitializer(createTestExternalData())
	saveData := initializer.InitializeNewGame()
	saveData.Statistics.MaxLevelReached = 50
	io.SaveGame(saveData)

	// セーブファイルが存在する
	if !io.Exists() {
		t.Error("セーブファイルが存在するべきです")
	}

	// リセット
	err := io.ResetSaveData()
	if err != nil {
		t.Fatalf("リセットに失敗: %v", err)
	}

	// セーブファイルが削除されている
	if io.Exists() {
		t.Error("リセット後はセーブファイルが存在しないべきです")
	}
}
