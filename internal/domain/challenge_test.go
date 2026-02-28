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
	// 4つの状態が定義されていることを確認
	statuses := []ChallengeStatus{ChallengeSuccess, ChallengeFail, ChallengeCancel, ChallengePerfect}
	if len(statuses) != 4 {
		t.Errorf("ChallengeStatusは4つの値を持つべき")
	}

	// それぞれが異なることを確認
	seen := make(map[ChallengeStatus]bool)
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("ChallengeStatusに重複値がある: %d", s)
		}
		seen[s] = true
	}
}

func TestChallengePerfect_Value(t *testing.T) {
	// ChallengePerfectは値3であること（iota末尾、後方互換）
	if ChallengePerfect != 3 {
		t.Errorf("ChallengePerfect = %d, want 3", ChallengePerfect)
	}
}

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		name   string
		status ChallengeStatus
		want   bool
	}{
		{"Success→true", ChallengeSuccess, true},
		{"Perfect→true", ChallengePerfect, true},
		{"Fail→false", ChallengeFail, false},
		{"Cancel→false", ChallengeCancel, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsSuccess()
			if got != tt.want {
				t.Errorf("ChallengeStatus(%d).IsSuccess() = %v, want %v", tt.status, got, tt.want)
			}
		})
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
		Score:          95,
		WPM:            80.0,
		CompletionTime: 5 * time.Second,
		Status:         ChallengeSuccess,
	}

	if output.Score != 95 {
		t.Errorf("Score = %d, want 95", output.Score)
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
