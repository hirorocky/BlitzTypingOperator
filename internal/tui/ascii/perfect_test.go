package ascii

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/tui/styles"
)

// TestPerfectRenderer_RenderPerfect はPERFECT!アートがレンダリングされることを確認します。
func TestPerfectRenderer_RenderPerfect(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewPerfectRenderer(gs)

	result := renderer.RenderPerfect()

	if result == "" {
		t.Error("RenderPerfect()が空文字列を返した")
	}

	// ASCIIアートはブロック文字で構成されるため、行数で確認
	lines := strings.Split(result, "\n")
	if len(lines) < renderer.GetHeight() {
		t.Errorf("RenderPerfect()の行数が不足: got %d, want >= %d", len(lines), renderer.GetHeight())
	}
}

// TestPerfectRenderer_GetWidth は幅が正の値であることを確認します。
func TestPerfectRenderer_GetWidth(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewPerfectRenderer(gs)

	width := renderer.GetWidth()
	if width <= 0 {
		t.Errorf("GetWidth() = %d, want > 0", width)
	}
}

// TestPerfectRenderer_GetHeight は高さが正の値であることを確認します。
func TestPerfectRenderer_GetHeight(t *testing.T) {
	gs := styles.NewGameStyles()
	renderer := NewPerfectRenderer(gs)

	height := renderer.GetHeight()
	if height <= 0 {
		t.Errorf("GetHeight() = %d, want > 0", height)
	}
}
