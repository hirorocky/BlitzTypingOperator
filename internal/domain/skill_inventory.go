// Package domain はゲームのドメインモデルを定義します。
// このファイルはSkillInventoryを定義します。

package domain

// SkillOwnership はスキルの保有情報を表す構造体です。
type SkillOwnership struct {
	// Owned はスキルを保有しているかどうかを示します。
	Owned bool
}

// NewSkillOwnership は新しいSkillOwnershipを作成します。
func NewSkillOwnership() *SkillOwnership {
	return &SkillOwnership{
		Owned: false,
	}
}

// SkillInventory はスキルをTypeIDごとにユニーク管理する構造体です。
type SkillInventory struct {
	// skills はSkillTypeID → 保有情報のマップ
	skills map[string]*SkillOwnership
}

// NewSkillInventory は新しいSkillInventoryを作成します。
func NewSkillInventory() *SkillInventory {
	return &SkillInventory{
		skills: make(map[string]*SkillOwnership),
	}
}

// AddSkill はスキルを追加します。
// TypeIDが空の場合は何もしません。
func (inv *SkillInventory) AddSkill(typeID string) {
	if typeID == "" {
		return
	}

	ownership, exists := inv.skills[typeID]
	if !exists {
		ownership = NewSkillOwnership()
		inv.skills[typeID] = ownership
	}

	ownership.Owned = true
}

// GetOwnedSkills は保有している全SkillTypeIDと保有情報を返します。
// 返されるマップは内部状態のコピーです。
func (inv *SkillInventory) GetOwnedSkills() map[string]*SkillOwnership {
	result := make(map[string]*SkillOwnership, len(inv.skills))
	for typeID, ownership := range inv.skills {
		ownershipCopy := &SkillOwnership{
			Owned: ownership.Owned,
		}
		result[typeID] = ownershipCopy
	}
	return result
}

// HasSkill は指定TypeIDのスキルを保有しているかを返します。
func (inv *SkillInventory) HasSkill(typeID string) bool {
	ownership, exists := inv.skills[typeID]
	if !exists {
		return false
	}
	return ownership.Owned
}
