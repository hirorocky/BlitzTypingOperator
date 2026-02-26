// Package combat はバトル関連のユースケースを提供します。
package combat

import (
	"math/rand"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/typing"
)

// ===== 受け入れ基準 12: マナ消費はバトルロジック層で1回だけ行われる =====
// ApplySkillEffectはマナを消費しない（Echo/DoubleCast対応のため呼び出し側で消費）

func TestApplySkillEffect_マナ消費はApplySkillEffectの責務外(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.PrepareForBattle()
	player.Mana = 5

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	// Act
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ApplySkillEffect内ではマナを消費しない
	if state.Player.Mana != 5 {
		t.Errorf("ApplySkillEffect内でマナが消費されています: got %d, want 5", state.Player.Mana)
	}
}

// ===== 受け入れ基準 13: 効果適用時にManaGain分のマナが獲得される =====

func TestApplySkillEffect_ManaGain獲得(t *testing.T) {
	// Arrange
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 0
	player.PrepareForBattle()

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	// Act
	engine.ApplySkillEffect(state, agent, skill, typingResult)

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

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	// Act
	engine.ApplySkillEffect(state, agent, skill, typingResult)

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

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	// Act
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: ManaCost=0なのでマナ変動なし
	if state.Player.Mana != 5 {
		t.Errorf("ManaCost=0のスキル使用でマナが変動しています: got %d, want 5", state.Player.Mana)
	}
}

// ===== 受け入れ基準 13: ApplySkillEffect内ではManaGainのみ発生する =====

func TestApplySkillEffect_ManaGainのみ発生_消費は呼び出し側(t *testing.T) {
	// Arrange: ManaCost=1, ManaGain=1のスキル → ApplySkillEffect内ではGainのみ
	player := domain.NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.PrepareForBattle()
	player.Mana = 1

	skillType := domain.SkillType{
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
	skill := domain.NewSkillFromType(skillType, nil)

	agent := createTestAgent()
	state := createManaTestState(player)

	engine := NewBattleEngine(nil)
	engine.SetRng(rand.New(rand.NewSource(42)))

	typingResult := &typing.TypingResult{
		Completed: true,
		WPM:       60.0,
		Score:     100,
	}

	// Act
	engine.ApplySkillEffect(state, agent, skill, typingResult)

	// Assert: 消費なし + 1獲得 = Mana 2
	if state.Player.Mana != 2 {
		t.Errorf("ApplySkillEffect後のManaが期待値と異なります: got %d, want 2", state.Player.Mana)
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
