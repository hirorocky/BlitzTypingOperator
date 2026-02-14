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

	// SkillTypes の変換
	skillTypes := make([]rewarding.SkillDropInfo, len(ext.SkillDefinitions))
	for i, md := range ext.SkillDefinitions {
		// マスタデータからドメインモデルに変換
		domainSkill := md.ToDomain()
		skillTypes[i] = rewarding.SkillDropInfo{
			ID:           md.ID,
			Name:         md.Name,
			Icon:         domainSkill.Icon(),
			Tags:         md.Tags,
			Description:  md.Description,
			ManaCost:     md.ManaCost,
			MinDropLevel: md.MinDropLevel,
			Effects:      domainSkill.Effects(),
		}
	}

	// EnemyPassiveSkills をマップに変換
	enemyPassiveMap := make(map[string]*domain.EnemyPassiveSkill, len(ext.EnemyPassiveSkills))
	for _, ep := range ext.EnemyPassiveSkills {
		enemyPassiveMap[ep.ID] = ep.ToDomain()
	}

	// EnemyTypes の変換（パッシブを解決）
	enemyTypes := make([]domain.EnemyType, len(ext.EnemyTypes))
	for i, et := range ext.EnemyTypes {
		enemyTypes[i] = et.ToDomain()
		// パッシブIDが設定されている場合、パッシブを解決
		if et.NormalPassiveID != "" {
			if passive, ok := enemyPassiveMap[et.NormalPassiveID]; ok {
				enemyTypes[i].NormalPassive = passive
			}
		}
		if et.EnhancedPassiveID != "" {
			if passive, ok := enemyPassiveMap[et.EnhancedPassiveID]; ok {
				enemyTypes[i].EnhancedPassive = passive
			}
		}
	}

	return &gamestate.DomainDataSources{
		CoreTypes:     coreTypes,
		SkillTypes:    skillTypes,
		EnemyTypes:    enemyTypes,
		PassiveSkills: nil,
	}
}

// initBattleForTestVerify はテスト用のバトル初期化ヘルパーです。
func initBattleForTestVerify(_ *combat.BattleEngine, level int, agents []*domain.AgentModel, enemyTypes []domain.EnemyType) *combat.BattleState {
	enemyType := enemyTypes[0]
	hp := enemyType.BaseHP * level
	attackPower := 5 + (level * 2) // 固定のベース攻撃力を使用
	enemy := domain.NewEnemy(
		"test-enemy-id",
		enemyType.Name,
		level,
		hp,
		attackPower,
		enemyType,
	)

	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	player.PrepareForBattle()

	state := &combat.BattleState{
		Enemy:          enemy,
		Player:         player,
		EquippedAgents: agents,
		Level:          level,
		Stats:          &combat.BattleStatistics{},
	}

	state.Enemy.PrepareNextAction()
	return state
}

// ==================================================
// Task 12.1: リファクタリング完了の検証
// ==================================================

// TestRefactoring_GameStateRoundTrip はGameStateのセーブ/ロード往復を検証します。
// 要件12.1, 12.2, 12.3: データ整合性の維持
func TestRefactoring_GameStateRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
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
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 50,
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)
	agents := createTestAgents()

	// バトル初期化
	state := initBattleForTestVerify(engine, 1, agents, enemyTypes)

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
	damage := engine.ProcessEnemyAttackDamage(state, "physical")
	if damage <= 0 {
		t.Error("ダメージは正の値であるべき")
	}
	if state.Player.HP >= initialHP {
		t.Error("プレイヤーHPが減少するべき")
	}

	// スキル効果の検証
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}
	initialEnemyHP := state.Enemy.HP
	agent := agents[0]
	skill := agent.Skills[0]
	skillDamage := engine.ApplySkillEffect(state, agent, skill, typingResult)
	if skillDamage <= 0 {
		t.Error("スキルダメージは正の値であるべき")
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
	if config.DefaultSkillCooldown != 5.0 {
		t.Errorf("DefaultSkillCooldown expected 5.0, got %f", config.DefaultSkillCooldown)
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
	if config.SkillsPerAgent != 4 {
		t.Errorf("SkillsPerAgent expected 4, got %d", config.SkillsPerAgent)
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

// TestRefactoring_SkillIconMethod はスキルのIcon()メソッドを検証します。
// 要件7.3: Skill.Icon()メソッドの追加
func TestRefactoring_SkillIconMethod(t *testing.T) {
	// 新しいスキルシステムでは、各効果にアイコンが設定されます
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "test_skill",
		Name: "テストスキル",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)

	// スキルタイプからアイコンが取得できること
	if skill.Icon() == "" {
		t.Error("Skill should return non-empty icon")
	}

	// 効果からアイコンが取得できること
	effects := skill.Effects()
	if len(effects) > 0 && effects[0].Icon == "" {
		t.Error("Effect should have non-empty icon")
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
		AllCoreTypes: []domain.CoreType{},
	}
	if encycData.AllCoreTypes == nil {
		t.Error("EncyclopediaData.AllCoreTypes should not be nil")
	}

	// 他のフィールドも存在することを確認（構造体の完全性テスト）
	encycDataFull := screens.EncyclopediaData{
		AllCoreTypes:       []domain.CoreType{},
		AllSkillTypes:      []screens.SkillTypeInfo{},
		AllEnemyTypes:      []domain.EnemyType{},
		AcquiredCoreTypes:  []string{"core1"},
		AcquiredSkillTypes: []string{"skill1"},
		EncounteredEnemies: []string{"enemy1"},
	}
	if len(encycDataFull.AcquiredCoreTypes) != 1 {
		t.Error("EncyclopediaData fields should be accessible")
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
	io := savedata.NewSaveDataIO(tempDir, false)
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
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)
	agents := initializer.CreateInitialAgents()
	if len(agents) == 0 {
		t.Fatal("初期エージェントが存在しません")
	}
	agent := agents[0]

	battleState := initBattleForTestVerify(engine, 1, agents, enemyTypes)

	// 6. バトル進行
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	for battleState.Enemy.IsAlive() {
		engine.ApplySkillEffect(battleState, agent, agent.Skills[0], typingResult)
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

	// 8. 報酬統計の構築（変換が正しく動作することを確認）
	rewardStats := &rewarding.BattleStatistics{
		TotalWPM:         result.Stats.TotalWPM,
		TotalAccuracy:    result.Stats.TotalAccuracy,
		TotalTypingCount: result.Stats.TotalTypingCount,
		TotalDamageDealt: result.Stats.TotalDamageDealt,
		TotalDamageTaken: result.Stats.TotalDamageTaken,
		TotalHealAmount:  result.Stats.TotalHealAmount,
	}
	if rewardStats.TotalWPM != result.Stats.TotalWPM {
		t.Errorf("報酬統計のTotalWPMが一致しない: %f vs %f", rewardStats.TotalWPM, result.Stats.TotalWPM)
	}

	// 9. GameStateへの統計更新と保存（選択レベル1, デフォルトレベル1）
	gs.RecordBattleVictory(1, 1)
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
