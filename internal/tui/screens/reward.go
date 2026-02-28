// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/ascii"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/rewarding"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== 報酬画面 ====================

// RewardScreen は報酬画面を表します。

type RewardScreen struct {
	result           *rewarding.RewardResult
	styles           *styles.GameStyles
	width            int
	height           int
	pendingTutorials []domain.TutorialDef
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
	case "enter":
		if s.hasPendingTutorials() {
			// PendingTutorial時: 最初のチュートリアルを開く
			tut := s.pendingTutorials[0]
			return s, func() tea.Msg {
				return CompleteTutorialMsg{TutorialID: tut.ID}
			}
		}
		// PendingTutorialなし: ホーム画面へ遷移
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case " ", "esc":
		// Space/Escでホーム画面へ（PendingTutorialの有無に関わらず）
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	}

	return s, nil
}

// sectionBoxWidth はセクションボックスの幅です。
const sectionBoxWidth = 60

// View は画面をレンダリングします。
// 撃破報酬→ランクアップ→タイピング統計の縦並びレイアウトです。
func (s *RewardScreen) View() string {
	var sections []string

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary).
		Align(lipgloss.Center).
		Width(s.width)

	sections = append(sections, titleStyle.Render("🎉 バトル勝利！ 🎉"))

	// 撃破報酬セクション（ドロップアイテムまたはHP増加がある場合のみ）
	if s.hasDefeatReward() {
		sections = append(sections, s.renderDefeatRewardSection())
	}

	// ランクアップセクション（ランク解放時のみ）
	if s.result.RankUnlocked {
		sections = append(sections, s.renderRankUpSection())
	}

	// チュートリアル誘導セクション（PendingTutorialがある場合のみ）
	if s.hasPendingTutorials() {
		sections = append(sections, s.renderTutorialGuidanceSection())
	}

	// タイピング統計セクション（常に表示）
	sections = append(sections, s.renderTypingStatsSection())

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	if s.hasPendingTutorials() {
		sections = append(sections, hintStyle.Render("Enter: チュートリアルを見る  Space/Esc: スキップ"))
	} else {
		sections = append(sections, hintStyle.Render("Enter: 続行"))
	}

	return strings.Join(sections, "\n\n")
}

// hasDefeatReward は撃破報酬があるかを返します。
func (s *RewardScreen) hasDefeatReward() bool {
	return len(s.result.DroppedCores) > 0 ||
		len(s.result.DroppedSkills) > 0 ||
		s.result.HPGain > 0
}

// renderDefeatRewardSection は撃破報酬セクションをレンダリングします。
func (s *RewardScreen) renderDefeatRewardSection() string {
	var items []string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary)

	items = append(items, titleStyle.Render("🎁 撃破報酬"))
	items = append(items, "")

	itemStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	// コアドロップ
	for _, core := range s.result.DroppedCores {
		items = append(items, itemStyle.Render(fmt.Sprintf("  コア: %s", core.Name)))
	}

	// スキルドロップ
	for _, skill := range s.result.DroppedSkills {
		skillInfo := fmt.Sprintf("  スキル: %s %s", skill.Icon(), skill.Name())
		items = append(items, itemStyle.Render(skillInfo))

		// チェイン効果の詳細
		if skill.HasChainEffect() {
			items = append(items, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Render(fmt.Sprintf("    + %s", skill.ChainEffect.Description)))
		}
	}

	// HP増加
	if s.result.HPGain > 0 {
		hpStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(styles.ColorBuff)
		items = append(items, hpStyle.Render(fmt.Sprintf("  💚 最大HP +%d (%d → %d)",
			s.result.HPGain, s.result.PreviousMaxHP, s.result.NewMaxHP)))
	}

	content := strings.Join(items, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(sectionBoxWidth).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// renderRankUpSection はランクアップセクションをレンダリングします。
func (s *RewardScreen) renderRankUpSection() string {
	var items []string

	// RANK UP タイトル
	rankUpTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ascii.GoldColor).
		Align(lipgloss.Center).
		Width(sectionBoxWidth - 4). // パディング分を差し引く
		Render("✨ RANK UP! ✨")

	items = append(items, rankUpTitle)

	// ランク数字を影付きで表示
	shadowRenderer := ascii.NewShadowNumbers()
	rankArt := shadowRenderer.RenderShadowNumber(s.result.NewRank)

	rankContent := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(sectionBoxWidth - 4).
		Render(rankArt)

	items = append(items, rankContent)

	// ランク変化の説明
	descStyle := lipgloss.NewStyle().
		Foreground(ascii.DarkGoldColor).
		Align(lipgloss.Center).
		Width(sectionBoxWidth - 4)

	items = append(items, descStyle.Render(fmt.Sprintf("ランク %d → %d",
		s.result.PreviousRank, s.result.NewRank)))

	// ランクアップ報酬アイテム
	if s.hasRankUpRewards() {
		items = append(items, "")

		rewardTitleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(ascii.GoldColor)
		items = append(items, rewardTitleStyle.Render("ランクアップ報酬:"))

		rewardItemStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

		for _, core := range s.result.RankUpRewardCores {
			items = append(items, rewardItemStyle.Render(fmt.Sprintf("  コア: %s", core.Name)))
		}

		for _, skill := range s.result.RankUpRewardSkills {
			items = append(items, rewardItemStyle.Render(fmt.Sprintf("  スキル: %s %s", skill.Icon(), skill.Name())))
		}

		for _, effect := range s.result.RankUpRewardChainEffects {
			items = append(items, rewardItemStyle.Render(fmt.Sprintf("  チェイン効果: %s", effect.Description)))
		}
	}

	content := strings.Join(items, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ascii.GoldColor).
		Padding(1).
		Width(sectionBoxWidth).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// hasRankUpRewards はランクアップ報酬アイテムがあるかを返します。
func (s *RewardScreen) hasRankUpRewards() bool {
	return len(s.result.RankUpRewardCores) > 0 ||
		len(s.result.RankUpRewardSkills) > 0 ||
		len(s.result.RankUpRewardChainEffects) > 0
}

// renderTutorialGuidanceSection はチュートリアル誘導セクションをレンダリングします。
func (s *RewardScreen) renderTutorialGuidanceSection() string {
	var items []string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(ascii.GoldColor)

	items = append(items, titleStyle.Render("🔓 新機能が解放されました！"))
	items = append(items, "")

	featureStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)
	for _, tut := range s.pendingTutorials {
		items = append(items, featureStyle.Render(fmt.Sprintf("  📖 %s", tut.Title)))
	}

	items = append(items, "")
	guideStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Italic(true)
	items = append(items, guideStyle.Render("  Enterでチュートリアルを見る"))

	content := strings.Join(items, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ascii.GoldColor).
		Padding(1).
		Width(sectionBoxWidth).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// renderTypingStatsSection はタイピング統計セクションをレンダリングします。
func (s *RewardScreen) renderTypingStatsSection() string {
	var items []string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary)

	items = append(items, titleStyle.Render("📊 タイピング統計"))
	items = append(items, "")

	itemStyle := lipgloss.NewStyle().Foreground(styles.ColorSecondary)

	if s.result.Stats != nil {
		avgWPM := s.result.Stats.GetAverageWPM()
		avgScore := s.result.Stats.GetAverageScore()

		items = append(items, itemStyle.Render(fmt.Sprintf("  平均WPM: %.1f", avgWPM)))
		items = append(items, itemStyle.Render(fmt.Sprintf("  平均スコア: %.0f", avgScore)))
		items = append(items, itemStyle.Render(fmt.Sprintf("  総ダメージ: %d", s.result.Stats.TotalDamageDealt)))
		items = append(items, itemStyle.Render(fmt.Sprintf("  被ダメージ: %d", s.result.Stats.TotalDamageTaken)))
		if s.result.Stats.TotalHealAmount > 0 {
			items = append(items, itemStyle.Render(fmt.Sprintf("  回復量: %d", s.result.Stats.TotalHealAmount)))
		}
	} else {
		items = append(items, itemStyle.Render("  統計データなし"))
	}

	content := strings.Join(items, "\n")

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(sectionBoxWidth).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// SetPendingTutorials はPendingTutorial情報を設定します。
// ランクアップで新たにPendingTutorialになったチュートリアルを渡します。
func (s *RewardScreen) SetPendingTutorials(tutorials []domain.TutorialDef) {
	s.pendingTutorials = tutorials
}

// hasPendingTutorials はPendingTutorialがあるかを返します。
func (s *RewardScreen) hasPendingTutorials() bool {
	return len(s.pendingTutorials) > 0
}

// SetSize はウィンドウサイズを設定します。
func (s *RewardScreen) SetSize(width, height int) {
	s.width = width
	s.height = height
}
