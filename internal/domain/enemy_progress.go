// Package domain はゲームのドメインモデルを定義します。
package domain

// HP成長に関する定数
const (
	// HPGainPerFirstDefeat は初撃破時のHP増加量です。
	HPGainPerFirstDefeat = 10

	// HPGainPerLevel は新記録レベル更新時の1レベルあたりHP増加量です。
	HPGainPerLevel = 10

	// SelectableLevelBonus は撃破済み敵の選択可能レベル上限ボーナスです。
	// 選択可能最大レベル = 最大撃破レベル + SelectableLevelBonus
	SelectableLevelBonus = 5
)

// EnemyDefeatRecord は敵1体の撃破記録を表すエンティティです。
type EnemyDefeatRecord struct {
	// Defeated は撃破済みフラグです。
	Defeated bool

	// MaxDefeatedLevel は撃破済み最大レベルです（未撃破なら0）。
	MaxDefeatedLevel int
}

// EnemyProgress は敵進行状況を管理するドメインモデルです。
// PlayerMaxHPは保持しません（PlayerModelが唯一の権威）。
type EnemyProgress struct {
	// CurrentRank は現在解放済みランクです（1から開始）。
	CurrentRank int

	// DefeatRecords は敵タイプIDごとの撃破記録です。
	DefeatRecords map[string]EnemyDefeatRecord
}

// NewEnemyProgress は新しいEnemyProgressを作成します。
// CurrentRankは1から開始します。
func NewEnemyProgress() *EnemyProgress {
	return &EnemyProgress{
		CurrentRank:   1,
		DefeatRecords: make(map[string]EnemyDefeatRecord),
	}
}

// IsDefeated は指定した敵タイプが撃破済みかどうかを返します。
func (p *EnemyProgress) IsDefeated(enemyTypeID string) bool {
	record, exists := p.DefeatRecords[enemyTypeID]
	return exists && record.Defeated
}

// GetMaxDefeatedLevel は指定した敵タイプの撃破済み最大レベルを返します。
// 未撃破の場合は0を返します。
func (p *EnemyProgress) GetMaxDefeatedLevel(enemyTypeID string) int {
	record, exists := p.DefeatRecords[enemyTypeID]
	if !exists {
		return 0
	}
	return record.MaxDefeatedLevel
}

// GetSelectableLevelRange は敵の選択可能レベル範囲を返します。
// defaultLevel は未撃破時のデフォルトレベルです。
//
// 返り値:
//   - min: 選択可能最小レベル（常に1または defaultLevel）
//   - max: 選択可能最大レベル（未撃破: defaultLevel、撃破済み: MaxDefeatedLevel + 5）
func (p *EnemyProgress) GetSelectableLevelRange(enemyTypeID string, defaultLevel int) (min, max int) {
	if !p.IsDefeated(enemyTypeID) {
		// 未撃破: デフォルトレベルのみ
		return defaultLevel, defaultLevel
	}

	// 撃破済み: 1 から (最大撃破レベル + ボーナス) まで
	maxDefeatedLevel := p.GetMaxDefeatedLevel(enemyTypeID)
	return 1, maxDefeatedLevel + SelectableLevelBonus
}

// RecordDefeat は撃破を記録し、HP増加量を返します。
// HP適用は呼び出し側が担当します。
//
// HP増加ルール:
//   - 初撃破: +10
//   - 新記録レベル: +(新記録レベル - 旧最大レベル) * 10
//   - 同レベル以下: 0
//
// 返り値:
//   - hpGain: 獲得HP増加量
//   - rankUnlocked: 新ランク解放フラグ（現在は常にfalse、呼び出し側で判定）
func (p *EnemyProgress) RecordDefeat(enemyTypeID string, level int) (hpGain int, rankUnlocked bool) {
	record, exists := p.DefeatRecords[enemyTypeID]

	if !exists {
		// 初撃破
		p.DefeatRecords[enemyTypeID] = EnemyDefeatRecord{
			Defeated:         true,
			MaxDefeatedLevel: level,
		}
		return HPGainPerFirstDefeat, false
	}

	if level <= record.MaxDefeatedLevel {
		// 同レベル以下: HP増加なし
		return 0, false
	}

	// 新記録レベル: 差分 * 10
	hpGain = (level - record.MaxDefeatedLevel) * HPGainPerLevel
	record.MaxDefeatedLevel = level
	p.DefeatRecords[enemyTypeID] = record

	return hpGain, false
}
