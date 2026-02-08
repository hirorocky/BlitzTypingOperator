---
description: 観点を分けてコードレビューを行う。
allowed-tools: Read, Bash, Glob, Grep
argument-hint: [<branch> | <commit1>..<commit2> | --staged]
---

# コードレビュー

## 概要

4つの観点別エキスパートエージェントを並列起動し、コードレビューを実施する。
各エージェントは内部でcodexにレビューを委譲する。

## 実行ステップ

### 1. レビュー対象の取得

`$ARGUMENTS` を解析してdiffを取得する:

- 引数なし → `git diff main...HEAD`
- `--staged` → `git diff --cached`
- `<commit1>..<commit2>` 形式 → `git diff <commit1>..<commit2>`
- その他 → `git diff $ARGUMENTS...HEAD`（ブランチ名として扱う）

```bash
# diffを取得
DIFF=$(git diff main...HEAD)  # 引数に応じて変更
```

### 2. フィーチャードキュメントの検出

`docs/project/` にディレクトリがあるか確認し、あれば `feature.md` を読み込む。
現在のブランチ名からフィーチャー名を推測し、`docs/project/{branch-name}/feature.md` を探す。

### 3. 4エージェント並列起動

以下の4つのTask toolを**1つのメッセージで並列起動**する。
各エージェントにはMODE、REVIEW_TARGET、FEATURE_DOCを渡す。

```
Task(subagent_type="spec-compliance", prompt="
MODE: review
REVIEW_TARGET:
{diff全文}

FEATURE_DOC:
{フィーチャードキュメント（あれば）}
")

Task(subagent_type="architecture", prompt="
MODE: review
REVIEW_TARGET:
{diff全文}

FEATURE_DOC:
{フィーチャードキュメント（あれば）}
")

Task(subagent_type="go-expert", prompt="
MODE: review
REVIEW_TARGET:
{diff全文}

FEATURE_DOC:
{フィーチャードキュメント（あれば）}
")

Task(subagent_type="test-quality", prompt="
MODE: review
REVIEW_TARGET:
{diff全文}

FEATURE_DOC:
{フィーチャードキュメント（あれば）}
")
```

### 4. 結果の統合と提示

各エージェントの結果を以下の形式で統合して提示する:

```
## コードレビュー結果

### 仕様準拠（spec-compliance）
{spec-complianceの結果}

### プロジェクト構造（architecture）
{architectureの結果}

### Go言語品質（go-expert）
{go-expertの結果}

### テスト品質（test-quality）
{test-qualityの結果}

---
**サマリー**
- MUST: {件数}件
- SHOULD: {件数}件
- NICE: {件数}件
- GOOD: {件数}件
```
