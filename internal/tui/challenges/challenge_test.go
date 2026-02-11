package challenges_test

import (
	"testing"

	"hirorocky/type-battle/internal/domain"
	"hirorocky/type-battle/internal/tui/challenges"

	// サブパッケージのinit()でRegisterを実行
	_ "hirorocky/type-battle/internal/tui/challenges/defense"
	_ "hirorocky/type-battle/internal/tui/challenges/shape"
	_ "hirorocky/type-battle/internal/tui/challenges/standard"
)

func TestNew_スタンダードタイプの生成(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"hello", "world"},
		AutoCorrectCount: 0,
	}

	challenge := challenges.New(domain.ChallengeTypeStandard, input)
	if challenge == nil {
		t.Fatal("ChallengeModelがnilです")
	}

	// 初期状態ではResult()はnilであるべき
	if challenge.Result() != nil {
		t.Error("初期状態でResult()がnon-nilです")
	}
}

func TestNew_未知のtypeIDはスタンダードにフォールバック(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"test"},
		AutoCorrectCount: 0,
	}

	challenge := challenges.New("unknown_type", input)
	if challenge == nil {
		t.Fatal("未知のtypeIDでChallengeModelがnilです")
	}

	// 初期状態ではResult()はnilであるべき
	if challenge.Result() != nil {
		t.Error("初期状態でResult()がnon-nilです")
	}
}

func TestNew_全3タイプが生成可能(t *testing.T) {
	types := []domain.ChallengeTypeID{
		domain.ChallengeTypeStandard,
		domain.ChallengeTypeShape,
		domain.ChallengeTypeDefense,
	}

	for _, typeID := range types {
		t.Run(string(typeID), func(t *testing.T) {
			input := domain.ChallengeInput{
				Difficulty:       domain.DifficultyRateStandard,
				Words:            []string{"test"},
				AutoCorrectCount: 0,
			}

			challenge := challenges.New(typeID, input)
			if challenge == nil {
				t.Fatalf("typeID=%s でChallengeModelがnilです", typeID)
			}
		})
	}
}

func TestChallengeModel_インターフェース準拠(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"test"},
		AutoCorrectCount: 0,
	}

	// ChallengeModel インターフェースを満たすことを型レベルで確認
	var _ challenges.ChallengeModel = challenges.New(domain.ChallengeTypeStandard, input)
	var _ challenges.ChallengeModel = challenges.New(domain.ChallengeTypeShape, input)
	var _ challenges.ChallengeModel = challenges.New(domain.ChallengeTypeDefense, input)
}

func TestDefenseProvider_ディフェンスタイプのみ実装(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"test"},
		AutoCorrectCount: 0,
	}

	// ディフェンスタイプはDefenseProviderを実装
	defenseChallenge := challenges.New(domain.ChallengeTypeDefense, input)
	dp, ok := defenseChallenge.(challenges.DefenseProvider)
	if !ok {
		t.Fatal("ディフェンスタイプがDefenseProviderを実装していません")
	}

	// 初期状態の防御率は0
	if dp.DefenseRate() != 0.0 {
		t.Errorf("初期防御率 = %f, want 0.0", dp.DefenseRate())
	}

	// スタンダードはDefenseProviderを実装しない
	standardChallenge := challenges.New(domain.ChallengeTypeStandard, input)
	if _, ok := standardChallenge.(challenges.DefenseProvider); ok {
		t.Error("スタンダードタイプがDefenseProviderを実装すべきではありません")
	}

	// シェイプもDefenseProviderを実装しない
	shapeChallenge := challenges.New(domain.ChallengeTypeShape, input)
	if _, ok := shapeChallenge.(challenges.DefenseProvider); ok {
		t.Error("シェイプタイプがDefenseProviderを実装すべきではありません")
	}
}

func TestDefenseProvider_CompleteByAttack(t *testing.T) {
	input := domain.ChallengeInput{
		Difficulty:       domain.DifficultyRateStandard,
		Words:            []string{"test"},
		AutoCorrectCount: 0,
	}

	challenge := challenges.New(domain.ChallengeTypeDefense, input)
	dp := challenge.(challenges.DefenseProvider)

	// 完了前はResult()がnil
	if challenge.Result() != nil {
		t.Fatal("CompleteByAttack前にResult()がnon-nilです")
	}

	// 攻撃で自動終了
	dp.CompleteByAttack()

	result := challenge.Result()
	if result == nil {
		t.Fatal("CompleteByAttack後にResult()がnilです")
	}
	if result.Status != domain.ChallengeSuccess {
		t.Errorf("Status = %d, want ChallengeSuccess(%d)", result.Status, domain.ChallengeSuccess)
	}

	// 二重呼び出しはno-op
	dp.CompleteByAttack()
	if challenge.Result() != result {
		t.Error("CompleteByAttackの二重呼び出しでResultが変化しました")
	}
}
