package domain_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// === 受け入れ基準 1: 新規ゲーム開始時、全機能はLocked状態で初期化される ===
func TestNewFeatureUnlockState_AllFeaturesLocked(t *testing.T) {
	state := domain.NewFeatureUnlockState()

	features := []domain.FeatureID{
		domain.FeatureDefenseSkill,
		domain.FeatureAgentCustomization,
		domain.FeatureChainEffect,
		domain.FeatureManaSystem,
	}

	for _, f := range features {
		status := state.GetStatus(f)
		if status != domain.FeatureLocked {
			t.Errorf("新規状態で %s が %v、期待: %v", f, status, domain.FeatureLocked)
		}
	}
}

// === 受け入れ基準 8: 状態遷移は単調増加のみ ===
func TestFeatureUnlockState_MonotonicTransition(t *testing.T) {
	testCases := []struct {
		name     string
		from     domain.FeatureStatus
		to       domain.FeatureStatus
		canApply bool
	}{
		{"Locked→PendingTutorial", domain.FeatureLocked, domain.FeaturePendingTutorial, true},
		{"PendingTutorial→Unlocked", domain.FeaturePendingTutorial, domain.FeatureUnlocked, true},
		{"Locked→Unlocked", domain.FeatureLocked, domain.FeatureUnlocked, true},
		{"Unlocked→Locked（ダウングレード）", domain.FeatureUnlocked, domain.FeatureLocked, false},
		{"Unlocked→PendingTutorial（ダウングレード）", domain.FeatureUnlocked, domain.FeaturePendingTutorial, false},
		{"PendingTutorial→Locked（ダウングレード）", domain.FeaturePendingTutorial, domain.FeatureLocked, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := domain.CanTransition(tc.from, tc.to)
			if result != tc.canApply {
				t.Errorf("CanTransition(%v, %v) = %v, 期待: %v", tc.from, tc.to, result, tc.canApply)
			}
		})
	}
}

// === 受け入れ基準 9: 不正な遷移要求はエラーを返す ===
func TestFeatureUnlockState_InvalidTransitionReturnsError(t *testing.T) {
	state := domain.NewFeatureUnlockState()

	// Locked状態でCompleteTutorialを試みる（PendingTutorial→Unlockedの遷移に相当）
	err := state.TransitionTo(domain.FeatureDefenseSkill, domain.FeatureUnlocked)
	if err != nil {
		// Locked→Unlockedは許可される（単調増加）
		t.Errorf("Locked→Unlockedでエラーが返った: %v", err)
	}

	// Unlocked状態でダウングレードを試みる
	err = state.TransitionTo(domain.FeatureDefenseSkill, domain.FeatureLocked)
	if err == nil {
		t.Error("Unlocked→Lockedでエラーが返らなかった")
	}
}

// === 受け入れ基準 8: 未知のFeatureIDに対する操作 ===
func TestFeatureUnlockState_UnknownFeatureID(t *testing.T) {
	state := domain.NewFeatureUnlockState()

	status := state.GetStatus("unknown_feature")
	if status != domain.FeatureLocked {
		t.Errorf("未知のFeatureIDのステータスが %v、期待: %v (Locked)", status, domain.FeatureLocked)
	}
}
