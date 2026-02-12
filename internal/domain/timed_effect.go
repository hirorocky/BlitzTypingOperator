// Package domain はゲームのドメインモデルを定義します。
package domain

// TimedEffect は時限効果（バフ/デバフ）の定義を表す値オブジェクトです。
// 重複判定のIDおよび効果列・効果値の定義として使用されます。
// パッシブ効果・チェイン効果と並ぶ第三の効果種別です。
type TimedEffect struct {
	// ID は時限効果の一意識別子です（例: "st_str_buff_lv1"）。
	ID string

	// Name は時限効果の表示名です。
	Name string

	// Description は時限効果の説明文です。
	Description string

	// Column は効果列です（例: ColSTRMultiplier, ColDamageMultiplier）。
	Column EffectColumn

	// Value は効果値です。
	Value float64
}
