// Package ascii はASCIIアート描画機能を提供します。

package ascii

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// TestNewASCIINumbers はASCIINumberRendererの作成をテストします。
func TestNewASCIINumbers(t *testing.T) {
	renderer := NewASCIINumbers()
	if renderer == nil {
		t.Error("NewASCIINumbers()がnilを返しました")
	}
}

// TestASCIINumbersRenderDigit は単一の数字（0-9）のレンダリングをテストします。

func TestASCIINumbersRenderDigit(t *testing.T) {
	renderer := NewASCIINumbers()

	// 全ての数字（0-9）をテスト
	for digit := 0; digit <= 9; digit++ {
		result := renderer.RenderDigit(digit)
		if len(result) == 0 {
			t.Errorf("RenderDigit(%d)が空の結果を返しました", digit)
			continue
		}

		// 3-5行のASCIIアートであることを確認
		if len(result) < 3 || len(result) > 5 {
			t.Errorf("RenderDigit(%d)の行数が想定外です: %d行（3-5行を想定）", digit, len(result))
		}
	}
}

// TestASCIINumbersRenderDigitInvalid は無効な数字の処理をテストします。
func TestASCIINumbersRenderDigitInvalid(t *testing.T) {
	renderer := NewASCIINumbers()

	// 範囲外の数字はnilまたは空を返す
	resultNeg := renderer.RenderDigit(-1)
	if resultNeg != nil {
		t.Error("RenderDigit(-1)がnilではない値を返しました")
	}

	result10 := renderer.RenderDigit(10)
	if result10 != nil {
		t.Error("RenderDigit(10)がnilではない値を返しました")
	}
}

// TestASCIINumbersRenderNumber は複数桁の数値レンダリングをテストします。

func TestASCIINumbersRenderNumber(t *testing.T) {
	renderer := NewASCIINumbers()

	tests := []struct {
		number   int
		expected bool // 結果が空でないことを確認
	}{
		{0, true},
		{5, true},
		{12, true},
		{123, true},
		{999, true},
	}

	for _, tt := range tests {
		result := renderer.RenderNumber(tt.number, lipgloss.Color("#FFFFFF"))
		if result == "" && tt.expected {
			t.Errorf("RenderNumber(%d)が空文字列を返しました", tt.number)
		}

		// 出力が複数行であることを確認
		lines := strings.Split(result, "\n")
		if len(lines) < 3 {
			t.Errorf("RenderNumber(%d)の行数が少なすぎます: %d行", tt.number, len(lines))
		}
	}
}

// TestASCIINumbersRenderNumberNegative は負数の処理をテストします。

func TestASCIINumbersRenderNumberNegative(t *testing.T) {
	renderer := NewASCIINumbers()

	result := renderer.RenderNumber(-5, lipgloss.Color("#FFFFFF"))
	expected := renderer.RenderNumber(0, lipgloss.Color("#FFFFFF"))

	// 負数は0として扱われることを確認
	if result != expected {
		t.Error("負数が0として扱われていません")
	}
}

// TestASCIINumbersRenderNumberLarge は大きな数値の処理をテストします。

func TestASCIINumbersRenderNumberLarge(t *testing.T) {
	renderer := NewASCIINumbers()

	result := renderer.RenderNumber(1000, lipgloss.Color("#FFFFFF"))
	// 999+の表示が含まれていることを確認
	if result == "" {
		t.Error("RenderNumber(1000)が空文字列を返しました")
	}

	result1500 := renderer.RenderNumber(1500, lipgloss.Color("#FFFFFF"))
	// 1000と1500で同じ結果（999+）になることを確認
	if result != result1500 {
		t.Error("1000以上の数値が統一されていません")
	}
}

// TestASCIINumbersDigitWidthConsistency は全数字の幅が一定であることを確認します。
func TestASCIINumbersDigitWidthConsistency(t *testing.T) {
	renderer := NewASCIINumbers()

	var expectedWidth int
	for digit := 0; digit <= 9; digit++ {
		result := renderer.RenderDigit(digit)
		if result == nil {
			continue
		}

		// 最初の行の幅を基準にする
		currentWidth := len([]rune(result[0]))
		if expectedWidth == 0 {
			expectedWidth = currentWidth
		}

		if currentWidth != expectedWidth {
			t.Errorf("数字%dの幅(%d)が基準幅(%d)と異なります", digit, currentWidth, expectedWidth)
		}
	}
}
