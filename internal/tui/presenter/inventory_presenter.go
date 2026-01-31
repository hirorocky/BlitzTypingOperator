package presenter

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/slot"
)

// InventoryProviderAdapter はAgentSlotManagerをInventoryProviderインターフェースに適合させます。
// v3.0.0以降はAgentSlotManagerを使用してエージェント構成を管理します。
type InventoryProviderAdapter struct {
	slotMgr *slot.AgentSlotManager
}

// NewInventoryProviderAdapter は新しいInventoryProviderAdapterを作成します。
// v3.0.0以降はAgentSlotManagerを使用します。
func NewInventoryProviderAdapter(slotMgr *slot.AgentSlotManager) *InventoryProviderAdapter {
	return &InventoryProviderAdapter{
		slotMgr: slotMgr,
	}
}

// GetCores はコア一覧を返します。
// v3.0.0では空のリストを返します（コアはスロットに設定済みのもののみ）。
func (a *InventoryProviderAdapter) GetCores() []*domain.CoreModel {
	// v3.0.0: コアインスタンスは保持しない（TypeID+Levelで管理）
	return []*domain.CoreModel{}
}

// GetModules はモジュール一覧を返します。
// v3.0.0では空のリストを返します（スキルはスロットに設定済みのもののみ）。
func (a *InventoryProviderAdapter) GetModules() []*domain.ModuleModel {
	// v3.0.0: スキルインスタンスは保持しない（TypeID+ChainEffectで管理）
	return []*domain.ModuleModel{}
}

// GetAgents はエージェント一覧を返します。
// v3.0.0ではBuildAgentsForBattleの結果を返します。
func (a *InventoryProviderAdapter) GetAgents() []*domain.AgentModel {
	return a.slotMgr.BuildAgentsForBattle()
}

// GetEquippedAgents は装備中のエージェント一覧を返します。
// v3.0.0ではBuildAgentsForBattleの結果を返します（全スロットが装備扱い）。
func (a *InventoryProviderAdapter) GetEquippedAgents() []*domain.AgentModel {
	return a.slotMgr.BuildAgentsForBattle()
}

// AddAgent はエージェントを追加します。
// v3.0.0では何もしません（エージェントはスロットに直接設定）。
func (a *InventoryProviderAdapter) AddAgent(_ *domain.AgentModel) error {
	// v3.0.0: エージェントはスロットに直接設定するため、この操作は無効
	return nil
}

// RemoveCore はコアをインベントリから削除します。
// v3.0.0では何もしません（コアはユニーク管理）。
func (a *InventoryProviderAdapter) RemoveCore(_ string) error {
	// v3.0.0: コアはユニーク管理のため、削除操作は無効
	return nil
}

// RemoveModule はモジュールをインベントリから削除します。
// v3.0.0では何もしません（スキルはユニーク管理）。
func (a *InventoryProviderAdapter) RemoveModule(_ string) error {
	// v3.0.0: スキルはユニーク管理のため、削除操作は無効
	return nil
}

// EquipAgent はエージェントを装備します。
// v3.0.0では何もしません（スロットへの設定はAgentSlotManager経由）。
func (a *InventoryProviderAdapter) EquipAgent(_ int, _ *domain.AgentModel) error {
	// v3.0.0: スロットへの設定はAgentSlotManager.SetCore/SetSkillで行う
	return nil
}

// UnequipAgent は装備を解除します。
// v3.0.0では何もしません（スロットのクリアはAgentSlotManager経由）。
func (a *InventoryProviderAdapter) UnequipAgent(_ int) error {
	// v3.0.0: スロットのクリアはAgentSlotManager.ClearCoreで行う
	return nil
}
