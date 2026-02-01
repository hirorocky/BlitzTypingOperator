// Package service はドメインサービスを提供します。
// 複数のドメインオブジェクトを組み合わせる純粋なビジネスロジックを配置します。
package service

import "hirorocky/type-battle/internal/domain"

// CalculateStats はコア特性からステータス値を計算します。
// 各ステータス: 基礎値(100) × ステータス重み
// 結果は整数に切り捨てられます。
func CalculateStats(coreType domain.CoreType) domain.Stats {
	// 各ステータスの重みを取得（未設定の場合は0.0扱い）
	strWeight := coreType.StatWeights["STR"]
	intWeight := coreType.StatWeights["INT"]
	wilWeight := coreType.StatWeights["WIL"]
	lukWeight := coreType.StatWeights["LUK"]

	// 計算式: 基礎値(100) × 重み
	baseValue := float64(domain.BaseStatValue)

	return domain.Stats{
		STR: int(baseValue * strWeight),
		INT: int(baseValue * intWeight),
		WIL: int(baseValue * wilWeight),
		LUK: int(baseValue * lukWeight),
	}
}
