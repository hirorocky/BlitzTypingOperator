// Package domain はゲームのドメインモデルを定義します。
package domain

// HP計算に使用する定数
// 新仕様: 初期最大HP = 1000、敵撃破による成長で増加
// 旧仕様（互換性用）: MaxHP = 装備中エージェントのコアレベル平均 × HP係数 + 基礎HP

const (
	HPCoefficient = 10.0 // 旧仕様：レベル平均に掛ける係数
	BaseHP        = 100  // 旧仕様：基礎HP値
	InitialMaxHP  = 1000 // 新仕様：初期最大HP
)

// PlayerModel はゲーム内のプレイヤーエンティティを表す構造体です。
// プレイヤーはHP（敵の攻撃対象）とバフ・デバフ状態（一時的なステータス効果）を持ちます。

type PlayerModel struct {
	// HP はプレイヤーの現在HP値です。

	HP int

	// MaxHP はプレイヤーの最大HP値です。

	MaxHP int

	// TempHP は一時HPです（オーバーヒール等で付与）。
	// ダメージを受けるとTempHPから先に消費されます。
	TempHP int

	// EffectTable はプレイヤーに適用されているステータス効果テーブルです。
	// バフ/デバフ/コア特性/モジュールパッシブなどの効果を集約します。

	EffectTable *EffectTable
}

// NewPlayer は新しいPlayerModelを作成します。
// 初期状態ではHP/MaxHPは0で、エージェント装備後にRecalculateHPで計算されます。
func NewPlayer() *PlayerModel {
	return &PlayerModel{
		HP:          0,
		MaxHP:       0,
		EffectTable: NewEffectTable(),
	}
}

// CalculateMaxHP は最大HPを計算します。
// 新仕様では初期最大HP（1000）を返し、成長はIncreaseMaxHPで管理します。
// 注意: この関数はレガシーコードとの互換性のために残されています。
func CalculateMaxHP(agents []*AgentModel) int {
	if len(agents) == 0 {
		return BaseHP
	}
	// 新仕様: エージェント装備時は初期最大HPを返す
	return InitialMaxHP
}

// RecalculateHP は装備エージェントに基づいてMaxHPを再計算し、HPを全回復します。

func (p *PlayerModel) RecalculateHP(agents []*AgentModel) {
	p.MaxHP = CalculateMaxHP(agents)
	p.HP = p.MaxHP
}

// FullHeal はHPを最大値まで回復します。
func (p *PlayerModel) FullHeal() {
	p.HP = p.MaxHP
}

// TakeDamage はダメージを受けてHPを減少させます。
// TempHPがある場合は先に消費されます。HPは0未満にはなりません。
func (p *PlayerModel) TakeDamage(damage int) {
	// TempHPから先に消費
	if p.TempHP > 0 {
		if damage <= p.TempHP {
			p.TempHP -= damage
			return
		}
		damage -= p.TempHP
		p.TempHP = 0
	}

	p.HP -= damage
	if p.HP < 0 {
		p.HP = 0
	}
}

// Heal はHPを回復します。
// HPはMaxHPを超えません。
func (p *PlayerModel) Heal(amount int) {
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
}

// HealWithOverheal はHPを回復し、オーバーヒール分をTempHPに変換します。
// TempHPの上限はMaxHPの50%です。
func (p *PlayerModel) HealWithOverheal(amount int) int {
	// まず通常回復
	hpBefore := p.HP
	p.HP += amount
	overflow := 0

	if p.HP > p.MaxHP {
		overflow = p.HP - p.MaxHP
		p.HP = p.MaxHP
	}

	// 超過分をTempHPに変換（上限はMaxHPの50%）
	tempHPCap := p.MaxHP / 2
	if overflow > 0 {
		p.TempHP += overflow
		if p.TempHP > tempHPCap {
			p.TempHP = tempHPCap
		}
	}

	// 実際に回復した量（通常回復+TempHP）を返す
	healed := p.HP - hpBefore + overflow
	return healed
}

// IsAlive はプレイヤーが生存しているかどうかを返します。
// HP > 0 の場合に生存とみなします。
func (p *PlayerModel) IsAlive() bool {
	return p.HP > 0
}

// PrepareForBattle はバトル開始時の準備を行います。

func (p *PlayerModel) PrepareForBattle() {
	p.FullHeal()
	// EffectTableもリセット（バトル間で効果を持ち越さない）
	p.EffectTable = NewEffectTable()
}

// GetHPPercentage はHPの残り割合を0.0〜1.0で返します。
func (p *PlayerModel) GetHPPercentage() float64 {
	if p.MaxHP == 0 {
		return 0.0
	}
	return float64(p.HP) / float64(p.MaxHP)
}
