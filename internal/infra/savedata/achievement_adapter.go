// Package savedata は実績データの永続化アダプターを提供します。
// achievementパッケージからの依存を解消するための変換関数を含みます。
package savedata

// AchievementStateToSaveData は解除済み実績のリストをセーブデータ形式に変換します。
// achievementパッケージの内部状態をpersistenceパッケージのセーブデータ型に変換します。
func AchievementStateToSaveData(unlockedIDs []string) *AchievementsSaveData {
	unlocked := make([]string, len(unlockedIDs))
	copy(unlocked, unlockedIDs)

	return &AchievementsSaveData{
		Unlocked: unlocked,
		Progress: make(map[string]int),
	}
}

// SaveDataToAchievementState はセーブデータから解除済み実績のリストを抽出します。
// persistenceパッケージのセーブデータ型からachievementパッケージで使用可能な形式に変換します。
func SaveDataToAchievementState(data *AchievementsSaveData) []string {
	if data == nil {
		return []string{}
	}

	unlocked := make([]string, len(data.Unlocked))
	copy(unlocked, data.Unlocked)
	return unlocked
}
