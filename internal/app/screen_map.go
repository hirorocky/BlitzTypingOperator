// Package app は BlitzTypingOperator TUIゲームの画面マップを提供します。
// ScreenMapは循環的複雑度を削減するため、シーンごとの画面操作を委譲します。
package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// ScreenGetter は画面のViewメソッドを持つインターフェースです
type ScreenGetter interface {
	View() string
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
}

// ScreenMap はシーンごとの画面インスタンスを管理します。
// 循環的複雑度を削減するため、renderCurrentSceneのswitch分岐を画面マップに委譲します。
type ScreenMap struct {
	model   *RootModel
	screens map[Scene]func() ScreenGetter
}

// NewScreenMap は新しいScreenMapを作成します。
func NewScreenMap(model *RootModel) *ScreenMap {
	sm := &ScreenMap{
		model:   model,
		screens: make(map[Scene]func() ScreenGetter),
	}
	sm.registerScreens()
	return sm
}

// registerScreens は全ての画面を登録します。
func (sm *ScreenMap) registerScreens() {
	sm.screens[SceneHome] = func() ScreenGetter {
		return sm.model.homeScreen
	}
	sm.screens[SceneBattleSelect] = func() ScreenGetter {
		return sm.model.battleSelectScreen
	}
	sm.screens[SceneBattle] = func() ScreenGetter {
		return sm.model.battleScreen
	}
	sm.screens[SceneAgentManagement] = func() ScreenGetter {
		return sm.model.agentManagementScreen
	}
	sm.screens[SceneEncyclopedia] = func() ScreenGetter {
		return sm.model.encyclopediaScreen
	}
	sm.screens[SceneAchievement] = func() ScreenGetter {
		return sm.model.statsAchievementsScreen
	}
	sm.screens[SceneSettings] = func() ScreenGetter {
		return sm.model.settingsScreen
	}
	sm.screens[SceneReward] = func() ScreenGetter {
		return sm.model.rewardScreen
	}
}

// GetScreen は指定されたシーンの画面を返します。
func (sm *ScreenMap) GetScreen(scene Scene) ScreenGetter {
	if getter, ok := sm.screens[scene]; ok {
		return getter()
	}
	return nil
}

// RenderScene は指定されたシーンのビューをレンダリングします。
func (sm *ScreenMap) RenderScene(scene Scene) string {
	screen := sm.GetScreen(scene)
	if screen != nil {
		return screen.View()
	}
	return sm.renderPlaceholder(scene)
}

// ForwardMessage は指定されたシーンの画面にメッセージを転送します。
func (sm *ScreenMap) ForwardMessage(scene Scene, msg tea.Msg) tea.Cmd {
	screen := sm.GetScreen(scene)
	if screen != nil {
		_, cmd := screen.Update(msg)
		return cmd
	}
	return nil
}

// MapCount は登録されている画面の数を返します。
func (sm *ScreenMap) MapCount() int {
	return len(sm.screens)
}

// renderPlaceholder はプレースホルダー画面をレンダリングします。
func (sm *ScreenMap) renderPlaceholder(scene Scene) string {
	title := sm.model.styles.Text.Title.Render("BlitzTypingOperator")
	info := sm.model.styles.Text.Subtle.Render(scene.String() + " (準備中)")
	hint := sm.model.styles.Text.Subtle.Render("Esc: ホームに戻る  q: 終了")
	return title + "\n\n" + info + "\n\n" + hint
}
