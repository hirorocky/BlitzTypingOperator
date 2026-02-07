package challenges

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSymbolStorm_正しい入力で進捗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
	}

	c := newSymbolStormChallengeForTest(input, 42)
	sc := c.(*symbolStormChallenge)
	text := sc.text

	// 全文字を正しく入力
	for _, ch := range text {
		c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}

	result := c.Result()
	if result == nil {
		t.Fatal("全文字入力後にResult()がnilです")
	}
	if result.Status != domain.ChallengeSuccess {
		t.Errorf("Status = %d, want ChallengeSuccess", result.Status)
	}
	if result.Accuracy != 1.0 {
		t.Errorf("Accuracy = %f, want 1.0", result.Accuracy)
	}
}

func TestSymbolStorm_ESCでキャンセル(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
	}

	c := newSymbolStormChallengeForTest(input, 42)

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	result := c.Result()
	if result == nil {
		t.Fatal("ESC後にResult()がnilです")
	}
	if result.Status != domain.ChallengeCancel {
		t.Errorf("Status = %d, want ChallengeCancel", result.Status)
	}
}

func TestSymbolStorm_タイムアウトで失敗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
	}

	c := newSymbolStormChallengeForTest(input, 42)
	sc := c.(*symbolStormChallenge)

	// タイムアウトをシミュレート
	sc.startTime = time.Now().Add(-sc.timeLimit - time.Second)

	c, _ = c.Update(symbolStormTickMsg{})
	result := c.Result()
	if result == nil {
		t.Fatal("タイムアウト後にResult()がnilです")
	}
	if result.Status != domain.ChallengeFail {
		t.Errorf("Status = %d, want ChallengeFail", result.Status)
	}
}

func TestSymbolStorm_テキストは記号を含む(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
	}

	c := newSymbolStormChallengeForTest(input, 42)
	sc := c.(*symbolStormChallenge)

	// テキストが空でないこと
	if len(sc.text) == 0 {
		t.Fatal("テキストが空です")
	}

	// テキストが記号を含むこと
	hasSymbol := false
	symbols := "!@#$%^&*()_+-=[]{}|;':\",./<>?`~"
	for _, ch := range sc.text {
		for _, s := range symbols {
			if ch == s {
				hasSymbol = true
				break
			}
		}
		if hasSymbol {
			break
		}
	}
	if !hasSymbol {
		t.Errorf("テキストに記号が含まれていません: %s", sc.text)
	}
}

func TestSymbolStorm_DifficultyRateで文字数が変動(t *testing.T) {
	// 低難易度: 4-6文字
	inputLow := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(60),
	}
	cLow := newSymbolStormChallengeForTest(inputLow, 42)
	scLow := cLow.(*symbolStormChallenge)
	if len(scLow.text) < 4 || len(scLow.text) > 8 {
		t.Errorf("低難易度のテキスト長 = %d, want 4-8: text=%s", len(scLow.text), scLow.text)
	}

	// 高難易度: 10-16文字
	inputHigh := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(180),
	}
	cHigh := newSymbolStormChallengeForTest(inputHigh, 42)
	scHigh := cHigh.(*symbolStormChallenge)
	if len(scHigh.text) < 8 || len(scHigh.text) > 18 {
		t.Errorf("高難易度のテキスト長 = %d, want 8-18: text=%s", len(scHigh.text), scHigh.text)
	}
}

func TestSymbolStorm_AutoCorrect(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		AutoCorrectCount: 1,
	}

	c := newSymbolStormChallengeForTest(input, 42)
	_ = c.(*symbolStormChallenge)

	// 誤入力（AutoCorrectで無視）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'~'}})
	sc := c.(*symbolStormChallenge)
	if sc.currentIndex != 1 {
		t.Errorf("AutoCorrect後のcurrentIndex = %d, want 1", sc.currentIndex)
	}
}

func TestSymbolStorm_View出力(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
	}

	c := newSymbolStormChallengeForTest(input, 42)
	view := c.View()

	if view == "" {
		t.Error("View()が空です")
	}
}
