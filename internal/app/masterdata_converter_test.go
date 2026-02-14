package app

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/infra/masterdata"
	"hirorocky/type-battle/internal/usecase/rewarding"
)

// TestResolveSkillTimedEffects はスキルのEffectColumnSpecにTimedEffectからColumn/Valueが解決されることをテストします。
func TestResolveSkillTimedEffects(t *testing.T) {
	// Arrange: TimedEffectマップを作成
	timedEffects := map[string]domain.TimedEffect{
		"st_str_buff_lv1": {
			ID:     "st_str_buff_lv1",
			Name:   "STR25%UP",
			Column: domain.ColSTRMultiplier,
			Value:  0.25,
		},
		"st_defense_buff_lv1": {
			ID:     "st_defense_buff_lv1",
			Name:   "被ダメ10%軽減",
			Column: domain.ColDamageCut,
			Value:  0.1,
		},
	}

	// Arrange: ColumnSpecにTimedEffectIDのみ設定されたスキル
	skills := []rewarding.SkillDropInfo{
		{
			ID:   "str_buff_lv1",
			Name: "気合い溜め",
			Effects: []domain.SkillEffect{
				{
					// HP回復効果（ColumnSpecなし）
					Target: domain.TargetSelf,
					HPFormula: &domain.HPFormula{
						Base:     10,
						StatCoef: 0.1,
						StatRef:  "WIL",
					},
				},
				{
					// バフ効果（Column/ValueはTimedEffectから解決される）
					Target: domain.TargetSelf,
					ColumnSpec: &domain.EffectColumnSpec{
						Duration:      10.0,
						TimedEffectID: "st_str_buff_lv1",
					},
				},
			},
		},
		{
			ID:   "defense_buff_lv1",
			Name: "防御強化1",
			Effects: []domain.SkillEffect{
				{
					Target: domain.TargetSelf,
					ColumnSpec: &domain.EffectColumnSpec{
						Duration:      10.0,
						TimedEffectID: "st_defense_buff_lv1",
					},
				},
			},
		},
	}

	// Act
	ResolveSkillTimedEffects(skills, timedEffects)

	// Assert: str_buff_lv1のバフ効果が解決されていること
	strBuff := skills[0].Effects[1]
	if strBuff.ColumnSpec.Column != domain.ColSTRMultiplier {
		t.Errorf("Column: got %s, want %s", strBuff.ColumnSpec.Column, domain.ColSTRMultiplier)
	}
	if strBuff.ColumnSpec.Value != 0.25 {
		t.Errorf("Value: got %f, want 0.25", strBuff.ColumnSpec.Value)
	}
	if strBuff.ColumnSpec.Duration != 10.0 {
		t.Errorf("Duration: got %f, want 10.0", strBuff.ColumnSpec.Duration)
	}

	// Assert: HP回復効果は影響なし
	if skills[0].Effects[0].ColumnSpec != nil {
		t.Error("HP回復効果のColumnSpecはnilのままであるべき")
	}

	// Assert: defense_buff_lv1の効果が解決されていること
	defBuff := skills[1].Effects[0]
	if defBuff.ColumnSpec.Column != domain.ColDamageCut {
		t.Errorf("Column: got %s, want %s", defBuff.ColumnSpec.Column, domain.ColDamageCut)
	}
	if defBuff.ColumnSpec.Value != 0.1 {
		t.Errorf("Value: got %f, want 0.1", defBuff.ColumnSpec.Value)
	}
}

// TestResolveEnemyActionTimedEffects は敵行動のEffectColumn/EffectValueがTimedEffectから解決されることをテストします（AC13）。
func TestResolveEnemyActionTimedEffects(t *testing.T) {
	// Arrange: TimedEffectマップ
	timedEffects := map[string]domain.TimedEffect{
		"st_goblin_rage": {
			ID:     "st_goblin_rage",
			Name:   "怒り",
			Column: domain.ColDamageMultiplier,
			Value:  1.3,
		},
		"st_skeleton_bone_dust": {
			ID:     "st_skeleton_bone_dust",
			Name:   "骨粉",
			Column: domain.ColDamageCut,
			Value:  -0.15,
		},
	}

	// Arrange: 敵タイプに解決済みアクションを設定
	enemyTypes := []domain.EnemyType{
		{
			ID:   "goblin",
			Name: "ゴブリン",
		},
	}
	enemyTypes[0].SetResolvedActions(
		[]domain.EnemyAction{
			{
				ID:            "act_goblin_attack",
				ActionType:    domain.EnemyActionAttack,
				TimedEffectID: "", // 攻撃行動にはTimedEffectIDなし
			},
			{
				ID:            "act_goblin_buff",
				ActionType:    domain.EnemyActionBuff,
				TimedEffectID: "st_goblin_rage",
				Duration:      10.0,
			},
		},
		[]domain.EnemyAction{
			{
				ID:            "act_skeleton_debuff",
				ActionType:    domain.EnemyActionDebuff,
				TimedEffectID: "st_skeleton_bone_dust",
				Duration:      6.0,
			},
		},
	)

	// Act
	ResolveEnemyActionTimedEffects(enemyTypes, timedEffects)

	// Assert: 攻撃行動は影響なし
	attackAction := enemyTypes[0].GetNormalActions()[0]
	if attackAction.EffectColumn != "" {
		t.Errorf("攻撃行動のEffectColumnは空であるべき: got %s", attackAction.EffectColumn)
	}

	// Assert: バフ行動のEffectColumn/EffectValueが解決されていること
	buffAction := enemyTypes[0].GetNormalActions()[1]
	if buffAction.EffectColumn != domain.ColDamageMultiplier {
		t.Errorf("EffectColumn: got %s, want %s", buffAction.EffectColumn, domain.ColDamageMultiplier)
	}
	if buffAction.EffectValue != 1.3 {
		t.Errorf("EffectValue: got %f, want 1.3", buffAction.EffectValue)
	}

	// Assert: デバフ行動のEffectColumn/EffectValueが解決されていること
	debuffAction := enemyTypes[0].GetEnhancedActions()[0]
	if debuffAction.EffectColumn != domain.ColDamageCut {
		t.Errorf("EffectColumn: got %s, want %s", debuffAction.EffectColumn, domain.ColDamageCut)
	}
	if debuffAction.EffectValue != -0.15 {
		t.Errorf("EffectValue: got %f, want -0.15", debuffAction.EffectValue)
	}
}

// TestConvertTimedEffectsWithColumnValue はTimedEffectDataからColumn/Value付きのドメインモデルへの変換をテストします（AC10）。
func TestConvertTimedEffectsWithColumnValue(t *testing.T) {
	// Arrange
	timedEffectData := []masterdata.TimedEffectData{
		{
			ID:           "st_str_buff_lv1",
			Name:         "STR25%UP",
			Description:  "STRを25%上昇させる",
			EffectColumn: "str_mult",
			Value:        0.25,
		},
		{
			ID:           "st_goblin_rage",
			Name:         "怒り",
			Description:  "ゴブリンの攻撃力が上昇する",
			EffectColumn: "damage_mult",
			Value:        1.3,
		},
	}

	// Act
	result := ConvertTimedEffects(timedEffectData)

	// Assert
	if len(result) != 2 {
		t.Fatalf("結果の数: got %d, want 2", len(result))
	}

	strBuff, ok := result["st_str_buff_lv1"]
	if !ok {
		t.Fatal("st_str_buff_lv1が結果に含まれていない")
	}
	if strBuff.Column != domain.ColSTRMultiplier {
		t.Errorf("Column: got %s, want %s", strBuff.Column, domain.ColSTRMultiplier)
	}
	if strBuff.Value != 0.25 {
		t.Errorf("Value: got %f, want 0.25", strBuff.Value)
	}

	rage, ok := result["st_goblin_rage"]
	if !ok {
		t.Fatal("st_goblin_rageが結果に含まれていない")
	}
	if rage.Column != domain.ColDamageMultiplier {
		t.Errorf("Column: got %s, want %s", rage.Column, domain.ColDamageMultiplier)
	}
	if rage.Value != 1.3 {
		t.Errorf("Value: got %f, want 1.3", rage.Value)
	}
}
