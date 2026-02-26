package domain

import "fmt"

// FeatureID はゲーム機能の識別子です。
type FeatureID string

const (
	// ディフェンススキル（ランク2で解放）
	FeatureDefenseSkill FeatureID = "defense_skill"
	// エージェントカスタマイズ（ランク3で解放）
	FeatureAgentCustomization FeatureID = "agent_customization"
	// チェイン効果（ランク4で解放）
	FeatureChainEffect FeatureID = "chain_effect"
	// マナシステム（ランク5で解放）
	FeatureManaSystem FeatureID = "mana_system"
	// 潜在効果（ランク6で解放）
	FeatureLatentEffect FeatureID = "latent_effect"
)

// FeatureStatus は機能の解放状態を表します。
// Locked → PendingTutorial → Unlocked の単調増加のみ許可されます。
type FeatureStatus uint8

const (
	// 未解放
	FeatureLocked FeatureStatus = iota
	// チュートリアル未完了（解放処理済みだがまだ有効化されていない）
	FeaturePendingTutorial
	// 解放済み（チュートリアル完了後に有効化）
	FeatureUnlocked
)

// UnlockRule はランクに応じた機能解放ルールを定義します。
type UnlockRule struct {
	FeatureID  FeatureID
	UnlockRank int
	TutorialID string
}

// TutorialDef はチュートリアルの定義です。
type TutorialDef struct {
	ID             string
	Title          string
	Pages          []string
	DefaultVisible bool
	FeatureID      FeatureID
}

// FeatureUnlockState は全機能の解放状態を管理します。
type FeatureUnlockState struct {
	features map[FeatureID]FeatureStatus
}

// NewFeatureUnlockState は全機能がLockedの初期状態を作成します。
func NewFeatureUnlockState() FeatureUnlockState {
	return FeatureUnlockState{
		features: make(map[FeatureID]FeatureStatus),
	}
}

// NewFeatureUnlockStateFrom は指定された状態マップから FeatureUnlockState を作成します。
func NewFeatureUnlockStateFrom(features map[FeatureID]FeatureStatus) FeatureUnlockState {
	copied := make(map[FeatureID]FeatureStatus, len(features))
	for k, v := range features {
		copied[k] = v
	}
	return FeatureUnlockState{features: copied}
}

// GetStatus は指定機能の現在の状態を返します。
// 未登録の機能はLockedとして扱います。
func (s FeatureUnlockState) GetStatus(id FeatureID) FeatureStatus {
	status, ok := s.features[id]
	if !ok {
		return FeatureLocked
	}
	return status
}

// CanTransition は指定された状態遷移が許可されるかを判定します。
// 単調増加（Locked < PendingTutorial < Unlocked）のみ許可されます。
func CanTransition(from, to FeatureStatus) bool {
	return to > from
}

// TransitionTo は指定機能の状態を遷移させます。
// 単調増加に違反する場合はエラーを返します。
func (s *FeatureUnlockState) TransitionTo(id FeatureID, to FeatureStatus) error {
	current := s.GetStatus(id)
	if !CanTransition(current, to) {
		return fmt.Errorf("不正な状態遷移: %s の %d → %d", id, current, to)
	}
	s.features[id] = to
	return nil
}

// Clone はディープコピーを返します。
func (s FeatureUnlockState) Clone() FeatureUnlockState {
	return NewFeatureUnlockStateFrom(s.features)
}

// AllFeatures は内部の状態マップのコピーを返します。
func (s FeatureUnlockState) AllFeatures() map[FeatureID]FeatureStatus {
	copied := make(map[FeatureID]FeatureStatus, len(s.features))
	for k, v := range s.features {
		copied[k] = v
	}
	return copied
}
