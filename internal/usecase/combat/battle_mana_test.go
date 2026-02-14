// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// ===== 受け入れ基準 12: スキル使用成功時にManaCost分のマナが消費される =====

func TestApplySkillEffect_ManaCost消費(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 5
	player.PrepareForBattle()
	player.Mana = 5 // PrepareForBattle後に再設定

	moduleType := domain.SkillType{
		ID:       "magic_skill",
		Name:     "テスト魔法",
		ManaCost: 1,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 100, StatCoef: 0, StatRef: "INT"},
				Probability: 1.0,
				Icon:        "💥",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act
	engine.ApplySkillEffect(state, agent, module, typingResult)

	// Assert
	if state.Player.Mana != 4 {
		t.Errorf("マナ消費後のManaが期待値と異なります: got %d, want 4", state.Player.Mana)
	}
}

// ===== 受け入れ基準 13: 効果適用時にManaGain分のマナが獲得される =====

func TestApplySkillEffect_ManaGain獲得(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 0
	player.PrepareForBattle()

	moduleType := domain.SkillType{
		ID:   "martial_skill",
		Name: "テスト格闘",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				ManaGain:    1,
				Icon:        "👊",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act
	engine.ApplySkillEffect(state, agent, module, typingResult)

	// Assert
	if state.Player.Mana != 1 {
		t.Errorf("マナ獲得後のManaが期待値と異なります: got %d, want 1", state.Player.Mana)
	}
}

// ===== 受け入れ基準 10: マナ獲得は確率判定に従う =====

func TestApplySkillEffect_ManaGain確率で不発(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 0
	player.PrepareForBattle()

	moduleType := domain.SkillType{
		ID:   "martial_skill",
		Name: "テスト格闘",
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 0.0, // 絶対不発
				ManaGain:    1,
				Icon:        "👊",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act
	engine.ApplySkillEffect(state, agent, module, typingResult)

	// Assert: 確率0なので効果不発、マナ獲得もなし
	if state.Player.Mana != 0 {
		t.Errorf("不発時にマナが変動しています: got %d, want 0", state.Player.Mana)
	}
}

// ===== 受け入れ基準 14: ManaCost=0のスキルはマナに影響しない =====

func TestApplySkillEffect_ManaCostゼロ_マナ不変(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 5
	player.PrepareForBattle()
	player.Mana = 5

	moduleType := domain.SkillType{
		ID:       "normal_skill",
		Name:     "通常攻撃",
		ManaCost: 0,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 100, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act
	engine.ApplySkillEffect(state, agent, module, typingResult)

	// Assert: ManaCost=0なのでマナ変動なし
	if state.Player.Mana != 5 {
		t.Errorf("ManaCost=0のスキル使用でマナが変動しています: got %d, want 5", state.Player.Mana)
	}
}

// ===== 受け入れ基準 13: 消費→効果→獲得の順序 =====

func TestApplySkillEffect_消費と獲得の順序(t *testing.T) {
	// Arrange: ManaCost=1で消費し、ManaGain=1で獲得 → 結果変わらず
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 1
	player.PrepareForBattle()
	player.Mana = 1

	moduleType := domain.SkillType{
		ID:       "combo_skill",
		Name:     "コンボスキル",
		ManaCost: 1,
		Effects: []domain.SkillEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 50, StatCoef: 0, StatRef: "STR"},
				Probability: 1.0,
				ManaGain:    1,
				Icon:        "🔥",
			},
		},
	}
	module := domain.NewSkillFromType(moduleType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            60.0,
		Accuracy:       1.0,
		SpeedFactor:    1.0,
		AccuracyFactor: 1.0,
	}

	// Act
	engine.ApplySkillEffect(state, agent, module, typingResult)

	// Assert: 1消費 + 1獲得 = 結果1のまま
	if state.Player.Mana != 1 {
		t.Errorf("消費→獲得後のManaが期待値と異なります: got %d, want 1", state.Player.Mana)
	}
}

// ===== ヘルパー関数 =====

func createTestAgent() *domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test_core",
		Name:        "テストコア",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low"},
	}
	core := domain.NewCoreWithTypeID("core_001", coreType, domain.PassiveSkill{})
	return domain.NewAgent("agent_001", core, nil)
}

func createManaTestState(player *domain.PlayerModel) *BattleState {
	enemyType := domain.EnemyType{ID: "test_enemy", Name: "テスト敵", BaseHP: 10000}
	enemy := domain.NewEnemy("test_enemy", "テスト敵", 1, 10000, 100, enemyType)
	return &BattleState{
		Enemy:  enemy,
		Player: player,
		Stats:  &BattleStatistics{},
	}
}
