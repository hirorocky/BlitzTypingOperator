# ランクアップ報酬

## 概要
ランクアップ時にマスタデータで定義された報酬アイテム（コア・スキル・チェイン効果）を獲得できるようにする。
報酬画面をリデザインし、撃破報酬・ランクアップ・タイピング統計を縦並びレイアウトに変更する。

## 受け入れ基準
1. バトル勝利時に同ランク全敵が撃破され次ランクが解放された場合、マスタデータ（`rank_rewards.json`）で定義された報酬アイテムが自動的にインベントリに追加される
2. 報酬画面にランクアップ報酬の内容（コア名・スキル名・チェイン効果名）がJSON定義順で表示される
3. ランクアップしなかった場合、ランクアップセクションは表示されない
4. 報酬画面のレイアウトが以下の順序で縦並びに表示される（セクション順序固定）:
   a. 撃破報酬セクション（ドロップアイテム1つ + HPアップ）※内容がある場合のみ表示
   b. ランクアップセクション（RANKUP ASCII + ランクアップ報酬アイテム一覧）※ランクアップ時のみ表示
   c. バトル中のタイピング統計セクション（常に表示）
5. `rank_rewards.json`で報酬が定義されていないランクでランクアップした場合、RANKUP ASCIIのみ表示し報酬アイテムは表示しない
6. ランクアップ報酬で獲得したアイテムはインベントリに正しく追加される（コア→CoreInventory、スキル→SkillInventory、チェイン効果→ChainEffectInventory）
7. ランクアップ報酬で獲得したスキルにはチェイン効果を付与しない
8. 撃破報酬もランクアップもHP増加もない場合、タイピング統計セクションのみ表示される

## 設計
### 撃破報酬の前提
- 撃破報酬のドロップアイテムは**コア・スキル・チェイン効果のいずれか1つ、もしくは無し**が基本
- HP増加も発生しない場合がある（同レベル以下の再撃破時）
- 撃破報酬セクションはドロップアイテムまたはHP増加のいずれかがある場合のみ表示

### 変更対象
- `internal/infra/masterdata/loader.go`: ランクアップ報酬データのロード関数`LoadRankRewards()`追加
- `internal/infra/masterdata/data/`: `rank_rewards.json`追加
- `internal/app/masterdata_converter.go`: ランクアップ報酬のドメイン型変換`ConvertRankRewards()`追加
- `internal/app/root_model.go`: バトル結果処理でランクアップ報酬を計算・インベントリ追加
- `internal/usecase/rewarding/reward.go`: `RewardResult`にランクアップ報酬フィールド追加（`RankUpRewardCores`、`RankUpRewardSkills`、`RankUpRewardChainEffects`）。`RewardCalculator`に`SetRankRewards()`メソッド追加（`SetChainEffectPool()`と同パターン）
- `internal/tui/screens/reward.go`: 報酬画面の縦並びリデザイン＋ランクアップ報酬表示＋セクション条件付き表示

### 新規作成
- `internal/infra/masterdata/data/rank_rewards.json`: ランクアップ報酬マスタデータ
- `internal/domain/rank_reward.go`: `RankReward`ドメインモデル（軽量VO。Rank + RewardItems）

### ドメインモデル
```go
// RankReward はランクアップ報酬を表すVO
type RankReward struct {
    Rank  int              // 対象ランク
    Items []RankRewardItem // 報酬アイテムリスト
}

// RankRewardItem はランクアップ報酬の個別アイテム
type RankRewardItem struct {
    Category string // "core", "skill", "chain_effect"
    TypeID   string // アイテムのTypeID
}
```

### データフロー
1. 起動時: `rank_rewards.json` → `masterdata.LoadRankRewards()` → `app.ConvertRankRewards()` → `map[int]domain.RankReward`（ランク→報酬）
2. バトル勝利時: EnemyProgressがランク解放判定を実施 → 新ランク解放時に`RewardCalculator`がマスタデータから報酬を取得 → `RewardResult.RankUpRewardCores`/`RankUpRewardSkills`/`RankUpRewardChainEffects`に格納
3. `root_model.go`: `RewardResult`からランクアップ報酬をインベントリに追加（既存`AddRewardsToInventory()`を拡張）
4. `reward.go`（TUI）: 各セクションの表示条件を判定し、内容があるセクションのみ描画

### `rank_rewards.json`の形式
```json
{
  "rank_rewards": [
    {
      "rank": 2,
      "rewards": [
        { "category": "core", "type_id": "magic_balance" },
        { "category": "skill", "type_id": "heal_lv1" },
        { "category": "chain_effect", "type_id": "chain_heal" }
      ]
    }
  ]
}
```

許可されるcategory値: `"core"`, `"skill"`, `"chain_effect"`

### エラーハンドリング
- 未知`type_id`: ログ警告してスキップ
- 未知`category`: ログ警告してスキップ
- 不正ランク（0以下）: ログ警告してスキップ
- 空rewards配列: 報酬なし（正常系、受け入れ基準5に該当）

### 報酬画面の新レイアウト

**全セクション表示時（ドロップ + ランクアップ + 統計）:**
```
┌─────────────────────────────────┐
│     撃破報酬                      │
│  コア: xxx                       │
│  HPアップ: +10 (1000 → 1010)     │
├─────────────────────────────────┤
│     ✨ RANK UP! ✨               │
│      [影付き数字ASCII]            │
│     ランク 1 → 2                  │
│     ランクアップ報酬:              │
│       コア: xxx                   │
│       スキル: xxx                 │
├─────────────────────────────────┤
│     タイピング統計                 │
│  平均WPM / 正確性 / ダメージ      │
└─────────────────────────────────┘
         Enter: 続行
```

**撃破報酬なし + ランクアップあり:**
```
┌─────────────────────────────────┐
│     ✨ RANK UP! ✨               │
│      [影付き数字ASCII]            │
│     ランク 1 → 2                  │
│     ランクアップ報酬:              │
│       スキル: xxx                 │
├─────────────────────────────────┤
│     タイピング統計                 │
│  平均WPM / 正確性 / ダメージ      │
└─────────────────────────────────┘
         Enter: 続行
```

**統計のみ（撃破報酬なし + ランクアップなし）:**
```
┌─────────────────────────────────┐
│     タイピング統計                 │
│  平均WPM / 正確性 / ダメージ      │
└─────────────────────────────────┘
         Enter: 続行
```

### セクション表示条件
| セクション | 表示条件 |
|-----------|---------|
| 撃破報酬 | ドロップアイテム >= 1 または HPGain > 0 |
| ランクアップ | RankUnlocked == true |
| タイピング統計 | 常に表示 |

## メモ
- 撃破報酬のドロップはコア/スキル/チェイン効果のいずれか1つが基本（敵タイプのDropItemCategory/DropItemTypeIDで決定）
- `RewardResult`のランクアップ報酬フィールドは既存パターンに合わせて型別に分離（`RankUpRewardCores []*domain.CoreModel`、`RankUpRewardSkills []*domain.SkillModel`、`RankUpRewardChainEffects`）
- `RewardCalculator.SetRankRewards()`は既存の`SetChainEffectPool()`と同じManagerパターン
- 現在の報酬画面は左右並び（統計+ドロップ）だが、新レイアウトは全て縦並びにする
