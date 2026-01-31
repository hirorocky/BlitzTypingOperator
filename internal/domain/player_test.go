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

// TestPlayerModel_最大HP計算 はエージェント装備時のMaxHPを確認します。
// 新仕様: 初期最大HPは1000固定、成長は敵撃破で増加
func TestPlayerModel_最大HP計算(t *testing.T) {
	tests := []struct {
		name          string
		agentCount    int
		expectedMaxHP int
	}{
		{
			name:          "エージェント1体装備",
			agentCount:    1,
			expectedMaxHP: 1000, // 初期最大HP
		},
		{
			name:          "エージェント3体装備",
			agentCount:    3,
			expectedMaxHP: 1000, // 初期最大HP
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// エージェント数分のダミーレベルリストを作成
			levels := make([]int, tt.agentCount)
			for i := range levels {
				levels[i] = 1 // レベルはもはや使用されない
			}
			agents := createTestAgents(levels)
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
// 新仕様: MaxHPは初期値1000固定

func TestPlayerModel_HP再計算(t *testing.T) {
	player := NewPlayer()

	// 初期状態
	if player.MaxHP != 0 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 0", player.MaxHP)
	}

	// エージェントを装備
	agents1 := createTestAgents([]int{1})
	player.RecalculateHP(agents1)

	// 新仕様: 初期最大HP = 1000
	if player.MaxHP != 1000 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 1000", player.MaxHP)
	}
	if player.HP != 1000 {
		t.Errorf("HPも最大値に設定されるべき: got %d, want 1000", player.HP)
	}

	// エージェントを追加装備
	agents2 := createTestAgents([]int{1, 1})
	player.RecalculateHP(agents2)

	// 新仕様: MaxHPはエージェント数に関わらず初期最大HP = 1000
	if player.MaxHP != 1000 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 1000", player.MaxHP)
	}
	// HPは新しいMaxHPで初期化される
	if player.HP != 1000 {
		t.Errorf("HPが期待値と異なります: got %d, want 1000", player.HP)
	}
}

// TestPlayerModel_バトル開始時全回復 はバトル開始時にHPが全回復することを確認します。

func TestPlayerModel_バトル開始時全回復(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{1})
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
	agents := createTestAgents([]int{1})
	player.RecalculateHP(agents)

	// 新仕様: MaxHP = 1000

	// ダメージを受ける
	player.TakeDamage(30)
	if player.HP != 970 {
		t.Errorf("HP減少後の値が期待値と異なります: got %d, want 970", player.HP)
	}

	// 回復
	player.Heal(20)
	if player.HP != 990 {
		t.Errorf("HP回復後の値が期待値と異なります: got %d, want 990", player.HP)
	}

	// 過剰回復（MaxHPを超えない）
	player.Heal(100)
	if player.HP != 1000 {
		t.Errorf("HPがMaxHPを超えています: got %d, want 1000", player.HP)
	}

	// 致死ダメージ（HPは0以下にならない）
	player.TakeDamage(1500)
	if player.HP != 0 {
		t.Errorf("HPが0未満になっています: got %d, want 0", player.HP)
	}
}

// TestPlayerModel_生存確認 はプレイヤーの生存確認を確認します。
func TestPlayerModel_生存確認(t *testing.T) {
	player := NewPlayer()
	agents := createTestAgents([]int{1})
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
	agents := createTestAgents([]int{1})
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
	// 旧仕様の定数（互換性用）
	if HPCoefficient != 10.0 {
		t.Errorf("HPCoefficientが期待値と異なります: got %f, want 10.0", HPCoefficient)
	}
	if BaseHP != 100 {
		t.Errorf("BaseHPが期待値と異なります: got %d, want 100", BaseHP)
	}
	// 新仕様の初期最大HP
	if InitialMaxHP != 1000 {
		t.Errorf("InitialMaxHPが期待値と異なります: got %d, want 1000", InitialMaxHP)
	}
}

// createTestAgents はテスト用のエージェントを作成するヘルパー関数です。
// 注意: 旧実装ではエージェントレベルがコアレベルから導出されていましたが、
// 新実装ではエージェントにはレベル概念がありません。
// このテストはレガシーコードとの互換性のために残されています。
func createTestAgents(levels []int) []*AgentModel {
	agents := make([]*AgentModel, len(levels))

	coreType := CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := PassiveSkill{ID: "test_skill"}

	modules := make([]*ModuleModel, 4)
	for i := 0; i < 4; i++ {
		modules[i] = NewModuleFromType(ModuleType{
			ID:          "mod",
			Name:        "テスト",
			Icon:        "⚔️",
			Tags:        []string{"physical_low"},
			Description: "テスト",
			Effects: []ModuleEffect{
				{
					Target:      TargetEnemy,
					HPFormula:   &HPFormula{Base: 0, StatCoef: 1.0, StatRef: "STR"},
					Probability: 1.0,
				},
			},
		}, nil)
	}

	for i := range levels {
		core := NewCoreWithTypeID("test", coreType, passiveSkill)
		agents[i] = NewAgent("agent_test", core, modules)
	}

	return agents
}
