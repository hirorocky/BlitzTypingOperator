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

// ConvertEnemyPassiveSkills はmasterdata.EnemyPassiveSkillDataのスライスをIDマップに変換します。
func ConvertEnemyPassiveSkills(skills []masterdata.EnemyPassiveSkillData) map[string]*domain.EnemyPassiveSkill {
	result := make(map[string]*domain.EnemyPassiveSkill, len(skills))
	for _, s := range skills {
		result[s.ID] = s.ToDomain()
	}
	return result
}

// ConvertEnemyTypesWithPassives は敵タイプを変換し、パッシブスキルを解決します。
func ConvertEnemyTypesWithPassives(
	types []masterdata.EnemyTypeData,
	passives []masterdata.EnemyPassiveSkillData,
) []domain.EnemyType {
	// パッシブスキルをマップに変換
	passiveMap := ConvertEnemyPassiveSkills(passives)

	result := make([]domain.EnemyType, len(types))
	for i, t := range types {
		result[i] = t.ToDomain()

		// パッシブIDが設定されている場合、パッシブを解決
		if t.NormalPassiveID != "" {
			if passive, ok := passiveMap[t.NormalPassiveID]; ok {
				result[i].NormalPassive = passive
			}
		}
		if t.EnhancedPassiveID != "" {
			if passive, ok := passiveMap[t.EnhancedPassiveID]; ok {
				result[i].EnhancedPassive = passive
			}
		}
	}
	return result
}

// ConvertEnemyActions はmasterdata.EnemyActionDataのスライスをドメインモデルに変換します。
func ConvertEnemyActions(actions []masterdata.EnemyActionData) []domain.EnemyAction {
	result := make([]domain.EnemyAction, len(actions))
	for i, a := range actions {
		result[i] = a.ToDomain()
	}
	return result
}

// ResolveEnemyTypeActions は敵タイプの行動パターンIDを実際のEnemyActionに解決します。
func ResolveEnemyTypeActions(enemyTypes []domain.EnemyType, actions []domain.EnemyAction) {
	actionMap := make(map[string]domain.EnemyAction)
	for _, action := range actions {
		actionMap[action.ID] = action
	}

	for i := range enemyTypes {
		var normalActions []domain.EnemyAction
		for _, id := range enemyTypes[i].NormalActionPatternIDs {
			if action, ok := actionMap[id]; ok {
				normalActions = append(normalActions, action)
			}
		}

		var enhancedActions []domain.EnemyAction
		for _, id := range enemyTypes[i].EnhancedActionPatternIDs {
			if action, ok := actionMap[id]; ok {
				enhancedActions = append(enhancedActions, action)
			}
		}

		enemyTypes[i].SetResolvedActions(normalActions, enhancedActions)
	}
}

// ConvertCoreTypes はmasterdata.CoreTypeDataのスライスをdomain.CoreTypeのスライスに変換します。
func ConvertCoreTypes(types []masterdata.CoreTypeData) []domain.CoreType {
	result := make([]domain.CoreType, len(types))
	for i, t := range types {
		result[i] = t.ToDomain()
	}
	return result
}

// ConvertSkillTypes はmasterdata.ModuleDefinitionDataのスライスをreward.ModuleDropInfoのスライスに変換します。
func ConvertSkillTypes(types []masterdata.ModuleDefinitionData) []rewarding.ModuleDropInfo {
	result := make([]rewarding.ModuleDropInfo, len(types))
	for i, t := range types {
		// マスターデータからドメインモデルのSkillTypeを取得し、そこからEffectsを使う
		moduleType := t.ToDomainType()

		result[i] = rewarding.ModuleDropInfo{
			ID:              t.ID,
			Name:            t.Name,
			Icon:            t.Icon,
			Tags:            t.Tags,
			Description:     t.Description,
			CooldownSeconds: t.CooldownSeconds,
			MinDropLevel:    t.MinDropLevel,
			DifficultyRate:  t.DifficultyRate,
			ChallengeType:   domain.ChallengeTypeID(t.Challenge.Type),
			Effects:         moduleType.Effects,
		}
	}
	return result
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

	// 敵タイプとパッシブスキルを変換（パッシブを解決）
	enemyTypes := ConvertEnemyTypesWithPassives(ext.EnemyTypes, ext.EnemyPassiveSkills)
	coreTypes := ConvertCoreTypes(ext.CoreTypes)
	moduleTypes := ConvertSkillTypes(ext.ModuleDefinitions)

	return enemyTypes, coreTypes, moduleTypes
}

// ConvertChainEffects はmasterdata.ChainEffectDataのスライスをrewarding.ChainEffectDefinitionのスライスに変換します。
func ConvertChainEffects(effects []masterdata.ChainEffectData) []rewarding.ChainEffectDefinition {
	result := make([]rewarding.ChainEffectDefinition, len(effects))
	for i, e := range effects {
		result[i] = rewarding.ChainEffectDefinition{
			ID:           e.ID,
			Name:         e.Name,
			Category:     e.Category,
			EffectType:   e.ToDomainEffectType(),
			MinValue:     e.MinValue,
			MaxValue:     e.MaxValue,
			MinDropLevel: e.MinDropLevel,
		}
	}
	return result
}

// ConvertPassiveSkills はmasterdata.PassiveSkillDataのスライスをdomain.PassiveSkillのマップに変換します。
// キーはパッシブスキルのIDです。
func ConvertPassiveSkills(skills []masterdata.PassiveSkillData) map[string]domain.PassiveSkill {
	result := make(map[string]domain.PassiveSkill, len(skills))
	for _, s := range skills {
		result[s.ID] = s.ToDomain()
	}
	return result
}

// ConvertChainEffectsToMap はmasterdata.ChainEffectDataのスライスをdomain.ChainEffectのマップに変換します。
// キーはチェイン効果のIDです。画面表示用に使用します。
func ConvertChainEffectsToMap(effects []masterdata.ChainEffectData) map[string]domain.ChainEffect {
	result := make(map[string]domain.ChainEffect, len(effects))
	for _, e := range effects {
		result[e.ID] = e.ToDomainChainEffect()
	}
	return result
}
