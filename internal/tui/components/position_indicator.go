// Package components は共通UIコンポーネントを提供します。
package components

import "strings"

// 位置インジケーターで使用する記号
const (
	IndicatorSelected   = "■"  // 選択中
	IndicatorUnselected = "□"  // 未選択
	IndicatorSpacing    = "  " // 横方向の間隔
)

// RenderHorizontalIndicator は横向きの位置インジケーターを描画します。
// 出力例（total=3, current=0）: "■  □  □"
// current は0から始まるインデックスです。
func RenderHorizontalIndicator(total, current int) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current >= total {
		current = total - 1
	}

	var parts []string
	for i := 0; i < total; i++ {
		if i == current {
			parts = append(parts, IndicatorSelected)
		} else {
			parts = append(parts, IndicatorUnselected)
		}
	}
	return strings.Join(parts, IndicatorSpacing)
}

// 敵選択インジケーターで使用する記号
const (
	EnemyIndicatorSelectedUndefeated   = "■" // 選択中・未撃破
	EnemyIndicatorUnselectedUndefeated = "□" // 未選択・未撃破
	EnemyIndicatorSelectedDefeated     = "★" // 選択中・撃破済み
	EnemyIndicatorUnselectedDefeated   = "☆" // 未選択・撃破済み
	EnemyIndicatorArrowLeft            = "◀" // 左矢印
	EnemyIndicatorArrowRight           = "▶" // 右矢印
)

// RenderEnemySelectionIndicator は敵選択インジケーターを描画します。
// 未撃破: ■（選択中）/ □（未選択）
// 撃破済み: ★（選択中）/ ☆（未選択）
// 矢印付きで表示: ◀  ■  □  ☆  ▶
//
// total: 敵の総数
// current: 現在選択中のインデックス（0から開始）
// defeatedFlags: 各敵の撃破済みフラグ（インデックス順）
func RenderEnemySelectionIndicator(total, current int, defeatedFlags []bool) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current >= total {
		current = total - 1
	}

	var parts []string
	parts = append(parts, EnemyIndicatorArrowLeft)

	for i := 0; i < total; i++ {
		isDefeated := i < len(defeatedFlags) && defeatedFlags[i]
		isSelected := i == current

		var symbol string
		if isSelected {
			if isDefeated {
				symbol = EnemyIndicatorSelectedDefeated
			} else {
				symbol = EnemyIndicatorSelectedUndefeated
			}
		} else {
			if isDefeated {
				symbol = EnemyIndicatorUnselectedDefeated
			} else {
				symbol = EnemyIndicatorUnselectedUndefeated
			}
		}
		parts = append(parts, symbol)
	}

	parts = append(parts, EnemyIndicatorArrowRight)
	return strings.Join(parts, IndicatorSpacing)
}

// RenderVerticalIndicator は縦向きの位置インジケーターを描画します。
// 出力例（total=3, current=1）:
//
//	□
//	■
//	□
//
// current は0から始まるインデックスです。
func RenderVerticalIndicator(total, current int) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current >= total {
		current = total - 1
	}

	var lines []string
	for i := 0; i < total; i++ {
		if i == current {
			lines = append(lines, IndicatorSelected)
		} else {
			lines = append(lines, IndicatorUnselected)
		}
	}
	return strings.Join(lines, "\n")
}
