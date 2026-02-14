// Package domain はゲームのドメインモデルを定義します。
package domain

// defaultSkillIcon はスキルのデフォルトアイコンを返します。
const defaultSkillIcon = "•"

// SkillType はスキルの種別（タイプ）を定義する構造体です。
// 外部データファイル（skills.json）から読み込まれ、ゲーム内のスキル種別を定義します。
// 各スキルは複数の効果（Effects）を持ち、使用時に各効果が確率で発動します。
type SkillType struct {
	// ID はスキル種別の一意識別子です。
	ID string

	// Name はスキルの表示名です（日本語）。
	Name string

	// Icon はスキルのアイコン（絵文字）です。
	Icon string

	// Tags はスキルのタグリストです。
	// コア特性との互換性チェックに使用されます。
	Tags []string

	// Description はスキルの効果説明です。
	Description string

	// CooldownSeconds はスキルのクールダウン時間（秒）です。
	CooldownSeconds float64

	// DifficultyRate はタイピングの難易度です（50-200、100=標準）。
	DifficultyRate int

	// ChallengeType はこのスキルに対応するチャレンジタイプのIDです。
	ChallengeType ChallengeTypeID

	// ChallengeOptions はチャレンジ固有のオプションです。
	ChallengeOptions map[string]string

	// MinDropLevel はこのスキルがドロップする最低敵レベルです。
	MinDropLevel int

	// ManaCost はスキル使用時に消費するマナ量です。0の場合はマナ消費なし。
	ManaCost int

	// Effects はこのスキルが持つ効果のリストです。
	// 使用時に各効果が確率（Probability + LUK補正）で発動します。
	Effects []SkillEffect
}

// HasTag は指定されたタグがこのスキルタイプに含まれているかを返します。
func (t SkillType) HasTag(tag string) bool {
	for _, myTag := range t.Tags {
		if myTag == tag {
			return true
		}
	}
	return false
}

// IsCompatibleWithCore はこのスキルが指定されたコアに装備可能かを判定します。
// スキルのタグのうち1つでもコアの許可タグに含まれていれば互換性ありとみなします。
func (t SkillType) IsCompatibleWithCore(core *CoreModel) bool {
	for _, tag := range t.Tags {
		if core.IsTagAllowed(tag) {
			return true
		}
	}
	return false
}
