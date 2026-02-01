package components

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

func TestFormatStats(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "test_core",
		Name: "テストコア",
		StatWeights: map[string]float64{
			"STR": 1.2,
			"INT": 1.0,
			"WIL": 0.8,
			"LUK": 1.0,
		},
	}

	result := FormatStats(coreType)

	// 固定順序（STR, INT, WIL, LUK）であることを確認
	if !strings.Contains(result, "STR: 120") {
		t.Errorf("STR value not found in result: %s", result)
	}
	if !strings.Contains(result, "INT: 100") {
		t.Errorf("INT value not found in result: %s", result)
	}
	if !strings.Contains(result, "WIL: 80") {
		t.Errorf("WIL value not found in result: %s", result)
	}
	if !strings.Contains(result, "LUK: 100") {
		t.Errorf("LUK value not found in result: %s", result)
	}

	// 順序が STR -> INT -> WIL -> LUK であることを確認
	strPos := strings.Index(result, "STR")
	intPos := strings.Index(result, "INT")
	wilPos := strings.Index(result, "WIL")
	lukPos := strings.Index(result, "LUK")

	if !(strPos < intPos && intPos < wilPos && wilPos < lukPos) {
		t.Errorf("Stats order is incorrect. Expected STR < INT < WIL < LUK, got positions: STR=%d, INT=%d, WIL=%d, LUK=%d",
			strPos, intPos, wilPos, lukPos)
	}
}

func TestFormatStatsFromStats(t *testing.T) {
	stats := domain.Stats{
		STR: 150,
		INT: 80,
		WIL: 100,
		LUK: 120,
	}

	result := FormatStatsFromStats(stats)

	expected := "STR: 150  INT: 80  WIL: 100  LUK: 120"
	if result != expected {
		t.Errorf("FormatStatsFromStats() = %q, want %q", result, expected)
	}
}

func TestFormatStatsMultiline(t *testing.T) {
	coreType := domain.CoreType{
		ID:   "test_core",
		Name: "テストコア",
		StatWeights: map[string]float64{
			"STR": 1.5,
			"INT": 0.5,
			"WIL": 1.0,
			"LUK": 0.8,
		},
	}

	result := FormatStatsMultiline(coreType)

	// 複数行であることを確認
	lines := strings.Split(result, "\n")
	if len(lines) != 4 {
		t.Errorf("FormatStatsMultiline() should have 4 lines, got %d", len(lines))
	}

	// 各行の値を確認
	if !strings.HasPrefix(lines[0], "STR:") {
		t.Errorf("First line should start with STR:, got %s", lines[0])
	}
	if !strings.HasPrefix(lines[1], "INT:") {
		t.Errorf("Second line should start with INT:, got %s", lines[1])
	}
	if !strings.HasPrefix(lines[2], "WIL:") {
		t.Errorf("Third line should start with WIL:, got %s", lines[2])
	}
	if !strings.HasPrefix(lines[3], "LUK:") {
		t.Errorf("Fourth line should start with LUK:, got %s", lines[3])
	}
}

func TestFormatStats_ZeroWeights(t *testing.T) {
	coreType := domain.CoreType{
		ID:          "empty_core",
		Name:        "空コア",
		StatWeights: map[string]float64{},
	}

	result := FormatStats(coreType)

	// すべて0になるはず
	expected := "STR: 0  INT: 0  WIL: 0  LUK: 0"
	if result != expected {
		t.Errorf("FormatStats() with empty weights = %q, want %q", result, expected)
	}
}
