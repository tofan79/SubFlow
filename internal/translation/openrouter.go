package translation

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

const openRouterDefaultURL = "https://openrouter.ai/api"

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
		baseURL = openRouterDefaultURL
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

func (p *OpenRouterProvider) MaxBatchSize() int { return 20 }

func (p *OpenRouterProvider) EstimateCost(charCount int) float64 {
	if charCount <= 0 {
		return 0
	}
	tokenCount := charCount / 4
	return p.estimateTokenCost(tokenCount)
}

func (p *OpenRouterProvider) estimateTokenCost(tokenCount int) float64 {
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

func (p *OpenRouterProvider) Translate(ctx context.Context, batch []string, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrTrnAPIKeyErr(p.Name())
	}

	systemPrompt := p.buildSystemPrompt(opts)
	userPrompt := p.buildUserPrompt(batch)

	reqBody := openRouterRequest{
		Model: p.model,
		Messages: []openRouterMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to marshal request", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to create request", err)
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
		return nil, pipeline.ErrTrnTimeoutErr(p.Name(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, p.handleError(resp)
	}

	var result openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to decode response", err)
	}

	if len(result.Choices) == 0 {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "no response from OpenRouter", nil)
	}

	return p.parseResponse(result.Choices[0].Message.Content, len(batch))
}

func (p *OpenRouterProvider) buildSystemPrompt(opts Opts) string {
	var sb strings.Builder
	sb.WriteString("You are a professional subtitle translator. ")
	sb.WriteString(fmt.Sprintf("Translate from %s to %s. ", opts.SourceLang, opts.TargetLang))

	if opts.ContentMode != "" {
		sb.WriteString(fmt.Sprintf("Content type: %s. ", opts.ContentMode))
	}

	sb.WriteString("Rules:\n")
	sb.WriteString("- Preserve line numbers exactly as given\n")
	sb.WriteString("- Keep translations natural and conversational\n")
	sb.WriteString("- Maintain the same number of lines\n")
	sb.WriteString("- Output format: [N] translated text\n")

	if len(opts.Glossary) > 0 {
		sb.WriteString("\nGlossary (use these translations):\n")
		for _, term := range opts.Glossary {
			sb.WriteString(fmt.Sprintf("- %s → %s\n", term.SourceTerm, term.TargetTerm))
		}
	}

	return sb.String()
}

func (p *OpenRouterProvider) buildUserPrompt(batch []string) string {
	var sb strings.Builder
	sb.WriteString("Translate these lines:\n\n")
	for i, line := range batch {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, line))
	}
	return sb.String()
}

func (p *OpenRouterProvider) parseResponse(content string, expectedCount int) ([]string, error) {
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

	return result, nil
}

func (p *OpenRouterProvider) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return pipeline.ErrTrnRateLimitErr(p.Name(), 60)
	case http.StatusUnauthorized, http.StatusForbidden:
		return pipeline.ErrTrnAPIKeyErr(p.Name())
	case http.StatusPaymentRequired:
		return pipeline.ErrTrnQuotaExceedErr(p.Name())
	default:
		return pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("OpenRouter API error %d: %s", resp.StatusCode, string(body)), nil)
	}
}

type openRouterRequest struct {
	Model       string              `json:"model"`
	Messages    []openRouterMessage `json:"messages"`
	Temperature float64             `json:"temperature,omitempty"`
}

type openRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

var _ Provider = (*OpenRouterProvider)(nil)
