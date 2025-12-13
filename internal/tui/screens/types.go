// Package screens はTUIゲームの画面を提供します。
package screens

import (
	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Screenインターフェース ====================

// Screen は画面の共通インターフェースです。
// tea.Modelを埋め込み、SetSizeとGetTitleメソッドを追加しています。
// 新規画面はこのインターフェースを実装することで、一貫性のある画面操作を保証します。

type Screen interface {
	tea.Model

	// SetSize は画面サイズを設定します。
	SetSize(width, height int)

	// GetTitle は画面のタイトルを返します。
	GetTitle() string
}

// BaseScreen は共通の画面機能を提供する基底構造体です。
// 各画面はこの構造体を埋め込むことで、SetSize/GetTitle/GetSizeの共通実装を継承できます。

type BaseScreen struct {
	width  int
	height int
	title  string
}

// NewBaseScreen は新しいBaseScreenを作成します。
func NewBaseScreen(title string) *BaseScreen {
	return &BaseScreen{
		title: title,
	}
}

// SetSize は画面サイズを設定します。
func (b *BaseScreen) SetSize(width, height int) {
	b.width = width
	b.height = height
}

// GetTitle は画面のタイトルを返します。
func (b *BaseScreen) GetTitle() string {
	return b.title
}

// GetSize は現在の画面サイズを返します。
func (b *BaseScreen) GetSize() (width, height int) {
	return b.width, b.height
}

// HandleWindowSizeMsg はWindowSizeMsgを処理してサイズを更新します。
// 各画面のUpdateメソッドで呼び出すことで、サイズ更新ロジックを共通化できます。
func (b *BaseScreen) HandleWindowSizeMsg(msg tea.WindowSizeMsg) {
	b.width = msg.Width
	b.height = msg.Height
}

// ==================== 共有型定義 ====================

// EncyclopediaData は図鑑データです。
type EncyclopediaData struct {
	AllCoreTypes        []domain.CoreType
	AllModuleTypes      []ModuleTypeInfo
	AllEnemyTypes       []domain.EnemyType
	AcquiredCoreTypes   []string
	AcquiredModuleTypes []string
	EncounteredEnemies  []string
}

// ModuleTypeInfo はモジュールタイプ情報です。
type ModuleTypeInfo struct {
	ID          string
	Name        string
	Category    domain.ModuleCategory
	Level       int
	Description string
}

// SettingsData は設定データです。
type SettingsData struct {
	Keybinds    map[string]string
	SoundVolume int
	Difficulty  string
}

// TypingStatsData はタイピング統計データです。
type TypingStatsData struct {
	MaxWPM               int
	AverageWPM           float64
	PerfectAccuracyCount int
	TotalCharacters      int
}

// BattleStatsData はバトル統計データです。
type BattleStatsData struct {
	TotalBattles    int
	Wins            int
	Losses          int
	MaxLevelReached int
}

// AchievementData は実績データです。
type AchievementData struct {
	ID          string
	Name        string
	Description string
	Achieved    bool
}

// StatsData は統計データです。
type StatsData struct {
	TypingStats  TypingStatsData
	BattleStats  BattleStatsData
	Achievements []AchievementData
}
