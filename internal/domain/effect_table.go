// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"math/rand"
	"time"

	"hirorocky/type-battle/internal/config"
)

// EffectTable は効果の表全体を管理します。
// 効果列モデルに基づき、パッシブスキル、チェイン効果、バフ、デバフを
// 統一的に管理し、コンテキストに応じた効果集計を行います。
type EffectTable struct {
	// Entries は全ての効果エントリです。
	Entries []EffectEntry

	// rng は確率判定用の乱数生成器です。
	rng *rand.Rand
}

// NewEffectTable は新しい EffectTable を生成します。
func NewEffectTable() *EffectTable {
	return &EffectTable{
		Entries: make([]EffectEntry, 0),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewEffectTableWithSeed は指定されたシードで EffectTable を生成します（テスト用）。
func NewEffectTableWithSeed(seed int64) *EffectTable {
	return &EffectTable{
		Entries: make([]EffectEntry, 0),
		rng:     rand.New(rand.NewSource(seed)),
	}
}

// ========== エントリ追加メソッド ==========

// AddEntry はエントリを追加します。
func (t *EffectTable) AddEntry(entry EffectEntry) {
	t.Entries = append(t.Entries, entry)
}

// AddBuff はバフをエントリとして追加します。
// 同一IDかつ同一SourceTypeのエントリが存在する場合、Duration加算（上限99.9秒）を行います。
func (t *EffectTable) AddBuff(id, name string, duration float64, values map[EffectColumn]float64) {
	t.addTimedEffect(SourceBuff, id, name, duration, values)
}

// AddDebuff はデバフをエントリとして追加します。
// 同一IDかつ同一SourceTypeのエントリが存在する場合、Duration加算（上限99.9秒）を行います。
func (t *EffectTable) AddDebuff(id, name string, duration float64, values map[EffectColumn]float64) {
	t.addTimedEffect(SourceDebuff, id, name, duration, values)
}

// addTimedEffect は時限効果（バフ/デバフ）を追加または更新します。
func (t *EffectTable) addTimedEffect(st EffectSourceType, id, name string, duration float64, values map[EffectColumn]float64) {
	if existing := t.findBySourceIDAndType(id, st); existing != nil && existing.Duration != nil {
		newDuration := *existing.Duration + duration
		if newDuration > config.MaxStatusDuration {
			newDuration = config.MaxStatusDuration
		}
		*existing.Duration = newDuration
		return
	}
	d := duration
	t.AddEntry(EffectEntry{
		SourceType: st,
		SourceID:   id,
		Name:       name,
		Duration:   &d,
		Values:     values,
	})
}

// findBySourceIDAndType はSourceIDとSourceTypeで既存エントリを検索します。
func (t *EffectTable) findBySourceIDAndType(id string, st EffectSourceType) *EffectEntry {
	for i := range t.Entries {
		if t.Entries[i].SourceID == id && t.Entries[i].SourceType == st {
			return &t.Entries[i]
		}
	}
	return nil
}

// ========== 集計メソッド ==========

// Aggregate は Context を受けて有効な効果のみを集計します。
func (t *EffectTable) Aggregate(ctx *EffectContext) EffectResult {
	result := NewEffectResult()

	for i := range t.Entries {
		entry := &t.Entries[i]

		// 既に発動済みの OneShot はスキップ
		if entry.OneShot && entry.Triggered {
			continue
		}

		// 有効条件の判定
		if !entry.IsEnabled(ctx) {
			continue
		}

		// 確率判定
		if entry.Probability > 0 && t.rng.Float64() >= entry.Probability {
			continue
		}

		// OneShot なら発動済みフラグを立てる
		if entry.OneShot {
			entry.Triggered = true
		}

		// 有効なソースとして記録
		result.ActiveSources = append(result.ActiveSources, entry.Name)

		// 数値型効果の集計
		for col, val := range entry.Values {
			t.aggregateValue(&result, col, val)
		}

		// bool型効果の集計
		for col, flag := range entry.Flags {
			if flag {
				t.aggregateFlag(&result, col)
			}
		}
	}

	return result
}

// aggregateValue は数値型効果を集計します。
func (t *EffectTable) aggregateValue(result *EffectResult, col EffectColumn, val float64) {
	switch col {
	case ColDamageBonus:
		result.DamageBonus += int(val)
	case ColDamageMultiplier:
		result.DamageMultiplier *= val
	case ColLifeSteal:
		if val > result.LifeSteal {
			result.LifeSteal = val
		}
	case ColDamageCut:
		if val > result.DamageCut {
			result.DamageCut = val
		}
	case ColEvasion:
		if val > result.Evasion {
			result.Evasion = val
		}
	case ColReflect:
		if val > result.Reflect {
			result.Reflect = val
		}
	case ColRegen:
		result.Regen += val
	case ColHealBonus:
		result.HealBonus += int(val)
	case ColHealMultiplier:
		result.HealMultiplier *= val
	case ColTimeExtend:
		result.TimeExtend += val
	case ColAutoCorrect:
		result.AutoCorrect += int(val)
	case ColTypingDifficulty:
		result.TypingDifficulty *= val
	case ColCooldownReduce:
		// 加算集計（正=短縮、負=延長）
		result.CooldownReduce += val
	case ColBuffExtend:
		result.BuffExtend += val
	case ColDebuffExtend:
		result.DebuffExtend += val
	case ColDoubleCast:
		if val > result.DoubleCast {
			result.DoubleCast = val
		}
	case ColSTRBonus:
		result.STRBonus += int(val)
	case ColSTRMultiplier:
		result.STRMultiplier += val // 増加率として加算（0.25 = +25%）
	case ColINTBonus:
		result.INTBonus += int(val)
	case ColINTMultiplier:
		result.INTMultiplier += val // 増加率として加算
	case ColWILBonus:
		result.WILBonus += int(val)
	case ColWILMultiplier:
		result.WILMultiplier += val // 増加率として加算
	case ColLUKBonus:
		result.LUKBonus += int(val)
	case ColLUKMultiplier:
		result.LUKMultiplier += val // 増加率として加算
	case ColCritRate:
		result.CritRate += val
	}
}

// aggregateFlag はbool型効果を集計します。
func (t *EffectTable) aggregateFlag(result *EffectResult, col EffectColumn) {
	switch col {
	case ColArmorPierce:
		result.ArmorPierce = true
	case ColOverheal:
		result.Overheal = true
	}
}

// ========== 逆引きメソッド ==========

// FindBySourceType はソース種別でエントリを検索します。
func (t *EffectTable) FindBySourceType(st EffectSourceType) []EffectEntry {
	var results []EffectEntry
	for _, e := range t.Entries {
		if e.SourceType == st {
			results = append(results, e)
		}
	}
	return results
}

// FindBySourceID はソースIDでエントリを検索します。
func (t *EffectTable) FindBySourceID(id string) *EffectEntry {
	for i := range t.Entries {
		if t.Entries[i].SourceID == id {
			return &t.Entries[i]
		}
	}
	return nil
}

// FindByAgentIndex はエージェント番号でエントリを検索します。
func (t *EffectTable) FindByAgentIndex(idx int) []EffectEntry {
	var results []EffectEntry
	for _, e := range t.Entries {
		if e.SourceIndex == idx {
			results = append(results, e)
		}
	}
	return results
}

// ========== 時間経過処理 ==========

// Tick は時間経過を処理（期限切れエントリを削除）します。
func (t *EffectTable) Tick(deltaSeconds float64) {
	remaining := make([]EffectEntry, 0, len(t.Entries))
	for i := range t.Entries {
		entry := &t.Entries[i]
		if entry.Duration != nil {
			*entry.Duration -= deltaSeconds
			if *entry.Duration <= 0 {
				continue // 期限切れ
			}
		}
		remaining = append(remaining, *entry)
	}
	t.Entries = remaining
}

// UpdateDurations は Tick のエイリアスです（既存コード互換性用）。
func (t *EffectTable) UpdateDurations(deltaSeconds float64) {
	t.Tick(deltaSeconds)
}

// ========== 削除メソッド ==========

// RemoveBySourceID は指定IDのエントリを削除します。
func (t *EffectTable) RemoveBySourceID(id string) bool {
	for i, e := range t.Entries {
		if e.SourceID == id {
			t.Entries = append(t.Entries[:i], t.Entries[i+1:]...)
			return true
		}
	}
	return false
}

// RemoveBySourceType は指定種別のエントリを全削除します。
func (t *EffectTable) RemoveBySourceType(st EffectSourceType) int {
	count := 0
	remaining := make([]EffectEntry, 0, len(t.Entries))
	for _, e := range t.Entries {
		if e.SourceType != st {
			remaining = append(remaining, e)
		} else {
			count++
		}
	}
	t.Entries = remaining
	return count
}

// Clear は全てのエントリをクリアします。
func (t *EffectTable) Clear() {
	t.Entries = make([]EffectEntry, 0)
}

// ========== バフ/デバフ延長 ==========

// ExtendBuffs は全てのバフの持続時間を延長します。
func (t *EffectTable) ExtendBuffs(seconds float64) {
	if seconds <= 0 {
		return
	}
	for i := range t.Entries {
		if t.Entries[i].SourceType == SourceBuff && t.Entries[i].Duration != nil {
			*t.Entries[i].Duration += seconds
		}
	}
}

// ExtendDebuffs は全てのデバフの持続時間を延長します。
func (t *EffectTable) ExtendDebuffs(seconds float64) {
	if seconds <= 0 {
		return
	}
	for i := range t.Entries {
		if t.Entries[i].SourceType == SourceDebuff && t.Entries[i].Duration != nil {
			*t.Entries[i].Duration += seconds
		}
	}
}

// ExtendBuffDurations は ExtendBuffs のエイリアスです（既存コード互換性用）。
func (t *EffectTable) ExtendBuffDurations(seconds float64) {
	t.ExtendBuffs(seconds)
}

// ExtendDebuffDurations は ExtendDebuffs のエイリアスです（既存コード互換性用）。
func (t *EffectTable) ExtendDebuffDurations(seconds float64) {
	t.ExtendDebuffs(seconds)
}

// ========== ユーティリティ ==========

// GetActiveBuffs はアクティブなバフのリストを取得します。
func (t *EffectTable) GetActiveBuffs() []EffectEntry {
	return t.FindBySourceType(SourceBuff)
}

// GetActiveDebuffs はアクティブなデバフのリストを取得します。
func (t *EffectTable) GetActiveDebuffs() []EffectEntry {
	return t.FindBySourceType(SourceDebuff)
}

// GetPassiveSkills はパッシブスキルのリストを取得します。
func (t *EffectTable) GetPassiveSkills() []EffectEntry {
	return t.FindBySourceType(SourcePassive)
}

// GetChainEffects はチェイン効果のリストを取得します。
func (t *EffectTable) GetChainEffects() []EffectEntry {
	return t.FindBySourceType(SourceChain)
}

// Count は登録されているエントリ数を返します。
func (t *EffectTable) Count() int {
	return len(t.Entries)
}

// HasDebuffs はデバフが存在するかを判定します。
func (t *EffectTable) HasDebuffs() bool {
	for _, e := range t.Entries {
		if e.SourceType == SourceDebuff {
			return true
		}
	}
	return false
}

// ResetOneShots は全てのOneShotフラグをリセットします（バトル再開時など）。
func (t *EffectTable) ResetOneShots() {
	for i := range t.Entries {
		t.Entries[i].Triggered = false
	}
}
