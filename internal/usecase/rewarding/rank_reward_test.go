package rewarding

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestCalculateGuaranteedRewardWithProgress_RankUpRewards はランクアップ時に報酬が設定されることをテストします。
func TestCalculateGuaranteedRewardWithProgress_RankUpRewards(t *testing.T) {
	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
		{
			ID:           "magic_balance",
			Name:         "魔法バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"magic"},
			StatWeights:  map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 1.0, "LUK": 1.0},
		},
	}

	skillTypes := []SkillDropInfo{
		{
			ID:           "heal_lv1",
			Name:         "応急手当Lv1",
			Icon:         "💚",
			Tags:         []string{"heal"},
			MinDropLevel: 1,
		},
	}

	calculator := NewRewardCalculator(coreTypes, skillTypes, nil)

	// ランクアップ報酬を設定（ランク2到達時にコアとスキルを報酬）
	rankRewards := map[int]domain.RankReward{
		2: {
			Rank: 2,
			Items: []domain.RankRewardItem{
				{Category: "core", TypeID: "magic_balance"},
				{Category: "skill", TypeID: "heal_lv1"},
			},
		},
	}
	calculator.SetRankRewards(rankRewards)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

	// ランク1に敵1体のみ（撃破でランクアップ）
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

	result := calculator.CalculateGuaranteedRewardWithProgress(stats, 1, enemyType, progress, player, enemyTypes)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if !result.RankUnlocked {
		t.Fatal("ランク解放されるべき")
	}
	// ランクアップ報酬のコアが設定されている
	if len(result.RankUpRewardCores) != 1 {
		t.Errorf("ランクアップ報酬コア数が期待と異なる: got %d, want 1", len(result.RankUpRewardCores))
	}
	// ランクアップ報酬のスキルが設定されている
	if len(result.RankUpRewardSkills) != 1 {
		t.Errorf("ランクアップ報酬スキル数が期待と異なる: got %d, want 1", len(result.RankUpRewardSkills))
	}
}

// TestCalculateGuaranteedRewardWithProgress_RankUpNoRewards は報酬未定義ランクでランクアップ時に報酬が空であることをテストします。
func TestCalculateGuaranteedRewardWithProgress_RankUpNoRewards(t *testing.T) {
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
	// ランクアップ報酬なし（空マップ）
	calculator.SetRankRewards(map[int]domain.RankReward{})

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

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

	result := calculator.CalculateGuaranteedRewardWithProgress(stats, 1, enemyType, progress, player, enemyTypes)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if !result.RankUnlocked {
		t.Fatal("ランク解放されるべき")
	}
	// ランクアップ報酬は空
	if len(result.RankUpRewardCores) != 0 {
		t.Errorf("ランクアップ報酬コアは空であるべき: got %d", len(result.RankUpRewardCores))
	}
	if len(result.RankUpRewardSkills) != 0 {
		t.Errorf("ランクアップ報酬スキルは空であるべき: got %d", len(result.RankUpRewardSkills))
	}
}

// TestCalculateGuaranteedRewardWithProgress_RankUpSkillNoChainEffect はランクアップ報酬のスキルにチェイン効果がないことをテストします。
func TestCalculateGuaranteedRewardWithProgress_RankUpSkillNoChainEffect(t *testing.T) {
	skillTypes := []SkillDropInfo{
		{
			ID:           "heal_lv1",
			Name:         "応急手当Lv1",
			Icon:         "💚",
			Tags:         []string{"heal"},
			MinDropLevel: 1,
		},
	}

	coreTypes := []domain.CoreType{
		{
			ID:           "attack_balance",
			Name:         "攻撃バランス",
			MinDropLevel: 1,
			AllowedTags:  []string{"physical_low"},
			StatWeights:  map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
	}

	calculator := NewRewardCalculator(coreTypes, skillTypes, nil)

	// チェイン効果プールを設定（通常ドロップではチェイン効果が付与される）
	pool := NewChainEffectPool([]ChainEffectDefinition{
		{
			ID:           "damage_amp",
			Name:         "ダメージアンプ",
			Category:     "attack",
			EffectType:   domain.ChainEffectDamageAmp,
			MinValue:     10,
			MaxValue:     30,
			MinDropLevel: 1,
		},
	})
	pool.SetNoEffectProbability(0.0)
	calculator.SetChainEffectPool(pool)

	// ランクアップ報酬にスキルを設定
	rankRewards := map[int]domain.RankReward{
		2: {
			Rank: 2,
			Items: []domain.RankRewardItem{
				{Category: "skill", TypeID: "heal_lv1"},
			},
		},
	}
	calculator.SetRankRewards(rankRewards)

	stats := &BattleStatistics{
		TotalWPM:         80.0,
		TotalAccuracy:    0.95,
		TotalTypingCount: 10,
	}

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

	result := calculator.CalculateGuaranteedRewardWithProgress(stats, 1, enemyType, progress, player, enemyTypes)

	if result == nil {
		t.Fatal("報酬結果がnilであってはならない")
	}
	if len(result.RankUpRewardSkills) != 1 {
		t.Fatalf("ランクアップ報酬スキルが1つであるべき: got %d", len(result.RankUpRewardSkills))
	}
	// ランクアップ報酬のスキルにチェイン効果がないことを確認
	if result.RankUpRewardSkills[0].HasChainEffect() {
		t.Error("ランクアップ報酬のスキルにはチェイン効果を付与すべきでない")
	}
}

// TestAddRewardsToInventory_RankUpRewards はランクアップ報酬がインベントリに正しく追加されることをテストします。
func TestAddRewardsToInventory_RankUpRewards(t *testing.T) {
	// Arrange: ランクアップ報酬としてコアとスキルを設定
	core := domain.NewCoreWithTypeID("magic_balance", domain.CoreType{
		ID:          "magic_balance",
		Name:        "魔法バランス",
		StatWeights: map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 1.0, "LUK": 1.0},
	}, domain.PassiveSkill{})

	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "heal_lv1",
		Name: "応急手当Lv1",
		Icon: "💚",
		Tags: []string{"heal"},
	}, nil) // チェイン効果なし

	result := &RewardResult{
		IsVictory:          true,
		RankUnlocked:       true,
		RankUpRewardCores:  []*domain.CoreModel{core},
		RankUpRewardSkills: []*domain.SkillModel{skill},
		DroppedCores:       make([]*domain.CoreModel, 0),
		DroppedSkills:      make([]*domain.SkillModel, 0),
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert
	if !coreInv.HasCore("magic_balance") {
		t.Error("ランクアップ報酬コアがCoreInventoryに追加されるべき")
	}
	if !skillInv.HasSkill("heal_lv1") {
		t.Error("ランクアップ報酬スキルがSkillInventoryに追加されるべき")
	}
}

// TestAddRewardsToInventory_DropsAndRankUpRewardsMixed は撃破報酬とランクアップ報酬が両方ある場合に
// 全てインベントリに追加されることをテストします。
func TestAddRewardsToInventory_DropsAndRankUpRewardsMixed(t *testing.T) {
	// Arrange: 撃破報酬としてコア、ランクアップ報酬としてスキルを設定
	droppedCore := domain.NewCoreWithTypeID("attack_balance", domain.CoreType{
		ID:          "attack_balance",
		Name:        "攻撃バランス",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
	}, domain.PassiveSkill{})

	rankUpSkill := domain.NewSkillFromType(domain.SkillType{
		ID:   "heal_lv1",
		Name: "応急手当Lv1",
		Icon: "💚",
		Tags: []string{"heal"},
	}, nil)

	result := &RewardResult{
		IsVictory:          true,
		RankUnlocked:       true,
		DroppedCores:       []*domain.CoreModel{droppedCore},
		DroppedSkills:      make([]*domain.SkillModel, 0),
		RankUpRewardCores:  make([]*domain.CoreModel, 0),
		RankUpRewardSkills: []*domain.SkillModel{rankUpSkill},
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert: 撃破報酬とランクアップ報酬の両方が追加される
	if !coreInv.HasCore("attack_balance") {
		t.Error("撃破報酬コアがCoreInventoryに追加されるべき")
	}
	if !skillInv.HasSkill("heal_lv1") {
		t.Error("ランクアップ報酬スキルがSkillInventoryに追加されるべき")
	}
}

// TestAddRewardsToInventory_RankUpChainEffect はランクアップ報酬のチェイン効果がChainEffectInventoryに追加されることをテストします。
func TestAddRewardsToInventory_RankUpChainEffect(t *testing.T) {
	effect := domain.NewChainEffect("damage_bonus", domain.ChainEffectDamageBonus, 25)

	result := &RewardResult{
		IsVictory:                true,
		RankUnlocked:             true,
		DroppedCores:             make([]*domain.CoreModel, 0),
		DroppedSkills:            make([]*domain.SkillModel, 0),
		RankUpRewardCores:        make([]*domain.CoreModel, 0),
		RankUpRewardSkills:       make([]*domain.SkillModel, 0),
		RankUpRewardChainEffects: []domain.ChainEffect{effect},
	}

	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()
	chainEffectInv := domain.NewChainEffectInventory()

	// Act
	AddRewardsToInventory(result, coreInv, skillInv, chainEffectInv)

	// Assert
	if !chainEffectInv.HasChainEffect("damage_bonus") {
		t.Error("ランクアップ報酬チェイン効果がChainEffectInventoryに追加されるべき")
	}
}

// TestRewardResult_NoDropsNoHPGain は撃破報酬もHP増加もない場合のRewardResultをテストします。
func TestRewardResult_NoDropsNoHPGain(t *testing.T) {
	result := &RewardResult{
		IsVictory:        true,
		ShowRewardScreen: true,
		DroppedCores:     make([]*domain.CoreModel, 0),
		DroppedSkills:    make([]*domain.SkillModel, 0),
		HPGain:           0,
		RankUnlocked:     false,
	}

	hasDrops := len(result.DroppedCores) > 0 || len(result.DroppedSkills) > 0
	hasDefeatReward := hasDrops || result.HPGain > 0

	if hasDefeatReward {
		t.Error("撃破報酬がない場合、hasDefeatRewardはfalseであるべき")
	}
}
