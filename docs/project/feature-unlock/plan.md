# 実装計画: 機能解放システム

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル | 状態 |
|---|------|-----------|---------------|------|
| 1 | 全機能がLocked初期化 | Go test | `internal/domain/feature_unlock_test.go` | 作成済 |
| 2-5 | ランク到達→PendingTutorial | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 6 | 冪等性 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 7 | ランク跨ぎ | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 8 | 単調増加 | Go test | `internal/domain/feature_unlock_test.go` | 作成済 |
| 9 | 不正遷移エラー | Go test | `internal/domain/feature_unlock_test.go` + Manager | 作成済 |
| 10 | 報酬画面チュートリアル誘導 | TUI test | セッション内 | 手順記載済 |
| 11 | Enter→チュートリアル遷移 | TUI test | セッション内 | 手順記載済 |
| 12 | チュートリアル完了→Unlocked | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 13 | ランク昇順で順次表示 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 14 | 全完了→ホーム遷移 | TUI test | セッション内 | 手順記載済 |
| 15 | PendingなしのRewardは従来通り | TUI test | セッション内 | 手順記載済 |
| 16 | defense_skillゲート | TUI test | セッション内 | 手順記載済 |
| 17 | agent_customizationゲート | TUI test | セッション内 | 手順記載済 |
| 18 | chain_effectゲート | TUI test | セッション内 | 手順記載済 |
| 19 | mana_systemゲート | TUI test | セッション内 | 手順記載済 |
| 20 | TIPSメニュー追加 | TUI test | セッション内 | 手順記載済 |
| 21 | 初期チュートリアル閲覧可能 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 22-23 | 解放済みチュートリアルのみ表示 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 24 | TIPS再閲覧は状態遷移なし | TUI test | セッション内 | 手順記載済 |
| 25-26 | セーブ/ロード永続化 | Go test | `internal/integration_test/feature_unlock_test.go` | 作成済 |
| 27 | 旧セーブ互換 | Go test | `internal/integration_test/feature_unlock_test.go` | 作成済 |
| 28 | 再ロック防止 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |
| 29 | 未知FeatureID保持 | Go test | `internal/integration_test/feature_unlock_test.go` | 作成済 |
| 30 | Reconcile処理 | Go test | `internal/usecase/unlocking/feature_unlock_manager_test.go` | 作成済 |

## 実装タスク

### タスク1: Domain層 - FeatureUnlockState（型定義と状態遷移）
- **対象**: `internal/domain/feature_unlock.go`（新規）
- **内容**:
  - `FeatureID` 型（string）、定数（FeatureDefenseSkill, FeatureAgentCustomization, FeatureChainEffect, FeatureManaSystem）
  - `FeatureStatus` 型（uint8）、定数（FeatureLocked, FeaturePendingTutorial, FeatureUnlocked）
  - `UnlockRule` 構造体（FeatureID, UnlockRank, TutorialID）
  - `TutorialDef` 構造体（ID, Title, Pages, DefaultVisible, FeatureID）
  - `FeatureUnlockState` 構造体（features map[FeatureID]FeatureStatus）
  - `NewFeatureUnlockState()`, `GetStatus()`, `CanTransition()`, `TransitionTo()`
- **関連テスト**: 基準1, 8, 9
- **状態**: 未着手

### タスク2: Usecase層 - FeatureUnlockManager（コアロジック）
- **対象**: `internal/usecase/unlocking/feature_unlock_manager.go`（新規）
- **内容**:
  - `Manager` 構造体（rules, tutorials, state, lastAppliedRank, pendingQueue）
  - `UnlockDelta` 構造体（NewPendingFeatures, QueuedTutorials）
  - `NewManager()`, `ApplyRank()`, `Reconcile()`, `NextPendingTutorial()`, `CompleteTutorial()`
  - `IsUnlocked()`, `IsPendingOrUnlocked()`, `ListVisibleTutorials()`, `Snapshot()`
  - map/sliceのディープコピー、ランク昇順の決定的順序
  - **注意**: セーブデータ型（infra/savedata）への依存は禁止。Snapshotはdomain型のみを返す
- **関連テスト**: 基準2-7, 9, 12, 13, 21-23, 28, 30
- **状態**: 未着手

### タスク3: Infra層 - マスタデータ（JSON + ローダー）
- **対象**:
  - `internal/infra/masterdata/data/feature_unlocks.json`（新規）
  - `internal/infra/masterdata/data/tutorials.json`（新規）
  - `internal/infra/masterdata/loader.go`（変更）
- **内容**:
  - `feature_unlocks.json`: ランク→機能マッピング定義
  - `tutorials.json`: チュートリアル内容定義（ID, タイトル, ページ, default_visible, feature_id）
  - `FeatureUnlockData`, `TutorialData` DTO構造体
  - `LoadFeatureUnlocks()`, `LoadTutorials()` メソッド
  - `ExternalData` に `FeatureUnlocks`, `Tutorials` フィールド追加
  - `LoadAllExternalData()` に統合
  - バリデーション: feature_id重複、tutorial_id重複、参照不整合チェック
- **関連テスト**: インフラ層テスト
- **状態**: 未着手

### タスク4: Infra層 - セーブデータ拡張
- **対象**: `internal/infra/savedata/savedata.go`（変更）
- **内容**:
  - `FeatureUnlockSave` 構造体（Features map[string]string, LastAppliedRank int）
  - `SaveData` に `FeatureUnlock *FeatureUnlockSave` フィールド追加（`json:"feature_unlock,omitempty"`）
  - `NewSaveData()` の更新（nil初期化、omitemptyで旧セーブ互換）
- **関連テスト**: 基準25, 27, 29
- **状態**: 未着手

### タスク5: App層 - マスタデータ変換 + セーブデータアダプター
- **対象**:
  - `internal/app/masterdata_converter.go`（変更）
  - `internal/app/feature_unlock_adapter.go`（新規）
- **内容**:
  - `ConvertFeatureUnlocks()`: FeatureUnlockData → []domain.UnlockRule（バリデーション含む）
  - `ConvertTutorials()`: TutorialData → []domain.TutorialDef（バリデーション含む）
  - `FeatureUnlockSnapshotToSave()`: domain.FeatureUnlockState → savedata.FeatureUnlockSave（app層でinfra↔domain変換）
  - `FeatureUnlockSaveToSnapshot()`: savedata.FeatureUnlockSave → domain.FeatureUnlockState
  - `RebuildFeatureUnlockFromRank()`: 旧セーブ互換用（ランクからUnlockState再構築）
- **関連テスト**: 変換バリデーションテスト、基準27
- **状態**: 未着手

### タスク6: App層 - RootModel統合（初期化 + セーブ/ロード）
- **対象**:
  - `internal/app/root_model.go`（変更）
- **内容**:
  - RootModelに `unlockManager *unlocking.Manager` フィールド追加
  - `NewRootModel()`: マスタデータ変換→セーブデータ復元（app層で判定: FeatureUnlock==nil時はRebuildFromRank）→Manager初期化→Reconcile実行
  - セーブ処理: `ToSaveData()`呼び出し後にapp層でFeatureUnlockSnapshotToSave()を追加
  - **注意**: オーケストレーション（復元/Reconcile/セーブ統合）はすべてapp層で実行。usecase/sessionには変更不要
- **関連テスト**: 基準25-30
- **状態**: 未着手

### タスク7: App層 - シーン基盤（Scene/Router/Factory/Map）
- **対象**:
  - `internal/app/scene.go`（変更）
  - `internal/app/scene_router.go`（変更）
  - `internal/app/screen_factory.go`（変更）
  - `internal/app/screen_map.go`（変更）
- **内容**:
  - `SceneTips`, `SceneTutorial` 追加
  - ルーティングマップに `"tips"`, `"tutorial"` 追加
  - `CreateTipsScreen()`, `CreateTutorialScreen()` ファクトリメソッド
  - ScreenMapに新シーンのマッピング
- **関連テスト**: シーン遷移テスト
- **状態**: 未着手

### タスク8: TUI層 - チュートリアル画面
- **対象**: `internal/tui/screens/tutorial.go`（新規）
- **内容**:
  - `TutorialScreen` 構造体（tutorial domain.TutorialDef, mode, currentPage）
  - `TutorialMode` 型（UnlockFlow / TipsView）
  - ページ送り（←→またはEnter）
  - UnlockFlowモード: 最終ページでEnter/Escで `CompleteTutorialMsg` 送出（screens内でメッセージ定義）
  - TipsViewモード: Escで `ChangeSceneMsg` 送出（TIPS一覧に戻る）
  - **注意**: app型への依存禁止。メッセージはscreens内で定義し、app側でハンドリング
- **関連テスト**: 基準12, 24
- **状態**: 未着手

### タスク9: TUI層 - TIPS画面 + プレゼンター
- **対象**:
  - `internal/tui/screens/tips.go`（新規）
  - `internal/tui/presenter/tutorial_presenter.go`（新規）
- **内容**:
  - `TipsScreen`: チュートリアル一覧表示（リスト選択UI）、選択→`OpenTutorialMsg`送出
  - `TutorialPresenter`: domain.TutorialDefリスト→表示用データ整形（タイトル、説明等）
  - **注意**: TipsScreenはdomain型のTutorialDefリストを受け取る。Providerインターフェースは不要（データはapp層から注入）
- **関連テスト**: 基準20-24
- **状態**: 未着手

### タスク10: App層 - メッセージハンドラー + バトル結果フロー
- **対象**:
  - `internal/app/message_handlers.go`（変更）
  - `internal/app/root_model.go`（変更）
- **内容**:
  - `handleBattleResultMsg()` 変更: ランクアップ時にApplyRank呼び出し、PendingTutorial情報をRewardScreenに渡す
  - `CompleteTutorialMsg` ハンドラー: CompleteTutorial → セーブ → 次Tutorial or Home遷移
  - `OpenTutorialMsg` ハンドラー: TipsViewモードでTutorialScreen遷移
  - Reward→Tutorial→Homeの遷移オーケストレーション
- **関連テスト**: 基準10-15, 26, 30
- **状態**: 完了

### タスク11: TUI層 - 報酬画面変更
- **対象**: `internal/tui/screens/reward.go`（変更）
- **内容**:
  - PendingTutorial情報（機能名リスト）をRewardScreenに渡す（RewardResultへの混入は避け、別フィールドで注入）
  - チュートリアル誘導セクションの追加（解放される機能名 + 「Enterで詳細を見る」キーガイド）
  - PendingTutorial時: Enter→`OpenTutorialMsg`、Space/Esc→ホーム遷移
  - PendingTutorialなし時: 従来通り（Enter/Space/Escでホーム遷移）
- **関連テスト**: 基準10, 11, 15
- **状態**: 未着手

### タスク12: TUI層 - ホーム画面変更
- **対象**: `internal/tui/screens/home.go`（変更）
- **内容**:
  - メニューに「TIPS」項目追加
  - `agent_customization` 未解放時: エージェント管理系メニューを非表示
  - `FeatureUnlockProvider` インターフェース定義（`IsUnlocked(FeatureID) bool`）をscreens内で定義
  - NewHomeScreenの引数にProvider追加
- **関連テスト**: 基準17, 20
- **状態**: 未着手

### タスク13: TUI層 - バトル画面 + エージェントカスタマイズ画面ゲート
- **対象**:
  - `internal/tui/screens/battle_logic.go`（変更）
  - `internal/tui/screens/battle_view.go`（変更）
  - `internal/tui/screens/agent_customization.go`（変更）
- **内容**:
  - `defense_skill` 未解放時: ディフェンスタイプスキルをスキル選択リストからフィルタ
  - `mana_system` 未解放時: マナ表示（⭐）非表示、ManaCost判定スキップ（全スキル使用可能扱い）
  - `chain_effect` 未解放時: チェイン効果スロット非表示
  - ゲート判定は `FeatureUnlockProvider` インターフェース経由で一元的に参照
- **関連テスト**: 基準16, 18, 19
- **状態**: 未着手

## TUIテスト手順

### 基準10,11,14: 報酬画面→チュートリアル→ホーム遷移フロー
```
# ランク2到達直後のセーブデータで起動
# （テストデータはタスク10実装時に作成）
launch_tui(command="go run ./cmd/BlitzTypingOperator -save testdata/rank2_just_reached.json", mode="buffer", dimensions="160x45")

# 起動後、報酬画面を表示（ランクアップ発生済み、PendingTutorialあり）
# 2秒待機（初回コンパイル考慮）
send_keys("", delay=2)

# 報酬画面にチュートリアル誘導セクションが表示されることを確認
assert_contains("新機能が解放されました")
assert_contains("ディフェンススキル")
assert_contains("Enterで詳細を見る")

# Enterでチュートリアルに遷移
send_keys("\r", delay=1, inter_key_delay=0.05)

# チュートリアル画面にタイトルと内容が表示される
assert_contains("ディフェンススキル")

# チュートリアル完了（Enter）
send_keys("\r", delay=1, inter_key_delay=0.05)

# ホーム画面に戻ることを確認
assert_contains("TIPS")
assert_contains("バトル")

close_session()
```

### 基準15: PendingなしのRewardは従来通り
```
# 全機能解放済みの状態でバトル勝利→報酬画面
launch_tui(command="go run ./cmd/BlitzTypingOperator -save testdata/all_unlocked.json", mode="buffer", dimensions="160x45")
send_keys("", delay=2)

# 報酬画面にチュートリアル誘導が表示されないこと
# 「新機能が解放されました」が含まれないことを確認
# Enterでホーム遷移（従来通り）
send_keys("\r", delay=1, inter_key_delay=0.05)
assert_contains("TIPS")

close_session()
```

### 基準16: defense_skillゲート（バトル中）
```
# 新規ゲーム（全機能Locked）で起動
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("", delay=2)

# バトル選択→バトル開始
# ホーム画面からバトル選択へ
send_keys("jj\r", delay=1, inter_key_delay=0.05)  # バトル選択まで移動

# バトル画面でスキル一覧を確認
# → ディフェンスタイプスキルが表示されないこと
# （「ディフェンス」「defense」系のスキル名が含まれない）

close_session()
```

### 基準17: agent_customizationゲート（ホーム画面）
```
# 新規ゲーム（全機能Locked）で起動
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("", delay=2)

# ホーム画面のメニューを確認
# 「エージェントカスタマイズ」が含まれないこと
# 「TIPS」が含まれること
assert_contains("TIPS")
assert_contains("バトル")
# 注: エージェントカスタマイズが非表示であることはTUI表示の目視確認

close_session()
```

### 基準19: mana_systemゲート（バトル中）
```
# 新規ゲーム（mana_system未解放）でバトル開始
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("", delay=2)

# バトル画面に遷移
# → マナ表示（⭐）が非表示であること
# → ManaCost>0のスキルもマナ制限なく選択可能であること

close_session()
```

### 基準20,24: TIPS画面（閲覧と再閲覧）
```
# 新規ゲームで起動
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("", delay=2)

# ホーム画面でTIPSメニューを選択
# TIPSのメニュー位置まで移動してEnter
send_keys("j\r", delay=1, inter_key_delay=0.05)

# TIPS画面に初期チュートリアル（基本操作）が表示される
assert_contains("基本操作")

# チュートリアルを選択して閲覧
send_keys("\r", delay=1, inter_key_delay=0.05)
# チュートリアル内容が表示される
assert_contains("基本操作")

# Escで一覧に戻る（状態遷移なし）
send_keys("\x1b", delay=1, inter_key_delay=0.05)
# TIPS一覧に戻る
assert_contains("基本操作")

# もう一度Escでホームに戻る
send_keys("\x1b", delay=1, inter_key_delay=0.05)
assert_contains("TIPS")

close_session()
```

## 進捗ログ

- タスク1-9: 前セッションで完了（Domain/Usecase/Infra/App/TUI各層の基盤実装）
- タスク10: App層メッセージハンドラー完了
  - CompleteTutorialMsg/OpenTutorialMsgハンドラーをmessage_handlers.goに追加
  - handleBattleResultにApplyRank呼び出し + RewardScreen.SetPendingTutorials追加
  - reward.goにpendingTutorialsフィールド、チュートリアル誘導セクション、Enter→CompleteTutorialMsg遷移を追加
  - prepareSceneTransitionに"tips"ケースを追加
  - 全テスト合格
