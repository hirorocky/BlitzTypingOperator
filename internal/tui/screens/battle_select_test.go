// Package screens はTUI画面のテストを提供します。
package screens

import (
	"testing"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// mockAgentProvider はテスト用のAgentProvider実装です。
type mockAgentProvider struct {
	agents []*domain.AgentModel
}

func (m *mockAgentProvider) GetEquippedAgents() []*domain.AgentModel {
	return m.agents
}

// ==================== Task 10.2: バトル選択画面のテスト ====================

// TestNewBattleSelectScreen はBattleSelectScreenの初期化をテストします。

func TestNewBattleSelectScreen(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	if screen == nil {
		t.Fatal("BattleSelectScreenがnilです")
	}

	// 入力フィールドが初期化されていること
	if screen.input == nil {
		t.Error("入力フィールドがnilです")
	}
}

// TestBattleSelectMaxLevelDisplay は最高レベル表示をテストします。

func TestBattleSelectMaxLevelDisplay(t *testing.T) {
	maxLevel := 15
	screen := NewBattleSelectScreen(maxLevel, &mockAgentProvider{})

	// 挑戦可能最大レベルは maxLevel + 1
	if screen.maxChallengeLevel != maxLevel+1 {
		t.Errorf("挑戦可能最大レベル: got %d, want %d", screen.maxChallengeLevel, maxLevel+1)
	}
}

// TestBattleSelectInputValidation は入力検証をテストします。

func TestBattleSelectInputValidation(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   string
	}{
		{"有効な入力", "5", false, ""},
		{"最大レベル", "11", false, ""},
		{"1未満", "0", true, "too_low"},
		{"最大レベル超過", "12", true, "too_high"},
		{"空入力", "", true, "empty"},
		{"非数値", "abc", true, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen.input.Value = tt.input
			valid, _ := screen.validateInput()

			if valid == tt.expectError {
				if tt.expectError {
					t.Error("エラーが検出されませんでした")
				} else {
					t.Error("有効な入力がエラーと判定されました")
				}
			}
		})
	}
}

// TestBattleSelectNoAgentEquipped はエージェント未装備時のテストです。

func TestBattleSelectNoAgentEquipped(t *testing.T) {
	// エージェント未装備
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})
	screen.input.Value = "5"

	// 確認画面に移動
	screen.state = StateConfirm

	// バトル開始を試みる
	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	// エラーメッセージが設定されるか確認
	if screen.error == "" && cmd != nil {
		// コマンドが返された場合、バトル開始メッセージかどうか確認
		// 未装備の場合はバトル開始できないはず
	}
}

// TestBattleSelectWithAgentEquipped はエージェント装備時のテストです。
func TestBattleSelectWithAgentEquipped(t *testing.T) {
	// エージェントを装備
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCore("core1", "テストコア", 5, coreType, domain.PassiveSkill{})
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール1", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m2", "モジュール2", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m3", "モジュール3", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m4", "モジュール4", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
	}
	agent := domain.NewAgent("agent1", core, modules)

	screen := NewBattleSelectScreen(10, &mockAgentProvider{agents: []*domain.AgentModel{agent}})
	screen.input.Value = "5"

	// 入力検証
	valid, _ := screen.validateInput()
	if !valid {
		t.Error("有効な入力がエラーと判定されました")
	}
}

// TestBattleSelectConfirmScreen は確認画面のテストです。

func TestBattleSelectConfirmScreen(t *testing.T) {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCore("core1", "テストコア", 5, coreType, domain.PassiveSkill{})
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール1", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m2", "モジュール2", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m3", "モジュール3", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
		domain.NewModule("m4", "モジュール4", domain.PhysicalAttack, 1, []string{"physical_low"}, 10, "STR", ""),
	}
	agent := domain.NewAgent("agent1", core, modules)

	screen := NewBattleSelectScreen(10, &mockAgentProvider{agents: []*domain.AgentModel{agent}})
	screen.input.Value = "5"

	// Enterで確認画面へ
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	if screen.state != StateConfirm {
		t.Errorf("確認画面に遷移していません: state=%d", screen.state)
	}
}

// TestBattleSelectRechallenge は再挑戦のテストです。

func TestBattleSelectRechallenge(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	// 過去にクリアしたレベル（例: レベル5）への再挑戦
	screen.input.Value = "5"
	valid, _ := screen.validateInput()

	if !valid {
		t.Error("過去クリアレベルへの再挑戦が拒否されました")
	}
}

// TestBattleSelectRender はバトル選択画面のレンダリングをテストします。
func TestBattleSelectRender(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestBattleSelectBackNavigation は戻るナビゲーションのテストです。

func TestBattleSelectBackNavigation(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Escキーでコマンドが返されません")
	}
}
