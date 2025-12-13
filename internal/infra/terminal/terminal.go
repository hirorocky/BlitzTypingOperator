// Package terminal は BlitzTypingOperator TUIゲームのターミナル環境検証を提供します。
// 最小要件（140x40文字）に基づいて、ターミナルサイズの検出、
// 検証、警告メッセージの生成を実装します。
package terminal

import (
	"fmt"
)

const (
	// MinTerminalWidth は最小必須ターミナル幅（文字数）です。
	// 以下のゲームインターフェースの適切な表示を保証します：
	// - 敵情報、プレイヤーHP、モジュールリストを含むバトル画面
	// - インベントリ表示を含むエージェント管理画面
	MinTerminalWidth = 140

	// MinTerminalHeight は最小必須ターミナル高さ（文字数）です。
	// 以下の適切な表示を保証します：
	// - HPバーとステータス情報
	// - タイピングチャレンジテキストと進捗
	// - メニュー項目とナビゲーション
	MinTerminalHeight = 40
)

// TerminalSizeError はターミナルサイズが要件を満たさない場合のエラーを表します。
type TerminalSizeError struct {
	CurrentWidth   int
	CurrentHeight  int
	RequiredWidth  int
	RequiredHeight int
}

// Error は現在のサイズと必要なサイズを含む説明的なエラーメッセージを返します。
func (e *TerminalSizeError) Error() string {
	return fmt.Sprintf(
		"terminal size too small: current %dx%d, required at least %dx%d",
		e.CurrentWidth, e.CurrentHeight,
		e.RequiredWidth, e.RequiredHeight,
	)
}

// CheckTerminalSize はターミナルサイズが最小要件を満たしているか検証します。
// 有効な場合はnilを、ターミナルが小さすぎる場合はTerminalSizeErrorを返します。
func CheckTerminalSize(width, height int) error {
	if width < MinTerminalWidth || height < MinTerminalHeight {
		return &TerminalSizeError{
			CurrentWidth:   width,
			CurrentHeight:  height,
			RequiredWidth:  MinTerminalWidth,
			RequiredHeight: MinTerminalHeight,
		}
	}
	return nil
}

// TerminalState は現在のターミナルサイズと検証状態を保持します。
type TerminalState struct {
	Width  int
	Height int
}

// NewTerminalState は指定されたサイズで新しいTerminalStateを作成します。
func NewTerminalState(width, height int) *TerminalState {
	return &TerminalState{
		Width:  width,
		Height: height,
	}
}

// IsValid はターミナルサイズが最小要件を満たしている場合にtrueを返します。
func (t *TerminalState) IsValid() bool {
	return CheckTerminalSize(t.Width, t.Height) == nil
}

// WarningMessage はターミナルが小さすぎる場合に警告メッセージを返し、
// サイズが有効な場合は空文字列を返します。
func (t *TerminalState) WarningMessage() string {
	if t.IsValid() {
		return ""
	}

	var widthWarning, heightWarning string

	if t.Width < MinTerminalWidth {
		widthWarning = fmt.Sprintf("Width: %d (requires %d)", t.Width, MinTerminalWidth)
	}

	if t.Height < MinTerminalHeight {
		heightWarning = fmt.Sprintf("Height: %d (requires %d)", t.Height, MinTerminalHeight)
	}

	msg := "Terminal size is too small.\n\n"

	if widthWarning != "" {
		msg += widthWarning + "\n"
	}
	if heightWarning != "" {
		msg += heightWarning + "\n"
	}

	msg += "\n" + FormatRecommendedSize()

	return msg
}

// FormatRecommendedSize は推奨ターミナルサイズのフォーマット済み文字列を返します。
func FormatRecommendedSize() string {
	return fmt.Sprintf("Please resize your terminal to at least %dx%d characters.",
		MinTerminalWidth, MinTerminalHeight)
}
