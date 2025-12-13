// Package enemy は敵生成システムを提供します。
// レベルに応じた敵の生成、ステータス計算を担当します。
package spawning

import (
	"fmt"
	"math/rand"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"

	"github.com/google/uuid"
)

// 敵生成関連の定数
const (
	// MaxEnemyLevel は敵の最大レベルです。
	MaxEnemyLevel = 100

	// MinEnemyLevel は敵の最小レベルです。
	MinEnemyLevel = 1

	// MinAttackInterval は敵の最低攻撃間隔です。
	MinAttackInterval = config.MinEnemyAttackInterval

	// AttackPowerPerLevel はレベルあたりの攻撃力上昇値です。
	AttackPowerPerLevel = 2

	// AttackIntervalReductionPerLevel はレベルあたりの攻撃間隔短縮（ミリ秒）です。
	AttackIntervalReductionPerLevel = 50
)

// EnemyGenerator はドメイン型を使用した敵生成を担当する構造体です。
type EnemyGenerator struct {
	// enemyTypes は敵タイプ定義リストです（ドメイン型）。
	enemyTypes []domain.EnemyType

	// rng は乱数生成器です。
	rng *rand.Rand
}

// NewEnemyGenerator はドメイン型を使用する新しいEnemyGeneratorを作成します。
func NewEnemyGenerator(enemyTypes []domain.EnemyType) *EnemyGenerator {
	return &EnemyGenerator{
		enemyTypes: enemyTypes,
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Generate は指定レベルの敵をドメイン型から生成します。
func (g *EnemyGenerator) Generate(level int) *domain.EnemyModel {
	// レベルをクランプ
	level = g.clampLevel(level)

	// 敵タイプがない場合はデフォルトの敵を生成
	if len(g.enemyTypes) == 0 {
		return g.generateDefaultEnemy(level)
	}

	// 敵タイプからランダムに選択
	selectedType := g.enemyTypes[g.rng.Intn(len(g.enemyTypes))]

	// レベルに応じたステータス計算
	hp := g.calculateHP(selectedType.BaseHP, level)
	attackPower := g.calculateAttackPower(selectedType.BaseAttackPower, level)
	attackInterval := g.calculateAttackInterval(selectedType.BaseAttackInterval, level)

	// 敵モデルを作成
	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("%s Lv.%d", selectedType.Name, level),
		level,
		hp,
		attackPower,
		attackInterval,
		selectedType,
	)
}

// GenerateWithType は指定された敵タイプで敵を生成します。
func (g *EnemyGenerator) GenerateWithType(level int, typeID string) *domain.EnemyModel {
	level = g.clampLevel(level)

	// 指定されたタイプを検索
	var selectedType *domain.EnemyType
	for i := range g.enemyTypes {
		if g.enemyTypes[i].ID == typeID {
			selectedType = &g.enemyTypes[i]
			break
		}
	}

	// タイプが見つからない場合はランダムに生成
	if selectedType == nil {
		return g.Generate(level)
	}

	hp := g.calculateHP(selectedType.BaseHP, level)
	attackPower := g.calculateAttackPower(selectedType.BaseAttackPower, level)
	attackInterval := g.calculateAttackInterval(selectedType.BaseAttackInterval, level)

	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("%s Lv.%d", selectedType.Name, level),
		level,
		hp,
		attackPower,
		attackInterval,
		*selectedType,
	)
}

// GetEnemyTypes は全ての敵タイプをドメイン型で返します。
func (g *EnemyGenerator) GetEnemyTypes() []domain.EnemyType {
	return g.enemyTypes
}

// SetSeed は乱数シードを設定します（テスト用）。
func (g *EnemyGenerator) SetSeed(seed int64) {
	g.rng = rand.New(rand.NewSource(seed))
}

// generateDefaultEnemy はデフォルトの敵を生成します。
func (g *EnemyGenerator) generateDefaultEnemy(level int) *domain.EnemyModel {
	defaultType := domain.EnemyType{
		ID:                 "default",
		Name:               "敵",
		BaseHP:             50,
		BaseAttackPower:    5,
		BaseAttackInterval: 3000 * time.Millisecond,
		AttackType:         "physical",
	}

	hp := g.calculateHP(defaultType.BaseHP, level)
	attackPower := g.calculateAttackPower(defaultType.BaseAttackPower, level)
	attackInterval := g.calculateAttackInterval(defaultType.BaseAttackInterval, level)

	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("敵 Lv.%d", level),
		level,
		hp,
		attackPower,
		attackInterval,
		defaultType,
	)
}

// calculateHP はレベルに応じたHPを計算します。
func (g *EnemyGenerator) calculateHP(baseHP int, level int) int {
	return baseHP * level
}

// calculateAttackPower はレベルに応じた攻撃力を計算します。
func (g *EnemyGenerator) calculateAttackPower(baseAttackPower int, level int) int {
	return baseAttackPower + (level * AttackPowerPerLevel)
}

// calculateAttackInterval はレベルに応じた攻撃間隔を計算します。
func (g *EnemyGenerator) calculateAttackInterval(baseInterval time.Duration, level int) time.Duration {
	reduction := time.Duration(level*AttackIntervalReductionPerLevel) * time.Millisecond
	interval := baseInterval - reduction

	if interval < MinAttackInterval {
		interval = MinAttackInterval
	}

	return interval
}

// clampLevel はレベルを有効範囲内にクランプします。
func (g *EnemyGenerator) clampLevel(level int) int {
	if level < MinEnemyLevel {
		return MinEnemyLevel
	}
	if level > MaxEnemyLevel {
		return MaxEnemyLevel
	}
	return level
}

// IsMaxLevelBattle は最高レベルのバトルかどうかを判定します。
func (g *EnemyGenerator) IsMaxLevelBattle(level int) bool {
	return level >= MaxEnemyLevel
}

// GetMaxLevel は最大レベルを返します。
func (g *EnemyGenerator) GetMaxLevel() int {
	return MaxEnemyLevel
}

// GetMinLevel は最小レベルを返します。
func (g *EnemyGenerator) GetMinLevel() int {
	return MinEnemyLevel
}

// GetEnemyTypeCount は敵タイプの数を返します。
func (g *EnemyGenerator) GetEnemyTypeCount() int {
	return len(g.enemyTypes)
}
