// Package app は BlitzTypingOperator TUIゲームのシーンルーティング機能を提供します。
package app

// SceneRouter はシーン遷移を管理します。
// シーン名からSceneへのマッピングを提供し、RootModelからシーン遷移ロジックを分離します。
type SceneRouter struct {
	// routes はシーン名からSceneへのマッピング
	routes map[string]Scene
}

// NewSceneRouter は新しいSceneRouterを作成します。
// 全てのシーンルートを初期化済みの状態で返します。
func NewSceneRouter() *SceneRouter {
	return &SceneRouter{
		routes: map[string]Scene{
			"home":               SceneHome,
			"battle_select":      SceneBattleSelect,
			"battle":             SceneBattle,
			"agent_management":   SceneAgentManagement,
			"encyclopedia":       SceneEncyclopedia,
			"stats_achievements": SceneAchievement,
			"settings":           SceneSettings,
			"reward":             SceneReward,
		},
	}
}

// Route はシーン名に基づいてシーンを返します。
// 未知のシーン名の場合はデフォルトとしてSceneHomeを返します。
func (r *SceneRouter) Route(sceneName string) Scene {
	if scene, ok := r.routes[sceneName]; ok {
		return scene
	}
	// 未知のシーン名の場合はデフォルトでホームを返す
	return SceneHome
}
