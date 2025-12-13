// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ==================== Task 10.6: 統計・実績画面 ====================

// StatsTab は統計・実績画面のタブを表します。
type StatsTab int

const (
	// TabTypingStats はタイピング統計タブです。
	TabTypingStats StatsTab = iota
	// TabBattleStats はバトル統計タブです。
	TabBattleStats
	// TabAchievements は実績タブです。
	TabAchievements
)

// StatsAchievementsScreen は統計・実績画面を表します。

type StatsAchievementsScreen struct {
	data          *StatsData
	currentTab    StatsTab
	selectedIndex int
	styles        *styles.GameStyles
	width         int
	height        int
}

// NewStatsAchievementsScreen は新しいStatsAchievementsScreenを作成します。
func NewStatsAchievementsScreen(data *StatsData) *StatsAchievementsScreen {
	return &StatsAchievementsScreen{
		data:          data,
		currentTab:    TabTypingStats,
		selectedIndex: 0,
		styles:        styles.NewGameStyles(),
		width:         140,
		height:        40,
	}
}

// Init は画面の初期化を行います。
func (s *StatsAchievementsScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理します。
func (s *StatsAchievementsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
func (s *StatsAchievementsScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "left", "h":
		s.prevTab()
	case "right", "l":
		s.nextTab()
	case "up", "k":
		s.moveUp()
	case "down", "j":
		s.moveDown()
	}

	return s, nil
}

// prevTab は前のタブに移動します。
func (s *StatsAchievementsScreen) prevTab() {
	if s.currentTab > TabTypingStats {
		s.currentTab--
		s.selectedIndex = 0
	}
}

// nextTab は次のタブに移動します。
func (s *StatsAchievementsScreen) nextTab() {
	if s.currentTab < TabAchievements {
		s.currentTab++
		s.selectedIndex = 0
	}
}

// moveUp は選択を上に移動します。
func (s *StatsAchievementsScreen) moveUp() {
	if s.selectedIndex > 0 {
		s.selectedIndex--
	}
}

// moveDown は選択を下に移動します。
func (s *StatsAchievementsScreen) moveDown() {
	maxIndex := len(s.data.Achievements)
	if s.currentTab == TabAchievements && s.selectedIndex < maxIndex-1 {
		s.selectedIndex++
	}
}

// getAchievementCompletionRate は実績のコンプリート率を返します。
func (s *StatsAchievementsScreen) getAchievementCompletionRate() int {
	if len(s.data.Achievements) == 0 {
		return 0
	}
	achieved := 0
	for _, ach := range s.data.Achievements {
		if ach.Achieved {
			achieved++
		}
	}
	return achieved * 100 / len(s.data.Achievements)
}

// View は画面をレンダリングします。
func (s *StatsAchievementsScreen) View() string {
	var builder strings.Builder

	// タイトル
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorPrimary).
		Align(lipgloss.Center).
		Width(s.width)

	builder.WriteString(titleStyle.Render("統計/実績"))
	builder.WriteString("\n\n")

	// タブバー
	builder.WriteString(s.renderTabBar())
	builder.WriteString("\n\n")

	// メインコンテンツ
	builder.WriteString(s.renderMainContent())
	builder.WriteString("\n\n")

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	hints := "←/→: タブ切替  ↑/↓: 選択  Esc: 戻る"
	builder.WriteString(hintStyle.Render(hints))

	return builder.String()
}

// renderTabBar はタブバーをレンダリングします。
func (s *StatsAchievementsScreen) renderTabBar() string {
	tabs := []string{"タイピング統計", "バトル統計", "実績"}

	var tabItems []string
	for i, tab := range tabs {
		style := lipgloss.NewStyle().Padding(0, 2)
		prefix := "  "
		if StatsTab(i) == s.currentTab {
			prefix = "> "
			style = style.
				Bold(true).
				Foreground(styles.ColorSelectedFg).
				Background(styles.ColorSelectedBg)
		} else {
			style = style.Foreground(styles.ColorSubtle)
		}
		tabItems = append(tabItems, style.Render(prefix+tab))
	}

	tabBar := lipgloss.JoinHorizontal(lipgloss.Center, tabItems...)
	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(tabBar)
}

// renderMainContent はメインコンテンツをレンダリングします。
func (s *StatsAchievementsScreen) renderMainContent() string {
	switch s.currentTab {
	case TabTypingStats:
		return s.renderTypingStats()
	case TabBattleStats:
		return s.renderBattleStats()
	case TabAchievements:
		return s.renderAchievements()
	}
	return ""
}

// renderTypingStats はタイピング統計をレンダリングします。

func (s *StatsAchievementsScreen) renderTypingStats() string {
	panel := components.NewInfoPanel("タイピング統計")
	panel.AddItem("最高WPM", fmt.Sprintf("%d WPM", s.data.TypingStats.MaxWPM))
	panel.AddItem("平均WPM", fmt.Sprintf("%.1f WPM", s.data.TypingStats.AverageWPM))
	panel.AddItem("100%正確性達成", fmt.Sprintf("%d回", s.data.TypingStats.PerfectAccuracyCount))
	panel.AddItem("総タイプ文字数", fmt.Sprintf("%d文字", s.data.TypingStats.TotalCharacters))

	content := panel.Render(50)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(60).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// renderBattleStats はバトル統計をレンダリングします。

func (s *StatsAchievementsScreen) renderBattleStats() string {
	winRate := float64(0)
	if s.data.BattleStats.TotalBattles > 0 {
		winRate = float64(s.data.BattleStats.Wins) * 100 / float64(s.data.BattleStats.TotalBattles)
	}

	panel := components.NewInfoPanel("バトル統計")
	panel.AddItem("総バトル数", fmt.Sprintf("%d戦", s.data.BattleStats.TotalBattles))
	panel.AddItem("勝利数", fmt.Sprintf("%d勝", s.data.BattleStats.Wins))
	panel.AddItem("敗北数", fmt.Sprintf("%d敗", s.data.BattleStats.Losses))
	panel.AddItem("勝率", fmt.Sprintf("%.1f%%", winRate))
	panel.AddItem("到達最高レベル", fmt.Sprintf("Lv.%d", s.data.BattleStats.MaxLevelReached))

	content := panel.Render(50)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(60).
		Render(content)

	return lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box)
}

// renderAchievements は実績一覧をレンダリングします。

func (s *StatsAchievementsScreen) renderAchievements() string {
	var builder strings.Builder

	// コンプリート率
	rate := s.getAchievementCompletionRate()
	rateStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ColorSecondary)

	rateText := fmt.Sprintf("達成率: %d%% (%d/%d)", rate, s.countAchieved(), len(s.data.Achievements))
	builder.WriteString(lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(rateStyle.Render(rateText)))
	builder.WriteString("\n\n")

	// 実績リスト
	achievementList := s.renderAchievementList()
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1).
		Width(80).
		Render(achievementList)

	builder.WriteString(lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center).
		Render(box))

	return builder.String()
}

// renderAchievementList は実績リストをレンダリングします。
func (s *StatsAchievementsScreen) renderAchievementList() string {
	var items []string

	for i, ach := range s.data.Achievements {
		style := lipgloss.NewStyle()
		if i == s.selectedIndex {
			style = style.Bold(true).Foreground(styles.ColorPrimary)
		} else if !ach.Achieved {
			style = style.Foreground(styles.ColorSubtle)
		}

		// ステータスアイコン
		status := "[ ]"
		if ach.Achieved {
			status = "[x]"
		}

		line := fmt.Sprintf("%s %s - %s", status, ach.Name, ach.Description)
		items = append(items, style.Render(line))
	}

	return strings.Join(items, "\n")
}

// countAchieved は達成済み実績の数を返します。
func (s *StatsAchievementsScreen) countAchieved() int {
	count := 0
	for _, ach := range s.data.Achievements {
		if ach.Achieved {
			count++
		}
	}
	return count
}
