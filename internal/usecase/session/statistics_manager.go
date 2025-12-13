package session

// StatisticsManager はゲームの統計情報を管理する構造体です。
// タイピング統計とバトル統計を一元管理します。
type StatisticsManager struct {
	// typing はタイピング関連の統計です。
	typing *TypingStatistics

	// battle はバトル関連の統計です。
	battle *BattleStatisticsData
}

// TypingStatistics はタイピング統計を表す構造体です。
type TypingStatistics struct {
	// MaxWPM は最高WPMです。
	MaxWPM int

	// TotalWPM はWPMの累計です（平均計算用）。
	TotalWPM float64

	// TotalSessions はタイピングセッション数です。
	TotalSessions int

	// PerfectAccuracyCount は100%正確性を達成した回数です。
	PerfectAccuracyCount int

	// TotalCharacters は総タイプ文字数です。
	TotalCharacters int

	// TotalCorrectCharacters は正しくタイプした文字数です。
	TotalCorrectCharacters int

	// TotalMissedCharacters はミスした文字数です。
	TotalMissedCharacters int
}

// BattleStatisticsData はバトル統計を表す構造体です。
type BattleStatisticsData struct {
	// TotalBattles は総バトル数です。
	TotalBattles int

	// Wins は勝利数です。
	Wins int

	// Losses は敗北数です。
	Losses int

	// MaxLevelReached は到達した最高レベルです。
	MaxLevelReached int

	// TotalEnemiesDefeated は倒した敵の総数です。
	TotalEnemiesDefeated int

	// TotalDamageDealt は与えた総ダメージです。
	TotalDamageDealt int

	// TotalDamageTaken は受けた総ダメージです。
	TotalDamageTaken int

	// TotalHealingDone は回復した総量です。
	TotalHealingDone int
}

// NewStatisticsManager は新しいStatisticsManagerを作成します。
func NewStatisticsManager() *StatisticsManager {
	return &StatisticsManager{
		typing: &TypingStatistics{},
		battle: &BattleStatisticsData{},
	}
}

// Typing はタイピング統計を返します。
func (m *StatisticsManager) Typing() *TypingStatistics {
	return m.typing
}

// Battle はバトル統計を返します。
func (m *StatisticsManager) Battle() *BattleStatisticsData {
	return m.battle
}

// RecordTypingResult はタイピング結果を記録します。
func (m *StatisticsManager) RecordTypingResult(wpm int, accuracy float64, characters int, correct int, missed int) {
	m.typing.TotalSessions++
	m.typing.TotalWPM += float64(wpm)
	m.typing.TotalCharacters += characters
	m.typing.TotalCorrectCharacters += correct
	m.typing.TotalMissedCharacters += missed

	if wpm > m.typing.MaxWPM {
		m.typing.MaxWPM = wpm
	}

	if accuracy >= 100.0 {
		m.typing.PerfectAccuracyCount++
	}
}

// RecordBattleResult はバトル結果を記録します。
func (m *StatisticsManager) RecordBattleResult(victory bool, level int) {
	m.battle.TotalBattles++
	if victory {
		m.battle.Wins++
		m.battle.TotalEnemiesDefeated++
		if level > m.battle.MaxLevelReached {
			m.battle.MaxLevelReached = level
		}
	} else {
		m.battle.Losses++
	}
}

// RecordDamageDealt は与えたダメージを記録します。
func (m *StatisticsManager) RecordDamageDealt(damage int) {
	m.battle.TotalDamageDealt += damage
}

// RecordDamageTaken は受けたダメージを記録します。
func (m *StatisticsManager) RecordDamageTaken(damage int) {
	m.battle.TotalDamageTaken += damage
}

// RecordHealing は回復量を記録します。
func (m *StatisticsManager) RecordHealing(amount int) {
	m.battle.TotalHealingDone += amount
}

// GetAverageWPM は平均WPMを返します。
func (m *StatisticsManager) GetAverageWPM() float64 {
	if m.typing.TotalSessions == 0 {
		return 0
	}
	return m.typing.TotalWPM / float64(m.typing.TotalSessions)
}

// GetAccuracyRate は全体の正確性を返します。
func (m *StatisticsManager) GetAccuracyRate() float64 {
	total := m.typing.TotalCorrectCharacters + m.typing.TotalMissedCharacters
	if total == 0 {
		return 0
	}
	return float64(m.typing.TotalCorrectCharacters) / float64(total) * 100
}

// GetWinRate は勝率を返します。
func (m *StatisticsManager) GetWinRate() float64 {
	if m.battle.TotalBattles == 0 {
		return 0
	}
	return float64(m.battle.Wins) / float64(m.battle.TotalBattles) * 100
}

// RecordTypingStats はバトル中のタイピング統計を記録します（簡易版）。
func (m *StatisticsManager) RecordTypingStats(wpm float64, accuracy float64) {
	m.typing.TotalSessions++
	m.typing.TotalWPM += wpm
	if int(wpm) > m.typing.MaxWPM {
		m.typing.MaxWPM = int(wpm)
	}
	// 正確性からCorrect/Missedを概算
	correctCount := int(accuracy)
	m.typing.TotalCorrectCharacters += correctCount
	m.typing.TotalMissedCharacters += 100 - correctCount
	if accuracy >= 100.0 {
		m.typing.PerfectAccuracyCount++
	}
}

// StatisticsSaveData はセーブ用の統計データです。
type StatisticsSaveData struct {
	TotalBattles         int
	Victories            int
	Defeats              int
	MaxLevelReached      int
	HighestWPM           float64
	AverageWPM           float64
	PerfectAccuracyCount int
	TotalCharactersTyped int
}

// LoadFromSaveData はセーブデータから統計を復元します。
func (m *StatisticsManager) LoadFromSaveData(data *StatisticsSaveData) {
	if data == nil {
		return
	}
	m.battle.TotalBattles = data.TotalBattles
	m.battle.Wins = data.Victories
	m.battle.Losses = data.Defeats
	m.battle.MaxLevelReached = data.MaxLevelReached
	m.typing.MaxWPM = int(data.HighestWPM)
	// 平均WPMからTotalWPMを逆算
	if data.AverageWPM > 0 {
		m.typing.TotalSessions = 1
		m.typing.TotalWPM = data.AverageWPM
	}
	m.typing.PerfectAccuracyCount = data.PerfectAccuracyCount
	m.typing.TotalCharacters = data.TotalCharactersTyped
}
