// Package app は BlitzTypingOperator TUIゲームのシーン管理機能を提供します。
package app

// Scene はゲーム内の各画面を表す列挙型です。
// ホーム画面、バトル画面、エージェント管理画面などのシーンを定義します。
type Scene int

const (
	// SceneHome はホーム画面（メインメニュー）を表します。
	// 4つの主要機能（エージェント管理、バトル選択、図鑑、統計/実績）へのアクセスを提供します。
	SceneHome Scene = iota

	// SceneBattle はバトル画面を表します。
	// リアルタイムでの敵との戦闘、タイピングチャレンジ、モジュール効果の適用が行われます。
	SceneBattle

	// SceneBattleSelect はバトル選択画面を表します。
	// プレイヤーが挑戦するレベルを選択します。
	SceneBattleSelect

	// SceneAgentManagement はエージェント管理画面を表します。
	// コア/モジュール一覧、エージェント合成、装備管理を含みます。
	SceneAgentManagement

	// SceneEncyclopedia は図鑑画面を表します。
	// コア図鑑、モジュール図鑑、敵図鑑を表示します。
	SceneEncyclopedia

	// SceneAchievement は統計・実績画面を表します。
	// タイピング統計、バトル統計、実績一覧を表示します。
	SceneAchievement

	// SceneSettings は設定画面を表します。
	// キーバインド設定などを変更可能です。
	SceneSettings

	// SceneReward は報酬画面を表します。
	// バトル勝利後のドロップアイテムとバトル統計を表示します。
	SceneReward

	// SceneAgentCustomization はエージェントカスタマイズ画面を表します。
	// 3つのエージェントスロットに対してコアとスキルを自由に付け替えできます。
	SceneAgentCustomization

	// SceneInventory はインベントリ画面を表します。
	// 保有コアとスキルの一覧を確認できます。
	SceneInventory
)

// String はシーンの文字列表現を返します。
func (s Scene) String() string {
	switch s {
	case SceneHome:
		return "Home"
	case SceneBattle:
		return "Battle"
	case SceneBattleSelect:
		return "BattleSelect"
	case SceneAgentManagement:
		return "AgentManagement"
	case SceneEncyclopedia:
		return "Encyclopedia"
	case SceneAchievement:
		return "Achievement"
	case SceneSettings:
		return "Settings"
	case SceneReward:
		return "Reward"
	case SceneAgentCustomization:
		return "AgentCustomization"
	case SceneInventory:
		return "Inventory"
	default:
		return "Unknown"
	}
}

// ChangeSceneMsg はシーン遷移を要求するBubbleteaメッセージです。
// RootModelのUpdateメソッドで処理され、指定されたシーンへの遷移が実行されます。
type ChangeSceneMsg struct {
	Scene Scene
}
