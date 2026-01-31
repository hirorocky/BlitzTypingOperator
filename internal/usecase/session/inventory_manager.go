package session

import (
	"hirorocky/type-battle/internal/domain"
)

// InventoryManager はゲーム全体のインベントリを統合管理する構造体です。
// コアとスキルの管理を担当します。
// エージェントの管理はAgentSlotManagerが一元的に行います。
type InventoryManager struct {
	// cores はコアインベントリです（TypeIDごとの最大レベル管理）。
	cores *domain.CoreInventory

	// skills はスキルインベントリです（TypeIDごとの保有状態とチェイン効果管理）。
	skills *domain.SkillInventory
}

// NewInventoryManager は新しいInventoryManagerを作成します。
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		cores:  domain.NewCoreInventory(),
		skills: domain.NewSkillInventory(),
	}
}

// Cores はコアインベントリを返します。
func (m *InventoryManager) Cores() *domain.CoreInventory {
	return m.cores
}

// Skills はスキルインベントリを返します。
func (m *InventoryManager) Skills() *domain.SkillInventory {
	return m.skills
}

// AddCore はコアをインベントリに追加します。
// 同一TypeIDの場合、より高いレベルのみ保存されます。
func (m *InventoryManager) AddCore(typeID string, level int) bool {
	return m.cores.AddCore(typeID, level)
}

// AddSkill はスキルをインベントリに追加します。
func (m *InventoryManager) AddSkill(typeID string, chainEffectID string) {
	m.skills.AddSkill(typeID, chainEffectID)
}

// GetOwnedCores は保有しているコアのTypeIDとレベルのマップを返します。
func (m *InventoryManager) GetOwnedCores() map[string]int {
	return m.cores.GetOwnedCores()
}

// GetOwnedSkills は保有しているスキルの情報を返します。
func (m *InventoryManager) GetOwnedSkills() map[string]*domain.SkillOwnership {
	return m.skills.GetOwnedSkills()
}

// HasCore は指定TypeIDのコアを保有しているかを返します。
func (m *InventoryManager) HasCore(typeID string) bool {
	return m.cores.HasCore(typeID)
}

// HasSkill は指定TypeIDのスキルを保有しているかを返します。
func (m *InventoryManager) HasSkill(typeID string) bool {
	return m.skills.HasSkill(typeID)
}

// GetCoreMaxLevel は指定TypeIDのコアの最大レベルを返します。
func (m *InventoryManager) GetCoreMaxLevel(typeID string) int {
	return m.cores.GetMaxLevel(typeID)
}
