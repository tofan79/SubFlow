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

const ollamaDefaultURL = "http://localhost:11434"

type OllamaProvider struct {
	model      string
	httpClient *http.Client
	baseURL    string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	if baseURL == "" {
		baseURL = ollamaDefaultURL
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

func (p *OllamaProvider) MaxBatchSize() int { return 10 }

func (p *OllamaProvider) EstimateCost(charCount int) float64 {
	return 0.0
}

func (p *OllamaProvider) Translate(ctx context.Context, batch []string, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	prompt := p.buildPrompt(batch, opts)

	reqBody := ollamaRequest{
		Model: p.model,
		Messages: []ollamaMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to marshal request", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL+"/api/chat", bytes.NewReader(bodyBytes))
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

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to decode response", err)
	}

	if result.Message.Content == "" {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "no response from Ollama", nil)
	}

	return p.parseResponse(result.Message.Content, len(batch))
}

func (p *OllamaProvider) buildPrompt(batch []string, opts Opts) string {
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

func (p *OllamaProvider) parseResponse(content string, expectedCount int) ([]string, error) {
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

func (p *OllamaProvider) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusNotFound:
		return pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("Ollama model not found: %s", p.model), nil)
	default:
		return pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("Ollama API error %d: %s", resp.StatusCode, string(body)), nil)
	}
}

type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

var _ Provider = (*OllamaProvider)(nil)
