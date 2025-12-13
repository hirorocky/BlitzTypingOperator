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
// セーブデータの形式が変更された場合にインクリメントします。
// v2.0.0: ID化最適化（フルオブジェクト → ID参照）
const CurrentSaveDataVersion = "2.0.0"

// SaveFileName はセーブファイル名です。
const SaveFileName = "save.json"

// TempSaveFileName は一時セーブファイル名です（原子的書き込み用）。
const TempSaveFileName = "save.json.tmp"

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

	// Statistics は統計データです。

	Statistics *StatisticsSaveData `json:"statistics"`

	// Achievements は実績データです。

	Achievements *AchievementsSaveData `json:"achievements"`

	// Settings はゲーム設定です。

	Settings *SettingsSaveData `json:"settings"`
}

// PlayerSaveData はプレイヤーのセーブデータです。
type PlayerSaveData struct {
	// EquippedAgentIDs は装備中のエージェントIDリストです（スロット番号順）。
	// 空きスロットは空文字列で表現されます。

	EquippedAgentIDs [3]string `json:"equipped_agent_ids"`
}

// CoreInstanceSave はコアインスタンスの軽量セーブデータです。
// フルオブジェクト保存ではなく、ID+Type+Levelのみ保存し、
// ロード時にマスタデータからステータスを再計算します。
type CoreInstanceSave struct {
	// ID はコアインスタンスの一意識別子です。
	ID string `json:"id"`

	// CoreTypeID はコア特性ID（例: "all_rounder"）です。
	// マスタデータ（cores.json）から参照します。
	CoreTypeID string `json:"core_type_id"`

	// Level はコアのレベルです。
	// ステータスはレベルとコア特性から再計算されます。
	Level int `json:"level"`
}

// AgentInstanceSave はエージェントインスタンスの軽量セーブデータです。
// コア情報を直接保持し、インベントリに依存せずにロード時に再構築します。
type AgentInstanceSave struct {
	// ID はエージェントインスタンスの一意識別子です。
	ID string `json:"id"`

	// Core はエージェントが使用するコアの情報です。
	// エージェント合成時にコアは消費されるため、インベントリとは別に保持します。
	Core CoreInstanceSave `json:"core"`

	// ModuleIDs はモジュール定義IDリストです（4つ）。
	// マスタデータ（modules.json）から参照します。
	ModuleIDs []string `json:"module_ids"`
}

// InventorySaveData はインベントリのセーブデータです。
// ID化最適化により、フルオブジェクトではなくID参照を保存します。
type InventorySaveData struct {
	// CoreInstances は所持コアのID+Type+Levelリストです。
	CoreInstances []CoreInstanceSave `json:"core_instances"`

	// ModuleCounts は所持モジュールの個数マップです（モジュールID → 個数）。
	// モジュールにはインスタンス固有データがないため、IDとカウントで十分です。
	ModuleCounts map[string]int `json:"module_counts"`

	// AgentInstances は所持エージェントの参照リストです。
	AgentInstances []AgentInstanceSave `json:"agent_instances"`

	// MaxCoreSlots はコアの最大所持数です。
	MaxCoreSlots int `json:"max_core_slots"`

	// MaxModuleSlots はモジュールの最大所持数です。
	MaxModuleSlots int `json:"max_module_slots"`

	// MaxAgentSlots はエージェントの最大所持数です。

	MaxAgentSlots int `json:"max_agent_slots"`
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

// NewSaveData は新しいセーブデータを作成します。
func NewSaveData() *SaveData {
	return &SaveData{
		Version:   CurrentSaveDataVersion,
		Timestamp: time.Now(),
		Player: &PlayerSaveData{
			EquippedAgentIDs: [3]string{},
		},
		Inventory: &InventorySaveData{
			CoreInstances:  make([]CoreInstanceSave, 0),
			ModuleCounts:   make(map[string]int),
			AgentInstances: make([]AgentInstanceSave, 0),
			MaxCoreSlots:   100,
			MaxModuleSlots: 200,
			MaxAgentSlots:  20,
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
}

// NewSaveDataIO は新しいSaveDataIOを作成します。
func NewSaveDataIO(saveDir string) *SaveDataIO {
	return &SaveDataIO{
		saveDir: saveDir,
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
	tmpPath := filepath.Join(io.saveDir, TempSaveFileName)
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
	savePath := filepath.Join(io.saveDir, SaveFileName)
	if err := os.Rename(tmpPath, savePath); err != nil {
		return fmt.Errorf("セーブファイルのリネームに失敗: %w", err)
	}

	return nil
}

// LoadGame はセーブデータをファイルから読み込みます。

func (io *SaveDataIO) LoadGame() (*SaveData, error) {
	savePath := filepath.Join(io.saveDir, SaveFileName)

	// メインのセーブファイルを読み込み
	data, err := io.loadFromFile(savePath)
	if err == nil {
		return data, nil
	}

	// メインファイルの読み込みに失敗した場合、バックアップから復元を試みる
	for i := 1; i <= MaxBackupCount; i++ {
		backupPath := filepath.Join(io.saveDir, fmt.Sprintf("save.json.bak%d", i))
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

// loadFromFile は指定されたファイルからセーブデータを読み込みます。
func (io *SaveDataIO) loadFromFile(filePath string) (*SaveData, error) {
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイル読み込みに失敗: %w", err)
	}

	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("JSONパースに失敗: %w", err)
	}

	// バージョンチェック
	if data.Version == "" {
		return nil, fmt.Errorf("セーブデータのバージョンが不正です")
	}

	// 将来的なバージョンマイグレーションはここで実装

	return &data, nil
}

// LoadFromBackup は指定したバックアップインデックスからセーブデータを読み込みます。
func (io *SaveDataIO) LoadFromBackup(backupIndex int) (*SaveData, error) {
	if backupIndex < 1 || backupIndex > MaxBackupCount {
		return nil, fmt.Errorf("不正なバックアップインデックス: %d", backupIndex)
	}

	backupPath := filepath.Join(io.saveDir, fmt.Sprintf("save.json.bak%d", backupIndex))
	return io.loadFromFile(backupPath)
}

// RotateBackups はバックアップファイルをローテーションします。

// save.json → save.json.bak1 → save.json.bak2 → save.json.bak3 (削除)
func (io *SaveDataIO) RotateBackups() error {
	// 古いバックアップを削除（存在しない場合は無視）
	bak3 := filepath.Join(io.saveDir, "save.json.bak3")
	_ = os.Remove(bak3)

	// バックアップをシフト
	for i := MaxBackupCount - 1; i >= 1; i-- {
		oldPath := filepath.Join(io.saveDir, fmt.Sprintf("save.json.bak%d", i))
		newPath := filepath.Join(io.saveDir, fmt.Sprintf("save.json.bak%d", i+1))
		if _, err := os.Stat(oldPath); err == nil {
			_ = os.Rename(oldPath, newPath)
		}
	}

	// 現在のセーブファイルをbak1にコピー
	savePath := filepath.Join(io.saveDir, SaveFileName)
	bak1 := filepath.Join(io.saveDir, "save.json.bak1")
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
	savePath := filepath.Join(io.saveDir, SaveFileName)
	_ = os.Remove(savePath)

	// バックアップファイルも削除
	for i := 1; i <= MaxBackupCount; i++ {
		backupPath := filepath.Join(io.saveDir, fmt.Sprintf("save.json.bak%d", i))
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

// Exists はセーブファイルが存在するかどうかを確認します。
func (io *SaveDataIO) Exists() bool {
	savePath := filepath.Join(io.saveDir, SaveFileName)
	_, err := os.Stat(savePath)
	return err == nil
}
