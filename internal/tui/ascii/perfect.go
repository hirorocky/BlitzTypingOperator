package ascii

import (
	"strings"

	"hirorocky/type-battle/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// PERFECT!のASCIIアート
var perfectArt = []string{
	"██████╗ ███████╗██████╗ ███████╗███████╗ ██████╗████████╗██╗",
	"██╔══██╗██╔════╝██╔══██╗██╔════╝██╔════╝██╔════╝╚══██╔══╝██║",
	"██████╔╝█████╗  ██████╔╝█████╗  █████╗  ██║        ██║   ██║",
	"██╔═══╝ ██╔══╝  ██╔══██╗██╔══╝  ██╔══╝  ██║        ██║   ╚═╝",
	"██║     ███████╗██║  ██║██║     ███████╗╚██████╗   ██║   ██╗",
	"╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝     ╚══════╝ ╚═════╝   ╚═╝   ╚═╝",
}

// PerfectRenderer はパーフェクトASCIIアート描画を提供するインターフェースです。
type PerfectRenderer interface {
	// RenderPerfect はPERFECT!のASCIIアートを描画します。
	RenderPerfect() string

	// GetWidth は幅（文字数）を返します。
	GetWidth() int

	// GetHeight は高さ（行数）を返します。
	GetHeight() int
}

// perfectRenderer はPerfectRendererの実装です。
type perfectRenderer struct {
	styles *styles.GameStyles
	width  int
	height int
}

// NewPerfectRenderer は新しいPerfectRendererを作成します。
func NewPerfectRenderer(gs *styles.GameStyles) PerfectRenderer {
	maxWidth := 0
	for _, line := range perfectArt {
		lineWidth := len([]rune(line))
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return &perfectRenderer{
		styles: gs,
		width:  maxWidth,
		height: len(perfectArt),
	}
}

// RenderPerfect はPERFECT!のASCIIアートを描画します。
func (r *perfectRenderer) RenderPerfect() string {
	style := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Bold(true)

	return style.Render(strings.Join(perfectArt, "\n"))
}

// GetWidth は幅（文字数）を返します。
func (r *perfectRenderer) GetWidth() int {
	return r.width
}

// GetHeight は高さ（行数）を返します。
func (r *perfectRenderer) GetHeight() int {
	return r.height
}
