// Package app は BlitzTypingOperator TUIゲームのRootModelを提供します。
package app

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/savedata"
)

// FeatureUnlockSnapshotToSave はFeatureUnlockStateをセーブデータ形式に変換します。
func FeatureUnlockSnapshotToSave(snap domain.FeatureUnlockState) savedata.FeatureUnlockSave {
	features := snap.AllFeatures()
	featureStrings := make(map[string]string, len(features))
	for id, status := range features {
		featureStrings[string(id)] = featureStatusToString(status)
	}
	return savedata.FeatureUnlockSave{
		Features: featureStrings,
	}
}

// FeatureUnlockSaveToSnapshot はセーブデータからFeatureUnlockStateを復元します。
func FeatureUnlockSaveToSnapshot(save savedata.FeatureUnlockSave) domain.FeatureUnlockState {
	features := make(map[domain.FeatureID]domain.FeatureStatus, len(save.Features))
	for id, statusStr := range save.Features {
		features[domain.FeatureID(id)] = stringToFeatureStatus(statusStr)
	}
	return domain.NewFeatureUnlockStateFrom(features)
}

// RebuildFeatureUnlockFromRank は旧セーブデータ互換用に、ランクから解放状態を再構築します。
// 該当ランク以下の機能はUnlocked扱いにします。
func RebuildFeatureUnlockFromRank(rules []domain.UnlockRule, currentRank int) domain.FeatureUnlockState {
	features := make(map[domain.FeatureID]domain.FeatureStatus)
	for _, rule := range rules {
		if rule.UnlockRank <= currentRank {
			features[rule.FeatureID] = domain.FeatureUnlocked
		}
	}
	return domain.NewFeatureUnlockStateFrom(features)
}

func featureStatusToString(s domain.FeatureStatus) string {
	switch s {
	case domain.FeatureLocked:
		return "locked"
	case domain.FeaturePendingTutorial:
		return "pending_tutorial"
	case domain.FeatureUnlocked:
		return "unlocked"
	default:
		return "locked"
	}
}

func stringToFeatureStatus(s string) domain.FeatureStatus {
	switch s {
	case "locked":
		return domain.FeatureLocked
	case "pending_tutorial":
		return domain.FeaturePendingTutorial
	case "unlocked":
		return domain.FeatureUnlocked
	default:
		return domain.FeatureLocked
	}
}
