# チェイン効果の持続仕様変更

## 概要
チェイン効果の有効期間を「リキャスト中のみ」から「次のスキル使用まで」に変更する。
現在はリキャスト完了時にチェイン効果が破棄されるため、リキャスト後に他エージェントがスキルを使用してもチェイン効果が発動しない。この変更により、チェイン効果をより戦略的に活用できるようになる。

## 受け入れ基準
1. スキル使用後のチェイン効果がリキャスト完了後も待機状態を維持する
2. リキャスト完了後に他エージェントがスキルを使用した場合、待機中チェイン効果が発動条件に基づいて発動する
3. 同エージェントが次にチェイン効果付きスキルを使用した際、新しいチェイン効果で既存のチェイン効果が上書きされる
4. 同エージェントがチェイン効果なしスキルを使用した際、既存の待機中チェイン効果が削除される
5. チェイン効果発動時にEffectTableへの登録が正常に行われる（従来通り）
6. UI: リキャスト中でなくても待機中チェイン効果が黄色背景で表示される
7. handleChallengeComplete内で「チェイン効果発動→新チェイン効果登録」の順序が保証され、前のチェイン効果が次のスキル効果計算前に消えない

## 設計
### 変更対象
- `internal/tui/screens/battle_logic.go`:
  - `UpdateRecasts()`: リキャスト完了時のチェイン効果破棄処理を削除（`ExpireEffectsForAgent`呼び出しとそのループの削除）。`completedAgents`変数が未使用になるため、`s.recastManager.UpdateRecast(delta)`の戻り値を受け取らない形に変更
  - `UpdateRecasts()`の関数コメントを更新（チェイン効果破棄の記述を削除）
  - `startAgentRecast()`: `module.ChainEffect == nil`の場合に既存のチェイン効果を削除する処理を追加
- `internal/usecase/combat/chain/chain_effect_manager.go`:
  - `ExpireEffectsForAgent()`メソッドを削除し、代わりに`ClearEffectForAgent(agentIndex int)`メソッドを追加（戻り値不要のシンプルな削除メソッド。`startAgentRecast`でチェイン効果なしスキル使用時に呼び出す）
  - パッケージコメント・型コメントの「リキャスト期間中」に関する記述を「次のスキル使用まで」に更新
- `internal/usecase/combat/chain/chain_effect_manager_test.go`:
  - `TestExpireEffectsForAgent`テストを削除
- `internal/tui/screens/battle_test.go`:
  - `TestBattleScreenRecastCompletionExpiresChainEffect`（1091-1135行）を新仕様に合わせて修正（「リキャスト完了後もチェイン効果が待機状態を維持する」ことを検証するテストに変更）
- `internal/tui/screens/battle_view.go`:
  - チェイン効果の表示条件は変更不要（`pendingChain != nil`で判定しており、リキャスト状態は条件に含まれていない。効果が破棄されなくなることで自然に表示が維持される）

### 新規作成
- `internal/usecase/combat/chain/chain_effect_manager.go`: `ClearEffectForAgent(agentIndex int)` — 指定エージェントの待機中チェイン効果を削除する。`startAgentRecast`でチェイン効果なしスキル使用時に呼び出す

### データフロー

**現在のフロー:**
```
スキル使用 → RegisterChainEffect → リキャスト中（待機） → リキャスト完了 → ExpireEffectsForAgent（破棄）
```

**変更後のフロー:**
```
スキル使用 → RegisterChainEffect → リキャスト中（待機） → リキャスト完了（待機を維持）
  → 他エージェントスキル使用: CheckAndTrigger（発動→削除）
  → 同エージェント次スキル使用（チェイン効果あり）: RegisterChainEffect（上書き）
  → 同エージェント次スキル使用（チェイン効果なし）: ClearEffectForAgent（削除）
```

**スキル使用時の実行順序（変更なし）:**
```
1. startAgentRecast(agentIndex, module)           ← スキル選択直後: 新チェイン効果を登録（上書き）or 既存チェイン効果を削除
2. タイピングチャレンジ実行
3. handleChallengeComplete:
   3a. triggerChainEffects(agentIndex, effectFlags) ← 他エージェントの待機中チェイン効果を発動
   3b. スキル効果パイプライン実行
```

## メモ
- `battle_view.go`のUI表示条件(`isChainActive`)は現在`pendingChain != nil && pendingChain.Effect.Type == slot.Module.ChainEffect.Type`であり、リキャスト状態をチェックしていない。そのため、`ExpireEffectsForAgent`呼び出しの削除だけで自然にUI表示が維持される
- `startAgentRecast`はスキル選択直後（チャレンジ開始前）に呼ばれ、`triggerChainEffects`は`handleChallengeComplete`内で呼ばれる。`startAgentRecast`は使用エージェント自身のチェイン効果のみを操作し、`triggerChainEffects`は他エージェントのチェイン効果を発動するため、互いに干渉しない
- `ClearAll()`はバトル終了時に呼ばれる既存処理のため影響なし
- `UpdateRecasts()`で`completedAgents`変数が未使用になるコンパイルエラーに注意。戻り値を受け取らない形`s.recastManager.UpdateRecast(delta)`に変更する
