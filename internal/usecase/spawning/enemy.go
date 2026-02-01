// Package enemy は敵生成システムを提供します。
// レベルに応じた敵の生成、ステータス計算を担当します。
package spawning

import (
	"fmt"
	"math/rand"
	"time"

	"hirorocky/type-battle/internal/domain"

	"github.com/google/uuid"
)

// 敵生成関連の定数
const (
	// MaxEnemyLevel は敵の最大レベルです。
	MaxEnemyLevel = 100

	// MinEnemyLevel は敵の最小レベルです。
	MinEnemyLevel = 1

	// AttackPowerPerLevel はレベルあたりの攻撃力上昇値です。
	AttackPowerPerLevel = 2

	// DefaultBaseAttackPower はデフォルトの基礎攻撃力です。
	DefaultBaseAttackPower = 5
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
	attackPower := g.calculateAttackPowerFromActions(selectedType, level)

	// 敵モデルを作成
	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("%s Lv.%d", selectedType.Name, level),
		level,
		hp,
		attackPower,
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
	attackPower := g.calculateAttackPowerFromActions(*selectedType, level)

	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("%s Lv.%d", selectedType.Name, level),
		level,
		hp,
		attackPower,
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
		ID:     "default",
		Name:   "敵",
		BaseHP: 50,
	}

	hp := g.calculateHP(defaultType.BaseHP, level)
	attackPower := g.calculateAttackPower(DefaultBaseAttackPower, level)

	return domain.NewEnemy(
		uuid.New().String(),
		fmt.Sprintf("敵 Lv.%d", level),
		level,
		hp,
		attackPower,
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

// calculateAttackPowerFromActions は行動パターンから攻撃力を計算します。
// 最初の攻撃行動のDamageBase + Level * DamagePerLevelを返します。
func (g *EnemyGenerator) calculateAttackPowerFromActions(enemyType domain.EnemyType, level int) int {
	// 解決済み行動パターンから最初の攻撃行動を取得
	for _, action := range enemyType.ResolvedNormalActions {
		if action.ActionType == domain.EnemyActionAttack {
			damage := action.DamageBase + float64(level)*action.DamagePerLevel
			return int(damage)
		}
	}
	// 攻撃行動がない場合はデフォルト値
	return g.calculateAttackPower(DefaultBaseAttackPower, level)
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
