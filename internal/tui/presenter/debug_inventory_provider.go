package presenter

import (
	"fmt"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/usecase/slot"
)

// DebugInventoryProvider はデバッグモード用のInventoryProviderです。
// マスターデータから全CoreType/SkillType/ChainEffectを提供し、
// 任意のパラメータでコア・スキルを作成できます。
type DebugInventoryProvider struct {
	coreTypes     []masterdata.CoreTypeData
	skillTypes    []masterdata.SkillDefinitionData
	chainEffects  []masterdata.ChainEffectData
	passiveSkills map[string]domain.PassiveSkill

	// 作成されたエージェント（メモリ上で管理）
	agents         []*domain.AgentModel
	equippedAgents [3]*domain.AgentModel

	// スロットマネージャーへの参照（装備エージェント取得用）
	slotManager *slot.AgentSlotManager
}

// NewDebugInventoryProvider は新しいDebugInventoryProviderを作成します。
func NewDebugInventoryProvider(
	coreTypes []masterdata.CoreTypeData,
	skillTypes []masterdata.SkillDefinitionData,
	chainEffects []masterdata.ChainEffectData,
	passiveSkills map[string]domain.PassiveSkill,
) *DebugInventoryProvider {
	return &DebugInventoryProvider{
		coreTypes:     coreTypes,
		skillTypes:    skillTypes,
		chainEffects:  chainEffects,
		passiveSkills: passiveSkills,
		agents:        make([]*domain.AgentModel, 0),
	}
}

// SetSlotManager はスロットマネージャーを設定します。
// 装備エージェント取得でAgentSlotManagerを使用するために必要です。
func (p *DebugInventoryProvider) SetSlotManager(slotMgr *slot.AgentSlotManager) {
	p.slotManager = slotMgr
}

// ==================== InventoryProvider インターフェース実装 ====================

// GetCores はデバッグモードでは空のスライスを返します。
// デバッグモードではCoreType選択UIを使用するため。
func (p *DebugInventoryProvider) GetCores() []*domain.CoreModel {
	return nil
}

// GetSkills はデバッグモードでは空のスライスを返します。
// デバッグモードではSkillType選択UIを使用するため。
func (p *DebugInventoryProvider) GetSkills() []*domain.SkillModel {
	return nil
}

// GetAgents はエージェント一覧を返します。
func (p *DebugInventoryProvider) GetAgents() []*domain.AgentModel {
	return p.agents
}

// GetEquippedAgents は装備中のエージェント一覧を返します。
// slotManagerがある場合はそちらから取得します。
func (p *DebugInventoryProvider) GetEquippedAgents() []*domain.AgentModel {
	// slotManagerがある場合はそちらを優先
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

// RemoveSkill はデバッグモードでは何もしません（スキルは無限）。
func (p *DebugInventoryProvider) RemoveSkill(id string) error {
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

// GetSkillTypes はすべてのSkillTypeを返します（デバッグモード専用）。
func (p *DebugInventoryProvider) GetSkillTypes() []masterdata.SkillDefinitionData {
	return p.skillTypes
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

// CreateSkillFromType はSkillTypeとChainEffectからSkillModelを作成します。
func (p *DebugInventoryProvider) CreateSkillFromType(typeID string, chainEffect *domain.ChainEffect) *domain.SkillModel {
	for _, mt := range p.skillTypes {
		if mt.ID == typeID {
			skillType := mt.ToDomainType()
			return domain.NewSkillFromType(skillType, chainEffect)
		}
	}
	return nil
}
