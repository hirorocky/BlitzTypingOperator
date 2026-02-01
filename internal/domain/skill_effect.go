// Package domain はゲームのドメインモデルを定義します。
package domain

import "math/rand"

// EffectTarget はスキル効果の対象を表します。
type EffectTarget string

const (
	// TargetSelf は自分自身を対象とする効果です。
	TargetSelf EffectTarget = "self"

	// TargetEnemy は敵を対象とする効果です。
	TargetEnemy EffectTarget = "enemy"

	// TargetBoth は自分と敵の両方を対象とする効果です。
	TargetBoth EffectTarget = "both"
)

// HPFormula はHP増減の計算式を表します。
// 計算式: Base + (StatCoef × ステータス値)
type HPFormula struct {
	// Base は固定値です。正の値は回復、負の値はダメージ/自傷を表します。
	Base float64

	// StatCoef はステータス係数です。
	StatCoef float64

	// StatRef は参照するステータスです（STR, INT, WIL）。
	StatRef string
}

// EffectColumnSpec はEffectColumn効果の仕様を表します。
// バフ/デバフなどの継続効果に使用されます。
type EffectColumnSpec struct {
	// Column は効果列です。
	Column EffectColumn

	// Value は効果値です。
	Value float64

	// Duration は持続時間（秒）です。0の場合は即時効果です。
	Duration float64
}

// BaseLUK はLUK補正の基準値です。
const BaseLUK = 10

// SkillEffect はスキルの1つの効果を表します。
// 各スキルは複数のSkillEffectを持つことができます。
type SkillEffect struct {
	// Target は効果の対象です。
	Target EffectTarget

	// HPFormula はHP増減の計算式です。nilの場合はHP増減なし。
	HPFormula *HPFormula

	// ColumnSpec はEffectColumn効果の仕様です。nilの場合はカラム効果なし。
	ColumnSpec *EffectColumnSpec

	// Probability はベース発動確率（0.0-1.0）です。
	Probability float64

	// LUKFactor はLUKが確率に与える影響です。
	// 正の値: LUKが高いほど確率上昇
	// 負の値: LUKが高いほど確率下降
	// 0: LUKの影響なし
	LUKFactor float64

	// Icon は表示用アイコンです。
	Icon string
}

// CalculateHPChange はステータス値からHP変化量を計算します。
// 正の値は回復、負の値はダメージを表します。
func (e *SkillEffect) CalculateHPChange(stats Stats) int {
	if e.HPFormula == nil {
		return 0
	}

	var statValue int
	switch e.HPFormula.StatRef {
	case "STR":
		statValue = stats.STR
	case "INT":
		statValue = stats.INT
	case "WIL":
		statValue = stats.WIL
	default:
		statValue = 0
	}

	return int(e.HPFormula.Base + e.HPFormula.StatCoef*float64(statValue))
}

// AdjustedProbability はLUK補正を適用した発動確率を計算します。
// 計算式: Probability + (LUK - 10) × LUKFactor
// 結果は0.0-1.0の範囲にクランプされます。
func (e *SkillEffect) AdjustedProbability(luk int) float64 {
	adjustedProb := e.Probability + float64(luk-BaseLUK)*e.LUKFactor

	// 0.0-1.0の範囲にクランプ
	if adjustedProb < 0.0 {
		return 0.0
	}
	if adjustedProb > 1.0 {
		return 1.0
	}
	return adjustedProb
}

// ShouldTrigger はLUK補正を考慮して効果が発動するかを判定します。
func (e *SkillEffect) ShouldTrigger(luk int, rng *rand.Rand) bool {
	adjustedProb := e.AdjustedProbability(luk)
	if adjustedProb >= 1.0 {
		return true
	}
	if adjustedProb <= 0.0 {
		return false
	}
	return rng.Float64() < adjustedProb
}

// IsHPEffect はHP増減効果を持つかを判定します。
func (e *SkillEffect) IsHPEffect() bool {
	return e.HPFormula != nil
}

// IsColumnEffect はEffectColumn効果を持つかを判定します。
func (e *SkillEffect) IsColumnEffect() bool {
	return e.ColumnSpec != nil
}

// IsDamageEffect は敵へのダメージ効果かを判定します。
func (e *SkillEffect) IsDamageEffect() bool {
	if e.HPFormula == nil {
		return false
	}
	// 敵対象かつHP減少（係数が正の場合、敵のHPを減らす = ダメージ）
	return e.Target == TargetEnemy && (e.HPFormula.Base > 0 || e.HPFormula.StatCoef > 0)
}

// IsHealEffect は自分への回復効果かを判定します。
func (e *SkillEffect) IsHealEffect() bool {
	if e.HPFormula == nil {
		return false
	}
	// 自分対象かつHP増加
	return e.Target == TargetSelf && (e.HPFormula.Base > 0 || e.HPFormula.StatCoef > 0)
}

// IsBuffEffect は自分へのバフ効果かを判定します。
func (e *SkillEffect) IsBuffEffect() bool {
	if e.ColumnSpec == nil {
		return false
	}
	// 自分対象かつカラム効果（バフ = 自分へのステータス強化）
	return e.Target == TargetSelf
}

// IsDebuffEffect は敵へのデバフ効果かを判定します。
func (e *SkillEffect) IsDebuffEffect() bool {
	if e.ColumnSpec == nil {
		return false
	}
	// 敵対象かつカラム効果（デバフ = 敵へのステータス弱化）
	return e.Target == TargetEnemy
}
