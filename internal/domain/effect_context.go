// Package domain はゲームのドメインモデルを定義します。
package domain

// EffectEventType はイベントの種類を表します。
// 効果の発動条件を判定する際に使用します。
type EffectEventType string

const (
	// EventNone はイベントなしを表します。
	EventNone EffectEventType = ""

	// EventOnDamageDealt はダメージを与えた時のイベントを表します。
	EventOnDamageDealt EffectEventType = "on_damage_dealt"

	// EventOnDamageRecv はダメージを受けた時のイベントを表します。
	EventOnDamageRecv EffectEventType = "on_damage_recv"

	// EventOnHeal は回復した時のイベントを表します。
	EventOnHeal EffectEventType = "on_heal"

	// EventOnSkillUse はスキル使用時のイベントを表します。
	EventOnSkillUse EffectEventType = "on_skill_use"

	// EventOnTypingDone はタイピング完了時のイベントを表します。
	EventOnTypingDone EffectEventType = "on_typing_done"

	// EventOnTypingMiss はタイピングミス時のイベントを表します。
	EventOnTypingMiss EffectEventType = "on_typing_miss"

	// EventOnBattleStart はバトル開始時のイベントを表します。
	EventOnBattleStart EffectEventType = "on_battle_start"

	// EventOnTimeout は時間切れ時のイベントを表します。
	EventOnTimeout EffectEventType = "on_timeout"
)

// EffectContext は EnableCondition の判定に使う状態を表します。
// バトル中の各種状態を保持し、効果の有効/無効判定に使用します。
type EffectContext struct {
	// ========== プレイヤー状態 ==========

	// PlayerHP はプレイヤー現在HPです。
	PlayerHP int

	// PlayerMaxHP はプレイヤー最大HPです。
	PlayerMaxHP int

	// PlayerHPPercent はプレイヤーHP割合 (0.0〜1.0) です。
	PlayerHPPercent float64

	// ========== 敵状態 ==========

	// EnemyHP は敵現在HPです。
	EnemyHP int

	// EnemyMaxHP は敵最大HPです。
	EnemyMaxHP int

	// EnemyHPPercent は敵HP割合 (0.0〜1.0) です。
	EnemyHPPercent float64

	// EnemyHasDebuff は敵にデバフがあるかを表します。
	EnemyHasDebuff bool

	// ========== タイピング結果 ==========

	// Accuracy は正確性 (0.0〜1.0) です。
	Accuracy float64

	// WPM はWords Per Minuteです。
	WPM float64

	// ComboCount はミスなし連続回数です。
	ComboCount int

	// MissCount は累計ミス回数です。
	MissCount int

	// ========== チェイン発動判定用 ==========

	// LastSkillAgent は最後にスキルを使ったエージェント番号です（-1: なし）。
	LastSkillAgent int

	// CurrentAgent は現在評価中のエージェント番号です。
	CurrentAgent int

	// ========== イベント情報 ==========

	// EventType は現在発生中のイベントです。
	EventType EffectEventType

	// DamageDealt は与えたダメージ（EventOnDamageDealt時）です。
	DamageDealt int

	// DamageReceived は受けたダメージ（EventOnDamageRecv時）です。
	DamageReceived int

	// HealAmount は回復量（EventOnHeal時）です。
	HealAmount int

	// ========== スキル情報 ==========

	// IsDamageSkill はダメージ効果を持つスキルかを表します。
	IsDamageSkill bool

	// IsHealSkill は回復効果を持つスキルかを表します。
	IsHealSkill bool

	// IsPhysical は物理攻撃かどうかを表します。
	IsPhysical bool

	// HasBuffDebuffEffect はバフまたはデバフ効果を持つスキルかを表します。
	HasBuffDebuffEffect bool

	// ========== 状態カウンタ ==========

	// SameAttackCount は同種攻撃を受けた回数です。
	SameAttackCount int

	// UsesRemaining はスキルごとの残り使用回数です。
	UsesRemaining map[string]int
}

// NewEffectContext はバトル状態からコンテキストを生成します。
func NewEffectContext(playerHP, playerMaxHP, enemyHP, enemyMaxHP int) *EffectContext {
	var playerPercent, enemyPercent float64
	if playerMaxHP > 0 {
		playerPercent = float64(playerHP) / float64(playerMaxHP)
	}
	if enemyMaxHP > 0 {
		enemyPercent = float64(enemyHP) / float64(enemyMaxHP)
	}

	return &EffectContext{
		PlayerHP:        playerHP,
		PlayerMaxHP:     playerMaxHP,
		PlayerHPPercent: playerPercent,
		EnemyHP:         enemyHP,
		EnemyMaxHP:      enemyMaxHP,
		EnemyHPPercent:  enemyPercent,
		LastSkillAgent:  -1,
		UsesRemaining:   make(map[string]int),
	}
}

// UpdateHP はHP情報を更新します。
func (ctx *EffectContext) UpdateHP(playerHP, enemyHP int) {
	ctx.PlayerHP = playerHP
	if ctx.PlayerMaxHP > 0 {
		ctx.PlayerHPPercent = float64(playerHP) / float64(ctx.PlayerMaxHP)
	}
	ctx.EnemyHP = enemyHP
	if ctx.EnemyMaxHP > 0 {
		ctx.EnemyHPPercent = float64(enemyHP) / float64(ctx.EnemyMaxHP)
	}
}

// SetTypingResult はタイピング結果を設定します。
func (ctx *EffectContext) SetTypingResult(accuracy float64, wpm float64, combo int) {
	ctx.Accuracy = accuracy
	ctx.WPM = wpm
	ctx.ComboCount = combo
}

// SetEvent はイベント情報を設定します。
func (ctx *EffectContext) SetEvent(eventType EffectEventType) {
	ctx.EventType = eventType
}

// SetDamageDealt は与ダメージイベントを設定します。
func (ctx *EffectContext) SetDamageDealt(damage int) {
	ctx.EventType = EventOnDamageDealt
	ctx.DamageDealt = damage
}

// SetDamageReceived は被ダメージイベントを設定します。
func (ctx *EffectContext) SetDamageReceived(damage int) {
	ctx.EventType = EventOnDamageRecv
	ctx.DamageReceived = damage
}

// SetHeal は回復イベントを設定します。
func (ctx *EffectContext) SetHeal(amount int) {
	ctx.EventType = EventOnHeal
	ctx.HealAmount = amount
}

// SetSkillUse はスキル使用イベントを設定します。
func (ctx *EffectContext) SetSkillUse(agentIndex int, isDamage, isHeal, isPhysical, hasBuffDebuff bool) {
	ctx.EventType = EventOnSkillUse
	ctx.LastSkillAgent = agentIndex
	ctx.IsDamageSkill = isDamage
	ctx.IsHealSkill = isHeal
	ctx.IsPhysical = isPhysical
	ctx.HasBuffDebuffEffect = hasBuffDebuff
}

// Clone はコンテキストのコピーを作成します。
func (ctx *EffectContext) Clone() *EffectContext {
	clone := *ctx
	clone.UsesRemaining = make(map[string]int)
	for k, v := range ctx.UsesRemaining {
		clone.UsesRemaining[k] = v
	}
	return &clone
}
