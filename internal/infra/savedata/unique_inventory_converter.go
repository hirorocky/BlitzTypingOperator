// Package savedata はセーブデータの永続化を担当します。
// このファイルはユニークインベントリとエージェントスロットの変換関数を提供します。

package savedata

// ==================== CoreInventory 変換 ====================

// ConvertCoreInventoryToSave はドメインのCoreInventory形式からセーブ形式に変換します。
// coresはTypeIDリスト形式です。
func ConvertCoreInventoryToSave(cores []string) *CoreInventorySave {
	result := &CoreInventorySave{
		Cores: make([]string, len(cores)),
	}

	copy(result.Cores, cores)

	return result
}

// ConvertSaveToCoreInventory はセーブ形式からドメインのCoreInventory形式に変換します。
// TypeIDリスト形式を返します。
// validTypeIDsに存在しないTypeIDは無視されます（マスタデータ整合性検証）。
// saveがnilの場合は空のスライスを返します。
func ConvertSaveToCoreInventory(save *CoreInventorySave, validTypeIDs map[string]bool) []string {
	if save == nil {
		return []string{}
	}

	result := make([]string, 0, len(save.Cores))

	for _, typeID := range save.Cores {
		// マスタデータに存在するTypeIDのみ復元
		if validTypeIDs[typeID] {
			result = append(result, typeID)
		}
	}

	return result
}

// ==================== SkillInventory 変換 ====================

// ConvertSkillInventoryToSave はドメインのSkillInventory形式からセーブ形式に変換します。
// skillsはTypeIDリスト形式です。
func ConvertSkillInventoryToSave(skills []string) *SkillInventorySave {
	result := &SkillInventorySave{
		Skills: make([]string, len(skills)),
	}
	copy(result.Skills, skills)
	return result
}

// ConvertSaveToSkillInventory はセーブ形式からドメインのSkillInventory形式に変換します。
// validTypeIDsに存在しないTypeIDは無視されます（マスタデータ整合性検証）。
// saveがnilの場合は空のスライスを返します。
func ConvertSaveToSkillInventory(save *SkillInventorySave, validTypeIDs map[string]bool) []string {
	if save == nil {
		return []string{}
	}

	result := make([]string, 0, len(save.Skills))
	for _, typeID := range save.Skills {
		// マスタデータに存在するTypeIDのみ復元
		if validTypeIDs[typeID] {
			result = append(result, typeID)
		}
	}

	return result
}

// ==================== ChainEffectInventory 変換 ====================

// ConvertChainEffectInventoryToSave はドメインのChainEffectInventory形式からセーブ形式に変換します。
// chainEffectsはTypeIDリスト形式です。
func ConvertChainEffectInventoryToSave(chainEffects []string) *ChainEffectInventorySave {
	result := &ChainEffectInventorySave{
		ChainEffects: make([]string, len(chainEffects)),
	}
	copy(result.ChainEffects, chainEffects)
	return result
}

// ConvertSaveToChainEffectInventory はセーブ形式からドメインのChainEffectInventory形式に変換します。
// validTypeIDsに存在しないTypeIDは無視されます（マスタデータ整合性検証）。
// saveがnilの場合は空のスライスを返します。
func ConvertSaveToChainEffectInventory(save *ChainEffectInventorySave, validTypeIDs map[string]bool) []string {
	if save == nil {
		return []string{}
	}

	result := make([]string, 0, len(save.ChainEffects))
	for _, typeID := range save.ChainEffects {
		// マスタデータに存在するTypeIDのみ復元
		if validTypeIDs[typeID] {
			result = append(result, typeID)
		}
	}

	return result
}

// ==================== AgentSlot 変換 ====================

// SkillSlotData はスキルスロットの変換用データ構造です。
type SkillSlotData struct {
	TypeID string
}

// ChainEffectSlotData はチェイン効果スロットの変換用データ構造です。
type ChainEffectSlotData struct {
	TypeID string
}

// ConvertAgentSlotToSave はドメインのAgentSlot形式からセーブ形式に変換します。
func ConvertAgentSlotToSave(coreTypeID string, skills [4]struct {
	TypeID string
}, chainEffects [4]struct {
	TypeID string
}) AgentSlotSave {
	result := AgentSlotSave{
		CoreTypeID: coreTypeID,
	}

	for i, skill := range skills {
		result.Skills[i] = SkillSlotSaveCfg{
			TypeID: skill.TypeID,
		}
	}

	for i, ce := range chainEffects {
		result.ChainEffects[i] = ChainEffectSlotSaveCfg{
			TypeID: ce.TypeID,
		}
	}

	return result
}

// ConvertSaveToAgentSlot はセーブ形式からドメインのAgentSlot形式に変換します。
// validCoreTypeIDs、validSkillTypeIDs、validChainEffectTypeIDsに存在しないTypeIDは無視されます。
// コアが無効な場合はスロット全体が空として扱われます。
func ConvertSaveToAgentSlot(
	save AgentSlotSave,
	validCoreTypeIDs map[string]bool,
	validSkillTypeIDs map[string]bool,
	validChainEffectTypeIDs map[string]bool,
) (coreTypeID string, skills [4]SkillSlotData, chainEffects [4]ChainEffectSlotData) {
	// コアが空の場合または無効な場合はスロット全体が空
	if save.CoreTypeID == "" || !validCoreTypeIDs[save.CoreTypeID] {
		return "", [4]SkillSlotData{}, [4]ChainEffectSlotData{}
	}

	coreTypeID = save.CoreTypeID

	// スキルの変換（無効なTypeIDは空スロットとして扱う）
	for i, skillSave := range save.Skills {
		if skillSave.TypeID != "" && validSkillTypeIDs[skillSave.TypeID] {
			skills[i] = SkillSlotData(skillSave)
		}
	}

	// チェイン効果の変換（無効なTypeIDは空スロットとして扱う）
	for i, ceSave := range save.ChainEffects {
		if ceSave.TypeID != "" && validChainEffectTypeIDs[ceSave.TypeID] {
			chainEffects[i] = ChainEffectSlotData(ceSave)
		}
	}

	return coreTypeID, skills, chainEffects
}

// ==================== EnemyProgress 変換 ====================

// DefeatRecordInput はドメインからセーブ形式への変換入力データです。
type DefeatRecordInput struct {
	Defeated         bool
	MaxDefeatedLevel int
}

// ConvertEnemyProgressToSave はドメインのEnemyProgress形式からセーブ形式に変換します。
// 敵進行データのセーブ形式への変換。
func ConvertEnemyProgressToSave(currentRank int, defeatRecords map[string]struct {
	Defeated         bool
	MaxDefeatedLevel int
}) *EnemyProgressSave {
	result := &EnemyProgressSave{
		CurrentRank:   currentRank,
		DefeatRecords: make(map[string]DefeatRecordSave, len(defeatRecords)),
	}

	for typeID, record := range defeatRecords {
		result.DefeatRecords[typeID] = DefeatRecordSave{
			Defeated:         record.Defeated,
			MaxDefeatedLevel: record.MaxDefeatedLevel,
		}
	}

	return result
}

// ConvertSaveToEnemyProgress はセーブ形式からドメインのEnemyProgress形式に変換します。
// 敵進行データのドメイン形式への変換。
// saveがnilの場合はデフォルト値（rank=1, 空のrecords）を返します。
func ConvertSaveToEnemyProgress(save *EnemyProgressSave) (currentRank int, defeatRecords map[string]DefeatRecordInput) {
	if save == nil {
		return 1, make(map[string]DefeatRecordInput)
	}

	defeatRecords = make(map[string]DefeatRecordInput, len(save.DefeatRecords))

	for typeID, recordSave := range save.DefeatRecords {
		defeatRecords[typeID] = DefeatRecordInput(recordSave)
	}

	return save.CurrentRank, defeatRecords
}
