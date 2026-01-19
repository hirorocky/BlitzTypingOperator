package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/slot"
)

// TestInventoryProviderAdapter はInventoryProviderAdapterの基本動作をテストします。
func TestInventoryProviderAdapter(t *testing.T) {
	// v3.0.0: 新システムのAgentSlotManagerを使用
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	slotMgr := slot.NewAgentSlotManager(
		coreInv,
		skillInv,
		map[string]domain.CoreType{},
		map[string]domain.SkillType{},
		map[string]domain.PassiveSkill{},
	)

	adapter := NewInventoryProviderAdapter(slotMgr)

	if adapter == nil {
		t.Fatal("NewInventoryProviderAdapter returned nil")
	}

	// コア取得（v3.0.0では空スライス）
	cores := adapter.GetCores()
	if cores == nil {
		t.Error("GetCores returned nil")
	}

	// モジュール取得（v3.0.0では空スライス）
	modules := adapter.GetModules()
	if modules == nil {
		t.Error("GetModules returned nil")
	}

	// エージェント取得（スロットが空なので空スライス）
	agents := adapter.GetAgents()
	if agents == nil {
		t.Error("GetAgents returned nil")
	}
}

// TestInventoryProviderAdapter_WithData はデータがある場合のテストです。
func TestInventoryProviderAdapter_WithData(t *testing.T) {
	// v3.0.0: コアとスキルをインベントリに追加
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// コアを追加
	coreInv.AddCore("warrior", 5)

	// スキルを追加
	skillInv.AddSkill("slash", "")

	// マスタデータを設定
	coreTypes := map[string]domain.CoreType{
		"warrior": {
			ID:             "warrior",
			Name:           "ウォリアー",
			AllowedTags:    []string{"attack"},
			PassiveSkillID: "",
		},
	}
	skillTypes := map[string]domain.SkillType{
		"slash": {
			ID:   "slash",
			Name: "斬撃",
			Tags: []string{"attack"},
		},
	}

	slotMgr := slot.NewAgentSlotManager(
		coreInv,
		skillInv,
		coreTypes,
		skillTypes,
		map[string]domain.PassiveSkill{},
	)

	// スロットにコアを設定
	_ = slotMgr.SetCore(0, "warrior", 5)
	_ = slotMgr.SetSkill(0, 0, "slash", "")

	adapter := NewInventoryProviderAdapter(slotMgr)

	// v3.0.0: GetCores/GetModulesは空スライスを返す
	cores := adapter.GetCores()
	if len(cores) != 0 {
		t.Errorf("Expected 0 cores, got %d", len(cores))
	}

	modules := adapter.GetModules()
	if len(modules) != 0 {
		t.Errorf("Expected 0 modules, got %d", len(modules))
	}

	// エージェントはBuildAgentsForBattleで構築される
	agents := adapter.GetAgents()
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}
}

// TestInventoryProviderAdapter_AddAgent はエージェント追加をテストします。
// v3.0.0では、この操作は無効化されています（スロットへの直接設定を使用）。
func TestInventoryProviderAdapter_AddAgent(t *testing.T) {
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	slotMgr := slot.NewAgentSlotManager(
		coreInv,
		skillInv,
		map[string]domain.CoreType{},
		map[string]domain.SkillType{},
		map[string]domain.PassiveSkill{},
	)

	adapter := NewInventoryProviderAdapter(slotMgr)

	initialCount := len(adapter.GetAgents())

	// v3.0.0: AddAgentは何もしない
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

	// v3.0.0: AddAgentは無効化されているので、カウントは変わらない
	if len(adapter.GetAgents()) != initialCount {
		t.Errorf("Expected %d agents, got %d", initialCount, len(adapter.GetAgents()))
	}
}
