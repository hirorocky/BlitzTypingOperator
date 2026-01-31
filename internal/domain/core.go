// Package domain はゲームのドメインモデルを定義します。
// コア、モジュール、エージェント、敵、プレイヤーなどのエンティティとそのビジネスルールを含みます。
package domain

// BaseStatValue はステータス計算で使用する基礎値です。
// ステータス = 基礎値 × 重み
const BaseStatValue = 100

// Stats はゲーム内のステータス値を表す構造体です。
// 各ステータスはコアのレベルと特性の重みによって計算されます。
type Stats struct {
	// STR は物理攻撃力を表します。
	// 物理攻撃モジュールのダメージ計算に使用されます。
	STR int

	// INT は魔法攻撃力を表します。
	// 攻撃魔法モジュールとデバフモジュールの効果計算に使用されます。
	INT int

	// WIL は意志力を表します。
	// 回復魔法モジュールとバフモジュールの効果計算に使用されます。
	WIL int

	// LUK は運を表します。
	// 確率系効果の発動率に影響します。
	// コアレベルでは変化せず、stat_weightsの影響のみ受けます。
	LUK int
}

// Total はステータスの合計値を返します。
func (s Stats) Total() int {
	return s.STR + s.INT + s.WIL + s.LUK
}

// CoreType はコアの特性（タイプ）を定義する構造体です。
// 外部データファイル（cores.json）から読み込まれ、ゲーム内のコア特性を定義します。
// 例: 攻撃バランス、パラディン、オールラウンダー、ヒーラー
type CoreType struct {
	// ID はコア特性の一意識別子です。
	ID string

	// Name はコア特性の表示名です（日本語）。
	Name string

	// StatWeights はステータス計算に使用する重みのマップです。
	// キーは "STR", "INT", "WIL", "LUK" で、値は重み係数（例: 1.2）です。
	StatWeights map[string]float64

	// PassiveSkillID はこのコア特性に紐づくパッシブスキルのIDです。
	PassiveSkillID string

	// AllowedTags はこのコア特性に装備可能なモジュールタグのリストです。
	// 例: ["physical_low", "magic_low"]
	AllowedTags []string

	// MinDropLevel はこのコア特性がドロップする最低敵レベルです。
	// このレベル未満の敵からはこの特性のコアはドロップしません。
	MinDropLevel int
}

// CoreModel はゲーム内のコアエンティティを表す構造体です。
// コアはエージェント合成時の中核となる素材で、特性と固定ステータスを持ちます。
// TypeIDのみで同一性が判定されます。
type CoreModel struct {
	// ID はコアインスタンスの一意識別子です。
	// 後方互換性のために残されています。新規コードではTypeIDを使用してください。
	ID string

	// TypeID はコア特性ID（マスタデータ参照用）です。
	// セーブデータにはTypeIDのみが保存されます。
	TypeID string

	// Name はコアの表示名です。
	Name string

	// Type はコアの特性（タイプ）です。
	Type CoreType

	// Stats はコアのステータス値です。
	// 特性の重みから計算されます（100 × 重み）。
	Stats Stats

	// PassiveSkill はこのコアに紐づくパッシブスキルです。
	PassiveSkill PassiveSkill

	// AllowedTags はこのコアに装備可能なモジュールタグのリストです。
	// 通常はType.AllowedTagsと同じですが、直接参照用にコピーされます。
	AllowedTags []string
}

// Equals はコアの同一性を判定します。
// TypeIDが同じ場合に等価とみなします。
func (c *CoreModel) Equals(other *CoreModel) bool {
	if other == nil {
		return false
	}
	return c.TypeID == other.TypeID
}

// CalculateStats はコア特性からステータス値を計算します。
// 各ステータス: 基礎値(100) × ステータス重み
// 結果は整数に切り捨てられます。
func CalculateStats(coreType CoreType) Stats {
	// 各ステータスの重みを取得（未設定の場合は0.0扱い）
	strWeight := coreType.StatWeights["STR"]
	intWeight := coreType.StatWeights["INT"]
	wilWeight := coreType.StatWeights["WIL"]
	lukWeight := coreType.StatWeights["LUK"]

	// 計算式: 基礎値(100) × 重み
	baseValue := float64(BaseStatValue)

	return Stats{
		STR: int(baseValue * strWeight),
		INT: int(baseValue * intWeight),
		WIL: int(baseValue * wilWeight),
		LUK: int(baseValue * lukWeight),
	}
}

// IsTagAllowed は指定されたタグがこのコアに許可されているかを判定します。
// モジュール装備時の互換性チェックに使用されます。
func (c *CoreModel) IsTagAllowed(tag string) bool {
	for _, allowedTag := range c.AllowedTags {
		if allowedTag == tag {
			return true
		}
	}
	return false
}

// NewCoreWithTypeID はTypeIDベースでCoreModelを作成します。
// ステータスはコア特性の重みから自動計算されます（100 × 重み）。
// NameはType.Name形式で自動生成されます。
// IDは後方互換性のためにTypeIDと同じ値が設定されますが、新規コードでは使用しないでください。
func NewCoreWithTypeID(typeID string, coreType CoreType, passiveSkill PassiveSkill) *CoreModel {
	// ステータスを重みベースで自動計算
	stats := CalculateStats(coreType)

	// AllowedTagsをコピー（スライスの参照共有を避ける）
	allowedTags := make([]string, len(coreType.AllowedTags))
	copy(allowedTags, coreType.AllowedTags)

	// Nameを自動生成（コア特性名のみ）
	name := coreType.Name

	return &CoreModel{
		ID:           typeID, // 後方互換性のため
		TypeID:       typeID,
		Name:         name,
		Type:         coreType,
		Stats:        stats,
		PassiveSkill: passiveSkill,
		AllowedTags:  allowedTags,
	}
}
