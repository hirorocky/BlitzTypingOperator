// Package components はTUI共通コンポーネントを提供します。
// hp_display_test.go はHP表示コンポーネントのテストを含みます。
package components

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/tui/styles"
)

// ==================== RenderHP関数のテスト ====================

// TestRenderHP_Basic は基本的なHP表示をテストします。
// 要件 7.1: RenderHP()関数でHPバー、数値表示、色分けロジックを共通化する
func TestRenderHP_Basic(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(80, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_FullHP は満タンHPの表示をテストします。
func TestRenderHP_FullHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(100, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_ZeroHP はHP 0の表示をテストします。
func TestRenderHP_ZeroHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(0, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_NegativeHP は負のHPの表示をテストします（0として扱う）。
func TestRenderHP_NegativeHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(-10, 100, 20, gs)

	// 結果が空でないこと（負の値は0として扱われる）
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_LowHP は低HP（25%未満）の表示をテストします。
func TestRenderHP_LowHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(20, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_MediumHP は中HP（25%〜50%）の表示をテストします。
func TestRenderHP_MediumHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(30, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_HighHP は高HP（50%以上）の表示をテストします。
func TestRenderHP_HighHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHP(60, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHP()が空文字列を返しました")
	}
}

// TestRenderHP_VariousBarWidths は異なるバー幅での表示をテストします。
func TestRenderHP_VariousBarWidths(t *testing.T) {
	gs := styles.NewGameStyles()

	testCases := []struct {
		barWidth int
	}{
		{barWidth: 10},
		{barWidth: 20},
		{barWidth: 50},
	}

	for _, tc := range testCases {
		result := RenderHP(50, 100, tc.barWidth, gs)
		if result == "" {
			t.Errorf("RenderHP()がbarWidth=%dで空文字列を返しました", tc.barWidth)
		}
	}
}

// ==================== RenderHPWithLabel関数のテスト ====================

// TestRenderHPWithLabel_Basic は基本的なラベル付きHP表示をテストします。
// 要件 7.2: RenderHPWithLabel()でラベル付き表示をサポートする
func TestRenderHPWithLabel_Basic(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHPWithLabel("HP", 80, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHPWithLabel()が空文字列を返しました")
	}

	// ラベルが含まれていること
	if !strings.Contains(result, "HP") {
		t.Error("RenderHPWithLabel()の結果にラベルが含まれていません")
	}
}

// TestRenderHPWithLabel_CustomLabel はカスタムラベルの表示をテストします。
func TestRenderHPWithLabel_CustomLabel(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHPWithLabel("プレイヤーHP", 50, 100, 20, gs)

	// ラベルが含まれていること
	if !strings.Contains(result, "プレイヤーHP") {
		t.Error("RenderHPWithLabel()の結果にカスタムラベルが含まれていません")
	}
}

// TestRenderHPWithLabel_EmptyLabel は空ラベルの表示をテストします。
func TestRenderHPWithLabel_EmptyLabel(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHPWithLabel("", 80, 100, 20, gs)

	// 空ラベルでも動作すること
	if result == "" {
		t.Error("RenderHPWithLabel()が空ラベルで空文字列を返しました")
	}
}

// TestRenderHPWithLabel_IncludesValue はHP値が含まれることをテストします。
func TestRenderHPWithLabel_IncludesValue(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHPWithLabel("HP", 75, 100, 20, gs)

	// HP値が含まれていること
	if !strings.Contains(result, "75") {
		t.Error("RenderHPWithLabel()の結果に現在HP値が含まれていません")
	}
	if !strings.Contains(result, "100") {
		t.Error("RenderHPWithLabel()の結果に最大HP値が含まれていません")
	}
}

// TestRenderHPWithLabel_ZeroHP はHP 0のラベル付き表示をテストします。
func TestRenderHPWithLabel_ZeroHP(t *testing.T) {
	gs := styles.NewGameStyles()

	result := RenderHPWithLabel("HP", 0, 100, 20, gs)

	// 結果が空でないこと
	if result == "" {
		t.Error("RenderHPWithLabel()がHP=0で空文字列を返しました")
	}
}

// TestRenderHPWithLabel_NilStyles はスタイルがnilの場合のテストです。
func TestRenderHPWithLabel_NilStyles(t *testing.T) {
	// nilスタイルでパニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderHPWithLabel()がnilスタイルでパニックしました: %v", r)
		}
	}()

	result := RenderHPWithLabel("HP", 80, 100, 20, nil)

	// nilの場合でも何かしら返すか、空文字列を返す
	// パニックしなければOK
	_ = result
}

// TestRenderHP_NilStyles はRenderHPでスタイルがnilの場合のテストです。
func TestRenderHP_NilStyles(t *testing.T) {
	// nilスタイルでパニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderHP()がnilスタイルでパニックしました: %v", r)
		}
	}()

	result := RenderHP(80, 100, 20, nil)

	// nilの場合でも何かしら返すか、空文字列を返す
	// パニックしなければOK
	_ = result
}

// TestRenderHP_ZeroMaxHP は最大HP 0の場合のテストです。
func TestRenderHP_ZeroMaxHP(t *testing.T) {
	gs := styles.NewGameStyles()

	// ゼロ除算を避けて動作すること
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderHP()が最大HP=0でパニックしました: %v", r)
		}
	}()

	result := RenderHP(0, 0, 20, gs)

	// パニックせず何かを返すこと
	_ = result
}

// TestRenderHPWithLabel_ZeroMaxHP はラベル付きで最大HP 0の場合のテストです。
func TestRenderHPWithLabel_ZeroMaxHP(t *testing.T) {
	gs := styles.NewGameStyles()

	// ゼロ除算を避けて動作すること
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderHPWithLabel()が最大HP=0でパニックしました: %v", r)
		}
	}()

	result := RenderHPWithLabel("HP", 0, 0, 20, gs)

	// パニックせず何かを返すこと
	_ = result
}
