package integration_test

import (
	"encoding/json"
	"testing"

	"hirorocky/type-battle/internal/app"
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/savedata"
	"hirorocky/type-battle/internal/usecase/unlocking"
)

// === 受け入れ基準 25,27: セーブ/ロードで機能解放状態が保持される ===
func TestFeatureUnlock_SaveLoadRoundTrip(t *testing.T) {
	rules := []domain.UnlockRule{
		{FeatureID: domain.FeatureDefenseSkill, UnlockRank: 2, TutorialID: "tut_defense"},
		{FeatureID: domain.FeatureAgentCustomization, UnlockRank: 3, TutorialID: "tut_agent"},
	}
	tutorials := []domain.TutorialDef{
		{ID: "tut_defense", Title: "ディフェンス", Pages: []string{"説明"}, FeatureID: domain.FeatureDefenseSkill},
		{ID: "tut_agent", Title: "エージェント", Pages: []string{"説明"}, FeatureID: domain.FeatureAgentCustomization},
	}

	// 1. Managerを作成し、ランク2到達→チュートリアル完了
	state := domain.NewFeatureUnlockState()
	mgr, _ := unlocking.NewManager(rules, tutorials, state)
	mgr.ApplyRank(2)
	mgr.CompleteTutorial("tut_defense")

	// 2. Snapshotを取得してセーブデータに変換
	snap := mgr.Snapshot()
	savedData := app.FeatureUnlockSnapshotToSave(snap)

	// 3. JSONにシリアライズ→デシリアライズ
	jsonBytes, err := json.Marshal(savedData)
	if err != nil {
		t.Fatalf("JSONシリアライズ失敗: %v", err)
	}

	var loadedSave savedata.FeatureUnlockSave
	if err := json.Unmarshal(jsonBytes, &loadedSave); err != nil {
		t.Fatalf("JSONデシリアライズ失敗: %v", err)
	}

	// 4. セーブデータからSnapshotを復元
	restoredSnap := app.FeatureUnlockSaveToSnapshot(loadedSave)

	// 5. 復元したSnapshotでManagerを再作成
	mgr2, _ := unlocking.NewManager(rules, tutorials, restoredSnap)

	// 6. 状態が保持されていることを確認
	if !mgr2.IsUnlocked(domain.FeatureDefenseSkill) {
		t.Error("ロード後にdefense_skillがUnlockedでない")
	}
	if mgr2.IsUnlocked(domain.FeatureAgentCustomization) {
		t.Error("ロード後にagent_customizationがUnlockedになっている（期待: Locked）")
	}
}

// === 受け入れ基準 27: 旧セーブデータからの復元 ===
func TestFeatureUnlock_OldSaveDataMigration(t *testing.T) {
	rules := []domain.UnlockRule{
		{FeatureID: domain.FeatureDefenseSkill, UnlockRank: 2, TutorialID: "tut_defense"},
		{FeatureID: domain.FeatureAgentCustomization, UnlockRank: 3, TutorialID: "tut_agent"},
		{FeatureID: domain.FeatureChainEffect, UnlockRank: 4, TutorialID: "tut_chain"},
	}
	tutorials := []domain.TutorialDef{
		{ID: "tut_defense", Title: "ディフェンス", Pages: []string{"説明"}, FeatureID: domain.FeatureDefenseSkill},
		{ID: "tut_agent", Title: "エージェント", Pages: []string{"説明"}, FeatureID: domain.FeatureAgentCustomization},
		{ID: "tut_chain", Title: "チェイン", Pages: []string{"説明"}, FeatureID: domain.FeatureChainEffect},
	}

	// 旧セーブ: FeatureUnlockフィールドなし、CurrentRank=3
	// → ランク3以下の機能はUnlocked扱い
	currentRank := 3
	snap := app.RebuildFeatureUnlockFromRank(rules, currentRank)

	mgr, _ := unlocking.NewManager(rules, tutorials, snap)

	if !mgr.IsUnlocked(domain.FeatureDefenseSkill) {
		t.Error("旧セーブ復元: defense_skill（ランク2）がUnlockedでない")
	}
	if !mgr.IsUnlocked(domain.FeatureAgentCustomization) {
		t.Error("旧セーブ復元: agent_customization（ランク3）がUnlockedでない")
	}
	if mgr.IsUnlocked(domain.FeatureChainEffect) {
		t.Error("旧セーブ復元: chain_effect（ランク4）がUnlockedになっている（期待: Locked）")
	}
}

// === 受け入れ基準 29: 未知のFeatureIDがセーブデータに含まれる場合、無視して保持 ===
func TestFeatureUnlock_UnknownFeatureIDPreserved(t *testing.T) {
	// 未知のFeatureIDを含むセーブデータ
	savedData := savedata.FeatureUnlockSave{
		Features: map[string]string{
			string(domain.FeatureDefenseSkill): "unlocked",
			"future_feature":                   "unlocked",
		},
	}

	snap := app.FeatureUnlockSaveToSnapshot(savedData)

	// 未知IDの状態は保持される
	if snap.GetStatus("future_feature") != domain.FeatureUnlocked {
		t.Error("未知のFeatureIDがセーブデータから保持されなかった")
	}
}

// === 受け入れ基準 6,7 統合: ランク1→4まで一気にApplyし、その後再度ApplyRank(4)で冪等 ===
func TestFeatureUnlock_FullFlowWithIdempotency(t *testing.T) {
	rules := []domain.UnlockRule{
		{FeatureID: domain.FeatureDefenseSkill, UnlockRank: 2, TutorialID: "tut_defense"},
		{FeatureID: domain.FeatureAgentCustomization, UnlockRank: 3, TutorialID: "tut_agent"},
		{FeatureID: domain.FeatureChainEffect, UnlockRank: 4, TutorialID: "tut_chain"},
		{FeatureID: domain.FeatureManaSystem, UnlockRank: 5, TutorialID: "tut_mana"},
	}
	tutorials := []domain.TutorialDef{
		{ID: "tut_defense", Title: "ディフェンス", Pages: []string{"説明"}, FeatureID: domain.FeatureDefenseSkill},
		{ID: "tut_agent", Title: "エージェント", Pages: []string{"説明"}, FeatureID: domain.FeatureAgentCustomization},
		{ID: "tut_chain", Title: "チェイン", Pages: []string{"説明"}, FeatureID: domain.FeatureChainEffect},
		{ID: "tut_mana", Title: "マナ", Pages: []string{"説明"}, FeatureID: domain.FeatureManaSystem},
	}

	state := domain.NewFeatureUnlockState()
	mgr, _ := unlocking.NewManager(rules, tutorials, state)

	// ランク4まで一気に
	delta, _ := mgr.ApplyRank(4)
	if len(delta.NewPendingFeatures) != 3 {
		t.Errorf("ApplyRank(4): Pending数 = %d, 期待: 3", len(delta.NewPendingFeatures))
	}

	// チュートリアルを順次完了
	for i := 0; i < 3; i++ {
		tutID, ok := mgr.NextPendingTutorial()
		if !ok {
			t.Fatalf("チュートリアル%d: NextPendingTutorialがfalse", i+1)
		}
		mgr.CompleteTutorial(tutID)
	}

	// 再度ApplyRank(4) → 冪等
	delta2, _ := mgr.ApplyRank(4)
	if len(delta2.NewPendingFeatures) != 0 {
		t.Errorf("再ApplyRank(4): Pending数 = %d, 期待: 0", len(delta2.NewPendingFeatures))
	}

	// ランク5到達
	delta3, _ := mgr.ApplyRank(5)
	if len(delta3.NewPendingFeatures) != 1 {
		t.Errorf("ApplyRank(5): Pending数 = %d, 期待: 1", len(delta3.NewPendingFeatures))
	}
}
