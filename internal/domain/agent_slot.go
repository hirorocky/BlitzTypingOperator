// Package domain はゲームのドメインモデルを定義します。
// このファイルはAgentSlotを定義します。

package domain

// ==================== SkillSlotConfig 値オブジェクト ====================

// SkillSlotConfig はスキルスロットの構成を表す値オブジェクトです。
// スキルTypeIDとチェイン効果IDを保持します。
type SkillSlotConfig struct {
	// TypeID はスキルTypeID（空の場合は未設定）
	TypeID string

	// ChainEffectID はチェイン効果ID（なしの場合は空文字列）
	ChainEffectID string
}

// IsEmpty はスキルスロットが空かどうかを返します。
// TypeIDが空の場合に空とみなします。
func (c *SkillSlotConfig) IsEmpty() bool {
	return c.TypeID == ""
}

// Clear はスキルスロットをクリアします。
func (c *SkillSlotConfig) Clear() {
	c.TypeID = ""
	c.ChainEffectID = ""
}

// ==================== AgentSlot ドメインモデル ====================

// MaxSkillSlotCount はエージェント1体あたりの最大スキルスロット数です。
const MaxSkillSlotCount = 4

// AgentSlot はエージェントスロット1つの構成を表すドメインモデルです。
// スロットはコア（TypeID）と最大4つのスキルで構成されます。
type AgentSlot struct {
	// CoreTypeID はコアTypeID（空の場合はスロット空）
	CoreTypeID string

	// Skills はスキルスロット構成（最大4つ）
	Skills [MaxSkillSlotCount]SkillSlotConfig
}

// NewAgentSlot は新しい空のAgentSlotを作成します。
func NewAgentSlot() *AgentSlot {
	return &AgentSlot{
		CoreTypeID: "",
		Skills:     [MaxSkillSlotCount]SkillSlotConfig{},
	}
}

// IsEmpty はスロットが空かどうかを返します。
// コア未設定の場合に空とみなします。
func (s *AgentSlot) IsEmpty() bool {
	return s.CoreTypeID == ""
}

// GetSkillCount は設定されているスキル数を返します。
func (s *AgentSlot) GetSkillCount() int {
	count := 0
	for i := range s.Skills {
		if !s.Skills[i].IsEmpty() {
			count++
		}
	}
	return count
}

// Clear はスロットをクリアします。
// コアとすべてのスキルが初期化されます。
func (s *AgentSlot) Clear() {
	s.CoreTypeID = ""
	for i := range s.Skills {
		s.Skills[i].Clear()
	}
}

// SetCore はスロットにコアを設定します。
// 検証は行いません。AgentSlotManagerを通じて設定してください。
func (s *AgentSlot) SetCore(typeID string) {
	s.CoreTypeID = typeID
}

// SetSkill はスロットのスキルを設定します。
// skillSlotは0-3の範囲である必要があります。範囲外の場合は何もしません。
// 検証は行いません。AgentSlotManagerを通じて設定してください。
func (s *AgentSlot) SetSkill(skillSlot int, typeID string, chainEffectID string) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.Skills[skillSlot].TypeID = typeID
	s.Skills[skillSlot].ChainEffectID = chainEffectID
}

// ClearSkill は指定スキルスロットをクリアします。
// skillSlotは0-3の範囲である必要があります。範囲外の場合は何もしません。
func (s *AgentSlot) ClearSkill(skillSlot int) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.Skills[skillSlot].Clear()
}

// GetSkill は指定スキルスロットの構成を返します。
// skillSlotは0-3の範囲である必要があります。範囲外の場合はnilを返します。
func (s *AgentSlot) GetSkill(skillSlot int) *SkillSlotConfig {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return nil
	}
	return &s.Skills[skillSlot]
}
