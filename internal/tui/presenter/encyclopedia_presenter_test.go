package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/usecase/session"
)

// TestCreateEncyclopediaData は図鑑データ作成をテストします。
func TestCreateEncyclopediaData(t *testing.T) {
	gs := session.NewGameStateForTest()

	data := CreateEncyclopediaData(gs)

	if data == nil {
		t.Fatal("CreateEncyclopediaData returned nil")
	}

	// 基本的なデータ構造が存在すること（内容はマスタデータ依存）
	// AllCoreTypes, AllSkillTypes, AllEnemyTypes が nil でないこと
	if data.AllCoreTypes == nil {
		t.Error("AllCoreTypes is nil")
	}

	if data.AllSkillTypes == nil {
		t.Error("AllSkillTypes is nil")
	}

	if data.AllEnemyTypes == nil {
		t.Error("AllEnemyTypes is nil")
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
