// Package domain はゲームのドメインモデルを定義します。
package domain

// ModuleCategory はモジュールのカテゴリを表す型です。
// カテゴリによって効果の種類と参照するステータスが決まります。
type ModuleCategory string

const (
	// PhysicalAttack は物理攻撃カテゴリを表します。
	// STR参照、敵に物理ダメージを与えます。
	PhysicalAttack ModuleCategory = "physical_attack"

	// MagicAttack は魔法攻撃カテゴリを表します。
	// MAG参照、敵に魔法ダメージを与えます。
	MagicAttack ModuleCategory = "magic_attack"

	// Heal は回復カテゴリを表します。
	// MAG参照、プレイヤーのHPを回復します。
	Heal ModuleCategory = "heal"

	// Buff はバフカテゴリを表します。
	// SPD参照、プレイヤーに有利な効果を付与します。
	Buff ModuleCategory = "buff"

	// Debuff はデバフカテゴリを表します。
	// SPD参照、敵に不利な効果を付与します。
	Debuff ModuleCategory = "debuff"
)

// String はModuleCategoryの日本語表示名を返します。
func (c ModuleCategory) String() string {
	switch c {
	case PhysicalAttack:
		return "物理攻撃"
	case MagicAttack:
		return "魔法攻撃"
	case Heal:
		return "回復"
	case Buff:
		return "バフ"
	case Debuff:
		return "デバフ"
	default:
		return "不明"
	}
}

// Icon はモジュールカテゴリのアイコンを返します。
func (c ModuleCategory) Icon() string {
	switch c {
	case PhysicalAttack:
		return "⚔"
	case MagicAttack:
		return "✦"
	case Heal:
		return "♥"
	case Buff:
		return "▲"
	case Debuff:
		return "▼"
	default:
		return "•"
	}
}

// GetLevelSuffix はレベルに応じた接尾辞（low, mid, high）を返します。
// レベル1はlow、レベル2はmid、レベル3以上はhighです。
func GetLevelSuffix(level int) string {
	if level <= 1 {
		return "low"
	} else if level == 2 {
		return "mid"
	}
	return "high"
}

// ModuleModel はゲーム内のモジュール（スキル）エンティティを表す構造体です。
// モジュールはエージェント合成時にコアに装備され、バトル中に使用可能なスキルになります。
type ModuleModel struct {
	// ID はモジュールインスタンスの一意識別子です。
	ID string

	// Name はモジュールの表示名です（日本語）。
	Name string

	// Category はモジュールのカテゴリです（物理攻撃、魔法攻撃、回復、バフ、デバフ）。
	Category ModuleCategory

	// Level はモジュールのレベルです。
	Level int

	// Tags はモジュールのタグリストです。
	// コア特性との互換性チェックに使用されます。
	Tags []string

	// BaseEffect はモジュールの基礎効果値です。
	BaseEffect float64

	// StatRef は効果計算時に参照するステータスです（STR, MAG, SPD, LUK）。
	StatRef string

	// Description はモジュールの効果説明です。
	Description string
}

// NewModule は新しいModuleModelを作成します。
// Tagsはコピーされ、元のスライスとの参照共有を避けます。
func NewModule(id, name string, category ModuleCategory, level int, tags []string, baseEffect float64, statRef, description string) *ModuleModel {
	// Tagsをコピー（スライスの参照共有を避ける）
	tagsCopy := make([]string, len(tags))
	copy(tagsCopy, tags)

	return &ModuleModel{
		ID:          id,
		Name:        name,
		Category:    category,
		Level:       level,
		Tags:        tagsCopy,
		BaseEffect:  baseEffect,
		StatRef:     statRef,
		Description: description,
	}
}

// HasTag は指定されたタグがこのモジュールに含まれているかを返します。
func (m *ModuleModel) HasTag(tag string) bool {
	for _, t := range m.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// IsCompatibleWithCore はこのモジュールが指定されたコアに装備可能かを判定します。
// モジュールのタグのうち1つでもコアの許可タグに含まれていれば互換性ありとみなします。
func (m *ModuleModel) IsCompatibleWithCore(core *CoreModel) bool {
	for _, tag := range m.Tags {
		if core.IsTagAllowed(tag) {
			return true
		}
	}
	return false
}
