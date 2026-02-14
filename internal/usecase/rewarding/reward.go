// Package reward はドロップ・報酬システムを提供します。
// バトル勝利時の報酬計算、コア/スキルのドロップ判定を担当します。
package rewarding

import (
	"log/slog"
	"math/rand"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// ==================== チェイン効果プール ====================

// デフォルトのチェイン効果なし確率（30%）
const DefaultNoEffectProbability = 0.3

// ChainEffectDefinition はチェイン効果定義の構造体です。
// マスタデータからロードされたチェイン効果情報を保持します。
type ChainEffectDefinition struct {
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

// ChainEffectPool はチェイン効果のプール管理を担当する構造体です。
// スキル入手時にランダムなチェイン効果を生成します。
type ChainEffectPool struct {
	// Effects は利用可能なチェイン効果定義のリストです。
	Effects []ChainEffectDefinition

	// rng は乱数生成器です。
	rng *rand.Rand

	// noEffectProbability はチェイン効果なしの確率です。
	noEffectProbability float64
}

// NewChainEffectPool は新しいChainEffectPoolを作成します。
func NewChainEffectPool(effects []ChainEffectDefinition) *ChainEffectPool {
	return &ChainEffectPool{
		Effects:             effects,
		rng:                 rand.New(rand.NewSource(time.Now().UnixNano())),
		noEffectProbability: DefaultNoEffectProbability,
	}
}

// SetNoEffectProbability はチェイン効果なしの確率を設定します。
// 値は0.0から1.0の範囲で指定します。
func (p *ChainEffectPool) SetNoEffectProbability(probability float64) {
	if probability < 0.0 {
		probability = 0.0
	} else if probability > 1.0 {
		probability = 1.0
	}
	p.noEffectProbability = probability
}

// GenerateRandomEffect はランダムなチェイン効果を生成します。
// noEffectProbabilityの確率でnilを返します。
func (p *ChainEffectPool) GenerateRandomEffect() *domain.ChainEffect {
	// チェイン効果がない場合
	if len(p.Effects) == 0 {
		return nil
	}

	// チェイン効果なしの確率判定
	if p.rng.Float64() < p.noEffectProbability {
		return nil
	}

	// ランダムにチェイン効果を選択
	selected := p.Effects[p.rng.Intn(len(p.Effects))]

	// 効果値をmin-max範囲内でランダムに決定
	value := selected.MinValue
	if selected.MaxValue > selected.MinValue {
		value = selected.MinValue + p.rng.Float64()*(selected.MaxValue-selected.MinValue)
		// 整数に丸める（効果値は通常整数で表現される）
		value = float64(int(value + 0.5))
	}

	effect := domain.NewChainEffectWithTemplate(
		selected.ID,
		selected.EffectType,
		value,
		selected.Description,
		selected.ShortDescription,
	)
	return &effect
}

// BattleStatistics はバトル統計を表す構造体です。
type BattleStatistics struct {
	// TotalWPM はWPMの合計値です。
	TotalWPM float64

	// TotalAccuracy は正確性の合計値です。
	TotalAccuracy float64

	// ClearTime はクリア時間です。
	ClearTime time.Duration

	// TotalTypingCount は総タイピング回数です。
	TotalTypingCount int

	// TotalDamageDealt は与えた総ダメージです。
	TotalDamageDealt int

	// TotalDamageTaken は受けた総ダメージです。
	TotalDamageTaken int

	// TotalHealAmount は総回復量です。
	TotalHealAmount int
}

// GetAverageWPM は平均WPMを返します。
func (s *BattleStatistics) GetAverageWPM() float64 {
	if s.TotalTypingCount == 0 {
		return 0
	}
	return s.TotalWPM / float64(s.TotalTypingCount)
}

// GetAverageAccuracy は平均正確性を返します。
func (s *BattleStatistics) GetAverageAccuracy() float64 {
	if s.TotalTypingCount == 0 {
		return 0
	}
	return s.TotalAccuracy / float64(s.TotalTypingCount)
}

// RewardResult は報酬計算結果を表す構造体です。
type RewardResult struct {
	// IsVictory は勝利かどうかです。
	IsVictory bool

	// ShowRewardScreen は報酬画面を表示すべきかどうかです。
	ShowRewardScreen bool

	// Stats はバトル統計です。
	Stats *BattleStatistics

	// DroppedCores はドロップしたコアのリストです。
	DroppedCores []*domain.CoreModel

	// DroppedSkills はドロップしたスキルのリストです。
	DroppedSkills []*domain.SkillModel

	// EnemyLevel は撃破した敵のレベルです。
	EnemyLevel int

	// HPGain は獲得最大HP（敵撃破によるHP成長）です。
	HPGain int

	// RankUnlocked は新ランク解放フラグです。
	RankUnlocked bool

	// PreviousMaxHP はHP増加前の最大HP値です（表示用）。
	PreviousMaxHP int

	// NewMaxHP はHP増加後の最大HP値です（表示用）。
	NewMaxHP int

	// PreviousRank はランク増加前のランク値です（表示用）。
	PreviousRank int

	// NewRank はランク増加後のランク値です（表示用）。
	NewRank int

	// RankUpRewardCores はランクアップ報酬のコアリストです。
	RankUpRewardCores []*domain.CoreModel

	// RankUpRewardSkills はランクアップ報酬のスキルリストです。
	RankUpRewardSkills []*domain.SkillModel

	// RankUpRewardChainEffects はランクアップ報酬のチェイン効果リストです。
	RankUpRewardChainEffects []domain.ChainEffect
}

// InventoryWarning はインベントリ警告を表す構造体です。
// 注: 新しいインベントリシステムでは容量制限がないため、この構造体は後方互換性のために残されています。
type InventoryWarning struct {
	// WarningMessage は警告メッセージです。
	WarningMessage string

	// SuggestDiscard は破棄を促すかどうかです（現在は未使用）。
	SuggestDiscard bool
}

// TempStorage は一時保管を表す構造体です。
type TempStorage struct {
	// Cores は一時保管中のコアリストです。
	Cores []*domain.CoreModel

	// Skills は一時保管中のスキルリストです。
	Skills []*domain.SkillModel
}

// AddCore はコアを一時保管に追加します。
func (s *TempStorage) AddCore(core *domain.CoreModel) {
	s.Cores = append(s.Cores, core)
}

// AddSkill はスキルを一時保管に追加します。
func (s *TempStorage) AddSkill(skill *domain.SkillModel) {
	s.Skills = append(s.Skills, skill)
}

// RetrieveCores は一時保管中のコアを全て取り出します。
func (s *TempStorage) RetrieveCores() []*domain.CoreModel {
	cores := s.Cores
	s.Cores = nil
	return cores
}

// RetrieveSkills は一時保管中のスキルを全て取り出します。
func (s *TempStorage) RetrieveSkills() []*domain.SkillModel {
	skills := s.Skills
	s.Skills = nil
	return skills
}

// HasItems は一時保管にアイテムがあるかどうかを返します。
func (s *TempStorage) HasItems() bool {
	return len(s.Cores) > 0 || len(s.Skills) > 0
}

// SkillDropInfo はスキルドロップに必要な情報を持つ構造体です。
type SkillDropInfo struct {
	// ID はスキルの一意識別子です。
	ID string

	// Name はスキルの表示名です。
	Name string

	// Icon はスキルのアイコン（絵文字）です。
	Icon string

	// Tags はスキルのタグリストです。
	Tags []string

	// Description はスキルの効果説明です。
	Description string

	// CooldownSeconds はスキルのクールダウン時間（秒）です。
	CooldownSeconds float64

	// DifficultyRate はタイピングの難易度です（50-200、100=標準）。
	DifficultyRate int

	// ChallengeType はチャレンジタイプのIDです。
	ChallengeType domain.ChallengeTypeID

	// ManaCost はスキル使用時に消費するマナ量です。
	ManaCost int

	// MinDropLevel はこのスキルがドロップする最低敵レベルです。
	MinDropLevel int

	// Effects はスキルの効果リストです。
	Effects []domain.SkillEffect
}

// ToSkillType はSkillDropInfoをドメインモデルのSkillTypeに変換します。
func (m *SkillDropInfo) ToSkillType() domain.SkillType {
	// Tagsをコピー（スライスの参照共有を避ける）
	tagsCopy := make([]string, len(m.Tags))
	copy(tagsCopy, m.Tags)

	// Effectsをコピー
	effectsCopy := make([]domain.SkillEffect, len(m.Effects))
	copy(effectsCopy, m.Effects)

	return domain.SkillType{
		ID:              m.ID,
		Name:            m.Name,
		Icon:            m.Icon,
		Tags:            tagsCopy,
		Description:     m.Description,
		CooldownSeconds: m.CooldownSeconds,
		DifficultyRate:  m.DifficultyRate,
		ChallengeType:   m.ChallengeType,
		ManaCost:        m.ManaCost,
		MinDropLevel:    m.MinDropLevel,
		Effects:         effectsCopy,
	}
}

// ToDomain はSkillDropInfoをドメインモデルのSkillModelに変換します。
func (m *SkillDropInfo) ToDomain() *domain.SkillModel {
	skillType := m.ToSkillType()
	return domain.NewSkillFromType(skillType, nil)
}

// ToDomainWithChainEffect はチェイン効果付きでドメインモデルに変換します。
func (m *SkillDropInfo) ToDomainWithChainEffect(chainEffect *domain.ChainEffect) *domain.SkillModel {
	skillType := m.ToSkillType()
	return domain.NewSkillFromType(skillType, chainEffect)
}

// RewardCalculator はドメイン型を使用した報酬計算を担当する構造体です。
type RewardCalculator struct {
	// coreTypes はコア特性定義リストです（ドメイン型）。
	coreTypes []domain.CoreType

	// skillTypes はスキル定義リストです。
	skillTypes []SkillDropInfo

	// passiveSkills はパッシブスキル定義マップです。
	passiveSkills map[string]domain.PassiveSkill

	// chainEffectPool はチェイン効果プールです。
	chainEffectPool *ChainEffectPool

	// rankRewards はランクアップ報酬マップです（ランク番号→報酬）。
	rankRewards map[int]domain.RankReward

	// rng は乱数生成器です。
	rng *rand.Rand
}

// NewRewardCalculator はドメイン型を使用する新しいRewardCalculatorを作成します。
func NewRewardCalculator(
	coreTypes []domain.CoreType,
	skillTypes []SkillDropInfo,
	passiveSkills map[string]domain.PassiveSkill,
) *RewardCalculator {
	return &RewardCalculator{
		coreTypes:     coreTypes,
		skillTypes:    skillTypes,
		passiveSkills: passiveSkills,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// SetChainEffectPool はチェイン効果プールを設定します。
func (c *RewardCalculator) SetChainEffectPool(pool *ChainEffectPool) {
	c.chainEffectPool = pool
}

// SetRankRewards はランクアップ報酬マップを設定します。
func (c *RewardCalculator) SetRankRewards(rewards map[int]domain.RankReward) {
	c.rankRewards = rewards
}

// GetChainEffectPool はチェイン効果プールを取得します。
func (c *RewardCalculator) GetChainEffectPool() *ChainEffectPool {
	return c.chainEffectPool
}

// CreateRewardResult は報酬結果を作成します。
func (c *RewardCalculator) CreateRewardResult(isVictory bool, stats *BattleStatistics, enemyLevel int) *RewardResult {
	result := &RewardResult{
		IsVictory:     isVictory,
		Stats:         stats,
		EnemyLevel:    enemyLevel,
		DroppedCores:  make([]*domain.CoreModel, 0),
		DroppedSkills: make([]*domain.SkillModel, 0),
	}

	if !isVictory {
		result.ShowRewardScreen = false
		return result
	}

	result.ShowRewardScreen = true
	return result
}

// GetEligibleCoreTypes は指定レベルでドロップ可能なコア特性を返します。
func (c *RewardCalculator) GetEligibleCoreTypes(enemyLevel int) []domain.CoreType {
	eligible := make([]domain.CoreType, 0)
	for _, coreType := range c.coreTypes {
		if coreType.MinDropLevel <= enemyLevel {
			eligible = append(eligible, coreType)
		}
	}
	return eligible
}

// GetEligibleSkillTypes は指定レベルでドロップ可能なスキルを返します。
func (c *RewardCalculator) GetEligibleSkillTypes(enemyLevel int) []SkillDropInfo {
	eligible := make([]SkillDropInfo, 0)
	for _, skillType := range c.skillTypes {
		if skillType.MinDropLevel <= enemyLevel {
			eligible = append(eligible, skillType)
		}
	}
	return eligible
}

// CheckInventoryFull はインベントリの満杯状態をチェックします。
// 注: 新しいインベントリシステムでは容量制限がないため、常に空の警告を返します。
func (c *RewardCalculator) CheckInventoryFull(
	coreInv *domain.CoreInventory,
	skillInv *domain.SkillInventory,
) *InventoryWarning {
	// 新しいシステムでは容量制限がないため、警告は不要
	return &InventoryWarning{}
}

// CreateTempStorage は一時保管を作成します。
func (c *RewardCalculator) CreateTempStorage() *TempStorage {
	return &TempStorage{
		Cores:  make([]*domain.CoreModel, 0),
		Skills: make([]*domain.SkillModel, 0),
	}
}

// CalculateGuaranteedReward は確定ドロップを計算します。
// 敵タイプのDropItemCategoryとDropItemTypeIDに基づいてアイテムを決定します。
// ドロップ設定がない場合はpanicします。
func (c *RewardCalculator) CalculateGuaranteedReward(
	stats *BattleStatistics,
	enemyLevel int,
	enemyType domain.EnemyType,
) *RewardResult {
	// ドロップ設定の確認
	if enemyType.DropItemCategory == "" || enemyType.DropItemTypeID == "" {
		panic("敵タイプ " + enemyType.ID + " にドロップ設定がありません")
	}

	result := &RewardResult{
		IsVictory:        true,
		ShowRewardScreen: true,
		Stats:            stats,
		EnemyLevel:       enemyLevel,
		DroppedCores:     make([]*domain.CoreModel, 0),
		DroppedSkills:    make([]*domain.SkillModel, 0),
	}

	// 確定ドロップ処理
	switch enemyType.DropItemCategory {
	case "core":
		core := c.RollCoreDropWithTypeID(enemyType.DropItemTypeID, enemyLevel)
		if core == nil {
			panic("敵タイプ " + enemyType.ID + " のコアTypeID " + enemyType.DropItemTypeID + " が見つかりません")
		}
		result.DroppedCores = append(result.DroppedCores, core)

	case "skill":
		skill := c.RollSkillDropWithTypeID(enemyType.DropItemTypeID, enemyLevel)
		if skill == nil {
			panic("敵タイプ " + enemyType.ID + " のスキルTypeID " + enemyType.DropItemTypeID + " が見つかりません")
		}
		result.DroppedSkills = append(result.DroppedSkills, skill)

	default:
		panic("敵タイプ " + enemyType.ID + " のドロップカテゴリ " + enemyType.DropItemCategory + " が不正です")
	}

	return result
}

// RollCoreDropWithTypeID は指定されたTypeIDのコアを生成します。
// コアレベルは敵レベルと同じになります。
func (c *RewardCalculator) RollCoreDropWithTypeID(typeID string, enemyLevel int) *domain.CoreModel {
	// 指定されたTypeIDのコア特性を検索
	var selectedType *domain.CoreType
	for i := range c.coreTypes {
		if c.coreTypes[i].ID == typeID {
			selectedType = &c.coreTypes[i]
			break
		}
	}

	if selectedType == nil {
		return nil
	}

	// パッシブスキルを取得
	passiveSkill := domain.PassiveSkill{}
	if c.passiveSkills != nil {
		if skill, ok := c.passiveSkills[selectedType.PassiveSkillID]; ok {
			passiveSkill = skill
		}
	}

	// コアをインスタンス化（TypeIDベース、レベルなし）
	return domain.NewCoreWithTypeID(
		selectedType.ID,
		*selectedType,
		passiveSkill,
	)
}

// RollSkillDropWithTypeID は指定されたTypeIDのスキルを生成します。
// 敵レベルに応じたチェイン効果をランダムに選択します。
func (c *RewardCalculator) RollSkillDropWithTypeID(typeID string, enemyLevel int) *domain.SkillModel {
	// 指定されたTypeIDのスキルを検索
	var selectedType *SkillDropInfo
	for i := range c.skillTypes {
		if c.skillTypes[i].ID == typeID {
			selectedType = &c.skillTypes[i]
			break
		}
	}

	if selectedType == nil {
		return nil
	}

	// チェイン効果を生成（プールが設定されている場合）
	var chainEffect *domain.ChainEffect
	if c.chainEffectPool != nil {
		chainEffect = c.generateLevelBasedChainEffect(enemyLevel)
	}

	// スキルをインスタンス化
	return selectedType.ToDomainWithChainEffect(chainEffect)
}

// generateLevelBasedChainEffect は敵レベルに応じたチェイン効果を生成します。
// 敵レベル以下のMinDropLevelを持つ効果からランダムに選択します。
// 必ずチェイン効果を返します（該当する効果がない場合はpanic）。
func (c *RewardCalculator) generateLevelBasedChainEffect(enemyLevel int) *domain.ChainEffect {
	if c.chainEffectPool == nil || len(c.chainEffectPool.Effects) == 0 {
		panic("チェイン効果プールが設定されていません")
	}

	// 敵レベル以下のチェイン効果をフィルタリング
	eligibleEffects := make([]ChainEffectDefinition, 0)
	for _, effect := range c.chainEffectPool.Effects {
		if effect.MinDropLevel <= enemyLevel {
			eligibleEffects = append(eligibleEffects, effect)
		}
	}

	if len(eligibleEffects) == 0 {
		panic("敵レベル に対応するチェイン効果がありません")
	}

	// ランダムにチェイン効果を選択
	selected := eligibleEffects[c.rng.Intn(len(eligibleEffects))]

	// 効果値を決定（高レベルほど高い値が出やすい）
	levelFactor := float64(enemyLevel) / 100.0
	valueRange := selected.MaxValue - selected.MinValue
	levelBonus := valueRange * levelFactor * 0.3 // 最大30%のボーナス
	baseValue := selected.MinValue + c.rng.Float64()*(valueRange-levelBonus) + levelBonus
	value := float64(int(baseValue + 0.5))

	if value > selected.MaxValue {
		value = selected.MaxValue
	}

	effect := domain.NewChainEffectWithTemplate(
		selected.ID,
		selected.EffectType,
		value,
		selected.Description,
		selected.ShortDescription,
	)
	return &effect
}

// CalculateGuaranteedRewardWithProgress はHP成長報酬付きで確定ドロップを計算します。
// 敵タイプのDropItemCategoryとDropItemTypeIDに基づいてアイテムを決定し、
// EnemyProgressに撃破を記録し、PlayerModelにHP成長を適用します。
func (c *RewardCalculator) CalculateGuaranteedRewardWithProgress(
	stats *BattleStatistics,
	enemyLevel int,
	enemyType domain.EnemyType,
	progress *domain.EnemyProgress,
	player *domain.PlayerModel,
	enemyTypes map[string]domain.EnemyType,
) *RewardResult {
	// 基本の報酬計算
	result := c.CalculateGuaranteedReward(stats, enemyLevel, enemyType)

	// 変更前の値を記録
	previousMaxHP := player.MaxHP
	previousRank := progress.CurrentRank

	// HP成長を計算
	hpGain, _ := progress.RecordDefeat(enemyType.ID, enemyLevel)

	// HP成長をPlayerModelに適用
	if hpGain > 0 {
		player.IncreaseMaxHP(hpGain)
	}

	// ランク解放チェック
	rankUnlocked := false
	if checkRankComplete(progress, enemyTypes) {
		progress.CurrentRank++
		rankUnlocked = true
	}

	result.HPGain = hpGain
	result.RankUnlocked = rankUnlocked
	result.PreviousMaxHP = previousMaxHP
	result.NewMaxHP = player.MaxHP
	result.PreviousRank = previousRank
	result.NewRank = progress.CurrentRank

	// ランクアップ報酬の解決
	if rankUnlocked && c.rankRewards != nil {
		c.resolveRankUpRewards(result, progress.CurrentRank)
	}

	return result
}

// checkRankComplete はランク内全敵撃破チェックを行います。
func checkRankComplete(progress *domain.EnemyProgress, enemyTypes map[string]domain.EnemyType) bool {
	currentRank := progress.CurrentRank

	hasEnemiesInRank := false
	for _, et := range enemyTypes {
		if et.Rank == currentRank {
			hasEnemiesInRank = true
			if !progress.IsDefeated(et.ID) {
				return false
			}
		}
	}

	return hasEnemiesInRank
}

// resolveRankUpRewards はランクアップ報酬をRewardResultに格納します。
func (c *RewardCalculator) resolveRankUpRewards(result *RewardResult, newRank int) {
	reward, ok := c.rankRewards[newRank]
	if !ok {
		return
	}

	for _, item := range reward.Items {
		switch item.Category {
		case "core":
			core := c.RollCoreDropWithTypeID(item.TypeID, 1)
			if core != nil {
				result.RankUpRewardCores = append(result.RankUpRewardCores, core)
			}
		case "skill":
			// ランクアップ報酬のスキルにはチェイン効果を付与しない
			var selectedType *SkillDropInfo
			for i := range c.skillTypes {
				if c.skillTypes[i].ID == item.TypeID {
					selectedType = &c.skillTypes[i]
					break
				}
			}
			if selectedType != nil {
				result.RankUpRewardSkills = append(result.RankUpRewardSkills, selectedType.ToDomain())
			}
		case "chain_effect":
			if c.chainEffectPool != nil {
				for _, def := range c.chainEffectPool.Effects {
					if def.ID == item.TypeID {
						effect := domain.NewChainEffectWithTemplate(
							def.ID,
							def.EffectType,
							def.MinValue,
							def.Description,
							def.ShortDescription,
						)
						result.RankUpRewardChainEffects = append(result.RankUpRewardChainEffects, effect)
						break
					}
				}
			}
		}
	}
}

// AddRewardsToInventory はドロップしたアイテムをインベントリに追加します。
// 新しいインベントリシステムではTypeIDベースでユニーク管理されるため、
// 同一TypeIDのより高いレベルのみが保存されます（コア）。
// スキルは保有状態とチェイン効果バリエーションが追加されます。
func AddRewardsToInventory(
	result *RewardResult,
	coreInv *domain.CoreInventory,
	skillInv *domain.SkillInventory,
	chainEffectInv *domain.ChainEffectInventory,
) *InventoryWarning {
	// コアをインベントリに追加（TypeIDのみ）
	for _, core := range result.DroppedCores {
		updated := coreInv.AddCore(core.Type.ID)
		if updated {
			slog.Info("コアをインベントリに追加",
				slog.String("core_type_id", core.Type.ID),
			)
		}
	}

	// スキル（スキル）をインベントリに追加
	for _, skill := range result.DroppedSkills {
		skillInv.AddSkill(skill.TypeID)
		slog.Info("スキルをインベントリに追加",
			slog.String("skill_type_id", skill.TypeID),
		)

		// チェイン効果を独立インベントリに追加
		if chainEffectInv != nil && skill.HasChainEffect() {
			added := chainEffectInv.AddChainEffect(skill.ChainEffect.ID)
			if added {
				slog.Info("チェイン効果をインベントリに追加",
					slog.String("chain_effect_id", skill.ChainEffect.ID),
				)
			}
		}
	}

	// ランクアップ報酬のコアをインベントリに追加
	for _, core := range result.RankUpRewardCores {
		updated := coreInv.AddCore(core.Type.ID)
		if updated {
			slog.Info("ランクアップ報酬コアをインベントリに追加",
				slog.String("core_type_id", core.Type.ID),
			)
		}
	}

	// ランクアップ報酬のスキルをインベントリに追加
	for _, skill := range result.RankUpRewardSkills {
		skillInv.AddSkill(skill.TypeID)
		slog.Info("ランクアップ報酬スキルをインベントリに追加",
			slog.String("skill_type_id", skill.TypeID),
		)
	}

	// ランクアップ報酬のチェイン効果をインベントリに追加
	if chainEffectInv != nil {
		for _, effect := range result.RankUpRewardChainEffects {
			added := chainEffectInv.AddChainEffect(effect.ID)
			if added {
				slog.Info("ランクアップ報酬チェイン効果をインベントリに追加",
					slog.String("chain_effect_id", effect.ID),
				)
			}
		}
	}

	// 新しいシステムでは容量制限がないため、警告は不要
	return &InventoryWarning{}
}
