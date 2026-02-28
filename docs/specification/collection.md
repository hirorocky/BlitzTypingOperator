# Collection System

## 概要

コレクションシステムは図鑑と実績を管理するドメインです。
コア/スキル/チェイン効果/敵の図鑑登録、実績の達成判定と通知を担当します。

**実装**: `/internal/achievement/achievement.go`

## 要件

### REQ-COLLECTION-1: タイピング実績
**種別**: Event-Driven

When タイピング結果が記録される, the collection system shall check:
- WPMマイルストーン: 50, 80, 100, 120
- スコア100達成（パーフェクト）

**受け入れ基準**:
1. 各閾値到達時に実績解除
2. 一度解除した実績は永続
3. 達成時に通知を返却

### REQ-COLLECTION-2: バトル実績
**種別**: Event-Driven

When バトル結果が記録される, the collection system shall check:
- 敵撃破数: 10, 50, 100, 500体
- 最高レベル: 10, 25, 50, 100
- ノーダメージクリア

**受け入れ基準**:
1. 累積統計に基づいて判定
2. ノーダメージは単一バトル単位
3. 達成時に通知を返却

### REQ-COLLECTION-3: 実績通知
**種別**: Event-Driven

When 実績が解除される, the collection system shall:
- 実績ID、名前、説明を含む通知を生成
- UI層への通知配信

**受け入れ基準**:
1. 既に解除済みの実績は通知しない
2. 複数同時解除時は複数通知

### REQ-COLLECTION-4: 達成率表示
**種別**: Ubiquitous

The collection system shall calculate completion rate:
- 達成率 = 解除済み実績数 / 全実績数

**受け入れ基準**:
1. 0.0〜1.0の範囲で表現
2. 画面上ではパーセント表示

### REQ-COLLECTION-5: セーブ/ロード
**種別**: Event-Driven

When ゲームをセーブ/ロードする, the collection system shall:
- 解除済み実績IDリストを永続化
- ロード時に解除状態を復元

**受け入れ基準**:
1. 進捗情報も保存（将来拡張用）
2. 存在しないIDは無視

### REQ-COLLECTION-6: 図鑑システム
**種別**: Ubiquitous

The collection system shall track collected items:
- コア図鑑: 入手済みコアの一覧
- スキル図鑑: 入手済みスキルの一覧
- チェイン効果図鑑: 入手済みチェイン効果の一覧
- 敵図鑑: エンカウント済み敵の一覧

**受け入れ基準**:
1. 入手・エンカウント時に図鑑登録
2. 図鑑画面で一覧表示・詳細確認

## 仕様

### AchievementDefinition

**責務**: 実績の定義（メタデータ）

**フィールド**:
- ID: 一意識別子
- Name: 表示名
- Description: 説明文
- Category: カテゴリ（typing/battle）

### AchievementManager

**責務**: 実績の管理と達成判定

**インターフェース**:
- 入力: タイピング結果、バトル結果
- 出力: AchievementNotification[]

**ルール**:
1. tryUnlockで解除可否を判定
2. 既に解除済みならnilを返す
3. 解除時に通知を生成

### 実績一覧

**タイピング実績**:
| ID | 名前 | 条件 |
|----|------|------|
| wpm_50 | タイピング見習い | WPM 50達成 |
| wpm_80 | タイピング上手 | WPM 80達成 |
| wpm_100 | タイピングマスター | WPM 100達成 |
| wpm_120 | タイピングレジェンド | WPM 120達成 |
| perfect_accuracy | 完璧主義者 | スコア100達成 |

**バトル実績**:
| ID | 名前 | 条件 |
|----|------|------|
| defeat_10 | 新米ハンター | 敵10体撃破 |
| defeat_50 | 歴戦の戦士 | 敵50体撃破 |
| defeat_100 | 百戦錬磨 | 敵100体撃破 |
| defeat_500 | 伝説の勇者 | 敵500体撃破 |
| level_10 | 探索者 | レベル10到達 |
| level_25 | 冒険者 | レベル25到達 |
| level_50 | 英雄 | レベル50到達 |
| level_100 | 覇者 | レベル100到達 |
| no_damage | 無傷の勝利 | ノーダメージ勝利 |

### 図鑑システム

**責務**: コア/スキル/チェイン効果/敵の収集状況を追跡

**ルール**:
1. 入手・エンカウント済みアイテムはGameStateで管理
2. 所持コア/スキル/チェイン効果はインベントリから参照
3. 図鑑画面で一覧表示
4. コレクション完成度を達成率で表示

## 関連ドメイン

- **Typing**: WPM/スコアに基づく実績トリガー
- **Battle**: 敵撃破/レベル/ノーダメージに基づく実績トリガー
- **Game Loop**: 実績状態の永続化、エンカウント敵の追跡
- **Agent**: インベントリからのコア/スキル/チェイン効果取得
