package standard

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hirorocky/type-battle/internal/config"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges"

	tea "github.com/charmbracelet/bubbletea"
)

// standardTickMsg はスタンダードチャレンジ専用のtickメッセージです。
type standardTickMsg struct{}

// standardChallenge はスタンダードタイプ（物理攻撃スキル向け）のチャレンジです。
type standardChallenge struct {
	input  domain.ChallengeInput
	result *domain.ChallengeOutput

	// テキストと進捗
	text         string
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
	challenges.Register(domain.ChallengeTypeStandard, newStandardChallenge)
}

func newStandardChallenge(input domain.ChallengeInput) challenges.ChallengeModel {
	return newStandardChallengeWithRng(input, rand.New(rand.NewSource(time.Now().UnixNano())))
}

// newStandardChallengeForTest はテスト用にシードを指定してチャレンジを生成します。
func newStandardChallengeForTest(input domain.ChallengeInput, seed int64) challenges.ChallengeModel {
	return newStandardChallengeWithRng(input, rand.New(rand.NewSource(seed)))
}

func newStandardChallengeWithRng(input domain.ChallengeInput, rng *rand.Rand) challenges.ChallengeModel {
	diffRate := int(input.Difficulty.Clamp())
	minLen, maxLen := config.GetTextLengthForRate(diffRate)
	timeLimitMS := config.GetTimeLimitForRate(diffRate)

	// 時間延長を適用
	timeLimit := time.Duration(timeLimitMS)*time.Millisecond + time.Duration(input.TimeExtendSec*float64(time.Second))

	// 単語選択
	text := selectWord(input.Words, minLen, maxLen, rng)

	return &standardChallenge{
		input:                    input,
		text:                     text,
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

func (c *standardChallenge) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return standardTickMsg{}
	})
}

func (c *standardChallenge) Update(msg tea.Msg) (challenges.ChallengeModel, tea.Cmd) {
	if c.result != nil {
		return c, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return c.handleKeyInput(msg)
	case standardTickMsg:
		return c.handleTick()
	}

	return c, nil
}

func (c *standardChallenge) handleKeyInput(msg tea.KeyMsg) (challenges.ChallengeModel, tea.Cmd) {
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

func (c *standardChallenge) processCharInput(input rune) (challenges.ChallengeModel, tea.Cmd) {
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
		// AutoCorrect: ミスを無視して次に進む
		if c.autoCorrectRemaining > 0 {
			c.autoCorrectRemaining--
			c.correctCount++ // AutoCorrectは正解扱い
			c.currentIndex++
			c.lastMistake = false
		} else {
			c.mistakeCount++
			c.lastMistake = true
			c.mistakePositions[c.currentIndex] = true

			// MistakeTimeExtend: 初回ミス時に時間延長（1回/チャレンジ）
			if c.mistakeTimeExtendSec > 0 && !c.mistakeTimeExtendUsed {
				c.timeLimit += time.Duration(c.mistakeTimeExtendSec * float64(time.Second))
				c.mistakeTimeExtendUsed = true
			}
		}
	}

	// 全文字入力完了チェック（タイムアウトと同フレーム時に成功を優先）
	if c.currentIndex >= len(c.text) {
		if c.mistakeCount == 0 {
			c.complete(domain.ChallengePerfect)
		} else {
			c.complete(domain.ChallengeSuccess)
		}
		return c, nil
	}

	return c, nil
}

func (c *standardChallenge) handleTick() (challenges.ChallengeModel, tea.Cmd) {
	elapsed := time.Since(c.startTime)

	if elapsed >= c.timeLimit {
		// RetryOnTimeout: 再挑戦
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
				return standardTickMsg{}
			})
		}

		c.complete(domain.ChallengeFail)
		return c, nil
	}

	return c, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return standardTickMsg{}
	})
}

func (c *standardChallenge) complete(status domain.ChallengeStatus) {
	completionTime := time.Since(c.startTime)

	accuracy := 1.0
	if c.totalInputCount > 0 {
		accuracy = float64(c.correctCount) / float64(c.totalInputCount)
	}

	wpm := 0.0
	if completionTime.Seconds() > 0 {
		wpm = (float64(c.correctCount) / completionTime.Seconds() * 60) / 5
	}

	score := int(accuracy * 100)

	c.result = &domain.ChallengeOutput{
		Score:          score,
		WPM:            wpm,
		CompletionTime: completionTime,
		Status:         status,
	}
}

func (c *standardChallenge) View() string {
	if c.result != nil {
		return ""
	}

	var b strings.Builder

	// テキスト表示: 入力済み / 現在位置 / 未入力
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

	// AutoCorrect残り
	if c.autoCorrectRemaining > 0 {
		fmt.Fprintf(&b, "ミス無視: %d回", c.autoCorrectRemaining)
		b.WriteString("\n")
	}

	return b.String()
}

func (c *standardChallenge) Result() *domain.ChallengeOutput {
	return c.result
}

// selectWord は辞書から指定文字数範囲の単語を選択します。
// 候補がなければ最も近い長さの単語を選択します。
func selectWord(words []string, minLen, maxLen int, rng *rand.Rand) string {
	if len(words) == 0 {
		return "default"
	}

	// 範囲内の候補を収集
	var candidates []string
	for _, w := range words {
		if len(w) >= minLen && len(w) <= maxLen {
			candidates = append(candidates, w)
		}
	}

	if len(candidates) > 0 {
		return candidates[rng.Intn(len(candidates))]
	}

	// 候補がなければ全単語から選択
	return words[rng.Intn(len(words))]
}
