# 実装計画: チェイン効果の持続仕様変更

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル | 状態 |
|---|------|-----------|---------------|------|
| 1 | リキャスト完了後もチェイン効果が待機状態を維持する | Go test | `internal/tui/screens/battle_test.go` (`TestBattleScreenRecastCompletionPersistsChainEffect`) | PASS |
| 2 | リキャスト完了後に他エージェントがスキル使用でチェイン効果が発動する | Go test | `internal/tui/screens/battle_test.go` (`TestBattleScreenChainEffectTriggersAfterRecastCompletion`) | PASS |
| 3 | 同エージェントがチェイン効果付きスキル使用時に上書きされる | Go test | `internal/usecase/combat/chain/chain_effect_manager_test.go` (`TestMultipleEffectsFromSameAgent`) | PASS |
| 4 | チェイン効果なしスキル使用時に既存チェイン効果が削除される | Go test | `internal/usecase/combat/chain/chain_effect_manager_test.go` (`TestClearEffectForAgent`, `TestClearEffectForAgentNoEffect`) + `internal/tui/screens/battle_test.go` (`TestBattleScreenStartRecastClearsChainEffectForNoChainSkill`) | PASS |
| 5 | チェイン効果発動時にEffectTableへの登録が正常に行われる | Go test | 既存テスト群 | PASS |
| 6 | リキャスト中でなくても待機中チェイン効果が黄色背景で表示される | TUI test | セッション内 | 手順記載済 |
| 7 | handleChallengeComplete内の実行順序が保証される | Go test | 既存テスト群 + コード構造 | PASS |

## 実装タスク

### タスク1: ClearEffectForAgentメソッドの追加
- **対象**: `internal/usecase/combat/chain/chain_effect_manager.go`
- **内容**: `ClearEffectForAgent(agentIndex int)` メソッドを追加。指定エージェントの待機中チェイン効果を削除する
- **関連テスト**: 基準4（`TestClearEffectForAgent`, `TestClearEffectForAgentNoEffect`）
- **状態**: 完了

### タスク2: ExpireEffectsForAgentの削除とコメント更新
- **対象**: `internal/usecase/combat/chain/chain_effect_manager.go`, `internal/usecase/combat/chain/chain_effect_manager_test.go`
- **内容**:
  - `ExpireEffectsForAgent()` メソッドを削除
  - `TestExpireEffectsForAgent` テストを削除
  - パッケージコメント・型コメントの「リキャスト期間中」を「次のスキル使用まで」に更新
- **関連テスト**: 基準1（コンパイルが通る前提条件）
- **状態**: 完了

### タスク3: UpdateRecastsからチェイン効果破棄処理を削除
- **対象**: `internal/tui/screens/battle_logic.go`
- **内容**:
  - `UpdateRecasts()` から `ExpireEffectsForAgent` 呼び出しとループを削除
  - `completedAgents` 変数が未使用になるため、`s.recastManager.UpdateRecast(delta)` の戻り値を受け取らない形に変更
  - 関数コメントからチェイン効果破棄の記述を削除
- **関連テスト**: 基準1（`TestBattleScreenRecastCompletionPersistsChainEffect`）、基準2（`TestBattleScreenChainEffectTriggersAfterRecastCompletion`）
- **状態**: 完了

### タスク4: startAgentRecastにチェイン効果クリア処理を追加
- **対象**: `internal/tui/screens/battle_logic.go`
- **内容**: `startAgentRecast()` で `module.ChainEffect == nil` かつ `chainEffectManager != nil` の場合に `ClearEffectForAgent(agentIndex)` を呼び出す
- **関連テスト**: 基準4（`TestBattleScreenStartRecastClearsChainEffectForNoChainSkill`）
- **状態**: 完了

## TUIテスト手順

### 基準6: リキャスト中でなくても待機中チェイン効果が黄色背景で表示される
```
launch_tui(command="go run ./cmd/BlitzTypingOperator -debug -save testdata/saves/chain_effect_agents.json", mode="buffer", dimensions="160x45")
# ホーム画面→バトル選択→バトル開始
send_keys("jjj", delay=2.0)  # 初回はコンパイル待ち。バトル選択まで移動
send_keys("\r")               # バトル選択へ
send_keys("\r")               # バトル開始

# チェイン効果付きスキルを使用（スロット1を選択してタイピング完了）
send_keys("1")                # スロット1選択
# （タイピング完了後、チェイン効果がバッジ表示される）

# リキャスト完了まで待機
# リキャスト完了後もチェイン効果バッジが黄色背景で表示されていることを確認
capture_screen()
# 黄色背景のチェイン効果バッジの存在を目視確認
```
注: このTUIテストはゲーム操作の手動介在が必要なため、完全自動化が難しい。主にGoテストで検証し、TUIテストは補助的に使用する。
フィクスチャ: `testdata/saves/chain_effect_agents.json`（全スキルにチェイン効果付き）

## 進捗ログ
- タスク1: `ClearEffectForAgent(agentIndex int)` メソッドを追加。`delete(m.pendingEffects, agentIndex)` の1行で実装
- タスク2: `ExpireEffectsForAgent` メソッドとテストを削除。パッケージコメントを「次のスキル使用まで」に更新
- タスク3: `UpdateRecasts()` から `ExpireEffectsForAgent` 呼び出しを削除。戻り値も未使用のため受け取らない形に変更
- タスク4: `startAgentRecast()` で `ChainEffect == nil` の場合に `ClearEffectForAgent` を呼ぶ分岐を追加
