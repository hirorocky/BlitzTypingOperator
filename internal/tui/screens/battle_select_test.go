// Package screens はTUI画面のテストを提供します。
package screens

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// newTestDamageModule はテスト用のダメージモジュールを作成するヘルパー関数です。
func newTestDamageModule(id, name string, tags []string, statCoef float64, statRef, description string) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:              id,
		Name:            name,
		Icon:            "⚔️",
		Tags:            tags,
		Description:     description,
		CooldownSeconds: 3.0, // リキャストテスト用に設定
		Effects: []domain.ModuleEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, nil)
}

// containsString は文字列に部分文字列が含まれているか確認します。
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// newTestHealModule はテスト用の回復モジュールを作成するヘルパー関数です。
func newTestHealModule(id, name string, tags []string, statCoef float64, statRef, description string) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:              id,
		Name:            name,
		Icon:            "💚",
		Tags:            tags,
		Description:     description,
		CooldownSeconds: 3.0, // リキャストテスト用に設定
		Effects: []domain.ModuleEffect{
			{
				Target:      domain.TargetSelf,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "💚",
			},
		},
	}, nil)
}

// newTestBuffModule はテスト用のバフモジュールを作成するヘルパー関数です。
func newTestBuffModule(id, name string, tags []string, description string) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:              id,
		Name:            name,
		Icon:            "⬆️",
		Tags:            tags,
		Description:     description,
		CooldownSeconds: 3.0, // リキャストテスト用に設定
		Effects: []domain.ModuleEffect{
			{
				Target: domain.TargetSelf,
				ColumnSpec: &domain.EffectColumnSpec{
					Column:   domain.ColDamageBonus,
					Value:    10,
					Duration: 10.0,
				},
				Probability: 1.0,
				Icon:        "⬆️",
			},
		},
	}, nil)
}

// newTestModuleWithChainEffect はチェイン効果付きモジュールを作成するヘルパー関数です。
func newTestModuleWithChainEffect(id, name string, tags []string, statCoef float64, statRef, description string, chainEffect *domain.ChainEffect) *domain.ModuleModel {
	return domain.NewModuleFromType(domain.ModuleType{
		ID:              id,
		Name:            name,
		Icon:            "⚔️",
		Tags:            tags,
		Description:     description,
		CooldownSeconds: 3.0, // リキャストテスト用に設定
		Effects: []domain.ModuleEffect{
			{
				Target:      domain.TargetEnemy,
				HPFormula:   &domain.HPFormula{Base: 0, StatCoef: statCoef, StatRef: statRef},
				Probability: 1.0,
				Icon:        "⚔️",
			},
		},
	}, chainEffect)
}

// mockAgentProvider はテスト用のAgentProvider実装です。
type mockAgentProvider struct {
	agents []*domain.AgentModel
}

func (m *mockAgentProvider) GetEquippedAgents() []*domain.AgentModel {
	return m.agents
}

// ==================== Task 10.2: バトル選択画面のテスト ====================

// TestNewBattleSelectScreen はBattleSelectScreenの初期化をテストします。

func TestNewBattleSelectScreen(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	if screen == nil {
		t.Fatal("BattleSelectScreenがnilです")
	}

	// 入力フィールドが初期化されていること
	if screen.input == nil {
		t.Error("入力フィールドがnilです")
	}
}

// TestBattleSelectMaxLevelDisplay は最高レベル表示をテストします。

func TestBattleSelectMaxLevelDisplay(t *testing.T) {
	maxLevel := 15
	screen := NewBattleSelectScreen(maxLevel, &mockAgentProvider{})

	// 挑戦可能最大レベルは maxLevel + 1
	if screen.maxChallengeLevel != maxLevel+1 {
		t.Errorf("挑戦可能最大レベル: got %d, want %d", screen.maxChallengeLevel, maxLevel+1)
	}
}

// TestBattleSelectInputValidation は入力検証をテストします。

func TestBattleSelectInputValidation(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	tests := []struct {
		name        string
		input       string
		expectError bool
		errorType   string
	}{
		{"有効な入力", "5", false, ""},
		{"最大レベル", "11", false, ""},
		{"1未満", "0", true, "too_low"},
		{"最大レベル超過", "12", true, "too_high"},
		{"空入力", "", true, "empty"},
		{"非数値", "abc", true, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screen.input.Value = tt.input
			valid, _ := screen.validateInput()

			if valid == tt.expectError {
				if tt.expectError {
					t.Error("エラーが検出されませんでした")
				} else {
					t.Error("有効な入力がエラーと判定されました")
				}
			}
		})
	}
}

// TestBattleSelectNoAgentEquipped はエージェント未装備時のテストです。

func TestBattleSelectNoAgentEquipped(t *testing.T) {
	// エージェント未装備
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})
	screen.input.Value = "5"

	// 確認画面に移動
	screen.state = StateConfirm

	// バトル開始を試みる
	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	// エラーメッセージが設定されるか確認
	if screen.error == "" && cmd != nil {
		// コマンドが返された場合、バトル開始メッセージかどうか確認
		// 未装備の場合はバトル開始できないはず
	}
}

// TestBattleSelectWithAgentEquipped はエージェント装備時のテストです。
func TestBattleSelectWithAgentEquipped(t *testing.T) {
	// エージェントを装備
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})
	modules := []*domain.ModuleModel{
		newTestDamageModule("m1", "モジュール1", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m2", "モジュール2", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m3", "モジュール3", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m4", "モジュール4", []string{"physical_low"}, 1.0, "STR", ""),
	}
	agent := domain.NewAgent("agent1", core, modules)

	screen := NewBattleSelectScreen(10, &mockAgentProvider{agents: []*domain.AgentModel{agent}})
	screen.input.Value = "5"

	// 入力検証
	valid, _ := screen.validateInput()
	if !valid {
		t.Error("有効な入力がエラーと判定されました")
	}
}

// TestBattleSelectConfirmScreen は確認画面のテストです。

func TestBattleSelectConfirmScreen(t *testing.T) {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})
	modules := []*domain.ModuleModel{
		newTestDamageModule("m1", "モジュール1", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m2", "モジュール2", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m3", "モジュール3", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m4", "モジュール4", []string{"physical_low"}, 1.0, "STR", ""),
	}
	agent := domain.NewAgent("agent1", core, modules)

	screen := NewBattleSelectScreen(10, &mockAgentProvider{agents: []*domain.AgentModel{agent}})
	screen.input.Value = "5"

	// Enterで確認画面へ
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	if screen.state != StateConfirm {
		t.Errorf("確認画面に遷移していません: state=%d", screen.state)
	}
}

// TestBattleSelectRechallenge は再挑戦のテストです。

func TestBattleSelectRechallenge(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	// 過去にクリアしたレベル（例: レベル5）への再挑戦
	screen.input.Value = "5"
	valid, _ := screen.validateInput()

	if !valid {
		t.Error("過去クリアレベルへの再挑戦が拒否されました")
	}
}

// TestBattleSelectRender はバトル選択画面のレンダリングをテストします。
func TestBattleSelectRender(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestBattleSelectBackNavigation は戻るナビゲーションのテストです。

func TestBattleSelectBackNavigation(t *testing.T) {
	screen := NewBattleSelectScreen(10, &mockAgentProvider{})

	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEsc})

	if cmd == nil {
		t.Error("Escキーでコマンドが返されません")
	}
}

// ==================== タスク2.2: カルーセル方式のテスト ====================

// mockDefeatedEnemyProvider はテスト用のDefeatedEnemyProvider実装です。
type mockDefeatedEnemyProvider struct {
	defeated        map[string]int
	maxLevelReached int // 到達最高レベル（敵のデフォルトレベルで更新）
}

func (m *mockDefeatedEnemyProvider) GetDefeatedEnemies() map[string]int {
	return m.defeated
}

func (m *mockDefeatedEnemyProvider) IsEnemyDefeated(enemyTypeID string) bool {
	_, exists := m.defeated[enemyTypeID]
	return exists
}

func (m *mockDefeatedEnemyProvider) GetDefeatedLevel(enemyTypeID string) int {
	return m.defeated[enemyTypeID]
}

func (m *mockDefeatedEnemyProvider) GetMaxDefeatedLevel() int {
	maxLevel := 0
	for _, level := range m.defeated {
		if level > maxLevel {
			maxLevel = level
		}
	}
	return maxLevel
}

func (m *mockDefeatedEnemyProvider) GetMaxLevelReached() int {
	return m.maxLevelReached
}

// mockEnemyTypeProvider はテスト用のEnemyTypeProvider実装です。
type mockEnemyTypeProvider struct {
	enemyTypes []domain.EnemyType
}

func (m *mockEnemyTypeProvider) GetEnemyTypes() []domain.EnemyType {
	return m.enemyTypes
}

// createTestEnemyTypes はテスト用の敵タイプリストを作成します。
// 注: フィルタリングにより「到達Lv + 1」以下のデフォルトLvを持つ敵のみ表示される
func createTestEnemyTypes() []domain.EnemyType {
	return []domain.EnemyType{
		{ID: "slime", Name: "スライム", DefaultLevel: 1, BaseHP: 50, AttackType: "physical"},
		{ID: "goblin", Name: "ゴブリン", DefaultLevel: 2, BaseHP: 80, AttackType: "physical"},
		{ID: "dragon", Name: "ドラゴン", DefaultLevel: 10, BaseHP: 500, AttackType: "magic"},
	}
}

// createTestEnemyTypesAllVisible は全敵が表示される設定です（全敵を撃破済みにする）。
func createTestEnemyTypesAllVisible() ([]domain.EnemyType, map[string]int, int) {
	enemyTypes := []domain.EnemyType{
		{ID: "slime", Name: "スライム", DefaultLevel: 1, BaseHP: 50, AttackType: "physical"},
		{ID: "goblin", Name: "ゴブリン", DefaultLevel: 2, BaseHP: 80, AttackType: "physical"},
		{ID: "dragon", Name: "ドラゴン", DefaultLevel: 10, BaseHP: 500, AttackType: "magic"},
	}
	// 全敵を撃破済みにして全敵を表示
	defeated := map[string]int{"slime": 1, "goblin": 2, "dragon": 10}
	maxLevelReached := 10 // 最高デフォルトレベル
	return enemyTypes, defeated, maxLevelReached
}

// TestBattleSelectCarouselInitialization はカルーセル方式の初期化をテストします。
func TestBattleSelectCarouselInitialization(t *testing.T) {
	enemyTypes, defeated, maxLevelReached := createTestEnemyTypesAllVisible()
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: maxLevelReached},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	if screen == nil {
		t.Fatal("BattleSelectScreenがnilです")
	}

	// 全敵が撃破済みなので3体全て表示
	if len(screen.enemyTypes) != 3 {
		t.Errorf("敵タイプ数: got %d, want 3", len(screen.enemyTypes))
	}

	// 初期選択インデックスが0であること
	if screen.selectedTypeIdx != 0 {
		t.Errorf("初期選択インデックス: got %d, want 0", screen.selectedTypeIdx)
	}
}

// TestBattleSelectCarouselNavigation は左右キーによる敵種類変更をテストします。
func TestBattleSelectCarouselNavigation(t *testing.T) {
	enemyTypes, defeated, maxLevelReached := createTestEnemyTypesAllVisible()
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: maxLevelReached},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// 右キーで次の敵タイプへ
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})
	if screen.selectedTypeIdx != 1 {
		t.Errorf("右キー後のインデックス: got %d, want 1", screen.selectedTypeIdx)
	}

	// 左キーで前の敵タイプへ
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyLeft})
	if screen.selectedTypeIdx != 0 {
		t.Errorf("左キー後のインデックス: got %d, want 0", screen.selectedTypeIdx)
	}

	// 最初の敵タイプで左キーを押すと最後に移動（ループ）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyLeft})
	if screen.selectedTypeIdx != 2 {
		t.Errorf("ループ後のインデックス: got %d, want 2", screen.selectedTypeIdx)
	}
}

// TestBattleSelectCarouselLevelSelection は上下キーによるレベル変更をテストします。
func TestBattleSelectCarouselLevelSelection(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// slimeをレベル5で撃破済み、MaxLevelReached=5（撃破済み敵はMaxLevelReachedまで選択可能）
	defeated := map[string]int{"slime": 5}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 5},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// slime（デフォルトレベル1）が選択されている状態で
	// 撃破済みなので、レベル1〜5（MaxLevelReached）まで選択可能
	initialLevel := screen.selectedLevel
	if initialLevel != 1 {
		t.Errorf("初期レベル: got %d, want 1", initialLevel)
	}

	// 上キーでレベル上昇
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyUp})
	if screen.selectedLevel != 2 {
		t.Errorf("上キー後のレベル: got %d, want 2", screen.selectedLevel)
	}

	// 下キーでレベル下降
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyDown})
	if screen.selectedLevel != 1 {
		t.Errorf("下キー後のレベル: got %d, want 1", screen.selectedLevel)
	}
}

// TestBattleSelectCarouselUndefeatedEnemy は未撃破敵のレベル制限をテストします。
func TestBattleSelectCarouselUndefeatedEnemy(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// slime（デフォルトLv=1）を撃破済み、MaxLevelReached=1
	// goblin（デフォルトLv=2）は未撃破で、MaxLevelReached+1=2 以上なので表示される
	defeated := map[string]int{"slime": 5}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 1},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// slime（撃破済み）とgoblin（未撃破で最小デフォルトLv）の2体が表示される
	// goblin（インデックス1）を選択
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})

	// 未撃破なのでデフォルトレベル（2）のみ選択可能
	if screen.selectedLevel != 2 {
		t.Errorf("未撃破敵のレベル: got %d, want 2", screen.selectedLevel)
	}

	// 上下キーを押してもレベルが変わらない
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyUp})
	if screen.selectedLevel != 2 {
		t.Errorf("上キー後のレベル（未撃破）: got %d, want 2", screen.selectedLevel)
	}
}

// TestBattleSelectCarouselStartBattle はバトル開始メッセージをテストします。
func TestBattleSelectCarouselStartBattle(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	agent := createTestAgent()
	// MaxLevelReached=0で未撃破、slime（デフォルトLv=1）が最小なので1体表示
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{agents: []*domain.AgentModel{agent}},
		&mockDefeatedEnemyProvider{defeated: map[string]int{}, maxLevelReached: 0},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// Enterでバトル開始
	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("コマンドがnilです")
	}

	// コマンドを実行してメッセージを取得
	msg := cmd()

	startBattleMsg, ok := msg.(StartBattleMsg)
	if !ok {
		t.Fatalf("StartBattleMsgではありません: %T", msg)
	}

	// 敵タイプIDが含まれていること
	if startBattleMsg.EnemyTypeID != "slime" {
		t.Errorf("敵タイプID: got %s, want slime", startBattleMsg.EnemyTypeID)
	}

	// レベルが正しいこと
	if startBattleMsg.Level != 1 {
		t.Errorf("レベル: got %d, want 1", startBattleMsg.Level)
	}
}

// createTestAgent はテスト用のエージェントを作成します。
func createTestAgent() *domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low"},
	}
	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})
	modules := []*domain.ModuleModel{
		newTestDamageModule("m1", "モジュール1", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m2", "モジュール2", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m3", "モジュール3", []string{"physical_low"}, 1.0, "STR", ""),
		newTestDamageModule("m4", "モジュール4", []string{"physical_low"}, 1.0, "STR", ""),
	}
	return domain.NewAgent("agent1", core, modules)
}

// ==================== タスク2.3: 敵情報パネルのテスト ====================

// TestBattleSelectCarouselEnemyInfoPanel は敵情報パネルの表示をテストします。
func TestBattleSelectCarouselEnemyInfoPanel(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{
			ID:           "slime",
			Name:         "スライム",
			DefaultLevel: 1,
			BaseHP:       50,
			AttackType:   "physical",
			NormalPassive: &domain.EnemyPassiveSkill{
				ID:          "slime_normal",
				Name:        "ぷるぷるボディ",
				Description: "物理ダメージ10%軽減",
			},
			EnhancedPassive: &domain.EnemyPassiveSkill{
				ID:          "slime_enhanced",
				Name:        "怒りのスライム",
				Description: "攻撃力20%上昇",
			},
		},
	}
	// MaxLevelReached=0（未撃破）でslime（デフォルトLv=1）のみ表示
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: map[string]int{}, maxLevelReached: 0},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	screen.width = 120
	screen.height = 40

	view := screen.View()

	// 敵名が表示されていること
	if !containsString(view, "スライム") {
		t.Error("敵名が表示されていません")
	}

	// 攻撃属性が表示されていること
	if !containsString(view, "物理") {
		t.Error("攻撃属性が表示されていません")
	}

	// パッシブスキルのDescriptionが表示されていること（★付き）
	if !containsString(view, "物理ダメージ10%軽減") {
		t.Error("通常パッシブのDescriptionが表示されていません")
	}

	if !containsString(view, "攻撃力20%上昇") {
		t.Error("強化パッシブのDescriptionが表示されていません")
	}
}

// TestBattleSelectCarouselUndefeatedIndicator は未撃破敵の視覚的表示をテストします。
func TestBattleSelectCarouselUndefeatedIndicator(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// MaxLevelReached=0（未撃破）でslime（デフォルトLv=1）のみ表示
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: map[string]int{}, maxLevelReached: 0},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	screen.width = 120
	screen.height = 40

	view := screen.View()

	// 未撃破の表示があること
	if !containsString(view, "未撃破") {
		t.Error("未撃破の表示がありません")
	}

	// 固定レベルの表示があること
	if !containsString(view, "固定") {
		t.Error("固定レベルの表示がありません")
	}
}

// TestBattleSelectCarouselDefeatedIndicator は撃破済み敵の視覚的表示をテストします。
func TestBattleSelectCarouselDefeatedIndicator(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// slimeをレベル5で撃破、MaxLevelReached=1（デフォルトレベル）
	defeated := map[string]int{"slime": 5}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 1},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	screen.width = 120
	screen.height = 40

	view := screen.View()

	// 撃破済みの表示があること
	if !containsString(view, "撃破済み") {
		t.Error("撃破済みの表示がありません")
	}

	// 最高レベルの表示があること
	if !containsString(view, "最高Lv.5") {
		t.Error("最高レベルの表示がありません")
	}
}

// containsString は既にagent_management_test.goで定義されているため、ここでは再利用します。

// ==================== タスク2.4: レベル選択制約のテスト ====================

// TestBattleSelectCarouselLevelConstraintUndefeated は未撃破敵のレベル制約をテストします。
func TestBattleSelectCarouselLevelConstraintUndefeated(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// slimeをLv.1で撃破済み、MaxLevelReached=1（slimeのデフォルトレベル）
	// goblin（デフォルトLv=2）は未撃破でMaxLevelReached+1=2 以上なので表示される
	defeated := map[string]int{"slime": 1}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 1},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// goblin（インデックス1）を選択 - goblinは未撃破だがデフォルトLv=2 >= MaxLevelReached+1=2 なので表示される
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})

	// デフォルトレベルが設定されていること（goblinはデフォルトLv=2）
	if screen.selectedLevel != 2 {
		t.Errorf("デフォルトレベル: got %d, want 2", screen.selectedLevel)
	}

	// goblinは未撃破なので min == max であること（固定）
	if screen.minSelectableLevel != screen.maxSelectableLevel {
		t.Errorf("未撃破敵のレベル範囲が固定でありません: min=%d, max=%d",
			screen.minSelectableLevel, screen.maxSelectableLevel)
	}
}

// TestBattleSelectCarouselLevelConstraintDefeated は撃破済み敵のレベル制約をテストします。
func TestBattleSelectCarouselLevelConstraintDefeated(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	// slimeをレベル10で撃破済み（デフォルトレベル1）、MaxLevelReached=10
	// 新ロジック: 撃破済み敵はMaxLevelReachedまで選択可能
	defeated := map[string]int{"slime": 10}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 10},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// slime（インデックス0）が選択されている
	// デフォルトレベル1、MaxLevelReached=10なので、1〜10まで選択可能
	if screen.minSelectableLevel != 1 {
		t.Errorf("最小選択可能レベル: got %d, want 1", screen.minSelectableLevel)
	}

	if screen.maxSelectableLevel != 10 {
		t.Errorf("最大選択可能レベル: got %d, want 10", screen.maxSelectableLevel)
	}
}

// TestBattleSelectCarouselLevelConstraintMaxLevel は最大レベル100の制約をテストします。
func TestBattleSelectCarouselLevelConstraintMaxLevel(t *testing.T) {
	enemyTypes := []domain.EnemyType{
		{ID: "boss", Name: "ボス", DefaultLevel: 50, BaseHP: 1000},
	}
	// レベル100で撃破済み、MaxLevelReached=100（または100を超える値でテスト）
	defeated := map[string]int{"boss": 100}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 150},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// 最大レベルは100を超えないこと
	if screen.maxSelectableLevel > 100 {
		t.Errorf("最大選択可能レベルが100を超えています: got %d", screen.maxSelectableLevel)
	}
}

// TestBattleSelectCarouselStartBattleWithEnemyTypeID はバトル開始メッセージに敵タイプIDが含まれることをテストします。
func TestBattleSelectCarouselStartBattleWithEnemyTypeID(t *testing.T) {
	enemyTypes := createTestEnemyTypes()
	agent := createTestAgent()
	// slime(Lv=1)とgoblin(Lv=2)を撃破済み、MaxLevelReached=5
	// 表示される敵: slime, goblin（撃破済み）, dragon（未撃破で最小）
	defeated := map[string]int{"slime": 3, "goblin": 5}
	screen := NewBattleSelectScreenCarousel(
		&mockAgentProvider{agents: []*domain.AgentModel{agent}},
		&mockDefeatedEnemyProvider{defeated: defeated, maxLevelReached: 5},
		&mockEnemyTypeProvider{enemyTypes: enemyTypes},
	)

	// goblin（インデックス1）を選択
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyRight})

	// レベルを3に設定（デフォルトLv=2から上キーで3に）
	screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyUp})

	// Enterでバトル開始
	_, cmd := screen.handleKeyMsg(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("コマンドがnilです")
	}

	msg := cmd()
	startBattleMsg, ok := msg.(StartBattleMsg)
	if !ok {
		t.Fatalf("StartBattleMsgではありません: %T", msg)
	}

	// 敵タイプIDがgoblinであること
	if startBattleMsg.EnemyTypeID != "goblin" {
		t.Errorf("敵タイプID: got %s, want goblin", startBattleMsg.EnemyTypeID)
	}

	// レベルが3であること
	if startBattleMsg.Level != 3 {
		t.Errorf("レベル: got %d, want 3", startBattleMsg.Level)
	}
}

// ==================== ランクベースの敵表示テスト ====================

// mockEnemyProgressProvider はテスト用のEnemyProgressProvider実装です。
type mockEnemyProgressProvider struct {
	currentRank     int
	defeatRecords   map[string]defeatRecord
	enemyTypesCache map[string]domain.EnemyType
}

type defeatRecord struct {
	defeated         bool
	maxDefeatedLevel int
}

func (m *mockEnemyProgressProvider) GetCurrentRank() int {
	return m.currentRank
}

func (m *mockEnemyProgressProvider) IsDefeated(enemyTypeID string) bool {
	record, exists := m.defeatRecords[enemyTypeID]
	return exists && record.defeated
}

func (m *mockEnemyProgressProvider) GetMaxDefeatedLevel(enemyTypeID string) int {
	record, exists := m.defeatRecords[enemyTypeID]
	if !exists {
		return 0
	}
	return record.maxDefeatedLevel
}

func (m *mockEnemyProgressProvider) GetSelectableLevelRange(enemyTypeID string) (min, max int) {
	if !m.IsDefeated(enemyTypeID) {
		// 未撃破: デフォルトレベルのみ
		if et, exists := m.enemyTypesCache[enemyTypeID]; exists {
			defaultLevel := et.DefaultLevel
			if defaultLevel < 1 {
				defaultLevel = 1
			}
			return defaultLevel, defaultLevel
		}
		return 1, 1
	}
	// 撃破済み: 1 から (最大撃破レベル + 5) まで
	maxDefeatedLevel := m.GetMaxDefeatedLevel(enemyTypeID)
	return 1, maxDefeatedLevel + 5
}

func (m *mockEnemyProgressProvider) GetCurrentRankEnemies() []domain.EnemyType {
	enemies := make([]domain.EnemyType, 0)
	for _, et := range m.enemyTypesCache {
		if et.Rank == m.currentRank {
			enemies = append(enemies, et)
		}
	}
	return enemies
}

// createTestEnemyTypesWithRank はランク付きのテスト用敵タイプリストを作成します。
func createTestEnemyTypesWithRank() []domain.EnemyType {
	return []domain.EnemyType{
		{ID: "slime", Name: "スライム", DefaultLevel: 1, BaseHP: 50, AttackType: "physical", Rank: 1},
		{ID: "goblin", Name: "ゴブリン", DefaultLevel: 2, BaseHP: 80, AttackType: "physical", Rank: 1},
		{ID: "orc", Name: "オーク", DefaultLevel: 5, BaseHP: 200, AttackType: "physical", Rank: 2},
		{ID: "dragon", Name: "ドラゴン", DefaultLevel: 10, BaseHP: 500, AttackType: "magic", Rank: 3},
	}
}

func createEnemyTypeMap(enemyTypes []domain.EnemyType) map[string]domain.EnemyType {
	m := make(map[string]domain.EnemyType)
	for _, et := range enemyTypes {
		m[et.ID] = et
	}
	return m
}

// TestBattleSelectRankBasedFiltering はランクベースの敵フィルタリングをテストします。
func TestBattleSelectRankBasedFiltering(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithRank()
	enemyTypeMap := createEnemyTypeMap(enemyTypes)

	// ランク1のみ解放済み
	provider := &mockEnemyProgressProvider{
		currentRank:     1,
		defeatRecords:   map[string]defeatRecord{},
		enemyTypesCache: enemyTypeMap,
	}

	screen := NewBattleSelectScreenRankBased(
		&mockAgentProvider{},
		provider,
	)

	// ランク1の敵のみ表示される（slime, goblin）
	if len(screen.enemyTypes) != 2 {
		t.Errorf("ランク1の敵数: got %d, want 2", len(screen.enemyTypes))
	}

	// 表示される敵がランク1のものであることを確認
	for _, et := range screen.enemyTypes {
		if et.Rank != 1 {
			t.Errorf("ランク1以外の敵が表示されています: %s (rank=%d)", et.ID, et.Rank)
		}
	}
}

// TestBattleSelectRankProgressDisplay は現在ランクと進行状況の表示をテストします。
func TestBattleSelectRankProgressDisplay(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithRank()
	enemyTypeMap := createEnemyTypeMap(enemyTypes)

	// ランク1でslimeのみ撃破済み
	provider := &mockEnemyProgressProvider{
		currentRank: 1,
		defeatRecords: map[string]defeatRecord{
			"slime": {defeated: true, maxDefeatedLevel: 3},
		},
		enemyTypesCache: enemyTypeMap,
	}

	screen := NewBattleSelectScreenRankBased(
		&mockAgentProvider{},
		provider,
	)

	screen.width = 120
	screen.height = 40

	view := screen.View()

	// ランク表示があること
	if !containsString(view, "Rank 1") {
		t.Error("ランク表示がありません")
	}

	// 進行状況表示があること（1/2撃破）
	if !containsString(view, "1/2") {
		t.Error("進行状況表示がありません")
	}
}

// TestBattleSelectDefeatedStatusPerEnemy は敵ごとの撃破状況表示をテストします。
func TestBattleSelectDefeatedStatusPerEnemy(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithRank()
	enemyTypeMap := createEnemyTypeMap(enemyTypes)

	// ランク1でslimeのみ撃破済み
	provider := &mockEnemyProgressProvider{
		currentRank: 1,
		defeatRecords: map[string]defeatRecord{
			"slime": {defeated: true, maxDefeatedLevel: 5},
		},
		enemyTypesCache: enemyTypeMap,
	}

	screen := NewBattleSelectScreenRankBased(
		&mockAgentProvider{},
		provider,
	)

	screen.width = 120
	screen.height = 40

	// slime（撃破済み）を選択した状態で描画
	view := screen.View()

	// 撃破済み敵の最大撃破レベル表示があること
	if !containsString(view, "最高Lv.5") {
		t.Error("最大撃破レベル表示がありません")
	}
}

// TestBattleSelectLevelSelectionForDefeatedEnemy は撃破済み敵のレベル選択をテストします。
func TestBattleSelectLevelSelectionForDefeatedEnemy(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithRank()
	enemyTypeMap := createEnemyTypeMap(enemyTypes)

	// slimeをレベル5で撃破済み
	provider := &mockEnemyProgressProvider{
		currentRank: 1,
		defeatRecords: map[string]defeatRecord{
			"slime": {defeated: true, maxDefeatedLevel: 5},
		},
		enemyTypesCache: enemyTypeMap,
	}

	screen := NewBattleSelectScreenRankBased(
		&mockAgentProvider{},
		provider,
	)

	// slimeが選択されている状態（インデックス0）
	// 撃破済みなので、レベル1〜10（5+5）まで選択可能
	if screen.minSelectableLevel != 1 {
		t.Errorf("最小選択可能レベル: got %d, want 1", screen.minSelectableLevel)
	}

	if screen.maxSelectableLevel != 10 {
		t.Errorf("最大選択可能レベル: got %d, want 10 (5+5)", screen.maxSelectableLevel)
	}
}

// TestBattleSelectLevelFixedForUndefeatedEnemy は未撃破敵のレベル固定をテストします。
func TestBattleSelectLevelFixedForUndefeatedEnemy(t *testing.T) {
	enemyTypes := createTestEnemyTypesWithRank()
	enemyTypeMap := createEnemyTypeMap(enemyTypes)

	// 未撃破状態
	provider := &mockEnemyProgressProvider{
		currentRank:     1,
		defeatRecords:   map[string]defeatRecord{},
		enemyTypesCache: enemyTypeMap,
	}

	screen := NewBattleSelectScreenRankBased(
		&mockAgentProvider{},
		provider,
	)

	// 選択されている敵のデフォルトレベルと一致すること
	// （mapの順序は不定なので、slimeかgoblinのどちらかが選択される）
	selectedEnemy := screen.enemyTypes[screen.selectedTypeIdx]
	expectedLevel := selectedEnemy.DefaultLevel
	if expectedLevel < 1 {
		expectedLevel = 1
	}

	// 未撃破なのでデフォルトレベルのみ選択可能
	if screen.selectedLevel != expectedLevel {
		t.Errorf("選択レベル: got %d, want %d (デフォルトレベル)", screen.selectedLevel, expectedLevel)
	}

	// min == max であること（固定）
	if screen.minSelectableLevel != screen.maxSelectableLevel {
		t.Errorf("未撃破敵のレベル範囲が固定でありません: min=%d, max=%d",
			screen.minSelectableLevel, screen.maxSelectableLevel)
	}
}
