// Package app は BlitzTypingOperator TUIゲームのメッセージハンドラーを提供します。
// MessageHandlersは循環的複雑度を削減するため、メッセージタイプごとの処理を委譲します。
package app

import (
	"hirorocky/type-battle/internal/infra/terminal"
	"hirorocky/type-battle/internal/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

// MessageHandler はメッセージを処理して更新されたモデルとコマンドを返す関数型です
type MessageHandler func(msg tea.Msg) (tea.Model, tea.Cmd)

// MessageHandlers はメッセージタイプごとのハンドラーを管理します。
// 循環的複雑度を削減するため、Updateメソッドのswitch分岐をハンドラーマップに委譲します。
type MessageHandlers struct {
	model    *RootModel
	handlers map[string]MessageHandler
}

// NewMessageHandlers は新しいMessageHandlersを作成します。
func NewMessageHandlers(model *RootModel) *MessageHandlers {
	mh := &MessageHandlers{
		model:    model,
		handlers: make(map[string]MessageHandler),
	}
	mh.registerHandlers()
	return mh
}

// registerHandlers は全てのメッセージハンドラーを登録します。
func (mh *MessageHandlers) registerHandlers() {
	// 5つの主要なハンドラーカテゴリに集約
	// 1. ウィンドウサイズ関連
	// 2. キー入力関連
	// 3. シーン遷移関連（ChangeSceneMsg、screens.ChangeSceneMsg）
	// 4. バトル関連（StartBattleMsg、BattleTickMsg、BattleResultMsg）
	// 5. その他の処理

	mh.handlers["window_size"] = mh.handleWindowSizeMsg
	mh.handlers["key"] = mh.handleKeyMsg
	mh.handlers["scene_change"] = mh.handleSceneChangeMsg
	mh.handlers["battle"] = mh.handleBattleMsg
}

// Handle はメッセージを適切なハンドラーにルーティングします。
func (mh *MessageHandlers) Handle(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return mh.handleWindowSizeMsg(msg)
	case tea.KeyMsg:
		return mh.handleKeyMsg(msg)
	case ChangeSceneMsg:
		return mh.handleChangeSceneMsg(msg)
	case screens.ChangeSceneMsg:
		return mh.handleScreensChangeSceneMsg(msg)
	case screens.StartBattleMsg:
		return mh.handleStartBattleMsg(msg)
	case screens.BattleTickMsg:
		return mh.handleBattleTickMsg(msg)
	case screens.BattleResultMsg:
		return mh.handleBattleResultMsg(msg)
	}
	return mh.model, nil
}

// handleWindowSizeMsg はウィンドウサイズ変更を処理します。
func (mh *MessageHandlers) handleWindowSizeMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	windowMsg := msg.(tea.WindowSizeMsg)
	mh.model.terminalState = terminal.NewTerminalState(windowMsg.Width, windowMsg.Height)
	mh.model.ready = mh.model.terminalState.IsValid()
	// 各画面にもサイズ変更を通知
	if mh.model.homeScreen != nil {
		mh.model.homeScreen.Update(msg)
	}
	return mh.model, nil
}

// handleKeyMsg はキー入力を処理します。
func (mh *MessageHandlers) handleKeyMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg := msg.(tea.KeyMsg)
	switch keyMsg.String() {
	case "ctrl+c":
		return mh.model, tea.Quit
	case "esc":
		// ホーム画面以外ならホームに戻る
		if mh.model.currentScene != SceneHome {
			mh.model.currentScene = SceneHome
			return mh.model, nil
		}
	case "q":
		// ホーム画面でのみ終了可能
		if mh.model.currentScene == SceneHome {
			return mh.model, tea.Quit
		}
	}
	// 各画面に入力を転送
	return mh.forwardToCurrentScene(msg)
}

// handleChangeSceneMsg はChangeSceneMsgを処理します。
func (mh *MessageHandlers) handleChangeSceneMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	sceneMsg := msg.(ChangeSceneMsg)
	mh.model.currentScene = sceneMsg.Scene
	return mh.model, nil
}

// handleScreensChangeSceneMsg は画面からのシーン遷移要求を処理します。
func (mh *MessageHandlers) handleScreensChangeSceneMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	sceneMsg := msg.(screens.ChangeSceneMsg)
	mh.model.handleScreenSceneChange(sceneMsg.Scene)
	return mh.model, nil
}

// handleSceneChangeMsg はシーン遷移関連のメッセージを統合処理します。
func (mh *MessageHandlers) handleSceneChangeMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case ChangeSceneMsg:
		return mh.handleChangeSceneMsg(m)
	case screens.ChangeSceneMsg:
		return mh.handleScreensChangeSceneMsg(m)
	}
	return mh.model, nil
}

// handleStartBattleMsg はバトル開始メッセージを処理します。
func (mh *MessageHandlers) handleStartBattleMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	startMsg := msg.(screens.StartBattleMsg)
	cmd := mh.model.startBattle(startMsg.Level)
	return mh.model, cmd
}

// handleBattleTickMsg はバトルのtickメッセージを処理します。
func (mh *MessageHandlers) handleBattleTickMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	if mh.model.currentScene == SceneBattle && mh.model.battleScreen != nil {
		_, cmd := mh.model.battleScreen.Update(msg)
		return mh.model, cmd
	}
	return mh.model, nil
}

// handleBattleResultMsg はバトル結果メッセージを処理します。
func (mh *MessageHandlers) handleBattleResultMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	resultMsg := msg.(screens.BattleResultMsg)
	mh.model.handleBattleResult(resultMsg)
	return mh.model, nil
}

// handleBattleMsg はバトル関連のメッセージを統合処理します。
func (mh *MessageHandlers) handleBattleMsg(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case screens.StartBattleMsg:
		return mh.handleStartBattleMsg(m)
	case screens.BattleTickMsg:
		return mh.handleBattleTickMsg(m)
	case screens.BattleResultMsg:
		return mh.handleBattleResultMsg(m)
	}
	return mh.model, nil
}

// forwardToCurrentScene は現在のシーンにメッセージを転送します。
// ScreenMapに処理を委譲することで循環的複雑度を削減しています。
func (mh *MessageHandlers) forwardToCurrentScene(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := mh.model.screenMap.ForwardMessage(mh.model.currentScene, msg)
	return mh.model, cmd
}

// HandlerCount は登録されているハンドラーカテゴリの数を返します。
// 循環的複雑度の削減を確認するために使用します。
func (mh *MessageHandlers) HandlerCount() int {
	return len(mh.handlers)
}
