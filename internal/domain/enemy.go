// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"time"
)

// EnhanceThreshold は敵が強化フェーズに移行するHP割合の閾値です（50%）。
const EnhanceThreshold = 0.5

// EnemyPhase は敵のフェーズを表す型です。
type EnemyPhase int

const (
	// PhaseNormal は通常フェーズです（HP50%以上）
	PhaseNormal EnemyPhase = 0

	// PhaseEnhanced は強化フェーズです（HP50%以下、特殊攻撃解禁）
	PhaseEnhanced EnemyPhase = 1
)

// String はEnemyPhaseの日本語表示名を返します。
func (p EnemyPhase) String() string {
	switch p {
	case PhaseNormal:
		return "通常"
	case PhaseEnhanced:
		return "強化"
	default:
		return "不明"
	}
}

// EnemyWaitMode は敵の待機状態を表す型です。
// 敵はチャージ中（行動実行待ち）またはディフェンス中（防御状態）のどちらかの状態しか持ちません。
type EnemyWaitMode int

const (
	// WaitModeNone は待機状態ではない（行動実行中）ことを示します。
	WaitModeNone EnemyWaitMode = iota

	// WaitModeCharging はチャージ中（行動実行待ち）であることを示します。
	WaitModeCharging

	// WaitModeDefending はディフェンス中（防御状態）であることを示します。
	WaitModeDefending
)

// String はEnemyWaitModeの日本語表示名を返します。
func (m EnemyWaitMode) String() string {
	switch m {
	case WaitModeNone:
		return "なし"
	case WaitModeCharging:
		return "チャージ中"
	case WaitModeDefending:
		return "ディフェンス中"
	default:
		return "不明"
	}
}

// EnemyType は敵の種類（タイプ）を定義する構造体です。
// 外部データファイル（enemies.json）から読み込まれます。
type EnemyType struct {
	// ID は敵タイプの一意識別子です。
	ID string

	// Name は敵タイプの表示名です（日本語）。
	Name string

	// BaseHP は敵の基礎HP値です。
	BaseHP int

	// BaseAttackPower は敵の基礎攻撃力です。
	BaseAttackPower int

	// AttackType は攻撃属性（physical / magic）です。
	AttackType string

	// ASCIIArt は敵の外観（ASCIIアート）です。
	ASCIIArt string

	// ========== 拡張フィールド ==========

	// DefaultLevel はデフォルトレベル（1〜100）です。未撃破時はこのレベルのみ選択可能。
	DefaultLevel int

	// NormalActionPatternIDs は通常状態での行動パターンIDの配列です。
	NormalActionPatternIDs []string

	// EnhancedActionPatternIDs は強化状態での行動パターンIDの配列です。空の場合は通常パターンを継続。
	EnhancedActionPatternIDs []string

	// ResolvedNormalActions は解決済みの通常行動パターンです（ランタイムで設定）。
	ResolvedNormalActions []EnemyAction

	// ResolvedEnhancedActions は解決済みの強化行動パターンです（ランタイムで設定）。
	ResolvedEnhancedActions []EnemyAction

	// NormalPassive は通常状態で適用されるパッシブスキルです。
	NormalPassive *EnemyPassiveSkill

	// EnhancedPassive は強化状態で適用されるパッシブスキルです。
	EnhancedPassive *EnemyPassiveSkill

	// DropItemCategory はドロップアイテムのカテゴリ（"core" または "module"）です。
	DropItemCategory string

	// DropItemTypeID はドロップアイテムのTypeIDです。
	DropItemTypeID string

	// ========== ボルテージシステム ==========

	// VoltageRisePer10s は10秒間でのボルテージ上昇量です。
	// 0の場合はボルテージが上昇しません。デフォルト値は10（infra層で設定）。
	VoltageRisePer10s float64

	// ========== ランクシステム ==========

	// Rank は敵のランクです（1から開始）。
	// 同じランクの敵を全て撃破すると次のランクが解放されます。
	Rank int
}

// IsValidDefaultLevel はデフォルトレベルが有効範囲（1〜100）かどうかを判定します。
func (e EnemyType) IsValidDefaultLevel() bool {
	return e.DefaultLevel >= 1 && e.DefaultLevel <= 100
}

// HasValidNormalActionPattern は通常行動パターンが有効（最低1つの行動を持つ）かどうかを判定します。
func (e EnemyType) HasValidNormalActionPattern() bool {
	return len(e.NormalActionPatternIDs) > 0 || len(e.ResolvedNormalActions) > 0
}

// SetResolvedActions は行動IDから解決された行動パターンを設定します。
func (e *EnemyType) SetResolvedActions(normalActions, enhancedActions []EnemyAction) {
	e.ResolvedNormalActions = normalActions
	e.ResolvedEnhancedActions = enhancedActions
}

// GetNormalActions は通常行動パターンを返します。
func (e EnemyType) GetNormalActions() []EnemyAction {
	return e.ResolvedNormalActions
}

// GetEnhancedActions は強化行動パターンを返します。
func (e EnemyType) GetEnhancedActions() []EnemyAction {
	return e.ResolvedEnhancedActions
}

// IsValidDropItemCategory はドロップカテゴリが有効かどうかを判定します。
func (e EnemyType) IsValidDropItemCategory() bool {
	return e.DropItemCategory == "core" || e.DropItemCategory == "module"
}

// GetVoltageRisePer10s は10秒あたりのボルテージ上昇量を返します。
// 負の値が設定されている場合は0を返します。
func (e EnemyType) GetVoltageRisePer10s() float64 {
	if e.VoltageRisePer10s < 0 {
		return 0
	}
	return e.VoltageRisePer10s
}

// IsValidRank はランクが有効範囲（1以上）かどうかを判定します。
func (e EnemyType) IsValidRank() bool {
	return e.Rank >= 1
}

// EnemyModel はゲーム内の敵エンティティを表す構造体です。
type EnemyModel struct {
	// ID は敵インスタンスの一意識別子です。
	ID string

	// Name は敵の表示名です。
	Name string

	// Level は敵のレベルです。
	Level int

	// HP は敵の現在HP値です。
	HP int

	// MaxHP は敵の最大HP値です。
	MaxHP int

	// AttackPower は敵の攻撃力です。
	AttackPower int

	// Type は敵の種類（タイプ）です。
	Type EnemyType

	// Phase は敵の現在フェーズです。
	Phase EnemyPhase

	// EffectTable は敵に適用されているステータス効果テーブルです。
	// 敵自身のバフやプレイヤーからのデバフを管理します。
	EffectTable *EffectTable

	// ========== 行動管理フィールド ==========

	// ActionIndex は現在の行動パターンインデックスです。
	ActionIndex int

	// ActivePassiveID は現在適用中のパッシブスキルIDです（解除時に使用）。
	ActivePassiveID string

	// ========== 待機状態管理フィールド ==========

	// WaitMode は敵の現在の待機状態を示します（チャージ中/ディフェンス中）。
	WaitMode EnemyWaitMode

	// ChargeStartTime はチャージ開始時刻です。
	ChargeStartTime time.Time

	// CurrentChargeTime は現在の行動のチャージタイムです。
	CurrentChargeTime time.Duration

	// PendingAction はチャージ後に実行する行動です。
	PendingAction *EnemyAction

	// DefenseStartTime はディフェンス開始時刻です。
	DefenseStartTime time.Time

	// DefenseDuration はディフェンス持続時間です。
	DefenseDuration time.Duration

	// ActiveDefenseType は発動中のディフェンス種別です。
	ActiveDefenseType EnemyDefenseType

	// DefenseValue は軽減率/回避率です。
	DefenseValue float64

	// ========== ボルテージシステム ==========

	// Voltage は現在のボルテージ値です（100.0 = 100%）。
	// 時間経過で上昇し、プレイヤーのダメージ乗算に使用されます。
	Voltage float64
}

// NewEnemy は新しいEnemyModelを作成します。
// 初期状態は通常フェーズ（PhaseNormal）で、行動インデックスは0、ボルテージは100.0です。
func NewEnemy(id, name string, level, hp, attackPower int, enemyType EnemyType) *EnemyModel {
	return &EnemyModel{
		ID:              id,
		Name:            name,
		Level:           level,
		HP:              hp,
		MaxHP:           hp,
		AttackPower:     attackPower,
		Type:            enemyType,
		Phase:           PhaseNormal,
		EffectTable:     NewEffectTable(),
		ActionIndex:     0,     // 行動インデックス初期化
		ActivePassiveID: "",    // パッシブID初期化
		Voltage:         100.0, // ボルテージ初期化（100% = 等倍）
	}
}

// TakeDamage はダメージを受けてHPを減少させます。
// HPは0未満にはなりません。
func (e *EnemyModel) TakeDamage(damage int) {
	e.HP -= damage
	if e.HP < 0 {
		e.HP = 0
	}
}

// IsAlive は敵が生存しているかどうかを返します。
// HP > 0 の場合に生存とみなします。
func (e *EnemyModel) IsAlive() bool {
	return e.HP > 0
}

// GetHPPercentage はHPの残り割合を0.0〜1.0で返します。
func (e *EnemyModel) GetHPPercentage() float64 {
	if e.MaxHP == 0 {
		return 0.0
	}
	return float64(e.HP) / float64(e.MaxHP)
}

// ShouldTransitionToEnhanced は強化フェーズに移行すべきかどうかを判定します。
// 既にPhaseEnhancedの場合はfalseを返します。
func (e *EnemyModel) ShouldTransitionToEnhanced() bool {
	if e.Phase == PhaseEnhanced {
		return false
	}
	return e.GetHPPercentage() <= EnhanceThreshold
}

// TransitionToEnhanced は強化フェーズに移行します。
func (e *EnemyModel) TransitionToEnhanced() {
	e.Phase = PhaseEnhanced
}

// CheckAndTransitionPhase はHPをチェックし、必要に応じてフェーズ移行を実行します。
// フェーズ移行した場合はtrueを返します。
func (e *EnemyModel) CheckAndTransitionPhase() bool {
	if e.ShouldTransitionToEnhanced() {
		e.TransitionToEnhanced()
		return true
	}
	return false
}

// IsEnhanced は現在強化フェーズかどうかを返します。
func (e *EnemyModel) IsEnhanced() bool {
	return e.Phase == PhaseEnhanced
}

// GetPhaseString は現在のフェーズの表示文字列を返します。
func (e *EnemyModel) GetPhaseString() string {
	return e.Phase.String()
}

// ========== 行動管理メソッド ==========

// GetCurrentPattern は現在のフェーズに対応する行動パターンを返します。
// 強化フェーズで強化パターンが空の場合は通常パターンを継続します。
func (e *EnemyModel) GetCurrentPattern() []EnemyAction {
	if e.Phase == PhaseEnhanced && len(e.Type.ResolvedEnhancedActions) > 0 {
		return e.Type.ResolvedEnhancedActions
	}
	return e.Type.ResolvedNormalActions
}

// GetCurrentAction は現在実行すべき行動を返します。
// 行動パターンが空の場合はデフォルトの攻撃行動を返します。
func (e *EnemyModel) GetCurrentAction() EnemyAction {
	pattern := e.GetCurrentPattern()
	if len(pattern) == 0 {
		// 行動パターンが空の場合はデフォルトの攻撃行動
		return EnemyAction{
			ActionType: EnemyActionAttack,
			AttackType: e.Type.AttackType,
		}
	}
	return pattern[e.ActionIndex]
}

// AdvanceActionIndex は行動インデックスを1つ進めます（ループ対応）。
// 現在の行動パターンの長さに基づいてループします。
func (e *EnemyModel) AdvanceActionIndex() {
	pattern := e.GetCurrentPattern()
	if len(pattern) == 0 {
		return // 空パターンの場合は何もしない
	}
	e.ActionIndex = (e.ActionIndex + 1) % len(pattern)
}

// ResetActionIndex は行動インデックスを0にリセットします。
// フェーズ遷移時に使用します。
func (e *EnemyModel) ResetActionIndex() {
	e.ActionIndex = 0
}

// ========== 次回行動管理メソッド ==========

// GetNextAction は次回実行予定の行動を返します。
// PendingActionがnilの場合はnilを返します。
func (e *EnemyModel) GetNextAction() *EnemyAction {
	return e.PendingAction
}

// SetNextAction は次回実行予定の行動を設定します。
func (e *EnemyModel) SetNextAction(action *EnemyAction) {
	e.PendingAction = action
}

// PrepareNextAction は現在の行動パターンから次の行動を取得し、PendingActionに設定します。
// この関数はチャージ開始前に呼ばれることを想定しています。
func (e *EnemyModel) PrepareNextAction() {
	action := e.GetCurrentAction()
	e.PendingAction = &action
	e.CurrentChargeTime = action.ChargeTime
}

// ========== チャージ状態管理メソッド ==========

// StartCharging はチャージを開始します。
func (e *EnemyModel) StartCharging(action EnemyAction, now time.Time) {
	e.WaitMode = WaitModeCharging
	e.ChargeStartTime = now
	e.CurrentChargeTime = action.ChargeTime
	e.PendingAction = &action
}

// GetChargeProgress はチャージ進捗（0.0〜1.0）を返します。
func (e *EnemyModel) GetChargeProgress(now time.Time) float64 {
	if e.WaitMode != WaitModeCharging || e.CurrentChargeTime == 0 {
		return 0
	}
	elapsed := now.Sub(e.ChargeStartTime)
	progress := float64(elapsed) / float64(e.CurrentChargeTime)
	if progress > 1.0 {
		return 1.0
	}
	return progress
}

// GetChargeRemainingTime はチャージ残り時間を返します。
func (e *EnemyModel) GetChargeRemainingTime(now time.Time) time.Duration {
	if e.WaitMode != WaitModeCharging {
		return 0
	}
	elapsed := now.Sub(e.ChargeStartTime)
	remaining := e.CurrentChargeTime - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// IsChargeComplete はチャージ完了かどうかを返します。
func (e *EnemyModel) IsChargeComplete(now time.Time) bool {
	if e.WaitMode != WaitModeCharging {
		return false
	}
	return now.Sub(e.ChargeStartTime) >= e.CurrentChargeTime
}

// ExecuteChargedAction はチャージ完了した行動を実行可能状態にします。
// 実行後は行動インデックスを進めます。
func (e *EnemyModel) ExecuteChargedAction() *EnemyAction {
	if e.PendingAction == nil {
		return nil
	}
	action := e.PendingAction
	e.WaitMode = WaitModeNone
	e.PendingAction = nil
	e.AdvanceActionIndex()
	return action
}

// CancelCharge はチャージをキャンセルします。
func (e *EnemyModel) CancelCharge() {
	e.WaitMode = WaitModeNone
	e.PendingAction = nil
}

// ========== ディフェンス状態管理メソッド ==========

// StartDefense はディフェンスを開始します。
func (e *EnemyModel) StartDefense(defenseType EnemyDefenseType, value float64, duration time.Duration, now time.Time) {
	e.WaitMode = WaitModeDefending
	e.DefenseStartTime = now
	e.DefenseDuration = duration
	e.ActiveDefenseType = defenseType
	e.DefenseValue = value
}

// IsDefenseActive はディフェンスが有効かどうかを返します。
func (e *EnemyModel) IsDefenseActive(now time.Time) bool {
	if e.WaitMode != WaitModeDefending {
		return false
	}
	return now.Sub(e.DefenseStartTime) < e.DefenseDuration
}

// GetDefenseRemainingTime はディフェンス残り時間を返します。
func (e *EnemyModel) GetDefenseRemainingTime(now time.Time) time.Duration {
	if e.WaitMode != WaitModeDefending {
		return 0
	}
	elapsed := now.Sub(e.DefenseStartTime)
	remaining := e.DefenseDuration - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// EndDefense はディフェンスを終了し、行動インデックスを進めます。
func (e *EnemyModel) EndDefense() {
	e.WaitMode = WaitModeNone
	e.ActiveDefenseType = ""
	e.DefenseValue = 0
	e.AdvanceActionIndex()
}

// GetDefenseTypeName はディフェンス種別の表示名を返します。
func (e *EnemyModel) GetDefenseTypeName() string {
	switch e.ActiveDefenseType {
	case DefensePhysicalCut:
		return "物理防御"
	case DefenseMagicCut:
		return "魔法防御"
	case DefenseDebuffEvade:
		return "デバフ回避"
	default:
		return "防御"
	}
}

// ========== 敵行動データ構造 ==========

// EnemyActionType は敵の行動タイプを表す列挙型です。
type EnemyActionType int

const (
	// EnemyActionAttack は攻撃行動です。
	EnemyActionAttack EnemyActionType = iota

	// EnemyActionBuff は自己バフ行動です。
	EnemyActionBuff

	// EnemyActionDebuff はプレイヤーへのデバフ行動です。
	EnemyActionDebuff

	// EnemyActionDefense はディフェンス行動です。
	EnemyActionDefense
)

// String はEnemyActionTypeの日本語表示名を返します。
func (t EnemyActionType) String() string {
	switch t {
	case EnemyActionAttack:
		return "攻撃"
	case EnemyActionBuff:
		return "バフ"
	case EnemyActionDebuff:
		return "デバフ"
	case EnemyActionDefense:
		return "ディフェンス"
	default:
		return "不明"
	}
}

// EnemyDefenseType はディフェンス行動の種類を表す列挙型です。
type EnemyDefenseType string

const (
	// DefensePhysicalCut は物理ダメージ軽減です。
	DefensePhysicalCut EnemyDefenseType = "physical_cut"

	// DefenseMagicCut は魔法ダメージ軽減です。
	DefenseMagicCut EnemyDefenseType = "magic_cut"

	// DefenseDebuffEvade はデバフ回避です。
	DefenseDebuffEvade EnemyDefenseType = "debuff_evade"
)

// EnemyAction は敵の個別行動を定義する値オブジェクトです。
type EnemyAction struct {
	// ========== 共通フィールド ==========

	// ID は行動の一意識別子です。
	ID string

	// Name は行動の表示名です。
	Name string

	// ActionType は行動タイプ（攻撃、バフ、デバフ、ディフェンス）です。
	ActionType EnemyActionType

	// ChargeTime はチャージタイム（行動決定から実行までの時間）です。
	ChargeTime time.Duration

	// ========== 攻撃行動用フィールド ==========

	// AttackType は攻撃行動時の攻撃属性（"physical" または "magic"）です。
	AttackType string

	// DamageBase は基礎ダメージ (a) です。ダメージ = a + Lv * b
	DamageBase float64

	// DamagePerLevel はレベル係数 (b) です。ダメージ = a + Lv * b
	DamagePerLevel float64

	// Element は攻撃の属性（"fire", "water", "dark"等）です。空文字は無属性。
	Element string

	// ========== バフ/デバフ行動用フィールド ==========

	// EffectType はバフ/デバフ行動時の効果種別です（例: "damage_mult", "cooldown_reduce"）。
	EffectType string

	// EffectValue はバフ/デバフ行動時の効果値です。
	EffectValue float64

	// Duration はバフ/デバフ/ディフェンスの持続時間（秒）です。
	Duration float64

	// ========== ディフェンス行動用フィールド ==========

	// DefenseType はディフェンスの種類（物理軽減/魔法軽減/デバフ回避）です。
	DefenseType EnemyDefenseType

	// ReductionRate は物理/魔法ダメージの軽減割合（0.0〜1.0）です。
	ReductionRate float64

	// EvadeRate はデバフ回避率（0.0〜1.0）です。
	EvadeRate float64
}

// IsAttack は攻撃行動かどうかを判定します。
func (a EnemyAction) IsAttack() bool {
	return a.ActionType == EnemyActionAttack
}

// IsBuff は自己バフ行動かどうかを判定します。
func (a EnemyAction) IsBuff() bool {
	return a.ActionType == EnemyActionBuff
}

// IsDebuff はデバフ行動かどうかを判定します。
func (a EnemyAction) IsDebuff() bool {
	return a.ActionType == EnemyActionDebuff
}

// IsDefense はディフェンス行動かどうかを判定します。
func (a EnemyAction) IsDefense() bool {
	return a.ActionType == EnemyActionDefense
}

// CalculateDamage はレベルに応じたダメージを計算します。
// ダメージ = DamageBase + Level * DamagePerLevel
func (a EnemyAction) CalculateDamage(level int) int {
	damage := a.DamageBase + float64(level)*a.DamagePerLevel
	if damage < 1 {
		return 1
	}
	return int(damage)
}

// GetChargeTimeMs はチャージタイムをミリ秒で返します。
func (a EnemyAction) GetChargeTimeMs() int64 {
	return a.ChargeTime.Milliseconds()
}

// ========== 敵パッシブスキルデータ構造 ==========

// EnemyPassiveSkill は敵用パッシブスキルを定義する構造体です。
// 敵の状態（通常/強化）に紐づき、EffectTableを通じて効果を適用します。
type EnemyPassiveSkill struct {
	// ID はパッシブスキルの一意識別子です。
	ID string

	// Name はパッシブスキルの表示名です。
	Name string

	// Description はパッシブスキルの説明文です。
	Description string

	// Effects は効果値のマップです（EffectColumnをキー、float64を値）。
	Effects map[EffectColumn]float64
}

// ToEntry はEnemyPassiveSkillをEffectTableに登録可能なEffectEntryに変換します。
// パッシブスキルは永続効果（Duration=nil）として登録されます。
func (p *EnemyPassiveSkill) ToEntry() EffectEntry {
	// 効果値をコピーして新しいマップを作成
	values := make(map[EffectColumn]float64)
	for k, v := range p.Effects {
		values[k] = v
	}

	return EffectEntry{
		SourceType: SourcePassive,
		SourceID:   p.ID,
		Name:       p.Name,
		Duration:   nil, // 永続効果
		Values:     values,
	}
}

// ========== ボルテージ関連メソッド ==========

// GetVoltage は現在のボルテージ値を返します。
func (e *EnemyModel) GetVoltage() float64 {
	return e.Voltage
}

// SetVoltage はボルテージ値を設定します。
func (e *EnemyModel) SetVoltage(voltage float64) {
	e.Voltage = voltage
}

// GetVoltageMultiplier はダメージ乗算用の倍率を返します（ボルテージ/100）。
// 例: ボルテージ100.0 -> 1.0倍、150.0 -> 1.5倍
func (e *EnemyModel) GetVoltageMultiplier() float64 {
	return e.Voltage / 100.0
}
