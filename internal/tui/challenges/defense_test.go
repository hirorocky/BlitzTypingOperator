package challenges

import (
	"testing"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

func TestDefense_正しい入力で防御率が上昇(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"abc"},
	}

	c := newDefenseChallengeForTest(input, 42)
	dp := c.(DefenseProvider)

	// 初期防御率は0
	if dp.DefenseRate() != 0.0 {
		t.Errorf("初期防御率 = %f, want 0.0", dp.DefenseRate())
	}

	// 'a'を入力
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	dp = c.(DefenseProvider)
	if dp.DefenseRate() <= 0.0 {
		t.Errorf("入力後の防御率 = %f, want > 0.0", dp.DefenseRate())
	}
}

func TestDefense_全文字入力で次の単語が生成(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"ab", "cd"},
	}

	c := newDefenseChallengeForTest(input, 42)
	dc := c.(*defenseChallenge)
	firstText := dc.text

	// 全文字入力
	for _, ch := range firstText {
		c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
	}

	dc = c.(*defenseChallenge)
	// 次の単語が生成されている
	if dc.currentIndex != 0 {
		t.Errorf("次の単語生成後のcurrentIndex = %d, want 0", dc.currentIndex)
	}
	// チャレンジは終了していない（制限時間がないため）
	if c.Result() != nil {
		t.Error("全文字入力後にResult()がnon-nil（ディフェンスタイプは敵攻撃まで終了しない）")
	}
}

func TestDefense_ESCでキャンセル(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newDefenseChallengeForTest(input, 42)

	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyEsc})
	result := c.Result()
	if result == nil {
		t.Fatal("ESC後にResult()がnilです")
	}
	if result.Status != domain.ChallengeCancel {
		t.Errorf("Status = %d, want ChallengeCancel", result.Status)
	}
}

func TestDefense_CompleteByAttackで自動終了(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"abcde"},
	}

	c := newDefenseChallengeForTest(input, 42)
	_ = c.(DefenseProvider)

	// いくつか文字を入力して防御率を上げる
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	dp := c.(DefenseProvider)
	defenseRate := dp.DefenseRate()

	// 敵攻撃で自動終了
	dp.CompleteByAttack()

	result := c.Result()
	if result == nil {
		t.Fatal("CompleteByAttack後にResult()がnilです")
	}
	if result.Status != domain.ChallengeSuccess {
		t.Errorf("Status = %d, want ChallengeSuccess", result.Status)
	}
	if result.Accuracy != defenseRate {
		t.Errorf("Accuracy = %f, want %f（防御率がAccuracyに反映されるべき）", result.Accuracy, defenseRate)
	}
}

func TestDefense_AutoCorrectでも防御率が上昇(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"abc"},
		AutoCorrectCount: 1,
	}

	c := newDefenseChallengeForTest(input, 42)
	_ = c.(DefenseProvider)

	// 誤入力（AutoCorrectで無視→防御率上昇）
	c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	dp := c.(DefenseProvider)
	if dp.DefenseRate() <= 0.0 {
		t.Errorf("AutoCorrect後の防御率 = %f, want > 0.0", dp.DefenseRate())
	}
}

func TestDefense_防御率は最大1_0(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(50), // 低難易度 = 上昇量が大きい
		Words:      []string{"ab", "cd", "ef", "gh", "ij"},
	}

	c := newDefenseChallengeForTest(input, 42)
	dp := c.(DefenseProvider)

	// 大量に入力
	for i := 0; i < 100; i++ {
		dc := c.(*defenseChallenge)
		if dc.currentIndex < len(dc.text) {
			ch := rune(dc.text[dc.currentIndex])
			c, _ = c.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}})
			dp = c.(DefenseProvider)
		}
	}

	if dp.DefenseRate() > 1.0 {
		t.Errorf("防御率 = %f, want <= 1.0", dp.DefenseRate())
	}
}

func TestDefense_DifficultyRateで上昇量が変動(t *testing.T) {
	// 低難易度: 1文字あたりの上昇量が大きい
	inputLow := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(60),
		Words:      []string{"abcde"},
	}
	cLow := newDefenseChallengeForTest(inputLow, 42)

	// 高難易度: 1文字あたりの上昇量が小さい
	inputHigh := domain.ChallengeInput{
		Difficulty: domain.DifficultyRate(180),
		Words:      []string{"abcde"},
	}
	cHigh := newDefenseChallengeForTest(inputHigh, 42)

	// 1文字入力
	cLow, _ = cLow.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	cHigh, _ = cHigh.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	dpLow := cLow.(DefenseProvider)
	dpHigh := cHigh.(DefenseProvider)

	if dpLow.DefenseRate() <= dpHigh.DefenseRate() {
		t.Errorf("低難易度(%f) <= 高難易度(%f)、低難易度の方が上昇量が大きいはず",
			dpLow.DefenseRate(), dpHigh.DefenseRate())
	}
}

func TestDefense_View出力(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newDefenseChallengeForTest(input, 42)
	view := c.View()

	if view == "" {
		t.Error("View()が空です")
	}
}

func TestDefense_制限時間なし(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty: domain.DifficultyRateStandard,
		Words:      []string{"hello"},
	}

	c := newDefenseChallengeForTest(input, 42)

	// Init()はtea.Tickを返さない（制限時間がないため）
	cmd := c.Init()
	if cmd != nil {
		t.Error("ディフェンスタイプのInit()はnilを返すべき（制限時間なし）")
	}
}
