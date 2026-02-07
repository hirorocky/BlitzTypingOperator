package challenges

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// symbolStormTickMsg はシンボルストームチャレンジ専用のtickメッセージです。
type symbolStormTickMsg struct{}

// 基本記号（Shift不要）
var basicSymbols = []rune{'<', '>', '(', ')', '[', ']', '{', '}', '-', '+', '=', '.', ',', '/', '\\', '|', '_'}

// Shift+数字記号（高難易度で割合が増える）
var shiftSymbols = []rune{'!', '@', '#', '$', '%', '^', '&', '*'}

// symbolStormChallenge はシンボルストームタイプ（魔法攻撃スキル向け）のチャレンジです。
type symbolStormChallenge struct {
	input  domain.ChallengeInput
	result *domain.ChallengeOutput

	// テキストと進捗
	text         string
	pattern      string // ASCIIアート表示用パターン
	currentIndex int

	// 入力統計
	correctCount    int
	totalInputCount int
	mistakeCount    int

	// 時間管理
	startTime time.Time
	timeLimit time.Duration

	// AutoCorrect
	autoCorrectRemaining int

	// MistakeTimeExtend（1回/チャレンジ）
	mistakeTimeExtendSec  float64
	mistakeTimeExtendUsed bool

	// RetryOnTimeout
	retryOnTimeout           bool
	retryTimeLimitMultiplier float64

	// 乱数生成器
	rng *rand.Rand
}

func init() {
	Register(domain.ChallengeTypeSymbolStorm, newSymbolStormChallenge)
}

func newSymbolStormChallenge(input domain.ChallengeInput) ChallengeModel {
	return newSymbolStormChallengeWithRng(input, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// newSymbolStormChallengeForTest はテスト用にシードを指定してチャレンジを生成します。
func newSymbolStormChallengeForTest(input domain.ChallengeInput, seed int64) ChallengeModel {
	return newSymbolStormChallengeWithRng(input, rand.New(rand.NewSource(seed)))
}

func newSymbolStormChallengeWithRng(input domain.ChallengeInput, rng *rand.Rand) ChallengeModel {
	diffRate := int(input.Difficulty.Clamp())
	timeLimitMS := config.GetTimeLimitForRate(diffRate)
	timeLimit := time.Duration(timeLimitMS)*time.Millisecond + time.Duration(input.TimeExtendSec*float64(time.Second))

	text, pattern := generateSymbolPattern(diffRate, rng)

	return &symbolStormChallenge{
		input:                    input,
		text:                     text,
		pattern:                  pattern,
		startTime:                time.Now(),
		timeLimit:                timeLimit,
		autoCorrectRemaining:     input.AutoCorrectCount,
		mistakeTimeExtendSec:     input.MistakeTimeExtendSec,
		retryOnTimeout:           input.RetryOnTimeout,
		retryTimeLimitMultiplier: input.RetryTimeLimitMultiplier,
		rng:                      rng,
	}
}

// generateSymbolPattern はDifficultyRateに応じた記号パターンを生成します。
func generateSymbolPattern(diffRate int, rng *rand.Rand) (text string, pattern string) {
	// 文字数を決定: rate=50→4-6, rate=100→6-8, rate=200→10-16
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

	// Shift記号の割合: rate=50→10%, rate=100→25%, rate=200→50%
	shiftRatio := 0.1
	if diffRate > 50 {
		shiftRatio = 0.1 + float64(diffRate-50)/150.0*0.4
	}
	if shiftRatio > 0.5 {
		shiftRatio = 0.5
	}

	// 記号列を生成
	chars := make([]rune, length)
	for i := range chars {
		if rng.Float64() < shiftRatio {
			chars[i] = shiftSymbols[rng.Intn(len(shiftSymbols))]
		} else {
			chars[i] = basicSymbols[rng.Intn(len(basicSymbols))]
		}
	}

	text = string(chars)
	pattern = formatAsPattern(chars)
	return
}

// formatAsPattern は記号列をASCIIアートパターン風に整形します。
func formatAsPattern(chars []rune) string {
	var b strings.Builder
	lineLen := 4
	if len(chars) > 10 {
		lineLen = 5
	}

	for i, ch := range chars {
		if i > 0 && i%lineLen == 0 {
			b.WriteRune('\n')
		}
		// 記号間にスペースを入れて見やすくする
		if i%lineLen != 0 {
			b.WriteRune(' ')
		}
		b.WriteRune(ch)
	}
	return b.String()
}

func (c *symbolStormChallenge) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return symbolStormTickMsg{}
	})
}

func (c *symbolStormChallenge) Update(msg tea.Msg) (ChallengeModel, tea.Cmd) {
	if c.result != nil {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyInput(msg)
	case symbolStormTickMsg:
		return c.handleTick()
	}

	return c, nil
}

func (c *symbolStormChallenge) handleKeyInput(msg tea.KeyMsg) (ChallengeModel, tea.Cmd) {
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
	}

	return c, nil
}

func (c *symbolStormChallenge) processCharInput(input rune) (ChallengeModel, tea.Cmd) {
	if c.currentIndex >= len(c.text) {
		return c, nil
	}

	expected := rune(c.text[c.currentIndex])
	c.totalInputCount++

	if input == expected {
		c.correctCount++
		c.currentIndex++
	} else {
		if c.autoCorrectRemaining > 0 {
			c.autoCorrectRemaining--
			c.correctCount++
			c.currentIndex++
		} else {
			c.mistakeCount++

			if c.mistakeTimeExtendSec > 0 && !c.mistakeTimeExtendUsed {
				c.timeLimit += time.Duration(c.mistakeTimeExtendSec * float64(time.Second))
				c.mistakeTimeExtendUsed = true
			}
		}
	}

	if c.currentIndex >= len(c.text) {
		c.complete(domain.ChallengeSuccess)
		return c, nil
	}

	return c, nil
}

func (c *symbolStormChallenge) handleTick() (ChallengeModel, tea.Cmd) {
	elapsed := time.Since(c.startTime)

	if elapsed >= c.timeLimit {
		if c.retryOnTimeout {
			c.retryOnTimeout = false
			c.currentIndex = 0
			c.correctCount = 0
			c.totalInputCount = 0
			c.mistakeCount = 0
			c.timeLimit = time.Duration(float64(c.timeLimit) * c.retryTimeLimitMultiplier)
			c.startTime = time.Now()
			c.mistakeTimeExtendUsed = false

			return c, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return symbolStormTickMsg{}
			})
		}

		c.complete(domain.ChallengeFail)
		return c, nil
	}

	return c, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return symbolStormTickMsg{}
	})
}

func (c *symbolStormChallenge) complete(status domain.ChallengeStatus) {
	completionTime := time.Since(c.startTime)

	accuracy := 1.0
	if c.totalInputCount > 0 {
		accuracy = float64(c.correctCount) / float64(c.totalInputCount)
	}

	wpm := 0.0
	if completionTime.Seconds() > 0 {
		wpm = (float64(c.correctCount) / completionTime.Seconds() * 60) / 5
	}

	speedFactor := 1.0
	if completionTime.Seconds() > 0 && status == domain.ChallengeSuccess {
		speedFactor = c.timeLimit.Seconds() / completionTime.Seconds()
		if speedFactor > 2.0 {
			speedFactor = 2.0
		}
	} else if status != domain.ChallengeSuccess {
		speedFactor = 0
	}

	c.result = &domain.ChallengeOutput{
		Accuracy:       accuracy,
		SpeedFactor:    speedFactor,
		WPM:            wpm,
		CompletionTime: completionTime,
		Status:         status,
	}
}

func (c *symbolStormChallenge) View() string {
	if c.result != nil {
		return ""
	}

	var b strings.Builder

	// パターン表示（入力済みの文字は色分け）
	textIdx := 0
	for _, ch := range c.pattern {
		if ch == '\n' || ch == ' ' {
			b.WriteRune(ch)
			continue
		}

		if textIdx < c.currentIndex {
			fmt.Fprintf(&b, "\033[35m%c\033[0m", ch) // マゼンタ: 入力済み
		} else if textIdx == c.currentIndex {
			fmt.Fprintf(&b, "\033[7;35m%c\033[0m", ch) // 反転マゼンタ: 現在位置
		} else {
			fmt.Fprintf(&b, "\033[90m%c\033[0m", ch) // グレー: 未入力
		}
		textIdx++
	}
	b.WriteString("\n")

	// 残り時間バー
	remaining := c.timeLimit - time.Since(c.startTime)
	if remaining < 0 {
		remaining = 0
	}
	progress := float64(remaining) / float64(c.timeLimit)
	barWidth := 20
	filled := int(progress * float64(barWidth))
	b.WriteString("[")
	b.WriteString(strings.Repeat("█", filled))
	b.WriteString(strings.Repeat("░", barWidth-filled))
	fmt.Fprintf(&b, "] %.1fs", remaining.Seconds())
	b.WriteString("\n")

	if c.autoCorrectRemaining > 0 {
		fmt.Fprintf(&b, "ミス無視: %d回", c.autoCorrectRemaining)
		b.WriteString("\n")
	}

	return b.String()
}

func (c *symbolStormChallenge) Result() *domain.ChallengeOutput {
	return c.result
}
