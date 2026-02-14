// Package masterdata のテスト - データファイルの存在と内容検証
package masterdata

import (
	"embed"
	"strings"
	"testing"
)

//go:embed testdata/*
var testData embed.FS

// createTestLoader は埋め込みデータを使用するDataLoaderを作成します。
func createTestLoader() *DataLoader {
	return NewEmbeddedDataLoader(EmbeddedData, "data")
}

// createTestDataLoader はテスト用データを使用するDataLoaderを作成します。
func createTestDataLoader() *DataLoader {
	return NewEmbeddedDataLoader(testData, "testdata")
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

// TestSkillsJSONExists はskills.jsonの存在と内容を検証します。
func TestSkillsJSONExists(t *testing.T) {
	loader := createTestLoader()

	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("skills.jsonの読み込みに失敗: %v", err)
	}

	// 各効果タイプのスキルが存在すること
	effectTypeCount := make(map[string]int)
	effectTypes := []string{"damage", "heal", "buff", "debuff"}
	for _, et := range effectTypes {
		effectTypeCount[et] = 0
	}

	for _, m := range skills {
		if err := ValidateSkillDefinitionData(m); err != nil {
			t.Errorf("スキルのバリデーションに失敗: %v", err)
		}
		// スキルの効果をドメインモデルに変換して判定
		domainSkill := m.ToDomain()
		for _, effect := range domainSkill.Effects() {
			if effect.IsDamageEffect() {
				effectTypeCount["damage"]++
			}
			if effect.IsHealEffect() {
				effectTypeCount["heal"]++
			}
			if effect.IsBuffEffect() {
				effectTypeCount["buff"]++
			}
			if effect.IsDebuffEffect() {
				effectTypeCount["debuff"]++
			}
		}
	}

	// 各効果タイプにスキルが存在することを確認
	for et, count := range effectTypeCount {
		if count == 0 {
			t.Errorf("%s 効果を持つスキルがありません", et)
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
	}
}

// TestWordsJSONExists はwords.jsonの存在と内容を検証します。
// テスト用のwords.jsonを使用して、本番データの変更に影響されないようにします。
func TestWordsJSONExists(t *testing.T) {
	loader := createTestDataLoader()

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
	if externalData.SkillDefinitions == nil || len(externalData.SkillDefinitions) == 0 {
		t.Error("SkillDefinitionsが空です")
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
		requiredStats := []string{"STR", "INT", "WIL", "LUK"}
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

// TestCorePassiveSkillIDsExist はコアのpassive_skill_idがpassive_skills.jsonに存在することを検証します。
func TestCorePassiveSkillIDsExist(t *testing.T) {
	loader := createTestLoader()

	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("cores.jsonの読み込みに失敗: %v", err)
	}

	passiveSkills, err := loader.LoadPassiveSkills()
	if err != nil {
		t.Fatalf("passive_skills.jsonの読み込みに失敗: %v", err)
	}

	// パッシブスキルIDのマップを作成
	passiveSkillIDs := make(map[string]bool)
	for _, ps := range passiveSkills {
		passiveSkillIDs[ps.ID] = true
	}

	// 各コアのpassive_skill_idがパッシブスキルに存在するか検証
	for _, ct := range coreTypes {
		if ct.PassiveSkillID == "" {
			t.Errorf("コア %s にpassive_skill_idが設定されていません", ct.ID)
			continue
		}
		if !passiveSkillIDs[ct.PassiveSkillID] {
			t.Errorf("コア %s のpassive_skill_id %s がpassive_skills.jsonに存在しません", ct.ID, ct.PassiveSkillID)
		}
	}
}

// TestSkillTagsMatchCoreAllowedTags はスキルのタグがコアの許可タグと適合することを検証します。
// 注: Requirement 5.18により、初期段階では高レベルスキルを装備可能な特化コアは用意されていません。
// Lv1スキル（_low タグ）のみが初期コアで使用可能であることを検証します。
func TestSkillTagsMatchCoreAllowedTags(t *testing.T) {
	loader := createTestLoader()

	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("cores.jsonの読み込みに失敗: %v", err)
	}

	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("skills.jsonの読み込みに失敗: %v", err)
	}

	// コアの許可タグを収集
	allowedTags := make(map[string]bool)
	for _, ct := range coreTypes {
		for _, tag := range ct.AllowedTags {
			allowedTags[tag] = true
		}
	}

	// _low タグを持つスキルのみが初期コアで使用可能であることを確認
	// 高難度スキル（_mid, _high タグ）はゲーム進行で追加されるコアで使用可能になる想定
	for _, m := range skills {
		// _low タグを持つスキルのみ検証
		hasLowTag := false
		for _, tag := range m.Tags {
			if strings.HasSuffix(tag, "_low") {
				hasLowTag = true
				break
			}
		}
		if !hasLowTag {
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
			t.Errorf("基本スキル %s のタグ %v がどのコアの許可タグにも含まれていません", m.ID, m.Tags)
		}
	}
}

// ===== 受け入れ基準 6: shapeチャレンジタイプの全スキルにmana_cost:1が設定されている =====

func TestShapeSkillsHaveManaCost(t *testing.T) {
	loader := createTestLoader()

	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("skills.jsonの読み込みに失敗: %v", err)
	}

	shapeCount := 0
	for _, s := range skills {
		if s.Challenge.Type == "shape" {
			shapeCount++
			if s.ManaCost != 1 {
				t.Errorf("shapeスキル %s のManaCostが1ではありません: got %d", s.ID, s.ManaCost)
			}
		}
	}

	if shapeCount == 0 {
		t.Error("shapeチャレンジタイプのスキルが存在しません")
	}
}

// ===== 受け入れ基準 9: mana_gain付きスキルが存在する =====

func TestManaGainSkillsExist(t *testing.T) {
	loader := createTestLoader()

	skills, err := loader.LoadSkillDefinitions()
	if err != nil {
		t.Fatalf("skills.jsonの読み込みに失敗: %v", err)
	}

	manaGainCount := 0
	for _, s := range skills {
		for _, e := range s.Effects {
			if e.ManaGain > 0 {
				manaGainCount++
				break
			}
		}
	}

	if manaGainCount == 0 {
		t.Error("mana_gain付きのスキルが存在しません")
	}
}
