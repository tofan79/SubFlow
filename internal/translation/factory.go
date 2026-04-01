package translation

import (
	"fmt"
)

func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case ProviderDeepL:
		return NewDeepLProvider(cfg.APIKey), nil
	case ProviderOpenAI:
		return NewOpenAIProvider(cfg.APIKey, cfg.Model, cfg.BaseURL), nil
	case ProviderAnthropic:
		return NewAnthropicProvider(cfg.APIKey, cfg.Model), nil
	case ProviderGemini:
		return NewGeminiProvider(cfg.APIKey, cfg.Model), nil
	case ProviderOllama:
		return NewOllamaProvider(cfg.BaseURL, cfg.Model), nil
	case ProviderOpenRouter:
		return NewOpenRouterProvider(cfg.APIKey, cfg.Model, cfg.BaseURL), nil
	default:
		return nil, fmt.Errorf("unknown translation provider: %s", cfg.Provider)
	}
}
