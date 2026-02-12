// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"fmt"
	"sort"
	"strings"
)

// DescribeSingleEffect はEffectColumnとvalueから効果説明を生成します。
// 敵の行動予測表示などで使用されます。
func DescribeSingleEffect(col EffectColumn, value float64) string {
	desc := describeColumn(col, value)
	if desc == "" {
		return "効果"
	}
	return desc
}

// DescribeEffectValues は効果マップから説明文を生成します。
// 複数の効果がある場合は " / " で区切られます。
func DescribeEffectValues(values map[EffectColumn]float64) string {
	if len(values) == 0 {
		return "効果"
	}

	// 効果列をソートして安定した出力順序を保証
	columns := make([]EffectColumn, 0, len(values))
	for col := range values {
		columns = append(columns, col)
	}
	sort.Slice(columns, func(i, j int) bool {
		return string(columns[i]) < string(columns[j])
	})

	descriptions := make([]string, 0, len(columns))
	for _, col := range columns {
		value := values[col]
		desc := describeColumn(col, value)
		if desc != "" {
			descriptions = append(descriptions, desc)
		}
	}

	if len(descriptions) == 0 {
		return "効果"
	}

	return strings.Join(descriptions, " / ")
}

// describeColumn は単一の効果列から説明を生成します。
func describeColumn(col EffectColumn, value float64) string {
	switch col {
	// 攻撃強化系
	case ColDamageBonus:
		if value >= 0 {
			return fmt.Sprintf("ダメージ+%.0f", value)
		}
		return fmt.Sprintf("ダメージ%.0f", value)
	case ColDamageMultiplier:
		return formatMultiplier("ダメージ", value)
	case ColArmorPierce:
		if value > 0 {
			return "防御貫通"
		}
		return ""
	case ColLifeSteal:
		return fmt.Sprintf("HP吸収%.0f%%", value*100)

	// 防御強化系
	case ColDamageCut:
		return formatDamageCutEffect(value)
	case ColEvasion:
		return fmt.Sprintf("回避%.0f%%", value*100)
	case ColReflect:
		return fmt.Sprintf("反射%.0f%%", value*100)
	case ColRegen:
		return fmt.Sprintf("HP回復%.0f/s", value)

	// 回復強化系
	case ColHealBonus:
		if value >= 0 {
			return fmt.Sprintf("回復+%.0f", value)
		}
		return fmt.Sprintf("回復%.0f", value)
	case ColHealMultiplier:
		return formatMultiplier("回復", value)
	case ColOverheal:
		if value > 0 {
			return "超過回復"
		}
		return ""

	// タイピング系
	case ColTimeExtend:
		if value >= 0 {
			return fmt.Sprintf("時間+%.1f秒", value)
		}
		return fmt.Sprintf("時間%.1f秒", value)
	case ColAutoCorrect:
		return fmt.Sprintf("ミス無視%.0f回", value)
	case ColTypingDifficulty:
		return formatMultiplier("難易度", value)

	// リキャスト系
	case ColCooldownReduce:
		return formatCooldownEffect(value)

	// バフ/デバフ延長系
	case ColBuffExtend:
		return fmt.Sprintf("バフ延長+%.1f秒", value)
	case ColDebuffExtend:
		return fmt.Sprintf("デバフ延長+%.1f秒", value)

	// 特殊系
	case ColDoubleCast:
		return fmt.Sprintf("2回発動%.0f%%", value*100)

	// ステータス系
	case ColCritRate:
		return fmt.Sprintf("クリ率+%.0f%%", value*100)
	case ColSTRBonus:
		return fmt.Sprintf("STR+%.0f", value)
	case ColINTBonus:
		return fmt.Sprintf("INT+%.0f", value)
	case ColWILBonus:
		return fmt.Sprintf("WIL+%.0f", value)
	case ColLUKBonus:
		return fmt.Sprintf("LUK+%.0f", value)
	case ColSTRMultiplier:
		return formatIncrementRate("STR", value)
	case ColINTMultiplier:
		return formatIncrementRate("INT", value)
	case ColWILMultiplier:
		return formatIncrementRate("WIL", value)
	case ColLUKMultiplier:
		return formatIncrementRate("LUK", value)

	default:
		return ""
	}
}

// formatMultiplier は倍率を「XX%UP」または「XX%ダウン」形式に変換します。
func formatMultiplier(statName string, value float64) string {
	if value >= 1.0 {
		percent := (value - 1.0) * 100
		return fmt.Sprintf("%s%.0f%%UP", statName, percent)
	}
	percent := (1.0 - value) * 100
	return fmt.Sprintf("%s%.0f%%ダウン", statName, percent)
}

// formatIncrementRate は増加率（0.25 = +25%）を「XX%UP」または「XX%ダウン」形式に変換します。
// ステータス倍率系（STR/INT/WIL/LUK Multiplier）で使用されます。
func formatIncrementRate(statName string, value float64) string {
	if value >= 0 {
		return fmt.Sprintf("%s%.0f%%UP", statName, value*100)
	}
	return fmt.Sprintf("%s%.0f%%ダウン", statName, -value*100)
}

// formatCooldownEffect はクールダウン効果を説明に変換します。
func formatCooldownEffect(value float64) string {
	if value > 0 {
		return fmt.Sprintf("CD%.0f%%短縮", value*100)
	}
	return fmt.Sprintf("CD%.0f%%延長", -value*100)
}

// formatDamageCutEffect は被ダメ軽減/増加を説明に変換します。
func formatDamageCutEffect(value float64) string {
	if value >= 0 {
		return fmt.Sprintf("被ダメ%.0f%%軽減", value*100)
	}
	return fmt.Sprintf("被ダメ%.0f%%増加", -value*100)
}
