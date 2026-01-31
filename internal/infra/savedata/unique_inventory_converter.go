// Package savedata はセーブデータの永続化を担当します。
// このファイルはユニークインベントリとエージェントスロットの変換関数を提供します。

package savedata

// ==================== CoreInventory 変換 ====================

// ConvertCoreInventoryToSave はドメインのCoreInventory形式からセーブ形式に変換します。
// coresは map[TypeID]MaxLevel 形式です。
func ConvertCoreInventoryToSave(cores map[string]int) *CoreInventorySave {
	result := &CoreInventorySave{
		Cores: make(map[string]int, len(cores)),
	}

	for typeID, level := range cores {
		result.Cores[typeID] = level
	}

	return result
}

// ConvertSaveToCoreInventory はセーブ形式からドメインのCoreInventory形式に変換します。
// validTypeIDsに存在しないTypeIDは無視されます（マスタデータ整合性検証）。
// saveがnilの場合は空のマップを返します。
func ConvertSaveToCoreInventory(save *CoreInventorySave, validTypeIDs map[string]bool) map[string]int {
	result := make(map[string]int)

	if save == nil {
		return result
	}

	for typeID, level := range save.Cores {
		// マスタデータに存在するTypeIDのみ復元
		if validTypeIDs[typeID] {
			result[typeID] = level
		}
	}

	return result
}

// ==================== SkillInventory 変換 ====================

// ConvertSkillInventoryToSave はドメインのSkillInventory形式からセーブ形式に変換します。
// skillsは map[TypeID][]ChainEffectID 形式です。
func ConvertSkillInventoryToSave(skills map[string][]string) *SkillInventorySave {
	result := &SkillInventorySave{
		Skills: make(map[string][]string, len(skills)),
	}

	for typeID, chainEffects := range skills {
		// チェイン効果リストをコピー
		chainsCopy := make([]string, len(chainEffects))
		copy(chainsCopy, chainEffects)
		result.Skills[typeID] = chainsCopy
	}

	return result
}

// ConvertSaveToSkillInventory はセーブ形式からドメインのSkillInventory形式に変換します。
// validTypeIDsに存在しないTypeIDは無視されます（マスタデータ整合性検証）。
// saveがnilの場合は空のマップを返します。
func ConvertSaveToSkillInventory(save *SkillInventorySave, validTypeIDs map[string]bool) map[string][]string {
	result := make(map[string][]string)

	if save == nil {
		return result
	}

	for typeID, chainEffects := range save.Skills {
		// マスタデータに存在するTypeIDのみ復元
		if validTypeIDs[typeID] {
			chainsCopy := make([]string, len(chainEffects))
			copy(chainsCopy, chainEffects)
			result[typeID] = chainsCopy
		}
	}

	return result
}

// ==================== AgentSlot 変換 ====================

// SkillSlotData はスキルスロットの変換用データ構造です。
type SkillSlotData struct {
	TypeID        string
	ChainEffectID string
}

// ConvertAgentSlotToSave はドメインのAgentSlot形式からセーブ形式に変換します。
func ConvertAgentSlotToSave(coreTypeID string, coreLevel int, skills [4]struct {
	TypeID        string
	ChainEffectID string
}) AgentSlotSave {
	result := AgentSlotSave{
		CoreTypeID: coreTypeID,
		CoreLevel:  coreLevel,
	}

	for i, skill := range skills {
		result.Skills[i] = SkillSlotSaveCfg{
			TypeID:        skill.TypeID,
			ChainEffectID: skill.ChainEffectID,
		}
	}

	return result
}

// ConvertSaveToAgentSlot はセーブ形式からドメインのAgentSlot形式に変換します。
// validCoreTypeIDsまたはvalidSkillTypeIDsに存在しないTypeIDは無視されます。
// コアが無効な場合はスロット全体が空として扱われます。
func ConvertSaveToAgentSlot(
	save AgentSlotSave,
	validCoreTypeIDs map[string]bool,
	validSkillTypeIDs map[string]bool,
) (coreTypeID string, coreLevel int, skills [4]SkillSlotData) {
	// コアが空の場合または無効な場合はスロット全体が空
	if save.CoreTypeID == "" || !validCoreTypeIDs[save.CoreTypeID] {
		return "", 0, [4]SkillSlotData{}
	}

	coreTypeID = save.CoreTypeID
	coreLevel = save.CoreLevel

	// スキルの変換（無効なTypeIDは空スロットとして扱う）
	for i, skillSave := range save.Skills {
		if skillSave.TypeID != "" && validSkillTypeIDs[skillSave.TypeID] {
			skills[i] = SkillSlotData(skillSave)
		}
		// 無効なTypeIDは初期値（空）のまま
	}

	return coreTypeID, coreLevel, skills
}
