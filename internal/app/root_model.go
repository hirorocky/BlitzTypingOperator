// Package app は BlitzTypingOperator TUIゲームのRootModelを提供します。
// RootModelはゲーム全体の状態管理とシーンルーティングを担当します。
package app

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/infra/startup"
	"hirorocky/type-battle/internal/infra/terminal"
	"hirorocky/type-battle/internal/tui/presenter"
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/rewarding"
	gamestate "hirorocky/type-battle/internal/usecase/session"
	"hirorocky/type-battle/internal/usecase/slot"
	"hirorocky/type-battle/internal/usecase/typing"

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
	homeScreen               *screens.HomeScreen
	battleSelectScreen       *screens.BattleSelectScreenRankBased
	battleScreen             *screens.BattleScreen
	agentCustomizationScreen *screens.AgentCustomizationScreen
	inventoryScreen          *screens.InventoryScreen
	encyclopediaScreen       *screens.EncyclopediaScreen
	statsAchievementsScreen  *screens.StatsAchievementsScreen
	settingsScreen           *screens.SettingsScreen
	rewardScreen             *screens.RewardScreen

	// パッシブスキル定義（バトル開始時に BattleEngine へ渡す）
	passiveSkills map[string]domain.PassiveSkill

	// タイピング辞書（words.jsonからロード）
	typingDictionary *typing.Dictionary

	// デバッグモードフラグ
	debugMode bool

	// インベントリプロバイダー（装備エージェント取得用）
	invProvider InventoryProvider

	// 外部データ（デバッグモードで使用）
	externalData *masterdata.ExternalData

	// チェイン効果データ（デバッグモードで使用）
	chainEffects []masterdata.ChainEffectData

	// slotManager はエージェントスロット管理を担当します
	slotManager *slot.AgentSlotManager

	// invManager はユニークインベントリ管理を担当します
	invManager *inventory.InventoryManager

	// coreTypes はコアTypeマスタデータです
	coreTypes map[string]domain.CoreType

	// skillTypes はスキルTypeマスタデータです
	skillTypes map[string]domain.SkillType

	// enemyTypes は敵Typeマスタデータです（報酬計算のランクチェックに使用）
	enemyTypes map[string]domain.EnemyType
}

// NewRootModel はデフォルトの初期状態で新しいRootModelを作成します。
// 初期シーンはSceneHome（ホーム画面）に設定されます。
// セーブデータが存在する場合は自動的にロードします。
// 外部データファイル（data/）から敵タイプ等を読み込みます。
//
// dataDir: 外部データディレクトリのパス（空の場合は埋め込みデータを使用）
// embeddedFS: 埋め込みファイルシステム（dataDir が空の場合に使用）
// debugMode: デバッグモードを有効化（全コア・モジュール・チェイン効果を選択可能）
func NewRootModel(dataDir string, embeddedFS fs.FS, debugMode bool, saveFilePath string) *RootModel {
	// セーブディレクトリを決定（デバッグモードでは専用のセーブファイルを使用）
	homeDir, _ := os.UserHomeDir()
	saveDir := filepath.Join(homeDir, ".BlitzTypingOperator")
	saveDataIO := savedata.NewSaveDataIO(saveDir, debugMode)

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

	// チェイン効果データをロード
	chainEffects, _ := dataLoader.LoadChainEffects()

	// masterdata → domain型への変換（app層で変換を行う）
	var domainSources *gamestate.DomainDataSources
	var passiveSkills map[string]domain.PassiveSkill
	var typingDict *typing.Dictionary
	if loadErr == nil && externalData != nil {
		enemyTypes, coreTypes, moduleTypes := ConvertExternalDataToDomain(externalData)
		// 敵行動パターンを解決（IDからEnemyActionオブジェクトへ）
		enemyActions := ConvertEnemyActions(externalData.EnemyActions)
		ResolveEnemyTypeActions(enemyTypes, enemyActions)
		// 時限効果の解決（TimedEffectIDからColumn/Valueへ）
		timedEffects := ConvertTimedEffects(externalData.TimedEffects)
		ResolveModuleTimedEffects(moduleTypes, timedEffects)
		ResolveEnemyActionTimedEffects(enemyTypes, timedEffects)
		passiveSkills = ConvertPassiveSkills(externalData.PassiveSkills)
		chainEffectDefs := ConvertChainEffects(chainEffects)
		domainSources = &gamestate.DomainDataSources{
			CoreTypes:              coreTypes,
			SkillTypes:             moduleTypes,
			EnemyTypes:             enemyTypes,
			PassiveSkills:          passiveSkills,
			ChainEffectDefinitions: chainEffectDefs,
		}
		// タイピング辞書を変換
		if externalData.TypingDictionary != nil {
			typingDict = &typing.Dictionary{
				Easy:   externalData.TypingDictionary.Easy,
				Medium: externalData.TypingDictionary.Medium,
				Hard:   externalData.TypingDictionary.Hard,
			}
		}
	}

	// セーブデータをロードまたは新規作成
	var gs *gamestate.GameState
	var statusMessage string
	var loadedSaveData *savedata.SaveData

	if saveFilePath != "" {
		// 指定されたセーブファイルからロード
		var loadErr error
		loadedSaveData, loadErr = savedata.LoadSaveDataFromFile(saveFilePath)
		if loadErr == nil {
			gs = gamestate.GameStateFromSaveData(loadedSaveData, domainSources)
			statusMessage = "指定セーブファイルをロードしました"
			slog.Info("指定セーブファイルをロード", slog.String("path", saveFilePath))
		} else {
			// ロード失敗時は新規ゲームにフォールバック
			slog.Error("指定セーブファイルの読み込みに失敗", slog.String("path", saveFilePath), slog.Any("error", loadErr))
			initializer := startup.NewNewGameInitializer(externalData)
			loadedSaveData = initializer.InitializeNewGame()
			gs = gamestate.GameStateFromSaveData(loadedSaveData, domainSources)
			statusMessage = "セーブファイルの読み込みに失敗しました。新規ゲームを開始します"
		}
	} else if saveDataIO.Exists() {
		var loadErr error
		loadedSaveData, loadErr = saveDataIO.LoadGame()
		if loadErr == nil {
			gs = gamestate.GameStateFromSaveData(loadedSaveData, domainSources)
			statusMessage = "セーブデータをロードしました"
		} else {
			// セーブデータの読み込みに失敗した場合、新規ゲームを初期化
			initializer := startup.NewNewGameInitializer(externalData)
			loadedSaveData = initializer.InitializeNewGame()
			gs = gamestate.GameStateFromSaveData(loadedSaveData, domainSources)
			statusMessage = "セーブデータの読み込みに失敗しました。新規ゲームを開始します"
		}
	} else {
		// セーブデータが存在しない場合、新規ゲームを初期化（マスタデータ参照）
		initializer := startup.NewNewGameInitializer(externalData)
		loadedSaveData = initializer.InitializeNewGame()
		gs = gamestate.GameStateFromSaveData(loadedSaveData, domainSources)
		statusMessage = "新規ゲームを開始します"
	}

	// 外部データで敵生成器と報酬計算器を更新（ドメイン型を使用）
	if domainSources != nil {
		gs.UpdateEnemyGenerator(domainSources.EnemyTypes)
		gs.UpdateRewardCalculator(domainSources.CoreTypes, domainSources.SkillTypes, domainSources.PassiveSkills)

		// チェイン効果プールを設定（UpdateRewardCalculatorで新しいRewardCalculatorが作成されるため再設定が必要）
		if len(domainSources.ChainEffectDefinitions) > 0 {
			chainEffectPool := rewarding.NewChainEffectPool(domainSources.ChainEffectDefinitions)
			gs.RewardCalculator().SetChainEffectPool(chainEffectPool)
		}
	}

	// 新システムのマネージャーを作成
	invManager := inventory.NewInventoryManager()

	// セーブデータからユニークインベントリを復元
	if loadedSaveData != nil && loadedSaveData.Inventory != nil {
		// UniqueCoresを復元（TypeIDリスト形式）
		if loadedSaveData.Inventory.UniqueCores != nil && loadedSaveData.Inventory.UniqueCores.Cores != nil {
			for _, typeID := range loadedSaveData.Inventory.UniqueCores.Cores {
				invManager.AddCore(typeID)
			}
		}

		// UniqueSkillsを復元
		if loadedSaveData.Inventory.UniqueSkills != nil && loadedSaveData.Inventory.UniqueSkills.Skills != nil {
			for typeID, chainVariations := range loadedSaveData.Inventory.UniqueSkills.Skills {
				if len(chainVariations) == 0 {
					// チェイン効果がない場合でもスキルを保有状態にする
					invManager.AddSkill(typeID, "")
				} else {
					// 各チェイン効果バリエーションを復元
					for _, chainID := range chainVariations {
						invManager.AddSkill(typeID, chainID)
					}
				}
			}
		}
	}

	// デバッグモード: 全コアと全スキルを所持
	if debugMode && externalData != nil {
		// 全コアを追加
		for _, ct := range externalData.CoreTypes {
			invManager.AddCore(ct.ID)
		}
		// 全スキルを追加
		for _, mt := range externalData.ModuleDefinitions {
			invManager.AddSkill(mt.ID, "")
		}
		// 全チェイン効果バリエーションを追加
		for _, ce := range chainEffects {
			// 各スキルに対してチェイン効果を追加
			for _, mt := range externalData.ModuleDefinitions {
				invManager.AddSkill(mt.ID, ce.ID)
			}
		}
		slog.Info("デバッグモード: 全コア・スキルを追加",
			slog.Int("cores", len(externalData.CoreTypes)),
			slog.Int("skills", len(externalData.ModuleDefinitions)),
		)
	}

	// コアType、スキルType、敵Typeをマップに変換
	coreTypesMap := make(map[string]domain.CoreType)
	skillTypesMap := make(map[string]domain.SkillType)
	enemyTypesMap := make(map[string]domain.EnemyType)
	if domainSources != nil {
		for _, ct := range domainSources.CoreTypes {
			coreTypesMap[ct.ID] = ct
		}
		for _, mt := range domainSources.SkillTypes {
			// ModuleDropInfoからSkillTypeへの変換
			skillTypesMap[mt.ID] = domain.SkillType{
				ID:              mt.ID,
				Name:            mt.Name,
				Icon:            mt.Icon,
				Tags:            mt.Tags,
				Description:     mt.Description,
				CooldownSeconds: mt.CooldownSeconds,
				DifficultyRate:  mt.DifficultyRate,
				ChallengeType:   mt.ChallengeType,
				MinDropLevel:    mt.MinDropLevel,
				Effects:         mt.Effects,
			}
		}
		for _, et := range domainSources.EnemyTypes {
			enemyTypesMap[et.ID] = et
		}
	}

	// ChainEffectsマップを作成（AgentSlotManager用）
	chainEffectsMap := ConvertChainEffectsToMap(chainEffects)

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(
		invManager.Cores(),
		invManager.Skills(),
		coreTypesMap,
		skillTypesMap,
		passiveSkills,
		chainEffectsMap,
	)

	// セーブデータからエージェントスロットを復元
	if loadedSaveData != nil && loadedSaveData.Player != nil {
		for i, agentSlotSave := range loadedSaveData.Player.AgentSlots {
			if agentSlotSave.CoreTypeID == "" {
				continue
			}
			// コアを設定
			if err := slotManager.SetCore(i, agentSlotSave.CoreTypeID); err != nil {
				slog.Warn("スロット復元時にコア設定に失敗",
					slog.Int("slot", i),
					slog.String("coreTypeID", agentSlotSave.CoreTypeID),
					slog.Any("error", err),
				)
				continue
			}
			// スキルを設定
			for j, skillSlotSave := range agentSlotSave.Skills {
				if skillSlotSave.TypeID == "" {
					continue
				}
				if err := slotManager.SetSkill(i, j, skillSlotSave.TypeID, skillSlotSave.ChainEffectID); err != nil {
					slog.Warn("スロット復元時にスキル設定に失敗",
						slog.Int("slot", i),
						slog.Int("skillSlot", j),
						slog.String("skillTypeID", skillSlotSave.TypeID),
						slog.Any("error", err),
					)
				}
			}
		}
	}

	// インベントリプロバイダーを作成（デバッグモードに応じて切り替え）
	var invProvider InventoryProvider
	var debugInvProvider *presenter.DebugInventoryProvider
	if debugMode && externalData != nil {
		// デバッグモード: 全CoreType/SkillType/ChainEffectを選択可能
		passiveSkills := ConvertPassiveSkills(externalData.PassiveSkills)
		debugInvProvider = presenter.NewDebugInventoryProvider(
			externalData.CoreTypes,
			externalData.ModuleDefinitions,
			chainEffects,
			passiveSkills,
		)

		// デバッグモードでもAgentSlotManagerを使用
		debugInvProvider.SetSlotManager(slotManager)

		invProvider = debugInvProvider
	} else {
		// 通常モード: AgentSlotManagerを使用
		invProvider = presenter.NewInventoryProviderAdapter(slotManager)
	}

	// ScreenFactoryを作成
	screenFactory := NewScreenFactory(gs)

	// ホーム画面を初期化
	homeScreen := screenFactory.CreateHomeScreen(gs.GetMaxLevelReached(), invProvider)
	homeScreen.SetStatusMessage(statusMessage)
	// スロット準備状態プロバイダーを設定（バトル選択の有効/無効判定に使用）
	homeScreen.SetSlotProvider(slotManager)
	// ランク進行情報を設定（ホーム画面での進行状況表示用）
	homeScreen.SetCurrentRank(gs.EnemyProgress().CurrentRank)
	homeScreen.SetMaxHP(gs.Player().MaxHP)

	// バトル選択画面を初期化（ランクベース方式）
	progressAdapter := NewEnemyProgressAdapter(gs, enemyTypesMap, coreTypesMap, skillTypesMap)
	battleSelectScreen := screenFactory.CreateBattleSelectScreenRankBased(invProvider, progressAdapter)

	// 図鑑画面を初期化
	encyclopediaScreen := screenFactory.CreateEncyclopediaScreen()

	// 統計・実績画面を初期化
	statsAchievementsScreen := screenFactory.CreateStatsAchievementsScreen()

	// 設定画面を初期化
	settingsScreen := screenFactory.CreateSettingsScreen()

	// エージェントカスタマイズ画面を初期化
	agentCustomizationScreen := screenFactory.CreateAgentCustomizationScreen(
		invManager,
		slotManager,
		coreTypesMap,
		skillTypesMap,
		passiveSkills,
		chainEffectsMap,
	)

	// インベントリ画面を初期化
	inventoryScreen := screenFactory.CreateInventoryScreen(
		invManager,
		slotManager,
		coreTypesMap,
		skillTypesMap,
	)

	model := &RootModel{
		ready:                    false,
		currentScene:             SceneHome,
		gameState:                gs,
		styles:                   styles.NewGameStyles(),
		saveDataIO:               saveDataIO,
		statusMessage:            statusMessage,
		sceneRouter:              NewSceneRouter(),
		screenFactory:            screenFactory,
		homeScreen:               homeScreen,
		battleSelectScreen:       battleSelectScreen,
		agentCustomizationScreen: agentCustomizationScreen,
		inventoryScreen:          inventoryScreen,
		encyclopediaScreen:       encyclopediaScreen,
		statsAchievementsScreen:  statsAchievementsScreen,
		settingsScreen:           settingsScreen,
		passiveSkills:            passiveSkills,
		typingDictionary:         typingDict,
		debugMode:                debugMode,
		invProvider:              invProvider,
		externalData:             externalData,
		chainEffects:             chainEffects,
		// 新システムのマネージャー
		slotManager: slotManager,
		invManager:  invManager,
		coreTypes:   coreTypesMap,
		skillTypes:  skillTypesMap,
		enemyTypes:  enemyTypesMap,
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
		// 勝利時：統計を記録し、デフォルトレベル1で更新
		const defaultLevel = 1
		m.gameState.RecordBattleVictory(result.Level, defaultLevel)

		// 撃破済み敵情報の記録とHP成長は CalculateGuaranteedRewardWithProgress 内で処理される

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

		// 確定報酬を計算（敵タイプのドロップ設定に基づき、HP/ランク進行も反映）
		rewardResult := m.gameState.RewardCalculator().CalculateGuaranteedRewardWithProgress(
			rewardStats,
			result.Level,
			*result.EnemyType,
			m.gameState.EnemyProgress(),
			m.gameState.Player(),
			m.enemyTypes,
		)

		// 報酬をインベントリに追加（RootModelのinvManagerに直接追加）
		if m.invManager != nil {
			rewarding.AddRewardsToInventory(
				rewardResult,
				m.invManager.Cores(),
				m.invManager.Skills(),
			)
		}

		// 報酬画面を作成
		m.rewardScreen = screens.NewRewardScreen(rewardResult)

		// 報酬画面へ遷移
		m.currentScene = SceneReward
	} else {
		// 敗北時：統計を記録
		m.gameState.RecordBattleDefeat(result.Level)

		// ホーム画面の進行状況を更新してホームに戻る
		m.homeScreen.SetMaxLevelReached(m.gameState.GetMaxLevelReached())
		m.homeScreen.SetCurrentRank(m.gameState.EnemyProgress().CurrentRank)
		m.homeScreen.SetMaxHP(m.gameState.Player().MaxHP)
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
	m.appendNewSchemaToSaveData(saveData)
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

// appendNewSchemaToSaveData は現行スキーマ（UniqueCores, UniqueSkills, AgentSlots）を
// セーブデータに追加します。
func (m *RootModel) appendNewSchemaToSaveData(saveData *savedata.SaveData) {
	if m.invManager == nil || m.slotManager == nil {
		return
	}

	// UniqueCoresを追加（TypeIDリスト形式）
	if saveData.Inventory.UniqueCores == nil {
		saveData.Inventory.UniqueCores = &savedata.CoreInventorySave{
			Cores: make([]string, 0),
		}
	}
	ownedCores := m.invManager.Cores().GetOwnedCores()
	saveData.Inventory.UniqueCores.Cores = ownedCores

	// UniqueSkillsを追加（SkillTypeID → チェイン効果IDリスト）
	if saveData.Inventory.UniqueSkills == nil {
		saveData.Inventory.UniqueSkills = &savedata.SkillInventorySave{
			Skills: make(map[string][]string),
		}
	}
	ownedSkills := m.invManager.Skills().GetOwnedSkills()
	for typeID, ownership := range ownedSkills {
		saveData.Inventory.UniqueSkills.Skills[typeID] = ownership.GetChainVariations()
	}

	// AgentSlotsを追加（3スロットの構成）
	slots := m.slotManager.GetSlots()
	for i, agentSlot := range slots {
		if agentSlot == nil || agentSlot.IsEmpty() {
			saveData.Player.AgentSlots[i] = savedata.AgentSlotSave{}
			continue
		}

		slotSave := savedata.AgentSlotSave{
			CoreTypeID: agentSlot.CoreTypeID,
		}

		// スキルスロット構成を保存
		for j := range domain.MaxSkillSlotCount {
			skillConfig := agentSlot.GetSkill(j)
			if skillConfig != nil && !skillConfig.IsEmpty() {
				slotSave.Skills[j] = savedata.SkillSlotSaveCfg{
					TypeID:        skillConfig.TypeID,
					ChainEffectID: skillConfig.ChainEffectID,
				}
			}
		}

		saveData.Player.AgentSlots[i] = slotSave
	}
}

// handleSaveRequest はメニューからのセーブ要求を処理します。
func (m *RootModel) handleSaveRequest() {
	if m.saveDataIO == nil {
		m.homeScreen.SetStatusMessage("セーブ機能が無効です")
		return
	}

	saveData := m.gameState.ToSaveData()
	m.appendNewSchemaToSaveData(saveData)

	if err := m.saveDataIO.SaveGame(saveData); err != nil {
		slog.Error("セーブに失敗",
			slog.Any("error", err),
		)
		m.homeScreen.SetStatusMessage("セーブに失敗しました")
	} else {
		m.homeScreen.SetStatusMessage("セーブしました")
	}
}

// startBattle はバトルを開始します。
// enemyTypeID が空でない場合は指定された敵タイプで生成し、空の場合はランダム生成します。
func (m *RootModel) startBattle(level int, enemyTypeID string) tea.Cmd {
	// 敵を生成（タイプが指定されている場合はそのタイプで、なければランダム）
	var enemy *domain.EnemyModel
	if enemyTypeID != "" {
		enemy = m.gameState.EnemyGenerator().GenerateWithType(level, enemyTypeID)
	} else {
		enemy = m.gameState.EnemyGenerator().Generate(level)
	}

	// AgentSlotManagerからエージェントを取得
	agents := m.invProvider.GetEquippedAgents()

	// プレイヤーを準備（エージェントリストを渡す）
	m.gameState.PreparePlayerForBattle(agents)
	player := m.gameState.Player()

	// バトル画面を作成（JSONからロードした辞書を渡す）
	m.battleScreen = screens.NewBattleScreen(enemy, player, agents, m.typingDictionary)

	// パッシブスキル定義を設定（EffectTable登録に使用）
	if m.passiveSkills != nil {
		m.battleScreen.SetPassiveSkills(m.passiveSkills)
	}

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
		// ホーム画面の進行状況を更新
		m.homeScreen.SetMaxLevelReached(m.gameState.GetMaxLevelReached())
		m.homeScreen.SetCurrentRank(m.gameState.EnemyProgress().CurrentRank)
		m.homeScreen.SetMaxHP(m.gameState.Player().MaxHP)
		// バトル選択メニューの有効/無効状態を更新
		m.homeScreen.RefreshMenuState()
	case "battle_select":
		// バトル選択画面を再初期化してリセット（ランクベース方式）
		progressAdapter := NewEnemyProgressAdapter(m.gameState, m.enemyTypes, m.coreTypes, m.skillTypes)
		m.battleSelectScreen = m.screenFactory.CreateBattleSelectScreenRankBased(
			m.invProvider,
			progressAdapter,
		)
	case "encyclopedia":
		// 最新の図鑑データで画面を再初期化
		m.encyclopediaScreen = m.screenFactory.CreateEncyclopediaScreen()
	case "stats_achievements":
		// 最新の統計データで画面を再初期化
		m.statsAchievementsScreen = m.screenFactory.CreateStatsAchievementsScreen()
	case "inventory":
		// インベントリデータを最新状態に更新
		m.inventoryScreen.RefreshData()
	}
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
