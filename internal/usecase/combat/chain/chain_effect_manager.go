// Package chain はチェイン効果管理機能を提供します。
// モジュール使用後、次のスキル使用まで待機し、他エージェントの行動で発動する追加効果を管理します。
package chain

import (
	"sort"

	"hirorocky/type-battle/internal/domain"
)

// SkillEffectFlags はモジュールが持つ効果の種別を表します。
type SkillEffectFlags struct {
	// HasDamage はダメージ効果を持つかを表します。
	HasDamage bool

	// HasHeal は回復効果を持つかを表します。
	HasHeal bool

	// HasBuff はバフ効果を持つかを表します。
	HasBuff bool

	// HasDebuff はデバフ効果を持つかを表します。
	HasDebuff bool
}

// PendingChainEffect は待機中のチェイン効果を表す構造体です。
type PendingChainEffect struct {
	// AgentIndex はこの効果を登録したエージェントのインデックスです。
	AgentIndex int

	// Effect はチェイン効果の内容です。
	Effect domain.ChainEffect

	// SourceModule はこの効果を発生させたモジュールのIDです。
	SourceModule string
}

// TriggeredChainEffect は発動したチェイン効果を表す構造体です。
type TriggeredChainEffect struct {
	// Effect はチェイン効果の内容です。
	Effect domain.ChainEffect

	// EffectValue は効果値です。
	EffectValue float64

	// Message は発動メッセージです。
	Message string

	// SourceAgentIndex は効果を登録したエージェントのインデックスです。
	SourceAgentIndex int
}

// ChainEffectManager はチェイン効果の管理を担当する構造体です。
// モジュール使用時にチェイン効果を待機状態として登録し、
// 他エージェントのモジュール使用時に発動条件をチェックして発動します。
type ChainEffectManager struct {
	// pendingEffects はエージェントインデックスごとの待機中チェイン効果です。
	// 1エージェントにつき1つの待機中効果のみ保持します。
	pendingEffects map[int]*PendingChainEffect
}

// NewChainEffectManager は新しいChainEffectManagerを作成します。
func NewChainEffectManager() *ChainEffectManager {
	return &ChainEffectManager{
		pendingEffects: make(map[int]*PendingChainEffect),
	}
}

// RegisterChainEffect はモジュール使用時にチェイン効果を待機状態として登録します。
// 既に同一エージェントの効果がある場合は上書きします。
// effectがnilの場合は何もしません。
func (m *ChainEffectManager) RegisterChainEffect(agentIndex int, effect *domain.ChainEffect, sourceModule string) {
	if effect == nil {
		return
	}

	m.pendingEffects[agentIndex] = &PendingChainEffect{
		AgentIndex:   agentIndex,
		Effect:       *effect,
		SourceModule: sourceModule,
	}
}

// CheckAndTrigger は他エージェントのモジュール使用時に発動条件をチェックし、発動を実行します。
// usingAgentIndexはモジュールを使用したエージェントのインデックスです。
// effectFlagsは使用したモジュールが持つ効果の種別です。
// 発動した効果のリストを返します。
func (m *ChainEffectManager) CheckAndTrigger(usingAgentIndex int, effectFlags SkillEffectFlags) []TriggeredChainEffect {
	triggered := make([]TriggeredChainEffect, 0)
	expiredAgents := make([]int, 0)

	for agentIndex, pending := range m.pendingEffects {
		// 同一エージェントのチェイン効果は発動しない
		if agentIndex == usingAgentIndex {
			continue
		}

		// 効果種別マッチングをチェック
		if !m.isEffectTriggeredBy(pending.Effect.Type, effectFlags) {
			continue
		}

		// 発動
		triggered = append(triggered, TriggeredChainEffect{
			Effect:           pending.Effect,
			EffectValue:      pending.Effect.Value,
			Message:          pending.Effect.Description,
			SourceAgentIndex: agentIndex,
		})

		// 発動後は削除対象
		expiredAgents = append(expiredAgents, agentIndex)
	}

	// 発動した効果を削除
	for _, agentIndex := range expiredAgents {
		delete(m.pendingEffects, agentIndex)
	}

	return triggered
}

// isEffectTriggeredBy はチェイン効果がモジュールの効果種別によって発動するかを判定します。
func (m *ChainEffectManager) isEffectTriggeredBy(effectType domain.ChainEffectType, flags SkillEffectFlags) bool {
	effectCategory := effectType.Category()

	switch effectCategory {
	case domain.ChainEffectCategoryAttack:
		// 攻撃強化効果はダメージモジュールで発動
		return flags.HasDamage

	case domain.ChainEffectCategoryHeal:
		// 回復強化効果は回復モジュールで発動
		return flags.HasHeal

	case domain.ChainEffectCategoryDefense:
		// 防御強化効果は任意のモジュールで発動
		return true

	case domain.ChainEffectCategoryTyping:
		// タイピング強化効果は任意のモジュールで発動
		return true

	case domain.ChainEffectCategoryRecast:
		// リキャスト強化効果は任意のモジュールで発動
		return true

	case domain.ChainEffectCategoryEffectExtend:
		// 効果延長カテゴリは効果種別で判定
		switch effectType {
		case domain.ChainEffectBuffExtend, domain.ChainEffectBuffDuration:
			return flags.HasBuff
		case domain.ChainEffectDebuffExtend, domain.ChainEffectDebuffDuration:
			return flags.HasDebuff
		}
		return false

	case domain.ChainEffectCategorySpecial:
		// 特殊効果は任意のモジュールで発動
		return true

	default:
		return false
	}
}

// ClearEffectForAgent は指定エージェントの待機中チェイン効果を削除します。
// チェイン効果のないスキル使用時に、以前のチェイン効果をクリアするために使用します。
func (m *ChainEffectManager) ClearEffectForAgent(agentIndex int) {
	delete(m.pendingEffects, agentIndex)
}

// GetPendingEffects は待機中チェイン効果を取得します。
// エージェントインデックス順にソートされて返されます。
func (m *ChainEffectManager) GetPendingEffects() []*PendingChainEffect {
	effects := make([]*PendingChainEffect, 0, len(m.pendingEffects))
	for _, pending := range m.pendingEffects {
		effects = append(effects, pending)
	}

	// インデックス順にソート
	sort.Slice(effects, func(i, j int) bool {
		return effects[i].AgentIndex < effects[j].AgentIndex
	})

	return effects
}

// HasPendingEffect は指定エージェントに待機中チェイン効果があるかを返します。
func (m *ChainEffectManager) HasPendingEffect(agentIndex int) bool {
	_, exists := m.pendingEffects[agentIndex]
	return exists
}

// GetPendingEffectForAgent は指定エージェントの待機中チェイン効果を取得します。
// 存在しない場合はnilを返します。
func (m *ChainEffectManager) GetPendingEffectForAgent(agentIndex int) *PendingChainEffect {
	return m.pendingEffects[agentIndex]
}

// ClearAll は全ての待機中チェイン効果をクリアします。
func (m *ChainEffectManager) ClearAll() {
	m.pendingEffects = make(map[int]*PendingChainEffect)
}
