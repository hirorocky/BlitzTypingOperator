// Package integration_test は統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/infra/startup"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/rewarding"
	"hirorocky/type-battle/internal/usecase/typing"
)

// createTestExternalData はテスト用の外部データを作成します。
func createTestExternalData() *masterdata.ExternalData {
	return &masterdata.ExternalData{
		CoreTypes: []masterdata.CoreTypeData{
			{
				ID:             "all_rounder",
				Name:           "オールラウンダー",
				AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
				StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
				PassiveSkillID: "adaptability",
				MinDropLevel:   1,
			},
		},
		ModuleDefinitions: []masterdata.ModuleDefinitionData{
			{
				ID:           "physical_strike_lv1",
				Name:         "物理打撃Lv1",
				Icon:         "⚔️",
				Tags:         []string{"physical_low"},
				Description:  "物理ダメージを与える基本攻撃",
				MinDropLevel: 1,
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "enemy",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 1.0, StatRef: "STR"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:           "fireball_lv1",
				Name:         "ファイアボールLv1",
				Icon:         "🔥",
				Tags:         []string{"magic_low"},
				Description:  "魔法ダメージを与える基本魔法",
				MinDropLevel: 1,
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "enemy",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 1.0, StatRef: "INT"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:           "heal_lv1",
				Name:         "ヒールLv1",
				Icon:         "💚",
				Tags:         []string{"heal_low"},
				Description:  "HPを回復する基本回復魔法",
				MinDropLevel: 1,
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "self",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 0.8, StatRef: "INT"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:           "attack_buff_lv1",
				Name:         "攻撃バフLv1",
				Icon:         "⬆️",
				Tags:         []string{"buff_low"},
				Description:  "一時的に攻撃力を上昇させる",
				MinDropLevel: 1,
				Effects: []masterdata.SkillEffectData{
					{
						Target: "self",
						EffectColumn: &masterdata.EffectColumnData{
							Duration:      10.0,
							TimedEffectID: "st_str_buff_lv1",
						},
						Probability: 1.0,
					},
				},
			},
		},
		EnemyTypes: []masterdata.EnemyTypeData{
			{
				ID:               "slime",
				Name:             "スライム",
				BaseHP:           50,
				DropItemCategory: "core",
				DropItemTypeID:   "all_rounder",
			},
		},
		PassiveSkills: []masterdata.PassiveSkillData{
			{
				ID:          "ps_combo_master",
				Name:        "コンボマスター",
				Description: "連続タイピングでダメージ増加",
			},
		},
		FirstAgents: []masterdata.FirstAgentData{
			{
				ID:         "agent_first_1",
				CoreTypeID: "all_rounder",
				CoreLevel:  1,
				Modules: []masterdata.FirstAgentModuleData{
					{TypeID: "physical_strike_lv1"},
				},
			},
			{
				ID:         "agent_first_2",
				CoreTypeID: "all_rounder",
				CoreLevel:  1,
				Modules: []masterdata.FirstAgentModuleData{
					{TypeID: "heal_lv1"},
				},
			},
			{
				ID:         "agent_first_3",
				CoreTypeID: "all_rounder",
				CoreLevel:  1,
				Modules: []masterdata.FirstAgentModuleData{
					{TypeID: "attack_buff_lv1"},
				},
			},
		},
	}
}

// createTestRewardCalculator はテスト用のRewardCalculatorを作成します。
func createTestRewardCalculator() *rewarding.RewardCalculator {
	coreTypes := []domain.CoreType{
		{
			ID:   "all_rounder",
			Name: "オールラウンダー",
			StatWeights: map[string]float64{
				"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0,
			},
			PassiveSkillID: "balanced_power",
			AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low"},
			MinDropLevel:   1,
		},
	}

	moduleTypes := []rewarding.ModuleDropInfo{
		{
			ID:           "physical_attack_1",
			Name:         "物理打撃Lv1",
			Icon:         "⚔️",
			Tags:         []string{"physical_low"},
			Description:  "物理ダメージを与える",
			MinDropLevel: 1,
			Effects: []domain.SkillEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
					Probability: 1.0,
					Icon:        "⚔️",
				},
			},
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{
		"balanced_power": {
			ID:          "balanced_power",
			Name:        "バランスフォース",
			Description: "全ステータスがバランス良く成長",
		},
	}

	return rewarding.NewRewardCalculator(coreTypes, moduleTypes, passiveSkills)
}

// initBattleForTestE2E はテスト用のバトル初期化ヘルパーです。
func initBattleForTestE2E(_ *combat.BattleEngine, level int, agents []*domain.AgentModel, enemyTypes []domain.EnemyType) *combat.BattleState {
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
// Task 15.4: ゲームループE2Eテスト
// ==================================================

func TestE2E_NewGameFlow(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	// セーブデータがない場合は新規ゲーム開始
	if !io.Exists() {
		saveData := initializer.InitializeNewGame()

		// 初期エージェントが装備されている（AgentSlots使用）
		equippedCount := 0
		for _, slot := range saveData.Player.AgentSlots {
			if slot.CoreTypeID != "" {
				equippedCount++
			}
		}
		if equippedCount == 0 {
			t.Error("初期エージェントが装備されているべきです")
		}

		// セーブ
		err := io.SaveGame(saveData)
		if err != nil {
			t.Fatalf("セーブに失敗: %v", err)
		}
	}

	// 再起動シミュレート：ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 状態が保持されている（AgentSlots使用）
	hasEquipped := false
	for _, slot := range loadedData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			hasEquipped = true
			break
		}
	}
	if !hasEquipped {
		t.Error("装備エージェントが復元されるべきです")
	}
}

func TestE2E_BattleVictoryFlow(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	// 新規ゲーム開始
	saveData := initializer.InitializeNewGame()

	// ホーム画面（シミュレート）- 装備エージェントを取得（ドメインオブジェクトを直接作成）
	agents := initializer.CreateInitialAgents()
	if len(agents) == 0 {
		t.Fatal("エージェントがいません")
	}

	// バトル選択画面（シミュレート）- レベル1を選択
	battleLevel := 1

	// バトル開始
	enemyTypes := []domain.EnemyType{
		{
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 50,
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)
	battleState := initBattleForTestE2E(engine, battleLevel, agents, enemyTypes)

	// バトル進行：プレイヤーが攻撃して敵を倒す
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	// 敵を倒すまで攻撃を繰り返す
	for battleState.Enemy.IsAlive() {
		agent := agents[0]
		module := agent.Modules[0] // 物理攻撃
		engine.ApplySkillEffect(battleState, agent, module, typingResult)
		engine.RecordTypingResult(battleState, typingResult)
	}

	// 勝敗判定
	ended, result := engine.CheckBattleEnd(battleState)
	if !ended {
		t.Error("バトルが終了するべきです")
	}
	if !result.IsVictory {
		t.Error("勝利であるべきです")
	}

	// 報酬計算
	rewardCalc := createTestRewardCalculator()
	// バトル統計を作成
	battleStats := &rewarding.BattleStatistics{
		TotalWPM:         result.Stats.TotalWPM,
		TotalAccuracy:    result.Stats.TotalAccuracy,
		TotalTypingCount: result.Stats.TotalTypingCount,
	}
	// 敵タイプを作成（確定ドロップ用）
	enemyType := domain.EnemyType{
		ID:               "slime",
		Name:             "スライム",
		DropItemCategory: "core",
		DropItemTypeID:   "all_rounder",
	}
	rewards := rewardCalc.CalculateGuaranteedReward(battleStats, battleLevel, enemyType)

	// 報酬画面（シミュレート）- WPM、正確性を表示
	avgWPM := result.Stats.GetAverageWPM()
	if avgWPM == 0 {
		t.Error("平均WPMが計算されるべきです")
	}

	// 報酬をユニークインベントリに追加
	for _, c := range rewards.DroppedCores {
		saveData.Inventory.UniqueCores.Cores = append(saveData.Inventory.UniqueCores.Cores, c.TypeID)
	}
	for _, m := range rewards.DroppedModules {
		saveData.Inventory.UniqueSkills.Skills = append(
			saveData.Inventory.UniqueSkills.Skills, m.TypeID)
		if m.ChainEffect != nil {
			saveData.Inventory.UniqueChainEffects.ChainEffects = append(
				saveData.Inventory.UniqueChainEffects.ChainEffects, m.ChainEffect.ID)
		}
	}

	// 統計更新
	saveData.Statistics.TotalBattles++
	saveData.Statistics.Victories++
	if battleLevel > saveData.Statistics.MaxLevelReached {
		saveData.Statistics.MaxLevelReached = battleLevel
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// 状態確認
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	if loadedData.Statistics.TotalBattles != 1 {
		t.Errorf("TotalBattles expected 1, got %d", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.MaxLevelReached != 1 {
		t.Errorf("MaxLevelReached expected 1, got %d", loadedData.Statistics.MaxLevelReached)
	}
}

func TestE2E_AgentSynthesisFlow(t *testing.T) {
	// エージェント合成フロー
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	// 追加アイテム付きで新規ゲーム開始
	saveData := initializer.CreateNewGameWithExtraItems()

	// ユニークコアとスキルがインベントリにある
	if len(saveData.Inventory.UniqueCores.Cores) == 0 {
		t.Fatal("コアがありません")
	}
	if len(saveData.Inventory.UniqueSkills.Skills) < 1 {
		t.Fatalf("スキルがありません: got %d", len(saveData.Inventory.UniqueSkills.Skills))
	}

	// テスト用にドメインオブジェクトを作成（マスタデータから初期エージェントを使用）
	firstAgents := initializer.CreateInitialAgents()
	if len(firstAgents) == 0 {
		t.Fatal("初期エージェントの作成に失敗しました")
	}
	firstAgent := firstAgents[0]
	core := firstAgent.Core
	selectedModules := firstAgent.Modules

	if len(selectedModules) < 1 {
		t.Fatalf("初期モジュールが1個以上必要です: got %d", len(selectedModules))
	}

	// エージェント合成
	newAgent := domain.NewAgent("new_agent_1", core, selectedModules)

	// 合成後の状態確認
	if len(newAgent.Modules) != len(selectedModules) {
		t.Errorf("エージェントは%d個のモジュールを持つべきです", len(selectedModules))
	}

	// AgentSlotsの最後のスロットを置換（初期状態では3スロット全て埋まっている）
	slotIdx := 2 // 最後のスロットを使用

	// スキルスロット構成を作成
	var skillSlots [4]savedata.SkillSlotSaveCfg
	var chainEffectSlots [4]savedata.ChainEffectSlotSaveCfg
	for i, m := range newAgent.Modules {
		if i >= 4 {
			break
		}
		skillSlots[i] = savedata.SkillSlotSaveCfg{
			TypeID: m.TypeID,
		}
		if m.ChainEffect != nil {
			chainEffectSlots[i] = savedata.ChainEffectSlotSaveCfg{
				TypeID: m.ChainEffect.ID,
			}
		}
	}
	saveData.Player.AgentSlots[slotIdx] = savedata.AgentSlotSave{
		CoreTypeID:   newAgent.Core.TypeID,
		Skills:       skillSlots,
		ChainEffects: chainEffectSlots,
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロードして確認
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 新しいエージェントが保存されている（AgentSlots形式）
	found := false
	for _, slot := range loadedData.Player.AgentSlots {
		if slot.CoreTypeID == newAgent.Core.TypeID {
			found = true
			break
		}
	}
	if !found {
		t.Error("合成したエージェントが保存されているべきです")
	}
}

func TestE2E_ProgressionFlow(t *testing.T) {
	// ゲーム進行フロー：複数バトル→レベル上昇
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	saveData := initializer.InitializeNewGame()
	// ドメインオブジェクトを直接作成
	agents := initializer.CreateInitialAgents()

	enemyTypes := []domain.EnemyType{
		{
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 20, // 弱めに設定
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            80,
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		AccuracyFactor: 0.95,
	}

	// 5回バトルして進行
	for level := 1; level <= 5; level++ {
		battleState := initBattleForTestE2E(engine, level, agents, enemyTypes)

		// 敵を倒す
		for battleState.Enemy.IsAlive() {
			agent := agents[0]
			module := agent.Modules[0]
			engine.ApplySkillEffect(battleState, agent, module, typingResult)
		}

		// 勝利確認
		ended, result := engine.CheckBattleEnd(battleState)
		if !ended || !result.IsVictory {
			t.Errorf("レベル%dのバトルで勝利するべきです", level)
		}

		// 統計更新
		saveData.Statistics.TotalBattles++
		saveData.Statistics.Victories++
		if level > saveData.Statistics.MaxLevelReached {
			saveData.Statistics.MaxLevelReached = level
		}
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// 状態確認
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	if loadedData.Statistics.TotalBattles != 5 {
		t.Errorf("TotalBattles expected 5, got %d", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.MaxLevelReached != 5 {
		t.Errorf("MaxLevelReached expected 5, got %d", loadedData.Statistics.MaxLevelReached)
	}
}

func TestE2E_SaveQuitRestartLoad(t *testing.T) {

	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	// ゲーム開始（セッション1）
	saveData := initializer.InitializeNewGame()
	saveData.Statistics.TotalBattles = 15
	saveData.Statistics.Victories = 12
	saveData.Statistics.MaxLevelReached = 8
	saveData.Statistics.HighestWPM = 150.5

	// セーブして終了
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// 再起動シミュレート（新しいIOインスタンス）
	io2 := savedata.NewSaveDataIO(tempDir, false)

	// ロード
	loadedData, err := io2.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 状態が完全に復元されている
	if loadedData.Statistics.TotalBattles != 15 {
		t.Errorf("TotalBattles expected 15, got %d", loadedData.Statistics.TotalBattles)
	}
	if loadedData.Statistics.Victories != 12 {
		t.Errorf("Victories expected 12, got %d", loadedData.Statistics.Victories)
	}
	if loadedData.Statistics.MaxLevelReached != 8 {
		t.Errorf("MaxLevelReached expected 8, got %d", loadedData.Statistics.MaxLevelReached)
	}
	if loadedData.Statistics.HighestWPM != 150.5 {
		t.Errorf("HighestWPM expected 150.5, got %f", loadedData.Statistics.HighestWPM)
	}

	// エージェントスロットも復元されている
	hasAgent := false
	for _, slot := range loadedData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			hasAgent = true
			break
		}
	}
	if !hasAgent {
		t.Error("エージェントが復元されるべきです")
	}
}

func TestE2E_DefeatAndRetry(t *testing.T) {
	// 敗北→リトライフロー
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)
	initializer := startup.NewNewGameInitializer(createTestExternalData())

	saveData := initializer.InitializeNewGame()
	// ドメインオブジェクトを直接作成
	agents := initializer.CreateInitialAgents()

	enemyTypes := []domain.EnemyType{
		{
			ID:     "dragon",
			Name:   "ドラゴン",
			BaseHP: 1000, // 強い敵
		},
	}
	engine := combat.NewBattleEngine(enemyTypes)

	// 強い敵とバトル
	battleState := initBattleForTestE2E(engine, 10, agents, enemyTypes)

	// 敵の攻撃を受け続けて敗北
	for battleState.Player.IsAlive() {
		engine.ProcessEnemyAttackDamage(battleState, "physical")
	}

	// 敗北確認
	ended, result := engine.CheckBattleEnd(battleState)
	if !ended {
		t.Error("バトルが終了するべきです")
	}
	if result.IsVictory {
		t.Error("敗北であるべきです")
	}

	// 敗北時は報酬なし、統計は敗北カウント
	saveData.Statistics.TotalBattles++
	saveData.Statistics.Defeats++

	// セーブ（MaxLevelReachedは更新されない）
	io.SaveGame(saveData)

	// ロードして確認
	loadedData, _ := io.LoadGame()
	if loadedData.Statistics.Defeats != 1 {
		t.Errorf("Defeats expected 1, got %d", loadedData.Statistics.Defeats)
	}
	if loadedData.Statistics.MaxLevelReached != 0 {
		t.Error("敗北後のMaxLevelReachedは0のままであるべきです")
	}
}
