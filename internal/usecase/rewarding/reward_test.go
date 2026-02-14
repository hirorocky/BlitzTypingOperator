// Package reward はドロップ・報酬システムのテストを提供します。

package rewarding

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// newTestSkill はテスト用ダメージスキルを作成するヘルパー関数です。
func newTestSkill(id, name string, tags []string, statCoef float64, statRef, description string) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
}

// newTestSkillWithChainEffect はチェイン効果付きスキルを作成するヘルパー関数です。
func newTestSkillWithChainEffect(id, name string, tags []string, statCoef float64, statRef, description string, chainEffect *domain.ChainEffect) *domain.SkillModel {
	return domain.NewSkillFromType(domain.SkillType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, chainEffect)
}

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
	if len(result.DroppedCores) > 0 || len(result.DroppedSkills) > 0 {
		t.Error("敗北時はドロップがないべき")
	}
}

// TestInventoryFull_Warning は新システムでは容量制限がないことを確認します。
func TestInventoryFull_Warning(t *testing.T) {
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// 新システムでは容量制限がないため、いくつ追加しても警告は出ない
	coreInv.AddCore("core_type_1")
	coreInv.AddCore("core_type_2")
	coreInv.AddCore("core_type_3")

	calculator := NewRewardCalculator(nil, nil, nil)

	// 新システムでは常に空の警告を返す
	warning := calculator.CheckInventoryFull(coreInv, skillInv)

	if warning.WarningMessage != "" {
		t.Error("新システムでは警告メッセージは空であるべき")
	}
}

// TestInventoryFull_TempStorage は一時保管機能をテストします。
func TestInventoryFull_TempStorage(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	// ドロップしたアイテムを一時保管
	droppedCore := domain.NewCoreWithTypeID("temp_core", domain.CoreType{}, domain.PassiveSkill{})
	droppedSkill := newTestSkill("temp_skill", "一時スキル", []string{}, 10.0, "STR", "テスト")

	storage := calculator.CreateTempStorage()
	storage.AddCore(droppedCore)
	storage.AddSkill(droppedSkill)

	if len(storage.Cores) != 1 {
		t.Errorf("一時保管コア数が期待と異なる: got %d, want 1", len(storage.Cores))
	}
	if len(storage.Skills) != 1 {
		t.Errorf("一時保管スキル数が期待と異なる: got %d, want 1", len(storage.Skills))
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

// TestInventoryFull_PromptDiscard は新システムでは破棄促進が不要であることを確認します。
func TestInventoryFull_PromptDiscard(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// いくつ追加してもユニーク管理なので容量制限なし
	coreInv.AddCore("core_type_1")
	coreInv.AddCore("core_type_2")

	warning := calculator.CheckInventoryFull(coreInv, skillInv)

	// 新システムでは破棄を促す必要がない
	if warning.SuggestDiscard {
		t.Error("新システムでは破棄促進は不要")
	}
}

// ==================== チェイン効果ランダム決定テスト ====================

// TestChainEffectPool_CreateFromSkillEffects はチェイン効果プールの作成をテストします。
func TestChainEffectPool_CreateFromSkillEffects(t *testing.T) {
	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   30,
		},
		{
			ID:         "damage_cut",
			Name:       "ダメージカット",
			Category:   "defense",
			EffectType: domain.ChainEffectDamageCut,
			MinValue:   10,
			MaxValue:   30,
		},
	}

	pool := NewChainEffectPool(skillEffects)

	if pool == nil {
		t.Fatal("チェイン効果プールがnilであってはならない")
	}
	if len(pool.Effects) != 2 {
		t.Errorf("チェイン効果数が期待と異なる: got %d, want 2", len(pool.Effects))
	}
}

// TestChainEffectPool_GenerateRandomEffect はランダムなチェイン効果生成をテストします。
func TestChainEffectPool_GenerateRandomEffect(t *testing.T) {
	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   30,
		},
	}

	pool := NewChainEffectPool(skillEffects)

	// 複数回生成して値が範囲内であることを確認
	for i := 0; i < 50; i++ {
		effect := pool.GenerateRandomEffect()
		if effect == nil {
			continue // nilチェイン効果もあり得る
		}
		if effect.Value < 10 || effect.Value > 30 {
			t.Errorf("効果値が範囲外: got %.0f, want 10-30", effect.Value)
		}
		if effect.Type != domain.ChainEffectDamageAmp {
			t.Errorf("効果タイプが期待と異なる: got %s, want %s", effect.Type, domain.ChainEffectDamageAmp)
		}
	}
}

// TestChainEffectPool_GenerateWithNilProbability はチェイン効果なしの確率をテストします。
func TestChainEffectPool_GenerateWithNilProbability(t *testing.T) {
	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   30,
		},
	}

	pool := NewChainEffectPool(skillEffects)

	// nilチェイン効果確率を100%に設定
	pool.SetNoEffectProbability(1.0)

	for i := 0; i < 10; i++ {
		effect := pool.GenerateRandomEffect()
		if effect != nil {
			t.Error("nil確率100%でチェイン効果がnilであるべき")
		}
	}

	// nil確率を0%に設定
	pool.SetNoEffectProbability(0.0)

	foundNonNil := false
	for i := 0; i < 10; i++ {
		effect := pool.GenerateRandomEffect()
		if effect != nil {
			foundNonNil = true
			break
		}
	}
	if !foundNonNil {
		t.Error("nil確率0%でチェイン効果が生成されるべき")
	}
}

// TestSkillDropInfo_ToDomainWithRandomChainEffect はチェイン効果付きドメイン変換をテストします。
func TestSkillDropInfo_ToDomainWithRandomChainEffect(t *testing.T) {
	dropInfo := SkillDropInfo{
		ID:          "physical_lv1",
		Name:        "物理攻撃Lv1",
		Icon:        "⚔️",
		Tags:        []string{"physical_low"},
		Description: "テスト",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 10.0, StatCoef: 1.0, StatRef: "STR"},
				Probability: 1.0,
			},
		},
	}

	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 20)

	skill := dropInfo.ToDomainWithChainEffect(&effect)

	if skill == nil {
		t.Fatal("スキルがnilであってはならない")
	}
	if skill.ChainEffect == nil {
		t.Error("チェイン効果が設定されるべき")
	}
	if skill.ChainEffect.Type != domain.ChainEffectDamageAmp {
		t.Errorf("チェイン効果タイプが期待と異なる: got %s, want %s", skill.ChainEffect.Type, domain.ChainEffectDamageAmp)
	}
	if skill.ChainEffect.Value != 20 {
		t.Errorf("チェイン効果値が期待と異なる: got %.0f, want 20", skill.ChainEffect.Value)
	}
}

// ==================== タスク11.2: スキル入手処理更新テスト ====================

// TestAddRewardsToInventory_WithChainEffect はチェイン効果付きスキルがインベントリに追加されることをテストします。
func TestAddRewardsToInventory_WithChainEffect(t *testing.T) {
	// チェイン効果付きスキルを作成
	effect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageAmp, 25)
	skill := newTestSkillWithChainEffect(
		"physical_lv1",
		"物理攻撃Lv1",
		[]string{"physical_low"},
		10.0,
		"STR",
		"テスト",
		&effect,
	)

	// 報酬結果を作成
	result := &RewardResult{
		IsVictory:     true,
		DroppedSkills: []*domain.SkillModel{skill},
	}

	// インベントリを作成（新システム）
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// インベントリに追加
	warning := AddRewardsToInventory(result, coreInv, skillInv, nil)

	if warning.WarningMessage != "" {
		t.Error("新システムでは警告メッセージは空であるべき")
	}

	// スキルインベントリにスキルが追加されたことを確認
	if !skillInv.HasSkill("physical_lv1") {
		t.Error("スキルがインベントリに追加されるべき")
	}
}

// TestChainEffectPool_MultipleEffectTypes は複数のチェイン効果タイプからランダム選択されることをテストします。
func TestChainEffectPool_MultipleEffectTypes(t *testing.T) {
	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   30,
		},
		{
			ID:         "damage_cut",
			Name:       "ダメージカット",
			Category:   "defense",
			EffectType: domain.ChainEffectDamageCut,
			MinValue:   10,
			MaxValue:   30,
		},
		{
			ID:         "heal_amp",
			Name:       "ヒールアンプ",
			Category:   "heal",
			EffectType: domain.ChainEffectHealAmp,
			MinValue:   15,
			MaxValue:   35,
		},
	}

	pool := NewChainEffectPool(skillEffects)
	pool.SetNoEffectProbability(0.0)

	// 複数回生成して複数タイプが選択されることを確認
	typeCounts := make(map[domain.ChainEffectType]int)

	for i := 0; i < 100; i++ {
		effect := pool.GenerateRandomEffect()
		if effect != nil {
			typeCounts[effect.Type]++
		}
	}

	// 最低2種類は選択されているはず（確率的に）
	if len(typeCounts) < 2 {
		t.Errorf("複数のチェイン効果タイプが選択されるべき: got %d types", len(typeCounts))
	}
}

// TestChainEffectPool_EmptyEffects は空のチェイン効果プールでnilが返ることをテストします。
func TestChainEffectPool_EmptyEffects(t *testing.T) {
	pool := NewChainEffectPool(nil)

	effect := pool.GenerateRandomEffect()

	if effect != nil {
		t.Error("空のプールではnilが返るべき")
	}
}

// ==================== タスク5.1: 確定ドロップの基本ロジックテスト ====================

// TestCalculateGuaranteedReward_EnemyWithDropCategory は敵にドロップカテゴリ設定がある場合に確定ドロップすることをテストします。
func TestCalculateGuaranteedReward_EnemyWithDropCategory(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	calculator := NewRewardCalculator(coreTypes, skillTypes, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// コアドロップ設定の敵タイプ
	enemyType := domain.EnemyType{
		ID:               "slime",
		Name:             "スライム",
		DropItemCategory: "core",
		DropItemTypeID:   "attack_balance",
	}

	result := calculator.CalculateGuaranteedReward(stats, 10, enemyType)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if !result.IsVictory {
		t.Error("勝利フラグがtrueであるべき")
	}

	// 必ず1つのアイテムがドロップすること
	totalItems := len(result.DroppedCores) + len(result.DroppedSkills)
	if totalItems != 1 {
		t.Errorf("確定ドロップで1つのアイテムがドロップすべき: got %d", totalItems)
	}

	// コアがドロップすること
	if len(result.DroppedCores) != 1 {
		t.Errorf("コアがドロップすべき: got %d cores", len(result.DroppedCores))
	}

	// ドロップしたコアがTypeIDに対応していること
	if len(result.DroppedCores) > 0 {
		core := result.DroppedCores[0]
		if core.Type.ID != "attack_balance" {
			t.Errorf("コアTypeIDが期待と異なる: got %s, want attack_balance", core.Type.ID)
		}
	}
}

// TestCalculateGuaranteedReward_SkillDrop はスキルドロップ設定の敵からスキルがドロップすることをテストします。
func TestCalculateGuaranteedReward_SkillDrop(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
			Effects: []domain.SkillEffect{
				{Target: domain.TargetEnemy, Probability: 1.0},
			},
		},
	}

	calculator := NewRewardCalculator(nil, skillTypes, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// スキルドロップ設定の敵タイプ
	enemyType := domain.EnemyType{
		ID:               "goblin",
		Name:             "ゴブリン",
		DropItemCategory: "skill",
		DropItemTypeID:   "physical_lv1",
	}

	result := calculator.CalculateGuaranteedReward(stats, 10, enemyType)

	// 必ず1つのアイテムがドロップすること
	totalItems := len(result.DroppedCores) + len(result.DroppedSkills)
	if totalItems != 1 {
		t.Errorf("確定ドロップで1つのアイテムがドロップすべき: got %d", totalItems)
	}

	// スキルがドロップすること
	if len(result.DroppedSkills) != 1 {
		t.Errorf("スキルがドロップすべき: got %d skills", len(result.DroppedSkills))
	}

	// ドロップしたスキルがTypeIDに対応していること
	if len(result.DroppedSkills) > 0 {
		skill := result.DroppedSkills[0]
		if skill.TypeID != "physical_lv1" {
			t.Errorf("スキルTypeIDが期待と異なる: got %s, want physical_lv1", skill.TypeID)
		}
	}
}

// TestCalculateGuaranteedReward_PanicOnMissingDropConfig はドロップ設定がない場合にpanicすることをテストします。
func TestCalculateGuaranteedReward_PanicOnMissingDropConfig(t *testing.T) {
	calculator := NewRewardCalculator(nil, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// ドロップ設定がない敵タイプ
	enemyType := domain.EnemyType{
		ID:               "unknown_enemy",
		Name:             "不明な敵",
		DropItemCategory: "", // 空
		DropItemTypeID:   "", // 空
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("ドロップ設定がない場合にpanicすべき")
		}
	}()

	calculator.CalculateGuaranteedReward(stats, 10, enemyType)
}

// TestCalculateGuaranteedReward_PanicOnInvalidTypeID は不正なTypeIDの場合にpanicすることをテストします。
func TestCalculateGuaranteedReward_PanicOnInvalidTypeID(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}

	calculator := NewRewardCalculator(coreTypes, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// 存在しないTypeID
	enemyType := domain.EnemyType{
		ID:               "unknown_enemy",
		Name:             "不明な敵",
		DropItemCategory: "core",
		DropItemTypeID:   "non_existent_core",
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("存在しないTypeIDの場合にpanicすべき")
		}
	}()

	calculator.CalculateGuaranteedReward(stats, 10, enemyType)
}

// ==================== タスク5.2: コアドロップの品質計算テスト ====================

// TestRollCoreDropWithTypeID_GeneratesCorrectType は指定したTypeIDのコアが生成されることをテストします。
func TestRollCoreDropWithTypeID_GeneratesCorrectType(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		},
		{
			ID:           "healer",
			Name:         "ヒーラー",
			MinDropLevel: 3,
			AllowedTags:  []string{"heal_low"},
			StatWeights:  map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 0.8, "LUK": 1.2},
		},
	}

	calculator := NewRewardCalculator(coreTypes, nil, nil)

	// attack_balance を指定
	core := calculator.RollCoreDropWithTypeID("attack_balance", 10)

	if core == nil {
		t.Fatal("コアがnilであってはならない")
	}
	if core.Type.ID != "attack_balance" {
		t.Errorf("コアTypeIDが期待と異なる: got %s, want attack_balance", core.Type.ID)
	}
	if core.Type.Name != "攻撃バランス" {
		t.Errorf("コアType.Nameが期待と異なる: got %s, want 攻撃バランス", core.Type.Name)
	}
	// core.Name はType名
	expectedName := "攻撃バランス"
	if core.Name != expectedName {
		t.Errorf("コア名が期待と異なる: got %s, want %s", core.Name, expectedName)
	}
}

// TestRollCoreDropWithTypeID_InvalidTypeID は存在しないTypeIDでnilを返すことをテストします。
func TestRollCoreDropWithTypeID_InvalidTypeID(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}

	calculator := NewRewardCalculator(coreTypes, nil, nil)

	core := calculator.RollCoreDropWithTypeID("non_existent_core", 10)

	if core != nil {
		t.Error("存在しないTypeIDの場合はnilを返すべき")
	}
}

// ==================== タスク5.3: スキルドロップの品質計算テスト ====================

// TestRollSkillDropWithTypeID_GeneratesCorrectType は指定したTypeIDのスキルが生成されることをテストします。
func TestRollSkillDropWithTypeID_GeneratesCorrectType(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			Icon:         "⚔️",
			MinDropLevel: 1,
		},
		{
			ID:           "heal_lv1",
			Name:         "応急手当",
			Icon:         "💚",
			MinDropLevel: 1,
		},
	}

	calculator := NewRewardCalculator(nil, skillTypes, nil)

	// physical_lv1 を指定
	skill := calculator.RollSkillDropWithTypeID("physical_lv1", 10)

	if skill == nil {
		t.Fatal("スキルがnilであってはならない")
	}
	if skill.TypeID != "physical_lv1" {
		t.Errorf("スキルTypeIDが期待と異なる: got %s, want physical_lv1", skill.TypeID)
	}
	if skill.Name() != "物理攻撃Lv1" {
		t.Errorf("スキル名が期待と異なる: got %s, want 物理攻撃Lv1", skill.Name())
	}
}

// TestRollSkillDropWithTypeID_ChainEffectWithPool はチェイン効果プールがある場合にチェイン効果が付与されることをテストします。
func TestRollSkillDropWithTypeID_ChainEffectWithPool(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   30,
		},
	}

	pool := NewChainEffectPool(skillEffects)
	pool.SetNoEffectProbability(0.0) // チェイン効果を必ず付与

	calculator := NewRewardCalculator(nil, skillTypes, nil)
	calculator.SetChainEffectPool(pool)

	skill := calculator.RollSkillDropWithTypeID("physical_lv1", 10)

	if skill == nil {
		t.Fatal("スキルがnilであってはならない")
	}
	if !skill.HasChainEffect() {
		t.Error("チェイン効果プールがある場合はチェイン効果が付与されるべき")
	}
}

// TestRollSkillDropWithTypeID_HighLevelBetterChainEffect は高レベル敵ほど高品質チェイン効果の確率が上がることをテストします。
func TestRollSkillDropWithTypeID_HighLevelBetterChainEffect(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	skillEffects := []ChainEffectDefinition{
		{
			ID:         "damage_amp",
			Name:       "ダメージアンプ",
			Category:   "attack",
			EffectType: domain.ChainEffectDamageAmp,
			MinValue:   10,
			MaxValue:   50,
		},
	}

	pool := NewChainEffectPool(skillEffects)
	pool.SetNoEffectProbability(0.0) // チェイン効果を必ず付与

	calculator := NewRewardCalculator(nil, skillTypes, nil)
	calculator.SetChainEffectPool(pool)

	// 低レベル敵（レベル10）のチェイン効果値の平均
	lowLevelSum := 0.0
	lowLevelCount := 100
	for i := 0; i < lowLevelCount; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 10)
		if skill != nil && skill.HasChainEffect() {
			lowLevelSum += skill.ChainEffect.Value
		}
	}
	lowLevelAvg := lowLevelSum / float64(lowLevelCount)

	// 高レベル敵（レベル100）のチェイン効果値の平均
	highLevelSum := 0.0
	highLevelCount := 100
	for i := 0; i < highLevelCount; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 100)
		if skill != nil && skill.HasChainEffect() {
			highLevelSum += skill.ChainEffect.Value
		}
	}
	highLevelAvg := highLevelSum / float64(highLevelCount)

	// 高レベル敵からのスキルのチェイン効果値の平均が高いことを確認
	if highLevelAvg <= lowLevelAvg {
		t.Errorf("高レベル敵からのスキルのチェイン効果値の平均が低レベル敵より高くなるべき: lowLevelAvg=%.2f, highLevelAvg=%.2f", lowLevelAvg, highLevelAvg)
	}
}

// TestRollSkillDropWithTypeID_AlwaysHasChainEffect はスキルに必ずチェイン効果がつくことをテストします。
func TestRollSkillDropWithTypeID_AlwaysHasChainEffect(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	skillEffects := []ChainEffectDefinition{
		{
			ID:           "damage_amp",
			Name:         "ダメージアンプ",
			Category:     "attack",
			EffectType:   domain.ChainEffectDamageAmp,
			MinValue:     10,
			MaxValue:     30,
			MinDropLevel: 1,
		},
	}

	pool := NewChainEffectPool(skillEffects)

	calculator := NewRewardCalculator(nil, skillTypes, nil)
	calculator.SetChainEffectPool(pool)

	// 低レベル敵（レベル1）でも100%チェイン効果がつく
	for i := 0; i < 10; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 1)
		if skill == nil {
			t.Error("スキルがnilであるべきではない")
			continue
		}
		if !skill.HasChainEffect() {
			t.Error("スキルには必ずチェイン効果がつくべき")
		}
	}

	// 高レベル敵（レベル100）でも100%チェイン効果がつく
	for i := 0; i < 10; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 100)
		if skill == nil {
			t.Error("スキルがnilであるべきではない")
			continue
		}
		if !skill.HasChainEffect() {
			t.Error("スキルには必ずチェイン効果がつくべき")
		}
	}
}

// TestRollSkillDropWithTypeID_ChainEffectLevelFiltering はチェイン効果のMinDropLevelでフィルタリングされることをテストします。
func TestRollSkillDropWithTypeID_ChainEffectLevelFiltering(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	skillEffects := []ChainEffectDefinition{
		{
			ID:           "damage_bonus",
			Name:         "ダメージボーナス",
			Category:     "attack",
			EffectType:   domain.ChainEffectDamageBonus,
			MinValue:     10,
			MaxValue:     50,
			MinDropLevel: 1, // レベル1からドロップ
		},
		{
			ID:           "double_cast",
			Name:         "ダブルキャスト",
			Category:     "special",
			EffectType:   domain.ChainEffectDoubleCast,
			MinValue:     10,
			MaxValue:     25,
			MinDropLevel: 10, // レベル10からドロップ
		},
	}

	pool := NewChainEffectPool(skillEffects)

	calculator := NewRewardCalculator(nil, skillTypes, nil)
	calculator.SetChainEffectPool(pool)

	// レベル1の敵からはdamage_bonusのみドロップ可能
	for i := 0; i < 20; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 1)
		if skill == nil || !skill.HasChainEffect() {
			t.Error("スキルにはチェイン効果があるべき")
			continue
		}
		if skill.ChainEffect.Type != domain.ChainEffectDamageBonus {
			t.Errorf("レベル1の敵からはdamage_bonusのみドロップすべき: got %s", skill.ChainEffect.Type)
		}
	}

	// レベル10以上の敵からは両方ドロップ可能
	foundDamageBonus := false
	foundDoubleCast := false
	for i := 0; i < 100; i++ {
		skill := calculator.RollSkillDropWithTypeID("physical_lv1", 10)
		if skill == nil || !skill.HasChainEffect() {
			t.Error("スキルにはチェイン効果があるべき")
			continue
		}
		if skill.ChainEffect.Type == domain.ChainEffectDamageBonus {
			foundDamageBonus = true
		}
		if skill.ChainEffect.Type == domain.ChainEffectDoubleCast {
			foundDoubleCast = true
		}
	}
	if !foundDamageBonus || !foundDoubleCast {
		t.Errorf("レベル10以上の敵からは両方のチェイン効果がドロップすべき: foundDamageBonus=%v, foundDoubleCast=%v", foundDamageBonus, foundDoubleCast)
	}
}

// TestRollSkillDropWithTypeID_InvalidTypeID は存在しないTypeIDでnilを返すことをテストします。
func TestRollSkillDropWithTypeID_InvalidTypeID(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	calculator := NewRewardCalculator(nil, skillTypes, nil)

	skill := calculator.RollSkillDropWithTypeID("non_existent_skill", 10)

	if skill != nil {
		t.Error("存在しないTypeIDの場合はnilを返すべき")
	}
}

// TestRollSkillDropWithTypeID_NoChainEffectPool はチェイン効果プールがない場合にチェイン効果なしのスキルが生成されることをテストします。
func TestRollSkillDropWithTypeID_NoChainEffectPool(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "physical_lv1",
			Name:         "物理攻撃Lv1",
			MinDropLevel: 1,
		},
	}

	// チェイン効果プールなし
	calculator := NewRewardCalculator(nil, skillTypes, nil)

	skill := calculator.RollSkillDropWithTypeID("physical_lv1", 10)

	if skill == nil {
		t.Fatal("スキルがnilであってはならない")
	}
	if skill.HasChainEffect() {
		t.Error("チェイン効果プールがない場合はチェイン効果がnilであるべき")
	}
}

// ==================== タスク5.2: HP成長報酬テスト ====================

// TestRewardResult_HPGain はRewardResultにHP成長情報が含まれることをテストします。
func TestRewardResult_HPGain(t *testing.T) {
	result := &RewardResult{
		IsVictory:        true,
		ShowRewardScreen: true,
		HPGain:           10,
		RankUnlocked:     false,
	}

	if result.HPGain != 10 {
		t.Errorf("HPGainが期待と異なる: got %d, want 10", result.HPGain)
	}
	if result.RankUnlocked {
		t.Error("RankUnlockedがfalseであるべき")
	}
}

// TestCalculateGuaranteedRewardWithProgress はHP成長報酬付きの報酬計算をテストします。
func TestCalculateGuaranteedRewardWithProgress(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}

	calculator := NewRewardCalculator(coreTypes, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	slimeType := domain.EnemyType{
		ID:               "slime",
		Name:             "スライム",
		DropItemCategory: "core",
		DropItemTypeID:   "attack_balance",
		Rank:             1,
	}

	batType := domain.EnemyType{
		ID:               "bat",
		Name:             "コウモリ",
		DropItemCategory: "core",
		DropItemTypeID:   "attack_balance",
		Rank:             1,
	}

	// EnemyProgressManagerを作成（ランク1に2体の敵）
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := map[string]domain.EnemyType{
		"slime": slimeType,
		"bat":   batType,
	}

	// HP成長報酬付きで報酬計算（1体目）
	result := calculator.CalculateGuaranteedRewardWithProgress(stats, 1, slimeType, progress, player, enemyTypes)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if result.HPGain != domain.HPGainPerFirstDefeat {
		t.Errorf("HPGainが期待と異なる: got %d, want %d", result.HPGain, domain.HPGainPerFirstDefeat)
	}
	if result.RankUnlocked {
		t.Error("1体目撃破でランク解放されるべきではない（ランク内に2体いる）")
	}

	// PlayerのMaxHPが増加していることを確認
	expectedMaxHP := domain.InitialMaxHP + domain.HPGainPerFirstDefeat
	if player.MaxHP != expectedMaxHP {
		t.Errorf("PlayerのMaxHPが期待と異なる: got %d, want %d", player.MaxHP, expectedMaxHP)
	}

	// EnemyProgressに撃破が記録されていることを確認
	if !progress.IsDefeated("slime") {
		t.Error("敵が撃破済みとして記録されていない")
	}
}

// TestCalculateGuaranteedRewardWithProgress_RankUnlock はランク解放をテストします。
func TestCalculateGuaranteedRewardWithProgress_RankUnlock(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}

	calculator := NewRewardCalculator(coreTypes, nil, nil)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// ランク1に敵1体のみ
	enemyType := domain.EnemyType{
		ID:               "slime",
		Name:             "スライム",
		DropItemCategory: "core",
		DropItemTypeID:   "attack_balance",
		Rank:             1,
	}

	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := map[string]domain.EnemyType{
		"slime": enemyType,
	}

	// 唯一の敵を撃破
	result := calculator.CalculateGuaranteedRewardWithProgress(stats, 1, enemyType, progress, player, enemyTypes)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if !result.RankUnlocked {
		t.Error("唯一の敵撃破でランク解放されるべき")
	}

	// ランクが進行していることを確認
	if progress.CurrentRank != 2 {
		t.Errorf("ランクが進行していない: got %d, want 2", progress.CurrentRank)
	}
}
