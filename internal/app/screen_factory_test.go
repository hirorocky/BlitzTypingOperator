// Package app は BlitzTypingOperator TUIゲームの画面生成テストを提供します。
package app

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
)

// mockInventoryProvider はテスト用のInventoryProviderモック
type mockInventoryProvider struct{}

func (m *mockInventoryProvider) GetCores() []*domain.CoreModel {
	return []*domain.CoreModel{}
}

func (m *mockInventoryProvider) GetModules() []*domain.ModuleModel {
	return []*domain.ModuleModel{}
}

func (m *mockInventoryProvider) GetAgents() []*domain.AgentModel {
	return []*domain.AgentModel{}
}

func (m *mockInventoryProvider) GetEquippedAgents() []*domain.AgentModel {
	return []*domain.AgentModel{}
}

func (m *mockInventoryProvider) AddAgent(agent *domain.AgentModel) error {
	return nil
}

func (m *mockInventoryProvider) RemoveCore(id string) error {
	return nil
}

func (m *mockInventoryProvider) RemoveModule(id string) error {
	return nil
}

func (m *mockInventoryProvider) EquipAgent(slot int, agent *domain.AgentModel) error {
	return nil
}

func (m *mockInventoryProvider) UnequipAgent(slot int) error {
	return nil
}

// TestNewScreenFactory は新しいScreenFactoryが正しく初期化されることを検証します
func TestNewScreenFactory(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())
	if factory == nil {
		t.Fatal("NewScreenFactory() returned nil")
	}
}

// TestScreenFactory_CreateHomeScreen はホーム画面の生成を検証します
func TestScreenFactory_CreateHomeScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	mockProvider := &mockInventoryProvider{}
	screen := factory.CreateHomeScreen(1, mockProvider)
	if screen == nil {
		t.Fatal("CreateHomeScreen() returned nil")
	}
}

// TestScreenFactory_CreateBattleSelectScreen はバトル選択画面の生成を検証します
func TestScreenFactory_CreateBattleSelectScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	mockProvider := &mockInventoryProvider{}
	screen := factory.CreateBattleSelectScreen(1, mockProvider)
	if screen == nil {
		t.Fatal("CreateBattleSelectScreen() returned nil")
	}
}

// TestScreenFactory_CreateAgentManagementScreen はエージェント管理画面の生成を検証します
func TestScreenFactory_CreateAgentManagementScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	mockProvider := &mockInventoryProvider{}
	screen := factory.CreateAgentManagementScreen(mockProvider)
	if screen == nil {
		t.Fatal("CreateAgentManagementScreen() returned nil")
	}
}

// TestScreenFactory_CreateEncyclopediaScreen は図鑑画面の生成を検証します
func TestScreenFactory_CreateEncyclopediaScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	screen := factory.CreateEncyclopediaScreen()
	if screen == nil {
		t.Fatal("CreateEncyclopediaScreen() returned nil")
	}
}

// TestScreenFactory_CreateStatsAchievementsScreen は統計・実績画面の生成を検証します
func TestScreenFactory_CreateStatsAchievementsScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	screen := factory.CreateStatsAchievementsScreen()
	if screen == nil {
		t.Fatal("CreateStatsAchievementsScreen() returned nil")
	}
}

// TestScreenFactory_CreateSettingsScreen は設定画面の生成を検証します
func TestScreenFactory_CreateSettingsScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	screen := factory.CreateSettingsScreen()
	if screen == nil {
		t.Fatal("CreateSettingsScreen() returned nil")
	}
}

// TestScreenFactory_GameStateReference はGameState参照を保持することを検証します
func TestScreenFactory_GameStateReference(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	factory := NewScreenFactory(model.GameState())

	if factory.gameState == nil {
		t.Fatal("ScreenFactory should hold GameState reference")
	}
}
