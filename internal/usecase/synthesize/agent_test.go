// Package agent はエージェント管理機能を提供します。
// コア特性とモジュールの互換性検証、エージェント合成、装備管理を担当します。

package synthesize

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// ==================== コア特性とモジュールタグ互換性検証テスト（Task 5.1） ====================

// TestValidateModuleCompatibility はコア特性とモジュールタグの互換性検証をテストします。

func TestValidateModuleCompatibility(t *testing.T) {
	// 攻撃バランスコア（physical_low, magic_low を許可）
	attackType := domain.CoreType{
		ID:          "attack_balance",
		Name:        "攻撃バランス",
		StatWeights: map[string]float64{"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	core := domain.NewCore("core_001", "攻撃コア", 5, attackType, passiveSkill)

	// 互換性のあるモジュール
	compatibleModule := domain.NewModule(
		"m1", "物理打撃Lv1", domain.PhysicalAttack, 1,
		[]string{"physical_low"}, 10.0, "STR", "",
	)

	// 互換性のないモジュール
	incompatibleModule := domain.NewModule(
		"m2", "ヒールLv2", domain.Heal, 2,
		[]string{"heal_mid"}, 16.0, "MAG", "",
	)

	manager := NewAgentManager(nil, nil)

	if !manager.ValidateModuleCompatibility(core, compatibleModule) {
		t.Error("互換性のあるモジュールが装備不可と判定された")
	}

	if manager.ValidateModuleCompatibility(core, incompatibleModule) {
		t.Error("互換性のないモジュールが装備可と判定された")
	}
}

// TestGetAllowedTags はコア特性の許可タグリスト取得をテストします。

func TestGetAllowedTags(t *testing.T) {
	healerType := domain.CoreType{
		ID:          "healer",
		Name:        "ヒーラー",
		StatWeights: map[string]float64{"STR": 0.5, "MAG": 1.5, "SPD": 0.8, "LUK": 1.2},
		AllowedTags: []string{"heal_low", "heal_mid", "heal_high"},
	}
	passiveSkill := domain.PassiveSkill{ID: "healing_aura", Name: "ヒーリングオーラ"}
	core := domain.NewCore("core_001", "ヒーラーコア", 10, healerType, passiveSkill)

	manager := NewAgentManager(nil, nil)
	tags := manager.GetAllowedTags(core)

	if len(tags) != 3 {
		t.Errorf("許可タグ数: 期待 3, 実際 %d", len(tags))
	}

	// heal_midが含まれていることを確認
	found := false
	for _, tag := range tags {
		if tag == "heal_mid" {
			found = true
			break
		}
	}
	if !found {
		t.Error("heal_mid タグが許可タグリストに含まれていない")
	}
}

// ==================== エージェント合成機能テスト（Task 5.2） ====================

// TestSynthesizeAgent はエージェント合成処理をテストします。

func TestSynthesizeAgent(t *testing.T) {
	// インベントリをセットアップ
	coreInv := domain.NewCoreInventory(10)
	moduleInv := domain.NewModuleInventory(20)

	// コアを追加
	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "adaptability", Name: "適応力"}
	core := domain.NewCore("core_001", "オールラウンダーコア", 5, coreType, passiveSkill)
	coreInv.Add(core)

	// モジュールを追加
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "物理打撃Lv1", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "ファイアボールLv1", domain.MagicAttack, 1, []string{"magic_low"}, 12.0, "MAG", ""),
		domain.NewModule("m3", "ヒールLv1", domain.Heal, 1, []string{"heal_low"}, 8.0, "MAG", ""),
		domain.NewModule("m4", "攻撃バフLv1", domain.Buff, 1, []string{"buff_low"}, 5.0, "SPD", ""),
	}
	for _, m := range modules {
		moduleInv.Add(m)
	}

	manager := NewAgentManager(coreInv, moduleInv)
	moduleIDs := []string{"m1", "m2", "m3", "m4"}

	agent, err := manager.SynthesizeAgent("core_001", moduleIDs)
	if err != nil {
		t.Errorf("エージェント合成に失敗: %v", err)
	}
	if agent == nil {
		t.Fatal("合成されたエージェントがnil")
	}

	// エージェントのレベルがコアのレベルと同じことを確認
	if agent.Level != 5 {
		t.Errorf("エージェントレベル: 期待 5, 実際 %d", agent.Level)
	}

	// 素材が消費されていることを確認
	if coreInv.Count() != 0 {
		t.Error("コアが消費されていない")
	}
	if moduleInv.Count() != 0 {
		t.Error("モジュールが消費されていない")
	}

	// エージェントがインベントリに追加されていることを確認
	agents := manager.GetAgents()
	if len(agents) != 1 {
		t.Errorf("エージェントインベントリのエージェント数: 期待 1, 実際 %d", len(agents))
	}
}

// TestSynthesizeAgent_IncompatibleModule は互換性のないモジュールでの合成拒否をテストします。

func TestSynthesizeAgent_IncompatibleModule(t *testing.T) {
	coreInv := domain.NewCoreInventory(10)
	moduleInv := domain.NewModuleInventory(20)

	// 攻撃バランスコア（physical_low, magic_lowのみ許可）
	coreType := domain.CoreType{
		ID:          "attack_balance",
		Name:        "攻撃バランス",
		StatWeights: map[string]float64{"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "balanced_stance", Name: "バランス構え"}
	core := domain.NewCore("core_001", "攻撃コア", 5, coreType, passiveSkill)
	coreInv.Add(core)

	// 互換性のあるモジュールと互換性のないモジュール
	moduleInv.Add(domain.NewModule("m1", "物理打撃", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""))
	moduleInv.Add(domain.NewModule("m2", "ファイアボール", domain.MagicAttack, 1, []string{"magic_low"}, 12.0, "MAG", ""))
	moduleInv.Add(domain.NewModule("m3", "ヒールLv2", domain.Heal, 2, []string{"heal_mid"}, 16.0, "MAG", "")) // 互換性なし
	moduleInv.Add(domain.NewModule("m4", "攻撃バフ", domain.Buff, 1, []string{"buff_low"}, 5.0, "SPD", ""))    // 互換性なし

	manager := NewAgentManager(coreInv, moduleInv)

	_, err := manager.SynthesizeAgent("core_001", []string{"m1", "m2", "m3", "m4"})
	if err == nil {
		t.Error("互換性のないモジュールでの合成がエラーにならなかった")
	}
}

// TestSynthesizeAgent_NotEnoughModules はモジュールが4個未満での合成拒否をテストします。

func TestSynthesizeAgent_NotEnoughModules(t *testing.T) {
	coreInv := domain.NewCoreInventory(10)
	moduleInv := domain.NewModuleInventory(20)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	core := domain.NewCore("core_001", "コア", 5, coreType, passiveSkill)
	coreInv.Add(core)

	moduleInv.Add(domain.NewModule("m1", "物理打撃", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""))
	moduleInv.Add(domain.NewModule("m2", "ファイアボール", domain.MagicAttack, 1, []string{"magic_low"}, 12.0, "MAG", ""))

	manager := NewAgentManager(coreInv, moduleInv)

	_, err := manager.SynthesizeAgent("core_001", []string{"m1", "m2"})
	if err == nil {
		t.Error("モジュール不足での合成がエラーにならなかった")
	}
}

// TestGetSynthesisPreview は合成プレビューをテストします。

func TestGetSynthesisPreview(t *testing.T) {
	coreInv := domain.NewCoreInventory(10)
	moduleInv := domain.NewModuleInventory(20)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	core := domain.NewCore("core_001", "オールラウンダーコア", 10, coreType, passiveSkill)
	coreInv.Add(core)

	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "物理打撃", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "ファイアボール", domain.MagicAttack, 1, []string{"magic_low"}, 12.0, "MAG", ""),
		domain.NewModule("m3", "ヒール", domain.Heal, 1, []string{"heal_low"}, 8.0, "MAG", ""),
		domain.NewModule("m4", "攻撃バフ", domain.Buff, 1, []string{"buff_low"}, 5.0, "SPD", ""),
	}
	for _, m := range modules {
		moduleInv.Add(m)
	}

	manager := NewAgentManager(coreInv, moduleInv)
	moduleIDs := []string{"m1", "m2", "m3", "m4"}

	preview, err := manager.GetSynthesisPreview("core_001", moduleIDs)
	if err != nil {
		t.Errorf("プレビュー取得に失敗: %v", err)
	}
	if preview == nil {
		t.Fatal("プレビューがnil")
	}

	// プレビューにステータス情報が含まれていることを確認
	if preview.Level != 10 {
		t.Errorf("プレビューレベル: 期待 10, 実際 %d", preview.Level)
	}
	if preview.CoreName != "オールラウンダーコア" {
		t.Errorf("プレビューコア名: 期待 オールラウンダーコア, 実際 %s", preview.CoreName)
	}
}

// ==================== エージェント装備機能テスト（Task 5.3） ====================

// TestEquipAgent はエージェント装備処理をテストします。

func TestEquipAgent(t *testing.T) {
	manager := NewAgentManager(nil, nil)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	core := domain.NewCore("core_001", "コア", 10, coreType, passiveSkill)
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m3", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m4", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
	}
	agent := domain.NewAgent("agent_001", core, modules)
	manager.AddAgent(agent)

	player := domain.NewPlayer()

	err := manager.EquipAgent(0, "agent_001", player)
	if err != nil {
		t.Errorf("エージェント装備に失敗: %v", err)
	}

	equipped := manager.GetEquippedAgents()
	if len(equipped) != 1 {
		t.Errorf("装備エージェント数: 期待 1, 実際 %d", len(equipped))
	}

	// プレイヤーのHPが再計算されていることを確認
	if player.MaxHP == 0 {
		t.Error("プレイヤーのMaxHPが再計算されていない")
	}
}

// TestEquipAgent_MaxSlots は3スロット制限をテストします。

func TestEquipAgent_MaxSlots(t *testing.T) {
	manager := NewAgentManager(nil, nil)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m3", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m4", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
	}

	for i := 0; i < 4; i++ {
		core := domain.NewCore(
			"core_00"+string(rune('1'+i)),
			"コア",
			10,
			coreType,
			passiveSkill,
		)
		agent := domain.NewAgent("agent_00"+string(rune('1'+i)), core, modules)
		manager.AddAgent(agent)
	}

	player := domain.NewPlayer()

	// 3体まで装備可能
	for i := 0; i < 3; i++ {
		err := manager.EquipAgent(i, "agent_00"+string(rune('1'+i)), player)
		if err != nil {
			t.Errorf("スロット%dへの装備に失敗: %v", i, err)
		}
	}

	// 4スロット目は存在しないのでエラー
	err := manager.EquipAgent(3, "agent_004", player)
	if err == nil {
		t.Error("4つ目のスロットへの装備がエラーにならなかった")
	}
}

// TestUnequipAgent はエージェント装備解除処理をテストします。

func TestUnequipAgent(t *testing.T) {
	manager := NewAgentManager(nil, nil)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	core := domain.NewCore("core_001", "コア", 10, coreType, passiveSkill)
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m3", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m4", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
	}
	agent := domain.NewAgent("agent_001", core, modules)
	manager.AddAgent(agent)

	player := domain.NewPlayer()

	manager.EquipAgent(0, "agent_001", player)
	err := manager.UnequipAgent(0, player)
	if err != nil {
		t.Errorf("エージェント装備解除に失敗: %v", err)
	}

	equipped := manager.GetEquippedAgents()
	if len(equipped) != 0 {
		t.Errorf("装備解除後のエージェント数: 期待 0, 実際 %d", len(equipped))
	}
}

// TestEquipAgent_RecalculateHP は装備変更時のHP再計算をテストします。

func TestEquipAgent_RecalculateHP(t *testing.T) {
	manager := NewAgentManager(nil, nil)

	coreType := domain.CoreType{
		ID:          "all_rounder",
		Name:        "オールラウンダー",
		StatWeights: map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	passiveSkill := domain.PassiveSkill{ID: "test", Name: "テスト"}
	modules := []*domain.ModuleModel{
		domain.NewModule("m1", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m2", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m3", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
		domain.NewModule("m4", "モジュール", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", ""),
	}

	// レベル10のエージェント
	core1 := domain.NewCore("core_001", "コア1", 10, coreType, passiveSkill)
	agent1 := domain.NewAgent("agent_001", core1, modules)
	manager.AddAgent(agent1)

	// レベル20のエージェント
	core2 := domain.NewCore("core_002", "コア2", 20, coreType, passiveSkill)
	agent2 := domain.NewAgent("agent_002", core2, modules)
	manager.AddAgent(agent2)

	player := domain.NewPlayer()

	// 最初のエージェントを装備
	manager.EquipAgent(0, "agent_001", player)
	hp1 := player.MaxHP

	// 2つ目のエージェントを装備
	manager.EquipAgent(1, "agent_002", player)
	hp2 := player.MaxHP

	// 2体装備時のHPは、1体の時より高いはず（平均レベルが15になる）
	if hp2 <= hp1 {
		t.Errorf("2体装備時のHP(%d)が1体装備時のHP(%d)以下", hp2, hp1)
	}
}
