// Package config は設定定数を提供します。
package config

import (
	"testing"
	"time"
)

// TestBattleConstants はバトル設定定数のテストです。
func TestBattleConstants(t *testing.T) {
	t.Run("BattleTickIntervalが100msである", func(t *testing.T) {
		expected := 100 * time.Millisecond
		if BattleTickInterval != expected {
			t.Errorf("BattleTickIntervalが期待値と異なります: got %v, want %v", BattleTickInterval, expected)
		}
	})

	t.Run("DefaultModuleCooldownが5.0である", func(t *testing.T) {
		expected := 5.0
		if DefaultModuleCooldown != expected {
			t.Errorf("DefaultModuleCooldownが期待値と異なります: got %f, want %f", DefaultModuleCooldown, expected)
		}
	})

	t.Run("AccuracyPenaltyThresholdが0.5である", func(t *testing.T) {
		expected := 0.5
		if AccuracyPenaltyThreshold != expected {
			t.Errorf("AccuracyPenaltyThresholdが期待値と異なります: got %f, want %f", AccuracyPenaltyThreshold, expected)
		}
	})

	t.Run("MinEnemyAttackIntervalが500msである", func(t *testing.T) {
		expected := 500 * time.Millisecond
		if MinEnemyAttackInterval != expected {
			t.Errorf("MinEnemyAttackIntervalが期待値と異なります: got %v, want %v", MinEnemyAttackInterval, expected)
		}
	})
}

// TestEffectDurationConstants は効果持続時間定数のテストです。
func TestEffectDurationConstants(t *testing.T) {
	t.Run("BuffDurationが10.0である", func(t *testing.T) {
		expected := 10.0
		if BuffDuration != expected {
			t.Errorf("BuffDurationが期待値と異なります: got %f, want %f", BuffDuration, expected)
		}
	})

	t.Run("DebuffDurationが8.0である", func(t *testing.T) {
		expected := 8.0
		if DebuffDuration != expected {
			t.Errorf("DebuffDurationが期待値と異なります: got %f, want %f", DebuffDuration, expected)
		}
	})
}

// TestInventoryConstants はインベントリ定数のテストです。
func TestInventoryConstants(t *testing.T) {
	t.Run("MaxAgentEquipSlotsが3である", func(t *testing.T) {
		expected := 3
		if MaxAgentEquipSlots != expected {
			t.Errorf("MaxAgentEquipSlotsが期待値と異なります: got %d, want %d", MaxAgentEquipSlots, expected)
		}
	})

	t.Run("ModulesPerAgentが4である", func(t *testing.T) {
		expected := 4
		if ModulesPerAgent != expected {
			t.Errorf("ModulesPerAgentが期待値と異なります: got %d, want %d", ModulesPerAgent, expected)
		}
	})
}
