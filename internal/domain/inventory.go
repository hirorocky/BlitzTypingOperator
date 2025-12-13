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

// ==================== モジュールインベントリ ====================

// ModuleInventory はモジュールのインベントリを管理する構造体です。

type ModuleInventory struct {
	// modules はモジュールのマップ（ID → ModuleModel）です。
	modules map[string]*ModuleModel

	// maxSlots はモジュールの最大保持数です。
	maxSlots int
}

// NewModuleInventory は新しいModuleInventoryを作成します。
func NewModuleInventory(maxSlots int) *ModuleInventory {
	return &ModuleInventory{
		modules:  make(map[string]*ModuleModel),
		maxSlots: maxSlots,
	}
}

// Add はモジュールをインベントリに追加します。
// 上限に達している場合はエラーを返します。

func (inv *ModuleInventory) Add(module *ModuleModel) error {
	if len(inv.modules) >= inv.maxSlots {
		return fmt.Errorf("モジュールインベントリが満杯です（上限: %d）", inv.maxSlots)
	}
	inv.modules[module.ID] = module
	return nil
}

// Remove はモジュールをインベントリから削除します。

func (inv *ModuleInventory) Remove(id string) *ModuleModel {
	module, exists := inv.modules[id]
	if !exists {
		return nil
	}
	delete(inv.modules, id)
	return module
}

// Get は指定されたIDのモジュールを取得します。
func (inv *ModuleInventory) Get(id string) *ModuleModel {
	return inv.modules[id]
}

// Count はインベントリ内のモジュール数を返します。
func (inv *ModuleInventory) Count() int {
	return len(inv.modules)
}

// MaxSlots はモジュールの最大保持数を返します。
func (inv *ModuleInventory) MaxSlots() int {
	return inv.maxSlots
}

// IsFull はインベントリが満杯かどうかを返します。
func (inv *ModuleInventory) IsFull() bool {
	return len(inv.modules) >= inv.maxSlots
}

// List は全てのモジュールをリストで返します。

func (inv *ModuleInventory) List() []*ModuleModel {
	result := make([]*ModuleModel, 0, len(inv.modules))
	for _, module := range inv.modules {
		result = append(result, module)
	}
	return result
}

// FilterByCategory は指定されたカテゴリでフィルタリングします。

func (inv *ModuleInventory) FilterByCategory(category ModuleCategory) []*ModuleModel {
	result := make([]*ModuleModel, 0)
	for _, module := range inv.modules {
		if module.Category == category {
			result = append(result, module)
		}
	}
	return result
}

// FilterByLevel は指定されたレベルでフィルタリングします。

func (inv *ModuleInventory) FilterByLevel(level int) []*ModuleModel {
	result := make([]*ModuleModel, 0)
	for _, module := range inv.modules {
		if module.Level == level {
			result = append(result, module)
		}
	}
	return result
}

// FilterByTag はタグでフィルタリングします。
func (inv *ModuleInventory) FilterByTag(tag string) []*ModuleModel {
	result := make([]*ModuleModel, 0)
	for _, module := range inv.modules {
		if module.HasTag(tag) {
			result = append(result, module)
		}
	}
	return result
}

// FilterCompatibleWithCore はコアに装備可能なモジュールのみをフィルタリングします。

func (inv *ModuleInventory) FilterCompatibleWithCore(core *CoreModel) []*ModuleModel {
	result := make([]*ModuleModel, 0)
	for _, module := range inv.modules {
		if module.IsCompatibleWithCore(core) {
			result = append(result, module)
		}
	}
	return result
}

// SortByLevel はレベルでソートしたモジュールリストを返します。

func (inv *ModuleInventory) SortByLevel(ascending bool) []*ModuleModel {
	result := inv.List()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].Level < result[j].Level
		}
		return result[i].Level > result[j].Level
	})
	return result
}

// SortByCategory はカテゴリでソートしたモジュールリストを返します。

func (inv *ModuleInventory) SortByCategory(ascending bool) []*ModuleModel {
	result := inv.List()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].Category < result[j].Category
		}
		return result[i].Category > result[j].Category
	})
	return result
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
