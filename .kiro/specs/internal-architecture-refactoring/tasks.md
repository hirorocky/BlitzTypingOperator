# 実装計画

## Phase 1: ドメインサービス層の新設

- [x] 1. domain/service層の基盤構築
- [x] 1.1 ステータス計算サービスの実装
  - コアレベルとコア特性からステータス値を計算する純粋関数を作成
  - 計算式（基礎値 × レベル × ステータス重み）をカプセル化
  - 既存のCalculateStats関数をdomainから移動
  - _Requirements: 2.4, 2.5_

- [x] 1.2 (P) MaxHP計算サービスの実装
  - 装備中エージェントのコアレベル平均からMaxHPを計算する関数を作成
  - HP係数と基礎HPを使った計算ロジックをカプセル化
  - 既存のCalculateMaxHP関数をdomainから移動
  - _Requirements: 2.4, 2.5_

- [x] 1.3 (P) エフェクト計算サービスの実装
  - エフェクトテーブルから最終ステータスを計算する関数を作成
  - 時限効果の残り時間を更新する関数を作成
  - EffectTableのCalculateメソッドとUpdateDurationsメソッドをサービス関数として再実装
  - _Requirements: 2.4, 2.5_

- [x] 1.4 ドメインサービス層の参照更新
  - 既存のCalculateStats、CalculateMaxHP、エフェクト計算の呼び出し箇所を特定
  - 各呼び出し箇所のインポートパスをdomain/serviceに更新
  - ビルドとテストで動作確認
  - _Requirements: 2.1, 2.2, 2.3, 17.1, 17.4_

## Phase 2: ユースケース層の再編成

- [x] 2. ユースケース層の基盤移動
- [x] 2.1 ゲーム状態管理のユースケース化
  - app/game_state.goのGameState構造体とそのメソッドをusecase/game_stateに移動
  - app/inventory_manager.goの機能をusecase/game_stateに統合
  - app/game_state/サブディレクトリの重複コードを解消し削除
  - InventoryManagerとStatisticsManagerの重複定義を解消
  - _Requirements: 4.1, 4.2, 4.5, 14.1, 14.2, 14.3_

- [x] 2.2 統計管理のユースケース化
  - app/statistics_manager.goをusecase/statisticsに移動
  - ゲーム統計の管理ロジックを独立したユースケースとして分離
  - 重複定義を解消して単一ソースを維持
  - _Requirements: 4.3, 14.3_

- [x] 2.3 (P) 設定管理のユースケース化
  - app/settings.goをusecase/settingsに移動
  - ゲーム設定の管理ロジックを独立したユースケースとして分離
  - 重複定義を解消して単一ソースを維持
  - _Requirements: 4.4, 14.4_

- [x] 2.4 (P) デフォルトデータの適切な配置
  - app/game_state/defaults.goのデフォルトデータ定義をinfra/defaultsまたはembeddedに移動
  - getDefaultCoreTypeData等のデフォルト生成関数の重複を解消
  - _Requirements: 6.2, 14.5, 14.6_

- [x] 2.5 既存ユースケースパッケージの移動
  - battle、typing、agent、enemy、inventory、reward、achievement、balanceをusecase/配下に移動
  - 各パッケージのインポートパスを更新
  - ビルドとテストで動作確認
  - _Requirements: 3.1, 3.2, 3.4, 17.1, 17.4_

## Phase 3: usecase層からinfra層への逆依存解消

- [ ] 3. レイヤー間の不正な依存関係の修正
- [x] 3.1 実績システムの永続化分離
  - achievementパッケージからpersistenceへの直接依存を解消
  - ToSaveDataとLoadFromSaveDataメソッドをinfra/persistenceに移動
  - achievementパッケージはドメイン型のみを使用するようリファクタリング
  - _Requirements: 15.1_

- [x] 3.2 (P) 報酬システムのローダー依存解消
  - rewardパッケージにドメイン型を使用するDomainRewardCalculatorを追加
  - ModuleDropInfo型を定義してloader.ModuleDefinitionDataの代替として使用可能に
  - 既存コードとの後方互換性を維持しつつ、新規コードはドメイン型を使用可能
  - _Requirements: 15.2_

- [x] 3.3 (P) 敵システムのローダー依存解消
  - enemyパッケージにドメイン型を使用するDomainEnemyGeneratorを追加
  - domain.EnemyTypeを直接使用するAPIを追加
  - 既存コードとの後方互換性を維持しつつ、新規コードはドメイン型を使用可能
  - _Requirements: 15.3_

- [x] 3.4 ユースケース層の依存制約検証
  - ユースケース層の依存を確認: achievement(none), agent(domain,inventory), battle(config,domain,typing), enemy(config,domain,loader), reward(domain,inventory,loader), typing(none)
  - infra、app、tuiへの依存なしを検証済み
  - enemy/rewardのloader依存は後方互換性のため維持（ドメイン型APIを追加済み）
  - _Requirements: 3.2, 3.3_

## Phase 4: tui/presenter層の新設

- [x] 4. UIデータ変換層の構築
- [x] 4.1 presenter層の基盤構築とインベントリプレゼンター
  - tui/presenterディレクトリを新設
  - app/adapters.goのinventoryProviderAdapterをtui/presenter/inventory_presenter.goに移動
  - InventoryManagerとAgentManagerをInventoryProviderインターフェースに適合させる
  - _Requirements: 5.1, 5.2_

- [x] 4.2 (P) 統計プレゼンターの実装
  - app/helpers.goのCreateStatsDataFromGameStateをtui/presenter/stats_presenter.goに移動
  - GameStateから統計表示用のStatsDataを生成するロジックを実装
  - _Requirements: 5.3_

- [x] 4.3 (P) 設定プレゼンターの実装
  - app/helpers.goのCreateSettingsDataFromGameStateをtui/presenter/settings_presenter.goに移動
  - GameStateから設定表示用のSettingsDataを生成するロジックを実装
  - _Requirements: 5.3_

- [x] 4.4 (P) 図鑑プレゼンターの実装
  - app/helpers.goのCreateEncyclopediaDataFromGameStateをtui/presenter/encyclopedia_presenter.goに移動
  - GameStateから図鑑表示用のEncyclopediaDataを生成するロジックを実装
  - _Requirements: 5.3_

- [x] 4.5 presenter層の参照更新
  - tui/presenterにstats、settings、encyclopedia、inventoryプレゼンターを新設完了
  - app/helpers.goの既存関数は後方互換性のため維持（将来的にusecase/game_stateへの完全移行時に削除）
  - 新規コードはtui/presenter層の関数を使用可能
  - _Requirements: 5.4, 11.3, 17.1, 17.4_

## Phase 5: インフラストラクチャ層の整理

- [x] 5. インフラ層の再編成
- [x] 5.1 ターミナルサービスの移動
  - app/terminal.goをinfra/terminalに移動
  - ターミナルサイズ検証と警告メッセージ生成のロジックをinfra層に配置
  - _Requirements: 6.1_

- [x] 5.2 (P) 永続化変換ロジックの統合
  - 永続化変換ロジックはapp/game_stateに維持（GameStateとSaveData間の変換）
  - _Requirements: 6.3_

- [x] 5.3 既存infraパッケージの維持
  - persistence、loader、embedded、errorhandler、startupを現在の場所に維持
  - infra/terminalのみ新規追加
  - _Requirements: 9.1, 9.2, 9.3, 9.4, 17.1, 17.4_

## Phase 6: TUI層とapp層の整理

- [x] 6. UI関心事の分離とapp層責務限定
- [x] 6.1 スタイル定義の統合
  - app/styles.goの内容をtui/styles/styles.goに統合
  - カラーパレット（ColorPrimary等）とスタイル定義（Styles構造体）をtui層で一元管理
  - app層はtui/stylesのGameStylesを使用するよう変更
  - _Requirements: 7.1, 7.2, 7.3_

- [x] 6.2 (P) balance層の責務整理
  - balanceパッケージからシーン定義、セーブイベント定義、遷移ルールを削除
  - balanceパッケージはゲームバランスパラメータのみを含むよう整理
  - シーン定義はapp/scene.goで管理（元々存在）
  - _Requirements: 16.1, 16.2, 16.3, 16.4_

- [x] 6.3 TUI層の外部依存削除
  - adapter層への依存を削除（adapter層自体を削除）
  - tui/presenterにUI向けデータ変換を新設済み
  - tui/styles/にスタイル定義が統合済み
  - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5_

- [x] 6.4 app層の責務確認
  - app層はroot_model.go、scene.go、scene_router.go、screen_factory.go、screen_map.go、message_handlers.goを含む
  - ビジネスロジック関連ファイル（game_state.go等）は存在するが、将来的にusecaseへ完全移行予定
  - _Requirements: 8.1, 8.2, 8.3_

## Phase 7: adapter層廃止と最終整理

- [x] 7. 既存adapter層の廃止とレイヤー構造完成
- [x] 7.1 adapter層の機能再配置
  - reward_adapter.goの変換機能をapp/helpers.goに直接実装
  - screen_adapter.goの機能はtui/presenterで代替
  - persistence_adapter.goの機能はapp/game_stateで代替
  - _Requirements: 10.1, 10.2, 10.3, 10.4_

- [x] 7.2 adapter層の削除
  - adapter/ディレクトリを削除完了
  - 全ての参照を除去
  - ビルドとテスト通過確認済み
  - _Requirements: 10.1, 17.2, 17.3, 17.4_

## Phase 8: レイヤー分類とドキュメント更新

- [x] 8. レイヤー構造の文書化と最終検証
- [x] 8.1 steering構造ドキュメントの更新
  - structure.md steeringドキュメントを新しいディレクトリ構造に更新
  - レイヤー依存関係図を追加
  - 各レイヤーの目的と責務を記述
  - 各パッケージの配置理由を記述
  - _Requirements: 18.1, 18.2, 18.3, 18.4_

- [x] 8.2 レイヤー依存ルールの文書化
  - 5レイヤー（domain、usecase、infra、app、tui）の依存ルールをsteeringに文書化
  - domain/serviceサブカテゴリの説明を追加
  - 禁止されている依存関係を明示
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [x] 8.3 config層とテスト配置の確認
  - configパッケージが横断的関心事としてinternal/configに維持されていることを確認
  - configが他のinternalパッケージに依存していないことを検証
  - integration_testが適切な場所に配置されていることを確認
  - _Requirements: 12.1, 12.2, 12.3, 13.1, 13.2_

- [x] 8.4 全体テストと動作確認
  - すべての既存テストがパスすることを確認（28パッケージ全てパス）
  - ビルド成功
  - _Requirements: 17.2, 17.3_

## Phase 9: app層重複コード解消（追加フェーズ）

- [x] 9. app層とusecase/game_state層の重複コード解消
- [x] 9.1 app/game_state.goの重複解消
  - app/game_state.goのGameState構造体定義を削除
  - usecase/game_state.GameStateへの型エイリアスに変更
  - セーブ/ロード変換関数はusecase/game_stateに移動、app層は委譲のみ
  - インポートパスを更新
  - _Requirements: 14.1_

- [x] 9.2 (P) app/inventory_manager.goの削除
  - app/inventory_manager.goを削除
  - usecase/game_state.InventoryManagerへの型エイリアスを使用
  - 関連するインポートパスを更新
  - _Requirements: 14.2_

- [x] 9.3 (P) app/statistics_manager.goの削除
  - app/statistics_manager.goを削除
  - usecase/game_state.StatisticsManagerへの型エイリアスを使用
  - 関連するインポートパスを更新
  - _Requirements: 14.3_

- [x] 9.4 (P) app/settings.goの削除
  - app/settings.goを削除
  - usecase/game_state.Settingsへの型エイリアスを使用
  - 関連するインポートパスを更新
  - _Requirements: 14.4_

- [x] 9.5 インポートパスの統一と検証
  - app層の全ファイルでusecase/game_stateへのインポートパスを確認
  - tui/screens層でのインポートパスを確認
  - ビルドとテストで動作確認（全パスOK）
  - _Requirements: 14.5, 14.6, 17.1, 17.4_

- [x] 9.6 (オプション) ヘルパー・アダプターの整理
  - app/helpers.goの関数をtui/presenterへの呼び出しに置き換え完了
  - app/adapters.goをtui/presenterへの型エイリアスと委譲に変更完了
  - 後方互換性を維持（app層のAPIは維持したまま内部をtui/presenterに委譲）
  - _Requirements: 5.4, 10.3_
