// Package challenges はタイピングチャレンジのモデルとファクトリを提供します。
package challenges

import (
	"log"

	"hirorocky/type-battle/internal/domain"

	tea "github.com/charmbracelet/bubbletea"
)

// ChallengeModel はタイピングチャレンジの共通インターフェースです。
// tea.Modelを埋め込まず、型付きのUpdate戻り値で型安全性を確保します。
type ChallengeModel interface {
	// Init はチャレンジの初期化コマンドを返します。
	Init() tea.Cmd

	// Update はメッセージを受け取りチャレンジの状態を更新します。
	Update(tea.Msg) (ChallengeModel, tea.Cmd)

	// View はチャレンジの描画結果を返します。
	View() string

	// Result はチャレンジ完了時に非nilの結果を返します。
	// 未完了の場合はnilを返します。
	Result() *domain.ChallengeOutput
}

// DefenseProvider はディフェンスタイプが実装するオプショナルインターフェースです。
// バトル画面がリアルタイムの防御率を参照するために使用します。
type DefenseProvider interface {
	// DefenseRate は現在の防御率（0.0-1.0）を返します。
	DefenseRate() float64

	// CompleteByAttack は敵攻撃時にチャレンジを自動終了させます。
	// 既に終了済み（Cancel含む）の場合はno-opです。
	CompleteByAttack()
}

// constructor はチャレンジモデルのコンストラクタ関数型です。
type constructor func(input domain.ChallengeInput) ChallengeModel

// registry はtypeID→コンストラクタのレジストリです。
var registry = map[domain.ChallengeTypeID]constructor{}

// Register はチャレンジタイプをレジストリに登録します。
func Register(typeID domain.ChallengeTypeID, ctor constructor) {
	registry[typeID] = ctor
}

// New は指定されたtypeIDに対応するChallengeModelを生成します。
// 未知のtypeIDの場合はwarnログを出力し、スタンダードにフォールバックします。
func New(typeID domain.ChallengeTypeID, input domain.ChallengeInput) ChallengeModel {
	if ctor, ok := registry[typeID]; ok {
		return ctor(input)
	}

	log.Printf("WARN: 未知のチャレンジタイプ '%s'、スタンダードにフォールバックします", typeID)

	if ctor, ok := registry[domain.ChallengeTypeStandard]; ok {
		return ctor(input)
	}

	// スタンダードも未登録の場合（起動時バグ）
	log.Printf("ERROR: スタンダードチャレンジが未登録です")
	return nil
}
