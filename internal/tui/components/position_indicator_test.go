package components

import (
	"strings"
	"testing"
)

func TestRenderHorizontalIndicator(t *testing.T) {
	tests := []struct {
		name     string
		total    int
		current  int
		expected string
	}{
		{
			name:     "3項目で最初を選択",
			total:    3,
			current:  0,
			expected: "■  □  □",
		},
		{
			name:     "3項目で中央を選択",
			total:    3,
			current:  1,
			expected: "□  ■  □",
		},
		{
			name:     "3項目で最後を選択",
			total:    3,
			current:  2,
			expected: "□  □  ■",
		},
		{
			name:     "1項目のみ",
			total:    1,
			current:  0,
			expected: "■",
		},
		{
			name:     "5項目で3番目を選択",
			total:    5,
			current:  2,
			expected: "□  □  ■  □  □",
		},
		{
			name:     "total=0の場合は空文字",
			total:    0,
			current:  0,
			expected: "",
		},
		{
			name:     "currentが範囲外（負）の場合は0に補正",
			total:    3,
			current:  -1,
			expected: "■  □  □",
		},
		{
			name:     "currentが範囲外（大）の場合は最大に補正",
			total:    3,
			current:  5,
			expected: "□  □  ■",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderHorizontalIndicator(tt.total, tt.current)
			if result != tt.expected {
				t.Errorf("RenderHorizontalIndicator(%d, %d) = %q, want %q",
					tt.total, tt.current, result, tt.expected)
			}
		})
	}
}

func TestRenderVerticalIndicator(t *testing.T) {
	tests := []struct {
		name         string
		total        int
		current      int
		expectedLine int // 選択されている行番号
	}{
		{
			name:         "3項目で最初を選択",
			total:        3,
			current:      0,
			expectedLine: 0,
		},
		{
			name:         "3項目で中央を選択",
			total:        3,
			current:      1,
			expectedLine: 1,
		},
		{
			name:         "3項目で最後を選択",
			total:        3,
			current:      2,
			expectedLine: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderVerticalIndicator(tt.total, tt.current)
			lines := strings.Split(result, "\n")

			if len(lines) != tt.total {
				t.Errorf("RenderVerticalIndicator(%d, %d) has %d lines, want %d",
					tt.total, tt.current, len(lines), tt.total)
				return
			}

			for i, line := range lines {
				if i == tt.expectedLine {
					if line != IndicatorSelected {
						t.Errorf("Line %d should be selected (■), got %q", i, line)
					}
				} else {
					if line != IndicatorUnselected {
						t.Errorf("Line %d should be unselected (□), got %q", i, line)
					}
				}
			}
		})
	}
}

func TestRenderVerticalIndicator_EdgeCases(t *testing.T) {
	// total=0の場合は空文字
	result := RenderVerticalIndicator(0, 0)
	if result != "" {
		t.Errorf("RenderVerticalIndicator(0, 0) = %q, want empty string", result)
	}

	// total=1の場合は1行のみ
	result = RenderVerticalIndicator(1, 0)
	if result != "■" {
		t.Errorf("RenderVerticalIndicator(1, 0) = %q, want ■", result)
	}
}

func TestRenderEnemySelectionIndicator(t *testing.T) {
	tests := []struct {
		name          string
		total         int
		current       int
		defeatedFlags []bool
		expected      string
	}{
		{
			name:          "3体で最初を選択、全て未撃破",
			total:         3,
			current:       0,
			defeatedFlags: []bool{false, false, false},
			expected:      "◀  ■  □  □  ▶",
		},
		{
			name:          "3体で中央を選択、全て未撃破",
			total:         3,
			current:       1,
			defeatedFlags: []bool{false, false, false},
			expected:      "◀  □  ■  □  ▶",
		},
		{
			name:          "3体で最初を選択、最初が撃破済み",
			total:         3,
			current:       0,
			defeatedFlags: []bool{true, false, false},
			expected:      "◀  ★  □  □  ▶",
		},
		{
			name:          "3体で2番目を選択、1番目と3番目が撃破済み",
			total:         3,
			current:       1,
			defeatedFlags: []bool{true, false, true},
			expected:      "◀  ☆  ■  ☆  ▶",
		},
		{
			name:          "3体で最後を選択、全て撃破済み",
			total:         3,
			current:       2,
			defeatedFlags: []bool{true, true, true},
			expected:      "◀  ☆  ☆  ★  ▶",
		},
		{
			name:          "1体のみ、未撃破",
			total:         1,
			current:       0,
			defeatedFlags: []bool{false},
			expected:      "◀  ■  ▶",
		},
		{
			name:          "1体のみ、撃破済み",
			total:         1,
			current:       0,
			defeatedFlags: []bool{true},
			expected:      "◀  ★  ▶",
		},
		{
			name:          "defeatedFlagsが空の場合は全て未撃破扱い",
			total:         3,
			current:       1,
			defeatedFlags: []bool{},
			expected:      "◀  □  ■  □  ▶",
		},
		{
			name:          "defeatedFlagsがtotalより少ない場合は残りは未撃破扱い",
			total:         3,
			current:       2,
			defeatedFlags: []bool{true},
			expected:      "◀  ☆  □  ■  ▶",
		},
		{
			name:          "total=0の場合は空文字",
			total:         0,
			current:       0,
			defeatedFlags: []bool{},
			expected:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderEnemySelectionIndicator(tt.total, tt.current, tt.defeatedFlags)
			if result != tt.expected {
				t.Errorf("RenderEnemySelectionIndicator(%d, %d, %v) = %q, want %q",
					tt.total, tt.current, tt.defeatedFlags, result, tt.expected)
			}
		})
	}
}
