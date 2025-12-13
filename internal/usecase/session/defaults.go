package session

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/rewarding"
)

// GetDefaultCoreTypes はデフォルトのコア特性をドメイン型で返します。
func GetDefaultCoreTypes() []domain.CoreType {
	return []domain.CoreType{
		{
			ID:             "all_rounder",
			Name:           "オールラウンダー",
			AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
			StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
			PassiveSkillID: "balance_mastery",
			MinDropLevel:   1,
		},
		{
			ID:             "attacker",
			Name:           "攻撃バランス",
			AllowedTags:    []string{"physical_low", "physical_mid", "magic_low", "magic_mid"},
			StatWeights:    map[string]float64{"STR": 1.2, "MAG": 1.2, "SPD": 0.8, "LUK": 0.8},
			PassiveSkillID: "attack_boost",
			MinDropLevel:   1,
		},
		{
			ID:             "healer",
			Name:           "ヒーラー",
			AllowedTags:    []string{"heal_low", "heal_mid", "magic_low", "buff_low"},
			StatWeights:    map[string]float64{"STR": 0.8, "MAG": 1.4, "SPD": 0.9, "LUK": 0.9},
			PassiveSkillID: "heal_boost",
			MinDropLevel:   5,
		},
		{
			ID:             "tank",
			Name:           "タンク",
			AllowedTags:    []string{"physical_low", "buff_low", "buff_mid"},
			StatWeights:    map[string]float64{"STR": 1.1, "MAG": 0.7, "SPD": 0.7, "LUK": 1.5},
			PassiveSkillID: "defense_boost",
			MinDropLevel:   3,
		},
	}
}

// GetDefaultModuleDropInfos はデフォルトのモジュール定義をドメイン型で返します。
func GetDefaultModuleDropInfos() []rewarding.ModuleDropInfo {
	return []rewarding.ModuleDropInfo{
		{ID: "mod_slash", Name: "斬撃", Category: domain.PhysicalAttack, Level: 1, Tags: []string{"physical_low"}, BaseEffect: 10.0, StatRef: "STR", Description: "基本的な物理攻撃", MinDropLevel: 1},
		{ID: "mod_thrust", Name: "突き", Category: domain.PhysicalAttack, Level: 1, Tags: []string{"physical_low"}, BaseEffect: 8.0, StatRef: "STR", Description: "素早い物理攻撃", MinDropLevel: 1},
		{ID: "mod_fireball", Name: "火球", Category: domain.MagicAttack, Level: 1, Tags: []string{"magic_low", "fire"}, BaseEffect: 12.0, StatRef: "MAG", Description: "火属性の魔法攻撃", MinDropLevel: 1},
		{ID: "mod_ice", Name: "氷結", Category: domain.MagicAttack, Level: 1, Tags: []string{"magic_low", "ice"}, BaseEffect: 11.0, StatRef: "MAG", Description: "氷属性の魔法攻撃", MinDropLevel: 1},
		{ID: "mod_heal", Name: "ヒール", Category: domain.Heal, Level: 1, Tags: []string{"heal_low"}, BaseEffect: 15.0, StatRef: "MAG", Description: "基本的な回復魔法", MinDropLevel: 1},
		{ID: "mod_attack_up", Name: "攻撃力アップ", Category: domain.Buff, Level: 1, Tags: []string{"buff_low"}, BaseEffect: 5.0, StatRef: "LUK", Description: "攻撃力を上昇させる", MinDropLevel: 1},
		{ID: "mod_defense_up", Name: "防御アップ", Category: domain.Buff, Level: 1, Tags: []string{"buff_low"}, BaseEffect: 4.0, StatRef: "LUK", Description: "防御力を上昇させる", MinDropLevel: 1},
		// レベル2モジュール
		{ID: "mod_heavy_slash", Name: "強斬撃", Category: domain.PhysicalAttack, Level: 2, Tags: []string{"physical_mid"}, BaseEffect: 20.0, StatRef: "STR", Description: "強力な物理攻撃", MinDropLevel: 5},
		{ID: "mod_blizzard", Name: "ブリザード", Category: domain.MagicAttack, Level: 2, Tags: []string{"magic_mid", "ice"}, BaseEffect: 22.0, StatRef: "MAG", Description: "氷属性の範囲魔法", MinDropLevel: 5},
		{ID: "mod_cure", Name: "キュア", Category: domain.Heal, Level: 2, Tags: []string{"heal_mid"}, BaseEffect: 30.0, StatRef: "MAG", Description: "中級回復魔法", MinDropLevel: 5},
	}
}

// GetDefaultPassiveSkills はデフォルトのパッシブスキルを返します。
func GetDefaultPassiveSkills() map[string]domain.PassiveSkill {
	return map[string]domain.PassiveSkill{
		"balanced_stats": {
			ID:          "balanced_stats",
			Name:        "バランス",
			Description: "全ステータスにバランスよくボーナス",
		},
		"attack_boost": {
			ID:          "attack_boost",
			Name:        "攻撃ブースト",
			Description: "攻撃力にボーナスを得る",
		},
		"heal_boost": {
			ID:          "heal_boost",
			Name:        "回復ブースト",
			Description: "回復効果にボーナスを得る",
		},
		"defense_boost": {
			ID:          "defense_boost",
			Name:        "防御ブースト",
			Description: "防御力にボーナスを得る",
		},
	}
}

// GetDefaultCoreType はIDからデフォルトのコア特性を検索します。
// 見つからない場合はデフォルトのオールラウンダーを返します。
func GetDefaultCoreType(coreTypeID string) domain.CoreType {
	coreTypes := GetDefaultCoreTypes()
	for _, ct := range coreTypes {
		if ct.ID == coreTypeID {
			return ct
		}
	}
	// デフォルト
	if len(coreTypes) > 0 {
		return coreTypes[0]
	}
	return domain.CoreType{
		ID:             "all_rounder",
		Name:           "オールラウンダー",
		AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
		StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		PassiveSkillID: "balance_mastery",
		MinDropLevel:   1,
	}
}

// GetDefaultPassiveSkill はコア特性IDからデフォルトのパッシブスキルを検索します。
// 見つからない場合はデフォルトのパッシブスキルを返します。
func GetDefaultPassiveSkill(coreTypeID string) domain.PassiveSkill {
	passiveSkills := GetDefaultPassiveSkills()

	skillID := coreTypeID + "_skill"
	if skill, ok := passiveSkills[skillID]; ok {
		return skill
	}
	for _, skill := range passiveSkills {
		if skill.ID == coreTypeID || skill.ID == skillID {
			return skill
		}
	}
	return domain.PassiveSkill{
		ID:          "default_skill",
		Name:        "バランス",
		Description: "バランスの取れた能力",
	}
}

// GetDefaultModuleDropInfo はIDからデフォルトのモジュール定義を検索します。
// 見つからない場合はnilを返します。
func GetDefaultModuleDropInfo(moduleID string) *rewarding.ModuleDropInfo {
	moduleDefs := GetDefaultModuleDropInfos()
	for i := range moduleDefs {
		if moduleDefs[i].ID == moduleID {
			return &moduleDefs[i]
		}
	}
	return nil
}

// FindCoreType はコア特性リストから指定IDのコア特性を検索します。
func FindCoreType(coreTypes []domain.CoreType, coreTypeID string) domain.CoreType {
	for _, ct := range coreTypes {
		if ct.ID == coreTypeID {
			return ct
		}
	}
	if len(coreTypes) > 0 {
		return coreTypes[0]
	}
	return domain.CoreType{
		ID:             "all_rounder",
		Name:           "オールラウンダー",
		AllowedTags:    []string{"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low"},
		StatWeights:    map[string]float64{"STR": 1.0, "MAG": 1.0, "SPD": 1.0, "LUK": 1.0},
		PassiveSkillID: "balance_mastery",
		MinDropLevel:   1,
	}
}

// FindPassiveSkill はパッシブスキルマップから指定コア特性に対応するスキルを検索します。
func FindPassiveSkill(passiveSkills map[string]domain.PassiveSkill, coreTypeID string) domain.PassiveSkill {
	skillID := coreTypeID + "_skill"
	if skill, ok := passiveSkills[skillID]; ok {
		return skill
	}
	for _, skill := range passiveSkills {
		if skill.ID == coreTypeID || skill.ID == skillID {
			return skill
		}
	}
	return domain.PassiveSkill{
		ID:          "default_skill",
		Name:        "バランス",
		Description: "バランスの取れた能力",
	}
}

// FindModuleDropInfo はモジュール定義リストから指定IDのモジュール定義を検索します。
func FindModuleDropInfo(moduleDefs []rewarding.ModuleDropInfo, moduleID string) *rewarding.ModuleDropInfo {
	for i := range moduleDefs {
		if moduleDefs[i].ID == moduleID {
			return &moduleDefs[i]
		}
	}
	return nil
}
