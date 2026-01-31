// Package slot はエージェントスロット管理のユースケースを提供します。
// 3つのエージェントスロットの管理とバトルシステムとの連携を担当します。

package slot

import (
	"errors"
	"slices"

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

	// ErrSlotIndexOutOfRange は無効なスロットインデックスの場合のエラー
	ErrSlotIndexOutOfRange = errors.New("slot index out of range")

	// ErrCoreNotSet はコア未設定スロットにスキルを設定しようとした場合のエラー
	ErrCoreNotSet = errors.New("core not set")

	// ErrSkillSlotIndexOutOfRange は無効なスキルスロットインデックスの場合のエラー
	ErrSkillSlotIndexOutOfRange = errors.New("skill slot index out of range")

	// ErrSlotLocked はバトル中にスロット変更を試みた場合のエラー
	ErrSlotLocked = errors.New("slot is locked during battle")

	// ErrSkillAlreadyEquipped は既に他のスロットで装備済みのスキルを設定しようとした場合のエラー
	ErrSkillAlreadyEquipped = errors.New("skill already equipped in another slot")
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

	// chainEffects はChainEffectマスタデータへの参照（ID→ChainEffect）
	chainEffects map[string]domain.ChainEffect

	// locked はバトル中のスロット変更禁止フラグ
	locked bool
}

// NewAgentSlotManager は新しいAgentSlotManagerを作成します。
func NewAgentSlotManager(
	coreInv *domain.CoreInventory,
	skillInv *domain.SkillInventory,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
	passiveSkills map[string]domain.PassiveSkill,
	chainEffects map[string]domain.ChainEffect,
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
		chainEffects:  chainEffects,
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
// コア変更時に互換性のないスキルは自動削除されます。
// ロック中はErrSlotLockedを返します。
func (m *AgentSlotManager) SetCore(slot int, typeID string) error {
	// ロックチェック
	if m.locked {
		return ErrSlotLocked
	}

	// スロットインデックスの検証
	if slot < 0 || slot >= MaxAgentSlotCount {
		return ErrSlotIndexOutOfRange
	}

	// インベントリ保有確認
	if !m.coreInv.HasCore(typeID) {
		return ErrCoreNotOwned
	}

	// スロットにコアを設定
	targetSlot := m.slots[slot]
	targetSlot.SetCore(typeID)

	// 互換性チェックを実行し、互換性のないスキルを自動削除
	m.removeIncompatibleSkills(slot)

	return nil
}

// ClearCore はスロットのコアをクリアします。
// コアをクリアすると、スキルも全て削除されます。
// ロック中はErrSlotLockedを返します。
func (m *AgentSlotManager) ClearCore(slot int) error {
	// ロックチェック
	if m.locked {
		return ErrSlotLocked
	}

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
	for i := range domain.MaxSkillSlotCount {
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
		if slices.Contains(coreType.AllowedTags, skillTag) {
			return true
		}
	}
	return false
}

// SetSkill はスロットのスキルを設定します。
// 同じスキルIDは全エージェントを通じて1つしか装備できません。
// ロック中はErrSlotLockedを返します。
func (m *AgentSlotManager) SetSkill(slot int, skillSlot int, typeID string, chainEffectID string) error {
	// ロックチェック
	if m.locked {
		return ErrSlotLocked
	}

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

	// 全スロットでのスキル重複チェック（同じスロット位置への上書きは許可）
	if m.isSkillEquippedElsewhere(slot, skillSlot, typeID) {
		return ErrSkillAlreadyEquipped
	}

	// スキルを設定
	targetSlot.SetSkill(skillSlot, typeID, chainEffectID)

	return nil
}

// isSkillEquippedElsewhere は指定スキルが他のスロットで装備されているかをチェックします。
// 同じ位置（slot, skillSlot）への上書きは重複とみなしません。
func (m *AgentSlotManager) isSkillEquippedElsewhere(targetSlot int, targetSkillSlot int, typeID string) bool {
	for slotIdx := range MaxAgentSlotCount {
		agentSlot := m.slots[slotIdx]
		if agentSlot == nil || agentSlot.IsEmpty() {
			continue
		}

		for skillSlotIdx := range domain.MaxSkillSlotCount {
			// 同じ位置への上書きはスキップ
			if slotIdx == targetSlot && skillSlotIdx == targetSkillSlot {
				continue
			}

			skillConfig := agentSlot.GetSkill(skillSlotIdx)
			if skillConfig != nil && !skillConfig.IsEmpty() && skillConfig.TypeID == typeID {
				return true
			}
		}
	}
	return false
}

// ClearSkill はスロットの指定スキルスロットをクリアします。
// ロック中はErrSlotLockedを返します。
func (m *AgentSlotManager) ClearSkill(slot int, skillSlot int) error {
	// ロックチェック
	if m.locked {
		return ErrSlotLocked
	}

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

// ==================== バトル連携機能 ====================

// IsSlotReady はスロットがバトルに使用可能かを返します。
// コアが設定されている場合に使用可能です。
func (m *AgentSlotManager) IsSlotReady(slot int) bool {
	if slot < 0 || slot >= MaxAgentSlotCount {
		return false
	}
	return !m.slots[slot].IsEmpty()
}

// GetReadySlotCount はバトルに使用可能なスロット数を返します。
func (m *AgentSlotManager) GetReadySlotCount() int {
	count := 0
	for i := range MaxAgentSlotCount {
		if m.IsSlotReady(i) {
			count++
		}
	}
	return count
}

// ValidateSkillCompatibility はスキルがスロットのコアと互換性があるかを検証します。
// スロットが空の場合やスキルTypeが見つからない場合はfalseを返します。
func (m *AgentSlotManager) ValidateSkillCompatibility(slot int, skillTypeID string) bool {
	if slot < 0 || slot >= MaxAgentSlotCount {
		return false
	}

	targetSlot := m.slots[slot]
	if targetSlot.IsEmpty() {
		return false
	}

	// コアTypeを取得
	coreType, exists := m.coreTypes[targetSlot.CoreTypeID]
	if !exists {
		return false
	}

	// スキルTypeを取得
	skillType, exists := m.skillTypes[skillTypeID]
	if !exists {
		return false
	}

	return m.isSkillCompatibleWithCoreType(skillType, coreType)
}

// GetCompatibleSkills は指定スロットのコアと互換性のあるスキル一覧を返します。
// インベントリに保有しているスキルのうち、互換性のあるもののTypeIDリストを返します。
// スロットが空の場合は空のリストを返します。
func (m *AgentSlotManager) GetCompatibleSkills(slot int) []string {
	result := []string{}

	if slot < 0 || slot >= MaxAgentSlotCount {
		return result
	}

	targetSlot := m.slots[slot]
	if targetSlot.IsEmpty() {
		return result
	}

	// コアTypeを取得
	coreType, exists := m.coreTypes[targetSlot.CoreTypeID]
	if !exists {
		return result
	}

	// 保有スキルを取得
	ownedSkills := m.skillInv.GetOwnedSkills()

	// 各スキルの互換性をチェック
	for typeID := range ownedSkills {
		skillType, exists := m.skillTypes[typeID]
		if !exists {
			continue
		}

		if m.isSkillCompatibleWithCoreType(skillType, coreType) {
			result = append(result, typeID)
		}
	}

	return result
}

// BuildAgentsForBattle はバトル用のAgentModelスライスを構築します。
// 空スロットはバトルに含めません。
func (m *AgentSlotManager) BuildAgentsForBattle() []*domain.AgentModel {
	agents := make([]*domain.AgentModel, 0, MaxAgentSlotCount)

	for i := range MaxAgentSlotCount {
		if !m.IsSlotReady(i) {
			continue
		}

		agent := m.buildAgentFromSlot(i)
		if agent != nil {
			agents = append(agents, agent)
		}
	}

	return agents
}

// buildAgentFromSlot はスロットからAgentModelを構築します。
func (m *AgentSlotManager) buildAgentFromSlot(slot int) *domain.AgentModel {
	targetSlot := m.slots[slot]
	if targetSlot.IsEmpty() {
		return nil
	}

	// CoreTypeを取得
	coreType, exists := m.coreTypes[targetSlot.CoreTypeID]
	if !exists {
		return nil
	}

	// PassiveSkillを取得
	passiveSkill := m.passiveSkills[coreType.PassiveSkillID]

	// CoreModelを作成（レベルなし）
	core := domain.NewCoreWithTypeID(
		targetSlot.CoreTypeID,
		coreType,
		passiveSkill,
	)

	// スキルを構築
	modules := make([]*domain.ModuleModel, 0, domain.MaxSkillSlotCount)
	for i := range domain.MaxSkillSlotCount {
		skillConfig := targetSlot.GetSkill(i)
		if skillConfig == nil || skillConfig.IsEmpty() {
			continue
		}

		module := m.buildModuleFromConfig(skillConfig)
		if module != nil {
			modules = append(modules, module)
		}
	}

	// エージェントIDを生成（スロット番号ベース）
	agentID := "agent_slot_" + string(rune('0'+slot))

	return domain.NewAgent(agentID, core, modules)
}

// buildModuleFromConfig はスキルスロット構成からModuleModelを構築します。
func (m *AgentSlotManager) buildModuleFromConfig(config *domain.SkillSlotConfig) *domain.ModuleModel {
	skillType, exists := m.skillTypes[config.TypeID]
	if !exists {
		return nil
	}

	// チェイン効果の取得
	var chainEffect *domain.ChainEffect = nil
	if config.ChainEffectID != "" && m.chainEffects != nil {
		if ce, exists := m.chainEffects[config.ChainEffectID]; exists {
			chainEffect = &ce
		}
	}

	return domain.NewSkillFromType(skillType, chainEffect)
}

// ==================== バトル中のスロット変更禁止機能 ====================

// Lock はスロットをロックし、バトル中の変更を禁止します。
// バトル開始時に呼び出してください。
func (m *AgentSlotManager) Lock() {
	m.locked = true
}

// Unlock はスロットのロックを解除し、変更を許可します。
// バトル終了時に呼び出してください。
func (m *AgentSlotManager) Unlock() {
	m.locked = false
}

// IsLocked はスロットがロックされているかどうかを返します。
func (m *AgentSlotManager) IsLocked() bool {
	return m.locked
}
