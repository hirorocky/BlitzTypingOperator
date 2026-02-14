package domain

// RankReward はランクアップ報酬を表すVOです。
// 特定ランクに到達した際に獲得できる報酬アイテムリストを保持します。
type RankReward struct {
	// Rank は対象ランクです。
	Rank int

	// Items は報酬アイテムリストです。
	Items []RankRewardItem
}

// RankRewardItem はランクアップ報酬の個別アイテムです。
type RankRewardItem struct {
	// Category はアイテムカテゴリです（"core", "skill", "chain_effect"）。
	Category string

	// TypeID はアイテムのTypeIDです。
	TypeID string
}
