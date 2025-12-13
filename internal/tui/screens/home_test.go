// Package screens はTUI画面のテストを提供します。
package screens

import (
	"testing"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.1: ホーム画面のテスト ====================

// TestNewHomeScreen はHomeScreenの初期化をテストします。

func TestNewHomeScreen(t *testing.T) {
	screen := NewHomeScreen(0, nil)

	if screen == nil {
		t.Fatal("HomeScreenがnilです")
	}

	// 初期状態で4つのメニューアイテムがあること

	if len(screen.menu.Items) != 5 { // 4メニュー + 設定
		t.Errorf("メニューアイテム数が不正: got %d, want 5", len(screen.menu.Items))
	}
}

// TestHomeScreenMenuItems はメニューアイテムをテストします。

func TestHomeScreenMenuItems(t *testing.T) {
	screen := NewHomeScreen(0, nil)

	expectedItems := []string{
		"agent_management",
		"battle_select",
		"encyclopedia",
		"stats_achievements",
		"settings",
	}

	for i, expected := range expectedItems {
		if i >= len(screen.menu.Items) {
			t.Errorf("メニューアイテム%dが存在しません", i)
			continue
		}
		if screen.menu.Items[i].Value != expected {
			t.Errorf("メニューアイテム%d: got %s, want %s", i, screen.menu.Items[i].Value, expected)
		}
	}
}

// TestHomeScreenNavigation はメニューナビゲーションをテストします。

func TestHomeScreenNavigation(t *testing.T) {
	// モックAgentProviderを使用してバトル選択を有効化
	screen := NewHomeScreen(0, &mockAgentProvider{agents: []*domain.AgentModel{{Level: 1}}})

	// 下キーで移動
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyDown})
	if screen.menu.SelectedIndex != 1 {
		t.Errorf("下キー後の選択インデックス: got %d, want 1", screen.menu.SelectedIndex)
	}

	// jキーで移動（vim形式）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if screen.menu.SelectedIndex != 2 {
		t.Errorf("jキー後の選択インデックス: got %d, want 2", screen.menu.SelectedIndex)
	}

	// 上キーで移動
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyUp})
	if screen.menu.SelectedIndex != 1 {
		t.Errorf("上キー後の選択インデックス: got %d, want 1", screen.menu.SelectedIndex)
	}

	// kキーで移動（vim形式）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if screen.menu.SelectedIndex != 0 {
		t.Errorf("kキー後の選択インデックス: got %d, want 0", screen.menu.SelectedIndex)
	}
}

// TestHomeScreenEnterSelection はEnterキーによる選択をテストします。

func TestHomeScreenEnterSelection(t *testing.T) {
	screen := NewHomeScreen(0, nil)

	// バトル選択を選択
	screen.menu.SelectedIndex = 1 // battle_select

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enterキーでコマンドが返されません")
	}
}

// TestHomeScreenProgressDisplay は進行状況表示をテストします。

func TestHomeScreenProgressDisplay(t *testing.T) {
	maxLevel := 15
	screen := NewHomeScreen(maxLevel, nil)

	rendered := screen.View()

	// レンダリング結果が空でないこと
	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}

	// 最高レベルの情報が含まれていることを確認
	// （実際のレンダリング内容はUI実装に依存）
}

// TestHomeScreenRender はホーム画面のレンダリングをテストします。
func TestHomeScreenRender(t *testing.T) {
	screen := NewHomeScreen(10, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// ==================== Task 4.1: ASCIIアートロゴ統合のテスト ====================

// TestHomeScreenHasASCIILogo はASCIIロゴが表示されることをテストします。

func TestHomeScreenHasASCIILogo(t *testing.T) {
	screen := NewHomeScreen(10, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// ロゴの特徴的な文字が含まれることを確認（フィグレット風）
	// BlitzTypingOperatorロゴは「╔╦╗」などの文字を使用
	if !containsAny(rendered, "╔", "╗", "╚", "╝") {
		t.Error("ASCIIアートロゴが表示されていません")
	}
}

// TestHomeScreenHasLevelASCII はレベルがASCII数字で表示されることをテストします。

func TestHomeScreenHasLevelASCII(t *testing.T) {
	screen := NewHomeScreen(15, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// ASCII数字の特徴的な文字が含まれることを確認
	// 数字は「█」を使用
	if !containsAny(rendered, "█") {
		t.Error("ASCII数字が表示されていません")
	}
}

// TestHomeScreenHasSubtitle はサブタイトルが表示されることをテストします。
func TestHomeScreenHasSubtitle(t *testing.T) {
	screen := NewHomeScreen(10, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// サブタイトルが含まれることを確認
	if !containsAny(rendered, "Terminal Typing Battle Game") {
		t.Error("サブタイトルが表示されていません")
	}
}

// containsAny は文字列に指定したいずれかの部分文字列が含まれるかを確認します。
func containsAny(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if contains(s, substr) {
			return true
		}
	}
	return false
}

// contains は文字列に部分文字列が含まれるかを確認します。
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

// findSubstring は文字列内で部分文字列を探します。
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ==================== Task 4.2: 左右分割レイアウトのテスト ====================

// TestHomeScreenHasLeftRightLayout は左右分割レイアウトをテストします。

func TestHomeScreenHasLeftRightLayout(t *testing.T) {
	screen := NewHomeScreen(10, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// メインメニューが含まれることを確認
	if !containsAny(rendered, "メインメニュー") {
		t.Error("メインメニューが表示されていません")
	}

	// 進行状況が含まれることを確認
	if !containsAny(rendered, "進行状況") {
		t.Error("進行状況パネルが表示されていません")
	}
}

// TestHomeScreenHasKeyHelp は操作キーヘルプが表示されることをテストします。

func TestHomeScreenHasKeyHelp(t *testing.T) {
	screen := NewHomeScreen(10, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 操作キーヘルプが含まれることを確認
	if !containsAny(rendered, "Enter", "選択", "終了") {
		t.Error("操作キーヘルプが表示されていません")
	}
}

// ==================== Task 4.3: 進行状況パネルのテスト ====================

// TestHomeScreenShowsEquippedAgentsWithCard は装備エージェント一覧がカード形式で表示されることをテストします。

func TestHomeScreenShowsEquippedAgentsWithCard(t *testing.T) {
	// mockAgentProviderを使用して装備エージェントがある状態を作成
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCore("core1", "テストコア", 5, coreType, domain.PassiveSkill{})
	agent := &domain.AgentModel{
		ID:    "agent1",
		Level: 5,
		Core:  core,
	}
	provider := &mockAgentProvider{agents: []*domain.AgentModel{agent}}
	screen := NewHomeScreen(10, provider)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 装備中エージェントセクションが含まれることを確認
	if !containsAny(rendered, "装備中エージェント") {
		t.Error("装備中エージェントセクションが表示されていません")
	}

	// スロット表示が含まれることを確認
	if !containsAny(rendered, "スロット") {
		t.Error("スロット表示がありません")
	}
}

// TestHomeScreenShowsMaxLevel は到達最高レベルセクションが表示されることをテストします。

func TestHomeScreenShowsMaxLevel(t *testing.T) {
	screen := NewHomeScreen(15, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 到達最高レベルセクションが含まれることを確認
	if !containsAny(rendered, "到達最高レベル") {
		t.Error("到達最高レベルセクションが表示されていません")
	}
}

// TestHomeScreenEmptySlots はエージェント未装備時に空スロットが表示されることをテストします。
func TestHomeScreenEmptySlots(t *testing.T) {
	// 1つのエージェントのみ装備している状態（残り2スロットが空）
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCore("core1", "テストコア", 5, coreType, domain.PassiveSkill{})
	agent := &domain.AgentModel{
		ID:    "agent1",
		Level: 5,
		Core:  core,
	}
	provider := &mockAgentProvider{agents: []*domain.AgentModel{agent}}
	screen := NewHomeScreen(5, provider)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 空きスロット表示が含まれることを確認
	if !containsAny(rendered, "(空)", "(未装備)") {
		t.Error("空スロット表示がありません")
	}
}

// ==================== Task 4.4: 装備なし時の誘導メッセージとバトル無効化のテスト ====================

// TestHomeScreenBattleDisabledWhenNoAgent は装備エージェントがない場合にバトル選択が無効化されることをテストします。

func TestHomeScreenBattleDisabledWhenNoAgent(t *testing.T) {
	screen := NewHomeScreen(5, nil)
	screen.width = 120
	screen.height = 40

	// 装備エージェントがない場合、バトル選択メニューが無効化されていること
	for _, item := range screen.menu.Items {
		if item.Value == "battle_select" {
			if !item.Disabled {
				t.Error("バトル選択メニューが無効化されていません")
			}
			return
		}
	}
	t.Error("バトル選択メニューが見つかりません")
}

// TestHomeScreenGuidanceMessageWhenNoAgent は装備エージェントがない場合に誘導メッセージが表示されることをテストします。

func TestHomeScreenGuidanceMessageWhenNoAgent(t *testing.T) {
	screen := NewHomeScreen(5, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 誘導メッセージが含まれることを確認
	if !containsAny(rendered, "エージェント管理", "装備") {
		t.Error("誘導メッセージが表示されていません")
	}
}
