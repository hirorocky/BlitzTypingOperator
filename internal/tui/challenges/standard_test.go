package challenges

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestStandard_正しい入力で進捗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"abc"},
	}

	c := newStandardChallengeForTest(input, 42)

	// 'a' を入力
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	sc := c.(*standardChallenge)
	if sc.currentIndex != 1 {
		t.Errorf("currentIndex = %d, want 1", sc.currentIndex)
	}

	// 'b' を入力
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	sc = c.(*standardChallenge)
	if sc.currentIndex != 2 {
		t.Errorf("currentIndex = %d, want 2", sc.currentIndex)
	}

	// 'c' を入力 → 完了
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
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

func TestStandard_誤入力でAccuracyが低下(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"ab"},
	}

	c := newStandardChallengeForTest(input, 42)

	// 'x' を入力（誤入力）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	sc := c.(*standardChallenge)
	if sc.currentIndex != 0 {
		t.Errorf("誤入力後のcurrentIndex = %d, want 0", sc.currentIndex)
	}
	if sc.mistakeCount != 1 {
		t.Errorf("mistakeCount = %d, want 1", sc.mistakeCount)
	}

	// 正しく入力して完了
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	result := c.Result()
	if result == nil {
		t.Fatal("Result()がnilです")
	}
	// 正解2回 / 総入力3回 = 0.666...
	if result.Accuracy < 0.6 || result.Accuracy > 0.7 {
		t.Errorf("Accuracy = %f, want ~0.67", result.Accuracy)
	}
}

func TestStandard_ESCでキャンセル(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newStandardChallengeForTest(input, 42)

	// ESCキー
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	result := c.Result()
	if result == nil {
		t.Fatal("ESC後にResult()がnilです")
	}
	if result.Status != domain.ChallengeCancel {
		t.Errorf("Status = %d, want ChallengeCancel", result.Status)
	}
}

func TestStandard_タイムアウトで失敗(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newStandardChallengeForTest(input, 42)
	sc := c.(*standardChallenge)

	// 開始時刻を過去に設定してタイムアウトをシミュレート
	sc.startTime = time.Now().Add(-sc.timeLimit - time.Second)

	// tickMsgを送信
	c, _ = c.Update(standardTickMsg{})
	result := c.Result()
	if result == nil {
		t.Fatal("タイムアウト後にResult()がnilです")
	}
	if result.Status != domain.ChallengeFail {
		t.Errorf("Status = %d, want ChallengeFail", result.Status)
	}
}

func TestStandard_AutoCorrectでミス無視(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"ab"},
		AutoCorrectCount: 1,
	}

	c := newStandardChallengeForTest(input, 42)

	// 'x' を入力（誤入力だがAutoCorrectで無視）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	sc := c.(*standardChallenge)
	if sc.currentIndex != 1 {
		t.Errorf("AutoCorrect後のcurrentIndex = %d, want 1（ミス無視で進む）", sc.currentIndex)
	}
	if sc.autoCorrectRemaining != 0 {
		t.Errorf("autoCorrectRemaining = %d, want 0", sc.autoCorrectRemaining)
	}

	// 次の誤入力はAutoCorrect使い切り済みなので通常ミス
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	sc = c.(*standardChallenge)
	if sc.currentIndex != 1 {
		t.Errorf("AutoCorrect使い切り後のcurrentIndex = %d, want 1", sc.currentIndex)
	}
}

func TestStandard_タイムアウトと全文字入力が同時なら成功優先(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"a"},
	}

	c := newStandardChallengeForTest(input, 42)
	sc := c.(*standardChallenge)

	// タイムアウト直前に設定
	sc.startTime = time.Now().Add(-sc.timeLimit)

	// 正しい入力（完了＝成功が優先される）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	result := c.Result()
	if result == nil {
		t.Fatal("Result()がnilです")
	}
	if result.Status != domain.ChallengeSuccess {
		t.Errorf("Status = %d, want ChallengeSuccess（成功が優先されるべき）", result.Status)
	}
}

func TestStandard_MistakeTimeExtend(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:           domain.DifficultyRateStandard,
		Words:                []string{"abc"},
		MistakeTimeExtendSec: 2.0,
	}

	c := newStandardChallengeForTest(input, 42)
	sc := c.(*standardChallenge)
	originalTimeLimit := sc.timeLimit

	// 誤入力で時間延長
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	sc = c.(*standardChallenge)
	if sc.timeLimit != originalTimeLimit+2*time.Second {
		t.Errorf("timeLimit = %v, want %v", sc.timeLimit, originalTimeLimit+2*time.Second)
	}

	// 2回目の誤入力では延長しない（1回/チャレンジ）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	sc = c.(*standardChallenge)
	if sc.timeLimit != originalTimeLimit+2*time.Second {
		t.Errorf("2回目のtimeLimit = %v, want %v（2回目は延長しない）", sc.timeLimit, originalTimeLimit+2*time.Second)
	}
}

func TestStandard_RetryOnTimeout(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:               domain.DifficultyRateStandard,
		Words:                    []string{"hello"},
		RetryOnTimeout:           true,
		RetryTimeLimitMultiplier: 0.5,
	}

	c := newStandardChallengeForTest(input, 42)
	sc := c.(*standardChallenge)
	originalTimeLimit := sc.timeLimit

	// タイムアウトをシミュレート
	sc.startTime = time.Now().Add(-sc.timeLimit - time.Second)

	// tickで再挑戦
	c, _ = c.Update(standardTickMsg{})
	sc = c.(*standardChallenge)

	// resultはまだnil（再挑戦中）
	if c.Result() != nil {
		t.Error("再挑戦中はResult()がnilであるべきです")
	}

	// 再挑戦の制限時間は元の50%
	expectedTimeLimit := time.Duration(float64(originalTimeLimit) * 0.5)
	if sc.timeLimit != expectedTimeLimit {
		t.Errorf("再挑戦の制限時間 = %v, want %v", sc.timeLimit, expectedTimeLimit)
	}

	// 再挑戦は1回のみ
	if sc.retryOnTimeout {
		t.Error("再挑戦後はretryOnTimeoutがfalseになるべきです")
	}
}

func TestStandard_View出力(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newStandardChallengeForTest(input, 42)
	view := c.View()

	if view == "" {
		t.Error("View()が空です")
	}
}

func TestStandard_DifficultyRateに応じた単語選択(t *testing.T) {
	// 低難易度: 短い単語
	inputLow := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(60),
		Words:      []string{"ab", "abc", "longword12345"},
	}
	cLow := newStandardChallengeForTest(inputLow, 42)
	scLow := cLow.(*standardChallenge)
	if len(scLow.text) > 8 {
		t.Errorf("低難易度で長すぎるテキスト: len=%d, text=%s", len(scLow.text), scLow.text)
	}

	// 高難易度: 長い単語
	inputHigh := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(180),
		Words:      []string{"ab", "abc", "longword12345"},
	}
	cHigh := newStandardChallengeForTest(inputHigh, 42)
	scHigh := cHigh.(*standardChallenge)
	if len(scHigh.text) < 3 {
		t.Errorf("高難易度でテキストが選択されていません: text=%s", scHigh.text)
	}
}
