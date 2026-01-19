// Package domain はゲームのドメインモデルを定義します。
// このファイルはSkillInventoryを定義します。

package domain

// SkillOwnership はスキルの保有情報を表す構造体です。
// スキルの保有状態と取得済みチェイン効果バリエーションを管理します。
type SkillOwnership struct {
	// Owned はスキルを保有しているかどうかを示します。
	Owned bool

	// ChainVariations は取得済みチェイン効果IDのセットです。
	// map[string]bool でセット表現を行います。
	ChainVariations map[string]bool
}

// NewSkillOwnership は新しいSkillOwnershipを作成します。
func NewSkillOwnership() *SkillOwnership {
	return &SkillOwnership{
		Owned:           false,
		ChainVariations: make(map[string]bool),
	}
}

// AddChainVariation はチェイン効果バリエーションを追加します。
// 空文字列の場合は追加しません。
func (o *SkillOwnership) AddChainVariation(chainEffectID string) {
	// 空文字列は追加しない
	if chainEffectID == "" {
		return
	}
	o.ChainVariations[chainEffectID] = true
}

// HasChainVariation は指定のチェイン効果バリエーションを保有しているかを返します。
func (o *SkillOwnership) HasChainVariation(chainEffectID string) bool {
	return o.ChainVariations[chainEffectID]
}

// GetChainVariations は保有しているチェイン効果IDのリストを返します。
func (o *SkillOwnership) GetChainVariations() []string {
	result := make([]string, 0, len(o.ChainVariations))
	for chainID := range o.ChainVariations {
		result = append(result, chainID)
	}
	return result
}

// SkillInventory はスキルをTypeIDごとにユニーク管理する構造体です。
// 各TypeIDに対して保有状態とチェイン効果バリエーションを分離して追跡します。
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

// AddSkill はスキルを追加し、チェイン効果を登録します。
// chainEffectID: チェイン効果ID（なしの場合は空文字列）
// TypeIDが空の場合は何もしません。
func (inv *SkillInventory) AddSkill(typeID string, chainEffectID string) {
	// 検証: TypeIDが空でないこと
	if typeID == "" {
		return
	}

	// 既存のエントリを取得または新規作成
	ownership, exists := inv.skills[typeID]
	if !exists {
		ownership = NewSkillOwnership()
		inv.skills[typeID] = ownership
	}

	// 保有状態をtrueに設定
	ownership.Owned = true

	// チェイン効果がある場合は追加
	if chainEffectID != "" {
		ownership.AddChainVariation(chainEffectID)
	}
}

// GetOwnedSkills は保有している全SkillTypeIDと利用可能なチェイン効果を返します。
// 返されるマップは内部状態のコピーです。
func (inv *SkillInventory) GetOwnedSkills() map[string]*SkillOwnership {
	result := make(map[string]*SkillOwnership, len(inv.skills))
	for typeID, ownership := range inv.skills {
		// SkillOwnershipのコピーを作成
		ownershipCopy := &SkillOwnership{
			Owned:           ownership.Owned,
			ChainVariations: make(map[string]bool, len(ownership.ChainVariations)),
		}
		for chainID := range ownership.ChainVariations {
			ownershipCopy.ChainVariations[chainID] = true
		}
		result[typeID] = ownershipCopy
	}
	return result
}

// GetChainVariations は指定TypeIDで利用可能なチェイン効果IDリストを返します。
// 未保有の場合は空のリストを返します。
func (inv *SkillInventory) GetChainVariations(typeID string) []string {
	ownership, exists := inv.skills[typeID]
	if !exists {
		return []string{}
	}
	return ownership.GetChainVariations()
}

// HasSkill は指定TypeIDのスキルを保有しているかを返します。
func (inv *SkillInventory) HasSkill(typeID string) bool {
	ownership, exists := inv.skills[typeID]
	if !exists {
		return false
	}
	return ownership.Owned
}

// HasChainVariation は指定のチェイン効果バリエーションを保有しているかを返します。
func (inv *SkillInventory) HasChainVariation(typeID string, chainEffectID string) bool {
	ownership, exists := inv.skills[typeID]
	if !exists {
		return false
	}
	return ownership.HasChainVariation(chainEffectID)
}
