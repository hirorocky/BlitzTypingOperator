// Package domain はゲームのドメインモデルを定義します。
package domain

// EffectColumn は効果の種類（列）を定義します。
// 効果テーブルの列として、各効果の種別を表します。
type EffectColumn string

const (
	// ========== 攻撃強化系 ==========

	// ColDamageBonus は加算ダメージ（固定値）を表します。
	ColDamageBonus EffectColumn = "damage_bonus"

	// ColDamageMultiplier はダメージ倍率を表します。
	ColDamageMultiplier EffectColumn = "damage_mult"

	// ColArmorPierce は防御貫通（bool）を表します。
	ColArmorPierce EffectColumn = "armor_pierce"

	// ColLifeSteal はHP吸収率（%）を表します。
	ColLifeSteal EffectColumn = "life_steal"

	// ========== 防御強化系 ==========

	// ColDamageCut は被ダメ軽減（%）を表します。
	ColDamageCut EffectColumn = "damage_cut"

	// ColEvasion は回避率（%）を表します。
	ColEvasion EffectColumn = "evasion"

	// ColReflect は反射率（%）を表します。
	ColReflect EffectColumn = "reflect"

	// ColRegen は継続回復（/秒）を表します。
	ColRegen EffectColumn = "regen"

	// ========== 回復強化系 ==========

	// ColHealBonus は加算回復（固定値）を表します。
	ColHealBonus EffectColumn = "heal_bonus"

	// ColHealMultiplier は回復倍率を表します。
	ColHealMultiplier EffectColumn = "heal_mult"

	// ColOverheal は超過回復→一時HP（bool）を表します。
	ColOverheal EffectColumn = "overheal"

	// ========== タイピング系 ==========

	// ColTimeExtend は時間延長（秒）を表します。
	ColTimeExtend EffectColumn = "time_extend"

	// ColAutoCorrect はミス無視回数を表します。
	ColAutoCorrect EffectColumn = "auto_correct"

	// ColTypingDifficulty はタイピング難易度の乗算修正値を表します。
	ColTypingDifficulty EffectColumn = "typing_difficulty"

	// ========== リキャスト系 ==========

	// ColCooldownReduce はCD短縮率（%）を表します。
	ColCooldownReduce EffectColumn = "cooldown_reduce"

	// ========== バフ/デバフ延長系 ==========

	// ColBuffExtend はバフ延長（秒）を表します。
	ColBuffExtend EffectColumn = "buff_extend"

	// ColDebuffExtend はデバフ延長（秒）を表します。
	ColDebuffExtend EffectColumn = "debuff_extend"

	// ========== 特殊系 ==========

	// ColDoubleCast は2回発動確率（%）を表します。
	ColDoubleCast EffectColumn = "double_cast"

	// ========== ステータス系 ==========

	// ColCritRate はクリティカル率加算（%）を表します。
	ColCritRate EffectColumn = "crit_rate"

	// ColSTRBonus はSTR加算（固定値）を表します。
	ColSTRBonus EffectColumn = "str_bonus"

	// ColINTBonus はINT加算（固定値）を表します。
	ColINTBonus EffectColumn = "int_bonus"

	// ColWILBonus はWIL加算（固定値）を表します。
	ColWILBonus EffectColumn = "wil_bonus"

	// ColLUKBonus はLUK加算（固定値）を表します。
	ColLUKBonus EffectColumn = "luk_bonus"

	// ColSTRMultiplier はSTR倍率を表します。
	ColSTRMultiplier EffectColumn = "str_mult"

	// ColINTMultiplier はINT倍率を表します。
	ColINTMultiplier EffectColumn = "int_mult"

	// ColWILMultiplier はWIL倍率を表します。
	ColWILMultiplier EffectColumn = "wil_mult"

	// ColLUKMultiplier はLUK倍率を表します。
	ColLUKMultiplier EffectColumn = "luk_mult"
)

// AggregationType は列ごとの集計方法を表します。
type AggregationType int

const (
	// AggAdd は加算集計を表します: Σ(values)
	AggAdd AggregationType = iota

	// AggMult は乗算集計を表します: Π(values)
	AggMult

	// AggMax は最大値集計を表します: max(values)
	AggMax

	// AggOr はOR集計を表します: any(values)
	AggOr
)

// ColumnAggregation は各列の集計方法を定義します。
var ColumnAggregation = map[EffectColumn]AggregationType{
	// 攻撃強化系
	ColDamageBonus:      AggAdd,  // 固定ダメージは加算
	ColDamageMultiplier: AggMult, // 倍率は乗算
	ColArmorPierce:      AggOr,   // 1つでも true なら貫通
	ColLifeSteal:        AggMax,  // 最大の吸収率を採用

	// 防御強化系
	ColDamageCut: AggMax, // 最大の軽減率を採用
	ColEvasion:   AggMax, // 最大の回避率を採用
	ColReflect:   AggMax, // 最大の反射率を採用
	ColRegen:     AggAdd, // 回復量は加算

	// 回復強化系
	ColHealBonus:      AggAdd,  // 固定回復は加算
	ColHealMultiplier: AggMult, // 倍率は乗算
	ColOverheal:       AggOr,   // 1つでも true なら有効

	// タイピング系
	ColTimeExtend:       AggAdd,  // 延長時間は加算
	ColAutoCorrect:      AggAdd,  // ミス無視回数は加算
	ColTypingDifficulty: AggMult, // 難易度倍率は乗算

	// リキャスト系
	ColCooldownReduce: AggMax, // 最大の短縮率を採用

	// バフ/デバフ延長系
	ColBuffExtend:   AggAdd, // 延長秒数は加算
	ColDebuffExtend: AggAdd, // 延長秒数は加算

	// 特殊系
	ColDoubleCast: AggMax, // 最大の確率を採用

	// ステータス系
	ColCritRate:      AggAdd,  // クリティカル率は加算
	ColSTRBonus:      AggAdd,  // STR加算は加算
	ColINTBonus:      AggAdd,  // INT加算は加算
	ColWILBonus:      AggAdd,  // WIL加算は加算
	ColLUKBonus:      AggAdd,  // LUK加算は加算
	ColSTRMultiplier: AggMult, // STR倍率は乗算
	ColINTMultiplier: AggMult, // INT倍率は乗算
	ColWILMultiplier: AggMult, // WIL倍率は乗算
	ColLUKMultiplier: AggMult, // LUK倍率は乗算
}

// ColumnDefault は集計の初期値を返します。
func ColumnDefault(col EffectColumn) float64 {
	switch ColumnAggregation[col] {
	case AggMult:
		return 1.0 // 乗算の初期値は1
	default:
		return 0.0 // 加算/最大値/ORの初期値は0
	}
}
