// Package styles はTUIスタイリングのスキルアイコン機能を提供します。

package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// 効果タイプ別のデフォルトカラー
var effectColors = map[string]lipgloss.Color{
	"damage":  ColorDamage, // ダメージ系は赤
	"heal":    ColorHPHigh, // 回復系は緑
	"buff":    ColorBuff,   // バフ系は青
	"debuff":  ColorDebuff, // デバフ系はピンク
	"default": ColorSubtle, // その他
}

// GetEffectColor は効果タイプに対応するカラーを返します。
func GetEffectColor(effectType string) lipgloss.Color {
	if color, ok := effectColors[effectType]; ok {
		return color
	}
	return effectColors["default"]
}

// RenderColoredIcon はカラー付きアイコンを返します。
func RenderColoredIcon(icon string, color lipgloss.Color) string {
	style := lipgloss.NewStyle().Foreground(color)
	return style.Render(icon)
}

// RenderSkillIcon はスキルのアイコンをカラー付きで描画します。
// effectType: "damage", "heal", "buff", "debuff" のいずれか
func RenderSkillIcon(icon string, effectType string) string {
	color := GetEffectColor(effectType)
	return RenderColoredIcon(icon, color)
}

// RenderIcons は複数のアイコンを描画します。
func RenderIcons(icons []string, color lipgloss.Color) []string {
	result := make([]string, len(icons))
	for i, icon := range icons {
		result[i] = RenderColoredIcon(icon, color)
	}
	return result
}
