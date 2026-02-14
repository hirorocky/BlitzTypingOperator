// Package config はアプリケーション設定を管理します。

package config

// ==================================================
// ゲームバランス設定
// ==================================================

// BalanceConfig はゲームバランスの設定を保持する構造体です。
type BalanceConfig struct {
	// HPCoefficient はプレイヤーHP計算の係数です。

	HPCoefficient float64

	// EnemyAttackPowerScale は敵攻撃力のレベルスケーリング係数です。

	EnemyAttackPowerScale float64

	// EnemyAttackIntervalScale は敵攻撃間隔のレベルスケーリング係数です。

	EnemyAttackIntervalScale float64

	// MinAttackIntervalMS は敵攻撃間隔の最小値（ミリ秒）です。

	MinAttackIntervalMS int

	// CoreDropRate はコアのドロップ確率です。

	CoreDropRate float64

	// SkillDropRate はスキルのドロップ確率です。
	SkillDropRate float64

	// MaxLevel は敵の最大レベルです。

	MaxLevel int

	// MaxEquippedAgents は装備可能なエージェント数です。
	MaxEquippedAgents int

	// SkillsPerAgent はエージェントあたりのスキル数です。
	SkillsPerAgent int

	// TextLengthByDifficulty は難易度ごとのテキスト長さ範囲[min, max]です。

	TextLengthByDifficulty map[int][2]int

	// TimeLimitByDifficulty は難易度ごとの制限時間（ミリ秒）です。
	TimeLimitByDifficulty map[int]int
}

// DefaultBalanceConfig はデフォルトのゲームバランス設定を返します。
func DefaultBalanceConfig() *BalanceConfig {
	return &BalanceConfig{
		HPCoefficient:            20.0, // 適度なHP
		EnemyAttackPowerScale:    1.1,  // レベルごとに10%増加
		EnemyAttackIntervalScale: 0.95, // レベルごとに5%短縮
		MinAttackIntervalMS:      1000, // 最小1秒
		CoreDropRate:             0.5,  // 50%
		SkillDropRate:            0.7,  // 70%
		MaxLevel:                 100,
		MaxEquippedAgents:        3,
		SkillsPerAgent:           4,
		TextLengthByDifficulty: map[int][2]int{
			1: {3, 6},   // Easy
			2: {7, 11},  // Medium
			3: {12, 20}, // Hard
		},
		TimeLimitByDifficulty: map[int]int{
			1: 10000, // Easy: 10秒
			2: 7000,  // Medium: 7秒
			3: 5000,  // Hard: 5秒
		},
	}
}

// BalanceConfigOption は設定のカスタマイズオプションです。
type BalanceConfigOption func(*BalanceConfig)

// NewBalanceConfig はカスタム設定を持つBalanceConfigを作成します。
func NewBalanceConfig(opts ...BalanceConfigOption) *BalanceConfig {
	config := DefaultBalanceConfig()
	for _, opt := range opts {
		opt(config)
	}
	return config
}

// WithHPCoefficient はHP係数を設定するオプションです。
func WithHPCoefficient(coeff float64) BalanceConfigOption {
	return func(c *BalanceConfig) {
		c.HPCoefficient = coeff
	}
}

// WithCoreDropRate はコアドロップ率を設定するオプションです。
func WithCoreDropRate(rate float64) BalanceConfigOption {
	return func(c *BalanceConfig) {
		c.CoreDropRate = rate
	}
}

// WithSkillDropRate はスキルドロップ率を設定するオプションです。
func WithSkillDropRate(rate float64) BalanceConfigOption {
	return func(c *BalanceConfig) {
		c.SkillDropRate = rate
	}
}

// ==================================================
// 計算メソッド
// ==================================================

// CalculateEnemyAttackPower は敵の攻撃力を計算します。

func (c *BalanceConfig) CalculateEnemyAttackPower(baseAttackPower int, level int) int {
	// 基礎攻撃力 × (スケーリング係数 ^ (レベル - 1))
	scale := 1.0
	for i := 1; i < level; i++ {
		scale *= c.EnemyAttackPowerScale
	}
	return int(float64(baseAttackPower) * scale)
}

// CalculateEnemyAttackInterval は敵の攻撃間隔を計算します。

func (c *BalanceConfig) CalculateEnemyAttackInterval(baseIntervalMS int, level int) int {
	// 基礎間隔 × (スケーリング係数 ^ (レベル - 1))
	scale := 1.0
	for i := 1; i < level; i++ {
		scale *= c.EnemyAttackIntervalScale
	}
	interval := int(float64(baseIntervalMS) * scale)

	// 最小間隔を保証
	if interval < c.MinAttackIntervalMS {
		return c.MinAttackIntervalMS
	}
	return interval
}

// GetTextLengthRange は指定難易度のテキスト長さ範囲を返します。

func (c *BalanceConfig) GetTextLengthRange(difficulty int) (min, max int) {
	if range_, exists := c.TextLengthByDifficulty[difficulty]; exists {
		return range_[0], range_[1]
	}
	// デフォルト（Easy）
	return 3, 6
}

// GetTimeLimit は指定難易度の制限時間を返します（ミリ秒）。
func (c *BalanceConfig) GetTimeLimit(difficulty int) int {
	if limit, exists := c.TimeLimitByDifficulty[difficulty]; exists {
		return limit
	}
	// デフォルト（Easy）
	return 10000
}

// ==================================================
// DifficultyRate ベースの連続関数（50-200スケール）
// ==================================================

// GetTextLengthForRate はDifficultyRate(50-200)に基づくテキスト長さ範囲を返します。
// 線形補間: rate=50→3-6文字、rate=100→7-11文字、rate=200→12-20文字
func GetTextLengthForRate(rate int) (min, max int) {
	if rate <= 50 {
		return 3, 6
	}
	if rate >= 200 {
		return 12, 20
	}

	// 50-100: 線形補間（3-6 → 7-11）
	if rate <= 100 {
		t := float64(rate-50) / 50.0
		minLen := 3.0 + t*(7.0-3.0)
		maxLen := 6.0 + t*(11.0-6.0)
		return int(minLen), int(maxLen)
	}

	// 100-200: 線形補間（7-11 → 12-20）
	t := float64(rate-100) / 100.0
	minLen := 7.0 + t*(12.0-7.0)
	maxLen := 11.0 + t*(20.0-11.0)
	return int(minLen), int(maxLen)
}

// GetTimeLimitForRate はDifficultyRate(50-200)に基づく制限時間をミリ秒で返します。
// 線形補間: rate=50→12000ms、rate=100→8000ms、rate=200→4000ms
func GetTimeLimitForRate(rate int) int {
	if rate <= 50 {
		return 12000
	}
	if rate >= 200 {
		return 4000
	}

	// 50-200: 線形補間（12000ms → 4000ms）
	t := float64(rate-50) / 150.0
	limit := 12000.0 - t*(12000.0-4000.0)
	return int(limit)
}
