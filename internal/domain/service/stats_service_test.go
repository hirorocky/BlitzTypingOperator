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
			"INT": 1.0,
			"WIL": 1.0,
			"LUK": 1.0,
		},
	}

	// 重みベースでテスト
	stats := CalculateStats(coreType)

	// 基礎値(100) × 重み(1.0) = 100
	if stats.STR != 100 {
		t.Errorf("STR expected 100, got %d", stats.STR)
	}
	if stats.INT != 100 {
		t.Errorf("INT expected 100, got %d", stats.INT)
	}
	if stats.WIL != 100 {
		t.Errorf("WIL expected 100, got %d", stats.WIL)
	}
	if stats.LUK != 100 {
		t.Errorf("LUK expected 100, got %d", stats.LUK)
	}
}

// TestCalculateStats_WithWeights は重み付きステータス計算をテストします。
func TestCalculateStats_WithWeights(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "weighted-type",
		Name: "重み付き特性",
		StatWeights: map[string]float64{
			"STR": 1.5, // 50%増加
			"INT": 0.5, // 50%減少
			"WIL": 2.0, // 2倍
			"LUK": 0.5, // 50%
		},
	}

	// 重みベースでテスト
	stats := CalculateStats(coreType)

	// 基礎値(100) × 重み
	expected := map[string]int{
		"STR": 150, // 100 × 1.5 = 150
		"INT": 50,  // 100 × 0.5 = 50
		"WIL": 200, // 100 × 2.0 = 200
		"LUK": 50,  // 100 × 0.5 = 50
	}

	if stats.STR != expected["STR"] {
		t.Errorf("STR expected %d, got %d", expected["STR"], stats.STR)
	}
	if stats.INT != expected["INT"] {
		t.Errorf("INT expected %d, got %d", expected["INT"], stats.INT)
	}
	if stats.WIL != expected["WIL"] {
		t.Errorf("WIL expected %d, got %d", expected["WIL"], stats.WIL)
	}
	if stats.LUK != expected["LUK"] {
		t.Errorf("LUK expected %d, got %d", expected["LUK"], stats.LUK)
	}
}

// TestCalculateStats_TypicalCoreType は典型的なコア特性での計算をテストします。
func TestCalculateStats_TypicalCoreType(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "attack-balance",
		Name: "攻撃バランス",
		StatWeights: map[string]float64{
			"STR": 1.2,
			"INT": 0.8,
			"WIL": 1.0,
			"LUK": 1.0,
		},
	}

	// 重みベースでテスト
	stats := CalculateStats(coreType)

	// 基礎値(100) × 重み
	if stats.STR != 120 {
		t.Errorf("STR expected 120, got %d", stats.STR)
	}
	if stats.INT != 80 {
		t.Errorf("INT expected 80, got %d", stats.INT)
	}
	if stats.WIL != 100 {
		t.Errorf("WIL expected 100, got %d", stats.WIL)
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
			"INT": 1.0,
			"WIL": 0.0,
			"LUK": 1.0,
		},
	}

	stats := CalculateStats(coreType)

	// ゼロ重みは0になる
	if stats.STR != 0 {
		t.Errorf("STR expected 0, got %d", stats.STR)
	}
	if stats.WIL != 0 {
		t.Errorf("WIL expected 0, got %d", stats.WIL)
	}
	// 通常重みは計算される
	// INT: 100 × 1.0 = 100
	if stats.INT != 100 {
		t.Errorf("INT expected 100, got %d", stats.INT)
	}
	// LUK: 100 × 1.0 = 100
	if stats.LUK != 100 {
		t.Errorf("LUK expected 100, got %d", stats.LUK)
	}
}
