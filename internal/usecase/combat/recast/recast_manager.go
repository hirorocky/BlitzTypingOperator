// Package recast はリキャスト管理機能を提供します。
// スキル使用後のエージェント全体のクールダウンを管理します。
package recast

import (
	"sort"
	"time"
)

// RecastState はエージェントのリキャスト状態を表す構造体です。
type RecastState struct {
	// AgentIndex はエージェントのインデックス（0-2）です。
	AgentIndex int

	// RemainingSeconds はリキャスト残り秒数です。
	RemainingSeconds float64

	// TotalSeconds はリキャスト総秒数です。
	TotalSeconds float64
}

// RecastManager はエージェントのリキャスト状態を管理する構造体です。
// エージェントがスキルを使用すると、そのエージェントの全スキルが
// 一定時間使用不可能になります。
type RecastManager struct {
	// states はエージェントインデックスごとのリキャスト状態マップです。
	states map[int]*RecastState
}

// NewRecastManager は新しいRecastManagerを作成します。
func NewRecastManager() *RecastManager {
	return &RecastManager{
		states: make(map[int]*RecastState),
	}
}

// StartRecast はエージェントのリキャストを開始します。
// 既にリキャスト中の場合は上書きされます。
// durationが0以下の場合はリキャストを開始しません。
func (m *RecastManager) StartRecast(agentIndex int, duration time.Duration) {
	if duration <= 0 {
		return
	}

	seconds := duration.Seconds()
	m.states[agentIndex] = &RecastState{
		AgentIndex:       agentIndex,
		RemainingSeconds: seconds,
		TotalSeconds:     seconds,
	}
}

// UpdateRecast はリキャスト時間を更新し、完了したエージェントのインデックスを返します。
// deltaはフレーム間の経過時間です。
func (m *RecastManager) UpdateRecast(delta time.Duration) []int {
	completed := make([]int, 0)
	deltaSeconds := delta.Seconds()

	for agentIndex, state := range m.states {
		state.RemainingSeconds -= deltaSeconds

		if state.RemainingSeconds <= 0 {
			completed = append(completed, agentIndex)
		}
	}

	// 完了したエージェントの状態を削除
	for _, agentIndex := range completed {
		delete(m.states, agentIndex)
	}

	// インデックス順にソート
	sort.Ints(completed)

	return completed
}

// IsAgentReady はエージェントが使用可能かを判定します。
// リキャスト中でなければ使用可能です。
func (m *RecastManager) IsAgentReady(agentIndex int) bool {
	_, exists := m.states[agentIndex]
	return !exists
}

// GetRecastState は指定エージェントのリキャスト状態を取得します。
// リキャスト中でなければnilを返します。
func (m *RecastManager) GetRecastState(agentIndex int) *RecastState {
	return m.states[agentIndex]
}

// GetAllRecastStates は全てのリキャスト状態をエージェントインデックス順で返します。
func (m *RecastManager) GetAllRecastStates() []*RecastState {
	states := make([]*RecastState, 0, len(m.states))
	for _, state := range m.states {
		states = append(states, state)
	}

	// インデックス順にソート
	sort.Slice(states, func(i, j int) bool {
		return states[i].AgentIndex < states[j].AgentIndex
	})

	return states
}

// GetProgress はエージェントのリキャスト進捗（0.0-1.0）を返します。
// 0.0は開始直後、1.0は完了または未リキャストを示します。
func (m *RecastManager) GetProgress(agentIndex int) float64 {
	state := m.states[agentIndex]
	if state == nil {
		return 1.0 // リキャストしていない = 完了
	}

	if state.TotalSeconds <= 0 {
		return 1.0
	}

	elapsed := state.TotalSeconds - state.RemainingSeconds
	return elapsed / state.TotalSeconds
}

// CancelRecast は指定エージェントのリキャストをキャンセルします。
func (m *RecastManager) CancelRecast(agentIndex int) {
	delete(m.states, agentIndex)
}

// CancelAllRecasts は全エージェントのリキャストをキャンセルします。
func (m *RecastManager) CancelAllRecasts() {
	m.states = make(map[int]*RecastState)
}

// ReduceAllRecasts は全エージェントのリキャスト残り時間を短縮します。
// reduction が残り時間より大きい場合でも0未満にはなりません。
func (m *RecastManager) ReduceAllRecasts(reduction time.Duration) {
	seconds := reduction.Seconds()
	for _, state := range m.states {
		state.RemainingSeconds -= seconds
		if state.RemainingSeconds < 0 {
			state.RemainingSeconds = 0
		}
	}
}
