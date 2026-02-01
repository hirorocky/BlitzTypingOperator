// Package enemy は敵生成システムのテストを提供します。

package spawning

import (
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"
)

// TestEnemyStats_HPCalculation はレベルに応じたHP計算をテストします。
func TestEnemyStats_HPCalculation(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// レベル1
	enemy1 := generator.Generate(1)
	if enemy1 == nil {
		t.Fatal("敵の生成に失敗")
	}

	// レベル10
	enemy10 := generator.Generate(10)
	if enemy10 == nil {
		t.Fatal("敵の生成に失敗")
	}

	// レベルが高いほどHPが高い
	if enemy10.MaxHP <= enemy1.MaxHP {
		t.Errorf("レベル10の敵はレベル1より高いHPを持つべき: Lv1=%d, Lv10=%d", enemy1.MaxHP, enemy10.MaxHP)
	}

	// HP計算式: BaseHP * level
	expectedHP1 := 50 * 1
	if enemy1.MaxHP != expectedHP1 {
		t.Errorf("レベル1のHP計算が不正: got %d, want %d", enemy1.MaxHP, expectedHP1)
	}

	expectedHP10 := 50 * 10
	if enemy10.MaxHP != expectedHP10 {
		t.Errorf("レベル10のHP計算が不正: got %d, want %d", enemy10.MaxHP, expectedHP10)
	}
}

// TestEnemyStats_AttackPowerCalculation はレベルに応じた攻撃力計算をテストします。
func TestEnemyStats_AttackPowerCalculation(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// レベル1
	enemy1 := generator.Generate(1)

	// レベル20
	enemy20 := generator.Generate(20)

	// レベルが高いほど攻撃力が高い
	if enemy20.AttackPower <= enemy1.AttackPower {
		t.Errorf("レベル20の敵はレベル1より高い攻撃力を持つべき: Lv1=%d, Lv20=%d",
			enemy1.AttackPower, enemy20.AttackPower)
	}

	// 攻撃力計算式: BaseAttackPower + (level * 2)
	expectedAttack1 := 5 + (1 * 2)
	if enemy1.AttackPower != expectedAttack1 {
		t.Errorf("レベル1の攻撃力計算が不正: got %d, want %d", enemy1.AttackPower, expectedAttack1)
	}

	expectedAttack20 := 5 + (20 * 2)
	if enemy20.AttackPower != expectedAttack20 {
		t.Errorf("レベル20の攻撃力計算が不正: got %d, want %d", enemy20.AttackPower, expectedAttack20)
	}
}

// TestEnemyVariation_RandomSelection は敵タイプからのランダム選択をテストします。
func TestEnemyVariation_RandomSelection(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
		{
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 80,
			Rank:   1,
		},
		{
			ID:     "skeleton",
			Name:   "スケルトン",
			BaseHP: 70,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// 複数回生成して異なるタイプが出ることを確認
	typeOccurrences := make(map[string]int)
	for i := 0; i < 100; i++ {
		enemy := generator.Generate(10)
		typeOccurrences[enemy.Type.ID]++
	}

	// 3種類の敵タイプがあるので、100回生成すれば全種類が出現するはず（確率的）
	if len(typeOccurrences) < 2 {
		t.Errorf("ランダム選択が機能していない可能性: 出現タイプ数=%d", len(typeOccurrences))
	}
}

// TestEnemyVariation_SameLevelMultipleTypes は同レベルでの複数バリエーション対応をテストします。
func TestEnemyVariation_SameLevelMultipleTypes(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
		{
			ID:     "goblin",
			Name:   "ゴブリン",
			BaseHP: 80,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// 同じレベルでも異なるタイプが出ることを確認
	level := 5
	seenSlime := false
	seenGoblin := false

	for i := 0; i < 50; i++ {
		enemy := generator.Generate(level)
		if enemy.Type.ID == "slime" {
			seenSlime = true
		}
		if enemy.Type.ID == "goblin" {
			seenGoblin = true
		}
		if seenSlime && seenGoblin {
			break
		}
	}

	if !seenSlime || !seenGoblin {
		t.Error("同レベルで複数のバリエーションが出現すべき")
	}
}

// TestEnemyLevel_Maximum はレベル上限をテストします。
func TestEnemyLevel_Maximum(t *testing.T) {
	if MaxEnemyLevel != 100 {
		t.Errorf("レベル上限が100であるべき: got %d", MaxEnemyLevel)
	}

	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// レベル上限を超えた値を指定しても上限でクランプされる
	enemy := generator.Generate(150)
	if enemy.Level > MaxEnemyLevel {
		t.Errorf("レベルが上限を超えている: got %d, max %d", enemy.Level, MaxEnemyLevel)
	}
}

// TestEnemyLevel_MaxLevelDefeat は最高レベル敵撃破時のゲームクリア判定をテストします。
func TestEnemyLevel_MaxLevelDefeat(t *testing.T) {
	generator := NewEnemyGenerator(nil)

	// レベル100の敵撃破はゲームクリア
	if !generator.IsMaxLevelBattle(MaxEnemyLevel) {
		t.Error("レベル100はゲームクリア対象であるべき")
	}

	// レベル99はゲームクリアではない
	if generator.IsMaxLevelBattle(99) {
		t.Error("レベル99はゲームクリア対象ではない")
	}
}

// TestEnemyLevel_ValidRange は有効なレベル範囲をテストします。
func TestEnemyLevel_ValidRange(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// 最低レベルは1
	enemy := generator.Generate(0)
	if enemy.Level < 1 {
		t.Errorf("レベルは最低1であるべき: got %d", enemy.Level)
	}

	// レベル-1も1にクランプ
	enemy = generator.Generate(-5)
	if enemy.Level < 1 {
		t.Errorf("負のレベルは1にクランプされるべき: got %d", enemy.Level)
	}
}

// TestEnemyGeneration_StatsScaling はステータススケーリングをテストします。
func TestEnemyGeneration_StatsScaling(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}
	generator := NewEnemyGenerator(enemyTypes)

	// レベルが上がるとステータスも上がることを確認
	previousHP := 0
	previousAttack := 0

	for level := 1; level <= 50; level += 10 {
		enemy := generator.Generate(level)

		if enemy.MaxHP <= previousHP && level > 1 {
			t.Errorf("レベル%dでHPが前レベルより低い: got %d, prev %d", level, enemy.MaxHP, previousHP)
		}
		if enemy.AttackPower <= previousAttack && level > 1 {
			t.Errorf("レベル%dで攻撃力が前レベルより低い: got %d, prev %d", level, enemy.AttackPower, previousAttack)
		}

		previousHP = enemy.MaxHP
		previousAttack = enemy.AttackPower
	}
}

// TestEnemyGenerator_GenerateWithType は特定タイプで敵を生成するテストです。
func TestEnemyGenerator_GenerateWithType(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "dragon",
			Name:   "ドラゴン",
			BaseHP: 500,
			Rank:   1,
		},
	}

	gen := NewEnemyGenerator(enemyTypes)
	gen.SetSeed(42) // 再現可能な結果のために固定シード

	enemy := gen.GenerateWithType(10, "dragon")
	if enemy == nil {
		t.Fatal("GenerateWithType returned nil")
	}

	// ドラゴンタイプの名前が含まれていること
	if enemy.Type.ID != "dragon" {
		t.Errorf("Expected dragon type, got %s", enemy.Type.ID)
	}

	// HP = BaseHP * level = 500 * 10 = 5000
	if enemy.HP != 5000 {
		t.Errorf("Expected HP 5000, got %d", enemy.HP)
	}
}

// TestEnemyGenerator_GetEnemyTypes は敵タイプ取得をテストします。
func TestEnemyGenerator_GetEnemyTypes(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "type1", Name: "タイプ1"},
		{ID: "type2", Name: "タイプ2"},
	}

	gen := NewEnemyGenerator(enemyTypes)

	types := gen.GetEnemyTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 enemy types, got %d", len(types))
	}
}

// ==================== Task 6.1: 敵生成フロー統合テスト ====================

// TestEnemyGenerator_GenerateWithType_ActionPatternIntegration は行動パターン付き敵タイプからの生成をテストします。
func TestEnemyGenerator_GenerateWithType_ActionPatternIntegration(t *testing.T) {
	// 行動パターンとパッシブを持つ敵タイプを定義
	normalActions := []domain.EnemyAction{
		{
			ID:             "act_attack",
			Name:           "攻撃",
			ActionType:     domain.EnemyActionAttack,
			AttackType:     "physical",
			DamageBase:     5.0,
			DamagePerLevel: 1.0,
			ChargeTime:     2 * time.Second,
		},
	}
	enhancedActions := []domain.EnemyAction{
		{
			ID:             "act_attack_enhanced",
			Name:           "強化攻撃",
			ActionType:     domain.EnemyActionAttack,
			AttackType:     "physical",
			DamageBase:     10.0,
			DamagePerLevel: 2.0,
			ChargeTime:     3 * time.Second,
		},
	}
	normalPassive := &domain.EnemyPassiveSkill{
		ID:          "passive_normal",
		Name:        "通常パッシブ",
		Description: "通常時のパッシブ効果",
		Effects: map[domain.EffectColumn]float64{
			domain.ColDamageMultiplier: 1.1,
		},
	}

	enemyTypes := []domain.EnemyType{
		{
			ID:                      "test_enemy",
			Name:                    "テスト敵",
			BaseHP:                  100,
			Rank:                    1,
			ResolvedNormalActions:   normalActions,
			ResolvedEnhancedActions: enhancedActions,
			NormalPassive:           normalPassive,
			DropItemCategory:        "core",
			DropItemTypeID:          "attack_balance",
		},
	}

	gen := NewEnemyGenerator(enemyTypes)

	// 敵タイプ指定で生成
	enemy := gen.GenerateWithType(5, "test_enemy")

	if enemy == nil {
		t.Fatal("敵の生成に失敗")
	}

	// ActionIndexが0で初期化されていること
	if enemy.ActionIndex != 0 {
		t.Errorf("ActionIndexは0で初期化されるべき: got %d", enemy.ActionIndex)
	}

	// 行動パターンが正しくコピーされていること
	if len(enemy.Type.ResolvedNormalActions) != 1 {
		t.Errorf("通常行動パターンが正しくコピーされていない: got %d, want 1", len(enemy.Type.ResolvedNormalActions))
	}
	if len(enemy.Type.ResolvedEnhancedActions) != 1 {
		t.Errorf("強化行動パターンが正しくコピーされていない: got %d, want 1", len(enemy.Type.ResolvedEnhancedActions))
	}

	// パッシブスキルが正しくコピーされていること
	if enemy.Type.NormalPassive == nil {
		t.Error("通常パッシブがコピーされていない")
	} else if enemy.Type.NormalPassive.ID != "passive_normal" {
		t.Errorf("通常パッシブIDが不正: got %s", enemy.Type.NormalPassive.ID)
	}

	// GetCurrentActionで最初の行動が取得できること
	action := enemy.GetCurrentAction()
	if action.ID != "act_attack" {
		t.Errorf("最初の行動が不正: got %s, want act_attack", action.ID)
	}
}

// TestEnemyGenerator_GenerateWithType_PassedFromBattleSelect はバトル選択画面から渡されるEnemyTypeIDの受け渡しをテストします。
func TestEnemyGenerator_GenerateWithType_PassedFromBattleSelect(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
		{
			ID:     "dragon",
			Name:   "ドラゴン",
			BaseHP: 200,
			Rank:   1,
		},
	}

	gen := NewEnemyGenerator(enemyTypes)

	// dragonを指定して生成
	enemy := gen.GenerateWithType(10, "dragon")

	if enemy == nil {
		t.Fatal("敵の生成に失敗")
	}

	// 正しいタイプが選択されていること
	if enemy.Type.ID != "dragon" {
		t.Errorf("指定した敵タイプが生成されていない: got %s, want dragon", enemy.Type.ID)
	}

	// レベルが正しく設定されていること
	if enemy.Level != 10 {
		t.Errorf("レベルが正しく設定されていない: got %d, want 10", enemy.Level)
	}

	// 名前にレベルが含まれていること
	expectedName := "ドラゴン Lv.10"
	if enemy.Name != expectedName {
		t.Errorf("名前が不正: got %s, want %s", enemy.Name, expectedName)
	}
}

// TestEnemyGenerator_GenerateWithType_InvalidID は存在しないTypeIDの場合のフォールバックをテストします。
func TestEnemyGenerator_GenerateWithType_InvalidID(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:     "slime",
			Name:   "スライム",
			BaseHP: 50,
			Rank:   1,
		},
	}

	gen := NewEnemyGenerator(enemyTypes)
	gen.SetSeed(42) // ランダム選択を固定

	// 存在しないIDを指定
	enemy := gen.GenerateWithType(5, "invalid_id")

	if enemy == nil {
		t.Fatal("存在しないIDでも敵が生成されるべき（フォールバック）")
	}

	// 既存の敵タイプからランダムに選択されること
	if enemy.Type.ID != "slime" {
		t.Errorf("フォールバックでslimeが選択されるべき: got %s", enemy.Type.ID)
	}
}
