# Implementation Plan

## Task Overview

本タスクリストは、BlitzTypingOperatorのコードベースリファクタリングを実装するためのものです。既存機能の動作を変更せず、内部品質を改善します。

---

## Tasks

### Phase 1: デッドコード削除と定数化

- [x] 1. デッドコード削除
- [x] 1.1 (P) appパッケージのデッドコードを削除する
  - `internal/app/app.go`ファイルを完全に削除する
  - `app_test.go`から旧Model構造体関連のテスト（TestNewApp, TestAppImplementsTeaModel, TestAppInit, TestAppUpdate, TestAppView, TestAppViewContainsGameTitle）を削除する
  - RootModelが正常に動作することを確認する
  - _Requirements: 1.1, 1.2_

- [x] 1.2 (P) domainパッケージの未使用メソッドを削除する
  - `agent.go`から未使用メソッド（GetModule, GetModuleCount, GetCoreName）を削除する
  - `GetCoreTypeName()`は使用中のため保持する
  - _Requirements: 1.3, 1.6_

- [x] 1.3 (P) moduleとeffect_tableの未使用メソッドを削除する
  - `module.go`から未使用メソッド（IsAttack, IsSupport, TargetsEnemy, TargetsPlayer, GetCategoryTag, DefaultStatRef）を削除する
  - `effect_table.go`から未使用メソッド（IsPermanent, IsExpired, FindByID, Clear）を削除する
  - `GetRowsBySource()`は使用中のため保持する
  - _Requirements: 1.4, 1.5, 1.6_

- [x] 2. 定数の一元管理
- [x] 2.1 定数パッケージを作成する
  - `internal/config/constants.go`を新規作成する
  - バトル設定定数（BattleTickInterval: 100ms, DefaultModuleCooldown: 5.0, AccuracyPenaltyThreshold: 0.5, MinEnemyAttackInterval: 500ms）を定義する
  - 効果持続時間定数（BuffDuration: 10.0, DebuffDuration: 8.0）を定義する
  - インベントリ定数（MaxAgentEquipSlots: 3, ModulesPerAgent: 4）を定義する
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [x] 2.2 ハードコード値を定数参照に置換する
  - battle.go内のマジックナンバーを定数参照に変更する
  - effect関連ファイルの持続時間を定数参照に変更する
  - エージェント/インベントリ関連のスロット数を定数参照に変更する
  - 全テストが通過することを確認する
  - _Requirements: 3.5_

### Phase 2: エラーハンドリング改善

- [x] 3. 構造化ログによるエラーハンドリング
- [x] 3.1 slogを使用したエラーログ出力を実装する
  - log/slogパッケージをインポートする
  - `invManager.AddCore()`のエラーハンドリングを追加し、エラー内容とコアIDをログ出力する
  - `agentMgr.AddAgent()`のエラーハンドリングを追加し、エラー内容とエージェント情報をログ出力する
  - `agentMgr.EquipAgent()`のエラーハンドリングを追加し、エラー内容とスロット情報をログ出力する
  - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 3.2 セーブ処理のエラーハンドリングを改善する
  - `saveDataIO.SaveGame()`がエラーを返した際にslog.Errorでログ出力する
  - `_`でエラーを無視している箇所を全て構造化ログ出力に置換する
  - エラー発生時もアプリケーションが継続動作することを確認する
  - _Requirements: 2.5, 2.6_

### Phase 3: root_model.goの分割

- [x] 4. root_model.goのファイル分割
- [x] 4.1 シーンルーティングロジックを分離する
  - `internal/app/scene_router.go`を新規作成する
  - SceneRouter構造体とハンドラーマップを定義する
  - シーン遷移ロジックをroot_model.goから移動する
  - Route()メソッドでシーン名に基づくルーティングを実装する
  - _Requirements: 4.2, 11.1_

- [x] 4.2 画面生成ロジックを分離する
  - `internal/app/screen_factory.go`を新規作成する
  - ScreenFactory構造体を定義する
  - 各画面の生成メソッド（CreateHomeScreen, CreateBattleSelectScreen, CreateEncyclopediaScreen, CreateStatsScreen等）を実装する
  - root_model.goからの画面生成ロジックを移動する
  - _Requirements: 4.3_

- [x] 4.3 アダプターとヘルパー関数を分離する
  - `internal/app/adapters.go`を新規作成し、inventoryProviderAdapter等を移動する
  - `internal/app/helpers.go`を新規作成し、createStatsDataFromGameState等のヘルパー関数を移動する
  - root_model.goが300行以下になることを確認する
  - 既存テストが全て通過することを確認する
  - _Requirements: 4.4, 4.5, 4.6, 4.7_

### Phase 4: game_state.goの分割

- [x] 5. game_stateパッケージの作成
- [x] 5.1 game_stateディレクトリ構造を作成する
  - `internal/app/game_state/`ディレクトリを新規作成する
  - `state.go`にGameState構造体本体とアクセサメソッドを移動する
  - パッケージとして正しくビルドできることを確認する
  - _Requirements: 5.1, 5.2_

- [x] 5.2 永続化ロジックを分離する
  - `internal/app/game_state/persistence.go`を新規作成する
  - ToSaveData/FromSaveDataメソッドを移動する
  - セーブデータの後方互換性を維持する
  - 既存セーブデータの読み込みテストを実施する
  - _Requirements: 5.3, 5.5, 12.1, 12.2, 12.3_

- [x] 5.3 デフォルト値生成を統合する
  - `internal/app/game_state/defaults.go`を新規作成する
  - getDefaultCoreTypeData、getDefaultModuleDefinitionData、getDefaultPassiveSkills等を統合する
  - root_model.goとgame_state.goで重複していたデフォルトデータ生成を単一ソースに一元化する
  - _Requirements: 5.4, 7.4_

### Phase 5: battle.goの分割

- [x] 6. バトル画面のファイル分割
- [x] 6.1 UIレンダリングとゲームロジックを分離する
  - バトルUIのレンダリングロジックを別ファイルに抽出する
  - バトルロジック（ダメージ計算、状態遷移等）を分離する
  - メインのbattle.goが500行以下になることを確認する
  - バトル進行が正常に動作することを確認する
  - _Requirements: 6.1, 6.2, 6.3_

### Phase 6: コード重複解消と共通コンポーネント

- [x] 7. 共通コンポーネントの作成
- [x] 7.1 (P) HP表示コンポーネントを作成する
  - `internal/tui/components/hp_display.go`を新規作成する
  - RenderHP()関数でHPバー、数値表示、色分けロジックを共通化する
  - RenderHPWithLabel()でラベル付き表示をサポートする
  - battle.goの重複HP表示コードをこのコンポーネントに置換する
  - _Requirements: 7.1, 7.2_

- [x] 7.2 (P) モジュールアイコンメソッドを追加する
  - `domain/module.go`のModuleCategoryにIcon()メソッドを追加する
  - カテゴリ（PhysicalAttack, MagicAttack, Heal, Buff, Debuff）ごとにアイコン文字を返す
  - battle.goのgetModuleIconロジックをこのメソッドに置換する
  - _Requirements: 7.3_

### Phase 7: 型名整理

- [x] 8. 型名のリネーム
- [x] 8.1 本番用途の型名を適切な名前に変更する
  - EncyclopediaTestDataをEncyclopediaDataにリネームする
  - StatsTestDataをStatsDataにリネームする
  - 参照箇所を全て更新する
  - `*TestData`という名前の本番用型が存在しないことを確認する
  - _Requirements: 8.1, 8.2, 8.3_

### Phase 8: Screenインターフェースの導入

- [x] 9. 画面インターフェースの定義
- [x] 9.1 Screenインターフェースを作成する
  - `internal/tui/screens/types.go`を新規作成する
  - tea.Modelを埋め込んだScreenインターフェースを定義する
  - SetSize(width, height int)とGetTitle() stringメソッドを定義する
  - BaseScreen構造体で共通機能を提供する
  - _Requirements: 9.1, 9.2, 9.3_

- [x]* 9.2 既存画面のインターフェース準拠を確認する
  - 新規画面作成時にScreenインターフェースを実装するガイドラインを確立する
  - 既存画面の一部をScreenインターフェース準拠に更新する（任意）
  - _Requirements: 9.4_

### Phase 9: データ変換層の作成

- [x] 10. アダプター層の構築
- [x] 10.1 adapterディレクトリを作成する
  - `internal/adapter/`ディレクトリを新規作成する
  - パッケージ構造を初期化する
  - _Requirements: 10.1_

- [x] 10.2 (P) 永続化アダプターを実装する
  - `internal/adapter/persistence_adapter.go`を作成する
  - SaveData <-> GameState変換ロジックを集約する
  - 後方互換性を保証する変換処理を実装する
  - _Requirements: 10.2, 12.1, 12.2, 12.3_

- [x] 10.3 (P) 画面データアダプターを実装する
  - `internal/adapter/screen_adapter.go`を作成する
  - ToStatsData、ToEncyclopediaData、ToSettingsData関数を実装する
  - GameStateから各画面用データへの変換を一元管理する
  - _Requirements: 10.3_

- [x] 10.4 (P) 報酬アダプターを実装する
  - `internal/adapter/reward_adapter.go`を作成する
  - ConvertBattleStatsToRewardStats関数を実装する
  - バトル統計から報酬用統計への変換を担当する
  - _Requirements: 10.4_

- [x] 10.5 アダプター層を統合する
  - root_model.goおよび関連ファイルからアダプター層を利用するように更新する
  - 全てのデータ変換がadapter層を経由することを確認する
  - _Requirements: 10.5_

### Phase 10: 循環的複雑度の削減

- [x] 11. メッセージハンドリングの改善
- [x] 11.1 ハンドラーマップによる分岐委譲を実装する
  - RootModel.Updateメソッドのメッセージハンドリングをハンドラーマップに置換する
  - シーンレンダリングを画面マップによる委譲に変更する
  - Update内のswitch分岐数が5以下になることを確認する
  - _Requirements: 11.1, 11.2, 11.3_

### Phase 11: 最終検証

- [x] 12. リファクタリング完了の検証
- [x] 12.1 全体テストと動作確認を実施する
  - 全ユニットテストが通過することを確認する
  - 既存セーブデータが正常に読み込めることを確認する
  - 実際のゲームプレイで全シーン遷移が動作することを確認する
  - バトル進行に異常がないことを確認する
  - 外部から観測可能な動作に変更がないことを確認する
  - _Requirements: 12.1, 12.2, 12.3_

---

## Requirements Coverage

| 要件ID | 対応タスク |
|--------|-----------|
| 1.1, 1.2 | 1.1 |
| 1.3, 1.6 | 1.2 |
| 1.4, 1.5, 1.6 | 1.3 |
| 2.1, 2.2, 2.3, 2.4 | 3.1 |
| 2.5, 2.6 | 3.2 |
| 3.1, 3.2, 3.3, 3.4 | 2.1 |
| 3.5 | 2.2 |
| 4.2, 11.1 | 4.1 |
| 4.3 | 4.2 |
| 4.4, 4.5, 4.6, 4.7 | 4.3 |
| 5.1, 5.2 | 5.1 |
| 5.3, 5.5, 12.1, 12.2, 12.3 | 5.2 |
| 5.4, 7.4 | 5.3 |
| 6.1, 6.2, 6.3 | 6.1 |
| 7.1, 7.2 | 7.1 |
| 7.3 | 7.2 |
| 8.1, 8.2, 8.3 | 8.1 |
| 9.1, 9.2, 9.3 | 9.1 |
| 9.4 | 9.2 |
| 10.1 | 10.1 |
| 10.2, 12.1, 12.2, 12.3 | 10.2 |
| 10.3 | 10.3 |
| 10.4 | 10.4 |
| 10.5 | 10.5 |
| 11.1, 11.2, 11.3 | 11.1 |
| 12.1, 12.2, 12.3 | 12.1 |
