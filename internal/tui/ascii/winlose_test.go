// Package ascii はASCIIアート描画機能を提供します。

package ascii

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/tui/styles"
)

// TestNewWinLoseRenderer はWinLoseRendererの作成をテストします。
func TestNewWinLoseRenderer(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)
	if renderer == nil {
		t.Error("NewWinLoseRenderer()がnilを返しました")
	}
}

// TestWinLoseRendererRenderWin は勝利時のASCIIアート描画をテストします。

func TestWinLoseRendererRenderWin(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)

	result := renderer.RenderWin()
	if result == "" {
		t.Error("RenderWin()が空文字列を返しました")
	}

	// 出力が複数行であることを確認
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("RenderWin()の行数が少なすぎます: %d行（3行以上必要）", len(lines))
	}
}

// TestWinLoseRendererRenderLose は敗北時のASCIIアート描画をテストします。

func TestWinLoseRendererRenderLose(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)

	result := renderer.RenderLose()
	if result == "" {
		t.Error("RenderLose()が空文字列を返しました")
	}

	// 出力が複数行であることを確認
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Errorf("RenderLose()の行数が少なすぎます: %d行（3行以上必要）", len(lines))
	}
}

// TestWinLoseRendererDifferentResults はWINとLOSEが異なることをテストします。
func TestWinLoseRendererDifferentResults(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)

	winResult := renderer.RenderWin()
	loseResult := renderer.RenderLose()

	// WINとLOSEは異なる結果であることを確認
	if winResult == loseResult {
		t.Error("RenderWin()とRenderLose()が同じ結果を返しました")
	}
}

// TestWinLoseRendererNoColorMode はモノクロモードでの動作をテストします。
func TestWinLoseRendererNoColorMode(t *testing.T) {
	gs := styles.NewGameStylesWithNoColor()
	renderer := NewWinLoseRenderer(gs)

	winResult := renderer.RenderWin()
	if winResult == "" {
		t.Error("モノクロモードでRenderWin()が空文字列を返しました")
	}

	loseResult := renderer.RenderLose()
	if loseResult == "" {
		t.Error("モノクロモードでRenderLose()が空文字列を返しました")
	}
}

// TestWinLoseRendererGetWidth はWIN/LOSEの幅取得をテストします。
func TestWinLoseRendererGetWidth(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)

	width := renderer.GetWidth()
	if width <= 0 {
		t.Errorf("GetWidth()が0以下の値を返しました: %d", width)
	}
}

// TestWinLoseRendererGetHeight はWIN/LOSEの高さ取得をテストします。
func TestWinLoseRendererGetHeight(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewWinLoseRenderer(gs)

	height := renderer.GetHeight()
	if height <= 0 {
		t.Errorf("GetHeight()が0以下の値を返しました: %d", height)
	}

	// エージェントエリア内に収まるサイズ（10行以下）を想定
	if height > 10 {
		t.Errorf("GetHeight()が大きすぎます: %d行（10行以下を想定）", height)
	}
}
