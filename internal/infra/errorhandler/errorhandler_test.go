// Package errorhandler はエラー処理とログ機能を担当します。

package errorhandler

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ==================================================
// Task 13.1: 入力検証とエラーハンドリングテスト
// ==================================================

func TestValidateLevel_InvalidInputs(t *testing.T) {

	tests := []struct {
		name        string
		level       int
		maxLevel    int
		expectError bool
	}{
		{"level 0", 0, 10, true},
		{"negative level", -1, 10, true},
		{"level above max", 12, 10, true},
		{"valid level", 5, 10, false},
		{"level equals max+1", 11, 10, false}, // 挑戦可能は max+1 まで
		{"level 1", 1, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLevel(tt.level, tt.maxLevel)
			if tt.expectError && err == nil {
				t.Error("エラーが返されるべきです")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが返されるべきではありません: %v", err)
			}
		})
	}
}

func TestValidateAgentSlot_InvalidInputs(t *testing.T) {

	tests := []struct {
		name        string
		slot        int
		expectError bool
	}{
		{"negative slot", -1, true},
		{"slot 0", 0, false},
		{"slot 1", 1, false},
		{"slot 2", 2, false},
		{"slot 3", 3, true}, // 最大3スロット (0, 1, 2)
		{"slot 100", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAgentSlot(tt.slot)
			if tt.expectError && err == nil {
				t.Error("エラーが返されるべきです")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが返されるべきではありません: %v", err)
			}
		})
	}
}

func TestValidatePositiveInt(t *testing.T) {

	tests := []struct {
		name        string
		value       int
		fieldName   string
		expectError bool
	}{
		{"negative value", -1, "HP", true},
		{"zero value", 0, "HP", true},
		{"positive value", 1, "HP", false},
		{"large positive", 1000, "damage", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositiveInt(tt.value, tt.fieldName)
			if tt.expectError && err == nil {
				t.Error("エラーが返されるべきです")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが返されるべきではありません: %v", err)
			}
		})
	}
}

func TestValidateNonNegativeInt(t *testing.T) {

	tests := []struct {
		name        string
		value       int
		fieldName   string
		expectError bool
	}{
		{"negative value", -1, "HP", true},
		{"zero value", 0, "HP", false},
		{"positive value", 1, "HP", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonNegativeInt(tt.value, tt.fieldName)
			if tt.expectError && err == nil {
				t.Error("エラーが返されるべきです")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが返されるべきではありません: %v", err)
			}
		})
	}
}

func TestValidateString(t *testing.T) {

	tests := []struct {
		name        string
		value       string
		fieldName   string
		expectError bool
	}{
		{"empty string", "", "ID", true},
		{"whitespace only", "   ", "ID", true},
		{"valid string", "core_001", "ID", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.value, tt.fieldName)
			if tt.expectError && err == nil {
				t.Error("エラーが返されるべきです")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが返されるべきではありません: %v", err)
			}
		})
	}
}

func TestGameError_Message(t *testing.T) {

	err := NewGameError(ErrInvalidInput, "レベルは1以上である必要があります")
	if err.Error() == "" {
		t.Error("エラーメッセージが返されるべきです")
	}
	if !strings.Contains(err.Error(), "レベルは1以上") {
		t.Error("エラーメッセージに詳細が含まれるべきです")
	}
}

func TestRecoverFromPanic(t *testing.T) {

	var recovered error

	func() {
		defer func() {
			if r := recover(); r != nil {
				recovered = RecoverToError(r)
			}
		}()
		panic("test panic")
	}()

	if recovered == nil {
		t.Error("パニックからエラーが回復されるべきです")
	}
}

// ==================================================
// Task 13.2: デバッグモードとログ機能テスト
// ==================================================

func TestDebugMode_Toggle(t *testing.T) {

	SetDebugMode(true)
	if !IsDebugMode() {
		t.Error("デバッグモードがtrueであるべきです")
	}

	SetDebugMode(false)
	if IsDebugMode() {
		t.Error("デバッグモードがfalseであるべきです")
	}
}

func TestLogger_WriteToFile(t *testing.T) {

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "debug.log")

	logger := NewLogger(logPath)
	defer logger.Close()

	// ログを書き込み
	logger.Error("テストエラー")
	logger.Info("テスト情報")
	logger.Debug("デバッグ情報")

	// ファイルが作成されていることを確認
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("ログファイルが作成されるべきです")
	}

	// ファイル内容を確認
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗: %v", err)
	}

	if !strings.Contains(string(content), "テストエラー") {
		t.Error("エラーログが含まれるべきです")
	}
	if !strings.Contains(string(content), "テスト情報") {
		t.Error("情報ログが含まれるべきです")
	}
}

func TestLogException(t *testing.T) {

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "debug.log")

	logger := NewLogger(logPath)
	defer logger.Close()

	// 例外をログに記録
	err := errors.New("予期しないエラー")
	logger.LogException(err, "バトル処理中")

	// ファイル内容を確認
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("ログファイルの読み込みに失敗: %v", err)
	}

	if !strings.Contains(string(content), "予期しないエラー") {
		t.Error("例外メッセージが含まれるべきです")
	}
	if !strings.Contains(string(content), "バトル処理中") {
		t.Error("コンテキスト情報が含まれるべきです")
	}
}

func TestLogger_DebugModeOnly(t *testing.T) {

	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "debug.log")

	logger := NewLogger(logPath)
	defer logger.Close()

	// デバッグモードOFFでDebugログは書き込まれない
	SetDebugMode(false)
	logger.Debug("デバッグOFF時のログ")

	// デバッグモードONでDebugログが書き込まれる
	SetDebugMode(true)
	logger.Debug("デバッグON時のログ")

	content, _ := os.ReadFile(logPath)

	if strings.Contains(string(content), "デバッグOFF時のログ") {
		t.Error("デバッグモードOFF時はDebugログが書き込まれるべきではありません")
	}
	// 注: デバッグON時のログは書き込まれることを確認
	if !strings.Contains(string(content), "デバッグON時のログ") {
		t.Error("デバッグモードON時はDebugログが書き込まれるべきです")
	}
}

// ==================================================
// エラー型のテスト
// ==================================================

func TestGameErrorType(t *testing.T) {
	tests := []struct {
		errorType ErrorType
		expected  string
	}{
		{ErrInvalidInput, "無効な入力"},
		{ErrSaveLoad, "セーブ/ロードエラー"},
		{ErrBattle, "バトルエラー"},
		{ErrUnexpected, "予期しないエラー"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			err := NewGameError(tt.errorType, "詳細")
			if err.Type != tt.errorType {
				t.Errorf("エラータイプが一致しません: got %v, want %v", err.Type, tt.errorType)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("元のエラー")
	wrapped := WrapError(originalErr, "追加コンテキスト")

	if !strings.Contains(wrapped.Error(), "元のエラー") {
		t.Error("元のエラーメッセージが含まれるべきです")
	}
	if !strings.Contains(wrapped.Error(), "追加コンテキスト") {
		t.Error("追加コンテキストが含まれるべきです")
	}
}
