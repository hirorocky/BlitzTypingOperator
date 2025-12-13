// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestSourceType_定数の確認 はSourceType定数が正しく定義されていることを確認します。
func TestSourceType_定数の確認(t *testing.T) {
	tests := []struct {
		sourceType SourceType
		expected   string
	}{
		{SourceCore, "Core"},
		{SourceModule, "Module"},
		{SourceBuff, "Buff"},
		{SourceDebuff, "Debuff"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if string(tt.sourceType) != tt.expected {
				t.Errorf("SourceTypeが期待値と異なります: got %s, want %s", tt.sourceType, tt.expected)
			}
		})
	}
}

// TestStatModifiers_ゼロ値 はStatModifiersのゼロ値が適切であることを確認します。
func TestStatModifiers_ゼロ値(t *testing.T) {
	m := StatModifiers{}

	// 加算値のゼロ値は0
	if m.STR_Add != 0 || m.MAG_Add != 0 || m.SPD_Add != 0 || m.LUK_Add != 0 {
		t.Error("加算値のゼロ値が0ではありません")
	}

	// 乗算値のゼロ値は0（計算時に1.0として扱う必要がある）
	if m.STR_Mult != 0 || m.MAG_Mult != 0 || m.SPD_Mult != 0 || m.LUK_Mult != 0 {
		t.Error("乗算値のゼロ値が0ではありません")
	}

	// 特殊効果のゼロ値は0
	if m.CDReduction != 0 || m.TypingTimeExt != 0 || m.DamageReduction != 0 {
		t.Error("特殊効果のゼロ値が0ではありません")
	}
}

// TestStatModifiers_値設定 はStatModifiersに値を設定できることを確認します。
func TestStatModifiers_値設定(t *testing.T) {
	m := StatModifiers{
		STR_Add:         10,
		STR_Mult:        1.2,
		MAG_Add:         5,
		MAG_Mult:        1.1,
		SPD_Add:         3,
		SPD_Mult:        0.9,
		LUK_Add:         2,
		LUK_Mult:        1.0,
		CDReduction:     0.1,
		TypingTimeExt:   2.0,
		DamageReduction: 0.15,
		CritRate:        0.05,
		PhysicalEvade:   0.1,
		MagicEvade:      0.08,
	}

	if m.STR_Add != 10 {
		t.Errorf("STR_Addが期待値と異なります: got %d, want 10", m.STR_Add)
	}
	if m.STR_Mult != 1.2 {
		t.Errorf("STR_Multが期待値と異なります: got %f, want 1.2", m.STR_Mult)
	}
	if m.DamageReduction != 0.15 {
		t.Errorf("DamageReductionが期待値と異なります: got %f, want 0.15", m.DamageReduction)
	}
}

// TestEffectRow_フィールドの確認 はEffectRow構造体のフィールドが正しく設定されることを確認します。
func TestEffectRow_フィールドの確認(t *testing.T) {
	duration := 5.0
	row := EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Name:       "攻撃UP",
		Duration:   &duration,
		Modifiers: StatModifiers{
			STR_Add: 10,
		},
	}

	if row.ID != "buff_001" {
		t.Errorf("IDが期待値と異なります: got %s, want buff_001", row.ID)
	}
	if row.SourceType != SourceBuff {
		t.Errorf("SourceTypeが期待値と異なります: got %s, want Buff", row.SourceType)
	}
	if row.Name != "攻撃UP" {
		t.Errorf("Nameが期待値と異なります: got %s, want 攻撃UP", row.Name)
	}
	if *row.Duration != 5.0 {
		t.Errorf("Durationが期待値と異なります: got %f, want 5.0", *row.Duration)
	}
	if row.Modifiers.STR_Add != 10 {
		t.Errorf("Modifiers.STR_Addが期待値と異なります: got %d, want 10", row.Modifiers.STR_Add)
	}
}

// TestEffectRow_永続効果 はCore/Moduleの永続効果（Duration=nil）を確認します。
func TestEffectRow_永続効果(t *testing.T) {
	row := EffectRow{
		ID:         "core_001",
		SourceType: SourceCore,
		Name:       "攻撃特化",
		Duration:   nil, // 永続効果
		Modifiers: StatModifiers{
			STR_Add: 15,
		},
	}

	if row.Duration != nil {
		t.Error("永続効果のDurationはnilであるべきです")
	}
}

// TestEffectTable_新規作成 はNewEffectTableで空のテーブルが作成されることを確認します。
func TestEffectTable_新規作成(t *testing.T) {
	table := NewEffectTable()

	if table == nil {
		t.Error("NewEffectTableがnilを返しました")
	}
	if len(table.Rows) != 0 {
		t.Errorf("新規テーブルの行数が0ではありません: got %d", len(table.Rows))
	}
}

// TestEffectTable_行の追加 はAddRowで行を追加できることを確認します。

func TestEffectTable_行の追加(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	row := EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Name:       "攻撃UP",
		Duration:   &duration,
		Modifiers: StatModifiers{
			STR_Add: 10,
		},
	}

	table.AddRow(row)

	if len(table.Rows) != 1 {
		t.Errorf("行の追加後の行数が1ではありません: got %d", len(table.Rows))
	}
	if table.Rows[0].ID != "buff_001" {
		t.Errorf("追加された行のIDが異なります: got %s, want buff_001", table.Rows[0].ID)
	}
}

// TestEffectTable_行の削除 はRemoveRowで行を削除できることを確認します。

func TestEffectTable_行の削除(t *testing.T) {
	table := NewEffectTable()

	duration1 := 5.0
	duration2 := 3.0
	table.AddRow(EffectRow{ID: "buff_001", SourceType: SourceBuff, Duration: &duration1})
	table.AddRow(EffectRow{ID: "buff_002", SourceType: SourceBuff, Duration: &duration2})

	if len(table.Rows) != 2 {
		t.Errorf("行の追加後の行数が2ではありません: got %d", len(table.Rows))
	}

	table.RemoveRow("buff_001")

	if len(table.Rows) != 1 {
		t.Errorf("行の削除後の行数が1ではありません: got %d", len(table.Rows))
	}
	if table.Rows[0].ID != "buff_002" {
		t.Errorf("残っている行のIDが異なります: got %s, want buff_002", table.Rows[0].ID)
	}
}

// TestEffectTable_時限更新 はUpdateDurationsで時限効果の残り時間を更新できることを確認します。
func TestEffectTable_時限更新(t *testing.T) {
	table := NewEffectTable()

	duration1 := 5.0
	duration2 := 2.0
	table.AddRow(EffectRow{ID: "buff_001", SourceType: SourceBuff, Duration: &duration1})
	table.AddRow(EffectRow{ID: "buff_002", SourceType: SourceBuff, Duration: &duration2})
	table.AddRow(EffectRow{ID: "core_001", SourceType: SourceCore, Duration: nil}) // 永続効果

	// 1秒経過
	table.UpdateDurations(1.0)

	// buff_001は4秒残り
	if *table.Rows[0].Duration != 4.0 {
		t.Errorf("buff_001のDurationが期待値と異なります: got %f, want 4.0", *table.Rows[0].Duration)
	}
	// buff_002は1秒残り
	if *table.Rows[1].Duration != 1.0 {
		t.Errorf("buff_002のDurationが期待値と異なります: got %f, want 1.0", *table.Rows[1].Duration)
	}
	// core_001は永続（変化なし）
	if table.Rows[2].Duration != nil {
		t.Error("永続効果のDurationが変化しています")
	}
}

// TestEffectTable_期限切れ削除 はUpdateDurationsで期限切れの効果が削除されることを確認します。
func TestEffectTable_期限切れ削除(t *testing.T) {
	table := NewEffectTable()

	duration1 := 5.0
	duration2 := 1.0
	table.AddRow(EffectRow{ID: "buff_001", SourceType: SourceBuff, Duration: &duration1})
	table.AddRow(EffectRow{ID: "buff_002", SourceType: SourceBuff, Duration: &duration2})

	// 2秒経過（buff_002は期限切れ）
	table.UpdateDurations(2.0)

	if len(table.Rows) != 1 {
		t.Errorf("期限切れ後の行数が1ではありません: got %d", len(table.Rows))
	}
	if table.Rows[0].ID != "buff_001" {
		t.Errorf("残っている行のIDが異なります: got %s, want buff_001", table.Rows[0].ID)
	}
}

// TestEffectTable_最終ステータス計算_加算のみ は加算のみの効果でステータスが正しく計算されることを確認します。
func TestEffectTable_最終ステータス計算_加算のみ(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{
		ID:         "core_001",
		SourceType: SourceCore,
		Duration:   nil,
		Modifiers:  StatModifiers{STR_Add: 10, MAG_Add: 5},
	})
	table.AddRow(EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Add: 5},
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// STR: 100 + 10 + 5 = 115
	if finalStats.STR != 115 {
		t.Errorf("STRが期待値と異なります: got %d, want 115", finalStats.STR)
	}
	// MAG: 100 + 5 = 105
	if finalStats.MAG != 105 {
		t.Errorf("MAGが期待値と異なります: got %d, want 105", finalStats.MAG)
	}
	// SPD, LUK: 変化なし
	if finalStats.SPD != 100 {
		t.Errorf("SPDが期待値と異なります: got %d, want 100", finalStats.SPD)
	}
	if finalStats.LUK != 100 {
		t.Errorf("LUKが期待値と異なります: got %d, want 100", finalStats.LUK)
	}
}

// TestEffectTable_最終ステータス計算_乗算のみ は乗算のみの効果でステータスが正しく計算されることを確認します。
func TestEffectTable_最終ステータス計算_乗算のみ(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Mult: 1.2}, // 20%増加
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// STR: 100 × 1.2 = 120
	if finalStats.STR != 120 {
		t.Errorf("STRが期待値と異なります: got %d, want 120", finalStats.STR)
	}
}

// TestEffectTable_最終ステータス計算_加算乗算順序 は加算→乗算の順序で計算されることを確認します。
func TestEffectTable_最終ステータス計算_加算乗算順序(t *testing.T) {
	table := NewEffectTable()

	duration1 := 5.0
	duration2 := 5.0
	table.AddRow(EffectRow{
		ID:         "core_001",
		SourceType: SourceCore,
		Duration:   nil,
		Modifiers:  StatModifiers{STR_Add: 18},
	})
	table.AddRow(EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Duration:   &duration1,
		Modifiers:  StatModifiers{STR_Mult: 1.2}, // 20%増加
	})
	table.AddRow(EffectRow{
		ID:         "buff_002",
		SourceType: SourceBuff,
		Duration:   &duration2,
		Modifiers:  StatModifiers{STR_Add: 0, STR_Mult: 1.0}, // 効果なし
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// STR: (100 + 18) × 1.2 × 1.0 = 118 × 1.2 = 141（切り捨て）
	if finalStats.STR != 141 {
		t.Errorf("STRが期待値と異なります: got %d, want 141", finalStats.STR)
	}
}

// TestEffectTable_複数乗算の積 は複数の乗算効果が掛け合わされることを確認します。
func TestEffectTable_複数乗算の積(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Mult: 1.2}, // 20%増加
	})
	table.AddRow(EffectRow{
		ID:         "buff_002",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Mult: 1.5}, // 50%増加
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// STR: 100 × 1.2 × 1.5 = 179 (浮動小数点計算で1.7999...となり切り捨てで179)
	if finalStats.STR != 179 {
		t.Errorf("STRが期待値と異なります: got %d, want 179", finalStats.STR)
	}
}

// TestEffectTable_特殊効果の集計 は特殊効果が正しく集計されることを確認します。
func TestEffectTable_特殊効果の集計(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{
		ID:         "buff_001",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{CDReduction: 0.1, DamageReduction: 0.1},
	})
	table.AddRow(EffectRow{
		ID:         "buff_002",
		SourceType: SourceBuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{CDReduction: 0.05, CritRate: 0.03},
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// CDReduction: 0.1 + 0.05 = 0.15 (浮動小数点の比較は許容誤差を使用)
	epsilon := 0.0001
	if abs(finalStats.CDReduction-0.15) > epsilon {
		t.Errorf("CDReductionが期待値と異なります: got %f, want 0.15", finalStats.CDReduction)
	}
	// DamageReduction: 0.1
	if abs(finalStats.DamageReduction-0.1) > epsilon {
		t.Errorf("DamageReductionが期待値と異なります: got %f, want 0.1", finalStats.DamageReduction)
	}
	// CritRate: 0.03
	if abs(finalStats.CritRate-0.03) > epsilon {
		t.Errorf("CritRateが期待値と異なります: got %f, want 0.03", finalStats.CritRate)
	}
}

// abs は浮動小数点の絶対値を返すヘルパー関数です。
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// TestEffectTable_空テーブル計算 は空のテーブルでも正しく計算されることを確認します。
func TestEffectTable_空テーブル計算(t *testing.T) {
	table := NewEffectTable()

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// 効果なしなのでベースステータスと同じ
	if finalStats.STR != 100 {
		t.Errorf("STRが期待値と異なります: got %d, want 100", finalStats.STR)
	}
	if finalStats.MAG != 100 {
		t.Errorf("MAGが期待値と異なります: got %d, want 100", finalStats.MAG)
	}
}

// TestEffectTable_ソース種別でフィルタ はGetRowsBySourceでソース種別でフィルタできることを確認します。
func TestEffectTable_ソース種別でフィルタ(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{ID: "core_001", SourceType: SourceCore, Duration: nil})
	table.AddRow(EffectRow{ID: "buff_001", SourceType: SourceBuff, Duration: &duration})
	table.AddRow(EffectRow{ID: "buff_002", SourceType: SourceBuff, Duration: &duration})
	table.AddRow(EffectRow{ID: "debuff_001", SourceType: SourceDebuff, Duration: &duration})

	buffs := table.GetRowsBySource(SourceBuff)
	if len(buffs) != 2 {
		t.Errorf("Buffの行数が期待値と異なります: got %d, want 2", len(buffs))
	}

	cores := table.GetRowsBySource(SourceCore)
	if len(cores) != 1 {
		t.Errorf("Coreの行数が期待値と異なります: got %d, want 1", len(cores))
	}
}

// TestEffectTable_デバフの効果 はデバフがステータスを減少させることを確認します。

func TestEffectTable_デバフの効果(t *testing.T) {
	table := NewEffectTable()

	duration := 5.0
	table.AddRow(EffectRow{
		ID:         "debuff_001",
		SourceType: SourceDebuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Add: -10}, // STR減少
	})
	table.AddRow(EffectRow{
		ID:         "debuff_002",
		SourceType: SourceDebuff,
		Duration:   &duration,
		Modifiers:  StatModifiers{STR_Mult: 0.8}, // 20%減少
	})

	baseStats := Stats{STR: 100, MAG: 100, SPD: 100, LUK: 100}
	finalStats := table.Calculate(baseStats)

	// STR: (100 - 10) × 0.8 = 90 × 0.8 = 72
	if finalStats.STR != 72 {
		t.Errorf("STRが期待値と異なります: got %d, want 72", finalStats.STR)
	}
}

// TestFinalStats_フィールドの確認 はFinalStats構造体のフィールドが正しく設定されることを確認します。
func TestFinalStats_フィールドの確認(t *testing.T) {
	fs := FinalStats{
		STR:             100,
		MAG:             80,
		SPD:             70,
		LUK:             60,
		CDReduction:     0.1,
		TypingTimeExt:   2.0,
		DamageReduction: 0.15,
		CritRate:        0.05,
		PhysicalEvade:   0.1,
		MagicEvade:      0.08,
	}

	if fs.STR != 100 {
		t.Errorf("STRが期待値と異なります: got %d, want 100", fs.STR)
	}
	if fs.CDReduction != 0.1 {
		t.Errorf("CDReductionが期待値と異なります: got %f, want 0.1", fs.CDReduction)
	}
}
