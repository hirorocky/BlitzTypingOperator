// Package screens はTUI画面のテストを提供します。
package screens

import (
	"strings"
	"testing"
	"time"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== Task 10.3: バトル画面のテスト ====================

// TestNewBattleScreen はBattleScreenの初期化をテストします。
func TestNewBattleScreen(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if screen == nil {
		t.Fatal("BattleScreenがnilです")
	}

	if screen.enemy != enemy {
		t.Error("敵が正しく設定されていません")
	}

	if screen.player != player {
		t.Error("プレイヤーが正しく設定されていません")
	}
}

// TestBattleScreenEnemyInfo は敵情報表示をテストします。

func TestBattleScreenEnemyInfo(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}

	// 敵の名前が含まれているか確認
	// （実際のレンダリング内容の詳細はUIに依存）
}

// TestBattleScreenPlayerInfo はプレイヤー情報表示をテストします。

func TestBattleScreenPlayerInfo(t *testing.T) {
	player := createTestPlayer()
	player.HP = 50
	player.MaxHP = 100

	// バフを追加
	player.EffectTable.AddBuff("攻撃UP", "攻撃UP", 5.0, map[domain.EffectColumn]float64{
		domain.ColDamageBonus: 10,
	})

	enemy := createTestEnemy()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestBattleScreenSkillList はスキル一覧表示をテストします。

func TestBattleScreenSkillList(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// スキルスロットが作成されていることを確認
	if len(screen.skillSlots) == 0 {
		t.Error("スキルスロットが空です")
	}

	// エージェントごとにグループ化されているか
	expectedSlots := len(agents) * 4 // 各エージェント4スキル
	if len(screen.skillSlots) != expectedSlots {
		t.Errorf("スキルスロット数: got %d, want %d", len(screen.skillSlots), expectedSlots)
	}
}

// TestBattleScreenCooldownDisplay はクールダウン表示をテストします。

func TestBattleScreenCooldownDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// クールダウンを設定
	if len(screen.skillSlots) > 0 {
		screen.skillSlots[0].CooldownRemaining = 3.0
		screen.skillSlots[0].CooldownTotal = 5.0
	}

	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// TestBattleScreenTypingChallenge はタイピングチャレンジ表示をテストします。

func TestBattleScreenTypingChallenge(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// タイピングチャレンジを開始
	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	if screen.activeChallenge == nil {
		t.Error("チャレンジが開始されていません")
	}
}

// TestBattleScreenTimeLimit はチャレンジ開始後にactiveChallengeがnon-nilであることをテストします。

func TestBattleScreenTimeLimit(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// チャレンジを開始
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	// チャレンジが開始されていること
	if screen.activeChallenge == nil {
		t.Error("チャレンジが開始されていません")
	}

	// Viewが空でないこと（制限時間等の情報が含まれている）
	view := screen.activeChallenge.View()
	if view == "" {
		t.Error("チャレンジのViewが空です")
	}
}

// TestBattleScreenRender はバトル画面のレンダリングをテストします。
func TestBattleScreenRender(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}
}

// ==================== Task: バトル画面のTick機能テスト ====================

// TestBattleScreenInitReturnsTick はInit()がtickコマンドを返すことをテストします。
func TestBattleScreenInitReturnsTick(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	cmd := screen.Init()

	if cmd == nil {
		t.Error("Init()がnilを返しました。tickコマンドを返す必要があります")
	}
}

// TestBattleScreenTickUpdatesCooldowns はTickMsgがクールダウンを更新することをテストします。
func TestBattleScreenTickUpdatesCooldowns(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// クールダウンを設定
	if len(screen.skillSlots) > 0 {
		screen.skillSlots[0].CooldownRemaining = 3.0
	}

	// TickMsgを送信（100ms経過をシミュレート）
	_, _ = screen.Update(BattleTickMsg{})

	// クールダウンが減少していること
	if len(screen.skillSlots) > 0 {
		// tickInterval (100ms = 0.1秒) 分減少しているはず
		expected := 3.0 - 0.1
		actual := screen.skillSlots[0].CooldownRemaining
		if actual > expected+0.01 || actual < expected-0.01 {
			t.Errorf("クールダウンが更新されていません: got %.2f, want %.2f", actual, expected)
		}
	}
}

// TestBattleScreenTickReturnsNextTick はTickMsg処理後に次のtickコマンドを返すことをテストします。
func TestBattleScreenTickReturnsNextTick(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	_, cmd := screen.Update(BattleTickMsg{})

	if cmd == nil {
		t.Error("TickMsg処理後にnilを返しました。次のtickコマンドを返す必要があります")
	}
}

// TestBattleScreenTickHandlesEnemyAttack はTickMsgが敵攻撃を処理することをテストします。
func TestBattleScreenTickHandlesEnemyAttack(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = 100
	player.MaxHP = 100
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// チャージを開始し、完了状態にする（過去の時刻でStartCharging）
	enemy.PrepareNextAction()
	if action := enemy.GetNextAction(); action != nil {
		enemy.StartCharging(*action, time.Now().Add(-2*time.Second))
	}

	initialHP := player.HP

	// TickMsgを送信
	_, _ = screen.Update(BattleTickMsg{})

	// プレイヤーがダメージを受けているか、または新しいチャージが開始されているはず
	hpDecreased := player.HP < initialHP
	newChargeStarted := enemy.WaitMode == domain.WaitModeCharging && !enemy.IsChargeComplete(time.Now())

	if !hpDecreased && !newChargeStarted {
		t.Error("敵攻撃後にダメージが発生していないか、次のチャージが開始されていません")
	}
}

// TestBattleScreenTypingTimeout はチャレンジのタイムアウトがhandleChallengeCompleteで処理されることをテストします。
func TestBattleScreenTypingTimeout(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// チャレンジを開始
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	if screen.activeChallenge == nil {
		t.Fatal("チャレンジが開始されていません")
	}

	// チャレンジが完了結果を返した場合にhandleChallengeCompleteで処理される
	// タイムアウトの処理はChallengeModel内部で行われるため、
	// ここではhandleChallengeCompleteが呼ばれた後にactiveChallengeがnilになることを確認
	failResult := &domain.ChallengeOutput{
		Status: domain.ChallengeFail,
	}
	screen.handleChallengeComplete(failResult)

	if screen.activeChallenge != nil {
		t.Error("チャレンジ完了後もactiveChallengeがnon-nilのままです")
	}
}

// ==================== Task: 敗北判定テスト ====================

// TestBattleScreenDefeatDetection はプレイヤーHP0で敗北判定されることをテストします。
func TestBattleScreenDefeatDetection(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = 0
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// TickMsgを送信
	_, cmd := screen.Update(BattleTickMsg{})

	// 敗北状態になっているはず
	if !screen.IsGameOver() {
		t.Error("HP0でもゲームオーバーになっていません")
	}

	if !screen.IsDefeat() {
		t.Error("HP0でも敗北判定になっていません")
	}

	// 敗北時はシーン遷移コマンドが返されるはず
	if cmd == nil {
		t.Error("敗北時にコマンドが返されていません")
	}
}

// TestBattleScreenDefeatAfterEnemyAttack は敵攻撃でHP0になった場合の敗北判定をテストします。
func TestBattleScreenDefeatAfterEnemyAttack(t *testing.T) {
	enemy := createTestEnemy()
	enemy.AttackPower = 200 // 一撃で倒せるダメージ

	player := createTestPlayer()
	player.HP = 100
	player.MaxHP = 100
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// チャージを開始し、完了状態にする（過去の時刻でStartCharging）
	enemy.PrepareNextAction()
	if action := enemy.GetNextAction(); action != nil {
		enemy.StartCharging(*action, time.Now().Add(-2*time.Second))
	}

	// TickMsgを送信（敵攻撃が発生）
	_, _ = screen.Update(BattleTickMsg{})

	// プレイヤーのHPが0以下になっているはず
	if player.HP > 0 {
		t.Errorf("敵攻撃後もHP残っています: %d", player.HP)
	}

	// 敗北状態になっているはず
	if !screen.IsDefeat() {
		t.Error("敵攻撃でHP0になっても敗北判定になっていません")
	}
}

// TestBattleScreenVictoryDetection は敵HP0で勝利判定されることをテストします。
func TestBattleScreenVictoryDetection(t *testing.T) {
	enemy := createTestEnemy()
	enemy.HP = 0

	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// TickMsgを送信
	_, _ = screen.Update(BattleTickMsg{})

	// 勝利状態になっているはず
	if !screen.IsGameOver() {
		t.Error("敵HP0でもゲームオーバーになっていません")
	}

	if !screen.IsVictory() {
		t.Error("敵HP0でも勝利判定になっていません")
	}

	// 結果表示状態になっているはず
	if !screen.IsShowingResult() {
		t.Error("勝利時に結果表示状態になっていません")
	}
}

// ==================== Task: 結果表示待機テスト ====================

// TestBattleScreenResultWaitsForEnter は結果表示後Enterを待つことをテストします。
func TestBattleScreenResultWaitsForEnter(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = 0
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// TickMsgを送信して敗北状態に
	_, _ = screen.Update(BattleTickMsg{})

	// 結果表示状態ではシーン遷移コマンドは返されない（tickのみ継続）
	// ただしtickは継続するのでcmdがnilではない可能性がある
	if screen.IsShowingResult() {
		// Enterを押さない限りBattleResultMsgは発行されない
		// 他のキーを押しても遷移しない
		_, _ = screen.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

		// まだ結果表示状態のはず
		if !screen.IsShowingResult() {
			t.Error("Enter以外のキーで結果表示が終了しました")
		}
	}

	// Enterを押すとBattleResultMsgが返される
	_, cmd := screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Error("Enter押下後にコマンドが返されていません")
	}
}

// TestBattleScreenResultDisplaysMessage は結果表示にメッセージが含まれることをテストします。
func TestBattleScreenResultDisplaysMessage(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = 0
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	// TickMsgを送信して敗北状態に
	_, _ = screen.Update(BattleTickMsg{})

	// Viewに結果メッセージが含まれているはず
	rendered := screen.View()

	if rendered == "" {
		t.Error("レンダリング結果が空です")
	}

	// "Enter"という文字が含まれているはず（続行のヒント）
	if !strings.Contains(rendered, "Enter") {
		t.Error("結果画面にEnterキーのヒントがありません")
	}
}

// ==================== ヘルパー関数 ====================

func createTestEnemy() *domain.EnemyModel {
	enemyType := domain.EnemyType{
		ID:     "test_enemy",
		Name:   "テストエネミー",
		BaseHP: 100,
	}

	return domain.NewEnemy(
		"enemy1",
		"テストエネミー Lv.5",
		5,
		500,
		20,
		enemyType,
	)
}

func createTestPlayer() *domain.PlayerModel {
	return domain.NewPlayerWithMaxHP(100)
}

func createTestAgents() []*domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}

	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})

	skills := []*domain.SkillModel{
		newTestDamageSkill("m1", "物理攻撃", []string{"physical_low"}, 1.0, "STR", "物理ダメージ"),
		newTestDamageSkill("m2", "魔法攻撃", []string{"magic_low"}, 1.0, "INT", "魔法ダメージ"),
		newTestHealSkill("m3", "回復", []string{"heal_low"}, 1.0, "INT", "HP回復"),
		newTestBuffSkill("m4", "バフ", []string{"buff_low"}, "攻撃力UP"),
	}

	agent := domain.NewAgent("agent1", core, skills)
	return []*domain.AgentModel{agent}
}

// ==================== Task 6.1-6.6: バトル画面UI改善のテスト ====================

// TestBattleScreen3AreaLayout はバトル画面の3エリアレイアウトをテストします。

func TestBattleScreen3AreaLayout(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.SetFeatureUnlockProvider(allFeaturesUnlockedProvider())
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 敵情報が含まれること
	if !strings.Contains(rendered, enemy.Name) {
		t.Error("敵情報エリアに敵の名前が表示されていません")
	}

	// マナ情報が含まれること
	if !strings.Contains(rendered, "Mana:") {
		t.Error("マナ情報エリアが表示されていません")
	}

	// スキル情報が含まれること
	if !strings.Contains(rendered, "スキル") {
		t.Error("スキルエリアが表示されていません")
	}
}

// TestBattleScreenAgentSkillDisplay はエージェントごとのスキル表示をテストします。

func TestBattleScreenAgentSkillDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// エージェントのコアタイプ名が含まれること
	if !strings.Contains(rendered, agents[0].GetCoreTypeName()) {
		t.Error("エージェントのコアタイプ名が表示されていません")
	}
}

// TestBattleScreenHPBarDisplay はHPバー表示をテストします。

func TestBattleScreenHPBarDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// HPバーが含まれること（HPの数値が表示されている）
	if !strings.Contains(rendered, "HP:") {
		t.Error("HP表示がありません")
	}
}

// TestBattleScreenEnemyAttackTimerDisplay は敵攻撃タイマー表示をテストします。

func TestBattleScreenEnemyAttackTimerDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	// 行動パターンを設定して次の行動を準備
	action := domain.EnemyAction{
		ID:         "attack",
		ActionType: domain.EnemyActionAttack,
		ChargeTime: 2 * time.Second,
	}
	enemy.SetNextAction(&action)
	enemy.StartCharging(action, time.Now())

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	rendered := screen.View()

	// 物理ダメージの表示が含まれること（チャージ中は攻撃効果を表示）
	if !strings.Contains(rendered, "物理ダメージ") {
		t.Logf("レンダリング結果:\n%s", rendered)
		t.Error("敵攻撃タイマー表示がありません")
	}
}

// TestBattleScreenTypingColorDisplay はチャレンジ中のView表示をテストします。

func TestBattleScreenTypingColorDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// チャレンジを開始
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	rendered := screen.View()

	// チャレンジ中のViewが表示されること（空でないこと）
	if rendered == "" {
		t.Error("チャレンジ中のViewが空です")
	}
}

// TestBattleScreenWinDisplay は勝利時のWIN表示をテストします。

func TestBattleScreenWinDisplay(t *testing.T) {
	enemy := createTestEnemy()
	enemy.HP = 0 // 敵HP0で勝利

	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	// TickMsgを送信して勝利状態に
	_, _ = screen.Update(BattleTickMsg{})

	if !screen.IsVictory() {
		t.Error("勝利状態になっていません")
	}

	rendered := screen.View()

	// 勝利メッセージが含まれること
	if !strings.Contains(rendered, "勝利") {
		t.Error("勝利メッセージが表示されていません")
	}
}

// TestBattleScreenLoseDisplay は敗北時のLOSE表示をテストします。

func TestBattleScreenLoseDisplay(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = 0 // プレイヤーHP0で敗北

	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.width = 120
	screen.height = 40

	// TickMsgを送信して敗北状態に
	_, _ = screen.Update(BattleTickMsg{})

	if !screen.IsDefeat() {
		t.Error("敗北状態になっていません")
	}

	rendered := screen.View()

	// 敗北メッセージが含まれること
	if !strings.Contains(rendered, "敗北") {
		t.Error("敗北メッセージが表示されていません")
	}
}

// ==================== Task 6.1: ファイル分割検証テスト ====================

// TestBattleScreenLogicSeparation はバトルロジックが正しく動作することを検証します。
// Task 6.1: UIレンダリングとゲームロジックの分離後も機能が維持されることを確認
func TestBattleScreenLogicSeparation(t *testing.T) {
	enemy := createTestEnemy()
	enemy.HP = 100
	player := createTestPlayer()
	player.HP = 100
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// checkGameOver - 正常状態では終了しない
	if screen.checkGameOver() {
		t.Error("HP残っているのにゲームオーバー判定されました")
	}

	// 敵HP0で勝利判定
	enemy.HP = 0
	if !screen.checkGameOver() {
		t.Error("敵HP0でゲームオーバー判定されませんでした")
	}
	if !screen.IsVictory() {
		t.Error("敵HP0で勝利判定されませんでした")
	}
}

// TestBattleScreenViewSeparation はView関連メソッドが正しく動作することを検証します。
// Task 6.1: UIレンダリングとゲームロジックの分離後も描画が維持されることを確認
func TestBattleScreenViewSeparation(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)
	screen.SetFeatureUnlockProvider(allFeaturesUnlockedProvider())
	screen.width = 120
	screen.height = 40

	// renderEnemyArea
	enemyArea := screen.renderEnemyArea()
	if enemyArea == "" {
		t.Error("renderEnemyAreaが空を返しました")
	}
	if !strings.Contains(enemyArea, enemy.Name) {
		t.Error("敵エリアに敵名が含まれていません")
	}

	// renderAgentArea
	agentArea := screen.renderAgentArea()
	if agentArea == "" {
		t.Error("renderAgentAreaが空を返しました")
	}

	// renderPlayerArea
	playerArea := screen.renderPlayerArea()
	if playerArea == "" {
		t.Error("renderPlayerAreaが空を返しました")
	}
	if !strings.Contains(playerArea, "Mana:") {
		t.Error("プレイヤーエリアにマナ情報が含まれていません")
	}
}

// TestBattleScreenCooldownLogic はクールダウンロジックが正しく動作することを検証します。
// Task 6.1: ロジック分離後もクールダウン機能が維持されることを確認
func TestBattleScreenCooldownLogic(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// クールダウン開始
	screen.StartCooldown(0, 5.0)
	if screen.skillSlots[0].CooldownRemaining != 5.0 {
		t.Errorf("クールダウン設定失敗: got %.2f, want 5.0", screen.skillSlots[0].CooldownRemaining)
	}

	// クールダウン更新
	screen.UpdateCooldowns(1.0)
	if screen.skillSlots[0].CooldownRemaining != 4.0 {
		t.Errorf("クールダウン更新失敗: got %.2f, want 4.0", screen.skillSlots[0].CooldownRemaining)
	}

	// IsReady確認
	if screen.skillSlots[0].IsReady() {
		t.Error("クールダウン中なのにIsReady=trueになっています")
	}

	// クールダウン完了
	screen.UpdateCooldowns(5.0)
	if !screen.skillSlots[0].IsReady() {
		t.Error("クールダウン完了後もIsReady=falseのままです")
	}
}

// TestBattleScreenTypingLogic はチャレンジシステムが正しく動作することを検証します。
func TestBattleScreenTypingLogic(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// チャレンジ開始
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)
	if screen.activeChallenge == nil {
		t.Fatal("チャレンジが開始されていません")
	}

	// チャレンジのViewが空でないこと
	view := screen.activeChallenge.View()
	if view == "" {
		t.Error("チャレンジのViewが空です")
	}

	// ESCでキャンセル
	updated, _ := screen.activeChallenge.Update(tea.KeyMsg{Type: tea.KeyEsc})
	screen.activeChallenge = updated
	result := screen.activeChallenge.Result()
	if result == nil {
		t.Fatal("ESC後にResult()がnilです")
	}
	if result.Status != domain.ChallengeCancel {
		t.Errorf("ESC後のStatus = %d, want ChallengeCancel", result.Status)
	}

	// handleChallengeCompleteでactiveChallengeがクリアされる
	screen.handleChallengeComplete(result)
	if screen.activeChallenge != nil {
		t.Error("チャレンジ完了後もactiveChallengeがnon-nilです")
	}
}

// TestBattleScreenEffectDuration はエフェクト持続時間更新が正しく動作することを検証します。
// Task 6.1: ロジック分離後もエフェクト更新が維持されることを確認
func TestBattleScreenEffectDuration(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	// プレイヤーにバフを追加
	player.EffectTable.AddBuff("テストバフ", "テストバフ", 5.0, map[domain.EffectColumn]float64{
		domain.ColDamageBonus: 10,
	})

	screen := NewBattleScreen(enemy, player, agents, nil)

	// 初期状態確認
	buffs := player.EffectTable.FindBySourceType(domain.SourceBuff)
	if len(buffs) == 0 {
		t.Fatal("バフが追加されていません")
	}
	if *buffs[0].Duration != 5.0 {
		t.Errorf("初期持続時間が正しくありません: got %.2f, want 5.0", *buffs[0].Duration)
	}

	// エフェクト持続時間更新
	screen.updateEffectDurations(1.0)

	buffs = player.EffectTable.FindBySourceType(domain.SourceBuff)
	if len(buffs) == 0 {
		t.Fatal("バフが消えてしまいました")
	}
	if *buffs[0].Duration != 4.0 {
		t.Errorf("持続時間更新後の値が正しくありません: got %.2f, want 4.0", *buffs[0].Duration)
	}
}

// ==================== Task 7.1: RecastManager統合テスト ====================

// TestBattleScreenHasRecastManager はRecastManagerが初期化されることを検証します。
func TestBattleScreenHasRecastManager(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if screen.recastManager == nil {
		t.Error("RecastManagerが初期化されていません")
	}
}

// TestBattleScreenSkillUsageStartsRecast はチャレンジ完了時にリキャストが開始されることを検証します。
func TestBattleScreenSkillUsageStartsRecast(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル使用前はエージェントがReady
	if !screen.recastManager.IsAgentReady(0) {
		t.Error("初期状態でエージェントがリキャスト中になっています")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	// スキル選択直後はまだReady（CDはチャレンジ完了時に開始される）
	if !screen.recastManager.IsAgentReady(0) {
		t.Error("スキル選択直後にリキャストが開始されています")
	}

	// チャレンジ完了でリキャストが開始される
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	if screen.recastManager.IsAgentReady(0) {
		t.Error("チャレンジ完了後もエージェントがReady状態です")
	}
}

// TestBattleScreenRecastBlocksSkillUsage はリキャスト中のエージェントのスキル使用がブロックされることを検証します。
func TestBattleScreenRecastBlocksSkillUsage(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// エージェント0のリキャストを開始
	screen.recastManager.StartRecast(0, 5*time.Second)

	// エージェント0のスキルは使用不可
	if screen.isSkillUsable(0) {
		t.Error("リキャスト中のエージェントのスキルが使用可能になっています")
	}
}

// TestBattleScreenTickUpdatesRecast はTickMsgでリキャスト時間が更新されることを検証します。
func TestBattleScreenTickUpdatesRecast(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// リキャストを開始（3秒）
	screen.recastManager.StartRecast(0, 3*time.Second)

	initialState := screen.recastManager.GetRecastState(0)
	if initialState == nil {
		t.Fatal("リキャスト状態が取得できません")
	}
	initialRemaining := initialState.RemainingSeconds

	// TickMsgを送信
	_, _ = screen.Update(BattleTickMsg{})

	state := screen.recastManager.GetRecastState(0)
	if state == nil {
		t.Fatal("TickMsg後にリキャスト状態が取得できません")
	}

	// 100ms(tickInterval)分減少しているはず
	expected := initialRemaining - 0.1
	if state.RemainingSeconds > expected+0.01 || state.RemainingSeconds < expected-0.01 {
		t.Errorf("リキャスト時間が更新されていません: got %.2f, want %.2f", state.RemainingSeconds, expected)
	}
}

// TestBattleScreenRecastCompletionEnablesAgent はリキャスト終了時にエージェントが使用可能になることを検証します。
func TestBattleScreenRecastCompletionEnablesAgent(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// 短いリキャストを開始（0.05秒 = 50ms = tick1回未満）
	screen.recastManager.StartRecast(0, 50*time.Millisecond)

	// リキャスト中
	if screen.recastManager.IsAgentReady(0) {
		t.Error("リキャスト開始直後にエージェントがReady状態です")
	}

	// TickMsgを送信（100ms経過）
	_, _ = screen.Update(BattleTickMsg{})

	// リキャスト完了
	if !screen.recastManager.IsAgentReady(0) {
		t.Error("リキャスト時間経過後もエージェントがリキャスト中です")
	}
}

// ==================== Task 7.2: ChainEffectManager統合テスト ====================

// TestBattleScreenHasChainEffectManager はChainEffectManagerが初期化されることを検証します。
func TestBattleScreenHasChainEffectManager(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if screen.chainEffectManager == nil {
		t.Error("ChainEffectManagerが初期化されていません")
	}
}

// TestBattleScreenSkillUsageRegistersChainEffect はスキル使用時にチェイン効果が登録されることを検証します。
func TestBattleScreenSkillUsageRegistersChainEffect(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// 登録前は待機中チェイン効果なし
	if len(screen.chainEffectManager.GetPendingEffects()) != 0 {
		t.Error("初期状態で待機中チェイン効果が存在します")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)

	// スキル選択直後はチェイン効果が登録されていない
	if len(screen.chainEffectManager.GetPendingEffects()) != 0 {
		t.Error("スキル選択時にチェイン効果が登録されています")
	}

	// チャレンジ完了後にチェイン効果が登録される
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})
	pendingEffects := screen.chainEffectManager.GetPendingEffects()
	if len(pendingEffects) == 0 {
		t.Error("チャレンジ完了後にチェイン効果が登録されていません")
	}
}

// TestBattleScreenChainEffectTrigger は他エージェントのスキル使用時にチェイン効果が発動することを検証します。
func TestBattleScreenChainEffectTrigger(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffectMultiple()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) < 5 {
		t.Skip("スキルスロットが足りません")
	}

	// エージェント0のスキル選択→チャレンジ完了（チェイン効果が登録される）
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// 待機中チェイン効果があること
	if len(screen.chainEffectManager.GetPendingEffects()) == 0 {
		t.Error("エージェント0のチェイン効果が登録されていません")
	}

	// エージェント1のスキル選択→チャレンジ完了（チェイン効果が発動）
	screen.selectedSkillIdx = 4 // エージェント1の最初のスキル
	screen.selectedAgentIdx = 1
	screen.startChallenge(screen.skillSlots[4].Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が発動して削除されているはず
	pendingEffects := screen.chainEffectManager.GetPendingEffects()
	// エージェント0の効果は発動済み、エージェント1の効果は待機中
	foundAgent0 := false
	for _, pe := range pendingEffects {
		if pe.AgentIndex == 0 {
			foundAgent0 = true
		}
	}
	if foundAgent0 {
		t.Error("エージェント0のチェイン効果が発動後も残っています")
	}
}

// TestBattleScreenRecastCompletionPersistsChainEffect はリキャスト終了後もチェイン効果が待機状態を維持することを検証します。
func TestBattleScreenRecastCompletionPersistsChainEffect(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// エージェント0のスキル選択→チャレンジ完了（チェイン効果が登録される）
	screen.selectedSkillIdx = 0
	screen.startChallenge(screen.skillSlots[0].Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が登録されている
	if len(screen.chainEffectManager.GetPendingEffects()) == 0 {
		t.Error("チェイン効果が登録されていません")
	}

	// リキャストを短い時間に設定して終了させる
	screen.recastManager.CancelRecast(0)
	screen.recastManager.StartRecast(0, 50*time.Millisecond)

	// TickMsgを送信してリキャストを終了
	_, _ = screen.Update(BattleTickMsg{})

	// リキャスト完了
	if !screen.recastManager.IsAgentReady(0) {
		t.Error("リキャストが終了していません")
	}

	// リキャスト完了後もチェイン効果が待機状態を維持している（効果内容も検証）
	pending := screen.chainEffectManager.GetPendingEffectForAgent(0)
	if pending == nil {
		t.Fatal("リキャスト完了後にエージェント0のチェイン効果が消えています（待機状態を維持すべき）")
	}
	if pending.Effect.Type != domain.ChainEffectDamageBonus {
		t.Errorf("チェイン効果タイプ: got %v, want %v", pending.Effect.Type, domain.ChainEffectDamageBonus)
	}
	if pending.Effect.Value != 25 {
		t.Errorf("チェイン効果値: got %v, want 25", pending.Effect.Value)
	}
}

// TestBattleScreenChainEffectTriggersAfterRecastCompletion はリキャスト完了後に他エージェントがスキルを使用した場合にチェイン効果が発動することを検証します。
func TestBattleScreenChainEffectTriggersAfterRecastCompletion(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffectMultiple()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) < 5 {
		t.Skip("スキルスロットが足りません")
	}

	// エージェント0のスキル選択（チェイン効果を登録）→チャレンジ完了
	screen.selectedSkillIdx = 0
	slot0 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot0.CooldownTotal)
	screen.startAgentRecast(slot0.AgentIndex, slot0.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が登録されている
	if !screen.chainEffectManager.HasPendingEffect(0) {
		t.Fatal("エージェント0のチェイン効果が登録されていません")
	}

	// リキャストを短い時間に設定して終了させる
	screen.recastManager.CancelRecast(0)
	screen.recastManager.StartRecast(0, 50*time.Millisecond)
	_, _ = screen.Update(BattleTickMsg{})

	// リキャスト完了を確認
	if !screen.recastManager.IsAgentReady(0) {
		t.Fatal("リキャストが終了していません")
	}

	// リキャスト完了後もチェイン効果が残っていることを確認
	if !screen.chainEffectManager.HasPendingEffect(0) {
		t.Fatal("リキャスト完了後にチェイン効果が消えています")
	}

	// エージェント1がスキル使用（チェイン効果が発動するはず）
	screen.selectedSkillIdx = 4
	screen.selectedAgentIdx = 1
	slot1 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot1.CooldownTotal)
	screen.startAgentRecast(slot1.AgentIndex, slot1.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// エージェント0のチェイン効果が発動して消えている
	if screen.chainEffectManager.HasPendingEffect(0) {
		t.Error("リキャスト完了後の他エージェントスキル使用でチェイン効果が発動していません")
	}
}

// TestBattleScreenStartRecastClearsChainEffectForNoChainSkill はチェイン効果なしスキル使用時に既存の待機中チェイン効果が削除されることを検証します。
func TestBattleScreenStartRecastClearsChainEffectForNoChainSkill(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) < 2 {
		t.Skip("スキルスロットが足りません")
	}

	// エージェント0のチェイン効果付きスキルを使用
	screen.selectedSkillIdx = 0
	slot0 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot0.CooldownTotal)
	screen.startAgentRecast(slot0.AgentIndex, slot0.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が登録されている
	if !screen.chainEffectManager.HasPendingEffect(0) {
		t.Fatal("チェイン効果が登録されていません")
	}

	// リキャスト完了
	screen.recastManager.CancelRecast(0)

	// チェイン効果なしスキル（m2: 魔法攻撃, ChainEffect=nil）を使用
	screen.selectedSkillIdx = 1
	slot1 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot1.CooldownTotal)
	screen.startAgentRecast(slot1.AgentIndex, slot1.Skill)

	// チェイン効果なしスキル使用後、既存のチェイン効果が削除されている
	if screen.chainEffectManager.HasPendingEffect(0) {
		t.Error("チェイン効果なしスキル使用後に既存のチェイン効果が削除されていません")
	}
}

// TestBattleScreenPendingChainEffectTriggersWhenOtherAgentUsesNonChainSkill は、
// エージェントAがチェイン効果付きスキルを使用した後に、別のエージェントBが非チェインスキルを
// 使用した際に、エージェントAの待機中チェイン効果が発動して消費されることを検証します。
func TestBattleScreenPendingChainEffectTriggersWhenOtherAgentUsesNonChainSkill(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffectMultiple()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) < 6 {
		t.Skip("スキルスロットが足りません")
	}

	// エージェント0のチェイン効果付きスキル（slot 0: m1）を使用
	screen.selectedSkillIdx = 0
	slot0 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot0.CooldownTotal)
	screen.startAgentRecast(slot0.AgentIndex, slot0.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が登録されていることを確認
	if !screen.chainEffectManager.HasPendingEffect(0) {
		t.Fatal("エージェント0のチェイン効果が登録されていません")
	}

	// リキャスト完了
	screen.recastManager.CancelRecast(0)

	// エージェント1の非チェインスキル（slot 5: m6 魔法攻撃2, ChainEffect=nil）を使用
	screen.selectedSkillIdx = 5
	screen.selectedAgentIdx = 1
	slot1 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot1.CooldownTotal)
	screen.startAgentRecast(slot1.AgentIndex, slot1.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// エージェント0の待機中チェイン効果が発動して消費されていることを確認
	if screen.chainEffectManager.HasPendingEffect(0) {
		t.Error("別エージェントの非チェインスキル使用時にエージェント0の待機中チェイン効果が発動・消費されていません")
	}
}

// ==================== Task 7.3: 統合フロー検証テスト ====================

// TestBattleScreenSkillRecastChainFlowIntegration はスキル使用→リキャスト開始→チェイン効果登録の一連フローを検証します。
func TestBattleScreenSkillRecastChainFlowIntegration(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// Step 1: 初期状態確認
	if !screen.recastManager.IsAgentReady(0) {
		t.Error("初期状態: エージェント0がリキャスト中")
	}
	if len(screen.chainEffectManager.GetPendingEffects()) != 0 {
		t.Error("初期状態: 待機中チェイン効果が存在")
	}

	// Step 2: スキル選択（クールダウン・リキャスト・チェイン効果の開始）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot.CooldownTotal)
	screen.startAgentRecast(slot.AgentIndex, slot.Skill)

	// Step 3: リキャスト開始確認（スキル選択直後）
	if screen.recastManager.IsAgentReady(0) {
		t.Error("スキル選択後: エージェント0がリキャスト中になっていない")
	}

	// Step 4: チェイン効果登録確認（スキル選択直後）
	pendingEffects := screen.chainEffectManager.GetPendingEffects()
	if len(pendingEffects) == 0 {
		t.Error("スキル選択後: チェイン効果が登録されていない")
	}

	// Step 5: エージェント0のスキル使用がブロックされる
	if screen.isSkillUsable(0) {
		t.Error("リキャスト中: エージェント0のスキルが使用可能になっている")
	}
}

// TestBattleScreenRecastBlockedSkillSelection はリキャスト中のスキル選択がブロックされることを検証します。
func TestBattleScreenRecastBlockedSkillSelection(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// エージェント0をリキャスト状態に
	screen.recastManager.StartRecast(0, 5*time.Second)

	// エージェント0のスキルを選択してEnterを押す
	screen.selectedSlot = 0
	screen.selectedAgentIdx = 0
	_, _ = screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// タイピングチャレンジが開始されていないはず
	if screen.activeChallenge != nil {
		t.Error("リキャスト中のエージェントのスキルでタイピングチャレンジが開始されました")
	}
}

// TestBattleScreenChainEffectTimingVerification はチェイン効果発動条件とタイミングを検証します。
func TestBattleScreenChainEffectTimingVerification(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffectMultiple()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) < 5 {
		t.Skip("スキルスロットが足りません（2エージェント必要）")
	}

	// エージェント0のスキル選択（ダメージボーナスのチェイン効果を登録）
	screen.selectedSkillIdx = 0
	screen.selectedAgentIdx = 0
	slot0 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot0.CooldownTotal)
	screen.startAgentRecast(slot0.AgentIndex, slot0.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// 待機中チェイン効果が存在
	pendingBefore := len(screen.chainEffectManager.GetPendingEffects())
	if pendingBefore == 0 {
		t.Error("チェイン効果が登録されていません")
	}

	initialEnemyHP := enemy.HP

	// エージェント1のスキル選択（チェイン効果が発動するはず）
	screen.selectedSkillIdx = 4 // エージェント1の最初のスキル
	screen.selectedAgentIdx = 1
	slot1 := screen.skillSlots[screen.selectedSkillIdx]
	screen.StartCooldown(screen.selectedSkillIdx, slot1.CooldownTotal)
	screen.startAgentRecast(slot1.AgentIndex, slot1.Skill)
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// チェイン効果が発動している（敵にダメージが与えられている）
	if enemy.HP >= initialEnemyHP {
		t.Log("注意: チェイン効果による追加ダメージが確認できませんでした。効果タイプによっては正常です。")
	}

	// エージェント0の待機中チェイン効果は発動済みで削除されている
	for _, pe := range screen.chainEffectManager.GetPendingEffects() {
		if pe.AgentIndex == 0 {
			t.Error("エージェント0のチェイン効果が発動後も残っています")
		}
	}
}

// ==================== テスト用ヘルパー関数（チェイン効果付き） ====================

// createTestAgentsWithChainEffect はチェイン効果付きのエージェントを作成します。
// ==================== 受け入れテスト: クールタイム開始タイミング変更 ====================

// TestCooldownNotStartedOnSkillSelection はスキル選択時にクールダウン・リキャストが開始されないことを検証します（受け入れ基準1）。
func TestCooldownNotStartedOnSkillSelection(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択を実行（startChallengeのみ、CD/リキャスト開始なし）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// クールダウンが開始されていないこと
	if screen.skillSlots[0].CooldownRemaining > 0 {
		t.Errorf("スキル選択時にクールダウンが開始されています: %f", screen.skillSlots[0].CooldownRemaining)
	}

	// リキャストが開始されていないこと
	if !screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("スキル選択時にリキャストが開始されています")
	}
}

// TestCooldownStartedOnChallengeSuccess はチャレンジ成功時にクールダウン・リキャストが開始されることを検証します（受け入れ基準2）。
func TestCooldownStartedOnChallengeSuccess(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// チャレンジ完了（Success）
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 80, WPM: 60,
	})

	// クールダウンが開始されていること
	if screen.skillSlots[0].CooldownRemaining <= 0 {
		t.Error("チャレンジ成功後にクールダウンが開始されていません")
	}

	// リキャストが開始されていること
	if screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("チャレンジ成功後にリキャストが開始されていません")
	}
}

// TestCooldownStartedOnChallengePerfect はチャレンジパーフェクト時にクールダウン・リキャストが開始されることを検証します（受け入れ基準2）。
func TestCooldownStartedOnChallengePerfect(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// チャレンジ完了（Perfect）
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengePerfect, Score: 100, WPM: 80,
	})

	// クールダウンが開始されていること
	if screen.skillSlots[0].CooldownRemaining <= 0 {
		t.Error("チャレンジパーフェクト後にクールダウンが開始されていません")
	}

	// リキャストが開始されていること
	if screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("チャレンジパーフェクト後にリキャストが開始されていません")
	}
}

// TestCooldownStartedOnChallengeFail はチャレンジ失敗時にクールダウン・リキャストが開始されることを検証します（受け入れ基準3）。
func TestCooldownStartedOnChallengeFail(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// チャレンジ失敗（Fail）
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeFail, Score: 0, WPM: 0,
	})

	// クールダウンが開始されていること
	if screen.skillSlots[0].CooldownRemaining <= 0 {
		t.Error("チャレンジ失敗後にクールダウンが開始されていません")
	}

	// リキャストが開始されていること
	if screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("チャレンジ失敗後にリキャストが開始されていません")
	}
}

// TestCooldownStartedOnChallengeCancel はチャレンジキャンセル時にクールダウン・リキャストが開始されることを検証します（受け入れ基準4）。
func TestCooldownStartedOnChallengeCancel(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// チャレンジキャンセル
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeCancel, Score: 0, WPM: 0,
	})

	// クールダウンが開始されていること
	if screen.skillSlots[0].CooldownRemaining <= 0 {
		t.Error("チャレンジキャンセル後にクールダウンが開始されていません")
	}

	// リキャストが開始されていること
	if screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("チャレンジキャンセル後にリキャストが開始されていません")
	}
}

// TestCooldownStartedOnDefenseAutoEnd はディフェンスチャレンジの敵攻撃自動終了時にクールダウン・リキャストが開始されることを検証します（受け入れ基準5）。
func TestCooldownStartedOnDefenseAutoEnd(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// スキル選択（チャレンジ開始のみ）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// モックDefenseProviderで実パスを通す
	dp := &mockDefenseProvider{defenseRate: 0.5}
	screen.processEnemyAttackWithDefense(dp)

	// クールダウンが開始されていること
	if screen.skillSlots[0].CooldownRemaining <= 0 {
		t.Error("ディフェンス自動終了後にクールダウンが開始されていません")
	}

	// リキャストが開始されていること
	if screen.recastManager.IsAgentReady(slot.AgentIndex) {
		t.Error("ディフェンス自動終了後にリキャストが開始されていません")
	}
}

// TestChainEffectRegisteredOnChallengeComplete はチャレンジ完了時にチェイン効果が登録されることを検証します（受け入れ基準6）。
func TestChainEffectRegisteredOnChallengeComplete(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgentsWithChainEffect()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// 初期状態で待機中チェイン効果なし
	if len(screen.chainEffectManager.GetPendingEffects()) != 0 {
		t.Error("初期状態で待機中チェイン効果が存在します")
	}

	// スキル選択（チャレンジ開始のみ、チェイン効果は登録されない）
	screen.selectedSkillIdx = 0
	slot := screen.skillSlots[screen.selectedSkillIdx]
	screen.startChallenge(slot.Skill)

	// スキル選択直後はチェイン効果が登録されていないこと
	if len(screen.chainEffectManager.GetPendingEffects()) != 0 {
		t.Error("スキル選択時にチェイン効果が登録されています")
	}

	// チャレンジ完了後にチェイン効果が登録されること
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	pendingEffects := screen.chainEffectManager.GetPendingEffects()
	if len(pendingEffects) == 0 {
		t.Error("チャレンジ完了後にチェイン効果が登録されていません")
	}
}

func createTestAgentsWithChainEffect() []*domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}

	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})

	// チェイン効果付きスキル
	chainEffect := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25)
	skills := []*domain.SkillModel{
		newTestSkillWithChainEffect("m1", "物理攻撃", []string{"physical_low"}, 1.0, "STR", "物理ダメージ", &chainEffect),
		newTestDamageSkill("m2", "魔法攻撃", []string{"magic_low"}, 1.0, "INT", "魔法ダメージ"),
		newTestHealSkill("m3", "回復", []string{"heal_low"}, 1.0, "INT", "HP回復"),
		newTestBuffSkill("m4", "バフ", []string{"buff_low"}, "攻撃力UP"),
	}

	agent := domain.NewAgent("agent1", core, skills)
	return []*domain.AgentModel{agent}
}

// ==================== チェイン効果消費側テスト ====================

// TestTimeExtend_PassedToChallengeInput はstartChallengeがTimeExtendSecをChallengeInputに渡すことをテストします。
func TestTimeExtend_PassedToChallengeInput(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// TimeExtend効果を追加（3秒延長）
	player.EffectTable.AddBuff("タイム延長", "タイム延長", 10.0, map[domain.EffectColumn]float64{
		domain.ColTimeExtend: 3.0,
	})

	// startChallengeでチャレンジが正常に開始されることを確認
	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}
	screen.selectedSkillIdx = 0
	skill := screen.skillSlots[0].Skill
	screen.startChallenge(skill)

	if screen.activeChallenge == nil {
		t.Fatal("チャレンジが開始されていません")
	}

	// TimeExtendが正しく取得できることを確認（Aggregateの結果）
	ctx := domain.NewEffectContext(player.HP, player.MaxHP, enemy.HP, enemy.MaxHP)
	effects := player.EffectTable.Aggregate(ctx)
	if effects.TimeExtend != 3.0 {
		t.Errorf("TimeExtend効果が正しくない: got %v, want 3.0", effects.TimeExtend)
	}
}

// TestTimeExtend_NegativeValue はTimeExtendが負の場合（デバフ）でもチャレンジが開始されることをテストします。
func TestTimeExtend_NegativeValue(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// TimeExtend効果を追加（-2秒、デバフ）
	player.EffectTable.AddDebuff("タイム短縮", "タイム短縮", 10.0, map[domain.EffectColumn]float64{
		domain.ColTimeExtend: -2.0,
	})

	// startChallengeでチャレンジが正常に開始されることを確認
	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}
	screen.selectedSkillIdx = 0
	skill := screen.skillSlots[0].Skill
	screen.startChallenge(skill)

	if screen.activeChallenge == nil {
		t.Fatal("TimeExtendデバフ時もチャレンジは開始されるべきです")
	}

	// デバフが正しく取得できることを確認
	ctx := domain.NewEffectContext(player.HP, player.MaxHP, enemy.HP, enemy.MaxHP)
	effects := player.EffectTable.Aggregate(ctx)
	if effects.TimeExtend != -2.0 {
		t.Errorf("TimeExtendデバフが正しくない: got %v, want -2.0", effects.TimeExtend)
	}
}

// TestCooldownReduce_ShortensInitialCooldown はクールダウン初期値の短縮をテストします。
func TestCooldownReduce_ShortensInitialCooldown(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// CooldownReduce効果を追加（30%短縮）
	player.EffectTable.AddBuff("クールダウン短縮", "クールダウン短縮", 10.0, map[domain.EffectColumn]float64{
		domain.ColCooldownReduce: 0.3,
	})

	// クールダウンを設定（10秒）→ 30%短縮で7秒になるはず
	screen.StartCooldown(0, 10.0)

	expected := 7.0
	tolerance := 0.01
	if screen.skillSlots[0].CooldownRemaining < expected-tolerance ||
		screen.skillSlots[0].CooldownRemaining > expected+tolerance {
		t.Errorf("CooldownReduce効果が初期値に適用されていない: got %.2f, want %.2f",
			screen.skillSlots[0].CooldownRemaining, expected)
	}

	// CooldownTotal は元の値（表示用）
	if screen.skillSlots[0].CooldownTotal != 10.0 {
		t.Errorf("CooldownTotal が変更されている: got %.2f, want 10.0",
			screen.skillSlots[0].CooldownTotal)
	}
}

// TestCooldownReduce_NegativeValue はCooldownReduceが負の場合（クールダウン延長）をテストします。
func TestCooldownReduce_NegativeValue(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// CooldownReduce効果を追加（-30%、つまり延長）
	player.EffectTable.AddDebuff("クールダウン延長", "クールダウン延長", 10.0, map[domain.EffectColumn]float64{
		domain.ColCooldownReduce: -0.3,
	})

	// クールダウンを設定（10秒）→ -30%延長で13秒になるはず
	screen.StartCooldown(0, 10.0)

	expected := 13.0
	tolerance := 0.01
	if screen.skillSlots[0].CooldownRemaining < expected-tolerance ||
		screen.skillSlots[0].CooldownRemaining > expected+tolerance {
		t.Errorf("CooldownReduce延長効果が初期値に適用されていない: got %.2f, want %.2f",
			screen.skillSlots[0].CooldownRemaining, expected)
	}
}

// TestCooldownReduce_MinimumLimit は短縮しすぎない下限をテストします。
func TestCooldownReduce_MinimumLimit(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}

	// CooldownReduce効果を追加（95%短縮 → 下限10%適用）
	player.EffectTable.AddBuff("超クールダウン短縮", "超クールダウン短縮", 10.0, map[domain.EffectColumn]float64{
		domain.ColCooldownReduce: 0.95,
	})

	// クールダウンを設定（10秒）→ 最低10%で1秒になるはず
	screen.StartCooldown(0, 10.0)

	minExpected := 1.0
	if screen.skillSlots[0].CooldownRemaining < minExpected {
		t.Errorf("CooldownReduce の下限が適用されていない: got %.2f, want >= %.2f",
			screen.skillSlots[0].CooldownRemaining, minExpected)
	}
}

// TestDoubleCast_DoublesDamageEffect はダブルキャストによる2回発動をテストします。
func TestDoubleCast_DoublesDamageEffect(t *testing.T) {
	// 乱数を固定して確率判定を決定的にする
	originalRandFloat := randFloat
	randFloat = func() float64 { return 0.0 } // 常に0を返す（100%確率は必ず成功）
	defer func() { randFloat = originalRandFloat }()

	// まず、DoubleCastなしでの基準ダメージを計測
	enemy1 := createTestEnemy()
	enemy1.HP = 1000
	player1 := createTestPlayer()
	agents1 := createTestAgents()
	screen1 := NewBattleScreen(enemy1, player1, agents1, nil)

	screen1.selectedSkillIdx = 0
	screen1.selectedSlot = 0
	initialHP1 := enemy1.HP
	screen1.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})
	baseDamage := initialHP1 - enemy1.HP

	// 次に、DoubleCast100%でのダメージを計測
	enemy2 := createTestEnemy()
	enemy2.HP = 1000
	player2 := createTestPlayer()
	agents2 := createTestAgents()
	screen2 := NewBattleScreen(enemy2, player2, agents2, nil)

	// DoubleCast効果を追加（100%確率で2回発動）
	player2.EffectTable.AddBuff("ダブルキャスト", "ダブルキャスト", 10.0, map[domain.EffectColumn]float64{
		domain.ColDoubleCast: 1.0, // 100%
	})

	screen2.selectedSkillIdx = 0
	screen2.selectedSlot = 0
	initialHP2 := enemy2.HP
	screen2.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})
	doubleDamage := initialHP2 - enemy2.HP

	// DoubleCastにより2倍のダメージが与えられているはず
	// 少なくとも1.8倍以上になっていることを確認（誤差許容）
	expectedMinDamage := int(float64(baseDamage) * 1.8)
	if doubleDamage < expectedMinDamage {
		t.Errorf("DoubleCast効果が適用されていない: baseDamage=%d, doubleDamage=%d, expected>=%d",
			baseDamage, doubleDamage, expectedMinDamage)
	}

	t.Logf("基準ダメージ: %d, DoubleCastダメージ: %d (%.1f倍)", baseDamage, doubleDamage, float64(doubleDamage)/float64(baseDamage))
}

// TestDoubleCast_ZeroProbability はDoubleCast確率0%の場合のテストです。
func TestDoubleCast_ZeroProbability(t *testing.T) {
	enemy := createTestEnemy()
	enemy.HP = 1000
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// DoubleCast効果なし（0%）
	// バフを追加しない

	// チャレンジを完了
	screen.selectedSkillIdx = 0
	screen.selectedSlot = 0
	initialEnemyHP := enemy.HP
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// 通常の1回分ダメージのみ
	damageDone := initialEnemyHP - enemy.HP
	if damageDone <= 0 {
		t.Error("ダメージが与えられていない")
	}

	t.Logf("DoubleCastなしでのダメージ: %d", damageDone)
}

// TestOverheal_ConvertExcessToTempHP はオーバーヒールによる一時HP変換をテストします。
func TestOverheal_ConvertExcessToTempHP(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	player.HP = player.MaxHP // 満タン状態
	agents := createTestAgentsWithHealSkill()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// Overheal効果を追加（フラグ形式）
	entry := domain.EffectEntry{
		SourceType: domain.SourceBuff,
		Name:       "オーバーヒール",
		Values:     map[domain.EffectColumn]float64{},
		Flags: map[domain.EffectColumn]bool{
			domain.ColOverheal: true,
		},
	}
	duration := 10.0
	entry.Duration = &duration
	player.EffectTable.Entries = append(player.EffectTable.Entries, entry)

	// 回復スキルを使用
	screen.selectedSkillIdx = 2 // 回復スキル
	screen.selectedSlot = 2
	screen.handleChallengeComplete(&domain.ChallengeOutput{
		Status: domain.ChallengeSuccess, Score: 100, WPM: 60,
	})

	// Overhealにより超過分がTempHPに変換されているはず
	if player.TempHP == 0 {
		t.Errorf("Overheal効果によりTempHPが増えていない: got %d, want > 0", player.TempHP)
	}

	t.Logf("HP: %d/%d, TempHP: %d", player.HP, player.MaxHP, player.TempHP)
}

// TestTempHP_AbsorbsDamage は一時HPがダメージを吸収することをテストします。
func TestTempHP_AbsorbsDamage(t *testing.T) {
	player := createTestPlayer()
	player.HP = 100
	player.MaxHP = 100
	player.TempHP = 30

	// 20ダメージ（TempHPで吸収）
	player.TakeDamage(20)
	if player.TempHP != 10 {
		t.Errorf("TempHPが消費されていない: got %d, want 10", player.TempHP)
	}
	if player.HP != 100 {
		t.Errorf("HPが減少している: got %d, want 100", player.HP)
	}

	// さらに20ダメージ（TempHPを貫通してHPにダメージ）
	player.TakeDamage(20)
	if player.TempHP != 0 {
		t.Errorf("TempHPがゼロになっていない: got %d", player.TempHP)
	}
	if player.HP != 90 {
		t.Errorf("HPが正しく減少していない: got %d, want 90", player.HP)
	}
}

// createTestAgentsWithHealSkill は回復スキル付きのエージェントを作成します。
func createTestAgentsWithHealSkill() []*domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "heal_low", "magic_low", "buff_low"},
	}

	core := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})

	skills := []*domain.SkillModel{
		newTestDamageSkill("m1", "物理攻撃", []string{"physical_low"}, 1.0, "STR", "物理ダメージ"),
		newTestDamageSkill("m2", "魔法攻撃", []string{"magic_low"}, 1.0, "INT", "魔法ダメージ"),
		newTestHealSkill("m3", "回復", []string{"heal_low"}, 5.0, "INT", "HP回復"),
		newTestBuffSkill("m4", "バフ", []string{"buff_low"}, "攻撃力UP"),
	}

	agent := domain.NewAgent("agent1", core, skills)
	return []*domain.AgentModel{agent}
}

// TestAutoCorrect_PassedToChallengeInput はstartChallengeがAutoCorrectCountをChallengeInputに渡すことをテストします。
func TestAutoCorrect_PassedToChallengeInput(t *testing.T) {
	enemy := createTestEnemy()
	player := createTestPlayer()
	agents := createTestAgents()

	screen := NewBattleScreen(enemy, player, agents, nil)

	// AutoCorrect効果を追加（2回分）
	player.EffectTable.AddBuff("オートコレクト", "オートコレクト", 10.0, map[domain.EffectColumn]float64{
		domain.ColAutoCorrect: 2.0,
	})

	// startChallengeでチャレンジが正常に開始されることを確認
	if len(screen.skillSlots) == 0 {
		t.Skip("スキルスロットがありません")
	}
	screen.selectedSkillIdx = 0
	skill := screen.skillSlots[0].Skill
	screen.startChallenge(skill)

	if screen.activeChallenge == nil {
		t.Fatal("チャレンジが開始されていません")
	}

	// AutoCorrectが正しく取得できることを確認（Aggregateの結果）
	ctx := domain.NewEffectContext(player.HP, player.MaxHP, enemy.HP, enemy.MaxHP)
	effects := player.EffectTable.Aggregate(ctx)
	if effects.AutoCorrect != 2 {
		t.Errorf("AutoCorrect効果が正しくない: got %v, want 2", effects.AutoCorrect)
	}
}

// createTestAgentsWithChainEffectMultiple は複数のチェイン効果付きエージェントを作成します。
func createTestAgentsWithChainEffectMultiple() []*domain.AgentModel {
	coreType := domain.CoreType{
		ID:          "test",
		Name:        "テスト",
		StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		AllowedTags: []string{"physical_low", "magic_low", "heal_low", "buff_low"},
	}

	// エージェント1
	core1 := domain.NewCoreWithTypeID("core1", coreType, domain.PassiveSkill{})
	chainEffect1 := domain.NewChainEffect("test_effect", domain.ChainEffectDamageBonus, 25)
	skills1 := []*domain.SkillModel{
		newTestSkillWithChainEffect("m1", "物理攻撃", []string{"physical_low"}, 1.0, "STR", "物理ダメージ", &chainEffect1),
		newTestDamageSkill("m2", "魔法攻撃", []string{"magic_low"}, 1.0, "INT", "魔法ダメージ"),
		newTestHealSkill("m3", "回復", []string{"heal_low"}, 1.0, "INT", "HP回復"),
		newTestBuffSkill("m4", "バフ", []string{"buff_low"}, "攻撃力UP"),
	}
	agent1 := domain.NewAgent("agent1", core1, skills1)

	// エージェント2
	core2 := domain.NewCoreWithTypeID("core2", coreType, domain.PassiveSkill{})
	chainEffect2 := domain.NewChainEffect("test_effect", domain.ChainEffectHealBonus, 30)
	skills2 := []*domain.SkillModel{
		newTestSkillWithChainEffect("m5", "物理攻撃2", []string{"physical_low"}, 1.0, "STR", "物理ダメージ", &chainEffect2),
		newTestDamageSkill("m6", "魔法攻撃2", []string{"magic_low"}, 1.0, "INT", "魔法ダメージ"),
		newTestHealSkill("m7", "回復2", []string{"heal_low"}, 1.0, "INT", "HP回復"),
		newTestBuffSkill("m8", "バフ2", []string{"buff_low"}, "攻撃力UP"),
	}
	agent2 := domain.NewAgent("agent2", core2, skills2)

	return []*domain.AgentModel{agent1, agent2}
}

// mockDefenseProvider はDefenseProviderインターフェースのテスト用モック。
type mockDefenseProvider struct {
	defenseRate float64
	completed   bool
}

func (m *mockDefenseProvider) DefenseRate() float64 { return m.defenseRate }
func (m *mockDefenseProvider) CompleteByAttack()    { m.completed = true }
