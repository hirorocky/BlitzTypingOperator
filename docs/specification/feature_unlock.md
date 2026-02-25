# 機能解放システム

## 概要

ランクアップに応じてゲーム機能を段階的に解放するシステム。初回プレイ時の情報過多を防ぎ、ランク進行に伴い新機能をチュートリアル付きで導入する。解放状態はセーブデータで永続化し、メインメニューの「TIPS」から閲覧済みチュートリアルを再閲覧可能。

**実装**: `/internal/domain/feature_unlock.go`、`/internal/usecase/unlocking/`

## 要件

### REQ-UNLOCK-1: 状態モデル
**種別**: Ubiquitous

機能の解放状態は3段階で管理する:
```
Locked ──(ランクアップ)──→ PendingTutorial ──(チュートリアル完了)──→ Unlocked
```

**受け入れ基準**:
1. 新規ゲーム開始時、全機能は `Locked` 状態で初期化
2. 状態遷移は単調増加のみ（Locked → PendingTutorial → Unlocked）。ダウングレード不可
3. 不正な遷移要求はエラーを返す

### REQ-UNLOCK-2: ランクに応じた解放
**種別**: Event-Driven

When ランクアップが発生する, the system shall:
- 該当ランク以下の未解放機能を `PendingTutorial` に遷移

**解放ランク定義**（`feature_unlocks.json`で管理）:

| ランク | 解放される機能 |
|--------|---------------|
| 2 | agent_customization（エージェントカスタマイズ） |
| 3 | defense_skill（ディフェンススキル） |
| 4 | chain_effect（チェイン効果） |
| 5 | mana_system（マナシステム） |
| 6 | latent_effect（潜在効果） |

**受け入れ基準**:
1. 解放処理は冪等（同一ランクで再度呼ばれても二重解放されない）
2. ランクを一気に跨いだ場合（例: 1→4）、中間ランクの機能も漏れなく解放

### REQ-UNLOCK-3: チュートリアルフロー
**種別**: Event-Driven

When `PendingTutorial` の機能がある場合, the system shall:
- 報酬画面にチュートリアル誘導セクションを表示
- Enterでチュートリアル画面へ遷移
- チュートリアル完了後に機能を `Unlocked` に遷移

**受け入れ基準**:
1. 報酬画面でEnterを押すと最初の `PendingTutorial` 機能のチュートリアルに遷移
2. チュートリアル画面でEnterまたはEscで対象機能が `Unlocked` に遷移（UnlockFlowモード）
3. 複数機能が同時に `PendingTutorial` の場合、ランク昇順で1つずつ順にチュートリアルを表示
4. 全チュートリアル完了後、ホーム画面に遷移
5. `PendingTutorial` がない場合の報酬画面は従来通り（Enter/Space/Escでホーム遷移）

### REQ-UNLOCK-4: 機能ゲート
**種別**: Ubiquitous

`Locked` と `PendingTutorial` の両方を「未解放」として扱い（`Unlocked` のみが「解放済み」）、UIレベルでゲート制御する。

**受け入れ基準**:
1. `defense_skill` 未解放時: バトル中のディフェンスタイプスキルが非表示かつ使用不可
2. `agent_customization` 未解放時: ホームメニューの「エージェント管理」が非表示
3. `chain_effect` 未解放時: エージェントカスタマイズ画面でチェイン効果スロットが非表示
4. `mana_system` 未解放時: バトル中のマナ表示が非表示、ManaCostによるスキル制限が無効（全スキル使用可能）
5. `latent_effect` 未解放時: パーフェクトタイピング時でもIsLatent=trueの効果が発動しない（パーフェクト演出は表示される）

### REQ-UNLOCK-5: TIPS画面
**種別**: State-Based

ホームメニューの「TIPS」からチュートリアルを再閲覧可能にする。

**受け入れ基準**:
1. ホームメニューに「TIPS」項目を表示
2. 基本操作チュートリアル（`default_visible: true`）は初期状態から閲覧可能
3. 解放済み（`Unlocked`）機能のチュートリアルはTIPSから再閲覧可能
4. 未解放機能のチュートリアルはTIPSに表示されない
5. TIPS画面でのチュートリアル閲覧は再閲覧モード（TipsViewモード、状態遷移は発生しない）

### REQ-UNLOCK-6: セーブ/ロード
**種別**: Event-Driven

**受け入れ基準**:
1. 機能解放状態（各FeatureIDのステータス）がセーブデータに永続化される
2. セーブはチュートリアル完了時に実行
3. 旧セーブデータ（FeatureUnlockフィールドなし）のロード時、CurrentRankから解放状態を再構築（該当ランク以下の機能は `Unlocked` 扱い）
4. 一度 `Unlocked` になった機能は、マスタデータ変更でも再ロックされない
5. 未知のFeatureIDがセーブデータに含まれる場合、無視して保持
6. 起動時にReconcile処理を実行し、マスタデータ追加分を反映

## 仕様

### FeatureID / FeatureStatus

**責務**: 機能の識別と解放状態の表現

**定義**:
```go
type FeatureID string // "defense_skill", "agent_customization", "chain_effect", "mana_system", "latent_effect"

const (
    FeatureLocked          FeatureStatus = iota // 未解放
    FeaturePendingTutorial                      // チュートリアル未完了
    FeatureUnlocked                             // 解放済み
)
```

### FeatureUnlockState

**責務**: 全機能の解放状態を管理する値オブジェクト

**ルール**:
1. 未登録の機能は `Locked` として扱う
2. 状態遷移は `CanTransition` で検証（単調増加のみ）
3. Clone/AllFeaturesはディープコピーを返す

### Manager（usecase/unlocking）

**責務**: 機能解放の状態遷移ロジックを管理

**インターフェース**:
- `ApplyRank(rank)`: 指定ランクまでの未解放機能をPendingTutorial化（冪等）
- `Reconcile(currentRank)`: マスタデータ追加時の既存ユーザー対応
- `NextPendingTutorial()`: キュー先頭のチュートリアルIDを返す
- `CompleteTutorial(tutorialID)`: チュートリアル完了→Unlocked遷移
- `IsUnlocked(featureID)`: 解放済み判定
- `ListVisibleTutorials()`: TIPS画面に表示可能なチュートリアル一覧（ディープコピー）
- `Snapshot()`: 現在の状態のディープコピー

### TutorialScreen

**責務**: チュートリアルの表示（ページ送り対応）

**モード**:
- **UnlockFlow**: 報酬画面からの導線。Enter/Escいずれでも機能が `Unlocked` に遷移
- **TipsView**: TIPS画面からの再閲覧。状態遷移は発生せず、Escで一覧に戻る

### マスタデータ

**feature_unlocks.json**: ランク→機能→チュートリアルのマッピング
**tutorials.json**: チュートリアル内容（ID、タイトル、ページ、初期表示フラグ、対応機能ID）

### アダプター（app/feature_unlock_adapter.go）

**責務**: domain（FeatureUnlockState）とinfra（FeatureUnlockSave）の変換

**関数**:
- `FeatureUnlockSnapshotToSave`: FeatureUnlockState → FeatureUnlockSave
- `FeatureUnlockSaveToSnapshot`: FeatureUnlockSave → FeatureUnlockState
- `RebuildFeatureUnlockFromRank`: ランクから解放状態を再構築（旧セーブ互換）

## 関連ドメイン

- **Game Loop**: バトル勝利時のApplyRank呼び出し、シーン遷移オーケストレーション
- **Battle**: 機能ゲート（ディフェンススキル、マナシステム）
- **Agent**: 機能ゲート（エージェントカスタマイズ、チェイン効果）
