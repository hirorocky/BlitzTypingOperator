// Package components はTUI共通コンポーネントを提供します。
package components

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

// PassiveSkillNotification はパッシブスキルの効果を表示するコンポーネントです。
// エージェントのコア特性から付与されるパッシブスキルの効果を
// コンパクトな形式または詳細な形式で表示します。
type PassiveSkillNotification struct {
	skill      *domain.PassiveSkill
	gameStyles *styles.GameStyles
}

// NewPassiveSkillNotification は新しいPassiveSkillNotificationを作成します。
// skillがnilの場合でも有効な通知を作成します（表示は空になります）。
func NewPassiveSkillNotification(skill *domain.PassiveSkill) *PassiveSkillNotification {
	return &PassiveSkillNotification{
		skill:      skill,
		gameStyles: styles.NewGameStyles(),
	}
}

// GetName はパッシブスキルの名前を返します。
// skillがnilの場合は空文字列を返します。
func (n *PassiveSkillNotification) GetName() string {
	if n.skill == nil {
		return ""
	}
	return n.skill.Name
}

// GetDescription はパッシブスキルの説明を返します。
// skillがnilの場合は空文字列を返します。
func (n *PassiveSkillNotification) GetDescription() string {
	if n.skill == nil {
		return ""
	}
	return n.skill.Description
}

// GetShortDescription はパッシブスキルの短い説明を返します。
// ShortDescriptionが設定されていない場合はDescriptionを返します。
// skillがnilの場合は空文字列を返します。
func (n *PassiveSkillNotification) GetShortDescription() string {
	if n.skill == nil {
		return ""
	}
	if n.skill.ShortDescription != "" {
		return n.skill.ShortDescription
	}
	return n.skill.Description
}

// GetEffects はパッシブスキルの効果マップを返します。
// skillがnilの場合はnilを返します。
func (n *PassiveSkillNotification) GetEffects() map[domain.EffectColumn]float64 {
	if n.skill == nil {
		return nil
	}
	return n.skill.Effects
}

// HasActiveEffects はこの通知が有効なパッシブスキルを持っているかを返します。
func (n *PassiveSkillNotification) HasActiveEffects() bool {
	return n.skill != nil
}

// RenderCompact はコンパクトな形式でパッシブスキルをレンダリングします。
// skillがnilの場合は空文字列を返します。
func (n *PassiveSkillNotification) RenderCompact() string {
	if n.skill == nil {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(styles.ColorBuff).
		Bold(true)

	return style.Render(fmt.Sprintf("★ %s", n.skill.Name))
}

// RenderDetail は詳細な形式でパッシブスキルをレンダリングします。
// 名前、説明、効果リストを含む複数行の出力を返します。
// skillがnilの場合は「パッシブスキルなし」を返します。
func (n *PassiveSkillNotification) RenderDetail(width int) string {
	if n.skill == nil {
		return lipgloss.NewStyle().
			Foreground(styles.ColorSubtle).
			Render("パッシブスキルなし")
	}

	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.ColorBuff).
		Bold(true)
	builder.WriteString(titleStyle.Render(fmt.Sprintf("★ %s", n.skill.Name)))
	builder.WriteString("\n")

	// 説明
	descStyle := lipgloss.NewStyle().Foreground(styles.ColorSubtle)
	builder.WriteString(descStyle.Render(n.skill.Description))
	builder.WriteString("\n")

	// 効果リスト
	effects := n.RenderEffectsList()
	if len(effects) > 0 {
		builder.WriteString("\n")
		labelStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
		builder.WriteString(labelStyle.Render("効果:"))
		builder.WriteString("\n")
		for _, effect := range effects {
			builder.WriteString("  " + effect + "\n")
		}
	}

	return builder.String()
}

// RenderEffectsList は効果のリストを文字列スライスとして返します。
// 各要素は「効果名: 値」形式です。
// skillがnilまたは効果がない場合は空のスライスを返します。
func (n *PassiveSkillNotification) RenderEffectsList() []string {
	if n.skill == nil {
		return []string{}
	}

	effectsMap := n.GetEffects()
	if effectsMap == nil {
		return []string{}
	}

	effects := make([]string, 0)
	valueStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)

	// 効果名のマッピング
	columnNames := map[domain.EffectColumn]string{
		domain.ColDamageBonus:      "ダメージ",
		domain.ColDamageMultiplier: "ダメージ倍率",
		domain.ColDamageCut:        "被ダメ軽減",
		domain.ColEvasion:          "回避率",
		domain.ColTimeExtend:       "入力時間",
		domain.ColCooldownReduce:   "CD短縮",
		domain.ColCritRate:         "クリ率",
		domain.ColSTRBonus:         "STR",
		domain.ColINTBonus:         "INT",
		domain.ColWILBonus:         "WIL",
		domain.ColLUKBonus:         "LUK",
		domain.ColSTRMultiplier:    "STR倍率",
		domain.ColINTMultiplier:    "INT倍率",
		domain.ColWILMultiplier:    "WIL倍率",
		domain.ColLUKMultiplier:    "LUK倍率",
		domain.ColHealBonus:        "回復",
		domain.ColHealMultiplier:   "回復倍率",
		domain.ColLifeSteal:        "HP吸収",
	}

	// 乗算系の列（パーセント表示）
	multColumns := map[domain.EffectColumn]bool{
		domain.ColDamageMultiplier: true,
		domain.ColSTRMultiplier:    true,
		domain.ColINTMultiplier:    true,
		domain.ColWILMultiplier:    true,
		domain.ColLUKMultiplier:    true,
		domain.ColHealMultiplier:   true,
	}

	// パーセント系の列
	percentColumns := map[domain.EffectColumn]bool{
		domain.ColDamageCut:      true,
		domain.ColEvasion:        true,
		domain.ColCooldownReduce: true,
		domain.ColCritRate:       true,
		domain.ColLifeSteal:      true,
	}

	for col, val := range effectsMap {
		if val == 0 {
			continue
		}

		name, ok := columnNames[col]
		if !ok {
			name = string(col)
		}

		var formatted string
		if multColumns[col] {
			// 乗算系は1.0からの差分をパーセント表示
			percent := (val - 1.0) * 100
			formatted = fmt.Sprintf("%+.0f%%", percent)
		} else if percentColumns[col] {
			// パーセント系はそのまま100倍
			formatted = fmt.Sprintf("+%.0f%%", val*100)
		} else if col == domain.ColTimeExtend {
			// 時間延長は秒数
			formatted = fmt.Sprintf("+%.1f秒", val)
		} else {
			// その他は整数または小数
			if val == float64(int(val)) {
				formatted = fmt.Sprintf("%+.0f", val)
			} else {
				formatted = fmt.Sprintf("%+.1f", val)
			}
		}

		effects = append(effects, fmt.Sprintf("%s: %s", name, valueStyle.Render(formatted)))
	}

	return effects
}

// RenderBadge はバッジ形式でパッシブスキルをレンダリングします。
// 小さな領域に表示するための最小形式です。
// skillがnilの場合は空文字列を返します。
func (n *PassiveSkillNotification) RenderBadge() string {
	if n.skill == nil {
		return ""
	}

	style := lipgloss.NewStyle().
		Foreground(styles.ColorBuff).
		Background(lipgloss.Color("236")).
		Padding(0, 1)

	// 短縮名（最初の4文字 + ...）
	name := n.skill.Name
	if len([]rune(name)) > 6 {
		name = string([]rune(name)[:6]) + ".."
	}

	return style.Render("★" + name)
}
