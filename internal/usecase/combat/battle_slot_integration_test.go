// Package combat はバトルエンジンを提供します。
// このファイルはAgentSlotManagerとBattleEngineの統合テストを定義します。

package combat

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/slot"
)

// ==================== バトルシステム統合テスト（Task 8） ====================

// createTestInventories はテスト用のインベントリを作成します。
func createTestInventories() (*domain.CoreInventory, *domain.SkillInventory) {
	coreInv := domain.NewCoreInventory()
	coreInv.AddCore("all_rounder", 10)
	coreInv.AddCore("attacker", 5)

	skillInv := domain.NewSkillInventory()
	skillInv.AddSkill("slash", "chain_fire")
	skillInv.AddSkill("heal", "")
	skillInv.AddSkill("fireball", "chain_burn")

	return coreInv, skillInv
}

// createTestMasterData はテスト用のマスタデータを作成します。
func createTestMasterData() (map[string]domain.CoreType, map[string]domain.SkillType, map[string]domain.PassiveSkill) {
	coreTypes := map[string]domain.CoreType{
		"all_rounder": {
			ID:             "all_rounder",
			Name:           "オールラウンダー",
			StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
			AllowedTags:    []string{"physical", "magic", "heal"},
			PassiveSkillID: "ps_test",
		},
		"attacker": {
			ID:             "attacker",
			Name:           "アタッカー",
			StatWeights:    map[string]float64{"STR": 2.0, "INT": 0.5, "WIL": 0.5, "LUK": 1.0},
			AllowedTags:    []string{"physical"},
			PassiveSkillID: "ps_test",
		},
	}

	skillTypes := map[string]domain.SkillType{
		"slash": {
			ID:   "slash",
			Name: "スラッシュ",
			Tags: []string{"physical"},
			Effects: []domain.ModuleEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 10, StatCoef: 1.0, StatRef: "STR"},
					Probability: 1.0,
				},
			},
		},
		"heal": {
			ID:   "heal",
			Name: "ヒール",
			Tags: []string{"heal"},
			Effects: []domain.ModuleEffect{
				{
					Target:      domain.TargetSelf,
					HPFormula:   &domain.HPFormula{Base: 20, StatCoef: 1.0, StatRef: "WIL"},
					Probability: 1.0,
				},
			},
		},
		"fireball": {
			ID:   "fireball",
			Name: "ファイアボール",
			Tags: []string{"magic"},
			Effects: []domain.ModuleEffect{
				{
					Target:      domain.TargetEnemy,
					HPFormula:   &domain.HPFormula{Base: 15, StatCoef: 1.5, StatRef: "INT"},
					Probability: 1.0,
				},
			},
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{
		"ps_test": {
			ID:          "ps_test",
			Name:        "テストパッシブ",
			Description: "テスト用パッシブスキル",
		},
	}

	return coreTypes, skillTypes, passiveSkills
}

// TestBattleEngine_InitializeBattleWithSlotManager はAgentSlotManagerからのバトル初期化をテストします。
func TestBattleEngine_InitializeBattleWithSlotManager(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成し、スロットを設定
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// スロット0にコアとスキルを設定
	if err := slotManager.SetCore(0, "all_rounder", 5); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 0, "slash", "chain_fire"); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 1, "heal", ""); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}

	// BuildAgentsForBattleでエージェントを構築
	agents := slotManager.BuildAgentsForBattle()
	if len(agents) != 1 {
		t.Fatalf("エージェント数: 期待 1, 実際 %d", len(agents))
	}

	// バトルエンジンを作成
	enemyTypes := []domain.EnemyType{
		{
			ID:              "slime",
			Name:            "スライム",
			BaseHP:          50,
			BaseAttackPower: 5,
			AttackType:      "physical",
		},
	}
	engine := NewBattleEngine(enemyTypes)

	// バトル初期化
	state, err := engine.InitializeBattle(5, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	// 検証
	if state.Enemy == nil {
		t.Error("敵が生成されていない")
	}
	if len(state.EquippedAgents) != 1 {
		t.Errorf("装備エージェント数: 期待 1, 実際 %d", len(state.EquippedAgents))
	}
	if state.EquippedAgents[0].Core == nil {
		t.Error("エージェントのコアがnil")
	}
	if state.EquippedAgents[0].Core.Level != 5 {
		t.Errorf("エージェントレベル: 期待 5, 実際 %d", state.EquippedAgents[0].Core.Level)
	}
}

// TestBattleEngine_BuildAgentsForBattle_EmptySlotExclusion は空スロット除外をテストします。
func TestBattleEngine_BuildAgentsForBattle_EmptySlotExclusion(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// スロット0と2のみ設定（スロット1は空）
	if err := slotManager.SetCore(0, "all_rounder", 5); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}
	if err := slotManager.SetCore(2, "attacker", 3); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}

	// BuildAgentsForBattleでエージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 空スロットは除外されるため、エージェント数は2
	if len(agents) != 2 {
		t.Fatalf("エージェント数: 期待 2, 実際 %d", len(agents))
	}

	// 敵タイプを準備
	enemyTypes := []domain.EnemyType{
		{
			ID:              "slime",
			Name:            "スライム",
			BaseHP:          50,
			BaseAttackPower: 5,
			AttackType:      "physical",
		},
	}
	engine := NewBattleEngine(enemyTypes)

	// バトル初期化
	state, err := engine.InitializeBattle(5, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	// バトルに含まれるエージェント数を確認
	if len(state.EquippedAgents) != 2 {
		t.Errorf("装備エージェント数: 期待 2, 実際 %d", len(state.EquippedAgents))
	}
}

// TestBattleEngine_AllSlotsEmpty はすべてのスロットが空の場合のエラーをテストします。
func TestBattleEngine_AllSlotsEmpty(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成（スロットは空のまま）
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// BuildAgentsForBattleでエージェントを構築（空のリストが返る）
	agents := slotManager.BuildAgentsForBattle()
	if len(agents) != 0 {
		t.Fatalf("空スロットからエージェントが生成された: %d", len(agents))
	}

	// バトルエンジンを作成
	enemyTypes := []domain.EnemyType{
		{
			ID:              "slime",
			Name:            "スライム",
			BaseHP:          50,
			BaseAttackPower: 5,
			AttackType:      "physical",
		},
	}
	engine := NewBattleEngine(enemyTypes)

	// バトル初期化はエラーになるべき
	_, err := engine.InitializeBattle(5, agents)
	if err == nil {
		t.Error("空のエージェントリストでバトル初期化が成功した（エラーになるべき）")
	}
}

// ==================== バトル中のスロット変更禁止テスト（Task 8.2） ====================

// TestAgentSlotManager_LockDuringBattle はバトル中のスロット変更禁止をテストします。
func TestAgentSlotManager_LockDuringBattle(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// スロット0にコアを設定
	if err := slotManager.SetCore(0, "all_rounder", 5); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}

	// バトル用にロック
	slotManager.Lock()

	// ロック中はスロット変更ができないことを確認
	err := slotManager.SetCore(1, "attacker", 3)
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のコア設定: 期待 ErrSlotLocked, 実際 %v", err)
	}

	err = slotManager.SetSkill(0, 0, "slash", "")
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のスキル設定: 期待 ErrSlotLocked, 実際 %v", err)
	}

	err = slotManager.ClearCore(0)
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のコアクリア: 期待 ErrSlotLocked, 実際 %v", err)
	}

	err = slotManager.ClearSkill(0, 0)
	if err != slot.ErrSlotLocked {
		t.Errorf("ロック中のスキルクリア: 期待 ErrSlotLocked, 実際 %v", err)
	}

	// アンロック
	slotManager.Unlock()

	// アンロック後は変更可能
	err = slotManager.SetCore(1, "attacker", 3)
	if err != nil {
		t.Errorf("アンロック後のコア設定に失敗: %v", err)
	}
}

// TestAgentSlotManager_IsLocked はロック状態の取得をテストします。
func TestAgentSlotManager_IsLocked(t *testing.T) {
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// 初期状態ではアンロック
	if slotManager.IsLocked() {
		t.Error("初期状態でロックされている")
	}

	// ロック
	slotManager.Lock()
	if !slotManager.IsLocked() {
		t.Error("Lock後にIsLockedがfalse")
	}

	// アンロック
	slotManager.Unlock()
	if slotManager.IsLocked() {
		t.Error("Unlock後にIsLockedがtrue")
	}
}

// ==================== バトル統合テスト（Task 8.3） ====================

// TestBattleIntegration_AgentModelConstruction はスロット構成からのAgentModel構築をテストします。
func TestBattleIntegration_AgentModelConstruction(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// 3スロット全てにエージェントを設定
	if err := slotManager.SetCore(0, "all_rounder", 10); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 0, "slash", "chain_fire"); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 1, "heal", ""); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 2, "fireball", "chain_burn"); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}

	if err := slotManager.SetCore(1, "attacker", 5); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(1, 0, "slash", ""); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}

	// スロット2は空のまま

	// エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 検証: 設定されたスロットのみ
	if len(agents) != 2 {
		t.Fatalf("エージェント数: 期待 2, 実際 %d", len(agents))
	}

	// 最初のエージェントを検証
	agent0 := agents[0]
	if agent0.Core.Type.ID != "all_rounder" {
		t.Errorf("エージェント0のコアTypeID: 期待 all_rounder, 実際 %s", agent0.Core.Type.ID)
	}
	if agent0.Core.Level != 10 {
		t.Errorf("エージェント0のレベル: 期待 10, 実際 %d", agent0.Core.Level)
	}
	if len(agent0.Modules) != 3 {
		t.Errorf("エージェント0のスキル数: 期待 3, 実際 %d", len(agent0.Modules))
	}

	// 2番目のエージェントを検証
	agent1 := agents[1]
	if agent1.Core.Type.ID != "attacker" {
		t.Errorf("エージェント1のコアTypeID: 期待 attacker, 実際 %s", agent1.Core.Type.ID)
	}
	if agent1.Core.Level != 5 {
		t.Errorf("エージェント1のレベル: 期待 5, 実際 %d", agent1.Core.Level)
	}
	if len(agent1.Modules) != 1 {
		t.Errorf("エージェント1のスキル数: 期待 1, 実際 %d", len(agent1.Modules))
	}
}

// TestBattleIntegration_BattleWithSlotAgents はスロットから構築されたエージェントでバトルが動作することをテストします。
func TestBattleIntegration_BattleWithSlotAgents(t *testing.T) {
	// テスト用データを準備
	coreInv, skillInv := createTestInventories()
	coreTypes, skillTypes, passiveSkills := createTestMasterData()

	// AgentSlotManagerを作成
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)

	// スロット0にエージェントを設定
	if err := slotManager.SetCore(0, "all_rounder", 10); err != nil {
		t.Fatalf("コア設定に失敗: %v", err)
	}
	if err := slotManager.SetSkill(0, 0, "slash", ""); err != nil {
		t.Fatalf("スキル設定に失敗: %v", err)
	}

	// バトル用にロック
	slotManager.Lock()
	defer slotManager.Unlock()

	// エージェントを構築
	agents := slotManager.BuildAgentsForBattle()

	// 敵タイプを準備
	enemyTypes := []domain.EnemyType{
		{
			ID:              "slime",
			Name:            "スライム",
			BaseHP:          100,
			BaseAttackPower: 5,
			AttackType:      "physical",
			ResolvedNormalActions: []domain.EnemyAction{
				{
					Name:           "通常攻撃",
					ActionType:     domain.EnemyActionAttack,
					AttackType:     "physical",
					DamageBase:     10,
					DamagePerLevel: 1,
					ChargeTime:     2 * time.Second,
				},
			},
		},
	}
	engine := NewBattleEngine(enemyTypes)
	engine.SetPassiveSkills(passiveSkills)

	// バトル初期化
	state, err := engine.InitializeBattle(5, agents)
	if err != nil {
		t.Fatalf("バトル初期化に失敗: %v", err)
	}

	// バトル状態の検証
	if state.Player == nil || state.Enemy == nil {
		t.Fatal("プレイヤーまたは敵がnil")
	}

	// プレイヤーHPが設定されていることを確認
	if state.Player.HP <= 0 {
		t.Error("プレイヤーHPが0以下")
	}

	// 敵HPが設定されていることを確認
	if state.Enemy.HP <= 0 {
		t.Error("敵HPが0以下")
	}

	// バトル終了判定（初期状態では終了していない）
	ended, _ := engine.CheckBattleEnd(state)
	if ended {
		t.Error("バトル初期状態で終了判定されている")
	}
}
