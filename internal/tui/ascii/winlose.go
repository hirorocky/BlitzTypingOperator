// Package ascii はASCIIアート描画機能を提供します。
// 戦闘結果（WIN/LOSE）のASCIIアート描画を担当します。

package ascii

import (
	"strings"

	"hirorocky/type-battle/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// WINのASCIIアート

var winArt = []string{
	"██╗    ██╗██╗███╗   ██╗██╗",
	"██║    ██║██║████╗  ██║██║",
	"██║ █╗ ██║██║██╔██╗ ██║██║",
	"██║███╗██║██║██║╚██╗██║╚═╝",
	"╚███╔███╔╝██║██║ ╚████║██╗",
	" ╚══╝╚══╝ ╚═╝╚═╝  ╚═══╝╚═╝",
}

// LOSEのASCIIアート

var loseArt = []string{
	"██╗      ██████╗ ███████╗███████╗",
	"██║     ██╔═══██╗██╔════╝██╔════╝",
	"██║     ██║   ██║███████╗█████╗  ",
	"██║     ██║   ██║╚════██║██╔══╝  ",
	"███████╗╚██████╔╝███████║███████╗",
	"╚══════╝ ╚═════╝ ╚══════╝╚══════╝",
}

// WinLoseRenderer は戦闘結果のASCIIアート描画を提供するインターフェースです。
type WinLoseRenderer interface {
	// RenderWin は勝利時のASCIIアートを描画します。
	RenderWin() string

	// RenderLose は敗北時のASCIIアートを描画します。
	RenderLose() string

	// GetWidth は最大幅（文字数）を返します。
	GetWidth() int

	// GetHeight は高さ（行数）を返します。
	GetHeight() int
}

// winLoseRenderer はWinLoseRendererの実装です。
type winLoseRenderer struct {
	styles    *styles.GameStyles
	winWidth  int
	loseWidth int
	height    int
}

// NewWinLoseRenderer は新しいWinLoseRendererを作成します。
func NewWinLoseRenderer(gs *styles.GameStyles) WinLoseRenderer {
	// WINの幅を計算
	winMaxWidth := 0
	for _, line := range winArt {
		lineWidth := len([]rune(line))
		if lineWidth > winMaxWidth {
			winMaxWidth = lineWidth
		}
	}

	// LOSEの幅を計算
	loseMaxWidth := 0
	for _, line := range loseArt {
		lineWidth := len([]rune(line))
		if lineWidth > loseMaxWidth {
			loseMaxWidth = lineWidth
		}
	}

	return &winLoseRenderer{
		styles:    gs,
		winWidth:  winMaxWidth,
		loseWidth: loseMaxWidth,
		height:    len(winArt), // WINとLOSEは同じ高さ
	}
}

// RenderWin は勝利時のASCIIアートを描画します。

func (r *winLoseRenderer) RenderWin() string {
	// 緑色（ColorHPHigh = 成功色）でスタイリング
	style := lipgloss.NewStyle().
		Foreground(styles.ColorHPHigh).
		Bold(true)

	return style.Render(strings.Join(winArt, "\n"))
}

// RenderLose は敗北時のASCIIアートを描画します。

func (r *winLoseRenderer) RenderLose() string {
	// 赤色（ColorDamage = 失敗色）でスタイリング
	style := lipgloss.NewStyle().
		Foreground(styles.ColorDamage).
		Bold(true)

	return style.Render(strings.Join(loseArt, "\n"))
}

// GetWidth は最大幅（文字数）を返します。
// LOSEの方が長いのでLOSEの幅を返す
func (r *winLoseRenderer) GetWidth() int {
	if r.winWidth > r.loseWidth {
		return r.winWidth
	}
	return r.loseWidth
}

// GetHeight は高さ（行数）を返します。
func (r *winLoseRenderer) GetHeight() int {
	return r.height
}

// RenderWinCentered は勝利時のASCIIアートを中央揃えで描画します。
func (r *winLoseRenderer) RenderWinCentered(width int) string {
	rendered := r.RenderWin()
	return r.centerText(rendered, width, r.winWidth)
}

// RenderLoseCentered は敗北時のASCIIアートを中央揃えで描画します。
func (r *winLoseRenderer) RenderLoseCentered(width int) string {
	rendered := r.RenderLose()
	return r.centerText(rendered, width, r.loseWidth)
}

// centerText はテキストを中央揃えにします。
func (r *winLoseRenderer) centerText(text string, targetWidth, textWidth int) string {
	if textWidth >= targetWidth {
		return text
	}

	lines := strings.Split(text, "\n")
	padding := (targetWidth - textWidth) / 2
	paddingStr := strings.Repeat(" ", padding)

	var centered []string
	for _, line := range lines {
		centered = append(centered, paddingStr+line)
	}

	return strings.Join(centered, "\n")
}
