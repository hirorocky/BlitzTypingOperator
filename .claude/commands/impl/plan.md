---
description: 受け入れテストを作成し、実装計画（plan.md）を策定する
allowed-tools: Read, Write, Edit, MultiEdit, Bash, Glob, Grep, WebSearch, WebFetch, TaskCreate, TaskUpdate, TaskList, TaskGet
argument-hint: <feature-name>
---

# 実装計画

## 概要
フィーチャードキュメントの受け入れ基準に基づき、受け入れテストを作成し、実装計画（plan.md）を策定する。
ユーザーの承認なしに実装には進まない。

## 実行ステップ

### 1. コンテキスト読み込み
- `docs/project/$ARGUMENTS/feature.md` を読む（存在しなければエラー）
- `docs/specification/` から関連するドメイン仕様を読む
- 関連する既存コードを読み、現在の実装を理解する

### 2. 受け入れテストの設計
feature.mdの受け入れ基準ごとにテスト方法を決定する:

| テスト対象 | ツール | 場所 |
|-----------|--------|------|
| 画面表示・キー入力・遷移 | MCP tui-test (buffer mode) | セッション内で実行 |
| ドメインロジック・計算 | Go test | `*_test.go`（同ディレクトリ） |
| 複数層にまたがるフロー | Go test | `internal/integration_test/` |

### 3. 受け入れテストの作成
受け入れ基準ごとにテストを先に書く:

- **Goテスト**: テストファイルを作成し、テストが失敗することを確認（RED状態）
- **TUIテスト**: テスト手順をplan.mdに記載（実装後にセッション内で実行）

テスト作成後、`go test ./...` で既存テストに影響がないことを確認する。
（新規テストの失敗は想定内）

### 4. 実装計画の策定
受け入れテストと設計を分析し、TDDで実装するためのタスク分割を行う。
`docs/project/$ARGUMENTS/plan.md` を作成する:

```markdown
# 実装計画: {フィーチャー名}

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル | 状態 |
|---|------|-----------|---------------|------|
| 1 | {基準1} | Go test | {path} | 作成済 |
| 2 | {基準2} | TUI test | セッション内 | 手順記載済 |

## 実装タスク

### タスク1: {タスク名}
- **対象**: {変更/作成するファイル}
- **内容**: {何をするか}
- **関連テスト**: {どの受け入れ基準に対応するか}
- **状態**: 未着手

### タスク2: {タスク名}
...

## TUIテスト手順
（TUIテストが必要な受け入れ基準がある場合のみ）

### 基準{N}: {基準の内容}
\```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
{操作手順}
assert_contains("{期待する表示}")
\```

## 進捗ログ
（impl:exec が記録する。初期状態では空。）
```

タスク分割のガイドライン:
- 各タスクはTDDの1サイクル（RED-GREEN-REFACTOR）で完了できるサイズにする
- タスク間の依存関係を明示する（先に実装すべきものを上に配置）
- 各タスクに対応する受け入れ基準を明記する

### 5. エキスパートエージェントに計画レビューを依頼

test-qualityとarchitectureの2つのエキスパートエージェントをTask toolで**1つのメッセージで並列起動**し、計画をレビューする:

```
Task(subagent_type="test-quality", prompt="
MODE: review
REVIEW_TARGET:
{plan.md の全文}

FEATURE_DOC:
{feature.md の全文}
")

Task(subagent_type="architecture", prompt="
MODE: review
REVIEW_TARGET:
{plan.md の全文}

FEATURE_DOC:
{feature.md の全文}
")
```

レビュー結果に基づいて計画を修正する（MUST/SHOULDの指摘のみ）。

### 6. 結果の提示
以下をユーザーに提示する:
- 作成した受け入れテストの一覧
- plan.mdの内容
- エキスパートのレビュー結果と対応内容
- 次のステップ: `/impl:exec $ARGUMENTS`

**重要**: ユーザーの承認を待つ。承認なしに実装に進まない。

## 出力形式
- 受け入れテストの一覧（基準・種別・ファイル）
- plan.mdの全文
- エキスパートのフィードバックを観点ごとに引用形式で提示
- 次のコマンドをコードブロックで提示
