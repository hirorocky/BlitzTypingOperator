# Requirements Document

## Introduction

本仕様は、コアレベルシステムの廃止と敵ランクシステムの導入に関する要件を定義します。
コアからレベル概念を完全に削除し、エージェントのパラメータ計算方式を変更するとともに、
敵の進行システムをランクベースに再設計し、最大HPの成長をプレイヤーの戦績に連動させます。

## Requirements

### Requirement 1: コアレベルの廃止

**Objective:** プレイヤーとして、コアのレベル管理をなくしたい。エージェント構成をシンプルにするため。

#### Acceptance Criteria

1. The Agent System shall コアからレベルに関する全てのデータと機能を削除する
2. The Agent System shall コアをレベルなしで管理する（TypeIDのみで識別）
3. The Core Inventory shall 最大レベル追跡機能を削除し、TypeIDベースの保有フラグのみで管理する
4. When コアを取得する, the Agent System shall レベル情報を保存せずにTypeIDのみを記録する

### Requirement 2: エージェントステータス計算の変更

**Objective:** プレイヤーとして、コアの重みに基づく固定ステータスを持つエージェントを使いたい。直感的なパラメータ理解のため。

#### Acceptance Criteria

1. The Agent System shall エージェントの各ステータス（STR/INT/WIL/LUK）を「100 x コアの重み」で計算する
2. The Stats Service shall レベルパラメータを受け取らず、コア重みのみでステータスを計算する
3. The Agent System shall 全てのエージェントに対して同一の計算式を適用する

### Requirement 3: コアレベル関連機能の削除

**Objective:** 開発者として、コアレベルに関する全てのコードを削除したい。コードベースの簡潔化のため。

#### Acceptance Criteria

1. The Agent Slot Manager shall コアレベル選択機能を削除する
2. The Agent Customization Screen shall コアレベル選択UIを削除する
3. The Reward System shall コアドロップ時のレベル計算処理を削除する
4. The Save Data System shall コアのレベル情報の読み書きを削除する
5. The Encyclopedia Screen shall コアのレベル表示を削除する

### Requirement 4: プレイヤー最大HPの初期化

**Objective:** プレイヤーとして、ゲーム開始時に適切な最大HPを持ちたい。バランスの取れたゲーム体験のため。

#### Acceptance Criteria

1. The Player System shall 新規ゲーム開始時に最大HPを1000に設定する
2. The Save Data System shall 新規セーブデータ作成時に最大HPを1000に初期化する

### Requirement 5: 敵レベルの初期化

**Objective:** プレイヤーとして、全ての敵がレベル1から始まるようにしたい。公平な難易度スタートのため。

#### Acceptance Criteria

1. The Enemy System shall 全ての敵タイプのデフォルトレベルを1に設定する
2. When 敵データを読み込む, the Enemy System shall デフォルトレベルを1として扱う

### Requirement 6: 敵ランクシステムの導入

**Objective:** プレイヤーとして、段階的に強い敵に挑戦したい。達成感のある進行のため。

#### Acceptance Criteria

1. The Enemy System shall 敵をランク（1から開始）で分類する
2. The Battle Select Screen shall 現在のランクの敵のみを表示する
3. When 現在のランクの全ての敵を撃破する, the Battle Select Screen shall 次のランクの敵を解放する
4. The Enemy System shall 各敵タイプにランク情報を持たせる
5. The Save Data System shall 現在解放済みのランクを永続化する

### Requirement 7: 敵撃破による最大HP増加

**Objective:** プレイヤーとして、新しい敵を倒すことで強くなりたい。成長実感のため。

#### Acceptance Criteria

1. When 未撃破の敵を初撃破する, the Player System shall 最大HPを10増加させる
2. The Save Data System shall 撃破済み敵リストを永続化する
3. While 既に撃破済みの敵を再度倒す, the Player System shall 最大HPを増加させない（初撃破時の報酬のみ）

### Requirement 8: 撃破済み敵のレベル選択

**Objective:** プレイヤーとして、倒した敵を高レベルで再挑戦したい。やり込み要素のため。

#### Acceptance Criteria

1. While 敵が撃破済みの状態である, the Battle Select Screen shall レベルを選択可能にする
2. When 敵を初撃破する, the Battle Select Screen shall レベル5までを選択可能にする
3. When 敵を特定レベルで撃破する, the Battle Select Screen shall 撃破レベル+5までを選択可能にする
4. The Save Data System shall 敵ごとの撃破済み最大レベルを永続化する

### Requirement 9: 高レベル敵撃破による最大HP増加

**Objective:** プレイヤーとして、より高いレベルの敵に挑戦することで成長したい。高難度への動機付けのため。

#### Acceptance Criteria

1. When 新記録レベルで敵を撃破する, the Player System shall 最大HPを「(新記録レベル - 1) x 10」だけ増加させる
2. While 過去に撃破したレベル以下で敵を倒す, the Player System shall 最大HPを増加させない
3. The Save Data System shall 敵ごとの撃破済み最大レベルを更新する

### Requirement 10: 敵進行状況の表示

**Objective:** プレイヤーとして、敵の撃破状況を確認したい。進行状況の把握のため。

#### Acceptance Criteria

1. The Battle Select Screen shall 各敵の撃破状況（未撃破/撃破済み）を表示する
2. The Battle Select Screen shall 撃破済み敵の最大撃破レベルを表示する
3. The Battle Select Screen shall 現在のランクと進行状況を表示する

### Requirement 11: 後方互換性の破棄

**Objective:** 開発者として、古いセーブデータとの互換性を考慮しないで実装したい。開発速度向上のため。

#### Acceptance Criteria

1. The Save Data System shall コアレベル情報を含む旧形式のセーブデータを読み込まない
2. The Save Data System shall 新形式のセーブデータ構造のみをサポートする
3. If 旧形式のセーブデータを検出する, the Save Data System shall 新規ゲームとして開始する（またはエラーを表示する）
