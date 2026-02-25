// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"math/rand"
	"testing"
)

// === 受け入れ基準 7,8: IsLatentフラグ ===

// TestSkillEffect_IsLatent_DefaultFalse はIsLatentのデフォルト値がfalseであることを確認します。
func TestSkillEffect_IsLatent_DefaultFalse(t *testing.T) {
	effect := SkillEffect{
		Target:      TargetEnemy,
		Probability: 1.0,
	}
	if effect.IsLatent {
		t.Error("IsLatentのデフォルト値はfalseであるべき")
	}
}

// TestSkillEffect_IsLatent_SetTrue はIsLatent=trueが設定できることを確認します。
func TestSkillEffect_IsLatent_SetTrue(t *testing.T) {
	effect := SkillEffect{
		Target:      TargetEnemy,
		Probability: 1.0,
		IsLatent:    true,
	}
	if !effect.IsLatent {
		t.Error("IsLatent=trueが設定されているべき")
	}
}

// === 受け入れ基準 9: LUK補正式の修正 ===

// TestSkillEffect_AdjustedProbability_LUKCorrection はLUK補正が
// probability + LUK × luk_factor で計算されることを確認します。
func TestSkillEffect_AdjustedProbability_LUKCorrection(t *testing.T) {
	tests := []struct {
		name        string
		probability float64
		lukFactor   float64
		luk         int
		expected    float64
	}{
		{
			name:        "LUK=10, factor=0.01 → probability + 10*0.01 = 0.6",
			probability: 0.5,
			lukFactor:   0.01,
			luk:         10,
			expected:    0.6,
		},
		{
			name:        "LUK=0, factor=0.01 → probability + 0*0.01 = 0.5",
			probability: 0.5,
			lukFactor:   0.01,
			luk:         0,
			expected:    0.5,
		},
		{
			name:        "LUK=20, factor=0.02 → probability + 20*0.02 = 0.9",
			probability: 0.5,
			lukFactor:   0.02,
			luk:         20,
			expected:    0.9,
		},
		{
			name:        "結果が1.0を超えない（クランプ）",
			probability: 0.8,
			lukFactor:   0.05,
			luk:         20,
			expected:    1.0,
		},
		{
			name:        "結果が0.0を下回らない（クランプ）",
			probability: 0.1,
			lukFactor:   -0.02,
			luk:         20,
			expected:    0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			effect := SkillEffect{
				Probability: tt.probability,
				LUKFactor:   tt.lukFactor,
			}
			got := effect.AdjustedProbability(tt.luk)
			if got != tt.expected {
				t.Errorf("AdjustedProbability(%d) = %v, want %v", tt.luk, got, tt.expected)
			}
		})
	}
}

// TestSkillEffect_ShouldTrigger_LUKCorrection はShouldTriggerがLUK補正後の確率で判定することを確認します。
func TestSkillEffect_ShouldTrigger_LUKCorrection(t *testing.T) {
	// LUK=10, probability=0.5, luk_factor=0.05 → 0.5 + 10*0.05 = 1.0 → 確実に発動
	effect := SkillEffect{
		Probability: 0.5,
		LUKFactor:   0.05,
	}
	rng := rand.New(rand.NewSource(42))
	if !effect.ShouldTrigger(10, rng) {
		t.Error("補正後確率1.0の場合、必ず発動するべき")
	}
}

// === 受け入れ基準 12: FeatureLatentEffect定数 ===

// TestFeatureLatentEffect_Defined はFeatureLatentEffect定数が定義されていることを確認します。
func TestFeatureLatentEffect_Defined(t *testing.T) {
	if FeatureLatentEffect != "latent_effect" {
		t.Errorf("FeatureLatentEffect = %q, want %q", FeatureLatentEffect, "latent_effect")
	}
}
