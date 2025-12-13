package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/session"
)

// TestInventoryProviderAdapter はInventoryProviderAdapterの基本動作をテストします。
func TestInventoryProviderAdapter(t *testing.T) {
	gs := session.NewGameState()

	adapter := NewInventoryProviderAdapter(
		gs.Inventory(),
		gs.AgentManager(),
		gs.Player(),
	)

	if adapter == nil {
		t.Fatal("NewInventoryProviderAdapter returned nil")
	}

	// コア取得
	cores := adapter.GetCores()
	if cores == nil {
		t.Error("GetCores returned nil")
	}

	// モジュール取得
	modules := adapter.GetModules()
	if modules == nil {
		t.Error("GetModules returned nil")
	}

	// エージェント取得
	agents := adapter.GetAgents()
	if agents == nil {
		t.Error("GetAgents returned nil")
	}
}

// TestInventoryProviderAdapter_WithData はデータがある場合のテストです。
func TestInventoryProviderAdapter_WithData(t *testing.T) {
	gs := session.NewGameState()

	adapter := NewInventoryProviderAdapter(
		gs.Inventory(),
		gs.AgentManager(),
		gs.Player(),
	)

	// 初期データが含まれていることを確認（少なくとも1つ以上）
	cores := adapter.GetCores()
	if len(cores) < 1 {
		t.Errorf("Expected at least 1 core, got %d", len(cores))
	}

	modules := adapter.GetModules()
	if len(modules) < 1 {
		t.Errorf("Expected at least 1 module, got %d", len(modules))
	}
}

// TestInventoryProviderAdapter_AddAgent はエージェント追加をテストします。
func TestInventoryProviderAdapter_AddAgent(t *testing.T) {
	gs := session.NewGameState()

	adapter := NewInventoryProviderAdapter(
		gs.Inventory(),
		gs.AgentManager(),
		gs.Player(),
	)

	initialCount := len(adapter.GetAgents())

	// 新しいエージェントを作成して追加
	coreType := domain.CoreType{
		ID:          "test_type",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
	}
	core := domain.NewCore("test_core", "テストコア", 1, coreType, domain.PassiveSkill{})
	agent := domain.NewAgent("test_agent", core, nil)

	err := adapter.AddAgent(agent)
	if err != nil {
		t.Errorf("AddAgent failed: %v", err)
	}

	if len(adapter.GetAgents()) != initialCount+1 {
		t.Errorf("Expected %d agents, got %d", initialCount+1, len(adapter.GetAgents()))
	}
}
