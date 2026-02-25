# 潜在効果システム

## 概要
パーフェクトタイピング（ミスなし・制限時間以内）時に追加の効果を発動させる「潜在効果」システムを実装する。パーフェクト時にはASCIIアート演出でフィードバックし、スキルのeffectsにある潜在効果フラグ付きの効果が確率で発動する。この機能はfeature_unlocksで段階的に解放される。

## 受け入れ基準

### パーフェクトタイピング判定
1. タイピングチャレンジでChallengeStatus=Successかつ Accuracy >= 1.0（ミスなし）のとき「パーフェクト」と判定する
2. ディフェンスタイプのチャレンジではパーフェクト判定を行わない（Accuracyが防御率を表すため）
3. パーフェクト判定はコンボカウント判定（既存のAccuracy >= 1.0チェック）と同一タイミングで行う

### パーフェクト演出
4. パーフェクト時、バトル画面中央部（通常エージェントボックスが表示される領域）に "PERFECT!" のASCIIアートを一時的に表示する
5. 約0.2秒後にASCIIアートが消え、エージェントボックスの通常表示に戻る
6. パーフェクト演出は潜在効果の解放状態に関係なく、常に表示する

### 潜在効果の発動
7. SkillEffectに`IsLatent`フラグを追加する。`is_latent: true`の効果はパーフェクト時にのみ発動判定される
8. `is_latent: false`（デフォルト）の効果は従来通り常に発動判定される
9. 潜在効果の発動確率はProbability + LUK補正（`probability + LUK × luk_factor`）に従う。現在の実装は`(LUK - 10) × luk_factor`と誤ってBaseLUK=10を減算しているため、`LUK × luk_factor`に修正する
10. 1つのスキルに複数の潜在効果がある場合、各効果は独立して発動判定される
11. Echo/DoubleCastの追加発動時は、初回と同じパーフェクト状態を引き継ぐ（追加発動でも潜在効果が発動する）

### 機能解放
12. `latent_effect`をFeatureIDとしてfeature_unlocksに追加し、ランク6で解放する
13. 未解放時はパーフェクトを達成しても潜在効果は発動しない（パーフェクト演出は表示される）
14. 解放時にチュートリアルを表示する
15. 潜在効果の機能解放チェックはTUI層で行い、未解放時はisPerfect=falseとしてBattleEngineに渡す

### マスタデータ互換性
16. skills.jsonの各effectに`is_latent`フィールドを追加する。省略時のデフォルトはfalse
17. 既存のスキルデータは変更なしで動作する（後方互換性）

## 設計

### 変更対象
- `internal/domain/skill_effect.go`: SkillEffectに`IsLatent bool`フィールドを追加。`AdjustedProbability()`のLUK補正式を`(luk-BaseLUK)*LUKFactor`から`luk*LUKFactor`に修正し、`BaseLUK`定数を削除
- `internal/domain/feature_unlock.go`: `FeatureLatentEffect FeatureID = "latent_effect"`定数を追加
- `internal/usecase/typing/typing.go`: TypingResultに`IsPerfect bool`フィールドを追加
- `internal/usecase/combat/battle.go`: ApplySkillEffect内でIsLatentフラグによるフィルタリングを追加
- `internal/tui/screens/battle.go`: `showingPerfect bool`状態フィールド追加、`isLatentEffectUnlocked()`メソッド追加
- `internal/tui/screens/battle_logic.go`: handleChallengeComplete内でパーフェクト判定とisPerfect設定、パーフェクト演出トリガー
- `internal/tui/screens/battle_view.go`: パーフェクトASCIIアートのレンダリング（`showingPerfect`時にエージェントエリアを置き換え）
- `internal/infra/masterdata/loader.go`: SkillEffectData構造体にIsLatentフィールド追加
- `internal/app/masterdata_converter.go`: ResolveSkillTimedEffects内でIsLatentのマッピング追加
- `internal/infra/masterdata/data/skills.json`: 潜在効果を持たせたいスキルにis_latent: trueを追加
- `internal/infra/masterdata/data/feature_unlocks.json`: latent_effectエントリ追加
- `internal/infra/masterdata/data/tutorials.json`: 潜在効果チュートリアル追加

### 新規作成
- `internal/tui/ascii/perfect.go`: "PERFECT!" ASCIIアートの定義とレンダラー
- `internal/tui/ascii/perfect_test.go`: ASCIIアートのテスト

### データフロー

```
タイピングチャレンジ完了（ChallengeOutput）
  ↓
battle_logic.handleChallengeComplete()
  ├─ Accuracy >= 1.0 → isPerfect = true
  ├─ isPerfect → showingPerfect = true（演出開始）
  ├─ isLatentEffectUnlocked() = false → isPerfect = false にリセット
  ├─ TypingResult { IsPerfect: isPerfect, ... } を構築
  ↓
BattleEngine.ApplySkillEffect(state, agent, skill, typingResult)
  for effect in skill.Effects:
    ├─ effect.IsLatent && !typingResult.IsPerfect → スキップ
    ├─ effect.ShouldTrigger(LUK) → 確率判定
    └─ 効果適用（ダメージ計算、バフ/デバフ付与、マナ獲得）
```

## メモ
- 既存のAccuracyFactor乗算（calculateHPChange内）はそのまま残す。これは非パーフェクト時のペナルティとして機能しており、潜在効果とは独立した仕組み
- issueで「パーフェクト時に威力が上がる仮実装を置き換える」と記載があるが、現実装ではAccuracyFactor=1.0で「フルダメージ」になるだけで実質的なボーナスはない。潜在効果がパーフェクト時の真のボーナスとなる
- 解放ランクは6を想定（既存: 2=エージェント管理、3=ディフェンス、4=チェイン、5=マナ）
- パーフェクト演出の実装方式: BattleScreenに`showingPerfect bool`と`perfectTimer`を持たせ、バトルtickで時間経過をカウントして自動消去する（showingResultパターンを参考）
- 潜在効果を持つスキルの具体的なマスタデータ設計（どのスキルにis_latent=trueの効果を追加するか）は、実装フェーズで決定する
