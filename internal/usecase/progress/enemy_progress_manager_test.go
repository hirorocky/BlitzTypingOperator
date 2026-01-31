// Package progress は敵進行管理のユースケースを提供します。
package progress

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestNewEnemyProgressManager はEnemyProgressManagerの初期化をテストします。
func TestNewEnemyProgressManager(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	if manager == nil {
		t.Fatal("EnemyProgressManagerの作成に失敗しました")
	}
	if manager.progress != progress {
		t.Error("progressが正しく設定されていません")
	}
	if manager.player != player {
		t.Error("playerが正しく設定されていません")
	}
	if len(manager.enemyTypes) != len(enemyTypes) {
		t.Errorf("enemyTypesの数が一致しません: got=%d, want=%d", len(manager.enemyTypes), len(enemyTypes))
	}
}

// TestGetCurrentRankEnemies は現在ランクの敵リスト取得をテストします。
func TestGetCurrentRankEnemies(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 初期状態ではランク1の敵のみ取得可能
	enemies := manager.GetCurrentRankEnemies()

	rank1Count := 0
	for _, et := range enemyTypes {
		if et.Rank == 1 {
			rank1Count++
		}
	}

	if len(enemies) != rank1Count {
		t.Errorf("ランク1の敵数が一致しません: got=%d, want=%d", len(enemies), rank1Count)
	}

	// 全てがランク1であることを確認
	for _, e := range enemies {
		if e.Rank != 1 {
			t.Errorf("ランク1以外の敵が含まれています: %s (Rank=%d)", e.ID, e.Rank)
		}
	}
}

// TestGetSelectableLevelRange は選択可能レベル範囲取得をテストします。
func TestGetSelectableLevelRange(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 未撃破の敵：デフォルトレベルのみ選択可能
	minLevel, maxLevel := manager.GetSelectableLevelRange("slime")
	if minLevel != 1 || maxLevel != 1 {
		t.Errorf("未撃破敵の範囲が不正: min=%d, max=%d (want min=1, max=1)", minLevel, maxLevel)
	}

	// 撃破済みの敵をシミュレート
	progress.RecordDefeat("slime", 1)
	minLevel, maxLevel = manager.GetSelectableLevelRange("slime")
	// 撃破済み: 1 ~ (1 + 5) = 1 ~ 6
	if minLevel != 1 || maxLevel != 6 {
		t.Errorf("撃破済み敵の範囲が不正: min=%d, max=%d (want min=1, max=6)", minLevel, maxLevel)
	}

	// 高レベル撃破後
	progress.RecordDefeat("slime", 10)
	minLevel, maxLevel = manager.GetSelectableLevelRange("slime")
	// 撃破済み: 1 ~ (10 + 5) = 1 ~ 15
	if minLevel != 1 || maxLevel != 15 {
		t.Errorf("高レベル撃破後の範囲が不正: min=%d, max=%d (want min=1, max=15)", minLevel, maxLevel)
	}
}

// TestRecordVictory_FirstDefeat は初撃破時のHP増加をテストします。
func TestRecordVictory_FirstDefeat(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 初期MaxHP確認
	initialMaxHP := player.MaxHP
	if initialMaxHP != domain.InitialMaxHP {
		t.Fatalf("初期MaxHPが不正: got=%d, want=%d", initialMaxHP, domain.InitialMaxHP)
	}

	// 初撃破
	result := manager.RecordVictory("slime", 1)

	// HP増加量の確認
	if result.HPGain != domain.HPGainPerFirstDefeat {
		t.Errorf("初撃破のHP増加量が不正: got=%d, want=%d", result.HPGain, domain.HPGainPerFirstDefeat)
	}

	// PlayerModelのMaxHPが増加していることを確認
	if player.MaxHP != initialMaxHP+domain.HPGainPerFirstDefeat {
		t.Errorf("PlayerのMaxHPが正しく増加していません: got=%d, want=%d", player.MaxHP, initialMaxHP+domain.HPGainPerFirstDefeat)
	}

	// 撃破済みであることを確認
	if !progress.IsDefeated("slime") {
		t.Error("敵が撃破済みとしてマークされていません")
	}
}

// TestRecordVictory_HigherLevel は高レベル撃破時のHP増加をテストします。
func TestRecordVictory_HigherLevel(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 初撃破（レベル1）
	manager.RecordVictory("slime", 1)
	hpAfterFirst := player.MaxHP

	// 同レベルで再撃破（HP増加なし）
	result := manager.RecordVictory("slime", 1)
	if result.HPGain != 0 {
		t.Errorf("同レベル撃破でHP増加が発生: got=%d, want=0", result.HPGain)
	}
	if player.MaxHP != hpAfterFirst {
		t.Errorf("同レベル撃破でMaxHPが変化: got=%d, want=%d", player.MaxHP, hpAfterFirst)
	}

	// 高レベルで撃破（レベル5）
	result = manager.RecordVictory("slime", 5)
	// 差分 = 5 - 1 = 4、HP増加 = 4 * 10 = 40
	expectedGain := 4 * domain.HPGainPerLevel
	if result.HPGain != expectedGain {
		t.Errorf("高レベル撃破のHP増加量が不正: got=%d, want=%d", result.HPGain, expectedGain)
	}
	if player.MaxHP != hpAfterFirst+expectedGain {
		t.Errorf("高レベル撃破後のMaxHPが不正: got=%d, want=%d", player.MaxHP, hpAfterFirst+expectedGain)
	}
}

// TestCheckRankComplete はランク内全敵撃破チェックをテストします。
func TestCheckRankComplete(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 初期状態ではランク完了していない
	if manager.CheckRankComplete() {
		t.Error("初期状態でランク完了と判定されました")
	}

	// ランク1の敵を一部撃破（直接progressに記録してRecordVictoryのランク進行を避ける）
	progress.RecordDefeat("slime", 1)
	if manager.CheckRankComplete() {
		t.Error("一部撃破でランク完了と判定されました")
	}

	// ランク1の敵を全て撃破
	progress.RecordDefeat("bat", 1)

	// ランク完了チェック
	if !manager.CheckRankComplete() {
		t.Error("ランク1の敵を全撃破後もランク未完了と判定されました")
	}
}

// TestRecordVictory_RankUnlock はランク解放をテストします。
func TestRecordVictory_RankUnlock(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// ランク1の敵を取得
	rank1Enemies := manager.GetCurrentRankEnemies()

	// ランク1の敵を全て撃破（最後の1体以外）
	for i := 0; i < len(rank1Enemies)-1; i++ {
		result := manager.RecordVictory(rank1Enemies[i].ID, 1)
		if result.RankUnlocked {
			t.Error("最後の敵撃破前にランク解放が発生しました")
		}
	}

	// 初期ランク確認
	if progress.CurrentRank != 1 {
		t.Errorf("初期ランクが不正: got=%d, want=1", progress.CurrentRank)
	}

	// 最後の敵を撃破
	result := manager.RecordVictory(rank1Enemies[len(rank1Enemies)-1].ID, 1)
	if !result.RankUnlocked {
		t.Error("最後の敵撃破後にランク解放が発生しませんでした")
	}

	// ランク2に進行
	if progress.CurrentRank != 2 {
		t.Errorf("ランク解放後のランクが不正: got=%d, want=2", progress.CurrentRank)
	}

	// ランク2の敵が取得可能か確認
	rank2Enemies := manager.GetCurrentRankEnemies()
	for _, e := range rank2Enemies {
		if e.Rank != 2 {
			t.Errorf("ランク2以外の敵が取得されました: %s (Rank=%d)", e.ID, e.Rank)
		}
	}
}

// TestGetEnemyByID は敵タイプ取得をテストします。
func TestGetEnemyByID(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 存在する敵
	enemy, found := manager.GetEnemyByID("slime")
	if !found {
		t.Error("存在する敵が見つかりませんでした")
	}
	if enemy.ID != "slime" {
		t.Errorf("取得した敵のIDが不正: got=%s, want=slime", enemy.ID)
	}

	// 存在しない敵
	_, found = manager.GetEnemyByID("unknown")
	if found {
		t.Error("存在しない敵が見つかりました")
	}
}

// TestIsEnemyInCurrentRank は敵が現在ランク内かどうかのテストです。
func TestIsEnemyInCurrentRank(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// ランク1の敵は現在ランク内
	if !manager.IsEnemyInCurrentRank("slime") {
		t.Error("ランク1の敵が現在ランク内と判定されませんでした")
	}

	// ランク2の敵は現在ランク外
	if manager.IsEnemyInCurrentRank("goblin") {
		t.Error("ランク2の敵が現在ランク内と判定されました")
	}

	// 存在しない敵
	if manager.IsEnemyInCurrentRank("unknown") {
		t.Error("存在しない敵が現在ランク内と判定されました")
	}
}

// TestVictoryResult はVictoryResult構造体の確認をテストします。
func TestVictoryResult(t *testing.T) {
	progress := domain.NewEnemyProgress()
	player := domain.NewPlayerWithMaxHP(domain.InitialMaxHP)
	enemyTypes := createTestEnemyTypes()

	manager := NewEnemyProgressManager(progress, player, enemyTypes)

	// 初撃破
	result := manager.RecordVictory("slime", 1)

	if result.HPGain != domain.HPGainPerFirstDefeat {
		t.Errorf("HPGainが不正: got=%d, want=%d", result.HPGain, domain.HPGainPerFirstDefeat)
	}

	// NewMaxLevelの確認（初撃破なので MaxDefeatedLevel = 1 + 5 = 6）
	if result.NewMaxLevel != 6 {
		t.Errorf("NewMaxLevelが不正: got=%d, want=6", result.NewMaxLevel)
	}
}

// createTestEnemyTypes はテスト用の敵タイプリストを作成します。
func createTestEnemyTypes() map[string]domain.EnemyType {
	return map[string]domain.EnemyType{
		"slime": {
			ID:           "slime",
			Name:         "スライム",
			Rank:         1,
			DefaultLevel: 1,
		},
		"bat": {
			ID:           "bat",
			Name:         "コウモリ",
			Rank:         1,
			DefaultLevel: 1,
		},
		"goblin": {
			ID:           "goblin",
			Name:         "ゴブリン",
			Rank:         2,
			DefaultLevel: 3,
		},
		"orc": {
			ID:           "orc",
			Name:         "オーク",
			Rank:         2,
			DefaultLevel: 5,
		},
	}
}
