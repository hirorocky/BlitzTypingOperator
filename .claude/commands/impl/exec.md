---
description: 実装計画に沿ってTDDで実装する
allowed-tools: Read, Write, Edit, MultiEdit, Bash, Glob, Grep, WebSearch, WebFetch, TaskCreate, TaskUpdate, TaskList, TaskGet
argument-hint: <feature-name>
---

# 計画実行

## 概要
plan.mdの実装タスクに沿って、TDDサイクルで実装する。
進捗はplan.mdとTask toolの両方で追跡する。

## 実行ステップ

### 1. コンテキスト読み込み
- `docs/project/$ARGUMENTS/feature.md` を読む（存在しなければエラー）
- `docs/project/$ARGUMENTS/plan.md` を読む（存在しなければエラー）
- `docs/specification/` から関連するドメイン仕様を読む

### 2. タスク登録
plan.mdの「未着手」タスクをTaskCreateで登録する（最大10タスク）。
すでに「完了」のタスクは登録しない（途中再開に対応）。

### 3. TDD実装ループ
plan.mdのタスク順に、Kent BeckのTDDサイクルで実装する:

```
🔴 RED      → 失敗するユニットテストを書く
🟢 GREEN    → テストを通す最小限のコードを書く
🔵 REFACTOR → コード改善（テスト成功維持）
```

#### ルール
- **必ずテストを先に書く**（例外なし）
- **小さなステップ**で進める
- **全テスト実行**で回帰を確認（`go test ./...`）
- 意味のある単位でgit commitする

#### タスク完了時
- TaskUpdateでタスクを `completed` に更新
- plan.mdの該当タスクの状態を「完了」に更新
- plan.mdの進捗ログに実施内容を追記

### 4. 10タスク制限
10タスクを完了したら停止し、ユーザーに中間報告する。
未完了タスクがある場合は、再度 `/impl:exec $ARGUMENTS` で続行できることを案内する。

### 5. 完了報告
全タスク完了後（または10タスク完了後）:
- plan.mdの最終状態を表示
- 実装したファイルの一覧
- テスト結果のサマリー（`go test ./...`）
- 次のステップ: 未完了タスクがあれば `/impl:exec $ARGUMENTS`、全完了なら `/impl:review $ARGUMENTS`

## タスク管理
- TaskCreateで実装タスクを登録（最大10）
- 各タスクの開始時に `in_progress`、完了時に `completed` に更新
- activeFormで現在の作業内容を表示

## plan.md更新ルール
- タスク完了時: 状態を「完了」に更新
- 進捗ログ: 各タスクの完了内容を簡潔に記録
- 新たな課題発見時: 実装タスクセクションにタスクを追加

## テスト品質
- AAA（Arrange-Act-Assert）パターン
- 1テスト = 1振る舞い
- エッジケース考慮（nil、空、境界値、エラー）
- モックは最小限
