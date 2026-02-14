# 実装計画: ランクアップ報酬

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル | 状態 |
|---|------|-----------|---------------|------|
| 1 | ランクアップ時に報酬がインベントリに追加される | Go test | `internal/usecase/rewarding/rank_reward_test.go` | 作成済(RED) |
| 2 | 報酬画面にランクアップ報酬が表示される | TUI test | セッション内 | 手順記載済 |
| 3 | ランクアップしない場合セクション非表示 | TUI test | セッション内 | 手順記載済 |
| 4 | 報酬画面の縦並びレイアウト | TUI test | セッション内 | 手順記載済 |
| 5 | 報酬未定義ランクでランクアップ時は報酬なし | Go test | `internal/usecase/rewarding/rank_reward_test.go` | 作成済(RED) |
| 6 | ランクアップ報酬のインベントリ追加 | Go test | `internal/usecase/rewarding/rank_reward_test.go` | 作成済(RED) |
| 7 | ランクアップ報酬スキルにチェイン効果なし | Go test | `internal/usecase/rewarding/rank_reward_test.go` | 作成済(RED) |
| 8 | 報酬なし時は統計のみ表示 | TUI test | セッション内 | 手順記載済 |

### Goテスト一覧

**`internal/domain/rank_reward_test.go`** (3テスト):
- `TestRankReward_HasItems`: 3カテゴリのアイテム保持
- `TestRankReward_EmptyItems`: 空報酬リスト
- `TestRankRewardItem_Categories`: カテゴリ値テーブル駆動テスト

**`internal/usecase/rewarding/rank_reward_test.go`** (8テスト):
- `TestCalculateGuaranteedRewardWithProgress_RankUpRewards`: ランクアップ時にコア+スキル報酬が設定される（基準1）
- `TestCalculateGuaranteedRewardWithProgress_RankUpNoRewards`: 報酬未定義ランクで報酬が空（基準5）
- `TestCalculateGuaranteedRewardWithProgress_RankUpSkillNoChainEffect`: チェイン効果非付与（基準7）
- `TestAddRewardsToInventory_RankUpRewards`: ランクアップ報酬がインベントリに追加される（基準6）
- `TestAddRewardsToInventory_DropsAndRankUpRewardsMixed`: 撃破報酬+ランクアップ報酬の混合（基準6）
- `TestAddRewardsToInventory_RankUpChainEffect`: チェイン効果報酬のインベントリ追加（基準6）
- `TestRewardResult_NoDropsNoHPGain`: 撃破報酬なしケース（基準8）

## 実装タスク

### タスク1: RankRewardドメインモデル作成
- **対象**: `internal/domain/rank_reward.go`（新規作成）
- **内容**: `RankReward`と`RankRewardItem`の軽量VOを定義
- **関連テスト**: `internal/domain/rank_reward_test.go`（基準1前提）
- **状態**: 完了

### タスク2: マスタデータローダー追加
- **対象**: `internal/infra/masterdata/loader.go`（変更）、`internal/infra/masterdata/data/rank_rewards.json`（新規作成）
- **内容**: `RankRewardData`構造体定義、`LoadRankRewards()`関数追加、`ExternalData`に`RankRewards`フィールド追加、`LoadAllExternalData()`に統合（オプショナルファイル扱い: ファイルなし時は空配列）
- **関連テスト**: 基準1
- **状態**: 完了

### タスク3: マスタデータ変換関数追加
- **対象**: `internal/app/masterdata_converter.go`（変更）
- **内容**: `ConvertRankRewards()`関数追加。`[]masterdata.RankRewardData` → `map[int]domain.RankReward`変換。未知category/typeIDはログ警告してスキップ
- **関連テスト**: 基準1
- **状態**: 完了

### タスク4: RewardCalculatorにランクアップ報酬ロジック追加
- **対象**: `internal/usecase/rewarding/reward.go`（変更）
- **内容**:
  - `RewardResult`に`RankUpRewardCores`/`RankUpRewardSkills`/`RankUpRewardChainEffects`フィールド追加
  - `RewardCalculator`に`rankRewards`フィールドと`SetRankRewards()`メソッド追加
  - `CalculateGuaranteedRewardWithProgress()`内でランクアップ時にランク報酬を解決してRewardResultに格納
  - ランクアップ報酬のスキルにはチェイン効果を付与しない（`ToDomain()`でnilを渡す）
  - ランクアップ報酬のチェイン効果は`RankUpRewardChainEffects`に格納
- **関連テスト**: 基準1,5,7（`rank_reward_test.go`の3テスト）
- **状態**: 完了

### タスク5: AddRewardsToInventory拡張
- **対象**: `internal/usecase/rewarding/reward.go`（変更）
- **内容**: `AddRewardsToInventory()`でランクアップ報酬フィールド（`RankUpRewardCores`/`RankUpRewardSkills`/`RankUpRewardChainEffects`）もインベントリに追加
- **関連テスト**: 基準6（3テスト: RankUpRewards, Mixed, ChainEffect）
- **状態**: 完了

### タスク6: app層のバトル結果処理統合
- **対象**: `internal/app/root_model.go`（変更）
- **内容**:
  - 起動時に`rank_rewards.json`をロード→変換→`RewardCalculator.SetRankRewards()`で設定
  - `handleBattleResult()`は既存の`CalculateGuaranteedRewardWithProgress()`と`AddRewardsToInventory()`を呼ぶだけ（タスク4,5で対応済み）
- **関連テスト**: 基準1,6
- **状態**: 完了

### タスク7: 報酬画面の縦並びリデザイン
- **対象**: `internal/tui/screens/reward.go`（変更）
- **内容**:
  - `View()`を全面リデザイン: 撃破報酬→ランクアップ→タイピング統計の縦並びレイアウト
  - 撃破報酬セクション: ドロップアイテム1つ + HPアップ。内容がある場合のみ表示
  - ランクアップセクション: RANKUP ASCII + ランク変化 + ランクアップ報酬リスト。`RankUnlocked==true`の場合のみ表示
  - タイピング統計セクション: 常に表示
  - 全セクションなし（統計のみ）のケース対応
- **関連テスト**: 基準2,3,4,5,8（TUIテスト）
- **状態**: 完了

## TUIテスト手順

### 基準2,4: ランクアップ報酬表示 + 縦並びレイアウト
```
# ゲーム起動
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")

# ランク1の最後の未撃破敵を選択・撃破してランクアップを発生させる
# (実際の操作はバトル画面の操作に依存。テスト用セーブデータの準備が必要な場合あり)

# 報酬画面で確認:
# 1. 撃破報酬セクションが上部に表示される
assert_contains("撃破報酬")
# 2. ランクアップセクションが中央に表示される
assert_contains("RANK UP")
# 3. ランクアップ報酬アイテムが表示される
assert_contains("ランクアップ報酬")
# 4. タイピング統計セクションが下部に表示される
assert_contains("タイピング統計")
```

### 基準3: ランクアップなし時のセクション非表示
```
# ランクアップしないバトル（既撃破の敵を再撃破）を実行
# 報酬画面で確認:
assert_contains("タイピング統計")  # 統計は常に表示
# RANK UP が表示されないことを確認（画面全体をキャプチャして目視確認）
capture_screen()
```

### 基準5: 報酬未定義ランクでのランクアップ
```
# rank_rewards.jsonに該当ランクの定義がない状態でランクアップ
# 報酬画面で確認:
assert_contains("RANK UP")        # RANKUP ASCIIは表示
# ランクアップ報酬アイテムは表示されない（目視確認）
capture_screen()
```

### 基準8: 統計のみ表示
```
# 既撃破の敵を同レベル以下で再撃破（HP増加なし、ランクアップなし）
# 報酬画面で確認:
assert_contains("タイピング統計")
# 撃破報酬セクション・ランクアップセクションが表示されないことを目視確認
capture_screen()
```

## 進捗ログ
- タスク1: RankReward/RankRewardItem VOを`internal/domain/rank_reward.go`に作成。3テストGREEN。
- タスク2: `LoadRankRewards()`と`RankRewardData`をloader.goに追加。ExternalDataにRankRewardsフィールド追加。`rank_rewards.json`作成。
- タスク3: `ConvertRankRewards()`をmasterdata_converter.goに追加。未知category/typeIDのログ警告スキップ対応。
- タスク4: RewardResultにRankUpRewardCores/Skills/ChainEffectsフィールド追加。SetRankRewards()とresolveRankUpRewards()実装。基準1,5,7の3テストGREEN。
- タスク5: AddRewardsToInventory()にランクアップ報酬のインベントリ追加処理を実装。基準6の3テストGREEN。
- タスク6: root_model.goの起動時にrank_rewards.jsonをロード→変換→SetRankRewards()で設定。
- タスク7: 報酬画面を縦並びレイアウトにリデザイン。撃破報酬/ランクアップ/タイピング統計の3セクション。条件付き表示対応。
