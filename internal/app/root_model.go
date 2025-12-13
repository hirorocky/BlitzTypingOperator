// Package app は BlitzTypingOperator TUIゲームのRootModelを提供します。
// RootModelはゲーム全体の状態管理とシーンルーティングを担当します。
package app

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/infra/startup"
	"hirorocky/type-battle/internal/infra/terminal"
	"hirorocky/type-battle/internal/tui/presenter"
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/rewarding"
	gamestate "hirorocky/type-battle/internal/usecase/session"

	tea "github.com/charmbracelet/bubbletea"
)

// RootModel は BlitzTypingOperatorゲームのメインアプリケーション状態を表します。
// Bubbletea TUIフレームワークのtea.Modelインターフェースを実装し、
// ゲーム全体の状態管理とシーン間の遷移を統括します。
//
// Elm Architectureパターンに基づき、以下の責務を持ちます：
// - ゲーム全体の状態（GameState）の保持
// - 現在のシーン（画面）の管理
// - メッセージを現在のシーンへルーティング
// - シーン遷移メッセージの処理
// - 起動時のセーブデータロード
// - バトル終了時のオートセーブ
type RootModel struct {
	// ready はアプリケーションが初期化され、
	// ターミナルサイズが最小要件を満たしているかを示します
	ready bool

	// currentScene は現在表示中のシーンを表します
	currentScene Scene

	// gameState はゲーム全体の共有状態を保持します
	gameState *gamestate.GameState

	// terminalState は現在のターミナルサイズと検証状態を保持します
	terminalState *terminal.TerminalState

	// styles はアプリケーションのlipglossスタイルを保持します
	styles *styles.GameStyles

	// saveDataIO はセーブデータの読み書きを担当します
	saveDataIO *savedata.SaveDataIO

	// statusMessage はステータスメッセージ（セーブ/ロード結果など）です
	statusMessage string

	// sceneRouter はシーン遷移を管理します
	sceneRouter *SceneRouter

	// screenFactory は画面インスタンスを生成します
	screenFactory *ScreenFactory

	// messageHandlers はメッセージハンドリングを管理します
	messageHandlers *MessageHandlers

	// screenMap は画面のレンダリングと転送を管理します
	screenMap *ScreenMap

	// 各シーンの画面インスタンス
	homeScreen              *screens.HomeScreen
	battleSelectScreen      *screens.BattleSelectScreen
	battleScreen            *screens.BattleScreen
	agentManagementScreen   *screens.AgentManagementScreen
	encyclopediaScreen      *screens.EncyclopediaScreen
	statsAchievementsScreen *screens.StatsAchievementsScreen
	settingsScreen          *screens.SettingsScreen
	rewardScreen            *screens.RewardScreen
}

// NewRootModel はデフォルトの初期状態で新しいRootModelを作成します。
// 初期シーンはSceneHome（ホーム画面）に設定されます。
// セーブデータが存在する場合は自動的にロードします。
// 外部データファイル（data/）から敵タイプ等を読み込みます。
//
// dataDir: 外部データディレクトリのパス（空の場合は埋め込みデータを使用）
// embeddedFS: 埋め込みファイルシステム（dataDir が空の場合に使用）
func NewRootModel(dataDir string, embeddedFS fs.FS) *RootModel {
	// セーブディレクトリを決定
	homeDir, _ := os.UserHomeDir()
	saveDir := filepath.Join(homeDir, ".BlitzTypingOperator")
	saveDataIO := savedata.NewSaveDataIO(saveDir)

	// 外部データをロード
	var dataLoader *masterdata.DataLoader
	if dataDir != "" {
		// 外部ディレクトリから読み込み
		dataLoader = masterdata.NewDataLoader(dataDir)
	} else {
		// 埋め込みFSから読み込み
		dataLoader = masterdata.NewEmbeddedDataLoader(embeddedFS, "data")
	}
	externalData, loadErr := dataLoader.LoadAllExternalData()

	// masterdata → domain型への変換（app層で変換を行う）
	var domainSources *gamestate.DomainDataSources
	if loadErr == nil && externalData != nil {
		enemyTypes, coreTypes, moduleTypes := ConvertExternalDataToDomain(externalData)
		domainSources = &gamestate.DomainDataSources{
			CoreTypes:     coreTypes,
			ModuleTypes:   moduleTypes,
			EnemyTypes:    enemyTypes,
			PassiveSkills: gamestate.GetDefaultPassiveSkills(),
		}
	}

	// セーブデータをロードまたは新規作成
	var gs *gamestate.GameState
	var statusMessage string

	if saveDataIO.Exists() {
		saveData, err := saveDataIO.LoadGame()
		if err == nil {
			gs = gamestate.GameStateFromSaveData(saveData, domainSources)
			statusMessage = "セーブデータをロードしました"
		} else {
			// セーブデータの読み込みに失敗した場合、新規ゲームを初期化
			initializer := startup.NewNewGameInitializer(externalData)
			saveData := initializer.InitializeNewGame()
			gs = gamestate.GameStateFromSaveData(saveData, domainSources)
			statusMessage = "セーブデータの読み込みに失敗しました。新規ゲームを開始します"
		}
	} else {
		// セーブデータが存在しない場合、新規ゲームを初期化（マスタデータ参照）
		initializer := startup.NewNewGameInitializer(externalData)
		saveData := initializer.InitializeNewGame()
		gs = gamestate.GameStateFromSaveData(saveData, domainSources)
		statusMessage = "新規ゲームを開始します"
	}

	// 外部データで敵生成器と報酬計算器を更新（ドメイン型を使用）
	if domainSources != nil {
		gs.UpdateEnemyGenerator(domainSources.EnemyTypes)
		gs.UpdateRewardCalculator(domainSources.CoreTypes, domainSources.ModuleTypes, domainSources.PassiveSkills)
	}

	// インベントリプロバイダーアダプターを作成（複数画面で共有）
	invAdapter := presenter.NewInventoryProviderAdapter(
		gs.Inventory(),
		gs.AgentManager(),
		gs.Player(),
	)

	// ScreenFactoryを作成
	screenFactory := NewScreenFactory(gs)

	// ホーム画面を初期化
	homeScreen := screenFactory.CreateHomeScreen(gs.MaxLevelReached, invAdapter)
	homeScreen.SetStatusMessage(statusMessage)

	// バトル選択画面を初期化
	battleSelectScreen := screenFactory.CreateBattleSelectScreen(gs.MaxLevelReached, invAdapter)

	// エージェント管理画面を初期化
	agentManagementScreen := screenFactory.CreateAgentManagementScreen(invAdapter)

	// 図鑑画面を初期化
	encyclopediaScreen := screenFactory.CreateEncyclopediaScreen()

	// 統計・実績画面を初期化
	statsAchievementsScreen := screenFactory.CreateStatsAchievementsScreen()

	// 設定画面を初期化
	settingsScreen := screenFactory.CreateSettingsScreen()

	model := &RootModel{
		ready:                   false,
		currentScene:            SceneHome,
		gameState:               gs,
		styles:                  styles.NewGameStyles(),
		saveDataIO:              saveDataIO,
		statusMessage:           statusMessage,
		sceneRouter:             NewSceneRouter(),
		screenFactory:           screenFactory,
		homeScreen:              homeScreen,
		battleSelectScreen:      battleSelectScreen,
		agentManagementScreen:   agentManagementScreen,
		encyclopediaScreen:      encyclopediaScreen,
		statsAchievementsScreen: statsAchievementsScreen,
		settingsScreen:          settingsScreen,
	}

	// メッセージハンドラーと画面マップを初期化
	model.messageHandlers = NewMessageHandlers(model)
	model.screenMap = NewScreenMap(model)

	return model
}

// Init はアプリケーションを初期化し、初期コマンドを返します。
// これはBubbleteaプログラム開始時に一度だけ呼び出されます。
// 将来的にはセーブデータのロードやデータファイルの読み込みを行います。
func (m *RootModel) Init() tea.Cmd {
	// 将来的にセーブデータのロードなどを行う
	return nil
}

// Update は受信メッセージを処理し、モデルの状態を更新します。
// Elm Architectureのコアとなるメソッドで、すべての状態変更はここを通じて行われます。
// メッセージハンドラーに処理を委譲することで循環的複雑度を削減しています。
//
// 処理されるメッセージ：
// - tea.WindowSizeMsg: ターミナルサイズの変更
// - tea.KeyMsg: キーボード入力（終了操作など）
// - ChangeSceneMsg: シーン遷移要求
// - screens.ChangeSceneMsg: 各画面からのシーン遷移要求
// - screens.StartBattleMsg: バトル開始要求
// - screens.BattleTickMsg: バトルのtick更新
// - screens.BattleResultMsg: バトル結果
//
// 更新されたモデルと実行するコマンドを返します。
func (m *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// メッセージハンドラーに処理を委譲
	return m.messageHandlers.Handle(msg)
}

// handleBattleResult はバトル結果を処理します。
func (m *RootModel) handleBattleResult(result screens.BattleResultMsg) {
	stats := m.gameState.Statistics()

	// バトル統計を転送（勝敗に関わらず記録）
	if result.Stats != nil {
		// ダメージ統計を記録
		stats.RecordDamageDealt(result.Stats.TotalDamageDealt)
		stats.RecordDamageTaken(result.Stats.TotalDamageTaken)
		stats.RecordHealing(result.Stats.TotalHealAmount)

		// タイピング統計を記録（平均値を計算）
		if result.Stats.TotalTypingCount > 0 {
			avgWPM := result.Stats.TotalWPM / float64(result.Stats.TotalTypingCount)
			avgAccuracy := result.Stats.TotalAccuracy / float64(result.Stats.TotalTypingCount)
			stats.RecordTypingStats(avgWPM, avgAccuracy)
		}
	}

	// 敵図鑑を更新
	m.gameState.AddEncounteredEnemy(result.EnemyID)

	if result.Victory {
		// 勝利時：統計を記録し、最高レベルを更新
		m.gameState.RecordBattleVictory(result.Level)

		// ノーダメージ判定付きで実績チェック
		noDamage := result.Stats != nil && result.Stats.TotalDamageTaken == 0
		m.gameState.CheckBattleAchievementsWithNoDamage(noDamage)

		// バトル統計を変換
		rewardStats := &rewarding.BattleStatistics{
			TotalWPM:         result.Stats.TotalWPM,
			TotalAccuracy:    result.Stats.TotalAccuracy,
			TotalTypingCount: result.Stats.TotalTypingCount,
			TotalDamageDealt: result.Stats.TotalDamageDealt,
			TotalDamageTaken: result.Stats.TotalDamageTaken,
			TotalHealAmount:  result.Stats.TotalHealAmount,
		}

		// 報酬を計算（ドメイン型API）
		rewardResult := m.gameState.RewardCalculator().CalculateRewards(
			true,
			rewardStats,
			result.Level,
		)

		// 報酬をインベントリに追加
		m.gameState.AddRewardsToInventory(rewardResult)

		// 報酬画面を作成
		m.rewardScreen = screens.NewRewardScreen(rewardResult)

		// 報酬画面へ遷移
		m.currentScene = SceneReward
	} else {
		// 敗北時：統計を記録
		m.gameState.RecordBattleDefeat(result.Level)

		// ホーム画面の最高到達レベルを更新してホームに戻る
		m.homeScreen.SetMaxLevelReached(m.gameState.MaxLevelReached)
		m.currentScene = SceneHome
	}

	// オートセーブ（勝敗に関わらず実行）
	m.performAutoSave()

	m.battleScreen = nil
}

// performAutoSave はオートセーブを実行します。
func (m *RootModel) performAutoSave() {
	if m.saveDataIO == nil {
		return
	}

	saveData := m.gameState.ToSaveData()
	if err := m.saveDataIO.SaveGame(saveData); err != nil {
		slog.Error("オートセーブに失敗",
			slog.Any("error", err),
		)
		m.statusMessage = "オートセーブに失敗しました"
	} else {
		m.statusMessage = "オートセーブしました"
	}
	m.homeScreen.SetStatusMessage(m.statusMessage)
}

// startBattle はバトルを開始します。
func (m *RootModel) startBattle(level int) tea.Cmd {
	// 敵を生成
	enemy := m.gameState.EnemyGenerator().Generate(level)

	// GameStateからプレイヤーとエージェントを取得
	m.gameState.PreparePlayerForBattle()
	player := m.gameState.Player()
	agents := m.gameState.GetEquippedAgents()

	// バトル画面を作成
	m.battleScreen = screens.NewBattleScreen(enemy, player, agents)

	// シーンを切り替え
	m.currentScene = SceneBattle

	// バトル画面を初期化（tickコマンドを開始）
	return m.battleScreen.Init()
}

// handleScreenSceneChange は画面からのシーン遷移要求を処理します。
func (m *RootModel) handleScreenSceneChange(sceneName string) {
	// ホーム画面から別の画面に遷移する場合、ステータスメッセージをクリア
	if m.currentScene == SceneHome && sceneName != "home" {
		m.homeScreen.ClearStatusMessage()
	}

	// シーン固有の前処理を実行
	m.prepareSceneTransition(sceneName)

	// SceneRouterを使用してシーンを取得
	m.currentScene = m.sceneRouter.Route(sceneName)
}

// prepareSceneTransition はシーン遷移前の準備処理を実行します。
func (m *RootModel) prepareSceneTransition(sceneName string) {
	switch sceneName {
	case "home":
		// ホーム画面の最高到達レベルを更新
		m.homeScreen.SetMaxLevelReached(m.gameState.MaxLevelReached)
	case "battle_select":
		// バトル選択画面を再初期化してリセット
		invAdapter := m.createInventoryAdapter()
		m.battleSelectScreen = m.screenFactory.CreateBattleSelectScreen(
			m.gameState.MaxLevelReached,
			invAdapter,
		)
	case "encyclopedia":
		// 最新の図鑑データで画面を再初期化
		m.encyclopediaScreen = m.screenFactory.CreateEncyclopediaScreen()
	case "stats_achievements":
		// 最新の統計データで画面を再初期化
		m.statsAchievementsScreen = m.screenFactory.CreateStatsAchievementsScreen()
	}
}

// createInventoryAdapter はインベントリプロバイダーアダプターを作成します。
func (m *RootModel) createInventoryAdapter() *presenter.InventoryProviderAdapter {
	return presenter.NewInventoryProviderAdapter(
		m.gameState.Inventory(),
		m.gameState.AgentManager(),
		m.gameState.Player(),
	)
}

// View はアプリケーションの現在の状態を文字列としてレンダリングします。
// 現在のシーンに応じて適切な画面を描画します。
func (m *RootModel) View() string {
	// ターミナル状態がまだ設定されていない場合、ローディングメッセージを表示
	if m.terminalState == nil {
		return m.styles.Text.Subtle.Render("Loading...")
	}

	// ターミナルが小さすぎる場合、警告メッセージを表示
	if !m.terminalState.IsValid() {
		warning := m.styles.Text.Warning.Render(m.terminalState.WarningMessage())
		quitHint := m.styles.Text.Subtle.Render("Press q to quit.")
		return warning + "\n\n" + quitHint
	}

	// 現在のシーンに応じてビューを描画
	return m.renderCurrentScene()
}

// renderCurrentScene は現在のシーンに応じたビューを返します。
// ScreenMapに処理を委譲することで循環的複雑度を削減しています。
func (m *RootModel) renderCurrentScene() string {
	return m.screenMap.RenderScene(m.currentScene)
}

// GameState はゲーム全体の状態への参照を返します。
func (m *RootModel) GameState() *gamestate.GameState {
	return m.gameState
}

// CurrentScene は現在表示中のシーンを返します。
func (m *RootModel) CurrentScene() Scene {
	return m.currentScene
}

// ChangeScene は現在のシーンを指定されたシーンに変更します。
// シーン遷移時のバリデーションや前処理が必要な場合はこのメソッドで行います。
func (m *RootModel) ChangeScene(scene Scene) {
	m.currentScene = scene
}

// TerminalState は現在のターミナル状態への参照を返します。
// WindowSizeMsgを受信するまではnilが返されます。
func (m *RootModel) TerminalState() *terminal.TerminalState {
	return m.terminalState
}

// Styles はアプリケーションのスタイル設定への参照を返します。
func (m *RootModel) Styles() *styles.GameStyles {
	return m.styles
}

// IsReady はアプリケーションが使用可能な状態かどうかを返します。
// ターミナルサイズが最小要件を満たしている場合にtrueを返します。
func (m *RootModel) IsReady() bool {
	return m.ready
}
