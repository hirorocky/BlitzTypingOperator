// Package components はTUI共通コンポーネントを提供します。
package components

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

func createTestPassiveSkill() domain.PassiveSkill {
	return domain.PassiveSkill{
		ID:          "test_passive",
		Name:        "テストスキル",
		Description: "テスト効果を付与する",
		Effects: map[domain.EffectColumn]float64{
			domain.ColSTRMultiplier: 1.1,
			domain.ColDamageCut:     0.1,
		},
	}
}

func TestPassiveSkillNotification_NewPassiveSkillNotification(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	if notification == nil {
		t.Fatal("NewPassiveSkillNotification should return non-nil")
	}
}

func TestPassiveSkillNotification_NewPassiveSkillNotificationWithNil(t *testing.T) {
	notification := NewPassiveSkillNotification(nil)

	if notification == nil {
		t.Fatal("NewPassiveSkillNotification should return non-nil even for nil skill")
	}
}

func TestPassiveSkillNotification_GetName(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	name := notification.GetName()
	if name != "テストスキル" {
		t.Errorf("GetName() = %v, want テストスキル", name)
	}
}

func TestPassiveSkillNotification_GetNameWithNil(t *testing.T) {
	notification := NewPassiveSkillNotification(nil)

	name := notification.GetName()
	if name != "" {
		t.Errorf("GetName() for nil skill = %v, want empty string", name)
	}
}

func TestPassiveSkillNotification_GetDescription(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	desc := notification.GetDescription()
	if desc != "テスト効果を付与する" {
		t.Errorf("GetDescription() = %v, want テスト効果を付与する", desc)
	}
}

func TestPassiveSkillNotification_GetEffects(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	effects := notification.GetEffects()

	// Effects マップの値を確認
	expectedSTRMult := 1.1
	if effects[domain.ColSTRMultiplier] != expectedSTRMult {
		t.Errorf("GetEffects()[ColSTRMultiplier] = %v, want %v", effects[domain.ColSTRMultiplier], expectedSTRMult)
	}
}

func TestPassiveSkillNotification_RenderCompact(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	result := notification.RenderCompact()

	// スキル名が含まれていることを確認
	if !strings.Contains(result, "テストスキル") {
		t.Errorf("RenderCompact() should contain skill name, got %v", result)
	}
}

func TestPassiveSkillNotification_RenderDetail(t *testing.T) {
	skill := createTestPassiveSkill()
	notification := NewPassiveSkillNotification(&skill)

	result := notification.RenderDetail(40)

	// スキル名と説明が含まれていることを確認
	if !strings.Contains(result, "テストスキル") {
		t.Errorf("RenderDetail() should contain skill name, got %v", result)
	}
	if !strings.Contains(result, "テスト効果") {
		t.Errorf("RenderDetail() should contain description, got %v", result)
	}
}

func TestPassiveSkillNotification_RenderNilSkill(t *testing.T) {
	notification := NewPassiveSkillNotification(nil)

	result := notification.RenderCompact()

	// nilの場合は空文字列または特定のプレースホルダー
	if result != "" && !strings.Contains(result, "-") && !strings.Contains(result, "None") {
		t.Errorf("RenderCompact() for nil skill should return empty or placeholder, got %v", result)
	}
}

func TestPassiveSkillNotification_HasActiveEffects(t *testing.T) {
	tests := []struct {
		name       string
		skill      *domain.PassiveSkill
		wantActive bool
	}{
		{
			name:       "with_skill",
			skill:      &domain.PassiveSkill{ID: "test", Name: "Test"},
			wantActive: true,
		},
		{
			name:       "nil_skill",
			skill:      nil,
			wantActive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification := NewPassiveSkillNotification(tt.skill)
			got := notification.HasActiveEffects()
			if got != tt.wantActive {
				t.Errorf("HasActiveEffects() = %v, want %v", got, tt.wantActive)
			}
		})
	}
}

func TestPassiveSkillNotification_RenderEffectsList(t *testing.T) {
	skill := domain.PassiveSkill{
		ID:          "multi_effect",
		Name:        "複合スキル",
		Description: "複数の効果を付与",
		Effects: map[domain.EffectColumn]float64{
			domain.ColSTRMultiplier: 1.1,
			domain.ColINTBonus:      5,
			domain.ColDamageCut:     0.05,
			domain.ColCritRate:      0.1,
		},
	}
	notification := NewPassiveSkillNotification(&skill)

	result := notification.RenderEffectsList()

	// 効果リストが返されることを確認
	if len(result) == 0 {
		t.Error("RenderEffectsList() should return non-empty list for skill with effects")
	}
}
