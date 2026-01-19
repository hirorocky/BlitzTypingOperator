// Package domain はゲームのドメインモデルを定義します。
// このファイルはCoreInventoryを定義します。

package domain

// CoreInventory はコアをTypeIDごとに管理する構造体です。
// 各TypeIDに対して取得済み最大レベルのみを保存します。
type CoreInventory struct {
	// cores はCoreTypeID → 取得済み最大レベルのマップ
	cores map[string]int
}

// NewCoreInventory は新しいCoreInventoryを作成します。
func NewCoreInventory() *CoreInventory {
	return &CoreInventory{
		cores: make(map[string]int),
	}
}

// AddCore はコアを追加し、レベル比較で更新判定を行います。
// 新規TypeIDの場合、または既存の最大レベルより高い場合に更新します。
// 更新された場合はtrue、更新されなかった場合はfalseを返します。
// レベルが1未満の場合、またはTypeIDが空の場合はfalseを返します。
func (inv *CoreInventory) AddCore(typeID string, level int) bool {
	// 検証: TypeIDが空でないこと
	if typeID == "" {
		return false
	}

	// 検証: レベルが1以上の正整数であること
	if level < 1 {
		return false
	}

	// 既存の最大レベルを取得
	currentMax, exists := inv.cores[typeID]

	// 新規TypeIDまたは既存最大レベルより高い場合のみ更新
	if !exists || level > currentMax {
		inv.cores[typeID] = level
		return true
	}

	return false
}

// GetMaxLevel は指定TypeIDの取得済み最大レベルを返します。
// 未保有の場合は0を返します。
func (inv *CoreInventory) GetMaxLevel(typeID string) int {
	return inv.cores[typeID]
}

// GetOwnedCores は保有している全CoreTypeIDとその最大レベルを返します。
// 返されるマップは内部状態のコピーです。
func (inv *CoreInventory) GetOwnedCores() map[string]int {
	result := make(map[string]int, len(inv.cores))
	for typeID, level := range inv.cores {
		result[typeID] = level
	}
	return result
}

// HasCore は指定TypeIDのコアを保有しているかを返します。
func (inv *CoreInventory) HasCore(typeID string) bool {
	_, exists := inv.cores[typeID]
	return exists
}
