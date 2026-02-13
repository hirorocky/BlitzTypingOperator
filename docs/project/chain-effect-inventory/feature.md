# チェイン効果のスキル分離・独立インベントリ管理

## 概要
チェイン効果を現在のスキル紐付き管理から分離し、コア・スキルと並ぶ第3の独立インベントリとして管理する。
エージェントカスタマイズではコア→スキル→チェイン効果の3階層でスロット設定を行い、コアとチェイン効果に3エージェント全体のユニーク制約を追加する。

## 受け入れ基準

### チェイン効果インベントリ
1. ChainEffectInventoryがチェイン効果をTypeIDごとにユニーク管理できること
2. 既に保有しているTypeIDの追加は無視されること（falseを返す）
3. InventoryManagerがCores, Skills, ChainEffectsの3つのインベントリを統合管理すること

### コアのユニーク制約
4. 同一コアTypeIDを複数のエージェントスロットに設定しようとするとエラーが返ること
5. コアをクリアした後、他スロットでそのコアが設定可能になること

### チェイン効果のユニーク制約
6. 同一チェイン効果TypeIDを複数のスキルスロット（3エージェント全体）に設定しようとするとエラーが返ること
7. チェイン効果をクリアした後、他スロットでそのチェイン効果が設定可能になること

### 階層制約
8. スキルが設定されていないスキルスロットにチェイン効果を設定しようとするとエラーが返ること
9. スキルをクリアした際、同スキルスロットのチェイン効果が自動的にクリアされること
10. コアをクリアした際、スキルが自動削除され、さらにチェイン効果も連鎖的にクリアされること

### エージェントカスタマイズUI
11. カスタマイズ画面でコア、スキル、チェイン効果の3つを個別に選択・変更できること
12. スキル未設定のスキルスロットではチェイン効果の選択が不可であること
13. 他スロットで使用中のコア・チェイン効果が選択リストで区別表示されること

### セーブ/ロード
14. チェイン効果インベントリがセーブデータに正しく保存・復元されること
15. エージェントスロットのチェイン効果設定がセーブデータに正しく保存・復元されること

### 報酬システム
16. モジュールドロップ時、スキルがSkillInventoryに、チェイン効果がChainEffectInventoryにそれぞれ独立して追加されること

### 初期エージェント
17. 3つの初期エージェントがそれぞれ異なるコアを持つこと
18. 初期コアが初期スキルと互換性を持つこと

### インベントリ画面（スキル一覧タブからのチェイン効果UI除去）
19. スキル一覧タブの各項目にチェイン効果バリエーション数（`(N種)`表示）が表示されないこと
20. スキル詳細プレビューにチェイン効果バリエーションセクション（「チェイン効果バリエーション:」ラベルおよび個別チェインID表示）が表示されないこと
21. `SkillInventoryItem`構造体に`ChainCount`, `ChainVariations`フィールドが存在しないこと
22. `updateSkillList()`がチェイン効果関連のデータ取得（`GetChainVariations()`呼び出し等）を行わないこと
23. スキル詳細のEnterキー詳細表示切り替え（`showingSkillDetail`）がチェイン効果情報を含まないこと

### インベントリ画面（チェイン効果一覧タブの追加）
24. インベントリ画面にチェイン効果一覧タブが存在し、コア一覧・スキル一覧と独立して表示されること
25. チェイン効果一覧で各チェイン効果の名前・説明・装備状態が確認できること

### バトル構築
26. BuildAgentsForBattleが各スキルスロットに対応するチェイン効果を正しく解決してAgentModelを構築すること

### debugモード対応
27. debugモード起動時に全チェイン効果がChainEffectInventoryに追加されること
28. debugモードのエージェントカスタマイズで全チェイン効果が選択肢に表示されること
29. debugモードでもチェイン効果のユニーク制約（同一TypeIDの複数スロット設定禁止）が機能すること

## 設計

### 変更対象

#### domain層
- `domain/skill_inventory.go`:
  - `SkillOwnership`から`ChainVariations`フィールドを除去
  - `AddSkill`のシグネチャを`AddSkill(typeID string)`に変更（chainEffectID引数を除去）
  - `GetChainVariations`, `HasChainVariation`メソッドを除去

- `domain/agent_slot.go`:
  - `SkillSlotConfig`から`ChainEffectID`フィールドを除去
  - `AgentSlot`に`ChainEffects [MaxSkillSlotCount]ChainEffectSlotConfig`フィールドを追加
  - `SetChainEffect`, `ClearChainEffect`, `GetChainEffect`メソッドを追加
  - `ClearSkill`時に同インデックスのチェイン効果も自動クリア

#### usecase層
- `usecase/slot/agent_slot_manager.go`:
  - `SetCore`: コアの3エージェント全体ユニーク制約を追加
  - `SetSkill`: シグネチャから`chainEffectID`を除去
  - `SetChainEffect(slot, skillSlot int, chainEffectTypeID string) error`: 新規メソッド
  - `ClearChainEffect(slot, skillSlot int) error`: 新規メソッド
  - `ClearSkill`: 同スキルスロットのチェイン効果も自動クリア
  - `chainEffectInv *domain.ChainEffectInventory`フィールドを追加
  - エラー定義追加: `ErrCoreAlreadyEquipped`, `ErrChainEffectAlreadyEquipped`, `ErrChainEffectNotOwned`, `ErrSkillNotSetForChain`
  - `buildModuleFromConfig`のチェイン効果解決を`ChainEffects`配列から取得するように変更

- `usecase/inventory/inventory_manager.go`:
  - `chainEffects *domain.ChainEffectInventory`フィールドを追加
  - `AddChainEffect(typeID string) bool`メソッドを追加
  - `ChainEffects() *domain.ChainEffectInventory`アクセサを追加
  - コンストラクタにChainEffectInventory引数を追加

#### infra層
- `infra/savedata/savedata.go`:
  - `ChainEffectInventorySave`構造体を追加（`ChainEffects []string`）
  - `InventorySaveData`に`UniqueChainEffects *ChainEffectInventorySave`を追加
  - `SkillInventorySave`の形式変更: `Skills map[string][]string` → `Skills []string`（チェイン効果分離後はTypeIDリストのみ）
  - `SkillSlotSaveCfg`から`ChainEffectID`を除去
  - `AgentSlotSave`に`ChainEffects [4]ChainEffectSlotSaveCfg`を追加
  - `ChainEffectSlotSaveCfg`構造体を追加

- `infra/savedata/unique_inventory_converter.go`:
  - `ConvertChainEffectInventoryToSave` / `ConvertSaveToChainEffectInventory`関数を追加
  - `ConvertSkillInventoryToSave` / `ConvertSaveToSkillInventory`をチェイン効果分離に対応
  - `ConvertAgentSlotToSave` / `ConvertSaveToAgentSlot`をチェイン効果スロット対応

#### usecase/rewarding層
- `usecase/rewarding/reward.go`:
  - `AddRewardsToInventory`にChainEffectInventory引数を追加
  - スキルとチェイン効果を別々のインベントリに追加するように変更

#### tui層
- `tui/screens/inventory.go`:
  - `InventoryTab`に`TabChainEffectInventory`を追加（3タブ構成）
  - `ChainEffectInventoryItem`構造体を追加（Name, Description, TypeID, Equipped bool）
  - チェイン効果一覧の表示処理を追加
  - `SkillInventoryItem`から`ChainCount`, `ChainVariations`フィールドを除去
  - `renderSkillListItems()`: スキル名の後の`(N種)`チェイン効果数表示を除去
  - `renderSkillPreviewContent()`: 「チェイン効果バリエーション:」セクション全体を除去（ラベル、個別ID表示、「(なし)」表示、詳細表示切り替え含む）
  - `updateSkillList()`: `GetChainVariations()`呼び出しと関連データ設定を除去
  - `showingSkillDetail`フラグ: チェイン効果以外の用途がなければ除去を検討

- `tui/screens/agent_customization.go`:
  - チェイン効果選択モードの追加（`ModeChainSelect`を独立フローに変更）
  - `focusPosition`のカード内表示でチェイン効果スロットを表示
  - スキル未設定時のチェイン効果選択無効化
  - コア選択リストで他スロット使用中を表示
  - チェイン効果選択リストで他スロット使用中を表示

#### app層
- `app/screen_factory.go` 等:
  - AgentSlotManager/InventoryManagerのコンストラクタ変更に追従
  - ChainEffectInventoryの初期化・受け渡し

- `app/root_model.go`:
  - debugモード初期化処理を新インベントリ構造に対応
  - 全チェイン効果を`invManager.AddChainEffect(ce.ID)`で個別追加（旧: スキル経由の追加を廃止）
  - 旧方式の`invManager.AddSkill(mt.ID, ce.ID)`ループを除去し、`invManager.AddSkill(typeID)`とチェイン効果追加を分離

#### tui/presenter層
- `tui/presenter/debug_inventory_provider.go`:
  - `GetChainEffects()`が新ChainEffectInventoryからデータを返すよう変更
  - チェイン効果選択リスト生成時にChainEffectInventory経由で保有チェイン効果を提供

#### マスタデータ
- `infra/masterdata/data/first_agent.json`:
  - 3エージェントに異なるコアTypeIDを設定
  - 候補: `all_rounder`(physical_strike_lv1互換), `healer`(heal_lv1互換), `quick_recoverer`(str_buff_lv1互換)

### 新規作成
- `domain/chain_effect_inventory.go`: ChainEffectInventory構造体
  - `map[string]struct{}`でTypeIDの保有セットを管理
  - `AddChainEffect(typeID string) bool`
  - `HasChainEffect(typeID string) bool`
  - `GetOwnedChainEffects() []string`（ソート済み）
- `domain/chain_effect_inventory_test.go`: ユニットテスト

### データフロー

#### チェイン効果の取得フロー
```
敵撃破 → RewardCalculator
  → スキルドロップ → SkillInventory.AddSkill(typeID)
  → チェイン効果ドロップ → ChainEffectInventory.AddChainEffect(typeID)
```

#### エージェントカスタマイズフロー
```
カード選択 → コア選択（ユニーク制約あり）
  → スキルスロット選択 → スキル選択（ユニーク制約あり）
  → チェイン効果選択（ユニーク制約あり、スキル必須）
```

#### バトルビルドフロー
```
AgentSlotManager.BuildAgentsForBattle()
  → 各スロットのCoreTypeID, Skills[i], ChainEffects[i]を取得
  → ChainEffects[i].TypeIDでマスタデータからChainEffectを解決
  → SkillModelにChainEffectを設定してAgentModelを構築
```

## メモ
- ChainEffectInventoryはCoreInventoryと同じ`map[string]struct{}`パターンを採用（GoのSet表現の定石）
- セーブデータ後方互換性: 旧形式のSkillInventorySave（TypeID→ChainEffectIDリスト）読み込み時はスキルのみ復元し、チェイン効果はChainEffectInventoryに移行する変換処理が必要
- チェイン効果の値（Value）はマスタデータから解決される（インベントリにはTypeIDのみ保存）
- SkillInventoryのSimplification: ChainVariations除去後はSkillOwnership構造体自体が不要になる可能性があるが、将来の拡張性を考慮して構造は維持してもよい（または単純なmap[string]struct{}に変更）
- コアのユニーク制約追加により、CoreInventoryには最低3種類のコアが必要。初期に3種類提供することでゲーム開始時から制約を満たせる
