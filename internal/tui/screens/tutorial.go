package screens

import (
	"fmt"
	"strings"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TutorialMode はチュートリアル画面の動作モードを表します。
type TutorialMode int

const (
	// UnlockFlow は報酬導線モードです。Enter/Escいずれでも機能がUnlockedに遷移します。
	UnlockFlow TutorialMode = iota
	// TipsView はTIPS再閲覧モードです。状態遷移は発生しません。
	TipsView
)

// CompleteTutorialMsg はチュートリアル完了を通知するメッセージです。
// app層でハンドリングし、機能解放状態の遷移を行います。
type CompleteTutorialMsg struct {
	TutorialID string
}

// OpenTutorialMsg はチュートリアル画面を開くことを要求するメッセージです。
type OpenTutorialMsg struct {
	Tutorial domain.TutorialDef
	Mode     TutorialMode
}

// TutorialScreen はチュートリアル表示画面です。
type TutorialScreen struct {
	tutorial    domain.TutorialDef
	mode        TutorialMode
	currentPage int
	styles      *styles.GameStyles
}

// NewTutorialScreen は新しいチュートリアル画面を作成します。
func NewTutorialScreen(tutorial domain.TutorialDef, mode TutorialMode) *TutorialScreen {
	return &TutorialScreen{
		tutorial:    tutorial,
		mode:        mode,
		currentPage: 0,
		styles:      styles.NewGameStyles(),
	}
}

// Init は画面の初期化コマンドを返します。
func (s *TutorialScreen) Init() tea.Cmd {
	return nil
}

// Update はメッセージを処理し、画面の状態を更新します。
func (s *TutorialScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return s.handleKeyMsg(msg)
	}
	return s, nil
}

func (s *TutorialScreen) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", " ":
		return s.handleNext()
	case "right", "l":
		return s.handlePageForward()
	case "left", "h":
		return s.handlePrev()
	case "esc":
		return s.handleEsc()
	case "q":
		return s, tea.Quit
	}
	return s, nil
}

func (s *TutorialScreen) handleNext() (tea.Model, tea.Cmd) {
	if s.currentPage < len(s.tutorial.Pages)-1 {
		s.currentPage++
		return s, nil
	}
	// 最終ページ
	return s.handleComplete()
}

// handlePageForward はページ送りのみ行います（最終ページでは何もしない）。
func (s *TutorialScreen) handlePageForward() (tea.Model, tea.Cmd) {
	if s.currentPage < len(s.tutorial.Pages)-1 {
		s.currentPage++
	}
	return s, nil
}

func (s *TutorialScreen) handlePrev() (tea.Model, tea.Cmd) {
	if s.currentPage > 0 {
		s.currentPage--
	}
	return s, nil
}

func (s *TutorialScreen) handleEsc() (tea.Model, tea.Cmd) {
	if s.mode == UnlockFlow {
		// 報酬導線: Escでもチュートリアル完了として扱う
		return s.handleComplete()
	}
	// TIPS再閲覧: TIPS一覧に戻る
	return s, func() tea.Msg {
		return ChangeSceneMsg{Scene: "tips"}
	}
}

func (s *TutorialScreen) handleComplete() (tea.Model, tea.Cmd) {
	if s.mode == UnlockFlow {
		return s, func() tea.Msg {
			return CompleteTutorialMsg{TutorialID: s.tutorial.ID}
		}
	}
	// TipsViewモード: TIPS一覧に戻る
	return s, func() tea.Msg {
		return ChangeSceneMsg{Scene: "tips"}
	}
}

// View は画面を描画します。
func (s *TutorialScreen) View() string {
	gs := s.styles
	var sb strings.Builder

	// タイトル
	title := gs.Text.Title.Render(fmt.Sprintf("📖 %s", s.tutorial.Title))
	sb.WriteString(title)
	sb.WriteString("\n\n")

	// ページ内容
	if s.currentPage < len(s.tutorial.Pages) {
		content := s.tutorial.Pages[s.currentPage]
		contentStyle := lipgloss.NewStyle().
			Width(70).
			PaddingLeft(2)
		sb.WriteString(contentStyle.Render(content))
	}
	sb.WriteString("\n\n")

	// ページインジケーター
	totalPages := len(s.tutorial.Pages)
	if totalPages > 1 {
		pageInfo := fmt.Sprintf("ページ %d / %d", s.currentPage+1, totalPages)
		sb.WriteString(gs.Text.Subtle.Render(pageInfo))
		sb.WriteString("\n")
	}

	// 操作ガイド
	sb.WriteString("\n")
	if s.mode == UnlockFlow {
		if s.currentPage < totalPages-1 {
			sb.WriteString(gs.Text.Subtle.Render("Enter/→: 次へ  ←: 前へ"))
		} else {
			sb.WriteString(gs.Text.Subtle.Render("Enter: 完了  ←: 前へ"))
		}
	} else {
		if s.currentPage < totalPages-1 {
			sb.WriteString(gs.Text.Subtle.Render("Enter/→: 次へ  ←: 前へ  Esc: 戻る"))
		} else {
			sb.WriteString(gs.Text.Subtle.Render("Enter/Esc: 戻る  ←: 前へ"))
		}
	}

	return sb.String()
}
