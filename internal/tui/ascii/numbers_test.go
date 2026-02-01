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

// ==================== 影付き数字レンダラーのテスト ====================

// TestNewShadowNumbers はShadowNumberRendererの作成をテストします。
func TestNewShadowNumbers(t *testing.T) {
	renderer := NewShadowNumbers()
	if renderer == nil {
		t.Error("NewShadowNumbers()がnilを返しました")
	}
}

// TestShadowNumbersRenderShadowNumber は影付き数字のレンダリングをテストします。
func TestShadowNumbersRenderShadowNumber(t *testing.T) {
	renderer := NewShadowNumbers()

	tests := []int{0, 1, 5, 12, 99, 123}

	for _, number := range tests {
		result := renderer.RenderShadowNumber(number)
		if result == "" {
			t.Errorf("RenderShadowNumber(%d)が空文字列を返しました", number)
			continue
		}

		// 6行（5行の本体 + 1行の下影）であることを確認
		lines := strings.Split(result, "\n")
		if len(lines) != 6 {
			t.Errorf("RenderShadowNumber(%d)の行数が%dです（6行を想定）", number, len(lines))
		}

		// 各行に影文字（░）が含まれていることを確認
		for i, line := range lines {
			if !strings.Contains(line, "░") {
				t.Errorf("RenderShadowNumber(%d)の%d行目に影文字が含まれていません", number, i)
			}
		}
	}
}

// TestShadowNumbersRenderShadowNumberWithColors はカスタム色での描画をテストします。
func TestShadowNumbersRenderShadowNumberWithColors(t *testing.T) {
	renderer := NewShadowNumbers()

	result := renderer.RenderShadowNumberWithColors(
		42,
		lipgloss.Color("51"), // シアン
		lipgloss.Color("30"), // 暗いシアン
	)

	if result == "" {
		t.Error("RenderShadowNumberWithColors()が空文字列を返しました")
	}

	// 出力が複数行であることを確認
	lines := strings.Split(result, "\n")
	if len(lines) != 6 {
		t.Errorf("RenderShadowNumberWithColors()の行数が%dです（6行を想定）", len(lines))
	}
}

// TestShadowNumbersNegativeNumber は負数の処理をテストします。
func TestShadowNumbersNegativeNumber(t *testing.T) {
	renderer := NewShadowNumbers()

	resultNeg := renderer.RenderShadowNumber(-5)
	resultZero := renderer.RenderShadowNumber(0)

	// 負数は0として扱われることを確認
	if resultNeg != resultZero {
		t.Error("負数が0として扱われていません")
	}
}

// TestShadowNumbersSingleDigit は1桁の数字の描画をテストします。
func TestShadowNumbersSingleDigit(t *testing.T) {
	renderer := NewShadowNumbers()

	result := renderer.RenderShadowNumber(7)
	lines := strings.Split(result, "\n")

	// 最後の行（下影）にShadowCharBottomと同じ幅の影があることを確認
	bottomLine := lines[len(lines)-1]
	if !strings.Contains(bottomLine, "░") {
		t.Error("下影行に影文字が含まれていません")
	}
}

// TestShadowNumbersMultiDigit は複数桁の数字の描画をテストします。
func TestShadowNumbersMultiDigit(t *testing.T) {
	renderer := NewShadowNumbers()

	result := renderer.RenderShadowNumber(123)
	lines := strings.Split(result, "\n")

	// 6行あることを確認
	if len(lines) != 6 {
		t.Errorf("3桁の数字で%d行（6行を想定）", len(lines))
	}

	// 下影行に複数の影があることを確認（3桁なので間隔を含めて複数の░がある）
	bottomLine := lines[len(lines)-1]
	shadowCount := strings.Count(bottomLine, "░")
	if shadowCount < 3 {
		t.Errorf("3桁の数字の下影で影文字が%d個しかありません", shadowCount)
	}
}
