package presenter

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/screens"
	"hirorocky/type-battle/internal/usecase/session"
)

// CreateDefaultEncyclopediaData は図鑑のデフォルトデータを作成します。
func CreateDefaultEncyclopediaData() *screens.EncyclopediaData {
	coreTypes := []domain.CoreType{
		{
			ID:             "all_rounder",
			Name:           "オールラウンダー",
			StatWeights:    map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
			PassiveSkillID: "balance_mastery",
			AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
			MinDropLevel:   1,
		},
		{
			ID:             "attacker",
			Name:           "攻撃バランス",
			StatWeights:    map[string]float64{"STR": 1.2, "INT": 1.2, "WIL": 0.8, "LUK": 0.8},
			PassiveSkillID: "attack_boost",
			AllowedTags:    []string{"physical_low", "physical_mid", "magic_low", "magic_mid"},
			MinDropLevel:   1,
		},
		{
			ID:             "healer",
			Name:           "ヒーラー",
			StatWeights:    map[string]float64{"STR": 0.8, "INT": 1.4, "WIL": 0.9, "LUK": 0.9},
			PassiveSkillID: "heal_boost",
			AllowedTags:    []string{"heal_low", "heal_mid", "magic_low", "buff_low"},
			MinDropLevel:   5,
		},
		{
			ID:             "tank",
			Name:           "タンク",
			StatWeights:    map[string]float64{"STR": 1.1, "INT": 0.7, "WIL": 0.7, "LUK": 1.5},
			PassiveSkillID: "defense_boost",
			AllowedTags:    []string{"physical_low", "buff_low", "buff_mid"},
			MinDropLevel:   3,
		},
	}
	moduleTypes := []screens.ModuleTypeInfo{
		{ID: "physical_lv1", Name: "物理攻撃Lv1", Icon: "⚔️", Tags: []string{"physical_low"}, Description: "基本的な物理攻撃"},
		{ID: "magic_lv1", Name: "魔法攻撃Lv1", Icon: "💥", Tags: []string{"magic_low"}, Description: "基本的な魔法攻撃"},
		{ID: "heal_lv1", Name: "回復Lv1", Icon: "💚", Tags: []string{"heal_low"}, Description: "基本的な回復"},
		{ID: "buff_lv1", Name: "バフLv1", Icon: "💪", Tags: []string{"buff_low"}, Description: "味方を強化"},
		{ID: "debuff_lv1", Name: "デバフLv1", Icon: "💀", Tags: []string{"debuff_low"}, Description: "敵を弱体化"},
	}
	enemyTypes := []domain.EnemyType{
		{ID: "goblin", Name: "ゴブリン", BaseHP: 100, BaseAttackPower: 10, AttackType: "physical"},
		{ID: "orc", Name: "オーク", BaseHP: 200, BaseAttackPower: 15, AttackType: "physical"},
		{ID: "dragon", Name: "ドラゴン", BaseHP: 500, BaseAttackPower: 30, AttackType: "magic"},
	}

	return &screens.EncyclopediaData{
		AllCoreTypes:        coreTypes,
		AllModuleTypes:      moduleTypes,
		AllEnemyTypes:       enemyTypes,
		AcquiredCoreTypes:   []string{"all_rounder"},
		AcquiredModuleTypes: []string{"physical_lv1"},
		EncounteredEnemies:  []string{},
	}
}

// CreateEncyclopediaData はGameStateから図鑑データを生成します。
func CreateEncyclopediaData(gs *session.GameState) *screens.EncyclopediaData {
	// 基本データを取得
	baseData := CreateDefaultEncyclopediaData()

	// 所持コアタイプを取得（新システム: TypeIDベースのユニーク管理）
	acquiredCoreTypes := gs.Inventory().GetOwnedCores()

	// 所持スキルタイプを取得（新システム: TypeIDベースのユニーク管理）
	ownedSkills := gs.Inventory().GetOwnedSkills()
	acquiredModuleTypes := make([]string, 0, len(ownedSkills))
	for typeID := range ownedSkills {
		acquiredModuleTypes = append(acquiredModuleTypes, typeID)
	}

	return &screens.EncyclopediaData{
		AllCoreTypes:        baseData.AllCoreTypes,
		AllModuleTypes:      baseData.AllModuleTypes,
		AllEnemyTypes:       baseData.AllEnemyTypes,
		AcquiredCoreTypes:   acquiredCoreTypes,
		AcquiredModuleTypes: acquiredModuleTypes,
		EncounteredEnemies:  gs.GetEncounteredEnemies(),
	}
}
