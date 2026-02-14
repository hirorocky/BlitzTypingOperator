// Package domain はゲームのドメインモデルを定義します。
package domain

import (
	"encoding/json"
	"testing"
	"time"
)

// TestEnemyPhase_定数の確認 はEnemyPhase定数が正しく定義されていることを確認します。
func TestEnemyPhase_定数の確認(t *testing.T) {
	if PhaseNormal != 0 {
		t.Errorf("PhaseNormalが期待値と異なります: got %d, want 0", PhaseNormal)
	}
	if PhaseEnhanced != 1 {
		t.Errorf("PhaseEnhancedが期待値と異なります: got %d, want 1", PhaseEnhanced)
	}
}

// TestEnemyPhase_String はEnemyPhaseのString()メソッドが正しい表示名を返すことを確認します。
func TestEnemyPhase_String(t *testing.T) {
	if PhaseNormal.String() != "通常" {
		t.Errorf("PhaseNormal.String()が期待値と異なります: got %s, want 通常", PhaseNormal.String())
	}
	if PhaseEnhanced.String() != "強化" {
		t.Errorf("PhaseEnhanced.String()が期待値と異なります: got %s, want 強化", PhaseEnhanced.String())
	}
}

// TestEnemyType_フィールドの確認 はEnemyType構造体のフィールドが正しく設定されることを確認します。
func TestEnemyType_フィールドの確認(t *testing.T) {
	enemyType := EnemyType{
		ID:     "goblin",
		Name:   "ゴブリン",
		BaseHP: 100,
		Rank:   1,
	}

	if enemyType.ID != "goblin" {
		t.Errorf("IDが期待値と異なります: got %s, want goblin", enemyType.ID)
	}
	if enemyType.Name != "ゴブリン" {
		t.Errorf("Nameが期待値と異なります: got %s, want ゴブリン", enemyType.Name)
	}
	if enemyType.BaseHP != 100 {
		t.Errorf("BaseHPが期待値と異なります: got %d, want 100", enemyType.BaseHP)
	}
	if enemyType.Rank != 1 {
		t.Errorf("Rankが期待値と異なります: got %d, want 1", enemyType.Rank)
	}
}

// TestEnemyType_ランクシステム は敵ランクの検証機能をテストします。
func TestEnemyType_ランクシステム(t *testing.T) {
	tests := []struct {
		name      string
		rank      int
		wantValid bool
	}{
		{
			name:      "ランク1は有効",
			rank:      1,
			wantValid: true,
		},
		{
			name:      "ランク5は有効",
			rank:      5,
			wantValid: true,
		},
		{
			name:      "ランク0は無効",
			rank:      0,
			wantValid: false,
		},
		{
			name:      "負のランクは無効",
			rank:      -1,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemyType := EnemyType{
				ID:   "test_enemy",
				Name: "テスト敵",
				Rank: tt.rank,
			}

			if got := enemyType.IsValidRank(); got != tt.wantValid {
				t.Errorf("IsValidRank() = %v, want %v", got, tt.wantValid)
			}
		})
	}
}

// TestEnemyModel_フィールドの確認 はEnemyModel構造体のフィールドが正しく設定されることを確認します。

func TestEnemyModel_フィールドの確認(t *testing.T) {
	enemyType := EnemyType{
		ID:     "goblin",
		Name:   "ゴブリン",
		BaseHP: 100,
	}

	enemy := EnemyModel{
		ID:          "enemy_001",
		Name:        "ゴブリン兵士",
		Level:       5,
		HP:          150,
		MaxHP:       150,
		AttackPower: 15,
		Type:        enemyType,
		Phase:       PhaseNormal,
		EffectTable: NewEffectTable(),
	}

	if enemy.ID != "enemy_001" {
		t.Errorf("IDが期待値と異なります: got %s, want enemy_001", enemy.ID)
	}
	if enemy.Name != "ゴブリン兵士" {
		t.Errorf("Nameが期待値と異なります: got %s, want ゴブリン兵士", enemy.Name)
	}
	if enemy.Level != 5 {
		t.Errorf("Levelが期待値と異なります: got %d, want 5", enemy.Level)
	}
	if enemy.HP != 150 {
		t.Errorf("HPが期待値と異なります: got %d, want 150", enemy.HP)
	}
	if enemy.MaxHP != 150 {
		t.Errorf("MaxHPが期待値と異なります: got %d, want 150", enemy.MaxHP)
	}
	if enemy.AttackPower != 15 {
		t.Errorf("AttackPowerが期待値と異なります: got %d, want 15", enemy.AttackPower)
	}
	if enemy.Phase != PhaseNormal {
		t.Errorf("Phaseが期待値と異なります: got %d, want PhaseNormal", enemy.Phase)
	}
	if enemy.EffectTable == nil {
		t.Error("EffectTableがnilです")
	}
}

// TestNewEnemy_敵作成 はNewEnemy関数で敵が正しく作成されることを確認します。
func TestNewEnemy_敵作成(t *testing.T) {
	enemyType := EnemyType{
		ID:     "goblin",
		Name:   "ゴブリン",
		BaseHP: 100,
	}

	enemy := NewEnemy("enemy_001", "ゴブリン兵士", 5, 150, 15, enemyType)

	if enemy.ID != "enemy_001" {
		t.Errorf("IDが期待値と異なります: got %s, want enemy_001", enemy.ID)
	}
	if enemy.Phase != PhaseNormal {
		t.Error("初期状態はPhaseNormalであるべきです")
	}
	if enemy.EffectTable == nil {
		t.Error("EffectTableが初期化されていません")
	}
}

// TestEnemyModel_HP50以下でフェーズ変化 は敵のHPが50%以下でフェーズ変化するルールを確認します。

func TestEnemyModel_HP50以下でフェーズ変化(t *testing.T) {
	tests := []struct {
		name          string
		maxHP         int
		currentHP     int
		shouldEnhance bool
	}{
		{"HP100% (100/100)", 100, 100, false},
		{"HP60% (60/100)", 100, 60, false},
		{"HP51% (51/100)", 100, 51, false},
		{"HP50% (50/100)", 100, 50, true}, // 50%以下で強化
		{"HP49% (49/100)", 100, 49, true},
		{"HP10% (10/100)", 100, 10, true},
		{"HP0% (0/100)", 100, 0, true},
		{"HP50% (25/50)", 50, 25, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := EnemyModel{
				HP:    tt.currentHP,
				MaxHP: tt.maxHP,
				Phase: PhaseNormal,
			}

			shouldTransition := enemy.ShouldTransitionToEnhanced()

			if shouldTransition != tt.shouldEnhance {
				t.Errorf("ShouldTransitionToEnhancedの結果が期待値と異なります: got %v, want %v",
					shouldTransition, tt.shouldEnhance)
			}
		})
	}
}

// TestEnemyModel_フェーズ移行 は敵のフェーズ移行を確認します。

func TestEnemyModel_フェーズ移行(t *testing.T) {
	enemy := EnemyModel{
		HP:    50,
		MaxHP: 100,
		Phase: PhaseNormal,
	}

	// 移行前
	if enemy.Phase != PhaseNormal {
		t.Error("初期状態はPhaseNormalであるべきです")
	}

	// フェーズ移行
	enemy.TransitionToEnhanced()

	// 移行後
	if enemy.Phase != PhaseEnhanced {
		t.Error("移行後はPhaseEnhancedであるべきです")
	}
}

// TestEnemyModel_フェーズ移行は1回のみ はフェーズ移行が2回行われないことを確認します。
func TestEnemyModel_フェーズ移行は1回のみ(t *testing.T) {
	enemy := EnemyModel{
		HP:    30,
		MaxHP: 100,
		Phase: PhaseEnhanced, // 既に強化フェーズ
	}

	// 強化フェーズ中は再移行しない
	if enemy.ShouldTransitionToEnhanced() {
		t.Error("既にPhaseEnhancedの場合、ShouldTransitionToEnhancedはfalseを返すべきです")
	}
}

// TestEnemyModel_ダメージ受け はダメージを受けてHPが減少することを確認します。
func TestEnemyModel_ダメージ受け(t *testing.T) {
	enemy := EnemyModel{
		HP:    100,
		MaxHP: 100,
		Phase: PhaseNormal,
	}

	// ダメージを受ける
	enemy.TakeDamage(30)
	if enemy.HP != 70 {
		t.Errorf("ダメージ後のHPが期待値と異なります: got %d, want 70", enemy.HP)
	}

	// 致死ダメージ（HPは0以下にならない）
	enemy.TakeDamage(100)
	if enemy.HP != 0 {
		t.Errorf("HPが0未満になっています: got %d, want 0", enemy.HP)
	}
}

// TestEnemyModel_生存確認 は敵の生存確認を確認します。
func TestEnemyModel_生存確認(t *testing.T) {
	enemy := EnemyModel{
		HP:    100,
		MaxHP: 100,
	}

	// 生存状態
	if !enemy.IsAlive() {
		t.Error("HPが0より大きい場合は生存しているはずです")
	}

	// 死亡状態
	enemy.HP = 0
	if enemy.IsAlive() {
		t.Error("HP=0の場合は死亡しているはずです")
	}
}

// TestEnemyModel_HP割合取得 はHP割合の取得を確認します。
func TestEnemyModel_HP割合取得(t *testing.T) {
	tests := []struct {
		name     string
		hp       int
		maxHP    int
		expected float64
	}{
		{"100%", 100, 100, 1.0},
		{"50%", 50, 100, 0.5},
		{"0%", 0, 100, 0.0},
		{"75%", 75, 100, 0.75},
		{"MaxHP=0", 0, 0, 0.0}, // ゼロ除算対応
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := EnemyModel{HP: tt.hp, MaxHP: tt.maxHP}
			percentage := enemy.GetHPPercentage()

			if percentage != tt.expected {
				t.Errorf("HP割合が期待値と異なります: got %f, want %f", percentage, tt.expected)
			}
		})
	}
}

// TestEnemyModel_強化フェーズ判定 は現在強化フェーズかどうかを確認します。
func TestEnemyModel_強化フェーズ判定(t *testing.T) {
	enemy := EnemyModel{Phase: PhaseNormal}
	if enemy.IsEnhanced() {
		t.Error("PhaseNormalではIsEnhanced()はfalseを返すべきです")
	}

	enemy.Phase = PhaseEnhanced
	if !enemy.IsEnhanced() {
		t.Error("PhaseEnhancedではIsEnhanced()はtrueを返すべきです")
	}
}

// TestEnemyModel_EffectTable操作 は敵のEffectTableを操作できることを確認します。
func TestEnemyModel_EffectTable操作(t *testing.T) {
	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, EnemyType{ID: "test"})

	// バフを追加
	enemy.EffectTable.AddBuff("攻撃力UP", "攻撃力UP", 5.0, map[EffectColumn]float64{
		ColDamageBonus: 10,
	})

	if len(enemy.EffectTable.Entries) != 1 {
		t.Errorf("EffectTableのエントリ数が期待値と異なります: got %d, want 1", len(enemy.EffectTable.Entries))
	}
}

// TestEnemyModel_CheckAndTransitionPhase はHP変化後のフェーズ移行チェックを確認します。
func TestEnemyModel_CheckAndTransitionPhase(t *testing.T) {
	enemy := EnemyModel{
		HP:    60,
		MaxHP: 100,
		Phase: PhaseNormal,
	}

	// まだ移行しない
	transitioned := enemy.CheckAndTransitionPhase()
	if transitioned {
		t.Error("HP60%ではフェーズ移行しないはずです")
	}
	if enemy.Phase != PhaseNormal {
		t.Error("フェーズがまだNormalであるべきです")
	}

	// HP減少
	enemy.HP = 50

	// フェーズ移行
	transitioned = enemy.CheckAndTransitionPhase()
	if !transitioned {
		t.Error("HP50%ではフェーズ移行するはずです")
	}
	if enemy.Phase != PhaseEnhanced {
		t.Error("フェーズがEnhancedに変わるべきです")
	}
}

// TestEnhanceThreshold は強化フェーズ移行の閾値が正しい値であることを確認します。
func TestEnhanceThreshold(t *testing.T) {
	if EnhanceThreshold != 0.5 {
		t.Errorf("EnhanceThresholdが期待値と異なります: got %f, want 0.5", EnhanceThreshold)
	}
}

// ========== タスク1.1: 敵の行動データ構造のテスト ==========

// TestEnemyActionType_定数の確認 はEnemyActionType定数が正しく定義されていることを確認します。
func TestEnemyActionType_定数の確認(t *testing.T) {
	// 行動タイプは Attack, Buff, Debuff, Defense の4種類
	if EnemyActionAttack != 0 {
		t.Errorf("EnemyActionAttackが期待値と異なります: got %d, want 0", EnemyActionAttack)
	}
	if EnemyActionBuff != 1 {
		t.Errorf("EnemyActionBuffが期待値と異なります: got %d, want 1", EnemyActionBuff)
	}
	if EnemyActionDebuff != 2 {
		t.Errorf("EnemyActionDebuffが期待値と異なります: got %d, want 2", EnemyActionDebuff)
	}
	if EnemyActionDefense != 3 {
		t.Errorf("EnemyActionDefenseが期待値と異なります: got %d, want 3", EnemyActionDefense)
	}
}

// TestEnemyActionType_String はEnemyActionTypeのString()メソッドが正しい表示名を返すことを確認します。
func TestEnemyActionType_String(t *testing.T) {
	tests := []struct {
		actionType EnemyActionType
		expected   string
	}{
		{EnemyActionAttack, "攻撃"},
		{EnemyActionBuff, "バフ"},
		{EnemyActionDebuff, "デバフ"},
		{EnemyActionDefense, "ディフェンス"},
		{EnemyActionType(99), "不明"}, // 未定義値
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.actionType.String() != tt.expected {
				t.Errorf("String()が期待値と異なります: got %s, want %s", tt.actionType.String(), tt.expected)
			}
		})
	}
}

// TestEnemyAction_攻撃行動のフィールド は攻撃行動のフィールドが正しく設定されることを確認します。
func TestEnemyAction_攻撃行動のフィールド(t *testing.T) {
	// 物理攻撃
	physicalAttack := EnemyAction{
		ActionType: EnemyActionAttack,
		AttackType: "physical",
	}
	if physicalAttack.ActionType != EnemyActionAttack {
		t.Error("ActionTypeがAttackであるべきです")
	}
	if physicalAttack.AttackType != "physical" {
		t.Errorf("AttackTypeが期待値と異なります: got %s, want physical", physicalAttack.AttackType)
	}

	// 魔法攻撃
	magicAttack := EnemyAction{
		ActionType: EnemyActionAttack,
		AttackType: "magic",
	}
	if magicAttack.AttackType != "magic" {
		t.Errorf("AttackTypeが期待値と異なります: got %s, want magic", magicAttack.AttackType)
	}
}

// TestEnemyAction_バフ行動のフィールド は自己バフ行動のフィールドが正しく設定されることを確認します。
func TestEnemyAction_バフ行動のフィールド(t *testing.T) {
	buff := EnemyAction{
		ActionType:   EnemyActionBuff,
		EffectColumn: EffectColumn("attackUp"),
		EffectValue:  0.3,
		Duration:     10.0,
	}

	if buff.ActionType != EnemyActionBuff {
		t.Error("ActionTypeがSelfBuffであるべきです")
	}
	if buff.EffectColumn != EffectColumn("attackUp") {
		t.Errorf("EffectColumnが期待値と異なります: got %s, want attackUp", buff.EffectColumn)
	}
	if buff.EffectValue != 0.3 {
		t.Errorf("EffectValueが期待値と異なります: got %f, want 0.3", buff.EffectValue)
	}
	if buff.Duration != 10.0 {
		t.Errorf("Durationが期待値と異なります: got %f, want 10.0", buff.Duration)
	}
}

// TestEnemyAction_デバフ行動のフィールド はデバフ行動のフィールドが正しく設定されることを確認します。
func TestEnemyAction_デバフ行動のフィールド(t *testing.T) {
	debuff := EnemyAction{
		ActionType:   EnemyActionDebuff,
		EffectColumn: EffectColumn("defenseDown"),
		EffectValue:  0.2,
		Duration:     5.0,
	}

	if debuff.ActionType != EnemyActionDebuff {
		t.Error("ActionTypeがDebuffであるべきです")
	}
	if debuff.EffectColumn != EffectColumn("defenseDown") {
		t.Errorf("EffectColumnが期待値と異なります: got %s, want defenseDown", debuff.EffectColumn)
	}
	if debuff.EffectValue != 0.2 {
		t.Errorf("EffectValueが期待値と異なります: got %f, want 0.2", debuff.EffectValue)
	}
	if debuff.Duration != 5.0 {
		t.Errorf("Durationが期待値と異なります: got %f, want 5.0", debuff.Duration)
	}
}

// TestEnemyAction_IsAttack は攻撃行動かどうかを判定するヘルパーメソッドを確認します。
func TestEnemyAction_IsAttack(t *testing.T) {
	attack := EnemyAction{ActionType: EnemyActionAttack}
	buff := EnemyAction{ActionType: EnemyActionBuff}
	debuff := EnemyAction{ActionType: EnemyActionDebuff}

	if !attack.IsAttack() {
		t.Error("Attack行動でIsAttack()がtrueを返すべきです")
	}
	if buff.IsAttack() {
		t.Error("SelfBuff行動でIsAttack()がfalseを返すべきです")
	}
	if debuff.IsAttack() {
		t.Error("Debuff行動でIsAttack()がfalseを返すべきです")
	}
}

// TestEnemyAction_IsBuff はバフ行動かどうかを判定するヘルパーメソッドを確認します。
func TestEnemyAction_IsBuff(t *testing.T) {
	attack := EnemyAction{ActionType: EnemyActionAttack}
	buff := EnemyAction{ActionType: EnemyActionBuff}
	debuff := EnemyAction{ActionType: EnemyActionDebuff}

	if attack.IsBuff() {
		t.Error("Attack行動でIsBuff()がfalseを返すべきです")
	}
	if !buff.IsBuff() {
		t.Error("SelfBuff行動でIsBuff()がtrueを返すべきです")
	}
	if debuff.IsBuff() {
		t.Error("Debuff行動でIsBuff()がfalseを返すべきです")
	}
}

// TestEnemyAction_IsDebuff はデバフ行動かどうかを判定するヘルパーメソッドを確認します。
func TestEnemyAction_IsDebuff(t *testing.T) {
	attack := EnemyAction{ActionType: EnemyActionAttack}
	buff := EnemyAction{ActionType: EnemyActionBuff}
	debuff := EnemyAction{ActionType: EnemyActionDebuff}

	if attack.IsDebuff() {
		t.Error("Attack行動でIsDebuff()がfalseを返すべきです")
	}
	if buff.IsDebuff() {
		t.Error("SelfBuff行動でIsDebuff()がfalseを返すべきです")
	}
	if !debuff.IsDebuff() {
		t.Error("Debuff行動でIsDebuff()がtrueを返すべきです")
	}
}

// ========== タスク1.2: 敵パッシブスキルのデータ構造のテスト ==========

// TestEnemyPassiveSkill_フィールドの確認 はEnemyPassiveSkill構造体のフィールドが正しく設定されることを確認します。
func TestEnemyPassiveSkill_フィールドの確認(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:          "slime_normal",
		Name:        "ぷるぷるボディ",
		Description: "物理ダメージを10%軽減",
		Effects: map[EffectColumn]float64{
			ColDamageCut: 0.1,
		},
	}

	if passive.ID != "slime_normal" {
		t.Errorf("IDが期待値と異なります: got %s, want slime_normal", passive.ID)
	}
	if passive.Name != "ぷるぷるボディ" {
		t.Errorf("Nameが期待値と異なります: got %s, want ぷるぷるボディ", passive.Name)
	}
	expectedDesc := "物理ダメージを10%軽減"
	if passive.Description != expectedDesc {
		t.Errorf("Descriptionが期待値と異なります: got %s, want %s", passive.Description, expectedDesc)
	}
	if len(passive.Effects) != 1 {
		t.Errorf("Effectsの要素数が期待値と異なります: got %d, want 1", len(passive.Effects))
	}
	if passive.Effects[ColDamageCut] != 0.1 {
		t.Errorf("Effects[ColDamageCut]が期待値と異なります: got %f, want 0.1", passive.Effects[ColDamageCut])
	}
}

// TestEnemyPassiveSkill_複数効果 は複数の効果を持つパッシブスキルを確認します。
func TestEnemyPassiveSkill_複数効果(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:          "boss_enhanced",
		Name:        "狂戦士の怒り",
		Description: "攻撃力20%上昇、ライフスティール10%",
		Effects: map[EffectColumn]float64{
			ColDamageMultiplier: 1.2,
			ColLifeSteal:        0.1,
		},
	}

	if len(passive.Effects) != 2 {
		t.Errorf("Effectsの要素数が期待値と異なります: got %d, want 2", len(passive.Effects))
	}
	if passive.Effects[ColDamageMultiplier] != 1.2 {
		t.Errorf("Effects[ColDamageMultiplier]が期待値と異なります: got %f, want 1.2", passive.Effects[ColDamageMultiplier])
	}
	if passive.Effects[ColLifeSteal] != 0.1 {
		t.Errorf("Effects[ColLifeSteal]が期待値と異なります: got %f, want 0.1", passive.Effects[ColLifeSteal])
	}
}

// TestEnemyPassiveSkill_ToEntry はEffectEntryへの変換を確認します。
func TestEnemyPassiveSkill_ToEntry(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:          "slime_normal",
		Name:        "ぷるぷるボディ",
		Description: "物理ダメージを10%軽減",
		Effects: map[EffectColumn]float64{
			ColDamageCut: 0.1,
		},
	}

	entry := passive.ToEntry()

	// ソースタイプはパッシブであること
	if entry.SourceType != SourcePassive {
		t.Errorf("SourceTypeが期待値と異なります: got %s, want %s", entry.SourceType, SourcePassive)
	}
	// ソースIDはパッシブのIDと一致すること
	if entry.SourceID != passive.ID {
		t.Errorf("SourceIDが期待値と異なります: got %s, want %s", entry.SourceID, passive.ID)
	}
	// 表示名はパッシブの名前と一致すること
	if entry.Name != passive.Name {
		t.Errorf("Nameが期待値と異なります: got %s, want %s", entry.Name, passive.Name)
	}
	// 永続効果であること（Durationがnil）
	if entry.Duration != nil {
		t.Error("パッシブスキルは永続効果（Duration=nil）であるべきです")
	}
	// 効果値が正しく変換されていること
	if entry.Values[ColDamageCut] != 0.1 {
		t.Errorf("Values[ColDamageCut]が期待値と異なります: got %f, want 0.1", entry.Values[ColDamageCut])
	}
}

// TestEnemyPassiveSkill_ToEntry_複数効果 は複数効果のEffectEntry変換を確認します。
func TestEnemyPassiveSkill_ToEntry_複数効果(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:   "boss_enhanced",
		Name: "狂戦士の怒り",
		Effects: map[EffectColumn]float64{
			ColDamageMultiplier: 1.2,
			ColLifeSteal:        0.1,
		},
	}

	entry := passive.ToEntry()

	if len(entry.Values) != 2 {
		t.Errorf("Valuesの要素数が期待値と異なります: got %d, want 2", len(entry.Values))
	}
	if entry.Values[ColDamageMultiplier] != 1.2 {
		t.Errorf("Values[ColDamageMultiplier]が期待値と異なります: got %f, want 1.2", entry.Values[ColDamageMultiplier])
	}
	if entry.Values[ColLifeSteal] != 0.1 {
		t.Errorf("Values[ColLifeSteal]が期待値と異なります: got %f, want 0.1", entry.Values[ColLifeSteal])
	}
}

// TestEnemyPassiveSkill_EffectTableとの連携 はEffectTableにパッシブスキルを登録できることを確認します。
func TestEnemyPassiveSkill_EffectTableとの連携(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:   "slime_normal",
		Name: "ぷるぷるボディ",
		Effects: map[EffectColumn]float64{
			ColDamageCut: 0.1,
		},
	}

	table := NewEffectTable()
	table.AddEntry(passive.ToEntry())

	// エントリが追加されていること
	if len(table.Entries) != 1 {
		t.Errorf("エントリ数が期待値と異なります: got %d, want 1", len(table.Entries))
	}

	// パッシブを検索できること
	passives := table.FindBySourceType(SourcePassive)
	if len(passives) != 1 {
		t.Errorf("パッシブスキル数が期待値と異なります: got %d, want 1", len(passives))
	}

	// SourceIDで検索できること
	found := table.FindBySourceID("slime_normal")
	if found == nil {
		t.Error("SourceIDで検索できませんでした")
	}
}

// ========== タスク1.3: 敵タイプ新規フィールドのテスト ==========

// TestEnemyType_拡張フィールドの確認 はEnemyType拡張フィールドが正しく設定されることを確認します。
func TestEnemyType_拡張フィールドの確認(t *testing.T) {
	normalPassive := &EnemyPassiveSkill{
		ID:   "slime_normal",
		Name: "ぷるぷるボディ",
	}
	enhancedPassive := &EnemyPassiveSkill{
		ID:   "slime_enhanced",
		Name: "怒りのスライム",
	}
	normalAction := EnemyAction{
		ActionType: EnemyActionAttack,
		AttackType: "physical",
	}
	enhancedAction := EnemyAction{
		ActionType:   EnemyActionBuff,
		EffectColumn: EffectColumn("attackUp"),
		EffectValue:  0.3,
		Duration:     10.0,
	}

	enemyType := EnemyType{
		ID:                      "slime",
		Name:                    "スライム",
		BaseHP:                  50,
		ResolvedNormalActions:   []EnemyAction{normalAction},
		ResolvedEnhancedActions: []EnemyAction{normalAction, enhancedAction},
		NormalPassive:           normalPassive,
		EnhancedPassive:         enhancedPassive,
		DropItemCategory:        "core",
		DropItemTypeID:          "fire",
	}

	// 通常行動パターン
	if len(enemyType.ResolvedNormalActions) != 1 {
		t.Errorf("ResolvedNormalActionsの長さが期待値と異なります: got %d, want 1", len(enemyType.ResolvedNormalActions))
	}
	if enemyType.ResolvedNormalActions[0].ActionType != EnemyActionAttack {
		t.Error("ResolvedNormalActions[0]のActionTypeがAttackであるべきです")
	}

	// 強化行動パターン
	if len(enemyType.ResolvedEnhancedActions) != 2 {
		t.Errorf("ResolvedEnhancedActionsの長さが期待値と異なります: got %d, want 2", len(enemyType.ResolvedEnhancedActions))
	}

	// 通常パッシブ
	if enemyType.NormalPassive == nil {
		t.Error("NormalPassiveがnilです")
	}
	if enemyType.NormalPassive.ID != "slime_normal" {
		t.Errorf("NormalPassive.IDが期待値と異なります: got %s, want slime_normal", enemyType.NormalPassive.ID)
	}

	// 強化パッシブ
	if enemyType.EnhancedPassive == nil {
		t.Error("EnhancedPassiveがnilです")
	}
	if enemyType.EnhancedPassive.ID != "slime_enhanced" {
		t.Errorf("EnhancedPassive.IDが期待値と異なります: got %s, want slime_enhanced", enemyType.EnhancedPassive.ID)
	}

	// ドロップアイテム設定
	if enemyType.DropItemCategory != "core" {
		t.Errorf("DropItemCategoryが期待値と異なります: got %s, want core", enemyType.DropItemCategory)
	}
	if enemyType.DropItemTypeID != "fire" {
		t.Errorf("DropItemTypeIDが期待値と異なります: got %s, want fire", enemyType.DropItemTypeID)
	}
}

// TestEnemyType_ランク範囲 はランクの有効範囲を確認します。
func TestEnemyType_ランク範囲(t *testing.T) {
	tests := []struct {
		name     string
		rank     int
		expected bool // true = 有効、false = 無効
	}{
		{"ランク0（無効）", 0, false},
		{"ランク1（最小有効値）", 1, true},
		{"ランク5（中間値）", 5, true},
		{"負のランク（無効）", -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemyType := EnemyType{Rank: tt.rank}
			result := enemyType.IsValidRank()
			if result != tt.expected {
				t.Errorf("IsValidRank()が期待値と異なります: got %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestEnemyType_行動パターンバリデーション は行動パターンの最低1つの行動を保証するバリデーションを確認します。
func TestEnemyType_行動パターンバリデーション(t *testing.T) {
	tests := []struct {
		name     string
		pattern  []EnemyAction
		expected bool // true = 有効、false = 無効
	}{
		{"空パターン（無効）", []EnemyAction{}, false},
		{"1つの行動（有効）", []EnemyAction{{ActionType: EnemyActionAttack}}, true},
		{"複数の行動（有効）", []EnemyAction{
			{ActionType: EnemyActionAttack},
			{ActionType: EnemyActionBuff},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemyType := EnemyType{ResolvedNormalActions: tt.pattern}
			result := enemyType.HasValidNormalActionPattern()
			if result != tt.expected {
				t.Errorf("HasValidNormalActionPattern()が期待値と異なります: got %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestEnemyType_ドロップカテゴリバリデーション はドロップカテゴリの有効値を確認します。
func TestEnemyType_ドロップカテゴリバリデーション(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool // true = 有効、false = 無効
	}{
		{"core（有効）", "core", true},
		{"skill（有効）", "skill", true},
		{"空文字（無効）", "", false},
		{"不正な値（無効）", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemyType := EnemyType{DropItemCategory: tt.category}
			result := enemyType.IsValidDropItemCategory()
			if result != tt.expected {
				t.Errorf("IsValidDropItemCategory()が期待値と異なります: got %v, want %v", result, tt.expected)
			}
		})
	}
}

// ========== タスク1.4: 敵インスタンス行動管理機能のテスト ==========

// TestEnemyModel_行動管理フィールドの確認 はEnemyModelの行動管理フィールドを確認します。
func TestEnemyModel_行動管理フィールドの確認(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{ActionType: EnemyActionAttack, AttackType: "physical"},
			{ActionType: EnemyActionBuff, EffectColumn: EffectColumn("attackUp"), EffectValue: 0.3, Duration: 10.0},
		},
		ResolvedEnhancedActions: []EnemyAction{
			{ActionType: EnemyActionAttack, AttackType: "magic"},
		},
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// 行動インデックスは0で初期化
	if enemy.ActionIndex != 0 {
		t.Errorf("ActionIndexが期待値と異なります: got %d, want 0", enemy.ActionIndex)
	}

	// 適用中パッシブIDは空で初期化
	if enemy.ActivePassiveID != "" {
		t.Errorf("ActivePassiveIDが期待値と異なります: got %s, want empty", enemy.ActivePassiveID)
	}
}

// TestEnemyModel_GetCurrentAction は現在実行すべき行動を取得するメソッドを確認します。
func TestEnemyModel_GetCurrentAction(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{ActionType: EnemyActionAttack, AttackType: "physical"},
			{ActionType: EnemyActionBuff, EffectColumn: EffectColumn("attackUp")},
		},
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// 初期状態では最初の行動
	action := enemy.GetCurrentAction()
	if action.ActionType != EnemyActionAttack {
		t.Errorf("GetCurrentAction()が期待値と異なります: got %v, want Attack", action.ActionType)
	}
	if action.AttackType != "physical" {
		t.Errorf("AttackTypeが期待値と異なります: got %s, want physical", action.AttackType)
	}
}

// TestEnemyModel_AdvanceActionIndex は行動インデックスを進める（ループ対応）メソッドを確認します。
func TestEnemyModel_AdvanceActionIndex(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{ActionType: EnemyActionAttack},
			{ActionType: EnemyActionBuff},
			{ActionType: EnemyActionDebuff},
		},
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// 初期状態
	if enemy.ActionIndex != 0 {
		t.Errorf("初期ActionIndexが期待値と異なります: got %d, want 0", enemy.ActionIndex)
	}

	// 1回進める
	enemy.AdvanceActionIndex()
	if enemy.ActionIndex != 1 {
		t.Errorf("1回目進めた後のActionIndexが期待値と異なります: got %d, want 1", enemy.ActionIndex)
	}

	// 2回進める
	enemy.AdvanceActionIndex()
	if enemy.ActionIndex != 2 {
		t.Errorf("2回目進めた後のActionIndexが期待値と異なります: got %d, want 2", enemy.ActionIndex)
	}

	// 3回進める（ループして0に戻る）
	enemy.AdvanceActionIndex()
	if enemy.ActionIndex != 0 {
		t.Errorf("ループ後のActionIndexが期待値と異なります: got %d, want 0", enemy.ActionIndex)
	}
}

// TestEnemyModel_GetCurrentPattern はフェーズに応じた行動パターンを返すメソッドを確認します。
func TestEnemyModel_GetCurrentPattern(t *testing.T) {
	normalAction := EnemyAction{ActionType: EnemyActionAttack, AttackType: "physical"}
	enhancedAction := EnemyAction{ActionType: EnemyActionAttack, AttackType: "magic"}

	enemyType := EnemyType{
		ResolvedNormalActions:   []EnemyAction{normalAction},
		ResolvedEnhancedActions: []EnemyAction{enhancedAction},
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// 通常フェーズでは通常パターン
	pattern := enemy.GetCurrentPattern()
	if len(pattern) != 1 || pattern[0].AttackType != "physical" {
		t.Error("通常フェーズでは通常パターンを返すべきです")
	}

	// 強化フェーズでは強化パターン
	enemy.Phase = PhaseEnhanced
	pattern = enemy.GetCurrentPattern()
	if len(pattern) != 1 || pattern[0].AttackType != "magic" {
		t.Error("強化フェーズでは強化パターンを返すべきです")
	}
}

// TestEnemyModel_GetCurrentPattern_強化パターン空の場合 は強化パターンが空の場合に通常パターンを継続することを確認します。
func TestEnemyModel_GetCurrentPattern_強化パターン空の場合(t *testing.T) {
	normalAction := EnemyAction{ActionType: EnemyActionAttack, AttackType: "physical"}

	enemyType := EnemyType{
		ResolvedNormalActions:   []EnemyAction{normalAction},
		ResolvedEnhancedActions: []EnemyAction{}, // 空
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)
	enemy.Phase = PhaseEnhanced

	// 強化パターンが空の場合は通常パターンを継続
	pattern := enemy.GetCurrentPattern()
	if len(pattern) != 1 || pattern[0].AttackType != "physical" {
		t.Error("強化パターンが空の場合は通常パターンを継続すべきです")
	}
}

// TestEnemyModel_ResetActionIndex はフェーズ遷移時に行動インデックスをリセットすることを確認します。
func TestEnemyModel_ResetActionIndex(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{ActionType: EnemyActionAttack},
			{ActionType: EnemyActionBuff},
		},
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// インデックスを進める
	enemy.AdvanceActionIndex()
	if enemy.ActionIndex != 1 {
		t.Errorf("AdvanceActionIndex後のActionIndexが期待値と異なります: got %d, want 1", enemy.ActionIndex)
	}

	// リセット
	enemy.ResetActionIndex()
	if enemy.ActionIndex != 0 {
		t.Errorf("ResetActionIndex後のActionIndexが期待値と異なります: got %d, want 0", enemy.ActionIndex)
	}
}

// TestEnemyModel_行動パターン空の場合のGetCurrentAction は行動パターンが空の場合のデフォルト動作を確認します。
func TestEnemyModel_行動パターン空の場合のGetCurrentAction(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{}, // 空
	}

	enemy := NewEnemy("enemy_001", "テスト敵", 5, 100, 15, enemyType)

	// 空パターンの場合はデフォルト攻撃を返す
	action := enemy.GetCurrentAction()
	if action.ActionType != EnemyActionAttack {
		t.Error("空パターンの場合はデフォルトの攻撃行動を返すべきです")
	}
}

// ========== タスク1.5: JSONシリアライズ/デシリアライズの確認テスト ==========

// TestEnemyAction_JSONシリアライズ はEnemyActionのJSONシリアライズを確認します。
func TestEnemyAction_JSONシリアライズ(t *testing.T) {
	action := EnemyAction{
		ActionType:   EnemyActionBuff,
		EffectColumn: EffectColumn("attackUp"),
		EffectValue:  0.3,
		Duration:     10.0,
	}

	// シリアライズ
	data, err := json.Marshal(action)
	if err != nil {
		t.Fatalf("JSONシリアライズに失敗: %v", err)
	}

	// デシリアライズ
	var restored EnemyAction
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSONデシリアライズに失敗: %v", err)
	}

	// 検証
	if restored.ActionType != action.ActionType {
		t.Errorf("ActionTypeが一致しません: got %v, want %v", restored.ActionType, action.ActionType)
	}
	if restored.EffectColumn != action.EffectColumn {
		t.Errorf("EffectColumnが一致しません: got %s, want %s", restored.EffectColumn, action.EffectColumn)
	}
	if restored.EffectValue != action.EffectValue {
		t.Errorf("EffectValueが一致しません: got %f, want %f", restored.EffectValue, action.EffectValue)
	}
	if restored.Duration != action.Duration {
		t.Errorf("Durationが一致しません: got %f, want %f", restored.Duration, action.Duration)
	}
}

// TestEnemyPassiveSkill_JSONシリアライズ はEnemyPassiveSkillのJSONシリアライズを確認します。
func TestEnemyPassiveSkill_JSONシリアライズ(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:          "slime_normal",
		Name:        "ぷるぷるボディ",
		Description: "物理ダメージを10%軽減",
		Effects: map[EffectColumn]float64{
			ColDamageCut: 0.1,
		},
	}

	// シリアライズ
	data, err := json.Marshal(passive)
	if err != nil {
		t.Fatalf("JSONシリアライズに失敗: %v", err)
	}

	// デシリアライズ
	var restored EnemyPassiveSkill
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSONデシリアライズに失敗: %v", err)
	}

	// 検証
	if restored.ID != passive.ID {
		t.Errorf("IDが一致しません: got %s, want %s", restored.ID, passive.ID)
	}
	if restored.Name != passive.Name {
		t.Errorf("Nameが一致しません: got %s, want %s", restored.Name, passive.Name)
	}
	if restored.Effects[ColDamageCut] != passive.Effects[ColDamageCut] {
		t.Errorf("Effects[ColDamageCut]が一致しません: got %f, want %f", restored.Effects[ColDamageCut], passive.Effects[ColDamageCut])
	}
}

// TestEnemyType拡張フィールド_JSONシリアライズ はEnemyType拡張フィールドのJSONシリアライズを確認します。
func TestEnemyType拡張フィールド_JSONシリアライズ(t *testing.T) {
	enemyType := EnemyType{
		ID:     "slime",
		Name:   "スライム",
		BaseHP: 50,
		Rank:   1,
		ResolvedNormalActions: []EnemyAction{
			{ActionType: EnemyActionAttack, AttackType: "physical"},
		},
		ResolvedEnhancedActions: []EnemyAction{
			{ActionType: EnemyActionAttack, AttackType: "physical"},
			{ActionType: EnemyActionBuff, EffectColumn: EffectColumn("attackUp"), EffectValue: 0.3, Duration: 10.0},
		},
		NormalPassive: &EnemyPassiveSkill{
			ID:   "slime_normal",
			Name: "ぷるぷるボディ",
			Effects: map[EffectColumn]float64{
				ColDamageCut: 0.1,
			},
		},
		EnhancedPassive: &EnemyPassiveSkill{
			ID:   "slime_enhanced",
			Name: "怒りのスライム",
			Effects: map[EffectColumn]float64{
				ColDamageMultiplier: 1.2,
			},
		},
		DropItemCategory: "core",
		DropItemTypeID:   "fire",
	}

	// シリアライズ
	data, err := json.Marshal(enemyType)
	if err != nil {
		t.Fatalf("JSONシリアライズに失敗: %v", err)
	}

	// デシリアライズ
	var restored EnemyType
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSONデシリアライズに失敗: %v", err)
	}

	// 検証
	if restored.Rank != enemyType.Rank {
		t.Errorf("Rankが一致しません: got %d, want %d", restored.Rank, enemyType.Rank)
	}
	if len(restored.ResolvedNormalActions) != len(enemyType.ResolvedNormalActions) {
		t.Errorf("ResolvedNormalActionsの長さが一致しません: got %d, want %d", len(restored.ResolvedNormalActions), len(enemyType.ResolvedNormalActions))
	}
	if len(restored.ResolvedEnhancedActions) != len(enemyType.ResolvedEnhancedActions) {
		t.Errorf("ResolvedEnhancedActionsの長さが一致しません: got %d, want %d", len(restored.ResolvedEnhancedActions), len(enemyType.ResolvedEnhancedActions))
	}
	if restored.NormalPassive == nil {
		t.Error("NormalPassiveがnilになっています")
	} else if restored.NormalPassive.ID != enemyType.NormalPassive.ID {
		t.Errorf("NormalPassive.IDが一致しません: got %s, want %s", restored.NormalPassive.ID, enemyType.NormalPassive.ID)
	}
	if restored.EnhancedPassive == nil {
		t.Error("EnhancedPassiveがnilになっています")
	} else if restored.EnhancedPassive.ID != enemyType.EnhancedPassive.ID {
		t.Errorf("EnhancedPassive.IDが一致しません: got %s, want %s", restored.EnhancedPassive.ID, enemyType.EnhancedPassive.ID)
	}
	if restored.DropItemCategory != enemyType.DropItemCategory {
		t.Errorf("DropItemCategoryが一致しません: got %s, want %s", restored.DropItemCategory, enemyType.DropItemCategory)
	}
	if restored.DropItemTypeID != enemyType.DropItemTypeID {
		t.Errorf("DropItemTypeIDが一致しません: got %s, want %s", restored.DropItemTypeID, enemyType.DropItemTypeID)
	}
}

// ========== タスク7.1: ドメイン層のユニットテスト ==========

// TestEnemyAction_CalculateDamage はダメージ計算を確認します。
func TestEnemyAction_CalculateDamage(t *testing.T) {
	tests := []struct {
		name           string
		damageBase     float64
		damagePerLevel float64
		level          int
		expected       int
	}{
		{"レベル1基本攻撃", 10.0, 2.0, 1, 12},      // 10 + 1*2 = 12
		{"レベル10基本攻撃", 10.0, 2.0, 10, 30},    // 10 + 10*2 = 30
		{"レベル100強力攻撃", 50.0, 5.0, 100, 550}, // 50 + 100*5 = 550
		{"最低ダメージ保証", 0.0, 0.0, 1, 1},        // 0 + 0*1 = 0 → 最低1
		{"小数点ダメージ", 5.5, 1.5, 3, 10},        // 5.5 + 3*1.5 = 10.0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := EnemyAction{
				ActionType:     EnemyActionAttack,
				DamageBase:     tt.damageBase,
				DamagePerLevel: tt.damagePerLevel,
			}
			damage := action.CalculateDamage(tt.level)
			if damage != tt.expected {
				t.Errorf("CalculateDamage()が期待値と異なります: got %d, want %d", damage, tt.expected)
			}
		})
	}
}

// TestEnemyAction_GetChargeTimeMs はチャージタイムをミリ秒で取得する機能を確認します。
func TestEnemyAction_GetChargeTimeMs(t *testing.T) {
	tests := []struct {
		name       string
		chargeTime time.Duration
		expected   int64
	}{
		{"0秒", 0, 0},
		{"1秒", 1 * time.Second, 1000},
		{"3秒", 3 * time.Second, 3000},
		{"500ミリ秒", 500 * time.Millisecond, 500},
		{"2.5秒", 2500 * time.Millisecond, 2500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := EnemyAction{
				ChargeTime: tt.chargeTime,
			}
			result := action.GetChargeTimeMs()
			if result != tt.expected {
				t.Errorf("GetChargeTimeMs()が期待値と異なります: got %d, want %d", result, tt.expected)
			}
		})
	}
}

// TestEnemyAction_IsDefense はディフェンス行動判定を確認します。
func TestEnemyAction_IsDefense(t *testing.T) {
	defense := EnemyAction{ActionType: EnemyActionDefense}
	attack := EnemyAction{ActionType: EnemyActionAttack}

	if !defense.IsDefense() {
		t.Error("Defense行動でIsDefense()がtrueを返すべきです")
	}
	if attack.IsDefense() {
		t.Error("Attack行動でIsDefense()がfalseを返すべきです")
	}
}

// TestEnemyDefenseType_定数の確認 はEnemyDefenseType定数が正しく定義されていることを確認します。
func TestEnemyDefenseType_定数の確認(t *testing.T) {
	if DefensePhysicalCut != "physical_cut" {
		t.Errorf("DefensePhysicalCutが期待値と異なります: got %s, want physical_cut", DefensePhysicalCut)
	}
	if DefenseMagicCut != "magic_cut" {
		t.Errorf("DefenseMagicCutが期待値と異なります: got %s, want magic_cut", DefenseMagicCut)
	}
	if DefenseDebuffEvade != "debuff_evade" {
		t.Errorf("DefenseDebuffEvadeが期待値と異なります: got %s, want debuff_evade", DefenseDebuffEvade)
	}
}

// TestEnemyModel_AdvanceActionIndex_境界値 は行動パターンループの境界値を確認します。
func TestEnemyModel_AdvanceActionIndex_境界値(t *testing.T) {
	tests := []struct {
		name          string
		patternLen    int
		advanceTimes  int
		expectedIndex int
	}{
		{"1行動パターン_1回進める", 1, 1, 0},
		{"1行動パターン_3回進める", 1, 3, 0},
		{"3行動パターン_2回進める", 3, 2, 2},
		{"3行動パターン_3回進める_ループ", 3, 3, 0},
		{"3行動パターン_5回進める_複数ループ", 3, 5, 2},
		{"5行動パターン_10回進める", 5, 10, 0},
		{"5行動パターン_7回進める", 5, 7, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// パターンを作成
			pattern := make([]EnemyAction, tt.patternLen)
			for i := 0; i < tt.patternLen; i++ {
				pattern[i] = EnemyAction{ActionType: EnemyActionAttack}
			}

			enemyType := EnemyType{
				ResolvedNormalActions: pattern,
			}
			enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

			// 指定回数進める
			for i := 0; i < tt.advanceTimes; i++ {
				enemy.AdvanceActionIndex()
			}

			if enemy.ActionIndex != tt.expectedIndex {
				t.Errorf("ActionIndexが期待値と異なります: got %d, want %d", enemy.ActionIndex, tt.expectedIndex)
			}
		})
	}
}

// TestEnemyModel_GetCurrentAction_フェーズ遷移後 はフェーズ遷移後に正しいパターンから行動を取得することを確認します。
func TestEnemyModel_GetCurrentAction_フェーズ遷移後(t *testing.T) {
	normalAction := EnemyAction{
		ID:         "normal_attack",
		ActionType: EnemyActionAttack,
		AttackType: "physical",
	}
	enhancedAction := EnemyAction{
		ID:         "enhanced_attack",
		ActionType: EnemyActionAttack,
		AttackType: "magic",
	}

	enemyType := EnemyType{
		ResolvedNormalActions:   []EnemyAction{normalAction},
		ResolvedEnhancedActions: []EnemyAction{enhancedAction},
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

	// 通常フェーズで通常パターンの行動を取得
	action := enemy.GetCurrentAction()
	if action.ID != "normal_attack" {
		t.Errorf("通常フェーズで通常行動を返すべきです: got %s", action.ID)
	}

	// フェーズ遷移
	enemy.TransitionToEnhanced()

	// 強化フェーズで強化パターンの行動を取得
	action = enemy.GetCurrentAction()
	if action.ID != "enhanced_attack" {
		t.Errorf("強化フェーズで強化行動を返すべきです: got %s", action.ID)
	}
}

// TestEnemyModel_ResetActionIndex_フェーズ遷移シナリオ はフェーズ遷移時のインデックスリセットを確認します。
func TestEnemyModel_ResetActionIndex_フェーズ遷移シナリオ(t *testing.T) {
	normalActions := []EnemyAction{
		{ID: "normal_1", ActionType: EnemyActionAttack},
		{ID: "normal_2", ActionType: EnemyActionBuff},
		{ID: "normal_3", ActionType: EnemyActionDebuff},
	}
	enhancedActions := []EnemyAction{
		{ID: "enhanced_1", ActionType: EnemyActionAttack},
		{ID: "enhanced_2", ActionType: EnemyActionDefense},
	}

	enemyType := EnemyType{
		ResolvedNormalActions:   normalActions,
		ResolvedEnhancedActions: enhancedActions,
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

	// 通常フェーズで2回進める
	enemy.AdvanceActionIndex()
	enemy.AdvanceActionIndex()
	if enemy.ActionIndex != 2 {
		t.Fatalf("ActionIndexが2であるべきです: got %d", enemy.ActionIndex)
	}

	// フェーズ遷移とリセット
	enemy.TransitionToEnhanced()
	enemy.ResetActionIndex()

	// インデックスが0にリセットされていること
	if enemy.ActionIndex != 0 {
		t.Errorf("ResetActionIndex後はActionIndexが0であるべきです: got %d", enemy.ActionIndex)
	}

	// 強化パターンの最初の行動を取得
	action := enemy.GetCurrentAction()
	if action.ID != "enhanced_1" {
		t.Errorf("リセット後は強化パターンの最初の行動を返すべきです: got %s", action.ID)
	}
}

// TestEnemyPassiveSkill_空のEffects は空の効果マップを持つパッシブスキルを確認します。
func TestEnemyPassiveSkill_空のEffects(t *testing.T) {
	passive := EnemyPassiveSkill{
		ID:      "empty_passive",
		Name:    "空パッシブ",
		Effects: map[EffectColumn]float64{},
	}

	entry := passive.ToEntry()

	if len(entry.Values) != 0 {
		t.Errorf("Valuesが空であるべきです: got %d entries", len(entry.Values))
	}
	if entry.SourceType != SourcePassive {
		t.Error("SourceTypeがSourcePassiveであるべきです")
	}
}

// TestEnemyModel_チャージ状態管理 はチャージ状態管理機能を確認します。
func TestEnemyModel_チャージ状態管理(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{
				ID:         "attack",
				ActionType: EnemyActionAttack,
				ChargeTime: 3 * time.Second,
			},
		},
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)
	action := enemy.GetCurrentAction()

	now := time.Now()

	// チャージ開始
	enemy.StartCharging(action, now)

	if enemy.WaitMode != WaitModeCharging {
		t.Error("チャージ開始後はWaitModeChargingであるべきです")
	}
	if enemy.PendingAction == nil {
		t.Error("PendingActionが設定されるべきです")
	}

	// 進捗確認（開始直後）
	progress := enemy.GetChargeProgress(now)
	if progress != 0 {
		t.Errorf("開始直後は進捗0であるべきです: got %f", progress)
	}

	// 進捗確認（1.5秒後 = 50%）
	halfTime := now.Add(1500 * time.Millisecond)
	progress = enemy.GetChargeProgress(halfTime)
	if progress < 0.49 || progress > 0.51 {
		t.Errorf("1.5秒後は約50%%であるべきです: got %f", progress)
	}

	// 完了チェック（1.5秒後）
	if enemy.IsChargeComplete(halfTime) {
		t.Error("1.5秒後はチャージ完了ではないべきです")
	}

	// 完了チェック（3秒後）
	completeTime := now.Add(3 * time.Second)
	if !enemy.IsChargeComplete(completeTime) {
		t.Error("3秒後はチャージ完了であるべきです")
	}

	// 進捗確認（完了後は1.0）
	progress = enemy.GetChargeProgress(completeTime)
	if progress != 1.0 {
		t.Errorf("完了後は進捗1.0であるべきです: got %f", progress)
	}
}

// TestEnemyModel_ディフェンス状態管理 はディフェンス状態管理機能を確認します。
func TestEnemyModel_ディフェンス状態管理(t *testing.T) {
	enemyType := EnemyType{
		ResolvedNormalActions: []EnemyAction{
			{ID: "defense", ActionType: EnemyActionDefense},
		},
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

	now := time.Now()
	duration := 5 * time.Second

	// ディフェンス開始
	enemy.StartDefense(DefensePhysicalCut, 0.5, duration, now)

	if enemy.WaitMode != WaitModeDefending {
		t.Error("ディフェンス開始後はWaitModeDefendingであるべきです")
	}
	if enemy.ActiveDefenseType != DefensePhysicalCut {
		t.Error("ActiveDefenseTypeがDefensePhysicalCutであるべきです")
	}
	if enemy.DefenseValue != 0.5 {
		t.Errorf("DefenseValueが0.5であるべきです: got %f", enemy.DefenseValue)
	}

	// ディフェンス有効チェック（開始直後）
	if !enemy.IsDefenseActive(now) {
		t.Error("開始直後はディフェンス有効であるべきです")
	}

	// ディフェンス有効チェック（2秒後）
	midTime := now.Add(2 * time.Second)
	if !enemy.IsDefenseActive(midTime) {
		t.Error("2秒後はディフェンス有効であるべきです")
	}

	// ディフェンス有効チェック（6秒後 = 終了後）
	afterTime := now.Add(6 * time.Second)
	if enemy.IsDefenseActive(afterTime) {
		t.Error("6秒後はディフェンス無効であるべきです")
	}

	// 残り時間確認（2秒後）
	remaining := enemy.GetDefenseRemainingTime(midTime)
	expectedRemaining := 3 * time.Second
	if remaining != expectedRemaining {
		t.Errorf("残り時間が3秒であるべきです: got %v", remaining)
	}
}

// TestEnemyModel_ExecuteChargedAction はチャージ完了後の行動実行を確認します。
func TestEnemyModel_ExecuteChargedAction(t *testing.T) {
	actions := []EnemyAction{
		{ID: "action_1", ActionType: EnemyActionAttack},
		{ID: "action_2", ActionType: EnemyActionBuff},
	}
	enemyType := EnemyType{
		ResolvedNormalActions: actions,
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

	// 最初の行動でチャージ開始
	action := enemy.GetCurrentAction()
	enemy.StartCharging(action, time.Now())

	// チャージ完了後に行動実行
	executedAction := enemy.ExecuteChargedAction()

	if executedAction == nil {
		t.Fatal("ExecuteChargedActionはnilを返すべきではありません")
	}
	if executedAction.ID != "action_1" {
		t.Errorf("実行された行動がaction_1であるべきです: got %s", executedAction.ID)
	}
	if enemy.WaitMode != WaitModeNone {
		t.Error("実行後はWaitModeNoneであるべきです")
	}
	if enemy.PendingAction != nil {
		t.Error("実行後はPendingActionがnilであるべきです")
	}
	// 行動インデックスが進んでいること
	if enemy.ActionIndex != 1 {
		t.Errorf("実行後はActionIndexが1であるべきです: got %d", enemy.ActionIndex)
	}
}

// TestEnemyModel_EndDefense はディフェンス終了処理を確認します。
func TestEnemyModel_EndDefense(t *testing.T) {
	actions := []EnemyAction{
		{ID: "defense_1", ActionType: EnemyActionDefense},
		{ID: "attack_1", ActionType: EnemyActionAttack},
	}
	enemyType := EnemyType{
		ResolvedNormalActions: actions,
	}
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, enemyType)

	// ディフェンス開始
	enemy.StartDefense(DefenseMagicCut, 0.3, 5*time.Second, time.Now())

	// ディフェンス終了
	enemy.EndDefense()

	if enemy.WaitMode != WaitModeNone {
		t.Error("終了後はWaitModeNoneであるべきです")
	}
	if enemy.ActiveDefenseType != "" {
		t.Error("終了後はActiveDefenseTypeが空であるべきです")
	}
	if enemy.DefenseValue != 0 {
		t.Error("終了後はDefenseValueが0であるべきです")
	}
	// 行動インデックスが進んでいること
	if enemy.ActionIndex != 1 {
		t.Errorf("終了後はActionIndexが1であるべきです: got %d", enemy.ActionIndex)
	}
}

// TestEnemyWaitMode_String はEnemyWaitModeのString()メソッドを確認します。
func TestEnemyWaitMode_String(t *testing.T) {
	tests := []struct {
		mode     EnemyWaitMode
		expected string
	}{
		{WaitModeNone, "なし"},
		{WaitModeCharging, "チャージ中"},
		{WaitModeDefending, "ディフェンス中"},
		{EnemyWaitMode(99), "不明"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.mode.String() != tt.expected {
				t.Errorf("String()が期待値と異なります: got %s, want %s", tt.mode.String(), tt.expected)
			}
		})
	}
}

// TestEnemyModel_GetDefenseTypeName はディフェンス種別名取得を確認します。
func TestEnemyModel_GetDefenseTypeName(t *testing.T) {
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, EnemyType{})

	tests := []struct {
		defenseType EnemyDefenseType
		expected    string
	}{
		{DefensePhysicalCut, "物理防御"},
		{DefenseMagicCut, "魔法防御"},
		{DefenseDebuffEvade, "デバフ回避"},
		{"", "防御"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			enemy.ActiveDefenseType = tt.defenseType
			result := enemy.GetDefenseTypeName()
			if result != tt.expected {
				t.Errorf("GetDefenseTypeName()が期待値と異なります: got %s, want %s", result, tt.expected)
			}
		})
	}
}

// ========== ボルテージシステムのテスト ==========

// TestNewEnemy_ボルテージ初期化 は敵生成時にボルテージが100.0で初期化されることを確認します。
func TestNewEnemy_ボルテージ初期化(t *testing.T) {
	enemyType := EnemyType{
		ID:   "slime",
		Name: "スライム",
	}
	enemy := NewEnemy("enemy_001", "スライム", 1, 100, 10, enemyType)

	// ボルテージは100.0で初期化されるべき
	if enemy.Voltage != 100.0 {
		t.Errorf("ボルテージが100.0で初期化されるべきです: got %f", enemy.Voltage)
	}
}

// TestEnemyModel_GetVoltage はボルテージ取得メソッドを確認します。
func TestEnemyModel_GetVoltage(t *testing.T) {
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, EnemyType{})

	// 初期値の取得
	voltage := enemy.GetVoltage()
	if voltage != 100.0 {
		t.Errorf("GetVoltage()が100.0を返すべきです: got %f", voltage)
	}

	// 値を変更して取得
	enemy.Voltage = 150.0
	voltage = enemy.GetVoltage()
	if voltage != 150.0 {
		t.Errorf("GetVoltage()が150.0を返すべきです: got %f", voltage)
	}
}

// TestEnemyModel_SetVoltage はボルテージ設定メソッドを確認します。
func TestEnemyModel_SetVoltage(t *testing.T) {
	enemy := NewEnemy("test", "テスト敵", 1, 100, 10, EnemyType{})

	// ボルテージを設定
	enemy.SetVoltage(175.5)
	if enemy.Voltage != 175.5 {
		t.Errorf("SetVoltage()でボルテージが設定されるべきです: got %f, want 175.5", enemy.Voltage)
	}

	// 別の値を設定
	enemy.SetVoltage(200.0)
	if enemy.Voltage != 200.0 {
		t.Errorf("SetVoltage()でボルテージが更新されるべきです: got %f, want 200.0", enemy.Voltage)
	}
}

// TestEnemyModel_GetVoltageMultiplier はダメージ乗算用の倍率取得メソッドを確認します。
func TestEnemyModel_GetVoltageMultiplier(t *testing.T) {
	tests := []struct {
		name     string
		voltage  float64
		expected float64
	}{
		{"100%で等倍", 100.0, 1.0},
		{"150%で1.5倍", 150.0, 1.5},
		{"200%で2.0倍", 200.0, 2.0},
		{"50%で0.5倍", 50.0, 0.5},
		{"115.5%で1.155倍", 115.5, 1.155},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := NewEnemy("test", "テスト敵", 1, 100, 10, EnemyType{})
			enemy.SetVoltage(tt.voltage)

			multiplier := enemy.GetVoltageMultiplier()
			if multiplier != tt.expected {
				t.Errorf("GetVoltageMultiplier()が期待値と異なります: got %f, want %f", multiplier, tt.expected)
			}
		})
	}
}

// TestEnemyType_VoltageRisePer10s はEnemyTypeのボルテージ上昇率フィールドを確認します。
func TestEnemyType_VoltageRisePer10s(t *testing.T) {
	enemyType := EnemyType{
		ID:                "slime",
		Name:              "スライム",
		VoltageRisePer10s: 15.0, // 10秒あたり15ポイント上昇
	}

	if enemyType.VoltageRisePer10s != 15.0 {
		t.Errorf("VoltageRisePer10sが期待値と異なります: got %f, want 15.0", enemyType.VoltageRisePer10s)
	}
}

// TestEnemyType_GetVoltageRisePer10s は負の値が設定された場合に0を返すことを確認します。
func TestEnemyType_GetVoltageRisePer10s(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"正の値", 10.0, 10.0},
		{"ゼロ", 0.0, 0.0},
		{"負の値は0として扱う", -5.0, 0.0},
		{"大きな負の値", -100.0, 0.0},
		{"小数値", 7.5, 7.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemyType := EnemyType{
				ID:                "test",
				VoltageRisePer10s: tt.value,
			}

			result := enemyType.GetVoltageRisePer10s()
			if result != tt.expected {
				t.Errorf("GetVoltageRisePer10s()が期待値と異なります: got %f, want %f", result, tt.expected)
			}
		})
	}
}
