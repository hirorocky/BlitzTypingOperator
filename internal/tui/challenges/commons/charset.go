// Package commons はチャレンジ共通のテンプレート配置・文字セット生成を提供します。
package commons

import "math/rand"

// 4段階の文字セット定義

// ホームポジション（DifficultyRate 50-74）
var homePositionSet = []rune("asdfghjkl;")

// キーボード上段・中段・下段の英字（DifficultyRate 75-99）
var keyboardRowSet = []rune("abcdefghijklmnopqrstuvwxyz")

// 英数字（DifficultyRate 100-149）
var alphanumericSet = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

// 記号込み（DifficultyRate 150-200）
// 基本記号（Shift不要）+ Shift+数字記号
var symbolSet = append(
	append([]rune("abcdefghijklmnopqrstuvwxyz0123456789"),
		[]rune("<>()[]{}-+=/\\|_.,")...),
	[]rune("!@#$%^&*")...,
)

// GenerateChars はDifficultyRateに応じた文字列を生成します。
// 文字セットは4段階で切り替わり、文字数は難易度に連動します（最低1文字）。
func GenerateChars(diffRate int, rng *rand.Rand) []rune {
	charset := selectCharset(diffRate)
	length := calcLength(diffRate, rng)

	chars := make([]rune, length)
	for i := range chars {
		chars[i] = charset[rng.Intn(len(charset))]
	}
	return chars
}

// selectCharset はDifficultyRateに応じた文字セットを返します。
func selectCharset(diffRate int) []rune {
	switch {
	case diffRate < 75:
		return homePositionSet
	case diffRate < 100:
		return keyboardRowSet
	case diffRate < 150:
		return alphanumericSet
	default:
		return symbolSet
	}
}

// calcLength はDifficultyRateに応じた文字数を決定します。
// 既存のsymbol_stormの文字数ロジックを継承（最低1文字保証）。
func calcLength(diffRate int, rng *rand.Rand) int {
	var minLen, maxLen int
	if diffRate <= 50 {
		minLen, maxLen = 4, 6
	} else if diffRate <= 100 {
		t := float64(diffRate-50) / 50.0
		minLen = 4 + int(t*2)
		maxLen = 6 + int(t*2)
	} else if diffRate <= 200 {
		t := float64(diffRate-100) / 100.0
		minLen = 6 + int(t*4)
		maxLen = 8 + int(t*8)
	} else {
		minLen, maxLen = 10, 16
	}

	length := minLen + rng.Intn(maxLen-minLen+1)
	if length < 1 {
		length = 1
	}
	return length
}
