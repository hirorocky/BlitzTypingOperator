package presenter

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/session"
	"hirorocky/type-battle/internal/usecase/synthesize"
)

// InventoryProviderAdapter はInventoryManagerとAgentManagerをInventoryProviderインターフェースに適合させます。
// コア・モジュールの管理はInventoryManager、エージェント・装備の管理はAgentManagerが担当します。
type InventoryProviderAdapter struct {
	inv      *session.InventoryManager
	agentMgr *synthesize.AgentManager
	player   *domain.PlayerModel
}

// NewInventoryProviderAdapter は新しいInventoryProviderAdapterを作成します。
func NewInventoryProviderAdapter(
	inv *session.InventoryManager,
	agentMgr *synthesize.AgentManager,
	player *domain.PlayerModel,
) *InventoryProviderAdapter {
	return &InventoryProviderAdapter{
		inv:      inv,
		agentMgr: agentMgr,
		player:   player,
	}
}

// GetCores はコア一覧を返します。
func (a *InventoryProviderAdapter) GetCores() []*domain.CoreModel {
	return a.inv.GetCores()
}

// GetModules はモジュール一覧を返します。
func (a *InventoryProviderAdapter) GetModules() []*domain.ModuleModel {
	return a.inv.GetModules()
}

// GetAgents はエージェント一覧を返します。
func (a *InventoryProviderAdapter) GetAgents() []*domain.AgentModel {
	return a.agentMgr.GetAgents()
}

// GetEquippedAgents は装備中のエージェント一覧を返します。
func (a *InventoryProviderAdapter) GetEquippedAgents() []*domain.AgentModel {
	return a.agentMgr.GetEquippedAgents()
}

// AddAgent はエージェントを追加します。
func (a *InventoryProviderAdapter) AddAgent(ag *domain.AgentModel) error {
	return a.agentMgr.AddAgent(ag)
}

// RemoveCore はコアをインベントリから削除します。
func (a *InventoryProviderAdapter) RemoveCore(id string) error {
	return a.inv.RemoveCore(id)
}

// RemoveModule はモジュールをインベントリから削除します。
func (a *InventoryProviderAdapter) RemoveModule(id string) error {
	return a.inv.RemoveModule(id)
}

// EquipAgent はエージェントを装備します。
func (a *InventoryProviderAdapter) EquipAgent(slot int, agentModel *domain.AgentModel) error {
	return a.agentMgr.EquipAgent(slot, agentModel.ID, a.player)
}

// UnequipAgent は装備を解除します。
func (a *InventoryProviderAdapter) UnequipAgent(slot int) error {
	return a.agentMgr.UnequipAgent(slot, a.player)
}
