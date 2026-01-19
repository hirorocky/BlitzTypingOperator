package session

import (
	"hirorocky/type-battle/internal/domain"
)

// InventoryManager はゲーム全体のインベントリを統合管理する構造体です。
// コアとモジュールの管理を担当します。
// エージェントの管理はAgentManagerが一元的に行います。
type InventoryManager struct {
	// cores はコアインベントリです。
	cores *domain.CoreInventoryLegacy

	// modules はモジュールインベントリです。
	modules *domain.ModuleInventory
}

// NewInventoryManager は新しいInventoryManagerを作成します。
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		cores:   domain.NewCoreInventoryLegacy(100),
		modules: domain.NewModuleInventory(200),
	}
}

// Cores はコアインベントリを返します。
func (m *InventoryManager) Cores() *domain.CoreInventoryLegacy {
	return m.cores
}

// Modules はモジュールインベントリを返します。
func (m *InventoryManager) Modules() *domain.ModuleInventory {
	return m.modules
}

// AddCore はコアをインベントリに追加します。
func (m *InventoryManager) AddCore(core *domain.CoreModel) error {
	return m.cores.Add(core)
}

// AddModule はモジュールをインベントリに追加します。
func (m *InventoryManager) AddModule(module *domain.ModuleModel) error {
	return m.modules.Add(module)
}

// GetCores はコア一覧を返します。
func (m *InventoryManager) GetCores() []*domain.CoreModel {
	return m.cores.List()
}

// GetModules はモジュール一覧を返します。
func (m *InventoryManager) GetModules() []*domain.ModuleModel {
	return m.modules.List()
}

// RemoveCore はコアをインベントリから削除します。
func (m *InventoryManager) RemoveCore(id string) error {
	m.cores.Remove(id)
	return nil
}

// RemoveModule はモジュールをインベントリから削除します。
// 指定されたTypeIDを持つ最初のモジュールを削除します。
func (m *InventoryManager) RemoveModule(typeID string) error {
	m.modules.RemoveByTypeID(typeID)
	return nil
}

// SetMaxCoreSlots はコアの最大スロット数を設定します。
func (m *InventoryManager) SetMaxCoreSlots(slots int) {
	m.cores = domain.NewCoreInventoryLegacy(slots)
}

// SetMaxModuleSlots はモジュールの最大スロット数を設定します。
func (m *InventoryManager) SetMaxModuleSlots(slots int) {
	m.modules = domain.NewModuleInventory(slots)
}
