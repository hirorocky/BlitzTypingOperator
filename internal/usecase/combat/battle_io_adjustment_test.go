package combat

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// === 受け入れ基準 9,10,11: ダメージ計算式がHP変化量 = base + stat_coef × STAT ===

func TestApplySkillEffect_DamageFormula_BaseStatOnly(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "test", Name: "テスト", BaseHP: 1000},
	}
	engine := NewBattleEngine(enemyTypes)

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test_ps", Name: "テスト"}
	core := domain.NewCoreWithTypeID("core_test", coreType, passiveSkill)

	// base=10, statCoef=2.0, STR参照
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "sk_base",
		Name: "ベーススキル",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 10, StatCoef: 2.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
	agent := domain.NewAgent("agent_test", core, []*domain.SkillModel{skill})

	state := createTestBattleState(agent, 1000)
	result := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	effectResult := engine.ApplySkillEffect(state, agent, skill, result)
	// STR=100（StatWeights{"STR":1.0} → 100×1.0=100）
	// 期待値: base + statCoef * STR = 10 + 2.0 * 100 = 210
	expected := 210
	if effectResult.TotalDamage != expected {
		t.Errorf("ダメージ = base + stat_coef × STAT: got %d, want %d", effectResult.TotalDamage, expected)
	}
}

// === 受け入れ基準 14: SkillEffectResult に通常効果と潜在効果が分離記録される ===

func TestApplySkillEffect_SkillEffectResult_NormalAndLatentEffects(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "test", Name: "テスト", BaseHP: 1000},
	}
	engine := NewBattleEngine(enemyTypes)

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test_ps", Name: "テスト"}
	core := domain.NewCoreWithTypeID("core_test", coreType, passiveSkill)

	// 通常効果 + 潜在効果のスキル
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "sk_mixed",
		Name: "混合スキル",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				IsLatent:    false,
				Icon:        "⚔️",
			},
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 30, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				IsLatent:    true,
				Icon:        "✨",
			},
		},
	}, nil)
	agent := domain.NewAgent("agent_test", core, []*domain.SkillModel{skill})

	t.Run("パーフェクト時は通常効果と潜在効果の両方が記録される", func(t *testing.T) {
		state := createTestBattleState(agent, 1000)
		result := &typing.TypingResult{
			Completed: true,
			WPM:       60.0,
			Score:     100,
			IsPerfect: true,
		}

		effectResult := engine.ApplySkillEffect(state, agent, skill, result)

		if len(effectResult.NormalEffects) != 1 {
			t.Errorf("NormalEffects = %d件, want 1件", len(effectResult.NormalEffects))
		}
		if len(effectResult.LatentEffects) != 1 {
			t.Errorf("LatentEffects = %d件, want 1件", len(effectResult.LatentEffects))
		}

		if len(effectResult.NormalEffects) > 0 {
			if !effectResult.NormalEffects[0].TargetIsEnemy {
				t.Error("NormalEffects[0].TargetIsEnemy = false, want true")
			}
			if effectResult.NormalEffects[0].HPChange >= 0 {
				t.Errorf("NormalEffects[0].HPChange = %d, want negative (damage)", effectResult.NormalEffects[0].HPChange)
			}
		}
	})

	t.Run("非パーフェクト時は通常効果のみ記録される", func(t *testing.T) {
		state := createTestBattleState(agent, 1000)
		result := &typing.TypingResult{
			Completed: true,
			WPM:       60.0,
			Score:     90,
			IsPerfect: false,
		}

		effectResult := engine.ApplySkillEffect(state, agent, skill, result)

		if len(effectResult.NormalEffects) != 1 {
			t.Errorf("NormalEffects = %d件, want 1件", len(effectResult.NormalEffects))
		}
		if len(effectResult.LatentEffects) != 0 {
			t.Errorf("LatentEffects = %d件, want 0件（非パーフェクト）", len(effectResult.LatentEffects))
		}
	})
}

// === 受け入れ基準: BattleStatistics TotalScore/AverageScore ===

func TestBattleStatistics_TotalScoreとAverageScore(t *testing.T) {
	stats := &BattleStatistics{}

	// 3回タイピング結果を記録（Score=100, 80, 90）
	stats.TotalTypingCount = 3
	stats.TotalScore = 100 + 80 + 90
	stats.TotalWPM = 60.0 + 70.0 + 80.0

	if stats.TotalScore != 270 {
		t.Errorf("TotalScore = %d, want 270", stats.TotalScore)
	}

	avg := stats.GetAverageScore()
	if avg != 90.0 {
		t.Errorf("GetAverageScore() = %f, want 90.0", avg)
	}
}

func TestRecordTypingResult_ScoreをTotalScoreに蓄積(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "test", Name: "テスト", BaseHP: 100},
	}
	engine := NewBattleEngine(enemyTypes)

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test_ps", Name: "テスト"}
	core := domain.NewCoreWithTypeID("core_test", coreType, passiveSkill)
	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "sk_test",
		Name: "テスト",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{Target: domain.TargetEnemy, HPFormula: &domain.HPFormula{Base: 10}, Probability: 1.0},
		},
	}, nil)
	agent := domain.NewAgent("agent_test", core, []*domain.SkillModel{skill})
	state := createTestBattleState(agent, 100)

	// Score=95のタイピング結果を記録
	result := &typing.TypingResult{
		Completed: true,
		WPM:       80.0,
		Score:     95,
	}
	engine.RecordTypingResult(state, result)

	if state.Stats.TotalScore != 95 {
		t.Errorf("TotalScore = %d, want 95", state.Stats.TotalScore)
	}
	if state.Stats.TotalTypingCount != 1 {
		t.Errorf("TotalTypingCount = %d, want 1", state.Stats.TotalTypingCount)
	}
}

// === ダメージ計算がScoreに依存しない回帰テスト ===

func TestApplySkillEffect_ダメージはScoreに依存しない(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "test", Name: "テスト", BaseHP: 1000},
	}
	engine := NewBattleEngine(enemyTypes)

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test_ps", Name: "テスト"}
	core := domain.NewCoreWithTypeID("core_test", coreType, passiveSkill)

	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "sk_test",
		Name: "テスト",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
	agent := domain.NewAgent("agent_test", core, []*domain.SkillModel{skill})

	// Score=100のタイピング結果
	state1 := createTestBattleState(agent, 1000)
	highScoreResult := &typing.TypingResult{Completed: true, WPM: 60.0, Score: 100}
	effect1 := engine.ApplySkillEffect(state1, agent, skill, highScoreResult)

	// Score=30のタイピング結果（旧仕様ではAccuracyPenaltyが適用されていた）
	state2 := createTestBattleState(agent, 1000)
	lowScoreResult := &typing.TypingResult{Completed: true, WPM: 60.0, Score: 30}
	effect2 := engine.ApplySkillEffect(state2, agent, skill, lowScoreResult)

	// 新仕様ではダメージはScoreに依存しない（base + stat_coef × STATのみ）
	if effect1.TotalDamage != effect2.TotalDamage {
		t.Errorf("Score=%d → ダメージ=%d, Score=%d → ダメージ=%d: ダメージがScoreに依存しています",
			highScoreResult.Score, effect1.TotalDamage,
			lowScoreResult.Score, effect2.TotalDamage)
	}
}

// === calculateRawHPChange: EffectTable補正前の素の値 ===

func TestCalculateRawHPChange_EffectTable補正前の値(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "test", Name: "テスト", BaseHP: 1000},
	}
	engine := NewBattleEngine(enemyTypes)

	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test_ps", Name: "テスト"}
	core := domain.NewCoreWithTypeID("core_test", coreType, passiveSkill)

	skill := domain.NewSkillFromType(domain.SkillType{
		ID:   "sk_test",
		Name: "テスト",
		Icon: "⚔️",
		Tags: []string{"physical_low"},
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 10, StatCoef: 2.0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
	agent := domain.NewAgent("agent_test", core, []*domain.SkillModel{skill})

	// DamageMultiplierバフを付与した状態で効果適用
	state := createTestBattleState(agent, 1000)
	state.Player.EffectTable.AddBuff("dmg_buff", "ダメージ強化", 3, map[domain.EffectColumn]float64{
		domain.ColDamageMultiplier: 2.0, // ダメージ2倍
	})

	result := &typing.TypingResult{Completed: true, WPM: 60.0, Score: 100}
	effectResult := engine.ApplySkillEffect(state, agent, skill, result)

	// STR=100, base=10, statCoef=2.0 → 素のダメージ = 10 + 2.0*100 = 210
	// DamageMultiplier=2.0 → 実際のダメージ = 210 * 2.0 = 420
	// EffectDetailには補正前の素の値（210）が記録されるべき
	if effectResult.TotalDamage != 420 {
		t.Errorf("実際のダメージ = %d, want 420（DamageMultiplier適用後）", effectResult.TotalDamage)
	}
	if len(effectResult.NormalEffects) != 1 {
		t.Fatalf("NormalEffects = %d件, want 1件", len(effectResult.NormalEffects))
	}
	rawDamage := effectResult.NormalEffects[0].HPChange
	if rawDamage != -210 {
		t.Errorf("EffectDetail.HPChange = %d, want -210（EffectTable補正前の素の値）", rawDamage)
	}
}

// === BattleStatistics.GetAverageScore ゼロ除算テスト ===

func TestBattleStatistics_GetAverageScore_ゼロ除算回避(t *testing.T) {
	stats := &BattleStatistics{
		TotalScore:       0,
		TotalTypingCount: 0,
	}
	avg := stats.GetAverageScore()
	if avg != 0 {
		t.Errorf("TotalTypingCount=0のとき、GetAverageScore() = %f, want 0", avg)
	}
}

// === ヘルパー関数 ===

func createTestBattleState(agent *domain.AgentModel, enemyHP int) *BattleState {
	enemy := domain.NewEnemy("e1", "テスト敵", 1, enemyHP, 10, domain.EnemyType{
		ID:   "test",
		Name: "テスト",
	})
	player := domain.NewPlayerWithMaxHP(100)
	player.PrepareForBattle()

	return &BattleState{
		Enemy:          enemy,
		Player:         player,
		EquippedAgents: []*domain.AgentModel{agent},
		Level:          1,
		Stats:          &BattleStatistics{},
	}
}
