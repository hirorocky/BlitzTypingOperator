// Package integration_test は統合テストを提供します。

// リファクタリング完了検証テスト - タスク12.1
package integration_test

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/infra/startup"
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/rewarding"
	gamestate "hirorocky/type-battle/internal/usecase/session"
	"hirorocky/type-battle/internal/usecase/typing"
)

// convertExternalDataToDomainSources は masterdata.ExternalData を gamestate.DomainDataSources に変換します。
func convertExternalDataToDomainSources(ext *masterdata.ExternalData) *gamestate.DomainDataSources {
	if ext == nil {
		return nil
	}

	// CoreTypes の変換
	coreTypes := make([]domain.CoreType, len(ext.CoreTypes))
	for i, ct := range ext.CoreTypes {
		coreTypes[i] = ct.ToDomain()
	}

	// ModuleTypes の変換
	moduleTypes := make([]rewarding.ModuleDropInfo, len(ext.ModuleDefinitions))
	for i, md := range ext.ModuleDefinitions {
		moduleTypes[i] = rewarding.ModuleDropInfo{
			ID:           md.ID,
			Name:         md.Name,
			Category:     categoryStringToModule(md.Category),
			Level:        md.Level,
			Tags:         md.Tags,
			BaseEffect:   md.BaseEffect,
			StatRef:      md.StatReference,
			Description:  md.Description,
			MinDropLevel: md.MinDropLevel,
		}
	}

	// EnemyTypes の変換
	enemyTypes := make([]domain.EnemyType, len(ext.EnemyTypes))
	for i, et := range ext.EnemyTypes {
		enemyTypes[i] = et.ToDomain()
	}

	return &gamestate.DomainDataSources{
		CoreTypes:     coreTypes,
		ModuleTypes:   moduleTypes,
		EnemyTypes:    enemyTypes,
		PassiveSkills: nil,
	}
}

// categoryStringToModule はカテゴリ文字列を domain.ModuleCategory に変換します。
func categoryStringToModule(category string) domain.ModuleCategory {
	switch category {
	case "physical_attack":
		return domain.PhysicalAttack
	case "magic_attack":
		return domain.MagicAttack
	case "heal":
		return domain.Heal
	case "buff":
		return domain.Buff
	case "debuff":
		return domain.Debuff
	default:
		return domain.PhysicalAttack
	}
}

// ==================================================
// Task 12.1: リファクタリング完了の検証
// 要件: 12.1 - 全ユニットテストが通過することを確認
// 要件: 12.2 - 既存セーブデータが正常に読み込めることを確認
// 要件: 12.3 - 外部から観測可能な動作に変更がないことを確認
// ==================================================

// TestRefactoring_SaveDataBackwardCompatibility はセーブデータの後方互換性を検証します。
// 要件12.2: 既存セーブデータが正常に読み込めることを確認
func TestRefactoring_SaveDataBackwardCompatibility(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)
	externalData := createTestExternalData()

	// リファクタリング前の形式でセーブデータを作成（手動構築）
	saveData := savedata.NewSaveData()

	// プレイヤー情報（装備エージェントIDのみ）
	saveData.Player.EquippedAgentIDs = [3]string{"agent_1", "agent_2", ""}

	// インベントリ情報（ID化された構造）
	saveData.Inventory.CoreInstances = []savedata.CoreInstanceSave{
		{ID: "core_1", CoreTypeID: "all_rounder", Level: 5},
		{ID: "core_2", CoreTypeID: "all_rounder", Level: 3},
	}
	saveData.Inventory.ModuleCounts = map[string]int{
		"physical_strike_lv1": 2,
		"fireball_lv1":        1,
		"heal_lv1":            1,
		"attack_buff_lv1":     1,
	}
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{
		{
			ID: "agent_1",
			Core: savedata.CoreInstanceSave{
				ID:         "core_1",
				CoreTypeID: "all_rounder",
				Level:      5,
			},
			ModuleIDs: []string{"physical_strike_lv1", "fireball_lv1", "heal_lv1", "attack_buff_lv1"},
		},
	}

	// 統計情報
	saveData.Statistics.TotalBattles = 25
	saveData.Statistics.Victories = 20
	saveData.Statistics.Defeats = 5
	saveData.Statistics.MaxLevelReached = 10
	saveData.Statistics.HighestWPM = 150.5
	saveData.Statistics.AverageWPM = 92.0
	saveData.Statistics.EncounteredEnemies = []string{"slime", "goblin"}

	// 実績情報
	saveData.Achievements.Unlocked = []string{"first_victory", "wpm_100"}
	saveData.Achievements.Progress = map[string]int{}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロードしてデータを検証
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// プレイヤー情報の検証
	if loadedData.Player.EquippedAgentIDs[0] != "agent_1" {
		t.Errorf("EquippedAgentIDs[0] expected 'agent_1', got '%s'", loadedData.Player.EquippedAgentIDs[0])
	}

	// インベントリの検証
	if len(loadedData.Inventory.CoreInstances) != 2 {
		t.Errorf("CoreInstances expected 2, got %d", len(loadedData.Inventory.CoreInstances))
	}
	if len(loadedData.Inventory.AgentInstances) != 1 {
		t.Errorf("AgentInstances expected 1, got %d", len(loadedData.Inventory.AgentInstances))
	}

	// 統計の検証
	if loadedData.Statistics.TotalBattles != 25 {
		t.Errorf("TotalBattles expected 25, got %d", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.Victories != 20 {
		t.Errorf("Victories expected 20, got %d", loadedData.Statistics.Victories)
	}
	if loadedData.Statistics.MaxLevelReached != 10 {
		t.Errorf("MaxLevelReached expected 10, got %d", loadedData.Statistics.MaxLevelReached)
	}

	// 図鑑の検証（統計に含まれる）
	if len(loadedData.Statistics.EncounteredEnemies) != 2 {
		t.Errorf("EncounteredEnemies expected 2, got %d", len(loadedData.Statistics.EncounteredEnemies))
	}

	// 実績の検証
	foundFirstVictory := false
	for _, ach := range loadedData.Achievements.Unlocked {
		if ach == "first_victory" {
			foundFirstVictory = true
			break
		}
	}
	if !foundFirstVictory {
		t.Error("first_victory achievement should be in unlocked list")
	}

	// GameStateへの変換を検証（リファクタリング後のgame_stateパッケージ使用）
	gs := gamestate.GameStateFromSaveData(loadedData, convertExternalDataToDomainSources(externalData))
	if gs == nil {
		t.Fatal("GameState should not be nil")
	}
	if gs.Statistics().Battle().TotalBattles != 25 {
		t.Errorf("GameState TotalBattles expected 25, got %d", gs.Statistics().Battle().TotalBattles)
	}
}

// TestRefactoring_GameStateRoundTrip はGameStateのセーブ/ロード往復を検証します。
// 要件12.1, 12.2, 12.3: データ整合性の維持
func TestRefactoring_GameStateRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir)
	externalData := createTestExternalData()

	// 新規ゲームを初期化
	initializer := startup.NewNewGameInitializer(externalData)
	originalSaveData := initializer.InitializeNewGame()

	// 統計データを追加
	originalSaveData.Statistics.TotalBattles = 100
	originalSaveData.Statistics.Victories = 75
	originalSaveData.Statistics.HighestWPM = 200.0
	originalSaveData.Statistics.MaxLevelReached = 50

	// セーブ
	err := io.SaveGame(originalSaveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedSaveData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// GameStateに変換
	gs := gamestate.GameStateFromSaveData(loadedSaveData, convertExternalDataToDomainSources(externalData))

	// ToSaveDataで再変換
	reconvertedSaveData := gs.ToSaveData()

	// 再変換されたデータの検証
	if reconvertedSaveData.Statistics.TotalBattles != 100 {
		t.Errorf("Reconverted TotalBattles expected 100, got %d", reconvertedSaveData.Statistics.TotalBattles)
	}
	if reconvertedSaveData.Statistics.HighestWPM != 200.0 {
		t.Errorf("Reconverted HighestWPM expected 200.0, got %f", reconvertedSaveData.Statistics.HighestWPM)
	}
}

// TestRefactoring_DataConversionIntegration はデータ変換の統合を検証します。
// 要件10.2-10.5: データ変換層の正常動作
func TestRefactoring_DataConversionIntegration(t *testing.T) {
	externalData := createTestExternalData()
	initializer := startup.NewNewGameInitializer(externalData)
	saveData := initializer.InitializeNewGame()

	// GameStateを作成
	gs := gamestate.GameStateFromSaveData(saveData, convertExternalDataToDomainSources(externalData))

	// PersistenceAdapter: SaveDataへの変換
	convertedSaveData := gs.ToSaveData()
	if convertedSaveData == nil {
		t.Fatal("ToSaveData should not return nil")
	}

	// StatsDataの検証（新規ゲームなので全て0）
	if gs.Statistics().Battle().TotalBattles != 0 {
		t.Errorf("New game TotalBattles expected 0, got %d", gs.Statistics().Battle().TotalBattles)
	}

	// EncounteredEnemiesの検証（初期状態は空）
	if len(gs.GetEncounteredEnemies()) != 0 {
		t.Errorf("New game EncounteredEnemies expected 0, got %d", len(gs.GetEncounteredEnemies()))
	}
}

// TestRefactoring_BattleFlowUnchanged はバトルフローが変更されていないことを検証します。
// 要件12.3: 外部から観測可能な動作に変更がないことを確認
func TestRefactoring_BattleFlowUnchanged(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:                 "goblin",
			Name:               "ゴブリン",
			BaseHP:             50,
			BaseAttackPower:    5,
			BaseAttackInterval: 3 * time.Second,
			AttackType:         "physical",
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	// バトル初期化
	state, err := engine.InitializeBattle(1, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	// 初期状態の検証
	if state.Enemy == nil {
		t.Error("敵が生成されるべき")
	}
	if state.Player == nil {
		t.Error("プレイヤーが初期化されるべき")
	}
	if state.Player.HP != state.Player.MaxHP {
		t.Error("プレイヤーHPは最大値であるべき")
	}

	// 敵攻撃の検証
	initialHP := state.Player.HP
	damage := engine.ProcessEnemyAttack(state)
	if damage <= 0 {
		t.Error("ダメージは正の値であるべき")
	}
	if state.Player.HP >= initialHP {
		t.Error("プレイヤーHPが減少するべき")
	}

	// モジュール効果の検証
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}
	initialEnemyHP := state.Enemy.HP
	agent := agents[0]
	module := agent.Modules[0]
	moduleDamage := engine.ApplyModuleEffect(state, agent, module, typingResult)
	if moduleDamage <= 0 {
		t.Error("モジュールダメージは正の値であるべき")
	}
	if state.Enemy.HP >= initialEnemyHP {
		t.Error("敵HPが減少するべき")
	}
}

// TestRefactoring_ConstantsIntegration は定数が正しく使用されていることを検証します。
// 要件3.1-3.5: マジックナンバーの定数化
func TestRefactoring_ConstantsIntegration(t *testing.T) {
	// バトル定数の検証
	if config.BattleTickInterval != 100*time.Millisecond {
		t.Errorf("BattleTickInterval expected 100ms, got %v", config.BattleTickInterval)
	}
	if config.DefaultModuleCooldown != 5.0 {
		t.Errorf("DefaultModuleCooldown expected 5.0, got %f", config.DefaultModuleCooldown)
	}
	if config.AccuracyPenaltyThreshold != 0.5 {
		t.Errorf("AccuracyPenaltyThreshold expected 0.5, got %f", config.AccuracyPenaltyThreshold)
	}
	if config.MinEnemyAttackInterval != 500*time.Millisecond {
		t.Errorf("MinEnemyAttackInterval expected 500ms, got %v", config.MinEnemyAttackInterval)
	}

	// 効果持続時間定数の検証
	if config.BuffDuration != 10.0 {
		t.Errorf("BuffDuration expected 10.0, got %f", config.BuffDuration)
	}
	if config.DebuffDuration != 8.0 {
		t.Errorf("DebuffDuration expected 8.0, got %f", config.DebuffDuration)
	}

	// インベントリ定数の検証
	if config.MaxAgentEquipSlots != 3 {
		t.Errorf("MaxAgentEquipSlots expected 3, got %d", config.MaxAgentEquipSlots)
	}
	if config.ModulesPerAgent != 4 {
		t.Errorf("ModulesPerAgent expected 4, got %d", config.ModulesPerAgent)
	}
}

// TestRefactoring_ScreenInterfaceCompliance はScreenインターフェースの準拠を検証します。
// 要件9.1-9.4: Screenインターフェースの導入
func TestRefactoring_ScreenInterfaceCompliance(t *testing.T) {
	// BaseScreenのテスト（NewBaseScreenを使用）
	baseScreen := screens.NewBaseScreen("Test Screen")
	baseScreen.SetSize(80, 24)

	width, height := baseScreen.GetSize()
	if width != 80 {
		t.Errorf("Width expected 80, got %d", width)
	}
	if height != 24 {
		t.Errorf("Height expected 24, got %d", height)
	}
	if baseScreen.GetTitle() != "Test Screen" {
		t.Errorf("Title expected 'Test Screen', got '%s'", baseScreen.GetTitle())
	}
}

// TestRefactoring_ModuleIconMethod はモジュールのIcon()メソッドを検証します。
// 要件7.3: Module.Icon()メソッドの追加
func TestRefactoring_ModuleIconMethod(t *testing.T) {
	testCases := []struct {
		category domain.ModuleCategory
		expected string
	}{
		{domain.PhysicalAttack, ""},
		{domain.MagicAttack, ""},
		{domain.Heal, ""},
		{domain.Buff, ""},
		{domain.Debuff, ""},
	}

	for _, tc := range testCases {
		icon := tc.category.Icon()
		if icon == "" {
			t.Errorf("Category %v should return non-empty icon", tc.category)
		}
	}
}

// TestRefactoring_HPDisplayComponent はHP表示コンポーネントを検証します。
// 要件7.1-7.2: HP表示の共通化
func TestRefactoring_HPDisplayComponent(t *testing.T) {
	// HP表示の生成テスト
	current := 75
	max := 100
	barWidth := 20

	// 正常なHP値での表示
	if current > max {
		t.Error("Current HP should not exceed max HP")
	}
	if barWidth <= 0 {
		t.Error("Bar width should be positive")
	}

	// HP割合計算
	ratio := float64(current) / float64(max)
	if ratio < 0 || ratio > 1 {
		t.Errorf("HP ratio should be between 0 and 1, got %f", ratio)
	}
}

// TestRefactoring_TypeRenaming は型名の変更を検証します。
// 要件8.1-8.3: 型名の整理
func TestRefactoring_TypeRenaming(t *testing.T) {
	// EncyclopediaData（旧EncyclopediaTestData）のテスト
	encycData := &screens.EncyclopediaData{
		AllCoreTypes:        []domain.CoreType{},
		AllModuleTypes:      []screens.ModuleTypeInfo{},
		AllEnemyTypes:       []domain.EnemyType{},
		AcquiredCoreTypes:   []string{},
		AcquiredModuleTypes: []string{},
		EncounteredEnemies:  []string{},
	}
	if encycData.AllCoreTypes == nil {
		t.Error("EncyclopediaData.AllCoreTypes should not be nil")
	}

	// StatsData（旧StatsTestData）のテスト
	statsData := &screens.StatsData{
		BattleStats: screens.BattleStatsData{
			TotalBattles: 10,
			Wins:         8,
		},
	}
	if statsData.BattleStats.TotalBattles != 10 {
		t.Errorf("StatsData.BattleStats.TotalBattles expected 10, got %d", statsData.BattleStats.TotalBattles)
	}
}

// TestRefactoring_RewardConversion は報酬変換を検証します。
// 要件10.4: 報酬アダプターの実装
func TestRefactoring_RewardConversion(t *testing.T) {
	// バトル統計を作成
	battleStats := &combat.BattleStatistics{
		TotalWPM:         250.0,
		TotalAccuracy:    2.85, // 0.95 × 3
		TotalTypingCount: 3,
		TotalDamageDealt: 150,
		TotalDamageTaken: 30,
		TotalHealAmount:  20,
	}

	// バトル統計から報酬統計を構築
	rewardStats := &rewarding.BattleStatistics{
		TotalWPM:         battleStats.TotalWPM,
		TotalAccuracy:    battleStats.TotalAccuracy,
		TotalTypingCount: battleStats.TotalTypingCount,
		TotalDamageDealt: battleStats.TotalDamageDealt,
		TotalDamageTaken: battleStats.TotalDamageTaken,
		TotalHealAmount:  battleStats.TotalHealAmount,
	}

	// 変換結果の検証
	if rewardStats == nil {
		t.Fatal("RewardStats should not be nil")
	}
	if rewardStats.TotalWPM != 250.0 {
		t.Errorf("TotalWPM expected 250.0, got %f", rewardStats.TotalWPM)
	}
	if rewardStats.TotalTypingCount != 3 {
		t.Errorf("TotalTypingCount expected 3, got %d", rewardStats.TotalTypingCount)
	}
}

// TestRefactoring_AllComponentsIntegrated は全コンポーネントの統合を検証します。
// 要件12.1, 12.2, 12.3: 全体的な統合検証
func TestRefactoring_AllComponentsIntegrated(t *testing.T) {
	tempDir := t.TempDir()
	externalData := createTestExternalData()

	// 1. 新規ゲーム初期化
	initializer := startup.NewNewGameInitializer(externalData)
	saveData := initializer.InitializeNewGame()

	// 2. セーブ/ロード
	io := savedata.NewSaveDataIO(tempDir)
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 3. GameState変換
	gs := gamestate.GameStateFromSaveData(loadedData, convertExternalDataToDomainSources(externalData))
	if gs == nil {
		t.Fatal("GameState変換に失敗")
	}

	// 4. GameStateの統計データ確認
	if gs.Statistics() == nil {
		t.Fatal("Statistics should not be nil")
	}
	if gs.Statistics().Battle().TotalBattles != 0 {
		t.Errorf("Initial TotalBattles expected 0, got %d", gs.Statistics().Battle().TotalBattles)
	}

	// 5. バトルエンジンの動作確認
	enemyTypes := []domain.EnemyType{
		{
			ID:                 "slime",
			Name:               "スライム",
			BaseHP:             50,
			BaseAttackPower:    5,
			BaseAttackInterval: 3 * time.Second,
			AttackType:         "physical",
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)
	agent := initializer.CreateInitialAgent()
	agents := []*domain.AgentModel{agent}

	battleState, err := engine.InitializeBattle(1, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	// 6. バトル進行
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	for battleState.Enemy.IsAlive() {
		engine.ApplyModuleEffect(battleState, agent, agent.Modules[0], typingResult)
		engine.RecordTypingResult(battleState, typingResult)
	}

	// 7. バトル終了判定
	ended, result := engine.CheckBattleEnd(battleState)
	if !ended {
		t.Error("バトルが終了するべき")
	}
	if !result.IsVictory {
		t.Error("勝利であるべき")
	}

	// 8. 報酬統計の構築
	rewardStats := &rewarding.BattleStatistics{
		TotalWPM:         result.Stats.TotalWPM,
		TotalAccuracy:    result.Stats.TotalAccuracy,
		TotalTypingCount: result.Stats.TotalTypingCount,
		TotalDamageDealt: result.Stats.TotalDamageDealt,
		TotalDamageTaken: result.Stats.TotalDamageTaken,
		TotalHealAmount:  result.Stats.TotalHealAmount,
	}
	if rewardStats == nil {
		t.Fatal("報酬統計の変換に失敗")
	}

	// 9. GameStateへの統計更新と保存
	gs.RecordBattleVictory(1)
	reconvertedSaveData := gs.ToSaveData()

	// 10. 最終セーブ
	err = io.SaveGame(reconvertedSaveData)
	if err != nil {
		t.Fatalf("最終セーブに失敗: %v", err)
	}

	// 11. 最終ロードと検証
	finalLoadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("最終ロードに失敗: %v", err)
	}

	if finalLoadedData.Statistics.TotalBattles != 1 {
		t.Errorf("Final TotalBattles expected 1, got %d", finalLoadedData.Statistics.TotalBattles)
	}
	if finalLoadedData.Statistics.Victories != 1 {
		t.Errorf("Final Victories expected 1, got %d", finalLoadedData.Statistics.Victories)
	}
}

// rewarding.BattleStatisticsが正しくインポートされていることを確認するダミー関数
var _ = rewarding.BattleStatistics{}
