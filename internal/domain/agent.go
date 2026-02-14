// Package domain はゲームのドメインモデルを定義します。
package domain

// AgentModel はゲーム内のエージェントエンティティを表す構造体です。
// エージェントは1つのコアと1〜4つのスキルで構成され、バトル中にプレイヤーを支援します。

type AgentModel struct {
	// ID はエージェントインスタンスの一意識別子です。
	ID string

	// Core はエージェントの核となるコアです。
	// エージェントのステータスはこのコアから導出されます。
	Core *CoreModel

	// Skills はエージェントに装備されているスキル（スキル）のリストです。
	// エージェントは1〜4つのスキルを装備できます。
	Skills []*SkillModel

	// BaseStats はエージェントの基礎ステータス値です。
	// コアのステータスから導出され、スキル効果計算の基準となります。
	// バフ/デバフ等の効果はEffectTableを通じて適用されます。
	BaseStats Stats
}

// NewAgent は新しいAgentModelを作成します。
// 基礎ステータスはコアのステータスからコピーされます。
// skillsはコピーされ、元のスライスとの参照共有を避けます。
func NewAgent(id string, core *CoreModel, skills []*SkillModel) *AgentModel {
	// スキルリストをコピー（スライスの参照共有を避ける）
	skillsCopy := make([]*SkillModel, len(skills))
	copy(skillsCopy, skills)

	return &AgentModel{
		ID:        id,
		Core:      core,
		Skills:    skillsCopy,
		BaseStats: core.Stats, // 基礎ステータスはコアから導出
	}
}

// GetCoreTypeName はコア特性の名前を返します。
func (a *AgentModel) GetCoreTypeName() string {
	if a.Core == nil {
		return ""
	}
	return a.Core.Type.Name
}
