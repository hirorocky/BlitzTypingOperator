// Package domain はゲームのドメインモデルを定義します。
// このファイルはAgentSlotの単体テストを定義します。

package domain

import "testing"

// ==================== SkillSlotConfig テスト ====================

func TestSkillSlotConfig_IsEmpty(t *testing.T) {
	tests := []struct {
		name   string
		config SkillSlotConfig
		want   bool
	}{
		{
			name:   "TypeIDが空の場合は空スロット",
			config: SkillSlotConfig{TypeID: "", ChainEffectID: ""},
			want:   true,
		},
		{
			name:   "TypeIDが設定されている場合は空でない",
			config: SkillSlotConfig{TypeID: "skill_001", ChainEffectID: ""},
			want:   false,
		},
		{
			name:   "TypeIDとChainEffectIDが両方設定されている場合は空でない",
			config: SkillSlotConfig{TypeID: "skill_001", ChainEffectID: "chain_001"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.IsEmpty()
			if got != tt.want {
				t.Errorf("SkillSlotConfig.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkillSlotConfig_Clear(t *testing.T) {
	config := SkillSlotConfig{TypeID: "skill_001", ChainEffectID: "chain_001"}

	config.Clear()

	if config.TypeID != "" {
		t.Errorf("Clear後のTypeIDは空であるべき、got %q", config.TypeID)
	}
	if config.ChainEffectID != "" {
		t.Errorf("Clear後のChainEffectIDは空であるべき、got %q", config.ChainEffectID)
	}
}

// ==================== AgentSlot テスト ====================

func TestAgentSlot_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		slot *AgentSlot
		want bool
	}{
		{
			name: "コア未設定の場合は空スロット",
			slot: NewAgentSlot(),
			want: true,
		},
		{
			name: "CoreTypeIDが空の場合は空スロット",
			slot: &AgentSlot{CoreTypeID: ""},
			want: true,
		},
		{
			name: "コアが設定されている場合は空でない",
			slot: &AgentSlot{CoreTypeID: "core_001"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slot.IsEmpty()
			if got != tt.want {
				t.Errorf("AgentSlot.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentSlot_GetSkillCount(t *testing.T) {
	tests := []struct {
		name string
		slot *AgentSlot
		want int
	}{
		{
			name: "全スキルスロットが空の場合は0",
			slot: NewAgentSlot(),
			want: 0,
		},
		{
			name: "1つのスキルが設定されている場合は1",
			slot: &AgentSlot{
				CoreTypeID: "core_001",
				Skills: [4]SkillSlotConfig{
					{TypeID: "skill_001", ChainEffectID: ""},
					{},
					{},
					{},
				},
			},
			want: 1,
		},
		{
			name: "複数のスキルが設定されている場合はその数",
			slot: &AgentSlot{
				CoreTypeID: "core_001",
				Skills: [4]SkillSlotConfig{
					{TypeID: "skill_001", ChainEffectID: "chain_001"},
					{TypeID: "skill_002", ChainEffectID: ""},
					{TypeID: "skill_003", ChainEffectID: "chain_002"},
					{},
				},
			},
			want: 3,
		},
		{
			name: "全スキルスロットが設定されている場合は4",
			slot: &AgentSlot{
				CoreTypeID: "core_001",
				Skills: [4]SkillSlotConfig{
					{TypeID: "skill_001", ChainEffectID: ""},
					{TypeID: "skill_002", ChainEffectID: ""},
					{TypeID: "skill_003", ChainEffectID: ""},
					{TypeID: "skill_004", ChainEffectID: ""},
				},
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.slot.GetSkillCount()
			if got != tt.want {
				t.Errorf("AgentSlot.GetSkillCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgentSlot_Clear(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
		Skills: [4]SkillSlotConfig{
			{TypeID: "skill_001", ChainEffectID: "chain_001"},
			{TypeID: "skill_002", ChainEffectID: ""},
			{},
			{},
		},
	}

	slot.Clear()

	if slot.CoreTypeID != "" {
		t.Errorf("Clear後のCoreTypeIDは空であるべき、got %q", slot.CoreTypeID)
	}
	if !slot.IsEmpty() {
		t.Error("Clear後のスロットは空であるべき")
	}
	if slot.GetSkillCount() != 0 {
		t.Errorf("Clear後のスキル数は0であるべき、got %d", slot.GetSkillCount())
	}
	for i, skillSlot := range slot.Skills {
		if !skillSlot.IsEmpty() {
			t.Errorf("Clear後のSkills[%d]は空であるべき", i)
		}
	}
}

func TestNewAgentSlot(t *testing.T) {
	slot := NewAgentSlot()

	if slot == nil {
		t.Fatal("NewAgentSlotはnilを返すべきではない")
	}
	if !slot.IsEmpty() {
		t.Error("新規作成されたスロットは空であるべき")
	}
	if slot.CoreTypeID != "" {
		t.Errorf("新規スロットのCoreTypeIDは空であるべき、got %q", slot.CoreTypeID)
	}
	if slot.GetSkillCount() != 0 {
		t.Errorf("新規スロットのスキル数は0であるべき、got %d", slot.GetSkillCount())
	}
}

func TestAgentSlot_SetCore(t *testing.T) {
	slot := NewAgentSlot()

	slot.SetCore("core_001")

	if slot.CoreTypeID != "core_001" {
		t.Errorf("SetCore後のCoreTypeID = %q, want %q", slot.CoreTypeID, "core_001")
	}
	if slot.IsEmpty() {
		t.Error("コア設定後のスロットは空であるべきではない")
	}
}

func TestAgentSlot_SetSkill(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
	}

	slot.SetSkill(0, "skill_001", "chain_001")
	slot.SetSkill(2, "skill_003", "")

	if slot.Skills[0].TypeID != "skill_001" {
		t.Errorf("Skills[0].TypeID = %q, want %q", slot.Skills[0].TypeID, "skill_001")
	}
	if slot.Skills[0].ChainEffectID != "chain_001" {
		t.Errorf("Skills[0].ChainEffectID = %q, want %q", slot.Skills[0].ChainEffectID, "chain_001")
	}
	if slot.Skills[2].TypeID != "skill_003" {
		t.Errorf("Skills[2].TypeID = %q, want %q", slot.Skills[2].TypeID, "skill_003")
	}
	if slot.Skills[2].ChainEffectID != "" {
		t.Errorf("Skills[2].ChainEffectID = %q, want %q", slot.Skills[2].ChainEffectID, "")
	}
	if slot.GetSkillCount() != 2 {
		t.Errorf("スキル数 = %d, want %d", slot.GetSkillCount(), 2)
	}
}

func TestAgentSlot_SetSkill_IgnoresInvalidIndex(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
	}

	// 無効なインデックスの設定は無視される
	slot.SetSkill(-1, "skill_001", "")
	slot.SetSkill(4, "skill_002", "")

	if slot.GetSkillCount() != 0 {
		t.Errorf("無効なインデックスへの設定後、スキル数 = %d, want %d", slot.GetSkillCount(), 0)
	}
}

func TestAgentSlot_ClearSkill(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
		Skills: [4]SkillSlotConfig{
			{TypeID: "skill_001", ChainEffectID: "chain_001"},
			{TypeID: "skill_002", ChainEffectID: ""},
			{},
			{},
		},
	}

	slot.ClearSkill(0)

	if !slot.Skills[0].IsEmpty() {
		t.Error("ClearSkill後のスキルスロットは空であるべき")
	}
	if slot.GetSkillCount() != 1 {
		t.Errorf("ClearSkill後のスキル数 = %d, want %d", slot.GetSkillCount(), 1)
	}
	// コア設定は変更されない
	if slot.IsEmpty() {
		t.Error("スキルクリアでスロット全体が空になるべきではない")
	}
}

func TestAgentSlot_ClearSkill_IgnoresInvalidIndex(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
		Skills: [4]SkillSlotConfig{
			{TypeID: "skill_001", ChainEffectID: ""},
			{},
			{},
			{},
		},
	}

	// 無効なインデックスのクリアは無視される
	slot.ClearSkill(-1)
	slot.ClearSkill(4)

	if slot.GetSkillCount() != 1 {
		t.Errorf("無効なインデックスのクリア後、スキル数 = %d, want %d", slot.GetSkillCount(), 1)
	}
}

func TestAgentSlot_GetSkill(t *testing.T) {
	slot := &AgentSlot{
		CoreTypeID: "core_001",
		Skills: [4]SkillSlotConfig{
			{TypeID: "skill_001", ChainEffectID: "chain_001"},
			{TypeID: "skill_002", ChainEffectID: ""},
			{},
			{},
		},
	}

	tests := []struct {
		name        string
		index       int
		wantTypeID  string
		wantChainID string
		wantNil     bool
	}{
		{
			name:        "有効なインデックス（設定済み）",
			index:       0,
			wantTypeID:  "skill_001",
			wantChainID: "chain_001",
			wantNil:     false,
		},
		{
			name:        "有効なインデックス（空スロット）",
			index:       2,
			wantTypeID:  "",
			wantChainID: "",
			wantNil:     false,
		},
		{
			name:    "無効なインデックス（負）",
			index:   -1,
			wantNil: true,
		},
		{
			name:    "無効なインデックス（範囲外）",
			index:   4,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slot.GetSkill(tt.index)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetSkill(%d)はnilを返すべき", tt.index)
				}
				return
			}
			if got == nil {
				t.Fatalf("GetSkill(%d)はnilを返すべきではない", tt.index)
			}
			if got.TypeID != tt.wantTypeID {
				t.Errorf("GetSkill(%d).TypeID = %q, want %q", tt.index, got.TypeID, tt.wantTypeID)
			}
			if got.ChainEffectID != tt.wantChainID {
				t.Errorf("GetSkill(%d).ChainEffectID = %q, want %q", tt.index, got.ChainEffectID, tt.wantChainID)
			}
		})
	}
}
