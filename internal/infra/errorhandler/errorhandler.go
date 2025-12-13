// Package errorhandler はエラー処理とログ機能を担当します。
// 入力検証、エラーハンドリング、デバッグモード、ログ機能を提供します。

package errorhandler

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ==================================================
// エラー型定義
// ==================================================

// ErrorType はエラーの種類を表す型です。
type ErrorType int

const (
	// ErrInvalidInput は無効な入力値エラーを表します。
	ErrInvalidInput ErrorType = iota
	// ErrSaveLoad はセーブ/ロード関連エラーを表します。
	ErrSaveLoad
	// ErrBattle はバトル処理中のエラーを表します。
	ErrBattle
	// ErrUnexpected は予期しないエラーを表します。
	ErrUnexpected
)

// GameError はゲーム固有のエラー型です。

type GameError struct {
	// Type はエラーの種類です。
	Type ErrorType
	// Message はエラーメッセージです。
	Message string
	// Details はエラーの詳細情報です。
	Details string
	// Cause は元のエラー（ラップ時）です。
	Cause error
}

// Error はerrorインターフェースを実装します。
func (e *GameError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

// Unwrap は元のエラーを返します。
func (e *GameError) Unwrap() error {
	return e.Cause
}

// NewGameError は新しいGameErrorを作成します。
func NewGameError(errorType ErrorType, details string) *GameError {
	var message string
	switch errorType {
	case ErrInvalidInput:
		message = "無効な入力"
	case ErrSaveLoad:
		message = "セーブ/ロードエラー"
	case ErrBattle:
		message = "バトルエラー"
	case ErrUnexpected:
		message = "予期しないエラー"
	default:
		message = "エラー"
	}

	return &GameError{
		Type:    errorType,
		Message: message,
		Details: details,
	}
}

// WrapError は既存のエラーをラップして新しいGameErrorを作成します。
func WrapError(err error, context string) *GameError {
	return &GameError{
		Type:    ErrUnexpected,
		Message: context,
		Details: err.Error(),
		Cause:   err,
	}
}

// ==================================================
// 入力検証 (Requirement 19.4)
// ==================================================

// ValidateLevel はレベル入力の検証を行います。

func ValidateLevel(level int, maxLevelReached int) error {
	if level <= 0 {
		return NewGameError(ErrInvalidInput, "レベルは1以上である必要があります")
	}
	// 挑戦可能な最大レベルは maxLevelReached + 1
	maxAllowed := maxLevelReached + 1
	if level > maxAllowed {
		return NewGameError(ErrInvalidInput,
			fmt.Sprintf("レベル%dは挑戦できません（最大: %d）", level, maxAllowed))
	}
	return nil
}

// ValidateAgentSlot はエージェントスロット番号の検証を行います。

func ValidateAgentSlot(slot int) error {
	// 最大3スロット (0, 1, 2)
	if slot < 0 || slot > 2 {
		return NewGameError(ErrInvalidInput,
			fmt.Sprintf("エージェントスロットは0〜2である必要があります（指定: %d）", slot))
	}
	return nil
}

// ValidatePositiveInt は正の整数の検証を行います。

func ValidatePositiveInt(value int, fieldName string) error {
	if value <= 0 {
		return NewGameError(ErrInvalidInput,
			fmt.Sprintf("%sは正の値である必要があります（指定: %d）", fieldName, value))
	}
	return nil
}

// ValidateNonNegativeInt は非負の整数の検証を行います。

func ValidateNonNegativeInt(value int, fieldName string) error {
	if value < 0 {
		return NewGameError(ErrInvalidInput,
			fmt.Sprintf("%sは0以上である必要があります（指定: %d）", fieldName, value))
	}
	return nil
}

// ValidateString は非空文字列の検証を行います。

func ValidateString(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return NewGameError(ErrInvalidInput,
			fmt.Sprintf("%sは空であってはなりません", fieldName))
	}
	return nil
}

// ==================================================
// パニック回復 (Requirement 19.4: ゲームクラッシュ防止)
// ==================================================

// RecoverToError はパニックからエラーを回復します。

func RecoverToError(r interface{}) error {
	if r == nil {
		return nil
	}
	switch v := r.(type) {
	case error:
		return v
	case string:
		return NewGameError(ErrUnexpected, v)
	default:
		return NewGameError(ErrUnexpected, fmt.Sprintf("%v", v))
	}
}

// ==================================================
// デバッグモード (Requirement 19.6)
// ==================================================

var (
	debugMode   bool
	debugModeMu sync.RWMutex
)

// SetDebugMode はデバッグモードを設定します。

func SetDebugMode(enabled bool) {
	debugModeMu.Lock()
	defer debugModeMu.Unlock()
	debugMode = enabled
}

// IsDebugMode はデバッグモードかどうかを返します。
func IsDebugMode() bool {
	debugModeMu.RLock()
	defer debugModeMu.RUnlock()
	return debugMode
}

// ==================================================
// ログ機能 (Requirement 19.6)
// ==================================================

// LogLevel はログレベルを表す型です。
type LogLevel int

const (
	LogLevelError LogLevel = iota
	LogLevelInfo
	LogLevelDebug
)

// Logger はログ機能を提供する構造体です。

type Logger struct {
	filePath string
	file     *os.File
	mu       sync.Mutex
}

// NewLogger は新しいLoggerを作成します。
func NewLogger(filePath string) *Logger {
	logger := &Logger{
		filePath: filePath,
	}
	// ファイルを開く（追記モード）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		logger.file = file
	}
	return logger
}

// Close はログファイルを閉じます。
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}
}

// log は指定されたレベルでログを書き込みます。
func (l *Logger) log(level LogLevel, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return
	}

	levelStr := ""
	switch level {
	case LogLevelError:
		levelStr = "ERROR"
	case LogLevelInfo:
		levelStr = "INFO"
	case LogLevelDebug:
		levelStr = "DEBUG"
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelStr, message)
	_, _ = l.file.WriteString(logLine)
}

// Error はエラーレベルのログを書き込みます。
func (l *Logger) Error(message string) {
	l.log(LogLevelError, message)
}

// Info は情報レベルのログを書き込みます。
func (l *Logger) Info(message string) {
	l.log(LogLevelInfo, message)
}

// Debug はデバッグレベルのログを書き込みます。
// デバッグモードがONの場合のみ書き込まれます。
func (l *Logger) Debug(message string) {
	if !IsDebugMode() {
		return
	}
	l.log(LogLevelDebug, message)
}

// LogException は例外（エラー）をログに記録します。

func (l *Logger) LogException(err error, context string) {
	message := fmt.Sprintf("例外発生 [%s]: %v", context, err)
	l.Error(message)
}

// ==================================================
// グローバルロガー
// ==================================================

var (
	globalLogger *Logger
	loggerOnce   sync.Once
)

// GetGlobalLogger はグローバルロガーを取得します。
// ロガーが未初期化の場合はnilを返します。
func GetGlobalLogger() *Logger {
	return globalLogger
}

// InitGlobalLogger はグローバルロガーを初期化します。
func InitGlobalLogger(logPath string) {
	loggerOnce.Do(func() {
		globalLogger = NewLogger(logPath)
	})
}

// LogError はグローバルロガーでエラーを記録します。
func LogError(message string) {
	if globalLogger != nil {
		globalLogger.Error(message)
	}
}

// LogInfo はグローバルロガーで情報を記録します。
func LogInfo(message string) {
	if globalLogger != nil {
		globalLogger.Info(message)
	}
}

// LogDebug はグローバルロガーでデバッグ情報を記録します。
func LogDebug(message string) {
	if globalLogger != nil {
		globalLogger.Debug(message)
	}
}

// LogException はグローバルロガーで例外を記録します。
func LogExceptionGlobal(err error, context string) {
	if globalLogger != nil {
		globalLogger.LogException(err, context)
	}
}
