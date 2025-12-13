// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"time"
)

// EnhanceThreshold は敵が強化フェーズに移行するHP割合の閾値です（50%）。
const EnhanceThreshold = 0.5

// EnemyPhase は敵のフェーズを表す型です。
type EnemyPhase int

const (
	// PhaseNormal は通常フェーズです（HP50%以上）
	PhaseNormal EnemyPhase = 0

	// PhaseEnhanced は強化フェーズです（HP50%以下、特殊攻撃解禁）
	PhaseEnhanced EnemyPhase = 1
)

// String はEnemyPhaseの日本語表示名を返します。
func (p EnemyPhase) String() string {
	switch p {
	case PhaseNormal:
		return "通常"
	case PhaseEnhanced:
		return "強化"
	default:
		return "不明"
	}
}

// EnemyType は敵の種類（タイプ）を定義する構造体です。
// 外部データファイル（enemies.json）から読み込まれます。
type EnemyType struct {
	// ID は敵タイプの一意識別子です。
	ID string

	// Name は敵タイプの表示名です（日本語）。
	Name string

	// BaseHP は敵の基礎HP値です。
	BaseHP int

	// BaseAttackPower は敵の基礎攻撃力です。
	BaseAttackPower int

	// BaseAttackInterval は敵の基礎攻撃間隔です。
	BaseAttackInterval time.Duration

	// AttackType は攻撃属性（physical / magic）です。
	AttackType string

	// ASCIIArt は敵の外観（ASCIIアート）です。
	ASCIIArt string
}

// EnemyModel はゲーム内の敵エンティティを表す構造体です。
type EnemyModel struct {
	// ID は敵インスタンスの一意識別子です。
	ID string

	// Name は敵の表示名です。
	Name string

	// Level は敵のレベルです。
	Level int

	// HP は敵の現在HP値です。
	HP int

	// MaxHP は敵の最大HP値です。
	MaxHP int

	// AttackPower は敵の攻撃力です。
	AttackPower int

	// AttackInterval は敵の攻撃間隔です。
	AttackInterval time.Duration

	// Type は敵の種類（タイプ）です。
	Type EnemyType

	// Phase は敵の現在フェーズです。
	Phase EnemyPhase

	// EffectTable は敵に適用されているステータス効果テーブルです。
	// 敵自身のバフやプレイヤーからのデバフを管理します。
	EffectTable *EffectTable
}

// NewEnemy は新しいEnemyModelを作成します。
// 初期状態は通常フェーズ（PhaseNormal）です。
func NewEnemy(id, name string, level, hp, attackPower int, attackInterval time.Duration, enemyType EnemyType) *EnemyModel {
	return &EnemyModel{
		ID:             id,
		Name:           name,
		Level:          level,
		HP:             hp,
		MaxHP:          hp,
		AttackPower:    attackPower,
		AttackInterval: attackInterval,
		Type:           enemyType,
		Phase:          PhaseNormal,
		EffectTable:    NewEffectTable(),
	}
}

// TakeDamage はダメージを受けてHPを減少させます。
// HPは0未満にはなりません。
func (e *EnemyModel) TakeDamage(damage int) {
	e.HP -= damage
	if e.HP < 0 {
		e.HP = 0
	}
}

// IsAlive は敵が生存しているかどうかを返します。
// HP > 0 の場合に生存とみなします。
func (e *EnemyModel) IsAlive() bool {
	return e.HP > 0
}

// GetHPPercentage はHPの残り割合を0.0〜1.0で返します。
func (e *EnemyModel) GetHPPercentage() float64 {
	if e.MaxHP == 0 {
		return 0.0
	}
	return float64(e.HP) / float64(e.MaxHP)
}

// ShouldTransitionToEnhanced は強化フェーズに移行すべきかどうかを判定します。
// 既にPhaseEnhancedの場合はfalseを返します。
func (e *EnemyModel) ShouldTransitionToEnhanced() bool {
	if e.Phase == PhaseEnhanced {
		return false
	}
	return e.GetHPPercentage() <= EnhanceThreshold
}

// TransitionToEnhanced は強化フェーズに移行します。
func (e *EnemyModel) TransitionToEnhanced() {
	e.Phase = PhaseEnhanced
}

// CheckAndTransitionPhase はHPをチェックし、必要に応じてフェーズ移行を実行します。
// フェーズ移行した場合はtrueを返します。
func (e *EnemyModel) CheckAndTransitionPhase() bool {
	if e.ShouldTransitionToEnhanced() {
		e.TransitionToEnhanced()
		return true
	}
	return false
}

// IsEnhanced は現在強化フェーズかどうかを返します。
func (e *EnemyModel) IsEnhanced() bool {
	return e.Phase == PhaseEnhanced
}

// GetPhaseString は現在のフェーズの表示文字列を返します。
func (e *EnemyModel) GetPhaseString() string {
	return e.Phase.String()
}
