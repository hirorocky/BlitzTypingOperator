// Package components は共通UIコンポーネントを提供します。
package components

import (
	"fmt"

	"hirorocky/type-battle/internal/domain"
)

// FormatStats はコアタイプのステータスを固定順序（STR, INT, WIL, LUK）でフォーマットします。
// 出力例: "STR: 120  INT: 100  WIL: 80  LUK: 100"
func FormatStats(coreType domain.CoreType) string {
	stats := domain.CalculateStats(coreType)
	return fmt.Sprintf("STR: %d  INT: %d  WIL: %d  LUK: %d",
		stats.STR, stats.INT, stats.WIL, stats.LUK)
}

// FormatStatsFromStats はStats構造体から直接フォーマットします。
// すでに計算済みのStatsがある場合に使用します。
func FormatStatsFromStats(stats domain.Stats) string {
	return fmt.Sprintf("STR: %d  INT: %d  WIL: %d  LUK: %d",
		stats.STR, stats.INT, stats.WIL, stats.LUK)
}

// FormatStatsMultiline はステータスを複数行でフォーマットします。
// 出力例:
//
//	STR: 120
//	INT: 100
//	WIL:  80
//	LUK: 100
func FormatStatsMultiline(coreType domain.CoreType) string {
	stats := domain.CalculateStats(coreType)
	return fmt.Sprintf("STR: %3d\nINT: %3d\nWIL: %3d\nLUK: %3d",
		stats.STR, stats.INT, stats.WIL, stats.LUK)
}
