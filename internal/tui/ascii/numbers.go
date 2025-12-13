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
