// Package app はmasterdata型からdomain型への変換ヘルパー関数を提供します。
// usecase層からinfra層への依存を解消するため、変換処理はapp層で行います。
package app

import (
	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/usecase/rewarding"
)

// ConvertEnemyTypes はmasterdata.EnemyTypeDataのスライスをdomain.EnemyTypeのスライスに変換します。
func ConvertEnemyTypes(types []masterdata.EnemyTypeData) []domain.EnemyType {
	result := make([]domain.EnemyType, len(types))
	for i, t := range types {
		result[i] = t.ToDomain()
	}
	return result
}

// ConvertCoreTypes はmasterdata.CoreTypeDataのスライスをdomain.CoreTypeのスライスに変換します。
func ConvertCoreTypes(types []masterdata.CoreTypeData) []domain.CoreType {
	result := make([]domain.CoreType, len(types))
	for i, t := range types {
		result[i] = t.ToDomain()
	}
	return result
}

// ConvertModuleTypes はmasterdata.ModuleDefinitionDataのスライスをreward.ModuleDropInfoのスライスに変換します。
func ConvertModuleTypes(types []masterdata.ModuleDefinitionData) []rewarding.ModuleDropInfo {
	result := make([]rewarding.ModuleDropInfo, len(types))
	for i, t := range types {
		result[i] = rewarding.ModuleDropInfo{
			ID:           t.ID,
			Name:         t.Name,
			Category:     convertCategory(t.Category),
			Level:        t.Level,
			Tags:         t.Tags,
			BaseEffect:   t.BaseEffect,
			StatRef:      t.StatReference,
			Description:  t.Description,
			MinDropLevel: t.MinDropLevel,
		}
	}
	return result
}

// convertCategory はカテゴリ文字列をdomain.ModuleCategoryに変換します。
func convertCategory(cat string) domain.ModuleCategory {
	switch cat {
	case "physical_attack":
		return domain.PhysicalAttack
	case "magic_attack":
		return domain.MagicAttack
	case "heal":
		return domain.Heal
	case "buff":
		return domain.Buff
	case "debuff":
		return domain.Debuff
	default:
		return domain.PhysicalAttack
	}
}

// ConvertExternalDataToDomain はExternalDataから全てのドメイン型データを変換します。
func ConvertExternalDataToDomain(ext *masterdata.ExternalData) (
	[]domain.EnemyType,
	[]domain.CoreType,
	[]rewarding.ModuleDropInfo,
) {
	if ext == nil {
		return nil, nil, nil
	}

	enemyTypes := ConvertEnemyTypes(ext.EnemyTypes)
	coreTypes := ConvertCoreTypes(ext.CoreTypes)
	moduleTypes := ConvertModuleTypes(ext.ModuleDefinitions)

	return enemyTypes, coreTypes, moduleTypes
}
