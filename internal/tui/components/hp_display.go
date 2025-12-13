// Package components はTUI共通コンポーネントを提供します。
// hp_display.go はHP表示ロジックの共通化を担当します。
// 要件 7.1, 7.2: HP表示コンポーネント
package components

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// HPバー色分けの閾値（styles.goと同じ値）
const (
	// hpHighThreshold はHP高（緑）の閾値（50%以上）
	hpHighThreshold = 0.50
	// hpMediumThreshold はHP中（黄）の閾値（25%以上）
	hpMediumThreshold = 0.25
)

// RenderHP はHP表示を共通フォーマットでレンダリングします。
// 要件 7.1: HPバー、数値表示、色分けロジックを共通化する
//
// パラメータ:
//   - current: 現在のHP値
//   - max: 最大HP値
//   - barWidth: HPバーの幅（文字数）
//   - gs: ゲームスタイル（nilの場合は新規作成）
//
// 戻り値: レンダリングされたHP表示文字列
func RenderHP(current, max int, barWidth int, gs *styles.GameStyles) string {
	// nilチェック
	if gs == nil {
		gs = styles.NewGameStyles()
	}

	// 境界値の正規化
	if max <= 0 {
		max = 1 // ゼロ除算を避ける
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	// HP割合を計算
	percentage := float64(current) / float64(max)

	// HPバーを構築
	hpBar := renderHPBar(percentage, barWidth, gs)

	return hpBar
}

// RenderHPWithLabel はラベル付きHP表示をレンダリングします。
// 要件 7.2: ラベル付き表示をサポートする
//
// パラメータ:
//   - label: 表示ラベル（例: "HP", "プレイヤーHP"）
//   - current: 現在のHP値
//   - max: 最大HP値
//   - barWidth: HPバーの幅（文字数）
//   - gs: ゲームスタイル（nilの場合は新規作成）
//
// 戻り値: レンダリングされたラベル付きHP表示文字列
func RenderHPWithLabel(label string, current, max int, barWidth int, gs *styles.GameStyles) string {
	// nilチェック
	if gs == nil {
		gs = styles.NewGameStyles()
	}

	// 境界値の正規化
	if max <= 0 {
		max = 1 // ゼロ除算を避ける
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	// HP割合を計算
	percentage := float64(current) / float64(max)

	// HPバーを構築
	hpBar := renderHPBar(percentage, barWidth, gs)

	// HP値を構築
	hpValue := fmt.Sprintf(" %d/%d", current, max)

	// ラベル付きで結合
	var result strings.Builder
	if label != "" {
		result.WriteString(label)
		result.WriteString(": ")
	}
	result.WriteString(hpBar)
	result.WriteString(hpValue)

	return result.String()
}

// renderHPBar はHPバーを描画する内部関数です。
func renderHPBar(percentage float64, barWidth int, gs *styles.GameStyles) string {
	// バー幅の正規化
	if barWidth < 2 {
		barWidth = 2
	}

	// バー内部の幅（ボーダー分を除く）
	innerWidth := barWidth - 2
	if innerWidth < 1 {
		innerWidth = 1
	}

	// 塗りつぶし部分の幅を計算
	filledWidth := int(float64(innerWidth) * percentage)
	if filledWidth < 0 {
		filledWidth = 0
	}
	if filledWidth > innerWidth {
		filledWidth = innerWidth
	}
	emptyWidth := innerWidth - filledWidth

	// HP割合に応じた色を選択
	fillColor := getHPColor(percentage)

	// バーを構築
	var bar strings.Builder
	bar.WriteString("[")

	// 塗りつぶし部分と空白部分を描画
	filledStyle := lipgloss.NewStyle().Background(fillColor)
	emptyStyle := lipgloss.NewStyle().Background(styles.ColorSubtle)
	bar.WriteString(filledStyle.Render(strings.Repeat(" ", filledWidth)))
	bar.WriteString(emptyStyle.Render(strings.Repeat(" ", emptyWidth)))

	bar.WriteString("]")

	return bar.String()
}

// getHPColor はHP割合に応じた色を返します。
func getHPColor(percentage float64) lipgloss.Color {
	if percentage > hpHighThreshold {
		return styles.ColorHPHigh // 緑（50%以上）
	} else if percentage > hpMediumThreshold {
		return styles.ColorHPMedium // 黄（25%以上50%未満）
	}
	return styles.ColorHPLow // 赤（25%未満）
}
