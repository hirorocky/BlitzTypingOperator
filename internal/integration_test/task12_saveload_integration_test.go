// Package integration_test はタスク12の統合テストを提供します。

package integration_test

import (
	"os"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/savedata"
)

// ==================================================
// Task 12.1: セーブ・ロードの統合テスト
// ==================================================

// TestSaveLoad_NewFormatPersistence は新形式セーブデータの保存・読み込みを検証します。
func TestSaveLoad_NewFormatPersistence(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	// 新形式のセーブデータを作成
	saveData := savedata.NewSaveData()

	// コアインスタンス（TypeID + Level形式）
	saveData.Inventory.CoreInstances = []savedata.CoreInstanceSave{
		{CoreTypeID: "attack_balance", Level: 5},
		{CoreTypeID: "paladin", Level: 3},
		{CoreTypeID: "all_rounder", Level: 1},
	}

	// モジュールインスタンス（TypeID + ChainEffect形式）
	saveData.Inventory.ModuleInstances = []savedata.ModuleInstanceSave{
		{
			TypeID: "physical_lv1",
			ChainEffect: &savedata.ChainEffectSave{
				Type:  "damage_bonus",
				Value: 20.0,
			},
		},
		{
			TypeID: "magic_lv1",
			ChainEffect: &savedata.ChainEffectSave{
				Type:  "heal_bonus",
				Value: 15.0,
			},
		},
		{
			TypeID:      "heal_lv1",
			ChainEffect: nil, // チェイン効果なし
		},
	}

	// エージェントインスタンス
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{
		{
			ID: "agent-1",
			Core: savedata.CoreInstanceSave{
				CoreTypeID: "attack_balance",
				Level:      5,
			},
			Skills: []savedata.SkillInstanceSave{
				{TypeID: "physical_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "damage_amp", Value: 25.0}},
				{TypeID: "magic_lv1"},
				{TypeID: "heal_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "buff_extend", Value: 3.0}},
				{TypeID: "buff_lv1"},
			},
		},
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// バージョン確認
	if loadedData.Version != savedata.CurrentSaveDataVersion {
		t.Errorf("バージョン expected %s, got %s", savedata.CurrentSaveDataVersion, loadedData.Version)
	}

	// コアインスタンスの検証
	if len(loadedData.Inventory.CoreInstances) != 3 {
		t.Fatalf("コアインスタンス数 expected 3, got %d", len(loadedData.Inventory.CoreInstances))
	}
	if loadedData.Inventory.CoreInstances[0].CoreTypeID != "attack_balance" {
		t.Error("コア1のTypeIDが正しく復元されていません")
	}
	if loadedData.Inventory.CoreInstances[0].Level != 5 {
		t.Error("コア1のLevelが正しく復元されていません")
	}

	// モジュールインスタンスの検証
	if len(loadedData.Inventory.ModuleInstances) != 3 {
		t.Fatalf("モジュールインスタンス数 expected 3, got %d", len(loadedData.Inventory.ModuleInstances))
	}
	mod1 := loadedData.Inventory.ModuleInstances[0]
	if mod1.TypeID != "physical_lv1" {
		t.Error("モジュール1のTypeIDが正しく復元されていません")
	}
	if mod1.ChainEffect == nil || mod1.ChainEffect.Type != "damage_bonus" {
		t.Error("モジュール1のChainEffectが正しく復元されていません")
	}
	if mod1.ChainEffect.Value != 20.0 {
		t.Errorf("モジュール1のChainEffect Value expected 20.0, got %f", mod1.ChainEffect.Value)
	}

	// チェイン効果なしモジュールの検証
	mod3 := loadedData.Inventory.ModuleInstances[2]
	if mod3.ChainEffect != nil {
		t.Error("モジュール3はチェイン効果なしであるべきです")
	}

	// エージェントインスタンスの検証
	if len(loadedData.Inventory.AgentInstances) != 1 {
		t.Fatalf("エージェントインスタンス数 expected 1, got %d", len(loadedData.Inventory.AgentInstances))
	}
	agent := loadedData.Inventory.AgentInstances[0]
	if agent.ID != "agent-1" {
		t.Error("エージェントIDが正しく復元されていません")
	}
	if agent.Core.CoreTypeID != "attack_balance" {
		t.Error("エージェントのコアTypeIDが正しく復元されていません")
	}
	if len(agent.Skills) != 4 {
		t.Fatalf("Modules数 expected 4, got %d", len(agent.Skills))
	}
	if agent.Skills[0].ChainEffect == nil || agent.Skills[0].ChainEffect.Type != "damage_amp" {
		t.Error("エージェントのモジュール1のチェイン効果が正しく復元されていません")
	}
	if agent.Skills[1].ChainEffect != nil {
		t.Error("エージェントのモジュール2はチェイン効果なしであるべきです")
	}
}

// TestSaveLoad_ChainEffectPersistence はChainEffectの永続化・復元を詳細に検証します。
func TestSaveLoad_ChainEffectPersistence(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	// 全種類のチェイン効果タイプをテスト
	chainEffectTypes := []string{
		"damage_bonus", "heal_bonus", "buff_extend", "debuff_extend",
		"damage_amp", "armor_pierce", "life_steal",
		"damage_cut", "evasion", "reflect", "regen",
		"heal_amp", "overheal",
		"time_extend", "auto_correct",
		"cooldown_reduce",
		"buff_duration", "debuff_duration",
		"double_cast",
	}

	saveData := savedata.NewSaveData()
	for i, effectType := range chainEffectTypes {
		saveData.Inventory.ModuleInstances = append(saveData.Inventory.ModuleInstances,
			savedata.ModuleInstanceSave{
				TypeID: "module_" + effectType,
				ChainEffect: &savedata.ChainEffectSave{
					Type:  effectType,
					Value: float64(i+1) * 5.0,
				},
			})
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// 全てのチェイン効果タイプが正しく復元されていることを確認
	if len(loadedData.Inventory.ModuleInstances) != len(chainEffectTypes) {
		t.Fatalf("モジュール数 expected %d, got %d", len(chainEffectTypes), len(loadedData.Inventory.ModuleInstances))
	}

	for i, effectType := range chainEffectTypes {
		mod := loadedData.Inventory.ModuleInstances[i]
		if mod.ChainEffect == nil {
			t.Errorf("モジュール %d のChainEffectがnilです", i)
			continue
		}
		if mod.ChainEffect.Type != effectType {
			t.Errorf("モジュール %d のChainEffect.Type expected %s, got %s", i, effectType, mod.ChainEffect.Type)
		}
		expectedValue := float64(i+1) * 5.0
		if mod.ChainEffect.Value != expectedValue {
			t.Errorf("モジュール %d のChainEffect.Value expected %f, got %f", i, expectedValue, mod.ChainEffect.Value)
		}
	}
}

// TestSaveLoad_PassiveSkillDataIntegrity はPassiveSkillの永続化・復元を検証します。
// PassiveSkillはマスタデータから導出されるため、CoreTypeIDの永続化で間接的に検証します。
func TestSaveLoad_PassiveSkillDataIntegrity(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	saveData := savedata.NewSaveData()

	// コアインスタンス（各コアTypeにはパッシブスキルが紐づいている）
	saveData.Inventory.CoreInstances = []savedata.CoreInstanceSave{
		{CoreTypeID: "attack_balance", Level: 1}, // ps_perfect_rhythm
		{CoreTypeID: "paladin", Level: 2},        // ps_last_stand
		{CoreTypeID: "all_rounder", Level: 3},    // ps_buff_extender
		{CoreTypeID: "healer", Level: 4},         // ps_miracle_heal
		{CoreTypeID: "speedster", Level: 5},      // ps_overdrive
	}

	// エージェントにもコアを設定
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{
		{
			ID: "agent-ps-test",
			Core: savedata.CoreInstanceSave{
				CoreTypeID: "attack_balance",
				Level:      10,
			},
			Skills: []savedata.SkillInstanceSave{
				{TypeID: "m1"},
				{TypeID: "m2"},
				{TypeID: "m3"},
				{TypeID: "m4"},
			},
		},
	}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// コアのTypeIDとLevelが正しく復元されていることを確認
	// パッシブスキルはロード時にマスタデータから再導出される
	expectedCores := []struct {
		typeID string
		level  int
	}{
		{"attack_balance", 1},
		{"paladin", 2},
		{"all_rounder", 3},
		{"healer", 4},
		{"speedster", 5},
	}

	for i, expected := range expectedCores {
		if i >= len(loadedData.Inventory.CoreInstances) {
			t.Fatalf("コアインスタンス %d が存在しません", i)
		}
		core := loadedData.Inventory.CoreInstances[i]
		if core.CoreTypeID != expected.typeID {
			t.Errorf("コア %d のTypeID expected %s, got %s", i, expected.typeID, core.CoreTypeID)
		}
		if core.Level != expected.level {
			t.Errorf("コア %d のLevel expected %d, got %d", i, expected.level, core.Level)
		}
	}

	// エージェントのコアも確認
	if len(loadedData.Inventory.AgentInstances) != 1 {
		t.Fatal("エージェントインスタンスが復元されていません")
	}
	agentCore := loadedData.Inventory.AgentInstances[0].Core
	if agentCore.CoreTypeID != "attack_balance" {
		t.Errorf("エージェントのコアが正しく復元されていません: %+v", agentCore)
	}
}

// TestSaveLoad_OldFormatDetection は旧形式セーブデータ検出時の処理を検証します。
func TestSaveLoad_OldFormatDetection(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	// 旧形式（version空）のセーブデータを直接書き込み
	oldFormatData := `{
		"version": "",
		"player": {"equipped_agent_ids": ["", "", ""]},
		"inventory": {
			"core_instances": [],
			"module_counts": {"physical_lv1": 3},
			"agent_instances": [],
			"max_core_slots": 100,
			"max_module_slots": 200,
			"max_agent_slots": 20
		}
	}`

	savePath := tempDir + "/save.json"
	err := os.WriteFile(savePath, []byte(oldFormatData), 0644)
	if err != nil {
		t.Fatalf("テストファイルの書き込みに失敗: %v", err)
	}

	// ロード試行（バージョン空のためエラーになるべき）
	_, err = io.LoadGame()
	if err == nil {
		t.Error("旧形式セーブデータではエラーが返されるべきです")
	}

	// エラーメッセージにバージョン関連の情報が含まれていることを確認
	// （実際の挙動はsavedata.goの実装による）
}

// TestSaveLoad_DataVersionMigration はデータバージョン管理を検証します。
func TestSaveLoad_DataVersionMigration(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	// 現在のバージョンでセーブ
	saveData := savedata.NewSaveData()
	saveData.Inventory.ModuleInstances = []savedata.ModuleInstanceSave{
		{
			TypeID: "test_module",
			ChainEffect: &savedata.ChainEffectSave{
				Type:  "damage_bonus",
				Value: 10.0,
			},
		},
	}

	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// バージョンが現在のバージョンであることを確認
	if loadedData.Version != savedata.CurrentSaveDataVersion {
		t.Errorf("バージョン expected %s, got %s", savedata.CurrentSaveDataVersion, loadedData.Version)
	}

	// データが正しくロードされていることを確認
	if len(loadedData.Inventory.ModuleInstances) != 1 {
		t.Fatal("モジュールインスタンスが復元されていません")
	}
}

// TestSaveLoad_ComplexAgentData は複雑なエージェントデータの永続化を検証します。
func TestSaveLoad_ComplexAgentData(t *testing.T) {
	tempDir := t.TempDir()
	io := savedata.NewSaveDataIO(tempDir, false)

	saveData := savedata.NewSaveData()

	// 複数エージェント、複数チェイン効果
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{
		{
			ID: "agent-1",
			Core: savedata.CoreInstanceSave{
				CoreTypeID: "attack_balance",
				Level:      5,
			},
			Skills: []savedata.SkillInstanceSave{
				{TypeID: "physical_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "damage_amp", Value: 25.0}},
				{TypeID: "magic_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "armor_pierce", Value: 1.0}},
				{TypeID: "heal_lv1"},
				{TypeID: "buff_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "buff_duration", Value: 5.0}},
			},
		},
		{
			ID: "agent-2",
			Core: savedata.CoreInstanceSave{
				CoreTypeID: "healer",
				Level:      3,
			},
			Skills: []savedata.SkillInstanceSave{
				{TypeID: "heal_lv2", ChainEffect: &savedata.ChainEffectSave{Type: "heal_amp", Value: 30.0}},
				{TypeID: "buff_lv2", ChainEffect: &savedata.ChainEffectSave{Type: "regen", Value: 2.0}},
				{TypeID: "debuff_lv1", ChainEffect: &savedata.ChainEffectSave{Type: "debuff_duration", Value: 4.0}},
				{TypeID: "magic_lv2"},
			},
		},
		{
			ID: "agent-3",
			Core: savedata.CoreInstanceSave{
				CoreTypeID: "speedster",
				Level:      7,
			},
			Skills: []savedata.SkillInstanceSave{
				{TypeID: "physical_lv3", ChainEffect: &savedata.ChainEffectSave{Type: "life_steal", Value: 15.0}},
				{TypeID: "buff_lv3", ChainEffect: &savedata.ChainEffectSave{Type: "cooldown_reduce", Value: 20.0}},
				{TypeID: "debuff_lv2", ChainEffect: &savedata.ChainEffectSave{Type: "double_cast", Value: 10.0}},
				{TypeID: "heal_lv3", ChainEffect: &savedata.ChainEffectSave{Type: "overheal", Value: 1.0}},
			},
		},
	}

	// 装備エージェントIDを設定
	saveData.Player.EquippedAgentIDs = [3]string{"agent-1", "agent-2", "agent-3"}

	// セーブ
	err := io.SaveGame(saveData)
	if err != nil {
		t.Fatalf("セーブに失敗: %v", err)
	}

	// ロード
	loadedData, err := io.LoadGame()
	if err != nil {
		t.Fatalf("ロードに失敗: %v", err)
	}

	// エージェント数の確認
	if len(loadedData.Inventory.AgentInstances) != 3 {
		t.Fatalf("エージェント数 expected 3, got %d", len(loadedData.Inventory.AgentInstances))
	}

	// 各エージェントのデータ整合性を確認
	for i, agent := range loadedData.Inventory.AgentInstances {
		expectedID := saveData.Inventory.AgentInstances[i].ID
		if agent.ID != expectedID {
			t.Errorf("エージェント %d のID expected %s, got %s", i, expectedID, agent.ID)
		}

		// モジュール数の確認
		if len(agent.Skills) != 4 {
			t.Errorf("エージェント %d のModules数 expected 4, got %d", i, len(agent.Skills))
		}
	}

	// 装備エージェントIDの確認
	for i, expectedID := range saveData.Player.EquippedAgentIDs {
		if loadedData.Player.EquippedAgentIDs[i] != expectedID {
			t.Errorf("装備エージェント %d のID expected %s, got %s", i, expectedID, loadedData.Player.EquippedAgentIDs[i])
		}
	}
}

// TestDomainConversion_ChainEffect はドメインへの変換を検証するヘルパーテストです。
func TestDomainConversion_ChainEffect(t *testing.T) {
	// ChainEffectSaveからChainEffectへの変換を検証
	chainEffectSave := &savedata.ChainEffectSave{
		Type:  "damage_bonus",
		Value: 25.0,
	}

	// ドメインChainEffectへの変換（テンプレート付き）
	chainEffect := domain.NewChainEffectWithTemplate(
		"test_damage_bonus",
		domain.ChainEffectType(chainEffectSave.Type),
		chainEffectSave.Value,
		"次の攻撃のダメージ+%.0f%%",
		"次攻撃ダメ+%.0f%%",
	)

	if chainEffect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("ChainEffect.Type expected damage_bonus, got %s", chainEffect.Type)
	}
	if chainEffect.Value != 25.0 {
		t.Errorf("ChainEffect.Value expected 25.0, got %f", chainEffect.Value)
	}
	if chainEffect.Description == "" {
		t.Error("ChainEffect.Descriptionが生成されるべきです")
	}
}
