// Package masterdata のテスト - データファイルの存在と内容検証
package masterdata

import (
	"testing"
)

// createTestLoader は埋め込みデータを使用するDataLoaderを作成します。
func createTestLoader() *DataLoader {
	return NewEmbeddedDataLoader(EmbeddedData, "data")
}

// TestCoresJSONExists はcores.jsonの存在と内容を検証します。
func TestCoresJSONExists(t *testing.T) {
	loader := createTestLoader()

	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("cores.jsonの読み込みに失敗: %v", err)
	}

	// 最低4種類のコア特性が存在すること

	if len(coreTypes) < 4 {
		t.Errorf("コア特性の数が足りません: got %d, want >= 4", len(coreTypes))
	}

	// 必須のコア特性IDを検証
	requiredIDs := map[string]bool{
		"attack_balance": false,
		"paladin":        false,
		"all_rounder":    false,
		"healer":         false,
	}

	for _, ct := range coreTypes {
		if _, ok := requiredIDs[ct.ID]; ok {
			requiredIDs[ct.ID] = true
		}
		// 各コア特性のバリデーション
		if err := ValidateCoreTypeData(ct); err != nil {
			t.Errorf("コア特性のバリデーションに失敗: %v", err)
		}
	}

	for id, found := range requiredIDs {
		if !found {
			t.Errorf("必須のコア特性が見つかりません: %s", id)
		}
	}
}

// TestModulesJSONExists はmodules.jsonの存在と内容を検証します。
func TestModulesJSONExists(t *testing.T) {
	loader := createTestLoader()

	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("modules.jsonの読み込みに失敗: %v", err)
	}

	// 各カテゴリにLv1〜Lv3が存在すること

	categoryLevelCount := make(map[string]map[int]bool)
	categories := []string{"physical_attack", "magic_attack", "heal", "buff", "debuff"}
	for _, cat := range categories {
		categoryLevelCount[cat] = make(map[int]bool)
	}

	for _, m := range modules {
		if err := ValidateModuleDefinitionData(m); err != nil {
			t.Errorf("モジュールのバリデーションに失敗: %v", err)
		}
		if _, ok := categoryLevelCount[m.Category]; ok {
			categoryLevelCount[m.Category][m.Level] = true
		}
	}

	// 各カテゴリにLv1が存在することを確認
	for cat, levels := range categoryLevelCount {
		if !levels[1] {
			t.Errorf("%s カテゴリにLv1モジュールがありません", cat)
		}
	}
}

// TestEnemiesJSONExists はenemies.jsonの存在と内容を検証します。
func TestEnemiesJSONExists(t *testing.T) {
	loader := createTestLoader()

	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("enemies.jsonの読み込みに失敗: %v", err)
	}

	// 最低4種類の敵バリエーションが存在すること
	if len(enemyTypes) < 4 {
		t.Errorf("敵タイプの数が足りません: got %d, want >= 4", len(enemyTypes))
	}

	for _, et := range enemyTypes {
		if err := ValidateEnemyTypeData(et); err != nil {
			t.Errorf("敵タイプのバリデーションに失敗: %v", err)
		}

		// ASCIIアートが設定されていること

		if et.ASCIIArt == "" {
			t.Errorf("敵タイプにASCIIアートがありません: %s", et.ID)
		}
	}
}

// TestWordsJSONExists はwords.jsonの存在と内容を検証します。
func TestWordsJSONExists(t *testing.T) {
	loader := createTestLoader()

	dictionary, err := loader.LoadTypingDictionary()
	if err != nil {
		t.Fatalf("words.jsonの読み込みに失敗: %v", err)
	}

	// 各難易度に単語が存在すること

	if len(dictionary.Easy) == 0 {
		t.Error("Easy単語が空です")
	}
	if len(dictionary.Medium) == 0 {
		t.Error("Medium単語が空です")
	}
	if len(dictionary.Hard) == 0 {
		t.Error("Hard単語が空です")
	}

	// 単語の長さの検証

	for _, word := range dictionary.Easy {
		if len(word) > 6 {
			t.Errorf("Easy単語が長すぎます: %s (len=%d)", word, len(word))
		}
	}

	for _, word := range dictionary.Medium {
		if len(word) < 4 || len(word) > 12 {
			t.Errorf("Medium単語の長さが不適切です: %s (len=%d)", word, len(word))
		}
	}

	for _, word := range dictionary.Hard {
		if len(word) < 8 {
			t.Errorf("Hard単語が短すぎます: %s (len=%d)", word, len(word))
		}
	}
}

// TestAllDataFilesLoadable は全データファイルが正常にロードできることを検証します。
func TestAllDataFilesLoadable(t *testing.T) {
	loader := createTestLoader()

	externalData, err := loader.LoadAllExternalData()
	if err != nil {
		t.Fatalf("全外部データのロードに失敗: %v", err)
	}

	// 全てのデータが存在することを確認
	if externalData.CoreTypes == nil || len(externalData.CoreTypes) == 0 {
		t.Error("CoreTypesが空です")
	}
	if externalData.ModuleDefinitions == nil || len(externalData.ModuleDefinitions) == 0 {
		t.Error("ModuleDefinitionsが空です")
	}
	if externalData.EnemyTypes == nil || len(externalData.EnemyTypes) == 0 {
		t.Error("EnemyTypesが空です")
	}
	if externalData.TypingDictionary == nil {
		t.Error("TypingDictionaryがnilです")
	}
}

// TestCoreTypeStatWeightsAreValid はコア特性のステータス重みが有効な範囲内かを検証します。
func TestCoreTypeStatWeightsAreValid(t *testing.T) {
	loader := createTestLoader()

	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("cores.jsonの読み込みに失敗: %v", err)
	}

	for _, ct := range coreTypes {
		// 各ステータス重みが0より大きいこと
		requiredStats := []string{"STR", "MAG", "SPD", "LUK"}
		for _, stat := range requiredStats {
			weight, exists := ct.StatWeights[stat]
			if !exists {
				t.Errorf("%s: %s の重みが定義されていません", ct.ID, stat)
				continue
			}
			if weight <= 0 {
				t.Errorf("%s: %s の重みが不正です: %f", ct.ID, stat, weight)
			}
		}
	}
}

// TestModuleTagsMatchCoreAllowedTags はモジュールのタグがコアの許可タグと適合することを検証します。
// 注: Requirement 5.18により、初期段階では高レベルモジュールを装備可能な特化コアは用意されていません。
// Lv1モジュール（_low タグ）のみが初期コアで使用可能であることを検証します。
func TestModuleTagsMatchCoreAllowedTags(t *testing.T) {
	loader := createTestLoader()

	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("cores.jsonの読み込みに失敗: %v", err)
	}

	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("modules.jsonの読み込みに失敗: %v", err)
	}

	// コアの許可タグを収集
	allowedTags := make(map[string]bool)
	for _, ct := range coreTypes {
		for _, tag := range ct.AllowedTags {
			allowedTags[tag] = true
		}
	}

	// Lv1モジュール（_low タグ）のみが初期コアで使用可能であることを確認
	// 高レベルモジュール（_mid, _high タグ）はゲーム進行で追加されるコアで使用可能になる想定
	for _, m := range modules {
		// Lv1モジュールのみ検証
		if m.Level != 1 {
			continue
		}
		hasValidTag := false
		for _, tag := range m.Tags {
			if allowedTags[tag] {
				hasValidTag = true
				break
			}
		}
		if !hasValidTag && len(m.Tags) > 0 {
			t.Errorf("Lv1モジュール %s のタグ %v がどのコアの許可タグにも含まれていません", m.ID, m.Tags)
		}
	}
}
