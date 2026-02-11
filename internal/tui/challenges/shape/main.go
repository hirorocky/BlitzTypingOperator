package shape

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges"
	"hirorocky/type-battle/internal/tui/challenges/commons"

	tea "github.com/charmbracelet/bubbletea"
)

// shapeTickMsg はシェイプチャレンジ専用のtickメッセージです。
type shapeTickMsg struct{}

// shapeChallenge はシェイプタイプ（魔法攻撃スキル向け）のチャレンジです。
// 共通文字セットと形状テンプレートを使用してASCIIアートパターンを表示します。
type shapeChallenge struct {
	input  domain.ChallengeInput
	result *domain.ChallengeOutput

	// テキストと進捗
	text         string
	pattern      string // ASCIIアート表示用パターン
	currentIndex int

	// 入力統計
	correctCount     int
	totalInputCount  int
	mistakeCount     int
	lastMistake      bool
	mistakePositions map[int]bool

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
	challenges.Register(domain.ChallengeTypeShape, newShapeChallenge)
}

func newShapeChallenge(input domain.ChallengeInput) challenges.ChallengeModel {
	return newShapeChallengeWithRng(input, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// newShapeChallengeForTest はテスト用にシードを指定してチャレンジを生成します。
func newShapeChallengeForTest(input domain.ChallengeInput, seed int64) challenges.ChallengeModel {
	return newShapeChallengeWithRng(input, rand.New(rand.NewSource(seed)))
}

func newShapeChallengeWithRng(input domain.ChallengeInput, rng *rand.Rand) challenges.ChallengeModel {
	diffRate := int(input.Difficulty.Clamp())
	timeLimitMS := config.GetTimeLimitForRate(diffRate)
	timeLimit := time.Duration(timeLimitMS)*time.Millisecond + time.Duration(input.TimeExtendSec*float64(time.Second))

	// 共通文字セットから文字を生成
	chars := commons.GenerateChars(diffRate, rng)

	// オプションから形状名を取得（デフォルト: flame）
	shapeName := "flame"
	if input.ChallengeOptions != nil {
		if s, ok := input.ChallengeOptions["shape"]; ok && s != "" {
			shapeName = s
		}
	}

	// テンプレート選択とパターン生成
	tmpl := selectTemplate(shapeName, len(chars))
	pattern := commons.FormatAsPattern(tmpl, chars)
	text := string(chars)

	return &shapeChallenge{
		input:                    input,
		text:                     text,
		pattern:                  pattern,
		mistakePositions:         make(map[int]bool),
		startTime:                time.Now(),
		timeLimit:                timeLimit,
		autoCorrectRemaining:     input.AutoCorrectCount,
		mistakeTimeExtendSec:     input.MistakeTimeExtendSec,
		retryOnTimeout:           input.RetryOnTimeout,
		retryTimeLimitMultiplier: input.RetryTimeLimitMultiplier,
		rng:                      rng,
	}
}

// selectTemplate は形状名と文字数に応じたテンプレートを返します。
func selectTemplate(shapeName string, charCount int) string {
	switch shapeName {
	case "flame":
		return selectFlameTemplate(charCount)
	default:
		// 未知の形状はflameにフォールバック
		return selectFlameTemplate(charCount)
	}
}

func (c *shapeChallenge) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return shapeTickMsg{}
	})
}

func (c *shapeChallenge) Update(msg tea.Msg) (challenges.ChallengeModel, tea.Cmd) {
	if c.result != nil {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyInput(msg)
	case shapeTickMsg:
		return c.handleTick()
	}

	return c, nil
}

func (c *shapeChallenge) handleKeyInput(msg tea.KeyMsg) (challenges.ChallengeModel, tea.Cmd) {
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

func (c *shapeChallenge) processCharInput(input rune) (challenges.ChallengeModel, tea.Cmd) {
	if c.currentIndex >= len(c.text) {
		return c, nil
	}

	expected := rune(c.text[c.currentIndex])
	c.totalInputCount++

	if input == expected {
		c.correctCount++
		c.currentIndex++
		c.lastMistake = false
	} else {
		if c.autoCorrectRemaining > 0 {
			c.autoCorrectRemaining--
			c.correctCount++
			c.currentIndex++
			c.lastMistake = false
		} else {
			c.mistakeCount++
			c.lastMistake = true
			c.mistakePositions[c.currentIndex] = true

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

func (c *shapeChallenge) handleTick() (challenges.ChallengeModel, tea.Cmd) {
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
				return shapeTickMsg{}
			})
		}

		c.complete(domain.ChallengeFail)
		return c, nil
	}

	return c, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return shapeTickMsg{}
	})
}

func (c *shapeChallenge) complete(status domain.ChallengeStatus) {
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

func (c *shapeChallenge) View() string {
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
			if c.mistakePositions[textIdx] {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorMissed, ch, challenges.ColorReset)
			} else {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorCorrect, ch, challenges.ColorReset)
			}
		} else if textIdx == c.currentIndex {
			if c.lastMistake {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorMissCursor, ch, challenges.ColorReset)
			} else {
				fmt.Fprintf(&b, "%s%c%s", challenges.ColorCursor, ch, challenges.ColorReset)
			}
		} else {
			fmt.Fprintf(&b, "%s%c%s", challenges.ColorUntyped, ch, challenges.ColorReset)
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

func (c *shapeChallenge) Result() *domain.ChallengeOutput {
	return c.result
}
