// Package agent はエージェント管理機能を提供します。
// コア特性とモジュールの互換性検証、エージェント合成、装備管理を担当します。

package synthesize

import (
	"fmt"

	"hirorocky/type-battle/internal/domain"

	"github.com/google/uuid"
)

// MaxEquipmentSlots はエージェント装備スロットの最大数です。

const MaxEquipmentSlots = 3

// SynthesisPreview はエージェント合成プレビュー情報を表す構造体です。

type SynthesisPreview struct {
	// CoreName はコアの名前です。
	CoreName string

	// CoreTypeName はコア特性の名前です。
	CoreTypeName string

	// Level はエージェントのレベル（=コアレベル）です。
	Level int

	// BaseStats は基礎ステータスです。
	BaseStats domain.Stats

	// Modules は装備予定のモジュールリストです。
	Modules []*domain.ModuleModel

	// PassiveSkillName はパッシブスキルの名前です。
	PassiveSkillName string

	// PassiveSkillDesc はパッシブスキルの説明です。
	PassiveSkillDesc string
}

// AgentManager はエージェント管理を担当する構造体です。

type AgentManager struct {
	// coreInventory はコアインベントリです。
	coreInventory *domain.CoreInventory

	// moduleInventory はモジュールインベントリです。
	moduleInventory *domain.ModuleInventory

	// agentInventory はエージェントインベントリです。
	agentInventory *domain.AgentInventory

	// equippedAgents は装備中のエージェント（スロット番号 → エージェント）です。

	equippedAgents [MaxEquipmentSlots]*domain.AgentModel
}

// NewAgentManager は新しいAgentManagerを作成します。
// AgentInventoryは内部で作成・管理されます（最大20体）。
func NewAgentManager(
	coreInv *domain.CoreInventory,
	moduleInv *domain.ModuleInventory,
) *AgentManager {
	return &AgentManager{
		coreInventory:   coreInv,
		moduleInventory: moduleInv,
		agentInventory:  domain.NewAgentInventoryWithDefault(20),
		equippedAgents:  [MaxEquipmentSlots]*domain.AgentModel{},
	}
}

// InitializeWithDefaults は初期エージェントをセットアップします。
func (m *AgentManager) InitializeWithDefaults() {
	// 初期エージェントを作成（コアとモジュールが存在する場合のみ）
	cores := m.coreInventory.List()
	modules := m.moduleInventory.List()

	if len(cores) > 0 && len(modules) >= domain.ModuleSlotCount {
		// 最初のコアと最初の4つのモジュールで初期エージェントを合成
		core := cores[0]
		moduleIDs := make([]string, domain.ModuleSlotCount)
		for i := 0; i < domain.ModuleSlotCount; i++ {
			moduleIDs[i] = modules[i].ID
		}

		agent, err := m.SynthesizeAgent(core.ID, moduleIDs)
		if err == nil && agent != nil {
			// 初期エージェントを装備（スロット0）
			m.equippedAgents[0] = agent
		}
	}
}

// ==================== コア特性とモジュールタグ互換性検証（Task 5.1） ====================

// GetAllowedTags はコア特性の許可タグリストを取得します。

func (m *AgentManager) GetAllowedTags(core *domain.CoreModel) []string {
	return core.AllowedTags
}

// ValidateModuleCompatibility はモジュールがコアに装備可能かを判定します。

func (m *AgentManager) ValidateModuleCompatibility(core *domain.CoreModel, module *domain.ModuleModel) bool {
	return module.IsCompatibleWithCore(core)
}

// FilterCompatibleModules はコアに装備可能なモジュールのみをフィルタリングします。

func (m *AgentManager) FilterCompatibleModules(core *domain.CoreModel) []*domain.ModuleModel {
	if m.moduleInventory == nil {
		return nil
	}
	return m.moduleInventory.FilterCompatibleWithCore(core)
}

// ==================== エージェント合成機能（Task 5.2） ====================

// SynthesizeAgent はコアとモジュールからエージェントを合成します。

func (m *AgentManager) SynthesizeAgent(coreID string, moduleIDs []string) (*domain.AgentModel, error) {

	if len(moduleIDs) != domain.ModuleSlotCount {
		return nil, fmt.Errorf("モジュールが4個必要です（現在: %d個）", len(moduleIDs))
	}

	// コアを取得
	core := m.coreInventory.Get(coreID)
	if core == nil {
		return nil, fmt.Errorf("コアが見つかりません: %s", coreID)
	}

	// モジュールを取得し、互換性チェック
	modules := make([]*domain.ModuleModel, 0, domain.ModuleSlotCount)
	for _, moduleID := range moduleIDs {
		module := m.moduleInventory.Get(moduleID)
		if module == nil {
			return nil, fmt.Errorf("モジュールが見つかりません: %s", moduleID)
		}

		if !m.ValidateModuleCompatibility(core, module) {
			return nil, fmt.Errorf("モジュール '%s' はコア '%s' に装備できません", module.Name, core.Name)
		}

		modules = append(modules, module)
	}

	if m.agentInventory.IsFull() {
		return nil, fmt.Errorf("エージェントインベントリが満杯です")
	}

	// 新しいエージェントを作成
	agentID := uuid.New().String()
	agent := domain.NewAgent(agentID, core, modules)

	m.coreInventory.Remove(coreID)
	for _, moduleID := range moduleIDs {
		m.moduleInventory.Remove(moduleID)
	}

	// エージェントをインベントリに追加
	if err := m.agentInventory.Add(agent); err != nil {
		return nil, fmt.Errorf("エージェントの追加に失敗: %w", err)
	}

	return agent, nil
}

// GetSynthesisPreview は合成プレビュー情報を取得します。

func (m *AgentManager) GetSynthesisPreview(coreID string, moduleIDs []string) (*SynthesisPreview, error) {
	if len(moduleIDs) != domain.ModuleSlotCount {
		return nil, fmt.Errorf("モジュールが4個必要です（現在: %d個）", len(moduleIDs))
	}

	core := m.coreInventory.Get(coreID)
	if core == nil {
		return nil, fmt.Errorf("コアが見つかりません: %s", coreID)
	}

	modules := make([]*domain.ModuleModel, 0, domain.ModuleSlotCount)
	for _, moduleID := range moduleIDs {
		module := m.moduleInventory.Get(moduleID)
		if module == nil {
			return nil, fmt.Errorf("モジュールが見つかりません: %s", moduleID)
		}
		if !m.ValidateModuleCompatibility(core, module) {
			return nil, fmt.Errorf("モジュール '%s' はコア '%s' に装備できません", module.Name, core.Name)
		}
		modules = append(modules, module)
	}

	return &SynthesisPreview{
		CoreName:         core.Name,
		CoreTypeName:     core.Type.Name,
		Level:            core.Level,
		BaseStats:        core.Stats,
		Modules:          modules,
		PassiveSkillName: core.PassiveSkill.Name,
		PassiveSkillDesc: core.PassiveSkill.Description,
	}, nil
}

// ==================== エージェント装備機能（Task 5.3） ====================

// EquipAgent はエージェントを指定スロットに装備します。

func (m *AgentManager) EquipAgent(slot int, agentID string, player *domain.PlayerModel) error {
	// スロット番号チェック
	if slot < 0 || slot >= MaxEquipmentSlots {
		return fmt.Errorf("無効なスロット番号です: %d（有効範囲: 0-%d）", slot, MaxEquipmentSlots-1)
	}

	// エージェントを取得
	agent := m.agentInventory.Get(agentID)
	if agent == nil {
		return fmt.Errorf("エージェントが見つかりません: %s", agentID)
	}

	// 既に他のスロットに装備されている場合はエラー
	for i := 0; i < MaxEquipmentSlots; i++ {
		if m.equippedAgents[i] != nil && m.equippedAgents[i].ID == agentID && i != slot {
			return fmt.Errorf("エージェント '%s' は既にスロット %d に装備されています", agent.ID, i)
		}
	}

	// 装備
	m.equippedAgents[slot] = agent

	m.recalculatePlayerHP(player)

	return nil
}

// UnequipAgent はエージェントを装備解除します。

func (m *AgentManager) UnequipAgent(slot int, player *domain.PlayerModel) error {
	if slot < 0 || slot >= MaxEquipmentSlots {
		return fmt.Errorf("無効なスロット番号です: %d", slot)
	}

	m.equippedAgents[slot] = nil

	m.recalculatePlayerHP(player)

	return nil
}

// GetEquippedAgents は装備中のエージェントリストを返します。

func (m *AgentManager) GetEquippedAgents() []*domain.AgentModel {
	result := make([]*domain.AgentModel, 0, MaxEquipmentSlots)
	for _, agent := range m.equippedAgents {
		if agent != nil {
			result = append(result, agent)
		}
	}
	return result
}

// GetEquippedAgentAt は指定スロットの装備エージェントを返します。
func (m *AgentManager) GetEquippedAgentAt(slot int) *domain.AgentModel {
	if slot < 0 || slot >= MaxEquipmentSlots {
		return nil
	}
	return m.equippedAgents[slot]
}

// GetAgents は所持エージェント一覧を返します。
func (m *AgentManager) GetAgents() []*domain.AgentModel {
	return m.agentInventory.List()
}

// AddAgent はエージェントをインベントリに追加します。
func (m *AgentManager) AddAgent(agent *domain.AgentModel) error {
	return m.agentInventory.Add(agent)
}

// GetEquippedCount は装備中のエージェント数を返します。
func (m *AgentManager) GetEquippedCount() int {
	count := 0
	for _, agent := range m.equippedAgents {
		if agent != nil {
			count++
		}
	}
	return count
}

// HasEquippedAgent はエージェントが1体以上装備されているかを返します。

func (m *AgentManager) HasEquippedAgent() bool {
	return m.GetEquippedCount() > 0
}

// recalculatePlayerHP は装備エージェントに基づいてプレイヤーHPを再計算します。

func (m *AgentManager) recalculatePlayerHP(player *domain.PlayerModel) {
	agents := m.GetEquippedAgents()
	player.RecalculateHP(agents)
}

// GetAgentDetails はエージェントの詳細情報を取得します。

func (m *AgentManager) GetAgentDetails(agentID string) *domain.AgentModel {
	return m.agentInventory.Get(agentID)
}
