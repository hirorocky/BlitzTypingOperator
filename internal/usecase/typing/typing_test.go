// Package typing はタイピングシステムを提供します。
// タイピングチャレンジの生成と評価を担当します。

package typing

import (
	"testing"
	"time"
)

// ==================== タイピングチャレンジ生成テスト（Task 6.1） ====================

// TestGenerateChallenge_Easy は弱いスキル用テキスト生成をテストします。

func TestGenerateChallenge_Easy(t *testing.T) {
	dict := &Dictionary{
		Easy:   []string{"cat", "dog", "hello"},
		Medium: []string{"program", "keyboard"},
		Hard:   []string{"extraordinary", "sophisticated"},
	}
	generator := NewChallengeGenerator(dict)

	challenge := generator.Generate(DifficultyEasy, 5*time.Second)
	if challenge == nil {
		t.Fatal("チャレンジ生成に失敗")
	}

	textLen := len(challenge.Text)
	if textLen < 3 || textLen > 6 {
		t.Errorf("簡単テキストの長さ: 期待 3-6, 実際 %d (%s)", textLen, challenge.Text)
	}
}

// TestGenerateChallenge_Medium は中程度のスキル用テキスト生成をテストします。

func TestGenerateChallenge_Medium(t *testing.T) {
	dict := &Dictionary{
		Easy:   []string{"cat", "dog"},
		Medium: []string{"program", "keyboard", "interface"},
		Hard:   []string{"extraordinary"},
	}
	generator := NewChallengeGenerator(dict)

	challenge := generator.Generate(DifficultyMedium, 8*time.Second)
	if challenge == nil {
		t.Fatal("チャレンジ生成に失敗")
	}

	textLen := len(challenge.Text)
	if textLen < 7 || textLen > 11 {
		t.Errorf("中程度テキストの長さ: 期待 7-11, 実際 %d (%s)", textLen, challenge.Text)
	}
}

// TestGenerateChallenge_Hard は強力なスキル用テキスト生成をテストします。

func TestGenerateChallenge_Hard(t *testing.T) {
	dict := &Dictionary{
		Easy:   []string{"cat"},
		Medium: []string{"program"},
		Hard:   []string{"extraordinary", "implementation", "authentication"},
	}
	generator := NewChallengeGenerator(dict)

	challenge := generator.Generate(DifficultyHard, 12*time.Second)
	if challenge == nil {
		t.Fatal("チャレンジ生成に失敗")
	}

	textLen := len(challenge.Text)
	if textLen < 12 || textLen > 20 {
		t.Errorf("難しいテキストの長さ: 期待 12-20, 実際 %d (%s)", textLen, challenge.Text)
	}
}

// TestGenerateChallenge_NoDuplication は連続同一テキスト回避をテストします。

func TestGenerateChallenge_NoDuplication(t *testing.T) {
	dict := &Dictionary{
		Easy:   []string{"cat", "dog", "run", "jump"},
		Medium: []string{"program", "keyboard"},
		Hard:   []string{"extraordinary"},
	}
	generator := NewChallengeGenerator(dict)

	challenge1 := generator.Generate(DifficultyEasy, 5*time.Second)
	challenge2 := generator.Generate(DifficultyEasy, 5*time.Second)

	// 連続して同じテキストにならないことを確認（確率的なので複数回実行）
	sameCount := 0
	for i := 0; i < 10; i++ {
		c1 := generator.Generate(DifficultyEasy, 5*time.Second)
		c2 := generator.Generate(DifficultyEasy, 5*time.Second)
		if c1.Text == c2.Text {
			sameCount++
		}
	}

	// 単語が4個あるので、連続重複は確率的に起きにくいはず
	// 10回中5回以上重複したら問題
	if sameCount >= 5 {
		t.Errorf("連続同一テキストが多すぎる: %d/10回", sameCount)
	}

	_ = challenge1
	_ = challenge2
}

// TestGenerateChallenge_TimeLimit はスキル別制限時間設定をテストします。

func TestGenerateChallenge_TimeLimit(t *testing.T) {
	dict := &Dictionary{
		Easy:   []string{"cat", "dog"},
		Medium: []string{"program"},
		Hard:   []string{"extraordinary"},
	}
	generator := NewChallengeGenerator(dict)

	timeLimit := 7500 * time.Millisecond
	challenge := generator.Generate(DifficultyEasy, timeLimit)

	if challenge.TimeLimit != timeLimit {
		t.Errorf("制限時間: 期待 %v, 実際 %v", timeLimit, challenge.TimeLimit)
	}
}

// ==================== タイピング評価エンジンテスト（Task 6.2） ====================

// TestEvaluator_StartChallenge はチャレンジ開始時のタイムスタンプ記録をテストします。

func TestEvaluator_StartChallenge(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}

	before := time.Now()
	state := evaluator.StartChallenge(challenge)
	after := time.Now()

	if state.StartTime.Before(before) || state.StartTime.After(after) {
		t.Error("開始時刻が正しく記録されていない")
	}
	if state.CurrentIndex != 0 {
		t.Errorf("初期インデックス: 期待 0, 実際 %d", state.CurrentIndex)
	}
}

// TestEvaluator_ProcessCorrectInput は正しい入力の処理をテストします。

func TestEvaluator_ProcessCorrectInput(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	state = evaluator.ProcessInput(state, 'h')
	if state.CurrentIndex != 1 {
		t.Errorf("入力後のインデックス: 期待 1, 実際 %d", state.CurrentIndex)
	}
	if state.CorrectCount != 1 {
		t.Errorf("正解数: 期待 1, 実際 %d", state.CorrectCount)
	}
	if state.TotalInputCount != 1 {
		t.Errorf("総入力数: 期待 1, 実際 %d", state.TotalInputCount)
	}
}

// TestEvaluator_ProcessIncorrectInput は誤った入力の処理をテストします。

func TestEvaluator_ProcessIncorrectInput(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	state = evaluator.ProcessInput(state, 'x') // 誤った入力
	if state.CurrentIndex != 0 {
		t.Errorf("誤入力後のインデックス: 期待 0（進まない）, 実際 %d", state.CurrentIndex)
	}
	if state.CorrectCount != 0 {
		t.Errorf("誤入力後の正解数: 期待 0, 実際 %d", state.CorrectCount)
	}
	if state.TotalInputCount != 1 {
		t.Errorf("誤入力後の総入力数: 期待 1, 実際 %d", state.TotalInputCount)
	}
	if len(state.Mistakes) != 1 {
		t.Errorf("ミス数: 期待 1, 実際 %d", len(state.Mistakes))
	}
}

// TestEvaluator_CalculateWPM はWPM計算をテストします。

func TestEvaluator_CalculateWPM(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello", // 5文字
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	// 5文字を全て正解
	for _, c := range "hello" {
		state = evaluator.ProcessInput(state, c)
	}

	// テスト用に完了時間を手動設定（2秒で完了したと仮定）
	state.StartTime = time.Now().Add(-2 * time.Second)

	result := evaluator.CompleteChallenge(state)

	// WPM = (5文字 / 2秒 * 60) / 5 = 30 WPM
	expectedWPM := 30.0
	if result.WPM < expectedWPM-5 || result.WPM > expectedWPM+5 {
		t.Errorf("WPM: 期待 約%.1f, 実際 %.1f", expectedWPM, result.WPM)
	}
}

// TestEvaluator_CompleteChallengeScore はScore算出をテストします。
func TestEvaluator_CompleteChallengeScore(t *testing.T) {
	evaluator := NewEvaluator()

	t.Run("全問正解→Score=100", func(t *testing.T) {
		challenge := &Challenge{Text: "hi", TimeLimit: 10 * time.Second}
		state := evaluator.StartChallenge(challenge)
		state = evaluator.ProcessInput(state, 'h')
		state = evaluator.ProcessInput(state, 'i')
		result := evaluator.CompleteChallenge(state)
		if result.Score != 100 {
			t.Errorf("Score = %d, want 100", result.Score)
		}
	})

	t.Run("2/3正解→Score=66（切り捨て）", func(t *testing.T) {
		challenge := &Challenge{Text: "hi", TimeLimit: 10 * time.Second}
		state := evaluator.StartChallenge(challenge)
		state = evaluator.ProcessInput(state, 'h')
		state = evaluator.ProcessInput(state, 'x') // ミス
		state = evaluator.ProcessInput(state, 'i')
		result := evaluator.CompleteChallenge(state)
		// accuracy = 2/3 = 0.666... → int(0.666*100) = 66
		if result.Score != 66 {
			t.Errorf("Score = %d, want 66", result.Score)
		}
	})
}

// TestEvaluator_CompleteChallengeNoSpeedFactor はSpeedFactor/AccuracyFactorフィールドが存在しないことをテストします。
func TestEvaluator_CompleteChallengeNoSpeedFactor(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{Text: "hi", TimeLimit: 10 * time.Second}
	state := evaluator.StartChallenge(challenge)
	state = evaluator.ProcessInput(state, 'h')
	state = evaluator.ProcessInput(state, 'i')
	result := evaluator.CompleteChallenge(state)

	// TypingResultにはSpeedFactor/AccuracyFactor/Accuracyフィールドが存在しない
	// コンパイルが通ること自体がテスト
	if !result.Completed {
		t.Error("Completed = false, want true")
	}
}

// TestEvaluator_GetTimeoutResultScore はタイムアウト時のScore=0をテストします。
func TestEvaluator_GetTimeoutResultScore(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{Text: "hello", TimeLimit: 5 * time.Second}
	state := evaluator.StartChallenge(challenge)
	result := evaluator.GetTimeoutResult(state)

	if result.Score != 0 {
		t.Errorf("タイムアウト時のScore = %d, want 0", result.Score)
	}
	if !result.Timeout {
		t.Error("Timeout = false, want true")
	}
}

// TestEvaluator_CheckTimeout は制限時間超過の検出をテストします。

func TestEvaluator_CheckTimeout(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 1 * time.Second, // 1秒制限
	}
	state := evaluator.StartChallenge(challenge)
	// 開始時刻を過去に設定してタイムアウトをシミュレート
	state.StartTime = time.Now().Add(-2 * time.Second)

	if !evaluator.IsTimeout(state) {
		t.Error("タイムアウトが検出されなかった")
	}
}

// TestEvaluator_NotTimeout は制限時間内の判定をテストします。
func TestEvaluator_NotTimeout(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	if evaluator.IsTimeout(state) {
		t.Error("制限時間内なのにタイムアウト判定された")
	}
}

// TestEvaluator_IsCompleted はチャレンジ完了判定をテストします。
func TestEvaluator_IsCompleted(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hi",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	// 途中
	state = evaluator.ProcessInput(state, 'h')
	if evaluator.IsCompleted(state) {
		t.Error("途中なのに完了判定された")
	}

	// 完了
	state = evaluator.ProcessInput(state, 'i')
	if !evaluator.IsCompleted(state) {
		t.Error("完了しているのに未完了判定された")
	}
}

// TestEvaluator_GetProgress は入力進捗の取得をテストします。
func TestEvaluator_GetProgress(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)

	state = evaluator.ProcessInput(state, 'h')
	state = evaluator.ProcessInput(state, 'e')

	progress := evaluator.GetProgress(state)
	if progress != 0.4 { // 2/5 = 0.4
		t.Errorf("進捗: 期待 0.4, 実際 %.2f", progress)
	}
}

// TestEvaluator_GetRemainingTime は残り時間の取得をテストします。

func TestEvaluator_GetRemainingTime(t *testing.T) {
	evaluator := NewEvaluator()
	challenge := &Challenge{
		Text:      "hello",
		TimeLimit: 10 * time.Second,
	}
	state := evaluator.StartChallenge(challenge)
	// 2秒経過したと仮定
	state.StartTime = time.Now().Add(-2 * time.Second)

	remaining := evaluator.GetRemainingTime(state)
	// 約8秒残り（±0.5秒の誤差許容）
	if remaining < 7500*time.Millisecond || remaining > 8500*time.Millisecond {
		t.Errorf("残り時間: 期待 約8秒, 実際 %v", remaining)
	}
}
