package presenter

import (
	"fmt"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/usecase/slot"
)

// DebugInventoryProvider はデバッグモード用のInventoryProviderです。
// マスターデータから全CoreType/ModuleType/ChainEffectを提供し、
// 任意のパラメータでコア・モジュールを作成できます。
type DebugInventoryProvider struct {
	coreTypes     []masterdata.CoreTypeData
	moduleTypes   []masterdata.ModuleDefinitionData
	chainEffects  []masterdata.ChainEffectData
	passiveSkills map[string]domain.PassiveSkill

	// 作成されたエージェント（メモリ上で管理）
	agents         []*domain.AgentModel
	equippedAgents [3]*domain.AgentModel

	// v3.0.0: スロットマネージャーへの参照（装備エージェント取得用）
	slotManager *slot.AgentSlotManager
}

// NewDebugInventoryProvider は新しいDebugInventoryProviderを作成します。
func NewDebugInventoryProvider(
	coreTypes []masterdata.CoreTypeData,
	moduleTypes []masterdata.ModuleDefinitionData,
	chainEffects []masterdata.ChainEffectData,
	passiveSkills map[string]domain.PassiveSkill,
) *DebugInventoryProvider {
	return &DebugInventoryProvider{
		coreTypes:     coreTypes,
		moduleTypes:   moduleTypes,
		chainEffects:  chainEffects,
		passiveSkills: passiveSkills,
		agents:        make([]*domain.AgentModel, 0),
	}
}

// SetSlotManager はスロットマネージャーを設定します。
// v3.0.0: 装備エージェント取得でAgentSlotManagerを使用するために必要です。
func (p *DebugInventoryProvider) SetSlotManager(slotMgr *slot.AgentSlotManager) {
	p.slotManager = slotMgr
}

// ==================== InventoryProvider インターフェース実装 ====================

// GetCores はデバッグモードでは空のスライスを返します。
// デバッグモードではCoreType選択UIを使用するため。
func (p *DebugInventoryProvider) GetCores() []*domain.CoreModel {
	return nil
}

// GetModules はデバッグモードでは空のスライスを返します。
// デバッグモードではModuleType選択UIを使用するため。
func (p *DebugInventoryProvider) GetModules() []*domain.ModuleModel {
	return nil
}

// GetAgents はエージェント一覧を返します。
func (p *DebugInventoryProvider) GetAgents() []*domain.AgentModel {
	return p.agents
}

// GetEquippedAgents は装備中のエージェント一覧を返します。
// v3.0.0: slotManagerがある場合はそちらから取得します。
func (p *DebugInventoryProvider) GetEquippedAgents() []*domain.AgentModel {
	// v3.0.0: slotManagerがある場合はそちらを優先
	if p.slotManager != nil {
		return p.slotManager.BuildAgentsForBattle()
	}

	// 後方互換: 旧来のequippedAgents配列から取得
	result := make([]*domain.AgentModel, 0, 3)
	for _, agent := range p.equippedAgents {
		if agent != nil {
			result = append(result, agent)
		}
	}
	return result
}

// AddAgent はエージェントを追加します。
func (p *DebugInventoryProvider) AddAgent(agent *domain.AgentModel) error {
	p.agents = append(p.agents, agent)
	return nil
}

// RemoveCore はデバッグモードでは何もしません（コアは無限）。
func (p *DebugInventoryProvider) RemoveCore(id string) error {
	return nil
}

// RemoveModule はデバッグモードでは何もしません（モジュールは無限）。
func (p *DebugInventoryProvider) RemoveModule(id string) error {
	return nil
}

// EquipAgent はエージェントを装備します。
func (p *DebugInventoryProvider) EquipAgent(slot int, agent *domain.AgentModel) error {
	if slot < 0 || slot >= 3 {
		return fmt.Errorf("無効なスロット番号: %d", slot)
	}
	p.equippedAgents[slot] = agent
	return nil
}

// UnequipAgent は装備を解除します。
func (p *DebugInventoryProvider) UnequipAgent(slot int) error {
	if slot < 0 || slot >= 3 {
		return fmt.Errorf("無効なスロット番号: %d", slot)
	}
	p.equippedAgents[slot] = nil
	return nil
}

// ==================== デバッグモード専用メソッド ====================

// GetCoreTypes はすべてのCoreTypeを返します（デバッグモード専用）。
func (p *DebugInventoryProvider) GetCoreTypes() []masterdata.CoreTypeData {
	return p.coreTypes
}

// GetModuleTypes はすべてのModuleTypeを返します（デバッグモード専用）。
func (p *DebugInventoryProvider) GetModuleTypes() []masterdata.ModuleDefinitionData {
	return p.moduleTypes
}

// GetChainEffects はすべてのChainEffectを返します（デバッグモード専用）。
func (p *DebugInventoryProvider) GetChainEffects() []masterdata.ChainEffectData {
	return p.chainEffects
}

// CreateCoreFromType はCoreTypeからCoreModelを作成します。
func (p *DebugInventoryProvider) CreateCoreFromType(typeID string) *domain.CoreModel {
	for _, ct := range p.coreTypes {
		if ct.ID == typeID {
			coreType := ct.ToDomain()
			passiveSkill := p.passiveSkills[ct.PassiveSkillID]
			return domain.NewCoreWithTypeID(typeID, coreType, passiveSkill)
		}
	}
	return nil
}

// CreateModuleFromType はModuleTypeとChainEffectからModuleModelを作成します。
func (p *DebugInventoryProvider) CreateModuleFromType(typeID string, chainEffect *domain.ChainEffect) *domain.ModuleModel {
	for _, mt := range p.moduleTypes {
		if mt.ID == typeID {
			moduleType := mt.ToDomainType()
			return domain.NewModuleFromType(moduleType, chainEffect)
		}
	}
	return nil
}
