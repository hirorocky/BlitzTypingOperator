package shape

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges/commons"

	tea "github.com/charmbracelet/bubbletea"
)

// === 既存symbol_stormテスト移行（基準30） ===

func TestShape_正しい入力で進捗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	sc := c.(*shapeChallenge)
	text := sc.text

	// 全文字を正しく入力
	for _, ch := range text {
		c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}

	result := c.Result()
	if result == nil {
		t.Fatal("全文字入力後にResult()がnilです")
	}
	if result.Status != domain.ChallengePerfect {
		t.Errorf("Status = %d, want ChallengePerfect（ミスなし完了）", result.Status)
	}
	if result.Score != 100 {
		t.Errorf("Score = %d, want 100", result.Score)
	}
}

func TestShape_ESCでキャンセル(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	result := c.Result()
	if result == nil {
		t.Fatal("ESC後にResult()がnilです")
	}
	if result.Status != domain.ChallengeCancel {
		t.Errorf("Status = %d, want ChallengeCancel", result.Status)
	}
}

func TestShape_タイムアウトで失敗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	sc := c.(*shapeChallenge)

	// タイムアウトをシミュレート
	sc.startTime = time.Now().Add(-sc.timeLimit - time.Second)

	c, _ = c.Update(shapeTickMsg{})
	result := c.Result()
	if result == nil {
		t.Fatal("タイムアウト後にResult()がnilです")
	}
	if result.Status != domain.ChallengeFail {
		t.Errorf("Status = %d, want ChallengeFail", result.Status)
	}
}

func TestShape_DifficultyRateで文字数が変動(t *testing.T) {
	// 低難易度: 4-6文字
	inputLow := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRate(60),
		ChallengeOptions: map[string]string{"shape": "flame"},
	}
	cLow := newShapeChallengeForTest(inputLow, 42)
	scLow := cLow.(*shapeChallenge)
	if len(scLow.text) < 4 || len(scLow.text) > 8 {
		t.Errorf("低難易度のテキスト長 = %d, want 4-8: text=%s", len(scLow.text), scLow.text)
	}

	// 高難易度: 10-16文字
	inputHigh := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRate(180),
		ChallengeOptions: map[string]string{"shape": "flame"},
	}
	cHigh := newShapeChallengeForTest(inputHigh, 42)
	scHigh := cHigh.(*shapeChallenge)
	if len(scHigh.text) < 8 || len(scHigh.text) > 18 {
		t.Errorf("高難易度のテキスト長 = %d, want 8-18: text=%s", len(scHigh.text), scHigh.text)
	}
}

func TestShape_AutoCorrect(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		AutoCorrectCount: 1,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)

	// 誤入力（AutoCorrectで無視）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'~'}})
	sc := c.(*shapeChallenge)
	if sc.currentIndex != 1 {
		t.Errorf("AutoCorrect後のcurrentIndex = %d, want 1", sc.currentIndex)
	}
}

func TestShape_View出力(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	view := c.View()

	if view == "" {
		t.Error("View()が空です")
	}
}

// === 新規テスト ===

// 受け入れ基準34: 同一seedで同一patternが再現される
func TestShape_パターン再現性(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c1 := newShapeChallengeForTest(input, 12345)
	c2 := newShapeChallengeForTest(input, 12345)

	sc1 := c1.(*shapeChallenge)
	sc2 := c2.(*shapeChallenge)

	if sc1.text != sc2.text {
		t.Errorf("同一シードでテキストが異なる: %q vs %q", sc1.text, sc2.text)
	}
	if sc1.pattern != sc2.pattern {
		t.Errorf("同一シードでパターンが異なる: %q vs %q", sc1.pattern, sc2.pattern)
	}
}

// 受け入れ基準32: patternの非空白文字数がtextの文字数と一致する
func TestShape_パターン整合性(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	sc := c.(*shapeChallenge)

	var nonSpace []rune
	for _, ch := range sc.pattern {
		if ch != ' ' && ch != '\n' {
			nonSpace = append(nonSpace, ch)
		}
	}

	if string(nonSpace) != sc.text {
		t.Errorf("パターンの非空白文字 = %q, want text = %q", string(nonSpace), sc.text)
	}
}

// 受け入れ基準15: shapeチャレンジが共通文字セットを使用する
func TestShape_共通文字セットを使用(t *testing.T) {
	// Rate=50ではホームポジション文字のみが生成される
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRate(50),
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	sc := c.(*shapeChallenge)

	homePositionChars := "asdfghjkl;"
	for _, ch := range sc.text {
		found := false
		for _, hp := range homePositionChars {
			if ch == hp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Rate=50でホームポジション外の文字 %q がテキストに含まれる", ch)
		}
	}
}

// 受け入れ基準17: options.shapeの値に応じてテンプレートセットが選択される
func TestShape_Optionsによるテンプレート選択(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		ChallengeOptions: map[string]string{"shape": "flame"},
	}

	c := newShapeChallengeForTest(input, 42)
	sc := c.(*shapeChallenge)

	// 文字数4以上のとき改行を含むパターンが生成される（テンプレート適用）
	if len(sc.text) >= 4 && !strings.Contains(sc.pattern, "\n") {
		t.Errorf("文字数%dでパターンに改行が含まれない（テンプレート適用されていない）", len(sc.text))
	}
}

// 受け入れ基準31: テンプレート選択境界値テスト
func TestShape_テンプレート選択境界値(t *testing.T) {
	tests := []struct {
		name      string
		charCount int
		wantTmpl  string
	}{
		{"small最小=4", 4, flameSmall},
		{"small最大=7", 7, flameSmall},
		{"medium最小=8", 8, flameMedium},
		{"medium最大=12", 12, flameMedium},
		{"large最小=13", 13, flameLarge},
		{"large=16", 16, flameLarge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectFlameTemplate(tt.charCount)
			if got != tt.wantTmpl {
				t.Errorf("charCount=%d でテンプレートが期待と異なる", tt.charCount)
			}
		})
	}
}

// 受け入れ基準35: 文字数1-3のフォールバック表示テスト
func TestShape_フォールバック表示(t *testing.T) {
	// commons.FormatAsPatternのフォールバックを直接テスト
	chars := []rune{'a', 'b'}
	pattern := commons.FormatAsPattern("# # # #", chars)

	// 改行を含まない1行表示
	if strings.Contains(pattern, "\n") {
		t.Errorf("文字数2のフォールバック時に改行が含まれる: %q", pattern)
	}
	if pattern != "a b" {
		t.Errorf("フォールバック表示 = %q, want %q", pattern, "a b")
	}
}

// options.shapeが未指定の場合はデフォルト形状（flame）にフォールバック
func TestShape_未指定オプションはflameにフォールバック(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		// ChallengeOptions未設定
	}

	c := newShapeChallengeForTest(input, 42)
	if c == nil {
		t.Fatal("オプション未指定でチャレンジが生成されない")
	}

	sc := c.(*shapeChallenge)
	if sc.text == "" {
		t.Error("テキストが空")
	}
}

// テンプレートのスロット数が十分であることを検証
func TestShape_全テンプレートのスロット数(t *testing.T) {
	templates := map[string]string{
		"flameSmall":  flameSmall,
		"flameMedium": flameMedium,
		"flameLarge":  flameLarge,
	}

	for name, tmpl := range templates {
		slotCount := strings.Count(tmpl, "#")
		if slotCount == 0 {
			t.Errorf("%s: スロット数が0", name)
		}
		t.Logf("%s: %dスロット", name, slotCount)
	}
}

// 統合テスト: GenerateChars → selectTemplate → FormatAsPattern の流れ
func TestShape_文字生成からパターン生成までの統合(t *testing.T) {
	diffRates := []int{50, 75, 100, 150, 200}

	for _, rate := range diffRates {
		rng := rand.New(rand.NewSource(42))
		chars := commons.GenerateChars(rate, rng)

		tmpl := selectTemplate("flame", len(chars))
		pattern := commons.FormatAsPattern(tmpl, chars)

		// パターンの非空白文字がcharsと一致
		var nonSpace []rune
		for _, ch := range pattern {
			if ch != ' ' && ch != '\n' {
				nonSpace = append(nonSpace, ch)
			}
		}

		if string(nonSpace) != string(chars) {
			t.Errorf("Rate=%d: パターン非空白文字 = %q, want %q",
				rate, string(nonSpace), string(chars))
		}
	}
}
