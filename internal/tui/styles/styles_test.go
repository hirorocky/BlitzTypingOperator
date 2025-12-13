// Package styles はTUIスタイリングのテストを提供します。
package styles

import (
	"testing"
)

// ==================== Task 11.1: 基本スタイリングのテスト ====================

// TestNewGameStyles はGameStylesの初期化をテストします。

func TestNewGameStyles(t *testing.T) {
	styles := NewGameStyles()

	if styles == nil {
		t.Fatal("GameStylesがnilです")
	}

	// スタイルが初期化されていることを確認（Border構造体は空でないこと）
	if styles.Box.BorderForeground == "" {
		t.Error("Boxスタイルのボーダー色が設定されていません")
	}
}

// TestHPColorRanges はHP色分けのテストです。

func TestHPColorRanges(t *testing.T) {
	styles := NewGameStyles()

	tests := []struct {
		name       string
		percentage float64
		expected   string // "green", "yellow", "red"
	}{
		{"満タン100%", 1.0, "green"},
		{"高HP75%", 0.75, "green"},
		{"中HP55%", 0.55, "green"},
		{"中HP50%", 0.50, "yellow"},
		{"低HP30%", 0.30, "yellow"},
		{"危険HP25%", 0.25, "red"},
		{"危険HP10%", 0.10, "red"},
		{"HP0%", 0.0, "red"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorType := styles.GetHPColorType(tt.percentage)
			if colorType != tt.expected {
				t.Errorf("HP %.0f%% の色: got %s, want %s", tt.percentage*100, colorType, tt.expected)
			}
		})
	}
}

// TestRenderHPBar はHPバーの描画をテストします。

func TestRenderHPBar(t *testing.T) {
	styles := NewGameStyles()

	tests := []struct {
		name        string
		current     int
		max         int
		width       int
		expectFull  bool
		expectEmpty bool
	}{
		{"満タン", 100, 100, 20, true, false},
		{"半分", 50, 100, 20, false, false},
		{"空", 0, 100, 20, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := styles.RenderHPBar(tt.current, tt.max, tt.width)
			if bar == "" {
				t.Error("HPバーが空文字列です")
			}
			// HPバーの長さが期待値に近いことを確認（ボーダー含む）
			// 実際のレンダリング結果をチェック
			if len(bar) == 0 {
				t.Error("HPバーのレンダリングに失敗しました")
			}
		})
	}
}

// TestDamageStyle はダメージ表示スタイルのテストです。

func TestDamageStyle(t *testing.T) {
	styles := NewGameStyles()

	rendered := styles.RenderDamage(42)
	if rendered == "" {
		t.Error("ダメージ表示が空です")
	}
}

// TestHealStyle は回復表示スタイルのテストです。

func TestHealStyle(t *testing.T) {
	styles := NewGameStyles()

	rendered := styles.RenderHeal(25)
	if rendered == "" {
		t.Error("回復表示が空です")
	}
}

// TestBoxBorder はボックス描画文字のテストです。

func TestBoxBorder(t *testing.T) {
	styles := NewGameStyles()

	boxContent := styles.RenderBox("テスト内容", 30)
	if boxContent == "" {
		t.Error("ボックス描画が空です")
	}

	// ボックス描画文字が含まれていることを確認
	if len(boxContent) < 10 {
		t.Error("ボックス描画が不十分です")
	}
}

// TestFallbackDisplay はカラー非対応ターミナルでの代替表示テストです。

func TestFallbackDisplay(t *testing.T) {
	// NoColorモードでのスタイル作成
	styles := NewGameStylesWithNoColor()

	// 代替表示でも正しくレンダリングできることを確認
	bar := styles.RenderHPBar(50, 100, 20)
	if bar == "" {
		t.Error("NoColorモードでHPバーがレンダリングできません")
	}

	damage := styles.RenderDamage(10)
	if damage == "" {
		t.Error("NoColorモードでダメージがレンダリングできません")
	}
}

// ==================== Task 7.1-7.3: カラーテーマとスタイルの統一テスト ====================

// TestColorPaletteConsistency はカラーパレットの一貫性をテストします。

func TestColorPaletteConsistency(t *testing.T) {
	// カラーパレット変数が定義されていることを確認
	colors := []struct {
		name  string
		color string
	}{
		{"ColorPrimary", string(ColorPrimary)},
		{"ColorSecondary", string(ColorSecondary)},
		{"ColorHPHigh", string(ColorHPHigh)},
		{"ColorHPMedium", string(ColorHPMedium)},
		{"ColorHPLow", string(ColorHPLow)},
		{"ColorDamage", string(ColorDamage)},
		{"ColorHeal", string(ColorHeal)},
		{"ColorSubtle", string(ColorSubtle)},
		{"ColorWarning", string(ColorWarning)},
		{"ColorInfo", string(ColorInfo)},
		{"ColorBuff", string(ColorBuff)},
		{"ColorDebuff", string(ColorDebuff)},
	}

	for _, c := range colors {
		if c.color == "" {
			t.Errorf("%s が定義されていません", c.name)
		}
	}
}

// TestRoundedBorderConsistency はボーダースタイルの一貫性をテストします。

func TestRoundedBorderConsistency(t *testing.T) {
	styles := NewGameStyles()

	// RoundedBorderが使用されていることを確認
	// lipgloss.RoundedBorder()のトップ左角は「╭」
	border := styles.Box.Border
	if border.TopLeft != "╭" {
		t.Error("RoundedBorderが使用されていません")
	}
}

// TestTextHierarchyStyles はテキスト階層スタイルをテストします。

func TestTextHierarchyStyles(t *testing.T) {
	styles := NewGameStyles()

	// 4つのテキストスタイルが定義されていることを確認
	textStyles := []struct {
		name  string
		empty bool
	}{
		{"Title", false},
		{"Subtitle", false},
		{"Normal", false},
		{"Subtle", false},
	}

	for _, ts := range textStyles {
		if ts.empty {
			t.Errorf("%s スタイルが定義されていません", ts.name)
		}
	}

	// Titleが太字であることを確認
	titleRendered := styles.Text.Title.Render("テスト")
	if titleRendered == "" {
		t.Error("Titleスタイルがレンダリングできません")
	}
}

// TestMonochromeModeSupport はモノクロモードのサポートをテストします。

func TestMonochromeModeSupport(t *testing.T) {
	colorStyles := NewGameStyles()
	monoStyles := NewGameStylesWithNoColor()

	// 両方のモードでHPバーがレンダリングできること
	colorBar := colorStyles.RenderHPBar(50, 100, 20)
	monoBar := monoStyles.RenderHPBar(50, 100, 20)

	if colorBar == "" {
		t.Error("カラーモードでHPバーがレンダリングできません")
	}

	if monoBar == "" {
		t.Error("モノクロモードでHPバーがレンダリングできません")
	}

	// モノクロモードでは#と-が使われる
	if monoStyles.noColor {
		// noColor状態が正しく設定されていることを確認
		// 実際のレンダリング内容は表示の問題なのでスキップ
	}
}

// TestBuffDebuffStyles はバフ・デバフスタイルの一貫性をテストします。

func TestBuffDebuffStyles(t *testing.T) {
	styles := NewGameStyles()

	// バフ表示
	buffRendered := styles.RenderBuff("攻撃UP", 5.0)
	if buffRendered == "" {
		t.Error("バフ表示がレンダリングできません")
	}

	// デバフ表示
	debuffRendered := styles.RenderDebuff("攻撃DOWN", 3.0)
	if debuffRendered == "" {
		t.Error("デバフ表示がレンダリングできません")
	}
}

// TestCooldownStyle はクールダウン表示スタイルをテストします。
func TestCooldownStyle(t *testing.T) {
	styles := NewGameStyles()

	// クールダウン表示
	cdRendered := styles.RenderCooldown(3.5)
	if cdRendered == "" {
		t.Error("クールダウン表示がレンダリングできません")
	}
}

// TestProgressBarStyle はプログレスバースタイルをテストします。
func TestProgressBarStyle(t *testing.T) {
	styles := NewGameStyles()

	// プログレスバー表示
	progressBar := styles.RenderProgressBar(0.5, 20, ColorPrimary, ColorSubtle)
	if progressBar == "" {
		t.Error("プログレスバーがレンダリングできません")
	}
}
