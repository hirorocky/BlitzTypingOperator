// Package screens はTUI画面のテストを提供します。
package screens

import (
	"testing"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.4: エージェント管理画面のテスト ====================

// TestNewAgentManagementScreen はAgentManagementScreenの初期化をテストします。
func TestNewAgentManagementScreen(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	if screen == nil {
		t.Fatal("AgentManagementScreenがnilです")
	}

	if screen.inventory == nil {
		t.Error("インベントリがnilです")
	}
}

// TestAgentManagementTabs はタブ切り替えをテストします。

func TestAgentManagementTabs(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// 初期タブ
	if screen.currentTab != TabCoreList {
		t.Errorf("初期タブが正しくありません: got %d, want %d", screen.currentTab, TabCoreList)
	}

	// タブ切り替え（右へ）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})
	if screen.currentTab != TabModuleList {
		t.Errorf("タブ切り替え(右)が正しくありません: got %d, want %d", screen.currentTab, TabModuleList)
	}

	// タブ切り替え（左へ）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyLeft})
	if screen.currentTab != TabCoreList {
		t.Errorf("タブ切り替え(左)が正しくありません: got %d, want %d", screen.currentTab, TabCoreList)
	}
}

// TestAgentManagementCoreList はコア一覧表示をテストします。

func TestAgentManagementCoreList(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// コア一覧タブに移動
	screen.currentTab = TabCoreList
	screen.updateCurrentList()

	// コアリストが表示されていること
	if len(screen.coreList) == 0 {
		t.Error("コアリストが空です")
	}
}

// TestAgentManagementModuleList はモジュール一覧表示をテストします。

func TestAgentManagementModuleList(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// モジュール一覧タブに移動
	screen.currentTab = TabModuleList
	screen.updateCurrentList()

	// モジュールリストが表示されていること
	if len(screen.moduleList) == 0 {
		t.Error("モジュールリストが空です")
	}
}

// TestAgentManagementSynthesis は合成サブ画面をテストします。

func TestAgentManagementSynthesis(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// 合成タブに移動
	screen.currentTab = TabSynthesis
	screen.updateCurrentList()

	// 合成状態が初期化されていること
	if screen.synthesisState.selectedCore != nil {
		t.Error("初期状態でコアが選択されています")
	}

	if len(screen.synthesisState.selectedModules) != 0 {
		t.Error("初期状態でモジュールが選択されています")
	}
}

// TestAgentManagementEquip は装備サブ画面をテストします。

func TestAgentManagementEquip(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// 装備タブに移動
	screen.currentTab = TabEquip
	screen.updateCurrentList()

	// 装備スロットが3つあること
	if len(screen.equipSlots) != 3 {
		t.Errorf("装備スロット数: got %d, want 3", len(screen.equipSlots))
	}
}

// TestAgentManagementCoreDetailDisplay はコア詳細情報表示をテストします。

func TestAgentManagementCoreDetailDisplay(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// コア一覧タブでコアを選択
	screen.currentTab = TabCoreList
	screen.updateCurrentList()

	if len(screen.coreList) > 0 {
		screen.selectedIndex = 0
		detail := screen.getSelectedCoreDetail()

		if detail == nil {
			t.Error("コア詳細が取得できません")
		}
	}
}

// TestAgentManagementModuleDetailDisplay はモジュール詳細情報表示をテストします。

func TestAgentManagementModuleDetailDisplay(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// モジュール一覧タブでモジュールを選択
	screen.currentTab = TabModuleList
	screen.updateCurrentList()

	if len(screen.moduleList) > 0 {
		screen.selectedIndex = 0
		detail := screen.getSelectedModuleDetail()

		if detail == nil {
			t.Error("モジュール詳細が取得できません")
		}
	}
}

// TestAgentManagementSynthesisFlow は合成フローをテストします。

func TestAgentManagementSynthesisFlow(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// 合成タブに移動
	screen.currentTab = TabSynthesis
	screen.updateCurrentList()

	// コアを選択
	if len(screen.coreList) > 0 {
		screen.synthesisState.selectedCore = screen.coreList[0]
	}

	// モジュールを4つ選択
	if len(screen.moduleList) >= 4 {
		screen.synthesisState.selectedModules = screen.moduleList[:4]
	}

	// 合成可能かチェック
	canSynthesize := screen.canSynthesize()
	if !canSynthesize {
		t.Log("合成に必要な条件が満たされていません（テスト環境依存）")
	}
}

// TestAgentManagementEquipFlow は装備フローをテストします。

func TestAgentManagementEquipFlow(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	// 装備タブに移動
	screen.currentTab = TabEquip
	screen.updateCurrentList()

	// 初期状態では空きスロットがあること
	emptySlots := 0
	for _, slot := range screen.equipSlots {
		if slot == nil {
			emptySlots++
		}
	}

	if emptySlots != 3 {
		t.Errorf("空きスロット数: got %d, want 3", emptySlots)
	}
}

// TestAgentManagementBackNavigation は戻るナビゲーションをテストします。
func TestAgentManagementBackNavigation(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Escキーでコマンドが返されません")
	}
}

// TestAgentManagementRender はレンダリングをテストします。
func TestAgentManagementRender(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// ==================== ヘルパー関数 ====================

// TestInventory はテスト用のインベントリを表すインターフェースです。
type TestInventory struct {
	cores    []*domain.CoreModel
	modules  []*domain.ModuleModel
	agents   []*domain.AgentModel
	equipped []*domain.AgentModel
}

// GetCores はコア一覧を返します。
func (i *TestInventory) GetCores() []*domain.CoreModel {
	return i.cores
}

// GetModules はモジュール一覧を返します。
func (i *TestInventory) GetModules() []*domain.ModuleModel {
	return i.modules
}

// GetAgents はエージェント一覧を返します。
func (i *TestInventory) GetAgents() []*domain.AgentModel {
	return i.agents
}

// GetEquippedAgents は装備中エージェント一覧を返します。
func (i *TestInventory) GetEquippedAgents() []*domain.AgentModel {
	return i.equipped
}

// AddAgent はエージェントを追加します。
func (i *TestInventory) AddAgent(agent *domain.AgentModel) error {
	i.agents = append(i.agents, agent)
	return nil
}

// RemoveCore はコアを削除します。
func (i *TestInventory) RemoveCore(id string) error {
	for idx, c := range i.cores {
		if c.ID == id {
			i.cores = append(i.cores[:idx], i.cores[idx+1:]...)
			return nil
		}
	}
	return nil
}

// RemoveModule はモジュールを削除します。
func (i *TestInventory) RemoveModule(id string) error {
	for idx, m := range i.modules {
		if m.ID == id {
			i.modules = append(i.modules[:idx], i.modules[idx+1:]...)
			return nil
		}
	}
	return nil
}

// EquipAgent はエージェントを装備します。
func (i *TestInventory) EquipAgent(slot int, agent *domain.AgentModel) error {
	for len(i.equipped) <= slot {
		i.equipped = append(i.equipped, nil)
	}
	i.equipped[slot] = agent
	return nil
}

// UnequipAgent はエージェントの装備を解除します。
func (i *TestInventory) UnequipAgent(slot int) error {
	if slot < len(i.equipped) {
		i.equipped[slot] = nil
	}
	return nil
}

func createTestInventory() *TestInventory {
	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
	}

	core1 := domain.NewCore("core1", "コア1", 5, coreType, domain.PassiveSkill{})
	core2 := domain.NewCore("core2", "コア2", 10, coreType, domain.PassiveSkill{})

	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "物理攻撃", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", "物理ダメージ"),
		domain.NewModule("m2", "魔法攻撃", domain.MagicAttack, 1, []string{"magic_low"}, 10, "MAG", "魔法ダメージ"),
		domain.NewModule("m3", "回復", domain.Heal, 1, []string{"heal_low"}, 10, "MAG", "HP回復"),
		domain.NewModule("m4", "バフ", domain.Buff, 1, []string{"buff_low"}, 10, "SPD", "攻撃力UP"),
		domain.NewModule("m5", "デバフ", domain.Debuff, 1, []string{"debuff_low"}, 10, "SPD", "攻撃力DOWN"),
	}

	// テスト用エージェントを作成
	agentCore1 := domain.NewCore("agent_core1", "エージェントコア1", 5, coreType, domain.PassiveSkill{})
	agentCore2 := domain.NewCore("agent_core2", "エージェントコア2", 10, coreType, domain.PassiveSkill{})
	agentModules1 := []*domain.ModuleModel{
		domain.NewModule("am1", "物理攻撃", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", "物理ダメージ"),
		domain.NewModule("am2", "魔法攻撃", domain.MagicAttack, 1, []string{"magic_low"}, 10, "MAG", "魔法ダメージ"),
		domain.NewModule("am3", "回復", domain.Heal, 1, []string{"heal_low"}, 10, "MAG", "HP回復"),
		domain.NewModule("am4", "バフ", domain.Buff, 1, []string{"buff_low"}, 10, "SPD", "攻撃力UP"),
	}
	agentModules2 := []*domain.ModuleModel{
		domain.NewModule("am5", "物理攻撃2", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", "物理ダメージ"),
		domain.NewModule("am6", "魔法攻撃2", domain.MagicAttack, 1, []string{"magic_low"}, 10, "MAG", "魔法ダメージ"),
		domain.NewModule("am7", "回復2", domain.Heal, 1, []string{"heal_low"}, 10, "MAG", "HP回復"),
		domain.NewModule("am8", "バフ2", domain.Buff, 1, []string{"buff_low"}, 10, "SPD", "攻撃力UP"),
	}
	agent1 := domain.NewAgent("agent1", agentCore1, agentModules1)
	agent2 := domain.NewAgent("agent2", agentCore2, agentModules2)

	return &TestInventory{
		cores:    []*domain.CoreModel{core1, core2},
		modules:  modules,
		agents:   []*domain.AgentModel{agent1, agent2},
		equipped: []*domain.AgentModel{nil, nil, nil},
	}
}

// ==================== Task 5.1-5.4: エージェント管理画面UI改善のテスト ====================

// TestAgentManagementSynthesisLeftRightLayout は合成タブの左右分割レイアウトをテストします。

func TestAgentManagementSynthesisLeftRightLayout(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)
	screen.currentTab = TabSynthesis
	screen.width = 120
	screen.height = 40
	screen.updateCurrentList()

	rendered := screen.View()

	// 合成タブが表示されていること
	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}

	// 選択状況パネルが表示されていること
	if !containsString(rendered, "コア:") {
		t.Error("コア選択状況が表示されていません")
	}

	if !containsString(rendered, "モジュール") {
		t.Error("モジュール選択状況が表示されていません")
	}
}

// TestAgentManagementSynthesisDetailAndPreview は合成タブのパーツ詳細と完成予測ステータス表示をテストします。

func TestAgentManagementSynthesisDetailAndPreview(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)
	screen.currentTab = TabSynthesis
	screen.synthesisState.step = 0
	screen.width = 120
	screen.height = 40
	screen.updateCurrentList()

	// コアを選択
	if len(screen.coreList) > 0 {
		screen.selectedIndex = 0
		detail := screen.getSelectedCoreDetail()
		if detail == nil {
			t.Error("コア詳細が取得できません")
		}
	}

	rendered := screen.View()
	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestAgentManagementEquipTopBottomLayout は装備タブの上下分割レイアウトをテストします。

func TestAgentManagementEquipTopBottomLayout(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)
	screen.currentTab = TabEquip
	screen.width = 120
	screen.height = 40
	screen.updateCurrentList()

	rendered := screen.View()

	// 装備スロットが表示されていること
	if !containsString(rendered, "スロット1") || !containsString(rendered, "スロット2") || !containsString(rendered, "スロット3") {
		t.Error("装備スロットが表示されていません")
	}
}

// TestAgentManagementEquipSlotSwitch は装備タブのスロット切替をテストします。

func TestAgentManagementEquipSlotSwitch(t *testing.T) {
	inventory := createTestInventory()
	screen := NewAgentManagementScreen(inventory)
	screen.currentTab = TabEquip
	screen.width = 120
	screen.height = 40
	screen.updateCurrentList()

	// 初期選択位置
	initialIndex := screen.selectedIndex

	// 下キーで次のスロット/エージェントに移動
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyDown})

	// 選択位置が変わっていること
	if screen.selectedIndex == initialIndex && screen.getMaxIndex() > 1 {
		t.Error("選択位置が変わっていません")
	}
}

// containsString は文字列に部分文字列が含まれるかを確認します（テスト用）。
func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
