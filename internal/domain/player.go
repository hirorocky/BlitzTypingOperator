// Package domain はゲームのドメインモデルを定義します。
package domain

const (
	// InitialMaxHP はプレイヤーの初期最大HPです。敵撃破による成長で増加します。
	InitialMaxHP = 1000
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

	// Mana はプレイヤーの現在マナ値です。バトル中のリソースで、スキル使用に消費されます。
	Mana int

	// MaxMana はプレイヤーの最大マナ値です。
	MaxMana int

	// EffectTable はプレイヤーに適用されているステータス効果テーブルです。
	// バフ/デバフ/コア特性/モジュールパッシブなどの効果を集約します。

	EffectTable *EffectTable
}

// NewPlayerWithMaxHP は指定された最大HPでPlayerModelを作成します。
// 新規ゲーム開始時はInitialMaxHP（1000）を渡してください。
// セーブデータからの復元時は保存されたMaxHPを渡してください。
func NewPlayerWithMaxHP(maxHP int) *PlayerModel {
	return &PlayerModel{
		HP:          maxHP,
		MaxHP:       maxHP,
		EffectTable: NewEffectTable(),
	}
}

// IncreaseMaxHP は最大HPを増加させます。
// 敵撃破による成長時に使用します。
// 負の値やゼロは無視されます（MaxHPは減少しません）。
func (p *PlayerModel) IncreaseMaxHP(amount int) {
	if amount <= 0 {
		return
	}
	p.MaxHP += amount
}

// InitializeHP はプレイヤーのMaxHPを設定し、HPを全回復します。
// エージェント装備後の初期化に使用します。
func (p *PlayerModel) InitializeHP(maxHP int) {
	p.MaxHP = maxHP
	p.HP = maxHP
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

// ConsumeMana はマナを消費します。
// マナが不足している場合はfalseを返し、マナは変更されません。
// 負の値は無視されます（falseを返す）。0の場合は何も消費せず成功を返します。
func (p *PlayerModel) ConsumeMana(amount int) bool {
	if amount < 0 {
		return false
	}
	if amount > p.Mana {
		return false
	}
	p.Mana -= amount
	return true
}

// GainMana はマナを獲得します。MaxManaを超えないようクランプされます。
// 負の値は無視されます。
func (p *PlayerModel) GainMana(amount int) {
	if amount <= 0 {
		return
	}
	p.Mana += amount
	if p.Mana > p.MaxMana {
		p.Mana = p.MaxMana
	}
}

// PrepareForBattle はバトル開始時の準備を行います。

func (p *PlayerModel) PrepareForBattle() {
	p.FullHeal()
	p.Mana = 0
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
