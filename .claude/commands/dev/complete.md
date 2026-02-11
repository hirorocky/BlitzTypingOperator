---
description: 実装をcodexにレビューさせ、仕様を更新する
allowed-tools: Read, Write, Edit, Bash, Glob, Grep
argument-hint: <feature-name>
---

# 完了レビュー

## 概要

4つの観点別エキスパートエージェントを並列起動して実装をレビューし、問題がなければ仕様を更新してフィーチャーを完了する。
ユーザーの承認なしにファイル削除や仕様更新は行わない。

## 実行ステップ

### 1. テスト確認

```bash
go test ./...
```

全テストが通過することを確認する。失敗がある場合はレビューに進まない。

### 2. 変更内容の収集

```bash
git diff main HEAD
```

mainブランチからの全変更を収集する。

### 3. エキスパートエージェントに実装レビューを依頼

フィーチャードキュメントと変更差分を、4つのエキスパートエージェントにTask toolで**1つのメッセージで並列起動**してレビューを依頼する:

```
Task(subagent_type="spec-compliance", prompt="
MODE: review
REVIEW_TARGET:
{git diff main HEAD の全文}

FEATURE_DOC:
{docs/project/$ARGUMENTS/feature.md の内容}
")

Task(subagent_type="architecture", prompt="
MODE: review
REVIEW_TARGET:
{git diff main HEAD の全文}

FEATURE_DOC:
{docs/project/$ARGUMENTS/feature.md の内容}
")

Task(subagent_type="go-expert", prompt="
MODE: review
REVIEW_TARGET:
{git diff main HEAD の全文}

FEATURE_DOC:
{docs/project/$ARGUMENTS/feature.md の内容}
")

Task(subagent_type="test-quality", prompt="
MODE: review
REVIEW_TARGET:
{git diff main HEAD の全文}

FEATURE_DOC:
{docs/project/$ARGUMENTS/feature.md の内容}
")
```

### 4. レビュー結果に基づく修正とユーザーへの報告

#### 4a. レビュー結果の分析と自律修正

各エキスパートの指摘を重要度別に分類し、以下のルールで対応する:

| 重要度     | 対応                                                                         |
| ---------- | ---------------------------------------------------------------------------- |
| **MUST**   | 自分で修正する。修正方法が判断できない場合はユーザーへの相談事項に回す       |
| **SHOULD** | 修正が明確なものは自分で修正する。判断に迷うものはユーザーへの相談事項に回す |
| **NICE**   | 修正しない。報告のみ                                                         |
| **GOOD**   | 報告のみ                                                                     |

修正後、テストを再実行して全テストが通過することを確認する:

```bash
go test ./...
```

#### 4b. ユーザーへの報告

以下をユーザーに提示する:

```
## コードレビュー結果

### 修正済み
{自分で修正した指摘の一覧（重要度・対象・内容・修正内容）}

### 相談事項
{判断できなかった指摘の一覧（重要度・対象・内容・判断に迷った理由）}

### その他の指摘
{NICE・GOODの指摘の一覧}
```

- 「相談事項」がある場合はユーザーの判断を仰ぐ
- 「相談事項」がない場合でも、修正内容の確認のためユーザーの承認を待つ

**重要**: ユーザーの承認を待つ。

### 5. 承認後の処理

ユーザーが承認したら:

1. **仕様更新**: spec-complianceエージェントにupdateモードで仕様の最新化を依頼する

```
Task(subagent_type="spec-compliance", prompt="
MODE: update
CHANGES:
{git diff main HEAD の全文}

FEATURE_DOC:
{docs/project/$ARGUMENTS/feature.md の内容}
")
```

エージェントが返した更新内容を確認し、`docs/specification/` のファイルに反映する。

2. 完了メッセージを表示
   `docs/project/$ARGUMENTS/` ディレクトリを削除できることをユーザーに伝える。
