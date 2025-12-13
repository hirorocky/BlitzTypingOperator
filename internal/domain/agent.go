// Package domain はゲームのドメインモデルを定義します。
package domain

// AgentModel はゲーム内のエージェントエンティティを表す構造体です。
// エージェントは1つのコアと4つのモジュールで構成され、バトル中にプレイヤーを支援します。

type AgentModel struct {
	// ID はエージェントインスタンスの一意識別子です。
	ID string

	// Core はエージェントの核となるコアです。
	// エージェントのレベルとステータスはこのコアから導出されます。
	Core *CoreModel

	// Modules はエージェントに装備されているモジュール（スキル）のリストです。
	// エージェントは必ず4つのモジュールを装備します。
	Modules []*ModuleModel

	// Level はエージェントのレベルです。

	// エージェント自体の成長/レベリングはありません。
	Level int

	// BaseStats はエージェントの基礎ステータス値です。
	// コアのステータスから導出され、モジュール効果計算の基準となります。
	// バフ/デバフ等の効果はEffectTableを通じて適用されます。
	BaseStats Stats
}

// ModuleSlotCount はエージェント1体あたりのモジュールスロット数です。
const ModuleSlotCount = 4

// NewAgent は新しいAgentModelを作成します。
// エージェントのレベルはコアのレベルから自動的に導出されます。
// 基礎ステータスはコアのステータスからコピーされます。
// modulesはコピーされ、元のスライスとの参照共有を避けます。
func NewAgent(id string, core *CoreModel, modules []*ModuleModel) *AgentModel {
	// モジュールリストをコピー（スライスの参照共有を避ける）
	modulesCopy := make([]*ModuleModel, len(modules))
	copy(modulesCopy, modules)

	return &AgentModel{
		ID:        id,
		Core:      core,
		Modules:   modulesCopy,
		Level:     core.Level,
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
