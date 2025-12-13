// Package domain はゲームのドメインモデルを定義します。
package domain

// HP計算に使用する定数
// MaxHP = 装備中エージェントのコアレベル平均 × HP係数 + 基礎HP

const (
	HPCoefficient = 10.0 // レベル平均に掛ける係数
	BaseHP        = 100  // 基礎HP値
)

// PlayerModel はゲーム内のプレイヤーエンティティを表す構造体です。
// プレイヤーはHP（敵の攻撃対象）とバフ・デバフ状態（一時的なステータス効果）を持ちます。

type PlayerModel struct {
	// HP はプレイヤーの現在HP値です。

	HP int

	// MaxHP はプレイヤーの最大HP値です。

	MaxHP int

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

// CalculateMaxHP は装備中エージェントのコアレベル平均からMaxHPを計算します。

// エージェントが装備されていない場合は基礎HPを返します。
func CalculateMaxHP(agents []*AgentModel) int {
	if len(agents) == 0 {
		return BaseHP
	}

	totalLevel := 0
	for _, agent := range agents {
		totalLevel += agent.Level
	}

	// 平均レベル × HP係数 + 基礎HP
	avgLevel := float64(totalLevel) / float64(len(agents))
	return int(avgLevel*HPCoefficient) + BaseHP
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
// HPは0未満にはなりません。
func (p *PlayerModel) TakeDamage(damage int) {
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
