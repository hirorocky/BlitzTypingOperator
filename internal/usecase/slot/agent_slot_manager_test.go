// Package slot はエージェントスロット管理のユースケースを提供します。
// このファイルはAgentSlotManagerの単体テストを定義します。

package slot

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ==================== テスト用ヘルパー ====================

// テスト用のCoreTypeを作成
func createTestCoreType(id string, tags []string) domain.CoreType {
	return domain.CoreType{
		ID:             id,
		Name:           "テストコア",
		StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		PassiveSkillID: "passive_001",
		AllowedTags:    tags,
		MinDropLevel:   1,
	}
}

// テスト用のSkillTypeを作成
func createTestSkillType(id string, tags []string) domain.SkillType {
	return domain.SkillType{
		ID:              id,
		Name:            "テストスキル",
		Icon:            "•",
		Tags:            tags,
		Description:     "テスト用スキル",
		CooldownSeconds: 5.0,
		Difficulty:      1,
		MinDropLevel:    1,
		Effects:         []domain.SkillEffect{},
	}
}

// テスト用のPassiveSkillを作成
func createTestPassiveSkill(id string) domain.PassiveSkill {
	return domain.PassiveSkill{
		ID:          id,
		Name:        "テストパッシブ",
		Description: "テスト用パッシブスキル",
	}
}

// テスト用のインベントリとマスタデータを設定したAgentSlotManagerを作成
func createTestManager() (*AgentSlotManager, *domain.CoreInventory, *domain.SkillInventory) {
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	coreTypes := map[string]domain.CoreType{
		"core_001": createTestCoreType("core_001", []string{"physical", "magic"}),
		"core_002": createTestCoreType("core_002", []string{"physical"}),
		"core_003": createTestCoreType("core_003", []string{"heal"}),
	}

	skillTypes := map[string]domain.SkillType{
		"skill_001": createTestSkillType("skill_001", []string{"physical"}),
		"skill_002": createTestSkillType("skill_002", []string{"magic"}),
		"skill_003": createTestSkillType("skill_003", []string{"heal"}),
		"skill_004": createTestSkillType("skill_004", []string{"physical", "magic"}),
	}

	passiveSkills := map[string]domain.PassiveSkill{
		"passive_001": createTestPassiveSkill("passive_001"),
	}

	manager := NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills)
	return manager, coreInv, skillInv
}

// ==================== 5.1 基本構造のテスト ====================

func TestNewAgentSlotManager(t *testing.T) {
	manager, _, _ := createTestManager()

	if manager == nil {
		t.Fatal("NewAgentSlotManagerはnilを返すべきではない")
	}

	slots := manager.GetSlots()
	if len(slots) != MaxAgentSlotCount {
		t.Errorf("スロット数 = %d, want %d", len(slots), MaxAgentSlotCount)
	}

	for i, slot := range slots {
		if slot == nil {
			t.Errorf("slots[%d]はnilであるべきではない", i)
		}
		if !slot.IsEmpty() {
			t.Errorf("新規作成されたslots[%d]は空であるべき", i)
		}
	}
}

func TestAgentSlotManager_GetSlot(t *testing.T) {
	manager, _, _ := createTestManager()

	tests := []struct {
		name    string
		slot    int
		wantNil bool
	}{
		{
			name:    "スロット0を取得",
			slot:    0,
			wantNil: false,
		},
		{
			name:    "スロット1を取得",
			slot:    1,
			wantNil: false,
		},
		{
			name:    "スロット2を取得",
			slot:    2,
			wantNil: false,
		},
		{
			name:    "範囲外（負）",
			slot:    -1,
			wantNil: true,
		},
		{
			name:    "範囲外（上限超え）",
			slot:    3,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := manager.GetSlot(tt.slot)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetSlot(%d)はnilを返すべき", tt.slot)
				}
			} else {
				if got == nil {
					t.Errorf("GetSlot(%d)はnilを返すべきではない", tt.slot)
				}
			}
		})
	}
}

func TestAgentSlotManager_GetSlots_ReturnsAllSlots(t *testing.T) {
	manager, _, _ := createTestManager()

	slots := manager.GetSlots()

	if len(slots) != MaxAgentSlotCount {
		t.Errorf("GetSlots()の長さ = %d, want %d", len(slots), MaxAgentSlotCount)
	}

	// 各スロットがGetSlotと同じ参照を返すことを確認
	for i := 0; i < MaxAgentSlotCount; i++ {
		if slots[i] != manager.GetSlot(i) {
			t.Errorf("slots[%d]はGetSlot(%d)と同じ参照であるべき", i, i)
		}
	}
}

// ==================== 5.2 コア付け替え機能のテスト ====================

func TestAgentSlotManager_SetCore_Success(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加（最大レベル10）
	coreInv.AddCore("core_001", 10)

	// スロット0にコアを設定
	err := manager.SetCore(0, "core_001", 5)

	if err != nil {
		t.Fatalf("SetCoreはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	if slot.CoreTypeID != "core_001" {
		t.Errorf("CoreTypeID = %q, want %q", slot.CoreTypeID, "core_001")
	}
	if slot.CoreLevel != 5 {
		t.Errorf("CoreLevel = %d, want %d", slot.CoreLevel, 5)
	}
	if slot.IsEmpty() {
		t.Error("コア設定後のスロットは空であるべきではない")
	}
}

func TestAgentSlotManager_SetCore_MaxLevel(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加（最大レベル10）
	coreInv.AddCore("core_001", 10)

	// 最大レベルで設定
	err := manager.SetCore(0, "core_001", 10)

	if err != nil {
		t.Fatalf("最大レベルでのSetCoreはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	if slot.CoreLevel != 10 {
		t.Errorf("CoreLevel = %d, want %d", slot.CoreLevel, 10)
	}
}

func TestAgentSlotManager_SetCore_LevelOne(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加
	coreInv.AddCore("core_001", 10)

	// レベル1で設定
	err := manager.SetCore(0, "core_001", 1)

	if err != nil {
		t.Fatalf("レベル1でのSetCoreはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	if slot.CoreLevel != 1 {
		t.Errorf("CoreLevel = %d, want %d", slot.CoreLevel, 1)
	}
}

func TestAgentSlotManager_SetCore_NotOwned(t *testing.T) {
	manager, _, _ := createTestManager()

	// インベントリにコアを追加せずに設定を試みる
	err := manager.SetCore(0, "core_001", 5)

	if err == nil {
		t.Fatal("未保有コアでのSetCoreはエラーを返すべき")
	}
	if err != ErrCoreNotOwned {
		t.Errorf("err = %v, want %v", err, ErrCoreNotOwned)
	}

	slot := manager.GetSlot(0)
	if !slot.IsEmpty() {
		t.Error("エラー時にスロットは変更されるべきではない")
	}
}

func TestAgentSlotManager_SetCore_LevelTooHigh(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加（最大レベル10）
	coreInv.AddCore("core_001", 10)

	// 最大レベルより高いレベルで設定
	err := manager.SetCore(0, "core_001", 11)

	if err == nil {
		t.Fatal("最大レベルより高いレベルでのSetCoreはエラーを返すべき")
	}
	if err != ErrLevelOutOfRange {
		t.Errorf("err = %v, want %v", err, ErrLevelOutOfRange)
	}
}

func TestAgentSlotManager_SetCore_LevelTooLow(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加
	coreInv.AddCore("core_001", 10)

	// レベル0で設定
	err := manager.SetCore(0, "core_001", 0)

	if err == nil {
		t.Fatal("レベル0でのSetCoreはエラーを返すべき")
	}
	if err != ErrLevelOutOfRange {
		t.Errorf("err = %v, want %v", err, ErrLevelOutOfRange)
	}
}

func TestAgentSlotManager_SetCore_InvalidSlotIndex(t *testing.T) {
	manager, coreInv, _ := createTestManager()
	coreInv.AddCore("core_001", 10)

	tests := []struct {
		name string
		slot int
	}{
		{"負のインデックス", -1},
		{"範囲外のインデックス", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetCore(tt.slot, "core_001", 5)
			if err == nil {
				t.Fatal("無効なスロットインデックスでのSetCoreはエラーを返すべき")
			}
			if err != ErrSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSlotIndexOutOfRange)
			}
		})
	}
}

func TestAgentSlotManager_SetCore_SameCoreInMultipleSlots(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアをインベントリに追加
	coreInv.AddCore("core_001", 10)

	// 同一コアを複数スロットに設定
	err1 := manager.SetCore(0, "core_001", 5)
	err2 := manager.SetCore(1, "core_001", 7)
	err3 := manager.SetCore(2, "core_001", 10)

	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("同一コアを複数スロットに設定できるべき: err1=%v, err2=%v, err3=%v", err1, err2, err3)
	}

	// 各スロットが独立していることを確認
	if manager.GetSlot(0).CoreLevel != 5 {
		t.Errorf("slot0.CoreLevel = %d, want %d", manager.GetSlot(0).CoreLevel, 5)
	}
	if manager.GetSlot(1).CoreLevel != 7 {
		t.Errorf("slot1.CoreLevel = %d, want %d", manager.GetSlot(1).CoreLevel, 7)
	}
	if manager.GetSlot(2).CoreLevel != 10 {
		t.Errorf("slot2.CoreLevel = %d, want %d", manager.GetSlot(2).CoreLevel, 10)
	}
}

func TestAgentSlotManager_SetCore_RemovesIncompatibleSkills(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// core_001 は physical, magic を許可
	coreInv.AddCore("core_001", 10)
	// core_003 は heal のみを許可
	coreInv.AddCore("core_003", 10)

	// スキルをインベントリに追加
	skillInv.AddSkill("skill_001", "") // physical タグ
	skillInv.AddSkill("skill_002", "") // magic タグ
	skillInv.AddSkill("skill_003", "") // heal タグ

	// まずcore_001を設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// physicalとmagicスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("skill_001の設定でエラー: %v", err)
	}
	err = manager.SetSkill(0, 1, "skill_002", "")
	if err != nil {
		t.Fatalf("skill_002の設定でエラー: %v", err)
	}

	if manager.GetSlot(0).GetSkillCount() != 2 {
		t.Errorf("スキル設定後のスキル数 = %d, want %d", manager.GetSlot(0).GetSkillCount(), 2)
	}

	// core_003に変更（healのみ許可）
	err = manager.SetCore(0, "core_003", 5)
	if err != nil {
		t.Fatalf("コア変更でエラー: %v", err)
	}

	// 互換性のないスキル（physical, magic）は削除されるべき
	slot := manager.GetSlot(0)
	if slot.GetSkillCount() != 0 {
		t.Errorf("互換性のないスキルが削除されていない: スキル数 = %d, want %d", slot.GetSkillCount(), 0)
	}
}

func TestAgentSlotManager_ClearCore(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")

	// コアとスキルを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("SetSkillでエラー: %v", err)
	}

	// コアをクリア
	err = manager.ClearCore(0)
	if err != nil {
		t.Fatalf("ClearCoreはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	if !slot.IsEmpty() {
		t.Error("ClearCore後のスロットは空であるべき")
	}
	if slot.GetSkillCount() != 0 {
		t.Errorf("ClearCore後のスキル数 = %d, want %d", slot.GetSkillCount(), 0)
	}
}

func TestAgentSlotManager_ClearCore_InvalidSlotIndex(t *testing.T) {
	manager, _, _ := createTestManager()

	tests := []struct {
		name string
		slot int
	}{
		{"負のインデックス", -1},
		{"範囲外のインデックス", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ClearCore(tt.slot)
			if err == nil {
				t.Fatal("無効なスロットインデックスでのClearCoreはエラーを返すべき")
			}
			if err != ErrSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSlotIndexOutOfRange)
			}
		})
	}
}
