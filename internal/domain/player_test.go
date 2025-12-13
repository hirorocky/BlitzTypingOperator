// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestPlayerModel_フィールドの確認 はPlayerModel構造体のフィールドが正しく設定されることを確認します。
func TestPlayerModel_フィールドの確認(t *testing.T) {
	player := PlayerModel{
		HP:    100,
		MaxHP: 100,
	}

	if player.HP != 100 {
		t.Errorf("HPが期待値と異なります: got %d, want 100", player.HP)
	}
	if player.MaxHP != 100 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 100", player.MaxHP)
	}
}

// TestNewPlayer_プレイヤー作成 はNewPlayer関数でプレイヤーが正しく作成されることを確認します。
func TestNewPlayer_プレイヤー作成(t *testing.T) {
	player := NewPlayer()

	// 初期状態ではHP/MaxHPは0（エージェント装備後に計算）
	if player.HP != 0 {
		t.Errorf("初期HPが期待値と異なります: got %d, want 0", player.HP)
	}
	if player.MaxHP != 0 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 0", player.MaxHP)
	}
	if player.EffectTable == nil {
		t.Error("EffectTableがnilです")
	}
}

// TestPlayerModel_最大HP計算 は装備エージェントのコアレベル平均からMaxHPを計算することを確認します。

func TestPlayerModel_最大HP計算(t *testing.T) {
	tests := []struct {
		name          string
		agentLevels   []int
		expectedMaxHP int
	}{
		{
			name:          "レベル10のエージェント1体",
			agentLevels:   []int{10},
			expectedMaxHP: 200, // 10 × 10.0 + 100
		},
		{
			name:          "レベル10,20,30のエージェント3体",
			agentLevels:   []int{10, 20, 30},
			expectedMaxHP: 300, // (10+20+30)/3 × 10.0 + 100 = 20 × 10.0 + 100
		},
		{
			name:          "レベル1のエージェント1体",
			agentLevels:   []int{1},
			expectedMaxHP: 110, // 1 × 10.0 + 100
		},
		{
			name:          "レベル100のエージェント3体",
			agentLevels:   []int{100, 100, 100},
			expectedMaxHP: 1100, // 100 × 10.0 + 100
		},
		{
			name:          "レベル5,10のエージェント2体",
			agentLevels:   []int{5, 10},
			expectedMaxHP: 175, // (5+10)/2 × 10.0 + 100 = 7.5 × 10.0 + 100 = 175
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agents := createTestAgents(tt.agentLevels)
			maxHP := CalculateMaxHP(agents)

			if maxHP != tt.expectedMaxHP {
				t.Errorf("MaxHPが期待値と異なります: got %d, want %d", maxHP, tt.expectedMaxHP)
			}
		})
	}
}

// TestPlayerModel_エージェント未装備時のHP は装備エージェントがいない場合のMaxHP計算を確認します。
func TestPlayerModel_エージェント未装備時のHP(t *testing.T) {
	agents := []*AgentModel{}
	maxHP := CalculateMaxHP(agents)

	// エージェント未装備時は基礎HP(100)を返す
	if maxHP != BaseHP {
		t.Errorf("エージェント未装備時のMaxHPは基礎HP(%d)であるべきです: got %d", BaseHP, maxHP)
	}
}

// TestPlayerModel_HP再計算 は装備変更時のHP再計算を確認します。

func TestPlayerModel_HP再計算(t *testing.T) {
	player := NewPlayer()

	// 初期状態
	if player.MaxHP != 0 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 0", player.MaxHP)
	}

	// エージェントを装備（レベル10）
	agents1 := createTestAgents([]int{10})
	player.RecalculateHP(agents1)

	// 10 × 10.0 + 100 = 200
	if player.MaxHP != 200 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 200", player.MaxHP)
	}
	if player.HP != 200 {
		t.Errorf("HPも最大値に設定されるべき: got %d, want 200", player.HP)
	}

	// エージェントを追加装備（レベル10,20）
	agents2 := createTestAgents([]int{10, 20})
	player.RecalculateHP(agents2)

	// (10+20)/2 × 10.0 + 100 = 15 × 10.0 + 100 = 250
	if player.MaxHP != 250 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 250", player.MaxHP)
	}
	// HPは新しいMaxHPで初期化される
	if player.HP != 250 {
		t.Errorf("HPが期待値と異なります: got %d, want 250", player.HP)
	}
}

// TestPlayerModel_バトル開始時全回復 はバトル開始時にHPが全回復することを確認します。

func TestPlayerModel_バトル開始時全回復(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{10})
	player.RecalculateHP(agents)

	// ダメージを受けた状態にする
	player.HP = 50

	// バトル開始時の処理
	player.FullHeal()

	if player.HP != player.MaxHP {
		t.Errorf("HPが全回復していません: got %d, want %d", player.HP, player.MaxHP)
	}
}

// TestPlayerModel_HP増減 はHPの増減処理を確認します。
func TestPlayerModel_HP増減(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{10})
	player.RecalculateHP(agents)

	// MaxHP = 10 × 10 + 100 = 200

	// ダメージを受ける
	player.TakeDamage(30)
	if player.HP != 170 {
		t.Errorf("HP減少後の値が期待値と異なります: got %d, want 170", player.HP)
	}

	// 回復
	player.Heal(20)
	if player.HP != 190 {
		t.Errorf("HP回復後の値が期待値と異なります: got %d, want 190", player.HP)
	}

	// 過剰回復（MaxHPを超えない）
	player.Heal(100)
	if player.HP != 200 {
		t.Errorf("HPがMaxHPを超えています: got %d, want 200", player.HP)
	}

	// 致死ダメージ（HPは0以下にならない）
	player.TakeDamage(300)
	if player.HP != 0 {
		t.Errorf("HPが0未満になっています: got %d, want 0", player.HP)
	}
}

// TestPlayerModel_生存確認 はプレイヤーの生存確認を確認します。
func TestPlayerModel_生存確認(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{10})
	player.RecalculateHP(agents)

	// 生存状態
	if !player.IsAlive() {
		t.Error("HPが0より大きい場合は生存しているはずです")
	}

	// 死亡状態
	player.HP = 0
	if player.IsAlive() {
		t.Error("HP=0の場合は死亡しているはずです")
	}
}

// TestPlayerModel_バトル持ち越しなし はHPがバトル間で持ち越されないことを確認します。

func TestPlayerModel_バトル持ち越しなし(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{10})
	player.RecalculateHP(agents)

	// 前のバトルでダメージを受けた
	player.HP = 30

	// 新しいバトル開始
	player.PrepareForBattle()

	// HPは全回復しているはず
	if player.HP != player.MaxHP {
		t.Errorf("バトル開始時にHPが全回復していません: got %d, want %d", player.HP, player.MaxHP)
	}
}

// TestHPConstants はHP計算定数が正しい値であることを確認します。
func TestHPConstants(t *testing.T) {
	// HP係数はゲームバランス調整用の定数
	if HPCoefficient != 10.0 {
		t.Errorf("HPCoefficientが期待値と異なります: got %f, want 10.0", HPCoefficient)
	}
	// 基礎HPはゲームバランス調整用の定数
	if BaseHP != 100 {
		t.Errorf("BaseHPが期待値と異なります: got %d, want 100", BaseHP)
	}
}

// createTestAgents はテスト用のエージェントを作成するヘルパー関数です。
func createTestAgents(levels []int) []*AgentModel {
	agents := make([]*AgentModel, len(levels))

	coreType := CoreType{
		ID:          "test",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}

	modules := make([]*ModuleModel, 4)
	for i := 0; i < 4; i++ {
		modules[i] = NewModule("mod", "テスト", PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", "テスト")
	}

	for i, level := range levels {
		core := NewCore("core_test", "テストコア", level, coreType, passiveSkill)
		agents[i] = NewAgent("agent_test", core, modules)
	}

	return agents
}
