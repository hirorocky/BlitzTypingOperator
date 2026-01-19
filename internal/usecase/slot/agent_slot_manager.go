// Package slot はエージェントスロット管理のユースケースを提供します。
// 3つのエージェントスロットの管理とバトルシステムとの連携を担当します。

package slot

import (
	"hirorocky/type-battle/internal/domain"
)

// MaxAgentSlotCount はエージェントスロットの最大数です。
const MaxAgentSlotCount = 3

// AgentSlotManager は3つのエージェントスロットを管理するユースケースです。
// コア・スキルの付け替え操作を提供し、バトル用AgentModel構築を担当します。
type AgentSlotManager struct {
	// slots は3つの固定エージェントスロット
	slots [MaxAgentSlotCount]*domain.AgentSlot

	// coreInv はコアインベントリへの参照
	coreInv *domain.CoreInventory

	// skillInv はスキルインベントリへの参照
	skillInv *domain.SkillInventory

	// coreTypes はCoreTypeマスタデータへの参照
	coreTypes map[string]domain.CoreType

	// skillTypes はSkillTypeマスタデータへの参照
	skillTypes map[string]domain.SkillType

	// passiveSkills はPassiveSkillマスタデータへの参照
	passiveSkills map[string]domain.PassiveSkill
}

// NewAgentSlotManager は新しいAgentSlotManagerを作成します。
func NewAgentSlotManager(
	coreInv *domain.CoreInventory,
	skillInv *domain.SkillInventory,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
	passiveSkills map[string]domain.PassiveSkill,
) *AgentSlotManager {
	// 3つの空スロットを初期化
	var slots [MaxAgentSlotCount]*domain.AgentSlot
	for i := range slots {
		slots[i] = domain.NewAgentSlot()
	}

	return &AgentSlotManager{
		slots:         slots,
		coreInv:       coreInv,
		skillInv:      skillInv,
		coreTypes:     coreTypes,
		skillTypes:    skillTypes,
		passiveSkills: passiveSkills,
	}
}

// GetSlots は全スロットの構成を返します。
func (m *AgentSlotManager) GetSlots() [MaxAgentSlotCount]*domain.AgentSlot {
	return m.slots
}

// GetSlot は指定スロットの構成を返します。
// slotが範囲外の場合はnilを返します。
func (m *AgentSlotManager) GetSlot(slot int) *domain.AgentSlot {
	if slot < 0 || slot >= MaxAgentSlotCount {
		return nil
	}
	return m.slots[slot]
}
