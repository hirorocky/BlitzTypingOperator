# 実装計画: 潜在効果システム

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル/手順 | 状態 |
|---|------|-----------|-------------------|------|
| 1 | パーフェクト判定（Success + Accuracy>=1.0） | TUI test | セッション内 | 手順記載済 |
| 2 | ディフェンスタイプ除外 | TUI test | セッション内 | 手順記載済 |
| 3 | コンボカウントと同一タイミング | Go test（既存テストで保証） | - | 既存 |
| 4 | パーフェクトASCIIアート表示 | TUI test | セッション内 | 手順記載済 |
| 5 | 約0.5秒後にASCIIアート消去 | TUI test | セッション内 | 手順記載済 |
| 6 | 解放状態に関係なく演出表示 | TUI test | セッション内 | 手順記載済 |
| 7 | IsLatent=trueはパーフェクト時のみ発動 | Go test | `combat/battle_latent_effect_test.go` | 作成済 |
| 8 | IsLatent=falseは常に発動 | Go test | `combat/battle_latent_effect_test.go` | 作成済 |
| 9 | LUK補正式の修正 | Go test | `domain/skill_effect_latent_test.go` | 作成済 |
| 10 | 複数潜在効果の独立判定 | Go test | `combat/battle_latent_effect_test.go` | 作成済 |
| 11 | Echo/DoubleCastでパーフェクト引き継ぎ | Go test | `combat/battle_latent_effect_test.go` | 作成済 |
| 12 | FeatureLatentEffect定数定義 | Go test | `domain/skill_effect_latent_test.go` | 作成済 |
| 13 | 未解放時に潜在効果非発動 | TUI test | セッション内 | 手順記載済 |
| 14 | 解放時チュートリアル表示 | TUI test | セッション内 | 手順記載済 |
| 15 | TUI層でisPerfect=false制御 | Go test（TUI結合テストとして基準7,8で間接検証） | - | 間接検証 |
| 16 | is_latentフィールドのロード | Go test | `masterdata/loader_latent_test.go` | 作成済 |
| 17 | is_latent省略時デフォルトfalse | Go test | `masterdata/loader_latent_test.go` | 作成済 |

## 実装タスク

### タスク1: ドメイン層 - IsLatentフラグとLUK補正修正
- **対象**: `internal/domain/skill_effect.go`, `internal/domain/feature_unlock.go`
- **内容**:
  - `SkillEffect`に`IsLatent bool`フィールドを追加
  - `BaseLUK`定数を削除し、`AdjustedProbability()`の計算式を`probability + luk * lukFactor`に修正
  - `FeatureLatentEffect FeatureID = "latent_effect"`定数を追加
- **関連テスト**: 基準 7, 8, 9, 12
- **状態**: 完了

### タスク2: タイピング結果 - IsPerfectフィールド追加
- **対象**: `internal/usecase/typing/typing.go`
- **内容**:
  - `TypingResult`に`IsPerfect bool`フィールドを追加
- **関連テスト**: 基準 7, 8（間接）
- **状態**: 完了

### タスク3: BattleEngine - 潜在効果フィルタリング
- **対象**: `internal/usecase/combat/battle.go`
- **内容**:
  - `ApplySkillEffect()`のeffectsループ内で、`ShouldTrigger()`の前に`effect.IsLatent && !typingResult.IsPerfect → continue`を追加
- **関連テスト**: 基準 7, 8, 10, 11
- **状態**: 完了

### タスク4: マスタデータ - IsLatentフィールドのロードと変換
- **対象**: `internal/infra/masterdata/loader.go`, `internal/app/masterdata_converter.go`
- **内容**:
  - `SkillEffectData`に`IsLatent bool`フィールド追加（`json:"is_latent"`）
  - `ToDomain()`でIsLatentをマッピング
  - `ResolveSkillTimedEffects()`でIsLatentの保持を確認（既存コードで自動的に保持される）
- **関連テスト**: 基準 16, 17
- **状態**: 完了

### タスク5: マスタデータ - feature_unlocks.jsonとtutorials.json
- **対象**: `internal/infra/masterdata/data/feature_unlocks.json`, `internal/infra/masterdata/data/tutorials.json`
- **内容**:
  - feature_unlocks.jsonにランク6のlatent_effectエントリを追加
  - tutorials.jsonに潜在効果チュートリアルを追加
- **関連テスト**: 基準 12, 14
- **状態**: 完了

### タスク6: ASCIIアート - PERFECT!レンダラー
- **対象**: 新規 `internal/tui/ascii/perfect.go`, `internal/tui/ascii/perfect_test.go`
- **内容**:
  - "PERFECT!" ASCIIアートの定義
  - `PerfectRenderer`インターフェース（`RenderPerfect() string`, `GetWidth() int`, `GetHeight() int`）
  - `NewPerfectRenderer()`コンストラクタ
  - winlose.goパターンに従いlipglossスタイリング
- **関連テスト**: 基準 4（テストファイル内でレンダリング確認）
- **状態**: 完了

### タスク7: TUI層 - パーフェクト判定と演出
- **対象**: `internal/tui/screens/battle.go`, `internal/tui/screens/battle_logic.go`, `internal/tui/screens/battle_view.go`
- **内容**:
  - `BattleScreen`に`showingPerfect bool`, `perfectTimer int`フィールド追加
  - `isLatentEffectUnlocked()`メソッド追加
  - `perfectRenderer ascii.PerfectRenderer`フィールド追加
  - `handleChallengeComplete()`内: コンボカウント判定と同じ`Accuracy >= 1.0`チェックで`isPerfect`判定。ディフェンスタイプは除外。`showingPerfect = true`設定。未解放時は`isPerfect = false`にリセット。`typingResult.IsPerfect`を設定
  - `handleTick()`内: `showingPerfect`時にタイマーをカウントし、約0.5秒後に`showingPerfect = false`に
  - `View()`内: `showingPerfect`時にエージェントエリアの代わりにPERFECT! ASCIIアートを表示
- **関連テスト**: 基準 1, 2, 3, 4, 5, 6, 13, 15
- **状態**: 完了

### タスク8: マスタデータ - スキルへの潜在効果追加
- **対象**: `internal/infra/masterdata/data/skills.json`
- **内容**:
  - lv2/lv3の攻撃・魔法・回復スキルに潜在効果（追加ダメージ/回復）を付与
  - probability: 0.3-0.5, luk_factor: 0.02
- **関連テスト**: 基準 16
- **状態**: 完了

## TUIテスト手順

### 基準1,3: パーフェクト判定（Success + Accuracy>=1.0）
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# ホーム画面 → バトル選択
send_keys("\r", delay=3)
# 敵選択
send_keys("\r", delay=1)
# スキル選択 → Enter
send_keys("\r", delay=1)
# タイピングをミスなしで完了（表示テキストを正確に入力）
# → パーフェクト判定が行われ、PERFECT! が表示されること
```

### 基準4,5: パーフェクトASCIIアート表示と消去
```
# パーフェクト達成後（上記の続き）
# PERFECT! ASCIIアートが中央エリアに表示されていることを確認
assert_contains("PERFECT")
# 約0.5秒後にASCIIアートが消え、エージェントエリアに戻ることを確認
# （次のtickで確認）
```

### 基準2: ディフェンスタイプ除外
```
# ディフェンススキルを選択してチャレンジを完了
# → PERFECT!が表示されないことを確認
# ※ ディフェンスタイプではAccuracyが防御率を表すため
```

### 基準6,13: 解放状態による動作差異
```
# ランク6未満のセーブデータでバトルを開始
# パーフェクト達成時:
# - PERFECT! ASCIIアートは表示される（基準6）
# - 潜在効果は発動しない（基準13）

# ランク6以上のセーブデータでバトルを開始
# パーフェクト達成時:
# - PERFECT! ASCIIアートが表示される
# - 潜在効果が発動する（ダメージ増加等）
```

### 基準14: 解放時チュートリアル
```
# ランク5→6のランクアップ後
# 報酬画面にチュートリアル誘導が表示されること
# Enter → 潜在効果チュートリアルが表示されること
```

## 進捗ログ

- タスク1完了: SkillEffectにIsLatentフィールド追加、BaseLUK定数削除、AdjustedProbability()修正、FeatureLatentEffect定数追加
- タスク2完了: TypingResultにIsPerfectフィールド追加
- タスク3完了: ApplySkillEffect()にIsLatentフィルタリング追加
- タスク4完了: SkillEffectDataにIsLatentフィールド追加、ToDomain()でマッピング
- タスク5完了: feature_unlocks.jsonにランク6のlatent_effect追加、tutorials.jsonに潜在効果チュートリアル追加
- タスク6完了: PERFECT! ASCIIアートレンダラー作成（perfect.go, perfect_test.go）
- タスク7完了: BattleScreenにパーフェクト判定・演出・タイマー・解放チェック実装
- タスク8完了: lv2/lv3スキル5種に潜在効果追加（physical_strike_lv2/lv3, fireball_lv2/lv3, heal_lv2）
