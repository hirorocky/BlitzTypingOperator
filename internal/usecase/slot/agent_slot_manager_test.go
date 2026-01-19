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
