package rewrite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

type OpenAIProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewOpenAIProvider(apiKey, model, baseURL string) *OpenAIProvider {
	if model == "" {
		model = "gpt-4o-mini"
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	return &OpenAIProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (p *OpenAIProvider) Name() string { return ProviderOpenAI }

func (p *OpenAIProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	if strings.Contains(p.model, "gpt-4o-mini") {
		return float64(tokenCount) * 0.00000015
	}
	return float64(tokenCount) * 0.000005
}

func (p *OpenAIProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": BuildSystemPrompt(opts)},
			{"role": "user", "content": BuildUserPrompt(batch)},
		},
		"temperature": 0.4,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Choices) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from OpenAI", nil)
	}

	return parseRewriteResponse(result.Choices[0].Message.Content, len(batch))
}

var _ Provider = (*OpenAIProvider)(nil)

type AnthropicProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	if model == "" {
		model = "claude-3-haiku-20240307"
	}
	return &AnthropicProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    "https://api.anthropic.com",
	}
}

func (p *AnthropicProvider) Name() string { return ProviderAnthropic }

func (p *AnthropicProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	if strings.Contains(p.model, "haiku") {
		return float64(tokenCount) * 0.00000025
	}
	return float64(tokenCount) * 0.000003
}

func (p *AnthropicProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	reqBody := map[string]any{
		"model":      p.model,
		"max_tokens": 4096,
		"system":     BuildSystemPrompt(opts),
		"messages": []map[string]string{
			{"role": "user", "content": BuildUserPrompt(batch)},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Content) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from Anthropic", nil)
	}

	return parseRewriteResponse(result.Content[0].Text, len(batch))
}

var _ Provider = (*AnthropicProvider)(nil)

type GeminiProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewGeminiProvider(apiKey, model string) *GeminiProvider {
	if model == "" {
		model = "gemini-1.5-flash"
	}
	return &GeminiProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    "https://generativelanguage.googleapis.com/v1beta",
	}
}

func (p *GeminiProvider) Name() string { return ProviderGemini }

func (p *GeminiProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	if strings.Contains(p.model, "flash") {
		return float64(tokenCount) * 0.000000075
	}
	return float64(tokenCount) * 0.00000125
}

func (p *GeminiProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	prompt := BuildSystemPrompt(opts) + "\n\n" + BuildUserPrompt(batch)
	reqBody := map[string]any{
		"contents": []map[string]any{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string]any{
			"temperature":     0.4,
			"maxOutputTokens": 4096,
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, p.model, p.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from Gemini", nil)
	}

	return parseRewriteResponse(result.Candidates[0].Content.Parts[0].Text, len(batch))
}

var _ Provider = (*GeminiProvider)(nil)

type QwenProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewQwenProvider(apiKey, model, baseURL string) *QwenProvider {
	if model == "" {
		model = "qwen-turbo"
	}
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode"
	}
	return &QwenProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (p *QwenProvider) Name() string { return ProviderQwen }

func (p *QwenProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	return float64(tokenCount) * 0.0000008
}

func (p *QwenProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": BuildSystemPrompt(opts)},
			{"role": "user", "content": BuildUserPrompt(batch)},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Choices) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from Qwen", nil)
	}

	return parseRewriteResponse(result.Choices[0].Message.Content, len(batch))
}

var _ Provider = (*QwenProvider)(nil)

type XAIProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewXAIProvider(apiKey, model string) *XAIProvider {
	if model == "" {
		model = "grok-beta"
	}
	return &XAIProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    "https://api.x.ai",
	}
}

func (p *XAIProvider) Name() string { return ProviderXAI }

func (p *XAIProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	return float64(tokenCount) * 0.000005
}

func (p *XAIProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": BuildSystemPrompt(opts)},
			{"role": "user", "content": BuildUserPrompt(batch)},
		},
		"temperature": 0.4,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Choices) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from xAI", nil)
	}

	return parseRewriteResponse(result.Choices[0].Message.Content, len(batch))
}

var _ Provider = (*XAIProvider)(nil)

type OllamaProvider struct {
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}
	return &OllamaProvider{
		model:      model,
		httpClient: &http.Client{Timeout: 10 * time.Minute},
		baseURL:    strings.TrimRight(baseURL, "/"),
	}
}

func (p *OllamaProvider) Name() string { return ProviderOllama }

func (p *OllamaProvider) EstimateCost(tokenCount int) float64 {
	return 0.0
}

func (p *OllamaProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	prompt := BuildSystemPrompt(opts) + "\n\n" + BuildUserPrompt(batch)
	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"stream": false,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if result.Message.Content == "" {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from Ollama", nil)
	}

	return parseRewriteResponse(result.Message.Content, len(batch))
}

var _ Provider = (*OllamaProvider)(nil)

type OpenRouterProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
	baseURL    string
	siteURL    string
	siteName   string
}

func NewOpenRouterProvider(apiKey, model, baseURL string) *OpenRouterProvider {
	if model == "" {
		model = "openai/gpt-4o-mini"
	}
	if baseURL == "" {
		baseURL = "https://openrouter.ai/api"
	}
	return &OpenRouterProvider{
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: 5 * time.Minute},
		baseURL:    strings.TrimRight(baseURL, "/"),
		siteURL:    "https://subflow.app",
		siteName:   "SubFlow",
	}
}

func (p *OpenRouterProvider) Name() string { return ProviderOpenRouter }

func (p *OpenRouterProvider) EstimateCost(tokenCount int) float64 {
	if tokenCount <= 0 {
		return 0
	}
	switch {
	case strings.Contains(p.model, "gpt-4o-mini"):
		return float64(tokenCount) * 0.00000015
	case strings.Contains(p.model, "gpt-4o"):
		return float64(tokenCount) * 0.000005
	case strings.Contains(p.model, "claude-3-haiku"):
		return float64(tokenCount) * 0.00000025
	case strings.Contains(p.model, "claude-3-sonnet"):
		return float64(tokenCount) * 0.000003
	case strings.Contains(p.model, "claude-3-opus"):
		return float64(tokenCount) * 0.000015
	case strings.Contains(p.model, "gemini-1.5-flash"):
		return float64(tokenCount) * 0.000000075
	case strings.Contains(p.model, "gemini-1.5-pro"):
		return float64(tokenCount) * 0.00000125
	case strings.Contains(p.model, "llama"):
		return float64(tokenCount) * 0.0000002
	case strings.Contains(p.model, "mistral"):
		return float64(tokenCount) * 0.0000002
	default:
		return float64(tokenCount) * 0.000001
	}
}

func (p *OpenRouterProvider) Rewrite(ctx context.Context, batch []Input, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrRwtAPIKeyErr(p.Name())
	}

	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "system", "content": BuildSystemPrompt(opts)},
			{"role": "user", "content": BuildUserPrompt(batch)},
		},
		"temperature": 0.4,
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", p.siteURL)
	req.Header.Set("X-Title", p.siteName)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, pipeline.ErrRwtTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleHTTPError(resp, p.Name())
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "failed to decode response", err)
	}

	if len(result.Choices) == 0 {
		return nil, pipeline.NewError(pipeline.ErrRwtTimeout, "no response from OpenRouter", nil)
	}

	return parseRewriteResponse(result.Choices[0].Message.Content, len(batch))
}

var _ Provider = (*OpenRouterProvider)(nil)

func parseRewriteResponse(content string, expectedCount int) ([]string, error) {
	lines := strings.Split(content, "\n")
	result := make([]string, expectedCount)

	linePattern := regexp.MustCompile(`^\[(\d+)\]\s*(.*)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := linePattern.FindStringSubmatch(line)
		if len(matches) == 3 {
			idx, err := strconv.Atoi(matches[1])
			if err == nil && idx >= 1 && idx <= expectedCount {
				result[idx-1] = matches[2]
			}
		}
	}

	for i, r := range result {
		if r == "" {
			result[i] = ""
		}
	}

	return result, nil
}

func handleHTTPError(resp *http.Response, providerName string) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return pipeline.ErrRwtRateLimitErr(providerName)
	case http.StatusUnauthorized, http.StatusForbidden:
		return pipeline.ErrRwtAPIKeyErr(providerName)
	default:
		return pipeline.NewError(pipeline.ErrRwtTimeout, fmt.Sprintf("%s API error %d: %s", providerName, resp.StatusCode, string(body)), nil)
	}
}
