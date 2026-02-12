// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestEffectSourceType_定数の確認 はEffectSourceType定数が正しく定義されていることを確認します。
func TestEffectSourceType_定数の確認(t *testing.T) {
	tests := []struct {
		sourceType EffectSourceType
		expected   string
	}{
		{SourcePassive, "passive"},
		{SourceChain, "chain"},
		{SourceBuff, "buff"},
		{SourceDebuff, "debuff"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.sourceType) != tt.expected {
				t.Errorf("EffectSourceTypeが期待値と異なります: got %s, want %s", tt.sourceType, tt.expected)
			}
		})
	}
}

// TestEffectEntry_フィールドの確認 はEffectEntry構造体のフィールドが正しく設定されることを確認します。
func TestEffectEntry_フィールドの確認(t *testing.T) {
	duration := 5.0
	entry := EffectEntry{
		SourceType:  SourceBuff,
		SourceID:    "buff_001",
		SourceIndex: 0,
		Name:        "攻撃UP",
		Duration:    &duration,
		Values: map[EffectColumn]float64{
			ColDamageBonus: 10,
		},
	}

	if entry.SourceID != "buff_001" {
		t.Errorf("SourceIDが期待値と異なります: got %s, want buff_001", entry.SourceID)
	}
	if entry.SourceType != SourceBuff {
		t.Errorf("SourceTypeが期待値と異なります: got %s, want buff", entry.SourceType)
	}
	if entry.Name != "攻撃UP" {
		t.Errorf("Nameが期待値と異なります: got %s, want 攻撃UP", entry.Name)
	}
	if *entry.Duration != 5.0 {
		t.Errorf("Durationが期待値と異なります: got %f, want 5.0", *entry.Duration)
	}
	if entry.Values[ColDamageBonus] != 10 {
		t.Errorf("Values[ColDamageBonus]が期待値と異なります: got %f, want 10", entry.Values[ColDamageBonus])
	}
}

// TestEffectEntry_永続効果 はパッシブスキルの永続効果（Duration=nil）を確認します。
func TestEffectEntry_永続効果(t *testing.T) {
	entry := EffectEntry{
		SourceType:  SourcePassive,
		SourceID:    "passive_001",
		SourceIndex: 0,
		Name:        "攻撃特化",
		Duration:    nil, // 永続効果
		Values: map[EffectColumn]float64{
			ColDamageBonus: 15,
		},
	}

	if entry.Duration != nil {
		t.Error("永続効果のDurationはnilであるべきです")
	}
	if !entry.IsPermanent() {
		t.Error("IsPermanent()がtrueを返すべきです")
	}
}

// TestEffectTable_新規作成 はNewEffectTableで空のテーブルが作成されることを確認します。
func TestEffectTable_新規作成(t *testing.T) {
	table := NewEffectTable()

	if table == nil {
		t.Error("NewEffectTableがnilを返しました")
	}
	if len(table.Entries) != 0 {
		t.Errorf("新規テーブルのエントリ数が0ではありません: got %d", len(table.Entries))
	}
}

// TestEffectTable_エントリの追加 はAddEntryでエントリを追加できることを確認します。
func TestEffectTable_エントリの追加(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	entry := EffectEntry{
		SourceType:  SourceBuff,
		SourceID:    "buff_001",
		SourceIndex: 0,
		Name:        "攻撃UP",
		Duration:    &duration,
		Values: map[EffectColumn]float64{
			ColDamageBonus: 10,
		},
	}

	table.AddEntry(entry)

	if len(table.Entries) != 1 {
		t.Errorf("エントリ追加後のエントリ数が1ではありません: got %d", len(table.Entries))
	}
	if table.Entries[0].SourceID != "buff_001" {
		t.Errorf("追加されたエントリのSourceIDが異なります: got %s, want buff_001", table.Entries[0].SourceID)
	}
}

// TestEffectTable_AddBuff はAddBuffでバフを追加できることを確認します。
func TestEffectTable_AddBuff(t *testing.T) {
	table := NewEffectTable()

	values := map[EffectColumn]float64{
		ColDamageMultiplier: 1.2,
	}
	table.AddBuff("攻撃力UP", "攻撃力UP", 10.0, values)

	if len(table.Entries) != 1 {
		t.Errorf("バフ追加後のエントリ数が1ではありません: got %d", len(table.Entries))
	}
	if table.Entries[0].SourceType != SourceBuff {
		t.Errorf("SourceTypeがSourceBuffではありません: got %s", table.Entries[0].SourceType)
	}
	if table.Entries[0].Name != "攻撃力UP" {
		t.Errorf("Nameが期待値と異なります: got %s", table.Entries[0].Name)
	}
}

// TestEffectTable_AddDebuff はAddDebuffでデバフを追加できることを確認します。
func TestEffectTable_AddDebuff(t *testing.T) {
	table := NewEffectTable()

	values := map[EffectColumn]float64{
		ColTimeExtend: -2.0,
	}
	table.AddDebuff("タイピング時間短縮", "タイピング時間短縮", 8.0, values)

	if len(table.Entries) != 1 {
		t.Errorf("デバフ追加後のエントリ数が1ではありません: got %d", len(table.Entries))
	}
	if table.Entries[0].SourceType != SourceDebuff {
		t.Errorf("SourceTypeがSourceDebuffではありません: got %s", table.Entries[0].SourceType)
	}
}

// TestEffectTable_エントリの削除 はRemoveBySourceIDでエントリを削除できることを確認します。
func TestEffectTable_エントリの削除(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddBuff("buff_002", "buff_002", 3.0, map[EffectColumn]float64{ColDamageBonus: 5})

	if len(table.Entries) != 2 {
		t.Errorf("エントリ追加後のエントリ数が2ではありません: got %d", len(table.Entries))
	}

	table.RemoveBySourceID("buff_001")

	if len(table.Entries) != 1 {
		t.Errorf("エントリ削除後のエントリ数が1ではありません: got %d", len(table.Entries))
	}
	if table.Entries[0].SourceID != "buff_002" {
		t.Errorf("残っているエントリのSourceIDが異なります: got %s, want buff_002", table.Entries[0].SourceID)
	}
}

// TestEffectTable_時限更新 はTickで時限効果の残り時間を更新できることを確認します。
func TestEffectTable_時限更新(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddBuff("buff_002", "buff_002", 2.0, map[EffectColumn]float64{ColDamageBonus: 5})
	table.AddEntry(EffectEntry{
		SourceType: SourcePassive,
		SourceID:   "passive_001",
		Name:       "パッシブ",
		Duration:   nil, // 永続効果
	})

	// 1秒経過
	table.Tick(1.0)

	// buff_001は4秒残り
	if *table.Entries[0].Duration != 4.0 {
		t.Errorf("buff_001のDurationが期待値と異なります: got %f, want 4.0", *table.Entries[0].Duration)
	}
	// buff_002は1秒残り
	if *table.Entries[1].Duration != 1.0 {
		t.Errorf("buff_002のDurationが期待値と異なります: got %f, want 1.0", *table.Entries[1].Duration)
	}
	// passive_001は永続（変化なし）
	if table.Entries[2].Duration != nil {
		t.Error("永続効果のDurationが変化しています")
	}
}

// TestEffectTable_期限切れ削除 はTickで期限切れの効果が削除されることを確認します。
func TestEffectTable_期限切れ削除(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddBuff("buff_002", "buff_002", 1.0, map[EffectColumn]float64{ColDamageBonus: 5})

	// 2秒経過（buff_002は期限切れ）
	table.Tick(2.0)

	if len(table.Entries) != 1 {
		t.Errorf("期限切れ後のエントリ数が1ではありません: got %d", len(table.Entries))
	}
	if table.Entries[0].SourceID != "buff_001" {
		t.Errorf("残っているエントリのSourceIDが異なります: got %s, want buff_001", table.Entries[0].SourceID)
	}
}

// TestEffectTable_Aggregate_加算効果 は加算効果が正しく集計されることを確認します。
func TestEffectTable_Aggregate_加算効果(t *testing.T) {
	table := NewEffectTableWithSeed(42)

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddBuff("buff_002", "buff_002", 5.0, map[EffectColumn]float64{ColDamageBonus: 5})

	ctx := NewEffectContext(100, 100, 50, 100)
	result := table.Aggregate(ctx)

	// DamageBonus: 10 + 5 = 15
	if result.DamageBonus != 15 {
		t.Errorf("DamageBonusが期待値と異なります: got %d, want 15", result.DamageBonus)
	}
}

// TestEffectTable_Aggregate_乗算効果 は乗算効果が正しく集計されることを確認します。
func TestEffectTable_Aggregate_乗算効果(t *testing.T) {
	table := NewEffectTableWithSeed(42)

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageMultiplier: 1.2})
	table.AddBuff("buff_002", "buff_002", 5.0, map[EffectColumn]float64{ColDamageMultiplier: 1.5})

	ctx := NewEffectContext(100, 100, 50, 100)
	result := table.Aggregate(ctx)

	// DamageMultiplier: 1.0 * 1.2 * 1.5 = 1.8
	epsilon := 0.0001
	if abs(result.DamageMultiplier-1.8) > epsilon {
		t.Errorf("DamageMultiplierが期待値と異なります: got %f, want 1.8", result.DamageMultiplier)
	}
}

// TestEffectTable_Aggregate_最大値効果 は最大値効果が正しく集計されることを確認します。
func TestEffectTable_Aggregate_最大値効果(t *testing.T) {
	table := NewEffectTableWithSeed(42)

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageCut: 0.2})
	table.AddBuff("buff_002", "buff_002", 5.0, map[EffectColumn]float64{ColDamageCut: 0.3})

	ctx := NewEffectContext(100, 100, 50, 100)
	result := table.Aggregate(ctx)

	// DamageCut: max(0.2, 0.3) = 0.3
	epsilon := 0.0001
	if abs(result.DamageCut-0.3) > epsilon {
		t.Errorf("DamageCutが期待値と異なります: got %f, want 0.3", result.DamageCut)
	}
}

// TestEffectTable_Aggregate_複合効果 は複合効果が正しく集計されることを確認します。
func TestEffectTable_Aggregate_複合効果(t *testing.T) {
	table := NewEffectTableWithSeed(42)

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{
		ColDamageBonus:      10,
		ColDamageMultiplier: 1.2,
		ColDamageCut:        0.2,
	})
	table.AddBuff("buff_002", "buff_002", 5.0, map[EffectColumn]float64{
		ColDamageBonus:    5,
		ColCooldownReduce: 0.1,
	})

	ctx := NewEffectContext(100, 100, 50, 100)
	result := table.Aggregate(ctx)

	// DamageBonus: 10 + 5 = 15
	if result.DamageBonus != 15 {
		t.Errorf("DamageBonusが期待値と異なります: got %d, want 15", result.DamageBonus)
	}
	// DamageMultiplier: 1.0 * 1.2 = 1.2
	epsilon := 0.0001
	if abs(result.DamageMultiplier-1.2) > epsilon {
		t.Errorf("DamageMultiplierが期待値と異なります: got %f, want 1.2", result.DamageMultiplier)
	}
	// DamageCut: 0.2
	if abs(result.DamageCut-0.2) > epsilon {
		t.Errorf("DamageCutが期待値と異なります: got %f, want 0.2", result.DamageCut)
	}
	// CooldownReduce: 0.1
	if abs(result.CooldownReduce-0.1) > epsilon {
		t.Errorf("CooldownReduceが期待値と異なります: got %f, want 0.1", result.CooldownReduce)
	}
}

// TestEffectTable_空テーブル集計 は空のテーブルでも正しく集計されることを確認します。
func TestEffectTable_空テーブル集計(t *testing.T) {
	table := NewEffectTable()

	ctx := NewEffectContext(100, 100, 50, 100)
	result := table.Aggregate(ctx)

	// 効果なしなのでデフォルト値
	if result.DamageBonus != 0 {
		t.Errorf("DamageBonusが期待値と異なります: got %d, want 0", result.DamageBonus)
	}
	if result.DamageMultiplier != 1.0 {
		t.Errorf("DamageMultiplierが期待値と異なります: got %f, want 1.0", result.DamageMultiplier)
	}
}

// TestEffectTable_ソース種別でフィルタ はFindBySourceTypeでソース種別でフィルタできることを確認します。
func TestEffectTable_ソース種別でフィルタ(t *testing.T) {
	table := NewEffectTable()

	table.AddEntry(EffectEntry{SourceType: SourcePassive, SourceID: "passive_001", Duration: nil})
	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddBuff("buff_002", "buff_002", 5.0, map[EffectColumn]float64{ColDamageBonus: 5})
	table.AddDebuff("debuff_001", "debuff_001", 5.0, map[EffectColumn]float64{ColTimeExtend: -2})

	buffs := table.FindBySourceType(SourceBuff)
	if len(buffs) != 2 {
		t.Errorf("Buffのエントリ数が期待値と異なります: got %d, want 2", len(buffs))
	}

	passives := table.FindBySourceType(SourcePassive)
	if len(passives) != 1 {
		t.Errorf("Passiveのエントリ数が期待値と異なります: got %d, want 1", len(passives))
	}
}

// TestEffectResult_CalculateFinalDamage はダメージ計算が正しく行われることを確認します。
func TestEffectResult_CalculateFinalDamage(t *testing.T) {
	result := NewEffectResult()
	result.DamageBonus = 10
	result.DamageMultiplier = 1.2

	// baseDamage: 100, bonus: 10, multiplier: 1.2
	// (100 + 10) * 1.2 = 132
	finalDamage := result.CalculateFinalDamage(100)
	if finalDamage != 132 {
		t.Errorf("最終ダメージが期待値と異なります: got %d, want 132", finalDamage)
	}
}

// TestEffectResult_CalculateDamageReceived は被ダメージ計算が正しく行われることを確認します。
func TestEffectResult_CalculateDamageReceived(t *testing.T) {
	result := NewEffectResult()
	result.DamageCut = 0.3

	// rawDamage: 100, cut: 0.3
	// 100 * (1 - 0.3) = 70
	received := result.CalculateDamageReceived(100)
	if received != 70 {
		t.Errorf("被ダメージが期待値と異なります: got %d, want 70", received)
	}
}

// TestEffectTable_ExtendBuffs はバフの持続時間が延長されることを確認します。
func TestEffectTable_ExtendBuffs(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("buff_001", "buff_001", 5.0, map[EffectColumn]float64{ColDamageBonus: 10})
	table.AddDebuff("debuff_001", "debuff_001", 5.0, map[EffectColumn]float64{ColTimeExtend: -2})

	table.ExtendBuffs(3.0)

	// バフは8秒になる
	buffEntry := table.FindBySourceType(SourceBuff)[0]
	if *buffEntry.Duration != 8.0 {
		t.Errorf("バフのDurationが期待値と異なります: got %f, want 8.0", *buffEntry.Duration)
	}

	// デバフは変化なし
	debuffEntry := table.FindBySourceType(SourceDebuff)[0]
	if *debuffEntry.Duration != 5.0 {
		t.Errorf("デバフのDurationが変化しています: got %f, want 5.0", *debuffEntry.Duration)
	}
}

func TestEffectTable_Aggregate_TypingDifficulty(t *testing.T) {
	table := NewEffectTableWithSeed(42)

	// TypingDifficulty は AggMult（乗算集計）
	table.AddEntry(EffectEntry{
		SourceType: SourcePassive,
		SourceID:   "typing_diff_1",
		Name:       "難易度低下",
		Values: map[EffectColumn]float64{
			ColTypingDifficulty: 0.8, // 20%低下
		},
	})
	table.AddEntry(EffectEntry{
		SourceType: SourceBuff,
		SourceID:   "typing_diff_2",
		Name:       "難易度バフ",
		Duration:   func() *float64 { d := 10.0; return &d }(),
		Values: map[EffectColumn]float64{
			ColTypingDifficulty: 0.9, // 10%低下
		},
	})

	result := table.Aggregate(nil)

	// 初期値1.0 × 0.8 × 0.9 = 0.72
	expected := 0.72
	if abs(result.TypingDifficulty-expected) > 0.001 {
		t.Errorf("TypingDifficulty = %f, want %f", result.TypingDifficulty, expected)
	}
}

func TestNewEffectResult_TypingDifficulty初期値(t *testing.T) {
	result := NewEffectResult()
	if result.TypingDifficulty != 1.0 {
		t.Errorf("TypingDifficultyの初期値 = %f, want 1.0", result.TypingDifficulty)
	}
}

// abs は浮動小数点の絶対値を返すヘルパー関数です。
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
