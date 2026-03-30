package cost

import "unicode/utf8"

type Provider string

const (
	ProviderDeepL     Provider = "deepl"
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderGemini    Provider = "gemini"
	ProviderOllama    Provider = "ollama"
	ProviderGroq      Provider = "groq"
	ProviderDeepgram  Provider = "deepgram"
)

type Estimate struct {
	Provider     Provider
	InputCost    float64
	OutputCost   float64
	Total        float64
	TokensIn     int
	TokensOut    int
	CharCount    int
	AudioMinutes float64
}

const (
	deepLPricePerMillionChars    = 20.0
	openAIInputPerMillionTokens  = 2.50
	openAIOutputPerMillionTokens = 10.0
	anthropicInputPerMillion     = 3.0
	anthropicOutputPerMillion    = 15.0
	geminiInputPerMillion        = 0.075
	geminiOutputPerMillion       = 0.30
	groqWhisperPerMinute         = 0.02
	deepgramPerMinute            = 0.0043
)

func EstimateTranslation(provider Provider, charCount int) Estimate {
	est := Estimate{Provider: provider, CharCount: charCount}
	if provider == ProviderOllama || charCount <= 0 {
		return est
	}

	if provider == ProviderDeepL {
		est.InputCost = float64(charCount) / 1_000_000.0 * deepLPricePerMillionChars
		est.Total = est.InputCost
	}
	return est
}

func EstimateRewrite(provider Provider, tokensIn, tokensOut int) Estimate {
	est := Estimate{Provider: provider, TokensIn: tokensIn, TokensOut: tokensOut}
	if provider == ProviderOllama || (tokensIn <= 0 && tokensOut <= 0) {
		return est
	}

	switch provider {
	case ProviderOpenAI:
		est.InputCost = float64(tokensIn) / 1_000_000.0 * openAIInputPerMillionTokens
		est.OutputCost = float64(tokensOut) / 1_000_000.0 * openAIOutputPerMillionTokens
	case ProviderAnthropic:
		est.InputCost = float64(tokensIn) / 1_000_000.0 * anthropicInputPerMillion
		est.OutputCost = float64(tokensOut) / 1_000_000.0 * anthropicOutputPerMillion
	case ProviderGemini:
		est.InputCost = float64(tokensIn) / 1_000_000.0 * geminiInputPerMillion
		est.OutputCost = float64(tokensOut) / 1_000_000.0 * geminiOutputPerMillion
	}
	est.Total = est.InputCost + est.OutputCost
	return est
}

func EstimateASR(provider Provider, audioMinutes float64) Estimate {
	est := Estimate{Provider: provider, AudioMinutes: audioMinutes}
	if provider == ProviderOllama || audioMinutes <= 0 {
		return est
	}

	switch provider {
	case ProviderGroq:
		est.InputCost = audioMinutes * groqWhisperPerMinute
	case ProviderDeepgram:
		est.InputCost = audioMinutes * deepgramPerMinute
	}
	est.Total = est.InputCost
	return est
}

func CountTokens(text string) int {
	if text == "" {
		return 0
	}
	return utf8.RuneCountInString(text) / 4
}
