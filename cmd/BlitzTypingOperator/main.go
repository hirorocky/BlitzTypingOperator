// Package main は TypeBattleターミナルタイピングバトルゲームのエントリーポイントです。
//
// ゲームはBubbletea TUIフレームワークを使用し、Elm Architectureパターンに基づいて
// RootModelを中心としたイベント駆動型アーキテクチャで実装されています。
//
// 終了操作（qキーまたはCtrl+C）を行うと、tea.WithAltScreen()により
// ターミナルの状態が自動的に復元されます。
//
// コマンドライン引数:
//
//	-data <path>  外部データディレクトリのパス（省略時は埋め込みデータを使用）
package main

import (
	"flag"
	"fmt"
	"os"

	"hirorocky/type-battle/internal/app"
	"hirorocky/type-battle/internal/infra/masterdata"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// コマンドライン引数を解析
	dataDir := flag.String("data", "", "外部データディレクトリのパス（省略時は埋め込みデータを使用）")
	flag.Parse()

	// RootModelを作成 - ゲーム全体の状態管理とシーンルーティングを担当
	// 外部データディレクトリが指定されていない場合は埋め込みデータを使用
	model := app.NewRootModel(*dataDir, masterdata.EmbeddedData)

	// Bubbleteaプログラムを作成
	// tea.WithAltScreen(): 代替スクリーンバッファを使用し、
	// 終了時に元のターミナル状態を自動的に復元する
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
