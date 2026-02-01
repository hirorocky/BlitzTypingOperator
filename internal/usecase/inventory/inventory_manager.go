// Package inventory はインベントリ管理のユースケースを提供します。
// このファイルはInventoryManager（改良版）を定義します。

package inventory

import (
	"hirorocky/type-battle/internal/domain"
)

// InventoryManager はコアとスキルのユニークインベントリを統合管理します。
// CoreInventoryとSkillInventoryを統合し、報酬追加時のインベントリ更新を担当します。
type InventoryManager struct {
	// cores はコアインベントリ
	cores *domain.CoreInventory

	// skills はスキルインベントリ
	skills *domain.SkillInventory
}

// NewInventoryManager は新しいInventoryManagerを作成します。
// 空のCoreInventoryとSkillInventoryで初期化されます。
func NewInventoryManager() *InventoryManager {
	return &InventoryManager{
		cores:  domain.NewCoreInventory(),
		skills: domain.NewSkillInventory(),
	}
}

// NewInventoryManagerWithInventories は既存のインベントリを使用してInventoryManagerを作成します。
// nilが渡された場合は新規インベントリを作成します。
// セーブデータからの復元時に使用します。
func NewInventoryManagerWithInventories(
	cores *domain.CoreInventory,
	skills *domain.SkillInventory,
) *InventoryManager {
	if cores == nil {
		cores = domain.NewCoreInventory()
	}
	if skills == nil {
		skills = domain.NewSkillInventory()
	}

	return &InventoryManager{
		cores:  cores,
		skills: skills,
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

// AddCore はコアを追加します。
// 新規TypeIDの場合にtrueを返し、既に保有している場合はfalseを返します。
func (m *InventoryManager) AddCore(typeID string) bool {
	return m.cores.AddCore(typeID)
}

// AddSkill はスキルを追加します。
// chainEffectID: チェイン効果ID（なしの場合は空文字列）
func (m *InventoryManager) AddSkill(typeID string, chainEffectID string) {
	m.skills.AddSkill(typeID, chainEffectID)
}

// GetOwnedCoreTypes は保有コアTypeID一覧を返します。
func (m *InventoryManager) GetOwnedCoreTypes() []string {
	return m.cores.GetOwnedCores()
}

// GetOwnedSkillTypes は保有スキルTypeID一覧を返します。
func (m *InventoryManager) GetOwnedSkillTypes() []string {
	ownedSkills := m.skills.GetOwnedSkills()
	result := make([]string, 0, len(ownedSkills))
	for typeID := range ownedSkills {
		result = append(result, typeID)
	}
	return result
}
