package translation

import (
	"context"

	"github.com/subflow/subflow/internal/pipeline"
)

type Opts = pipeline.TranslationOpts

type GlossaryTerm = pipeline.GlossaryTerm

type Provider interface {
	Translate(ctx context.Context, batch []string, opts Opts) ([]string, error)
	EstimateCost(charCount int) float64
	MaxBatchSize() int
	Name() string
}

const (
	ProviderDeepL      = "deepl"
	ProviderOpenAI     = "openai"
	ProviderAnthropic  = "anthropic"
	ProviderGemini     = "gemini"
	ProviderOllama     = "ollama"
	ProviderOpenRouter = "openrouter"
)

type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}
