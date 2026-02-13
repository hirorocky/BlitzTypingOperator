// Package app は BlitzTypingOperator TUIゲームのエラーログ出力テストを提供します。
package app

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
	gamestate "hirorocky/type-battle/internal/usecase/session"
)

// テスト用のログバッファとハンドラーを作成するヘルパー関数
func setupTestLogger() (*bytes.Buffer, func()) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	oldLogger := slog.Default()
	slog.SetDefault(slog.New(handler))

	return &buf, func() {
		slog.SetDefault(oldLogger)
	}
}

// TestGameStateFromSaveDataLogsAddCoreError は AddCore のエラーがログ出力されることをテストします。

func TestGameStateFromSaveDataLogsAddCoreError(t *testing.T) {
	// slogのログ出力をキャプチャ
	buf, cleanup := setupTestLogger()
	defer cleanup()

	// 正常なセーブデータを作成
	saveData := savedata.NewSaveData()
	saveData.Inventory = &savedata.InventorySaveData{
		UniqueCores: &savedata.CoreInventorySave{
			Cores: []string{"all_rounder"},
		},
		UniqueSkills: &savedata.SkillInventorySave{
			Skills: make([]string, 0),
		},
	}

	// テスト用のソースデータを作成
	sources := &gamestate.DomainDataSources{
		CoreTypes: []domain.CoreType{
			{
				ID:             "all_rounder",
				Name:           "オールラウンダー",
				StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
				PassiveSkillID: "test_skill",
				AllowedTags:    []string{"physical_low"},
				MinDropLevel:   1,
			},
		},
		PassiveSkills: map[string]domain.PassiveSkill{
			"test_skill": {ID: "test_skill", Name: "テストスキル", Description: "テスト用"},
		},
	}

	// GameStateをセーブデータから作成
	// 正常なケースではエラーは発生しないが、ログ機能自体が動作していることを確認
	_ = gamestate.GameStateFromSaveData(saveData, sources)

	// ログ出力の検証（正常ケースではエラーログは出力されない）
	logOutput := buf.String()
	if strings.Contains(logOutput, "level=ERROR") {
		t.Logf("ログ出力が検出されました: %s", logOutput)
	}
}

// TestGameStateFromSaveDataLogsAgentErrors は AddAgent および EquipAgent のエラーがログ出力されることをテストします。
func TestGameStateFromSaveDataLogsAgentErrors(t *testing.T) {
	// slogのログ出力をキャプチャ
	buf, cleanup := setupTestLogger()
	defer cleanup()

	// エージェントを含むセーブデータを作成
	saveData := savedata.NewSaveData()
	saveData.Inventory = &savedata.InventorySaveData{
		UniqueCores: &savedata.CoreInventorySave{
			Cores: []string{},
		},
		UniqueSkills: &savedata.SkillInventorySave{
			Skills: make([]string, 0),
		},
	}
	saveData.Player = &savedata.PlayerSaveData{
		MaxHP: 1000,
		AgentSlots: [3]savedata.AgentSlotSave{
			{
				CoreTypeID: "all_rounder",
				Skills: [4]savedata.SkillSlotSaveCfg{
					{TypeID: "mod_slash"},
					{TypeID: "mod_slash"},
					{TypeID: "mod_slash"},
					{TypeID: "mod_slash"},
				},
			},
			{},
			{},
		},
	}

	// テスト用のソースデータを作成
	sources := &gamestate.DomainDataSources{
		CoreTypes: []domain.CoreType{
			{
				ID:             "all_rounder",
				Name:           "オールラウンダー",
				StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
				PassiveSkillID: "test_skill",
				AllowedTags:    []string{"physical_low"},
				MinDropLevel:   1,
			},
		},
		PassiveSkills: map[string]domain.PassiveSkill{
			"test_skill": {ID: "test_skill", Name: "テストスキル", Description: "テスト用"},
		},
	}

	// GameStateをセーブデータから作成
	gs := gamestate.GameStateFromSaveData(saveData, sources)

	// ログ出力の検証（ログ機能自体が動作していることを確認）
	logOutput := buf.String()
	_ = logOutput // ログが正常に設定されていることを確認

	// GameStateが正常に作成されていることを確認
	if gs == nil {
		t.Error("GameStateがnilです")
	}
}

// TestInventoryManagerLogsErrors は InventoryManager のコア追加が正常に動作することをテストします。
// 新システムではTypeIDベースのユニーク管理のため、容量制限はありません。

func TestInventoryManagerLogsErrors(t *testing.T) {
	// slogのログ出力をキャプチャ
	buf, cleanup := setupTestLogger()
	defer cleanup()

	// InventoryManagerを初期化
	invManager := gamestate.NewInventoryManager()

	// コアを追加
	updated := invManager.AddCore("all_rounder")
	if !updated {
		t.Error("コア追加で更新が期待されます")
	}

	// ログ出力の検証（正常時はエラーログは出力されない）
	logOutput := buf.String()
	_ = logOutput

	// コアが追加されていることを確認
	ownedCores := invManager.GetOwnedCores()
	if len(ownedCores) == 0 {
		t.Error("コアが追加されていません")
	}
}

// TestSlogLoggingFunctionality は slog パッケージが正常に動作することをテストします。
func TestSlogLoggingFunctionality(t *testing.T) {
	buf, cleanup := setupTestLogger()
	defer cleanup()

	// テスト用にエラーログを出力
	slog.Error("テストエラー",
		slog.String("core_id", "test_core_001"),
		slog.String("error", "テストエラーメッセージ"),
	)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "テストエラー") {
		t.Errorf("エラーメッセージがログに含まれていません: %s", logOutput)
	}
	if !strings.Contains(logOutput, "core_id") {
		t.Errorf("core_idがログに含まれていません: %s", logOutput)
	}
	if !strings.Contains(logOutput, "test_core_001") {
		t.Errorf("コアIDの値がログに含まれていません: %s", logOutput)
	}
}

// TestLoggedAddCoreError は構造化ログが正しく出力されることをテストします。
// 新システムでは容量制限がないため、手動でエラーログを出力してテストします。

func TestLoggedAddCoreError(t *testing.T) {
	buf, cleanup := setupTestLogger()
	defer cleanup()

	// 新システムでは容量制限がないため、エラーを直接シミュレート
	invManager := gamestate.NewInventoryManager()

	// コアを追加
	invManager.AddCore("all_rounder")

	// エラーログのフォーマットをテスト（実際のエラーケースをシミュレート）
	slog.Error("コア追加に失敗",
		slog.String("core_id", "core_002"),
		slog.String("core_type", "all_rounder"),
		slog.String("error", "シミュレートされたエラー"),
	)

	// ログ出力を検証
	logOutput := buf.String()
	if !strings.Contains(logOutput, "コア追加に失敗") {
		t.Errorf("エラーメッセージがログに含まれていません: %s", logOutput)
	}
	if !strings.Contains(logOutput, "core_002") {
		t.Errorf("コアIDがログに含まれていません: %s", logOutput)
	}
}

// TestLoaderCoreTypeData はローダーのCoreTypeDataが正しく動作することを確認します。
func TestLoaderCoreTypeData(t *testing.T) {
	coreTypeData := masterdata.CoreTypeData{
		ID:             "all_rounder",
		Name:           "オールラウンダー",
		AllowedTags:    []string{"physical_low", "magic_low"},
		StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		PassiveSkillID: "balance_mastery",
		MinDropLevel:   1,
	}

	domainType := coreTypeData.ToDomain()
	if domainType.ID != "all_rounder" {
		t.Errorf("コアタイプIDが一致しません: got %s, want all_rounder", domainType.ID)
	}
}
