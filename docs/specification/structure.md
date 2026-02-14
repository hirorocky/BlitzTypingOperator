# プロジェクト構造

## 組織哲学

5層レイヤードアーキテクチャを採用し、ドメインロジックとUIを明確に分離。
Elm Architectureパターンに基づくイベント駆動型設計で、状態管理を一元化。

## レイヤー構造

```
app層        ← 全ての層に依存可能（オーケストレーション）
    ↓
tui層        ← domain, usecase, config に依存
    ↓
usecase層    ← domain, domain/service, config に依存
    ↓
infra層      ← domain, config に依存
    ↓
domain層     ← 外部依存なし
    ↓
config       ← 横断的関心事（全層から参照可能）
```

### レイヤー間の依存ルール

| 層 | 許可された依存先 | 禁止されている依存先 |
|----|------------------|---------------------|
| domain（VO・エンティティ） | なし | domain/service, usecase, infra, tui, app |
| domain/service | domain | usecase, infra, tui, app |
| usecase | domain, domain/service, config | infra, tui, app |
| infra | domain, config | domain/service, usecase, tui, app |
| tui | domain, usecase, config | infra, app |
| app | 全ての層 | なし |

## ディレクトリパターン

### app層 - オーケストレーション
**場所**: `/internal/app/`
**目的**: Bubbleteaのtea.Model実装とシーン管理。他の全ての層をオーケストレーション
**含まれるファイル**:
- `root_model.go`: BubbleteaのModel実装
- `scene.go`: シーン列挙型とChangeSceneMsg
- `scene_router.go`: シーン名からSceneへの変換
- `screen_factory.go`: 画面インスタンスの生成
- `screen_map.go`: シーンと画面のマッピング
- `message_handlers.go`: Bubbleteaメッセージハンドリング
- `masterdata_converter.go`: masterdata→domain型変換ヘルパー。以下の変換関数を提供：
  - `ConvertTimedEffects`: timed_effects.jsonの時限効果定義をマップに変換
  - `ResolveModuleTimedEffects`: スキルのEffectColumnSpec.Column/ValueをTimedEffectから解決（ManaCost/ManaGainフィールド込み）
  - `ResolveEnemyActionTimedEffects`: 敵行動のEffectColumn/EffectValueをTimedEffectから解決
- `*_adapter.go`: 層間データ変換アダプター（例: enemy_progress_adapter.go）

### domain層 - ドメインモデル
**場所**: `/internal/domain/`
**目的**: VO、エンティティ。UIやインフラに依存しない純粋なドメイン層
**例**: `core.go`（コア特性）、`skill.go`（スキル）、`agent.go`（エージェント）、`enemy.go`（敵）

**主要なサブシステム**:
- **エンティティ**: core.go, skill.go, agent.go, enemy.go, player.go
- **スロットシステム**: agent_slot.go（AgentSlot、SkillSlotConfig、ChainEffectSlotConfig）- 3スロットのエージェント構成管理
- **インベントリ**:
  - core_inventory.go: コアの保有管理（TypeIDの保有フラグ）
  - skill_inventory.go: スキルのユニーク管理（TypeIDのみ）
  - chain_effect_inventory.go: チェイン効果の保有管理（TypeIDの保有フラグ）
- **敵進行システム**: enemy_progress.go - 敵撃破記録・ランク進行・選択可能レベル範囲
- **効果システム**: effect_table.go, effect_column.go, effect_context.go, effect_entry.go
  - EffectTableパターン: バフ、デバフ、パッシブ、チェイン効果を統一的に管理
  - 列指向設計: 効果種別を EffectColumn として定義
  - AddBuff/AddDebuff: 同一TimedEffectID時の重複判定・Duration加算（上限99.9秒）
- **時限効果**: timed_effect.go - マスタデータで定義される一時ステータス（バフ/デバフ）。ID、効果列、効果値を保持
- **効果説明**: effect_description.go（EffectColumn/効果値からUI向け説明テキストを生成）
- **チェイン効果**: chain_effect.go（スキル使用後、次のスキル使用まで待機し、他エージェントの行動で発動する追加効果）
- **パッシブスキル**: passive_skill.go（トリガー/効果タイプ定義）, passive_evaluator.go（条件判定・効果適用）
- **敵行動システム**: enemy.go 内に行動パターン（EnemyAction）、フェーズ遷移、チャージ/ディフェンス状態管理を含む。バフ/デバフ行動の効果は時限効果IDで参照

**サブパッケージ**:
- `/internal/domain/service/` - ドメインサービス
  - `stats_service.go`: ステータス計算（CoreType→Stats、基礎値100×重み）

### usecase層 - ユースケース
**場所**: `/internal/usecase/`
**目的**: ドメインオブジェクト + ドメインサービスを組み合わせたアプリケーション固有の処理フロー
**サブパッケージ**（動詞形でユースケースを表現）:
- `combat`: バトル実行
  - `combat/chain`: チェイン効果管理（ChainEffectManager）
  - `combat/recast`: リキャスト管理（RecastManager）
  - `combat/voltage`: ボルテージ管理（VoltageManager）- 時間経過によるダメージ乗算
- `typing`: タイピング評価
- `slot`: エージェントスロット管理（AgentSlotManager）- 3スロットのコア・スキル・チェイン効果付け替え
- `inventory`: インベントリ管理（InventoryManager）- コア・スキル・チェイン効果保有状態の統合管理
- `spawning`: 敵生成
- `rewarding`: 報酬計算・ドロップ
- `progress`: 敵進行管理（EnemyProgressManager）- 撃破状況・ランク進行・HP成長
- `achievement`: 実績解除
- `session`: セッション管理（統計・設定含む）

### infra層 - インフラストラクチャ
**場所**: `/internal/infra/`
**目的**: 外部リソース（ファイル、ターミナル等）とのやり取り
**サブパッケージ**:
- `infra/savedata/`: セーブ/ロード永続化
  - `savedata.go`: セーブデータ構造体（CoreInventorySave, SkillInventorySave, ChainEffectInventorySave, AgentSlotSave）
  - `unique_inventory_converter.go`: 各インベントリのセーブ/ロード変換関数
- `infra/masterdata/`: JSONマスタデータローダー＋埋め込みデータ（Go embed.FS）
  - timed_effects.json: 時限効果定義（ID、名前、説明、効果列、効果値）
  - skills.json: スキル定義（skill_typesキーで定義、ManaCost/ManaGainフィールド、各effect_columnにtimed_effect_idを参照）
  - enemy_actions.json: 敵行動定義（各バフ/デバフ行動にtimed_effect_idを参照）
  - chain_effects.json: チェイン効果定義（ID、名前、説明、効果）
- `infra/errorhandler/`: エラーハンドリング
- `infra/startup/`: 起動処理
- `infra/terminal/`: ターミナル環境検証

### tui層 - UI
**場所**: `/internal/tui/`
**目的**: 各シーンの画面実装、コンポーネント、スタイル、プレゼンター
**サブディレクトリ**:
- `screens/`: 各シーンの画面実装（Bubbleteaの`tea.Model`実装）
  - 画面タイプ: home, battle_select, battle, agent_customization, inventory, reward, encyclopedia, settings, stats_achievements
  - 大きな画面は分割: battle.go（状態）、battle_view.go（描画）、battle_logic.go（ロジック）
  - agent_customization: 3スロットエージェントのコア・スキル・チェイン効果付け替えUI
  - inventory: コア・スキル・チェイン効果の所持一覧表示（3タブ構成）
- `components/`: 再利用可能なUIコンポーネント
  - 基本コンポーネント: components.go
  - 専用コンポーネント: hp_display.go, recast_progress_bar.go, chain_effect_badge.go, passive_skill_notification.go
  - 表示ヘルパー: position_indicator.go（位置インジケーター）, stats_display.go（ステータス表示）
- `styles/`: lipglossスタイル定義（カラーパレット含む）
- `presenter/`: UI向けデータ変換（GameState→ViewModel）
- `ascii/`: ASCIIアート

### config - 横断的関心事
**場所**: `/internal/config/`
**目的**: マジックナンバーを一元管理。バトル設定、効果持続時間、インベントリ設定等
**含まれるファイル**:
- `constants.go`: 定数（`BattleTickInterval`, `DefaultModuleCooldown`, `MaxAgentEquipSlots`, `MaxStatusDuration` など）
- `balance.go`: ゲームバランス設定（旧usecase/balance）

### integration_test - 統合テスト
**場所**: `/internal/integration_test/`
**目的**: 複数層にまたがる統合テスト

### cmd - エントリーポイント
**場所**: `/cmd/BlitzTypingOperator/`
**目的**: `main.go`のみ。アプリケーション起動のみを担当

## 命名規則

- **ファイル**: snake_case（例: `game_state.go`, `battle_select.go`）
- **構造体**: PascalCase（例: `CoreModel`, `BattleScreen`）
- **関数**: PascalCase（エクスポート）/ camelCase（非エクスポート）
- **テスト**: `*_test.go`で同一ディレクトリに配置

## インポート組織

```go
import (
    // 標準ライブラリ
    "fmt"
    "time"

    // 外部パッケージ
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"

    // プロジェクト内パッケージ
    "hirorocky/type-battle/internal/domain"
    "hirorocky/type-battle/internal/tui/screens"
)
```

**パスエイリアス**: なし（Go標準のモジュールパスを使用）

## コード組織原則

1. **ドメイン層の独立性**: `/internal/domain/`は他の内部パッケージに依存しない
2. **ドメインサービスの分離**: 複数ドメインオブジェクトの組み合わせロジックは`domain/service/`に配置
3. **App層での型変換**: infra→domain型変換、データ解決はapp層で実施し、usecase層のinfra依存を排除
4. **画面の自己完結性**: 各画面は独立して動作可能。RootModelがルーティングを担当
5. **外部データ駆動**: ゲームコンテンツ（コア、モジュール、敵、時限効果）はJSONファイルで定義
6. **テストの同居**: テストファイルは実装と同じディレクトリに配置
7. **プレゼンター層の活用**: UI向けデータ変換は`tui/presenter/`で実装
8. **定数の一元管理**: マジックナンバーはconfigパッケージに集約
9. **ハンドラーマップパターン**: シーン遷移・メッセージ処理はマップ駆動で分岐

## ドメイン別仕様

各ドメインの詳細な要件・仕様は `docs/specification/` 配下を参照:

- `battle.md`: バトルシステム
- `gameloop.md`: ゲームループ・状態遷移
- `agent.md`: エージェント・スロットシステム
- `typing.md`: タイピング評価・入力処理
- `enemy.md`: 敵・ステージシステム
- `collection.md`: 図鑑・実績システム

---
_パターンを記述。新規ファイルがパターンに従えばドキュメントの更新は不要_
