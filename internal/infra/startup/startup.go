// Package startup は初回起動時の初期化処理を担当します。
// 新規ゲーム開始時の初期エージェントの提供を行います。

package startup

import (
	"log/slog"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
)

// NewGameInitializer は新規ゲーム初期化を担当する構造体です。
type NewGameInitializer struct {
	// externalData は外部マスタデータです。
	externalData *masterdata.ExternalData
}

// NewNewGameInitializer は新しいNewGameInitializerを作成します。
// externalData はマスタデータを含む外部データです。
func NewNewGameInitializer(externalData *masterdata.ExternalData) *NewGameInitializer {
	return &NewGameInitializer{
		externalData: externalData,
	}
}

// CreateInitialAgents は初期エージェントを作成します。
// マスタデータから初期エージェントを構築します。
func (i *NewGameInitializer) CreateInitialAgents() []*domain.AgentModel {
	if i.externalData == nil || len(i.externalData.FirstAgents) == 0 {
		slog.Error("初期エージェントデータがありません")
		return nil
	}

	agents := make([]*domain.AgentModel, 0, len(i.externalData.FirstAgents))

	for _, firstAgentData := range i.externalData.FirstAgents {
		// コア特性を検索
		var coreType domain.CoreType
		for _, ct := range i.externalData.CoreTypes {
			if ct.ID == firstAgentData.CoreTypeID {
				coreType = ct.ToDomain()
				break
			}
		}

		// パッシブスキルを検索
		var passiveSkill domain.PassiveSkill
		for _, ps := range i.externalData.PassiveSkills {
			if ps.ID == coreType.PassiveSkillID {
				passiveSkill = domain.PassiveSkill{
					ID:          ps.ID,
					Name:        ps.Name,
					Description: ps.Description,
				}
				break
			}
		}

		// コアを作成（レベルなし）
		core := domain.NewCoreWithTypeID(
			firstAgentData.CoreTypeID,
			coreType,
			passiveSkill,
		)

		// モジュールを作成
		modules := make([]*domain.SkillModel, 0, len(firstAgentData.Modules))
		for _, modData := range firstAgentData.Modules {
			// モジュール定義を検索
			var moduleDef *masterdata.ModuleDefinitionData
			for j := range i.externalData.ModuleDefinitions {
				if i.externalData.ModuleDefinitions[j].ID == modData.TypeID {
					moduleDef = &i.externalData.ModuleDefinitions[j]
					break
				}
			}
			if moduleDef == nil {
				slog.Warn("モジュール定義が見つかりません",
					slog.String("type_id", modData.TypeID),
				)
				continue
			}

			// チェイン効果を作成
			var chainEffect *domain.ChainEffect
			if modData.HasChainEffect() {
				// チェイン効果定義を検索して説明文テンプレートを取得
				var chainEffectDef *masterdata.ChainEffectData
				for j := range i.externalData.ChainEffects {
					if i.externalData.ChainEffects[j].EffectType == modData.ChainEffectType {
						chainEffectDef = &i.externalData.ChainEffects[j]
						break
					}
				}
				if chainEffectDef != nil {
					ce := domain.NewChainEffectWithTemplate(
						chainEffectDef.ID,
						convertChainEffectType(modData.ChainEffectType),
						modData.ChainEffectValue,
						chainEffectDef.Description,
						chainEffectDef.ShortDescription,
					)
					chainEffect = &ce
				} else {
					slog.Warn("チェイン効果定義が見つかりません",
						slog.String("effect_type", modData.ChainEffectType),
					)
				}
			}

			// モジュールを作成
			module := domain.NewSkillFromType(moduleDef.ToDomainType(), chainEffect)
			modules = append(modules, module)
		}

		agent := domain.NewAgent(firstAgentData.ID, core, modules)
		agents = append(agents, agent)
	}

	return agents
}

// convertChainEffectType は文字列をChainEffectTypeに変換します。
func convertChainEffectType(s string) domain.ChainEffectType {
	switch s {
	case "damage_bonus":
		return domain.ChainEffectDamageBonus
	case "heal_bonus":
		return domain.ChainEffectHealBonus
	case "buff_extend":
		return domain.ChainEffectBuffExtend
	case "debuff_extend":
		return domain.ChainEffectDebuffExtend
	case "damage_amp":
		return domain.ChainEffectDamageAmp
	case "armor_pierce":
		return domain.ChainEffectArmorPierce
	case "life_steal":
		return domain.ChainEffectLifeSteal
	case "damage_cut":
		return domain.ChainEffectDamageCut
	case "evasion":
		return domain.ChainEffectEvasion
	case "reflect":
		return domain.ChainEffectReflect
	case "regen":
		return domain.ChainEffectRegen
	case "heal_amp":
		return domain.ChainEffectHealAmp
	case "overheal":
		return domain.ChainEffectOverheal
	case "time_extend":
		return domain.ChainEffectTimeExtend
	case "auto_correct":
		return domain.ChainEffectAutoCorrect
	case "cooldown_reduce":
		return domain.ChainEffectCooldownReduce
	case "buff_duration":
		return domain.ChainEffectBuffDuration
	case "debuff_duration":
		return domain.ChainEffectDebuffDuration
	case "double_cast":
		return domain.ChainEffectDoubleCast
	default:
		return domain.ChainEffectDamageBonus
	}
}

// InitializeNewGame は新規ゲームを初期化してセーブデータを作成します。
// ユニークインベントリ + エージェントスロットシステム
func (i *NewGameInitializer) InitializeNewGame() *savedata.SaveData {
	// 基本のセーブデータを作成
	saveData := savedata.NewSaveData()

	// 初期コアとスキルをユニークインベントリに追加
	// コア：オールラウンダー（TypeIDリスト形式）
	// スキル：軽斬撃、応急手当、気合い溜め
	saveData.Inventory.UniqueCores.Cores = append(saveData.Inventory.UniqueCores.Cores, "all_rounder")

	saveData.Inventory.UniqueSkills.Skills["physical_strike_lv1"] = []string{}
	saveData.Inventory.UniqueSkills.Skills["heal_lv1"] = []string{}
	saveData.Inventory.UniqueSkills.Skills["str_buff_lv1"] = []string{}

	// 初期エージェントスロット構成
	// 3スロットにオールラウンダー + 1スキルずつ分散
	saveData.Player.AgentSlots[0] = savedata.AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]savedata.SkillSlotSaveCfg{
			{TypeID: "physical_strike_lv1"},
			{},
			{},
			{},
		},
	}
	saveData.Player.AgentSlots[1] = savedata.AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]savedata.SkillSlotSaveCfg{
			{TypeID: "heal_lv1"},
			{},
			{},
			{},
		},
	}
	saveData.Player.AgentSlots[2] = savedata.AgentSlotSave{
		CoreTypeID: "all_rounder",
		Skills: [4]savedata.SkillSlotSaveCfg{
			{TypeID: "str_buff_lv1"},
			{},
			{},
			{},
		},
	}

	return saveData
}

// CreateNewGameWithExtraItems は追加アイテム付きで新規ゲームを初期化します。
// デバッグや特殊条件での開始用
func (i *NewGameInitializer) CreateNewGameWithExtraItems() *savedata.SaveData {
	saveData := i.InitializeNewGame()

	// 追加のエージェントから情報を取得（初期エージェントと同じ構成を使用）
	extraAgents := i.CreateInitialAgents()
	if len(extraAgents) == 0 {
		return saveData
	}

	// 最初のエージェントから追加素材を作成
	extraAgent := extraAgents[0]

	// 追加のコアをユニークコアインベントリに追加
	saveData.Inventory.UniqueCores.Cores = append(
		saveData.Inventory.UniqueCores.Cores,
		extraAgent.Core.TypeID,
	)

	// 追加のスキルをユニークスキルインベントリに追加
	for _, skill := range extraAgent.Modules {
		chainEffectID := ""
		if skill.ChainEffect != nil {
			chainEffectID = string(skill.ChainEffect.Type)
		}
		if saveData.Inventory.UniqueSkills.Skills[skill.TypeID] == nil {
			saveData.Inventory.UniqueSkills.Skills[skill.TypeID] = []string{}
		}
		saveData.Inventory.UniqueSkills.Skills[skill.TypeID] = append(
			saveData.Inventory.UniqueSkills.Skills[skill.TypeID],
			chainEffectID,
		)
	}

	return saveData
}
