// Package domain はゲームのドメインモデルを定義します。
// このファイルはCoreInventoryを定義します。

package domain

// CoreInventory はコアをTypeIDごとに管理する構造体です。
// 各TypeIDの保有フラグのみを管理します（レベル概念なし）。
type CoreInventory struct {
	// cores はCoreTypeID → 保有フラグのマップ
	cores map[string]bool
}

// NewCoreInventory は新しいCoreInventoryを作成します。
func NewCoreInventory() *CoreInventory {
	return &CoreInventory{
		cores: make(map[string]bool),
	}
}

// AddCore はコアを追加します。
// 新規TypeIDの場合にtrueを返し、既に保有している場合はfalseを返します。
// TypeIDが空の場合はfalseを返します。
func (inv *CoreInventory) AddCore(typeID string) bool {
	// 検証: TypeIDが空でないこと
	if typeID == "" {
		return false
	}

	// 既に保有している場合はfalse
	if inv.cores[typeID] {
		return false
	}

	// 新規追加
	inv.cores[typeID] = true
	return true
}

// GetOwnedCores は保有している全CoreTypeIDを返します。
// 返されるスライスは内部状態のコピーです。
// 注意: 返されるスライスの順序は不定です。
func (inv *CoreInventory) GetOwnedCores() []string {
	result := make([]string, 0, len(inv.cores))
	for typeID := range inv.cores {
		result = append(result, typeID)
	}
	return result
}

// HasCore は指定TypeIDのコアを保有しているかを返します。
func (inv *CoreInventory) HasCore(typeID string) bool {
	return inv.cores[typeID]
}
