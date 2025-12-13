// Package presenter はUI用のデータ変換を提供します。
// usecase層のデータをtui/screens用のViewModelに変換します。
package presenter

import (
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/usecase/session"
)

// CreateStatsData はGameStateから統計データを生成します。
func CreateStatsData(gs *session.GameState) *screens.StatsData {
	stats := gs.Statistics()
	achievements := gs.Achievements()

	// 実績データを変換
	allAchievements := achievements.GetAllAchievements()
	achievementData := make([]screens.AchievementData, 0, len(allAchievements))
	for _, ach := range allAchievements {
		achievementData = append(achievementData, screens.AchievementData{
			ID:          ach.ID,
			Name:        ach.Name,
			Description: ach.Description,
			Achieved:    achievements.IsUnlocked(ach.ID),
		})
	}

	return &screens.StatsData{
		TypingStats: screens.TypingStatsData{
			MaxWPM:               stats.Typing().MaxWPM,
			AverageWPM:           stats.GetAverageWPM(),
			PerfectAccuracyCount: stats.Typing().PerfectAccuracyCount,
			TotalCharacters:      stats.Typing().TotalCharacters,
		},
		BattleStats: screens.BattleStatsData{
			TotalBattles:    stats.Battle().TotalBattles,
			Wins:            stats.Battle().Wins,
			Losses:          stats.Battle().Losses,
			MaxLevelReached: gs.MaxLevelReached,
		},
		Achievements: achievementData,
	}
}
