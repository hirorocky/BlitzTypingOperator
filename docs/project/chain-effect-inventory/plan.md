# 実装計画: チェイン効果のスキル分離・独立インベントリ管理

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル | 状態 |
|---|------|-----------|---------------|------|
| 1 | ChainEffectInventoryのユニーク管理 | Go test | `internal/domain/chain_effect_inventory_test.go` | 作成済 |
| 2 | 既保有TypeIDの追加は無視 | Go test | `internal/domain/chain_effect_inventory_test.go` | 作成済 |
| 3 | InventoryManagerの3インベントリ統合 | Go test | `internal/usecase/inventory/chain_effect_integration_test.go` | 作成済 |
| 4 | コアのユニーク制約（重複エラー） | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 5 | コアクリア後の再設定 | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 6 | チェイン効果のユニーク制約 | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 7 | チェイン効果クリア後の再設定 | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 8 | スキル未設定でチェイン効果設定不可 | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 9 | スキルクリアでチェイン効果自動クリア | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 10 | コアクリアで連鎖クリア | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 11 | カスタマイズ画面で3種個別選択 | TUI test | セッション内 | 手順記載済 |
| 12 | スキル未設定でチェイン効果選択不可 | TUI test | セッション内 | 手順記載済 |
| 13 | 使用中コア・チェイン効果の区別表示 | TUI test | セッション内 | 手順記載済 |
| 14 | チェイン効果インベントリのセーブ/ロード | Go test | `internal/infra/savedata/chain_effect_save_test.go` | 作成済 |
| 15 | エージェントスロットのチェイン効果セーブ/ロード | Go test | `internal/infra/savedata/chain_effect_save_test.go` | 作成済 |
| 16 | 報酬でのスキル/チェイン効果分離追加 | Go test | `internal/usecase/rewarding/chain_effect_reward_test.go` | 作成済 |
| 17 | 初期エージェントが異なるコアを持つ | Go test | `internal/infra/startup/chain_effect_startup_test.go` | 作成済 |
| 18 | 初期コアと初期スキルの互換性 | Go test | `internal/infra/startup/chain_effect_startup_test.go` | 作成済 |
| 19 | スキル一覧のチェイン効果数表示なし | TUI test | セッション内 | 手順記載済 |
| 20 | スキル詳細のチェイン効果セクションなし | TUI test | セッション内 | 手順記載済 |
| 21 | SkillInventoryItemのフィールド除去 | コンパイル検証 | タスク11で`ChainCount`/`ChainVariations`除去後にビルド成功で確認 | タスク内で確認 |
| 22 | updateSkillListのChainVariations除去 | コンパイル検証 | タスク3/11で`GetChainVariations`除去後にビルド成功で確認 | タスク内で確認 |
| 23 | スキル詳細切替のチェイン効果除去 | TUI test | セッション内（基準20と同時検証） | 手順記載済 |
| 24 | チェイン効果一覧タブの存在 | TUI test | セッション内 | 手順記載済 |
| 25 | チェイン効果一覧の名前・説明・装備状態 | TUI test | セッション内 | 手順記載済 |
| 26 | BuildAgentsForBattleのチェイン効果解決 | Go test | `internal/usecase/slot/chain_effect_slot_test.go` | 作成済 |
| 27 | debugモードで全チェイン効果がインベントリに追加 | TUI test | セッション内（debugモード起動 → インベントリ確認） | 手順記載済 |
| 28 | debugモードで全チェイン効果が選択肢に表示 | TUI test | セッション内（debugモード起動 → カスタマイズ確認） | 手順記載済 |
| 29 | debugモードでもユニーク制約が機能 | Go test | 基準4-10テストで暗黙カバー（制約はインベントリ由来に依存しない） | カバー済 |

## 実装順序と依存関係

```
タスク1 (ChainEffectInventory) ─┬─→ タスク4 (InventoryManager統合)
                                ├─→ タスク2 (AgentSlotチェイン効果スロット) ─┐
                                │                                           │
タスク5 (コアユニーク制約) ──────┼───────────────────────────────────────────┤
                                │                                           ↓
                                ├─→ タスク6a (チェイン効果基本管理)
                                │   タスク6b (ユニーク制約・BuildForBattle)
                                │
タスク3+15 (SkillInventory変更   │
         + 全呼び出し元修正) ────┤
                                │
タスク7 (セーブデータ構造) ──────┼─→ タスク8 (変換関数)
                                │
                                ├─→ タスク9 (報酬システム)
                                ├─→ タスク10 (初期エージェント)
                                │
                                ├─→ タスク11 (TUI スキルタブUI除去)
                                ├─→ タスク12 (TUI チェイン効果タブ追加)
                                └─→ タスク13 (TUI カスタマイズ対応)

タスク14 (app層統合) ← 全タスク完了後
```

**実施順序**: 1 → 2, 3+15, 5, 7 (並行可能) → 4, 6a → 6b → 8, 9, 10 → 11, 12, 13 → 14

## 実装タスク

### タスク1: ChainEffectInventoryドメインモデル作成
- **対象**: `internal/domain/chain_effect_inventory.go`（新規作成）
- **内容**:
  - `ChainEffectInventory`構造体（`map[string]struct{}`）
  - `NewChainEffectInventory()` コンストラクタ
  - `AddChainEffect(typeID string) bool`（空文字列は拒否してfalseを返す）
  - `HasChainEffect(typeID string) bool`
  - `GetOwnedChainEffects() []string`（ソート済み）
- **関連テスト**: 基準1, 2
- **状態**: 完了

### タスク2: AgentSlotにチェイン効果スロット追加
- **対象**: `internal/domain/agent_slot.go`
- **内容**:
  - `ChainEffectSlotConfig`構造体を追加（TypeID string）
  - `AgentSlot`に`ChainEffects [MaxSkillSlotCount]ChainEffectSlotConfig`を追加
  - `NewAgentSlot`で`ChainEffects`配列は値型のためゼロ値初期化（全TypeID=""）
  - `SetChainEffect(skillSlot int, typeID string)`メソッド追加
  - `ClearChainEffect(skillSlot int)`メソッド追加
  - `GetChainEffect(skillSlot int) *ChainEffectSlotConfig`メソッド追加
  - `ClearSkill`を修正: 同インデックスのチェイン効果も自動クリア
  - `SkillSlotConfig`から`ChainEffectID`フィールドを除去
- **関連テスト**: 基準8, 9, 10（ドメインレベル）
- **依存**: なし（ChainEffectSlotConfigは単純な構造体）
- **状態**: 完了

### タスク3+15: SkillInventoryからChainVariations除去 + 全呼び出し元修正
- **対象**: `internal/domain/skill_inventory.go` + 全パッケージの呼び出し元
- **内容**:
  - `SkillOwnership`から`ChainVariations`フィールドを除去
  - `AddSkill`のシグネチャを`AddSkill(typeID string)`に変更（breaking change）
  - `GetChainVariations`, `HasChainVariation`, `AddChainVariation`メソッドを除去
  - 既存テスト(`skill_inventory_test.go`)の更新
  - **同時に全呼び出し元を修正**（シグネチャ変更はbreaking changeのため）:
    - `AddSkill(typeID, chainEffectID)` → `AddSkill(typeID)`
    - `SetSkill(slot, skillSlot, typeID, chainEffectID)` → `SetSkill(slot, skillSlot, typeID)`
    - 各テストファイルの呼び出し修正
    - `GetChainVariations()`呼び出しの除去
  - 不要になったChainVariations関連テストの除去
- **関連テスト**: 基準21, 22
- **注意**: シグネチャ変更の影響範囲が広いため、一度に実施してビルドを通す
- **breaking changeチェックリスト**:
  - [ ] `domain/skill_inventory.go`: AddSkill, SkillOwnership
  - [ ] `domain/skill_inventory_test.go`: 全テスト関数
  - [ ] `usecase/slot/agent_slot_manager.go`: SetSkill
  - [ ] `usecase/slot/agent_slot_manager_test.go`: SetSkill呼び出し
  - [ ] `usecase/slot/chain_effect_slot_test.go`: AddSkill, SetSkill呼び出し
  - [ ] `usecase/inventory/inventory_manager.go`: AddSkill
  - [ ] `usecase/inventory/inventory_manager_test.go`: AddSkill呼び出し
  - [ ] `usecase/rewarding/reward.go`: AddRewardsToInventory内
  - [ ] `usecase/rewarding/reward_test.go`: テスト呼び出し
  - [ ] `app/root_model.go`: debugモード初期化
  - [ ] `infra/startup/startup.go`: 初期化処理
  - [ ] `tui/screens/agent_customization.go`: SetSkill呼び出し
  - [ ] `integration_test/`: 各統合テスト
  - [ ] 修正後に`go build ./...`と`go vet ./...`でビルド確認
- **状態**: 完了

### タスク4: InventoryManagerへのChainEffectInventory統合
- **対象**: `internal/usecase/inventory/inventory_manager.go`
- **内容**:
  - `chainEffects *domain.ChainEffectInventory`フィールドを追加
  - `AddChainEffect(typeID string) bool`メソッドを追加
  - `ChainEffects() *domain.ChainEffectInventory`アクセサを追加
  - `NewInventoryManagerWithInventories`に3番目の引数を追加
  - nil渡し時のデフォルトChainEffectInventory作成
- **関連テスト**: 基準3
- **依存**: タスク1
- **状態**: 完了

### タスク5: AgentSlotManagerにコアユニーク制約追加
- **対象**: `internal/usecase/slot/agent_slot_manager.go`
- **内容**:
  - `ErrCoreAlreadyEquipped`エラー定義
  - `SetCore`にユニーク制約を追加（他スロットで使用中のコアを拒否）
  - `ClearCore`で制約を解放
- **関連テスト**: 基準4, 5
- **依存**: なし（既存のSetCore/ClearCoreの修正）
- **状態**: 完了

### タスク6a: AgentSlotManagerにチェイン効果基本管理追加
- **対象**: `internal/usecase/slot/agent_slot_manager.go`
- **内容**:
  - `chainEffectInv *domain.ChainEffectInventory`フィールドを追加
  - `chainEffects map[string]domain.ChainEffect`フィールドを追加（マスタデータ参照用）
  - `NewAgentSlotManagerWithChainEffectInv`コンストラクタ
  - `ErrChainEffectNotOwned`, `ErrSkillNotSetForChain`エラー定義
  - `SetChainEffect(slot, skillSlot int, chainEffectTypeID string) error`
    - スキル未設定チェック
    - インベントリ保有チェック
  - `ClearChainEffect(slot, skillSlot int) error`
  - `SetSkill`のシグネチャから`chainEffectID`を除去
  - `ClearSkill`で同スロットのチェイン効果も自動クリア
  - `ClearCore`でスキル＆チェイン効果の連鎖クリア
- **関連テスト**: 基準8, 9, 10
- **依存**: タスク1, タスク2
- **状態**: 完了

### タスク6b: AgentSlotManagerにユニーク制約・BuildForBattle対応追加
- **対象**: `internal/usecase/slot/agent_slot_manager.go`
- **内容**:
  - `ErrChainEffectAlreadyEquipped`エラー定義
  - `SetChainEffect`にユニーク制約チェック追加（全スロット横断）
  - `buildModuleFromConfig`のチェイン効果解決を`AgentSlot.ChainEffects`配列から取得するように変更
    - `buildModuleFromConfig`は`AgentSlot.GetChainEffect(skillSlotIndex)`で`ChainEffectSlotConfig`を取得
    - `ChainEffectSlotConfig.TypeID`で`m.chainEffects`マスタデータマップから`domain.ChainEffect`を解決
- **関連テスト**: 基準6, 7, 26
- **依存**: タスク2（ChainEffectSlotConfig構造体）, タスク6a（SetChainEffect基本実装）
- **状態**: 完了

### タスク7: セーブデータ構造の更新
- **対象**: `internal/infra/savedata/savedata.go`
- **内容**:
  - `ChainEffectInventorySave`構造体を追加
  - `ChainEffectSlotSaveCfg`構造体を追加
  - `InventorySaveData`に`UniqueChainEffects *ChainEffectInventorySave`を追加
  - `AgentSlotSave`に`ChainEffects [4]ChainEffectSlotSaveCfg`を追加
  - `SkillInventorySave`の`Skills`を`map[string][]string` → `[]string`に変更
  - `SkillSlotSaveCfg`から`ChainEffectID`を除去
  - `NewSaveData()`を更新（UniqueChainEffectsの初期化）
  - `UniqueChainEffects`がnilの旧セーブデータ読込時にデフォルト値を補完
- **関連テスト**: 基準14, 15
- **依存**: なし（独立して実施可能）
- **状態**: 完了

### タスク8: セーブデータ変換関数の更新
- **対象**: `internal/infra/savedata/unique_inventory_converter.go`
- **内容**:
  - `ConvertChainEffectInventoryToSave` / `ConvertSaveToChainEffectInventory`関数を追加
  - `ConvertSkillInventoryToSave`を新フォーマット対応に変更
  - `ConvertSaveToSkillInventory`を新フォーマット対応に変更
  - `ConvertAgentSlotToSave` / `ConvertSaveToAgentSlot`にチェイン効果スロット対応追加
- **レイヤー責務の分離**:
  - infra層（このタスク）: セーブDTO ↔ ドメインオブジェクトの単純なフィールドマッピング
  - app層（タスク14）: 複数インベントリのセーブ/ロードオーケストレーション、セーブデータ全体の組み立て
  - converter関数はドメインロジックを持たず、純粋なデータ変換のみ行う
- **関連テスト**: 基準14, 15
- **依存**: タスク7, タスク1, タスク2, タスク3+15
- **状態**: 完了

### タスク9: 報酬システムのチェイン効果分離
- **対象**: `internal/usecase/rewarding/reward.go`（usecaseレイヤーに所属）
- **内容**:
  - 既存の`AddRewardsToInventory`を修正してChainEffectInventory引数を追加
    （後方互換のための旧関数残置は行わない — CLAUDE.md規約に従い一本化）
  - スキルは`SkillInventory.AddSkill(typeID)`に追加（chainEffectID引数なし）
  - チェイン効果は`ChainEffectInventory.AddChainEffect(chainEffectID)`に追加
  - 呼び出し元（app層）のシグネチャ変更はタスク14で対応
- **注意**: 報酬ロジックはusecaseレイヤーに正しく配置済み（`internal/usecase/rewarding/`）
- **関連テスト**: 基準16
- **依存**: タスク1, タスク3+15
- **状態**: 完了

### タスク10: 初期エージェントのコア分離
- **対象**:
  - `internal/infra/masterdata/data/first_agent.json`
  - `internal/infra/startup/startup.go`
- **内容**:
  - `first_agent.json`: 3エージェントに異なるコアを設定
    - agent_first_1: `all_rounder`（physical_strike_lv1互換）
    - agent_first_2: `healer`（heal_lv1互換）
    - agent_first_3: `quick_recoverer`（str_buff_lv1互換）
  - `InitializeNewGame`: 3スロットに異なるコアを設定
  - `InitializeNewGame`: コアインベントリに3種類のコアを追加
  - `InitializeNewGame`: SkillInventorySaveの新フォーマット対応
  - `InitializeNewGame`: ChainEffectInventorySaveの初期化追加
- **関連テスト**: 基準17, 18
- **依存**: タスク7（SkillInventorySaveの新フォーマット）
- **状態**: 完了

### タスク11: TUI - スキル一覧タブからチェイン効果UI除去
- **対象**: `internal/tui/screens/inventory.go`
- **内容**:
  - `SkillInventoryItem`から`ChainCount`, `ChainVariations`フィールドを除去
  - `renderSkillListItems()`: スキル名後の`(N種)`表示を除去
  - `renderSkillPreviewContent()`: 「チェイン効果バリエーション:」セクション全体を除去
  - `updateSkillList()`: `GetChainVariations()`呼び出しと関連データ設定を除去
  - `showingSkillDetail`フラグ: チェイン効果以外の用途がなければ除去
- **関連テスト**: 基準19, 20, 21, 22, 23
- **依存**: タスク3+15（GetChainVariations除去後）
- **状態**: 完了

### タスク12: TUI - チェイン効果一覧タブ追加
- **対象**: `internal/tui/screens/inventory.go`
- **内容**:
  - `InventoryTab`に`TabChainEffectInventory`を追加
  - `ChainEffectInventoryItem`構造体を追加（Name, Description, TypeID, Equipped bool）
  - チェイン効果一覧のリスト表示処理を追加
  - チェイン効果の装備状態チェック（AgentSlotManagerから参照）
  - タブ切替の操作処理を更新
- **関連テスト**: 基準24, 25
- **依存**: タスク1, タスク6a
- **状態**: 完了

### タスク13: TUI - エージェントカスタマイズのチェイン効果対応
- **対象**: `internal/tui/screens/agent_customization.go`
- **内容**:
  - チェイン効果スロット表示を`ChainEffects`配列からの参照に変更
  - チェイン効果選択モード（`ModeChainSelect`）を独立フローに変更
  - スキル未設定時のチェイン効果選択無効化
  - コア選択リストで他スロット使用中を「(装備中)」表示
  - チェイン効果選択リストで他スロット使用中を「(装備中)」表示
  - `SetSkill`の呼び出しからchainEffectID引数を除去
- **関連テスト**: 基準11, 12, 13
- **依存**: タスク6b
- **状態**: 完了

### タスク14: app層の統合更新 + debugモード対応
- **対象**: `internal/app/root_model.go`, `internal/app/screen_factory.go`, `internal/tui/presenter/debug_inventory_provider.go`
- **内容**:
  - `ChainEffectInventory`の初期化・受け渡し
  - `AgentSlotManager`のコンストラクタ変更に追従（`NewAgentSlotManagerWithChainEffectInv`を使用）
  - `InventoryManager`のコンストラクタ変更に追従（3引数版を使用）
  - `AddRewardsToInventory`の新シグネチャに追従
  - セーブ/ロード処理にチェイン効果インベントリを追加
  - **debugモード対応（基準27, 28, 29）**:
    - 旧方式の`invManager.AddSkill(mt.ID, ce.ID)`ループを除去
    - 全スキルは`invManager.AddSkill(mt.ID)`で追加（チェイン効果なし）
    - 全チェイン効果を`invManager.AddChainEffect(ce.ID)`で独立追加
    - `DebugInventoryProvider.GetChainEffects()`はマスタデータから全チェイン効果を返す（現行通り動作）
- **関連テスト**: 全基準（統合） + 基準27, 28, 29（debugモード）
- **依存**: 全タスク完了後
- **状態**: 完了

## TUIテスト手順

### 基準11: カスタマイズ画面で3種個別選択
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# 初期状態: 3エージェントがそれぞれコア+スキル1つ装備済み
# メニューからカスタマイズ画面へ遷移
send_keys("jj\r", delay=2000)  # エージェントカスタマイズを選択
# エージェントカードを選択
send_keys("\r", delay=500)
# コア選択が表示される
assert_contains("コア")
# スキル選択に移動
send_keys("\x1b", delay=300)  # 戻る
send_keys("j\r", delay=300)    # スキルスロットを選択
assert_contains("スキル")
# チェイン効果選択に移動
send_keys("\x1b", delay=300)
send_keys("jj\r", delay=300)   # チェイン効果スロットを選択
assert_contains("チェイン効果")
```

### 基準12: スキル未設定でチェイン効果選択不可
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# 初期状態: スロット0のスキルスロット1は未設定
send_keys("jj\r", delay=2000)  # カスタマイズ画面へ
send_keys("\r", delay=500)      # エージェント1を選択
# 未設定スキルスロットのチェイン効果スロットに移動
# → チェイン効果選択リストが表示されないこと（操作がブロックされる）
# オラクル: 画面上にチェイン効果選択リストが表示されず、フォーカスが移動しない
#           またはスキル未設定を示すメッセージが表示される
```

### 基準13: 使用中コア・チェイン効果の区別表示
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# 初期状態: 3エージェントがそれぞれ異なるコアを装備済み
send_keys("jj\r", delay=2000)  # カスタマイズ画面へ
send_keys("\r", delay=500)      # エージェント1を選択
# コア選択リストで他エージェントに装備中のコアが区別表示されること
assert_contains("装備中")
```

### 基準19: スキル一覧のチェイン効果数表示なし
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# 初期状態: physical_strike_lv1, heal_lv1, str_buff_lv1 を保有
send_keys("j\r", delay=2000)    # インベントリを選択
send_keys("\t", delay=300)       # スキル一覧タブへ移動
# スキル名の後に(N種)が表示されていないことを確認
assert_not_contains("種)")
```

### 基準20: スキル詳細のチェイン効果セクションなし
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("j\r", delay=2000)     # インベントリ
send_keys("\t", delay=300)        # スキルタブ
send_keys("\r", delay=300)        # スキル詳細表示（Enterで詳細トグル）
# チェイン効果バリエーションセクションが表示されていないこと
assert_not_contains("チェイン効果バリエーション")
```

### 基準24: チェイン効果一覧タブの存在
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
send_keys("j\r", delay=2000)     # インベントリ画面へ
send_keys("\t\t", delay=300)      # 3番目のタブへ
# チェイン効果タブが存在し表示されること
assert_contains("チェイン効果")
```

### 基準25: チェイン効果一覧の名前・説明・装備状態
```
# debugモードで起動（全チェイン効果が自動追加される）
launch_tui(command="go run ./cmd/BlitzTypingOperator -debug", mode="buffer", dimensions="160x45")
send_keys("j\r", delay=2000)     # インベントリ
send_keys("\t\t", delay=300)      # チェイン効果タブ
# オラクル: チェイン効果のマスタデータ名称が表示される
assert_contains("ダメージボーナス")  # chain_effects.jsonの名称
# 装備状態が確認できること
# オラクル: 装備済みのチェイン効果には装備マーカーが表示される
```

### 基準27-28: debugモードでの全チェイン効果選択
```
launch_tui(command="go run ./cmd/BlitzTypingOperator -debug", mode="buffer", dimensions="160x45")
# debugモードで起動 → 全チェイン効果がインベントリに追加されていること

# インベントリ画面でチェイン効果タブを確認（基準27）
send_keys("j\r", delay=2000)     # インベントリ画面へ
send_keys("\t\t", delay=300)      # チェイン効果タブ
# 全チェイン効果（damage_bonus, damage_amp, armor_pierce等）が表示されること
assert_contains("ダメージボーナス")

# カスタマイズ画面でチェイン効果が全て選択可能（基準28）
send_keys("\x1b", delay=300)      # ホームに戻る
send_keys("jj\r", delay=300)     # カスタマイズ画面へ
# スキル設定後、チェイン効果選択で全チェイン効果が表示されること
# 操作手順はUI実装後に確定
```

## 進捗ログ
- タスク1完了: ChainEffectInventoryドメインモデル作成（domain/chain_effect_inventory.go）
- タスク2完了: AgentSlotにChainEffectSlotConfig/ChainEffects配列追加、SetChainEffect/ClearChainEffect/GetChainEffectメソッド追加
- タスク3+15完了: SkillOwnershipからChainVariations除去、AddSkill/SetSkillシグネチャ簡略化、全呼び出し元修正（~25ファイル）
- タスク4完了: InventoryManagerにchainEffectsフィールド追加、ChainEffects()アクセサ、AddChainEffect()メソッド、NewInventoryManagerWithInventories 3引数化
- タスク5完了: SetCoreにコアユニーク制約追加（ErrCoreAlreadyEquipped）、isCoreEquippedElsewhere()
- タスク6a完了: AgentSlotManagerにchainEffectInvフィールド追加、SetChainEffect/ClearChainEffect実装
- タスク6b完了: チェイン効果ユニーク制約追加（ErrChainEffectAlreadyEquipped）、buildModuleFromConfigのチェイン効果解決
- タスク7完了: ChainEffectInventorySave/ChainEffectSlotSaveCfg追加、SkillInventorySave.Skillsを[]stringに変更、SkillSlotSaveCfgからChainEffectID除去
- タスク8完了: ChainEffect変換関数追加、AgentSlot変換にチェイン効果対応追加
- タスク9完了: AddRewardsToInventoryにChainEffectInventory引数追加、チェイン効果の独立追加ロジック実装、session.InventoryManagerにchainEffectsフィールド追加
- タスク10完了: first_agent.jsonの3エージェントに異なるコア設定（all_rounder/healer/quick_recoverer）、InitializeNewGameの3種コア対応
- タスク11完了: showingSkillDetail除去、Enter→詳細トグルの除去、ヒントテキスト更新
- タスク12完了: TabChainEffectInventory追加、ChainEffectInventoryItem構造体、チェイン効果タブのレンダリング
- タスク13完了: focusPosition拡張（0=コア,奇数=スキル,偶数=チェイン効果）、独立チェイン効果選択フロー、「(装備中)」表示、後方互換エイリアス削除
- タスク14完了: UniqueChainEffectsのセーブ/ロード対応、ChainEffectsスロットのセーブ/ロード、debugモードで全チェイン効果追加、SlotManagerへのChainEffectInventory設定
