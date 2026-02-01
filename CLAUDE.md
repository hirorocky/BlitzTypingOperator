# AI-DLC and Spec-Driven Development

Kiro-style Spec Driven Development implementation on AI-DLC (AI Development Life Cycle)

## Project Context

### Paths
- Steering: `.kiro/steering/`
- Specs: `.kiro/specs/`

### Steering vs Specification

**Steering** (`.kiro/steering/`) - Guide AI with project-wide rules and context
**Specs** (`.kiro/specs/`) - Formalize development process for individual features

### Active Specifications
- Check `.kiro/specs/` for active specifications
- Use `/kiro:spec-status [feature-name]` to check progress

## Development Guidelines
- Think in English, generate responses in Japanese. All Markdown content written to project files (e.g., requirements.md, design.md, tasks.md, research.md, validation reports) MUST be written in the target language configured for this specification (see spec.json.language).

## Coding Style

### コメントの言語
- コード内のコメントは日本語で記述すること

### コメントに含めてはいけない情報
以下のような実装プロセスに関する情報はコメントに含めないこと：
- 要件番号（例: `Requirement 12.1`, `Requirements 4.1-4.7`）
- タスク番号（例: `Task 7.2`, `Phase 1`）
- 仕様書への参照（例: `REQ-001`, `設計書3.2節`）
- バージョン情報（例: `v4.0.0:`, `v3.0.0以降は`）
- 「新仕様」「旧仕様」などの一時的な差分を示す表現

### コメントに含めるべき情報
- **何をするか**（処理の目的）
- **なぜそうするか**（設計意図、制約）
- **どう動くか**（複雑なロジックの説明）

### 良い例・悪い例

```go
// 悪い例（実装寄り）
// Requirement 11.15: HP50%以下で強化フェーズに移行
const EnhanceThreshold = 0.5

// 良い例（ドメインロジックの説明）
// 敵が強化フェーズに移行するHP割合の閾値（50%）
const EnhanceThreshold = 0.5
```

```go
// 悪い例（実装寄り）
// Requirements 6.8-6.16に基づいて設計されています。
type ModuleModel struct { ... }

// 良い例（型の説明）
// モジュールはエージェント合成時にコアに装備され、バトル中に使用可能なスキルになります。
type ModuleModel struct { ... }
```

### 残してはいけないコード
以下のようなコードは作成せず、発見した場合は削除すること：
- 後方互換性のためだけに残されたコード・フィールド・関数
- 「レガシー」「廃止予定」とコメントされたコード
- 使用されていない古い形式のデータ構造
- 新旧両方の実装が並存している状態

## Minimal Workflow
- Phase 0 (optional): `/kiro:steering`, `/kiro:steering-custom`
- Phase 1 (Specification):
  - `/kiro:spec-init "description"`
  - `/kiro:spec-requirements {feature}`
  - `/kiro:validate-gap {feature}` (optional: for existing codebase)
  - `/kiro:spec-design {feature} [-y]`
  - `/kiro:validate-design {feature}` (optional: design review)
  - `/kiro:spec-tasks {feature} [-y]`
- Phase 2 (Implementation): `/kiro:spec-impl {feature} [tasks]`
  - `/kiro:validate-impl {feature}` (optional: after implementation)
- Progress check: `/kiro:spec-status {feature}` (use anytime)

## Development Rules
- 3-phase approval workflow: Requirements → Design → Tasks → Implementation
- Human review required each phase; use `-y` only for intentional fast-track
- Keep steering current and verify alignment with `/kiro:spec-status`
- Follow the user's instructions precisely, and within that scope act autonomously: gather the necessary context and complete the requested work end-to-end in this run, asking questions only when essential information is missing or the instructions are critically ambiguous.

## Steering Configuration
- Load entire `.kiro/steering/` as project memory
- Default files: `product.md`, `tech.md`, `structure.md`
- Custom files are supported (managed via `/kiro:steering-custom`)

## TUI Testing (MCP: tui-test)
- 本プロジェクトは複雑なTUIを持つため、テスト時は必ず**バッファモード**を使用すること
- `launch_tui`実行時に `mode: "buffer"` を指定する
- バッファモードでは位置ベースのアサーション、カーソル追跡、領域抽出が可能

### 特殊キーの送信方法
| キー | 送信方法 | 備考 |
|------|----------|------|
| Enter | `\r` | `\n`は機能しない（bubbleteaはCRを使用） |
| ESC | `\x1b` | 画面遷移の「戻る」操作に使用 |
| Tab | `\t` | |
| Ctrl+X | `send_ctrl("x")` | 例: Ctrl+C は `send_ctrl("c")` |

### 使用例
```
# ゲーム起動（常に最新コードを使用）
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")

# メニュー操作
send_keys("j")      # 下に移動
send_keys("\r")     # Enter（選択）
send_keys("\x1b")   # ESC（戻る）
```

※ 初回起動時はコンパイルのため数秒かかる。delayを長め（2秒程度）に設定すると安定する。

