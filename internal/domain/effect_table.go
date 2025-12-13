// Package domain はゲームのドメインモデルを定義します。
package domain

// SourceType は効果のソース種別を表す型です。
// コア特性、モジュールパッシブ、バフ、デバフを区別します。
type SourceType string

const (
	// SourceCore はコア特性からの効果を表します（装備中常時有効、永続）
	SourceCore SourceType = "Core"

	// SourceModule はモジュールパッシブからの効果を表します（装備中常時有効、永続）
	SourceModule SourceType = "Module"

	// SourceBuff はバフ（プレイヤー有利効果）を表します（時限）
	SourceBuff SourceType = "Buff"

	// SourceDebuff はデバフ（プレイヤー不利効果）を表します（時限）
	SourceDebuff SourceType = "Debuff"
)

// StatModifiers はステータス修正値を表す構造体です。
// 加算値と乗算値、および特殊効果を持ちます。
type StatModifiers struct {
	// 基本ステータス（加算）
	STR_Add int
	MAG_Add int
	SPD_Add int
	LUK_Add int

	// 基本ステータス（乗算）
	// 1.0 = 変化なし、1.2 = 20%増加、0.8 = 20%減少
	// 0.0の場合は1.0として扱う（ゼロ値対応）
	STR_Mult float64
	MAG_Mult float64
	SPD_Mult float64
	LUK_Mult float64

	// 特殊効果
	CDReduction     float64 // クールダウン短縮率（0.1 = 10%短縮）
	TypingTimeExt   float64 // タイピング時間延長（秒数）
	DamageReduction float64 // ダメージ軽減率（0.1 = 10%軽減）
	CritRate        float64 // クリティカル率加算（0.05 = 5%）
	PhysicalEvade   float64 // 物理回避率加算（0.1 = 10%）
	MagicEvade      float64 // 魔法回避率加算（0.1 = 10%）
}

// EffectRow は効果テーブルの1行を表す構造体です。

type EffectRow struct {
	// ID は効果の一意識別子です（例: "core_001", "buff_a3f2"）
	ID string

	// SourceType は効果のソース種別です（Core, Module, Buff, Debuff）
	SourceType SourceType

	// Name は効果の表示名です
	Name string

	// Duration は残り秒数です（nil = 永続、Core/Moduleは永続）

	Duration *float64

	// Modifiers はステータス修正値です
	Modifiers StatModifiers
}

// EffectTable は効果テーブルを表す構造体です。
// コア特性、モジュールパッシブ、バフ、デバフの効果を表形式で管理します。
type EffectTable struct {
	Rows []EffectRow
}

// NewEffectTable は新しい空のEffectTableを作成します。
func NewEffectTable() *EffectTable {
	return &EffectTable{
		Rows: make([]EffectRow, 0),
	}
}

// AddRow は効果テーブルに行を追加します。
// バフ付与時、バトル開始時のコア/モジュールパッシブ登録に使用します。
func (t *EffectTable) AddRow(row EffectRow) {
	t.Rows = append(t.Rows, row)
}

// RemoveRow はIDを指定して行を削除します。
// 時限切れやバトル終了時に使用します。
func (t *EffectTable) RemoveRow(id string) {
	newRows := make([]EffectRow, 0, len(t.Rows))
	for _, row := range t.Rows {
		if row.ID != id {
			newRows = append(newRows, row)
		}
	}
	t.Rows = newRows
}

// UpdateDurations は時限効果の残り時間を更新します。
// 毎ティック呼び出され、残り時間が0以下になった行は自動削除されます。

func (t *EffectTable) UpdateDurations(deltaSeconds float64) {
	for i := range t.Rows {
		if t.Rows[i].Duration != nil {
			*t.Rows[i].Duration -= deltaSeconds
		}
	}
	// 期限切れの行を削除
	t.Rows = filterExpired(t.Rows)
}

// filterExpired は期限切れでない行のみを返すヘルパー関数です。
func filterExpired(rows []EffectRow) []EffectRow {
	newRows := make([]EffectRow, 0, len(rows))
	for _, row := range rows {
		// 永続効果（Duration == nil）または残り時間がある行を保持
		if row.Duration == nil || *row.Duration > 0 {
			newRows = append(newRows, row)
		}
	}
	return newRows
}

// Calculate は基礎ステータスに効果を適用して最終ステータスを計算します。
// 計算順序: 加算 → 乗算
func (t *EffectTable) Calculate(baseStats Stats) FinalStats {
	// 加算値の集計
	var addSTR, addMAG, addSPD, addLUK int
	// 乗算値の集計（初期値1.0）
	multSTR, multMAG, multSPD, multLUK := 1.0, 1.0, 1.0, 1.0
	// 特殊効果の集計
	var cdReduction, typingTimeExt, dmgReduction float64
	var critRate, physEvade, magEvade float64

	for _, row := range t.Rows {
		m := row.Modifiers

		// 加算値を集計
		addSTR += m.STR_Add
		addMAG += m.MAG_Add
		addSPD += m.SPD_Add
		addLUK += m.LUK_Add

		// 乗算値を集計（0.0は1.0として扱う）
		if m.STR_Mult != 0 {
			multSTR *= m.STR_Mult
		}
		if m.MAG_Mult != 0 {
			multMAG *= m.MAG_Mult
		}
		if m.SPD_Mult != 0 {
			multSPD *= m.SPD_Mult
		}
		if m.LUK_Mult != 0 {
			multLUK *= m.LUK_Mult
		}

		// 特殊効果を集計（加算）
		cdReduction += m.CDReduction
		typingTimeExt += m.TypingTimeExt
		dmgReduction += m.DamageReduction
		critRate += m.CritRate
		physEvade += m.PhysicalEvade
		magEvade += m.MagicEvade
	}

	// 最終ステータス計算: (基礎 + 加算) × 乗算
	return FinalStats{
		STR: int(float64(baseStats.STR+addSTR) * multSTR),
		MAG: int(float64(baseStats.MAG+addMAG) * multMAG),
		SPD: int(float64(baseStats.SPD+addSPD) * multSPD),
		LUK: int(float64(baseStats.LUK+addLUK) * multLUK),

		CDReduction:     cdReduction,
		TypingTimeExt:   typingTimeExt,
		DamageReduction: dmgReduction,
		CritRate:        critRate,
		PhysicalEvade:   physEvade,
		MagicEvade:      magEvade,
	}
}

// GetRowsBySource はソース種別で行をフィルタします。
func (t *EffectTable) GetRowsBySource(sourceType SourceType) []EffectRow {
	result := make([]EffectRow, 0)
	for _, row := range t.Rows {
		if row.SourceType == sourceType {
			result = append(result, row)
		}
	}
	return result
}

// FinalStats は効果適用後の最終ステータスを表す構造体です。
// 基本ステータス（整数）と特殊効果（浮動小数点）を含みます。
type FinalStats struct {
	// 基本ステータス
	STR int
	MAG int
	SPD int
	LUK int

	// 特殊効果
	CDReduction     float64 // クールダウン短縮率
	TypingTimeExt   float64 // タイピング時間延長
	DamageReduction float64 // ダメージ軽減率
	CritRate        float64 // クリティカル率
	PhysicalEvade   float64 // 物理回避率
	MagicEvade      float64 // 魔法回避率
}
