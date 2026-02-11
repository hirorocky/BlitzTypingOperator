# 実装計画: flame-pattern-display

## 受け入れテスト

| # | 基準 | テスト種別 | テストファイル/方法 | 状態 |
|---|------|-----------|-------------------|------|
| 1-5 | マスタデータ構造（パース、変換） | Go test | `masterdata/loader_challenge_test.go` | 作成済（build tag: flame_pattern→タスク1で除去） |
| 6 | 全16エントリ新形式変換 | Go test | 既存`loader_test.go`のLoadAllが通過 | タスク2後に確認 |
| 7-11 | テンプレートシステム | Go test | `commons/template_test.go` | 作成済（RED） |
| 12-14 | 文字セットシステム | Go test | `commons/charset_test.go` | 作成済（RED） |
| 15-19 | shapeチャレンジ適用 | Go test | `shape/main_test.go` | タスク5で作成 |
| 20-22 | ファイル構造 | Go compile | `go build ./...`成功 = パッケージ構造OK | タスク5で確認 |
| 23-25 | 入力・評価ロジック（変更なし） | Go test | 既存テスト移行後の通過で保証 | 既存 |
| 26-28 | View表示 | TUI test + Go test | セッション内 + `template_test.go`(基準28) | タスク8 |
| 29, 34 | 再現性 | Go test | `charset_test.go` + `shape/main_test.go` | 作成済（RED）/ タスク5 |
| 30 | 既存テスト通過 | Go test | 移行後の全テスト | タスク完了後 |
| 31 | テンプレート選択境界値 | Go test | `shape/main_test.go` | タスク5で作成 |
| 32, 35 | スロット整合性・フォールバック | Go test | `commons/template_test.go` | 作成済（RED） |
| 33 | 文字セット選択 | Go test | `commons/charset_test.go` | 作成済（RED） |
| 36-37 | ローダーパース・変換 | Go test | `masterdata/loader_challenge_test.go` | 作成済（build tag） |

## 既存テスト移行チェックリスト（基準30）

移行対象の既存テスト一覧。移行後も同等のテスト内容を保証する。

### symbol_storm_test.go → shape/main_test.go
| テスト名 | 検証内容 | 移行後の差分 |
|---------|---------|-------------|
| TestSymbolStorm_正しい入力で進捗 | 全文字入力→Success, Accuracy=1.0 | 構造体名変更のみ |
| TestSymbolStorm_ESCでキャンセル | ESC→Cancel | 構造体名変更のみ |
| TestSymbolStorm_タイムアウトで失敗 | タイムアウト→Fail | 構造体名変更のみ |
| TestSymbolStorm_テキストは記号を含む | テキストに記号が含まれる | 文字セットが4段階に変更。Rate>=150で記号テスト |
| TestSymbolStorm_DifficultyRateで文字数が変動 | 難易度で文字数変動 | 文字数範囲は同一ロジック継承 |
| TestSymbolStorm_AutoCorrect | AutoCorrectで次に進む | 構造体名変更のみ |
| TestSymbolStorm_View出力 | Viewが空でない | 構造体名変更のみ |

### standard_test.go → standard/main_test.go
既存テストをパッケージ移動。内容変更なし。

### defense_test.go → defense/main_test.go
既存テストをパッケージ移動。内容変更なし。

### challenge_test.go（rootパッケージ）
| テスト名 | 移行後の差分 |
|---------|-------------|
| TestNew_全3タイプが生成可能 | `ChallengeTypeSymbolStorm`→`ChallengeTypeShape` |
| TestChallengeModel_インターフェース準拠 | 同上 |
| TestDefenseProvider_ディフェンスタイプのみ実装 | 同上 |

## 実装タスク

### タスク1: マスタデータ構造の変更（domain + infra）
- **対象**: `internal/domain/challenge.go`, `internal/domain/skill.go`, `internal/infra/masterdata/loader.go`
- **内容**:
  1. `domain/challenge.go`: `ChallengeTypeSymbolStorm` → `ChallengeTypeShape`にリネーム
  2. `domain/challenge.go`: `ChallengeInput`に`ChallengeOptions map[string]string`フィールド追加
  3. `domain/skill.go`: `SkillType`に`ChallengeOptions map[string]string`フィールド追加
  4. `loader.go`: `ChallengeData`構造体（Type string + Options map[string]string）を新規定義
  5. `loader.go`: `SkillDefinitionData`の`ChallengeType string`を`Challenge ChallengeData`に変更
  6. `loader.go`: `ToDomainType()`を更新（ChallengeData→ChallengeTypeID + ChallengeOptions）
  7. `loader_challenge_test.go`のbuild tag除去
  8. 既存loader_test.goのmodules.json fixture更新（新形式に合わせる）
- **関連テスト**: 基準1-6, 36-37
- **状態**: 完了

### タスク2: modules.jsonの全エントリ変換
- **対象**: `internal/infra/masterdata/data/modules.json`
- **内容**:
  1. 全エントリの`"challenge_type": "xxx"`を`"challenge": {"type": "xxx"}`に変換
  2. symbol_storm→`{"type": "shape", "options": {"shape": "flame"}}`
  3. standard→`{"type": "standard"}`
  4. defense→`{"type": "defense"}`
- **関連テスト**: 基準6
- **依存**: タスク1完了後
- **状態**: 完了

### タスク3: commons/charset.go の実装
- **対象**: `internal/tui/challenges/commons/charset.go`
- **内容**:
  1. 4段階の文字セット定義（ホームポジション / キーボード1行 / 英数字 / 記号込み）
  2. `GenerateChars(diffRate, rng)`: DifficultyRateに応じた文字セット選択と文字数決定
  3. 文字数ロジックは現行symbol_stormから移植（最低1文字保証）
- **関連テスト**: 基準12-14, 29, 33
- **状態**: 完了

### タスク4: commons/template.go の実装
- **対象**: `internal/tui/challenges/commons/template.go`
- **内容**:
  1. `FormatAsPattern(tmpl, chars)`: テンプレートの#スロットにcharsを配置
  2. 文字数3以下のフォールバック（スペース区切り1行表示）
  3. 余剰スロットの空白置換
  4. 走査順は上→下、左→右で固定
- **関連テスト**: 基準7-11, 28, 32, 35
- **状態**: 完了

### タスク5: ディレクトリ再編成 + shape/flame.goの作成
- **対象**: `internal/tui/challenges/` 配下
- **内容**:
  1. `challenges/standard/main.go`の作成（standard.goから移動、パッケージ名変更）
  2. `challenges/standard/main_test.go`の作成（standard_test.goから移動）
  3. `challenges/defense/main.go`の作成（defense.goから移動、パッケージ名変更）
  4. `challenges/defense/main_test.go`の作成（defense_test.goから移動）
  5. `challenges/shape/main.go`の作成（symbol_storm.goからリファクタ、commons使用）
  6. `challenges/shape/flame.go`の作成（炎形テンプレート定義：小・中・大）
  7. `challenges/shape/main_test.go`の作成（移行チェックリスト参照 + 新規テスト）
  8. 各サブパッケージのinit()でRegister
  9. `challenge.go`にblank import追加
  10. 旧ファイル削除（standard.go, defense.go, symbol_storm.go, *_test.go）
- **shape/main_test.goに追加する新規テスト**:
  - Options→テンプレート選択の統合テスト（options.shape="flame"で炎形テンプレート使用）
  - テンプレート選択境界値（文字数ごとに正しいサイズの炎形テンプレートが選択される）
  - パターン再現性（同一seedで同一pattern）
  - パターン整合性（patternの非空白文字数 == textの文字数）
- **関連テスト**: 基準15-22, 30, 31, 34
- **依存**: タスク1, 3, 4完了後
- **状態**: 完了

### タスク6: battle_logic.goの更新
- **対象**: `internal/tui/screens/battle_logic.go`
- **内容**:
  1. `startChallenge()`で`module.Type.ChallengeOptions`を`ChallengeInput.ChallengeOptions`に渡す
- **関連テスト**: 基準15（shape/main_test.goの統合テストで間接検証）
- **依存**: タスク1完了後
- **状態**: 完了

### タスク7: ChallengeTypeSymbolStormの参照箇所更新
- **対象**: コードベース全体
- **内容**:
  1. `ChallengeTypeSymbolStorm`→`ChallengeTypeShape`のリネームに伴う全参照箇所の更新
  2. challenge_test.goの更新（移行チェックリスト参照）
  3. typing仕様（docs/specification/typing.md）の文字列参照更新
- **関連テスト**: 基準30
- **依存**: タスク1完了後
- **状態**: 完了

### タスク8: 全テスト通過 + TUI動作確認
- **対象**: ゲーム全体
- **内容**:
  1. `go test ./...`で全テスト通過確認
  2. `go build ./...`でコンパイル確認
  3. TUIテスト: ゲーム起動→バトル→shapeチャレンジの炎形パターン表示確認
  4. TUIテスト: standardチャレンジ・defenseチャレンジの動作確認
- **関連テスト**: 基準20-22, 26-28, 30
- **依存**: 全タスク完了後
- **状態**: 完了

## TUIテスト手順

### 基準26-28: View表示の確認
```
launch_tui(command="go run ./cmd/BlitzTypingOperator", mode="buffer", dimensions="160x45")
# ホーム画面からバトルセレクトへ
send_keys("\r", delay=2)
# 敵を選択してバトル開始
send_keys("\r", delay=1)
# バトル画面でスキル選択（魔法攻撃スキルを選択）
# → shapeチャレンジの炎形パターン表示を確認
capture_screen()
# ASCIIアートパターンが表示されていることを確認
# 文字色（ANSI）、残り時間バー、パターン構造を視覚確認
```

## 進捗ログ

### タスク1: マスタデータ構造の変更（domain + infra）
- `ChallengeTypeSymbolStorm` → `ChallengeTypeShape` にリネーム
- `ChallengeInput` に `ChallengeOptions map[string]string` フィールド追加
- `SkillType` に `ChallengeOptions` フィールド追加
- `ChallengeData` 構造体を新規定義、`SkillDefinitionData.ChallengeType` を `Challenge ChallengeData` に変更
- `ToDomainType()` を更新
- `loader_challenge_test.go` のbuild tag除去
- `masterdata_converter.go` の参照修正（`t.ChallengeType` → `t.Challenge.Type`）

### タスク2: modules.jsonの全エントリ変換
- 全18エントリを新形式に変換
- `"challenge_type": "standard"` → `"challenge": { "type": "standard" }`
- `"challenge_type": "symbol_storm"` → `"challenge": { "type": "shape", "options": { "shape": "flame" } }`
- `"challenge_type": "defense"` → `"challenge": { "type": "defense" }`

### タスク3: commons/charset.go の実装
- 4段階文字セット: homePosition(50-74), keyboardRow(75-99), alphanumeric(100-149), symbol(150+)
- `GenerateChars(diffRate, rng)` を実装
- 既存symbol_stormの文字数計算ロジックを移植

### タスク4: commons/template.go の実装
- `FormatAsPattern(tmpl, chars)` を実装
- #スロット置換、フォールバック表示（≤3文字）、余剰スロット空白化

### タスク5: ディレクトリ再編成 + shape/flame.go作成
- `standard/main.go`, `defense/main.go`, `shape/main.go` をサブパッケージとして作成
- `shape/flame.go`: 3サイズの炎形テンプレート（small: 4-7文字, medium: 8-12文字, large: 13+文字）
- 各サブパッケージのinit()でRegister、blank importは`screens/challenge_imports.go`に配置
- `challenge_test.go` を外部テストパッケージ（`package challenges_test`）に変更しインポートサイクル回避
- 旧ファイル（standard.go, defense.go, symbol_storm.go, *_test.go）を削除
- 新規テスト14件追加（shape/main_test.go）

### タスク6: battle_logic.goの更新
- `ChallengeOptions: module.Type.ChallengeOptions` を `ChallengeInput` に追加

### タスク7: ChallengeTypeSymbolStormの参照箇所更新
- `domain/challenge_test.go` の参照を更新

### タスク8: 全テスト通過 + TUI動作確認
- `go build ./...`: 全パッケージコンパイル成功
- `go test ./...`: 全テスト通過（20+ パッケージ）
- TUI動作確認: ゲーム起動→バトル選択→バトル開始→standardチャレンジ正常動作
- shapeチャレンジはセーブデータにshapeモジュールがないためTUIでの直接確認は不可（ユニットテストで網羅済）
