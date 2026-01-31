// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestStats_ゼロ値の確認 はStats構造体のゼロ値が正しいことを確認します。
func TestStats_ゼロ値の確認(t *testing.T) {
	stats := Stats{}

	if stats.STR != 0 {
		t.Errorf("STRのゼロ値が期待値と異なります: got %d, want 0", stats.STR)
	}
	if stats.INT != 0 {
		t.Errorf("INTのゼロ値が期待値と異なります: got %d, want 0", stats.INT)
	}
	if stats.WIL != 0 {
		t.Errorf("WILのゼロ値が期待値と異なります: got %d, want 0", stats.WIL)
	}
	if stats.LUK != 0 {
		t.Errorf("LUKのゼロ値が期待値と異なります: got %d, want 0", stats.LUK)
	}
}

// TestStats_値の設定 はStats構造体に値を設定できることを確認します。
func TestStats_値の設定(t *testing.T) {
	stats := Stats{
		STR: 10,
		INT: 20,
		WIL: 15,
		LUK: 5,
	}

	if stats.STR != 10 {
		t.Errorf("STRの値が期待値と異なります: got %d, want 10", stats.STR)
	}
	if stats.INT != 20 {
		t.Errorf("INTの値が期待値と異なります: got %d, want 20", stats.INT)
	}
	if stats.WIL != 15 {
		t.Errorf("WILの値が期待値と異なります: got %d, want 15", stats.WIL)
	}
	if stats.LUK != 5 {
		t.Errorf("LUKの値が期待値と異なります: got %d, want 5", stats.LUK)
	}
}

// TestCoreType_フィールドの確認 はCoreType構造体のフィールドが正しく設定されることを確認します。
func TestCoreType_フィールドの確認(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	if coreType.ID != "attack_balance" {
		t.Errorf("IDが期待値と異なります: got %s, want attack_balance", coreType.ID)
	}
	if coreType.Name != "攻撃バランス" {
		t.Errorf("Nameが期待値と異なります: got %s, want 攻撃バランス", coreType.Name)
	}
	if coreType.StatWeights["STR"] != 1.2 {
		t.Errorf("StatWeights[STR]が期待値と異なります: got %f, want 1.2", coreType.StatWeights["STR"])
	}
	if coreType.PassiveSkillID != "balanced_stance" {
		t.Errorf("PassiveSkillIDが期待値と異なります: got %s, want balanced_stance", coreType.PassiveSkillID)
	}
	if len(coreType.AllowedTags) != 2 {
		t.Errorf("AllowedTagsの長さが期待値と異なります: got %d, want 2", len(coreType.AllowedTags))
	}
	if coreType.MinDropLevel != 1 {
		t.Errorf("MinDropLevelが期待値と異なります: got %d, want 1", coreType.MinDropLevel)
	}
}

// TestPassiveSkill_フィールドの確認 はPassiveSkill構造体のフィールドが正しく設定されることを確認します。
func TestPassiveSkill_フィールドの確認(t *testing.T) {
	skill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	if skill.ID != "balanced_stance" {
		t.Errorf("IDが期待値と異なります: got %s, want balanced_stance", skill.ID)
	}
	if skill.Name != "バランススタンス" {
		t.Errorf("Nameが期待値と異なります: got %s, want バランススタンス", skill.Name)
	}
	if skill.Description != "物理と魔法のダメージをバランスよく強化する" {
		t.Errorf("Descriptionが期待値と異なります: got %s", skill.Description)
	}
}

// TestCoreModel_フィールドの確認 はCoreModel構造体のフィールドが正しく設定されることを確認します。
func TestCoreModel_フィールドの確認(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	passiveSkill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	core := CoreModel{
		ID:           "core_001",
		TypeID:       "attack_balance",
		Name:         "攻撃バランス",
		Type:         coreType,
		Stats:        Stats{STR: 120, INT: 100, WIL: 80, LUK: 100},
		PassiveSkill: passiveSkill,
		AllowedTags:  []string{"physical_low", "magic_low"},
	}

	if core.ID != "core_001" {
		t.Errorf("IDが期待値と異なります: got %s, want core_001", core.ID)
	}
	if core.TypeID != "attack_balance" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want attack_balance", core.TypeID)
	}
	if core.Name != "攻撃バランス" {
		t.Errorf("Nameが期待値と異なります: got %s, want 攻撃バランス", core.Name)
	}
	if core.Type.ID != "attack_balance" {
		t.Errorf("Type.IDが期待値と異なります: got %s, want attack_balance", core.Type.ID)
	}
	if core.Stats.STR != 120 {
		t.Errorf("Stats.STRが期待値と異なります: got %d, want 120", core.Stats.STR)
	}
	if core.PassiveSkill.ID != "balanced_stance" {
		t.Errorf("PassiveSkill.IDが期待値と異なります: got %s, want balanced_stance", core.PassiveSkill.ID)
	}
	if len(core.AllowedTags) != 2 {
		t.Errorf("AllowedTagsの長さが期待値と異なります: got %d, want 2", len(core.AllowedTags))
	}
}

// TestCalculateStats_重みベース計算 は重みベースでステータスが正しく計算されることを確認します。
// ステータス計算式: 100 × コアの重み
func TestCalculateStats_重みベース計算(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	stats := CalculateStats(coreType)

	// 100 × 重み = 期待値
	// STR: 100 × 1.2 = 120
	// INT: 100 × 1.0 = 100
	// WIL: 100 × 0.8 = 80
	// LUK: 100 × 1.0 = 100
	if stats.STR != 120 {
		t.Errorf("STRが期待値と異なります: got %d, want 120", stats.STR)
	}
	if stats.INT != 100 {
		t.Errorf("INTが期待値と異なります: got %d, want 100", stats.INT)
	}
	if stats.WIL != 80 {
		t.Errorf("WILが期待値と異なります: got %d, want 80", stats.WIL)
	}
	if stats.LUK != 100 {
		t.Errorf("LUKが期待値と異なります: got %d, want 100", stats.LUK)
	}
}

// TestCalculateStats_ヒーラー特性_重みベース はヒーラー特性でステータスが正しく計算されることを確認します。
func TestCalculateStats_ヒーラー特性_重みベース(t *testing.T) {
	coreType := CoreType{
		ID:             "healer",
		Name:           "ヒーラー",
		StatWeights:    map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 0.8, "LUK": 1.2},
		PassiveSkillID: "healing_aura",
		AllowedTags:    []string{"heal_mid", "heal_high"},
		MinDropLevel:   3,
	}

	stats := CalculateStats(coreType)

	// 100 × 重み = 期待値
	// STR: 100 × 0.5 = 50
	// INT: 100 × 1.5 = 150
	// WIL: 100 × 0.8 = 80
	// LUK: 100 × 1.2 = 120
	if stats.STR != 50 {
		t.Errorf("STRが期待値と異なります: got %d, want 50", stats.STR)
	}
	if stats.INT != 150 {
		t.Errorf("INTが期待値と異なります: got %d, want 150", stats.INT)
	}
	if stats.WIL != 80 {
		t.Errorf("WILが期待値と異なります: got %d, want 80", stats.WIL)
	}
	if stats.LUK != 120 {
		t.Errorf("LUKが期待値と異なります: got %d, want 120", stats.LUK)
	}
}

// TestCalculateStats_オールラウンダー特性_重みベース はオールラウンダー特性でステータスが正しく計算されることを確認します。
func TestCalculateStats_オールラウンダー特性_重みベース(t *testing.T) {
	coreType := CoreType{
		ID:             "all_rounder",
		Name:           "オールラウンダー",
		StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		PassiveSkillID: "versatile",
		AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
		MinDropLevel:   1,
	}

	stats := CalculateStats(coreType)

	// 100 × 重み = 期待値
	// STR, INT, WIL, LUK: 100 × 1.0 = 100
	if stats.STR != 100 {
		t.Errorf("STRが期待値と異なります: got %d, want 100", stats.STR)
	}
	if stats.INT != 100 {
		t.Errorf("INTが期待値と異なります: got %d, want 100", stats.INT)
	}
	if stats.WIL != 100 {
		t.Errorf("WILが期待値と異なります: got %d, want 100", stats.WIL)
	}
	if stats.LUK != 100 {
		t.Errorf("LUKが期待値と異なります: got %d, want 100", stats.LUK)
	}
}

// TestCalculateStats_パラディン特性_重みベース はパラディン特性でステータスが正しく計算されることを確認します。
func TestCalculateStats_パラディン特性_重みベース(t *testing.T) {
	coreType := CoreType{
		ID:             "paladin",
		Name:           "パラディン",
		StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.1, "WIL": 0.7, "LUK": 1.2},
		PassiveSkillID: "holy_protection",
		AllowedTags:    []string{"buff_mid", "heal_low"},
		MinDropLevel:   5,
	}

	stats := CalculateStats(coreType)

	// 100 × 重み = 期待値
	// STR: 100 × 1.0 = 100
	// INT: 100 × 1.1 = 110
	// WIL: 100 × 0.7 = 70
	// LUK: 100 × 1.2 = 120
	if stats.STR != 100 {
		t.Errorf("STRが期待値と異なります: got %d, want 100", stats.STR)
	}
	if stats.INT != 110 {
		t.Errorf("INTが期待値と異なります: got %d, want 110", stats.INT)
	}
	if stats.WIL != 70 {
		t.Errorf("WILが期待値と異なります: got %d, want 70", stats.WIL)
	}
	if stats.LUK != 120 {
		t.Errorf("LUKが期待値と異なります: got %d, want 120", stats.LUK)
	}
}

// TestCalculateStats_重み未設定_重みベース は重みが設定されていないステータスが0になることを確認します。
func TestCalculateStats_重み未設定_重みベース(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		StatWeights: map[string]float64{"STR": 1.0}, // INT, WIL, LUK は未設定
	}

	stats := CalculateStats(coreType)

	// STRは設定あり: 100 × 1.0 = 100
	if stats.STR != 100 {
		t.Errorf("STRが期待値と異なります: got %d, want 100", stats.STR)
	}
	// INT, WIL, LUKは未設定のため0（重み0.0扱い）
	if stats.INT != 0 {
		t.Errorf("INTは未設定のため0になるべきです: got %d", stats.INT)
	}
	if stats.WIL != 0 {
		t.Errorf("WILは未設定のため0になるべきです: got %d", stats.WIL)
	}
	if stats.LUK != 0 {
		t.Errorf("LUKは未設定のため0になるべきです: got %d", stats.LUK)
	}
}

// TestNewCoreWithTypeID_コア作成 はNewCoreWithTypeID関数でコアが正しく作成されることを確認します。
func TestNewCoreWithTypeID_コア作成(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	passiveSkill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	if core.TypeID != "attack_balance" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want attack_balance", core.TypeID)
	}
	// Nameはコア特性名のみ（レベル表記なし）
	expectedName := "攻撃バランス"
	if core.Name != expectedName {
		t.Errorf("Nameが期待値と異なります: got %s, want %s", core.Name, expectedName)
	}

	// ステータスが重みベースで自動計算されていることを確認
	// STR: 100 × 1.2 = 120
	if core.Stats.STR != 120 {
		t.Errorf("Stats.STRが期待値と異なります: got %d, want 120", core.Stats.STR)
	}
	// INT: 100 × 1.0 = 100
	if core.Stats.INT != 100 {
		t.Errorf("Stats.INTが期待値と異なります: got %d, want 100", core.Stats.INT)
	}

	// AllowedTagsがCoreTypeからコピーされていることを確認
	if len(core.AllowedTags) != 2 {
		t.Errorf("AllowedTagsの長さが期待値と異なります: got %d, want 2", len(core.AllowedTags))
	}
}

// TestCoreModel_IsTagAllowed_許可タグ はIsTagAllowedメソッドが許可タグを正しく判定することを確認します。
func TestCoreModel_IsTagAllowed_許可タグ(t *testing.T) {
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{"physical_low", "magic_low"},
	}

	if !core.IsTagAllowed("physical_low") {
		t.Error("physical_lowは許可タグのはずですがfalseが返されました")
	}
	if !core.IsTagAllowed("magic_low") {
		t.Error("magic_lowは許可タグのはずですがfalseが返されました")
	}
}

// TestCoreModel_IsTagAllowed_非許可タグ はIsTagAllowedメソッドが非許可タグを正しく拒否することを確認します。
func TestCoreModel_IsTagAllowed_非許可タグ(t *testing.T) {
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{"physical_low", "magic_low"},
	}

	if core.IsTagAllowed("heal_mid") {
		t.Error("heal_midは非許可タグのはずですがtrueが返されました")
	}
	if core.IsTagAllowed("buff_high") {
		t.Error("buff_highは非許可タグのはずですがtrueが返されました")
	}
}

// TestStats_Total はStats構造体の合計値を計算できることを確認します。
func TestStats_Total(t *testing.T) {
	stats := Stats{
		STR: 10,
		INT: 20,
		WIL: 15,
		LUK: 5,
	}

	total := stats.Total()
	expected := 50

	if total != expected {
		t.Errorf("Totalが期待値と異なります: got %d, want %d", total, expected)
	}
}

// TestCoreModel_IsTagAllowed_空リスト はAllowedTagsが空の場合に常にfalseを返すことを確認します。
func TestCoreModel_IsTagAllowed_空リスト(t *testing.T) {
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{},
	}

	if core.IsTagAllowed("physical_low") {
		t.Error("AllowedTagsが空の場合、falseを返すべきです")
	}
}

// TestCoreModel_IsTagAllowed_空文字タグ は空文字のタグに対して正しく判定することを確認します。
func TestCoreModel_IsTagAllowed_空文字タグ(t *testing.T) {
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{"physical_low", ""},
	}

	// 空文字がAllowedTagsに含まれている場合
	if !core.IsTagAllowed("") {
		t.Error("空文字タグがAllowedTagsに含まれているため、trueを返すべきです")
	}

	// 許可タグを持たない場合
	core2 := CoreModel{
		ID:          "core_002",
		AllowedTags: []string{"physical_low"},
	}
	if core2.IsTagAllowed("") {
		t.Error("空文字タグがAllowedTagsに含まれていないため、falseを返すべきです")
	}
}

// TestNewCoreWithTypeID_AllowedTagsの独立性 はNewCoreWithTypeIDで作成したコアのAllowedTagsが元のCoreTypeと独立していることを確認します。
func TestNewCoreWithTypeID_AllowedTagsの独立性(t *testing.T) {
	coreType := CoreType{
		ID:          "test",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"tag1", "tag2"},
	}

	passiveSkill := PassiveSkill{ID: "test_skill"}

	core := NewCoreWithTypeID("test", coreType, passiveSkill)

	// 元のAllowedTagsを変更
	coreType.AllowedTags[0] = "modified_tag"

	// CoreModelのAllowedTagsは影響を受けないはず
	if core.AllowedTags[0] != "tag1" {
		t.Errorf("CoreModelのAllowedTagsが変更されています: got %s, want tag1", core.AllowedTags[0])
	}
}

// TestBaseStatValue はBaseStatValue定数が正しい値であることを確認します。
func TestBaseStatValue(t *testing.T) {
	// 新仕様: 基礎値は100
	if BaseStatValue != 100 {
		t.Errorf("BaseStatValueが期待値と異なります: got %d, want 100", BaseStatValue)
	}
}

// ==================== CoreModel TypeIDベースリファクタリングテスト ====================

// TestCoreModel_TypeIDフィールドの確認 はCoreModelにTypeIDフィールドが存在することを確認します。
func TestCoreModel_TypeIDフィールドの確認(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
		MinDropLevel:   1,
	}

	passiveSkill := PassiveSkill{
		ID:          "balanced_stance",
		Name:        "バランススタンス",
		Description: "物理と魔法のダメージをバランスよく強化する",
	}

	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	if core.TypeID != "attack_balance" {
		t.Errorf("TypeIDが期待値と異なります: got %s, want attack_balance", core.TypeID)
	}
	// Nameはコア特性名のみ（レベル表記なし）
	expectedName := "攻撃バランス"
	if core.Name != expectedName {
		t.Errorf("Nameが期待値と異なります: got %s, want %s", core.Name, expectedName)
	}
}

// TestCoreModel_Equals_同一性判定 はEqualsメソッドがTypeIDのみで同一性を判定することを確認します。
func TestCoreModel_Equals_同一性判定(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
	}
	passiveSkill := PassiveSkill{ID: "balanced_stance"}

	healerType := CoreType{
		ID:             "healer",
		Name:           "ヒーラー",
		StatWeights:    map[string]float64{"STR": 0.5, "INT": 1.5, "WIL": 0.8, "LUK": 1.2},
		PassiveSkillID: "healing_aura",
		AllowedTags:    []string{"heal_mid", "heal_high"},
	}

	core1 := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)
	core2 := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)
	core3 := NewCoreWithTypeID("healer", healerType, passiveSkill)

	if !core1.Equals(core2) {
		t.Error("同じTypeIDのコアは等価であるべきです")
	}

	if core1.Equals(core3) {
		t.Error("異なるTypeIDのコアは等価でないべきです")
	}
}

// TestCoreModel_Equals_nilチェック はEqualsメソッドがnilを正しく処理することを確認します。
func TestCoreModel_Equals_nilチェック(t *testing.T) {
	coreType := CoreType{
		ID:          "attack_balance",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
	}
	passiveSkill := PassiveSkill{ID: "test"}
	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	if core.Equals(nil) {
		t.Error("nilとの比較はfalseを返すべきです")
	}
}

// TestCoreModel_ステータス重みベース計算 はNewCoreWithTypeIDでステータスが重みベースで正しく計算されることを確認します。
func TestCoreModel_ステータス重みベース計算(t *testing.T) {
	coreType := CoreType{
		ID:             "attack_balance",
		Name:           "攻撃バランス",
		StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.0, "WIL": 0.8, "LUK": 1.0},
		PassiveSkillID: "balanced_stance",
		AllowedTags:    []string{"physical_low", "magic_low"},
	}
	passiveSkill := PassiveSkill{ID: "balanced_stance"}

	core := NewCoreWithTypeID("attack_balance", coreType, passiveSkill)

	// ステータスが重みベースで正しく計算されていることを確認
	// STR: 100 × 1.2 = 120
	if core.Stats.STR != 120 {
		t.Errorf("Stats.STRが期待値と異なります: got %d, want 120", core.Stats.STR)
	}
	// INT: 100 × 1.0 = 100
	if core.Stats.INT != 100 {
		t.Errorf("Stats.INTが期待値と異なります: got %d, want 100", core.Stats.INT)
	}
	// WIL: 100 × 0.8 = 80
	if core.Stats.WIL != 80 {
		t.Errorf("Stats.WILが期待値と異なります: got %d, want 80", core.Stats.WIL)
	}
	// LUK: 100 × 1.0 = 100
	if core.Stats.LUK != 100 {
		t.Errorf("Stats.LUKが期待値と異なります: got %d, want 100", core.Stats.LUK)
	}
}
