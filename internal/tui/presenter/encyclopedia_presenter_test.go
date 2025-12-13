package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/usecase/session"
)

// TestCreateEncyclopediaData は図鑑データ作成をテストします。
func TestCreateEncyclopediaData(t *testing.T) {
	gs := session.NewGameState()

	data := CreateEncyclopediaData(gs)

	if data == nil {
		t.Fatal("CreateEncyclopediaData returned nil")
	}

	// コアタイプが含まれていること
	if len(data.AllCoreTypes) == 0 {
		t.Error("AllCoreTypes is empty")
	}

	// モジュールタイプが含まれていること
	if len(data.AllModuleTypes) == 0 {
		t.Error("AllModuleTypes is empty")
	}

	// 敵タイプが含まれていること
	if len(data.AllEnemyTypes) == 0 {
		t.Error("AllEnemyTypes is empty")
	}
}

// TestCreateDefaultEncyclopediaData はデフォルト図鑑データ作成をテストします。
func TestCreateDefaultEncyclopediaData(t *testing.T) {
	data := CreateDefaultEncyclopediaData()

	if data == nil {
		t.Fatal("CreateDefaultEncyclopediaData returned nil")
	}

	// 全コアタイプが4つ以上
	if len(data.AllCoreTypes) < 4 {
		t.Errorf("Expected at least 4 core types, got %d", len(data.AllCoreTypes))
	}
}
