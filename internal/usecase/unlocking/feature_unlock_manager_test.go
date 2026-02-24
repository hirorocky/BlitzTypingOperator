package unlocking_test

import (
	"fmt"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/unlocking"
)

// テスト用のUnlockRuleとTutorialDefを作成するヘルパー
func testRules() []domain.UnlockRule {
	return []domain.UnlockRule{
		{FeatureID: domain.FeatureDefenseSkill, UnlockRank: 2, TutorialID: "tut_defense"},
		{FeatureID: domain.FeatureAgentCustomization, UnlockRank: 3, TutorialID: "tut_agent"},
		{FeatureID: domain.FeatureChainEffect, UnlockRank: 4, TutorialID: "tut_chain"},
		{FeatureID: domain.FeatureManaSystem, UnlockRank: 5, TutorialID: "tut_mana"},
	}
}

func testTutorials() []domain.TutorialDef {
	return []domain.TutorialDef{
		{ID: "tut_basic", Title: "基本操作", Pages: []string{"基本操作の説明"}, DefaultVisible: true},
		{ID: "tut_defense", Title: "ディフェンススキル", Pages: []string{"ディフェンスの説明"}, DefaultVisible: false, FeatureID: domain.FeatureDefenseSkill},
		{ID: "tut_agent", Title: "エージェントカスタマイズ", Pages: []string{"カスタマイズの説明"}, DefaultVisible: false, FeatureID: domain.FeatureAgentCustomization},
		{ID: "tut_chain", Title: "チェイン効果", Pages: []string{"チェイン効果の説明"}, DefaultVisible: false, FeatureID: domain.FeatureChainEffect},
		{ID: "tut_mana", Title: "マナシステム", Pages: []string{"マナの説明"}, DefaultVisible: false, FeatureID: domain.FeatureManaSystem},
	}
}

func newTestManager() *unlocking.Manager {
	state := domain.NewFeatureUnlockState()
	mgr, _ := unlocking.NewManager(testRules(), testTutorials(), state)
	return mgr
}

// === 受け入れ基準 2-5: ランク到達時に対応機能がPendingTutorialに遷移する ===
func TestManager_ApplyRank_UnlocksCorrectFeatures(t *testing.T) {
	testCases := []struct {
		rank            int
		expectedPending []domain.FeatureID
	}{
		{1, nil},
		{2, []domain.FeatureID{domain.FeatureDefenseSkill}},
		{3, []domain.FeatureID{domain.FeatureDefenseSkill, domain.FeatureAgentCustomization}},
		{4, []domain.FeatureID{domain.FeatureDefenseSkill, domain.FeatureAgentCustomization, domain.FeatureChainEffect}},
		{5, []domain.FeatureID{domain.FeatureDefenseSkill, domain.FeatureAgentCustomization, domain.FeatureChainEffect, domain.FeatureManaSystem}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("rank%d", tc.rank), func(t *testing.T) {
			mgr := newTestManager()
			_, err := mgr.ApplyRank(tc.rank)
			if err != nil {
				t.Fatalf("ApplyRank(%d)でエラー: %v", tc.rank, err)
			}

			for _, f := range tc.expectedPending {
				if !mgr.IsPendingOrUnlocked(f) {
					t.Errorf("ランク%dでFeature %s がPendingTutorialまたはUnlockedでない", tc.rank, f)
				}
			}
		})
	}
}

// === 受け入れ基準 6: 解放処理は冪等 ===
func TestManager_ApplyRank_Idempotent(t *testing.T) {
	mgr := newTestManager()

	delta1, _ := mgr.ApplyRank(2)
	delta2, _ := mgr.ApplyRank(2)

	if len(delta1.NewPendingFeatures) == 0 {
		t.Error("最初のApplyRank(2)でPending機能が0件")
	}
	if len(delta2.NewPendingFeatures) != 0 {
		t.Errorf("2回目のApplyRank(2)でPending機能が%d件（期待: 0件）", len(delta2.NewPendingFeatures))
	}
}

// === 受け入れ基準 7: ランクを一気に跨いだ場合、中間ランクの機能も漏れなくPendingTutorial化 ===
func TestManager_ApplyRank_SkipRanks(t *testing.T) {
	mgr := newTestManager()

	delta, err := mgr.ApplyRank(4)
	if err != nil {
		t.Fatalf("ApplyRank(4)でエラー: %v", err)
	}

	// ランク2,3,4の機能が全てPendingになるべき
	expected := []domain.FeatureID{
		domain.FeatureDefenseSkill,
		domain.FeatureAgentCustomization,
		domain.FeatureChainEffect,
	}

	if len(delta.NewPendingFeatures) != len(expected) {
		t.Errorf("Pending機能数 = %d, 期待: %d", len(delta.NewPendingFeatures), len(expected))
	}
}

// === 受け入れ基準 13: 複数PendingTutorialはランク昇順で表示 ===
func TestManager_NextPendingTutorial_OrderedByRank(t *testing.T) {
	mgr := newTestManager()
	mgr.ApplyRank(4) // ランク2,3,4の機能がPending

	// 最初のPendingはランク2のdefense_skill
	tutID, ok := mgr.NextPendingTutorial()
	if !ok {
		t.Fatal("NextPendingTutorialがfalseを返した")
	}
	if tutID != "tut_defense" {
		t.Errorf("最初のPendingチュートリアル = %s, 期待: tut_defense", tutID)
	}
}

// === 受け入れ基準 12: チュートリアル完了で機能がUnlockedに遷移 ===
func TestManager_CompleteTutorial_UnlocksFeature(t *testing.T) {
	mgr := newTestManager()
	mgr.ApplyRank(2)

	_, err := mgr.CompleteTutorial("tut_defense")
	if err != nil {
		t.Fatalf("CompleteTutorialでエラー: %v", err)
	}

	if !mgr.IsUnlocked(domain.FeatureDefenseSkill) {
		t.Error("CompleteTutorial後にdefense_skillがUnlockedでない")
	}
}

// === 受け入れ基準 9: 不正な遷移要求（LockedのままCompleteTutorial） ===
func TestManager_CompleteTutorial_LockedFeature_ReturnsError(t *testing.T) {
	mgr := newTestManager()

	// defense_skillはLockedのまま
	_, err := mgr.CompleteTutorial("tut_defense")
	if err == nil {
		t.Error("Locked状態でCompleteTutorialがエラーを返さなかった")
	}
}

// === 受け入れ基準 21,22,23: ListVisibleTutorials ===
func TestManager_ListVisibleTutorials(t *testing.T) {
	mgr := newTestManager()

	// 初期状態: default_visibleのみ
	tutorials := mgr.ListVisibleTutorials()
	if len(tutorials) != 1 {
		t.Errorf("初期状態の表示可能チュートリアル数 = %d, 期待: 1", len(tutorials))
	}
	if tutorials[0].ID != "tut_basic" {
		t.Errorf("初期チュートリアル = %s, 期待: tut_basic", tutorials[0].ID)
	}

	// defense_skillをUnlockした後
	mgr.ApplyRank(2)
	mgr.CompleteTutorial("tut_defense")

	tutorials = mgr.ListVisibleTutorials()
	if len(tutorials) != 2 {
		t.Errorf("defense解放後の表示可能チュートリアル数 = %d, 期待: 2", len(tutorials))
	}
}

// === 受け入れ基準 30: Reconcile処理 ===
func TestManager_Reconcile_AppliesNewMasterData(t *testing.T) {
	// ランク3に到達済みだが、ランク2,3の機能がLockedのまま（旧セーブ想定）
	state := domain.NewFeatureUnlockState()
	mgr, _ := unlocking.NewManager(testRules(), testTutorials(), state)

	delta, err := mgr.Reconcile(3)
	if err != nil {
		t.Fatalf("Reconcileでエラー: %v", err)
	}

	if len(delta.NewPendingFeatures) != 2 {
		t.Errorf("Reconcile後のPending機能数 = %d, 期待: 2 (defense_skill, agent_customization)", len(delta.NewPendingFeatures))
	}
}

// === 受け入れ基準 28: Unlocked機能はReconcileでも再ロックされない ===
func TestManager_Reconcile_DoesNotRelockUnlocked(t *testing.T) {
	state := domain.NewFeatureUnlockState()
	// defense_skillを事前にUnlockedに設定
	state.TransitionTo(domain.FeatureDefenseSkill, domain.FeatureUnlocked)

	mgr, _ := unlocking.NewManager(testRules(), testTutorials(), state)
	mgr.Reconcile(3)

	if !mgr.IsUnlocked(domain.FeatureDefenseSkill) {
		t.Error("Reconcile後にdefense_skillがUnlockedでなくなった")
	}
}

// === Snapshot: ディープコピーの検証 ===
func TestManager_Snapshot_DeepCopy(t *testing.T) {
	mgr := newTestManager()
	mgr.ApplyRank(2)

	snap1 := mgr.Snapshot()

	// Snapshot取得後にManagerの状態を変更
	mgr.CompleteTutorial("tut_defense")

	snap2 := mgr.Snapshot()

	// snap1はCompleteTutorial前の状態を保持しているべき
	if snap1.GetStatus(domain.FeatureDefenseSkill) == snap2.GetStatus(domain.FeatureDefenseSkill) {
		t.Error("Snapshotがディープコピーでない（Managerの状態変更がSnapshotに影響）")
	}
}
