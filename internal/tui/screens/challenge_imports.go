package screens

// チャレンジサブパッケージの登録
// 各パッケージのinit()でchallenges.Register()が呼ばれ、チャレンジタイプが使用可能になる。
import (
	_ "hirorocky/type-battle/internal/tui/challenges/defense"
	_ "hirorocky/type-battle/internal/tui/challenges/shape"
	_ "hirorocky/type-battle/internal/tui/challenges/standard"
)
