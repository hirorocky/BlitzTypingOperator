package masterdata

import "embed"

// EmbeddedData は data/ ディレクトリ配下の全ファイルを埋め込みます。
// デフォルトデータの埋め込みと外部ファイルによる上書きをサポートします。
//
//go:embed data/*
var EmbeddedData embed.FS
