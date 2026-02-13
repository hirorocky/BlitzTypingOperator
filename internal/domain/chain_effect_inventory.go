// Package domain はゲームのドメインモデルを定義します。
// このファイルはチェイン効果インベントリを定義します。

package domain

import "sort"

// ChainEffectInventory はチェイン効果のユニークインベントリを管理します。
// TypeIDベースのセット管理で、同一TypeIDの重複登録を防止します。
type ChainEffectInventory struct {
	chainEffects map[string]struct{}
}

// NewChainEffectInventory は新しい空のChainEffectInventoryを作成します。
func NewChainEffectInventory() *ChainEffectInventory {
	return &ChainEffectInventory{
		chainEffects: make(map[string]struct{}),
	}
}

// AddChainEffect はチェイン効果をインベントリに追加します。
// 空文字列の場合はfalseを返します。
// 既に保有している場合はfalseを返します。
// 新規追加の場合はtrueを返します。
func (inv *ChainEffectInventory) AddChainEffect(typeID string) bool {
	if typeID == "" {
		return false
	}
	if _, exists := inv.chainEffects[typeID]; exists {
		return false
	}
	inv.chainEffects[typeID] = struct{}{}
	return true
}

// HasChainEffect は指定TypeIDのチェイン効果を保有しているかを返します。
func (inv *ChainEffectInventory) HasChainEffect(typeID string) bool {
	_, exists := inv.chainEffects[typeID]
	return exists
}

// GetOwnedChainEffects は保有チェイン効果TypeIDの一覧をソート済みで返します。
func (inv *ChainEffectInventory) GetOwnedChainEffects() []string {
	result := make([]string, 0, len(inv.chainEffects))
	for typeID := range inv.chainEffects {
		result = append(result, typeID)
	}
	sort.Strings(result)
	return result
}
