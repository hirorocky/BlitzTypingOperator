// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"strings"
	"testing"
)

// TestDescribeSingleEffect は単一効果の説明生成をテストします。
func TestDescribeSingleEffect(t *testing.T) {
	tests := []struct {
		name      string
		effectCol EffectColumn
		value     float64
		expected  string
	}{
		// DamageMultiplier
		{"DamageMultiplier UP", ColDamageMultiplier, 1.3, "ダメージ30%UP"},
		{"DamageMultiplier DOWN", ColDamageMultiplier, 0.7, "ダメージ30%ダウン"},
		{"DamageMultiplier 0%", ColDamageMultiplier, 1.0, "ダメージ0%UP"},

		// CooldownReduce
		{"CooldownReduce positive", ColCooldownReduce, 0.2, "CD20%短縮"},
		{"CooldownReduce negative", ColCooldownReduce, -0.2, "CD20%延長"},

		// DamageCut
		{"DamageCut positive", ColDamageCut, 0.15, "被ダメ15%軽減"},
		{"DamageCut negative", ColDamageCut, -0.15, "被ダメ15%増加"},

		// DamageBonus
		{"DamageBonus", ColDamageBonus, 10, "ダメージ+10"},

		// unknown
		{"unknown type", EffectColumn("unknown"), 1.0, "効果"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DescribeSingleEffect(tt.effectCol, tt.value)
			if result != tt.expected {
				t.Errorf("DescribeSingleEffect(%q, %v) = %q, want %q", tt.effectCol, tt.value, result, tt.expected)
			}
		})
	}
}

// TestDescribeEffectValues_SingleColumn は単一列の効果説明生成をテストします。
func TestDescribeEffectValues_SingleColumn(t *testing.T) {
	tests := []struct {
		name     string
		values   map[EffectColumn]float64
		expected string
	}{
		// 攻撃強化系
		{"DamageBonus positive", map[EffectColumn]float64{ColDamageBonus: 10}, "ダメージ+10"},
		{"DamageBonus negative", map[EffectColumn]float64{ColDamageBonus: -5}, "ダメージ-5"},
		{"DamageMultiplier UP", map[EffectColumn]float64{ColDamageMultiplier: 1.3}, "ダメージ30%UP"},
		{"DamageMultiplier DOWN", map[EffectColumn]float64{ColDamageMultiplier: 0.7}, "ダメージ30%ダウン"},
		{"ArmorPierce active", map[EffectColumn]float64{ColArmorPierce: 1}, "防御貫通"},
		{"LifeSteal", map[EffectColumn]float64{ColLifeSteal: 0.15}, "HP吸収15%"},

		// 防御強化系
		{"DamageCut positive", map[EffectColumn]float64{ColDamageCut: 0.2}, "被ダメ20%軽減"},
		{"DamageCut negative", map[EffectColumn]float64{ColDamageCut: -0.1}, "被ダメ10%増加"},
		{"Evasion", map[EffectColumn]float64{ColEvasion: 0.25}, "回避25%"},
		{"Reflect", map[EffectColumn]float64{ColReflect: 0.3}, "反射30%"},
		{"Regen", map[EffectColumn]float64{ColRegen: 5}, "HP回復5/s"},

		// 回復強化系
		{"HealBonus positive", map[EffectColumn]float64{ColHealBonus: 20}, "回復+20"},
		{"HealMultiplier", map[EffectColumn]float64{ColHealMultiplier: 1.2}, "回復20%UP"},

		// タイピング系
		{"TimeExtend positive", map[EffectColumn]float64{ColTimeExtend: 2.5}, "時間+2.5秒"},
		{"TimeExtend negative", map[EffectColumn]float64{ColTimeExtend: -1.0}, "時間-1.0秒"},
		{"AutoCorrect", map[EffectColumn]float64{ColAutoCorrect: 3}, "ミス無視3回"},

		// リキャスト系
		{"CooldownReduce positive", map[EffectColumn]float64{ColCooldownReduce: 0.2}, "CD20%短縮"},
		{"CooldownReduce negative", map[EffectColumn]float64{ColCooldownReduce: -0.15}, "CD15%延長"},

		// バフ/デバフ延長系
		{"BuffExtend", map[EffectColumn]float64{ColBuffExtend: 2.0}, "バフ延長+2.0秒"},
		{"DebuffExtend", map[EffectColumn]float64{ColDebuffExtend: 3.0}, "デバフ延長+3.0秒"},

		// 特殊系
		{"DoubleCast", map[EffectColumn]float64{ColDoubleCast: 0.15}, "2回発動15%"},

		// ステータス系
		{"CritRate", map[EffectColumn]float64{ColCritRate: 0.1}, "クリ率+10%"},
		{"STRBonus", map[EffectColumn]float64{ColSTRBonus: 5}, "STR+5"},
		{"STRMultiplier", map[EffectColumn]float64{ColSTRMultiplier: 0.25}, "STR25%UP"},
		{"STRMultiplier negative", map[EffectColumn]float64{ColSTRMultiplier: -0.2}, "STR20%ダウン"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DescribeEffectValues(tt.values)
			if result != tt.expected {
				t.Errorf("DescribeEffectValues(%v) = %q, want %q", tt.values, result, tt.expected)
			}
		})
	}
}

// TestDescribeEffectValues_MultipleColumns は複数列の効果説明生成をテストします。
func TestDescribeEffectValues_MultipleColumns(t *testing.T) {
	tests := []struct {
		name     string
		values   map[EffectColumn]float64
		contains []string
	}{
		{
			"DamageMultiplier and DamageCut",
			map[EffectColumn]float64{
				ColDamageMultiplier: 1.3,
				ColDamageCut:        0.2,
			},
			[]string{"ダメージ30%UP", "被ダメ20%軽減"},
		},
		{
			"Multiple stat bonuses",
			map[EffectColumn]float64{
				ColSTRBonus: 5,
				ColINTBonus: 3,
			},
			[]string{"STR+5", "INT+3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DescribeEffectValues(tt.values)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("DescribeEffectValues result %q should contain %q", result, expected)
				}
			}
			// 複数効果は " / " で区切られる
			if !strings.Contains(result, " / ") {
				t.Errorf("DescribeEffectValues result %q should contain ' / ' separator", result)
			}
		})
	}
}

// TestDescribeEffectValues_EmptyMap は空のマップの場合をテストします。
func TestDescribeEffectValues_EmptyMap(t *testing.T) {
	result := DescribeEffectValues(map[EffectColumn]float64{})
	if result != "効果" {
		t.Errorf("DescribeEffectValues(empty) = %q, want %q", result, "効果")
	}
}

// TestDescribeEffectValues_NilMap はnilマップの場合をテストします。
func TestDescribeEffectValues_NilMap(t *testing.T) {
	result := DescribeEffectValues(nil)
	if result != "効果" {
		t.Errorf("DescribeEffectValues(nil) = %q, want %q", result, "効果")
	}
}
