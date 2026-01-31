package screens

import "testing"

func TestStringDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"ASCII", "hello", 5},
		{"日本語", "こんにちは", 10},
		{"絵文字単独", "🔥", 2},
		{"絵文字VS16", "⚔️", 2},     // U+2694 + U+FE0F (lipgloss.Widthでは2幅)
		{"混在", "test 🔥 テスト", 14}, // 4 + 1 + 2 + 1 + 6 = 14
		{"空文字", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stringDisplayWidth(tt.input)
			if got != tt.expected {
				t.Errorf("stringDisplayWidth(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestTruncateToDisplayWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWidth int
		expected string
	}{
		{"短い文字列", "hello", 10, "hello"},
		{"ぴったり", "hello", 5, "hello"},
		{"切り詰め", "hello world", 5, "hello"},
		{"日本語切り詰め", "こんにちは", 6, "こんに"},
		{"絵文字保持", "🔥hello", 6, "🔥hell"},
		{"絵文字途中で切らない", "hello🔥", 6, "hello"},
		{"ゼロ幅", "hello", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateToDisplayWidth(tt.input, tt.maxWidth)
			if got != tt.expected {
				t.Errorf("truncateToDisplayWidth(%q, %d) = %q, want %q", tt.input, tt.maxWidth, got, tt.expected)
			}
		})
	}
}

func TestPadToDisplayWidth(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		targetWidth int
		expected    int // 結果の表示幅
	}{
		{"パディング追加", "hello", 10, 10},
		{"パディング不要", "hello", 5, 5},
		{"すでに超過", "hello world", 5, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padToDisplayWidth(tt.input, tt.targetWidth)
			gotWidth := stringDisplayWidth(got)
			if gotWidth != tt.expected {
				t.Errorf("padToDisplayWidth(%q, %d) width = %d, want %d", tt.input, tt.targetWidth, gotWidth, tt.expected)
			}
		})
	}
}
