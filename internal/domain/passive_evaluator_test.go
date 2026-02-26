// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// TestPassiveEvaluator_永続効果_バフエクステンダー は永続効果タイプのパッシブスキルが正しく評価されることを確認します。
func TestPassiveEvaluator_永続効果_バフエクステンダー(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_buff_extender",
		Name:        "バフエクステンダー",
		Description: "バフ効果時間+50%",
		TriggerType: PassiveTriggerPermanent,
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	ctx := &PassiveEvaluationContext{}
	result := EvaluatePassive(def, ctx)

	if !result.IsActive {
		t.Error("永続効果は常にActiveであるべきです")
	}
	if result.EffectMultiplier != 1.5 {
		t.Errorf("EffectMultiplierが期待値と異なります: got %f, want 1.5", result.EffectMultiplier)
	}
}

// TestPassiveEvaluator_条件付き_パーフェクトリズム_条件満たす は条件を満たした時に効果が発動することを確認します。
func TestPassiveEvaluator_条件付き_パーフェクトリズム_条件満たす(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_perfect_rhythm",
		Name:        "パーフェクトリズム",
		Description: "正確性100%でスキル効果1.5倍",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	ctx := &PassiveEvaluationContext{
		IsPerfect: true,
	}
	result := EvaluatePassive(def, ctx)

	if !result.IsActive {
		t.Error("パーフェクトの場合、パーフェクトリズムは発動するべきです")
	}
	if result.EffectMultiplier != 1.5 {
		t.Errorf("EffectMultiplierが期待値と異なります: got %f, want 1.5", result.EffectMultiplier)
	}
}

// TestPassiveEvaluator_条件付き_パーフェクトリズム_条件満たさない は条件を満たさない時に効果が発動しないことを確認します。
func TestPassiveEvaluator_条件付き_パーフェクトリズム_条件満たさない(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_perfect_rhythm",
		Name:        "パーフェクトリズム",
		Description: "正確性100%でスキル効果1.5倍",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	ctx := &PassiveEvaluationContext{
		IsPerfect: false,
	}
	result := EvaluatePassive(def, ctx)

	if result.IsActive {
		t.Error("非パーフェクトの場合、パーフェクトリズムは発動しないべきです")
	}
}

// TestPassiveEvaluator_条件付き_スピードブレイク_WPM80以上 はWPM条件を満たした時に効果が発動することを確認します。
func TestPassiveEvaluator_条件付き_スピードブレイク_WPM80以上(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_speed_break",
		Name:        "スピードブレイク",
		Description: "WPM80以上で25%追加ダメージ",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionWPMAbove,
			Value: 80,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.25,
	}

	tests := []struct {
		name     string
		wpm      float64
		expected bool
	}{
		{"WPM80ちょうど", 80, true},
		{"WPM100", 100, true},
		{"WPM79", 79, false},
		{"WPM0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				WPM: tt.wpm,
			}
			result := EvaluatePassive(def, ctx)
			if result.IsActive != tt.expected {
				t.Errorf("IsActiveが期待値と異なります: got %v, want %v", result.IsActive, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluator_条件付き_エンドゲームスペシャリスト は敵HP条件を満たした時に効果が発動することを確認します。
func TestPassiveEvaluator_条件付き_エンドゲームスペシャリスト(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_endgame_specialist",
		Name:        "エンドゲームスペシャリスト",
		Description: "敵HP30%以下で全ダメージ+25%",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionEnemyHPBelowPercent,
			Value: 30,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.25,
	}

	tests := []struct {
		name           string
		enemyHPPercent float64
		expected       bool
	}{
		{"敵HP30%", 30, true},
		{"敵HP25%", 25, true},
		{"敵HP31%", 31, false},
		{"敵HP100%", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				EnemyHPPercent: tt.enemyHPPercent,
			}
			result := EvaluatePassive(def, ctx)
			if result.IsActive != tt.expected {
				t.Errorf("IsActiveが期待値と異なります: got %v, want %v", result.IsActive, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluator_条件付き_ウィークポイント は敵デバフ状態で効果が発動することを確認します。
func TestPassiveEvaluator_条件付き_ウィークポイント(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_weak_point",
		Name:        "ウィークポイント",
		Description: "デバフ中の敵へダメージ+20%",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionEnemyHasDebuff,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.2,
	}

	tests := []struct {
		name           string
		enemyHasDebuff bool
		expected       bool
	}{
		{"敵デバフあり", true, true},
		{"敵デバフなし", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				EnemyHasDebuff: tt.enemyHasDebuff,
			}
			result := EvaluatePassive(def, ctx)
			if result.IsActive != tt.expected {
				t.Errorf("IsActiveが期待値と異なります: got %v, want %v", result.IsActive, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluator_条件付き_オーバードライブ はHP条件で複合効果が発動することを確認します。
func TestPassiveEvaluator_条件付き_オーバードライブ(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_overdrive",
		Name:        "オーバードライブ",
		Description: "HP50%以下でリキャスト-30%、被ダメ+20%",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionHPBelowPercent,
			Value: 50,
		},
		EffectType:  PassiveEffectSpecial,
		EffectValue: 0.3, // リキャスト短縮率
	}

	tests := []struct {
		name            string
		playerHPPercent float64
		expected        bool
	}{
		{"HP50%", 50, true},
		{"HP25%", 25, true},
		{"HP51%", 51, false},
		{"HP100%", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				PlayerHPPercent: tt.playerHPPercent,
			}
			result := EvaluatePassive(def, ctx)
			if result.IsActive != tt.expected {
				t.Errorf("IsActiveが期待値と異なります: got %v, want %v", result.IsActive, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluator_確率トリガー_ラストスタンド は確率トリガーの基本動作を確認します。
func TestPassiveEvaluator_確率トリガー_ラストスタンド(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_last_stand",
		Name:        "ラストスタンド",
		Description: "HP25%以下で30%の確率で被ダメージ1",
		TriggerType: PassiveTriggerProbability,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionHPBelowPercent,
			Value: 25,
		},
		EffectType:  PassiveEffectSpecial,
		EffectValue: 1,
		Probability: 0.3,
	}

	// HP条件を満たさない場合は確率判定なしで不発動
	ctx := &PassiveEvaluationContext{
		PlayerHPPercent: 50,
	}
	result := EvaluatePassive(def, ctx)
	if result.IsActive {
		t.Error("HP条件を満たさない場合は発動しないべきです")
	}
	if result.NeedsProbabilityCheck {
		t.Error("HP条件を満たさない場合は確率チェックも不要です")
	}

	// HP条件を満たす場合は確率チェックが必要
	ctx = &PassiveEvaluationContext{
		PlayerHPPercent: 20,
	}
	result = EvaluatePassive(def, ctx)
	if !result.NeedsProbabilityCheck {
		t.Error("HP条件を満たす場合は確率チェックが必要です")
	}
	if result.Probability != 0.3 {
		t.Errorf("Probabilityが期待値と異なります: got %f, want 0.3", result.Probability)
	}
}

// TestPassiveEvaluator_確率トリガー_カウンターチャージ は被ダメージ時の確率トリガーを確認します。
func TestPassiveEvaluator_確率トリガー_カウンターチャージ(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_counter_charge",
		Name:        "カウンターチャージ",
		Description: "被ダメージ時20%で次の攻撃2倍",
		TriggerType: PassiveTriggerProbability,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnDamageReceived,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 2.0,
		Probability: 0.2,
	}

	// イベント発生時のみ確率チェック
	ctx := &PassiveEvaluationContext{
		Event: PassiveEventDamageReceived,
	}
	result := EvaluatePassive(def, ctx)
	if !result.NeedsProbabilityCheck {
		t.Error("被ダメージイベント時は確率チェックが必要です")
	}

	// イベントなしの場合はスキップ
	ctx = &PassiveEvaluationContext{}
	result = EvaluatePassive(def, ctx)
	if result.NeedsProbabilityCheck {
		t.Error("イベントなしの場合は確率チェック不要です")
	}
}

// TestPassiveEvaluator_スタック型_コンボマスター はスタック型の基本動作を確認します。
func TestPassiveEvaluator_スタック型_コンボマスター(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_combo_master",
		Name:        "コンボマスター",
		Description: "ミスなし連続タイピングでダメージ累積+10%（最大+50%）",
		TriggerType: PassiveTriggerStack,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionNoMissStreak,
		},
		EffectType:     PassiveEffectModifier,
		EffectValue:    0.1,
		MaxStacks:      5,
		StackIncrement: 0.1,
	}

	tests := []struct {
		name          string
		currentStacks int
		expectedMult  float64
	}{
		{"スタック0", 0, 1.0},
		{"スタック1", 1, 1.1},
		{"スタック3", 3, 1.3},
		{"スタック5（最大）", 5, 1.5},
		{"スタック6（上限超過）", 6, 1.5}, // 最大で5スタック
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				CurrentStacks: tt.currentStacks,
			}
			result := EvaluatePassive(def, ctx)

			if result.EffectMultiplier != tt.expectedMult {
				t.Errorf("EffectMultiplierが期待値と異なります: got %f, want %f", result.EffectMultiplier, tt.expectedMult)
			}
		})
	}
}

// TestPassiveEvaluator_スタック型_アダプティブシールド はカウンター型の動作を確認します。
func TestPassiveEvaluator_スタック型_アダプティブシールド(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_adaptive_shield",
		Name:        "アダプティブシールド",
		Description: "同種攻撃3回目以降ダメージ-25%",
		TriggerType: PassiveTriggerStack,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionSameAttackCount,
			Value: 3, // 3回目から発動
		},
		EffectType:  PassiveEffectModifier,
		EffectValue: 0.25, // 25%軽減
	}

	tests := []struct {
		name        string
		attackCount int
		expected    bool
	}{
		{"1回目", 1, false},
		{"2回目", 2, false},
		{"3回目", 3, true},
		{"4回目", 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &PassiveEvaluationContext{
				SameAttackCount: tt.attackCount,
			}
			result := EvaluatePassive(def, ctx)
			if result.IsActive != tt.expected {
				t.Errorf("IsActiveが期待値と異なります: got %v, want %v", result.IsActive, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluator_反応型_ファーストストライク は戦闘開始時トリガーを確認します。
func TestPassiveEvaluator_反応型_ファーストストライク(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_first_strike",
		Name:        "ファーストストライク",
		Description: "戦闘開始時、最初のスキルが即発動",
		TriggerType: PassiveTriggerReactive,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnBattleStart,
		},
		EffectType:    PassiveEffectSpecial,
		UsesPerBattle: 1,
	}

	// 戦闘開始イベント
	ctx := &PassiveEvaluationContext{
		Event:         PassiveEventBattleStart,
		UsesRemaining: 1,
	}
	result := EvaluatePassive(def, ctx)
	if !result.IsActive {
		t.Error("戦闘開始時は発動するべきです")
	}

	// 使用回数切れ
	ctx = &PassiveEvaluationContext{
		Event:         PassiveEventBattleStart,
		UsesRemaining: 0,
	}
	result = EvaluatePassive(def, ctx)
	if result.IsActive {
		t.Error("使用回数切れの場合は発動しないべきです")
	}
}

// TestPassiveEvaluator_反応型_デバフアブソーバー はデバフ受け時の反応を確認します。
func TestPassiveEvaluator_反応型_デバフアブソーバー(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_debuff_absorber",
		Name:        "デバフアブソーバー",
		Description: "デバフ効果時間半減＋小回復",
		TriggerType: PassiveTriggerReactive,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnDebuffReceived,
		},
		EffectType:  PassiveEffectSpecial,
		EffectValue: 0.5, // 効果時間半減
	}

	// デバフ受けイベント
	ctx := &PassiveEvaluationContext{
		Event: PassiveEventDebuffReceived,
	}
	result := EvaluatePassive(def, ctx)
	if !result.IsActive {
		t.Error("デバフ受け時は発動するべきです")
	}

	// 別のイベント
	ctx = &PassiveEvaluationContext{
		Event: PassiveEventDamageReceived,
	}
	result = EvaluatePassive(def, ctx)
	if result.IsActive {
		t.Error("別のイベントでは発動しないべきです")
	}
}

// TestPassiveEvaluator_反応型_タイポリカバリー はタイピングミス時の反応を確認します。
func TestPassiveEvaluator_反応型_タイポリカバリー(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_typo_recovery",
		Name:        "タイポリカバリー",
		Description: "ミス時制限時間+1秒（1回/チャレンジ）",
		TriggerType: PassiveTriggerReactive,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnTypingMiss,
		},
		EffectType:    PassiveEffectModifier,
		EffectValue:   1.0, // 1秒追加
		UsesPerBattle: 1,
	}

	// タイピングミスイベント（使用回数あり）
	ctx := &PassiveEvaluationContext{
		Event:         PassiveEventTypingMiss,
		UsesRemaining: 1,
	}
	result := EvaluatePassive(def, ctx)
	if !result.IsActive {
		t.Error("タイピングミス時は発動するべきです")
	}

	// 使用回数切れ
	ctx = &PassiveEvaluationContext{
		Event:         PassiveEventTypingMiss,
		UsesRemaining: 0,
	}
	result = EvaluatePassive(def, ctx)
	if result.IsActive {
		t.Error("使用回数切れの場合は発動しないべきです")
	}
}

// TestPassiveEvent_定数の確認 はPassiveEvent定数が正しく定義されていることを確認します。
func TestPassiveEvent_定数の確認(t *testing.T) {
	tests := []struct {
		name     string
		event    PassiveEvent
		expected string
	}{
		{"なし", PassiveEventNone, "none"},
		{"戦闘開始", PassiveEventBattleStart, "battle_start"},
		{"スキル使用", PassiveEventSkillUse, "skill_use"},
		{"被ダメージ", PassiveEventDamageReceived, "damage_received"},
		{"回復", PassiveEventHeal, "heal"},
		{"バフ/デバフ使用", PassiveEventBuffDebuffUse, "buff_debuff_use"},
		{"物理攻撃", PassiveEventPhysicalAttack, "physical_attack"},
		{"タイピングミス", PassiveEventTypingMiss, "typing_miss"},
		{"時間切れ", PassiveEventTimeout, "timeout"},
		{"デバフ受け", PassiveEventDebuffReceived, "debuff_received"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.event) != tt.expected {
				t.Errorf("PassiveEventが期待値と異なります: got %s, want %s", tt.event, tt.expected)
			}
		})
	}
}

// TestPassiveEvaluationResult_フィールドの確認 はPassiveEvaluationResultのフィールドを確認します。
func TestPassiveEvaluationResult_フィールドの確認(t *testing.T) {
	result := PassiveEvaluationResult{
		IsActive:              true,
		EffectMultiplier:      1.5,
		EffectValue:           10,
		NeedsProbabilityCheck: true,
		Probability:           0.3,
	}

	if !result.IsActive {
		t.Error("IsActiveが期待値と異なります")
	}
	if result.EffectMultiplier != 1.5 {
		t.Errorf("EffectMultiplierが期待値と異なります: got %f, want 1.5", result.EffectMultiplier)
	}
	if result.EffectValue != 10 {
		t.Errorf("EffectValueが期待値と異なります: got %f, want 10", result.EffectValue)
	}
	if !result.NeedsProbabilityCheck {
		t.Error("NeedsProbabilityCheckが期待値と異なります")
	}
	if result.Probability != 0.3 {
		t.Errorf("Probabilityが期待値と異なります: got %f, want 0.3", result.Probability)
	}
}

// TestPassiveEvaluator_確率トリガー_モック化テスト は確率判定のモック化テストを行います。
func TestPassiveEvaluator_確率トリガー_モック化テスト(t *testing.T) {
	def := PassiveSkill{
		ID:          "ps_echo_skill",
		Name:        "エコースキル",
		Description: "15%の確率でスキル2回発動",
		TriggerType: PassiveTriggerProbability,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnSkillUse,
		},
		EffectType:  PassiveEffectSpecial,
		Probability: 0.15,
	}

	// スキル使用イベント
	ctx := &PassiveEvaluationContext{
		Event: PassiveEventSkillUse,
	}
	result := EvaluatePassive(def, ctx)

	// 確率チェックが必要であることを確認
	if !result.NeedsProbabilityCheck {
		t.Error("スキル使用イベント時は確率チェックが必要です")
	}
	if result.Probability != 0.15 {
		t.Errorf("Probabilityが期待値と異なります: got %f, want 0.15", result.Probability)
	}

	// 発動判定はユースケース層で行うので、ここでは確率値の返却を確認
	// モック化テスト: 確率0の場合は発動しない
	defNoChance := PassiveSkill{
		ID:          "test_no_chance",
		TriggerType: PassiveTriggerProbability,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnSkillUse,
		},
		Probability: 0,
	}
	resultNoChance := EvaluatePassive(defNoChance, ctx)
	// Probability 0 でも確率チェックが必要（実際の発動判定はユースケース層）
	if !resultNoChance.NeedsProbabilityCheck {
		t.Error("確率0でも確率チェックフラグは立つべきです")
	}
}

// TestPassiveEvaluator_複数パッシブスキル併存 は複数のパッシブスキルが独立して評価されることを確認します。
func TestPassiveEvaluator_複数パッシブスキル併存(t *testing.T) {
	// 永続効果
	buffExtender := PassiveSkill{
		ID:          "ps_buff_extender",
		TriggerType: PassiveTriggerPermanent,
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	// 条件付き効果
	perfectRhythm := PassiveSkill{
		ID:          "ps_perfect_rhythm",
		TriggerType: PassiveTriggerConditional,
		TriggerCondition: &TriggerCondition{
			Type:  TriggerConditionAccuracyEquals,
			Value: 100,
		},
		EffectType:  PassiveEffectMultiplier,
		EffectValue: 1.5,
	}

	// 同一コンテキストで両方を評価
	ctx := &PassiveEvaluationContext{
		IsPerfect: true,
	}

	result1 := EvaluatePassive(buffExtender, ctx)
	result2 := EvaluatePassive(perfectRhythm, ctx)

	// 両方が独立してアクティブ
	if !result1.IsActive {
		t.Error("buffExtenderはアクティブであるべきです")
	}
	if !result2.IsActive {
		t.Error("perfectRhythmはアクティブであるべきです（正確性100%）")
	}

	// 効果倍率の合成は上位層の責務だが、独立した値を返すことを確認
	if result1.EffectMultiplier != 1.5 {
		t.Errorf("buffExtenderのEffectMultiplierが期待値と異なります: got %f, want 1.5", result1.EffectMultiplier)
	}
	if result2.EffectMultiplier != 1.5 {
		t.Errorf("perfectRhythmのEffectMultiplierが期待値と異なります: got %f, want 1.5", result2.EffectMultiplier)
	}

	// 非パーフェクトの場合、perfectRhythmは不発動
	ctx2 := &PassiveEvaluationContext{
		IsPerfect: false,
	}
	result3 := EvaluatePassive(buffExtender, ctx2)
	result4 := EvaluatePassive(perfectRhythm, ctx2)

	if !result3.IsActive {
		t.Error("buffExtenderは常にアクティブであるべきです")
	}
	if result4.IsActive {
		t.Error("perfectRhythmは正確性95%では不発動であるべきです")
	}
}

// TestPassiveEvaluator_条件付きnil はTriggerConditionがnilの場合の処理を確認します。
func TestPassiveEvaluator_条件付きnil(t *testing.T) {
	def := PassiveSkill{
		ID:               "test_no_condition",
		TriggerType:      PassiveTriggerConditional,
		TriggerCondition: nil, // 条件なし
		EffectType:       PassiveEffectMultiplier,
		EffectValue:      1.5,
	}

	ctx := &PassiveEvaluationContext{}
	result := EvaluatePassive(def, ctx)

	// 条件がない場合は発動しない
	if result.IsActive {
		t.Error("TriggerConditionがnilの場合は発動しないべきです")
	}
}

// TestPassiveEvaluator_スタック型_条件nil はスタック型でTriggerConditionがnilの場合の処理を確認します。
func TestPassiveEvaluator_スタック型_条件nil(t *testing.T) {
	def := PassiveSkill{
		ID:               "test_stack_no_condition",
		TriggerType:      PassiveTriggerStack,
		TriggerCondition: nil, // 条件なし
		EffectValue:      0.1,
		MaxStacks:        5,
	}

	ctx := &PassiveEvaluationContext{
		CurrentStacks: 3,
	}
	result := EvaluatePassive(def, ctx)

	// 条件がない場合はデフォルト値
	if result.IsActive {
		t.Error("TriggerConditionがnilの場合は発動しないべきです")
	}
}

// TestPassiveEvaluator_反応型_条件nil は反応型でTriggerConditionがnilの場合の処理を確認します。
func TestPassiveEvaluator_反応型_条件nil(t *testing.T) {
	def := PassiveSkill{
		ID:               "test_reactive_no_condition",
		TriggerType:      PassiveTriggerReactive,
		TriggerCondition: nil, // 条件なし
		EffectType:       PassiveEffectSpecial,
	}

	ctx := &PassiveEvaluationContext{
		Event: PassiveEventBattleStart,
	}
	result := EvaluatePassive(def, ctx)

	// 条件がない場合は発動しない
	if result.IsActive {
		t.Error("TriggerConditionがnilの場合は発動しないべきです")
	}
}

// TestPassiveEvaluator_使用回数無制限 は使用回数制限なし（UsesPerBattle=0）の処理を確認します。
func TestPassiveEvaluator_使用回数無制限(t *testing.T) {
	def := PassiveSkill{
		ID:          "test_unlimited_uses",
		TriggerType: PassiveTriggerReactive,
		TriggerCondition: &TriggerCondition{
			Type: TriggerConditionOnDamageReceived,
		},
		EffectType:    PassiveEffectSpecial,
		UsesPerBattle: 0, // 無制限
	}

	ctx := &PassiveEvaluationContext{
		Event:         PassiveEventDamageReceived,
		UsesRemaining: 0, // 使用回数切れ状態でもUsesPerBattle=0なら発動
	}
	result := EvaluatePassive(def, ctx)

	// UsesPerBattle=0の場合は使用回数チェックなしで発動
	if !result.IsActive {
		t.Error("UsesPerBattle=0の場合は使用回数制限なしで発動するべきです")
	}
}
