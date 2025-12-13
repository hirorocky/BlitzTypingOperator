// Package service はドメインサービスを提供します。
// 複数のドメインオブジェクトを組み合わせる純粋なビジネスロジックを配置します。
package service

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestCalculateStats_Basic はステータス計算の基本動作をテストします。
func TestCalculateStats_Basic(t *testing.T) {
	// コア特性を用意
	coreType := domain.CoreType{
		ID:   "test-type",
		Name: "テスト特性",
		StatWeights: map[string]float64{
			"STR": 1.0,
			"MAG": 1.0,
			"SPD": 1.0,
			"LUK": 1.0,
		},
	}

	// レベル1でテスト
	stats := CalculateStats(1, coreType)

	// 基礎値(10) × レベル(1) × 重み(1.0) = 10
	if stats.STR != 10 {
		t.Errorf("STR expected 10, got %d", stats.STR)
	}
	if stats.MAG != 10 {
		t.Errorf("MAG expected 10, got %d", stats.MAG)
	}
	if stats.SPD != 10 {
		t.Errorf("SPD expected 10, got %d", stats.SPD)
	}
	if stats.LUK != 10 {
		t.Errorf("LUK expected 10, got %d", stats.LUK)
	}
}

// TestCalculateStats_WithWeights は重み付きステータス計算をテストします。
func TestCalculateStats_WithWeights(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "weighted-type",
		Name: "重み付き特性",
		StatWeights: map[string]float64{
			"STR": 1.5,  // 50%増加
			"MAG": 0.5,  // 50%減少
			"SPD": 2.0,  // 2倍
			"LUK": 0.25, // 25%
		},
	}

	// レベル2でテスト
	stats := CalculateStats(2, coreType)

	// 基礎値(10) × レベル(2) × 重み
	expected := map[string]int{
		"STR": 30, // 20 × 1.5 = 30
		"MAG": 10, // 20 × 0.5 = 10
		"SPD": 40, // 20 × 2.0 = 40
		"LUK": 5,  // 20 × 0.25 = 5
	}

	if stats.STR != expected["STR"] {
		t.Errorf("STR expected %d, got %d", expected["STR"], stats.STR)
	}
	if stats.MAG != expected["MAG"] {
		t.Errorf("MAG expected %d, got %d", expected["MAG"], stats.MAG)
	}
	if stats.SPD != expected["SPD"] {
		t.Errorf("SPD expected %d, got %d", expected["SPD"], stats.SPD)
	}
	if stats.LUK != expected["LUK"] {
		t.Errorf("LUK expected %d, got %d", expected["LUK"], stats.LUK)
	}
}

// TestCalculateStats_HighLevel は高レベルでの計算をテストします。
func TestCalculateStats_HighLevel(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "high-level-type",
		Name: "高レベルテスト",
		StatWeights: map[string]float64{
			"STR": 1.2,
			"MAG": 0.8,
			"SPD": 1.0,
			"LUK": 1.0,
		},
	}

	// レベル10でテスト
	stats := CalculateStats(10, coreType)

	// 基礎値(10) × レベル(10) = 100
	// STR: 100 × 1.2 = 120
	// MAG: 100 × 0.8 = 80
	if stats.STR != 120 {
		t.Errorf("STR expected 120, got %d", stats.STR)
	}
	if stats.MAG != 80 {
		t.Errorf("MAG expected 80, got %d", stats.MAG)
	}
	if stats.SPD != 100 {
		t.Errorf("SPD expected 100, got %d", stats.SPD)
	}
	if stats.LUK != 100 {
		t.Errorf("LUK expected 100, got %d", stats.LUK)
	}
}

// TestCalculateStats_ZeroWeight はゼロ重みをテストします。
func TestCalculateStats_ZeroWeight(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "zero-weight-type",
		Name: "ゼロ重み",
		StatWeights: map[string]float64{
			"STR": 0.0, // ゼロ
			"MAG": 1.0,
			"SPD": 0.0,
			"LUK": 1.0,
		},
	}

	stats := CalculateStats(5, coreType)

	// ゼロ重みは0になる
	if stats.STR != 0 {
		t.Errorf("STR expected 0, got %d", stats.STR)
	}
	if stats.SPD != 0 {
		t.Errorf("SPD expected 0, got %d", stats.SPD)
	}
	// 通常重みは計算される
	if stats.MAG != 50 {
		t.Errorf("MAG expected 50, got %d", stats.MAG)
	}
	if stats.LUK != 50 {
		t.Errorf("LUK expected 50, got %d", stats.LUK)
	}
}
