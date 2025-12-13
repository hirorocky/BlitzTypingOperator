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
	CoreTypes         []CoreTypeData
	ModuleDefinitions []ModuleDefinitionData
	EnemyTypes        []EnemyTypeData
	TypingDictionary  *TypingDictionary
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

// ==================== モジュール定義 ====================

// ModuleDefinitionData はmodules.jsonから読み込むモジュール定義データの構造体です。
type ModuleDefinitionData struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Level           int      `json:"level"`
	Tags            []string `json:"tags"`
	BaseEffect      float64  `json:"base_effect"`
	StatReference   string   `json:"stat_reference"`
	Description     string   `json:"description"`
	CooldownSeconds float64  `json:"cooldown_seconds"`
	Difficulty      int      `json:"difficulty"`
	MinDropLevel    int      `json:"min_drop_level"`
}

// modulesFileData はmodules.jsonのルート構造です。
type modulesFileData struct {
	Modules []ModuleDefinitionData `json:"modules"`
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

	return fileData.Modules, nil
}

// ToDomain はModuleDefinitionDataをドメインモデルのModuleModelに変換します。
func (m *ModuleDefinitionData) ToDomain() *domain.ModuleModel {
	// カテゴリ文字列をModuleCategoryに変換
	var category domain.ModuleCategory
	switch m.Category {
	case "physical_attack":
		category = domain.PhysicalAttack
	case "magic_attack":
		category = domain.MagicAttack
	case "heal":
		category = domain.Heal
	case "buff":
		category = domain.Buff
	case "debuff":
		category = domain.Debuff
	default:
		category = domain.PhysicalAttack // デフォルト
	}

	return domain.NewModule(
		m.ID,
		m.Name,
		category,
		m.Level,
		m.Tags,
		m.BaseEffect,
		m.StatReference,
		m.Description,
	)
}

// ==================== 敵タイプ定義 ====================

// EnemyTypeData はenemies.jsonから読み込む敵タイプデータの構造体です。
type EnemyTypeData struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	BaseHP               int    `json:"base_hp"`
	BaseAttackPower      int    `json:"base_attack_power"`
	BaseAttackIntervalMS int64  `json:"base_attack_interval_ms"`
	AttackType           string `json:"attack_type"`
	ASCIIArt             string `json:"ascii_art"`
	// 内部で計算されるフィールド
	BaseAttackInterval time.Duration `json:"-"`
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

	// ミリ秒をtime.Durationに変換
	for i := range fileData.EnemyTypes {
		fileData.EnemyTypes[i].BaseAttackInterval = time.Duration(fileData.EnemyTypes[i].BaseAttackIntervalMS) * time.Millisecond
	}

	return fileData.EnemyTypes, nil
}

// ToDomain はEnemyTypeDataをドメインモデルのEnemyTypeに変換します。
func (e *EnemyTypeData) ToDomain() domain.EnemyType {
	return domain.EnemyType{
		ID:                 e.ID,
		Name:               e.Name,
		BaseHP:             e.BaseHP,
		BaseAttackPower:    e.BaseAttackPower,
		BaseAttackInterval: e.BaseAttackInterval,
		AttackType:         e.AttackType,
		ASCIIArt:           e.ASCIIArt,
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

	dictionary, err := l.LoadTypingDictionary()
	if err != nil {
		return nil, fmt.Errorf("タイピング辞書のロードに失敗: %w", err)
	}

	return &ExternalData{
		CoreTypes:         coreTypes,
		ModuleDefinitions: modules,
		EnemyTypes:        enemyTypes,
		TypingDictionary:  dictionary,
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

// ValidateModuleDefinitionData はモジュール定義データのバリデーションを行います。
func ValidateModuleDefinitionData(data ModuleDefinitionData) error {
	if data.ID == "" {
		return fmt.Errorf("モジュールIDが空です")
	}
	if data.Name == "" {
		return fmt.Errorf("モジュール名が空です: ID=%s", data.ID)
	}
	if data.Category == "" {
		return fmt.Errorf("モジュールカテゴリが空です: ID=%s", data.ID)
	}
	if data.Level < 1 {
		return fmt.Errorf("モジュールレベルが不正です: ID=%s, Level=%d", data.ID, data.Level)
	}
	return nil
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
