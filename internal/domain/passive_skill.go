// Package domain はゲームのドメインモデルを定義します。
package domain

// PassiveTriggerType はパッシブスキルのトリガータイプを表す型です。
// パッシブスキルがいつ発動するかを定義します。
type PassiveTriggerType string

const (
	// PassiveTriggerPermanent は永続効果を表します。
	// 常時有効な効果（例: バフ効果時間+50%）
	PassiveTriggerPermanent PassiveTriggerType = "permanent"

	// PassiveTriggerConditional は条件付き効果を表します。
	// 特定の条件を満たしたときに効果が有効になる（例: 正確性100%でスキル効果1.5倍）
	PassiveTriggerConditional PassiveTriggerType = "conditional"

	// PassiveTriggerProbability は確率トリガーを表します。
	// 特定の条件で確率的に発動する（例: 被ダメージ時20%で次の攻撃2倍）
	PassiveTriggerProbability PassiveTriggerType = "probability"

	// PassiveTriggerStack はスタック型を表します。
	// 条件を満たすたびに効果が累積する（例: ミスなし連続でダメージ+10%、最大+50%）
	PassiveTriggerStack PassiveTriggerType = "stack"

	// PassiveTriggerReactive は反応型を表します。
	// 特定のイベントに反応して発動する（例: 戦闘開始時に最初のスキルが即発動）
	PassiveTriggerReactive PassiveTriggerType = "reactive"
)

// PassiveEffectType はパッシブスキルの効果タイプを表す型です。
// どのような効果を与えるかを定義します。
type PassiveEffectType string

const (
	// PassiveEffectModifier はステータス修正効果を表します。
	// 加算・乗算によるステータス変化
	PassiveEffectModifier PassiveEffectType = "modifier"

	// PassiveEffectMultiplier は効果倍率を表します。
	// スキル効果やダメージの倍率変更
	PassiveEffectMultiplier PassiveEffectType = "multiplier"

	// PassiveEffectSpecial は特殊効果を表します。
	// 上記に分類されない特殊な効果（例: スキル2回発動、被ダメージ固定化）
	PassiveEffectSpecial PassiveEffectType = "special"
)

// TriggerConditionType はトリガー条件の種別を表す型です。
type TriggerConditionType string

const (
	// TriggerConditionAccuracyEquals は正確性が特定値に一致する条件です。
	TriggerConditionAccuracyEquals TriggerConditionType = "accuracy_equals"

	// TriggerConditionWPMAbove はWPMが特定値以上の条件です。
	TriggerConditionWPMAbove TriggerConditionType = "wpm_above"

	// TriggerConditionHPBelowPercent はHPが特定割合以下の条件です。
	TriggerConditionHPBelowPercent TriggerConditionType = "hp_below_percent"

	// TriggerConditionEnemyHPBelowPercent は敵HPが特定割合以下の条件です。
	TriggerConditionEnemyHPBelowPercent TriggerConditionType = "enemy_hp_below_percent"

	// TriggerConditionEnemyHasDebuff は敵がデバフ状態の条件です。
	TriggerConditionEnemyHasDebuff TriggerConditionType = "enemy_has_debuff"

	// TriggerConditionOnSkillUse はスキル使用時の条件です。
	TriggerConditionOnSkillUse TriggerConditionType = "on_skill_use"

	// TriggerConditionOnDamageReceived は被ダメージ時の条件です。
	TriggerConditionOnDamageReceived TriggerConditionType = "on_damage_received"

	// TriggerConditionOnHeal は回復時の条件です。
	TriggerConditionOnHeal TriggerConditionType = "on_heal"

	// TriggerConditionOnBuffDebuffUse はバフ/デバフ使用時の条件です。
	TriggerConditionOnBuffDebuffUse TriggerConditionType = "on_buff_debuff_use"

	// TriggerConditionOnPhysicalAttack は物理攻撃時の条件です。
	TriggerConditionOnPhysicalAttack TriggerConditionType = "on_physical_attack"

	// TriggerConditionOnTypingMiss はタイピングミス時の条件です。
	TriggerConditionOnTypingMiss TriggerConditionType = "on_typing_miss"

	// TriggerConditionOnTimeout は時間切れ時の条件です。
	TriggerConditionOnTimeout TriggerConditionType = "on_timeout"

	// TriggerConditionOnDebuffReceived はデバフ受け時の条件です。
	TriggerConditionOnDebuffReceived TriggerConditionType = "on_debuff_received"

	// TriggerConditionOnBattleStart は戦闘開始時の条件です。
	TriggerConditionOnBattleStart TriggerConditionType = "on_battle_start"

	// TriggerConditionNoMissStreak はミスなし連続の条件です。
	TriggerConditionNoMissStreak TriggerConditionType = "no_miss_streak"

	// TriggerConditionSameAttackCount は同種攻撃カウントの条件です。
	TriggerConditionSameAttackCount TriggerConditionType = "same_attack_count"
)

// TriggerCondition はパッシブスキルの発動条件を表す構造体です。
type TriggerCondition struct {
	// Type は条件の種別です。
	Type TriggerConditionType

	// Value は条件の閾値です（例: 正確性100、HP25%など）。
	Value float64
}

// PassiveSkill はパッシブスキルを表す構造体です。
// コア特性に紐づくパッシブスキルの仕様を定義します。
type PassiveSkill struct {
	// ID はパッシブスキルの一意識別子です。
	ID string

	// Name はパッシブスキルの表示名です。
	Name string

	// Description はパッシブスキルの効果説明です。
	Description string

	// ShortDescription はパッシブスキルの短い効果説明です（最大16文字程度）。
	// UI上でコンパクトに表示する際に使用します。
	ShortDescription string

	// TriggerType はトリガータイプです。
	TriggerType PassiveTriggerType

	// TriggerCondition は発動条件です（永続効果の場合はnil）。
	TriggerCondition *TriggerCondition

	// EffectType は効果タイプです。
	EffectType PassiveEffectType

	// EffectValue は効果量です（倍率、加算値など）。
	EffectValue float64

	// Effects は効果値のマップです（EffectColumn → 値）。
	// StatModifiersの代替として使用します。
	Effects map[EffectColumn]float64

	// Probability は発動確率です（確率トリガーの場合のみ使用、0.0〜1.0）。
	Probability float64

	// MaxStacks はスタック型の最大スタック数です。
	MaxStacks int

	// StackIncrement はスタックごとの効果増分です。
	StackIncrement float64

	// UsesPerBattle はバトル中の使用回数制限です（0=無制限）。
	UsesPerBattle int
}

// IsPermanent は永続効果かどうかを判定します。
func (p PassiveSkill) IsPermanent() bool {
	return p.TriggerType == PassiveTriggerPermanent
}

// HasProbability は確率判定があるかどうかを判定します。
func (p PassiveSkill) HasProbability() bool {
	return p.Probability > 0
}

// IsStackable はスタック可能かどうかを判定します。
func (p PassiveSkill) IsStackable() bool {
	return p.TriggerType == PassiveTriggerStack && p.MaxStacks > 0
}

// ToEntry は PassiveSkill を EffectEntry に変換します。
// 効果テーブルに登録可能な形式に変換します。
func (p PassiveSkill) ToEntry() EffectEntry {
	entry := EffectEntry{
		SourceType:      SourcePassive,
		SourceID:        p.ID,
		Name:            p.Name,
		EnableCondition: p.buildCondition(),
		Values:          p.buildValues(),
		Flags:           p.buildFlags(),
		Probability:     p.Probability,
		MaxStacks:       p.MaxStacks,
		StackIncrement:  p.StackIncrement,
	}
	return entry
}

// buildCondition はトリガー条件を EnableCondition 関数に変換します。
func (p PassiveSkill) buildCondition() func(*EffectContext) bool {
	if p.TriggerCondition == nil {
		return nil // 常に有効
	}

	cond := p.TriggerCondition
	switch cond.Type {
	case TriggerConditionHPBelowPercent:
		threshold := cond.Value / 100.0 // パーセントを0-1に変換
		return func(ctx *EffectContext) bool {
			return ctx.PlayerHPPercent < threshold
		}

	case TriggerConditionAccuracyEquals:
		return func(ctx *EffectContext) bool {
			return ctx.Accuracy >= 1.0
		}

	case TriggerConditionWPMAbove:
		threshold := cond.Value
		return func(ctx *EffectContext) bool {
			return ctx.WPM >= threshold
		}

	case TriggerConditionEnemyHPBelowPercent:
		threshold := cond.Value / 100.0
		return func(ctx *EffectContext) bool {
			return ctx.EnemyHPPercent < threshold
		}

	case TriggerConditionEnemyHasDebuff:
		return func(ctx *EffectContext) bool {
			return ctx.EnemyHasDebuff
		}

	case TriggerConditionOnDamageReceived:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnDamageRecv
		}

	case TriggerConditionOnHeal:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnHeal
		}

	case TriggerConditionOnSkillUse:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnSkillUse
		}

	case TriggerConditionOnTypingMiss:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnTypingMiss
		}

	case TriggerConditionOnBattleStart:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnBattleStart
		}

	case TriggerConditionNoMissStreak:
		count := int(cond.Value)
		return func(ctx *EffectContext) bool {
			return ctx.ComboCount >= count
		}

	case TriggerConditionOnTimeout:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnTimeout
		}

	case TriggerConditionOnPhysicalAttack:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnSkillUse && ctx.IsPhysical
		}

	case TriggerConditionOnBuffDebuffUse:
		return func(ctx *EffectContext) bool {
			return ctx.EventType == EventOnSkillUse && ctx.HasBuffDebuffEffect
		}

	case TriggerConditionSameAttackCount:
		threshold := int(cond.Value)
		return func(ctx *EffectContext) bool {
			return ctx.SameAttackCount >= threshold
		}

	default:
		return nil // 常に有効
	}
}

// buildValues は効果を EffectColumn の map に変換します。
func (p PassiveSkill) buildValues() map[EffectColumn]float64 {
	// Effectsが設定されている場合はそれを使用
	if len(p.Effects) > 0 {
		return p.Effects
	}

	// EffectTypeとEffectValueから変換
	values := make(map[EffectColumn]float64)

	switch p.EffectType {
	case PassiveEffectMultiplier:
		// 効果倍率はダメージ倍率として扱う
		values[ColDamageMultiplier] = p.EffectValue
	case PassiveEffectModifier:
		// ステータス修正は固定ダメージボーナスとして扱う
		values[ColDamageBonus] = p.EffectValue
	case PassiveEffectSpecial:
		// 特殊効果は個別に判定
	}

	return values
}

// buildFlags は bool 型効果を map に変換します。
func (p PassiveSkill) buildFlags() map[EffectColumn]bool {
	// 現時点では bool 型効果なし
	return nil
}
