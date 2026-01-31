// Package masterdata はマスタデータのロード処理を提供します。
// コア特性、モジュール定義、敵タイプ定義、タイピング辞書などを
// JSONファイルから読み込みます。

package masterdata

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// DataLoader は外部データファイルのロードを担当する構造体です。
type DataLoader struct {
	// dataDir はデータファイルが格納されているディレクトリパスです。
	// 外部ディレクトリ使用時のみ設定されます。
	dataDir string
	// fileSystem は埋め込みまたはOSファイルシステムです。
	fileSystem fs.FS
	// subDir はfs.FS内でのサブディレクトリパスです（埋め込み時は "data"）。
	subDir string
}

// NewDataLoader は外部ディレクトリから読み込むDataLoaderを作成します。
func NewDataLoader(dataDir string) *DataLoader {
	return &DataLoader{
		dataDir:    dataDir,
		fileSystem: nil,
		subDir:     "",
	}
}

// NewEmbeddedDataLoader は埋め込みFSから読み込むDataLoaderを作成します。
// embeddedFS は embed.FS などの fs.FS 実装です。
// subDir は埋め込みFS内でのサブディレクトリパスです（通常は "data"）。
func NewEmbeddedDataLoader(embeddedFS fs.FS, subDir string) *DataLoader {
	return &DataLoader{
		dataDir:    "",
		fileSystem: embeddedFS,
		subDir:     subDir,
	}
}

// readFile はファイルを読み込むヘルパーメソッドです。
// 外部ディレクトリまたは埋め込みFSから読み込みます。
func (l *DataLoader) readFile(filename string) ([]byte, error) {
	if l.dataDir != "" {
		// 外部ディレクトリから読み込み
		filePath := filepath.Join(l.dataDir, filename)
		return os.ReadFile(filePath)
	}

	// 埋め込みFSから読み込み
	var filePath string
	if l.subDir != "" {
		filePath = l.subDir + "/" + filename
	} else {
		filePath = filename
	}
	return fs.ReadFile(l.fileSystem, filePath)
}

// ExternalData は外部データファイルから読み込んだ全データを格納する構造体です。
type ExternalData struct {
	CoreTypes          []CoreTypeData
	ModuleDefinitions  []ModuleDefinitionData
	EnemyTypes         []EnemyTypeData
	EnemyActions       []EnemyActionData
	EnemyPassiveSkills []EnemyPassiveSkillData
	PassiveSkills      []PassiveSkillData
	ChainEffects       []ChainEffectData
	TypingDictionary   *TypingDictionary
	FirstAgents        []FirstAgentData
}

// ==================== コア特性定義 ====================

// CoreTypeData はcores.jsonから読み込むコア特性データの構造体です。
type CoreTypeData struct {
	ID             string             `json:"id"`
	Name           string             `json:"name"`
	AllowedTags    []string           `json:"allowed_tags"`
	StatWeights    map[string]float64 `json:"stat_weights"`
	PassiveSkillID string             `json:"passive_skill_id"`
	MinDropLevel   int                `json:"min_drop_level"`
}

// coresFileData はcores.jsonのルート構造です。
type coresFileData struct {
	CoreTypes []CoreTypeData `json:"core_types"`
}

// LoadCoreTypes はcores.jsonからコア特性定義を読み込みます。

func (l *DataLoader) LoadCoreTypes() ([]CoreTypeData, error) {
	data, err := l.readFile("cores.json")
	if err != nil {
		return nil, fmt.Errorf("cores.jsonの読み込みに失敗: %w", err)
	}

	var fileData coresFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("cores.jsonのパースに失敗: %w", err)
	}

	return fileData.CoreTypes, nil
}

// ToDomain はCoreTypeDataをドメインモデルのCoreTypeに変換します。
func (c *CoreTypeData) ToDomain() domain.CoreType {
	// AllowedTagsをコピー
	allowedTags := make([]string, len(c.AllowedTags))
	copy(allowedTags, c.AllowedTags)

	// StatWeightsをコピー
	statWeights := make(map[string]float64)
	for k, v := range c.StatWeights {
		statWeights[k] = v
	}

	return domain.CoreType{
		ID:             c.ID,
		Name:           c.Name,
		StatWeights:    statWeights,
		PassiveSkillID: c.PassiveSkillID,
		AllowedTags:    allowedTags,
		MinDropLevel:   c.MinDropLevel,
	}
}

// ==================== スキル（モジュール）定義 ====================

// HPFormulaData はHP増減計算式のJSONデータ構造体です。
type HPFormulaData struct {
	Base     float64 `json:"base"`
	StatCoef float64 `json:"stat_coef"`
	StatRef  string  `json:"stat_ref"`
}

// EffectColumnData はEffectColumn効果のJSONデータ構造体です。
type EffectColumnData struct {
	Column   string  `json:"column"`
	Value    float64 `json:"value"`
	Duration float64 `json:"duration"`
}

// SkillEffectData はスキル効果のJSONデータ構造体です。
type SkillEffectData struct {
	Target       string            `json:"target"`
	HPFormula    *HPFormulaData    `json:"hp_formula,omitempty"`
	EffectColumn *EffectColumnData `json:"effect_column,omitempty"`
	Probability  float64           `json:"probability"`
	LUKFactor    float64           `json:"luk_factor"`
	Icon         string            `json:"icon"`
}

// ModuleEffectData はSkillEffectDataのエイリアスです。
// 後方互換性のために残されています。新規コードではSkillEffectDataを使用してください。
type ModuleEffectData = SkillEffectData

// SkillDefinitionData はskills.json（またはmodules.json）から読み込むスキル定義データの構造体です。
type SkillDefinitionData struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	Icon            string            `json:"icon"`
	Tags            []string          `json:"tags"`
	Description     string            `json:"description"`
	CooldownSeconds float64           `json:"cooldown_seconds"`
	Difficulty      int               `json:"difficulty"`
	MinDropLevel    int               `json:"min_drop_level"`
	Effects         []SkillEffectData `json:"effects"`
}

// ModuleDefinitionData はSkillDefinitionDataのエイリアスです。
// 後方互換性のために残されています。新規コードではSkillDefinitionDataを使用してください。
type ModuleDefinitionData = SkillDefinitionData

// modulesFileData はmodules.jsonのルート構造です。
type modulesFileData struct {
	ModuleTypes []ModuleDefinitionData `json:"module_types"`
}

// LoadModuleDefinitions はmodules.jsonからモジュール定義を読み込みます。

func (l *DataLoader) LoadModuleDefinitions() ([]ModuleDefinitionData, error) {
	data, err := l.readFile("modules.json")
	if err != nil {
		return nil, fmt.Errorf("modules.jsonの読み込みに失敗: %w", err)
	}

	var fileData modulesFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("modules.jsonのパースに失敗: %w", err)
	}

	return fileData.ModuleTypes, nil
}

// ToDomainType はSkillDefinitionDataをドメインモデルのSkillTypeに変換します。
func (m *SkillDefinitionData) ToDomainType() domain.SkillType {
	// Tagsをコピー（スライスの参照共有を避ける）
	tagsCopy := make([]string, len(m.Tags))
	copy(tagsCopy, m.Tags)

	// Effectsを変換
	effects := make([]domain.SkillEffect, len(m.Effects))
	for i, e := range m.Effects {
		effects[i] = e.ToDomain()
	}

	return domain.SkillType{
		ID:              m.ID,
		Name:            m.Name,
		Icon:            m.Icon,
		Tags:            tagsCopy,
		Description:     m.Description,
		CooldownSeconds: m.CooldownSeconds,
		Difficulty:      m.Difficulty,
		MinDropLevel:    m.MinDropLevel,
		Effects:         effects,
	}
}

// ToDomain はSkillEffectDataをドメインモデルのSkillEffectに変換します。
func (e *SkillEffectData) ToDomain() domain.SkillEffect {
	effect := domain.SkillEffect{
		Target:      convertEffectTarget(e.Target),
		Probability: e.Probability,
		LUKFactor:   e.LUKFactor,
		Icon:        e.Icon,
	}

	if e.HPFormula != nil {
		effect.HPFormula = &domain.HPFormula{
			Base:     e.HPFormula.Base,
			StatCoef: e.HPFormula.StatCoef,
			StatRef:  e.HPFormula.StatRef,
		}
	}

	if e.EffectColumn != nil {
		effect.ColumnSpec = &domain.EffectColumnSpec{
			Column:   domain.EffectColumn(e.EffectColumn.Column),
			Value:    e.EffectColumn.Value,
			Duration: e.EffectColumn.Duration,
		}
	}

	return effect
}

// convertEffectTarget は文字列をEffectTargetに変換します。
func convertEffectTarget(s string) domain.EffectTarget {
	switch s {
	case "self":
		return domain.TargetSelf
	case "enemy":
		return domain.TargetEnemy
	case "both":
		return domain.TargetBoth
	default:
		return domain.TargetEnemy
	}
}

// ToDomain はSkillDefinitionDataをドメインモデルのSkillModelに変換します。
// チェイン効果なしのスキルを作成します。
func (m *SkillDefinitionData) ToDomain() *domain.SkillModel {
	skillType := m.ToDomainType()
	return domain.NewSkillFromType(skillType, nil)
}

// ==================== 敵タイプ定義 ====================

// EnemyTypeData はenemies.jsonから読み込む敵タイプデータの構造体です。
type EnemyTypeData struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	BaseHP          int    `json:"base_hp"`
	BaseAttackPower int    `json:"base_attack_power"`
	AttackType      string `json:"attack_type"`
	ASCIIArt        string `json:"ascii_art"`

	// 拡張フィールド
	DefaultLevel             int      `json:"default_level"`
	NormalActionPatternIDs   []string `json:"normal_action_pattern"`
	EnhancedActionPatternIDs []string `json:"enhanced_action_pattern"`
	NormalPassiveID          string   `json:"normal_passive_id"`
	EnhancedPassiveID        string   `json:"enhanced_passive_id"`
	DropItemCategory         string   `json:"drop_item_category"`
	DropItemTypeID           string   `json:"drop_item_type_id"`

	// ボルテージシステム
	// VoltageRisePer10s は10秒あたりのボルテージ上昇量です。
	// 0の場合はボルテージが上昇しません。未設定時のデフォルト値は10です。
	VoltageRisePer10s *float64 `json:"voltage_rise_per_10s,omitempty"`

	// ランクシステム
	// Rank は敵のランクです（1から開始）。未設定時のデフォルト値は1です。
	Rank *int `json:"rank,omitempty"`
}

// enemiesFileData はenemies.jsonのルート構造です。
type enemiesFileData struct {
	EnemyTypes []EnemyTypeData `json:"enemy_types"`
}

// LoadEnemyTypes はenemies.jsonから敵タイプ定義を読み込みます。

func (l *DataLoader) LoadEnemyTypes() ([]EnemyTypeData, error) {
	data, err := l.readFile("enemies.json")
	if err != nil {
		return nil, fmt.Errorf("enemies.jsonの読み込みに失敗: %w", err)
	}

	var fileData enemiesFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("enemies.jsonのパースに失敗: %w", err)
	}

	return fileData.EnemyTypes, nil
}

// DefaultVoltageRisePer10s はボルテージ上昇量のデフォルト値（10ポイント/10秒）です。
const DefaultVoltageRisePer10s = 10.0

// DefaultEnemyRank は敵ランクのデフォルト値（ランク1）です。
const DefaultEnemyRank = 1

// ToDomain はEnemyTypeDataをドメインモデルのEnemyTypeに変換します。
// actionMap が指定された場合、行動パターンIDを解決します。
// VoltageRisePer10sが未設定の場合はデフォルト値（10）を適用します。
// Rankが未設定の場合はデフォルト値（1）を適用します。
func (e *EnemyTypeData) ToDomain() domain.EnemyType {
	return domain.EnemyType{
		ID:                       e.ID,
		Name:                     e.Name,
		BaseHP:                   e.BaseHP,
		BaseAttackPower:          e.BaseAttackPower,
		AttackType:               e.AttackType,
		ASCIIArt:                 e.ASCIIArt,
		DefaultLevel:             e.DefaultLevel,
		NormalActionPatternIDs:   e.NormalActionPatternIDs,
		EnhancedActionPatternIDs: e.EnhancedActionPatternIDs,
		DropItemCategory:         e.DropItemCategory,
		DropItemTypeID:           e.DropItemTypeID,
		VoltageRisePer10s:        e.GetVoltageRisePer10s(),
		Rank:                     e.GetRank(),
	}
}

// GetVoltageRisePer10s は10秒あたりのボルテージ上昇量を返します。
// 未設定の場合はデフォルト値を返します。
func (e *EnemyTypeData) GetVoltageRisePer10s() float64 {
	if e.VoltageRisePer10s == nil {
		return DefaultVoltageRisePer10s
	}
	return *e.VoltageRisePer10s
}

// GetRank は敵のランクを返します。
// 未設定の場合はデフォルト値（1）を返します。
func (e *EnemyTypeData) GetRank() int {
	if e.Rank == nil {
		return DefaultEnemyRank
	}
	return *e.Rank
}

// ==================== 敵行動定義 ====================

// EnemyActionData はenemy_actions.jsonから読み込む敵行動データの構造体です。
type EnemyActionData struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	ActionType     string  `json:"action_type"`
	AttackType     string  `json:"attack_type,omitempty"`
	DamageBase     float64 `json:"damage_base,omitempty"`
	DamagePerLevel float64 `json:"damage_per_level,omitempty"`
	Element        string  `json:"element,omitempty"`
	EffectType     string  `json:"effect_type,omitempty"`
	EffectValue    float64 `json:"effect_value,omitempty"`
	DefenseType    string  `json:"defense_type,omitempty"`
	ReductionRate  float64 `json:"reduction_rate,omitempty"`
	EvadeRate      float64 `json:"evade_rate,omitempty"`
	DurationSec    float64 `json:"duration_seconds,omitempty"`
	ChargeTimeMS   int64   `json:"charge_time_ms"`
}

// enemyActionsFileData はenemy_actions.jsonのルート構造です。
type enemyActionsFileData struct {
	EnemyActions []EnemyActionData `json:"enemy_actions"`
}

// LoadEnemyActions はenemy_actions.jsonから敵行動定義を読み込みます。
func (l *DataLoader) LoadEnemyActions() ([]EnemyActionData, error) {
	data, err := l.readFile("enemy_actions.json")
	if err != nil {
		return nil, fmt.Errorf("enemy_actions.jsonの読み込みに失敗: %w", err)
	}

	var fileData enemyActionsFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("enemy_actions.jsonのパースに失敗: %w", err)
	}

	return fileData.EnemyActions, nil
}

// ToDomain はEnemyActionDataをドメインモデルのEnemyActionに変換します。
func (a *EnemyActionData) ToDomain() domain.EnemyAction {
	action := domain.EnemyAction{
		ID:             a.ID,
		Name:           a.Name,
		ChargeTime:     time.Duration(a.ChargeTimeMS) * time.Millisecond,
		AttackType:     a.AttackType,
		DamageBase:     a.DamageBase,
		DamagePerLevel: a.DamagePerLevel,
		Element:        a.Element,
		EffectType:     a.EffectType,
		EffectValue:    a.EffectValue,
		Duration:       a.DurationSec,
		ReductionRate:  a.ReductionRate,
		EvadeRate:      a.EvadeRate,
	}

	// ActionTypeの変換
	switch a.ActionType {
	case "attack":
		action.ActionType = domain.EnemyActionAttack
	case "buff":
		action.ActionType = domain.EnemyActionBuff
	case "debuff":
		action.ActionType = domain.EnemyActionDebuff
	case "defense":
		action.ActionType = domain.EnemyActionDefense
	}

	// DefenseTypeの変換
	switch a.DefenseType {
	case "physical_cut":
		action.DefenseType = domain.DefensePhysicalCut
	case "magic_cut":
		action.DefenseType = domain.DefenseMagicCut
	case "debuff_evade":
		action.DefenseType = domain.DefenseDebuffEvade
	}

	return action
}

// ValidateEnemyActionData は敵行動データのバリデーションを行います。
func ValidateEnemyActionData(data EnemyActionData) error {
	if data.ID == "" {
		return fmt.Errorf("敵行動IDが空です")
	}
	if data.ActionType == "" {
		return fmt.Errorf("敵行動タイプが空です: ID=%s", data.ID)
	}
	validTypes := map[string]bool{"attack": true, "buff": true, "debuff": true, "defense": true}
	if !validTypes[data.ActionType] {
		return fmt.Errorf("敵行動タイプが不正です: ID=%s, ActionType=%s", data.ID, data.ActionType)
	}
	return nil
}

// ==================== パッシブスキル定義 ====================

// PassiveSkillData はpassive_skills.jsonから読み込むパッシブスキルデータの構造体です。
type PassiveSkillData struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	TriggerType      string                `json:"trigger_type"`
	TriggerCondition *TriggerConditionData `json:"trigger_condition"`
	EffectType       string                `json:"effect_type"`
	EffectValue      float64               `json:"effect_value"`
	Probability      float64               `json:"probability"`
	MaxStacks        int                   `json:"max_stacks"`
	StackIncrement   float64               `json:"stack_increment"`
	UsesPerBattle    int                   `json:"uses_per_battle"`
}

// TriggerConditionData はトリガー条件のJSONデータ構造体です。
type TriggerConditionData struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

// passiveSkillsFileData はpassive_skills.jsonのルート構造です。
type passiveSkillsFileData struct {
	PassiveSkills []PassiveSkillData `json:"passive_skills"`
}

// LoadPassiveSkills はpassive_skills.jsonからパッシブスキル定義を読み込みます。
func (l *DataLoader) LoadPassiveSkills() ([]PassiveSkillData, error) {
	data, err := l.readFile("passive_skills.json")
	if err != nil {
		return nil, fmt.Errorf("passive_skills.jsonの読み込みに失敗: %w", err)
	}

	var fileData passiveSkillsFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("passive_skills.jsonのパースに失敗: %w", err)
	}

	return fileData.PassiveSkills, nil
}

// ToDomain はPassiveSkillDataをドメインモデルのPassiveSkillに変換します。
func (p *PassiveSkillData) ToDomain() domain.PassiveSkill {
	skill := domain.PassiveSkill{
		ID:               p.ID,
		Name:             p.Name,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		TriggerType:      convertTriggerType(p.TriggerType),
		EffectType:       convertEffectType(p.EffectType),
		EffectValue:      p.EffectValue,
		Probability:      p.Probability,
		MaxStacks:        p.MaxStacks,
		StackIncrement:   p.StackIncrement,
		UsesPerBattle:    p.UsesPerBattle,
	}

	if p.TriggerCondition != nil {
		skill.TriggerCondition = &domain.TriggerCondition{
			Type:  convertTriggerConditionType(p.TriggerCondition.Type),
			Value: p.TriggerCondition.Value,
		}
	}

	return skill
}

// convertTriggerType は文字列をPassiveTriggerTypeに変換します。
func convertTriggerType(s string) domain.PassiveTriggerType {
	switch s {
	case "permanent":
		return domain.PassiveTriggerPermanent
	case "conditional":
		return domain.PassiveTriggerConditional
	case "probability":
		return domain.PassiveTriggerProbability
	case "stack":
		return domain.PassiveTriggerStack
	case "reactive":
		return domain.PassiveTriggerReactive
	default:
		return domain.PassiveTriggerPermanent
	}
}

// convertEffectType は文字列をPassiveEffectTypeに変換します。
func convertEffectType(s string) domain.PassiveEffectType {
	switch s {
	case "modifier":
		return domain.PassiveEffectModifier
	case "multiplier":
		return domain.PassiveEffectMultiplier
	case "special":
		return domain.PassiveEffectSpecial
	default:
		return domain.PassiveEffectModifier
	}
}

// convertTriggerConditionType は文字列をTriggerConditionTypeに変換します。
func convertTriggerConditionType(s string) domain.TriggerConditionType {
	switch s {
	case "accuracy_equals":
		return domain.TriggerConditionAccuracyEquals
	case "wpm_above":
		return domain.TriggerConditionWPMAbove
	case "hp_below_percent":
		return domain.TriggerConditionHPBelowPercent
	case "enemy_hp_below_percent":
		return domain.TriggerConditionEnemyHPBelowPercent
	case "enemy_has_debuff":
		return domain.TriggerConditionEnemyHasDebuff
	case "on_skill_use":
		return domain.TriggerConditionOnSkillUse
	case "on_damage_received":
		return domain.TriggerConditionOnDamageReceived
	case "on_heal":
		return domain.TriggerConditionOnHeal
	case "on_buff_debuff_use":
		return domain.TriggerConditionOnBuffDebuffUse
	case "on_physical_attack":
		return domain.TriggerConditionOnPhysicalAttack
	case "on_typing_miss":
		return domain.TriggerConditionOnTypingMiss
	case "on_timeout":
		return domain.TriggerConditionOnTimeout
	case "on_debuff_received":
		return domain.TriggerConditionOnDebuffReceived
	case "on_battle_start":
		return domain.TriggerConditionOnBattleStart
	case "no_miss_streak":
		return domain.TriggerConditionNoMissStreak
	case "same_attack_count":
		return domain.TriggerConditionSameAttackCount
	default:
		return domain.TriggerConditionAccuracyEquals
	}
}

// ==================== 敵パッシブスキル定義 ====================

// EnemyPassiveSkillData はenemy_passive_skills.jsonから読み込む敵パッシブスキルデータの構造体です。
type EnemyPassiveSkillData struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Effects     map[string]float64 `json:"effects"`
}

// enemyPassiveSkillsFileData はenemy_passive_skills.jsonのルート構造です。
type enemyPassiveSkillsFileData struct {
	EnemyPassiveSkills []EnemyPassiveSkillData `json:"enemy_passive_skills"`
}

// LoadEnemyPassiveSkills はenemy_passive_skills.jsonから敵パッシブスキル定義を読み込みます。
func (l *DataLoader) LoadEnemyPassiveSkills() ([]EnemyPassiveSkillData, error) {
	data, err := l.readFile("enemy_passive_skills.json")
	if err != nil {
		return nil, fmt.Errorf("enemy_passive_skills.jsonの読み込みに失敗: %w", err)
	}

	var fileData enemyPassiveSkillsFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("enemy_passive_skills.jsonのパースに失敗: %w", err)
	}

	return fileData.EnemyPassiveSkills, nil
}

// ToDomain はEnemyPassiveSkillDataをドメインモデルのEnemyPassiveSkillに変換します。
func (p *EnemyPassiveSkillData) ToDomain() *domain.EnemyPassiveSkill {
	effects := make(map[domain.EffectColumn]float64)
	for key, value := range p.Effects {
		effects[domain.EffectColumn(key)] = value
	}

	return &domain.EnemyPassiveSkill{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Effects:     effects,
	}
}

// ==================== チェイン効果定義 ====================

// ChainEffectData はchain_effects.jsonから読み込むチェイン効果データの構造体です。
type ChainEffectData struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ShortDescription string  `json:"short_description"`
	Category         string  `json:"category"`
	EffectType       string  `json:"effect_type"`
	MinValue         float64 `json:"min_value"`
	MaxValue         float64 `json:"max_value"`
	MinDropLevel     int     `json:"min_drop_level"`
}

// chainEffectsFileData はchain_effects.jsonのルート構造です。
type chainEffectsFileData struct {
	ChainEffects []ChainEffectData `json:"chain_effects"`
}

// LoadChainEffects はchain_effects.jsonからチェイン効果定義を読み込みます。
func (l *DataLoader) LoadChainEffects() ([]ChainEffectData, error) {
	data, err := l.readFile("chain_effects.json")
	if err != nil {
		return nil, fmt.Errorf("chain_effects.jsonの読み込みに失敗: %w", err)
	}

	var fileData chainEffectsFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("chain_effects.jsonのパースに失敗: %w", err)
	}

	return fileData.ChainEffects, nil
}

// ToDomainEffectType はChainEffectDataからドメインモデルのChainEffectTypeに変換します。
func (s *ChainEffectData) ToDomainEffectType() domain.ChainEffectType {
	return convertChainEffectType(s.EffectType)
}

// ToDomainCategory はChainEffectDataからドメインモデルのChainEffectCategoryに変換します。
func (s *ChainEffectData) ToDomainCategory() domain.ChainEffectCategory {
	return convertChainEffectCategory(s.Category)
}

// ToDomainChainEffect はChainEffectDataをドメインモデルのChainEffectに変換します。
// 効果値はMinValueを使用します（表示用のデフォルト値として）。
func (s *ChainEffectData) ToDomainChainEffect() domain.ChainEffect {
	return domain.NewChainEffectWithTemplate(
		s.ID,
		convertChainEffectType(s.EffectType),
		s.MinValue,
		s.Description,
		s.ShortDescription,
	)
}

// ToChainEffectDefinition はChainEffectDataをChainEffectDefinitionに変換します。
// rewardingパッケージのChainEffectPoolで使用されます。
func (s *ChainEffectData) ToChainEffectDefinition() ChainEffectDefinitionData {
	return ChainEffectDefinitionData{
		ID:               s.ID,
		Name:             s.Name,
		Description:      s.Description,
		ShortDescription: s.ShortDescription,
		Category:         s.Category,
		EffectType:       convertChainEffectType(s.EffectType),
		MinValue:         s.MinValue,
		MaxValue:         s.MaxValue,
		MinDropLevel:     s.MinDropLevel,
	}
}

// ChainEffectDefinitionData はチェイン効果定義のデータ構造体です。
// rewardingパッケージのChainEffectDefinitionと同等の構造を持ちます。
type ChainEffectDefinitionData struct {
	// ID はチェイン効果の一意識別子です。
	ID string

	// Name は表示名です。
	Name string

	// Description は説明文テンプレートです（%.0fをプレースホルダとして使用）。
	Description string

	// ShortDescription は短い説明文テンプレートです（%.0fをプレースホルダとして使用）。
	ShortDescription string

	// Category はカテゴリです（attack, defense, heal等）。
	Category string

	// EffectType はドメインモデルのChainEffectTypeです。
	EffectType domain.ChainEffectType

	// MinValue は効果値の最小値です。
	MinValue float64

	// MaxValue は効果値の最大値です。
	MaxValue float64

	// MinDropLevel はこのチェイン効果がドロップする最低敵レベルです。
	MinDropLevel int
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

// convertChainEffectCategory は文字列をChainEffectCategoryに変換します。
func convertChainEffectCategory(s string) domain.ChainEffectCategory {
	switch s {
	case "attack":
		return domain.ChainEffectCategoryAttack
	case "defense":
		return domain.ChainEffectCategoryDefense
	case "heal":
		return domain.ChainEffectCategoryHeal
	case "typing":
		return domain.ChainEffectCategoryTyping
	case "recast":
		return domain.ChainEffectCategoryRecast
	case "effect_extend":
		return domain.ChainEffectCategoryEffectExtend
	case "special":
		return domain.ChainEffectCategorySpecial
	default:
		return domain.ChainEffectCategorySpecial
	}
}

// ==================== タイピング辞書 ====================

// TypingDictionary はwords.jsonから読み込むタイピング辞書データの構造体です。

type TypingDictionary struct {
	Easy   []string `json:"easy"`
	Medium []string `json:"medium"`
	Hard   []string `json:"hard"`
}

// wordsFileData はwords.jsonのルート構造です。
type wordsFileData struct {
	Words TypingDictionary `json:"words"`
}

// LoadTypingDictionary はwords.jsonからタイピング辞書を読み込みます。

func (l *DataLoader) LoadTypingDictionary() (*TypingDictionary, error) {
	data, err := l.readFile("words.json")
	if err != nil {
		return nil, fmt.Errorf("words.jsonの読み込みに失敗: %w", err)
	}

	var fileData wordsFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("words.jsonのパースに失敗: %w", err)
	}

	return &fileData.Words, nil
}

// ==================== 初期エージェント定義 ====================

// FirstAgentData はfirst_agent.jsonから読み込む初期エージェントデータの構造体です。
type FirstAgentData struct {
	ID         string                 `json:"id"`
	CoreTypeID string                 `json:"core_type_id"`
	CoreLevel  int                    `json:"core_level"`
	Modules    []FirstAgentModuleData `json:"modules"`
}

// FirstAgentModuleData は初期エージェントのモジュールデータの構造体です。
// ChainEffectはモジュールと密結合しているため、同じ構造体で定義します。
type FirstAgentModuleData struct {
	TypeID           string  `json:"type_id"`
	ChainEffectType  string  `json:"chain_effect_type,omitempty"`
	ChainEffectValue float64 `json:"chain_effect_value,omitempty"`
}

// HasChainEffect はチェイン効果を持つかどうかを返します。
func (m *FirstAgentModuleData) HasChainEffect() bool {
	return m.ChainEffectType != ""
}

// firstAgentFileData はfirst_agent.jsonのルート構造です。
type firstAgentFileData struct {
	FirstAgents []FirstAgentData `json:"first_agents"`
}

// LoadFirstAgents はfirst_agent.jsonから初期エージェント定義を読み込みます。
func (l *DataLoader) LoadFirstAgents() ([]FirstAgentData, error) {
	data, err := l.readFile("first_agent.json")
	if err != nil {
		return nil, fmt.Errorf("first_agent.jsonの読み込みに失敗: %w", err)
	}

	var fileData firstAgentFileData
	if err := json.Unmarshal(data, &fileData); err != nil {
		return nil, fmt.Errorf("first_agent.jsonのパースに失敗: %w", err)
	}

	return fileData.FirstAgents, nil
}

// ==================== 全データ一括ロード ====================

// LoadAllExternalData は全ての外部データファイルを一括でロードします。

func (l *DataLoader) LoadAllExternalData() (*ExternalData, error) {
	coreTypes, err := l.LoadCoreTypes()
	if err != nil {
		return nil, fmt.Errorf("コア特性のロードに失敗: %w", err)
	}

	modules, err := l.LoadModuleDefinitions()
	if err != nil {
		return nil, fmt.Errorf("モジュール定義のロードに失敗: %w", err)
	}

	enemyTypes, err := l.LoadEnemyTypes()
	if err != nil {
		return nil, fmt.Errorf("敵タイプのロードに失敗: %w", err)
	}

	// 敵行動データのロード（オプショナル：ファイルが存在しない場合は空配列）
	enemyActions, err := l.LoadEnemyActions()
	if err != nil {
		// enemy_actions.jsonが存在しない場合は空配列を使用（後方互換性）
		enemyActions = []EnemyActionData{}
	}

	// 敵パッシブスキルデータのロード（オプショナル：ファイルが存在しない場合は空配列）
	enemyPassiveSkills, err := l.LoadEnemyPassiveSkills()
	if err != nil {
		// enemy_passive_skills.jsonが存在しない場合は空配列を使用（後方互換性）
		enemyPassiveSkills = []EnemyPassiveSkillData{}
	}

	passiveSkills, err := l.LoadPassiveSkills()
	if err != nil {
		return nil, fmt.Errorf("パッシブスキルのロードに失敗: %w", err)
	}

	chainEffects, err := l.LoadChainEffects()
	if err != nil {
		return nil, fmt.Errorf("チェイン効果のロードに失敗: %w", err)
	}

	dictionary, err := l.LoadTypingDictionary()
	if err != nil {
		return nil, fmt.Errorf("タイピング辞書のロードに失敗: %w", err)
	}

	firstAgents, err := l.LoadFirstAgents()
	if err != nil {
		return nil, fmt.Errorf("初期エージェントのロードに失敗: %w", err)
	}

	return &ExternalData{
		CoreTypes:          coreTypes,
		ModuleDefinitions:  modules,
		EnemyTypes:         enemyTypes,
		EnemyActions:       enemyActions,
		EnemyPassiveSkills: enemyPassiveSkills,
		PassiveSkills:      passiveSkills,
		ChainEffects:       chainEffects,
		TypingDictionary:   dictionary,
		FirstAgents:        firstAgents,
	}, nil
}

// ==================== バリデーション ====================

// ValidateCoreTypeData はコア特性データのバリデーションを行います。
func ValidateCoreTypeData(data CoreTypeData) error {
	if data.ID == "" {
		return fmt.Errorf("コア特性IDが空です")
	}
	if data.Name == "" {
		return fmt.Errorf("コア特性名が空です: ID=%s", data.ID)
	}
	if len(data.AllowedTags) == 0 {
		return fmt.Errorf("許可タグが空です: ID=%s", data.ID)
	}
	if len(data.StatWeights) == 0 {
		return fmt.Errorf("ステータス重みが空です: ID=%s", data.ID)
	}
	return nil
}

// ValidateSkillDefinitionData はスキル定義データのバリデーションを行います。
func ValidateSkillDefinitionData(data SkillDefinitionData) error {
	if data.ID == "" {
		return fmt.Errorf("スキルIDが空です")
	}
	if data.Name == "" {
		return fmt.Errorf("スキル名が空です: ID=%s", data.ID)
	}
	if len(data.Effects) == 0 {
		return fmt.Errorf("スキル効果が空です: ID=%s", data.ID)
	}
	return nil
}

// ValidateModuleDefinitionData はValidateSkillDefinitionDataの後方互換ラッパーです。
// 後方互換性のために残されています。新規コードではValidateSkillDefinitionDataを使用してください。
func ValidateModuleDefinitionData(data ModuleDefinitionData) error {
	return ValidateSkillDefinitionData(data)
}

// ValidateEnemyTypeData は敵タイプデータのバリデーションを行います。
func ValidateEnemyTypeData(data EnemyTypeData) error {
	if data.ID == "" {
		return fmt.Errorf("敵タイプIDが空です")
	}
	if data.Name == "" {
		return fmt.Errorf("敵タイプ名が空です: ID=%s", data.ID)
	}
	if data.BaseHP <= 0 {
		return fmt.Errorf("敵の基礎HPが不正です: ID=%s, BaseHP=%d", data.ID, data.BaseHP)
	}
	if data.BaseAttackPower <= 0 {
		return fmt.Errorf("敵の基礎攻撃力が不正です: ID=%s, BaseAttackPower=%d", data.ID, data.BaseAttackPower)
	}
	return nil
}
