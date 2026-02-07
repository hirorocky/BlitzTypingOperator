// Package domain はゲームのドメインモデルを定義します。
package domain

// EffectSourceType は効果のソース種別を表します。
// パッシブスキル、チェイン効果、バフ、デバフを区別します。
type EffectSourceType string

const (
	// SourcePassive はパッシブスキルからの効果を表します。
	SourcePassive EffectSourceType = "passive"

	// SourceChain はチェイン効果からの効果を表します。
	SourceChain EffectSourceType = "chain"

	// SourceBuff はバフからの効果を表します。
	SourceBuff EffectSourceType = "buff"

	// SourceDebuff はデバフからの効果を表します。
	SourceDebuff EffectSourceType = "debuff"
)

// EffectEntry は効果テーブルの1行（効果の1つのソース）を表します。
// パッシブスキル、チェイン効果、バフ、デバフなど全てのソースを統一的に扱います。
type EffectEntry struct {
	// ========== ソース管理列（逆引き用） ==========

	// SourceType はパッシブ/チェイン/バフ/デバフの種別です。
	SourceType EffectSourceType

	// SourceID は元定義のID（passive_skills.jsonのidなど）です。
	SourceID string

	// SourceIndex はエージェント番号などのインデックスです。
	SourceIndex int

	// ========== 表示・識別列 ==========

	// Name は効果の表示名です。
	Name string

	// ========== 有効条件列 ==========

	// EnableCondition は nil なら常に有効、非 nil なら Context で判定します。
	EnableCondition func(ctx *EffectContext) bool

	// ========== 時間管理列 ==========

	// Duration は残り時間（nil=永続）です。
	Duration *float64

	// ========== 効果値列 ==========

	// Values は数値型効果のマップです。
	Values map[EffectColumn]float64

	// Flags はbool型効果のマップです。
	Flags map[EffectColumn]bool

	// ========== 確率判定 ==========

	// Probability は発動確率（0.0=常に発動、0.5=50%など）です。
	Probability float64

	// ========== メタ情報 ==========

	// OneShot は一度発動したら消えるかを表します。
	OneShot bool

	// Triggered は既に発動済みか（OneShot用）を表します。
	Triggered bool

	// ========== スタック型パッシブ用 ==========

	// MaxStacks はスタック型パッシブの最大スタック数です（0=非スタック型）。
	MaxStacks int

	// StackIncrement はスタックごとの効果増分です。
	StackIncrement float64
}

// IsEnabled は現在のコンテキストで有効かを判定します。
func (e *EffectEntry) IsEnabled(ctx *EffectContext) bool {
	if e.EnableCondition == nil {
		return true
	}
	return e.EnableCondition(ctx)
}

// Clone はエントリのコピーを作成します。
func (e *EffectEntry) Clone() EffectEntry {
	clone := *e
	clone.Values = make(map[EffectColumn]float64)
	for k, v := range e.Values {
		clone.Values[k] = v
	}
	clone.Flags = make(map[EffectColumn]bool)
	for k, v := range e.Flags {
		clone.Flags[k] = v
	}
	// EnableCondition は参照コピー（関数なので深いコピーは不要）
	return clone
}

// IsPermanent は永続効果かどうかを判定します。
func (e *EffectEntry) IsPermanent() bool {
	return e.Duration == nil
}

// GetRemainingDuration は残り時間を取得します（永続の場合は-1を返す）。
func (e *EffectEntry) GetRemainingDuration() float64 {
	if e.Duration == nil {
		return -1
	}
	return *e.Duration
}

// EffectResult は列ごとの集計結果を表します。
// 全ての有効な効果を集計した最終的な効果値を保持します。
type EffectResult struct {
	// ========== 攻撃強化系 ==========

	// DamageBonus は加算ダメージです。
	DamageBonus int

	// DamageMultiplier はダメージ倍率です。
	DamageMultiplier float64

	// ArmorPierce は防御貫通が有効かを表します。
	ArmorPierce bool

	// LifeSteal はHP吸収率です。
	LifeSteal float64

	// ========== 防御強化系 ==========

	// DamageCut は被ダメ軽減率です。
	DamageCut float64

	// Evasion は回避率です。
	Evasion float64

	// Reflect は反射率です。
	Reflect float64

	// Regen は継続回復量です。
	Regen float64

	// ========== 回復強化系 ==========

	// HealBonus は加算回復です。
	HealBonus int

	// HealMultiplier は回復倍率です。
	HealMultiplier float64

	// Overheal は超過回復が有効かを表します。
	Overheal bool

	// ========== タイピング系 ==========

	// TimeExtend は時間延長（秒）です。
	TimeExtend float64

	// AutoCorrect はミス無視回数です。
	AutoCorrect int

	// TypingDifficulty はタイピング難易度の乗算修正値です（1.0が基準）。
	TypingDifficulty float64

	// ========== リキャスト系 ==========

	// CooldownReduce はCD短縮率です。
	CooldownReduce float64

	// ========== バフ/デバフ延長系 ==========

	// BuffExtend はバフ延長（秒）です。
	BuffExtend float64

	// DebuffExtend はデバフ延長（秒）です。
	DebuffExtend float64

	// ========== 特殊系 ==========

	// DoubleCast は2回発動確率です。
	DoubleCast float64

	// ========== ステータス系 ==========

	// STRBonus はSTR加算値です。
	STRBonus int

	// STRMultiplier はSTR倍率です。
	STRMultiplier float64

	// INTBonus はINT加算値です。
	INTBonus int

	// INTMultiplier はINT倍率です。
	INTMultiplier float64

	// WILBonus はWIL加算値です。
	WILBonus int

	// WILMultiplier はWIL倍率です。
	WILMultiplier float64

	// LUKBonus はLUK加算値です。
	LUKBonus int

	// LUKMultiplier はLUK倍率です。
	LUKMultiplier float64

	// CritRate はクリティカル率加算です。
	CritRate float64

	// ========== デバッグ用 ==========

	// ActiveSources は有効だったソース名のリストです。
	ActiveSources []string
}

// NewEffectResult は初期値を持つ EffectResult を生成します。
func NewEffectResult() EffectResult {
	return EffectResult{
		DamageMultiplier: 1.0,
		HealMultiplier:   1.0,
		TypingDifficulty: 1.0, // 乗算の初期値
		// ステータス系Multiplierは増加率として扱う（0.0 = +0%、0.25 = +25%）
		STRMultiplier: 0.0,
		INTMultiplier: 0.0,
		WILMultiplier: 0.0,
		LUKMultiplier: 0.0,
		ActiveSources: make([]string, 0),
	}
}

// CalculateFinalDamage は基礎ダメージに効果を適用します。
func (r *EffectResult) CalculateFinalDamage(baseDamage int) int {
	damage := float64(baseDamage)
	damage += float64(r.DamageBonus)
	damage *= r.DamageMultiplier
	if damage < 0 {
		damage = 0
	}
	return int(damage)
}

// CalculateFinalHeal は基礎回復量に効果を適用します。
func (r *EffectResult) CalculateFinalHeal(baseHeal int) int {
	heal := float64(baseHeal)
	heal += float64(r.HealBonus)
	heal *= r.HealMultiplier
	if heal < 0 {
		heal = 0
	}
	return int(heal)
}

// CalculateDamageReceived は受けるダメージを計算します。
func (r *EffectResult) CalculateDamageReceived(rawDamage int) int {
	damage := float64(rawDamage)
	damage *= (1.0 - r.DamageCut)
	if damage < 1 {
		damage = 1 // 最低1ダメージ
	}
	return int(damage)
}

// CalculateLifeStealHeal はライフスティールによる回復量を計算します。
func (r *EffectResult) CalculateLifeStealHeal(damage int) int {
	if r.LifeSteal <= 0 {
		return 0
	}
	return int(float64(damage) * r.LifeSteal)
}

// HasActiveEffects は有効な効果があるかを判定します。
func (r *EffectResult) HasActiveEffects() bool {
	return len(r.ActiveSources) > 0
}
