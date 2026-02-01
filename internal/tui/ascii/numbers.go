// Package ascii はASCIIアート描画機能を提供します。
// 数字のASCIIアート描画を担当します。

package ascii

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// 数字のASCIIアート定義（各数字は5行）

var asciiDigits = map[int][]string{
	0: {
		"█████",
		"█   █",
		"█   █",
		"█   █",
		"█████",
	},
	1: {
		"  ██ ",
		" ███ ",
		"  ██ ",
		"  ██ ",
		"█████",
	},
	2: {
		"█████",
		"    █",
		"█████",
		"█    ",
		"█████",
	},
	3: {
		"█████",
		"    █",
		"█████",
		"    █",
		"█████",
	},
	4: {
		"█   █",
		"█   █",
		"█████",
		"    █",
		"    █",
	},
	5: {
		"█████",
		"█    ",
		"█████",
		"    █",
		"█████",
	},
	6: {
		"█████",
		"█    ",
		"█████",
		"█   █",
		"█████",
	},
	7: {
		"█████",
		"    █",
		"   █ ",
		"  █  ",
		"  █  ",
	},
	8: {
		"█████",
		"█   █",
		"█████",
		"█   █",
		"█████",
	},
	9: {
		"█████",
		"█   █",
		"█████",
		"    █",
		"█████",
	},
}

// "+"記号のASCIIアート（999+表示用）
var asciiPlus = []string{
	"     ",
	"  █  ",
	"█████",
	"  █  ",
	"     ",
}

// ASCIINumberRenderer は数字のASCIIアート描画を提供するインターフェースです。
type ASCIINumberRenderer interface {
	// RenderNumber は指定された数値をASCIIアートで描画します。
	// number: 描画する数値（0以上）
	// color: 数字の色
	RenderNumber(number int, color lipgloss.Color) string

	// RenderDigit は単一の数字（0-9）を描画します。
	// 範囲外の数字の場合はnilを返します。
	RenderDigit(digit int) []string
}

// asciiNumbers はASCIINumberRendererの実装です。
type asciiNumbers struct {
	digitWidth  int
	digitHeight int
}

// NewASCIINumbers は新しいASCIINumberRendererを作成します。
func NewASCIINumbers() ASCIINumberRenderer {
	// 数字の幅と高さを取得（全数字で同じ）
	digitHeight := len(asciiDigits[0])
	digitWidth := len([]rune(asciiDigits[0][0]))

	return &asciiNumbers{
		digitWidth:  digitWidth,
		digitHeight: digitHeight,
	}
}

// RenderDigit は単一の数字（0-9）を描画します。

func (n *asciiNumbers) RenderDigit(digit int) []string {
	// 範囲外の数字はnilを返す
	if digit < 0 || digit > 9 {
		return nil
	}

	// 元のスライスを返す（コピーは呼び出し側で必要に応じて行う）
	return asciiDigits[digit]
}

// RenderNumber は指定された数値をASCIIアートで描画します。

func (n *asciiNumbers) RenderNumber(number int, color lipgloss.Color) string {

	if number < 0 {
		number = 0
	}

	showPlus := false
	if number >= 1000 {
		number = 999
		showPlus = true
	}

	// 数字を桁ごとに分解
	digits := n.getDigits(number)

	// 各行を構築
	lines := make([]string, n.digitHeight)
	spacing := " " // 数字間のスペース

	for row := 0; row < n.digitHeight; row++ {
		var rowBuilder strings.Builder
		for i, digit := range digits {
			if i > 0 {
				rowBuilder.WriteString(spacing)
			}
			rowBuilder.WriteString(asciiDigits[digit][row])
		}

		// 999+表示の場合、+を追加
		if showPlus {
			rowBuilder.WriteString(spacing)
			rowBuilder.WriteString(asciiPlus[row])
		}

		lines[row] = rowBuilder.String()
	}

	// スタイルを適用して結合
	style := lipgloss.NewStyle().Foreground(color)
	result := strings.Join(lines, "\n")
	return style.Render(result)
}

// getDigits は数値を桁ごとの配列に変換します。
func (n *asciiNumbers) getDigits(number int) []int {
	if number == 0 {
		return []int{0}
	}

	var digits []int
	for number > 0 {
		digits = append([]int{number % 10}, digits...)
		number /= 10
	}
	return digits
}

// 影付き数字レンダリング用の定数
const (
	ShadowCharRight  = "░"                   // 右影の文字
	ShadowCharBottom = "░░░░░░"              // 下影の文字（数字幅+右影分）
	GoldColor        = lipgloss.Color("220") // 本体色（明るいゴールド）
	DarkGoldColor    = lipgloss.Color("136") // 影色（暗いゴールド）
)

// ShadowNumberRenderer は影付き数字のASCIIアート描画を提供するインターフェースです。
type ShadowNumberRenderer interface {
	// RenderShadowNumber は指定された数値を影付きASCIIアートで描画します。
	// number: 描画する数値（0以上）
	RenderShadowNumber(number int) string

	// RenderShadowNumberWithColors は色を指定して影付きASCIIアートで描画します。
	// number: 描画する数値（0以上）
	// bodyColor: 本体の色
	// shadowColor: 影の色
	RenderShadowNumberWithColors(number int, bodyColor, shadowColor lipgloss.Color) string
}

// shadowNumbers はShadowNumberRendererの実装です。
type shadowNumbers struct {
	digitWidth  int
	digitHeight int
}

// NewShadowNumbers は新しいShadowNumberRendererを作成します。
func NewShadowNumbers() ShadowNumberRenderer {
	digitHeight := len(asciiDigits[0])
	digitWidth := len([]rune(asciiDigits[0][0]))

	return &shadowNumbers{
		digitWidth:  digitWidth,
		digitHeight: digitHeight,
	}
}

// RenderShadowNumber は指定された数値を影付きASCIIアートで描画します。
// デフォルトのゴールド系カラーを使用します。
func (s *shadowNumbers) RenderShadowNumber(number int) string {
	return s.RenderShadowNumberWithColors(number, GoldColor, DarkGoldColor)
}

// RenderShadowNumberWithColors は色を指定して影付きASCIIアートで描画します。
func (s *shadowNumbers) RenderShadowNumberWithColors(number int, bodyColor, shadowColor lipgloss.Color) string {
	if number < 0 {
		number = 0
	}

	// 数字を桁ごとに分解
	digits := s.getDigits(number)

	// スタイル
	bodyStyle := lipgloss.NewStyle().Foreground(bodyColor)
	shadowStyle := lipgloss.NewStyle().Foreground(shadowColor)

	// 6行（5行の本体 + 1行の下影）
	lines := make([]string, s.digitHeight+1)
	digitSpacing := "  " // 数字間のスペース

	// 本体部分（5行）+ 右影
	for row := 0; row < s.digitHeight; row++ {
		var parts []string
		for _, digit := range digits {
			body := asciiDigits[digit][row]
			parts = append(parts, bodyStyle.Render(body)+shadowStyle.Render(ShadowCharRight))
		}
		lines[row] = strings.Join(parts, digitSpacing)
	}

	// 下影（最終行）
	var bottomParts []string
	for range digits {
		bottomParts = append(bottomParts, shadowStyle.Render(ShadowCharBottom))
	}
	lines[s.digitHeight] = strings.Join(bottomParts, digitSpacing)

	return strings.Join(lines, "\n")
}

// getDigits は数値を桁ごとの配列に変換します。
func (s *shadowNumbers) getDigits(number int) []int {
	if number == 0 {
		return []int{0}
	}

	var digits []int
	for number > 0 {
		digits = append([]int{number % 10}, digits...)
		number /= 10
	}
	return digits
}
