// Package screens はTUIゲームの画面を提供します。
// battle_view.go はバトル画面のUIレンダリングロジックを担当します。
package screens

import (
	"fmt"
	"strings"
	"time"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/components"
	"hirorocky/type-battle/internal/tui/styles"
	"hirorocky/type-battle/internal/usecase/combat/recast"

	"github.com/charmbracelet/lipgloss"
)

// ==================== UIレンダリング ====================

// View は画面をレンダリングします。
// UI-Improvement Requirement 3.1: 3エリアレイアウト（敵情報、エージェント、プレイヤー情報）
func (s *BattleScreen) View() string {
	var builder strings.Builder

	// 上部: 敵情報エリア
	enemyArea := s.renderEnemyArea()
	builder.WriteString(enemyArea)
	builder.WriteString("\n")

	// 中央: エージェントエリア / タイピングエリア / 結果表示
	if s.showingResult {
		// 結果表示（WIN/LOSE ASCIIアート）
		resultArea := s.renderResultArea()
		builder.WriteString(resultArea)
	} else if s.activeChallenge != nil {
		// タイピングチャレンジ（ChallengeModelのViewに委譲）
		challengeArea := s.renderChallengeArea()
		builder.WriteString(challengeArea)
	} else {
		// エージェントエリア（3体横並びカード）
		agentArea := s.renderAgentArea()
		builder.WriteString(agentArea)
	}
	builder.WriteString("\n")

	// 下部: プレイヤー情報エリア
	playerArea := s.renderPlayerArea()
	builder.WriteString(playerArea)
	builder.WriteString("\n")

	// メッセージ
	if s.message != "" {
		msgStyle := lipgloss.NewStyle().
			Foreground(styles.ColorInfo).
			Align(lipgloss.Center).
			Width(s.width)
		builder.WriteString(msgStyle.Render(s.message))
		builder.WriteString("\n")
	}

	// ヒント
	hintStyle := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Align(lipgloss.Center).
		Width(s.width)

	var hint string
	if s.showingResult {
		hint = "Enter: 続ける"
	} else if s.activeChallenge != nil {
		hint = "タイピング中...  Esc: キャンセル"
	} else {
		hint = "←/→: エージェント切替  ↑/↓: スキル選択  Enter: 使用  Esc: 中断"
	}
	builder.WriteString(hintStyle.Render(hint))

	return builder.String()
}

// ==================== エリアレンダリング ====================

// renderEnemyArea は敵情報エリアをレンダリングします。
// UI-Improvement Requirement 3.1: 敵情報エリア
func (s *BattleScreen) renderEnemyArea() string {
	var builder strings.Builder

	// ボックス内の利用可能幅（ボックス幅 - パディング左右）
	contentWidth := s.width - 4 - 4 // Width(s.width-4) - Padding(2,2)

	// 1行目: 敵名（左）とボルテージ（右）
	nameStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorDamage)
	leftContent := nameStyle.Render(s.enemy.Name) + fmt.Sprintf(" Lv.%d", s.enemy.Level)
	rightContent := s.styles.RenderVoltage(s.enemy.GetVoltage())

	// 左側を左揃え、右側を右揃えで配置
	leftStyle := lipgloss.NewStyle().Width(contentWidth - lipgloss.Width(rightContent)).Align(lipgloss.Left)
	firstLine := lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(leftContent), rightContent)
	builder.WriteString(firstLine)
	builder.WriteString("\n")

	// パッシブスキル表示（フェーズに応じて通常または強化パッシブを表示）
	passiveStyle := lipgloss.NewStyle().Foreground(styles.ColorBuff)
	if s.enemy.IsEnhanced() {
		phaseStyle := lipgloss.NewStyle().Foreground(styles.ColorDamage).Bold(true)
		builder.WriteString(phaseStyle.Render("[強化フェーズ]"))
		// 強化パッシブを表示
		if s.enemy.Type.EnhancedPassive != nil {
			builder.WriteString("  ")
			builder.WriteString(passiveStyle.Render("★" + s.enemy.Type.EnhancedPassive.Description))
		}
	} else {
		// 通常パッシブを表示
		if s.enemy.Type.NormalPassive != nil {
			builder.WriteString(passiveStyle.Render("★" + s.enemy.Type.NormalPassive.Description))
		}
	}
	builder.WriteString("\n")

	// HP表示（UI改善: アニメーション付きHPバー + フローティングダメージ）
	hpBar := s.enemyHPBar.Render(s.styles, 50)
	displayHP := s.enemyHPBar.GetCurrentHP()
	hpValue := fmt.Sprintf(" %d/%d", displayHP, s.enemy.MaxHP)
	builder.WriteString("HP: ")
	builder.WriteString(hpBar)
	builder.WriteString(hpValue)

	// フローティングダメージ表示
	floatingTexts := s.floatingDamageManager.GetTextsForArea("enemy")
	if len(floatingTexts) > 0 {
		// 最新のフローティングテキストを表示
		text := floatingTexts[0]
		var floatStyle lipgloss.Style
		if text.IsHealing {
			floatStyle = lipgloss.NewStyle().Foreground(styles.ColorHPHigh).Bold(true)
		} else {
			floatStyle = lipgloss.NewStyle().Foreground(styles.ColorDamage).Bold(true)
		}
		builder.WriteString("  ")
		builder.WriteString(floatStyle.Render(text.Text))
	}
	builder.WriteString("\n")

	// 敵のバフ表示
	buffs := s.enemy.EffectTable.FindBySourceType(domain.SourceBuff)
	if len(buffs) > 0 {
		for _, buff := range buffs {
			if buff.Duration != nil {
				builder.WriteString(s.styles.RenderBuff(buff.Name, *buff.Duration))
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}

	// チャージ後行動
	icon, actionText, actionColor := s.getActionDisplay()
	actionStyle := lipgloss.NewStyle().Foreground(actionColor).Bold(true)
	builder.WriteString(actionStyle.Render(fmt.Sprintf("%s %s", icon, actionText)))
	builder.WriteString("\n")

	// 待機状態のプログレスバー（チャージ中/ディフェンス中で計算を分岐）
	now := time.Now()
	var remaining time.Duration
	var ratio float64
	var barColor lipgloss.Color

	switch s.enemy.WaitMode {
	case domain.WaitModeCharging:
		remaining = s.enemy.GetChargeRemainingTime(now)
		ratio = 1.0 - s.enemy.GetChargeProgress(now)
		barColor = styles.ColorSubtle
	case domain.WaitModeDefending:
		remaining = s.enemy.GetDefenseRemainingTime(now)
		// ディフェンス進捗を計算
		if s.enemy.DefenseDuration > 0 {
			elapsed := now.Sub(s.enemy.DefenseStartTime)
			progress := float64(elapsed) / float64(s.enemy.DefenseDuration)
			ratio = 1.0 - progress
		} else {
			ratio = 0
		}
		barColor = styles.ColorInfo // ディフェンス中は青
	default:
		remaining = 0
		ratio = 0
		barColor = styles.ColorSubtle
	}

	if remaining < 0 {
		remaining = 0
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	builder.WriteString(s.renderEnemyActionBar(remaining.Seconds(), ratio, barColor))

	// エリアボックス
	areaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorDamage).
		Padding(1, 2).
		Width(s.width - 4)

	title := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Render("─────────────────────────────────  ENEMY  ─────────────────────────────────")

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		areaStyle.Render(builder.String()),
	)
}

// renderAgentArea はエージェントエリア（3体横並びカード）をレンダリングします。
// タスク 9: リキャスト状態、チェイン効果、パッシブスキル表示を追加
func (s *BattleScreen) renderAgentArea() string {
	// 画面幅に基づいてカード幅を計算（余白を最小限に）
	// 外枠: border(2) + padding(4) = 6
	// カード間スペース: 1 × 2 = 2
	// カードボーダー: 2 × 3 = 6
	// 利用可能幅 = s.width - 4(外枠Width調整) - 6(外枠) - 2(カード間) - 6(カードボーダー) = s.width - 18
	cardWidth := (s.width - 18) / 3
	if cardWidth < 30 {
		cardWidth = 30 // 最小幅を確保
	}
	var cards []string

	for i := 0; i < 3; i++ {
		var cardContent strings.Builder
		isSelected := i == s.selectedAgentIdx

		// リキャスト状態を取得（枠色判定にも使用）
		var recastState *recast.RecastState
		if i < len(s.equippedAgents) {
			recastState = s.recastManager.GetRecastState(i)
		}

		if i < len(s.equippedAgents) {
			agent := s.equippedAgents[i]

			// エージェント名
			nameStyle := lipgloss.NewStyle().Bold(true)
			if isSelected {
				nameStyle = nameStyle.
					Foreground(styles.ColorSelectedFg).
					Background(styles.ColorSelectedBg)
			}
			cardContent.WriteString(nameStyle.Render(agent.GetCoreTypeName()))
			cardContent.WriteString("\n")

			// パッシブスキル表示（コア特性から）- ShortDescriptionを使用
			if agent.Core != nil && agent.Core.PassiveSkill.ID != "" {
				passiveNotification := components.NewPassiveSkillNotification(&agent.Core.PassiveSkill)
				shortDesc := passiveNotification.GetShortDescription()
				passiveStyle := lipgloss.NewStyle().
					Foreground(styles.ColorBuff).
					Bold(true)
				cardContent.WriteString(passiveStyle.Render(fmt.Sprintf("★ %s", shortDesc)))
				cardContent.WriteString("\n")
			}

			// リキャスト状態表示
			if recastState != nil {
				recastBar := components.NewRecastProgressBar()
				recastBar.SetProgress(recastState.RemainingSeconds, recastState.TotalSeconds)
				cardContent.WriteString(lipgloss.NewStyle().Foreground(styles.ColorWarning).Render("⏳ "))
				cardContent.WriteString(recastBar.RenderCompact(10))
				cardContent.WriteString("\n")
			}

			// エージェントのスキル一覧（2行表示）
			// 待機中チェイン効果を取得（発動中の強調表示判定用）
			pendingChain := s.chainEffectManager.GetPendingEffectForAgent(i)

			agentSkills := s.getSkillsForAgent(i)
			for j, slot := range agentSkills {
				// ディフェンススキル未解放時は非表示
				if !s.isDefenseSkillUnlocked() && slot.Skill.GetChallengeType() == domain.ChallengeTypeDefense {
					continue
				}

				isSkillSelected := isSelected && j == s.getSelectedSkillInAgent(i)

				// スキルアイコン
				icon := slot.Skill.Icon()

				// スキル名のスタイル
				var skillStyle lipgloss.Style
				if isSkillSelected {
					skillStyle = lipgloss.NewStyle().
						Bold(true).
						Foreground(styles.ColorSelectedFg).
						Background(styles.ColorSelectedBg)
				} else if !slot.IsReady() || recastState != nil {
					// クールダウン中またはリキャスト中は淡い色
					skillStyle = lipgloss.NewStyle().Foreground(styles.ColorSubtle)
				} else if s.isManaSystemUnlocked() && slot.Skill.Type.ManaCost > 0 && s.player.Mana < slot.Skill.Type.ManaCost {
					// マナ不足時は淡い色（マナシステム解放時のみ判定）
					skillStyle = lipgloss.NewStyle().Foreground(styles.ColorSubtle)
				} else {
					skillStyle = lipgloss.NewStyle().Foreground(styles.ColorSecondary)
				}

				prefix := "  "
				if isSkillSelected {
					prefix = "> "
				}

				// 1行目: プレフィックス + アイコン + スキル名
				cardContent.WriteString(skillStyle.Render(fmt.Sprintf("%s%s %s", prefix, icon, slot.Skill.Name())))
				cardContent.WriteString("\n")

				// 2行目: チェイン効果（あれば）または空行
				if slot.Skill.HasChainEffect() {
					chainBadge := components.NewChainEffectBadge(slot.Skill.ChainEffect)
					// このスキルのチェイン効果が発動中かチェック
					// リキャスト中 かつ 待機中チェイン効果がこのスキルのものなら発動中
					isChainActive := pendingChain != nil &&
						pendingChain.Effect.Type == slot.Skill.ChainEffect.Type
					cardContent.WriteString("    ") // インデント（prefixと同じ幅 + アイコン分）
					if isChainActive {
						cardContent.WriteString(chainBadge.RenderActive())
					} else {
						cardContent.WriteString(chainBadge.RenderWithValue())
					}
					cardContent.WriteString("\n")
				} else {
					// チェイン効果がなくても空行を出力（高さを揃えるため）
					cardContent.WriteString("\n")
				}
			}
		} else {
			// 空スロット
			cardContent.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSubtle).Render("(空)"))
		}

		// カードボックス - リキャスト状態で枠色を変更
		borderColor := styles.ColorSubtle
		if recastState != nil {
			// リキャスト中（クールダウン中）は黄色枠
			borderColor = styles.ColorWarning
		} else if isSelected {
			borderColor = styles.ColorPrimary
		}

		cardStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(0, 1).
			Width(cardWidth).
			Height(10) // 高さを詰める（待機中チェイン効果表示削除分）

		cards = append(cards, cardStyle.Render(cardContent.String()))
	}

	// カードを横に並べる（スペースを最小限に）
	agentCards := lipgloss.JoinHorizontal(lipgloss.Top, cards[0], " ", cards[1], " ", cards[2])

	// タイトル（枠なし）
	title := lipgloss.NewStyle().
		Foreground(styles.ColorSubtle).
		Render("────────────────────────────────  PLAYER  ────────────────────────────────")

	// エリア枠を削除し、カードのみを表示
	areaContent := lipgloss.NewStyle().
		Padding(1, 2).
		Width(s.width - 4).
		Align(lipgloss.Center).
		Render(agentCards)

	return lipgloss.JoinVertical(lipgloss.Center,
		title,
		areaContent,
	)
}

// renderPlayerArea はプレイヤー情報エリアをレンダリングします。
// UI-Improvement Requirement 3.1: プレイヤー情報エリア
func (s *BattleScreen) renderPlayerArea() string {
	var builder strings.Builder

	// マナ表示（マナシステム解放時のみ）
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorHPHigh)
	if s.isManaSystemUnlocked() {
		manaLabel := "Mana:"
		if s.player.Mana > 0 {
			manaStars := strings.Repeat("⭐", s.player.Mana)
			manaLabel += " " + manaStars
		}
		builder.WriteString(titleStyle.Render(manaLabel))
		builder.WriteString("\n")
	}

	// HP表示（UI改善: アニメーション付きHPバー + フローティングダメージ/回復）
	hpBar := s.playerHPBar.Render(s.styles, 50)
	displayHP := s.playerHPBar.GetCurrentHP()
	hpValue := fmt.Sprintf(" %d/%d", displayHP, s.player.MaxHP)
	builder.WriteString("HP: ")
	builder.WriteString(hpBar)
	builder.WriteString(hpValue)

	// フローティングダメージ/回復表示
	floatingTexts := s.floatingDamageManager.GetTextsForArea("player")
	if len(floatingTexts) > 0 {
		// 最新のフローティングテキストを表示
		text := floatingTexts[0]
		var floatStyle lipgloss.Style
		if text.IsHealing {
			floatStyle = lipgloss.NewStyle().Foreground(styles.ColorHPHigh).Bold(true)
		} else {
			floatStyle = lipgloss.NewStyle().Foreground(styles.ColorDamage).Bold(true)
		}
		builder.WriteString("  ")
		builder.WriteString(floatStyle.Render(text.Text))
	}
	builder.WriteString("\n")

	// バフ表示
	buffs := s.player.EffectTable.FindBySourceType(domain.SourceBuff)
	if len(buffs) > 0 {
		builder.WriteString("バフ: ")
		for _, buff := range buffs {
			if buff.Duration != nil {
				builder.WriteString(s.styles.RenderBuff(buff.Name, *buff.Duration))
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}

	// デバフ表示
	debuffs := s.player.EffectTable.FindBySourceType(domain.SourceDebuff)
	if len(debuffs) > 0 {
		builder.WriteString("デバフ: ")
		for _, debuff := range debuffs {
			if debuff.Duration != nil {
				builder.WriteString(s.styles.RenderDebuff(debuff.Name, *debuff.Duration))
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}

	// パッシブスキル表示（スタック型パッシブのダメージ倍率）
	passives := s.player.EffectTable.FindBySourceType(domain.SourcePassive)
	hasStackPassive := false
	for _, passive := range passives {
		if passive.MaxStacks > 0 {
			if !hasStackPassive {
				builder.WriteString("パッシブ: ")
				hasStackPassive = true
			}
			// スタック数を計算（コンボ数を使用、最大スタックでキャップ）
			stacks := s.comboCount
			if stacks > passive.MaxStacks {
				stacks = passive.MaxStacks
			}
			// ダメージ倍率を計算（スタック数 × 効果増分）
			bonusPercent := float64(stacks) * passive.StackIncrement * 100
			builder.WriteString(s.styles.RenderPassive(passive.Name, bonusPercent))
			builder.WriteString(" ")
		}
	}
	if hasStackPassive {
		builder.WriteString("\n")
	}

	// エリアボックス
	areaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorHPHigh).
		Padding(1, 2).
		Width(s.width - 4)

	return areaStyle.Render(builder.String())
}

// renderResultArea は結果表示（WIN/LOSE ASCIIアート）をレンダリングします。
// UI-Improvement Requirement 3.9: WIN/LOSE ASCIIアート表示
func (s *BattleScreen) renderResultArea() string {
	var resultArt string
	if s.victory {
		resultArt = s.winLoseRenderer.RenderWin()
	} else {
		resultArt = s.winLoseRenderer.RenderLose()
	}

	// 中央揃え
	centeredArt := lipgloss.NewStyle().
		Width(s.width - 8).
		Align(lipgloss.Center).
		Render(resultArt)

	// エリアボックス
	areaStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(2, 2).
		Width(s.width - 4)

	return areaStyle.Render(centeredArt)
}

// renderChallengeArea はチャレンジ中のタイピングエリアを描画します。
// ChallengeModelのView()に委譲し、ボックスで囲みます。
func (s *BattleScreen) renderChallengeArea() string {
	if s.activeChallenge == nil {
		return ""
	}

	challengeView := s.activeChallenge.View()

	// チャレンジViewをボックスで囲む
	challengeBox := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.ColorPrimary).
		Padding(1, 2).
		Render(challengeView)

	return lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(s.width).
		Render(challengeBox)
}

// ==================== プログレスバーレンダリング ====================

// renderEnemyActionBar は敵の次回行動までのプログレスバーを描画します。
func (s *BattleScreen) renderEnemyActionBar(remainingSeconds float64, ratio float64, barColor lipgloss.Color) string {
	barWidth := 40
	timeText := fmt.Sprintf("%.1fs", remainingSeconds)

	// 塗りつぶし部分の計算
	filledWidth := int(float64(barWidth) * ratio)
	if filledWidth < 0 {
		filledWidth = 0
	}
	if filledWidth > barWidth {
		filledWidth = barWidth
	}

	// プログレスバー文字列を構築
	filled := strings.Repeat("█", filledWidth)
	empty := strings.Repeat("░", barWidth-filledWidth)
	bar := filled + empty
	barRunes := []rune(bar)

	// テキストの位置を計算
	textStart := (barWidth - len(timeText)) / 2
	if textStart < 0 {
		textStart = 0
	}
	textEnd := textStart + len(timeText)
	if textEnd > barWidth {
		textEnd = barWidth
	}

	// スタイル定義（バー色はパラメータで指定）
	barStyle := lipgloss.NewStyle().Foreground(barColor)
	textStyle := lipgloss.NewStyle().Foreground(styles.ColorSelectedFg).Bold(true)
	bracketStyle := lipgloss.NewStyle().Foreground(barColor)

	// バーを3つの部分に分けてレンダリング（前半バー + テキスト + 後半バー）
	beforeText := string(barRunes[:textStart])
	afterText := string(barRunes[textEnd:])

	return bracketStyle.Render("[") +
		barStyle.Render(beforeText) +
		textStyle.Render(timeText) +
		barStyle.Render(afterText) +
		bracketStyle.Render("]")
}

// ==================== UIヘルパー ====================

// getSkillsForAgent は指定エージェントのスキルスロットを取得します。
func (s *BattleScreen) getSkillsForAgent(agentIdx int) []SkillSlot {
	var skills []SkillSlot
	for _, slot := range s.skillSlots {
		if slot.AgentIndex == agentIdx {
			skills = append(skills, slot)
		}
	}
	return skills
}

// getSelectedSkillInAgent は選択中エージェント内でのスキル選択位置を返します。
func (s *BattleScreen) getSelectedSkillInAgent(agentIdx int) int {
	if s.selectedAgentIdx != agentIdx {
		return -1
	}

	// 現在選択されているスロットがこのエージェントのものか確認
	if s.selectedSlot >= 0 && s.selectedSlot < len(s.skillSlots) {
		slot := s.skillSlots[s.selectedSlot]
		if slot.AgentIndex == agentIdx {
			// このエージェント内での相対位置を計算
			skillIdx := 0
			for i := 0; i < s.selectedSlot; i++ {
				if s.skillSlots[i].AgentIndex == agentIdx {
					skillIdx++
				}
			}
			return skillIdx
		}
	}
	return 0
}
