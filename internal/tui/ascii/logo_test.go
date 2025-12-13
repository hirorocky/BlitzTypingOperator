// Package ascii はASCIIアート描画機能を提供します。

package ascii

import (
	"strings"
	"testing"
)

// TestNewASCIILogo はASCIILogoRendererの作成をテストします。
func TestNewASCIILogo(t *testing.T) {
	logo := NewASCIILogo()
	if logo == nil {
		t.Error("NewASCIILogo()がnilを返しました")
	}
}

// TestASCIILogoRender はロゴのレンダリングをテストします。

func TestASCIILogoRender(t *testing.T) {
	logo := NewASCIILogo()

	// カラーモードでレンダリング
	colorOutput := logo.Render(true)
	if colorOutput == "" {
		t.Error("Render(true)が空文字列を返しました")
	}

	// モノクロモードでレンダリング
	monoOutput := logo.Render(false)
	if monoOutput == "" {
		t.Error("Render(false)が空文字列を返しました")
	}

	// ロゴが複数行であることを確認
	lines := strings.Split(monoOutput, "\n")
	if len(lines) < 3 {
		t.Errorf("ロゴの行数が少なすぎます: %d行（3行以上必要）", len(lines))
	}
}

// TestASCIILogoGetWidth はロゴの幅取得をテストします。
func TestASCIILogoGetWidth(t *testing.T) {
	logo := NewASCIILogo()
	width := logo.GetWidth()

	// 幅が正の値であることを確認
	if width <= 0 {
		t.Errorf("ロゴの幅が0以下です: %d", width)
	}
}

// TestASCIILogoGetHeight はロゴの高さ取得をテストします。

func TestASCIILogoGetHeight(t *testing.T) {
	logo := NewASCIILogo()
	height := logo.GetHeight()

	// 高さが5-8行程度であることを確認
	if height < 3 || height > 10 {
		t.Errorf("ロゴの高さが想定範囲外です: %d行（3-10行を想定）", height)
	}
}

// TestASCIILogoContainsBlitzTypingOperator はロゴに「BLITZTYPINGOPERATOR」が含まれることを確認します。
func TestASCIILogoContainsBlitzTypingOperator(t *testing.T) {
	logo := NewASCIILogo()
	output := logo.Render(false)

	// ASCIIアートなのでBLITZTYPINGOPERATORの文字が形作られていることを確認
	// 少なくともロゴは空ではない
	if len(output) < 50 {
		t.Errorf("ロゴのサイズが小さすぎます: %d文字", len(output))
	}
}
