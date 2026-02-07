// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"fmt"
	"time"
)

// ChallengeTypeID はチャレンジ種別の識別子です。
type ChallengeTypeID string

const (
	// ChallengeTypeStandard は逐次入力型のスタンダードタイプです（物理攻撃スキル向け）。
	ChallengeTypeStandard ChallengeTypeID = "standard"

	// ChallengeTypeSymbolStorm は記号パターン入力型のシンボルストームタイプです（魔法攻撃スキル向け）。
	ChallengeTypeSymbolStorm ChallengeTypeID = "symbol_storm"

	// ChallengeTypeDefense はリアルタイム防御型のディフェンスタイプです（防御スキル向け）。
	ChallengeTypeDefense ChallengeTypeID = "defense"
)

// DifficultyRate はタイピング難易度を表す値です（50-200、100=標準）。
type DifficultyRate int

const (
	// DifficultyRateMin は難易度の最小値です。
	DifficultyRateMin DifficultyRate = 50

	// DifficultyRateMax は難易度の最大値です。
	DifficultyRateMax DifficultyRate = 200

	// DifficultyRateStandard は標準難易度です。
	DifficultyRateStandard DifficultyRate = 100
)

// Clamp は難易度を有効範囲内に制限します。
func (d DifficultyRate) Clamp() DifficultyRate {
	if d < DifficultyRateMin {
		return DifficultyRateMin
	}
	if d > DifficultyRateMax {
		return DifficultyRateMax
	}
	return d
}

// ChallengeStatus はチャレンジの結果状態を表します。
type ChallengeStatus int

const (
	// ChallengeSuccess は全文字入力完了（成功）です。
	ChallengeSuccess ChallengeStatus = iota

	// ChallengeFail は制限時間超過（失敗）です。
	ChallengeFail

	// ChallengeCancel はESCキーによるキャンセルです。
	ChallengeCancel
)

// ChallengeInput はチャレンジへの入力パラメータです。
type ChallengeInput struct {
	// Difficulty は難易度（50-200、100=標準）です。
	Difficulty DifficultyRate

	// Words は辞書（optional）です。
	Words []string

	// TimeExtendSec はバフによる時間延長秒数です。ディフェンスタイプでは無視されます。
	TimeExtendSec float64

	// AutoCorrectCount はバフによるミス無視回数です。全タイプで有効です。
	AutoCorrectCount int

	// MistakeTimeExtendSec はps_typo_recovery: ミス時の時間延長秒数です。0で無効。ディフェンスタイプでは無視されます。
	MistakeTimeExtendSec float64

	// RetryOnTimeout はps_second_chance: タイムアウト時に再挑戦を許可します。ディフェンスタイプでは無視されます。
	RetryOnTimeout bool

	// RetryTimeLimitMultiplier はps_second_chance: 再挑戦時の制限時間倍率です（例: 0.5で半分）。
	RetryTimeLimitMultiplier float64
}

// Validate は入力値のバリデーションを行います。
func (i ChallengeInput) Validate() error {
	clamped := i.Difficulty.Clamp()
	if clamped != i.Difficulty {
		return fmt.Errorf("DifficultyRateが範囲外です: %d (有効範囲: %d-%d)", i.Difficulty, DifficultyRateMin, DifficultyRateMax)
	}
	if i.AutoCorrectCount < 0 {
		return fmt.Errorf("AutoCorrectCountが負値です: %d", i.AutoCorrectCount)
	}
	return nil
}

// ChallengeOutput はチャレンジの出力結果です。
type ChallengeOutput struct {
	// Accuracy は正確性（0.0-1.0）です。ディフェンスタイプでは最終防御率を表します。
	Accuracy float64

	// SpeedFactor は速度係数（上限2.0）です。
	SpeedFactor float64

	// WPM はWords Per Minuteです（パッシブスキル判定用）。
	WPM float64

	// CompletionTime はチャレンジ完了までの時間です（統計用）。
	CompletionTime time.Duration

	// Status はチャレンジの結果状態です。
	Status ChallengeStatus
}
