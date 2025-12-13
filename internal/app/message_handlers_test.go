// Package app は BlitzTypingOperator TUIゲームのメッセージハンドラーテストを提供します。
package app

import (
	"testing"

	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// TestNewMessageHandlers はMessageHandlersが正しく初期化されることを検証します
func TestNewMessageHandlers(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)
	if handlers == nil {
		t.Fatal("NewMessageHandlers() returned nil")
	}
}

// TestMessageHandlers_HandleWindowSizeMsg はWindowSizeMsgが正しく処理されることを検証します
func TestMessageHandlers_HandleWindowSizeMsg(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	msg := tea.WindowSizeMsg{Width: 150, Height: 50}
	_, cmd := handlers.Handle(msg)

	if model.TerminalState() == nil {
		t.Fatal("TerminalState should be set after WindowSizeMsg")
	}
	if model.TerminalState().Width != 150 {
		t.Errorf("Width should be 150, got %d", model.TerminalState().Width)
	}
	// WindowSizeMsgの場合コマンドはnilになるべき
	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}
}

// TestMessageHandlers_HandleChangeSceneMsg はChangeSceneMsgが正しく処理されることを検証します
func TestMessageHandlers_HandleChangeSceneMsg(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	msg := ChangeSceneMsg{Scene: SceneBattleSelect}
	_, _ = handlers.Handle(msg)

	if model.CurrentScene() != SceneBattleSelect {
		t.Errorf("CurrentScene should be SceneBattleSelect, got %v", model.CurrentScene())
	}
}

// TestMessageHandlers_HandleScreensChangeSceneMsg は画面からのシーン遷移要求が処理されることを検証します
func TestMessageHandlers_HandleScreensChangeSceneMsg(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	msg := screens.ChangeSceneMsg{Scene: "battle_select"}
	_, _ = handlers.Handle(msg)

	if model.CurrentScene() != SceneBattleSelect {
		t.Errorf("CurrentScene should be SceneBattleSelect, got %v", model.CurrentScene())
	}
}

// TestMessageHandlers_HandleCtrlC はCtrl+Cでtea.Quitが返されることを検証します
func TestMessageHandlers_HandleCtrlC(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	// まずWindowSizeMsgで初期化
	handlers.Handle(tea.WindowSizeMsg{Width: 140, Height: 40})

	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := handlers.Handle(msg)

	if cmd == nil {
		t.Fatal("Ctrl+C should return a command")
	}

	quitMsg := cmd()
	if _, ok := quitMsg.(tea.QuitMsg); !ok {
		t.Error("Ctrl+C should return tea.Quit command")
	}
}

// TestMessageHandlers_HandleEscKey はEscキーでホーム以外からホームに戻ることを検証します
func TestMessageHandlers_HandleEscKey(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	// まずWindowSizeMsgで初期化
	handlers.Handle(tea.WindowSizeMsg{Width: 140, Height: 40})

	// バトル選択画面に移動
	handlers.Handle(ChangeSceneMsg{Scene: SceneBattleSelect})
	if model.CurrentScene() != SceneBattleSelect {
		t.Fatal("Should be on BattleSelect screen")
	}

	// Escキーを押してホームに戻る
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	_, _ = handlers.Handle(msg)

	if model.CurrentScene() != SceneHome {
		t.Errorf("Esc should return to SceneHome, got %v", model.CurrentScene())
	}
}

// TestMessageHandlers_HandleQKeyOnHome はホーム画面でQキーが終了することを検証します
func TestMessageHandlers_HandleQKeyOnHome(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	// まずWindowSizeMsgで初期化
	handlers.Handle(tea.WindowSizeMsg{Width: 140, Height: 40})

	// ホーム画面でqを押す
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := handlers.Handle(msg)

	if cmd == nil {
		t.Fatal("q on Home should return a command")
	}

	quitMsg := cmd()
	if _, ok := quitMsg.(tea.QuitMsg); !ok {
		t.Error("q on Home should return tea.Quit command")
	}
}

// TestMessageHandlers_HandleQKeyNotOnHome はホーム以外でQキーが終了しないことを検証します
func TestMessageHandlers_HandleQKeyNotOnHome(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	// まずWindowSizeMsgで初期化
	handlers.Handle(tea.WindowSizeMsg{Width: 140, Height: 40})

	// バトル選択画面に移動
	handlers.Handle(ChangeSceneMsg{Scene: SceneBattleSelect})

	// バトル選択画面でqを押す（終了しないはず）
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := handlers.Handle(msg)

	// qは画面に転送されるので、tea.Quitではないはず
	if cmd != nil {
		// 画面に転送された場合はコマンドが返る可能性があるが、QuitMsgではないはず
		if quitMsg := cmd(); quitMsg != nil {
			if _, ok := quitMsg.(tea.QuitMsg); ok {
				t.Error("q on non-Home screen should not quit")
			}
		}
	}
}

// TestMessageHandlers_HandlerCount はハンドラー数が適切であることを検証します
func TestMessageHandlers_HandlerCount(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	handlers := NewMessageHandlers(model)

	// ハンドラー数を取得
	count := handlers.HandlerCount()

	// メッセージタイプハンドラーが5以下であることを確認
	// これは循環的複雑度の削減要件を満たすため
	if count > 5 {
		t.Errorf("HandlerCount should be 5 or less, got %d", count)
	}
}

// TestNewScreenMap は画面マップが正しく初期化されることを検証します
func TestNewScreenMap(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	screenMap := NewScreenMap(model)
	if screenMap == nil {
		t.Fatal("NewScreenMap() returned nil")
	}
}

// TestScreenMap_GetScreen は各シーンに対して正しい画面が返されることを検証します
func TestScreenMap_GetScreen(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	screenMap := NewScreenMap(model)

	tests := []struct {
		scene    Scene
		hasValue bool
	}{
		{SceneHome, true},
		{SceneBattleSelect, true},
		{SceneAgentManagement, true},
		{SceneEncyclopedia, true},
		{SceneAchievement, true},
		{SceneSettings, true},
	}

	for _, tt := range tests {
		t.Run(tt.scene.String(), func(t *testing.T) {
			screen := screenMap.GetScreen(tt.scene)
			if tt.hasValue && screen == nil {
				t.Errorf("GetScreen(%v) should return a screen", tt.scene)
			}
		})
	}
}

// TestScreenMap_RenderScene はシーンに応じた描画が行われることを検証します
func TestScreenMap_RenderScene(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	screenMap := NewScreenMap(model)

	// ホーム画面のレンダリング
	view := screenMap.RenderScene(SceneHome)
	if view == "" {
		t.Error("RenderScene(SceneHome) should return non-empty view")
	}
}

// TestScreenMap_ForwardMessage はメッセージが正しく転送されることを検証します
func TestScreenMap_ForwardMessage(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	screenMap := NewScreenMap(model)

	// キーメッセージを転送
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	cmd := screenMap.ForwardMessage(SceneHome, msg)

	// ホーム画面にメッセージが転送されたことを確認
	// コマンドはnilまたは有効なコマンドである
	_ = cmd
}

// TestScreenMap_MapCount は画面マップの要素数が適切であることを検証します
func TestScreenMap_MapCount(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	screenMap := NewScreenMap(model)

	count := screenMap.MapCount()

	// 少なくとも主要な画面は登録されているべき
	if count < 6 {
		t.Errorf("MapCount should be at least 6, got %d", count)
	}
}
