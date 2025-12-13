// Package styles はTUIゲームのスタイリングを提供します。
// ボックス描画、カラー表示、HPバーなどの視覚的表現を担当します。

package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HP色分けの閾値
const (
	// HPHighThreshold はHP高（緑）の閾値（50%以上）
	HPHighThreshold = 0.50
	// HPMediumThreshold はHP中（黄）の閾値（25%以上）
	HPMediumThreshold = 0.25
)

// カラーパレット
var (
	// 基本色
	ColorPrimary    = lipgloss.Color("#9B7BF7")
	ColorSecondary  = lipgloss.Color("#FAFAFA")
	ColorSelectedBg = lipgloss.Color("#5E43A4") // 選択中アイテムの背景（暗い紫）
	ColorSelectedFg = lipgloss.Color("#FFFFFF") // 選択中アイテムの前景（白）

	// HP色分け（Requirement 18.2）
	ColorHPHigh   = lipgloss.Color("#04B575") // 緑（50%以上）
	ColorHPMedium = lipgloss.Color("#FFB454") // 黄（25%以上50%未満）
	ColorHPLow    = lipgloss.Color("#FF4672") // 赤（25%未満）

	// ダメージ・回復色（Requirement 18.2）
	ColorDamage = lipgloss.Color("#FF4672") // 赤
	ColorHeal   = lipgloss.Color("#04B575") // 緑

	// その他
	ColorSubtle  = lipgloss.Color("#6C6C6C")
	ColorWarning = lipgloss.Color("#FFB454")
	ColorInfo    = lipgloss.Color("#00BFFF")

	// バフ・デバフ
	ColorBuff   = lipgloss.Color("#00BFFF") // 青
	ColorDebuff = lipgloss.Color("#FF69B4") // ピンク
)

// GameStyles はゲーム全体で使用するスタイルを保持します。

type GameStyles struct {
	// Box はボックス（枠線）スタイル
	Box BoxStyle

	// Text はテキストスタイル
	Text TextStyles

	// HP はHP関連スタイル
	HP HPStyles

	// Battle はバトル関連スタイル
	Battle BattleStyles

	// noColor はカラー無効モードかどうか
	noColor bool
}

// BoxStyle はボックス描画のスタイルを保持します。
type BoxStyle struct {
	Border           lipgloss.Border
	BorderStyle      lipgloss.Style
	TitleStyle       lipgloss.Style
	ContentPadding   int
	ContentMargin    int
	BorderForeground lipgloss.Color
}

// TextStyles はテキスト用スタイルを保持します。
type TextStyles struct {
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Normal   lipgloss.Style
	Subtle   lipgloss.Style
	Bold     lipgloss.Style
	Error    lipgloss.Style
	Success  lipgloss.Style
	Warning  lipgloss.Style
	Info     lipgloss.Style
}

// HPStyles はHP表示用スタイルを保持します。
type HPStyles struct {
	BarFilled   lipgloss.Style
	BarEmpty    lipgloss.Style
	BarBorder   lipgloss.Style
	ValueHigh   lipgloss.Style
	ValueMedium lipgloss.Style
	ValueLow    lipgloss.Style
}

// BattleStyles はバトル画面用スタイルを保持します。
type BattleStyles struct {
	Damage    lipgloss.Style
	Heal      lipgloss.Style
	Buff      lipgloss.Style
	Debuff    lipgloss.Style
	Cooldown  lipgloss.Style
	Available lipgloss.Style
}

// NewGameStyles は新しいGameStylesを作成します。
func NewGameStyles() *GameStyles {
	return createStyles(false)
}

// NewGameStylesWithNoColor はカラー無効モードでGameStylesを作成します。

func NewGameStylesWithNoColor() *GameStyles {
	return createStyles(true)
}

// createStyles はスタイルを作成するヘルパー関数です。
func createStyles(noColor bool) *GameStyles {
	gs := &GameStyles{
		noColor: noColor,
	}

	// ボックススタイル（Requirement 18.1）
	gs.Box = BoxStyle{
		Border: lipgloss.RoundedBorder(),
		BorderStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 1),
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary),
		ContentPadding:   1,
		ContentMargin:    0,
		BorderForeground: ColorPrimary,
	}

	// テキストスタイル
	gs.Text = TextStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1),
		Subtitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary),
		Normal: lipgloss.NewStyle().
			Foreground(ColorSecondary),
		Subtle: lipgloss.NewStyle().
			Foreground(ColorSubtle),
		Bold: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorSecondary),
		Error: lipgloss.NewStyle().
			Foreground(ColorDamage).
			Bold(true),
		Success: lipgloss.NewStyle().
			Foreground(ColorHeal).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(ColorWarning),
		Info: lipgloss.NewStyle().
			Foreground(ColorInfo),
	}

	// HPスタイル（Requirement 18.2, 18.4）
	gs.HP = HPStyles{
		BarFilled: lipgloss.NewStyle().
			Background(ColorHPHigh),
		BarEmpty: lipgloss.NewStyle().
			Background(ColorSubtle),
		BarBorder: lipgloss.NewStyle().
			Foreground(ColorSecondary),
		ValueHigh: lipgloss.NewStyle().
			Foreground(ColorHPHigh).
			Bold(true),
		ValueMedium: lipgloss.NewStyle().
			Foreground(ColorHPMedium).
			Bold(true),
		ValueLow: lipgloss.NewStyle().
			Foreground(ColorHPLow).
			Bold(true),
	}

	// バトルスタイル
	gs.Battle = BattleStyles{
		Damage: lipgloss.NewStyle().
			Foreground(ColorDamage).
			Bold(true),
		Heal: lipgloss.NewStyle().
			Foreground(ColorHeal).
			Bold(true),
		Buff: lipgloss.NewStyle().
			Foreground(ColorBuff),
		Debuff: lipgloss.NewStyle().
			Foreground(ColorDebuff),
		Cooldown: lipgloss.NewStyle().
			Foreground(ColorSubtle),
		Available: lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true),
	}

	return gs
}

// GetHPColorType はHP割合に応じた色タイプを返します。

func (gs *GameStyles) GetHPColorType(percentage float64) string {
	if percentage > HPHighThreshold {
		return "green"
	} else if percentage > HPMediumThreshold {
		return "yellow"
	}
	return "red"
}

// GetHPColor はHP割合に応じた色を返します。
func (gs *GameStyles) GetHPColor(percentage float64) lipgloss.Color {
	if percentage > HPHighThreshold {
		return ColorHPHigh
	} else if percentage > HPMediumThreshold {
		return ColorHPMedium
	}
	return ColorHPLow
}

// GetHPStyle はHP割合に応じたスタイルを返します。
func (gs *GameStyles) GetHPStyle(percentage float64) lipgloss.Style {
	if percentage > HPHighThreshold {
		return gs.HP.ValueHigh
	} else if percentage > HPMediumThreshold {
		return gs.HP.ValueMedium
	}
	return gs.HP.ValueLow
}

// RenderHPBar はHPバーを描画します。

func (gs *GameStyles) RenderHPBar(current, max, width int) string {
	if max <= 0 {
		max = 1
	}
	if current < 0 {
		current = 0
	}
	if current > max {
		current = max
	}

	percentage := float64(current) / float64(max)

	// バー内部の幅（ボーダー分を除く）
	innerWidth := width - 2
	if innerWidth < 1 {
		innerWidth = 1
	}

	// 塗りつぶし部分の幅
	filledWidth := int(float64(innerWidth) * percentage)
	emptyWidth := innerWidth - filledWidth

	// 色を選択
	fillColor := gs.GetHPColor(percentage)

	// バーを構築
	var bar strings.Builder
	bar.WriteString("[")

	// カラーモードまたは非カラーモードで描画
	if gs.noColor {

		bar.WriteString(strings.Repeat("#", filledWidth))
		bar.WriteString(strings.Repeat("-", emptyWidth))
	} else {
		filledStyle := lipgloss.NewStyle().Background(fillColor)
		emptyStyle := lipgloss.NewStyle().Background(ColorSubtle)
		bar.WriteString(filledStyle.Render(strings.Repeat(" ", filledWidth)))
		bar.WriteString(emptyStyle.Render(strings.Repeat(" ", emptyWidth)))
	}

	bar.WriteString("]")

	return bar.String()
}

// RenderHPBarWithValue はHPバーと数値を一緒に描画します。
func (gs *GameStyles) RenderHPBarWithValue(current, max, width int) string {
	percentage := float64(current) / float64(max)
	bar := gs.RenderHPBar(current, max, width)
	valueStyle := gs.GetHPStyle(percentage)
	value := valueStyle.Render(fmt.Sprintf("%d/%d", current, max))
	return fmt.Sprintf("%s %s", bar, value)
}

// RenderDamage はダメージ値を描画します。

func (gs *GameStyles) RenderDamage(amount int) string {
	if gs.noColor {
		return fmt.Sprintf("-%d", amount)
	}
	return gs.Battle.Damage.Render(fmt.Sprintf("-%d", amount))
}

// RenderHeal は回復値を描画します。

func (gs *GameStyles) RenderHeal(amount int) string {
	if gs.noColor {
		return fmt.Sprintf("+%d", amount)
	}
	return gs.Battle.Heal.Render(fmt.Sprintf("+%d", amount))
}

// RenderBox はボックスで囲んだコンテンツを描画します。

func (gs *GameStyles) RenderBox(content string, width int) string {
	style := lipgloss.NewStyle().
		Border(gs.Box.Border).
		BorderForeground(gs.Box.BorderForeground).
		Padding(0, 1).
		Width(width)

	return style.Render(content)
}

// RenderBoxWithTitle はタイトル付きボックスを描画します。
func (gs *GameStyles) RenderBoxWithTitle(title, content string, width int) string {
	titleRendered := gs.Box.TitleStyle.Render(title)
	boxContent := fmt.Sprintf("%s\n%s", titleRendered, content)
	return gs.RenderBox(boxContent, width)
}

// RenderBuff はバフ表示を描画します。
func (gs *GameStyles) RenderBuff(name string, remainingSeconds float64) string {
	if remainingSeconds > 0 {
		return gs.Battle.Buff.Render(fmt.Sprintf("[%s %.1fs]", name, remainingSeconds))
	}
	return gs.Battle.Buff.Render(fmt.Sprintf("[%s]", name))
}

// RenderDebuff はデバフ表示を描画します。
func (gs *GameStyles) RenderDebuff(name string, remainingSeconds float64) string {
	if remainingSeconds > 0 {
		return gs.Battle.Debuff.Render(fmt.Sprintf("[%s %.1fs]", name, remainingSeconds))
	}
	return gs.Battle.Debuff.Render(fmt.Sprintf("[%s]", name))
}

// RenderCooldown はクールダウン表示を描画します。
func (gs *GameStyles) RenderCooldown(seconds float64) string {
	return gs.Battle.Cooldown.Render(fmt.Sprintf("CD: %.1fs", seconds))
}

// RenderProgressBar は汎用プログレスバーを描画します。
func (gs *GameStyles) RenderProgressBar(progress float64, width int, filledColor, emptyColor lipgloss.Color) string {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	filledWidth := int(float64(width) * progress)
	emptyWidth := width - filledWidth

	filledStyle := lipgloss.NewStyle().Background(filledColor)
	emptyStyle := lipgloss.NewStyle().Background(emptyColor)

	return filledStyle.Render(strings.Repeat(" ", filledWidth)) +
		emptyStyle.Render(strings.Repeat(" ", emptyWidth))
}
