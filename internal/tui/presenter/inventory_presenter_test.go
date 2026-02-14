package presenter

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/slot"
)

// TestInventoryProviderAdapter はInventoryProviderAdapterの基本動作をテストします。
func TestInventoryProviderAdapter(t *testing.T) {
	// 新システムのAgentSlotManagerを使用
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	slotMgr := slot.NewAgentSlotManager(
		coreInv,
		skillInv,
		map[string]domain.CoreType{},
		map[string]domain.SkillType{},
		map[string]domain.PassiveSkill{},
		nil, // chainEffects
	)

	adapter := NewInventoryProviderAdapter(slotMgr)

	if adapter == nil {
		t.Fatal("NewInventoryProviderAdapter returned nil")
	}

	// コア取得（空スライス）
	cores := adapter.GetCores()
	if cores == nil {
		t.Error("GetCores returned nil")
	}

	// スキル取得（空スライス）
	skills := adapter.GetSkills()
	if skills == nil {
		t.Error("GetSkills returned nil")
	}

	// エージェント取得（スロットが空なので空スライス）
	agents := adapter.GetAgents()
	if agents == nil {
		t.Error("GetAgents returned nil")
	}
}

// TestInventoryProviderAdapter_WithData はデータがある場合のテストです。
func TestInventoryProviderAdapter_WithData(t *testing.T) {
	// コアとスキルをインベントリに追加
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// コアを追加
	coreInv.AddCore("warrior")

	// スキルを追加
	skillInv.AddSkill("slash")

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
		nil, // chainEffects
	)

	// スロットにコアを設定
	_ = slotMgr.SetCore(0, "warrior")
	_ = slotMgr.SetSkill(0, 0, "slash")

	adapter := NewInventoryProviderAdapter(slotMgr)

	// GetCores/GetSkillsは空スライスを返す
	cores := adapter.GetCores()
	if len(cores) != 0 {
		t.Errorf("Expected 0 cores, got %d", len(cores))
	}

	skills := adapter.GetSkills()
	if len(skills) != 0 {
		t.Errorf("Expected 0 skills, got %d", len(skills))
	}

	// エージェントはBuildAgentsForBattleで構築される
	agents := adapter.GetAgents()
	if len(agents) != 1 {
		t.Errorf("Expected 1 agent, got %d", len(agents))
	}
}

// TestInventoryProviderAdapter_AddAgent はエージェント追加をテストします。
// 、この操作は無効化されています（スロットへの直接設定を使用）。
func TestInventoryProviderAdapter_AddAgent(t *testing.T) {
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	slotMgr := slot.NewAgentSlotManager(
		coreInv,
		skillInv,
		map[string]domain.CoreType{},
		map[string]domain.SkillType{},
		map[string]domain.PassiveSkill{},
		nil, // chainEffects
	)

	adapter := NewInventoryProviderAdapter(slotMgr)

	initialCount := len(adapter.GetAgents())

	// AddAgentは何もしない
	coreType := domain.CoreType{
		ID:          "test_type",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
	}
	core := domain.NewCoreWithTypeID("test_core", coreType, domain.PassiveSkill{})
	agent := domain.NewAgent("test_agent", core, nil)

	err := adapter.AddAgent(agent)
	if err != nil {
		t.Errorf("AddAgent failed: %v", err)
	}

	// AddAgentは無効化されているので、カウントは変わらない
	if len(adapter.GetAgents()) != initialCount {
		t.Errorf("Expected %d agents, got %d", initialCount, len(adapter.GetAgents()))
	}
}
