package defense

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges"

	tea "github.com/charmbracelet/bubbletea"
)

// defenseChallenge はディフェンスタイプ（防御スキル向け）のチャレンジです。
// DefenseProviderインターフェースを実装し、リアルタイムの防御率を提供します。
type defenseChallenge struct {
	input  domain.ChallengeInput
	result *domain.ChallengeOutput

	// テキストと進捗
	text         string
	currentIndex int

	// 防御率（0.0-1.0）
	defenseRate  float64
	ratePerChar  float64 // 1文字あたりの防御率上昇量
	totalCorrect     int // 正解入力数の合計（統計用）
	totalInputs      int // 総入力数（統計用）
	lastMistake      bool
	mistakePositions map[int]bool

	// AutoCorrect
	autoCorrectRemaining int

	// 開始時刻（統計用）
	startTime time.Time

	// 乱数生成器
	rng   *rand.Rand
	words []string
}

func init() {
	challenges.Register(domain.ChallengeTypeDefense, newDefenseChallenge)
}

func newDefenseChallenge(input domain.ChallengeInput) challenges.ChallengeModel {
	return newDefenseChallengeWithRng(input, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// newDefenseChallengeForTest はテスト用にシードを指定してチャレンジを生成します。
func newDefenseChallengeForTest(input domain.ChallengeInput, seed int64) challenges.ChallengeModel {
	return newDefenseChallengeWithRng(input, rand.New(rand.NewSource(seed)))
}

func newDefenseChallengeWithRng(input domain.ChallengeInput, rng *rand.Rand) challenges.ChallengeModel {
	diffRate := int(input.Difficulty.Clamp())

	// 1文字あたりの防御率上昇量: 低難易度は大きく、高難易度は小さい
	// rate=50 → 0.10/文字, rate=100 → 0.05/文字, rate=200 → 0.02/文字
	ratePerChar := calculateRatePerChar(diffRate)

	words := input.Words
	if len(words) == 0 {
		words = []string{"default"}
	}

	text := words[rng.Intn(len(words))]

	return &defenseChallenge{
		input:                input,
		text:                 text,
		ratePerChar:          ratePerChar,
		mistakePositions:     make(map[int]bool),
		autoCorrectRemaining: input.AutoCorrectCount,
		startTime:            time.Now(),
		rng:                  rng,
		words:                words,
	}
}

// calculateRatePerChar はDifficultyRateに応じた1文字あたりの防御率上昇量を計算します。
func calculateRatePerChar(diffRate int) float64 {
	if diffRate <= 50 {
		return 0.10
	}
	if diffRate >= 200 {
		return 0.02
	}

	// 50-200: 線形補間（0.10 → 0.02）
	t := float64(diffRate-50) / 150.0
	return 0.10 - t*(0.10-0.02)
}

func (c *defenseChallenge) Init() tea.Cmd {
	// ディフェンスタイプは制限時間がないためtea.Tick不要
	return nil
}

func (c *defenseChallenge) Update(msg tea.Msg) (challenges.ChallengeModel, tea.Cmd) {
	if c.result != nil {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyInput(msg)
	}

	return c, nil
}

func (c *defenseChallenge) handleKeyInput(msg tea.KeyMsg) (challenges.ChallengeModel, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		c.result = &domain.ChallengeOutput{
			Status:         domain.ChallengeCancel,
			CompletionTime: time.Since(c.startTime),
		}
		return c, nil

	case tea.KeyRunes:
		if len(msg.Runes) == 0 {
			return c, nil
		}
		return c.processCharInput(msg.Runes[0])

	case tea.KeySpace:
		return c.processCharInput(' ')
	}

	return c, nil
}

func (c *defenseChallenge) processCharInput(input rune) (challenges.ChallengeModel, tea.Cmd) {
	if c.currentIndex >= len(c.text) {
		return c, nil
	}

	expected := rune(c.text[c.currentIndex])
	c.totalInputs++

	if input == expected {
		c.totalCorrect++
		c.currentIndex++
		c.increaseDefenseRate()
		c.lastMistake = false
	} else {
		// AutoCorrect: ミスを無視して防御率も上昇する
		if c.autoCorrectRemaining > 0 {
			c.autoCorrectRemaining--
			c.totalCorrect++
			c.currentIndex++
			c.increaseDefenseRate()
			c.lastMistake = false
		} else {
			c.lastMistake = true
			c.mistakePositions[c.currentIndex] = true
		}
	}

	// 全文字入力完了で次の単語を生成（ディフェンスタイプは敵攻撃まで継続）
	if c.currentIndex >= len(c.text) {
		c.nextWord()
	}

	return c, nil
}

func (c *defenseChallenge) increaseDefenseRate() {
	c.defenseRate += c.ratePerChar
	if c.defenseRate > 1.0 {
		c.defenseRate = 1.0
	}
}

func (c *defenseChallenge) nextWord() {
	c.currentIndex = 0
	c.lastMistake = false
	c.mistakePositions = make(map[int]bool)
	c.text = c.words[c.rng.Intn(len(c.words))]
}

func (c *defenseChallenge) View() string {
	if c.result != nil {
		return ""
	}

	var b strings.Builder

	// 防御率バー
	barWidth := 20
	filled := int(c.defenseRate * float64(barWidth))
	b.WriteString("防御率 [")
	b.WriteString(strings.Repeat("█", filled))
	b.WriteString(strings.Repeat("░", barWidth-filled))
	fmt.Fprintf(&b, "] %.0f%%", c.defenseRate*100)
	b.WriteString("\n\n")

	// テキスト表示
	for i, ch := range c.text {
		if i < c.currentIndex {
			if c.mistakePositions[i] {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorMissed, ch, challenges.ColorReset)
			} else {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorCorrect, ch, challenges.ColorReset)
			}
		} else if i == c.currentIndex {
			if c.lastMistake {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorMissCursor, ch, challenges.ColorReset)
			} else {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorCursor, ch, challenges.ColorReset)
			}
		} else {
			fmt.Fprintf(&b, "%s%c%s", challenges.ColorUntyped, ch, challenges.ColorReset)
		}
	}
	b.WriteString("\n")

	if c.autoCorrectRemaining > 0 {
		fmt.Fprintf(&b, "ミス無視: %d回", c.autoCorrectRemaining)
		b.WriteString("\n")
	}

	return b.String()
}

func (c *defenseChallenge) Result() *domain.ChallengeOutput {
	return c.result
}

// DefenseRate は現在の防御率（0.0-1.0）を返します。
func (c *defenseChallenge) DefenseRate() float64 {
	return c.defenseRate
}

// CompleteByAttack は敵攻撃時にチャレンジを自動終了させます。
// 既に終了済み（Cancel含む）の場合はno-opです。
func (c *defenseChallenge) CompleteByAttack() {
	if c.result != nil {
		return
	}
	c.result = &domain.ChallengeOutput{
		Accuracy:       c.defenseRate,
		Status:         domain.ChallengeSuccess,
		CompletionTime: time.Since(c.startTime),
	}
}
