// Package startup は初回起動時の初期化処理を担当します。
// 新規ゲーム開始時の初期コア、モジュール、エージェントの提供を行います。

package startup

import (
	"fmt"

	"github.com/google/uuid"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/infra/savedata"
)

// 初期モジュールID定義（マスタデータのmodules.jsonと一致させる）
var initialModuleIDs = []string{
	"physical_strike_lv1",
	"fireball_lv1",
	"heal_lv1",
	"attack_buff_lv1",
}

// 初期コア特性ID（マスタデータのcores.jsonと一致させる）
const initialCoreTypeID = "all_rounder"

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

// generateUUID は一意のUUIDを生成します。
func generateUUID() string {
	return uuid.New().String()
}

// CreateInitialCore は初期コアを作成します。

func (i *NewGameInitializer) CreateInitialCore() *domain.CoreModel {
	// マスタデータから"all_rounder"コア特性を検索
	var coreType domain.CoreType
	if i.externalData != nil {
		for _, ct := range i.externalData.CoreTypes {
			if ct.ID == initialCoreTypeID {
				coreType = ct.ToDomain()
				break
			}
		}
	}

	// 見つからない場合のフォールバック（外部データがない場合も含む）
	if coreType.ID == "" {
		coreType = domain.CoreType{
			ID:   initialCoreTypeID,
			Name: "オールラウンダー",
			StatWeights: map[string]float64{
				"STR": 1.0,
				"MAG": 1.0,
				"SPD": 1.0,
				"LUK": 1.0,
			},
			PassiveSkillID: "adaptability",
			AllowedTags: []string{
				"physical_low", "magic_low", "heal_low", "buff_low", "debuff_low",
			},
			MinDropLevel: 1,
		}
	}

	// パッシブスキルの定義（マスタデータのpassive_skill_idに基づく）
	passiveSkill := domain.PassiveSkill{
		ID:          coreType.PassiveSkillID,
		Name:        "適応力",
		Description: "全ステータスがバランス良く成長",
	}

	// レベル1のコアを作成
	return domain.NewCore(
		generateUUID(),
		"初期コア",
		1, // レベル1
		coreType,
		passiveSkill,
	)
}

// CreateInitialModules は初期モジュールを作成します。

func (i *NewGameInitializer) CreateInitialModules() []*domain.ModuleModel {
	modules := make([]*domain.ModuleModel, 0, len(initialModuleIDs))

	// マスタデータから初期モジュールを検索
	if i.externalData != nil {
		for _, moduleID := range initialModuleIDs {
			for _, md := range i.externalData.ModuleDefinitions {
				if md.ID == moduleID {
					modules = append(modules, md.ToDomain())
					break
				}
			}
		}
	}

	// 外部データがない場合やモジュールが見つからない場合のフォールバック
	if len(modules) == 0 {
		// デフォルトのモジュールを作成
		modules = []*domain.ModuleModel{
			domain.NewModule("physical_strike_lv1", "物理打撃Lv1", domain.PhysicalAttack, 1, []string{"physical_low"}, 10.0, "STR", "物理ダメージを与える基本攻撃"),
			domain.NewModule("fireball_lv1", "ファイアボールLv1", domain.MagicAttack, 1, []string{"magic_low"}, 12.0, "MAG", "魔法ダメージを与える基本魔法"),
			domain.NewModule("heal_lv1", "ヒールLv1", domain.Heal, 1, []string{"heal_low"}, 8.0, "MAG", "HPを回復する基本回復魔法"),
			domain.NewModule("attack_buff_lv1", "攻撃バフLv1", domain.Buff, 1, []string{"buff_low"}, 5.0, "SPD", "一時的に攻撃力を上昇させる"),
		}
	} else if len(modules) < len(initialModuleIDs) {
		// ログを出力（デバッグ用）
		fmt.Printf("警告: 一部の初期モジュールがマスタデータに見つかりませんでした (%d/%d)\n",
			len(modules), len(initialModuleIDs))
	}

	return modules
}

// CreateInitialAgent は初期エージェントを作成します。

func (i *NewGameInitializer) CreateInitialAgent() *domain.AgentModel {
	core := i.CreateInitialCore()
	modules := i.CreateInitialModules()

	return domain.NewAgent(
		generateUUID(),
		core,
		modules,
	)
}

// InitializeNewGame は新規ゲームを初期化してセーブデータを作成します。

// ID化最適化に対応：フルオブジェクトではなくID参照を保存
func (i *NewGameInitializer) InitializeNewGame() *savedata.SaveData {
	// 基本のセーブデータを作成
	saveData := savedata.NewSaveData()

	// 初期エージェントを作成
	initialAgent := i.CreateInitialAgent()

	// インベントリにエージェントを追加（コア情報を直接埋め込み）
	moduleIDs := make([]string, len(initialAgent.Modules))
	for idx, m := range initialAgent.Modules {
		moduleIDs[idx] = m.ID
	}
	saveData.Inventory.AgentInstances = []savedata.AgentInstanceSave{
		{
			ID: initialAgent.ID,
			Core: savedata.CoreInstanceSave{
				ID:         initialAgent.Core.ID,
				CoreTypeID: initialAgent.Core.Type.ID,
				Level:      initialAgent.Core.Level,
			},
			ModuleIDs: moduleIDs,
		},
	}

	// インベントリのコアは空（エージェントのコアはエージェント内に保持される）
	saveData.Inventory.CoreInstances = []savedata.CoreInstanceSave{}

	// コアとモジュールはエージェント合成で消費されるため、インベントリには追加しない
	// （エージェントに含まれているコアとモジュールは参照として保持される）

	// 初期エージェントを装備（スロット0に装備）
	saveData.Player.EquippedAgentIDs = [3]string{initialAgent.ID, "", ""}

	return saveData
}

// CreateNewGameWithExtraItems は追加アイテム付きで新規ゲームを初期化します。
// デバッグや特殊条件での開始用
// ID化最適化に対応：フルオブジェクトではなくID参照を保存
func (i *NewGameInitializer) CreateNewGameWithExtraItems() *savedata.SaveData {
	saveData := i.InitializeNewGame()

	// 追加のコアを作成してインベントリにID化して追加
	extraCore := i.CreateInitialCore()
	saveData.Inventory.CoreInstances = append(saveData.Inventory.CoreInstances, savedata.CoreInstanceSave{
		ID:         extraCore.ID,
		CoreTypeID: extraCore.Type.ID,
		Level:      extraCore.Level,
	})

	// 追加のモジュールを作成してインベントリにカウント追加
	extraModules := i.CreateInitialModules()
	for _, module := range extraModules {
		saveData.Inventory.ModuleCounts[module.ID]++
	}

	return saveData
}
