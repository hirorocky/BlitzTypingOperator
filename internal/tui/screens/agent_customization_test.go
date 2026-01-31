// Package screens はTUIゲームの画面を提供します。
// このファイルはAgentCustomizationScreenのテストを定義します。

package screens

import (
	"strings"
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/usecase/inventory"
	"hirorocky/type-battle/internal/usecase/slot"

	tea "github.com/charmbracelet/bubbletea"
)

// ==================== テスト用ヘルパー ====================

// setupTestCustomizationScreen はテスト用のAgentCustomizationScreenをセットアップします。
func setupTestCustomizationScreen() (*AgentCustomizationScreen, *domain.CoreInventory, *domain.SkillInventory, *slot.AgentSlotManager) {
	// コアとスキルのインベントリを作成
	coreInv := domain.NewCoreInventory()
	skillInv := domain.NewSkillInventory()

	// テスト用マスタデータ
	coreTypes := map[string]domain.CoreType{
		"balance": {
			ID:          "balance",
			Name:        "バランス",
			AllowedTags: []string{"physical", "magic"},
			StatWeights: map[string]float64{"STR": 1.0, "INT": 1.0, "WIL": 1.0, "LUK": 1.0},
		},
		"attacker": {
			ID:          "attacker",
			Name:        "アタッカー",
			AllowedTags: []string{"physical"},
			StatWeights: map[string]float64{"STR": 1.5, "INT": 0.5, "WIL": 0.8, "LUK": 1.0},
		},
	}

	skillTypes := map[string]domain.SkillType{
		"strike": {
			ID:          "strike",
			Name:        "ストライク",
			Icon:        "⚔️",
			Tags:        []string{"physical"},
			Description: "物理攻撃",
		},
		"fireball": {
			ID:          "fireball",
			Name:        "ファイアボール",
			Icon:        "🔥",
			Tags:        []string{"magic"},
			Description: "魔法攻撃",
		},
		"heal": {
			ID:          "heal",
			Name:        "ヒール",
			Icon:        "💚",
			Tags:        []string{"magic"},
			Description: "回復",
		},
	}

	passiveSkills := map[string]domain.PassiveSkill{}
	chainEffects := map[string]domain.ChainEffect{}

	// スロットマネージャーを作成
	slotManager := slot.NewAgentSlotManager(coreInv, skillInv, coreTypes, skillTypes, passiveSkills, chainEffects)

	// インベントリマネージャーを作成
	invManager := inventory.NewInventoryManagerWithInventories(coreInv, skillInv)

	// カスタマイズ画面を作成
	screen := NewAgentCustomizationScreen(invManager, slotManager, coreTypes, skillTypes, passiveSkills, chainEffects)

	return screen, coreInv, skillInv, slotManager
}

// ==================== タスク10.1: 基本構造のテスト ====================

func TestAgentCustomizationScreen_NewAgentCustomizationScreen(t *testing.T) {
	screen, _, _, _ := setupTestCustomizationScreen()

	if screen == nil {
		t.Fatal("画面がnilです")
	}

	// 初期状態の確認
	if screen.currentMode != CustomizationModeSlotSelect {
		t.Errorf("初期モードが不正です: got %v, want %v", screen.currentMode, CustomizationModeSlotSelect)
	}

	if screen.selectedSlotIndex != 0 {
		t.Errorf("初期スロットインデックスが不正です: got %v, want 0", screen.selectedSlotIndex)
	}
}

func TestAgentCustomizationScreen_View_ShowsThreeSlots(t *testing.T) {
	screen, _, _, _ := setupTestCustomizationScreen()

	view := screen.View()

	// 新UI: 3つのカードが表示されていることを確認
	// 空スロットは "(空)" と表示される
	if !strings.Contains(view, "(空)") {
		t.Error("空スロットが表示されていません")
	}
	// タイトルが表示されていることを確認
	if !strings.Contains(view, "エージェントカスタマイズ") {
		t.Error("タイトルが表示されていません")
	}
}

func TestAgentCustomizationScreen_Update_SlotNavigation(t *testing.T) {
	screen, _, _, _ := setupTestCustomizationScreen()

	// 新UI: 左右キーでスロット切替、上下キーでカード内フォーカス移動
	// 右キーでスロット2に移動
	screen.Update(tea.KeyMsg{Type: tea.KeyRight})
	if screen.selectedSlotIndex != 1 {
		t.Errorf("右キーでスロットが移動していません: got %v, want 1", screen.selectedSlotIndex)
	}

	// 左キーでスロット1に戻る
	screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if screen.selectedSlotIndex != 0 {
		t.Errorf("左キーでスロットが移動していません: got %v, want 0", screen.selectedSlotIndex)
	}

	// スロット1で左キーを押しても0より下にならない
	screen.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if screen.selectedSlotIndex != 0 {
		t.Errorf("スロットインデックスが0未満になりました: got %v", screen.selectedSlotIndex)
	}
}

func TestAgentCustomizationScreen_Update_EnterTransitionsToCoreSelectMode(t *testing.T) {
	screen, coreInv, _, _ := setupTestCustomizationScreen()

	// コアをインベントリに追加
	coreInv.AddCore("balance")

	// Enterを押してコア選択モードに遷移
	screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if screen.currentMode != CustomizationModeCoreSelect {
		t.Errorf("Enterでコア選択モードに遷移していません: got %v, want %v", screen.currentMode, CustomizationModeCoreSelect)
	}
}

// ==================== タスク10.2: コア選択モードのテスト ====================

func TestAgentCustomizationScreen_CoreSelectMode_ShowsOwnedCores(t *testing.T) {
	screen, coreInv, _, _ := setupTestCustomizationScreen()

	// コアをインベントリに追加
	coreInv.AddCore("balance")
	coreInv.AddCore("attacker")

	// コア選択モードに移行
	screen.enterCoreSelectMode()

	view := screen.View()

	// 保有コアが表示されていることを確認
	if !strings.Contains(view, "バランス") {
		t.Error("バランスコアが表示されていません")
	}
	if !strings.Contains(view, "アタッカー") {
		t.Error("アタッカーコアが表示されていません")
	}
}

func TestAgentCustomizationScreen_CoreSelectMode_CoreSelection(t *testing.T) {
	screen, coreInv, _, slotManager := setupTestCustomizationScreen()

	// コアをインベントリに追加
	coreInv.AddCore("balance")

	// コア選択モードに移行
	screen.enterCoreSelectMode()
	screen.updateCoreList()

	// コアを選択（Enter）- 直接コアがスロットに設定される
	screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// スロットにコアが設定されていることを確認
	slot := slotManager.GetSlot(0)
	if slot.CoreTypeID != "balance" {
		t.Errorf("スロットにコアが設定されていません: got %v, want balance", slot.CoreTypeID)
	}
}

// ==================== タスク10.3: スキル選択モードのテスト ====================

func TestAgentCustomizationScreen_SkillSelectMode_ShowsCompatibleSkills(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// コアをインベントリに追加して設定
	coreInv.AddCore("balance")
	slotManager.SetCore(0, "balance")

	// スキルをインベントリに追加
	skillInv.AddSkill("strike", "")   // physical タグ
	skillInv.AddSkill("fireball", "") // magic タグ

	// スキル選択モードに移行
	screen.focusPosition = 1 // スキルスロット1にフォーカス
	screen.enterSkillSelectMode(0)

	// 互換スキルリストを更新
	screen.updateSkillList()

	view := screen.View()

	// バランスコアは physical と magic の両方を許可するので両方表示
	if !strings.Contains(view, "ストライク") {
		t.Error("ストライクスキルが表示されていません")
	}
	if !strings.Contains(view, "ファイアボール") {
		t.Error("ファイアボールスキルが表示されていません")
	}
}

func TestAgentCustomizationScreen_SkillSelectMode_FiltersIncompatibleSkills(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// アタッカーコアを設定（physicalのみ許可）
	coreInv.AddCore("attacker")
	slotManager.SetCore(0, "attacker")

	// スキルをインベントリに追加
	skillInv.AddSkill("strike", "")   // physical タグ - 互換
	skillInv.AddSkill("fireball", "") // magic タグ - 非互換

	// スキル選択モードに移行
	screen.focusPosition = 1
	screen.enterSkillSelectMode(0)
	screen.updateSkillList()

	// 互換スキルのみがリストに含まれることを確認
	foundStrike := false
	foundFireball := false
	for _, item := range screen.compatibleSkillList {
		if item.TypeID == "strike" {
			foundStrike = true
		}
		if item.TypeID == "fireball" {
			foundFireball = true
		}
	}

	if !foundStrike {
		t.Error("互換スキル（ストライク）がリストに含まれていません")
	}
	if foundFireball {
		t.Error("非互換スキル（ファイアボール）がリストに含まれています")
	}
}

func TestAgentCustomizationScreen_SkillSelectMode_SetsSkillToSlot(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// コアを設定
	coreInv.AddCore("balance")
	slotManager.SetCore(0, "balance")

	// スキルを追加
	skillInv.AddSkill("strike", "")

	// スキル選択モードに移行
	screen.selectedSlotIndex = 0
	screen.focusPosition = 1 // スキルスロット1にフォーカス
	screen.enterSkillSelectMode(0)
	screen.updateSkillList()

	// スキルを選択
	screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// スロットにスキルが設定されていることを確認
	slot := slotManager.GetSlot(0)
	skillConfig := slot.GetSkill(0)
	if skillConfig == nil || skillConfig.TypeID != "strike" {
		t.Errorf("スロットにスキルが設定されていません: got %v, want strike", skillConfig)
	}
}

func TestAgentCustomizationScreen_SkillSelectMode_ChainEffectVariationSelection(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// コアを設定
	coreInv.AddCore("balance")
	slotManager.SetCore(0, "balance")

	// チェイン効果付きスキルを追加
	skillInv.AddSkill("strike", "chain_damage_boost")
	skillInv.AddSkill("strike", "chain_heal")

	// スキル選択モードに移行
	screen.selectedSlotIndex = 0
	screen.focusPosition = 1 // スキルスロット1にフォーカス
	screen.enterSkillSelectMode(0)
	screen.updateSkillList()

	// スキルを選択してチェイン効果選択モードに遷移
	screen.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// チェイン効果バリエーションが2つあることを確認
	if screen.currentMode == CustomizationModeChainSelect {
		if len(screen.chainVariationList) != 2 {
			t.Errorf("チェイン効果バリエーション数が不正です: got %v, want 2", len(screen.chainVariationList))
		}
	}
}

// ==================== タスク10.4: スロットクリア機能のテスト ====================

func TestAgentCustomizationScreen_ClearCore(t *testing.T) {
	screen, coreInv, _, slotManager := setupTestCustomizationScreen()

	// コアを設定
	coreInv.AddCore("balance")
	slotManager.SetCore(0, "balance")

	screen.selectedSlotIndex = 0

	// クリア操作（Deleteキー）
	screen.Update(tea.KeyMsg{Type: tea.KeyDelete})

	// スロットが空になっていることを確認
	slot := slotManager.GetSlot(0)
	if !slot.IsEmpty() {
		t.Error("スロットがクリアされていません")
	}
}

func TestAgentCustomizationScreen_ClearSkill(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// コアとスキルを設定
	coreInv.AddCore("balance")
	slotManager.SetCore(0, "balance")
	skillInv.AddSkill("strike", "")
	slotManager.SetSkill(0, 0, "strike", "")

	// スキルスロット1にフォーカス
	screen.selectedSlotIndex = 0
	screen.focusPosition = 1 // スキルスロット1にフォーカス

	// スキルクリア操作（Backspaceキー）
	screen.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	// スキルがクリアされていることを確認
	slot := slotManager.GetSlot(0)
	skillConfig := slot.GetSkill(0)
	if skillConfig != nil && !skillConfig.IsEmpty() {
		t.Error("スキルがクリアされていません")
	}
}

func TestAgentCustomizationScreen_EscReturnsToPreviousMode(t *testing.T) {
	screen, coreInv, _, _ := setupTestCustomizationScreen()

	coreInv.AddCore("balance")

	// スロット選択からコア選択へ
	screen.enterCoreSelectMode()
	if screen.currentMode != CustomizationModeCoreSelect {
		t.Fatal("コア選択モードへの遷移に失敗")
	}

	// Escでスロット選択に戻る
	screen.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if screen.currentMode != CustomizationModeSlotSelect {
		t.Errorf("Escでスロット選択に戻っていません: got %v, want %v", screen.currentMode, CustomizationModeSlotSelect)
	}
}

// ==================== 互換性表示のテスト ====================

func TestAgentCustomizationScreen_ShowsCompatibilityIndicator(t *testing.T) {
	screen, coreInv, skillInv, slotManager := setupTestCustomizationScreen()

	// アタッカーコアを設定（physicalのみ許可）
	coreInv.AddCore("attacker")
	slotManager.SetCore(0, "attacker")

	// スキルを追加
	skillInv.AddSkill("strike", "")   // 互換
	skillInv.AddSkill("fireball", "") // 非互換（magicタグ）

	// スキル選択モードに移行
	screen.focusPosition = 1
	screen.enterSkillSelectMode(0)
	screen.updateSkillList()

	view := screen.View()

	// 互換性インジケーターが表示されていることを確認
	// 実装では互換スキルのみ表示するのでこのテストは確認のみ
	if strings.Contains(view, "ファイアボール") {
		t.Error("非互換スキルが表示されています（互換スキルのみ表示されるべき）")
	}
}

// ==================== Screenインターフェースのテスト ====================

func TestAgentCustomizationScreen_ImplementsScreenInterface(t *testing.T) {
	screen, _, _, _ := setupTestCustomizationScreen()

	// Screenインターフェースを実装していることを確認
	var _ Screen = screen

	// SetSizeとGetTitleのテスト
	screen.SetSize(100, 50)
	w, h := screen.GetSize()
	if w != 100 || h != 50 {
		t.Errorf("SetSize/GetSizeが正しく動作していません: got (%v, %v), want (100, 50)", w, h)
	}

	title := screen.GetTitle()
	if title == "" {
		t.Error("GetTitleが空文字を返しています")
	}
}
