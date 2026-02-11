package shape

import "strings"

// 炎形テンプレート（小・中・大）
// テンプレートは # (スロット)・空白・改行のみで構成。
// View()のtextIdxカウントとの整合性を保証する。

// 小（8スロット）: charCount 4-7
var flameSmall = strings.Trim(`
  #
 # #
# # #
 # #
`, "\n")

// 中（14スロット）: charCount 8-12
var flameMedium = strings.Trim(`
   #
  # #
 # # #
# # # #
 # # #
  # #
`, "\n")

// 大（22スロット）: charCount 13+
var flameLarge = strings.Trim(`
    #
   # #
  # # #
 # # # #
# # # # #
# # # # #
 # # # #
  # # #
`, "\n")

// selectFlameTemplate は文字数に応じた炎形テンプレートを返します。
func selectFlameTemplate(charCount int) string {
	switch {
	case charCount <= 7:
		return flameSmall
	case charCount <= 12:
		return flameMedium
	default:
		return flameLarge
	}
}
