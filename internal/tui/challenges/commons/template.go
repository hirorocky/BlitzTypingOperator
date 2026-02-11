// Package commons はチャレンジ共通のテンプレート配置・文字セット生成を提供します。
package commons

import "strings"

// FormatAsPattern はテンプレートの#スロットにcharsを配置します。
// 文字数が3以下の場合はスペース区切りの1行表示にフォールバックします。
// 余剰スロットは空白に置換されます。
// 走査順は上→下、各行は左→右で固定です。
func FormatAsPattern(tmpl string, chars []rune) string {
	// 文字数3以下はフォールバック（テンプレートを使わず1行表示）
	if len(chars) <= 3 {
		var parts []string
		for _, ch := range chars {
			parts = append(parts, string(ch))
		}
		return strings.Join(parts, " ")
	}

	// テンプレートの#スロットにcharsを順番に配置
	var b strings.Builder
	charIdx := 0
	for _, ch := range tmpl {
		if ch == '#' {
			if charIdx < len(chars) {
				b.WriteRune(chars[charIdx])
				charIdx++
			} else {
				// 余剰スロットは空白に置換
				b.WriteRune(' ')
			}
		} else {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
