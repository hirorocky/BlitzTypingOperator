// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"testing"
)

// ==================== EnemyDefeatRecord テスト ====================

// TestEnemyDefeatRecord_新規作成 は新規作成されたEnemyDefeatRecordの初期状態を確認します。
func TestEnemyDefeatRecord_新規作成(t *testing.T) {
	record := EnemyDefeatRecord{}

	if record.Defeated {
		t.Error("新規作成されたレコードは未撃破であるべき")
	}
	if record.MaxDefeatedLevel != 0 {
		t.Errorf("新規作成されたレコードのMaxDefeatedLevelは0であるべき: got %d", record.MaxDefeatedLevel)
	}
}

// TestEnemyDefeatRecord_撃破記録 は撃破記録が正しく設定されることを確認します。
func TestEnemyDefeatRecord_撃破記録(t *testing.T) {
	record := EnemyDefeatRecord{
		Defeated:         true,
		MaxDefeatedLevel: 5,
	}

	if !record.Defeated {
		t.Error("撃破済みレコードはDefeated=trueであるべき")
	}
	if record.MaxDefeatedLevel != 5 {
		t.Errorf("MaxDefeatedLevel = %d, want 5", record.MaxDefeatedLevel)
	}
}

// ==================== EnemyProgress テスト ====================

// TestNewEnemyProgress_新規作成 はNewEnemyProgressで初期状態が正しく設定されることを確認します。
func TestNewEnemyProgress_新規作成(t *testing.T) {
	progress := NewEnemyProgress()

	if progress == nil {
		t.Fatal("NewEnemyProgressはnilを返すべきではない")
	}
	if progress.CurrentRank != 1 {
		t.Errorf("初期CurrentRank = %d, want 1", progress.CurrentRank)
	}
	if progress.DefeatRecords == nil {
		t.Error("DefeatRecordsはnilであるべきではない")
	}
	if len(progress.DefeatRecords) != 0 {
		t.Errorf("初期DefeatRecordsは空であるべき: len = %d", len(progress.DefeatRecords))
	}
}

// TestEnemyProgress_IsDefeated はIsDefeatedが撃破状況を正しく返すことを確認します。
func TestEnemyProgress_IsDefeated(t *testing.T) {
	progress := NewEnemyProgress()

	// 未撃破の敵
	if progress.IsDefeated("slime") {
		t.Error("未撃破の敵に対してIsDefeatedはfalseを返すべき")
	}

	// 撃破を記録
	progress.DefeatRecords["slime"] = EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 1}

	if !progress.IsDefeated("slime") {
		t.Error("撃破済みの敵に対してIsDefeatedはtrueを返すべき")
	}
}

// TestEnemyProgress_GetMaxDefeatedLevel はGetMaxDefeatedLevelが最大撃破レベルを正しく返すことを確認します。
func TestEnemyProgress_GetMaxDefeatedLevel(t *testing.T) {
	progress := NewEnemyProgress()

	// 未撃破の敵は0を返す
	if level := progress.GetMaxDefeatedLevel("slime"); level != 0 {
		t.Errorf("未撃破の敵のMaxDefeatedLevel = %d, want 0", level)
	}

	// 撃破を記録
	progress.DefeatRecords["slime"] = EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 10}

	if level := progress.GetMaxDefeatedLevel("slime"); level != 10 {
		t.Errorf("撃破済みの敵のMaxDefeatedLevel = %d, want 10", level)
	}
}

// TestEnemyProgress_GetSelectableLevelRange は選択可能レベル範囲を正しく計算することを確認します。
func TestEnemyProgress_GetSelectableLevelRange(t *testing.T) {
	progress := NewEnemyProgress()
	defaultLevel := 1

	tests := []struct {
		name        string
		enemyTypeID string
		setupRecord *EnemyDefeatRecord
		wantMin     int
		wantMax     int
	}{
		{
			name:        "未撃破の敵はデフォルトレベルのみ",
			enemyTypeID: "slime",
			setupRecord: nil,
			wantMin:     defaultLevel,
			wantMax:     defaultLevel,
		},
		{
			name:        "初撃破後はレベル1〜6が選択可能（撃破レベル1+5）",
			enemyTypeID: "goblin",
			setupRecord: &EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 1},
			wantMin:     1,
			wantMax:     6,
		},
		{
			name:        "レベル5撃破後はレベル1〜10が選択可能",
			enemyTypeID: "orc",
			setupRecord: &EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 5},
			wantMin:     1,
			wantMax:     10,
		},
		{
			name:        "レベル10撃破後はレベル1〜15が選択可能",
			enemyTypeID: "dragon",
			setupRecord: &EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 10},
			wantMin:     1,
			wantMax:     15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupRecord != nil {
				progress.DefeatRecords[tt.enemyTypeID] = *tt.setupRecord
			}

			min, max := progress.GetSelectableLevelRange(tt.enemyTypeID, defaultLevel)

			if min != tt.wantMin {
				t.Errorf("min = %d, want %d", min, tt.wantMin)
			}
			if max != tt.wantMax {
				t.Errorf("max = %d, want %d", max, tt.wantMax)
			}
		})
	}
}

// TestEnemyProgress_RecordDefeat_初撃破 は初撃破時のHP増加量が正しいことを確認します。
func TestEnemyProgress_RecordDefeat_初撃破(t *testing.T) {
	progress := NewEnemyProgress()

	hpGain, rankUnlocked := progress.RecordDefeat("slime", 1)

	// 初撃破: HP +10
	if hpGain != 10 {
		t.Errorf("初撃破のhpGain = %d, want 10", hpGain)
	}
	if !progress.IsDefeated("slime") {
		t.Error("撃破後はIsDefeated=trueであるべき")
	}
	if progress.GetMaxDefeatedLevel("slime") != 1 {
		t.Errorf("MaxDefeatedLevel = %d, want 1", progress.GetMaxDefeatedLevel("slime"))
	}
	// ランク解放は未テスト（敵タイプ情報が必要）
	_ = rankUnlocked
}

// TestEnemyProgress_RecordDefeat_高レベル撃破 は高レベル撃破時のHP増加量が正しいことを確認します。
func TestEnemyProgress_RecordDefeat_高レベル撃破(t *testing.T) {
	progress := NewEnemyProgress()

	// 初撃破
	progress.RecordDefeat("slime", 1)

	// レベル5で撃破: HP += (5-1) * 10 = 40
	hpGain, _ := progress.RecordDefeat("slime", 5)

	if hpGain != 40 {
		t.Errorf("高レベル撃破のhpGain = %d, want 40", hpGain)
	}
	if progress.GetMaxDefeatedLevel("slime") != 5 {
		t.Errorf("MaxDefeatedLevel = %d, want 5", progress.GetMaxDefeatedLevel("slime"))
	}
}

// TestEnemyProgress_RecordDefeat_低レベル再撃破 は低レベル再撃破時にHP増加がないことを確認します。
func TestEnemyProgress_RecordDefeat_低レベル再撃破(t *testing.T) {
	progress := NewEnemyProgress()

	// レベル5で初撃破
	progress.DefeatRecords["slime"] = EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 5}

	// レベル3で再撃破: HP増加なし
	hpGain, _ := progress.RecordDefeat("slime", 3)

	if hpGain != 0 {
		t.Errorf("低レベル再撃破のhpGain = %d, want 0", hpGain)
	}
	// MaxDefeatedLevelは変わらない
	if progress.GetMaxDefeatedLevel("slime") != 5 {
		t.Errorf("MaxDefeatedLevel = %d, want 5（変更されないべき）", progress.GetMaxDefeatedLevel("slime"))
	}
}

// TestEnemyProgress_RecordDefeat_同レベル再撃破 は同レベル再撃破時にHP増加がないことを確認します。
func TestEnemyProgress_RecordDefeat_同レベル再撃破(t *testing.T) {
	progress := NewEnemyProgress()

	// レベル5で初撃破
	progress.DefeatRecords["slime"] = EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 5}

	// レベル5で再撃破: HP増加なし
	hpGain, _ := progress.RecordDefeat("slime", 5)

	if hpGain != 0 {
		t.Errorf("同レベル再撃破のhpGain = %d, want 0", hpGain)
	}
}

// TestEnemyProgress_RecordDefeat_新記録更新 は新記録更新時のHP増加量が正しいことを確認します。
func TestEnemyProgress_RecordDefeat_新記録更新(t *testing.T) {
	progress := NewEnemyProgress()

	// レベル5で初撃破
	progress.DefeatRecords["slime"] = EnemyDefeatRecord{Defeated: true, MaxDefeatedLevel: 5}

	// レベル10で撃破: HP += (10-5) * 10 = 50
	hpGain, _ := progress.RecordDefeat("slime", 10)

	if hpGain != 50 {
		t.Errorf("新記録更新のhpGain = %d, want 50", hpGain)
	}
	if progress.GetMaxDefeatedLevel("slime") != 10 {
		t.Errorf("MaxDefeatedLevel = %d, want 10", progress.GetMaxDefeatedLevel("slime"))
	}
}
