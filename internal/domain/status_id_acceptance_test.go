// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"math"
	"testing"

	"hirorocky/type-battle/internal/config"
)

// ========== 受け入れ基準4: 同一SourceIDのバフ/デバフが重複追加されずDurationが加算される ==========

func TestAC4_AddBuff_同一ID_Duration加算(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColSTRMultiplier: 0.25}

	table.AddBuff("st_str_buff_lv1", "STR25%UP", 10.0, values)
	table.AddBuff("st_str_buff_lv1", "STR25%UP", 10.0, values)

	if len(table.Entries) != 1 {
		t.Fatalf("同一IDのバフが重複追加されました: got %d entries, want 1", len(table.Entries))
	}
	expected := 20.0
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Duration加算が不正: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

func TestAC4_AddDebuff_同一ID_Duration加算(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColDamageMultiplier: 0.9}

	table.AddDebuff("st_attack_debuff", "ダメージ10%ダウン", 10.0, values)
	table.AddDebuff("st_attack_debuff", "ダメージ10%ダウン", 8.0, values)

	if len(table.Entries) != 1 {
		t.Fatalf("同一IDのデバフが重複追加されました: got %d entries, want 1", len(table.Entries))
	}
	expected := 18.0
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Duration加算が不正: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

func TestAC4_AddBuff_同一ID_3回重複_Duration累積(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColDamageCut: 0.1}

	table.AddBuff("st_defense_buff", "被ダメ10%軽減", 10.0, values)
	table.AddBuff("st_defense_buff", "被ダメ10%軽減", 5.0, values)
	table.AddBuff("st_defense_buff", "被ダメ10%軽減", 3.0, values)

	if len(table.Entries) != 1 {
		t.Fatalf("同一IDのバフが重複追加されました: got %d entries, want 1", len(table.Entries))
	}
	expected := 18.0
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Duration加算が不正: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

// ========== 受け入れ基準5: Duration加算の上限は99.9秒 ==========

func TestAC5_AddBuff_Duration上限クランプ(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColSTRMultiplier: 0.25}

	table.AddBuff("st_str_buff", "STR UP", 50.0, values)
	table.AddBuff("st_str_buff", "STR UP", 60.0, values)

	expected := config.MaxStatusDuration
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Durationが上限にクランプされていません: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

func TestAC5_AddDebuff_Duration上限クランプ(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColDamageMultiplier: 0.9}

	table.AddDebuff("st_debuff", "デバフ", 99.0, values)
	table.AddDebuff("st_debuff", "デバフ", 5.0, values)

	expected := config.MaxStatusDuration
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Durationが上限にクランプされていません: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

func TestAC5_AddBuff_上限ぴったりの場合加算されない(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColSTRMultiplier: 0.25}

	table.AddBuff("st_buff", "バフ", config.MaxStatusDuration, values)
	table.AddBuff("st_buff", "バフ", 10.0, values)

	expected := config.MaxStatusDuration
	if math.Abs(*table.Entries[0].Duration-expected) > 0.001 {
		t.Errorf("Durationが上限を超えています: got %f, want %f", *table.Entries[0].Duration, expected)
	}
}

// ========== 受け入れ基準6: 異なるIDのバフ/デバフは別エントリとして共存する ==========

func TestAC6_AddBuff_異なるID_別エントリ共存(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("st_str_buff_lv1", "STR25%UP", 10.0, map[EffectColumn]float64{ColSTRMultiplier: 0.25})
	table.AddBuff("st_str_buff_lv2", "STR50%UP", 10.0, map[EffectColumn]float64{ColSTRMultiplier: 0.50})

	if len(table.Entries) != 2 {
		t.Errorf("異なるIDのバフが2エントリにならない: got %d", len(table.Entries))
	}
}

func TestAC6_AddDebuff_異なるID_別エントリ共存(t *testing.T) {
	table := NewEffectTable()

	table.AddDebuff("st_debuff_a", "デバフA", 10.0, map[EffectColumn]float64{ColDamageMultiplier: 0.9})
	table.AddDebuff("st_debuff_b", "デバフB", 10.0, map[EffectColumn]float64{ColDamageCut: -0.1})

	if len(table.Entries) != 2 {
		t.Errorf("異なるIDのデバフが2エントリにならない: got %d", len(table.Entries))
	}
}

// ========== 受け入れ基準7: パッシブスキル/チェイン効果は影響を受けない ==========

func TestAC7_パッシブはAddEntryで直接追加_重複チェックなし(t *testing.T) {
	table := NewEffectTable()

	table.AddEntry(EffectEntry{
		SourceType: SourcePassive,
		SourceID:   "ps_test",
		Name:       "テストパッシブ",
		Duration:   nil,
		Values:     map[EffectColumn]float64{ColDamageBonus: 10},
	})
	table.AddEntry(EffectEntry{
		SourceType: SourcePassive,
		SourceID:   "ps_test",
		Name:       "テストパッシブ",
		Duration:   nil,
		Values:     map[EffectColumn]float64{ColDamageBonus: 10},
	})

	if len(table.Entries) != 2 {
		t.Errorf("パッシブのAddEntryで重複チェックが誤って実行されています: got %d entries, want 2", len(table.Entries))
	}
}

func TestAC7_チェインはAddEntryで直接追加_重複チェックなし(t *testing.T) {
	table := NewEffectTable()

	table.AddEntry(EffectEntry{
		SourceType: SourceChain,
		SourceID:   "chain_test",
		Name:       "テストチェイン",
		Duration:   nil,
		Values:     map[EffectColumn]float64{ColDamageBonus: 5},
	})
	table.AddEntry(EffectEntry{
		SourceType: SourceChain,
		SourceID:   "chain_test",
		Name:       "テストチェイン",
		Duration:   nil,
		Values:     map[EffectColumn]float64{ColDamageBonus: 5},
	})

	if len(table.Entries) != 2 {
		t.Errorf("チェインのAddEntryで重複チェックが誤って実行されています: got %d entries, want 2", len(table.Entries))
	}
}

// ========== 追加テスト: 同一IDでもSourceTypeが異なれば別エントリ ==========

func TestAddBuff_同一ID_異なるSourceType_別エントリ(t *testing.T) {
	table := NewEffectTable()
	values := map[EffectColumn]float64{ColDamageCut: 0.1}

	table.AddBuff("st_damage_cut", "被ダメ軽減", 10.0, values)
	table.AddDebuff("st_damage_cut", "被ダメ増加", 10.0, values)

	if len(table.Entries) != 2 {
		t.Errorf("バフとデバフで同一IDが別エントリにならない: got %d", len(table.Entries))
	}
}

// ========== 追加テスト: 重複時の効果値保持 ==========

func TestAddBuff_同一ID_効果値は既存を保持(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("st_str_buff", "STR UP", 10.0, map[EffectColumn]float64{ColSTRMultiplier: 0.25})
	table.AddBuff("st_str_buff", "STR UP", 10.0, map[EffectColumn]float64{ColSTRMultiplier: 0.50})

	if math.Abs(table.Entries[0].Values[ColSTRMultiplier]-0.25) > 0.001 {
		t.Errorf("効果値が上書きされています: got %f, want 0.25", table.Entries[0].Values[ColSTRMultiplier])
	}
}

// ========== 追加テスト: AddBuff新シグネチャでIDとNameが分離されている ==========

func TestAddBuff_IDとNameの分離(t *testing.T) {
	table := NewEffectTable()

	table.AddBuff("st_str_buff_lv1", "STR25%UP", 10.0, map[EffectColumn]float64{ColSTRMultiplier: 0.25})

	entry := table.Entries[0]
	if entry.SourceID != "st_str_buff_lv1" {
		t.Errorf("SourceIDが不正: got %s, want st_str_buff_lv1", entry.SourceID)
	}
	if entry.Name != "STR25%UP" {
		t.Errorf("Nameが不正: got %s, want STR25%%UP", entry.Name)
	}
}
