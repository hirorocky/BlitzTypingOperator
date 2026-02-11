---
description: フィーチャー開発を開始し、codexと共同で計画・設計する
allowed-tools: Read, Write, Glob, Grep, Bash
argument-hint: <feature-description>
---

# フィーチャー開始

## 概要
観点別エキスパートエージェントを活用しながら、受け入れ基準と設計の全体像を含むフィーチャードキュメントを作成する。
ユーザーの承認なしに実装には進まない。

## 実行ステップ

### 1. フィーチャー名の生成
`$ARGUMENTS` からkebab-caseのフィーチャー名を生成する。
既存の `docs/project/` を確認し、名前の重複を避ける。

### 2. フィーチャーディレクトリの作成
`docs/project/{feature-name}/` ディレクトリを作成する。

### 3. 仕様の確認
`docs/specification/` から関連するドメイン仕様を読み、フィーチャーに必要な背景知識を得る。

### 4. エキスパートエージェントに設計提案を依頼

4つのエキスパートエージェントをTask toolで**1つのメッセージで並列起動**し、各観点から設計提案を得る:

```
Task(subagent_type="spec-compliance", prompt="
MODE: propose
FEATURE_DESCRIPTION:
$ARGUMENTS

{読み込んだ関連仕様の内容}
")

Task(subagent_type="architecture", prompt="
MODE: propose
FEATURE_DESCRIPTION:
$ARGUMENTS

{読み込んだ関連仕様の内容}
")

Task(subagent_type="go-expert", prompt="
MODE: propose
FEATURE_DESCRIPTION:
$ARGUMENTS

{読み込んだ関連仕様の内容}
")

Task(subagent_type="test-quality", prompt="
MODE: propose
FEATURE_DESCRIPTION:
$ARGUMENTS

{読み込んだ関連仕様の内容}
")
```

### 5. フィーチャードキュメントの作成
エキスパートエージェントの提案と自分の分析を総合して、`docs/project/{feature-name}/feature.md` を作成する:

```markdown
# {フィーチャー名}

## 概要
{何を作るか、なぜ作るか。2-3文。}

## 受け入れ基準
1. {テスト可能な基準}
2. ...

## 設計
### 変更対象
- {変更するファイル/パッケージと変更内容}

### 新規作成
- {新規作成するファイル/型/関数}

### データフロー
{処理の流れの概要}

## メモ
{設計判断、リスク、エッジケースなど。任意。}
```

- 受け入れ基準は具体的かつテスト可能にする。各基準は Go テストまたは TUI テストで検証できるように書く。
- 各エキスパートの提案を鵜呑みにせず、自分の分析と比較検討した上で記述する。

### 6. エキスパートエージェントに設計レビューを依頼

作成したフィーチャードキュメントの全文をレビュー対象として、4つのエキスパートエージェントをTask toolで**1つのメッセージで並列起動**する:

```
Task(subagent_type="spec-compliance", prompt="
MODE: review
REVIEW_TARGET:
{フィーチャードキュメントの全文}
")

Task(subagent_type="architecture", prompt="
MODE: review
REVIEW_TARGET:
{フィーチャードキュメントの全文}
")

Task(subagent_type="go-expert", prompt="
MODE: review
REVIEW_TARGET:
{フィーチャードキュメントの全文}
")

Task(subagent_type="test-quality", prompt="
MODE: review
REVIEW_TARGET:
{フィーチャードキュメントの全文}
")
```

### 7. 結果の提示
以下をユーザーに提示する:
- 作成したフィーチャードキュメントの内容
- 各エキスパートのレビュー結果
- 次のステップの案内:
  - フィードバックがある場合: `/dev:refine {feature-name}`
  - 承認して実装に進む場合: `/impl:plan {feature-name}`

**重要**: ユーザーの承認を待つ。承認なしに実装に進まない。

## 出力形式
- フィーチャー名とディレクトリパスを明示
- 各エキスパートのフィードバックを観点ごとに引用形式で提示
- 次のコマンドをコードブロックで提示
