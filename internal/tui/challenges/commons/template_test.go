package commons

import (
	"strings"
	"testing"
)

// === 受け入れテスト: テンプレートシステム ===

// 受け入れ基準9: テンプレートのスロット数が入力対象文字数以上であり、余剰スロットは空白に置換される
// 受け入れ基準28: パターンの非空白文字を除去した結果がtextと一致する（表示と入力の整合性）
// 受け入れ基準32: patternの非空白文字数がtextの文字数と一致する
func TestFormatAsPattern_スロット整合性(t *testing.T) {
	tmpl := "# # #\n # # \n# # #"
	chars := []rune{'a', 'b', 'c', 'd', 'e'}

	pattern := FormatAsPattern(tmpl, chars)

	// patternから空白・改行を除去した文字列が入力charsと一致する
	var nonSpace []rune
	for _, ch := range pattern {
		if ch != ' ' && ch != '\n' {
			nonSpace = append(nonSpace, ch)
		}
	}

	if string(nonSpace) != string(chars) {
		t.Errorf("パターンの非空白文字 = %q, want %q", string(nonSpace), string(chars))
	}
}

// 受け入れ基準9: 余剰スロットは空白に置換される
func TestFormatAsPattern_余剰スロットは空白(t *testing.T) {
	tmpl := "# # # #\n # # # \n# # # #"
	chars := []rune{'a', 'b'} // テンプレートより少ない文字数

	pattern := FormatAsPattern(tmpl, chars)

	// 非空白文字がcharsと一致する
	var nonSpace []rune
	for _, ch := range pattern {
		if ch != ' ' && ch != '\n' {
			nonSpace = append(nonSpace, ch)
		}
	}

	if string(nonSpace) != string(chars) {
		t.Errorf("余剰スロット後の非空白文字 = %q, want %q", string(nonSpace), string(chars))
	}
}

// 受け入れ基準11: 文字数が3以下の場合はテンプレートを使わず、スペース区切りの1行表示にフォールバック
func TestFormatAsPattern_フォールバック_文字数1(t *testing.T) {
	tmpl := "# #\n# #"
	chars := []rune{'x'}

	pattern := FormatAsPattern(tmpl, chars)

	// 改行を含まない1行表示
	if strings.Contains(pattern, "\n") {
		t.Errorf("文字数1のフォールバック時に改行が含まれる: %q", pattern)
	}
	// 文字が含まれる
	if !strings.ContainsRune(pattern, 'x') {
		t.Errorf("フォールバック時に文字が含まれない: %q", pattern)
	}
}

func TestFormatAsPattern_フォールバック_文字数3(t *testing.T) {
	tmpl := "# # #\n # # \n# # #"
	chars := []rune{'a', 'b', 'c'}

	pattern := FormatAsPattern(tmpl, chars)

	// 改行を含まない1行表示
	if strings.Contains(pattern, "\n") {
		t.Errorf("文字数3のフォールバック時に改行が含まれる: %q", pattern)
	}
	// スペース区切りで全文字が含まれる
	if pattern != "a b c" {
		t.Errorf("フォールバック表示 = %q, want %q", pattern, "a b c")
	}
}

// 受け入れ基準8: テンプレートは#（スロット）・空白・改行のみで構成される（他の文字があればエラー検出用）
func TestFormatAsPattern_テンプレート適用で改行あり(t *testing.T) {
	tmpl := "# #\n# #\n# #"
	chars := []rune{'a', 'b', 'c', 'd', 'e', 'f'}

	pattern := FormatAsPattern(tmpl, chars)

	// 文字数4以上でテンプレート適用 → 改行を含む
	if !strings.Contains(pattern, "\n") {
		t.Errorf("文字数4以上のパターンに改行が含まれない: %q", pattern)
	}
}

// 受け入れ基準10: 入力順はテンプレート上の走査順（上から下、各行は左から右）で固定
func TestFormatAsPattern_走査順が入力順(t *testing.T) {
	tmpl := " # \n# #\n # "
	chars := []rune{'1', '2', '3', '4'}

	pattern := FormatAsPattern(tmpl, chars)

	// 走査順に非空白文字を取得
	var order []rune
	for _, ch := range pattern {
		if ch != ' ' && ch != '\n' {
			order = append(order, ch)
		}
	}

	expected := []rune{'1', '2', '3', '4'}
	if string(order) != string(expected) {
		t.Errorf("走査順 = %q, want %q", string(order), string(expected))
	}
}
