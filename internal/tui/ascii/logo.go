// Package ascii はASCIIアート描画機能を提供します。
// ゲームロゴ、数字、WIN/LOSEなどのASCIIアートを担当します。

package ascii

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// BlitzTypingOperatorロゴのASCIIアート（フィグレット風）

var blitzTypingOperatorLogo = []string{
	"╔╗ ╦  ╦╔╦╗╔═╗╔╦╗╦ ╦╔═╗╦╔╗╔╔═╗╔═╗╔═╗╔═╗╦═╗╔═╗╔╦╗╔═╗╦═╗",
	"╠╩╗║  ║ ║ ╔═╝ ║ ╚╦╝╠═╝║║║║║ ╦║ ║╠═╝║╣ ╠╦╝╠═╣ ║ ║ ║╠╦╝",
	"╚═╝╩═╝╩ ╩ ╚═╝ ╩  ╩ ╩  ╩╝╚╝╚═╝╚═╝╩  ╚═╝╩╚═╩ ╩ ╩ ╚═╝╩╚═",
}

// ロゴのカラー定義
var (
	// LogoColorPrimary はロゴのプライマリカラー
	LogoColorPrimary = lipgloss.Color("#7D56F4")
	// LogoColorSecondary はロゴのセカンダリカラー
	LogoColorSecondary = lipgloss.Color("#FAFAFA")
)

// ASCIILogoRenderer はゲームロゴのASCIIアート描画を提供するインターフェースです。
type ASCIILogoRenderer interface {
	// Render はロゴをレンダリングします
	// colorMode: true=カラー、false=モノクロ
	Render(colorMode bool) string

	// GetWidth はロゴの幅（文字数）を返します
	GetWidth() int

	// GetHeight はロゴの高さ（行数）を返します
	GetHeight() int
}

// asciiLogo はASCIILogoRendererの実装です。
type asciiLogo struct {
	lines  []string
	width  int
	height int
}

// NewASCIILogo は新しいASCIILogoRendererを作成します。
func NewASCIILogo() ASCIILogoRenderer {
	// 幅は最長の行の長さを計算
	maxWidth := 0
	for _, line := range blitzTypingOperatorLogo {
		lineWidth := len([]rune(line))
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	return &asciiLogo{
		lines:  blitzTypingOperatorLogo,
		width:  maxWidth,
		height: len(blitzTypingOperatorLogo),
	}
}

// Render はロゴをレンダリングします。
// colorMode: true=カラー、false=モノクロ
func (l *asciiLogo) Render(colorMode bool) string {
	if colorMode {
		// カラーモード：プライマリカラーでスタイリング
		style := lipgloss.NewStyle().
			Foreground(LogoColorPrimary).
			Bold(true)
		return style.Render(strings.Join(l.lines, "\n"))
	}

	// モノクロモード：スタイリングなし
	return strings.Join(l.lines, "\n")
}

// GetWidth はロゴの幅（文字数）を返します。
func (l *asciiLogo) GetWidth() int {
	return l.width
}

// GetHeight はロゴの高さ（行数）を返します。
func (l *asciiLogo) GetHeight() int {
	return l.height
}

// RenderCentered はロゴを指定された幅で中央揃えしてレンダリングします。
func (l *asciiLogo) RenderCentered(width int, colorMode bool) string {
	rendered := l.Render(colorMode)
	lines := strings.Split(rendered, "\n")

	var centered []string
	for _, line := range lines {
		lineWidth := len([]rune(line))
		if lineWidth < width {
			padding := (width - lineWidth) / 2
			line = strings.Repeat(" ", padding) + line
		}
		centered = append(centered, line)
	}

	return strings.Join(centered, "\n")
}
