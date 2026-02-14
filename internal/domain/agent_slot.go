// Package domain はゲームのドメインモデルを定義します。
// このファイルはAgentSlotを定義します。

package domain

// ==================== SkillSlotConfig 値オブジェクト ====================

// SkillSlotConfig はスキルスロットの構成を表す値オブジェクトです。
type SkillSlotConfig struct {
	// TypeID はスキルTypeID（空の場合は未設定）
	TypeID string
}

// IsEmpty はスキルスロットが空かどうかを返します。
// TypeIDが空の場合に空とみなします。
func (c *SkillSlotConfig) IsEmpty() bool {
	return c.TypeID == ""
}

// Clear はスキルスロットをクリアします。
func (c *SkillSlotConfig) Clear() {
	c.TypeID = ""
}

// ==================== ChainEffectSlotConfig 値オブジェクト ====================

// ChainEffectSlotConfig はチェイン効果スロットの構成を表す値オブジェクトです。
type ChainEffectSlotConfig struct {
	// TypeID はチェイン効果TypeID（空の場合は未設定）
	TypeID string
}

// IsEmpty はチェイン効果スロットが空かどうかを返します。
func (c *ChainEffectSlotConfig) IsEmpty() bool {
	return c.TypeID == ""
}

// Clear はチェイン効果スロットをクリアします。
func (c *ChainEffectSlotConfig) Clear() {
	c.TypeID = ""
}

// ==================== AgentSlot ドメインモデル ====================

// MaxSkillSlotCount はエージェント1体あたりの最大スキルスロット数です。
const MaxSkillSlotCount = 4

// AgentSlot はエージェントスロット1つの構成を表すドメインモデルです。
// スロットはコア（TypeID）、最大4つのスキル、各スキルに対応するチェイン効果で構成されます。
type AgentSlot struct {
	// CoreTypeID はコアTypeID（空の場合はスロット空）
	CoreTypeID string

	// Skills はスキルスロット構成（最大4つ）
	Skills [MaxSkillSlotCount]SkillSlotConfig

	// ChainEffects はチェイン効果スロット構成（Skillsと同じインデックスで対応）
	ChainEffects [MaxSkillSlotCount]ChainEffectSlotConfig
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
// コア、すべてのスキル、すべてのチェイン効果が初期化されます。
func (s *AgentSlot) Clear() {
	s.CoreTypeID = ""
	for i := range s.Skills {
		s.Skills[i].Clear()
	}
	for i := range s.ChainEffects {
		s.ChainEffects[i].Clear()
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
func (s *AgentSlot) SetSkill(skillSlot int, typeID string) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.Skills[skillSlot].TypeID = typeID
}

// ClearSkill は指定スキルスロットをクリアします。
// 同インデックスのチェイン効果も自動的にクリアされます。
// skillSlotは0-3の範囲である必要があります。範囲外の場合は何もしません。
func (s *AgentSlot) ClearSkill(skillSlot int) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.Skills[skillSlot].Clear()
	s.ChainEffects[skillSlot].Clear()
}

// GetSkill は指定スキルスロットの構成を返します。
// skillSlotは0-3の範囲である必要があります。範囲外の場合はnilを返します。
func (s *AgentSlot) GetSkill(skillSlot int) *SkillSlotConfig {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return nil
	}
	return &s.Skills[skillSlot]
}

// SetChainEffect は指定スキルスロットにチェイン効果を設定します。
// skillSlotは0-3の範囲である必要があります。範囲外の場合は何もしません。
// 検証は行いません。AgentSlotManagerを通じて設定してください。
func (s *AgentSlot) SetChainEffect(skillSlot int, typeID string) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.ChainEffects[skillSlot].TypeID = typeID
}

// ClearChainEffect は指定スキルスロットのチェイン効果をクリアします。
// skillSlotは0-3の範囲である必要があります。範囲外の場合は何もしません。
func (s *AgentSlot) ClearChainEffect(skillSlot int) {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return
	}
	s.ChainEffects[skillSlot].Clear()
}

// GetChainEffect は指定スキルスロットのチェイン効果構成を返します。
// skillSlotは0-3の範囲である必要があります。範囲外の場合はnilを返します。
func (s *AgentSlot) GetChainEffect(skillSlot int) *ChainEffectSlotConfig {
	if skillSlot < 0 || skillSlot >= MaxSkillSlotCount {
		return nil
	}
	return &s.ChainEffects[skillSlot]
}
