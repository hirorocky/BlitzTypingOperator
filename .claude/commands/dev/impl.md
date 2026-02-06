---
description: 受け入れテスト駆動で実装する（ATDD）
allowed-tools: Read, Write, Edit, MultiEdit, Bash, Glob, Grep, WebSearch, WebFetch, TaskCreate, TaskUpdate, TaskList, TaskGet
argument-hint: <feature-name>
---

# ATDD実装

## 概要
フィーチャードキュメントの受け入れ基準に基づき、受け入れテスト駆動開発で実装する。

## 実行ステップ

### 1. コンテキスト読み込み
- `docs/project/$ARGUMENTS.md` を読む（存在しなければエラー）
- `docs/specification/` から関連するドメイン仕様を読む
- 受け入れ基準を抽出し、タスクとして登録する（TaskCreate）

### 2. 受け入れテストの作成
受け入れ基準ごとにテストを先に書く:
- **ドメインロジック**: Go統合テスト（`internal/integration_test/`）またはユニットテスト
- **UI動作**: MCP tui-test（buffer mode）でのTUIテスト

この時点で全テストは失敗する（RED状態）。

### 3. TDD実装ループ
各受け入れ基準に対して、Kent BeckのTDDサイクルで実装:

```
🔴 RED    → 失敗するユニットテストを書く
🟢 GREEN  → テストを通す最小限のコードを書く
🔵 REFACTOR → コード改善（テスト成功維持）
```

#### ルール
- **必ずテストを先に書く**（例外なし）
- **小さなステップ**で進める
- **全テスト実行**で回帰を確認
- 意味のある単位でgit commitする

### 4. 受け入れテストの自動実行
全実装が完了したら、受け入れテストを自動実行する:

1. `go test ./...` で全テスト実行
2. TUI受け入れテストがある場合は tui-test で実行
3. 受け入れ基準ごとの合否を一覧で報告:

```
## 受け入れテスト結果
| # | 基準 | 結果 |
|---|------|------|
| 1 | {基準1} | PASS / FAIL |
| 2 | {基準2} | PASS / FAIL |
```

### 5. 完了報告
- 実装したファイルの一覧
- テスト結果のサマリー
- 次のステップ: `/dev:complete $ARGUMENTS`

## タスク管理
- TaskCreateで受け入れ基準をタスクとして登録（最大10）
- 各タスクの開始時にin_progress、完了時にcompletedに更新
- activeFormで現在の作業内容を表示

## テスト品質
- AAA（Arrange-Act-Assert）パターン
- 1テスト = 1振る舞い
- エッジケース考慮（nil、空、境界値、エラー）
- モックは最小限
