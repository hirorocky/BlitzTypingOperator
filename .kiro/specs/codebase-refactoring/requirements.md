# Requirements Document

## Introduction

本ドキュメントは、BlitzTypingOperator プロジェクトのコードベースリファクタリングに関する要件を定義します。目的は、コードの保守性・可読性・拡張性を向上させ、技術的負債を解消することです。リファクタリングは既存機能の動作を変更せず、内部品質のみを改善します。

## Requirements

### Requirement 1: デッドコード削除

**Objective:** As a 開発者, I want 未使用のコードを削除する, so that コードベースが簡素化され、保守コストが削減される

#### Acceptance Criteria

1. When リファクタリングが完了した時, the codebase shall `internal/app/app.go` ファイルが存在しない状態となる
2. When リファクタリングが完了した時, the codebase shall `app_test.go` から旧 `Model` 構造体関連のテスト（TestNewApp, TestAppImplementsTeaModel, TestAppInit, TestAppUpdate, TestAppView, TestAppViewContainsGameTitle）が削除された状態となる
3. When リファクタリングが完了した時, the domain package shall `agent.go` から未使用メソッド（GetModule, GetModuleCount, GetCoreName）が削除された状態となる
4. When リファクタリングが完了した時, the domain package shall `module.go` から未使用メソッド（IsAttack, IsSupport, TargetsEnemy, TargetsPlayer, GetCategoryTag, DefaultStatRef）が削除された状態となる
5. When リファクタリングが完了した時, the domain package shall `effect_table.go` から未使用メソッド（IsPermanent, IsExpired, FindByID, Clear）が削除された状態となる
6. While デッドコードを削除する時, the refactoring shall 実運用で使用されている `GetCoreTypeName()` および `GetRowsBySource()` を保持する

### Requirement 2: エラーハンドリング改善

**Objective:** As a 開発者, I want エラーを適切にハンドリングする, so that デバッグが容易になり、静かな失敗によるデータ不整合を防止できる

#### Acceptance Criteria

1. The codebase shall `log/slog` パッケージを使用して構造化ログを出力する
2. When `invManager.AddCore()` がエラーを返した時, the system shall エラー内容とコアIDをログに出力する
3. When `agentMgr.AddAgent()` がエラーを返した時, the system shall エラー内容とエージェント情報をログに出力する
4. When `agentMgr.EquipAgent()` がエラーを返した時, the system shall エラー内容とスロット情報をログに出力する
5. When `saveDataIO.SaveGame()` がエラーを返した時, the system shall エラーレベルのログを出力する
6. The codebase shall `_` でエラーを無視する箇所を解消する

### Requirement 3: マジックナンバーの定数化

**Objective:** As a 開発者, I want マジックナンバーを定数として定義する, so that 変更が容易になり、値の意味が明確になる

#### Acceptance Criteria

1. The project shall `internal/config/constants.go` ファイルを持つ
2. The constants file shall バトル設定定数（BattleTickInterval: 100ms, DefaultModuleCooldown: 5.0, AccuracyPenaltyThreshold: 0.5, MinEnemyAttackInterval: 500ms）を定義する
3. The constants file shall 効果持続時間定数（BuffDuration: 10.0, DebuffDuration: 8.0）を定義する
4. The constants file shall インベントリ定数（MaxAgentEquipSlots: 3, ModulesPerAgent: 4）を定義する
5. When 定数化が完了した時, the codebase shall 各ファイル内のハードコード値を定数参照に置き換えた状態となる

### Requirement 4: ファイル分割（root_model.go）

**Objective:** As a 開発者, I want `root_model.go` を責務ごとに分割する, so that 各ファイルの責務が明確になり、保守性が向上する

#### Acceptance Criteria

1. The `internal/app/` directory shall `root_model.go`（RootModel構造体とInit/Update/Viewのみ）を持つ
2. The `internal/app/` directory shall `scene_router.go`（シーン遷移ロジック）を持つ
3. The `internal/app/` directory shall `screen_factory.go`（画面インスタンス生成）を持つ
4. The `internal/app/` directory shall `adapters.go`（inventoryProviderAdapter等のアダプター）を持つ
5. The `internal/app/` directory shall `helpers.go`（createStatsDataFromGameState等のヘルパー関数）を持つ
6. While ファイル分割を行う時, the refactoring shall 既存の機能動作を維持する
7. When 分割が完了した時, the `root_model.go` shall 300行以下となる

### Requirement 5: ファイル分割（game_state.go）

**Objective:** As a 開発者, I want `game_state.go` を責務ごとに分割する, so that 状態管理と変換ロジックが分離され、テストが容易になる

#### Acceptance Criteria

1. The project shall `internal/app/game_state/` ディレクトリを持つ
2. The `game_state/` directory shall `state.go`（GameState本体）を持つ
3. The `game_state/` directory shall `persistence.go`（ToSaveData, FromSaveData）を持つ
4. The `game_state/` directory shall `defaults.go`（getDefaultCoreTypeData等のデフォルト値生成）を持つ
5. While ファイル分割を行う時, the refactoring shall セーブデータの後方互換性を維持する

### Requirement 6: ファイル分割（battle.go）

**Objective:** As a 開発者, I want `battle.go` を責務ごとに分割する, so that UIレンダリングとゲームロジックが分離される

#### Acceptance Criteria

1. When 分割が完了した時, the battle screen shall UIレンダリング責務とゲームロジック責務が別ファイルに分離された状態となる
2. The refactoring shall バトル画面の既存動作を維持する
3. When 分割が完了した時, the main battle file shall 500行以下となる

### Requirement 7: コード重複の解消

**Objective:** As a 開発者, I want 重複コードを共通化する, so that DRY原則に従い、一貫性が確保される

#### Acceptance Criteria

1. The project shall `internal/tui/components/hp_display.go` を持ち、HP表示ロジックを共通化する
2. When HP表示をレンダリングする時, the system shall 共通の `RenderHP()` 関数を使用する
3. The `domain/module.go` shall `Icon()` メソッドを持ち、モジュールカテゴリアイコンを返す
4. When デフォルトデータを生成する時, the system shall 単一のソースからデフォルト値を取得する（`game_state.go` と `root_model.go` の重複解消）

### Requirement 8: 型名の整理

**Objective:** As a 開発者, I want 型名を適切な名前に変更する, so that コードの可読性が向上し、誤解を防ぐ

#### Acceptance Criteria

1. When リファクタリングが完了した時, the screens package shall `EncyclopediaTestData` を `EncyclopediaData` に改名した状態となる
2. When リファクタリングが完了した時, the screens package shall `StatsTestData` を `StatsData` に改名した状態となる
3. When リファクタリングが完了した時, the codebase shall `*TestData` という名前の本番用型が存在しない状態となる

### Requirement 9: Screenインターフェースの導入

**Objective:** As a 開発者, I want 明示的なScreenインターフェースを導入する, so that 画面実装の一貫性が保証され、拡張が容易になる

#### Acceptance Criteria

1. The project shall `internal/tui/screens/types.go` を持ち、Screenインターフェースを定義する
2. The Screen interface shall `tea.Model` を埋め込み、`SetSize(width, height int)` と `GetTitle() string` メソッドを持つ
3. The project shall `BaseScreen` 構造体を持ち、共通の画面機能を提供する
4. Where 新しい画面を作成する時, the screen shall Screenインターフェースを実装する

### Requirement 10: データ変換層の作成

**Objective:** As a 開発者, I want データ変換ロジックを専用層に集約する, so that ドメインモデルとUI/永続化層の境界が明確になる

#### Acceptance Criteria

1. The project shall `internal/adapter/` ディレクトリを持つ
2. The adapter directory shall `persistence_adapter.go`（SaveData <-> GameState変換）を持つ
3. The adapter directory shall `screen_adapter.go`（GameState -> 各種ScreenData変換）を持つ
4. The adapter directory shall `reward_adapter.go`（BattleStats -> RewardStats変換）を持つ
5. When データ変換を行う時, the system shall adapter層の関数を使用する

### Requirement 11: 循環的複雑度の削減

**Objective:** As a 開発者, I want 大きなswitch文を持つメソッドを改善する, so that コードの可読性と保守性が向上する

#### Acceptance Criteria

1. When メッセージハンドリングを行う時, the RootModel shall ハンドラーマップを使用して分岐を委譲する
2. When シーンレンダリングを行う時, the RootModel shall 画面マップを使用して描画を委譲する
3. When リファクタリングが完了した時, the Update method shall switch分岐数が5以下となる

### Requirement 12: 後方互換性の維持

**Objective:** As a ユーザー, I want リファクタリング後も既存のセーブデータが使用できる, so that ゲームの進行状況が失われない

#### Acceptance Criteria

1. While リファクタリングを行う時, the system shall 既存のセーブデータ形式を変更しない
2. When リファクタリング後にゲームを起動した時, the system shall 既存のセーブデータを正常に読み込む
3. The refactoring shall 外部から観測可能な動作を変更しない
