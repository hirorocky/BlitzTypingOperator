// Package inventory はインベントリ管理機能を提供します。
// コア、モジュール、エージェントの保管と管理を担当します。

package domain

import (
	"fmt"
	"sort"
)

// ==================== コアインベントリ ====================

// CoreInventory はコアのインベントリを管理する構造体です。

type CoreInventory struct {
	// cores はコアのマップ（ID → CoreModel）です。
	cores map[string]*CoreModel

	// maxSlots はコアの最大保持数です。
	maxSlots int
}

// NewCoreInventory は新しいCoreInventoryを作成します。
func NewCoreInventory(maxSlots int) *CoreInventory {
	return &CoreInventory{
		cores:    make(map[string]*CoreModel),
		maxSlots: maxSlots,
	}
}

// Add はコアをインベントリに追加します。
// 上限に達している場合はエラーを返します。

func (inv *CoreInventory) Add(core *CoreModel) error {
	if len(inv.cores) >= inv.maxSlots {
		return fmt.Errorf("コアインベントリが満杯です（上限: %d）", inv.maxSlots)
	}
	inv.cores[core.ID] = core
	return nil
}

// Remove はコアをインベントリから削除します。

func (inv *CoreInventory) Remove(id string) *CoreModel {
	core, exists := inv.cores[id]
	if !exists {
		return nil
	}
	delete(inv.cores, id)
	return core
}

// Get は指定されたIDのコアを取得します。
func (inv *CoreInventory) Get(id string) *CoreModel {
	return inv.cores[id]
}

// Count はインベントリ内のコア数を返します。
func (inv *CoreInventory) Count() int {
	return len(inv.cores)
}

// MaxSlots はコアの最大保持数を返します。
func (inv *CoreInventory) MaxSlots() int {
	return inv.maxSlots
}

// IsFull はインベントリが満杯かどうかを返します。
func (inv *CoreInventory) IsFull() bool {
	return len(inv.cores) >= inv.maxSlots
}

// List は全てのコアをリストで返します。

func (inv *CoreInventory) List() []*CoreModel {
	result := make([]*CoreModel, 0, len(inv.cores))
	for _, core := range inv.cores {
		result = append(result, core)
	}
	return result
}

// FilterByType は指定されたコア特性でフィルタリングします。

func (inv *CoreInventory) FilterByType(typeID string) []*CoreModel {
	result := make([]*CoreModel, 0)
	for _, core := range inv.cores {
		if core.Type.ID == typeID {
			result = append(result, core)
		}
	}
	return result
}

// FilterByLevelRange は指定されたレベル範囲でフィルタリングします。

func (inv *CoreInventory) FilterByLevelRange(minLevel, maxLevel int) []*CoreModel {
	result := make([]*CoreModel, 0)
	for _, core := range inv.cores {
		if core.Level >= minLevel && core.Level <= maxLevel {
			result = append(result, core)
		}
	}
	return result
}

// SortByLevel はレベルでソートしたコアリストを返します。

// ascending: trueなら昇順、falseなら降順
func (inv *CoreInventory) SortByLevel(ascending bool) []*CoreModel {
	result := inv.List()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].Level < result[j].Level
		}
		return result[i].Level > result[j].Level
	})
	return result
}

// SortByType は特性名でソートしたコアリストを返します。

func (inv *CoreInventory) SortByType(ascending bool) []*CoreModel {
	result := inv.List()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].Type.Name < result[j].Type.Name
		}
		return result[i].Type.Name > result[j].Type.Name
	})
	return result
}

// ==================== スキルインベントリ ====================

// SkillInventoryLegacy はスキルのインベントリを管理する構造体です。
// スキルはTypeIDとChainEffectの組み合わせで識別されるため、スライスで管理します。
// 注意: この構造体は後方互換性のために残されています。
// 新しいコードではTask 3で実装されるSkillInventory（ユニーク管理版）を使用してください。

type SkillInventoryLegacy struct {
	// skills はスキルのスライスです。
	skills []*SkillModel

	// maxSlots はスキルの最大保持数です。
	maxSlots int
}

// NewSkillInventoryLegacy は新しいSkillInventoryLegacyを作成します。
func NewSkillInventoryLegacy(maxSlots int) *SkillInventoryLegacy {
	return &SkillInventoryLegacy{
		skills:   make([]*SkillModel, 0),
		maxSlots: maxSlots,
	}
}

// Add はスキルをインベントリに追加します。
// 上限に達している場合はエラーを返します。

func (inv *SkillInventoryLegacy) Add(skill *SkillModel) error {
	if len(inv.skills) >= inv.maxSlots {
		return fmt.Errorf("スキルインベントリが満杯です（上限: %d）", inv.maxSlots)
	}
	inv.skills = append(inv.skills, skill)
	return nil
}

// Remove はスキルをインベントリから削除します。
// インデックスで指定します。無効なインデックスの場合はnilを返します。

func (inv *SkillInventoryLegacy) Remove(index int) *SkillModel {
	if index < 0 || index >= len(inv.skills) {
		return nil
	}
	skill := inv.skills[index]
	inv.skills = append(inv.skills[:index], inv.skills[index+1:]...)
	return skill
}

// RemoveByTypeID は指定されたTypeIDの最初のスキルを削除します。
// 後方互換性のためのメソッドです。

func (inv *SkillInventoryLegacy) RemoveByTypeID(typeID string) *SkillModel {
	for i, skill := range inv.skills {
		if skill.TypeID == typeID {
			return inv.Remove(i)
		}
	}
	return nil
}

// Get は指定されたインデックスのスキルを取得します。
func (inv *SkillInventoryLegacy) Get(index int) *SkillModel {
	if index < 0 || index >= len(inv.skills) {
		return nil
	}
	return inv.skills[index]
}

// GetByTypeID は指定されたTypeIDの最初のスキルを取得します。
// 後方互換性のためのメソッドです。
func (inv *SkillInventoryLegacy) GetByTypeID(typeID string) *SkillModel {
	for _, skill := range inv.skills {
		if skill.TypeID == typeID {
			return skill
		}
	}
	return nil
}

// Count はインベントリ内のスキル数を返します。
func (inv *SkillInventoryLegacy) Count() int {
	return len(inv.skills)
}

// MaxSlots はスキルの最大保持数を返します。
func (inv *SkillInventoryLegacy) MaxSlots() int {
	return inv.maxSlots
}

// IsFull はインベントリが満杯かどうかを返します。
func (inv *SkillInventoryLegacy) IsFull() bool {
	return len(inv.skills) >= inv.maxSlots
}

// List は全てのスキルをリストで返します。

func (inv *SkillInventoryLegacy) List() []*SkillModel {
	result := make([]*SkillModel, len(inv.skills))
	copy(result, inv.skills)
	return result
}

// FilterByDamageEffect はダメージ効果を持つスキルをフィルタリングします。
func (inv *SkillInventoryLegacy) FilterByDamageEffect() []*SkillModel {
	result := make([]*SkillModel, 0)
	for _, skill := range inv.skills {
		for _, effect := range skill.Effects() {
			if effect.IsDamageEffect() {
				result = append(result, skill)
				break
			}
		}
	}
	return result
}

// FilterByHealEffect は回復効果を持つスキルをフィルタリングします。
func (inv *SkillInventoryLegacy) FilterByHealEffect() []*SkillModel {
	result := make([]*SkillModel, 0)
	for _, skill := range inv.skills {
		for _, effect := range skill.Effects() {
			if effect.IsHealEffect() {
				result = append(result, skill)
				break
			}
		}
	}
	return result
}

// FilterByTag はタグでフィルタリングします。
func (inv *SkillInventoryLegacy) FilterByTag(tag string) []*SkillModel {
	result := make([]*SkillModel, 0)
	for _, skill := range inv.skills {
		if skill.HasTag(tag) {
			result = append(result, skill)
		}
	}
	return result
}

// FilterCompatibleWithCore はコアに装備可能なスキルのみをフィルタリングします。

func (inv *SkillInventoryLegacy) FilterCompatibleWithCore(core *CoreModel) []*SkillModel {
	result := make([]*SkillModel, 0)
	for _, skill := range inv.skills {
		if skill.IsCompatibleWithCore(core) {
			result = append(result, skill)
		}
	}
	return result
}

// ModuleInventory はSkillInventoryLegacyのエイリアスです。
// 後方互換性のために残されています。新規コードではSkillInventoryLegacyを使用してください。
type ModuleInventory = SkillInventoryLegacy

// NewModuleInventory はNewSkillInventoryLegacyの後方互換ラッパーです。
// 後方互換性のために残されています。
func NewModuleInventory(maxSlots int) *ModuleInventory {
	return NewSkillInventoryLegacy(maxSlots)
}

// ==================== エージェントインベントリ ====================

// AgentInventory はエージェントのインベントリを管理する構造体です。
// 装備状態はAgentManagerで一元管理されます。

type AgentInventory struct {
	// agents はエージェントのマップ（ID → AgentModel）です。
	agents map[string]*AgentModel

	// maxSlots はエージェントの最大保持数です。

	maxSlots int
}

// NewAgentInventory は新しいAgentInventoryを作成します。

func NewAgentInventory(maxSlots int) *AgentInventory {
	return &AgentInventory{
		agents:   make(map[string]*AgentModel),
		maxSlots: maxSlots,
	}
}

// NewAgentInventoryWithDefault はデフォルトの最低保持数（20体）を保証して作成します。
func NewAgentInventoryWithDefault(maxSlots int) *AgentInventory {
	if maxSlots < 20 {
		maxSlots = 20 // 最低20体を保証
	}
	return &AgentInventory{
		agents:   make(map[string]*AgentModel),
		maxSlots: maxSlots,
	}
}

// Add はエージェントをインベントリに追加します。
// 上限に達している場合はエラーを返します。

func (inv *AgentInventory) Add(agent *AgentModel) error {
	if len(inv.agents) >= inv.maxSlots {
		return fmt.Errorf("エージェントインベントリが満杯です（上限: %d）", inv.maxSlots)
	}
	inv.agents[agent.ID] = agent
	return nil
}

// Remove はエージェントをインベントリから削除します。
// 装備状態はAgentManagerで管理されているため、装備解除は別途行う必要があります。

func (inv *AgentInventory) Remove(id string) *AgentModel {
	agent, exists := inv.agents[id]
	if !exists {
		return nil
	}
	delete(inv.agents, id)
	return agent
}

// Get は指定されたIDのエージェントを取得します。
func (inv *AgentInventory) Get(id string) *AgentModel {
	return inv.agents[id]
}

// Count はインベントリ内のエージェント数を返します。
func (inv *AgentInventory) Count() int {
	return len(inv.agents)
}

// MaxSlots はエージェントの最大保持数を返します。
func (inv *AgentInventory) MaxSlots() int {
	return inv.maxSlots
}

// IsFull はインベントリが満杯かどうかを返します。
func (inv *AgentInventory) IsFull() bool {
	return len(inv.agents) >= inv.maxSlots
}

// List は全てのエージェントをリストで返します。
func (inv *AgentInventory) List() []*AgentModel {
	result := make([]*AgentModel, 0, len(inv.agents))
	for _, agent := range inv.agents {
		result = append(result, agent)
	}
	return result
}
