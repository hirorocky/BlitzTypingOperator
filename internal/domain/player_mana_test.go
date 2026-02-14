// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// ===== 受け入れ基準 1: PlayerModelにMana/MaxManaフィールドが追加される =====

func TestPlayerModel_マナフィールドが存在する(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10

	if player.Mana != 0 {
		t.Errorf("初期Manaが期待値と異なります: got %d, want 0", player.Mana)
	}
	if player.MaxMana != 10 {
		t.Errorf("MaxManaが期待値と異なります: got %d, want 10", player.MaxMana)
	}
}

// ===== 受け入れ基準 2: バトル開始時にManaが0に初期化される =====

func TestPlayerModel_PrepareForBattle_マナが0に初期化される(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 5 // バトル前にマナがある状態

	player.PrepareForBattle()

	if player.Mana != 0 {
		t.Errorf("PrepareForBattle後のManaが期待値と異なります: got %d, want 0", player.Mana)
	}
	// MaxManaは変更されない
	if player.MaxMana != 10 {
		t.Errorf("PrepareForBattle後のMaxManaが変更されています: got %d, want 10", player.MaxMana)
	}
}

// ===== 受け入れ基準 3: マナ回復時にMaxManaを超えないようクランプされる =====

func TestPlayerModel_GainMana_クランプされる(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 8

	player.GainMana(5) // 8 + 5 = 13 → 10にクランプ

	if player.Mana != 10 {
		t.Errorf("GainMana後のManaが期待値と異なります: got %d, want 10", player.Mana)
	}
}

func TestPlayerModel_GainMana_正常に加算される(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 3

	player.GainMana(2)

	if player.Mana != 5 {
		t.Errorf("GainMana後のManaが期待値と異なります: got %d, want 5", player.Mana)
	}
}

func TestPlayerModel_GainMana_MaxManaぴったり(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 7

	player.GainMana(3) // 7 + 3 = 10（ちょうどMaxMana）

	if player.Mana != 10 {
		t.Errorf("GainMana後のManaが期待値と異なります: got %d, want 10", player.Mana)
	}
}

// ===== 受け入れ基準 4: マナ消費時に0未満にならない =====

func TestPlayerModel_ConsumeMana_正常に消費される(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 5

	ok := player.ConsumeMana(3)

	if !ok {
		t.Error("マナ消費が成功するべきです")
	}
	if player.Mana != 2 {
		t.Errorf("ConsumeMana後のManaが期待値と異なります: got %d, want 2", player.Mana)
	}
}

func TestPlayerModel_ConsumeMana_不足時は失敗する(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 2

	ok := player.ConsumeMana(3)

	if ok {
		t.Error("マナ不足時は消費失敗するべきです")
	}
	// マナは変更されない
	if player.Mana != 2 {
		t.Errorf("消費失敗時のManaが変更されています: got %d, want 2", player.Mana)
	}
}

func TestPlayerModel_ConsumeMana_ゼロコストは常に成功(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 0

	ok := player.ConsumeMana(0)

	if !ok {
		t.Error("コスト0のマナ消費は常に成功するべきです")
	}
	if player.Mana != 0 {
		t.Errorf("コスト0消費後のManaが変更されています: got %d, want 0", player.Mana)
	}
}

func TestPlayerModel_ConsumeMana_ぴったり消費(t *testing.T) {
	player := NewPlayerWithMaxHP(1000)
	player.MaxMana = 10
	player.Mana = 3

	ok := player.ConsumeMana(3)

	if !ok {
		t.Error("ぴったりのマナ消費は成功するべきです")
	}
	if player.Mana != 0 {
		t.Errorf("ぴったり消費後のManaが期待値と異なります: got %d, want 0", player.Mana)
	}
}
