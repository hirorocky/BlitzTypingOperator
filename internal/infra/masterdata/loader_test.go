// Package masterdata はマスタデータのロード処理を提供します。
package masterdata

import (
	"os"
	"path/filepath"
	"testing"

	"hirorocky/type-battle/internal/domain"
)

// TestLoadCoreTypes はコア特性定義のロードをテストします。
func TestLoadCoreTypes(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	// テスト用cores.jsonを作成
	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			},
			{
				"id": "healer",
				"name": "ヒーラー",
				"allowed_tags": ["heal_mid", "heal_high"],
				"stat_weights": {"STR": 0.5, "MAG": 1.5, "SPD": 0.8, "LUK": 1.2},
				"passive_skill_id": "healing_aura",
				"min_drop_level": 3
			}
		]
	}`

	coresPath := filepath.Join(tmpDir, "cores.json")
	if err := os.WriteFile(coresPath, []byte(coresJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	// ロード実行
	loader := NewDataLoader(tmpDir)
	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("コア特性のロードに失敗: %v", err)
	}

	// 検証
	if len(coreTypes) != 2 {
		t.Errorf("コア特性の数が期待値と異なる: got %d, want 2", len(coreTypes))
	}

	// 攻撃バランスの検証
	if coreTypes[0].ID != "attack_balance" {
		t.Errorf("ID: got %s, want attack_balance", coreTypes[0].ID)
	}
	if coreTypes[0].Name != "攻撃バランス" {
		t.Errorf("Name: got %s, want 攻撃バランス", coreTypes[0].Name)
	}
	if len(coreTypes[0].AllowedTags) != 2 {
		t.Errorf("AllowedTags length: got %d, want 2", len(coreTypes[0].AllowedTags))
	}
	if coreTypes[0].StatWeights["STR"] != 1.2 {
		t.Errorf("StatWeights[STR]: got %f, want 1.2", coreTypes[0].StatWeights["STR"])
	}
	if coreTypes[0].MinDropLevel != 1 {
		t.Errorf("MinDropLevel: got %d, want 1", coreTypes[0].MinDropLevel)
	}
}

// TestLoadModuleDefinitions はモジュール定義のロードをテストします。
func TestLoadModuleDefinitions(t *testing.T) {
	tmpDir := t.TempDir()

	modulesJSON := `{
		"module_types": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"icon": "⚔️",
				"tags": ["physical_low"],
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 1.0, "stat_ref": "STR"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "⚔️"
					}
				]
			},
			{
				"id": "fireball_lv2",
				"name": "ファイアボールLv2",
				"icon": "🔥",
				"tags": ["magic_mid"],
				"description": "中級の魔法攻撃",
				"cooldown_seconds": 3.5,
				"difficulty": 2,
				"min_drop_level": 10,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 1.5, "stat_ref": "MAG"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "🔥"
					}
				]
			}
		]
	}`

	modulesPath := filepath.Join(tmpDir, "modules.json")
	if err := os.WriteFile(modulesPath, []byte(modulesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("モジュール定義のロードに失敗: %v", err)
	}

	if len(modules) != 2 {
		t.Errorf("モジュール数: got %d, want 2", len(modules))
	}

	// 物理打撃Lv1の検証
	if modules[0].ID != "physical_strike_lv1" {
		t.Errorf("ID: got %s, want physical_strike_lv1", modules[0].ID)
	}
	if len(modules[0].Effects) != 1 {
		t.Errorf("Effects length: got %d, want 1", len(modules[0].Effects))
	}
	if modules[0].Effects[0].Target != "enemy" {
		t.Errorf("Effects[0].Target: got %s, want enemy", modules[0].Effects[0].Target)
	}
}

// TestLoadEnemyTypes は敵タイプ定義のロードをテストします。
func TestLoadEnemyTypes(t *testing.T) {
	tmpDir := t.TempDir()

	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"attack_type": "physical",
				"ascii_art": "  ___\n /   \\\n|     |\n \\___|"
			},
			{
				"id": "goblin",
				"name": "ゴブリン",
				"base_hp": 80,
				"base_attack_power": 10,
				"attack_type": "physical",
				"ascii_art": "  /\\_/\\\n ( o o )\n  > ^ <"
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	if len(enemyTypes) != 2 {
		t.Errorf("敵タイプ数: got %d, want 2", len(enemyTypes))
	}

	// スライムの検証
	if enemyTypes[0].ID != "slime" {
		t.Errorf("ID: got %s, want slime", enemyTypes[0].ID)
	}
	if enemyTypes[0].Name != "スライム" {
		t.Errorf("Name: got %s, want スライム", enemyTypes[0].Name)
	}
	if enemyTypes[0].BaseHP != 50 {
		t.Errorf("BaseHP: got %d, want 50", enemyTypes[0].BaseHP)
	}
}

// TestLoadTypingDictionary はタイピング辞書のロードをテストします。
func TestLoadTypingDictionary(t *testing.T) {
	tmpDir := t.TempDir()

	wordsJSON := `{
		"words": {
			"easy": ["cat", "dog", "sun", "run"],
			"medium": ["function", "variable", "package"],
			"hard": ["implementation", "infrastructure"]
		}
	}`

	wordsPath := filepath.Join(tmpDir, "words.json")
	if err := os.WriteFile(wordsPath, []byte(wordsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	dictionary, err := loader.LoadTypingDictionary()
	if err != nil {
		t.Fatalf("タイピング辞書のロードに失敗: %v", err)
	}

	// easy単語の検証
	if len(dictionary.Easy) != 4 {
		t.Errorf("Easy単語数: got %d, want 4", len(dictionary.Easy))
	}
	if dictionary.Easy[0] != "cat" {
		t.Errorf("Easy[0]: got %s, want cat", dictionary.Easy[0])
	}

	// medium単語の検証
	if len(dictionary.Medium) != 3 {
		t.Errorf("Medium単語数: got %d, want 3", len(dictionary.Medium))
	}

	// hard単語の検証
	if len(dictionary.Hard) != 2 {
		t.Errorf("Hard単語数: got %d, want 2", len(dictionary.Hard))
	}
}

// TestLoadCoreTypesFileNotFound はファイルが存在しない場合のエラーをテストします。
func TestLoadCoreTypesFileNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	loader := NewDataLoader(tmpDir)
	_, err := loader.LoadCoreTypes()
	if err == nil {
		t.Error("ファイルが存在しない場合はエラーが返されるべき")
	}
}

// TestLoadCoreTypesInvalidJSON は不正なJSONの場合のエラーをテストします。
func TestLoadCoreTypesInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	// 不正なJSON
	invalidJSON := `{ invalid json }`
	coresPath := filepath.Join(tmpDir, "cores.json")
	if err := os.WriteFile(coresPath, []byte(invalidJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	_, err := loader.LoadCoreTypes()
	if err == nil {
		t.Error("不正なJSONの場合はエラーが返されるべき")
	}
}

// TestLoadAllExternalData は全外部データの一括ロードをテストします。
func TestLoadAllExternalData(t *testing.T) {
	tmpDir := t.TempDir()

	// cores.json
	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "cores.json"), []byte(coresJSON), 0644)

	// modules.json
	modulesJSON := `{
		"module_types": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"icon": "⚔️",
				"tags": ["physical_low"],
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 1.0, "stat_ref": "STR"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "⚔️"
					}
				]
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "modules.json"), []byte(modulesJSON), 0644)

	// enemies.json
	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"attack_type": "physical",
				"ascii_art": "  ___\n /   \\\n|     |"
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "enemies.json"), []byte(enemiesJSON), 0644)

	// passive_skills.json
	passiveSkillsJSON := `{
		"passive_skills": [
			{
				"id": "balanced_stance",
				"name": "バランススタンス",
				"description": "バランスの取れた姿勢"
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "passive_skills.json"), []byte(passiveSkillsJSON), 0644)

	// words.json
	wordsJSON := `{
		"words": {
			"easy": ["cat", "dog"],
			"medium": ["function"],
			"hard": ["implementation"]
		}
	}`
	os.WriteFile(filepath.Join(tmpDir, "words.json"), []byte(wordsJSON), 0644)

	// first_agent.json
	firstAgentJSON := `{
		"first_agents": [
			{
				"id": "agent_first",
				"core_type_id": "attack_balance",
				"core_level": 1,
				"modules": [
					{"type_id": "physical_strike_lv1"}
				]
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "first_agent.json"), []byte(firstAgentJSON), 0644)

	// chain_effects.json
	chainEffectsJSON := `{
		"chain_effects": [
			{
				"id": "damage_bonus",
				"name": "ダメージボーナス",
				"description": "次の攻撃のダメージ+%.0f%%",
				"short_description": "次攻撃ダメ+%.0f%%",
				"category": "attack",
				"min_value": 10,
				"max_value": 50,
				"value_step": 5
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "chain_effects.json"), []byte(chainEffectsJSON), 0644)

	loader := NewDataLoader(tmpDir)
	externalData, err := loader.LoadAllExternalData()
	if err != nil {
		t.Fatalf("全外部データのロードに失敗: %v", err)
	}

	if len(externalData.CoreTypes) != 1 {
		t.Errorf("CoreTypes: got %d, want 1", len(externalData.CoreTypes))
	}
	if len(externalData.ModuleDefinitions) != 1 {
		t.Errorf("ModuleDefinitions: got %d, want 1", len(externalData.ModuleDefinitions))
	}
	if len(externalData.EnemyTypes) != 1 {
		t.Errorf("EnemyTypes: got %d, want 1", len(externalData.EnemyTypes))
	}
	if len(externalData.PassiveSkills) != 1 {
		t.Errorf("PassiveSkills: got %d, want 1", len(externalData.PassiveSkills))
	}
	if externalData.TypingDictionary == nil {
		t.Error("TypingDictionary should not be nil")
	}
	if len(externalData.FirstAgents) == 0 {
		t.Error("FirstAgents should not be empty")
	}
}

// TestConvertToDomainCoreType はJSONデータからドメインモデルへの変換をテストします。
func TestConvertToDomainCoreType(t *testing.T) {
	tmpDir := t.TempDir()

	coresJSON := `{
		"core_types": [
			{
				"id": "attack_balance",
				"name": "攻撃バランス",
				"allowed_tags": ["physical_low", "magic_low"],
				"stat_weights": {"STR": 1.2, "MAG": 1.0, "SPD": 0.8, "LUK": 1.0},
				"passive_skill_id": "balanced_stance",
				"min_drop_level": 1
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "cores.json"), []byte(coresJSON), 0644)

	loader := NewDataLoader(tmpDir)
	coreTypes, err := loader.LoadCoreTypes()
	if err != nil {
		t.Fatalf("コア特性のロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainCoreType := coreTypes[0].ToDomain()

	if domainCoreType.ID != "attack_balance" {
		t.Errorf("ID: got %s, want attack_balance", domainCoreType.ID)
	}
	if domainCoreType.StatWeights["STR"] != 1.2 {
		t.Errorf("StatWeights[STR]: got %f, want 1.2", domainCoreType.StatWeights["STR"])
	}
}

// TestConvertToDomainSkillModel はモジュールデータからドメインモデルへの変換をテストします。
func TestConvertToDomainSkillModel(t *testing.T) {
	tmpDir := t.TempDir()

	modulesJSON := `{
		"module_types": [
			{
				"id": "physical_strike_lv1",
				"name": "物理打撃Lv1",
				"icon": "⚔️",
				"tags": ["physical_low"],
				"description": "基本的な物理攻撃",
				"cooldown_seconds": 2.0,
				"difficulty": 1,
				"min_drop_level": 1,
				"effects": [
					{
						"target": "enemy",
						"hp_formula": {"base": 0, "stat_coef": 1.0, "stat_ref": "STR"},
						"probability": 1.0,
						"luk_factor": 0,
						"icon": "⚔️"
					}
				]
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "modules.json"), []byte(modulesJSON), 0644)

	loader := NewDataLoader(tmpDir)
	modules, err := loader.LoadModuleDefinitions()
	if err != nil {
		t.Fatalf("モジュール定義のロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainModule := modules[0].ToDomain()

	if domainModule.TypeID != "physical_strike_lv1" {
		t.Errorf("TypeID: got %s, want physical_strike_lv1", domainModule.TypeID)
	}
	// 効果がダメージ効果であることを確認
	effects := domainModule.Effects()
	if len(effects) != 1 {
		t.Errorf("Effects length: got %d, want 1", len(effects))
	}
	if !effects[0].IsDamageEffect() {
		t.Error("Effect should be a damage effect")
	}
}

// TestLoadPassiveSkills はパッシブスキル定義のロードをテストします。
func TestLoadPassiveSkills(t *testing.T) {
	tmpDir := t.TempDir()

	passiveSkillsJSON := `{
		"passive_skills": [
			{
				"id": "ps_buff_extender",
				"name": "バフエクステンダー",
				"description": "バフ効果時間+50%",
				"trigger_type": "permanent",
				"effect_type": "multiplier",
				"effect_value": 1.5
			},
			{
				"id": "ps_perfect_rhythm",
				"name": "パーフェクトリズム",
				"description": "正確性100%でスキル効果1.5倍",
				"trigger_type": "conditional",
				"trigger_condition": {
					"type": "accuracy_equals",
					"value": 100
				},
				"effect_type": "multiplier",
				"effect_value": 1.5
			},
			{
				"id": "ps_last_stand",
				"name": "ラストスタンド",
				"description": "HP25%以下で30%の確率で被ダメージ1",
				"trigger_type": "probability",
				"trigger_condition": {
					"type": "hp_below_percent",
					"value": 25
				},
				"effect_type": "special",
				"effect_value": 1,
				"probability": 0.3
			},
			{
				"id": "ps_combo_master",
				"name": "コンボマスター",
				"description": "ミスなし連続タイピングでダメージ累積+10%（最大+50%）",
				"trigger_type": "stack",
				"trigger_condition": {
					"type": "no_miss_streak"
				},
				"effect_type": "modifier",
				"effect_value": 0.1,
				"max_stacks": 5,
				"stack_increment": 0.1
			},
			{
				"id": "ps_first_strike",
				"name": "ファーストストライク",
				"description": "戦闘開始時、最初のスキルが即発動",
				"trigger_type": "reactive",
				"trigger_condition": {
					"type": "on_battle_start"
				},
				"effect_type": "special",
				"uses_per_battle": 1
			}
		]
	}`

	passiveSkillsPath := filepath.Join(tmpDir, "passive_skills.json")
	if err := os.WriteFile(passiveSkillsPath, []byte(passiveSkillsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	passiveSkills, err := loader.LoadPassiveSkills()
	if err != nil {
		t.Fatalf("パッシブスキルのロードに失敗: %v", err)
	}

	if len(passiveSkills) != 5 {
		t.Errorf("パッシブスキル数: got %d, want 5", len(passiveSkills))
	}

	// バフエクステンダー（永続効果）の検証
	if passiveSkills[0].ID != "ps_buff_extender" {
		t.Errorf("ID: got %s, want ps_buff_extender", passiveSkills[0].ID)
	}
	if passiveSkills[0].TriggerType != "permanent" {
		t.Errorf("TriggerType: got %s, want permanent", passiveSkills[0].TriggerType)
	}
	if passiveSkills[0].EffectValue != 1.5 {
		t.Errorf("EffectValue: got %f, want 1.5", passiveSkills[0].EffectValue)
	}

	// パーフェクトリズム（条件付き）の検証
	if passiveSkills[1].TriggerCondition == nil {
		t.Error("TriggerCondition should not be nil")
	} else {
		if passiveSkills[1].TriggerCondition.Type != "accuracy_equals" {
			t.Errorf("TriggerCondition.Type: got %s, want accuracy_equals", passiveSkills[1].TriggerCondition.Type)
		}
	}

	// ラストスタンド（確率トリガー）の検証
	if passiveSkills[2].Probability != 0.3 {
		t.Errorf("Probability: got %f, want 0.3", passiveSkills[2].Probability)
	}

	// コンボマスター（スタック型）の検証
	if passiveSkills[3].MaxStacks != 5 {
		t.Errorf("MaxStacks: got %d, want 5", passiveSkills[3].MaxStacks)
	}
	if passiveSkills[3].StackIncrement != 0.1 {
		t.Errorf("StackIncrement: got %f, want 0.1", passiveSkills[3].StackIncrement)
	}

	// ファーストストライク（反応型）の検証
	if passiveSkills[4].UsesPerBattle != 1 {
		t.Errorf("UsesPerBattle: got %d, want 1", passiveSkills[4].UsesPerBattle)
	}
}

// TestConvertToDomainPassiveSkill はパッシブスキルデータからドメインモデルへの変換をテストします。
func TestConvertToDomainPassiveSkill(t *testing.T) {
	tmpDir := t.TempDir()

	passiveSkillsJSON := `{
		"passive_skills": [
			{
				"id": "ps_perfect_rhythm",
				"name": "パーフェクトリズム",
				"description": "正確性100%でスキル効果1.5倍",
				"trigger_type": "conditional",
				"trigger_condition": {
					"type": "accuracy_equals",
					"value": 100
				},
				"effect_type": "multiplier",
				"effect_value": 1.5
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "passive_skills.json"), []byte(passiveSkillsJSON), 0644)

	loader := NewDataLoader(tmpDir)
	passiveSkills, err := loader.LoadPassiveSkills()
	if err != nil {
		t.Fatalf("パッシブスキルのロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainPassiveSkill := passiveSkills[0].ToDomain()

	if domainPassiveSkill.ID != "ps_perfect_rhythm" {
		t.Errorf("ID: got %s, want ps_perfect_rhythm", domainPassiveSkill.ID)
	}
	if domainPassiveSkill.TriggerType != domain.PassiveTriggerConditional {
		t.Errorf("TriggerType: got %s, want %s", domainPassiveSkill.TriggerType, domain.PassiveTriggerConditional)
	}
	if domainPassiveSkill.TriggerCondition == nil {
		t.Error("TriggerCondition should not be nil")
	} else {
		if domainPassiveSkill.TriggerCondition.Type != domain.TriggerConditionAccuracyEquals {
			t.Errorf("TriggerCondition.Type: got %s, want %s", domainPassiveSkill.TriggerCondition.Type, domain.TriggerConditionAccuracyEquals)
		}
	}
	if domainPassiveSkill.EffectType != domain.PassiveEffectMultiplier {
		t.Errorf("EffectType: got %s, want %s", domainPassiveSkill.EffectType, domain.PassiveEffectMultiplier)
	}
}

// TestLoadChainEffects はチェイン効果定義のロードをテストします。
func TestLoadChainEffects(t *testing.T) {
	tmpDir := t.TempDir()

	chainEffectsJSON := `{
		"chain_effects": [
			{
				"id": "damage_bonus",
				"name": "ダメージボーナス",
				"description": "次の攻撃のダメージにボーナスを付与",
				"category": "attack",
				"effect_type": "damage_bonus",
				"min_value": 10,
				"max_value": 50
			},
			{
				"id": "heal_bonus",
				"name": "ヒールボーナス",
				"description": "次の回復量にボーナスを付与",
				"category": "heal",
				"effect_type": "heal_bonus",
				"min_value": 15,
				"max_value": 40
			},
			{
				"id": "damage_cut",
				"name": "ダメージカット",
				"description": "効果中の被ダメージを軽減",
				"category": "defense",
				"effect_type": "damage_cut",
				"min_value": 10,
				"max_value": 30
			},
			{
				"id": "time_extend",
				"name": "タイムエクステンド",
				"description": "効果中のタイピング制限時間を延長",
				"category": "typing",
				"effect_type": "time_extend",
				"min_value": 1,
				"max_value": 3
			},
			{
				"id": "cooldown_reduce",
				"name": "クールダウンリデュース",
				"description": "他エージェントのリキャスト時間を短縮",
				"category": "recast",
				"effect_type": "cooldown_reduce",
				"min_value": 10,
				"max_value": 30
			},
			{
				"id": "buff_duration",
				"name": "バフデュレーション",
				"description": "バフスキルの効果時間を延長",
				"category": "effect_extend",
				"effect_type": "buff_duration",
				"min_value": 1,
				"max_value": 5
			},
			{
				"id": "double_cast",
				"name": "ダブルキャスト",
				"description": "一定確率でスキルを2回発動",
				"category": "special",
				"effect_type": "double_cast",
				"min_value": 10,
				"max_value": 25
			}
		]
	}`

	chainEffectsPath := filepath.Join(tmpDir, "chain_effects.json")
	if err := os.WriteFile(chainEffectsPath, []byte(chainEffectsJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	chainEffects, err := loader.LoadChainEffects()
	if err != nil {
		t.Fatalf("チェイン効果のロードに失敗: %v", err)
	}

	if len(chainEffects) != 7 {
		t.Errorf("チェイン効果数: got %d, want 7", len(chainEffects))
	}

	// ダメージボーナス（攻撃カテゴリ）の検証
	if chainEffects[0].ID != "damage_bonus" {
		t.Errorf("ID: got %s, want damage_bonus", chainEffects[0].ID)
	}
	if chainEffects[0].Category != "attack" {
		t.Errorf("Category: got %s, want attack", chainEffects[0].Category)
	}
	if chainEffects[0].EffectType != "damage_bonus" {
		t.Errorf("EffectType: got %s, want damage_bonus", chainEffects[0].EffectType)
	}
	if chainEffects[0].MinValue != 10 {
		t.Errorf("MinValue: got %f, want 10", chainEffects[0].MinValue)
	}
	if chainEffects[0].MaxValue != 50 {
		t.Errorf("MaxValue: got %f, want 50", chainEffects[0].MaxValue)
	}

	// 各カテゴリの検証
	categories := map[string]bool{
		"attack":        false,
		"heal":          false,
		"defense":       false,
		"typing":        false,
		"recast":        false,
		"effect_extend": false,
		"special":       false,
	}
	for _, effect := range chainEffects {
		categories[effect.Category] = true
	}
	for cat, found := range categories {
		if !found {
			t.Errorf("カテゴリ %s が見つからない", cat)
		}
	}
}

// TestConvertToDomainChainEffectType はチェイン効果データからドメインモデルへの変換をテストします。
func TestConvertToDomainChainEffectType(t *testing.T) {
	tmpDir := t.TempDir()

	chainEffectsJSON := `{
		"chain_effects": [
			{
				"id": "damage_amp",
				"name": "ダメージアンプ",
				"description": "効果中の攻撃ダメージを増加",
				"category": "attack",
				"effect_type": "damage_amp",
				"min_value": 10,
				"max_value": 30
			}
		]
	}`
	os.WriteFile(filepath.Join(tmpDir, "chain_effects.json"), []byte(chainEffectsJSON), 0644)

	loader := NewDataLoader(tmpDir)
	chainEffects, err := loader.LoadChainEffects()
	if err != nil {
		t.Fatalf("チェイン効果のロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainEffectType := chainEffects[0].ToDomainEffectType()
	domainCategory := chainEffects[0].ToDomainCategory()

	if domainEffectType != domain.ChainEffectDamageAmp {
		t.Errorf("EffectType: got %s, want %s", domainEffectType, domain.ChainEffectDamageAmp)
	}
	if domainCategory != domain.ChainEffectCategoryAttack {
		t.Errorf("Category: got %s, want %s", domainCategory, domain.ChainEffectCategoryAttack)
	}
}

// ==================== ボルテージ設定テスト ====================

// TestLoadEnemyTypesWithVoltageRise はボルテージ上昇率の読み込みをテストします。
func TestLoadEnemyTypesWithVoltageRise(t *testing.T) {
	tmpDir := t.TempDir()

	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"attack_type": "physical",
				"ascii_art": "  ___",
				"voltage_rise_per_10s": 15.0
			},
			{
				"id": "boss",
				"name": "ボス",
				"base_hp": 500,
				"base_attack_power": 20,
				"attack_type": "physical",
				"ascii_art": "  BOSS",
				"voltage_rise_per_10s": 25.5
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	// スライムのボルテージ上昇率を検証
	if enemyTypes[0].GetVoltageRisePer10s() != 15.0 {
		t.Errorf("VoltageRisePer10s: got %f, want 15.0", enemyTypes[0].GetVoltageRisePer10s())
	}

	// ボスのボルテージ上昇率を検証（小数値対応）
	if enemyTypes[1].GetVoltageRisePer10s() != 25.5 {
		t.Errorf("VoltageRisePer10s: got %f, want 25.5", enemyTypes[1].GetVoltageRisePer10s())
	}
}

// TestLoadEnemyTypesVoltageDefaultValue はボルテージ上昇率未設定時のデフォルト値をテストします。
func TestLoadEnemyTypesVoltageDefaultValue(t *testing.T) {
	tmpDir := t.TempDir()

	// voltage_rise_per_10sフィールドが未設定のJSON
	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "slime",
				"name": "スライム",
				"base_hp": 50,
				"base_attack_power": 5,
				"attack_type": "physical",
				"ascii_art": "  ___"
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	// ToDomain変換してデフォルト値が適用されることを検証
	domainEnemy := enemyTypes[0].ToDomain()
	if domainEnemy.VoltageRisePer10s != 10.0 {
		t.Errorf("VoltageRisePer10s default: got %f, want 10.0", domainEnemy.VoltageRisePer10s)
	}
}

// TestLoadEnemyTypesVoltageZeroValue はボルテージ上昇率0（上昇なし）をテストします。
func TestLoadEnemyTypesVoltageZeroValue(t *testing.T) {
	tmpDir := t.TempDir()

	// voltage_rise_per_10sが0のJSON（ボルテージ上昇なし）
	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "calm_enemy",
				"name": "穏やかな敵",
				"base_hp": 100,
				"base_attack_power": 10,
				"attack_type": "physical",
				"ascii_art": "  :)",
				"voltage_rise_per_10s": 0
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	// 0が正しく読み込まれることを検証
	if enemyTypes[0].GetVoltageRisePer10s() != 0 {
		t.Errorf("VoltageRisePer10s: got %f, want 0", enemyTypes[0].GetVoltageRisePer10s())
	}

	// ToDomain変換で0が維持されることを検証
	domainEnemy := enemyTypes[0].ToDomain()
	if domainEnemy.VoltageRisePer10s != 0 {
		t.Errorf("VoltageRisePer10s after ToDomain: got %f, want 0", domainEnemy.VoltageRisePer10s)
	}
}

// TestEnemyTypeToDomainWithVoltage はToDomain変換時のボルテージ上昇率設定をテストします。
func TestEnemyTypeToDomainWithVoltage(t *testing.T) {
	tmpDir := t.TempDir()

	enemiesJSON := `{
		"enemy_types": [
			{
				"id": "fast_enemy",
				"name": "高速敵",
				"base_hp": 80,
				"base_attack_power": 12,
				"attack_type": "magic",
				"ascii_art": "  FAST",
				"voltage_rise_per_10s": 20.0
			}
		]
	}`

	enemiesPath := filepath.Join(tmpDir, "enemies.json")
	if err := os.WriteFile(enemiesPath, []byte(enemiesJSON), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}

	loader := NewDataLoader(tmpDir)
	enemyTypes, err := loader.LoadEnemyTypes()
	if err != nil {
		t.Fatalf("敵タイプのロードに失敗: %v", err)
	}

	// ドメインモデルに変換
	domainEnemy := enemyTypes[0].ToDomain()

	// 基本フィールドの検証
	if domainEnemy.ID != "fast_enemy" {
		t.Errorf("ID: got %s, want fast_enemy", domainEnemy.ID)
	}

	// ボルテージ上昇率が正しく変換されていることを検証
	if domainEnemy.VoltageRisePer10s != 20.0 {
		t.Errorf("VoltageRisePer10s: got %f, want 20.0", domainEnemy.VoltageRisePer10s)
	}
}
