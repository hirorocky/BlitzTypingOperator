// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestPlayerModel_フィールドの確認 はPlayerModel構造体のフィールドが正しく設定されることを確認します。
func TestPlayerModel_フィールドの確認(t *testing.T) {
	player := PlayerModel{
		HP:    100,
		MaxHP: 100,
	}

	if player.HP != 100 {
		t.Errorf("HPが期待値と異なります: got %d, want 100", player.HP)
	}
	if player.MaxHP != 100 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 100", player.MaxHP)
	}
}

// TestNewPlayerWithMaxHP_プレイヤー作成 はNewPlayerWithMaxHP関数でプレイヤーが正しく作成されることを確認します。
func TestNewPlayerWithMaxHP_プレイヤー作成(t *testing.T) {
	player := NewPlayerWithMaxHP(0)

	// 初期状態ではHP/MaxHPは0（エージェント装備後に計算）
	if player.HP != 0 {
		t.Errorf("初期HPが期待値と異なります: got %d, want 0", player.HP)
	}
	if player.MaxHP != 0 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 0", player.MaxHP)
	}
	if player.EffectTable == nil {
		t.Error("EffectTableがnilです")
	}
}

// TestPlayerModel_InitializeHP はHP初期化を確認します。
func TestPlayerModel_InitializeHP(t *testing.T) {
	player := NewPlayerWithMaxHP(0)

	// 初期状態
	if player.MaxHP != 0 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 0", player.MaxHP)
	}

	// HP初期化
	player.InitializeHP(InitialMaxHP)

	// 初期最大HP = 1000
	if player.MaxHP != 1000 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 1000", player.MaxHP)
	}
	if player.HP != 1000 {
		t.Errorf("HPも最大値に設定されるべき: got %d, want 1000", player.HP)
	}
}

// TestPlayerModel_バトル開始時全回復 はバトル開始時にHPが全回復することを確認します。

func TestPlayerModel_バトル開始時全回復(t *testing.T) {
	player := NewPlayerWithMaxHP(InitialMaxHP)

	// ダメージを受けた状態にする
	player.HP = 50

	// バトル開始時の処理
	player.FullHeal()

	if player.HP != player.MaxHP {
		t.Errorf("HPが全回復していません: got %d, want %d", player.HP, player.MaxHP)
	}
}

// TestPlayerModel_HP増減 はHPの増減処理を確認します。
func TestPlayerModel_HP増減(t *testing.T) {
	player := NewPlayerWithMaxHP(InitialMaxHP)

	// ダメージを受ける
	player.TakeDamage(30)
	if player.HP != 970 {
		t.Errorf("HP減少後の値が期待値と異なります: got %d, want 970", player.HP)
	}

	// 回復
	player.Heal(20)
	if player.HP != 990 {
		t.Errorf("HP回復後の値が期待値と異なります: got %d, want 990", player.HP)
	}

	// 過剰回復（MaxHPを超えない）
	player.Heal(100)
	if player.HP != 1000 {
		t.Errorf("HPがMaxHPを超えています: got %d, want 1000", player.HP)
	}

	// 致死ダメージ（HPは0以下にならない）
	player.TakeDamage(1500)
	if player.HP != 0 {
		t.Errorf("HPが0未満になっています: got %d, want 0", player.HP)
	}
}

// TestPlayerModel_生存確認 はプレイヤーの生存確認を確認します。
func TestPlayerModel_生存確認(t *testing.T) {
	player := NewPlayerWithMaxHP(InitialMaxHP)

	// 生存状態
	if !player.IsAlive() {
		t.Error("HPが0より大きい場合は生存しているはずです")
	}

	// 死亡状態
	player.HP = 0
	if player.IsAlive() {
		t.Error("HP=0の場合は死亡しているはずです")
	}
}

// TestPlayerModel_バトル持ち越しなし はHPがバトル間で持ち越されないことを確認します。

func TestPlayerModel_バトル持ち越しなし(t *testing.T) {
	player := NewPlayerWithMaxHP(InitialMaxHP)

	// 前のバトルでダメージを受けた
	player.HP = 30

	// 新しいバトル開始
	player.PrepareForBattle()

	// HPは全回復しているはず
	if player.HP != player.MaxHP {
		t.Errorf("バトル開始時にHPが全回復していません: got %d, want %d", player.HP, player.MaxHP)
	}
}

// TestInitialMaxHP は初期最大HP定数が正しい値であることを確認します。
func TestInitialMaxHP(t *testing.T) {
	// 初期最大HPは1000
	if InitialMaxHP != 1000 {
		t.Errorf("InitialMaxHPが期待値と異なります: got %d, want 1000", InitialMaxHP)
	}
}

// TestNewPlayerWithMaxHP_新規プレイヤーのMaxHP はNewPlayerWithMaxHPで初期MaxHPが正しく設定されることを確認します。
func TestNewPlayerWithMaxHP_新規プレイヤーのMaxHP(t *testing.T) {
	// 新規プレイヤーは初期最大HP（1000）で作成される
	player := NewPlayerWithMaxHP(InitialMaxHP)

	if player.MaxHP != 1000 {
		t.Errorf("初期MaxHPが期待値と異なります: got %d, want 1000", player.MaxHP)
	}
	if player.HP != 1000 {
		t.Errorf("初期HPが期待値と異なります: got %d, want 1000", player.HP)
	}
}

// TestPlayerModel_IncreaseMaxHP は敵撃破による最大HP増加を確認します。
func TestPlayerModel_IncreaseMaxHP(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)

	// 最大HPを10増加
	player.IncreaseMaxHP(10)

	if player.MaxHP != 1010 {
		t.Errorf("IncreaseMaxHP後のMaxHPが期待値と異なります: got %d, want 1010", player.MaxHP)
	}
	// 現在HPは変更されない
	if player.HP != 1000 {
		t.Errorf("IncreaseMaxHP後のHPが変更されています: got %d, want 1000", player.HP)
	}
}

// TestPlayerModel_IncreaseMaxHP_複数回増加 は複数回の最大HP増加を確認します。
func TestPlayerModel_IncreaseMaxHP_複数回増加(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)

	// 初撃破報酬: +10
	player.IncreaseMaxHP(10)
	// 高レベル撃破報酬: (5-1) x 10 = +40
	player.IncreaseMaxHP(40)

	if player.MaxHP != 1050 {
		t.Errorf("複数回IncreaseMaxHP後のMaxHPが期待値と異なります: got %d, want 1050", player.MaxHP)
	}
}

// TestPlayerModel_IncreaseMaxHP_負の値は無視 は負の値がIncreaseMaxHPに渡されても無視されることを確認します。
func TestPlayerModel_IncreaseMaxHP_負の値は無視(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)

	player.IncreaseMaxHP(-50)

	// MaxHPは減少しない
	if player.MaxHP != 1000 {
		t.Errorf("負の値でMaxHPが変更されています: got %d, want 1000", player.MaxHP)
	}
}

// TestPlayerModel_IncreaseMaxHP_ゼロは無視 はゼロがIncreaseMaxHPに渡されても何も変わらないことを確認します。
func TestPlayerModel_IncreaseMaxHP_ゼロは無視(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)

	player.IncreaseMaxHP(0)

	if player.MaxHP != 1000 {
		t.Errorf("ゼロでMaxHPが変更されています: got %d, want 1000", player.MaxHP)
	}
}
