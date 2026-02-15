package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TipsScreen はTIPS一覧画面です。
// 閲覧可能なチュートリアルの一覧を表示し、選択するとチュートリアル画面に遷移します。
type TipsScreen struct {
	tutorials   []domain.TutorialDef
	selectedIdx int
	styles      *styles.GameStyles
}

// NewTipsScreen は新しいTIPS画面を作成します。
func NewTipsScreen(tutorials []domain.TutorialDef) *TipsScreen {
	return &TipsScreen{
		tutorials:   tutorials,
		selectedIdx: 0,
		styles:      styles.NewGameStyles(),
	}
}

// Init は画面の初期化コマンドを返します。
func (s *TipsScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理し、画面の状態を更新します。
func (s *TipsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}
	return s, nil
}

func (s *TipsScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if s.selectedIdx < len(s.tutorials)-1 {
			s.selectedIdx++
		}
	case "k", "up":
		if s.selectedIdx > 0 {
			s.selectedIdx--
		}
	case "enter":
		if len(s.tutorials) > 0 {
			tut := s.tutorials[s.selectedIdx]
			return s, func() tea.Msg {
				return OpenTutorialMsg{Tutorial: tut, Mode: TipsView}
			}
		}
	case "esc":
		return s, func() tea.Msg {
			return ChangeSceneMsg{Scene: "home"}
		}
	case "q":
		return s, tea.Quit
	}
	return s, nil
}

// View は画面を描画します。
func (s *TipsScreen) View() string {
	gs := s.styles
	var sb strings.Builder

	title := gs.Text.Title.Render("📚 TIPS")
	sb.WriteString(title)
	sb.WriteString("\n\n")

	if len(s.tutorials) == 0 {
		sb.WriteString(gs.Text.Subtle.Render("表示可能なTIPSはありません"))
	} else {
		selectedStyle := lipgloss.NewStyle().Foreground(styles.ColorPrimary).Bold(true)
		for i, tut := range s.tutorials {
			cursor := "  "
			if i == s.selectedIdx {
				cursor = "▸ "
			}
			line := fmt.Sprintf("%s%s", cursor, tut.Title)
			if i == s.selectedIdx {
				sb.WriteString(selectedStyle.Render(line))
			} else {
				sb.WriteString(line)
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(gs.Text.Subtle.Render("Enter: 閲覧  Esc: ホームに戻る"))

	return sb.String()
}
