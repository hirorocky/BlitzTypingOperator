package unlocking

import (
	"fmt"
	"sort"

	"hirorocky/type-battle/internal/domain"
)

// UnlockDelta はApplyRank/Reconcileの結果を返す構造体です。
type UnlockDelta struct {
	// 今回新たにPendingTutorialになった機能IDリスト（ランク昇順）
	NewPendingFeatures []domain.FeatureID
	// チュートリアルキューに追加されたチュートリアルIDリスト
	QueuedTutorials []string
}

// Manager は機能解放状態を管理します。
type Manager struct {
	rules             []domain.UnlockRule
	tutorials         []domain.TutorialDef
	tutorialByID      map[string]domain.TutorialDef
	tutorialByFeature map[domain.FeatureID]string
	state             domain.FeatureUnlockState
	lastAppliedRank   int
	pendingQueue      []string // チュートリアルIDのキュー（ランク昇順）
}

// NewManager は新しいManagerを作成します。
// rulesはランク昇順にソートされます。
func NewManager(rules []domain.UnlockRule, tutorials []domain.TutorialDef, state domain.FeatureUnlockState) (*Manager, error) {
	// ルールをランク昇順にソート
	sortedRules := make([]domain.UnlockRule, len(rules))
	copy(sortedRules, rules)
	sort.Slice(sortedRules, func(i, j int) bool {
		return sortedRules[i].UnlockRank < sortedRules[j].UnlockRank
	})

	// チュートリアルのインデックスを構築
	tutorialByID := make(map[string]domain.TutorialDef, len(tutorials))
	tutorialByFeature := make(map[domain.FeatureID]string)
	for _, t := range tutorials {
		tutorialByID[t.ID] = t
		if t.FeatureID != "" {
			tutorialByFeature[t.FeatureID] = t.ID
		}
	}

	// 既存のPendingTutorialからキューを再構築
	var pendingQueue []string
	for _, r := range sortedRules {
		if state.GetStatus(r.FeatureID) == domain.FeaturePendingTutorial {
			pendingQueue = append(pendingQueue, r.TutorialID)
		}
	}

	return &Manager{
		rules:             sortedRules,
		tutorials:         tutorials,
		tutorialByID:      tutorialByID,
		tutorialByFeature: tutorialByFeature,
		state:             state.Clone(),
		lastAppliedRank:   0,
		pendingQueue:      pendingQueue,
	}, nil
}

// ApplyRank は指定ランクまでの未解放機能をPendingTutorialに遷移させます。
// 冪等: 同一ランクで再度呼んでも二重解放されません。
func (m *Manager) ApplyRank(rank int) (UnlockDelta, error) {
	delta := UnlockDelta{}

	for _, rule := range m.rules {
		if rule.UnlockRank > rank {
			continue
		}
		current := m.state.GetStatus(rule.FeatureID)
		if current != domain.FeatureLocked {
			continue
		}
		if err := m.state.TransitionTo(rule.FeatureID, domain.FeaturePendingTutorial); err != nil {
			return delta, fmt.Errorf("ApplyRank: %w", err)
		}
		delta.NewPendingFeatures = append(delta.NewPendingFeatures, rule.FeatureID)
		delta.QueuedTutorials = append(delta.QueuedTutorials, rule.TutorialID)
		m.pendingQueue = append(m.pendingQueue, rule.TutorialID)
	}

	if rank > m.lastAppliedRank {
		m.lastAppliedRank = rank
	}

	return delta, nil
}

// Reconcile は現在ランクに対して未適用の機能があればPendingTutorial化します。
// マスタデータ追加時の既存ユーザー対応に使用します。
func (m *Manager) Reconcile(currentRank int) (UnlockDelta, error) {
	return m.ApplyRank(currentRank)
}

// NextPendingTutorial はキューの先頭のチュートリアルIDを返します。
// キューが空の場合はfalseを返します。
func (m *Manager) NextPendingTutorial() (string, bool) {
	if len(m.pendingQueue) == 0 {
		return "", false
	}
	return m.pendingQueue[0], true
}

// CompleteTutorial は指定チュートリアルを完了し、対応機能をUnlockedに遷移させます。
// 対象機能がPendingTutorialでない場合はエラーを返します。
func (m *Manager) CompleteTutorial(tutorialID string) (domain.FeatureID, error) {
	tut, ok := m.tutorialByID[tutorialID]
	if !ok {
		return "", fmt.Errorf("チュートリアルが見つかりません: %s", tutorialID)
	}

	featureID := tut.FeatureID
	if featureID == "" {
		return "", fmt.Errorf("チュートリアル %s に対応する機能がありません", tutorialID)
	}

	current := m.state.GetStatus(featureID)
	if current != domain.FeaturePendingTutorial {
		return "", fmt.Errorf("機能 %s はPendingTutorial状態ではありません（現在: %d）", featureID, current)
	}

	if err := m.state.TransitionTo(featureID, domain.FeatureUnlocked); err != nil {
		return "", fmt.Errorf("CompleteTutorial: %w", err)
	}

	// キューから削除
	m.removeFromQueue(tutorialID)

	return featureID, nil
}

// IsUnlocked は指定機能がUnlocked状態かどうかを返します。
func (m *Manager) IsUnlocked(id domain.FeatureID) bool {
	return m.state.GetStatus(id) == domain.FeatureUnlocked
}

// IsPendingOrUnlocked は指定機能がPendingTutorialまたはUnlocked状態かどうかを返します。
func (m *Manager) IsPendingOrUnlocked(id domain.FeatureID) bool {
	status := m.state.GetStatus(id)
	return status == domain.FeaturePendingTutorial || status == domain.FeatureUnlocked
}

// ListVisibleTutorials はTIPS画面に表示可能なチュートリアルを返します。
// default_visibleなもの + Unlocked機能のチュートリアルが対象です。
// 返却値はディープコピーされるため、呼び出し側で安全に使用できます。
func (m *Manager) ListVisibleTutorials() []domain.TutorialDef {
	var result []domain.TutorialDef
	for _, t := range m.tutorials {
		if t.DefaultVisible || (t.FeatureID != "" && m.IsUnlocked(t.FeatureID)) {
			copied := t
			copied.Pages = make([]string, len(t.Pages))
			copy(copied.Pages, t.Pages)
			result = append(result, copied)
		}
	}
	return result
}

// Snapshot は現在の状態のディープコピーを返します。
func (m *Manager) Snapshot() domain.FeatureUnlockState {
	return m.state.Clone()
}

// removeFromQueue はキューから指定チュートリアルIDを削除します。
func (m *Manager) removeFromQueue(tutorialID string) {
	for i, id := range m.pendingQueue {
		if id == tutorialID {
			m.pendingQueue = append(m.pendingQueue[:i], m.pendingQueue[i+1:]...)
			return
		}
	}
}
