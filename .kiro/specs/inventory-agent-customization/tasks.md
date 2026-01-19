# 実装計画

## タスク一覧

- [x] 1. ユビキタス言語の変更（Module → Skill）
- [x] 1.1 (P) ドメインモデル名の変更
  - ModuleModelをSkillModelに名称変更、ModuleTypeをSkillTypeに変更
  - ModuleEffectをSkillEffectに変更
  - ModuleInventoryをSkillInventoryLegacyに名称変更
  - 後方互換性のためModule*はSkill*のエイリアスとして残存
  - _Requirements: 2.1_

- [x] 1.2 (P) セーブデータスキーマの変更
  - savedata.AgentInstanceSaveのModulesフィールドをSkillsに変更
  - ModuleInstanceSaveをSkillInstanceSaveに変更
  - ModuleEffectDataをSkillEffectDataに変更
  - JSONキーは後方互換性のため"modules"のまま
  - _Requirements: 2.1_

- [x] 2. コアインベントリのユニーク管理実装
- [x] 2.1 CoreInventoryドメインモデルの実装
  - CoreTypeIDをキーとし、取得済み最大レベルを値とするマップ構造で管理
  - AddCoreメソッド: 新規追加時はそのまま保存、既存TypeIDの場合はレベル比較して高い方を保持
  - GetMaxLevelメソッド: 指定TypeIDの最大レベル取得（未保有時は0を返却）
  - GetOwnedCoresメソッド: 全保有コア一覧を返却
  - HasCoreメソッド: 指定TypeIDの保有判定
  - レベルは1以上の正整数であることを検証
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [x] 2.2 CoreInventoryの単体テスト
  - 新規コア追加のテスト
  - より高いレベルのコア取得で最大レベルが更新されるテスト
  - 既存最大レベル以下のコア取得で変更されないテスト
  - 未保有TypeIDへのGetMaxLevelで0が返るテスト
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5_

- [x] 3. スキルインベントリのユニーク管理実装
- [x] 3.1 (P) SkillOwnership構造体の実装
  - 保有フラグとチェイン効果バリエーションセットを保持
  - チェイン効果IDはmap[string]boolでセット管理
  - _Requirements: 2.1, 2.3_

- [x] 3.2 (P) SkillInventoryドメインモデルの実装
  - SkillTypeIDをキーとし、SkillOwnershipを値とするマップ構造で管理
  - AddSkillメソッド: 初回取得で保有状態をtrue、チェイン効果があれば追加
  - GetOwnedSkillsメソッド: 全保有スキル情報を返却
  - GetChainVariationsメソッド: 指定TypeIDで利用可能なチェイン効果IDリストを返却
  - HasSkillメソッド: 指定TypeIDの保有判定
  - HasChainVariationメソッド: 指定チェイン効果バリエーションの保有判定
  - チェイン効果なしスキルは空文字列で表現
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

- [x] 3.3 SkillInventoryの単体テスト
  - 初回スキル取得で保有状態がtrueになるテスト
  - チェイン効果付きスキル取得でバリエーションが追加されるテスト
  - 同一スキル再取得で新しいチェイン効果が追加されるテスト
  - チェイン効果なしスキルの取得テスト
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

- [x] 4. エージェントスロットシステムの実装
- [x] 4.1 (P) SkillSlotConfig値オブジェクトの実装
  - スキルTypeIDとチェイン効果IDを保持
  - 空スロット判定メソッドの実装
  - _Requirements: 5.1, 5.2_

- [x] 4.2 (P) AgentSlotドメインモデルの実装
  - コアTypeIDと選択レベルを保持
  - 4つのSkillSlotConfigスロットを固定配列で保持
  - IsEmptyメソッド: コア未設定の場合trueを返却
  - GetSkillCountメソッド: 設定済みスキル数を返却
  - Clearメソッド: スロット全体をクリア
  - コア未設定時はスキルも全て空であるべき不変条件を維持
  - _Requirements: 3.1, 3.2, 3.3, 4.1, 5.1_

- [x] 4.3 AgentSlotの単体テスト
  - 空スロットの判定テスト
  - コア設定済みスロットの判定テスト
  - スキル数カウントのテスト
  - クリア操作のテスト
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 5. エージェントスロット管理ユースケースの実装
- [x] 5.1 AgentSlotManagerの基本構造実装
  - 3つの固定スロット配列を管理
  - CoreInventory、SkillInventoryへの参照を保持
  - マスタデータ（CoreType、SkillType、PassiveSkill）への参照を保持
  - NewAgentSlotManagerファクトリ関数の実装
  - GetSlots、GetSlotメソッドの実装
  - _Requirements: 3.1, 3.2_

- [x] 5.2 コア付け替え機能の実装
  - SetCoreメソッド: 指定スロットにコアを設定
  - インベントリ保有確認（未保有エラー）
  - レベル範囲確認（1〜最大レベル、範囲外エラー）
  - 設定時にエージェントのステータスを再計算
  - コア変更時にスキル互換性チェックを実行し、互換性のないスキルを自動削除
  - ClearCoreメソッド: コアをクリア（スキルも全削除）
  - 同一CoreTypeIDを複数スロットに設定可能
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 7.1, 7.2, 7.3_

- [x] 5.3 スキル付け替え機能の実装
  - SetSkillメソッド: 指定スロットのスキルスロットにスキルを設定
  - インベントリ保有確認（未保有エラー）
  - コア未設定エラーのチェック
  - タグマッチングによる互換性検証（コアの許可タグを持たない場合エラー）
  - チェイン効果バリエーションの保有確認（指定した場合）
  - ClearSkillメソッド: 指定スキルスロットをクリア
  - 同一スキルTypeIDを異なるエージェント間で装備可能
  - 同一スキルTypeIDを同一エージェント内の複数スロットに装備可能
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 6.1, 6.2, 6.3, 6.4_

- [x] 5.4 バトル連携機能の実装
  - IsSlotReadyメソッド: コア設定済みでバトル使用可能かを判定
  - GetReadySlotCountメソッド: 使用可能スロット数を返却
  - BuildAgentsForBattleメソッド: スロット構成からAgentModelスライスを構築
  - 空スロットはバトルに含めない
  - ValidateSkillCompatibilityメソッド: スキルとコアの互換性検証
  - GetCompatibleSkillsメソッド: 指定スロットのコアと互換性のあるスキル一覧を返却
  - _Requirements: 9.1, 9.2, 9.3_

- [x] 5.5 AgentSlotManagerの単体テスト
  - コア設定・クリアのテスト
  - レベル範囲外エラーのテスト
  - 未保有コア設定エラーのテスト
  - スキル設定・クリアのテスト
  - 互換性なしスキル設定エラーのテスト
  - コア変更時の互換性なしスキル自動削除テスト
  - BuildAgentsForBattleの生成結果テスト
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 5.7, 9.1, 9.2, 9.3_

- [x] 6. インベントリ統合管理ユースケースの実装
- [x] 6.1 (P) InventoryManager改の実装
  - CoreInventoryとSkillInventoryを統合管理
  - AddCoreメソッド: コア追加のラッパー
  - AddSkillメソッド: スキル追加のラッパー
  - Cores、Skillsメソッド: 各インベントリへのアクセサ
  - GetOwnedCoreTypes、GetOwnedSkillTypesメソッド: 保有TypeID一覧を返却
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

- [x] 6.2 (P) InventoryManager改の単体テスト
  - コア追加・取得のテスト
  - スキル追加・取得のテスト
  - _Requirements: 1.1, 1.5, 2.1, 2.6_

- [ ] 7. セーブデータスキーマの更新
- [ ] 7.1 新セーブデータ構造体の実装
  - CoreInventorySave: TypeID→最大レベルのマップ形式
  - SkillInventorySave: TypeID→チェイン効果IDリストのマップ形式
  - SkillSlotSaveCfg: スキルTypeIDとチェイン効果IDを保持
  - AgentSlotSave: コアTypeID、レベル、4つのスキルスロット構成を保持
  - InventorySaveDataにユニークコアとスキルを格納
  - PlayerSaveDataに3スロットのAgentSlotSaveを格納
  - _Requirements: 3.4_

- [ ] 7.2 セーブ/ロード機能の実装
  - 新スキーマでの永続化処理
  - ロード時のマスタデータ整合性検証（存在しないTypeIDは無視）
  - セーブデータバージョンを3.0.0に更新
  - _Requirements: 3.4_

- [ ] 7.3 セーブ/ロードの統合テスト
  - 新スキーマでの保存・復元テスト
  - マスタデータに存在しないTypeIDの無視テスト
  - _Requirements: 3.4_

- [ ] 8. バトルシステムとの統合
- [ ] 8.1 BattleEngineのエージェント参照変更
  - 装備エージェントリストの代わりにAgentSlotManagerからエージェント構成を参照
  - バトル開始時にBuildAgentsForBattleでAgentModelを構築
  - コア設定済みスロットのみをバトルに含める
  - _Requirements: 9.1, 9.2, 9.3_

- [ ] 8.2 バトル中のスロット変更禁止
  - バトル中はスロット構成の変更を許可しない制御を実装
  - _Requirements: 9.4_

- [ ] 8.3 バトル統合の統合テスト
  - スロット構成からのAgentModel構築テスト
  - 空スロット除外テスト
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [ ] 9. インベントリUI画面の実装
- [ ] 9.1 InventoryScreenの基本構造実装
  - InventoryManagerへの参照を保持
  - AgentSlotManagerへの参照を保持（装備状況表示用）
  - タブ切り替え（コア一覧/スキル一覧）の実装
  - _Requirements: 10.1, 10.2_

- [ ] 9.2 コア一覧表示の実装
  - 保有CoreTypeIDと取得済み最大レベルを一覧表示
  - 現在エージェントスロットに装備中のコアを示すマーカー表示
  - _Requirements: 10.1, 10.4_

- [ ] 9.3 スキル一覧表示の実装
  - 保有SkillTypeIDと利用可能なチェイン効果バリエーション数を一覧表示
  - スキル選択時にチェイン効果バリエーション一覧を詳細表示
  - 現在エージェントスロットに装備中のスキルを示すマーカー表示
  - _Requirements: 10.2, 10.3, 10.4_

- [ ] 10. エージェントカスタマイズUI画面の実装
- [ ] 10.1 AgentCustomizationScreenの基本構造実装
  - AgentSlotManagerへの参照を保持
  - 3つのエージェントスロットと現在の構成を表示
  - スロット選択モードの実装
  - _Requirements: 11.1, 11.2_

- [ ] 10.2 コア選択モードの実装
  - 保有CoreTypeID一覧を表示
  - レベル選択オプションを表示（1〜取得済み最大レベル）
  - コア選択・設定の操作フロー
  - _Requirements: 11.3_

- [ ] 10.3 スキル選択モードの実装
  - 現在のコアと互換性のあるスキルをフィルタリングして表示
  - チェイン効果バリエーションオプションを表示
  - スキル選択・設定の操作フロー
  - スキルとコアの互換性を視覚的に表示（互換/非互換マーカー）
  - _Requirements: 11.4, 11.5_

- [ ] 10.4 スロットクリア機能の実装
  - コアまたはスキルスロットをクリアして空にする操作
  - _Requirements: 11.6_

- [ ] 11. 既存システムの削除
- [ ] 11.1 合成システムの削除
  - synthesize.AgentManagerの削除
  - Synthesize、Equip、Unequipメソッドの削除
  - 合成によるエージェントインスタンス作成機能の削除
  - _Requirements: 8.1_

- [ ] 11.2 旧インベントリモデルの削除
  - 旧CoreInventory（インスタンス管理方式）の削除
  - 旧SkillInventory（スライス管理方式）の削除
  - AgentInventory（合成エージェント保存用）の削除
  - _Requirements: 8.2_

- [ ] 11.3 旧セーブデータ構造体の削除
  - CoreInstanceSaveの削除
  - ModuleInstanceSaveの削除
  - AgentInstanceSaveの削除
  - _Requirements: 8.3_

- [ ] 11.4 旧UI画面の削除
  - AgentManagementScreenの合成タブを削除
  - 合成関連のUIコンポーネントを削除
  - _Requirements: 8.1, 8.2_

- [ ] 12. 統合テストとシステム検証
- [ ] 12.1 インベントリ→スロット連携テスト
  - コア追加後のスロット設定フローテスト
  - スキル追加後のスロット設定フローテスト
  - _Requirements: 1.1, 1.2, 1.3, 2.1, 2.2, 2.3, 4.1, 5.2_

- [ ] 12.2 スロット→バトル連携テスト
  - BuildAgentsForBattleの生成AgentModel検証
  - スロット構成変更後のバトル反映テスト
  - _Requirements: 9.1, 9.2, 9.3_

- [ ] 12.3 E2Iカスタマイズフローテスト
  - スロット選択→コア設定→レベル選択→スキル設定→チェイン効果選択の一連フローテスト
  - 互換性フィルタリングの動作確認
  - _Requirements: 11.1, 11.2, 11.3, 11.4, 11.5, 11.6_
