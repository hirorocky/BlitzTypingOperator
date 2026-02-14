// Package startup は初回起動時の初期化処理を担当します。

package startup

import (
	"testing"

	"hirorocky/type-battle/internal/infra/masterdata"
)

// createTestExternalData はテスト用の外部データを作成します。
func createTestExternalData() *masterdata.ExternalData {
	return &masterdata.ExternalData{
		CoreTypes: []masterdata.CoreTypeData{
			{
				ID:             "all_rounder",
				Name:           "オールラウンダー",
				AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
				StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
				PassiveSkillID: "ps_combo_master",
				MinDropLevel:   1,
			},
			{
				ID:             "healer",
				Name:           "ヒーラー",
				AllowedTags:    []string{"heal_low", "heal_mid", "heal_high"},
				StatWeights:    map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 0.8, "LUK": 1.2},
				PassiveSkillID: "ps_miracle_heal",
				MinDropLevel:   3,
			},
			{
				ID:             "quick_recoverer",
				Name:           "クイックリカバラー",
				AllowedTags:    []string{"heal_low", "heal_mid", "buff_low"},
				StatWeights:    map[string]float64{"STR": 0.7, "INT": 1.2, "WIL": 1.1, "LUK": 1.0},
				PassiveSkillID: "ps_quick_recovery",
				MinDropLevel:   4,
			},
		},
		SkillDefinitions: []masterdata.SkillDefinitionData{
			{
				ID:          "physical_strike_lv1",
				Name:        "物理打撃Lv1",
				Icon:        "⚔️",
				Tags:        []string{"physical_low"},
				Description: "物理ダメージを与える基本攻撃",
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "enemy",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 1.0, StatRef: "STR"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:          "fireball_lv1",
				Name:        "ファイアボールLv1",
				Icon:        "🔥",
				Tags:        []string{"magic_low"},
				Description: "魔法ダメージを与える基本魔法",
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "enemy",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 1.2, StatRef: "MAG"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:          "heal_lv1",
				Name:        "ヒールLv1",
				Icon:        "💚",
				Tags:        []string{"heal_low"},
				Description: "HPを回復する基本回復魔法",
				Effects: []masterdata.SkillEffectData{
					{
						Target:      "self",
						HPFormula:   &masterdata.HPFormulaData{Base: 0, StatCoef: 0.8, StatRef: "MAG"},
						Probability: 1.0,
					},
				},
			},
			{
				ID:          "attack_buff_lv1",
				Name:        "攻撃バフLv1",
				Icon:        "⬆️",
				Tags:        []string{"buff_low"},
				Description: "一時的に攻撃力を上昇させる",
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
				ID:     "slime",
				Name:   "スライム",
				BaseHP: 50,
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
				Skills: []masterdata.FirstAgentSkillData{
					{TypeID: "physical_strike_lv1"},
				},
			},
			{
				ID:         "agent_first_2",
				CoreTypeID: "healer",
				CoreLevel:  1,
				Skills: []masterdata.FirstAgentSkillData{
					{TypeID: "heal_lv1"},
				},
			},
			{
				ID:         "agent_first_3",
				CoreTypeID: "quick_recoverer",
				CoreLevel:  1,
				Skills: []masterdata.FirstAgentSkillData{
					{TypeID: "attack_buff_lv1"},
				},
			},
		},
	}
}

// ==================================================
// Task 14.1: 新規ゲーム初期化テスト
// ==================================================

func TestNewGameInitializer_CreateInitialAgents(t *testing.T) {
	initializer := NewNewGameInitializer(createTestExternalData())

	agents := initializer.CreateInitialAgents()
	if agents == nil {
		t.Fatal("初期エージェントが作成されるべきです")
	}

	// 3体のエージェントが作成されること
	if len(agents) != 3 {
		t.Fatalf("初期エージェントは3体作成されるべきです: got %d", len(agents))
	}

	expectedCores := []string{"all_rounder", "healer", "quick_recoverer"}
	for i, agent := range agents {
		// エージェントがコアを持つこと
		if agent.Core == nil {
			t.Errorf("初期エージェント%dはコアを持つべきです", i+1)
			continue
		}

		// エージェントが1つのスキルを持つこと
		if len(agent.Skills) != 1 {
			t.Errorf("初期エージェント%dは1つのスキルを持つべきです: got %d", i+1, len(agent.Skills))
		}

		// 期待するコア特性であること
		if agent.Core.Type.ID != expectedCores[i] {
			t.Errorf("初期エージェント%dのコアは%sであるべきです: got %s", i+1, expectedCores[i], agent.Core.Type.ID)
		}
	}
}

func TestNewGameInitializer_InitializeNewGame(t *testing.T) {
	initializer := NewNewGameInitializer(createTestExternalData())

	saveData := initializer.InitializeNewGame()
	if saveData == nil {
		t.Fatal("新規ゲームデータが作成されるべきです")
	}

	// AgentSlotsに初期エージェントが3体設定されていること
	equippedCount := 0
	for _, slot := range saveData.Player.AgentSlots {
		if slot.CoreTypeID != "" {
			equippedCount++
		}
	}
	if equippedCount != 3 {
		t.Errorf("初期エージェントが3体装備されているべきです: got %d", equippedCount)
	}

	// 各スロットにコアとスキルが設定されていること
	for i, slot := range saveData.Player.AgentSlots {
		if slot.CoreTypeID == "" {
			t.Errorf("スロット%dにコアが設定されているべきです", i)
		}
		// 少なくとも1つのスキルが設定されていること
		hasSkill := false
		for _, skill := range slot.Skills {
			if skill.TypeID != "" {
				hasSkill = true
				break
			}
		}
		if !hasSkill {
			t.Errorf("スロット%dに少なくとも1つのスキルが設定されているべきです", i)
		}
	}
}

func TestNewGameInitializer_InitialStats(t *testing.T) {
	// 新規ゲーム開始時の統計情報がリセットされていること
	initializer := NewNewGameInitializer(createTestExternalData())

	saveData := initializer.InitializeNewGame()

	if saveData.Statistics.TotalBattles != 0 {
		t.Error("総バトル数は0であるべきです")
	}
	if saveData.Statistics.Victories != 0 {
		t.Error("勝利数は0であるべきです")
	}
	if saveData.Statistics.MaxLevelReached != 0 {
		t.Error("到達最高レベルは0であるべきです")
	}
}

func TestNewGameInitializer_InitialAchievements(t *testing.T) {
	// 新規ゲーム開始時の実績がリセットされていること
	initializer := NewNewGameInitializer(createTestExternalData())

	saveData := initializer.InitializeNewGame()

	if len(saveData.Achievements.Unlocked) != 0 {
		t.Error("解除済み実績は空であるべきです")
	}
}

func TestInitialAgent_SkillsCompatibleWithCore(t *testing.T) {
	// 初期エージェントのスキルがコアと互換性があること
	initializer := NewNewGameInitializer(createTestExternalData())

	agents := initializer.CreateInitialAgents()

	for agentIdx, agent := range agents {
		for i, skill := range agent.Skills {
			if !skill.IsCompatibleWithCore(agent.Core) {
				t.Errorf("エージェント%dのスキル%dはコアと互換性があるべきです", agentIdx+1, i)
			}
		}
	}
}

func TestNewGameInitializer_MultipleCalls(t *testing.T) {
	// 複数回呼び出しても毎回新しいセーブデータオブジェクトが作成されること
	initializer := NewNewGameInitializer(createTestExternalData())

	saveData1 := initializer.InitializeNewGame()
	saveData2 := initializer.InitializeNewGame()

	// 別のセーブデータオブジェクトが作成されていること
	if saveData1 == saveData2 {
		t.Error("異なる呼び出しで異なるセーブデータオブジェクトが作成されるべきです")
	}

	// 両方のセーブデータにAgentSlotsが設定されていること
	hasSlot1 := saveData1.Player.AgentSlots[0].CoreTypeID != ""
	hasSlot2 := saveData2.Player.AgentSlots[0].CoreTypeID != ""
	if !hasSlot1 {
		t.Error("saveData1にエージェントスロットが設定されているべきです")
	}
	if !hasSlot2 {
		t.Error("saveData2にエージェントスロットが設定されているべきです")
	}
}
