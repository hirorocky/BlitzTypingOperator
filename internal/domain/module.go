// Package domain はゲームのドメインモデルを定義します。
package domain

// SkillModel はゲーム内のスキルエンティティを表す構造体です。
// スキルはエージェント合成時にコアに装備され、バトル中に使用可能になります。
// TypeIDとChainEffectの組み合わせでインスタンスを識別します。
type SkillModel struct {
	// TypeID はスキル種別ID（マスタデータ参照用）です。
	// セーブデータにはTypeIDとChainEffectが保存されます。
	TypeID string

	// Type はスキルの種別（タイプ）です。
	// TypeIDからマスタデータを参照して取得されます。
	Type SkillType

	// ChainEffect はこのスキルインスタンスのチェイン効果です。
	// nilの場合はチェイン効果を持たないスキルです。
	ChainEffect *ChainEffect
}

// Name はスキルの表示名を返します。
func (m *SkillModel) Name() string {
	return m.Type.Name
}

// Tags はスキルのタグリストを返します。
func (m *SkillModel) Tags() []string {
	return m.Type.Tags
}

// Description はスキルの効果説明を返します。
func (m *SkillModel) Description() string {
	return m.Type.Description
}

// Icon はスキルのアイコンを返します。
func (m *SkillModel) Icon() string {
	if m.Type.Icon != "" {
		return m.Type.Icon
	}
	return defaultSkillIcon
}

// Effects はスキルの効果リストを返します。
func (m *SkillModel) Effects() []SkillEffect {
	return m.Type.Effects
}

// CooldownSeconds はスキルのクールダウン時間を返します。
func (m *SkillModel) CooldownSeconds() float64 {
	return m.Type.CooldownSeconds
}

// GetDifficultyRate はタイピングの難易度を返します。
func (m *SkillModel) GetDifficultyRate() DifficultyRate {
	return DifficultyRate(m.Type.DifficultyRate)
}

// GetChallengeType はこのスキルに対応するチャレンジタイプIDを返します。
func (m *SkillModel) GetChallengeType() ChallengeTypeID {
	return m.Type.ChallengeType
}

// HasTag は指定されたタグがこのスキルに含まれているかを返します。
func (m *SkillModel) HasTag(tag string) bool {
	return m.Type.HasTag(tag)
}

// IsCompatibleWithCore はこのスキルが指定されたコアに装備可能かを判定します。
// スキルのタグのうち1つでもコアの許可タグに含まれていれば互換性ありとみなします。
func (m *SkillModel) IsCompatibleWithCore(core *CoreModel) bool {
	for _, tag := range m.Type.Tags {
		if core.IsTagAllowed(tag) {
			return true
		}
	}
	return false
}

// HasChainEffect はこのスキルがチェイン効果を持っているかを返します。
func (m *SkillModel) HasChainEffect() bool {
	return m.ChainEffect != nil
}

// NewSkillFromType はSkillTypeからSkillModelを作成します。
// chainEffectはnilを許容します。
func NewSkillFromType(skillType SkillType, chainEffect *ChainEffect) *SkillModel {
	return &SkillModel{
		TypeID:      skillType.ID,
		Type:        skillType,
		ChainEffect: chainEffect,
	}
}
