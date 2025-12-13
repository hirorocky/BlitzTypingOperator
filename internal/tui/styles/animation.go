// Package styles はTUIゲームのアニメーションとフィードバック機能を提供します。
// タイピング表示、ダメージアニメーション、強調メッセージなどを担当します。

package styles

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// MessageType は強調メッセージの種類を表す型です。
type MessageType int

const (
	// MessageTypeSuccess は成功メッセージ（緑）
	MessageTypeSuccess MessageType = iota
	// MessageTypeInfo は情報メッセージ（青）
	MessageTypeInfo
	// MessageTypeWarning は警告メッセージ（黄）
	MessageTypeWarning
	// MessageTypeError はエラーメッセージ（赤）
	MessageTypeError
)

// タイピング表示用カラー（Requirement 18.6）
var (
	// ColorTypingCompleted は完了済み文字の色（緑）
	ColorTypingCompleted = lipgloss.Color("#04B575")
	// ColorTypingCurrent は入力中文字の色（背景ハイライト）
	ColorTypingCurrent = lipgloss.Color("#FFFFFF")
	// ColorTypingRemaining は未入力文字の色（薄グレー）
	ColorTypingRemaining = lipgloss.Color("#AAAAAA")
	// ColorTypingIncorrect は誤入力文字の色（赤）
	ColorTypingIncorrect = lipgloss.Color("#FF4672")
)

// TypingStyles はタイピング表示用スタイルを保持します。
type TypingStyles struct {
	Completed lipgloss.Style
	Current   lipgloss.Style
	Remaining lipgloss.Style
	Incorrect lipgloss.Style
}

// newTypingStyles はタイピング表示用スタイルを作成します。
func newTypingStyles() TypingStyles {
	return TypingStyles{
		Completed: lipgloss.NewStyle().
			Foreground(ColorTypingCompleted),
		Current: lipgloss.NewStyle().
			Background(ColorTypingCurrent).
			Foreground(lipgloss.Color("#000000")).
			Bold(true),
		Remaining: lipgloss.NewStyle().
			Foreground(ColorTypingRemaining),
		Incorrect: lipgloss.NewStyle().
			Foreground(ColorTypingIncorrect).
			Strikethrough(true),
	}
}

// RenderTypingCompleted は完了済み文字を描画します。

func (gs *GameStyles) RenderTypingCompleted(text string) string {
	ts := newTypingStyles()
	return ts.Completed.Render(text)
}

// RenderTypingCurrent は入力中文字を描画します。

func (gs *GameStyles) RenderTypingCurrent(text string) string {
	ts := newTypingStyles()
	return ts.Current.Render(text)
}

// RenderTypingRemaining は未入力文字を描画します。

func (gs *GameStyles) RenderTypingRemaining(text string) string {
	ts := newTypingStyles()
	return ts.Remaining.Render(text)
}

// RenderTypingIncorrect は誤入力文字を描画します。
func (gs *GameStyles) RenderTypingIncorrect(text string) string {
	ts := newTypingStyles()
	return ts.Incorrect.Render(text)
}

// RenderTypingChallenge はタイピングチャレンジ全体を描画します。

func (gs *GameStyles) RenderTypingChallenge(text string, currentIndex int, mistakes []int) string {
	if len(text) == 0 {
		return ""
	}

	ts := newTypingStyles()
	var result strings.Builder

	// 誤入力位置をマップ化
	mistakeMap := make(map[int]bool)
	for _, pos := range mistakes {
		mistakeMap[pos] = true
	}

	for i, char := range text {
		charStr := string(char)

		if i < currentIndex {
			// 完了済み（誤入力があった位置はマーク）
			if mistakeMap[i] {
				result.WriteString(ts.Incorrect.Render(charStr))
			} else {
				result.WriteString(ts.Completed.Render(charStr))
			}
		} else if i == currentIndex {
			// 入力中
			result.WriteString(ts.Current.Render(charStr))
		} else {
			// 未入力
			result.WriteString(ts.Remaining.Render(charStr))
		}
	}

	return result.String()
}

// GetDamageAnimationFrames はダメージアニメーションのフレームを返します。

func (gs *GameStyles) GetDamageAnimationFrames(damage int) []string {
	// シンプルなテキストアニメーション
	// 実際のアニメーションはUIレイヤーで時間経過で切り替える
	damageStyle := gs.Battle.Damage

	frames := []string{
		damageStyle.Render(fmt.Sprintf(" -%d ", damage)),
		damageStyle.Bold(true).Render(fmt.Sprintf("[ -%d ]", damage)),
		damageStyle.Render(fmt.Sprintf("< -%d >", damage)),
		damageStyle.Bold(true).Render(fmt.Sprintf("「-%d」", damage)),
	}

	return frames
}

// GetHealAnimationFrames は回復アニメーションのフレームを返します。
func (gs *GameStyles) GetHealAnimationFrames(heal int) []string {
	healStyle := gs.Battle.Heal

	frames := []string{
		healStyle.Render(fmt.Sprintf(" +%d ", heal)),
		healStyle.Bold(true).Render(fmt.Sprintf("[ +%d ]", heal)),
		healStyle.Render(fmt.Sprintf("< +%d >", heal)),
		healStyle.Bold(true).Render(fmt.Sprintf("「+%d」", heal)),
	}

	return frames
}

// RenderHighlightMessage は重要メッセージを強調表示します。

func (gs *GameStyles) RenderHighlightMessage(message string, msgType MessageType) string {
	var style lipgloss.Style

	switch msgType {
	case MessageTypeSuccess:
		style = lipgloss.NewStyle().
			Foreground(ColorHPHigh).
			Bold(true).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorHPHigh).
			Padding(0, 2)
	case MessageTypeInfo:
		style = lipgloss.NewStyle().
			Foreground(ColorInfo).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorInfo).
			Padding(0, 2)
	case MessageTypeWarning:
		style = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorWarning).
			Padding(0, 2)
	case MessageTypeError:
		style = lipgloss.NewStyle().
			Foreground(ColorDamage).
			Bold(true).
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorDamage).
			Padding(0, 2)
	default:
		style = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2)
	}

	return style.Render(message)
}

// RenderCooldownBar はクールダウンプログレスバーを描画します。

func (gs *GameStyles) RenderCooldownBar(remaining, total float64, width int) string {
	if total <= 0 {
		total = 1
	}
	if remaining < 0 {
		remaining = 0
	}
	if remaining > total {
		remaining = total
	}

	// クールダウンは「残り時間」なので、進捗は逆転
	// remaining=0 → 100%完了、remaining=total → 0%完了
	progress := 1.0 - (remaining / total)

	// バー内部の幅
	innerWidth := width - 2
	if innerWidth < 1 {
		innerWidth = 1
	}

	filledWidth := int(float64(innerWidth) * progress)
	emptyWidth := innerWidth - filledWidth

	// クールダウン完了（使用可能）は緑、クールダウン中はグレー
	var filledColor, emptyColor lipgloss.Color
	if progress >= 1.0 {
		filledColor = ColorHPHigh // 使用可能
		emptyColor = ColorSubtle
	} else {
		filledColor = ColorInfo // 回復中
		emptyColor = ColorSubtle
	}

	filledStyle := lipgloss.NewStyle().Background(filledColor)
	emptyStyle := lipgloss.NewStyle().Background(emptyColor)

	return "[" +
		filledStyle.Render(strings.Repeat(" ", filledWidth)) +
		emptyStyle.Render(strings.Repeat(" ", emptyWidth)) +
		"]"
}

// RenderCooldownBarWithTime はクールダウンバーと残り時間を描画します。

func (gs *GameStyles) RenderCooldownBarWithTime(remaining, total float64, width int) string {
	bar := gs.RenderCooldownBar(remaining, total, width)

	if remaining <= 0 {
		return fmt.Sprintf("%s %s", bar, gs.Text.Success.Render("READY"))
	}

	timeStr := gs.Battle.Cooldown.Render(fmt.Sprintf("%.1fs", remaining))
	return fmt.Sprintf("%s %s", bar, timeStr)
}

// AnimationState はアニメーション状態を管理する構造体です。

type AnimationState struct {
	// DamageAnimations は現在表示中のダメージアニメーション
	DamageAnimations []DamageAnimation
	// HealAnimations は現在表示中の回復アニメーション
	HealAnimations []HealAnimation
	// Messages は現在表示中のメッセージ
	Messages []TimedMessage
}

// DamageAnimation はダメージアニメーションの状態を表します。
type DamageAnimation struct {
	Amount      int
	FrameIndex  int
	RemainingMS int
	Position    Position
}

// HealAnimation は回復アニメーションの状態を表します。
type HealAnimation struct {
	Amount      int
	FrameIndex  int
	RemainingMS int
	Position    Position
}

// TimedMessage は時限付きメッセージを表します。
type TimedMessage struct {
	Text        string
	Type        MessageType
	RemainingMS int
}

// Position は画面上の位置を表します。
type Position struct {
	X int
	Y int
}

// NewAnimationState は新しいAnimationStateを作成します。
func NewAnimationState() *AnimationState {
	return &AnimationState{
		DamageAnimations: make([]DamageAnimation, 0),
		HealAnimations:   make([]HealAnimation, 0),
		Messages:         make([]TimedMessage, 0),
	}
}

// AddDamageAnimation はダメージアニメーションを追加します。
func (as *AnimationState) AddDamageAnimation(amount int, pos Position) {
	as.DamageAnimations = append(as.DamageAnimations, DamageAnimation{
		Amount:      amount,
		FrameIndex:  0,
		RemainingMS: 500, // 500msでアニメーション
		Position:    pos,
	})
}

// AddHealAnimation は回復アニメーションを追加します。
func (as *AnimationState) AddHealAnimation(amount int, pos Position) {
	as.HealAnimations = append(as.HealAnimations, HealAnimation{
		Amount:      amount,
		FrameIndex:  0,
		RemainingMS: 500,
		Position:    pos,
	})
}

// AddMessage はメッセージを追加します。
func (as *AnimationState) AddMessage(text string, msgType MessageType, durationMS int) {
	as.Messages = append(as.Messages, TimedMessage{
		Text:        text,
		Type:        msgType,
		RemainingMS: durationMS,
	})
}

// Update はアニメーション状態を更新します（deltaMS ミリ秒経過）。
func (as *AnimationState) Update(deltaMS int) {
	// ダメージアニメーションを更新
	activeDamage := make([]DamageAnimation, 0)
	for _, anim := range as.DamageAnimations {
		anim.RemainingMS -= deltaMS
		if anim.RemainingMS > 0 {
			// フレームインデックスを更新
			totalFrames := 4
			frameTime := 500 / totalFrames
			anim.FrameIndex = (500 - anim.RemainingMS) / frameTime
			if anim.FrameIndex >= totalFrames {
				anim.FrameIndex = totalFrames - 1
			}
			activeDamage = append(activeDamage, anim)
		}
	}
	as.DamageAnimations = activeDamage

	// 回復アニメーションを更新
	activeHeal := make([]HealAnimation, 0)
	for _, anim := range as.HealAnimations {
		anim.RemainingMS -= deltaMS
		if anim.RemainingMS > 0 {
			totalFrames := 4
			frameTime := 500 / totalFrames
			anim.FrameIndex = (500 - anim.RemainingMS) / frameTime
			if anim.FrameIndex >= totalFrames {
				anim.FrameIndex = totalFrames - 1
			}
			activeHeal = append(activeHeal, anim)
		}
	}
	as.HealAnimations = activeHeal

	// メッセージを更新
	activeMessages := make([]TimedMessage, 0)
	for _, msg := range as.Messages {
		msg.RemainingMS -= deltaMS
		if msg.RemainingMS > 0 {
			activeMessages = append(activeMessages, msg)
		}
	}
	as.Messages = activeMessages
}

// HasActiveAnimations はアクティブなアニメーションがあるかを返します。
func (as *AnimationState) HasActiveAnimations() bool {
	return len(as.DamageAnimations) > 0 || len(as.HealAnimations) > 0 || len(as.Messages) > 0
}

// ==================== AnimatedHPBar ====================

// AnimatedHPBar はアニメーション付きHPバーの状態を管理します。

type AnimatedHPBar struct {
	// CurrentDisplayHP は現在表示中のHP値（アニメーション用）
	CurrentDisplayHP float64

	// TargetHP は目標HP値
	TargetHP int

	// MaxHP は最大HP値
	MaxHP int

	// AnimationSpeed は1秒あたりのHP変化量
	AnimationSpeed float64

	// IsAnimating はアニメーション中かどうか
	IsAnimating bool
}

// デフォルトのアニメーション速度（1秒あたりのHP変化量）
const (
	// DefaultAnimationSpeed はデフォルトのアニメーション速度
	DefaultAnimationSpeed = 100.0
)

// NewAnimatedHPBar は新しいAnimatedHPBarを作成します。
func NewAnimatedHPBar(maxHP int) *AnimatedHPBar {
	return &AnimatedHPBar{
		CurrentDisplayHP: float64(maxHP),
		TargetHP:         maxHP,
		MaxHP:            maxHP,
		AnimationSpeed:   DefaultAnimationSpeed,
		IsAnimating:      false,
	}
}

// SetTarget は目標HP値を設定しアニメーションを開始します。
func (a *AnimatedHPBar) SetTarget(targetHP int) {
	// 境界値チェック
	if targetHP < 0 {
		targetHP = 0
	}
	if targetHP > a.MaxHP {
		targetHP = a.MaxHP
	}

	a.TargetHP = targetHP

	// 目標と現在表示が異なる場合はアニメーション開始
	if int(a.CurrentDisplayHP+0.5) != targetHP {
		a.IsAnimating = true
	}
}

// Update はアニメーションを更新します（deltaMS: 経過ミリ秒）。

func (a *AnimatedHPBar) Update(deltaMS int) {
	if !a.IsAnimating {
		return
	}

	// ミリ秒を秒に変換
	deltaSeconds := float64(deltaMS) / 1000.0

	// 変化量を計算
	change := a.AnimationSpeed * deltaSeconds

	targetFloat := float64(a.TargetHP)

	if a.CurrentDisplayHP > targetFloat {
		// ダメージ（減少）
		a.CurrentDisplayHP -= change
		if a.CurrentDisplayHP <= targetFloat {
			a.CurrentDisplayHP = targetFloat
			a.IsAnimating = false
		}
	} else if a.CurrentDisplayHP < targetFloat {
		// 回復（増加）
		a.CurrentDisplayHP += change
		if a.CurrentDisplayHP >= targetFloat {
			a.CurrentDisplayHP = targetFloat
			a.IsAnimating = false
		}
	} else {
		// 目標に到達
		a.IsAnimating = false
	}

	// 境界値制限
	if a.CurrentDisplayHP < 0 {
		a.CurrentDisplayHP = 0
	}
	if a.CurrentDisplayHP > float64(a.MaxHP) {
		a.CurrentDisplayHP = float64(a.MaxHP)
	}
}

// Render は現在の表示HPでHPバーを描画します。
func (a *AnimatedHPBar) Render(styles *GameStyles, width int) string {
	currentHP := a.GetCurrentHP()
	return styles.RenderHPBar(currentHP, a.MaxHP, width)
}

// GetCurrentHP は現在の表示HP（整数）を返します。
// 四捨五入して整数値を返します。
func (a *AnimatedHPBar) GetCurrentHP() int {
	return int(a.CurrentDisplayHP + 0.5)
}

// ForceComplete はアニメーションを強制完了し、表示HPを目標HPに即座に設定します。
func (a *AnimatedHPBar) ForceComplete() {
	a.CurrentDisplayHP = float64(a.TargetHP)
	a.IsAnimating = false
}

// ==================== FloatingDamageManager ====================

// FloatingText は浮遊テキストの状態を表します。

type FloatingText struct {
	Text        string
	IsHealing   bool   // true=回復（緑）、false=ダメージ（赤）
	RemainingMS int    // 残り表示時間（ミリ秒）
	YOffset     int    // Y方向オフセット（上方向に増加）
	TargetArea  string // "enemy", "player", "agent_{index}"
}

// FloatingDamageManager はフローティングダメージの状態を管理します。

type FloatingDamageManager struct {
	Texts []FloatingText
}

// フローティングテキストの設定
const (
	// FloatingTextDurationMS は表示時間（ミリ秒）

	FloatingTextDurationMS = 2500

	// FloatingTextRiseSpeed はY方向の移動速度（1秒あたりのピクセル数）

	FloatingTextRiseSpeed = 2
)

// NewFloatingDamageManager は新しいFloatingDamageManagerを作成します。
func NewFloatingDamageManager() *FloatingDamageManager {
	return &FloatingDamageManager{
		Texts: make([]FloatingText, 0),
	}
}

// AddDamage はダメージ表示を追加します。

func (m *FloatingDamageManager) AddDamage(amount int, targetArea string) {
	m.Texts = append(m.Texts, FloatingText{
		Text:        fmt.Sprintf("-%d", amount),
		IsHealing:   false,
		RemainingMS: FloatingTextDurationMS,
		YOffset:     0,
		TargetArea:  targetArea,
	})
}

// AddHeal は回復表示を追加します。

func (m *FloatingDamageManager) AddHeal(amount int, targetArea string) {
	m.Texts = append(m.Texts, FloatingText{
		Text:        fmt.Sprintf("+%d", amount),
		IsHealing:   true,
		RemainingMS: FloatingTextDurationMS,
		YOffset:     0,
		TargetArea:  targetArea,
	})
}

// Update は状態を更新します（deltaMS: 経過ミリ秒）。

func (m *FloatingDamageManager) Update(deltaMS int) {
	// 経過時間を秒に変換
	deltaSeconds := float64(deltaMS) / 1000.0

	activeTexts := make([]FloatingText, 0, len(m.Texts))
	for _, text := range m.Texts {
		text.RemainingMS -= deltaMS

		if text.RemainingMS > 0 {
			// Y方向への浮遊（上に移動）
			text.YOffset += int(float64(FloatingTextRiseSpeed) * deltaSeconds)
			activeTexts = append(activeTexts, text)
		}
	}
	m.Texts = activeTexts
}

// GetTextsForArea は指定エリアの表示テキストを取得します。

func (m *FloatingDamageManager) GetTextsForArea(targetArea string) []FloatingText {
	result := make([]FloatingText, 0)
	for _, text := range m.Texts {
		if text.TargetArea == targetArea {
			result = append(result, text)
		}
	}
	return result
}

// HasActiveTexts はアクティブな表示があるかを返します。
func (m *FloatingDamageManager) HasActiveTexts() bool {
	return len(m.Texts) > 0
}

// RenderFloatingText はフローティングテキストをスタイル付きで描画します。
func (m *FloatingDamageManager) RenderFloatingText(text FloatingText, styles *GameStyles) string {
	if text.IsHealing {
		return styles.Battle.Heal.Render(text.Text)
	}
	return styles.Battle.Damage.Render(text.Text)
}
