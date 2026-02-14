// Package savedata はセーブデータの永続化を担当します。
// 原子的書き込みパターンとバックアップによるセーブデータの整合性を保証します。

package savedata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CurrentSaveDataVersion は現在のセーブデータバージョンです。
// まだ正式版をリリースしていないので、リリースするまではv0.0.1で固定する。
const CurrentSaveDataVersion = "0.0.1"

// SaveFileName はセーブファイル名です。
const SaveFileName = "save.json"

// DebugSaveFileName はデバッグモード用セーブファイル名です。
const DebugSaveFileName = "debug_save.json"

// TempSaveFileName は一時セーブファイル名です（原子的書き込み用）。
const TempSaveFileName = "save.json.tmp"

// DebugTempSaveFileName はデバッグモード用一時セーブファイル名です。
const DebugTempSaveFileName = "debug_save.json.tmp"

// MaxBackupCount はバックアップの最大世代数です。

const MaxBackupCount = 3

// SaveData はゲームのセーブデータを表す構造体です。

type SaveData struct {
	// Version はセーブデータのバージョンです。
	// スキーマ変更時のマイグレーションに使用します。
	Version string `json:"version"`

	// Timestamp はセーブした日時です。
	Timestamp time.Time `json:"timestamp"`

	// Player はプレイヤーの状態です。
	// 注: バトル中のHP等は保存せず、装備エージェントから再計算
	Player *PlayerSaveData `json:"player"`

	// Inventory は所持アイテム（コア、モジュール、エージェント）です。

	Inventory *InventorySaveData `json:"inventory"`

	// EnemyProgress は敵進行データです。

	EnemyProgress *EnemyProgressSave `json:"enemy_progress"`

	// Statistics は統計データです。

	Statistics *StatisticsSaveData `json:"statistics"`

	// Achievements は実績データです。

	Achievements *AchievementsSaveData `json:"achievements"`

	// Settings はゲーム設定です。

	Settings *SettingsSaveData `json:"settings"`
}

// PlayerSaveData はプレイヤーのセーブデータです。
type PlayerSaveData struct {
	// MaxHP はプレイヤーの最大HPです。
	// 新規ゲーム開始時は1000、敵撃破による成長で増加します。
	MaxHP int `json:"max_hp"`

	// AgentSlots はエージェントスロット構成です（3スロット）。
	// コア・スキルの自由付け替え方式で管理します。
	AgentSlots [3]AgentSlotSave `json:"agent_slots"`
}

// ChainEffectSave はチェイン効果のセーブデータです。
type ChainEffectSave struct {
	// Type はチェイン効果の種別（damage_bonus, heal_bonus等）です。
	Type string `json:"type"`

	// Value は効果量です。
	Value float64 `json:"value"`
}

// SkillInstanceSave はスキルインスタンスのセーブデータです。
// v2.0.0ではTypeIDとChainEffectのペアとして永続化します。
// 同一TypeIDでも異なるChainEffectを持つことを許容します。
type SkillInstanceSave struct {
	// TypeID はスキル種別ID（マスタデータ参照用）です。
	TypeID string `json:"type_id"`

	// ChainEffect はこのスキルインスタンスのチェイン効果です。
	// nilの場合はチェイン効果なしとしてomitemptyで省略されます。
	ChainEffect *ChainEffectSave `json:"chain_effect,omitempty"`
}

// ModuleInstanceSave はSkillInstanceSaveのエイリアスです。
type ModuleInstanceSave = SkillInstanceSave

// CoreInventorySave はユニークコアインベントリのセーブデータです。
// TypeIDリスト形式で管理します。
type CoreInventorySave struct {
	// Cores は保有CoreTypeIDリストです。
	Cores []string `json:"cores"`
}

// SkillInventorySave はユニークスキルインベントリのセーブデータです。
// TypeIDリスト形式で管理します。
type SkillInventorySave struct {
	// Skills は保有SkillTypeIDリストです。
	Skills []string `json:"skills"`
}

// SkillSlotSaveCfg はスキルスロット構成のセーブデータです。
type SkillSlotSaveCfg struct {
	// TypeID はスキルTypeID（空の場合は未設定）です。
	TypeID string `json:"type_id,omitempty"`
}

// ChainEffectInventorySave はチェイン効果インベントリのセーブデータです。
type ChainEffectInventorySave struct {
	// ChainEffects は保有チェイン効果TypeIDリストです。
	ChainEffects []string `json:"chain_effects"`
}

// ChainEffectSlotSaveCfg はチェイン効果スロット構成のセーブデータです。
type ChainEffectSlotSaveCfg struct {
	// TypeID はチェイン効果TypeID（空の場合は未設定）です。
	TypeID string `json:"type_id,omitempty"`
}

// AgentSlotSave はエージェントスロットのセーブデータです。
// コアはTypeIDのみで管理します。
type AgentSlotSave struct {
	// CoreTypeID はコアTypeID（空の場合はスロット空）です。
	CoreTypeID string `json:"core_type_id,omitempty"`

	// Skills はスキルスロット構成（4つ）です。
	Skills [4]SkillSlotSaveCfg `json:"skills"`

	// ChainEffects はチェイン効果スロット構成（4つ）です。
	ChainEffects [4]ChainEffectSlotSaveCfg `json:"chain_effects"`
}

// InventorySaveData はインベントリのセーブデータです。
type InventorySaveData struct {
	// UniqueCores はユニークコアインベントリです。
	UniqueCores *CoreInventorySave `json:"unique_cores,omitempty"`

	// UniqueSkills はユニークスキルインベントリです。
	UniqueSkills *SkillInventorySave `json:"skills,omitempty"`

	// UniqueChainEffects はユニークチェイン効果インベントリです。
	UniqueChainEffects *ChainEffectInventorySave `json:"unique_chain_effects,omitempty"`
}

// StatisticsSaveData は統計のセーブデータです。
type StatisticsSaveData struct {
	// TotalBattles は総バトル数です。
	TotalBattles int `json:"total_battles"`

	// Victories は勝利数です。
	Victories int `json:"victories"`

	// Defeats は敗北数です。
	Defeats int `json:"defeats"`

	// MaxLevelReached は到達最高レベルです。

	MaxLevelReached int `json:"max_level_reached"`

	// HighestWPM は最高WPMです。
	HighestWPM float64 `json:"highest_wpm"`

	// AverageWPM は平均WPMです。
	AverageWPM float64 `json:"average_wpm"`

	// PerfectAccuracyCount は100%正確性達成回数です。
	PerfectAccuracyCount int `json:"perfect_accuracy_count"`

	// TotalCharactersTyped は総タイプ文字数です。
	TotalCharactersTyped int `json:"total_characters_typed"`

	// EncounteredEnemies はエンカウントした敵のIDリストです（敵図鑑用）。
	EncounteredEnemies []string `json:"encountered_enemies"`

	// DefeatedEnemies は撃破済み敵の情報です（敵タイプID→撃破最高レベル）。
	DefeatedEnemies map[string]int `json:"defeated_enemies,omitempty"`
}

// AchievementsSaveData は実績のセーブデータです。
type AchievementsSaveData struct {
	// Unlocked は解除済み実績IDリストです。
	Unlocked []string `json:"unlocked"`

	// Progress は実績の進捗状況です（実績ID→進捗値）。
	Progress map[string]int `json:"progress"`
}

// SettingsSaveData は設定のセーブデータです。
type SettingsSaveData struct {
	// KeyBindings はキーバインド設定です。
	KeyBindings map[string]string `json:"key_bindings"`
}

// EnemyProgressSave は敵進行データのセーブデータです。
type EnemyProgressSave struct {
	// CurrentRank は現在解放済みランクです（1から開始）。
	CurrentRank int `json:"current_rank"`

	// DefeatRecords は敵タイプIDごとの撃破記録です。
	DefeatRecords map[string]DefeatRecordSave `json:"defeat_records"`
}

// DefeatRecordSave は敵1体の撃破記録のセーブデータです。
type DefeatRecordSave struct {
	// Defeated は撃破済みフラグです。
	Defeated bool `json:"defeated"`

	// MaxDefeatedLevel は撃破済み最大レベルです（未撃破なら0）。
	MaxDefeatedLevel int `json:"max_defeated_level"`
}

// NewSaveData は新しいセーブデータを作成します。
func NewSaveData() *SaveData {
	return &SaveData{
		Version:   CurrentSaveDataVersion,
		Timestamp: time.Now(),
		Player: &PlayerSaveData{
			MaxHP:      1000,
			AgentSlots: [3]AgentSlotSave{},
		},
		Inventory: &InventorySaveData{
			UniqueCores: &CoreInventorySave{
				Cores: make([]string, 0),
			},
			UniqueSkills: &SkillInventorySave{
				Skills: make([]string, 0),
			},
			UniqueChainEffects: &ChainEffectInventorySave{
				ChainEffects: make([]string, 0),
			},
		},
		EnemyProgress: &EnemyProgressSave{
			CurrentRank:   1,
			DefeatRecords: make(map[string]DefeatRecordSave),
		},
		Statistics: &StatisticsSaveData{
			TotalBattles:         0,
			Victories:            0,
			Defeats:              0,
			MaxLevelReached:      0,
			HighestWPM:           0,
			AverageWPM:           0,
			PerfectAccuracyCount: 0,
			TotalCharactersTyped: 0,
		},
		Achievements: &AchievementsSaveData{
			Unlocked: make([]string, 0),
			Progress: make(map[string]int),
		},
		Settings: &SettingsSaveData{
			KeyBindings: make(map[string]string),
		},
	}
}

// SaveDataIO はセーブデータのI/Oを担当する構造体です。
type SaveDataIO struct {
	// saveDir はセーブファイルを保存するディレクトリパスです。
	saveDir string
	// saveFileName はセーブファイル名です（通常: save.json、デバッグ: debug_save.json）。
	saveFileName string
	// tempSaveFileName は一時セーブファイル名です。
	tempSaveFileName string
}

// NewSaveDataIO は新しいSaveDataIOを作成します。
// debugModeがtrueの場合、デバッグ専用のセーブファイル（debug_save.json）を使用します。
func NewSaveDataIO(saveDir string, debugMode bool) *SaveDataIO {
	saveFileName := SaveFileName
	tempFileName := TempSaveFileName
	if debugMode {
		saveFileName = DebugSaveFileName
		tempFileName = DebugTempSaveFileName
	}
	return &SaveDataIO{
		saveDir:          saveDir,
		saveFileName:     saveFileName,
		tempSaveFileName: tempFileName,
	}
}

// SaveGame はセーブデータをファイルに保存します。

func (io *SaveDataIO) SaveGame(data *SaveData) error {
	// タイムスタンプを更新
	data.Timestamp = time.Now()

	// JSONにシリアライズ
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("セーブデータのシリアライズに失敗: %w", err)
	}

	// ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(io.saveDir, 0755); err != nil {
		return fmt.Errorf("セーブディレクトリの作成に失敗: %w", err)
	}

	// 一時ファイルに書き込み
	tmpPath := filepath.Join(io.saveDir, io.tempSaveFileName)
	if err := os.WriteFile(tmpPath, jsonData, 0644); err != nil {
		return fmt.Errorf("一時ファイルへの書き込みに失敗: %w", err)
	}

	// 一時ファイルの検証（読み込んでパースできるか確認）
	tmpData, err := os.ReadFile(tmpPath)
	if err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("一時ファイルの検証読み込みに失敗: %w", err)
	}
	var validateData SaveData
	if err := json.Unmarshal(tmpData, &validateData); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("一時ファイルの検証パースに失敗: %w", err)
	}

	// バックアップローテーション（失敗してもセーブは続行）
	_ = io.RotateBackups()

	// 原子的リネーム
	savePath := filepath.Join(io.saveDir, io.saveFileName)
	if err := os.Rename(tmpPath, savePath); err != nil {
		return fmt.Errorf("セーブファイルのリネームに失敗: %w", err)
	}

	return nil
}

// LoadGame はセーブデータをファイルから読み込みます。

func (io *SaveDataIO) LoadGame() (*SaveData, error) {
	savePath := filepath.Join(io.saveDir, io.saveFileName)

	// メインのセーブファイルを読み込み
	data, err := io.loadFromFile(savePath)
	if err == nil {
		return data, nil
	}

	// メインファイルの読み込みに失敗した場合、バックアップから復元を試みる
	for i := 1; i <= MaxBackupCount; i++ {
		backupPath := filepath.Join(io.saveDir, io.backupFileName(i))
		data, err := io.loadFromFile(backupPath)
		if err == nil {
			// バックアップからの復元に成功
			// メインファイルを復元
			if jsonData, marshalErr := json.MarshalIndent(data, "", "  "); marshalErr == nil {
				_ = os.WriteFile(savePath, jsonData, 0644)
			}
			return data, nil
		}
	}

	return nil, fmt.Errorf("セーブデータのロードに失敗: %w", err)
}

// backupFileName はバックアップファイル名を生成します。
func (io *SaveDataIO) backupFileName(index int) string {
	return fmt.Sprintf("%s.bak%d", io.saveFileName, index)
}

// loadFromFile は指定されたファイルからセーブデータを読み込みます。
func (io *SaveDataIO) loadFromFile(filePath string) (*SaveData, error) {
	return loadAndParseSaveFile(filePath)
}

// loadAndParseSaveFile はファイルの読み込み、JSONパース、バージョン検証を行う共通処理です。
func loadAndParseSaveFile(filePath string) (*SaveData, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みに失敗: %w", err)
	}

	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("JSONパースに失敗: %w", err)
	}

	if err := ValidateSaveVersion(data.Version); err != nil {
		return nil, err
	}

	return &data, nil
}

// LoadSaveDataFromFile は指定パスからセーブデータを読み込みます。
// -saveフラグで指定されたテスト用セーブファイルのロードに使用します。
func LoadSaveDataFromFile(filePath string) (*SaveData, error) {
	data, err := loadAndParseSaveFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := ValidateSaveData(data); err != nil {
		return nil, err
	}

	return data, nil
}

// LoadFromBackup は指定したバックアップインデックスからセーブデータを読み込みます。
func (io *SaveDataIO) LoadFromBackup(backupIndex int) (*SaveData, error) {
	if backupIndex < 1 || backupIndex > MaxBackupCount {
		return nil, fmt.Errorf("不正なバックアップインデックス: %d", backupIndex)
	}

	backupPath := filepath.Join(io.saveDir, io.backupFileName(backupIndex))
	return io.loadFromFile(backupPath)
}

// RotateBackups はバックアップファイルをローテーションします。
// {saveFileName} → {saveFileName}.bak1 → {saveFileName}.bak2 → {saveFileName}.bak3 (削除)
func (io *SaveDataIO) RotateBackups() error {
	// 古いバックアップを削除（存在しない場合は無視）
	bak3 := filepath.Join(io.saveDir, io.backupFileName(MaxBackupCount))
	_ = os.Remove(bak3)

	// バックアップをシフト
	for i := MaxBackupCount - 1; i >= 1; i-- {
		oldPath := filepath.Join(io.saveDir, io.backupFileName(i))
		newPath := filepath.Join(io.saveDir, io.backupFileName(i+1))
		if _, err := os.Stat(oldPath); err == nil {
			_ = os.Rename(oldPath, newPath)
		}
	}

	// 現在のセーブファイルをbak1にコピー
	savePath := filepath.Join(io.saveDir, io.saveFileName)
	bak1 := filepath.Join(io.saveDir, io.backupFileName(1))
	if _, err := os.Stat(savePath); err == nil {
		data, err := os.ReadFile(savePath)
		if err == nil {
			_ = os.WriteFile(bak1, data, 0644)
		}
	}

	return nil
}

// ResetSaveData はセーブデータをリセットします。

func (io *SaveDataIO) ResetSaveData() error {
	// メインセーブファイルを削除
	savePath := filepath.Join(io.saveDir, io.saveFileName)
	_ = os.Remove(savePath)

	// バックアップファイルも削除
	for i := 1; i <= MaxBackupCount; i++ {
		backupPath := filepath.Join(io.saveDir, io.backupFileName(i))
		_ = os.Remove(backupPath)
	}

	return nil
}

// ValidateSaveData はセーブデータのバリデーションを行います。

func ValidateSaveData(data *SaveData) error {
	if data.Version == "" {
		return fmt.Errorf("セーブデータのバージョンが空です")
	}
	if data.Player == nil {
		return fmt.Errorf("プレイヤーデータがnilです")
	}
	if data.Inventory == nil {
		return fmt.Errorf("インベントリデータがnilです")
	}
	if data.Statistics == nil {
		return fmt.Errorf("統計データがnilです")
	}
	if data.Achievements == nil {
		return fmt.Errorf("実績データがnilです")
	}
	if data.Settings == nil {
		return fmt.Errorf("設定データがnilです")
	}
	return nil
}

// ValidateSaveVersion はセーブデータのバージョンが現在のバージョンと互換性があるか検証します。
func ValidateSaveVersion(version string) error {
	if version == "" {
		return fmt.Errorf("セーブデータのバージョンが空です")
	}

	if version != CurrentSaveDataVersion {
		return fmt.Errorf("セーブデータのバージョン %s はサポートされていません。v%s が必要です", version, CurrentSaveDataVersion)
	}

	return nil
}

// Exists はセーブファイルが存在するかどうかを確認します。
func (io *SaveDataIO) Exists() bool {
	savePath := filepath.Join(io.saveDir, io.saveFileName)
	_, err := os.Stat(savePath)
	return err == nil
}
