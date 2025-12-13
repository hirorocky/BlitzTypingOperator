// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestModuleCategory_定数の確認 はModuleCategory定数が正しく定義されていることを確認します。
func TestModuleCategory_定数の確認(t *testing.T) {
	tests := []struct {
		name     string
		category ModuleCategory
		expected string
	}{
		{"物理攻撃", PhysicalAttack, "physical_attack"},
		{"魔法攻撃", MagicAttack, "magic_attack"},
		{"回復", Heal, "heal"},
		{"バフ", Buff, "buff"},
		{"デバフ", Debuff, "debuff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.category) != tt.expected {
				t.Errorf("ModuleCategoryが期待値と異なります: got %s, want %s", tt.category, tt.expected)
			}
		})
	}
}

// TestModuleCategory_String はModuleCategoryのString()メソッドが正しい日本語名を返すことを確認します。
func TestModuleCategory_String(t *testing.T) {
	tests := []struct {
		category ModuleCategory
		expected string
	}{
		{PhysicalAttack, "物理攻撃"},
		{MagicAttack, "魔法攻撃"},
		{Heal, "回復"},
		{Buff, "バフ"},
		{Debuff, "デバフ"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.category.String() != tt.expected {
				t.Errorf("String()が期待値と異なります: got %s, want %s", tt.category.String(), tt.expected)
			}
		})
	}
}

// TestModuleModel_フィールドの確認 はModuleModel構造体のフィールドが正しく設定されることを確認します。
func TestModuleModel_フィールドの確認(t *testing.T) {
	module := ModuleModel{
		ID:          "fireball_lv1",
		Name:        "ファイアボール",
		Category:    MagicAttack,
		Level:       1,
		Tags:        []string{"magic_low"},
		BaseEffect:  10.0,
		StatRef:     "MAG",
		Description: "炎の魔法で敵に魔法ダメージを与える",
	}

	if module.ID != "fireball_lv1" {
		t.Errorf("IDが期待値と異なります: got %s, want fireball_lv1", module.ID)
	}
	if module.Name != "ファイアボール" {
		t.Errorf("Nameが期待値と異なります: got %s, want ファイアボール", module.Name)
	}
	if module.Category != MagicAttack {
		t.Errorf("Categoryが期待値と異なります: got %s, want magic_attack", module.Category)
	}
	if module.Level != 1 {
		t.Errorf("Levelが期待値と異なります: got %d, want 1", module.Level)
	}
	if len(module.Tags) != 1 || module.Tags[0] != "magic_low" {
		t.Errorf("Tagsが期待値と異なります: got %v, want [magic_low]", module.Tags)
	}
	if module.BaseEffect != 10.0 {
		t.Errorf("BaseEffectが期待値と異なります: got %f, want 10.0", module.BaseEffect)
	}
	if module.StatRef != "MAG" {
		t.Errorf("StatRefが期待値と異なります: got %s, want MAG", module.StatRef)
	}
	if module.Description != "炎の魔法で敵に魔法ダメージを与える" {
		t.Errorf("Descriptionが期待値と異なります: got %s", module.Description)
	}
}

// TestNewModule_モジュール作成 はNewModule関数でモジュールが正しく作成されることを確認します。
func TestNewModule_モジュール作成(t *testing.T) {
	module := NewModule(
		"physical_attack_lv1",
		"物理打撃",
		PhysicalAttack,
		1,
		[]string{"physical_low"},
		10.0,
		"STR",
		"物理攻撃で敵にダメージを与える",
	)

	if module.ID != "physical_attack_lv1" {
		t.Errorf("IDが期待値と異なります: got %s, want physical_attack_lv1", module.ID)
	}
	if module.Name != "物理打撃" {
		t.Errorf("Nameが期待値と異なります: got %s, want 物理打撃", module.Name)
	}
	if module.Category != PhysicalAttack {
		t.Errorf("Categoryが期待値と異なります: got %s, want physical_attack", module.Category)
	}
	if module.Level != 1 {
		t.Errorf("Levelが期待値と異なります: got %d, want 1", module.Level)
	}
	if module.BaseEffect != 10.0 {
		t.Errorf("BaseEffectが期待値と異なります: got %f, want 10.0", module.BaseEffect)
	}
	if module.StatRef != "STR" {
		t.Errorf("StatRefが期待値と異なります: got %s, want STR", module.StatRef)
	}
}

// TestNewModule_タグのコピー はNewModuleで作成したモジュールのTagsが元のスライスと独立していることを確認します。
func TestNewModule_タグのコピー(t *testing.T) {
	originalTags := []string{"magic_low", "fire"}
	module := NewModule(
		"fireball_lv1",
		"ファイアボール",
		MagicAttack,
		1,
		originalTags,
		10.0,
		"MAG",
		"炎の魔法で敵に魔法ダメージを与える",
	)

	// 元のタグを変更
	originalTags[0] = "modified_tag"

	// モジュールのTagsは影響を受けないはず
	if module.Tags[0] != "magic_low" {
		t.Errorf("ModuleModelのTagsが変更されています: got %s, want magic_low", module.Tags[0])
	}
}

// TestModuleModel_HasTag_タグ存在確認 はHasTagメソッドがタグの存在を正しく判定することを確認します。
func TestModuleModel_HasTag_タグ存在確認(t *testing.T) {
	module := ModuleModel{
		ID:   "test_module",
		Tags: []string{"physical_low", "fire"},
	}

	if !module.HasTag("physical_low") {
		t.Error("physical_lowタグが存在するはずですがfalseが返されました")
	}
	if !module.HasTag("fire") {
		t.Error("fireタグが存在するはずですがfalseが返されました")
	}
	if module.HasTag("magic_low") {
		t.Error("magic_lowタグは存在しないはずですがtrueが返されました")
	}
}

// TestModuleModel_HasTag_空タグリスト はTagsが空の場合に常にfalseを返すことを確認します。
func TestModuleModel_HasTag_空タグリスト(t *testing.T) {
	module := ModuleModel{
		ID:   "test_module",
		Tags: []string{},
	}

	if module.HasTag("physical_low") {
		t.Error("Tagsが空の場合、falseを返すべきです")
	}
}

// TestModuleModel_IsCompatibleWithCore はモジュールがコアに装備可能かを判定するメソッドをテストします。
func TestModuleModel_IsCompatibleWithCore(t *testing.T) {
	// 物理攻撃と魔法攻撃の低レベルモジュールを許可するコア
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{"physical_low", "magic_low"},
	}

	// 互換性のあるモジュール
	compatibleModule := ModuleModel{
		ID:   "physical_attack_lv1",
		Tags: []string{"physical_low"},
	}

	// 互換性のないモジュール
	incompatibleModule := ModuleModel{
		ID:   "heal_lv2",
		Tags: []string{"heal_mid"},
	}

	if !compatibleModule.IsCompatibleWithCore(&core) {
		t.Error("physical_lowタグを持つモジュールはコアと互換性があるはずです")
	}

	if incompatibleModule.IsCompatibleWithCore(&core) {
		t.Error("heal_midタグを持つモジュールはコアと互換性がないはずです")
	}
}

// TestModuleModel_IsCompatibleWithCore_複数タグ はモジュールが複数タグを持つ場合の互換性判定をテストします。
func TestModuleModel_IsCompatibleWithCore_複数タグ(t *testing.T) {
	core := CoreModel{
		ID:          "core_001",
		AllowedTags: []string{"physical_low", "magic_low"},
	}

	// 複数タグのうち1つがコアの許可タグに含まれる場合
	moduleWithMultipleTags := ModuleModel{
		ID:   "hybrid_lv1",
		Tags: []string{"physical_low", "fire"},
	}

	if !moduleWithMultipleTags.IsCompatibleWithCore(&core) {
		t.Error("1つでもコアの許可タグに含まれるタグがあれば互換性があるはずです")
	}

	// どのタグもコアの許可タグに含まれない場合
	moduleNoMatch := ModuleModel{
		ID:   "heal_lv1",
		Tags: []string{"heal_low", "light"},
	}

	if moduleNoMatch.IsCompatibleWithCore(&core) {
		t.Error("どのタグもコアの許可タグに含まれない場合、互換性がないはずです")
	}
}

// TestModuleCategory_Unknown_String は未知のカテゴリに対してString()が適切な値を返すことを確認します。
func TestModuleCategory_Unknown_String(t *testing.T) {
	unknownCategory := ModuleCategory("unknown")
	result := unknownCategory.String()
	if result != "不明" {
		t.Errorf("未知のカテゴリに対するString()が期待値と異なります: got %s, want 不明", result)
	}
}

// TestGetLevelSuffix はレベルに応じた接尾辞が正しく返されることを確認します。
func TestGetLevelSuffix(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{1, "low"},
		{2, "mid"},
		{3, "high"},
		{4, "high"}, // 3以上はすべてhigh
		{10, "high"},
		{0, "low"},  // 0以下はlow
		{-1, "low"}, // 負の値はlow
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.level)), func(t *testing.T) {
			result := GetLevelSuffix(tt.level)
			if result != tt.expected {
				t.Errorf("GetLevelSuffix(%d)が期待値と異なります: got %s, want %s", tt.level, result, tt.expected)
			}
		})
	}
}

// ==================== Task 7.2: Icon()メソッドのテスト ====================

// TestModuleCategory_Icon はModuleCategoryのIcon()メソッドが正しいアイコンを返すことを確認します。
// 要件 7.3: カテゴリごとにアイコン文字を返す
func TestModuleCategory_Icon(t *testing.T) {
	tests := []struct {
		category ModuleCategory
		expected string
	}{
		{PhysicalAttack, "⚔"},
		{MagicAttack, "✦"},
		{Heal, "♥"},
		{Buff, "▲"},
		{Debuff, "▼"},
	}

	for _, tt := range tests {
		t.Run(tt.category.String(), func(t *testing.T) {
			result := tt.category.Icon()
			if result != tt.expected {
				t.Errorf("Icon()が期待値と異なります: got %s, want %s", result, tt.expected)
			}
		})
	}
}

// TestModuleCategory_Icon_Unknown は未知のカテゴリに対してIcon()がデフォルト値を返すことを確認します。
func TestModuleCategory_Icon_Unknown(t *testing.T) {
	unknownCategory := ModuleCategory("unknown")
	result := unknownCategory.Icon()
	if result != "•" {
		t.Errorf("未知のカテゴリに対するIcon()が期待値と異なります: got %s, want •", result)
	}
}

// TestModuleCategory_Icon_Empty は空のカテゴリに対してIcon()がデフォルト値を返すことを確認します。
func TestModuleCategory_Icon_Empty(t *testing.T) {
	emptyCategory := ModuleCategory("")
	result := emptyCategory.Icon()
	if result != "•" {
		t.Errorf("空のカテゴリに対するIcon()が期待値と異なります: got %s, want •", result)
	}
}
