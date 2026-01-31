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

	// ChainEffectsは空マップで渡す（テスト用）
	chainEffects := map[string]domain.ChainEffect{}

	manager := NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills, chainEffects)
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
	for i := range MaxAgentSlotCount {
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

// ==================== 5.3 スキル付け替え機能のテスト ====================

func TestAgentSlotManager_SetSkill_Success(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)             // physical, magic を許可
	skillInv.AddSkill("skill_001", "chain_001") // physical タグ

	// コアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// スキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "chain_001")
	if err != nil {
		t.Fatalf("SetSkillはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	skillConfig := slot.GetSkill(0)
	if skillConfig.TypeID != "skill_001" {
		t.Errorf("TypeID = %q, want %q", skillConfig.TypeID, "skill_001")
	}
	if skillConfig.ChainEffectID != "chain_001" {
		t.Errorf("ChainEffectID = %q, want %q", skillConfig.ChainEffectID, "chain_001")
	}
}

func TestAgentSlotManager_SetSkill_WithoutChainEffect(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加（チェイン効果なし）
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "") // チェイン効果なし

	// コアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// チェイン効果なしでスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("チェイン効果なしでのSetSkillはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	skillConfig := slot.GetSkill(0)
	if skillConfig.TypeID != "skill_001" {
		t.Errorf("TypeID = %q, want %q", skillConfig.TypeID, "skill_001")
	}
	if skillConfig.ChainEffectID != "" {
		t.Errorf("ChainEffectID = %q, want %q", skillConfig.ChainEffectID, "")
	}
}

func TestAgentSlotManager_SetSkill_CoreNotSet(t *testing.T) {
	manager, _, skillInv := createTestManager()

	// スキルをインベントリに追加
	skillInv.AddSkill("skill_001", "")

	// コアを設定せずにスキル設定を試みる
	err := manager.SetSkill(0, 0, "skill_001", "")

	if err == nil {
		t.Fatal("コア未設定でのSetSkillはエラーを返すべき")
	}
	if err != ErrCoreNotSet {
		t.Errorf("err = %v, want %v", err, ErrCoreNotSet)
	}
}

func TestAgentSlotManager_SetSkill_NotOwned(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// コアを設定
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// インベントリにスキルを追加せずに設定を試みる
	err = manager.SetSkill(0, 0, "skill_001", "")

	if err == nil {
		t.Fatal("未保有スキルでのSetSkillはエラーを返すべき")
	}
	if err != ErrSkillNotOwned {
		t.Errorf("err = %v, want %v", err, ErrSkillNotOwned)
	}
}

func TestAgentSlotManager_SetSkill_IncompatibleWithCore(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// core_002 は physical のみを許可
	coreInv.AddCore("core_002", 10)
	// skill_002 は magic タグのみ
	skillInv.AddSkill("skill_002", "")

	// コアを設定
	err := manager.SetCore(0, "core_002", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// 互換性のないスキルを設定しようとする
	err = manager.SetSkill(0, 0, "skill_002", "")

	if err == nil {
		t.Fatal("互換性のないスキルでのSetSkillはエラーを返すべき")
	}
	if err != ErrSkillIncompatible {
		t.Errorf("err = %v, want %v", err, ErrSkillIncompatible)
	}
}

func TestAgentSlotManager_SetSkill_ChainVariationNotOwned(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "chain_001") // chain_001 のみ保有

	// コアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// 未保有のチェイン効果バリエーションを指定
	err = manager.SetSkill(0, 0, "skill_001", "chain_002")

	if err == nil {
		t.Fatal("未保有チェイン効果でのSetSkillはエラーを返すべき")
	}
	if err != ErrChainVariationNotOwned {
		t.Errorf("err = %v, want %v", err, ErrChainVariationNotOwned)
	}
}

func TestAgentSlotManager_SetSkill_InvalidSlotIndex(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	tests := []struct {
		name string
		slot int
	}{
		{"負のインデックス", -1},
		{"範囲外のインデックス", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetSkill(tt.slot, 0, "skill_001", "")
			if err == nil {
				t.Fatal("無効なスロットインデックスでのSetSkillはエラーを返すべき")
			}
			if err != ErrSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSlotIndexOutOfRange)
			}
		})
	}
}

func TestAgentSlotManager_SetSkill_InvalidSkillSlotIndex(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	tests := []struct {
		name      string
		skillSlot int
	}{
		{"負のインデックス", -1},
		{"範囲外のインデックス", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.SetSkill(0, tt.skillSlot, "skill_001", "")
			if err == nil {
				t.Fatal("無効なスキルスロットインデックスでのSetSkillはエラーを返すべき")
			}
			if err != ErrSkillSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSkillSlotIndexOutOfRange)
			}
		})
	}
}

func TestAgentSlotManager_SetSkill_SameSkillInDifferentAgents_ShouldFail(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")

	// 2つのスロットにコアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCore(0)でエラー: %v", err)
	}
	err = manager.SetCore(1, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCore(1)でエラー: %v", err)
	}

	// 最初のスロットにスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("スロット0へのSetSkillでエラー: %v", err)
	}

	// 同じスキルを別のエージェントに設定しようとするとエラーになるべき
	err = manager.SetSkill(1, 0, "skill_001", "")
	if err == nil {
		t.Fatal("同じスキルを異なるエージェントに設定するとエラーを返すべき")
	}
	if err != ErrSkillAlreadyEquipped {
		t.Errorf("err = %v, want %v", err, ErrSkillAlreadyEquipped)
	}

	// 最初のスロットのスキルは設定されているまま
	if manager.GetSlot(0).GetSkill(0).TypeID != "skill_001" {
		t.Error("スロット0にskill_001が設定されているべき")
	}
	// 2番目のスロットにはスキルが設定されていない
	if !manager.GetSlot(1).GetSkill(0).IsEmpty() {
		t.Error("スロット1のスキルは空であるべき")
	}
}

func TestAgentSlotManager_SetSkill_SameSkillInSameAgentMultipleSlots_ShouldFail(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")

	// コアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// 最初のスキルスロットにスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("スキルスロット0へのSetSkillでエラー: %v", err)
	}

	// 同じスキルを同一エージェント内の別スロットに設定しようとするとエラーになるべき
	err = manager.SetSkill(0, 1, "skill_001", "")
	if err == nil {
		t.Fatal("同じスキルを同一エージェント内の複数スロットに設定するとエラーを返すべき")
	}
	if err != ErrSkillAlreadyEquipped {
		t.Errorf("err = %v, want %v", err, ErrSkillAlreadyEquipped)
	}

	// 最初のスキルスロットにはスキルが設定されている
	slot := manager.GetSlot(0)
	if slot.GetSkill(0).TypeID != "skill_001" {
		t.Error("スキルスロット0にskill_001が設定されているべき")
	}
	// 2番目のスキルスロットは空
	if !slot.GetSkill(1).IsEmpty() {
		t.Error("スキルスロット1は空であるべき")
	}
	if slot.GetSkillCount() != 1 {
		t.Errorf("スキル数 = %d, want %d", slot.GetSkillCount(), 1)
	}
}

func TestAgentSlotManager_SetSkill_ReplaceInSameSlot_Success(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "") // physical
	skillInv.AddSkill("skill_002", "") // magic

	// コアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// 最初のスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("skill_001の設定でエラー: %v", err)
	}

	// 同じスロットに別のスキルを設定（上書き）
	err = manager.SetSkill(0, 0, "skill_002", "")
	if err != nil {
		t.Fatalf("同じスロットへの別スキル設定でエラー: %v", err)
	}

	// skill_002が設定されていることを確認
	slot := manager.GetSlot(0)
	if slot.GetSkill(0).TypeID != "skill_002" {
		t.Errorf("TypeID = %q, want %q", slot.GetSkill(0).TypeID, "skill_002")
	}
}

func TestAgentSlotManager_SetSkill_AfterClearCanEquipToOtherSlot(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")

	// 2つのスロットにコアを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCore(0)でエラー: %v", err)
	}
	err = manager.SetCore(1, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCore(1)でエラー: %v", err)
	}

	// スロット0にスキルを設定
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("スロット0へのSetSkillでエラー: %v", err)
	}

	// スロット0のスキルをクリア
	err = manager.ClearSkill(0, 0)
	if err != nil {
		t.Fatalf("ClearSkillでエラー: %v", err)
	}

	// スロット1に同じスキルを設定できるはず
	err = manager.SetSkill(1, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("クリア後に別スロットへのSetSkillでエラー: %v", err)
	}

	// スロット1にスキルが設定されていることを確認
	if manager.GetSlot(1).GetSkill(0).TypeID != "skill_001" {
		t.Error("スロット1にskill_001が設定されているべき")
	}
}

func TestAgentSlotManager_ClearSkill(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")
	skillInv.AddSkill("skill_004", "") // physical, magic タグ

	// コアとスキルを設定
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("SetSkillでエラー: %v", err)
	}
	err = manager.SetSkill(0, 1, "skill_004", "")
	if err != nil {
		t.Fatalf("SetSkill(skill_004)でエラー: %v", err)
	}

	// スキルをクリア
	err = manager.ClearSkill(0, 0)
	if err != nil {
		t.Fatalf("ClearSkillはエラーを返すべきではない: %v", err)
	}

	slot := manager.GetSlot(0)
	if !slot.GetSkill(0).IsEmpty() {
		t.Error("ClearSkill後のスキルスロットは空であるべき")
	}
	// 他のスキルスロットは影響を受けない
	if slot.GetSkill(1).IsEmpty() {
		t.Error("ClearSkillは他のスキルスロットに影響しないべき")
	}
	if slot.GetSkillCount() != 1 {
		t.Errorf("ClearSkill後のスキル数 = %d, want %d", slot.GetSkillCount(), 1)
	}
	// コア設定は変更されない
	if slot.IsEmpty() {
		t.Error("ClearSkillでスロット全体が空になるべきではない")
	}
}

func TestAgentSlotManager_ClearSkill_InvalidSlotIndex(t *testing.T) {
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
			err := manager.ClearSkill(tt.slot, 0)
			if err == nil {
				t.Fatal("無効なスロットインデックスでのClearSkillはエラーを返すべき")
			}
			if err != ErrSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSlotIndexOutOfRange)
			}
		})
	}
}

func TestAgentSlotManager_ClearSkill_InvalidSkillSlotIndex(t *testing.T) {
	manager, coreInv, _ := createTestManager()
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	tests := []struct {
		name      string
		skillSlot int
	}{
		{"負のインデックス", -1},
		{"範囲外のインデックス", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.ClearSkill(0, tt.skillSlot)
			if err == nil {
				t.Fatal("無効なスキルスロットインデックスでのClearSkillはエラーを返すべき")
			}
			if err != ErrSkillSlotIndexOutOfRange {
				t.Errorf("err = %v, want %v", err, ErrSkillSlotIndexOutOfRange)
			}
		})
	}
}

// ==================== 5.4 バトル連携機能のテスト ====================

func TestAgentSlotManager_IsSlotReady(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// 空スロットは使用不可
	if manager.IsSlotReady(0) {
		t.Error("空スロットはバトルに使用不可であるべき")
	}

	// コアを設定したスロットは使用可能
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	if !manager.IsSlotReady(0) {
		t.Error("コア設定済みスロットはバトルに使用可能であるべき")
	}

	// 他のスロットは依然として使用不可
	if manager.IsSlotReady(1) {
		t.Error("コア未設定のスロット1は使用不可であるべき")
	}
}

func TestAgentSlotManager_IsSlotReady_InvalidIndex(t *testing.T) {
	manager, _, _ := createTestManager()

	// 範囲外のインデックスは常にfalse
	if manager.IsSlotReady(-1) {
		t.Error("負のインデックスはfalseを返すべき")
	}
	if manager.IsSlotReady(3) {
		t.Error("範囲外のインデックスはfalseを返すべき")
	}
}

func TestAgentSlotManager_GetReadySlotCount(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// 全スロット空の場合
	if manager.GetReadySlotCount() != 0 {
		t.Errorf("全スロット空の場合のReadySlotCount = %d, want %d", manager.GetReadySlotCount(), 0)
	}

	// 1つのスロットにコアを設定
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	if manager.GetReadySlotCount() != 1 {
		t.Errorf("1スロット設定後のReadySlotCount = %d, want %d", manager.GetReadySlotCount(), 1)
	}

	// 2つ目のスロットにコアを設定
	err = manager.SetCore(1, "core_001", 7)
	if err != nil {
		t.Fatalf("SetCore(1)でエラー: %v", err)
	}

	if manager.GetReadySlotCount() != 2 {
		t.Errorf("2スロット設定後のReadySlotCount = %d, want %d", manager.GetReadySlotCount(), 2)
	}

	// 全スロットにコアを設定
	err = manager.SetCore(2, "core_001", 10)
	if err != nil {
		t.Fatalf("SetCore(2)でエラー: %v", err)
	}

	if manager.GetReadySlotCount() != 3 {
		t.Errorf("全スロット設定後のReadySlotCount = %d, want %d", manager.GetReadySlotCount(), 3)
	}
}

func TestAgentSlotManager_ValidateSkillCompatibility(t *testing.T) {
	manager, coreInv, _ := createTestManager()

	// core_001 は physical, magic を許可
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// skill_001 は physical タグ → 互換性あり
	if !manager.ValidateSkillCompatibility(0, "skill_001") {
		t.Error("physical タグのスキルはcore_001と互換性があるべき")
	}

	// skill_002 は magic タグ → 互換性あり
	if !manager.ValidateSkillCompatibility(0, "skill_002") {
		t.Error("magic タグのスキルはcore_001と互換性があるべき")
	}

	// skill_003 は heal タグ → 互換性なし
	if manager.ValidateSkillCompatibility(0, "skill_003") {
		t.Error("heal タグのスキルはcore_001と互換性がないべき")
	}
}

func TestAgentSlotManager_ValidateSkillCompatibility_EmptySlot(t *testing.T) {
	manager, _, _ := createTestManager()

	// 空スロットでの互換性チェックは常にfalse
	if manager.ValidateSkillCompatibility(0, "skill_001") {
		t.Error("空スロットでの互換性チェックはfalseを返すべき")
	}
}

func TestAgentSlotManager_GetCompatibleSkills(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// core_001 は physical, magic を許可
	coreInv.AddCore("core_001", 10)
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}

	// スキルをインベントリに追加
	skillInv.AddSkill("skill_001", "") // physical
	skillInv.AddSkill("skill_002", "") // magic
	skillInv.AddSkill("skill_003", "") // heal
	skillInv.AddSkill("skill_004", "") // physical, magic

	// 互換性のあるスキル一覧を取得
	compatibleSkills := manager.GetCompatibleSkills(0)

	// skill_001, skill_002, skill_004 が互換性あり
	if len(compatibleSkills) != 3 {
		t.Errorf("互換スキル数 = %d, want %d", len(compatibleSkills), 3)
	}

	// 含まれるべきスキルを確認
	skillSet := make(map[string]bool)
	for _, s := range compatibleSkills {
		skillSet[s] = true
	}

	if !skillSet["skill_001"] {
		t.Error("skill_001が互換スキルに含まれるべき")
	}
	if !skillSet["skill_002"] {
		t.Error("skill_002が互換スキルに含まれるべき")
	}
	if !skillSet["skill_004"] {
		t.Error("skill_004が互換スキルに含まれるべき")
	}
	if skillSet["skill_003"] {
		t.Error("skill_003は互換スキルに含まれないべき")
	}
}

func TestAgentSlotManager_GetCompatibleSkills_EmptySlot(t *testing.T) {
	manager, _, skillInv := createTestManager()

	// スキルをインベントリに追加
	skillInv.AddSkill("skill_001", "")

	// 空スロットでは空のリストを返す
	compatibleSkills := manager.GetCompatibleSkills(0)
	if len(compatibleSkills) != 0 {
		t.Errorf("空スロットでの互換スキル数 = %d, want %d", len(compatibleSkills), 0)
	}
}

func TestAgentSlotManager_BuildAgentsForBattle(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	coreInv.AddCore("core_002", 5)
	skillInv.AddSkill("skill_001", "")
	skillInv.AddSkill("skill_004", "") // physical, magic タグ

	// スロット0: core_001 + skill_001, skill_004
	err := manager.SetCore(0, "core_001", 5)
	if err != nil {
		t.Fatalf("SetCore(0)でエラー: %v", err)
	}
	err = manager.SetSkill(0, 0, "skill_001", "")
	if err != nil {
		t.Fatalf("SetSkill(0,0)でエラー: %v", err)
	}
	err = manager.SetSkill(0, 1, "skill_004", "")
	if err != nil {
		t.Fatalf("SetSkill(0,1)でエラー: %v", err)
	}

	// スロット1: 空

	// スロット2: core_002（スキルなし）
	err = manager.SetCore(2, "core_002", 3)
	if err != nil {
		t.Fatalf("SetCore(2)でエラー: %v", err)
	}

	// バトル用エージェントを構築
	agents := manager.BuildAgentsForBattle()

	// 空スロットは含まれないので2体
	if len(agents) != 2 {
		t.Errorf("エージェント数 = %d, want %d", len(agents), 2)
	}

	// 最初のエージェント（スロット0）を確認
	agent0 := agents[0]
	if agent0.Core == nil {
		t.Fatal("エージェント0のCoreがnilであるべきではない")
	}
	if agent0.Core.TypeID != "core_001" {
		t.Errorf("エージェント0のCoreTypeID = %q, want %q", agent0.Core.TypeID, "core_001")
	}
	if agent0.Level != 5 {
		t.Errorf("エージェント0のLevel = %d, want %d", agent0.Level, 5)
	}
	if len(agent0.Modules) != 2 {
		t.Errorf("エージェント0のスキル数 = %d, want %d", len(agent0.Modules), 2)
	}

	// 2番目のエージェント（スロット2）を確認
	agent2 := agents[1]
	if agent2.Core.TypeID != "core_002" {
		t.Errorf("エージェント2のCoreTypeID = %q, want %q", agent2.Core.TypeID, "core_002")
	}
	if agent2.Level != 3 {
		t.Errorf("エージェント2のLevel = %d, want %d", agent2.Level, 3)
	}
	if len(agent2.Modules) != 0 {
		t.Errorf("エージェント2のスキル数 = %d, want %d", len(agent2.Modules), 0)
	}
}

func TestAgentSlotManager_BuildAgentsForBattle_AllEmpty(t *testing.T) {
	manager, _, _ := createTestManager()

	// 全スロット空の場合は空のスライス
	agents := manager.BuildAgentsForBattle()

	if agents == nil {
		t.Fatal("BuildAgentsForBattleはnilを返すべきではない")
	}
	if len(agents) != 0 {
		t.Errorf("全スロット空の場合のエージェント数 = %d, want %d", len(agents), 0)
	}
}

func TestAgentSlotManager_BuildAgentsForBattle_AllFull(t *testing.T) {
	manager, coreInv, skillInv := createTestManager()

	// コアとスキルをインベントリに追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")
	skillInv.AddSkill("skill_002", "") // magic
	skillInv.AddSkill("skill_004", "") // physical, magic

	// 各スロットには異なるスキルを設定
	skills := []string{"skill_001", "skill_002", "skill_004"}

	// 全スロットにコアを設定
	for i := range 3 {
		err := manager.SetCore(i, "core_001", 5+i)
		if err != nil {
			t.Fatalf("SetCore(%d)でエラー: %v", i, err)
		}
		err = manager.SetSkill(i, 0, skills[i], "")
		if err != nil {
			t.Fatalf("SetSkill(%d,0)でエラー: %v", i, err)
		}
	}

	// バトル用エージェントを構築
	agents := manager.BuildAgentsForBattle()

	if len(agents) != 3 {
		t.Errorf("全スロット設定時のエージェント数 = %d, want %d", len(agents), 3)
	}

	// 各エージェントのレベルを確認
	for i, agent := range agents {
		expectedLevel := 5 + i
		if agent.Level != expectedLevel {
			t.Errorf("エージェント%dのLevel = %d, want %d", i, agent.Level, expectedLevel)
		}
	}
}

// ==================== ChainEffect適用のテスト ====================

func TestAgentSlotManager_BuildAgentsForBattle_WithChainEffect(t *testing.T) {
	// ChainEffectsを含むマネージャーを作成
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	coreTypes := map[string]domain.CoreType{
		"core_001": {
			ID:          "core_001",
			Name:        "テストコア",
			AllowedTags: []string{"physical"},
		},
	}

	skillTypes := map[string]domain.SkillType{
		"skill_001": {
			ID:   "skill_001",
			Name: "テストスキル",
			Tags: []string{"physical"},
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{}

	// ChainEffectsマップを作成
	chainEffects := map[string]domain.ChainEffect{
		"chain_damage_bonus": domain.NewChainEffectWithTemplate(
			domain.ChainEffectDamageBonus,
			50.0,
			"次の攻撃のダメージ+%.0f",
			"ダメージ+%.0f",
		),
	}

	manager := NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills, chainEffects)

	// インベントリにコアとスキルを追加
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "chain_damage_bonus")

	// スロットを設定
	if err := manager.SetCore(0, "core_001", 5); err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}
	if err := manager.SetSkill(0, 0, "skill_001", "chain_damage_bonus"); err != nil {
		t.Fatalf("SetSkillでエラー: %v", err)
	}

	// バトル用エージェントを構築
	agents := manager.BuildAgentsForBattle()

	if len(agents) != 1 {
		t.Fatalf("エージェント数 = %d, want %d", len(agents), 1)
	}

	// スキルにChainEffectが適用されていることを確認
	agent := agents[0]
	if len(agent.Modules) != 1 {
		t.Fatalf("スキル数 = %d, want %d", len(agent.Modules), 1)
	}

	skill := agent.Modules[0]
	if !skill.HasChainEffect() {
		t.Error("スキルにChainEffectが適用されていない")
	}

	if skill.ChainEffect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("ChainEffectType = %v, want %v", skill.ChainEffect.Type, domain.ChainEffectDamageBonus)
	}

	if skill.ChainEffect.Value != 50.0 {
		t.Errorf("ChainEffect.Value = %f, want %f", skill.ChainEffect.Value, 50.0)
	}
}

func TestAgentSlotManager_BuildAgentsForBattle_WithoutChainEffect(t *testing.T) {
	// ChainEffectsを含むマネージャーを作成
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	coreTypes := map[string]domain.CoreType{
		"core_001": {
			ID:          "core_001",
			Name:        "テストコア",
			AllowedTags: []string{"physical"},
		},
	}

	skillTypes := map[string]domain.SkillType{
		"skill_001": {
			ID:   "skill_001",
			Name: "テストスキル",
			Tags: []string{"physical"},
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{}
	chainEffects := map[string]domain.ChainEffect{}

	manager := NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills, chainEffects)

	// インベントリにコアとスキルを追加（ChainEffectなし）
	coreInv.AddCore("core_001", 10)
	skillInv.AddSkill("skill_001", "")

	// スロットを設定（ChainEffectID空）
	if err := manager.SetCore(0, "core_001", 5); err != nil {
		t.Fatalf("SetCoreでエラー: %v", err)
	}
	if err := manager.SetSkill(0, 0, "skill_001", ""); err != nil {
		t.Fatalf("SetSkillでエラー: %v", err)
	}

	// バトル用エージェントを構築
	agents := manager.BuildAgentsForBattle()

	if len(agents) != 1 {
		t.Fatalf("エージェント数 = %d, want %d", len(agents), 1)
	}

	// スキルにChainEffectが適用されていないことを確認
	agent := agents[0]
	if len(agent.Modules) != 1 {
		t.Fatalf("スキル数 = %d, want %d", len(agent.Modules), 1)
	}

	skill := agent.Modules[0]
	if skill.HasChainEffect() {
		t.Error("ChainEffectIDが空の場合はChainEffectが適用されないべき")
	}
}
