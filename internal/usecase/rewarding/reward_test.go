// Package reward はドロップ・報酬システムのテストを提供します。

package rewarding

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// TestBattleReward_Victory_ShowsRewardScreen は勝利時に報酬画面を表示することをテストします。
func TestBattleReward_Victory_ShowsRewardScreen(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.5,
		TotalAccuracy:    0.95,
		ClearTime:        2*time.Minute + 30*time.Second,
		TotalTypingCount: 15,
	}

	result := calculator.CreateRewardResult(true, stats, 10)

	if !result.IsVictory {
		t.Error("勝利時にIsVictoryがtrueであるべき")
	}
	if result.Stats == nil {
		t.Error("統計情報が設定されるべき")
	}
	if !result.ShowRewardScreen {
		t.Error("勝利時は報酬画面を表示すべき")
	}
}

// TestBattleReward_Victory_ShowsStatistics は勝利時にバトル統計を表示することをテストします。
func TestBattleReward_Victory_ShowsStatistics(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.5,
		TotalAccuracy:    0.95,
		ClearTime:        2*time.Minute + 30*time.Second,
		TotalTypingCount: 15,
	}

	result := calculator.CreateRewardResult(true, stats, 10)

	if result.Stats.TotalWPM != 80.5 {
		t.Errorf("WPMが期待値と異なる: got %f, want %f", result.Stats.TotalWPM, 80.5)
	}
	if result.Stats.TotalAccuracy != 0.95 {
		t.Errorf("正確性が期待値と異なる: got %f, want %f", result.Stats.TotalAccuracy, 0.95)
	}
	if result.Stats.ClearTime != 2*time.Minute+30*time.Second {
		t.Errorf("クリアタイムが期待値と異なる: got %v", result.Stats.ClearTime)
	}
}

// TestBattleReward_Defeat_NoRewardScreen は敗北時に報酬画面を表示しないことをテストします。
func TestBattleReward_Defeat_NoRewardScreen(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:      50.0,
		TotalAccuracy: 0.80,
		ClearTime:     3 * time.Minute,
	}

	result := calculator.CreateRewardResult(false, stats, 10)

	if result.IsVictory {
		t.Error("敗北時にIsVictoryがfalseであるべき")
	}
	if result.ShowRewardScreen {
		t.Error("敗北時は報酬画面を表示すべきでない")
	}
	if len(result.DroppedCores) > 0 || len(result.DroppedModules) > 0 {
		t.Error("敗北時はドロップがないべき")
	}
}

// TestCoreDrop_Judgment はコアドロップ判定が正しく動作することをテストします。
func TestCoreDrop_Judgment(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "test_core",
			Name:         "テストコア",
			MinDropLevel: 1,
			AllowedTags:  []string{"test"},
			StatWeights:  map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		},
	}
	calculator := NewRewardCalculator(coreTypes, nil, nil)

	// ドロップ率100%で確認
	calculator.SetCoreDropRate(1.0)

	droppedCore := calculator.RollCoreDrop(10)

	if droppedCore == nil {
		t.Error("ドロップ率100%でコアがドロップすべき")
	}
}

// TestCoreDrop_LevelInRange はコアレベルが敵レベル±範囲内であることをテストします。
func TestCoreDrop_LevelInRange(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "test_core",
			Name:         "テストコア",
			MinDropLevel: 1,
			AllowedTags:  []string{"test"},
			StatWeights:  map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		},
	}
	calculator := NewRewardCalculator(coreTypes, nil, nil)
	calculator.SetCoreDropRate(1.0)

	enemyLevel := 20
	levelRange := calculator.GetCoreLevelRange()

	// 複数回テストして範囲内であることを確認
	for i := 0; i < 50; i++ {
		droppedCore := calculator.RollCoreDrop(enemyLevel)
		if droppedCore == nil {
			continue
		}

		minLevel := enemyLevel - levelRange
		if minLevel < 1 {
			minLevel = 1
		}
		maxLevel := enemyLevel + levelRange

		if droppedCore.Level < minLevel || droppedCore.Level > maxLevel {
			t.Errorf("コアレベルが範囲外: got %d, want %d-%d", droppedCore.Level, minLevel, maxLevel)
		}
	}
}

// TestCoreDrop_MinDropLevel は特性別ドロップ最低敵レベル制限をテストします。
func TestCoreDrop_MinDropLevel(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "common_core",
			Name:         "一般コア",
			MinDropLevel: 1,
			AllowedTags:  []string{"test"},
			StatWeights:  map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		},
		{
			ID:           "rare_core",
			Name:         "レアコア",
			MinDropLevel: 10,
			AllowedTags:  []string{"test"},
			StatWeights:  map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		},
		{
			ID:           "epic_core",
			Name:         "エピックコア",
			MinDropLevel: 20,
			AllowedTags:  []string{"test"},
			StatWeights:  map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		},
	}
	calculator := NewRewardCalculator(coreTypes, nil, nil)
	calculator.SetCoreDropRate(1.0)

	// レベル5の敵からはMinDropLevel <= 5のコアのみドロップ可能
	eligibleTypes := calculator.GetEligibleCoreTypes(5)
	if len(eligibleTypes) != 1 {
		t.Errorf("レベル5では1種類のコアのみがドロップ可能: got %d", len(eligibleTypes))
	}
	if eligibleTypes[0].ID != "common_core" {
		t.Errorf("レベル5ではcommon_coreのみドロップ可能: got %s", eligibleTypes[0].ID)
	}

	// レベル15の敵からはMinDropLevel <= 15のコアがドロップ可能
	eligibleTypes = calculator.GetEligibleCoreTypes(15)
	if len(eligibleTypes) != 2 {
		t.Errorf("レベル15では2種類のコアがドロップ可能: got %d", len(eligibleTypes))
	}
}

// TestCoreDrop_InitialCoreTypes は初期コア特性のドロップ最低敵レベルをテストします。
func TestCoreDrop_InitialCoreTypes(t *testing.T) {
	// cores.jsonから読み込まれる初期設定を再現
	coreTypes := []domain.CoreType{
		{ID: "attack_balance", Name: "攻撃バランス", MinDropLevel: 1},
		{ID: "all_rounder", Name: "オールラウンダー", MinDropLevel: 1},
		{ID: "healer", Name: "ヒーラー", MinDropLevel: 3},
		{ID: "paladin", Name: "パラディン", MinDropLevel: 5},
	}
	calculator := NewRewardCalculator(coreTypes, nil, nil)

	// レベル1では2種類
	eligibleTypes := calculator.GetEligibleCoreTypes(1)
	if len(eligibleTypes) != 2 {
		t.Errorf("レベル1では2種類のコアがドロップ可能: got %d", len(eligibleTypes))
	}

	// レベル3では3種類
	eligibleTypes = calculator.GetEligibleCoreTypes(3)
	if len(eligibleTypes) != 3 {
		t.Errorf("レベル3では3種類のコアがドロップ可能: got %d", len(eligibleTypes))
	}

	// レベル5では4種類
	eligibleTypes = calculator.GetEligibleCoreTypes(5)
	if len(eligibleTypes) != 4 {
		t.Errorf("レベル5では4種類のコアがドロップ可能: got %d", len(eligibleTypes))
	}
}

// TestModuleDrop_Judgment はモジュールドロップ判定が正しく動作することをテストします。
func TestModuleDrop_Judgment(t *testing.T) {
	moduleTypes := []ModuleDropInfo{
		{
			ID:           "test_module",
			Name:         "テストモジュール",
			Category:     domain.PhysicalAttack,
			Level:        1,
			Tags:         []string{"physical_low"},
			MinDropLevel: 1,
		},
	}
	calculator := NewRewardCalculator(nil, moduleTypes, nil)

	// ドロップ率100%で確認
	calculator.SetModuleDropRate(1.0)

	droppedModules := calculator.RollModuleDrop(10, 2)

	if len(droppedModules) == 0 {
		t.Error("ドロップ率100%でモジュールがドロップすべき")
	}
}

// TestModuleDrop_MinDropLevel はカテゴリ×レベル別ドロップ最低敵レベル制限をテストします。
func TestModuleDrop_MinDropLevel(t *testing.T) {
	moduleTypes := []ModuleDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			Category:     domain.PhysicalAttack,
			Level:        1,
			Tags:         []string{"physical_low"},
			MinDropLevel: 1,
		},
		{
			ID:           "physical_lv2",
			Name:         "物理攻撃Lv2",
			Category:     domain.PhysicalAttack,
			Level:        2,
			Tags:         []string{"physical_mid"},
			MinDropLevel: 10,
		},
		{
			ID:           "physical_lv3",
			Name:         "物理攻撃Lv3",
			Category:     domain.PhysicalAttack,
			Level:        3,
			Tags:         []string{"physical_high"},
			MinDropLevel: 20,
		},
	}
	calculator := NewRewardCalculator(nil, moduleTypes, nil)

	// レベル5の敵からはMinDropLevel <= 5のモジュールのみドロップ可能
	eligibleTypes := calculator.GetEligibleModuleTypes(5)
	if len(eligibleTypes) != 1 {
		t.Errorf("レベル5では1種類のモジュールのみがドロップ可能: got %d", len(eligibleTypes))
	}

	// レベル15の敵からはMinDropLevel <= 15のモジュールがドロップ可能
	eligibleTypes = calculator.GetEligibleModuleTypes(15)
	if len(eligibleTypes) != 2 {
		t.Errorf("レベル15では2種類のモジュールがドロップ可能: got %d", len(eligibleTypes))
	}

	// レベル25の敵からは全モジュールがドロップ可能
	eligibleTypes = calculator.GetEligibleModuleTypes(25)
	if len(eligibleTypes) != 3 {
		t.Errorf("レベル25では3種類のモジュールがドロップ可能: got %d", len(eligibleTypes))
	}
}

// TestModuleDrop_HighLevelProgression は高レベルモジュールの段階的ドロップ設定をテストします。
func TestModuleDrop_HighLevelProgression(t *testing.T) {
	// modules.jsonから読み込まれる設定を再現
	moduleTypes := []ModuleDropInfo{
		{ID: "physical_lv1", Name: "物理打撃Lv1", MinDropLevel: 1},
		{ID: "physical_lv2", Name: "物理打撃Lv2", MinDropLevel: 10},
		{ID: "physical_lv3", Name: "物理打撃Lv3", MinDropLevel: 20},
		{ID: "magic_lv1", Name: "ファイアボールLv1", MinDropLevel: 1},
		{ID: "magic_lv2", Name: "ファイアボールLv2", MinDropLevel: 10},
		{ID: "magic_lv3", Name: "ファイアボールLv3", MinDropLevel: 20},
	}
	calculator := NewRewardCalculator(nil, moduleTypes, nil)

	// Lv1モジュールは敵Lv1以上でドロップ
	eligibleLv1 := calculator.GetEligibleModuleTypes(1)
	for _, m := range eligibleLv1 {
		if m.MinDropLevel > 1 {
			t.Errorf("レベル1で不正なモジュールがドロップ可能: %s", m.ID)
		}
	}

	// 高レベルほど高レベルの敵からのみドロップ
	for level := 1; level <= 30; level += 5 {
		eligible := calculator.GetEligibleModuleTypes(level)
		for _, m := range eligible {
			if m.MinDropLevel > level {
				t.Errorf("レベル%dでMinDropLevel=%dのモジュールがドロップ可能になっている", level, m.MinDropLevel)
			}
		}
	}
}

// TestInventoryFull_Warning はインベントリ満杯時に警告を表示することをテストします。
func TestInventoryFull_Warning(t *testing.T) {
	coreInv := domain.NewCoreInventory(2)
	moduleInv := domain.NewModuleInventory(2)

	// インベントリを満杯にする
	core1 := domain.NewCore("core1", "コア1", 1, domain.CoreType{}, domain.PassiveSkill{})
	core2 := domain.NewCore("core2", "コア2", 1, domain.CoreType{}, domain.PassiveSkill{})
	coreInv.Add(core1)
	coreInv.Add(core2)

	calculator := NewRewardCalculator(nil, nil, nil)

	// 満杯チェック
	warning := calculator.CheckInventoryFull(coreInv, moduleInv)

	if warning.CoreInventoryFull != true {
		t.Error("コアインベントリが満杯の場合、警告が出るべき")
	}
	if warning.WarningMessage == "" {
		t.Error("警告メッセージが設定されるべき")
	}
}

// TestInventoryFull_TempStorage は一時保管機能をテストします。
func TestInventoryFull_TempStorage(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	// ドロップしたアイテムを一時保管
	droppedCore := domain.NewCore("temp_core", "一時コア", 10, domain.CoreType{}, domain.PassiveSkill{})
	droppedModule := domain.NewModule("temp_module", "一時モジュール", domain.PhysicalAttack, 1, []string{}, 10.0, "STR", "テスト")

	storage := calculator.CreateTempStorage()
	storage.AddCore(droppedCore)
	storage.AddModule(droppedModule)

	if len(storage.Cores) != 1 {
		t.Errorf("一時保管コア数が期待と異なる: got %d, want 1", len(storage.Cores))
	}
	if len(storage.Modules) != 1 {
		t.Errorf("一時保管モジュール数が期待と異なる: got %d, want 1", len(storage.Modules))
	}

	// 後日受け取り
	retrievedCores := storage.RetrieveCores()
	if len(retrievedCores) != 1 {
		t.Errorf("受け取りコア数が期待と異なる: got %d, want 1", len(retrievedCores))
	}
	if len(storage.Cores) != 0 {
		t.Error("受け取り後は一時保管が空になるべき")
	}
}

// TestInventoryFull_PromptDiscard は不要アイテム破棄促進をテストします。
func TestInventoryFull_PromptDiscard(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	coreInv := domain.NewCoreInventory(2)
	core1 := domain.NewCore("core1", "コア1", 1, domain.CoreType{}, domain.PassiveSkill{})
	core2 := domain.NewCore("core2", "コア2", 1, domain.CoreType{}, domain.PassiveSkill{})
	coreInv.Add(core1)
	coreInv.Add(core2)

	moduleInv := domain.NewModuleInventory(10)

	warning := calculator.CheckInventoryFull(coreInv, moduleInv)

	if !warning.SuggestDiscard {
		t.Error("満杯時は破棄を促すべき")
	}
}
