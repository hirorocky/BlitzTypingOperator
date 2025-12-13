package terminal

import (
	"testing"
)

// TestMinTerminalWidth は最小ターミナル幅の定数を検証します
func TestMinTerminalWidth(t *testing.T) {
	if MinTerminalWidth != 140 {
		t.Errorf("MinTerminalWidth should be 140, got %d", MinTerminalWidth)
	}
}

// TestMinTerminalHeight は最小ターミナル高さの定数を検証します
func TestMinTerminalHeight(t *testing.T) {
	if MinTerminalHeight != 40 {
		t.Errorf("MinTerminalHeight should be 40, got %d", MinTerminalHeight)
	}
}

// TestCheckTerminalSize_ValidSize は有効なターミナルサイズが検証を通過することを検証します
func TestCheckTerminalSize_ValidSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"exact minimum", 140, 40},
		{"larger width", 200, 40},
		{"larger height", 140, 60},
		{"both larger", 200, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckTerminalSize(tt.width, tt.height)
			if err != nil {
				t.Errorf("CheckTerminalSize(%d, %d) returned error: %v", tt.width, tt.height, err)
			}
		})
	}
}

// TestCheckTerminalSize_InvalidSize は無効なターミナルサイズがエラーを返すことを検証します
func TestCheckTerminalSize_InvalidSize(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		height int
	}{
		{"width too small", 100, 40},
		{"height too small", 140, 30},
		{"both too small", 100, 30},
		{"zero width", 0, 40},
		{"zero height", 140, 0},
		{"negative width", -1, 40},
		{"negative height", 140, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckTerminalSize(tt.width, tt.height)
			if err == nil {
				t.Errorf("CheckTerminalSize(%d, %d) should return error but got nil", tt.width, tt.height)
			}
		})
	}
}

// TestTerminalSizeError_Width はエラーに幅の情報が含まれることを検証します
func TestTerminalSizeError_Width(t *testing.T) {
	err := CheckTerminalSize(100, 40)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	sizeErr, ok := err.(*TerminalSizeError)
	if !ok {
		t.Fatalf("expected *TerminalSizeError, got %T", err)
	}

	if sizeErr.CurrentWidth != 100 {
		t.Errorf("CurrentWidth should be 100, got %d", sizeErr.CurrentWidth)
	}
	if sizeErr.RequiredWidth != 140 {
		t.Errorf("RequiredWidth should be 140, got %d", sizeErr.RequiredWidth)
	}
}

// TestTerminalSizeError_Height はエラーに高さの情報が含まれることを検証します
func TestTerminalSizeError_Height(t *testing.T) {
	err := CheckTerminalSize(140, 30)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	sizeErr, ok := err.(*TerminalSizeError)
	if !ok {
		t.Fatalf("expected *TerminalSizeError, got %T", err)
	}

	if sizeErr.CurrentHeight != 30 {
		t.Errorf("CurrentHeight should be 30, got %d", sizeErr.CurrentHeight)
	}
	if sizeErr.RequiredHeight != 40 {
		t.Errorf("RequiredHeight should be 40, got %d", sizeErr.RequiredHeight)
	}
}

// TestTerminalSizeError_Error はエラーメッセージの形式を検証します
func TestTerminalSizeError_Error(t *testing.T) {
	err := CheckTerminalSize(100, 30)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	msg := err.Error()
	if msg == "" {
		t.Error("error message should not be empty")
	}

	// エラーメッセージは両方のサイズと要件を含むべき
	if len(msg) < 20 {
		t.Error("error message should be descriptive")
	}
}

// TestNewTerminalState は新しいターミナル状態を作成します
func TestNewTerminalState(t *testing.T) {
	state := NewTerminalState(140, 40)
	if state == nil {
		t.Fatal("NewTerminalState returned nil")
	}
	if state.Width != 140 {
		t.Errorf("Width should be 140, got %d", state.Width)
	}
	if state.Height != 40 {
		t.Errorf("Height should be 40, got %d", state.Height)
	}
}

// TestTerminalState_IsValid はターミナルサイズ検証メソッドを検証します
func TestTerminalState_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		width    int
		height   int
		expected bool
	}{
		{"valid size", 140, 40, true},
		{"larger size", 200, 60, true},
		{"width too small", 100, 40, false},
		{"height too small", 140, 30, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewTerminalState(tt.width, tt.height)
			if state.IsValid() != tt.expected {
				t.Errorf("IsValid() = %v, expected %v", state.IsValid(), tt.expected)
			}
		})
	}
}

// TestTerminalState_WarningMessage は無効なサイズに対して適切な警告を返します
func TestTerminalState_WarningMessage(t *testing.T) {
	state := NewTerminalState(100, 30)
	msg := state.WarningMessage()

	if msg == "" {
		t.Error("WarningMessage should not be empty for invalid size")
	}
}

// TestTerminalState_WarningMessage_ValidSize は有効なサイズに対して空を返します
func TestTerminalState_WarningMessage_ValidSize(t *testing.T) {
	state := NewTerminalState(140, 40)
	msg := state.WarningMessage()

	if msg != "" {
		t.Errorf("WarningMessage should be empty for valid size, got: %s", msg)
	}
}

// TestFormatRecommendedSize は推奨サイズメッセージを正しくフォーマットします
func TestFormatRecommendedSize(t *testing.T) {
	msg := FormatRecommendedSize()
	if msg == "" {
		t.Error("FormatRecommendedSize should not return empty string")
	}
	// 最小サイズを含むべき
	if len(msg) < 10 {
		t.Error("FormatRecommendedSize should be descriptive")
	}
}
