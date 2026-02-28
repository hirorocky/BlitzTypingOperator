// Package screens はTUIゲームの画面を提供します。
// battle_logic.go はバトル画面のゲームロジックを担当します。
package screens

import (
	"fmt"
	"math/rand"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/combat"
	"hirorocky/type-battle/internal/usecase/combat/chain"
	"hirorocky/type-battle/internal/usecase/typing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// randFloat は0.0〜1.0の乱数を返す関数です。
// テスト時にモック可能にするため変数として定義しています。
var randFloat = func() float64 {
	return rand.Float64()
}

// ==================== ゲームロジック: 状態判定 ====================

// checkGameOver は勝敗を判定します。
func (s *BattleScreen) checkGameOver() bool {
	// プレイヤー敗北
	if s.player.HP <= 0 {
		s.gameOver = true
		s.victory = false
		s.message = "敗北..."
		return true
	}

	// プレイヤー勝利
	if s.enemy.HP <= 0 {
		s.gameOver = true
		s.victory = true
		s.message = "勝利！"
		return true
	}

	return false
}

// createGameOverCmd はゲーム終了時のコマンドを作成します。
func (s *BattleScreen) createGameOverCmd() tea.Cmd {
	result := BattleResultMsg{
		Victory:   s.victory,
		Level:     s.enemy.Level,
		Stats:     s.battleState.Stats,
		EnemyID:   s.enemy.Type.ID,
		EnemyType: &s.enemy.Type,
	}
	return func() tea.Msg {
		return result
	}
}

// IsGameOver はゲームが終了したかを返します。
func (s *BattleScreen) IsGameOver() bool {
	return s.gameOver
}

// IsVictory は勝利したかを返します。
func (s *BattleScreen) IsVictory() bool {
	return s.gameOver && s.victory
}

// IsDefeat は敗北したかを返します。
func (s *BattleScreen) IsDefeat() bool {
	return s.gameOver && !s.victory
}

// IsShowingResult は結果表示中かを返します。
func (s *BattleScreen) IsShowingResult() bool {
	return s.showingResult
}

// ==================== ゲームロジック: 敵攻撃処理 ====================

// processEnemyAttack は敵のターンを処理します。
// ビジネスロジックは BattleEngine.ProcessEnemyTurn() に委譲し、UI更新のみを担当します。
func (s *BattleScreen) processEnemyAttack() {
	if s.battleEngine == nil || s.battleState == nil {
		s.processLegacyEnemyAttack()
		return
	}

	// ビジネスロジックはエンジンに委譲
	result := s.battleEngine.ProcessEnemyTurn(s.battleState)

	// UI更新
	s.updateUIAfterEnemyTurn(result)
}

// processLegacyEnemyAttack はフォールバック用の従来攻撃処理です。
func (s *BattleScreen) processLegacyEnemyAttack() {
	damage := s.enemy.AttackPower
	s.player.HP -= damage
	if s.player.HP < 0 {
		s.player.HP = 0
	}
	s.message = fmt.Sprintf("%sの攻撃！ %dダメージを受けた！", s.enemy.Name, damage)
	// 次の行動を準備してチャージ開始
	s.enemy.PrepareNextAction()
	if action := s.enemy.GetNextAction(); action != nil {
		s.enemy.StartCharging(*action, time.Now())
	}
	s.floatingDamageManager.AddDamage(damage, "player")
	s.playerHPBar.SetTarget(s.player.HP)
}

// updateUIAfterEnemyTurn は敵ターン結果に基づいてUIを更新します。
func (s *BattleScreen) updateUIAfterEnemyTurn(result combat.EnemyTurnResult) {
	switch result.ActionType {
	case domain.EnemyActionAttack:
		s.message = fmt.Sprintf("%sの攻撃！ %s", s.enemy.Name, result.Message)
		if result.Damage > 0 {
			s.floatingDamageManager.AddDamage(result.Damage, "player")
			s.playerHPBar.SetTarget(s.player.HP)
		}
		// ps_quick_recovery: 被ダメージ時にリキャスト短縮
		s.evaluateQuickRecovery()

	case domain.EnemyActionBuff:
		s.message = fmt.Sprintf("%sが%s！", s.enemy.Name, result.Message)

	case domain.EnemyActionDebuff:
		s.message = fmt.Sprintf("%sが%s", s.enemy.Name, result.Message)

	case domain.EnemyActionDefense:
		s.message = fmt.Sprintf("%sが%s！", s.enemy.Name, result.Message)

	default:
		s.message = "敵の行動"
	}

	// フェーズ変化をUIに反映
	if result.PhaseChanged {
		s.message += " [敵が強化フェーズに突入！]"
	}
}

// evaluateQuickRecovery はps_quick_recoveryの発動を評価し、リキャストを短縮します。
func (s *BattleScreen) evaluateQuickRecovery() {
	if s.battleEngine == nil || s.battleState == nil || s.recastManager == nil {
		return
	}

	for _, agent := range s.battleState.EquippedAgents {
		reduction := s.battleEngine.EvaluateQuickRecovery(s.battleState, agent)
		if reduction > 0 {
			s.recastManager.ReduceAllRecasts(time.Duration(reduction) * time.Second)
			s.message += " [クイックリカバリー発動！]"
		}
	}
}

// updateEffectDurations はバフ・デバフの持続時間を更新します。
func (s *BattleScreen) updateEffectDurations(deltaSeconds float64) {
	// プレイヤーのエフェクトを更新
	if s.player.EffectTable != nil {
		s.player.EffectTable.UpdateDurations(deltaSeconds)
	}

	// 敵のエフェクトを更新
	if s.enemy.EffectTable != nil {
		s.enemy.EffectTable.UpdateDurations(deltaSeconds)
	}
}

// ==================== ゲームロジック: クールダウン ====================

// UpdateCooldowns はクールダウンを更新します。
func (s *BattleScreen) UpdateCooldowns(deltaSeconds float64) {
	for i := range s.skillSlots {
		if s.skillSlots[i].CooldownRemaining > 0 {
			s.skillSlots[i].CooldownRemaining -= deltaSeconds
			if s.skillSlots[i].CooldownRemaining < 0 {
				s.skillSlots[i].CooldownRemaining = 0
			}
		}
	}
}

// StartCooldown はスキルのクールダウンを開始します。
// EffectTableからCooldownReduceを取得して初期値を短縮します。
func (s *BattleScreen) StartCooldown(slotIndex int, duration float64) {
	if slotIndex >= 0 && slotIndex < len(s.skillSlots) {
		reducedDuration := duration

		// CooldownReduceを取得して初期値を短縮
		if s.player != nil && s.player.EffectTable != nil {
			ctx := domain.NewEffectContext(s.player.HP, s.player.MaxHP, 0, 0)
			if s.enemy != nil {
				ctx = domain.NewEffectContext(s.player.HP, s.player.MaxHP, s.enemy.HP, s.enemy.MaxHP)
			}
			effects := s.player.EffectTable.Aggregate(ctx)

			// CooldownReduce を適用（正=短縮、負=延長）
			// 30%短縮の場合、CooldownReduce=0.3 → duration * (1 - 0.3) = 70%
			reducedDuration = duration * (1.0 - effects.CooldownReduce)

			// 最低10%は残す（極端な短縮対策）
			minDuration := duration * 0.1
			if reducedDuration < minDuration {
				reducedDuration = minDuration
			}
		}

		s.skillSlots[slotIndex].CooldownRemaining = reducedDuration
		s.skillSlots[slotIndex].CooldownTotal = duration // 表示用に元の値を保持
	}
}

// ==================== ゲームロジック: リキャスト管理 ====================

// UpdateRecasts はリキャスト時間を更新します。
func (s *BattleScreen) UpdateRecasts(deltaSeconds float64) {
	if s.recastManager == nil {
		return
	}

	// リキャスト時間を更新（deltaSecondsをtime.Durationに変換）
	delta := time.Duration(deltaSeconds * float64(time.Second))
	s.recastManager.UpdateRecast(delta)
}

// isSkillUsable は指定スロットのスキルが使用可能かを判定します。
// スキルのクールダウンとエージェントのリキャスト状態を両方チェックします。
func (s *BattleScreen) isSkillUsable(slotIndex int) bool {
	if slotIndex < 0 || slotIndex >= len(s.skillSlots) {
		return false
	}

	slot := s.skillSlots[slotIndex]

	// ディフェンススキル未解放時はディフェンスタイプを使用不可
	if !s.isDefenseSkillUnlocked() && slot.Skill.GetChallengeType() == domain.ChallengeTypeDefense {
		return false
	}

	// スキルのクールダウンチェック
	if !slot.IsReady() {
		return false
	}

	// エージェントのリキャストチェック
	if s.recastManager != nil && !s.recastManager.IsAgentReady(slot.AgentIndex) {
		return false
	}

	// マナ不足チェック（マナシステム未解放時はスキップ）
	if s.isManaSystemUnlocked() && slot.Skill.Type.ManaCost > 0 && s.player.Mana < slot.Skill.Type.ManaCost {
		return false
	}

	return true
}

// startCooldownAndRecast は選択中スキルのクールダウンとリキャストを一括で開始します。
// チャレンジ終了時に呼ぶことで、タイピング中はCDが進行しないようにします。
func (s *BattleScreen) startCooldownAndRecast() {
	if s.selectedSkillIdx < 0 || s.selectedSkillIdx >= len(s.skillSlots) {
		return
	}
	slot := s.skillSlots[s.selectedSkillIdx]
	s.StartCooldown(s.selectedSkillIdx, slot.CooldownTotal)
	s.startAgentRecast(slot.AgentIndex, slot.Skill)
}

// startAgentRecast はエージェントのリキャストを開始し、チェイン効果を登録します。
func (s *BattleScreen) startAgentRecast(agentIndex int, skill *domain.SkillModel) {
	if s.recastManager == nil {
		return
	}

	// スキルのクールダウン秒数を使用してリキャストを開始
	cooldownDuration := time.Duration(skill.CooldownSeconds() * float64(time.Second))
	s.recastManager.StartRecast(agentIndex, cooldownDuration)

	// チェイン効果の管理: スキルにチェイン効果があれば登録、なければ既存効果をクリア
	if s.chainEffectManager != nil {
		if skill.ChainEffect != nil {
			s.chainEffectManager.RegisterChainEffect(agentIndex, skill.ChainEffect, skill.TypeID)
		} else {
			s.chainEffectManager.ClearEffectForAgent(agentIndex)
		}
	}
}

// triggerChainEffects はスキル使用時に他エージェントのチェイン効果を発動します。
func (s *BattleScreen) triggerChainEffects(usingAgentIndex int, effectFlags chain.SkillEffectFlags) {
	if s.chainEffectManager == nil {
		return
	}

	// チェイン効果の発動をチェック
	triggered := s.chainEffectManager.CheckAndTrigger(usingAgentIndex, effectFlags)

	// 発動した効果を適用
	for _, effect := range triggered {
		s.applyTriggeredChainEffect(&effect)
	}
}

// chainEffectDuration はチェイン効果の持続時間（秒）です。
const chainEffectDuration = 10.0

// applyTriggeredChainEffect は発動したチェイン効果を適用します。
func (s *BattleScreen) applyTriggeredChainEffect(effect *chain.TriggeredChainEffect) {
	// 効果タイプに応じた処理
	switch effect.Effect.Type {
	case domain.ChainEffectDamageBonus:
		// 追加ダメージ（敵へのダメージ）- 即時適用
		bonusDamage := int(effect.EffectValue)
		if s.enemy != nil {
			s.enemy.HP -= bonusDamage
			if s.enemy.HP < 0 {
				s.enemy.HP = 0
			}
			s.floatingDamageManager.AddDamage(bonusDamage, "enemy")
			s.enemyHPBar.SetTarget(s.enemy.HP)
			s.message = fmt.Sprintf("チェイン発動！ %s (+%dダメージ)", effect.Message, bonusDamage)
		}

	case domain.ChainEffectHealBonus:
		// 追加回復 - 即時適用
		bonusHeal := int(effect.EffectValue)
		if s.player != nil {
			s.player.HP += bonusHeal
			if s.player.HP > s.player.MaxHP {
				s.player.HP = s.player.MaxHP
			}
			s.floatingDamageManager.AddHeal(bonusHeal, "player")
			s.playerHPBar.SetTarget(s.player.HP)
			s.message = fmt.Sprintf("チェイン発動！ %s (+%d回復)", effect.Message, bonusHeal)
		}

	case domain.ChainEffectBuffExtend, domain.ChainEffectBuffDuration:
		// バフ延長 - 即時適用
		if s.player != nil && s.player.EffectTable != nil {
			s.player.EffectTable.ExtendBuffDurations(effect.EffectValue)
			s.message = fmt.Sprintf("チェイン発動！ %s", effect.Message)
		}

	case domain.ChainEffectDebuffExtend, domain.ChainEffectDebuffDuration:
		// デバフ延長 - 即時適用
		if s.enemy != nil && s.enemy.EffectTable != nil {
			s.enemy.EffectTable.ExtendDebuffDurations(effect.EffectValue)
			s.message = fmt.Sprintf("チェイン発動！ %s", effect.Message)
		}

	default:
		// 持続効果は EffectTable に登録
		s.registerChainEffectToTable(effect)
		s.message = fmt.Sprintf("チェイン発動！ %s", effect.Message)
	}
}

// registerChainEffectToTable はチェイン効果を EffectTable に登録します。
func (s *BattleScreen) registerChainEffectToTable(effect *chain.TriggeredChainEffect) {
	if s.player == nil || s.player.EffectTable == nil {
		return
	}

	// チェイン効果の値を EffectColumn にマッピング
	values := make(map[domain.EffectColumn]float64)
	flags := make(map[domain.EffectColumn]bool)

	switch effect.Effect.Type {
	// 攻撃強化カテゴリ
	case domain.ChainEffectDamageAmp:
		values[domain.ColDamageMultiplier] = 1.0 + effect.EffectValue/100.0
	case domain.ChainEffectArmorPierce:
		flags[domain.ColArmorPierce] = true
	case domain.ChainEffectLifeSteal:
		values[domain.ColLifeSteal] = effect.EffectValue / 100.0

	// 防御強化カテゴリ
	case domain.ChainEffectDamageCut:
		values[domain.ColDamageCut] = effect.EffectValue / 100.0
	case domain.ChainEffectEvasion:
		values[domain.ColEvasion] = effect.EffectValue / 100.0
	case domain.ChainEffectReflect:
		values[domain.ColReflect] = effect.EffectValue / 100.0
	case domain.ChainEffectRegen:
		values[domain.ColRegen] = effect.EffectValue

	// 回復強化カテゴリ
	case domain.ChainEffectHealAmp:
		values[domain.ColHealMultiplier] = 1.0 + effect.EffectValue/100.0
	case domain.ChainEffectOverheal:
		flags[domain.ColOverheal] = true

	// タイピングカテゴリ
	case domain.ChainEffectTimeExtend:
		values[domain.ColTimeExtend] = effect.EffectValue
	case domain.ChainEffectAutoCorrect:
		values[domain.ColAutoCorrect] = effect.EffectValue

	// リキャストカテゴリ
	case domain.ChainEffectCooldownReduce:
		values[domain.ColCooldownReduce] = effect.EffectValue / 100.0

	// 特殊カテゴリ
	case domain.ChainEffectDoubleCast:
		values[domain.ColDoubleCast] = effect.EffectValue / 100.0
	}

	// EffectEntry を作成して登録
	duration := chainEffectDuration
	entry := domain.EffectEntry{
		SourceType:  domain.SourceChain,
		SourceID:    string(effect.Effect.Type),
		SourceIndex: effect.SourceAgentIndex,
		Name:        effect.Effect.Description,
		Duration:    &duration,
		Values:      values,
		Flags:       flags,
	}

	s.player.EffectTable.AddEntry(entry)
}

// ==================== ゲームロジック: チャレンジ ====================

// startChallenge はChallengeInputを構築してチャレンジを開始します。
// EffectTableからTimeExtend、AutoCorrect、TypingDifficultyを取得して適用します。
func (s *BattleScreen) startChallenge(skill *domain.SkillModel) tea.Cmd {
	// 新しいチャレンジ開始時に既存のリザルト表示をクリア
	s.resultDisplay = nil
	// EffectTableから効果を取得
	var effects domain.EffectResult
	if s.player != nil && s.player.EffectTable != nil {
		ctx := domain.NewEffectContext(s.player.HP, s.player.MaxHP, 0, 0)
		if s.enemy != nil {
			ctx = domain.NewEffectContext(s.player.HP, s.player.MaxHP, s.enemy.HP, s.enemy.MaxHP)
		}
		effects = s.player.EffectTable.Aggregate(ctx)
	} else {
		effects = domain.NewEffectResult()
	}

	// 基礎DifficultyRateにTypingDifficulty修正を適用
	baseDifficulty := float64(skill.GetDifficultyRate())
	adjustedDifficulty := baseDifficulty * effects.TypingDifficulty
	difficultyRate := domain.DifficultyRate(int(adjustedDifficulty)).Clamp()

	// パッシブスキル: ps_typo_recovery（ミス時の時間延長）
	mistakeTimeExtendSec := 0.0
	if s.battleEngine != nil && s.battleState != nil {
		slot := s.skillSlots[s.selectedSkillIdx]
		agent := slot.Agent
		mistakeTimeExtendSec = s.battleEngine.EvaluateTypoRecovery(s.battleState, agent)
	}

	// パッシブスキル: ps_second_chance（タイムアウト時の再挑戦）
	retryOnTimeout := false
	retryTimeLimitMultiplier := 0.5
	if s.battleEngine != nil && s.battleState != nil {
		slot := s.skillSlots[s.selectedSkillIdx]
		agent := slot.Agent
		if s.battleEngine.EvaluateSecondChance(s.battleState, agent) {
			retryOnTimeout = true
		}
	}

	input := domain.ChallengeInput{
		Difficulty:               difficultyRate,
		Words:                    s.dictionary,
		TimeExtendSec:            effects.TimeExtend,
		AutoCorrectCount:         effects.AutoCorrect,
		MistakeTimeExtendSec:     mistakeTimeExtendSec,
		RetryOnTimeout:           retryOnTimeout,
		RetryTimeLimitMultiplier: retryTimeLimitMultiplier,
		ChallengeOptions:         skill.Type.ChallengeOptions,
	}

	challengeType := skill.GetChallengeType()
	s.activeChallenge = challenges.New(challengeType, input)

	if s.activeChallenge != nil {
		return s.activeChallenge.Init()
	}
	return nil
}

// handleChallengeComplete はチャレンジ完了時の処理を行います。
// ChallengeOutputをTypingResultに変換し、スキル効果パイプラインを実行します。
func (s *BattleScreen) handleChallengeComplete(result *domain.ChallengeOutput) {
	// チャレンジをクリア
	s.activeChallenge = nil

	// クールダウンとリキャストを開始（全終了パターン共通）
	s.startCooldownAndRecast()

	// キャンセル時は効果を適用しない
	if result.Status == domain.ChallengeCancel {
		s.message = "タイピングキャンセル"
		return
	}

	// 失敗時は効果を適用しない
	if result.Status == domain.ChallengeFail {
		s.message = "タイムアウト！"
		return
	}

	// ChallengeOutput → TypingResult に変換
	typingResult := &typing.TypingResult{
		Completed:      true,
		WPM:            result.WPM,
		Score:          result.Score,
		CompletionTime: result.CompletionTime,
	}

	// コンボカウントの更新（ps_combo_master用）
	if result.Status == domain.ChallengePerfect {
		s.comboCount++
	} else {
		s.comboCount = 0
	}

	// パーフェクト判定（Status==ChallengePerfect で判定）
	isPerfect := result.Status == domain.ChallengePerfect

	slot := s.skillSlots[s.selectedSkillIdx]

	// 潜在効果の解放チェック（未解放時はisPerfect=falseにリセット）
	if isPerfect && !s.isLatentEffectUnlocked() {
		isPerfect = false
	}
	typingResult.IsPerfect = isPerfect

	// バトル統計に記録
	if s.battleEngine != nil && s.battleState != nil {
		s.battleEngine.RecordTypingResult(s.battleState, typingResult)
	}

	// スキル効果を適用
	agent := slot.Agent
	skill := slot.Skill
	agentIndex := slot.AgentIndex

	// スキルの効果フラグを取得
	effectFlags := getSkillEffectFlags(skill)

	// 他エージェントの待機中チェイン効果を発動（スキル効果適用前）
	s.triggerChainEffects(agentIndex, effectFlags)

	// DoubleCast判定
	doubleCastTriggered := false
	if s.player != nil && s.player.EffectTable != nil {
		ctx := domain.NewEffectContext(s.player.HP, s.player.MaxHP, 0, 0)
		if s.enemy != nil {
			ctx = domain.NewEffectContext(s.player.HP, s.player.MaxHP, s.enemy.HP, s.enemy.MaxHP)
		}
		effects := s.player.EffectTable.Aggregate(ctx)
		if effects.DoubleCast > 0 {
			if randFloat() < effects.DoubleCast {
				doubleCastTriggered = true
			}
		}
	}

	// ps_echo_skill判定（スキル2回発動）
	echoSkillRepeat := 1
	echoSkillTriggered := false
	if s.battleEngine != nil && s.battleState != nil {
		echoSkillRepeat = s.battleEngine.EvaluateEchoSkill(s.battleState, agent)
		if echoSkillRepeat > 1 {
			echoSkillTriggered = true
		}
	}

	// ps_miracle_heal判定（回復スキル時HP全回復）
	miracleHealTriggered := false
	if s.battleEngine != nil && s.battleState != nil {
		if s.battleEngine.EvaluateMiracleHeal(s.battleState, agent, skill) {
			miracleHealTriggered = true
		}
	}

	// マナ消費（Echo/DoubleCastの追加発動ではコストを消費しない）
	if skill.Type.ManaCost > 0 && s.player != nil {
		s.player.ConsumeMana(skill.Type.ManaCost)
	}

	// 効果結果を集約（ResultDisplayState用）
	combinedEffectResult := &combat.SkillEffectResult{}
	var effectAmount int
	if s.battleEngine != nil && s.battleState != nil {
		effectResult := s.battleEngine.ApplySkillEffectWithCombo(s.battleState, agent, skill, typingResult, s.comboCount)
		effectAmount = effectResult.TotalDamage
		combinedEffectResult.TotalDamage += effectResult.TotalDamage
		combinedEffectResult.TotalHeal += effectResult.TotalHeal
		combinedEffectResult.NormalEffects = append(combinedEffectResult.NormalEffects, effectResult.NormalEffects...)
		combinedEffectResult.LatentEffects = append(combinedEffectResult.LatentEffects, effectResult.LatentEffects...)

		for i := 1; i < echoSkillRepeat; i++ {
			additionalResult := s.battleEngine.ApplySkillEffectWithCombo(s.battleState, agent, skill, typingResult, s.comboCount)
			effectAmount += additionalResult.TotalDamage
			combinedEffectResult.TotalDamage += additionalResult.TotalDamage
			combinedEffectResult.TotalHeal += additionalResult.TotalHeal
			combinedEffectResult.NormalEffects = append(combinedEffectResult.NormalEffects, additionalResult.NormalEffects...)
			combinedEffectResult.LatentEffects = append(combinedEffectResult.LatentEffects, additionalResult.LatentEffects...)
		}

		if doubleCastTriggered {
			secondResult := s.battleEngine.ApplySkillEffectWithCombo(s.battleState, agent, skill, typingResult, s.comboCount)
			effectAmount += secondResult.TotalDamage
			combinedEffectResult.TotalDamage += secondResult.TotalDamage
			combinedEffectResult.TotalHeal += secondResult.TotalHeal
			combinedEffectResult.NormalEffects = append(combinedEffectResult.NormalEffects, secondResult.NormalEffects...)
			combinedEffectResult.LatentEffects = append(combinedEffectResult.LatentEffects, secondResult.LatentEffects...)
		}

		if miracleHealTriggered {
			s.player.HP = s.player.MaxHP
			s.playerHPBar.SetTarget(s.player.HP)
		}
	}

	// フローティングダメージ/回復とHPアニメーション
	if effectAmount > 0 {
		if effectFlags.HasDamage {
			s.floatingDamageManager.AddDamage(effectAmount, "enemy")
			s.enemyHPBar.SetTarget(s.enemy.HP)
		} else if effectFlags.HasHeal {
			s.floatingDamageManager.AddHeal(effectAmount, "player")
			s.playerHPBar.SetTarget(s.player.HP)
		}
	}

	// メッセージを表示
	s.message = s.formatEffectMessage(skill, effectAmount, typingResult, effectFlags)
	if s.comboCount > 0 {
		s.message += fmt.Sprintf(" [コンボ:%d]", s.comboCount)
	}
	if echoSkillTriggered {
		s.message += " [エコースキル発動！]"
	}
	if miracleHealTriggered {
		s.message += " [ミラクルヒール発動！]"
	}
	if doubleCastTriggered {
		s.message += " [ダブルキャスト発動！]"
	}

	// フェーズ変化をチェック
	if s.battleEngine != nil && s.battleState != nil {
		if s.battleEngine.CheckPhaseTransition(s.battleState) {
			s.battleEngine.SwitchEnemyPassive(s.battleState)
			s.message += " [敵が強化フェーズに突入！]"
		}
	}

	// リザルト表示を設定
	s.resultDisplay = &ResultDisplayState{
		AgentIndex:   agentIndex,
		Status:       result.Status,
		EffectResult: combinedEffectResult,
		StartTime:    time.Now(),
	}

	// カーソルを別エージェントに移動（複数エージェントの場合のみ）
	if len(s.equippedAgents) > 1 {
		s.moveToNextAvailableAgent(agentIndex)
	}
}

// processEnemyAttackWithDefense はディフェンスチャレンジ中の敵攻撃を処理します。
// 防御率を適用してダメージを軽減し、チャレンジを自動終了させます。
func (s *BattleScreen) processEnemyAttackWithDefense(dp challenges.DefenseProvider) {
	// 防御率を取得してチャレンジを自動終了
	defenseRate := dp.DefenseRate()
	dp.CompleteByAttack()
	s.activeChallenge = nil

	// クールダウンとリキャストを開始
	s.startCooldownAndRecast()

	if s.battleEngine == nil || s.battleState == nil {
		return
	}

	// 通常の敵ターン処理を実行
	result := s.battleEngine.ProcessEnemyTurn(s.battleState)

	// 攻撃ダメージに防御率を適用（DamageCutとして機能）
	if result.ActionType == domain.EnemyActionAttack && result.Damage > 0 && defenseRate > 0 {
		// 防御率分のダメージを軽減（既に適用済みのダメージを修正）
		reducedDamage := int(float64(result.Damage) * defenseRate)
		s.player.HP += reducedDamage // ProcessEnemyTurnで減らされた分を戻す
		if s.player.HP > s.player.MaxHP {
			s.player.HP = s.player.MaxHP
		}
		result.Damage -= reducedDamage
		result.Message += fmt.Sprintf(" (防御率%.0f%%で%d軽減)", defenseRate*100, reducedDamage)
	}

	// UI更新
	s.updateUIAfterEnemyTurn(result)
}

// formatEffectMessage は効果メッセージをフォーマットします。
func (s *BattleScreen) formatEffectMessage(skill *domain.SkillModel, effectAmount int, result *typing.TypingResult, flags chain.SkillEffectFlags) string {
	var action string
	if flags.HasDamage {
		action = fmt.Sprintf("%dダメージを与えた！", effectAmount)
	} else if flags.HasHeal {
		action = fmt.Sprintf("%d回復した！", effectAmount)
	} else if flags.HasBuff {
		action = fmt.Sprintf("%sを付与した！", skill.Name())
	} else if flags.HasDebuff {
		action = fmt.Sprintf("敵に%sを付与した！", skill.Name())
	} else {
		action = "効果を発動した！"
	}

	return fmt.Sprintf("%s (WPM:%.0f 正確性:%d%%)", action, result.WPM, result.Score)
}

// ==================== ゲームロジック: 行動表示 ====================

// getActionDisplay はチャージ後行動の表示情報を返します。

func (s *BattleScreen) getActionDisplay() (icon string, text string, color lipgloss.Color) {
	if s.battleState == nil {
		return "?", "不明", styles.ColorSubtle
	}

	// チャージ中の場合はチャージ状態を表示
	if s.enemy != nil && s.enemy.WaitMode == domain.WaitModeCharging {
		return s.getChargingActionDisplay()
	}

	// ディフェンス中の場合はディフェンス状態を表示
	if s.enemy != nil && s.enemy.WaitMode == domain.WaitModeDefending {
		return s.getDefenseActionDisplay()
	}

	action := s.enemy.GetNextAction()
	if action == nil {
		return "?", "不明", styles.ColorSubtle
	}

	switch action.ActionType {
	case domain.EnemyActionAttack:
		// 攻撃予告（赤色）- バフ反映のため毎回計算
		expectedDamage := s.battleEngine.GetExpectedDamage(s.battleState)
		if action.AttackType == "magic" {
			return "💥", fmt.Sprintf("魔法%dダメージ", expectedDamage), styles.ColorDamage
		}
		return "⚔️", fmt.Sprintf("物理%dダメージ", expectedDamage), styles.ColorDamage

	case domain.EnemyActionBuff:
		// 自己バフ予告（黄色）- 効果内容を表示
		effectDesc := domain.DescribeSingleEffect(action.EffectColumn, action.EffectValue)
		return "💪", effectDesc, styles.ColorWarning

	case domain.EnemyActionDebuff:
		// プレイヤーデバフ予告（青色）- 効果内容を表示
		effectDesc := domain.DescribeSingleEffect(action.EffectColumn, action.EffectValue)
		return "💀", effectDesc, styles.ColorInfo
	}

	return "?", "不明", styles.ColorSubtle
}

// getChargingActionDisplay はチャージ中の行動表示情報を返します。
// チャージ後行動の効果を表示します（行動名ではなく効果説明）。
func (s *BattleScreen) getChargingActionDisplay() (icon string, text string, color lipgloss.Color) {
	// チャージ後行動の効果を表示
	if action := s.enemy.PendingAction; action != nil {
		switch action.ActionType {
		case domain.EnemyActionAttack:
			expectedDamage := s.battleEngine.GetExpectedDamage(s.battleState)
			if action.AttackType == "magic" {
				return "💥", fmt.Sprintf("魔法ダメージ%d", expectedDamage), styles.ColorDamage
			}
			return "⚔️", fmt.Sprintf("物理ダメージ%d", expectedDamage), styles.ColorDamage
		case domain.EnemyActionBuff:
			effectDesc := domain.DescribeSingleEffect(action.EffectColumn, action.EffectValue)
			return "💪", effectDesc, styles.ColorWarning
		case domain.EnemyActionDebuff:
			effectDesc := domain.DescribeSingleEffect(action.EffectColumn, action.EffectValue)
			return "💀", effectDesc, styles.ColorInfo
		default:
			return "?", action.Name, styles.ColorSubtle
		}
	}
	return "?", "不明", styles.ColorSubtle
}

// getDefenseActionDisplay はディフェンス中の行動表示情報を返します。
func (s *BattleScreen) getDefenseActionDisplay() (icon string, text string, color lipgloss.Color) {
	// 現在発動中のディフェンス効果を表示
	switch s.enemy.ActiveDefenseType {
	case domain.DefensePhysicalCut:
		return "🛡️", fmt.Sprintf("物理ダメージ%.0f%%カット", s.enemy.DefenseValue*100), styles.ColorInfo
	case domain.DefenseMagicCut:
		return "🛡️", fmt.Sprintf("魔法ダメージ%.0f%%カット", s.enemy.DefenseValue*100), styles.ColorInfo
	case domain.DefenseDebuffEvade:
		return "🛡️", fmt.Sprintf("デバフ%.0f%%回避", s.enemy.DefenseValue*100), styles.ColorInfo
	default:
		return "🛡️", "防御中", styles.ColorInfo
	}
}

// ==================== ゲームロジック: スキル選択ナビゲーション ====================

// selectFirstSkillOfAgent は指定エージェントの最初のスキルを選択します。
func (s *BattleScreen) selectFirstSkillOfAgent(agentIdx int) {
	for i, slot := range s.skillSlots {
		if slot.AgentIndex == agentIdx {
			s.selectedSlot = i
			return
		}
	}
}

// moveToPrevSkillInAgent は現在のエージェント内で前のスキルに移動します。
func (s *BattleScreen) moveToPrevSkillInAgent() {
	if len(s.skillSlots) == 0 {
		return
	}

	currentAgentIdx := s.selectedAgentIdx
	agentSkills := s.getSkillIndicesForAgent(currentAgentIdx)

	if len(agentSkills) == 0 {
		return
	}

	// 現在のスキルの位置を見つける
	currentPos := 0
	for i, idx := range agentSkills {
		if idx == s.selectedSlot {
			currentPos = i
			break
		}
	}

	// 前のスキルに移動（ループ）
	newPos := currentPos - 1
	if newPos < 0 {
		newPos = len(agentSkills) - 1
	}
	s.selectedSlot = agentSkills[newPos]
}

// moveToNextSkillInAgent は現在のエージェント内で次のスキルに移動します。
func (s *BattleScreen) moveToNextSkillInAgent() {
	if len(s.skillSlots) == 0 {
		return
	}

	currentAgentIdx := s.selectedAgentIdx
	agentSkills := s.getSkillIndicesForAgent(currentAgentIdx)

	if len(agentSkills) == 0 {
		return
	}

	// 現在のスキルの位置を見つける
	currentPos := 0
	for i, idx := range agentSkills {
		if idx == s.selectedSlot {
			currentPos = i
			break
		}
	}

	// 次のスキルに移動（ループ）
	newPos := currentPos + 1
	if newPos >= len(agentSkills) {
		newPos = 0
	}
	s.selectedSlot = agentSkills[newPos]
}

// closeResultIfOnAgent はカーソルがリザルト表示中のエージェントに移動した場合にリザルトを閉じます。
func (s *BattleScreen) closeResultIfOnAgent(agentIdx int) {
	if s.resultDisplay != nil && s.resultDisplay.AgentIndex == agentIdx {
		s.resultDisplay = nil
	}
}

// moveToNextAvailableAgent はリザルト表示対象以外のエージェントにカーソルを移動します。
func (s *BattleScreen) moveToNextAvailableAgent(resultAgentIdx int) {
	nextIdx := (resultAgentIdx + 1) % len(s.equippedAgents)
	s.selectedAgentIdx = nextIdx
	s.selectFirstSkillOfAgent(nextIdx)
}

// getSkillIndicesForAgent は指定エージェントのスキルスロットのインデックスを返します。
func (s *BattleScreen) getSkillIndicesForAgent(agentIdx int) []int {
	var indices []int
	for i, slot := range s.skillSlots {
		if slot.AgentIndex == agentIdx {
			indices = append(indices, i)
		}
	}
	return indices
}

// getSkillEffectFlags はスキルが持つ効果の種別フラグを取得します。
func getSkillEffectFlags(skill *domain.SkillModel) chain.SkillEffectFlags {
	flags := chain.SkillEffectFlags{}

	for _, effect := range skill.Type.Effects {
		if effect.IsDamageEffect() {
			flags.HasDamage = true
		}
		if effect.IsHealEffect() {
			flags.HasHeal = true
		}
		if effect.IsBuffEffect() {
			flags.HasBuff = true
		}
		if effect.IsDebuffEffect() {
			flags.HasDebuff = true
		}
	}

	return flags
}
