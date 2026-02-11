package commons

import (
	"math/rand"
	"testing"
)

// === 受け入れテスト: 文字セットシステム ===

// ホームポジション文字セット
var homePositionChars = "asdfghjkl;"

// 受け入れ基準13: DifficultyRateに応じて4段階の文字セットから文字が生成される
func TestGenerateChars_ホームポジション_Rate50(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	chars := GenerateChars(50, rng)

	if len(chars) < 1 {
		t.Fatal("文字が生成されない")
	}

	// 全文字がホームポジションであること
	for _, ch := range chars {
		found := false
		for _, hp := range homePositionChars {
			if ch == hp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Rate=50でホームポジション外の文字 %q が生成された", ch)
		}
	}
}

func TestGenerateChars_ホームポジション_Rate74(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	chars := GenerateChars(74, rng)

	if len(chars) < 1 {
		t.Fatal("文字が生成されない")
	}

	for _, ch := range chars {
		found := false
		for _, hp := range homePositionChars {
			if ch == hp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Rate=74でホームポジション外の文字 %q が生成された", ch)
		}
	}
}

func TestGenerateChars_キーボード1行_Rate75(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	chars := GenerateChars(75, rng)

	if len(chars) < 1 {
		t.Fatal("文字が生成されない")
	}

	// 英字a-zのみ（キーボード1行混合）
	for _, ch := range chars {
		if ch < 'a' || ch > 'z' {
			// セミコロンも許容（ホームポジションに含まれるため）
			if ch != ';' {
				t.Errorf("Rate=75で期待外の文字 %q が生成された", ch)
			}
		}
	}
}

func TestGenerateChars_英数字_Rate100(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	chars := GenerateChars(100, rng)

	if len(chars) < 1 {
		t.Fatal("文字が生成されない")
	}

	// 英数字のみ（a-z, 0-9）
	for _, ch := range chars {
		isAlpha := (ch >= 'a' && ch <= 'z')
		isNum := (ch >= '0' && ch <= '9')
		if !isAlpha && !isNum {
			t.Errorf("Rate=100で英数字以外の文字 %q が生成された", ch)
		}
	}
}

func TestGenerateChars_記号込み_Rate150(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	// 複数回生成して記号が含まれる可能性を確認
	hasSymbol := false
	for i := 0; i < 20; i++ {
		chars := GenerateChars(150, rng)
		for _, ch := range chars {
			isAlpha := (ch >= 'a' && ch <= 'z')
			isNum := (ch >= '0' && ch <= '9')
			if !isAlpha && !isNum {
				hasSymbol = true
				break
			}
		}
		if hasSymbol {
			break
		}
	}

	if !hasSymbol {
		t.Error("Rate=150で20回生成しても記号が出現しない")
	}
}

// 受け入れ基準14: 文字数もDifficultyRateに連動する（最低1文字）
func TestGenerateChars_文字数は最低1(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	chars := GenerateChars(50, rng)

	if len(chars) < 1 {
		t.Error("文字数が0: 最低1文字を保証すべき")
	}
}

// 受け入れ基準29: 同一シードで同一結果が再現される
func TestGenerateChars_再現性(t *testing.T) {
	rng1 := rand.New(rand.NewSource(12345))
	chars1 := GenerateChars(100, rng1)

	rng2 := rand.New(rand.NewSource(12345))
	chars2 := GenerateChars(100, rng2)

	if string(chars1) != string(chars2) {
		t.Errorf("同一シードで異なる結果: %q vs %q", string(chars1), string(chars2))
	}
}
