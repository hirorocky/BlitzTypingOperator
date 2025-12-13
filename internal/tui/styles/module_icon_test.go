// Package styles はTUIスタイリングのテストを提供します。

package styles

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestGetModuleIcon はモジュールアイコンの取得をテストします。

func TestGetModuleIcon(t *testing.T) {
	tests := []struct {
		category domain.ModuleCategory
		expected string
	}{
		{domain.PhysicalAttack, "⚔"},
		{domain.MagicAttack, "✦"},
		{domain.Heal, "♥"},
		{domain.Buff, "▲"},
		{domain.Debuff, "▼"},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			icon := GetModuleIcon(tt.category)
			if icon != tt.expected {
				t.Errorf("GetModuleIcon(%s)が正しくありません: got %s, want %s", tt.category, icon, tt.expected)
			}
		})
	}
}

// TestGetModuleIconUnknown は不明なカテゴリの処理をテストします。
func TestGetModuleIconUnknown(t *testing.T) {
	icon := GetModuleIcon(domain.ModuleCategory("unknown"))
	if icon == "" {
		t.Error("不明なカテゴリで空のアイコンが返されました")
	}
}

// TestGetModuleIconColored はカラー付きアイコンの取得をテストします。
func TestGetModuleIconColored(t *testing.T) {
	gs := NewGameStyles()

	tests := []domain.ModuleCategory{
		domain.PhysicalAttack,
		domain.MagicAttack,
		domain.Heal,
		domain.Buff,
		domain.Debuff,
	}

	for _, category := range tests {
		t.Run(string(category), func(t *testing.T) {
			icon := GetModuleIconColored(category, gs)
			if icon == "" {
				t.Errorf("GetModuleIconColored(%s)が空文字列を返しました", category)
			}
		})
	}
}

// TestGetModuleIconsForAgent はエージェントのモジュールアイコンリスト取得をテストします。
func TestGetModuleIconsForAgent(t *testing.T) {
	categories := []domain.ModuleCategory{
		domain.PhysicalAttack,
		domain.PhysicalAttack,
		domain.Buff,
		domain.Heal,
	}

	icons := GetModuleIcons(categories)

	if len(icons) != 4 {
		t.Errorf("アイコン数が正しくありません: got %d, want 4", len(icons))
	}

	// 最初の2つは物理攻撃アイコン
	if icons[0] != "⚔" || icons[1] != "⚔" {
		t.Error("物理攻撃アイコンが正しくありません")
	}

	// 3番目はバフアイコン
	if icons[2] != "▲" {
		t.Error("バフアイコンが正しくありません")
	}

	// 4番目は回復アイコン
	if icons[3] != "♥" {
		t.Error("回復アイコンが正しくありません")
	}
}

// TestModuleIconMapping はアイコンマッピングの一貫性をテストします。
func TestModuleIconMapping(t *testing.T) {
	// 各カテゴリに対応するアイコンが一意であることを確認
	allCategories := []domain.ModuleCategory{
		domain.PhysicalAttack,
		domain.MagicAttack,
		domain.Heal,
		domain.Buff,
		domain.Debuff,
	}

	icons := make(map[string]domain.ModuleCategory)
	for _, cat := range allCategories {
		icon := GetModuleIcon(cat)
		if existing, ok := icons[icon]; ok {
			t.Errorf("アイコン %s が %s と %s で重複しています", icon, existing, cat)
		}
		icons[icon] = cat
	}
}
