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

	// スキルドロップ率は0〜1の範囲
	if config.SkillDropRate < 0 || config.SkillDropRate > 1 {
		t.Errorf("スキルドロップ率は0〜1の範囲であるべきです: got %f", config.SkillDropRate)
	}

	// スキルドロップ率 >= コアドロップ率（スキルの方がドロップしやすい）
	if config.SkillDropRate < config.CoreDropRate {
		t.Error("スキルドロップ率はコアドロップ率以上であるべきです")
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

func TestSkillsPerAgent(t *testing.T) {
	// エージェントあたりのスキル数
	config := DefaultBalanceConfig()

	if config.SkillsPerAgent != 4 {
		t.Errorf("エージェントあたりのスキル数は4であるべきです: got %d", config.SkillsPerAgent)
	}
}

// === DifficultyRate ベース連続関数テスト ===

func TestGetTextLengthForRate(t *testing.T) {
	tests := []struct {
		name    string
		rate    int
		wantMin int
		wantMax int
	}{
		{"最低難易度", 50, 3, 6},
		{"標準難易度", 100, 7, 11},
		{"最高難易度", 200, 12, 20},
		{"低難易度(75)", 75, 5, 8},
		{"高難易度(150)", 150, 9, 15},
		{"範囲下限未満", 30, 3, 6},
		{"範囲上限超過", 250, 12, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minLen, maxLen := GetTextLengthForRate(tt.rate)
			if minLen != tt.wantMin {
				t.Errorf("min = %d, want %d", minLen, tt.wantMin)
			}
			if maxLen != tt.wantMax {
				t.Errorf("max = %d, want %d", maxLen, tt.wantMax)
			}
		})
	}
}

func TestGetTextLengthForRate_単調増加(t *testing.T) {
	prevMin, prevMax := GetTextLengthForRate(50)
	for rate := 60; rate <= 200; rate += 10 {
		min, max := GetTextLengthForRate(rate)
		if min < prevMin {
			t.Errorf("rate=%d: min(%d) < prevMin(%d)、単調増加でない", rate, min, prevMin)
		}
		if max < prevMax {
			t.Errorf("rate=%d: max(%d) < prevMax(%d)、単調増加でない", rate, max, prevMax)
		}
		prevMin, prevMax = min, max
	}
}

func TestGetTimeLimitForRate(t *testing.T) {
	tests := []struct {
		name   string
		rate   int
		wantMS int
	}{
		{"最低難易度", 50, 12000},
		{"最高難易度", 200, 4000},
		{"範囲下限未満", 30, 12000},
		{"範囲上限超過", 250, 4000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeLimitForRate(tt.rate)
			if got != tt.wantMS {
				t.Errorf("GetTimeLimitForRate(%d) = %d, want %d", tt.rate, got, tt.wantMS)
			}
		})
	}
}

func TestGetTimeLimitForRate_単調減少(t *testing.T) {
	prevLimit := GetTimeLimitForRate(50)
	for rate := 60; rate <= 200; rate += 10 {
		limit := GetTimeLimitForRate(rate)
		if limit > prevLimit {
			t.Errorf("rate=%d: limit(%d) > prevLimit(%d)、単調減少でない", rate, limit, prevLimit)
		}
		prevLimit = limit
	}
}

func TestBalanceConfigCustomization(t *testing.T) {
	// 設定のカスタマイズが可能であること
	config := NewBalanceConfig(
		WithHPCoefficient(50.0),
		WithCoreDropRate(0.6),
		WithSkillDropRate(0.8),
	)

	if config.HPCoefficient != 50.0 {
		t.Errorf("カスタムHP係数: expected 50.0, got %f", config.HPCoefficient)
	}
	if config.CoreDropRate != 0.6 {
		t.Errorf("カスタムコアドロップ率: expected 0.6, got %f", config.CoreDropRate)
	}
	if config.SkillDropRate != 0.8 {
		t.Errorf("カスタムスキルドロップ率: expected 0.8, got %f", config.SkillDropRate)
	}
}
