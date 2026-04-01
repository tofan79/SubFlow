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

const geminiDefaultURL = "https://generativelanguage.googleapis.com/v1beta"

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
		baseURL:    geminiDefaultURL,
	}
}

func (p *GeminiProvider) Name() string { return ProviderGemini }

func (p *GeminiProvider) MaxBatchSize() int { return 20 }

func (p *GeminiProvider) EstimateCost(charCount int) float64 {
	if charCount <= 0 {
		return 0
	}
	if strings.Contains(p.model, "flash") {
		return float64(charCount) * 0.0000001
	}
	return float64(charCount) * 0.000002
}

func (p *GeminiProvider) Translate(ctx context.Context, batch []string, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrTrnAPIKeyErr(p.Name())
	}

	prompt := p.buildPrompt(batch, opts)

	reqBody := geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: geminiGenerationConfig{
			Temperature: 0.3,
			MaxTokens:   4096,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to marshal request", err)
	}

	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", p.baseURL, p.model, p.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

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

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to decode response", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "no response from Gemini", nil)
	}

	return p.parseResponse(result.Candidates[0].Content.Parts[0].Text, len(batch))
}

func (p *GeminiProvider) buildPrompt(batch []string, opts Opts) string {
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

	sb.WriteString("\nTranslate these lines:\n\n")
	for i, line := range batch {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, line))
	}

	return sb.String()
}

func (p *GeminiProvider) parseResponse(content string, expectedCount int) ([]string, error) {
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

func (p *GeminiProvider) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return pipeline.ErrTrnRateLimitErr(p.Name(), 60)
	case http.StatusUnauthorized, http.StatusForbidden:
		return pipeline.ErrTrnAPIKeyErr(p.Name())
	default:
		return pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("Gemini API error %d: %s", resp.StatusCode, string(body)), nil)
	}
}

type geminiRequest struct {
	Contents         []geminiContent        `json:"contents"`
	GenerationConfig geminiGenerationConfig `json:"generationConfig,omitempty"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenerationConfig struct {
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"maxOutputTokens,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

var _ Provider = (*GeminiProvider)(nil)
