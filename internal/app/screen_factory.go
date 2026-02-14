// Package app は BlitzTypingOperator TUIゲームの画面生成機能を提供します。
package app

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/presenter"
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/usecase/inventory"
	gamestate "hirorocky/type-battle/internal/usecase/session"
	"hirorocky/type-battle/internal/usecase/slot"
	"hirorocky/type-battle/internal/usecase/spawning"
)

// InventoryProvider は画面に必要なインベントリ操作を提供するインターフェースです。
type InventoryProvider interface {
	GetCores() []*domain.CoreModel
	GetSkills() []*domain.SkillModel
	GetAgents() []*domain.AgentModel
	GetEquippedAgents() []*domain.AgentModel
	AddAgent(agent *domain.AgentModel) error
	RemoveCore(id string) error
	RemoveSkill(id string) error
	EquipAgent(slot int, agent *domain.AgentModel) error
	UnequipAgent(slot int) error
}

// ScreenFactory は画面インスタンスを生成します。
// GameStateから必要なデータを取得して各画面を初期化します。
type ScreenFactory struct {
	gameState      *gamestate.GameState
	enemyGenerator *spawning.EnemyGenerator
}

// NewScreenFactory は新しいScreenFactoryを作成します。
func NewScreenFactory(gs *gamestate.GameState) *ScreenFactory {
	return &ScreenFactory{
		gameState:      gs,
		enemyGenerator: gs.EnemyGenerator(),
	}
}

// CreateHomeScreen はホーム画面を作成します。
func (f *ScreenFactory) CreateHomeScreen(maxLevelReached int, invProvider InventoryProvider) *screens.HomeScreen {
	return screens.NewHomeScreen(maxLevelReached, invProvider)
}

// CreateBattleSelectScreen はバトル選択画面を作成します（旧：入力フィールド方式）。
// 後方互換性のために残していますが、新規開発ではCreateBattleSelectScreenCarouselを使用してください。
func (f *ScreenFactory) CreateBattleSelectScreen(maxLevelReached int, invProvider InventoryProvider) *screens.BattleSelectScreen {
	return screens.NewBattleSelectScreen(maxLevelReached, invProvider)
}

// CreateBattleSelectScreenCarousel はカルーセル方式のバトル選択画面を作成します。
// 敵タイプを左右キーで選択し、レベルを上下キーで調整できます。
func (f *ScreenFactory) CreateBattleSelectScreenCarousel(
	invProvider InventoryProvider,
	defeatedProvider screens.DefeatedEnemyProvider,
) *screens.BattleSelectScreenCarousel {
	return screens.NewBattleSelectScreenCarousel(invProvider, defeatedProvider, f.enemyGenerator)
}

// CreateBattleSelectScreenRankBased はランクベースのバトル選択画面を作成します。
// 現在ランクの敵のみ表示され、ランク内全敵撃破でランクアップします。
func (f *ScreenFactory) CreateBattleSelectScreenRankBased(
	agentProvider screens.AgentProvider,
	progressProvider screens.EnemyProgressProvider,
) *screens.BattleSelectScreenRankBased {
	return screens.NewBattleSelectScreenRankBased(agentProvider, progressProvider)
}

// CreateEncyclopediaScreen は図鑑画面を作成します。
func (f *ScreenFactory) CreateEncyclopediaScreen() *screens.EncyclopediaScreen {
	encycData := presenter.CreateEncyclopediaData(f.gameState)
	return screens.NewEncyclopediaScreen(encycData)
}

// CreateStatsAchievementsScreen は統計・実績画面を作成します。
func (f *ScreenFactory) CreateStatsAchievementsScreen() *screens.StatsAchievementsScreen {
	statsData := presenter.CreateStatsData(f.gameState)
	return screens.NewStatsAchievementsScreen(statsData)
}

// CreateSettingsScreen は設定画面を作成します。
func (f *ScreenFactory) CreateSettingsScreen() *screens.SettingsScreen {
	settingsData := presenter.CreateSettingsData(f.gameState)
	return screens.NewSettingsScreen(settingsData)
}

// CreateAgentCustomizationScreen はエージェントカスタマイズ画面を作成します。
// コア・スキルの付け替えを行うための画面を生成します。
func (f *ScreenFactory) CreateAgentCustomizationScreen(
	invManager *inventory.InventoryManager,
	slotManager *slot.AgentSlotManager,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
	passiveSkills map[string]domain.PassiveSkill,
	chainEffects map[string]domain.ChainEffect,
) *screens.AgentCustomizationScreen {
	return screens.NewAgentCustomizationScreen(invManager, slotManager, coreTypes, skillTypes, passiveSkills, chainEffects)
}

// CreateInventoryScreen はインベントリ画面を作成します。
// コア・スキル・チェイン効果の一覧表示と装備状況の確認ができる画面を生成します。
func (f *ScreenFactory) CreateInventoryScreen(
	invManager *inventory.InventoryManager,
	slotManager *slot.AgentSlotManager,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
	chainEffects map[string]domain.ChainEffect,
) *screens.InventoryScreen {
	return screens.NewInventoryScreen(invManager, slotManager, coreTypes, skillTypes, chainEffects)
}
