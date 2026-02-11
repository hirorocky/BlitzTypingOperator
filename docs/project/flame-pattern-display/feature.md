# flame-pattern-display

## 概要
チャレンジの表示パターンと文字生成を共通部品として切り出し、複数のチャレンジタイプから再利用可能にする。具体的には:

1. **炎形テンプレートシステム**: ASCIIアートテンプレートによるパターン表示を共通化
2. **難易度ベースの文字セット**: DifficultyRateに応じた段階的な文字種選択を導入
3. **チャレンジ設定の拡張**: modules.jsonの`challenge_type`を`challenge`オブジェクトに変更し、type + optionsで柔軟に設定可能にする

現在のsymbol_stormチャレンジを「shape」タイプとしてリファクタし、options.shapeで形状（flame等）を切り替え可能にする。合わせてchallengesディレクトリをチャレンジ毎のフォルダ構造に再編成する。

## 受け入れ基準

### マスタデータ構造（modules.json）
1. `challenge_type`(文字列)が`challenge`(オブジェクト)に置き換わり、`type`と`options`を持つ
2. shapeチャレンジは`{"type": "shape", "options": {"shape": "flame"}}`で定義される
3. standard/defenseチャレンジは`{"type": "standard"}`/`{"type": "defense"}`で定義される（optionsなし）
4. ローダーが新形式を正しくパースできる（パーステスト）
5. `ToDomainType()`がChallengeTypeID + ChallengeOptionsを正しく変換できる（変換テスト）
6. 全16エントリを新形式に一括変換する（旧形式の後方互換性は不要。マスタデータは全数管理）

### テンプレートシステム（commons/template.go）
7. `FormatAsPattern`関数が`challenges/commons/template.go`に配置されている
8. テンプレートは `#`（スロット）・空白・改行のみで構成される（View()のtextIdxカウントとの整合性保証）
9. テンプレートのスロット数が入力対象文字数以上であり、余剰スロットは空白に置換される
10. 入力順はテンプレート上の走査順（上から下、各行は左から右）で固定
11. 文字数が3以下の場合はテンプレートを使わず、スペース区切りの1行表示にフォールバックする

### 文字セットシステム（commons/charset.go）
12. 文字セット定義・生成ロジックが`challenges/commons/charset.go`に配置されている
13. DifficultyRateに応じて4段階の文字セットから文字が生成される:
   - 50-74: ホームポジション（`asdfghjkl;`）
   - 75-99: キーボード1行（上段 or 下段もランダムに混合）
   - 100-149: 英数字全体（a-z, 0-9）
   - 150-200: 記号込み（既存のbasicSymbols + shiftSymbols）
14. 文字数もDifficultyRateに連動する（現行symbol_stormのロジックを継承、ただし最低1文字）

### shapeチャレンジへの適用（shape/）
15. shapeチャレンジが共通文字セットと共通FormatAsPatternを使用する
16. テンプレート選択ロジックはshapeパッケージが担当する（`shape.selectTemplate(shapeName, charCount)`）
17. options.shapeの値（"flame"等）に応じてテンプレートセットが選択される
18. 炎形テンプレート定義が`challenges/shape/flame.go`に配置されている
19. 記号が炎形のASCIIアートパターンに配置されて表示される（文字数4以上の場合）

### ファイル構造
20. challenges配下がチャレンジ毎のフォルダ構造に再編成されている
21. 共通コードが`challenges/commons/`に集約されている
22. 各チャレンジフォルダの`main.go`にChallengeModel実装が配置されている

### 入力・評価ロジック（変更なし）
23. 正誤判定・Accuracy・SpeedFactor・WPM計算は既存と同一
24. AutoCorrect・MistakeTimeExtend・RetryOnTimeoutの動作は既存と同一
25. ESCキャンセル・タイムアウト失敗の動作は既存と同一

### View表示
26. 文字色はchallenge.goの共通色定数（ColorCorrect/ColorCursor/ColorUntyped）を使用
27. 残り時間バー・AutoCorrect残り表示は既存と同一
28. パターンの非空白文字を除去した結果がtextと一致する（表示と入力の整合性）

### テスト
29. 同一シードで同一テンプレート・同一patternが再現される
30. 既存の7テストが全て通過する（テスト互換性）
31. テンプレート選択の境界値テスト
32. スロット・文字整合性テスト: patternの非空白文字数がtextの文字数と一致する
33. 文字セット選択テスト: 各DifficultyRate範囲で期待する文字セットが使用される
34. パターン再現性テスト: 同一seedで2回生成したpatternが完全一致する
35. 文字数1-3のフォールバック表示テスト
36. ローダーのパーステスト: 新JSON形式がChallengeDataに正しく変換される
37. ローダーの変換テスト: ToDomainType()がChallengeTypeID + ChallengeOptionsを正しく出力する

## 設計

### マスタデータ変更

#### `internal/infra/masterdata/data/modules.json`
`challenge_type`(文字列)を`challenge`(オブジェクト)に変更:

```json
// Before
"challenge_type": "symbol_storm"

// After（shapeチャレンジ）
"challenge": {
  "type": "shape",
  "options": {
    "shape": "flame"
  }
}

// After（standardチャレンジ）
"challenge": {
  "type": "standard"
}

// After（defenseチャレンジ）
"challenge": {
  "type": "defense"
}
```

#### `internal/infra/masterdata/loader.go`
- `SkillDefinitionData.ChallengeType string` → `Challenge ChallengeData`
- 新規: `ChallengeData`構造体（Type string + Options map[string]string）
- `ToDomainType()`: ChallengeDataからChallengeTypeID + ChallengeOptionsに変換

### ドメイン変更

#### `internal/domain/challenge.go`
- `ChallengeTypeSymbolStorm` → `ChallengeTypeShape`に変更
- `ChallengeInput`に`ChallengeOptions map[string]string`フィールドを追加

#### `internal/domain/skill.go`
- `SkillType`に`ChallengeOptions map[string]string`フィールドを追加

### ディレクトリ構造変更

```
internal/tui/challenges/
├── challenge.go              # ChallengeModel interface, registry, factory, 共通色定数
├── challenge_test.go         # factory/registryテスト
├── commons/
│   ├── template.go           # FormatAsPattern（テンプレート配置の共通ロジック）
│   ├── template_test.go
│   ├── charset.go            # 文字セット定義・GenerateChars
│   └── charset_test.go
├── standard/
│   ├── main.go               # standardチャレンジ（standard.goから移動）
│   └── main_test.go
├── shape/
│   ├── main.go               # shapeチャレンジ（symbol_storm.goからリファクタ）
│   ├── main_test.go
│   └── flame.go              # 炎形テンプレート定義（小・中・大）
└── defense/
    ├── main.go               # defenseチャレンジ（defense.goから移動）
    └── main_test.go
```

#### パッケージ構成
- `challenges` — インターフェース定義、レジストリ、ファクトリ、blank import集約
- `challenges/commons` — FormatAsPattern（テンプレート配置）・GenerateChars（文字セット生成）
- `challenges/standard` — package standard、init()でRegister
- `challenges/shape` — package shape、init()でRegister、テンプレート選択・テンプレート定義を内包
- `challenges/defense` — package defense、init()でRegister

### 変更対象

#### `internal/tui/challenges/challenge.go`
- `constructor`の引数型はそのまま（`domain.ChallengeInput`がoptionsを持つ）
- sub-packageからのRegister呼び出しに対応（公開関数のまま）

#### `internal/tui/challenges/symbol_storm.go` → 削除
- `challenges/shape/main.go`に移動・リファクタ
- `generateSymbolPattern`を共通部品の呼び出しに簡素化
- `formatAsPattern`・記号定数を削除（commons/に移動）
- options.shapeに応じたテンプレートセット選択を追加

#### `internal/tui/challenges/standard.go` → `challenges/standard/main.go`に移動
#### `internal/tui/challenges/defense.go` → `challenges/defense/main.go`に移動

#### `internal/tui/screens/battle_logic.go`
- `startChallenge()`関数内で、`module.Type.ChallengeOptions`を`ChallengeInput.ChallengeOptions`に渡す:
  ```go
  input.ChallengeOptions = module.Type.ChallengeOptions
  ```

#### sub-packageのimport登録
- sub-packageのinit()が自動登録するため、root packageからのblank importが必要
- blank importは`challenges/challenge.go`の先頭に集約し、サブパッケージの登録を一元管理する:
  ```go
  import (
      _ "hirorocky/type-battle/internal/tui/challenges/standard"
      _ "hirorocky/type-battle/internal/tui/challenges/shape"
      _ "hirorocky/type-battle/internal/tui/challenges/defense"
  )
  ```

### 新規作成

#### `internal/tui/challenges/commons/template.go`
テンプレートにcharsを配置する共通部品。テンプレート選択はここに含めない（各shapeパッケージが担当）。

```go
package commons

// FormatAsPattern はテンプレートの#スロットにcharsを配置する。
// 文字数が3以下の場合はスペース区切りの1行表示にフォールバック。
// 余剰スロットは空白に置換される。
func FormatAsPattern(tmpl string, chars []rune) string { ... }
```

#### `internal/tui/challenges/commons/charset.go`
難易度ベースの文字セット定義と生成ロジック。

```go
package commons

// GenerateChars はDifficultyRateに応じた文字列を生成する。
func GenerateChars(diffRate int, rng *rand.Rand) []rune { ... }
```

#### `internal/tui/challenges/shape/flame.go`
炎形テンプレート定義。shape packageが自身のテンプレート選択を完全にカプセル化する。

```go
package shape

// 炎形テンプレート（小・中・大）
var flameSmall = strings.Trim(`...`, "\n")
var flameMedium = strings.Trim(`...`, "\n")
var flameLarge = strings.Trim(`...`, "\n")

// selectFlameTemplate は文字数に応じた炎形テンプレートを返す。
func selectFlameTemplate(charCount int) string { ... }
```

※ テンプレート形状・スロット数は実装時にTUIテスト（buffer mode）で視覚確認しながら調整。
※ 新形状（thunder.go等）追加時は、shape package内にファイルと`selectXxxTemplate()`を追加し、`selectTemplate()`から分岐するだけ。commons側の修正は不要。

### データフロー（変更後）

```
modules.json
  └─ "challenge": {"type": "shape", "options": {"shape": "flame"}}

loader.ToDomainType()
  └─ SkillType { ChallengeType: "shape", ChallengeOptions: {"shape": "flame"} }

battle_logic.go: startChallenge()
  └─ ChallengeInput { ..., ChallengeOptions: {"shape": "flame"} }

challenges.New("shape", input)
  └─ shape package: newShapeChallenge(input)
       ├─ shapeName := input.ChallengeOptions["shape"]  // "flame"
       ├─ commons.GenerateChars(diffRate, rng)           // 共通: 文字セット+文字数
       ├─ tmpl := selectTemplate(shapeName, len(chars))  // shape内: テンプレート選択
       ├─ commons.FormatAsPattern(tmpl, chars)           // 共通: テンプレート配置
       └─ (text, pattern) を返す
```

### View()への影響

**変更不要**。理由:
- テンプレートは `#`/空白/改行のみで構成（受け入れ基準7）
- `#` は文字配置後に記号に置換され、余剰 `#` は空白に置換
- フォールバック（1行表示）もスペース区切りで既存と同構造
- 現行View()の `ch == '\n' || ch == ' '` スキップ判定がそのまま正しく動作

### テスト追加方針

`commons/template_test.go`:
1. **スロット整合性テスト**: `FormatAsPattern`の出力から空白・改行を除去した文字列が入力charsと一致
2. **フォールバックテスト**: 文字数1-3でテンプレートなしの1行表示になる
3. **パターン構造テスト**: 文字数4以上で改行を含む（テンプレート適用を検証）

`commons/charset_test.go`:
5. **文字セット段階テスト**: DiffRate 50/75/100/150/200で期待する文字種のみが生成される
6. **文字数範囲テスト**: DiffRateに応じた文字数範囲で生成される
7. **再現性テスト**: 同一シードで同一結果

`shape/main_test.go`:
- 既存symbol_stormの7テストを移行（テスト互換性）
- パターン再現性テスト（同一seedで同一pattern）
- options.shapeによるテンプレート切り替えテスト
- テンプレート選択境界値テスト: 文字数ごとに正しいサイズの炎形テンプレートが選択される

`internal/infra/masterdata/loader_test.go`:
- ローダーパーステスト: 新JSON形式（challenge オブジェクト）が正しくChallengeDataにパースされる
- ローダー変換テスト: ToDomainType()がChallengeTypeID + ChallengeOptionsを正しく出力する

## メモ

### 設計判断
- `challenge_type`(文字列)から`challenge`(オブジェクト)への変更で、チャレンジ設定の拡張性を確保
- `symbol_storm` → `shape`へのリネームで、形状パラメータによる多様な表示パターンに対応
- challenges配下をチャレンジ毎フォルダ構造にし、アセット（テンプレート等）をチャレンジと同居させる
- **commonsの責務限定**: commonsは汎用ロジック（FormatAsPattern, GenerateChars）のみ提供。テンプレート選択は各shapeパッケージがカプセル化する。これにより新形状追加時にcommonsの修正が不要
- **ChallengeOptions**: `map[string]string`で汎用的に渡し、型安全なアクセスは各チャレンジpackage内で行う。domain層にオプション専用型を定義しない（チャレンジ毎にオプションの意味が異なるため）
- Go sub-packageパターン: 各チャレンジはinit()でRegister、blank importは`challenge.go`に集約
- テンプレート文字種を `#`/空白/改行のみに制限することで、View()の `textIdx` カウントとの整合性を保証
- 文字数3以下はテンプレートなし（1文字の炎形は視覚的に意味がない）
- マスタデータは全数管理のため、旧形式の後方互換性は不要。全エントリを一括変換する

### エッジケース
- 文字数1-3: テンプレートなし、スペース区切り1行表示
- 余剰スロット: 空白化（View()はスキップ）
- 文字数がどのテンプレートにも収まらない: flameLargeにフォールバック
- raw string literalの先頭/末尾改行: `strings.Trim` で除去
- options.shapeが未指定の場合: デフォルト形状（flame）にフォールバック

### 仕様同期
- `ChallengeTypeSymbolStorm` → `ChallengeTypeShape`のリネームに伴い、`/dev:complete`時に`docs/specification/typing.md`を同期更新する
- `ChallengeInput`への`ChallengeOptions`フィールド追加も同様にtyping.mdに反映する

### 将来の拡張
- 雷形・渦形等のテンプレート追加は、`shape/`に新ファイル（thunder.go等）+ `selectTemplate()`への分岐追加で完了（commonsの修正不要）
- 形状をスキルやエフェクトで切り替える場合は、modules.jsonのoptionsを変更するだけ
- 新しいチャレンジタイプ追加は、新フォルダ + init()でRegisterのパターンを踏襲
