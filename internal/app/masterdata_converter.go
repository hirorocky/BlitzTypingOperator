// Package app はmasterdata型からdomain型への変換ヘルパー関数を提供します。
// usecase層からinfra層への依存を解消するため、変換処理はapp層で行います。
package app

import (
	"log/slog"

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

// ConvertSkillTypes はmasterdata.SkillDefinitionDataのスライスをreward.SkillDropInfoのスライスに変換します。
func ConvertSkillTypes(types []masterdata.SkillDefinitionData) []rewarding.SkillDropInfo {
	result := make([]rewarding.SkillDropInfo, len(types))
	for i, t := range types {
		// マスターデータからドメインモデルのSkillTypeを取得し、そこからEffectsを使う
		skillType := t.ToDomainType()

		result[i] = rewarding.SkillDropInfo{
			ID:              t.ID,
			Name:            t.Name,
			Icon:            t.Icon,
			Tags:            t.Tags,
			Description:     t.Description,
			CooldownSeconds: t.CooldownSeconds,
			MinDropLevel:    t.MinDropLevel,
			DifficultyRate:  t.DifficultyRate,
			ChallengeType:   domain.ChallengeTypeID(t.Challenge.Type),
			ManaCost:        t.ManaCost,
			Effects:         skillType.Effects,
		}
	}
	return result
}

// ConvertExternalDataToDomain はExternalDataから全てのドメイン型データを変換します。
func ConvertExternalDataToDomain(ext *masterdata.ExternalData) (
	[]domain.EnemyType,
	[]domain.CoreType,
	[]rewarding.SkillDropInfo,
) {
	if ext == nil {
		return nil, nil, nil
	}

	// 敵タイプとパッシブスキルを変換（パッシブを解決）
	enemyTypes := ConvertEnemyTypesWithPassives(ext.EnemyTypes, ext.EnemyPassiveSkills)
	coreTypes := ConvertCoreTypes(ext.CoreTypes)
	skillTypes := ConvertSkillTypes(ext.SkillDefinitions)

	return enemyTypes, coreTypes, skillTypes
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

// ConvertTimedEffects はmasterdata.TimedEffectDataのスライスをdomain.TimedEffectのマップに変換します。
// キーは時限効果のIDです。
func ConvertTimedEffects(types []masterdata.TimedEffectData) map[string]domain.TimedEffect {
	result := make(map[string]domain.TimedEffect, len(types))
	for _, t := range types {
		result[t.ID] = domain.TimedEffect{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Description,
			Column:      domain.EffectColumn(t.EffectColumn),
			Value:       t.Value,
		}
	}
	return result
}

// ResolveSkillTimedEffects はスキルのEffectColumnSpecにTimedEffectからColumn/Valueを解決します。
func ResolveSkillTimedEffects(skills []rewarding.SkillDropInfo, timedEffects map[string]domain.TimedEffect) {
	for i := range skills {
		for j := range skills[i].Effects {
			spec := skills[i].Effects[j].ColumnSpec
			if spec != nil && spec.TimedEffectID != "" {
				if te, ok := timedEffects[spec.TimedEffectID]; ok {
					skills[i].Effects[j].ColumnSpec.Column = te.Column
					skills[i].Effects[j].ColumnSpec.Value = te.Value
				} else {
					slog.Error("マスタデータ不整合: スキルの時限効果IDが見つかりません",
						slog.String("timedEffectID", spec.TimedEffectID),
						slog.String("skillID", skills[i].ID))
				}
			}
		}
	}
}

// ResolveEnemyActionTimedEffects は敵行動のバフ/デバフにTimedEffectからEffectColumn/EffectValueを解決します。
func ResolveEnemyActionTimedEffects(enemyTypes []domain.EnemyType, timedEffects map[string]domain.TimedEffect) {
	for i := range enemyTypes {
		resolveActionsTimedEffects(enemyTypes[i].ResolvedNormalActions, timedEffects)
		resolveActionsTimedEffects(enemyTypes[i].ResolvedEnhancedActions, timedEffects)
	}
}

// resolveActionsTimedEffects は行動リストのTimedEffectを解決するヘルパーです。
func resolveActionsTimedEffects(actions []domain.EnemyAction, timedEffects map[string]domain.TimedEffect) {
	for i := range actions {
		if actions[i].TimedEffectID != "" {
			if te, ok := timedEffects[actions[i].TimedEffectID]; ok {
				actions[i].EffectColumn = te.Column
				actions[i].EffectValue = te.Value
			} else {
				slog.Error("マスタデータ不整合: 敵行動の時限効果IDが見つかりません",
					slog.String("timedEffectID", actions[i].TimedEffectID),
					slog.String("actionID", actions[i].ID))
			}
		}
	}
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

// ConvertFeatureUnlocks はmasterdata.FeatureUnlockDataのスライスをdomain.UnlockRuleのスライスに変換します。
func ConvertFeatureUnlocks(data []masterdata.FeatureUnlockData) []domain.UnlockRule {
	result := make([]domain.UnlockRule, len(data))
	for i, d := range data {
		result[i] = domain.UnlockRule{
			FeatureID:  domain.FeatureID(d.FeatureID),
			UnlockRank: d.UnlockRank,
			TutorialID: d.TutorialID,
		}
	}
	return result
}

// ConvertTutorials はmasterdata.TutorialDataのスライスをdomain.TutorialDefのスライスに変換します。
func ConvertTutorials(data []masterdata.TutorialData) []domain.TutorialDef {
	result := make([]domain.TutorialDef, len(data))
	for i, d := range data {
		pages := make([]string, len(d.Pages))
		copy(pages, d.Pages)
		result[i] = domain.TutorialDef{
			ID:             d.ID,
			Title:          d.Title,
			Pages:          pages,
			DefaultVisible: d.DefaultVisible,
			FeatureID:      domain.FeatureID(d.FeatureID),
		}
	}
	return result
}

// ConvertRankRewards はmasterdata.RankRewardDataのスライスをmap[int]domain.RankRewardに変換します。
// キーはランク番号です。未知のcategoryや不正ランクはログ警告してスキップします。
func ConvertRankRewards(
	data []masterdata.RankRewardData,
	coreTypes []domain.CoreType,
	skillTypes []rewarding.SkillDropInfo,
	chainEffectDefs []rewarding.ChainEffectDefinition,
) map[int]domain.RankReward {
	// 存在確認用マップ
	coreMap := make(map[string]bool, len(coreTypes))
	for _, ct := range coreTypes {
		coreMap[ct.ID] = true
	}
	skillMap := make(map[string]bool, len(skillTypes))
	for _, st := range skillTypes {
		skillMap[st.ID] = true
	}
	chainEffectMap := make(map[string]bool, len(chainEffectDefs))
	for _, ce := range chainEffectDefs {
		chainEffectMap[ce.ID] = true
	}

	result := make(map[int]domain.RankReward)
	for _, rd := range data {
		if rd.Rank <= 0 {
			slog.Warn("不正なランクアップ報酬ランク、スキップします",
				slog.Int("rank", rd.Rank))
			continue
		}

		items := make([]domain.RankRewardItem, 0, len(rd.Rewards))
		for _, item := range rd.Rewards {
			switch item.Category {
			case "core":
				if !coreMap[item.TypeID] {
					slog.Warn("ランクアップ報酬の未知コアTypeID、スキップします",
						slog.String("typeID", item.TypeID),
						slog.Int("rank", rd.Rank))
					continue
				}
			case "skill":
				if !skillMap[item.TypeID] {
					slog.Warn("ランクアップ報酬の未知スキルTypeID、スキップします",
						slog.String("typeID", item.TypeID),
						slog.Int("rank", rd.Rank))
					continue
				}
			case "chain_effect":
				if !chainEffectMap[item.TypeID] {
					slog.Warn("ランクアップ報酬の未知チェイン効果TypeID、スキップします",
						slog.String("typeID", item.TypeID),
						slog.Int("rank", rd.Rank))
					continue
				}
			default:
				slog.Warn("ランクアップ報酬の未知カテゴリ、スキップします",
					slog.String("category", item.Category),
					slog.Int("rank", rd.Rank))
				continue
			}

			items = append(items, domain.RankRewardItem{
				Category: item.Category,
				TypeID:   item.TypeID,
			})
		}

		result[rd.Rank] = domain.RankReward{
			Rank:  rd.Rank,
			Items: items,
		}
	}
	return result
}
