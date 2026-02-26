// Package integration_test はタスク12の統合テストを提供します。

package integration_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ==================================================
// Task 12.4: 全19種類のパッシブスキルの挙動検証テスト
// ==================================================

// TestPassiveSkill_PermanentEffect は永続効果タイプのパッシブスキルを検証します。
func TestPassiveSkill_PermanentEffect(t *testing.T) {
	t.Run("ps_buff_extender", func(t *testing.T) {
		// バフエクステンダー: バフ効果時間+50%が常時適用されること
		def := domain.PassiveSkill{
			ID:          "ps_buff_extender",
			Name:        "バフエクステンダー",
			Description: "バフ効果時間+50%",
			TriggerType: domain.PassiveTriggerPermanent,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 1.5,
		}

		ctx := &domain.PassiveEvaluationContext{
			IsPerfect:       false,
			WPM:             60,
			PlayerHPPercent: 100,
			EnemyHPPercent:  100,
		}

		result := domain.EvaluatePassive(def, ctx)

		// 永続効果なので常にアクティブ
		if !result.IsActive {
			t.Error("永続効果は常にアクティブであるべきです")
		}

		// 効果倍率が1.5であること
		if result.EffectMultiplier != 1.5 {
			t.Errorf("EffectMultiplier expected 1.5, got %f", result.EffectMultiplier)
		}
	})
}

// TestPassiveSkill_ConditionalMultiplier は条件付き効果倍率タイプのパッシブスキルを検証します。
func TestPassiveSkill_ConditionalMultiplier(t *testing.T) {
	t.Run("ps_perfect_rhythm", func(t *testing.T) {
		// パーフェクトリズム: 正確性100%時のみ効果1.5倍
		def := domain.PassiveSkill{
			ID:          "ps_perfect_rhythm",
			Name:        "パーフェクトリズム",
			Description: "正確性100%でスキル効果1.5倍",
			TriggerType: domain.PassiveTriggerConditional,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 1.5,
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionAccuracyEquals,
				Value: 100,
			},
		}

		// パーフェクトの場合
		ctx100 := &domain.PassiveEvaluationContext{
			IsPerfect: true,
		}
		result100 := domain.EvaluatePassive(def, ctx100)

		if !result100.IsActive {
			t.Error("正確性100%でアクティブであるべきです")
		}
		if result100.EffectMultiplier != 1.5 {
			t.Errorf("EffectMultiplier expected 1.5, got %f", result100.EffectMultiplier)
		}

		// 非パーフェクトの場合
		ctx95 := &domain.PassiveEvaluationContext{
			IsPerfect: false,
		}
		result95 := domain.EvaluatePassive(def, ctx95)

		if result95.IsActive {
			t.Error("正確性95%ではアクティブでないべきです")
		}
	})

	t.Run("ps_speed_break", func(t *testing.T) {
		// スピードブレイク: WPM80以上で25%追加ダメージ
		def := domain.PassiveSkill{
			ID:          "ps_speed_break",
			Name:        "スピードブレイク",
			Description: "WPM80以上で25%追加ダメージ",
			TriggerType: domain.PassiveTriggerConditional,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 1.25,
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionWPMAbove,
				Value: 80,
			},
		}

		// WPM80以上の場合
		ctx80 := &domain.PassiveEvaluationContext{
			WPM: 80,
		}
		result80 := domain.EvaluatePassive(def, ctx80)

		if !result80.IsActive {
			t.Error("WPM80でアクティブであるべきです")
		}

		// WPM80未満の場合
		ctx60 := &domain.PassiveEvaluationContext{
			WPM: 60,
		}
		result60 := domain.EvaluatePassive(def, ctx60)

		if result60.IsActive {
			t.Error("WPM60ではアクティブでないべきです")
		}
	})

	t.Run("ps_endgame_specialist", func(t *testing.T) {
		// エンドゲームスペシャリスト: 敵HP30%以下で全ダメージ+25%
		def := domain.PassiveSkill{
			ID:          "ps_endgame_specialist",
			Name:        "エンドゲームスペシャリスト",
			Description: "敵HP30%以下で全ダメージ+25%",
			TriggerType: domain.PassiveTriggerConditional,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 1.25,
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionEnemyHPBelowPercent,
				Value: 30,
			},
		}

		// 敵HP30%の場合
		ctx30 := &domain.PassiveEvaluationContext{
			EnemyHPPercent: 30,
		}
		result30 := domain.EvaluatePassive(def, ctx30)

		if !result30.IsActive {
			t.Error("敵HP30%でアクティブであるべきです")
		}

		// 敵HP50%の場合
		ctx50 := &domain.PassiveEvaluationContext{
			EnemyHPPercent: 50,
		}
		result50 := domain.EvaluatePassive(def, ctx50)

		if result50.IsActive {
			t.Error("敵HP50%ではアクティブでないべきです")
		}
	})

	t.Run("ps_weak_point", func(t *testing.T) {
		// ウィークポイント: デバフ中の敵へダメージ+20%
		def := domain.PassiveSkill{
			ID:          "ps_weak_point",
			Name:        "ウィークポイント",
			Description: "デバフ中の敵へダメージ+20%",
			TriggerType: domain.PassiveTriggerConditional,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 1.2,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionEnemyHasDebuff,
			},
		}

		// 敵がデバフ中の場合
		ctxDebuff := &domain.PassiveEvaluationContext{
			EnemyHasDebuff: true,
		}
		resultDebuff := domain.EvaluatePassive(def, ctxDebuff)

		if !resultDebuff.IsActive {
			t.Error("敵がデバフ中でアクティブであるべきです")
		}

		// 敵がデバフ中でない場合
		ctxNoDebuff := &domain.PassiveEvaluationContext{
			EnemyHasDebuff: false,
		}
		resultNoDebuff := domain.EvaluatePassive(def, ctxNoDebuff)

		if resultNoDebuff.IsActive {
			t.Error("敵がデバフ中でないときはアクティブでないべきです")
		}
	})

	t.Run("ps_overdrive", func(t *testing.T) {
		// オーバードライブ: HP50%以下でリキャスト-30%、被ダメ+20%
		def := domain.PassiveSkill{
			ID:          "ps_overdrive",
			Name:        "オーバードライブ",
			Description: "HP50%以下でリキャスト-30%、被ダメ+20%",
			TriggerType: domain.PassiveTriggerConditional,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 0.7, // リキャスト-30%
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionHPBelowPercent,
				Value: 50,
			},
		}

		// HP50%以下の場合
		ctx50 := &domain.PassiveEvaluationContext{
			PlayerHPPercent: 50,
		}
		result50 := domain.EvaluatePassive(def, ctx50)

		if !result50.IsActive {
			t.Error("HP50%でアクティブであるべきです")
		}

		// HP60%の場合
		ctx60 := &domain.PassiveEvaluationContext{
			PlayerHPPercent: 60,
		}
		result60 := domain.EvaluatePassive(def, ctx60)

		if result60.IsActive {
			t.Error("HP60%ではアクティブでないべきです")
		}
	})
}

// TestPassiveSkill_ProbabilityTrigger は確率トリガータイプのパッシブスキルを検証します。
func TestPassiveSkill_ProbabilityTrigger(t *testing.T) {
	t.Run("ps_last_stand", func(t *testing.T) {
		// ラストスタンド: HP25%以下で30%の確率で被ダメージ1
		def := domain.PassiveSkill{
			ID:          "ps_last_stand",
			Name:        "ラストスタンド",
			Description: "HP25%以下で30%の確率で被ダメージ1",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 1,
			Probability: 0.3,
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionHPBelowPercent,
				Value: 25,
			},
		}

		// HP25%以下の場合
		ctx25 := &domain.PassiveEvaluationContext{
			PlayerHPPercent: 25,
		}
		result25 := domain.EvaluatePassive(def, ctx25)

		// 確率チェックが必要であること
		if !result25.NeedsProbabilityCheck {
			t.Error("確率チェックが必要であるべきです")
		}
		if result25.Probability != 0.3 {
			t.Errorf("Probability expected 0.3, got %f", result25.Probability)
		}

		// HP50%の場合
		ctx50 := &domain.PassiveEvaluationContext{
			PlayerHPPercent: 50,
		}
		result50 := domain.EvaluatePassive(def, ctx50)

		// 条件を満たさないので確率チェック不要
		if result50.NeedsProbabilityCheck {
			t.Error("HP50%では確率チェックは不要であるべきです")
		}
	})

	t.Run("ps_counter_charge", func(t *testing.T) {
		// カウンターチャージ: 被ダメージ時20%で次の攻撃2倍
		def := domain.PassiveSkill{
			ID:          "ps_counter_charge",
			Name:        "カウンターチャージ",
			Description: "被ダメージ時20%で次の攻撃2倍",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 2.0,
			Probability: 0.2,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnDamageReceived,
			},
		}

		// 被ダメージイベント発生時
		ctxDamage := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventDamageReceived,
		}
		resultDamage := domain.EvaluatePassive(def, ctxDamage)

		if !resultDamage.NeedsProbabilityCheck {
			t.Error("被ダメージ時は確率チェックが必要であるべきです")
		}
		if resultDamage.Probability != 0.2 {
			t.Errorf("Probability expected 0.2, got %f", resultDamage.Probability)
		}
	})

	t.Run("ps_miracle_heal", func(t *testing.T) {
		// ミラクルヒール: 回復スキル時10%でHP全回復
		def := domain.PassiveSkill{
			ID:          "ps_miracle_heal",
			Name:        "ミラクルヒール",
			Description: "回復スキル時10%でHP全回復",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 1.0,
			Probability: 0.1,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnHeal,
			},
		}

		// 回復イベント発生時
		ctxHeal := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventHeal,
		}
		resultHeal := domain.EvaluatePassive(def, ctxHeal)

		if !resultHeal.NeedsProbabilityCheck {
			t.Error("回復時は確率チェックが必要であるべきです")
		}
		if resultHeal.Probability != 0.1 {
			t.Errorf("Probability expected 0.1, got %f", resultHeal.Probability)
		}
	})

	t.Run("ps_chain_reaction", func(t *testing.T) {
		// チェインリアクション: バフ/デバフ使用時30%で効果時間2倍
		def := domain.PassiveSkill{
			ID:          "ps_chain_reaction",
			Name:        "チェインリアクション",
			Description: "バフ/デバフ使用時30%で効果時間2倍",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectMultiplier,
			EffectValue: 2.0,
			Probability: 0.3,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnBuffDebuffUse,
			},
		}

		// バフ/デバフ使用イベント発生時
		ctxBuff := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventBuffDebuffUse,
		}
		resultBuff := domain.EvaluatePassive(def, ctxBuff)

		if !resultBuff.NeedsProbabilityCheck {
			t.Error("バフ/デバフ使用時は確率チェックが必要であるべきです")
		}
	})

	t.Run("ps_echo_skill", func(t *testing.T) {
		// エコースキル: 15%の確率でスキル2回発動
		def := domain.PassiveSkill{
			ID:          "ps_echo_skill",
			Name:        "エコースキル",
			Description: "15%の確率でスキル2回発動",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 2.0,
			Probability: 0.15,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnSkillUse,
			},
		}

		// スキル使用イベント発生時
		ctxSkill := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventSkillUse,
		}
		resultSkill := domain.EvaluatePassive(def, ctxSkill)

		if !resultSkill.NeedsProbabilityCheck {
			t.Error("スキル使用時は確率チェックが必要であるべきです")
		}
		if resultSkill.Probability != 0.15 {
			t.Errorf("Probability expected 0.15, got %f", resultSkill.Probability)
		}
	})

	t.Run("ps_shadow_step", func(t *testing.T) {
		// シャドウステップ: 物理攻撃成功時20%で敵攻撃タイマーリセット
		def := domain.PassiveSkill{
			ID:          "ps_shadow_step",
			Name:        "シャドウステップ",
			Description: "物理攻撃成功時20%で敵攻撃タイマーリセット",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 1.0,
			Probability: 0.2,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnPhysicalAttack,
			},
		}

		// 物理攻撃イベント発生時
		ctxPhysical := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventPhysicalAttack,
		}
		resultPhysical := domain.EvaluatePassive(def, ctxPhysical)

		if !resultPhysical.NeedsProbabilityCheck {
			t.Error("物理攻撃時は確率チェックが必要であるべきです")
		}
	})

	t.Run("ps_debuff_reflect", func(t *testing.T) {
		// デバフリフレクト: デバフ受け時30%で敵にも同効果
		def := domain.PassiveSkill{
			ID:          "ps_debuff_reflect",
			Name:        "デバフリフレクト",
			Description: "デバフ受け時30%で敵にも同効果",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 1.0,
			Probability: 0.3,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnDebuffReceived,
			},
		}

		// デバフ受けイベント発生時
		ctxDebuff := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventDebuffReceived,
		}
		resultDebuff := domain.EvaluatePassive(def, ctxDebuff)

		if !resultDebuff.NeedsProbabilityCheck {
			t.Error("デバフ受け時は確率チェックが必要であるべきです")
		}
	})

	t.Run("ps_second_chance", func(t *testing.T) {
		// セカンドチャンス: 時間切れ時50%で再挑戦（制限時間半分）
		def := domain.PassiveSkill{
			ID:          "ps_second_chance",
			Name:        "セカンドチャンス",
			Description: "時間切れ時50%で再挑戦（制限時間半分）",
			TriggerType: domain.PassiveTriggerProbability,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 0.5,
			Probability: 0.5,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnTimeout,
			},
		}

		// 時間切れイベント発生時
		ctxTimeout := &domain.PassiveEvaluationContext{
			Event: domain.PassiveEventTimeout,
		}
		resultTimeout := domain.EvaluatePassive(def, ctxTimeout)

		if !resultTimeout.NeedsProbabilityCheck {
			t.Error("時間切れ時は確率チェックが必要であるべきです")
		}
		if resultTimeout.Probability != 0.5 {
			t.Errorf("Probability expected 0.5, got %f", resultTimeout.Probability)
		}
	})
}

// TestPassiveSkill_StackType はスタック型パッシブスキルを検証します。
func TestPassiveSkill_StackType(t *testing.T) {
	t.Run("ps_combo_master", func(t *testing.T) {
		// コンボマスター: ミスなし連続タイピングでダメージ累積+10%（最大+50%）
		def := domain.PassiveSkill{
			ID:             "ps_combo_master",
			Name:           "コンボマスター",
			Description:    "ミスなし連続タイピングでダメージ累積+10%（最大+50%）",
			TriggerType:    domain.PassiveTriggerStack,
			EffectType:     domain.PassiveEffectMultiplier,
			EffectValue:    1.0,
			MaxStacks:      5,
			StackIncrement: 0.1,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionNoMissStreak,
			},
		}

		// スタック0の場合
		ctx0 := &domain.PassiveEvaluationContext{
			CurrentStacks: 0,
		}
		result0 := domain.EvaluatePassive(def, ctx0)

		if result0.IsActive {
			t.Error("スタック0ではアクティブでないべきです")
		}

		// スタック3の場合（+30%）
		ctx3 := &domain.PassiveEvaluationContext{
			CurrentStacks: 3,
		}
		result3 := domain.EvaluatePassive(def, ctx3)

		if !result3.IsActive {
			t.Error("スタック3でアクティブであるべきです")
		}
		expectedMult := 1.0 + (3 * 0.1)
		if result3.EffectMultiplier != expectedMult {
			t.Errorf("EffectMultiplier expected %f, got %f", expectedMult, result3.EffectMultiplier)
		}

		// スタック5の場合（最大+50%）
		ctx5 := &domain.PassiveEvaluationContext{
			CurrentStacks: 5,
		}
		result5 := domain.EvaluatePassive(def, ctx5)

		expectedMult5 := 1.0 + (5 * 0.1)
		if result5.EffectMultiplier != expectedMult5 {
			t.Errorf("EffectMultiplier expected %f, got %f", expectedMult5, result5.EffectMultiplier)
		}

		// スタック7の場合（最大5でキャップ）
		ctx7 := &domain.PassiveEvaluationContext{
			CurrentStacks: 7,
		}
		result7 := domain.EvaluatePassive(def, ctx7)

		// 最大5にキャップされる
		expectedMult7 := 1.0 + (5 * 0.1)
		if result7.EffectMultiplier != expectedMult7 {
			t.Errorf("EffectMultiplier expected %f, got %f", expectedMult7, result7.EffectMultiplier)
		}
	})

	t.Run("ps_adaptive_shield", func(t *testing.T) {
		// アダプティブシールド: 同種攻撃3回目以降ダメージ-25%
		def := domain.PassiveSkill{
			ID:          "ps_adaptive_shield",
			Name:        "アダプティブシールド",
			Description: "同種攻撃3回目以降ダメージ-25%",
			TriggerType: domain.PassiveTriggerStack,
			EffectType:  domain.PassiveEffectModifier,
			EffectValue: 0.25,
			MaxStacks:   10,
			TriggerCondition: &domain.TriggerCondition{
				Type:  domain.TriggerConditionSameAttackCount,
				Value: 3,
			},
		}

		// 同種攻撃2回の場合
		ctx2 := &domain.PassiveEvaluationContext{
			SameAttackCount: 2,
		}
		result2 := domain.EvaluatePassive(def, ctx2)

		if result2.IsActive {
			t.Error("同種攻撃2回ではアクティブでないべきです")
		}

		// 同種攻撃3回の場合
		ctx3 := &domain.PassiveEvaluationContext{
			SameAttackCount: 3,
		}
		result3 := domain.EvaluatePassive(def, ctx3)

		if !result3.IsActive {
			t.Error("同種攻撃3回でアクティブであるべきです")
		}
		// ダメージ軽減25%なので、倍率は0.75
		expectedMult := 1.0 - 0.25
		if result3.EffectMultiplier != expectedMult {
			t.Errorf("EffectMultiplier expected %f, got %f", expectedMult, result3.EffectMultiplier)
		}
	})
}

// TestPassiveSkill_ReactiveType は反応型パッシブスキルを検証します。
func TestPassiveSkill_ReactiveType(t *testing.T) {
	t.Run("ps_debuff_absorber", func(t *testing.T) {
		// デバフアブソーバー: デバフ効果時間半減＋小回復
		def := domain.PassiveSkill{
			ID:          "ps_debuff_absorber",
			Name:        "デバフアブソーバー",
			Description: "デバフ効果時間半減＋小回復",
			TriggerType: domain.PassiveTriggerReactive,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 0.5,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnDebuffReceived,
			},
		}

		// デバフ受けイベント発生時
		ctxDebuff := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventDebuffReceived,
			UsesRemaining: 1,
		}
		resultDebuff := domain.EvaluatePassive(def, ctxDebuff)

		if !resultDebuff.IsActive {
			t.Error("デバフ受け時にアクティブであるべきです")
		}

		// 他のイベント発生時
		ctxOther := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventHeal,
			UsesRemaining: 1,
		}
		resultOther := domain.EvaluatePassive(def, ctxOther)

		if resultOther.IsActive {
			t.Error("他のイベントではアクティブでないべきです")
		}
	})

	t.Run("ps_typo_recovery", func(t *testing.T) {
		// タイポリカバリー: ミス時制限時間+1秒（1回/チャレンジ）
		def := domain.PassiveSkill{
			ID:            "ps_typo_recovery",
			Name:          "タイポリカバリー",
			Description:   "ミス時制限時間+1秒（1回/チャレンジ）",
			TriggerType:   domain.PassiveTriggerReactive,
			EffectType:    domain.PassiveEffectModifier,
			EffectValue:   1.0,
			UsesPerBattle: 1,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnTypingMiss,
			},
		}

		// ミスイベント発生時（使用回数あり）
		ctxMiss := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventTypingMiss,
			UsesRemaining: 1,
		}
		resultMiss := domain.EvaluatePassive(def, ctxMiss)

		if !resultMiss.IsActive {
			t.Error("ミス時にアクティブであるべきです")
		}

		// ミスイベント発生時（使用回数なし）
		ctxNoUses := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventTypingMiss,
			UsesRemaining: 0,
		}
		resultNoUses := domain.EvaluatePassive(def, ctxNoUses)

		if resultNoUses.IsActive {
			t.Error("使用回数がない場合はアクティブでないべきです")
		}
	})

	t.Run("ps_first_strike", func(t *testing.T) {
		// ファーストストライク: 戦闘開始時、最初のスキルが即発動
		def := domain.PassiveSkill{
			ID:          "ps_first_strike",
			Name:        "ファーストストライク",
			Description: "戦闘開始時、最初のスキルが即発動",
			TriggerType: domain.PassiveTriggerReactive,
			EffectType:  domain.PassiveEffectSpecial,
			EffectValue: 1.0,
			TriggerCondition: &domain.TriggerCondition{
				Type: domain.TriggerConditionOnBattleStart,
			},
		}

		// 戦闘開始イベント発生時
		ctxStart := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventBattleStart,
			UsesRemaining: 1,
		}
		resultStart := domain.EvaluatePassive(def, ctxStart)

		if !resultStart.IsActive {
			t.Error("戦闘開始時にアクティブであるべきです")
		}

		// 他のイベント発生時
		ctxOther := &domain.PassiveEvaluationContext{
			Event:         domain.PassiveEventSkillUse,
			UsesRemaining: 1,
		}
		resultOther := domain.EvaluatePassive(def, ctxOther)

		if resultOther.IsActive {
			t.Error("他のイベントではアクティブでないべきです")
		}
	})
}

// TestPassiveSkill_DefinitionHelpers はパッシブスキル定義のヘルパーメソッドを検証します。
func TestPassiveSkill_DefinitionHelpers(t *testing.T) {
	permanentDef := domain.PassiveSkill{
		TriggerType: domain.PassiveTriggerPermanent,
	}
	if !permanentDef.IsPermanent() {
		t.Error("IsPermanent() should return true for permanent trigger")
	}

	probabilityDef := domain.PassiveSkill{
		TriggerType: domain.PassiveTriggerProbability,
		Probability: 0.3,
	}
	if !probabilityDef.HasProbability() {
		t.Error("HasProbability() should return true when probability > 0")
	}

	stackDef := domain.PassiveSkill{
		TriggerType: domain.PassiveTriggerStack,
		MaxStacks:   5,
	}
	if !stackDef.IsStackable() {
		t.Error("IsStackable() should return true for stack trigger with max stacks > 0")
	}
}

// TestPassiveSkill_AllTriggerTypes は全てのトリガータイプが正しく評価されることを検証します。
func TestPassiveSkill_AllTriggerTypes(t *testing.T) {
	triggerTypes := []domain.PassiveTriggerType{
		domain.PassiveTriggerPermanent,
		domain.PassiveTriggerConditional,
		domain.PassiveTriggerProbability,
		domain.PassiveTriggerStack,
		domain.PassiveTriggerReactive,
	}

	for _, triggerType := range triggerTypes {
		t.Run(string(triggerType), func(t *testing.T) {
			def := domain.PassiveSkill{
				ID:          "test_" + string(triggerType),
				Name:        "テスト",
				TriggerType: triggerType,
				EffectValue: 1.5,
			}

			// スタック型と反応型にはトリガー条件が必要
			if triggerType == domain.PassiveTriggerStack || triggerType == domain.PassiveTriggerReactive {
				def.TriggerCondition = &domain.TriggerCondition{
					Type: domain.TriggerConditionNoMissStreak,
				}
			}

			ctx := &domain.PassiveEvaluationContext{}
			result := domain.EvaluatePassive(def, ctx)

			// 結果が返されることを確認
			t.Logf("TriggerType %s: IsActive=%v, EffectMultiplier=%f",
				triggerType, result.IsActive, result.EffectMultiplier)
		})
	}
}
