// Package achievement は実績システムを担当します。
// タイピング実績とバトル実績の管理、達成判定、通知処理を提供します。

package achievement

// ==================================================
// 実績ID定数
// ==================================================

const (
	// タイピング実績 (Requirements 15.4, 15.5)
	AchievementWPM50           = "wpm_50"
	AchievementWPM80           = "wpm_80"
	AchievementWPM100          = "wpm_100"
	AchievementWPM120          = "wpm_120"
	AchievementPerfectAccuracy = "perfect_accuracy"

	// バトル実績 (Requirements 15.6, 15.7, 15.8)
	AchievementDefeat10  = "defeat_10"
	AchievementDefeat50  = "defeat_50"
	AchievementDefeat100 = "defeat_100"
	AchievementDefeat500 = "defeat_500"
	AchievementLevel10   = "level_10"
	AchievementLevel25   = "level_25"
	AchievementLevel50   = "level_50"
	AchievementLevel100  = "level_100"
	AchievementNoDamage  = "no_damage"
)

// ==================================================
// 実績定義
// ==================================================

// AchievementDefinition は実績の定義を表す構造体です。
type AchievementDefinition struct {
	// ID は実績の一意識別子です。
	ID string

	// Name は実績の表示名です。
	Name string

	// Description は実績の説明文です。
	Description string

	// Category は実績のカテゴリです（typing / battle）。
	Category string
}

// allAchievements は全実績の定義リストです。
var allAchievements = []AchievementDefinition{
	// タイピング実績
	{AchievementWPM50, "タイピング見習い", "WPM 50 を達成する", "typing"},
	{AchievementWPM80, "タイピング上手", "WPM 80 を達成する", "typing"},
	{AchievementWPM100, "タイピングマスター", "WPM 100 を達成する", "typing"},
	{AchievementWPM120, "タイピングレジェンド", "WPM 120 を達成する", "typing"},
	{AchievementPerfectAccuracy, "完璧主義者", "100%正確性でクリアする", "typing"},

	// バトル実績
	{AchievementDefeat10, "新米ハンター", "敵を10体撃破する", "battle"},
	{AchievementDefeat50, "歴戦の戦士", "敵を50体撃破する", "battle"},
	{AchievementDefeat100, "百戦錬磨", "敵を100体撃破する", "battle"},
	{AchievementDefeat500, "伝説の勇者", "敵を500体撃破する", "battle"},
	{AchievementLevel10, "探索者", "レベル10に到達する", "battle"},
	{AchievementLevel25, "冒険者", "レベル25に到達する", "battle"},
	{AchievementLevel50, "英雄", "レベル50に到達する", "battle"},
	{AchievementLevel100, "覇者", "レベル100に到達する", "battle"},
	{AchievementNoDamage, "無傷の勝利", "ノーダメージでバトルに勝利する", "battle"},
}

// ==================================================
// 実績通知
// ==================================================

// AchievementNotification は実績達成通知を表す構造体です。

type AchievementNotification struct {
	// AchievementID は達成した実績のIDです。
	AchievementID string

	// Name は実績の表示名です。
	Name string

	// Description は実績の説明文です。
	Description string
}

// ==================================================
// 実績マネージャー
// ==================================================

// AchievementManager は実績の管理を担当する構造体です。

type AchievementManager struct {
	// unlocked は解除済み実績IDのマップです。
	unlocked map[string]bool
}

// NewAchievementManager は新しいAchievementManagerを作成します。
func NewAchievementManager() *AchievementManager {
	return &AchievementManager{
		unlocked: make(map[string]bool),
	}
}

// CheckTypingAchievements はタイピング成績を基に実績の解除をチェックします。

func (m *AchievementManager) CheckTypingAchievements(wpm float64, accuracy float64) []AchievementNotification {
	var notifications []AchievementNotification

	// WPMマイルストーン (Requirement 15.4)
	if wpm >= 50 {
		if n := m.tryUnlock(AchievementWPM50); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if wpm >= 80 {
		if n := m.tryUnlock(AchievementWPM80); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if wpm >= 100 {
		if n := m.tryUnlock(AchievementWPM100); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if wpm >= 120 {
		if n := m.tryUnlock(AchievementWPM120); n != nil {
			notifications = append(notifications, *n)
		}
	}

	// 100%正確性 (Requirement 15.5)
	if accuracy >= 100 {
		if n := m.tryUnlock(AchievementPerfectAccuracy); n != nil {
			notifications = append(notifications, *n)
		}
	}

	return notifications
}

// CheckBattleAchievements はバトル成績を基に実績の解除をチェックします。

func (m *AchievementManager) CheckBattleAchievements(totalDefeated int, maxLevel int, isNoDamage bool) []AchievementNotification {
	var notifications []AchievementNotification

	// 敵撃破数マイルストーン (Requirement 15.6)
	if totalDefeated >= 10 {
		if n := m.tryUnlock(AchievementDefeat10); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if totalDefeated >= 50 {
		if n := m.tryUnlock(AchievementDefeat50); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if totalDefeated >= 100 {
		if n := m.tryUnlock(AchievementDefeat100); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if totalDefeated >= 500 {
		if n := m.tryUnlock(AchievementDefeat500); n != nil {
			notifications = append(notifications, *n)
		}
	}

	// レベルマイルストーン (Requirement 15.7)
	if maxLevel >= 10 {
		if n := m.tryUnlock(AchievementLevel10); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if maxLevel >= 25 {
		if n := m.tryUnlock(AchievementLevel25); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if maxLevel >= 50 {
		if n := m.tryUnlock(AchievementLevel50); n != nil {
			notifications = append(notifications, *n)
		}
	}
	if maxLevel >= 100 {
		if n := m.tryUnlock(AchievementLevel100); n != nil {
			notifications = append(notifications, *n)
		}
	}

	// ノーダメージクリア (Requirement 15.8)
	if isNoDamage {
		if n := m.tryUnlock(AchievementNoDamage); n != nil {
			notifications = append(notifications, *n)
		}
	}

	return notifications
}

// tryUnlock は実績の解除を試み、成功時に通知を返します。
// 既に解除済みの場合はnilを返します。
func (m *AchievementManager) tryUnlock(achievementID string) *AchievementNotification {
	if m.unlocked[achievementID] {
		return nil // 既に解除済み
	}

	// 実績定義を検索
	var def *AchievementDefinition
	for _, a := range allAchievements {
		if a.ID == achievementID {
			def = &a
			break
		}
	}

	if def == nil {
		return nil // 定義が見つからない
	}

	// 解除
	m.unlocked[achievementID] = true

	return &AchievementNotification{
		AchievementID: def.ID,
		Name:          def.Name,
		Description:   def.Description,
	}
}

// IsUnlocked は指定した実績が解除済みかを返します。

func (m *AchievementManager) IsUnlocked(achievementID string) bool {
	return m.unlocked[achievementID]
}

// GetAllAchievements は全実績の定義リストを返します。
func (m *AchievementManager) GetAllAchievements() []AchievementDefinition {
	return allAchievements
}

// GetCompletionRate は実績の達成率（0.0〜1.0）を返します。

func (m *AchievementManager) GetCompletionRate() float64 {
	total := len(allAchievements)
	if total == 0 {
		return 0.0
	}
	unlocked := 0
	for _, a := range allAchievements {
		if m.unlocked[a.ID] {
			unlocked++
		}
	}
	return float64(unlocked) / float64(total)
}

// GetUnlockedCount は解除済み実績の数を返します。
func (m *AchievementManager) GetUnlockedCount() int {
	count := 0
	for _, a := range allAchievements {
		if m.unlocked[a.ID] {
			count++
		}
	}
	return count
}

// GetTotalCount は全実績の数を返します。
func (m *AchievementManager) GetTotalCount() int {
	return len(allAchievements)
}

// ==================================================
// セーブ/ロード（ドメイン型のみ使用）
// ==================================================

// GetUnlockedIDs は解除済み実績IDのリストを返します。
func (m *AchievementManager) GetUnlockedIDs() []string {
	unlocked := make([]string, 0, len(m.unlocked))
	for id, isUnlocked := range m.unlocked {
		if isUnlocked {
			unlocked = append(unlocked, id)
		}
	}
	return unlocked
}

// LoadFromUnlockedIDs は解除済み実績IDリストから状態を復元します。
func (m *AchievementManager) LoadFromUnlockedIDs(unlockedIDs []string) {
	m.unlocked = make(map[string]bool)
	for _, id := range unlockedIDs {
		m.unlocked[id] = true
	}
}
