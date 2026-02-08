# flame-pattern-display

## 概要
symbol_stormチャレンジの表示パターンを、現在の単純なグリッド配置から、実際の炎形ASCIIアートテンプレートに変更する。仕様（typing.md）では「記号や英数字を視覚的なパターン（炎形、雷形等）に配置」と記載されているが、現実装は4-5列のグリッド配置のみ。炎形テンプレートを導入し、視覚的なゲーム体験を向上させる。

## 受け入れ基準

### テンプレート表示
1. symbol_stormチャレンジで、記号が炎形のASCIIアートパターンに配置されて表示される
2. テンプレートは文字数レンジに応じて3種類（小: 4-6文字、中: 7-10文字、大: 11-16文字）から選択される
3. テンプレートのスロット数が入力対象文字数以上であり、先頭から必要数だけ文字が配置される。余剰スロットは空白に置換される（View()変更不要）
4. 入力順はテンプレート上の走査順（上から下、各行は左から右）で固定
5. テンプレートは `#`（スロット）・空白・改行のみで構成される。それ以外の文字は使用禁止（View()のtextIdxカウントとの整合性保証）

### 入力・評価ロジック（変更なし）
6. 入力対象のtext（記号列）は既存のgenerateSymbolPattern関数で生成され、変更なし
7. 正誤判定・Accuracy・SpeedFactor・WPM計算は既存と同一
8. AutoCorrect・MistakeTimeExtend・RetryOnTimeoutの動作は既存と同一
9. ESCキャンセル・タイムアウト失敗の動作は既存と同一

### View表示
10. 文字色はchallenge.goの共通色定数（ColorCorrect/ColorCursor/ColorUntyped）を使用（既存と同一）
11. 残り時間バー・AutoCorrect残り表示は既存と同一
12. パターンの非空白文字を除去した結果がtextと一致する（表示と入力の整合性）

### テスト
13. 同一シードで同一テンプレート・同一patternが再現される
14. 既存の7テストが全て通過する（テスト互換性）
15. テンプレート選択の境界値テスト: 文字数4,6,7,10,11,16でそれぞれ正しいサイズのテンプレートが選択される
16. スロット・文字整合性テスト: patternの非空白文字数がtextの文字数と一致する
17. パターン再現性テスト: 同一seedで2回生成したpatternが完全一致する

## 設計

### 変更対象

#### `internal/tui/challenges/symbol_storm.go`
- `formatAsPattern(chars []rune) string` を炎形テンプレートベースに置換
- テンプレートデータ（Go変数）とテンプレート選択関数を同ファイルに追加
- `symbolStormChallenge` 構造体のフィールド変更なし（text/pattern分離は維持）

### 新規作成
なし（symbol_storm.go内に追加）

### テンプレート設計

テンプレートは `#`（スロット）・空白・改行のみで構成されるASCIIアート文字列。`#` の走査順（上→下、左→右）が入力順になる。

raw string literalで定義し、先頭・末尾の余分な改行は `strings.Trim(..., "\n")` で除去する。

```go
// 炎形テンプレート（小: スロット数は実装時にデザイン調整）
var flameSmall = strings.Trim(`
   #
  # #
 #   #
  # #
`, "\n")

// 炎形テンプレート（中: スロット数は実装時にデザイン調整）
var flameMedium = strings.Trim(`
     #
    # #
   # # #
  #     #
   # # #
`, "\n")

// 炎形テンプレート（大: スロット数は実装時にデザイン調整）
var flameLarge = strings.Trim(`
       #
      # #
     # # #
    # # # #
   #       #
    # # # #
     # # #
`, "\n")
```

※ テンプレート形状・スロット数は実装時にTUIテスト（buffer mode）で視覚確認しながら調整。上記は概念例。

### テンプレート選択ロジック

```go
func selectFlameTemplate(charCount int) string {
    switch {
    case charCount <= smallMaxSlots:
        return flameSmall
    case charCount <= mediumMaxSlots:
        return flameMedium
    default:
        return flameLarge
    }
}
```

境界値は名前付き定数で定義。文字数が0や想定外の場合はflameLargeにフォールバック。

### formatAsPattern の置換

```go
func formatAsPattern(chars []rune) string {
    template := selectFlameTemplate(len(chars))

    var b strings.Builder
    b.Grow(len(template))
    charIdx := 0
    for _, r := range template {
        if r == '#' {
            if charIdx < len(chars) {
                b.WriteRune(chars[charIdx])
                charIdx++
            } else {
                // 余剰スロットは空白化（View()変更不要）
                b.WriteRune(' ')
            }
        } else {
            b.WriteRune(r)
        }
    }
    return b.String()
}
```

### データフロー（変更部分のみ）

```
generateSymbolPattern(diffRate, rng)
  ├─ 文字数・記号生成（変更なし）
  ├─ formatAsPattern(chars) ← ここだけ変更
  │    ├─ 文字数からテンプレート選択
  │    ├─ テンプレートの#スロットにcharsを順に配置
  │    └─ 余剰スロットは空白化
  └─ (text, pattern) を返す
```

### View()への影響

**変更不要**。理由:
- テンプレートは `#`/空白/改行のみで構成（受け入れ基準5）
- `#` は文字配置後に記号に置換され、余剰 `#` は空白に置換
- 結果として pattern には空白・改行・記号文字のみが含まれる
- 現行View()の `ch == '\n' || ch == ' '` スキップ判定がそのまま正しく動作

### テスト追加方針

`symbol_storm_test.go` に以下を追加:

1. **テンプレート選択境界値テスト**: `selectFlameTemplate` を文字数4,6,7,10,11,16で呼び、期待テンプレートと一致
2. **スロット整合性テスト**: `formatAsPattern` の出力から空白・改行を除去した文字列が入力charsと一致
3. **パターン再現性テスト**: 同一シードで `generateSymbolPattern` を2回呼び、pattern文字列が完全一致
4. **パターン構造テスト**: `formatAsPattern` の出力が改行を含む（グリッド形式からの脱却を検証）

既存テストの修正:
- `TestSymbolStorm_View出力`: 「空でない」に加えて改行を含むことを検証（炎形は複数行）

## メモ

### 設計判断
- テンプレートはGo変数として`symbol_storm.go`内に定義。embed.FSや別パッケージは過剰設計（テンプレートは3個のみ）
- Catalogインターフェースやメタデータファイルは不要。テンプレート選択は単純な文字数範囲マッチング
- `patterns/`サブパッケージは作らない。1ファイルで完結する変更量
- テンプレート形状の正確なデザインは実装時にTUIテスト（buffer mode）で視覚確認しながら調整
- テンプレート文字種を `#`/空白/改行のみに制限することで、View()の `textIdx` カウントとの整合性を保証

### エッジケース
- 余剰スロット: 空白化（View()はスキップ）
- 文字数がどのテンプレートにも収まらない: flameLargeにフォールバック
- raw string literalの先頭/末尾改行: `strings.Trim` で除去

### 将来の拡張
- 雷形・渦形等のテンプレート追加は、Go変数を追加しselectFlameTemplateの形状選択を拡張するだけ
- 形状をスキルやエフェクトで切り替える場合は、ChallengeInputにshapeIDを追加
