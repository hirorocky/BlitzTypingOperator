// Package service はドメインサービスを提供します。
// 複数のドメインオブジェクトを組み合わせる純粋なビジネスロジックを配置します。
package service

import "hirorocky/type-battle/internal/domain"

// CalculateStats はコアレベルとコア特性からステータス値を計算します。
// 計算式: 基礎値(10) × レベル × ステータス重み
// 結果は整数に切り捨てられます。
func CalculateStats(level int, coreType domain.CoreType) domain.Stats {
	// 各ステータスの重みを取得（未設定の場合はデフォルト0.0）
	strWeight := coreType.StatWeights["STR"]
	magWeight := coreType.StatWeights["MAG"]
	spdWeight := coreType.StatWeights["SPD"]
	lukWeight := coreType.StatWeights["LUK"]

	// 計算式: 基礎値 × レベル × 重み
	baseValue := float64(domain.BaseStatValue * level)

	return domain.Stats{
		STR: int(baseValue * strWeight),
		MAG: int(baseValue * magWeight),
		SPD: int(baseValue * spdWeight),
		LUK: int(baseValue * lukWeight),
	}
}
