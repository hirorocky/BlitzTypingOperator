// Package app は BlitzTypingOperator TUIゲームのRootModelを提供します。
package app

import (
	"sort"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/screens"
	gamestate "hirorocky/type-battle/internal/usecase/session"
)

// DefaultEnemyLevel は未撃破敵のデフォルトレベルです。
const DefaultEnemyLevel = 1

// EnemyProgressAdapter はGameStateとenemyTypesを組み合わせて
// screens.EnemyProgressProviderインターフェースを実装します。
type EnemyProgressAdapter struct {
	gameState  *gamestate.GameState
	enemyTypes map[string]domain.EnemyType
	coreTypes  map[string]domain.CoreType
	skillTypes map[string]domain.SkillType
}

// NewEnemyProgressAdapter は新しいEnemyProgressAdapterを作成します。
func NewEnemyProgressAdapter(
	gs *gamestate.GameState,
	enemyTypes map[string]domain.EnemyType,
	coreTypes map[string]domain.CoreType,
	skillTypes map[string]domain.SkillType,
) *EnemyProgressAdapter {
	return &EnemyProgressAdapter{
		gameState:  gs,
		enemyTypes: enemyTypes,
		coreTypes:  coreTypes,
		skillTypes: skillTypes,
	}
}

// GetCurrentRank は現在解放済みランクを返します。
func (a *EnemyProgressAdapter) GetCurrentRank() int {
	return a.gameState.EnemyProgress().CurrentRank
}

// GetMaxUnlockedRank は解放済み最大ランクを返します。
// プレイヤーの進行度で解放されているランクの最大値です。
func (a *EnemyProgressAdapter) GetMaxUnlockedRank() int {
	return a.gameState.EnemyProgress().CurrentRank
}

// GetCurrentRankEnemies は現在ランクの敵リストを返します。
// 結果はIDでソートされ、一貫した順序で返されます。
func (a *EnemyProgressAdapter) GetCurrentRankEnemies() []domain.EnemyType {
	currentRank := a.gameState.EnemyProgress().CurrentRank
	return a.GetEnemiesByRank(currentRank)
}

// GetEnemiesByRank は指定ランクの敵リストを返します。
// 結果はIDでソートされ、一貫した順序で返されます。
func (a *EnemyProgressAdapter) GetEnemiesByRank(rank int) []domain.EnemyType {
	var enemies []domain.EnemyType
	for _, et := range a.enemyTypes {
		if et.Rank == rank {
			enemies = append(enemies, et)
		}
	}
	// IDでソートして順序を固定
	sort.Slice(enemies, func(i, j int) bool {
		return enemies[i].ID < enemies[j].ID
	})
	return enemies
}

// IsDefeated は指定した敵が撃破済みかどうかを返します。
func (a *EnemyProgressAdapter) IsDefeated(enemyTypeID string) bool {
	return a.gameState.EnemyProgress().IsDefeated(enemyTypeID)
}

// GetMaxDefeatedLevel は指定した敵の撃破済み最大レベルを返します。
func (a *EnemyProgressAdapter) GetMaxDefeatedLevel(enemyTypeID string) int {
	return a.gameState.EnemyProgress().GetMaxDefeatedLevel(enemyTypeID)
}

// GetSelectableLevelRange は敵の選択可能レベル範囲を取得します。
// 未撃破の場合はレベル1のみ、撃破済みの場合は1〜(撃破最高レベル+ボーナス)です。
func (a *EnemyProgressAdapter) GetSelectableLevelRange(enemyTypeID string) (min, max int) {
	_, ok := a.enemyTypes[enemyTypeID]
	if !ok {
		return 1, 1
	}

	// 未撃破時はデフォルトレベル1を使用
	return a.gameState.EnemyProgress().GetSelectableLevelRange(enemyTypeID, DefaultEnemyLevel)
}

// GetDropItemName はドロップアイテムの日本語名を返します。
// category: "core" または "module"
// typeID: コアまたはスキルのTypeID
func (a *EnemyProgressAdapter) GetDropItemName(category, typeID string) string {
	switch category {
	case "core":
		if ct, ok := a.coreTypes[typeID]; ok {
			return ct.Name
		}
		return typeID
	case "module":
		if st, ok := a.skillTypes[typeID]; ok {
			return st.Name
		}
		return typeID
	default:
		return typeID
	}
}

// インターフェース準拠チェック
var _ screens.EnemyProgressProvider = (*EnemyProgressAdapter)(nil)
