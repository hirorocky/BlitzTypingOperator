// Package domain はゲームのドメインモデルを定義します。
package domain

// PassiveEvent はパッシブスキル評価時のイベントを表す型です。
type PassiveEvent string

const (
	// PassiveEventNone はイベントなしを表します。
	PassiveEventNone PassiveEvent = "none"

	// PassiveEventBattleStart は戦闘開始イベントを表します。
	PassiveEventBattleStart PassiveEvent = "battle_start"

	// PassiveEventSkillUse はスキル使用イベントを表します。
	PassiveEventSkillUse PassiveEvent = "skill_use"

	// PassiveEventDamageReceived は被ダメージイベントを表します。
	PassiveEventDamageReceived PassiveEvent = "damage_received"

	// PassiveEventHeal は回復イベントを表します。
	PassiveEventHeal PassiveEvent = "heal"

	// PassiveEventBuffDebuffUse はバフ/デバフ使用イベントを表します。
	PassiveEventBuffDebuffUse PassiveEvent = "buff_debuff_use"

	// PassiveEventPhysicalAttack は物理攻撃イベントを表します。
	PassiveEventPhysicalAttack PassiveEvent = "physical_attack"

	// PassiveEventTypingMiss はタイピングミスイベントを表します。
	PassiveEventTypingMiss PassiveEvent = "typing_miss"

	// PassiveEventTimeout は時間切れイベントを表します。
	PassiveEventTimeout PassiveEvent = "timeout"

	// PassiveEventDebuffReceived はデバフ受けイベントを表します。
	PassiveEventDebuffReceived PassiveEvent = "debuff_received"
)

// PassiveEvaluationContext はパッシブスキル評価時のコンテキストを表す構造体です。
// バトル状態、プレイヤー状態、敵状態など、評価に必要な情報を含みます。
type PassiveEvaluationContext struct {
	// IsPerfect はパーフェクトタイピング（ミスなし完了）かどうかです。
	IsPerfect bool

	// WPM は現在のWPM（Words Per Minute）です。
	WPM float64

	// PlayerHPPercent はプレイヤーのHP割合（0-100）です。
	PlayerHPPercent float64

	// EnemyHPPercent は敵のHP割合（0-100）です。
	EnemyHPPercent float64

	// EnemyHasDebuff は敵がデバフ状態かどうかです。
	EnemyHasDebuff bool

	// Event は現在のイベントです。
	Event PassiveEvent

	// CurrentStacks は現在のスタック数です。
	CurrentStacks int

	// SameAttackCount は同種攻撃のカウントです。
	SameAttackCount int

	// UsesRemaining はバトル中の残り使用回数です。
	UsesRemaining int
}

// PassiveEvaluationResult はパッシブスキル評価結果を表す構造体です。
type PassiveEvaluationResult struct {
	// IsActive はパッシブスキルが有効かどうかです。
	IsActive bool

	// EffectMultiplier は効果倍率です（1.0=変化なし）。
	EffectMultiplier float64

	// EffectValue は効果値です（ダメージ、回復量など）。
	EffectValue float64

	// NeedsProbabilityCheck は確率チェックが必要かどうかです。
	NeedsProbabilityCheck bool

	// Probability は発動確率です（0.0〜1.0）。
	Probability float64
}

// EvaluatePassive はパッシブスキル定義とコンテキストから評価結果を計算します。
func EvaluatePassive(def PassiveSkill, ctx *PassiveEvaluationContext) PassiveEvaluationResult {
	result := PassiveEvaluationResult{
		IsActive:         false,
		EffectMultiplier: 1.0,
		EffectValue:      def.EffectValue,
	}

	switch def.TriggerType {
	case PassiveTriggerPermanent:
		result = evaluatePermanent(def, ctx)
	case PassiveTriggerConditional:
		result = evaluateConditional(def, ctx)
	case PassiveTriggerProbability:
		result = evaluateProbability(def, ctx)
	case PassiveTriggerStack:
		result = evaluateStack(def, ctx)
	case PassiveTriggerReactive:
		result = evaluateReactive(def, ctx)
	}

	return result
}

// evaluatePermanent は永続効果タイプを評価します。
func evaluatePermanent(def PassiveSkill, _ *PassiveEvaluationContext) PassiveEvaluationResult {
	return PassiveEvaluationResult{
		IsActive:         true,
		EffectMultiplier: def.EffectValue, // 永続効果の場合、EffectValueが倍率
		EffectValue:      def.EffectValue,
	}
}

// evaluateConditional は条件付き効果タイプを評価します。
func evaluateConditional(def PassiveSkill, ctx *PassiveEvaluationContext) PassiveEvaluationResult {
	result := PassiveEvaluationResult{
		IsActive:         false,
		EffectMultiplier: 1.0,
		EffectValue:      def.EffectValue,
	}

	if def.TriggerCondition == nil {
		return result
	}

	conditionMet := checkCondition(def.TriggerCondition, ctx)
	if conditionMet {
		result.IsActive = true
		result.EffectMultiplier = def.EffectValue
	}

	return result
}

// evaluateProbability は確率トリガータイプを評価します。
func evaluateProbability(def PassiveSkill, ctx *PassiveEvaluationContext) PassiveEvaluationResult {
	result := PassiveEvaluationResult{
		IsActive:              false,
		EffectMultiplier:      def.EffectValue,
		EffectValue:           def.EffectValue,
		NeedsProbabilityCheck: false,
		Probability:           def.Probability,
	}

	// 条件がある場合はまず条件をチェック
	if def.TriggerCondition != nil {
		conditionMet := checkCondition(def.TriggerCondition, ctx)
		if !conditionMet {
			return result
		}
	}

	// 条件を満たした場合は確率チェックが必要
	result.NeedsProbabilityCheck = true

	return result
}

// evaluateStack はスタック型を評価します。
func evaluateStack(def PassiveSkill, ctx *PassiveEvaluationContext) PassiveEvaluationResult {
	result := PassiveEvaluationResult{
		IsActive:         false,
		EffectMultiplier: 1.0,
		EffectValue:      def.EffectValue,
	}

	// スタック型の場合、条件タイプによって挙動が異なる
	if def.TriggerCondition != nil {
		switch def.TriggerCondition.Type {
		case TriggerConditionNoMissStreak:
			// ミスなし連続の場合、スタック数に応じて効果倍率を計算
			stacks := ctx.CurrentStacks
			if stacks > def.MaxStacks {
				stacks = def.MaxStacks
			}
			result.IsActive = stacks > 0
			result.EffectMultiplier = 1.0 + (float64(stacks) * def.StackIncrement)
		case TriggerConditionSameAttackCount:
			// 同種攻撃カウントの場合、閾値以上で発動
			threshold := int(def.TriggerCondition.Value)
			if ctx.SameAttackCount >= threshold {
				result.IsActive = true
				result.EffectMultiplier = 1.0 - def.EffectValue // ダメージ軽減
			}
		}
	}

	return result
}

// evaluateReactive は反応型を評価します。
func evaluateReactive(def PassiveSkill, ctx *PassiveEvaluationContext) PassiveEvaluationResult {
	result := PassiveEvaluationResult{
		IsActive:         false,
		EffectMultiplier: 1.0,
		EffectValue:      def.EffectValue,
	}

	if def.TriggerCondition == nil {
		return result
	}

	// 使用回数制限のチェック
	if def.UsesPerBattle > 0 && ctx.UsesRemaining <= 0 {
		return result
	}

	// イベントとトリガー条件の一致をチェック
	eventMatches := checkEventMatch(def.TriggerCondition.Type, ctx.Event)
	if eventMatches {
		result.IsActive = true
		result.EffectMultiplier = def.EffectValue
	}

	return result
}

// checkCondition はトリガー条件が満たされているかをチェックします。
func checkCondition(cond *TriggerCondition, ctx *PassiveEvaluationContext) bool {
	switch cond.Type {
	case TriggerConditionAccuracyEquals:
		return ctx.IsPerfect
	case TriggerConditionWPMAbove:
		return ctx.WPM >= cond.Value
	case TriggerConditionHPBelowPercent:
		return ctx.PlayerHPPercent <= cond.Value
	case TriggerConditionEnemyHPBelowPercent:
		return ctx.EnemyHPPercent <= cond.Value
	case TriggerConditionEnemyHasDebuff:
		return ctx.EnemyHasDebuff
	case TriggerConditionOnDamageReceived:
		return ctx.Event == PassiveEventDamageReceived
	case TriggerConditionOnHeal:
		return ctx.Event == PassiveEventHeal
	case TriggerConditionOnSkillUse:
		return ctx.Event == PassiveEventSkillUse
	case TriggerConditionOnBuffDebuffUse:
		return ctx.Event == PassiveEventBuffDebuffUse
	case TriggerConditionOnPhysicalAttack:
		return ctx.Event == PassiveEventPhysicalAttack
	case TriggerConditionOnTypingMiss:
		return ctx.Event == PassiveEventTypingMiss
	case TriggerConditionOnTimeout:
		return ctx.Event == PassiveEventTimeout
	case TriggerConditionOnDebuffReceived:
		return ctx.Event == PassiveEventDebuffReceived
	case TriggerConditionOnBattleStart:
		return ctx.Event == PassiveEventBattleStart
	default:
		return false
	}
}

// checkEventMatch はイベントとトリガー条件タイプが一致するかをチェックします。
func checkEventMatch(condType TriggerConditionType, event PassiveEvent) bool {
	switch condType {
	case TriggerConditionOnBattleStart:
		return event == PassiveEventBattleStart
	case TriggerConditionOnDamageReceived:
		return event == PassiveEventDamageReceived
	case TriggerConditionOnHeal:
		return event == PassiveEventHeal
	case TriggerConditionOnSkillUse:
		return event == PassiveEventSkillUse
	case TriggerConditionOnBuffDebuffUse:
		return event == PassiveEventBuffDebuffUse
	case TriggerConditionOnPhysicalAttack:
		return event == PassiveEventPhysicalAttack
	case TriggerConditionOnTypingMiss:
		return event == PassiveEventTypingMiss
	case TriggerConditionOnTimeout:
		return event == PassiveEventTimeout
	case TriggerConditionOnDebuffReceived:
		return event == PassiveEventDebuffReceived
	default:
		return false
	}
}
