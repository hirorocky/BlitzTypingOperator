// Package slot はエージェントスロット管理のユースケースを提供します。
// 3つのエージェントスロットの管理とバトルシステムとの連携を担当します。

package slot

import (
	"errors"

	"hirorocky/type-battle/internal/domain"
)

// エラー定義
var (
	// ErrCoreNotOwned は保有していないコアを設定しようとした場合のエラー
	ErrCoreNotOwned = errors.New("core not owned")

	// ErrSkillNotOwned は保有していないスキルを設定しようとした場合のエラー
	ErrSkillNotOwned = errors.New("skill not owned")

	// ErrChainVariationNotOwned は保有していないチェイン効果を選択しようとした場合のエラー
	ErrChainVariationNotOwned = errors.New("chain variation not owned")

	// ErrSkillIncompatible はコアと互換性のないスキルを設定しようとした場合のエラー
	ErrSkillIncompatible = errors.New("skill incompatible with core")

	// ErrLevelOutOfRange は取得済み最大レベルを超えるレベルを選択しようとした場合のエラー
	ErrLevelOutOfRange = errors.New("level out of range")

	// ErrSlotIndexOutOfRange は無効なスロットインデックスの場合のエラー
	ErrSlotIndexOutOfRange = errors.New("slot index out of range")

	// ErrCoreNotSet はコア未設定スロットにスキルを設定しようとした場合のエラー
	ErrCoreNotSet = errors.New("core not set")

	// ErrSkillSlotIndexOutOfRange は無効なスキルスロットインデックスの場合のエラー
	ErrSkillSlotIndexOutOfRange = errors.New("skill slot index out of range")
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

// SetCore はスロットにコアを設定します。
// levelは1から取得済み最大レベルまでの任意のレベルを指定できます。
// コア変更時に互換性のないスキルは自動削除されます。
func (m *AgentSlotManager) SetCore(slot int, typeID string, level int) error {
	// スロットインデックスの検証
	if slot < 0 || slot >= MaxAgentSlotCount {
		return ErrSlotIndexOutOfRange
	}

	// インベントリ保有確認
	if !m.coreInv.HasCore(typeID) {
		return ErrCoreNotOwned
	}

	// レベル範囲確認（1〜最大レベル）
	maxLevel := m.coreInv.GetMaxLevel(typeID)
	if level < 1 || level > maxLevel {
		return ErrLevelOutOfRange
	}

	// スロットにコアを設定
	targetSlot := m.slots[slot]
	targetSlot.SetCore(typeID, level)

	// 互換性チェックを実行し、互換性のないスキルを自動削除
	m.removeIncompatibleSkills(slot)

	return nil
}

// ClearCore はスロットのコアをクリアします。
// コアをクリアすると、スキルも全て削除されます。
func (m *AgentSlotManager) ClearCore(slot int) error {
	// スロットインデックスの検証
	if slot < 0 || slot >= MaxAgentSlotCount {
		return ErrSlotIndexOutOfRange
	}

	// スロットをクリア
	m.slots[slot].Clear()

	return nil
}

// removeIncompatibleSkills はスロット内の互換性のないスキルを削除します。
func (m *AgentSlotManager) removeIncompatibleSkills(slot int) {
	targetSlot := m.slots[slot]

	// コアが未設定の場合は何もしない
	if targetSlot.IsEmpty() {
		return
	}

	// コアTypeを取得
	coreType, exists := m.coreTypes[targetSlot.CoreTypeID]
	if !exists {
		return
	}

	// 各スキルスロットをチェック
	for i := 0; i < domain.MaxSkillSlotCount; i++ {
		skillConfig := targetSlot.GetSkill(i)
		if skillConfig == nil || skillConfig.IsEmpty() {
			continue
		}

		// スキルTypeを取得
		skillType, exists := m.skillTypes[skillConfig.TypeID]
		if !exists {
			// マスタデータに存在しないスキルは削除
			targetSlot.ClearSkill(i)
			continue
		}

		// 互換性チェック
		if !m.isSkillCompatibleWithCoreType(skillType, coreType) {
			targetSlot.ClearSkill(i)
		}
	}
}

// isSkillCompatibleWithCoreType はスキルがコアTypeと互換性があるかを判定します。
func (m *AgentSlotManager) isSkillCompatibleWithCoreType(skillType domain.SkillType, coreType domain.CoreType) bool {
	// スキルのタグのうち1つでもコアの許可タグに含まれていれば互換性あり
	for _, skillTag := range skillType.Tags {
		for _, allowedTag := range coreType.AllowedTags {
			if skillTag == allowedTag {
				return true
			}
		}
	}
	return false
}

// SetSkill はスロットのスキルを設定します。
func (m *AgentSlotManager) SetSkill(slot int, skillSlot int, typeID string, chainEffectID string) error {
	// スロットインデックスの検証
	if slot < 0 || slot >= MaxAgentSlotCount {
		return ErrSlotIndexOutOfRange
	}

	// スキルスロットインデックスの検証
	if skillSlot < 0 || skillSlot >= domain.MaxSkillSlotCount {
		return ErrSkillSlotIndexOutOfRange
	}

	// コア未設定チェック
	targetSlot := m.slots[slot]
	if targetSlot.IsEmpty() {
		return ErrCoreNotSet
	}

	// インベントリ保有確認
	if !m.skillInv.HasSkill(typeID) {
		return ErrSkillNotOwned
	}

	// チェイン効果バリエーションの保有確認
	if chainEffectID != "" && !m.skillInv.HasChainVariation(typeID, chainEffectID) {
		return ErrChainVariationNotOwned
	}

	// スキルTypeを取得
	skillType, exists := m.skillTypes[typeID]
	if !exists {
		return ErrSkillNotOwned
	}

	// コアTypeを取得
	coreType, exists := m.coreTypes[targetSlot.CoreTypeID]
	if !exists {
		return ErrCoreNotSet
	}

	// タグマッチングによる互換性検証
	if !m.isSkillCompatibleWithCoreType(skillType, coreType) {
		return ErrSkillIncompatible
	}

	// スキルを設定
	targetSlot.SetSkill(skillSlot, typeID, chainEffectID)

	return nil
}

// ClearSkill はスロットの指定スキルスロットをクリアします。
func (m *AgentSlotManager) ClearSkill(slot int, skillSlot int) error {
	// スロットインデックスの検証
	if slot < 0 || slot >= MaxAgentSlotCount {
		return ErrSlotIndexOutOfRange
	}

	// スキルスロットインデックスの検証
	if skillSlot < 0 || skillSlot >= domain.MaxSkillSlotCount {
		return ErrSkillSlotIndexOutOfRange
	}

	// スキルスロットをクリア
	m.slots[slot].ClearSkill(skillSlot)

	return nil
}
