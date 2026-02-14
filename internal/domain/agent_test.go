// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// newTestDamageSkill はテスト用ダメージスキルを作成するヘルパー関数です。
func newTestDamageSkill(id, name string, tags []string, statCoef float64, statRef, description string) *SkillModel {
	return NewSkillFromType(SkillType{
		ID:          id,
		Name:        name,
		Icon:        "⚔️",
		Tags:        tags,
		Description: description,
		Effects: []SkillEffect{
			{
				Target:      TargetEnemy,
				HPFormula:   &HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
}

// TestAgentModel_フィールドの確認 はAgentModel構造体のフィールドが正しく設定されることを確認します。
func TestAgentModel_フィールドの確認(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	passiveSkill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	skills := []*SkillModel{
		newTestDamageSkill("mod_001", "物理打撃Lv1", []string{"physical_low"}, 1.0, "STR", "物理攻撃"),
		newTestDamageSkill("mod_002", "ファイアボールLv1", []string{"magic_low"}, 1.0, "MAG", "魔法攻撃"),
		newTestDamageSkill("mod_003", "物理打撃Lv1", []string{"physical_low"}, 1.0, "STR", "物理攻撃"),
		newTestDamageSkill("mod_004", "ファイアボールLv1", []string{"magic_low"}, 1.0, "MAG", "魔法攻撃"),
	}

	agent := AgentModel{
		ID:        "agent_001",
		Core:      core,
		Skills:    skills,
		BaseStats: core.Stats, // 基礎ステータス = コアのステータス
	}

	if agent.ID != "agent_001" {
		t.Errorf("IDが期待値と異なります: got %s, want agent_001", agent.ID)
	}
	if agent.Core.TypeID != "attack_balance" {
		t.Errorf("Core.TypeIDが期待値と異なります: got %s, want attack_balance", agent.Core.TypeID)
	}
	if len(agent.Skills) != 4 {
		t.Errorf("Skillsの長さが期待値と異なります: got %d, want 4", len(agent.Skills))
	}
	// STR: 100 × 1.2 = 120
	if agent.BaseStats.STR != 120 {
		t.Errorf("BaseStats.STRが期待値と異なります: got %d, want 120", agent.BaseStats.STR)
	}
}

// TestAgentModel_ステータス導出 はエージェントのステータスがコアから正しく導出されることを確認します。

func TestAgentModel_ステータス導出(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}
	core := NewCoreWithTypeID("test", coreType, passiveSkill)

	skills := []*SkillModel{
		newTestDamageSkill("mod_001", "テストスキル1", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_002", "テストスキル2", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_003", "テストスキル3", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_004", "テストスキル4", []string{"physical_low"}, 1.0, "STR", "テスト"),
	}

	agent := NewAgent("agent_test", core, skills)

	// ステータスがコアから導出されていることを確認
	// STR: 100 × 1.0 = 100
	if agent.BaseStats.STR != 100 {
		t.Errorf("BaseStats.STRがコアのステータスと一致しません: got %d, want %d", agent.BaseStats.STR, core.Stats.STR)
	}
}

// TestNewAgent_エージェント作成 はNewAgent関数でエージェントが正しく作成されることを確認します。
func TestNewAgent_エージェント作成(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	passiveSkill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	skills := []*SkillModel{
		newTestDamageSkill("mod_001", "物理打撃Lv1", []string{"physical_low"}, 1.0, "STR", "物理攻撃"),
		newTestDamageSkill("mod_002", "ファイアボールLv1", []string{"magic_low"}, 1.0, "MAG", "魔法攻撃"),
		newTestDamageSkill("mod_003", "物理打撃Lv1", []string{"physical_low"}, 1.0, "STR", "物理攻撃"),
		newTestDamageSkill("mod_004", "ファイアボールLv1", []string{"magic_low"}, 1.0, "MAG", "魔法攻撃"),
	}

	agent := NewAgent("agent_001", core, skills)

	if agent.ID != "agent_001" {
		t.Errorf("IDが期待値と異なります: got %s, want agent_001", agent.ID)
	}
	// 基礎ステータスはコアから導出される
	// STR: 100 × 1.2 = 120
	if agent.BaseStats.STR != 120 {
		t.Errorf("BaseStats.STRが期待値と異なります: got %d, want 120", agent.BaseStats.STR)
	}
}

// TestNewAgent_スキル数制約 はNewAgent関数がスキル数を検証することを確認します。
// エージェントは必ず4個のスキルを装備する必要があります
func TestNewAgent_スキル数確認(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}
	core := NewCoreWithTypeID("test", coreType, passiveSkill)

	skills := []*SkillModel{
		newTestDamageSkill("mod_001", "テストスキル1", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_002", "テストスキル2", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_003", "テストスキル3", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_004", "テストスキル4", []string{"physical_low"}, 1.0, "STR", "テスト"),
	}

	agent := NewAgent("agent_test", core, skills)

	// 4個のスキルが装備されていることを確認
	if len(agent.Skills) != 4 {
		t.Errorf("Skillsの長さが4でありません: got %d, want 4", len(agent.Skills))
	}
}

// TestAgentModel_基礎ステータス算出 は基礎ステータスがコアから正しく導出されることを確認します。

func TestAgentModel_基礎ステータス算出(t *testing.T) {
	tests := []struct {
		name        string
		coreType    CoreType
		expectedSTR int
		expectedINT int
		expectedWIL int
		expectedLUK int
	}{
		{
			name: "攻撃バランス型",
			coreType: CoreType{
				ID:          "attack_balance",
				Name:        "攻撃バランス",
				StatWeights: map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
				AllowedTags: []string{"physical_low"},
			},
			expectedSTR: 120, // 100 × 1.2
			expectedINT: 100, // 100 × 1.0
			expectedWIL: 80,  // 100 × 0.8
			expectedLUK: 100, // 100 × 1.0
		},
		{
			name: "ヒーラー型",
			coreType: CoreType{
				ID:          "healer",
				Name:        "ヒーラー",
				StatWeights: map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 0.8, "LUK": 1.2},
				AllowedTags: []string{"heal_low"},
			},
			expectedSTR: 50,  // 100 × 0.5
			expectedINT: 150, // 100 × 1.5
			expectedWIL: 80,  // 100 × 0.8
			expectedLUK: 120, // 100 × 1.2
		},
		{
			name: "オールラウンダー型",
			coreType: CoreType{
				ID:          "all_rounder",
				Name:        "オールラウンダー",
				StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
				AllowedTags: []string{"physical_low"},
			},
			expectedSTR: 100, // 100 × 1.0
			expectedINT: 100, // 100 × 1.0
			expectedWIL: 100, // 100 × 1.0
			expectedLUK: 100, // 100 × 1.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			passiveSkill := PassiveSkill{ID: "test_skill"}
			core := NewCoreWithTypeID(tt.coreType.ID, tt.coreType, passiveSkill)

			skills := make([]*SkillModel, 4)
			for i := 0; i < 4; i++ {
				skills[i] = newTestDamageSkill("mod", "テスト", []string{"physical_low"}, 1.0, "STR", "テスト")
			}

			agent := NewAgent("agent_test", core, skills)

			if agent.BaseStats.STR != tt.expectedSTR {
				t.Errorf("BaseStats.STRが期待値と異なります: got %d, want %d", agent.BaseStats.STR, tt.expectedSTR)
			}
			if agent.BaseStats.INT != tt.expectedINT {
				t.Errorf("BaseStats.INTが期待値と異なります: got %d, want %d", agent.BaseStats.INT, tt.expectedINT)
			}
			if agent.BaseStats.WIL != tt.expectedWIL {
				t.Errorf("BaseStats.WILが期待値と異なります: got %d, want %d", agent.BaseStats.WIL, tt.expectedWIL)
			}
			if agent.BaseStats.LUK != tt.expectedLUK {
				t.Errorf("BaseStats.LUKが期待値と異なります: got %d, want %d", agent.BaseStats.LUK, tt.expectedLUK)
			}
		})
	}
}

// TestAgentModel_Skills はエージェントから指定インデックスのスキルを直接取得できることを確認します。
func TestAgentModel_Skills(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}
	core := NewCoreWithTypeID("test", coreType, passiveSkill)

	skills := []*SkillModel{
		newTestDamageSkill("mod_001", "スキル1", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_002", "スキル2", []string{"physical_low"}, 1.5, "STR", "テスト"),
		newTestDamageSkill("mod_003", "スキル3", []string{"physical_low"}, 2.0, "STR", "テスト"),
		newTestDamageSkill("mod_004", "スキル4", []string{"physical_low"}, 2.5, "STR", "テスト"),
	}

	agent := NewAgent("agent_test", core, skills)

	// 正常系: 各インデックスのスキルを取得（直接アクセス）
	for i := 0; i < 4; i++ {
		skill := agent.Skills[i]
		if skill == nil {
			t.Errorf("インデックス%dのスキルがnilです", i)
			continue
		}
		if skill.TypeID != skills[i].TypeID {
			t.Errorf("インデックス%dのスキルTypeIDが異なります: got %s, want %s", i, skill.TypeID, skills[i].TypeID)
		}
	}

	// スキル数の確認
	if len(agent.Skills) != 4 {
		t.Errorf("スキル数が4でありません: got %d, want 4", len(agent.Skills))
	}
}

// TestAgentModel_スキルの独立性 はNewAgentで作成したエージェントのSkillsが元のスライスと独立していることを確認します。
func TestAgentModel_スキルの独立性(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}
	core := NewCoreWithTypeID("test", coreType, passiveSkill)

	originalSkills := []*SkillModel{
		newTestDamageSkill("mod_001", "スキル1", []string{"physical_low"}, 1.0, "STR", "テスト"),
		newTestDamageSkill("mod_002", "スキル2", []string{"physical_low"}, 1.5, "STR", "テスト"),
		newTestDamageSkill("mod_003", "スキル3", []string{"physical_low"}, 2.0, "STR", "テスト"),
		newTestDamageSkill("mod_004", "スキル4", []string{"physical_low"}, 2.5, "STR", "テスト"),
	}

	agent := NewAgent("agent_test", core, originalSkills)

	// 元のスライスを変更
	originalSkills[0] = newTestDamageSkill("mod_changed", "変更済み", []string{"physical_low"}, 9.9, "STR", "変更")

	// エージェントのスキルは影響を受けないはず
	if agent.Skills[0].TypeID == "mod_changed" {
		t.Error("AgentModelのSkillsが元のスライスの変更の影響を受けています")
	}
}
