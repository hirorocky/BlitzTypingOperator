package masterdata

import (
	"os"
	"path/filepath"
	"testing"
)

// === 受け入れテスト: ローダーのパース・変換 ===

// 受け入れ基準1-3: challenge_typeがchallengeオブジェクトに置き換わり、type+optionsを持つ
// 受け入れ基準4: ローダーが新形式を正しくパースできる
// 受け入れ基準36: ローダーのパーステスト
func TestLoadModuleDefinitions_新形式パース(t *testing.T) {
	tmpDir := t.TempDir()

	modulesJSON := `{
		"module_types": [
			{
				"id": "fireball_lv1",
				"name": "初級攻撃魔法",
				"icon": "💥",
				"tags": ["magic_low"],
				"description": "炎の魔法弾を放つ",
				"cooldown_seconds": 10.0,
				"difficulty": 80,
				"challenge": {
					"type": "shape",
					"options": {
						"shape": "flame"
					}
				},
				"min_drop_level": 1,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 1.2, "stat_ref": "INT"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "💥"
					}
				]
			},
			{
				"id": "physical_strike_lv1",
				"name": "軽斬撃",
				"icon": "🔪",
				"tags": ["physical_low"],
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 10.0,
				"difficulty": 80,
				"challenge": {
					"type": "standard"
				},
				"min_drop_level": 1,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 10.0, "stat_ref": "STR"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "🔪"
					}
				]
			},
			{
				"id": "defense_buff_lv1",
				"name": "防御強化1",
				"icon": "🧱",
				"tags": ["buff_low"],
				"description": "防御を強化する",
				"cooldown_seconds": 10.0,
				"difficulty": 80,
				"challenge": {
					"type": "defense"
				},
				"min_drop_level": 1,
				"effects": [
					{
						"target": "self",
						"effect_column": {"column": "damage_cut", "value": 0.1, "duration": 10.0},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "🧱"
					}
				]
			}
		]
	}`

	modulesPath := filepath.Join(tmpDir, "modules.json")
	if err := os.WriteFile(modulesPath, []byte(modulesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("モジュール定義のロードに失敗: %v", err)
	}

	if len(modules) != 3 {
		t.Fatalf("モジュール数: got %d, want 3", len(modules))
	}

	// shapeチャレンジのパース検証
	if modules[0].Challenge.Type != "shape" {
		t.Errorf("fireball Challenge.Type = %q, want %q", modules[0].Challenge.Type, "shape")
	}
	if modules[0].Challenge.Options["shape"] != "flame" {
		t.Errorf("fireball Challenge.Options[shape] = %q, want %q", modules[0].Challenge.Options["shape"], "flame")
	}

	// standardチャレンジのパース検証
	if modules[1].Challenge.Type != "standard" {
		t.Errorf("physical_strike Challenge.Type = %q, want %q", modules[1].Challenge.Type, "standard")
	}
	if len(modules[1].Challenge.Options) != 0 {
		t.Errorf("physical_strike Challenge.Options should be empty, got %v", modules[1].Challenge.Options)
	}

	// defenseチャレンジのパース検証
	if modules[2].Challenge.Type != "defense" {
		t.Errorf("defense_buff Challenge.Type = %q, want %q", modules[2].Challenge.Type, "defense")
	}
}

// 受け入れ基準5: ToDomainType()がChallengeTypeID + ChallengeOptionsを正しく変換できる
// 受け入れ基準37: ローダーの変換テスト
func TestToDomainType_ChallengeTypeIDとOptions(t *testing.T) {
	tests := []struct {
		name            string
		challengeData   ChallengeData
		wantType        string
		wantOptionsLen  int
		wantOptionKey   string
		wantOptionValue string
	}{
		{
			name:            "shapeチャレンジ",
			challengeData:   ChallengeData{Type: "shape", Options: map[string]string{"shape": "flame"}},
			wantType:        "shape",
			wantOptionsLen:  1,
			wantOptionKey:   "shape",
			wantOptionValue: "flame",
		},
		{
			name:           "standardチャレンジ",
			challengeData:  ChallengeData{Type: "standard"},
			wantType:       "standard",
			wantOptionsLen: 0,
		},
		{
			name:           "defenseチャレンジ",
			challengeData:  ChallengeData{Type: "defense"},
			wantType:       "defense",
			wantOptionsLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			module := &SkillDefinitionData{
				ID:              "test_skill",
				Name:            "テスト",
				Icon:            "🔥",
				Tags:            []string{"test"},
				Description:     "テスト用",
				CooldownSeconds: 10.0,
				DifficultyRate:  80,
				Challenge:       tt.challengeData,
				MinDropLevel:    1,
				Effects: []SkillEffectData{
					{
						Target:      "enemy",
						HPFormula:   &HPFormulaData{Base: 0, StatCoef: 1.0, StatRef: "STR"},
						Probability: 1.0,
						Icon:        "🔥",
					},
				},
			}

			domainType := module.ToDomainType()

			if string(domainType.ChallengeType) != tt.wantType {
				t.Errorf("ChallengeType = %q, want %q", domainType.ChallengeType, tt.wantType)
			}

			if len(domainType.ChallengeOptions) != tt.wantOptionsLen {
				t.Errorf("ChallengeOptions length = %d, want %d", len(domainType.ChallengeOptions), tt.wantOptionsLen)
			}

			if tt.wantOptionsLen > 0 {
				if domainType.ChallengeOptions[tt.wantOptionKey] != tt.wantOptionValue {
					t.Errorf("ChallengeOptions[%s] = %q, want %q",
						tt.wantOptionKey,
						domainType.ChallengeOptions[tt.wantOptionKey],
						tt.wantOptionValue)
				}
			}
		})
	}
}
