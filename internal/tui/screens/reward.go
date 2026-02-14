// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/tui/ascii"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/rewarding"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== 報酬画面 ====================

// RewardScreen は報酬画面を表します。

type RewardScreen struct {
	result *rewarding.RewardResult
	styles *styles.GameStyles
	width  int
	height int
}

// NewRewardScreen は新しいRewardScreenを作成します。
func NewRewardScreen(result *rewarding.RewardResult) *RewardScreen {
	return &RewardScreen{
		result: result,
		styles: styles.NewGameStyles(),
		width:  140,
		height: 40,
	}
}

// Init は画面の初期化を行います。
func (s *RewardScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *RewardScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}

	return s, nil
}

// handleKeyMsg はキーボード入力を処理します。
func (s *RewardScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		// ホーム画面へ遷移
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "esc":
		// Escでもホーム画面へ
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	}

	return s, nil
}

// View は画面をレンダリングします。
func (s *RewardScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("🎉 バトル勝利！ 🎉"))
	builder.WriteString("\n\n")

	// 敵レベル
	levelStyle := lipgloss.NewStyle().
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(levelStyle.Render(fmt.Sprintf("Lv.%d の敵を撃破！", s.result.EnemyLevel)))
	builder.WriteString("\n\n")

	// メインコンテンツ（統計とドロップ）
	builder.WriteString(s.renderMainContent())
	builder.WriteString("\n\n")

	// HP増加表示
	if s.result.HPGain > 0 {
		builder.WriteString(s.renderHPGain())
		builder.WriteString("\n\n")
	}

	// ランクアップ表示
	if s.result.RankUnlocked {
		builder.WriteString(s.renderRankUp())
		builder.WriteString("\n\n")
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(hintStyle.Render("Enter: 続行"))

	return builder.String()
}

// renderMainContent はメインコンテンツをレンダリングします。
func (s *RewardScreen) renderMainContent() string {
	// 左側：バトル統計、右側：ドロップアイテム
	statsBox := s.renderBattleStats()
	dropsBox := s.renderDrops()

	content := lipgloss.JoinHorizontal(lipgloss.Top, statsBox, "  ", dropsBox)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(content)
}

// renderBattleStats はバトル統計をレンダリングします。

func (s *RewardScreen) renderBattleStats() string {
	var items []string

	itemStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	// WPM
	if s.result.Stats != nil {
		avgWPM := s.result.Stats.GetAverageWPM()
		avgAccuracy := s.result.Stats.GetAverageAccuracy()

		items = append(items, itemStyle.Render(fmt.Sprintf("平均WPM: %.1f", avgWPM)))
		items = append(items, itemStyle.Render(fmt.Sprintf("平均正確性: %.1f%%", avgAccuracy)))
		items = append(items, itemStyle.Render(fmt.Sprintf("総ダメージ: %d", s.result.Stats.TotalDamageDealt)))
		items = append(items, itemStyle.Render(fmt.Sprintf("被ダメージ: %d", s.result.Stats.TotalDamageTaken)))
		if s.result.Stats.TotalHealAmount > 0 {
			items = append(items, itemStyle.Render(fmt.Sprintf("回復量: %d", s.result.Stats.TotalHealAmount)))
		}
	} else {
		items = append(items, itemStyle.Render("統計データなし"))
	}

	content := strings.Join(items, "\n")

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(35).
		Render(titleStyle.Render("📊 バトル統計") + "\n\n" + content)
}

// renderDrops はドロップアイテムをレンダリングします。

func (s *RewardScreen) renderDrops() string {
	var items []string

	// コアドロップ
	if len(s.result.DroppedCores) > 0 {
		coreStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorSecondary)
		items = append(items, coreStyle.Render("【コア】"))

		for _, core := range s.result.DroppedCores {
			coreInfo := fmt.Sprintf("  %s", core.Name)
			items = append(items, lipgloss.NewStyle().
				Foreground(styles.ColorSecondary).
				Render(coreInfo))
		}
		items = append(items, "")
	}

	// スキルドロップ
	if len(s.result.DroppedSkills) > 0 {
		skillStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorInfo)
		items = append(items, skillStyle.Render("【スキル】"))

		for _, skill := range s.result.DroppedSkills {
			// スキル基本情報
			skillInfo := fmt.Sprintf("  %s %s", skill.Icon(), skill.Name())

			items = append(items, lipgloss.NewStyle().
				Foreground(styles.ColorSecondary).
				Render(skillInfo))

			// チェイン効果の詳細説明を追加
			if skill.HasChainEffect() {
				items = append(items, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#FFFFFF")).
					Render(fmt.Sprintf("    + %s", skill.ChainEffect.Description)))
			}
		}
	}

	// ドロップなしの場合
	if len(s.result.DroppedCores) == 0 && len(s.result.DroppedSkills) == 0 {
		items = append(items, lipgloss.NewStyle().
			Foreground(styles.ColorSubtle).
			Render("ドロップアイテムなし"))
	}

	content := strings.Join(items, "\n")

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(35).
		Render(titleStyle.Render("🎁 ドロップ") + "\n\n" + content)
}

// renderHPGain はHP増加を表示します。
func (s *RewardScreen) renderHPGain() string {
	hpStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorBuff).
		Align(lipgloss.Center).
		Width(s.width)

	return hpStyle.Render(fmt.Sprintf("💚 最大HP +%d！ (%d → %d)",
		s.result.HPGain, s.result.PreviousMaxHP, s.result.NewMaxHP))
}

// renderRankUp はランクアップを大きく表示します。
func (s *RewardScreen) renderRankUp() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ascii.GoldColor).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("✨ RANK UP! ✨"))
	builder.WriteString("\n\n")

	// ランク数字を影付きで表示
	shadowRenderer := ascii.NewShadowNumbers()
	rankArt := shadowRenderer.RenderShadowNumber(s.result.NewRank)

	rankContent := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(s.width).
		Render(rankArt)

	builder.WriteString(rankContent)
	builder.WriteString("\n")

	// ランク変化の説明
	descStyle := lipgloss.NewStyle().
		Foreground(ascii.DarkGoldColor).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(descStyle.Render(fmt.Sprintf("ランク %d → %d",
		s.result.PreviousRank, s.result.NewRank)))

	return builder.String()
}

// SetSize はウィンドウサイズを設定します。
func (s *RewardScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}
