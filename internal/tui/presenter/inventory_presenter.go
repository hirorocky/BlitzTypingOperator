package presenter

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/slot"
)

// InventoryProviderAdapter はAgentSlotManagerをInventoryProviderインターフェースに適合させます。
type InventoryProviderAdapter struct {
	slotMgr *slot.AgentSlotManager
}

// NewInventoryProviderAdapter は新しいInventoryProviderAdapterを作成します。
func NewInventoryProviderAdapter(slotMgr *slot.AgentSlotManager) *InventoryProviderAdapter {
	return &InventoryProviderAdapter{
		slotMgr: slotMgr,
	}
}

// GetCores はコア一覧を返します。
// 空のリストを返します（コアはスロットに設定済みのもののみ）。
func (a *InventoryProviderAdapter) GetCores() []*domain.CoreModel {
	return []*domain.CoreModel{}
}

// GetSkills はスキル一覧を返します。
// 空のリストを返します（スキルはスロットに設定済みのもののみ）。
func (a *InventoryProviderAdapter) GetSkills() []*domain.SkillModel {
	return []*domain.SkillModel{}
}

// GetAgents はエージェント一覧を返します。
func (a *InventoryProviderAdapter) GetAgents() []*domain.AgentModel {
	return a.slotMgr.BuildAgentsForBattle()
}

// GetEquippedAgents は装備中のエージェント一覧を返します。
func (a *InventoryProviderAdapter) GetEquippedAgents() []*domain.AgentModel {
	return a.slotMgr.BuildAgentsForBattle()
}

// AddAgent はエージェントを追加します。
// スロットシステムでは使用されません。
func (a *InventoryProviderAdapter) AddAgent(_ *domain.AgentModel) error {
	return nil
}

// RemoveCore はコアをインベントリから削除します。
// ユニーク管理のため使用されません。
func (a *InventoryProviderAdapter) RemoveCore(_ string) error {
	return nil
}

// RemoveSkill はスキルをインベントリから削除します。
// ユニーク管理のため使用されません。
func (a *InventoryProviderAdapter) RemoveSkill(_ string) error {
	return nil
}

// EquipAgent はエージェントを装備します。
// スロットシステムでは使用されません。
func (a *InventoryProviderAdapter) EquipAgent(_ int, _ *domain.AgentModel) error {
	return nil
}

// UnequipAgent は装備を解除します。
// スロットシステムでは使用されません。
func (a *InventoryProviderAdapter) UnequipAgent(_ int) error {
	return nil
}
