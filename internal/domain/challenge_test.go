package domain

import (
	"testing"
	"time"
)

// === DifficultyRate ===

func TestDifficultyRateClamp(t *testing.T) {
	tests := []struct {
		name     string
		input    DifficultyRate
		expected DifficultyRate
	}{
		{"範囲内はそのまま", 100, 100},
		{"最小値", 50, 50},
		{"最大値", 200, 200},
		{"下限未満はClampされる", 30, 50},
		{"上限超過はClampされる", 250, 200},
		{"0はClampされる", 0, 50},
		{"負値はClampされる", -10, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Clamp()
			if result != tt.expected {
				t.Errorf("DifficultyRate(%d).Clamp() = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// === ChallengeStatus ===

func TestChallengeStatusValues(t *testing.T) {
	// 3つの状態が定義されていることを確認
	statuses := []ChallengeStatus{ChallengeSuccess, ChallengeFail, ChallengeCancel}
	if len(statuses) != 3 {
		t.Errorf("ChallengeStatusは3つの値を持つべき")
	}

	// それぞれが異なることを確認
	if ChallengeSuccess == ChallengeFail || ChallengeSuccess == ChallengeCancel || ChallengeFail == ChallengeCancel {
		t.Errorf("ChallengeStatusの値はそれぞれ異なるべき")
	}
}

// === ChallengeInput ===

func TestChallengeInputValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   ChallengeInput
		wantErr bool
	}{
		{
			name: "有効な入力",
			input: ChallengeInput{
				Difficulty:       100,
				AutoCorrectCount: 0,
			},
			wantErr: false,
		},
		{
			name: "最小DifficultyRate",
			input: ChallengeInput{
				Difficulty:       50,
				AutoCorrectCount: 0,
			},
			wantErr: false,
		},
		{
			name: "最大DifficultyRate",
			input: ChallengeInput{
				Difficulty:       200,
				AutoCorrectCount: 0,
			},
			wantErr: false,
		},
		{
			name: "範囲外のDifficultyRateはエラー",
			input: ChallengeInput{
				Difficulty:       30,
				AutoCorrectCount: 0,
			},
			wantErr: true,
		},
		{
			name: "負のAutoCorrectCountはエラー",
			input: ChallengeInput{
				Difficulty:       100,
				AutoCorrectCount: -1,
			},
			wantErr: true,
		},
		{
			name: "正のAutoCorrectCount",
			input: ChallengeInput{
				Difficulty:       100,
				AutoCorrectCount: 3,
			},
			wantErr: false,
		},
		{
			name: "Words付き有効な入力",
			input: ChallengeInput{
				Difficulty:       100,
				Words:            []string{"hello", "world"},
				AutoCorrectCount: 0,
			},
			wantErr: false,
		},
		{
			name: "全フィールド設定",
			input: ChallengeInput{
				Difficulty:               150,
				Words:                    []string{"test"},
				TimeExtendSec:            2.0,
				AutoCorrectCount:         2,
				MistakeTimeExtendSec:     1.5,
				RetryOnTimeout:           true,
				RetryTimeLimitMultiplier: 0.5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ChallengeInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// === ChallengeOutput ===

func TestChallengeOutputFields(t *testing.T) {
	output := ChallengeOutput{
		Accuracy:       0.95,
		SpeedFactor:    1.5,
		WPM:            80.0,
		CompletionTime: 5 * time.Second,
		Status:         ChallengeSuccess,
	}

	if output.Accuracy != 0.95 {
		t.Errorf("Accuracy = %f, want 0.95", output.Accuracy)
	}
	if output.SpeedFactor != 1.5 {
		t.Errorf("SpeedFactor = %f, want 1.5", output.SpeedFactor)
	}
	if output.WPM != 80.0 {
		t.Errorf("WPM = %f, want 80.0", output.WPM)
	}
	if output.CompletionTime != 5*time.Second {
		t.Errorf("CompletionTime = %v, want 5s", output.CompletionTime)
	}
	if output.Status != ChallengeSuccess {
		t.Errorf("Status = %v, want ChallengeSuccess", output.Status)
	}
}

// === ChallengeTypeID ===

func TestChallengeTypeIDConstants(t *testing.T) {
	// 3つのタイプが定義されていることを確認
	if ChallengeTypeStandard == "" {
		t.Error("ChallengeTypeStandard が空文字列")
	}
	if ChallengeTypeShape == "" {
		t.Error("ChallengeTypeShape が空文字列")
	}
	if ChallengeTypeDefense == "" {
		t.Error("ChallengeTypeDefense が空文字列")
	}

	// それぞれ異なる値であること
	if ChallengeTypeStandard == ChallengeTypeShape {
		t.Error("Standard と Shape が同じ値")
	}
	if ChallengeTypeStandard == ChallengeTypeDefense {
		t.Error("Standard と Defense が同じ値")
	}
}
