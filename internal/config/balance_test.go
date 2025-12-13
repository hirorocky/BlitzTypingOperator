package config

import (
	"testing"
)

// ==================================================
// ゲームバランスパラメータのテスト
// ==================================================

func TestHPCoefficient(t *testing.T) {

	config := DefaultBalanceConfig()

	// HP係数は正の値であること
	if config.HPCoefficient <= 0 {
		t.Error("HP係数は正の値であるべきです")
	}

	// HP係数の典型的な範囲（10〜100）
	if config.HPCoefficient < 10 || config.HPCoefficient > 100 {
		t.Errorf("HP係数は適切な範囲であるべきです: got %f", config.HPCoefficient)
	}
}

func TestEnemyAttackPowerScaling(t *testing.T) {

	config := DefaultBalanceConfig()

	// スケーリング係数は正の値
	if config.EnemyAttackPowerScale <= 0 {
		t.Error("敵攻撃力スケーリング係数は正の値であるべきです")
	}

	// レベル1とレベル10での攻撃力計算
	level1Attack := config.CalculateEnemyAttackPower(10, 1)
	level10Attack := config.CalculateEnemyAttackPower(10, 10)

	// 高レベルほど攻撃力が高い
	if level10Attack <= level1Attack {
		t.Error("高レベルの敵は高い攻撃力を持つべきです")
	}
}

func TestEnemyAttackIntervalScaling(t *testing.T) {

	config := DefaultBalanceConfig()

	// レベル1とレベル10での攻撃間隔計算
	level1Interval := config.CalculateEnemyAttackInterval(3000, 1)
	level10Interval := config.CalculateEnemyAttackInterval(3000, 10)

	// 高レベルほど攻撃間隔が短い
	if level10Interval >= level1Interval {
		t.Error("高レベルの敵は短い攻撃間隔を持つべきです")
	}

	// 最小間隔は保証される
	if level10Interval < config.MinAttackIntervalMS {
		t.Error("攻撃間隔は最小値を下回るべきではありません")
	}
}

func TestDropRates(t *testing.T) {

	config := DefaultBalanceConfig()

	// コアドロップ率は0〜1の範囲
	if config.CoreDropRate < 0 || config.CoreDropRate > 1 {
		t.Errorf("コアドロップ率は0〜1の範囲であるべきです: got %f", config.CoreDropRate)
	}

	// モジュールドロップ率は0〜1の範囲
	if config.ModuleDropRate < 0 || config.ModuleDropRate > 1 {
		t.Errorf("モジュールドロップ率は0〜1の範囲であるべきです: got %f", config.ModuleDropRate)
	}

	// モジュールドロップ率 >= コアドロップ率（モジュールの方がドロップしやすい）
	if config.ModuleDropRate < config.CoreDropRate {
		t.Error("モジュールドロップ率はコアドロップ率以上であるべきです")
	}
}

func TestTypingChallengeTextLength(t *testing.T) {

	config := DefaultBalanceConfig()

	// 難易度ごとのテキスト長さ範囲
	tests := []struct {
		name       string
		difficulty int // 1:Easy, 2:Medium, 3:Hard
		wantMinLen int
		wantMaxLen int
	}{
		{"Easy", 1, 3, 6},
		{"Medium", 2, 7, 11},
		{"Hard", 3, 12, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLen, maxLen := config.GetTextLengthRange(tt.difficulty)
			if minLen != tt.wantMinLen {
				t.Errorf("最小テキスト長: expected %d, got %d", tt.wantMinLen, minLen)
			}
			if maxLen != tt.wantMaxLen {
				t.Errorf("最大テキスト長: expected %d, got %d", tt.wantMaxLen, maxLen)
			}
		})
	}
}

func TestTypingChallengeTimeLimit(t *testing.T) {
	// 制限時間の設定
	config := DefaultBalanceConfig()

	// 難易度ごとの制限時間（ミリ秒）
	tests := []struct {
		name       string
		difficulty int
		wantMinMS  int
		wantMaxMS  int
	}{
		{"Easy", 1, 5000, 15000},
		{"Medium", 2, 3000, 10000},
		{"Hard", 3, 2000, 8000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeLimit := config.GetTimeLimit(tt.difficulty)
			if timeLimit < tt.wantMinMS || timeLimit > tt.wantMaxMS {
				t.Errorf("制限時間は%d〜%d msの範囲であるべきです: got %d", tt.wantMinMS, tt.wantMaxMS, timeLimit)
			}
		})
	}
}

func TestMaxLevel(t *testing.T) {

	config := DefaultBalanceConfig()

	if config.MaxLevel != 100 {
		t.Errorf("最大レベルは100であるべきです: got %d", config.MaxLevel)
	}
}

func TestMaxEquippedAgents(t *testing.T) {
	// 最大装備エージェント数
	config := DefaultBalanceConfig()

	if config.MaxEquippedAgents != 3 {
		t.Errorf("最大装備エージェント数は3であるべきです: got %d", config.MaxEquippedAgents)
	}
}

func TestModulesPerAgent(t *testing.T) {
	// エージェントあたりのモジュール数
	config := DefaultBalanceConfig()

	if config.ModulesPerAgent != 4 {
		t.Errorf("エージェントあたりのモジュール数は4であるべきです: got %d", config.ModulesPerAgent)
	}
}

func TestBalanceConfigCustomization(t *testing.T) {
	// 設定のカスタマイズが可能であること
	config := NewBalanceConfig(
		WithHPCoefficient(50.0),
		WithCoreDropRate(0.6),
		WithModuleDropRate(0.8),
	)

	if config.HPCoefficient != 50.0 {
		t.Errorf("カスタムHP係数: expected 50.0, got %f", config.HPCoefficient)
	}
	if config.CoreDropRate != 0.6 {
		t.Errorf("カスタムコアドロップ率: expected 0.6, got %f", config.CoreDropRate)
	}
	if config.ModuleDropRate != 0.8 {
		t.Errorf("カスタムモジュールドロップ率: expected 0.8, got %f", config.ModuleDropRate)
	}
}
