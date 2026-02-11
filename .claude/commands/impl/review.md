---
description: 受け入れテストを実行し、ユーザーフィードバックに基づいて修正する
allowed-tools: Read, Write, Edit, MultiEdit, Bash, Glob, Grep, TaskCreate, TaskUpdate, TaskList, TaskGet, mcp__tui-test__launch_tui, mcp__tui-test__send_keys, mcp__tui-test__send_ctrl, mcp__tui-test__capture_screen, mcp__tui-test__expect_text, mcp__tui-test__assert_contains, mcp__tui-test__assert_at_position, mcp__tui-test__get_cursor_position, mcp__tui-test__get_screen_region, mcp__tui-test__get_line, mcp__tui-test__close_session, mcp__tui-test__list_sessions
argument-hint: <feature-name>
---

# 実装レビュー

## 概要
受け入れテストを全て実行し、結果をユーザーに提示する。
ユーザーのフィードバックに基づいて修正を行い、全テスト通過とユーザー承認を得る。

## 実行ステップ

### 1. コンテキスト読み込み
- `docs/project/$ARGUMENTS/feature.md` を読む（存在しなければエラー）
- `docs/project/$ARGUMENTS/plan.md` を読む（存在しなければエラー）

### 2. Goテストの実行

```bash
go test ./...
```

全テスト結果を収集する。

### 3. TUIテストの実行
plan.mdにTUIテスト手順が記載されている場合、MCP tui-testで実行する:

```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
```

plan.mdの手順に沿ってTUI操作とアサーションを実行する。

**注意**: 初回起動時はGoのコンパイルのため数秒かかる。`launch_tui` 後の最初の操作では `delay` を長め（2秒程度）に設定すること。これを怠るとコンパイル完了前にキー入力が送信され、テストが不安定になる。

### 4. 受け入れテスト結果の報告
受け入れ基準ごとの合否を一覧で報告する:

```
## 受け入れテスト結果

| # | 基準 | テスト種別 | 結果 |
|---|------|-----------|------|
| 1 | {基準1} | Go test | ✅ PASS / ❌ FAIL |
| 2 | {基準2} | TUI test | ✅ PASS / ❌ FAIL |
```

失敗がある場合は、失敗の詳細（エラーメッセージ、期待値と実際値）も表示する。

### 5. ユーザーフィードバックの対応
ユーザーからフィードバックがあった場合:

1. フィードバック内容を分析
2. 修正を実施
3. 関連するテストを再実行して確認
4. plan.mdの進捗ログに修正内容を追記
5. ステップ2-4に戻り、再度全テスト結果を報告

フィードバックループは以下のいずれかで終了:
- ユーザーが承認した場合
- ユーザーが `/dev:complete $ARGUMENTS` への移行を指示した場合

### 6. 承認後の案内
次のステップを案内する: `/dev:complete $ARGUMENTS`

## 出力形式
- 受け入れテスト結果を表形式で提示
- 失敗テストの詳細をコードブロックで提示
- 修正内容のサマリーを変更ごとに提示
- 次のコマンドをコードブロックで提示
