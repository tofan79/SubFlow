package cost

import "testing"

func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	if a > b {
		return a-b < epsilon
	}
	return b-a < epsilon
}

func TestEstimate_Translation(t *testing.T) {
	t.Run("deepl pricing", func(t *testing.T) {
		est := EstimateTranslation(ProviderDeepL, 500_000)
		if est.Provider != ProviderDeepL {
			t.Fatalf("expected provider deepl, got %q", est.Provider)
		}
		if est.CharCount != 500_000 {
			t.Fatalf("expected char count 500000, got %d", est.CharCount)
		}
		if !almostEqual(est.InputCost, 10.0) || !almostEqual(est.Total, 10.0) {
			t.Fatalf("unexpected deepl cost: %+v", est)
		}
	})

	t.Run("ollama free", func(t *testing.T) {
		est := EstimateTranslation(ProviderOllama, 500_000)
		if est.InputCost != 0 || est.OutputCost != 0 || est.Total != 0 {
			t.Fatalf("expected zero cost, got %+v", est)
		}
	})
}

func TestEstimate_Rewrite(t *testing.T) {
	t.Run("openai pricing", func(t *testing.T) {
		est := EstimateRewrite(ProviderOpenAI, 1_000_000, 2_000_000)
		if !almostEqual(est.InputCost, 2.5) || !almostEqual(est.OutputCost, 20.0) || !almostEqual(est.Total, 22.5) {
			t.Fatalf("unexpected openai cost: %+v", est)
		}
	})

	t.Run("anthropic pricing", func(t *testing.T) {
		est := EstimateRewrite(ProviderAnthropic, 1_000_000, 2_000_000)
		if !almostEqual(est.InputCost, 3.0) || !almostEqual(est.OutputCost, 30.0) || !almostEqual(est.Total, 33.0) {
			t.Fatalf("unexpected anthropic cost: %+v", est)
		}
	})

	t.Run("ollama free", func(t *testing.T) {
		est := EstimateRewrite(ProviderOllama, 123, 456)
		if est.InputCost != 0 || est.OutputCost != 0 || est.Total != 0 {
			t.Fatalf("expected zero cost, got %+v", est)
		}
	})
}

func TestEstimate_ASR(t *testing.T) {
	t.Run("groq pricing", func(t *testing.T) {
		est := EstimateASR(ProviderGroq, 12.5)
		if !almostEqual(est.InputCost, 0.25) || !almostEqual(est.Total, 0.25) {
			t.Fatalf("unexpected groq cost: %+v", est)
		}
	})

	t.Run("deepgram pricing", func(t *testing.T) {
		est := EstimateASR(ProviderDeepgram, 10.0)
		if !almostEqual(est.InputCost, 0.043) || !almostEqual(est.Total, 0.043) {
			t.Fatalf("unexpected deepgram cost: %+v", est)
		}
	})
}

func TestEstimate_OllamaFree(t *testing.T) {
	translation := EstimateTranslation(ProviderOllama, 1_000_000)
	rewrite := EstimateRewrite(ProviderOllama, 1_000_000, 1_000_000)
	asr := EstimateASR(ProviderOllama, 60)

	if translation.Total != 0 || rewrite.Total != 0 || asr.Total != 0 {
		t.Fatalf("expected ollama to be free, got translation=%+v rewrite=%+v asr=%+v", translation, rewrite, asr)
	}
}
