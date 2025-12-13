// Package app は BlitzTypingOperator TUIゲームのシーンルーティングテストを提供します。
package app

import (
	"testing"
)

// TestNewSceneRouter は新しいSceneRouterが正しく初期化されることを検証します
func TestNewSceneRouter(t *testing.T) {
	router := NewSceneRouter()
	if router == nil {
		t.Fatal("NewSceneRouter() returned nil")
	}
}

// TestSceneRouter_RouteToHome はホームシーンへのルーティングを検証します
func TestSceneRouter_RouteToHome(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("home")
	if scene != SceneHome {
		t.Errorf("Route(\"home\") should return SceneHome, got %v", scene)
	}
}

// TestSceneRouter_RouteToBattleSelect はバトル選択シーンへのルーティングを検証します
func TestSceneRouter_RouteToBattleSelect(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("battle_select")
	if scene != SceneBattleSelect {
		t.Errorf("Route(\"battle_select\") should return SceneBattleSelect, got %v", scene)
	}
}

// TestSceneRouter_RouteToBattle はバトルシーンへのルーティングを検証します
func TestSceneRouter_RouteToBattle(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("battle")
	if scene != SceneBattle {
		t.Errorf("Route(\"battle\") should return SceneBattle, got %v", scene)
	}
}

// TestSceneRouter_RouteToAgentManagement はエージェント管理シーンへのルーティングを検証します
func TestSceneRouter_RouteToAgentManagement(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("agent_management")
	if scene != SceneAgentManagement {
		t.Errorf("Route(\"agent_management\") should return SceneAgentManagement, got %v", scene)
	}
}

// TestSceneRouter_RouteToEncyclopedia は図鑑シーンへのルーティングを検証します
func TestSceneRouter_RouteToEncyclopedia(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("encyclopedia")
	if scene != SceneEncyclopedia {
		t.Errorf("Route(\"encyclopedia\") should return SceneEncyclopedia, got %v", scene)
	}
}

// TestSceneRouter_RouteToStatsAchievements は統計・実績シーンへのルーティングを検証します
func TestSceneRouter_RouteToStatsAchievements(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("stats_achievements")
	if scene != SceneAchievement {
		t.Errorf("Route(\"stats_achievements\") should return SceneAchievement, got %v", scene)
	}
}

// TestSceneRouter_RouteToSettings は設定シーンへのルーティングを検証します
func TestSceneRouter_RouteToSettings(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("settings")
	if scene != SceneSettings {
		t.Errorf("Route(\"settings\") should return SceneSettings, got %v", scene)
	}
}

// TestSceneRouter_RouteToReward は報酬シーンへのルーティングを検証します
func TestSceneRouter_RouteToReward(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("reward")
	if scene != SceneReward {
		t.Errorf("Route(\"reward\") should return SceneReward, got %v", scene)
	}
}

// TestSceneRouter_RouteUnknown は未知のシーン名でデフォルトシーンを返すことを検証します
func TestSceneRouter_RouteUnknown(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("unknown_scene")
	if scene != SceneHome {
		t.Errorf("Route(\"unknown_scene\") should return SceneHome as default, got %v", scene)
	}
}

// TestSceneRouter_RouteEmptyString は空文字列でデフォルトシーンを返すことを検証します
func TestSceneRouter_RouteEmptyString(t *testing.T) {
	router := NewSceneRouter()
	scene := router.Route("")
	if scene != SceneHome {
		t.Errorf("Route(\"\") should return SceneHome as default, got %v", scene)
	}
}

// TestSceneRouter_AllRoutes はすべてのルートが正しくマッピングされることを検証します
func TestSceneRouter_AllRoutes(t *testing.T) {
	router := NewSceneRouter()

	tests := []struct {
		name      string
		sceneName string
		expected  Scene
	}{
		{"home", "home", SceneHome},
		{"battle_select", "battle_select", SceneBattleSelect},
		{"battle", "battle", SceneBattle},
		{"agent_management", "agent_management", SceneAgentManagement},
		{"encyclopedia", "encyclopedia", SceneEncyclopedia},
		{"stats_achievements", "stats_achievements", SceneAchievement},
		{"settings", "settings", SceneSettings},
		{"reward", "reward", SceneReward},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scene := router.Route(tt.sceneName)
			if scene != tt.expected {
				t.Errorf("Route(%q) = %v, expected %v", tt.sceneName, scene, tt.expected)
			}
		})
	}
}
