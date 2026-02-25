// Package masterdata はマスタデータのロード処理を提供します。
package masterdata

import (
	"os"
	"path/filepath"
	"testing"
)

// === 受け入れ基準 16,17: マスタデータ互換性 ===

// TestLoadSkillDefinitions_IsLatentField はis_latentフィールドのロードを確認します。
func TestLoadSkillDefinitions_IsLatentField(t *testing.T) {
	tmpDir := t.TempDir()

	skillsJSON := `{
		"skill_types": [
			{
				"id": "test_attack",
				"name": "テスト攻撃",
				"icon": "⚔️",
				"tags": ["physical_low"],
				"description": "テスト用",
				"cooldown_seconds": 3.0,
				"difficulty_rate": 100,
				"challenge": {"type": "standard"},
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 100, "stat_coef": 1.0, "stat_ref": "STR"},
						"probability": 1.0,
						"icon": "⚔️",
						"is_latent": true
					}
				]
			}
		]
	}`

	skillsPath := filepath.Join(tmpDir, "skills.json")
	if err := os.WriteFile(skillsPath, []byte(skillsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("スキル定義のロードに失敗: %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("スキル数: got %d, want 1", len(skills))
	}

	if len(skills[0].Effects) != 1 {
		t.Fatalf("エフェクト数: got %d, want 1", len(skills[0].Effects))
	}

	if !skills[0].Effects[0].IsLatent {
		t.Error("is_latent=true のエフェクトでIsLatent=trueであるべき")
	}
}

// TestLoadSkillDefinitions_IsLatentDefault はis_latent省略時にデフォルトfalseであることを確認します。
func TestLoadSkillDefinitions_IsLatentDefault(t *testing.T) {
	tmpDir := t.TempDir()

	// is_latentフィールドを省略
	skillsJSON := `{
		"skill_types": [
			{
				"id": "test_attack",
				"name": "テスト攻撃",
				"icon": "⚔️",
				"tags": ["physical_low"],
				"description": "テスト用",
				"cooldown_seconds": 3.0,
				"difficulty_rate": 100,
				"challenge": {"type": "standard"},
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 100, "stat_coef": 1.0, "stat_ref": "STR"},
						"probability": 1.0,
						"icon": "⚔️"
					}
				]
			}
		]
	}`

	skillsPath := filepath.Join(tmpDir, "skills.json")
	if err := os.WriteFile(skillsPath, []byte(skillsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("スキル定義のロードに失敗: %v", err)
	}

	if len(skills) != 1 {
		t.Fatalf("スキル数: got %d, want 1", len(skills))
	}

	if len(skills[0].Effects) != 1 {
		t.Fatalf("エフェクト数: got %d, want 1", len(skills[0].Effects))
	}

	if skills[0].Effects[0].IsLatent {
		t.Error("is_latent省略時にIsLatent=falseであるべき")
	}
}
