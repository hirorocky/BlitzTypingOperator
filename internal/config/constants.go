// Package config は設定定数を提供します。
// バトル設定、効果持続時間、インベントリ設定のマジックナンバーを一元管理します。
package config

import "time"

// バトル設定定数
const (
	// BattleTickInterval はバトル画面の更新間隔です。
	BattleTickInterval = 100 * time.Millisecond

	// DefaultModuleCooldown はモジュールのデフォルトクールダウン秒数です。
	DefaultModuleCooldown = 5.0

	// AccuracyPenaltyThreshold は正確性ペナルティ発生閾値です。
	// この値未満の正確性の場合、効果が半減します。
	AccuracyPenaltyThreshold = 0.5

	// MinEnemyAttackInterval は敵の最小攻撃間隔です。
	// 高レベルの敵でもこれ以上短い間隔では攻撃しません。
	MinEnemyAttackInterval = 500 * time.Millisecond
)

// 効果持続時間定数
const (
	// BuffDuration はバフのデフォルト持続時間（秒）です。
	BuffDuration = 10.0

	// DebuffDuration はデバフのデフォルト持続時間（秒）です。
	DebuffDuration = 8.0
)

// インベントリ設定定数
const (
	// MaxAgentEquipSlots はエージェント装備スロットの最大数です。
	MaxAgentEquipSlots = 3

	// ModulesPerAgent はエージェントあたりのモジュール数です。
	ModulesPerAgent = 4
)
