// Package masterdata はマスタデータのロード処理を提供します。
package masterdata

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// TestLoadCoreTypes はコア特性定義のロードをテストします。
func TestLoadCoreTypes(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用cores.jsonを作成
	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			},
			{
				"id": "healer",
				"name": "ヒーラー",
				"allowed_tags": ["heal_mid", "heal_high"],
				"stat_weights": {"STR": 0.5, "MAG": 1.5, "SPD": 0.8, "LUK": 1.2},
				"passive_skill_id": "healing_aura",
				"min_drop_level": 3
			}
		]
	}`

	coresPath := filepath.Join(tmpDir, "cores.json")
	if err := os.WriteFile(coresPath, []byte(coresJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	// ロード実行
	loader := NewDataLoader(tmpDir)
	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("コア特性のロードに失敗: %v", err)
	}

	// 検証
	if len(coreTypes) != 2 {
		t.Errorf("コア特性の数が期待値と異なる: got %d, want 2", len(coreTypes))
	}

	// 攻撃バランスの検証
	if coreTypes[0].ID != "attack_balance" {
		t.Errorf("ID: got %s, want attack_balance", coreTypes[0].ID)
	}
	if coreTypes[0].Name != "攻撃バランス" {
		t.Errorf("Name: got %s, want 攻撃バランス", coreTypes[0].Name)
	}
	if len(coreTypes[0].AllowedTags) != 2 {
		t.Errorf("AllowedTags length: got %d, want 2", len(coreTypes[0].AllowedTags))
	}
	if coreTypes[0].StatWeights["STR"] != 1.2 {
		t.Errorf("StatWeights[STR]: got %f, want 1.2", coreTypes[0].StatWeights["STR"])
	}
	if coreTypes[0].MinDropLevel != 1 {
		t.Errorf("MinDropLevel: got %d, want 1", coreTypes[0].MinDropLevel)
	}
}

// TestLoadModuleDefinitions はモジュール定義のロードをテストします。
func TestLoadModuleDefinitions(t *testing.T) {
	tmpDir := t.TempDir()

	modulesJSON := `{
		"modules": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"category": "physical_attack",
				"level": 1,
				"tags": ["physical_low"],
				"base_effect": 10.0,
				"stat_reference": "STR",
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1
			},
			{
				"id": "fireball_lv2",
				"name": "ファイアボールLv2",
				"category": "magic_attack",
				"level": 2,
				"tags": ["magic_mid"],
				"base_effect": 20.0,
				"stat_reference": "MAG",
				"description": "中級の魔法攻撃",
				"cooldown_seconds": 3.5,
				"difficulty": 2,
				"min_drop_level": 10
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

	if len(modules) != 2 {
		t.Errorf("モジュール数: got %d, want 2", len(modules))
	}

	// 物理打撃Lv1の検証
	if modules[0].ID != "physical_strike_lv1" {
		t.Errorf("ID: got %s, want physical_strike_lv1", modules[0].ID)
	}
	if modules[0].Category != "physical_attack" {
		t.Errorf("Category: got %s, want physical_attack", modules[0].Category)
	}
	if modules[0].Level != 1 {
		t.Errorf("Level: got %d, want 1", modules[0].Level)
	}
	if modules[0].BaseEffect != 10.0 {
		t.Errorf("BaseEffect: got %f, want 10.0", modules[0].BaseEffect)
	}
}

// TestLoadEnemyTypes は敵タイプ定義のロードをテストします。
func TestLoadEnemyTypes(t *testing.T) {
	tmpDir := t.TempDir()

	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"base_attack_interval_ms": 3000,
				"attack_type": "physical",
				"ascii_art": "  ___\n /   \\\n|     |\n \\___|"
			},
			{
				"id": "goblin",
				"name": "ゴブリン",
				"base_hp": 80,
				"base_attack_power": 10,
				"base_attack_interval_ms": 2500,
				"attack_type": "physical",
				"ascii_art": "  /\\_/\\\n ( o o )\n  > ^ <"
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	if len(enemyTypes) != 2 {
		t.Errorf("敵タイプ数: got %d, want 2", len(enemyTypes))
	}

	// スライムの検証
	if enemyTypes[0].ID != "slime" {
		t.Errorf("ID: got %s, want slime", enemyTypes[0].ID)
	}
	if enemyTypes[0].Name != "スライム" {
		t.Errorf("Name: got %s, want スライム", enemyTypes[0].Name)
	}
	if enemyTypes[0].BaseHP != 50 {
		t.Errorf("BaseHP: got %d, want 50", enemyTypes[0].BaseHP)
	}
	if enemyTypes[0].BaseAttackInterval != 3000*time.Millisecond {
		t.Errorf("BaseAttackInterval: got %v, want 3s", enemyTypes[0].BaseAttackInterval)
	}
}

// TestLoadTypingDictionary はタイピング辞書のロードをテストします。
func TestLoadTypingDictionary(t *testing.T) {
	tmpDir := t.TempDir()

	wordsJSON := `{
		"words": {
			"easy": ["cat", "dog", "sun", "run"],
			"medium": ["function", "variable", "package"],
			"hard": ["implementation", "infrastructure"]
		}
	}`

	wordsPath := filepath.Join(tmpDir, "words.json")
	if err := os.WriteFile(wordsPath, []byte(wordsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	dictionary, err := loader.LoadTypingDictionary()
	if err != nil {
		t.Fatalf("タイピング辞書のロードに失敗: %v", err)
	}

	// easy単語の検証
	if len(dictionary.Easy) != 4 {
		t.Errorf("Easy単語数: got %d, want 4", len(dictionary.Easy))
	}
	if dictionary.Easy[0] != "cat" {
		t.Errorf("Easy[0]: got %s, want cat", dictionary.Easy[0])
	}

	// medium単語の検証
	if len(dictionary.Medium) != 3 {
		t.Errorf("Medium単語数: got %d, want 3", len(dictionary.Medium))
	}

	// hard単語の検証
	if len(dictionary.Hard) != 2 {
		t.Errorf("Hard単語数: got %d, want 2", len(dictionary.Hard))
	}
}

// TestLoadCoreTypesFileNotFound はファイルが存在しない場合のエラーをテストします。
func TestLoadCoreTypesFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	loader := NewDataLoader(tmpDir)
	_, err := loader.LoadCoreTypes()
	if err == nil {
		t.Error("ファイルが存在しない場合はエラーが返されるべき")
	}
}

// TestLoadCoreTypesInvalidJSON は不正なJSONの場合のエラーをテストします。
func TestLoadCoreTypesInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// 不正なJSON
	invalidJSON := `{ invalid json }`
	coresPath := filepath.Join(tmpDir, "cores.json")
	if err := os.WriteFile(coresPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	_, err := loader.LoadCoreTypes()
	if err == nil {
		t.Error("不正なJSONの場合はエラーが返されるべき")
	}
}

// TestLoadAllExternalData は全外部データの一括ロードをテストします。
func TestLoadAllExternalData(t *testing.T) {
	tmpDir := t.TempDir()

	// cores.json
	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "cores.json"), []byte(coresJSON), 0644)

	// modules.json
	modulesJSON := `{
		"modules": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"category": "physical_attack",
				"level": 1,
				"tags": ["physical_low"],
				"base_effect": 10.0,
				"stat_reference": "STR",
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "modules.json"), []byte(modulesJSON), 0644)

	// enemies.json
	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"base_attack_interval_ms": 3000,
				"attack_type": "physical",
				"ascii_art": "  ___\n /   \\\n|     |"
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "enemies.json"), []byte(enemiesJSON), 0644)

	// words.json
	wordsJSON := `{
		"words": {
			"easy": ["cat", "dog"],
			"medium": ["function"],
			"hard": ["implementation"]
		}
	}`
	os.WriteFile(filepath.Join(tmpDir, "words.json"), []byte(wordsJSON), 0644)

	loader := NewDataLoader(tmpDir)
	externalData, err := loader.LoadAllExternalData()
	if err != nil {
		t.Fatalf("全外部データのロードに失敗: %v", err)
	}

	if len(externalData.CoreTypes) != 1 {
		t.Errorf("CoreTypes: got %d, want 1", len(externalData.CoreTypes))
	}
	if len(externalData.ModuleDefinitions) != 1 {
		t.Errorf("ModuleDefinitions: got %d, want 1", len(externalData.ModuleDefinitions))
	}
	if len(externalData.EnemyTypes) != 1 {
		t.Errorf("EnemyTypes: got %d, want 1", len(externalData.EnemyTypes))
	}
	if externalData.TypingDictionary == nil {
		t.Error("TypingDictionary should not be nil")
	}
}

// TestConvertToDomainCoreType はJSONデータからドメインモデルへの変換をテストします。
func TestConvertToDomainCoreType(t *testing.T) {
	tmpDir := t.TempDir()

	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "cores.json"), []byte(coresJSON), 0644)

	loader := NewDataLoader(tmpDir)
	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("コア特性のロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainCoreType := coreTypes[0].ToDomain()

	if domainCoreType.ID != "attack_balance" {
		t.Errorf("ID: got %s, want attack_balance", domainCoreType.ID)
	}
	if domainCoreType.StatWeights["STR"] != 1.2 {
		t.Errorf("StatWeights[STR]: got %f, want 1.2", domainCoreType.StatWeights["STR"])
	}
}

// TestConvertToDomainModuleModel はモジュールデータからドメインモデルへの変換をテストします。
func TestConvertToDomainModuleModel(t *testing.T) {
	tmpDir := t.TempDir()

	modulesJSON := `{
		"modules": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"category": "physical_attack",
				"level": 1,
				"tags": ["physical_low"],
				"base_effect": 10.0,
				"stat_reference": "STR",
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "modules.json"), []byte(modulesJSON), 0644)

	loader := NewDataLoader(tmpDir)
	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("モジュール定義のロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainModule := modules[0].ToDomain()

	if domainModule.ID != "physical_strike_lv1" {
		t.Errorf("ID: got %s, want physical_strike_lv1", domainModule.ID)
	}
	if domainModule.Category != domain.PhysicalAttack {
		t.Errorf("Category: got %v, want %v", domainModule.Category, domain.PhysicalAttack)
	}
}
