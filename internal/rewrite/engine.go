// Package rewrite provides Layer 2 AI rewrite functionality for subtitle enhancement.
package rewrite

import (
	"context"
	"fmt"
	"strings"

	"github.com/subflow/subflow/internal/pipeline"
)

// Input is an alias for pipeline.RewriteInput
type Input = pipeline.RewriteInput

// Opts is an alias for pipeline.RewriteOpts
type Opts = pipeline.RewriteOpts

// GlossaryTerm is an alias for pipeline.GlossaryTerm
type GlossaryTerm = pipeline.GlossaryTerm

// TonePreset is an alias for pipeline.TonePreset
type TonePreset = pipeline.TonePreset

const (
	ToneNatural   = pipeline.ToneNatural
	ToneFormal    = pipeline.ToneFormal
	ToneCasual    = pipeline.ToneCasual
	ToneCinematic = pipeline.ToneCinematic
)

// Provider constants
const (
	ProviderOpenAI     = "openai"
	ProviderAnthropic  = "anthropic"
	ProviderGemini     = "gemini"
	ProviderQwen       = "qwen"
	ProviderXAI        = "xai"
	ProviderOllama     = "ollama"
	ProviderOpenRouter = "openrouter"
)

// Provider defines the interface for rewrite services.
type Provider interface {
	Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error)
	EstimateCost(tokenCount int) float64
	Name() string
}

// Config holds configuration for creating a rewrite provider.
type Config struct {
	Provider string
	APIKey   string
	Model    string
	BaseURL  string
}

// Engine orchestrates Layer 2 rewrite operations.
type Engine struct {
	provider Provider
}

// NewEngine creates a new rewrite engine with the specified provider.
func NewEngine(cfg Config) (*Engine, error) {
	p, err := NewProvider(cfg)
	if err != nil {
		return nil, err
	}
	return &Engine{provider: p}, nil
}

// NewProvider creates a new rewrite provider based on configuration.
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case ProviderOpenAI:
		return NewOpenAIProvider(cfg.APIKey, cfg.Model, cfg.BaseURL), nil
	case ProviderAnthropic:
		return NewAnthropicProvider(cfg.APIKey, cfg.Model), nil
	case ProviderGemini:
		return NewGeminiProvider(cfg.APIKey, cfg.Model), nil
	case ProviderQwen:
		return NewQwenProvider(cfg.APIKey, cfg.Model, cfg.BaseURL), nil
	case ProviderXAI:
		return NewXAIProvider(cfg.APIKey, cfg.Model), nil
	case ProviderOllama:
		return NewOllamaProvider(cfg.BaseURL, cfg.Model), nil
	case ProviderOpenRouter:
		return NewOpenRouterProvider(cfg.APIKey, cfg.Model, cfg.BaseURL), nil
	default:
		return nil, fmt.Errorf("unknown rewrite provider: %s", cfg.Provider)
	}
}

// Rewrite processes a batch of inputs through the rewrite provider.
func (e *Engine) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if opts.MaxCharsPerLine <= 0 {
		opts.MaxCharsPerLine = 42
	}
	if opts.MaxLines <= 0 {
		opts.MaxLines = 2
	}
	if opts.MaxCPS <= 0 {
		opts.MaxCPS = 17.0
	}
	if opts.TonePreset == "" {
		opts.TonePreset = string(ToneNatural)
	}

	return e.provider.Rewrite(ctx, batch, opts)
}

// EstimateCost estimates the cost for rewriting the given token count.
func (e *Engine) EstimateCost(tokenCount int) float64 {
	return e.provider.EstimateCost(tokenCount)
}

// ProviderName returns the name of the underlying provider.
func (e *Engine) ProviderName() string {
	return e.provider.Name()
}

// BuildSystemPrompt constructs the system prompt for rewrite operations.
func BuildSystemPrompt(opts Opts) string {
	var sb strings.Builder

	sb.WriteString("You are a professional subtitle editor. ")
	sb.WriteString("Your task is to rewrite translations to be more natural and readable. ")

	switch TonePreset(opts.TonePreset) {
	case ToneNatural:
		sb.WriteString("\nTone: Natural - balanced, conversational, sounds like real dialogue. ")
	case ToneFormal:
		sb.WriteString("\nTone: Formal - proper grammar, complete sentences, follows EYD rules. ")
	case ToneCasual:
		sb.WriteString("\nTone: Casual - relaxed, may use informal pronouns (gue/lo), colloquial expressions. ")
	case ToneCinematic:
		sb.WriteString("\nTone: Cinematic - dramatic, strong word choice, emphasize emotions. ")
	default:
		sb.WriteString("\nTone: Natural - balanced, conversational, sounds like real dialogue. ")
	}

	sb.WriteString("\n\nRules:\n")
	sb.WriteString(fmt.Sprintf("- Maximum %d characters per line\n", opts.MaxCharsPerLine))
	sb.WriteString(fmt.Sprintf("- Maximum %d lines per subtitle\n", opts.MaxLines))
	sb.WriteString(fmt.Sprintf("- Target reading speed: %.1f characters per second or less\n", opts.MaxCPS))
	sb.WriteString("- Preserve the original meaning\n")
	sb.WriteString("- Keep speaker's personality and emotion\n")
	sb.WriteString("- Use natural line breaks for readability\n")
	sb.WriteString("- Output format: [N] rewritten text (preserve line numbers)\n")

	if len(opts.Glossary) > 0 {
		sb.WriteString("\nGlossary (must use these translations):\n")
		for _, term := range opts.Glossary {
			sb.WriteString(fmt.Sprintf("- %s → %s\n", term.SourceTerm, term.TargetTerm))
		}
	}

	return sb.String()
}

// BuildUserPrompt constructs the user prompt for rewrite operations.
func BuildUserPrompt(batch []Input) string {
	var sb strings.Builder

	sb.WriteString("Rewrite these translations to be more natural:\n\n")

	for i, input := range batch {
		sb.WriteString(fmt.Sprintf("[%d]\n", i+1))
		sb.WriteString(fmt.Sprintf("Source: %s\n", input.Source))
		sb.WriteString(fmt.Sprintf("Translation: %s\n", input.Translated))

		if input.Speaker != "" {
			sb.WriteString(fmt.Sprintf("Speaker: %s\n", input.Speaker))
		}
		if input.Emotion != "" {
			sb.WriteString(fmt.Sprintf("Emotion: %s\n", input.Emotion))
		}
		if input.Context != "" {
			sb.WriteString(fmt.Sprintf("Context: %s\n", input.Context))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// EstimateTokens estimates the token count for a batch of inputs.
func EstimateTokens(batch []Input) int {
	const charsPerToken = 4 // typical ratio for English/Indonesian mixed text
	totalChars := 0
	for _, input := range batch {
		totalChars += len(input.Source)
		totalChars += len(input.Translated)
		totalChars += len(input.Speaker)
		totalChars += len(input.Emotion)
		totalChars += len(input.Context)
	}
	return totalChars / charsPerToken
}
