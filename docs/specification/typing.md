# Typing System

## 概要

タイピングシステムはプレイヤーの入力評価を担当するドメインです。
チャレンジタイプごとに専用のタイピングミニゲームを定義し、共通のI/F（入力: 難易度・辞書、出力: スコア・スピード・成否）で統一します。

**実装**:
- ドメイン型: `/internal/domain/challenge.go`
- チャレンジ基盤: `/internal/tui/challenges/challenge.go`
- 各チャレンジ: `/internal/tui/challenges/standard/`, `shape/`, `defense/`
- 共通部品: `/internal/tui/challenges/commons/`（テンプレート配置・文字セット生成）

## 要件

### REQ-TYPING-1: チャレンジシステム基盤
**種別**: Ubiquitous

The typing system shall provide a challenge framework with:
- `ChallengeModel` インターフェース（型付きUpdate/View/Result）
- `challenges.New(typeID, input)` によるファクトリ生成
- 未知のtypeIDはスタンダードにフォールバック（warnログ出力）
- 新タイプ追加は `challenges/` にサブパッケージ + init()でのRegister呼び出し + blank importで完了

**受け入れ基準**:
1. ChallengeModelはInit/Update/View/Resultメソッドを持つ
2. Update戻り値はChallengeModel型（型アサーション不要）
3. Result()は完了時に非nilを返す
4. 未知のtypeIDでもパニックせずスタンダードにフォールバック

### REQ-TYPING-2: 難易度システム
**種別**: Ubiquitous

The typing system shall use DifficultyRate for difficulty scaling:
- DifficultyRate: 50-200（100=標準）。範囲外はClamp
- 文字数: DifficultyRateに応じた連続関数で決定（低→短い、高→長い）
- 制限時間: DifficultyRateに応じた連続関数で決定（低→長い、高→短い）
- EffectTableの`ColTypingDifficulty`（乗算集計）でバフ/デバフから動的制御

**受け入れ基準**:
1. DifficultyRate=100で標準的な文字数・制限時間
2. DifficultyRate=50で最も簡単（短い文字列、長い制限時間）
3. DifficultyRate=200で最も難しい（長い文字列、短い制限時間）
4. 範囲外の値はClampされる

### REQ-TYPING-3: チャレンジステータス
**種別**: Ubiquitous

The typing system shall use a 4-state challenge status:
- **Success** (=0): 全文字入力完了（ミスあり。タイムアウトと同フレーム競合時は成功優先）
- **Fail** (=1): 制限時間超過（タイムアウト）
- **Cancel** (=2): ESCキーによるキャンセル
- **Perfect** (=3): ミスなし全文字入力完了（Standard/Shapeのみ。Defenseでは発生しない）

**受け入れ基準**:
1. Success/Perfect時: 効果適用パイプライン実行
2. Fail時: 効果なし
3. Cancel時: 効果なし
4. ディフェンスタイプはFail状態にならない（制限時間なし）
5. ディフェンスタイプはPerfect状態にならない（常にSuccess）
6. IsSuccess()メソッドがSuccess/Perfect両方でtrueを返す

### REQ-TYPING-4: 入力評価
**種別**: State-Driven

While タイピングチャレンジ進行中, the typing system shall:
- 各入力文字の正誤判定
- 正解で進捗更新（緑文字表示）、誤りで誤字記録（カーソルが赤文字白背景に変化）
- ミスした位置は正解後も赤文字で表示（ミス履歴の可視化）
- リアルタイム進捗表示

**受け入れ基準**:
1. 正しい入力で次の文字へ進む
2. 誤った入力は誤字としてカウント
3. 進捗率を0.0〜1.0で計算

### REQ-TYPING-5: WPM計算
**種別**: Ubiquitous

The typing system shall calculate WPM as:
- WPM = (正しい文字数 / 完了時間(秒) x 60) / 5

**受け入れ基準**:
1. 標準的なWPM計算式に準拠
2. 完了時間0の場合は0を返す

### REQ-TYPING-6: スコア計算
**種別**: Ubiquitous

The typing system shall calculate score as:
- Standard/Shape: Score = int(正しい入力数 / 総入力数 × 100)（切り捨て）
- Defense: Score = int(最終防御率 × 100)（切り捨て）

**受け入れ基準**:
1. 0〜100の整数値
2. ミスなしの場合はScore=100
3. Cancel時はScore=0

## チャレンジタイプ

### スタンダード（standard）

物理攻撃スキル向けの基本タイピングチャレンジ。

**動作**:
- 英単語辞書からDifficultyRateに応じた文字数の単語を選択
- 1文字ずつ正誤判定。AutoCorrect対応
- 制限時間あり。DifficultyRateに応じた連続関数で決定
- MistakeTimeExtendSec（ミス時の時間延長）対応
- RetryOnTimeout（タイムアウト時の再挑戦）対応

**評価**:
- Score: int(正しい入力数 / 総入力数 × 100)
- WPM: パッシブスキル判定用
- ミスなし完了時: Status=ChallengePerfect

### シェイプ（shape）

魔法攻撃スキル向けのパターンチャレンジ。共通文字セットと形状テンプレートを使用してASCIIアートパターンを表示する。

**動作**:
- DifficultyRateに応じた4段階の文字セットから文字を生成（commons/charset.go）
  - 50-74: ホームポジション（`asdfghjkl;`）
  - 75-99: キーボード英字（`a-z`）
  - 100-149: 英数字（`a-z`, `0-9`）
  - 150-200: 記号込み（基本記号 + Shift+数字記号）
- 文字数もDifficultyRateに連動（低→4-6文字、高→10-16文字）
- `ChallengeOptions["shape"]`で形状を指定（デフォルト: "flame"）
- 形状テンプレート（`#`スロット・空白・改行で構成）に文字を配置（commons/template.go）
- 文字数3以下はテンプレートを使わずスペース区切りの1行表示にフォールバック
- テンプレートの走査順は上→下、各行は左→右で固定
- 制限時間・AutoCorrect・パッシブスキルの動作はスタンダードと同一

**テンプレート**:
- 炎形テンプレート（shape/flame.go）: 小（4-7文字）・中（8-12文字）・大（13+文字）
- 新形状はshapeパッケージ内にファイル追加で拡張可能

**評価**:
- スタンダードと共通（Score算出、ミスなし→ChallengePerfect）

### ディフェンス（defense）

防御スキル向けのリアルタイム防御チャレンジ。

**動作**:
- 英単語辞書から文字列を選択し、1文字ずつ入力（制限時間なし）
- 正しい文字を入力するたびに防御率（0%→100%）が上昇
- AutoCorrect発火時もミス無視として防御率が上昇（正しい入力と同等扱い）
- DifficultyRateに応じて1文字あたりの防御率上昇量が変動（低難易度=大きい）
- 全文字入力完了で次の単語が自動生成される（繰り返し）
- 敵の攻撃を受けた時点の防御率でダメージ軽減
  - 実効軽減率 = `damage_cut値 × DefenseRate`
- 敵攻撃後にチャレンジが自動終了（Status=Success, Score=int(攻撃時の防御率×100)）
- ESCキーでキャンセル可能（Status=Cancel、ダメージ軽減なし）

**特殊仕様**:
- 制限時間なし（Fail状態にならない）
- TimeExtend/RetryOnTimeout/MistakeTimeExtendは無効
- `DefenseProvider` インターフェースでバトル画面にリアルタイム防御率を公開
- `CompleteByAttack()` で敵攻撃時にバトル画面がチャレンジを自動終了

## 仕様

### ChallengeInput

**責務**: チャレンジへの共通入力パラメータ

**フィールド**:
- Difficulty: DifficultyRate（50-200）
- Words: 辞書（[]string、optional）
- TimeExtendSec: バフによる時間延長（ディフェンスでは無視）
- AutoCorrectCount: バフによるミス無視回数（全タイプで有効）
- MistakeTimeExtendSec: ミス時の時間延長秒数（ps_typo_recovery。ディフェンスでは無視）
- RetryOnTimeout: タイムアウト時の再挑戦許可（ps_second_chance。ディフェンスでは無視）
- RetryTimeLimitMultiplier: 再挑戦時の制限時間倍率
- ChallengeOptions: チャレンジ固有設定（map[string]string。shapeの場合: `{"shape": "flame"}`）

### ChallengeOutput

**責務**: チャレンジの出力結果

**フィールド**:
- Score: スコア（0-100の整数。Standard/Shape: int(accuracy×100)、Defense: int(defenseRate×100)）
- WPM: Words Per Minute
- CompletionTime: 完了時間
- Status: ChallengeStatus（Success/Fail/Cancel/Perfect）

### TypingResult

**責務**: チャレンジ結果からバトル層への入力データ

**実装**: `/internal/usecase/typing/typing.go`

**フィールド**:
- Completed: 完了フラグ
- WPM: Words Per Minute
- Score: スコア（0-100の整数）
- Timeout: タイムアウトフラグ
- IsPerfect: パーフェクトタイピングフラグ（Status==ChallengePerfect かつ非ディフェンス かつ`latent_effect`解放済みの場合 true）

### ChallengeModel

**責務**: チャレンジの共通インターフェース（tui層に配置）

**インターフェース**:
- `Init() tea.Cmd`
- `Update(tea.Msg) (ChallengeModel, tea.Cmd)` — 型付き戻り値
- `View() string`
- `Result() *ChallengeOutput` — 完了時に非nil

**設計**:
- tea.Modelを埋め込まず、型付きメソッドで定義
- 各チャレンジは独自のtea.Tickでタイマー管理（BattleTickMsgに依存しない）
- ファクトリは`map[ChallengeTypeID]constructor`で管理
- 各サブパッケージのinit()による自動登録（blank importはscreens/challenge_imports.goに集約）

### DefenseProvider

**責務**: ディフェンスタイプ専用のオプショナルインターフェース

**インターフェース**:
- `DefenseRate() float64` — 現在の防御率（0.0-1.0）
- `CompleteByAttack()` — 敵攻撃時にチャレンジを自動終了

**使用方法**:
- バトル画面が`activeChallenge.(DefenseProvider)`で型アサーション
- 既に終了済み（Cancel含む）の場合、CompleteByAttack()はno-op

## 関連ドメイン

- **Battle**: タイピング結果に基づくスキル効果計算、ディフェンス中の敵攻撃処理
- **Game Loop**: タイピング統計の記録
- **Collection**: WPM/スコアに基づく実績解除
- **Agent**: SkillType.ChallengeType/DifficultyRate/ChallengeOptionsの参照
