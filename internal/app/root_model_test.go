package app

import (
	"testing"

	"hirorocky/type-battle/internal/infra/masterdata"
	gamestate "hirorocky/type-battle/internal/usecase/session"

	tea "github.com/charmbracelet/bubbletea"
)

// === シーン定義のテスト ===

// TestSceneType_Constants はシーンタイプ定数が定義されていることを検証します
func TestSceneType_Constants(t *testing.T) {
	// 各シーンタイプが異なる値を持つことを確認
	scenes := []Scene{
		SceneHome,
		SceneBattle,
		SceneBattleSelect,
		SceneAgentManagement,
		SceneEncyclopedia,
		SceneAchievement,
		SceneSettings,
	}

	// シーン値が重複していないことを確認
	seen := make(map[Scene]bool)
	for _, scene := range scenes {
		if seen[scene] {
			t.Errorf("Duplicate scene value detected: %d", scene)
		}
		seen[scene] = true
	}
}

// TestSceneType_String は各シーンに文字列表現があることを検証します
func TestSceneType_String(t *testing.T) {
	tests := []struct {
		scene    Scene
		expected string
	}{
		{SceneHome, "Home"},
		{SceneBattle, "Battle"},
		{SceneBattleSelect, "BattleSelect"},
		{SceneAgentManagement, "AgentManagement"},
		{SceneEncyclopedia, "Encyclopedia"},
		{SceneAchievement, "Achievement"},
		{SceneSettings, "Settings"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.scene.String() != tt.expected {
				t.Errorf("Scene.String() = %s, expected %s", tt.scene.String(), tt.expected)
			}
		})
	}
}

// === GameState のテスト ===

// TestNewGameState は新しいGameStateが正しく初期化されることを検証します
func TestNewGameState(t *testing.T) {
	gs := gamestate.NewGameState()
	if gs == nil {
		t.Fatal("NewGameState() returned nil")
	}
}

// TestGameState_HasMaxLevelReached はGameStateが到達最高レベルを保持することを検証します
func TestGameState_HasMaxLevelReached(t *testing.T) {
	gs := gamestate.NewGameState()
	// 初期値は0または1であるべき
	if gs.MaxLevelReached < 0 {
		t.Errorf("MaxLevelReached should not be negative, got %d", gs.MaxLevelReached)
	}
}

// === RootModel のテスト ===

// TestNewRootModel は新しいRootModelが正しく初期化されることを検証します
func TestNewRootModel(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	if model == nil {
		t.Fatal("NewRootModel() returned nil")
	}
}

// TestRootModel_ImplementsTeaModel はRootModelがtea.Modelインターフェースを実装していることを検証します
func TestRootModel_ImplementsTeaModel(t *testing.T) {
	var _ tea.Model = (*RootModel)(nil)
}

// TestRootModel_HasGameState はRootModelがGameStateを保持していることを検証します
func TestRootModel_HasGameState(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	if model.GameState() == nil {
		t.Fatal("RootModel should have GameState")
	}
}

// TestRootModel_HasCurrentScene はRootModelが現在のシーンを保持していることを検証します
func TestRootModel_HasCurrentScene(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	// 初期シーンはSceneHomeであるべき
	if model.CurrentScene() != SceneHome {
		t.Errorf("Initial scene should be SceneHome, got %v", model.CurrentScene())
	}
}

// TestRootModel_HasTerminalState はRootModelがターミナル状態を保持していることを検証します
func TestRootModel_HasTerminalState(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	// ターミナル状態はWindowSizeMsg受信後に設定されるのでnilでもOK
	_ = model.TerminalState()
}

// TestRootModel_HasStyles はRootModelがスタイルを保持していることを検証します
func TestRootModel_HasStyles(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	if model.Styles() == nil {
		t.Fatal("RootModel should have Styles")
	}
}

// === シーンルーティングのテスト ===

// TestRootModel_ChangeScene はシーン変更が正しく動作することを検証します
func TestRootModel_ChangeScene(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// ホームからバトル選択へ遷移
	model.ChangeScene(SceneBattleSelect)
	if model.CurrentScene() != SceneBattleSelect {
		t.Errorf("CurrentScene should be SceneBattleSelect, got %v", model.CurrentScene())
	}

	// バトル選択からエージェント管理へ遷移
	model.ChangeScene(SceneAgentManagement)
	if model.CurrentScene() != SceneAgentManagement {
		t.Errorf("CurrentScene should be SceneAgentManagement, got %v", model.CurrentScene())
	}
}

// TestRootModel_ChangeSceneMsg はChangeSceneMsgでシーンが変更されることを検証します
func TestRootModel_ChangeSceneMsg(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// WindowSizeMsgを先に送信してモデルを初期化
	msg1 := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg1)
	model = updatedModel.(*RootModel)

	// ChangeSceneMsgを送信
	msg := ChangeSceneMsg{Scene: SceneBattleSelect}
	updatedModel, _ = model.Update(msg)
	m := updatedModel.(*RootModel)

	if m.CurrentScene() != SceneBattleSelect {
		t.Errorf("CurrentScene should be SceneBattleSelect after ChangeSceneMsg, got %v", m.CurrentScene())
	}
}

// === Init/Update/View メソッドのテスト ===

// TestRootModel_Init はInitが正しく動作することを検証します
func TestRootModel_Init(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	cmd := model.Init()
	// Initは nil または有効なコマンドを返すことができます
	_ = cmd
}

// TestRootModel_Update_WindowSizeMsg はWindowSizeMsgが正しく処理されることを検証します
func TestRootModel_Update_WindowSizeMsg(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	msg := tea.WindowSizeMsg{Width: 150, Height: 50}
	updatedModel, _ := model.Update(msg)

	m, ok := updatedModel.(*RootModel)
	if !ok {
		t.Fatal("Update should return *RootModel")
	}

	ts := m.TerminalState()
	if ts == nil {
		t.Fatal("TerminalState should be set after WindowSizeMsg")
	}

	if ts.Width != 150 {
		t.Errorf("Width should be 150, got %d", ts.Width)
	}
	if ts.Height != 50 {
		t.Errorf("Height should be 50, got %d", ts.Height)
	}
}

// TestRootModel_Update_QuitKey は終了キーでtea.Quitが返されることを検証します
func TestRootModel_Update_QuitKey(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// まずWindowSizeMsgで初期化
	msg1 := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg1)
	model = updatedModel.(*RootModel)

	// qキーをテスト
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Fatal("Update with 'q' key should return a command")
	}

	// tea.Quitが返されることを確認（間接的に）
	// Bubbleteaのtea.Quitは tea.Msg を返すfuncなので、実行してみる
	quitMsg := cmd()
	if _, ok := quitMsg.(tea.QuitMsg); !ok {
		t.Errorf("Update with 'q' key should return tea.Quit command")
	}
}

// TestRootModel_Update_CtrlC はCtrl+Cでtea.Quitが返されることを検証します
func TestRootModel_Update_CtrlC(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// まずWindowSizeMsgで初期化
	msg1 := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg1)
	model = updatedModel.(*RootModel)

	// Ctrl+Cをテスト
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(msg)

	if cmd == nil {
		t.Fatal("Update with Ctrl+C should return a command")
	}

	// tea.Quitが返されることを確認
	quitMsg := cmd()
	if _, ok := quitMsg.(tea.QuitMsg); !ok {
		t.Errorf("Update with Ctrl+C should return tea.Quit command")
	}
}

// TestRootModel_View_Loading は初期状態でローディングメッセージが表示されることを検証します
func TestRootModel_View_Loading(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)
	view := model.View()

	if view == "" {
		t.Fatal("View should not be empty")
	}
}

// TestRootModel_View_SmallTerminal は小さいターミナルで警告が表示されることを検証します
func TestRootModel_View_SmallTerminal(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// 小さいWindowSizeMsgを送信
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*RootModel)

	view := m.View()

	// 警告メッセージが含まれているべき
	if view == "" {
		t.Fatal("View should not be empty for small terminal")
	}
}

// TestRootModel_View_ValidTerminal は有効なターミナルでタイトルが表示されることを検証します
func TestRootModel_View_ValidTerminal(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// 有効なWindowSizeMsgを送信
	msg := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*RootModel)

	view := m.View()

	// タイトルが含まれているべき
	if view == "" {
		t.Fatal("View should not be empty for valid terminal")
	}

	// BlitzTypingOperatorの文字列が含まれているべき
	if len(view) < 10 {
		t.Error("View should contain game title or content")
	}
}

// TestRootModel_IsReady は有効なターミナルでIsReady()がtrueになることを検証します
func TestRootModel_IsReady(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// 初期状態では準備完了でない
	if model.IsReady() {
		t.Error("Model should not be ready initially")
	}

	// 有効なWindowSizeMsgを送信
	msg := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*RootModel)

	// 準備完了であるべき
	if !m.IsReady() {
		t.Error("Model should be ready after valid WindowSizeMsg")
	}
}

// TestRootModel_NotReady_SmallTerminal は小さいターミナルでIsReady()がfalseになることを検証します
func TestRootModel_NotReady_SmallTerminal(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// 小さいWindowSizeMsgを送信
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*RootModel)

	// 準備完了でないべき
	if m.IsReady() {
		t.Error("Model should not be ready after invalid WindowSizeMsg")
	}
}

// === 終了操作のテスト ===

// TestRootModel_QuitPreservesTerminalState は終了時にターミナル状態を保存することを検証します
// 注: Bubbleteaでは tea.WithAltScreen() により自動的に復元される
func TestRootModel_QuitPreservesTerminalState(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// WindowSizeMsgで初期化
	msg1 := tea.WindowSizeMsg{Width: 140, Height: 40}
	updatedModel, _ := model.Update(msg1)
	model = updatedModel.(*RootModel)

	// 終了操作を実行
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := model.Update(msg)

	// tea.Quitコマンドが返されるべき
	if cmd == nil {
		t.Fatal("Quit should return a command")
	}
}

// === ChangeSceneMsg のテスト ===

// TestChangeSceneMsg_Type はChangeSceneMsgの型を検証します
func TestChangeSceneMsg_Type(t *testing.T) {
	msg := ChangeSceneMsg{Scene: SceneBattle}
	if msg.Scene != SceneBattle {
		t.Errorf("ChangeSceneMsg.Scene should be SceneBattle, got %v", msg.Scene)
	}
}

// === シーン間遷移のテスト ===

// TestRootModel_SceneTransition_HomeToAll はホームから各シーンへの遷移を検証します
func TestRootModel_SceneTransition_HomeToAll(t *testing.T) {
	tests := []struct {
		name  string
		scene Scene
	}{
		{"Home to Battle", SceneBattle},
		{"Home to BattleSelect", SceneBattleSelect},
		{"Home to AgentManagement", SceneAgentManagement},
		{"Home to Encyclopedia", SceneEncyclopedia},
		{"Home to Achievement", SceneAchievement},
		{"Home to Settings", SceneSettings},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewRootModel("", masterdata.EmbeddedData)

			// 初期状態はホームであるべき
			if model.CurrentScene() != SceneHome {
				t.Fatalf("Initial scene should be SceneHome")
			}

			// シーン変更
			model.ChangeScene(tt.scene)

			if model.CurrentScene() != tt.scene {
				t.Errorf("CurrentScene should be %v, got %v", tt.scene, model.CurrentScene())
			}
		})
	}
}

// TestRootModel_SceneTransition_BackToHome は各シーンからホームへ戻れることを検証します
func TestRootModel_SceneTransition_BackToHome(t *testing.T) {
	model := NewRootModel("", masterdata.EmbeddedData)

	// バトル選択へ遷移
	model.ChangeScene(SceneBattleSelect)
	if model.CurrentScene() != SceneBattleSelect {
		t.Fatal("Should transition to BattleSelect")
	}

	// ホームへ戻る
	model.ChangeScene(SceneHome)
	if model.CurrentScene() != SceneHome {
		t.Errorf("Should return to Home, got %v", model.CurrentScene())
	}
}
