package translation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

const (
	deeplFreeURL = "https://api-free.deepl.com/v2/translate"
	deeplProURL  = "https://api.deepl.com/v2/translate"
)

type DeepLProvider struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
	semaphore  chan struct{}
}

func NewDeepLProvider(apiKey string) *DeepLProvider {
	baseURL := deeplFreeURL
	if strings.HasSuffix(apiKey, ":fx") {
		baseURL = deeplFreeURL
	} else if len(apiKey) > 0 {
		baseURL = deeplProURL
	}

	return &DeepLProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 2 * time.Minute},
		baseURL:    baseURL,
		semaphore:  make(chan struct{}, 10),
	}
}

func (p *DeepLProvider) Name() string { return ProviderDeepL }

func (p *DeepLProvider) MaxBatchSize() int { return 50 }

func (p *DeepLProvider) EstimateCost(charCount int) float64 {
	if charCount <= 0 {
		return 0
	}
	return float64(charCount) * 0.00002
}

func (p *DeepLProvider) Translate(ctx context.Context, batch []string, opts Opts) ([]string, error) {
	if len(batch) == 0 {
		return []string{}, nil
	}

	if strings.TrimSpace(p.apiKey) == "" {
		return nil, pipeline.ErrTrnAPIKeyErr(p.Name())
	}

	select {
	case p.semaphore <- struct{}{}:
		defer func() { <-p.semaphore }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	reqBody := deeplRequest{
		Text:        batch,
		SourceLang:  strings.ToUpper(opts.SourceLang),
		TargetLang:  strings.ToUpper(opts.TargetLang),
		TagHandling: "html",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to marshal request", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to create request", err)
	}

	req.Header.Set("Authorization", "DeepL-Auth-Key "+p.apiKey)
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

	var result deeplResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, "failed to decode response", err)
	}

	translations := make([]string, len(result.Translations))
	for i, t := range result.Translations {
		translations[i] = t.Text
	}

	if len(translations) != len(batch) {
		return nil, pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("expected %d translations, got %d", len(batch), len(translations)), nil)
	}

	return translations, nil
}

func (p *DeepLProvider) handleError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return pipeline.ErrTrnRateLimitErr(p.Name(), 60)
	case http.StatusForbidden, http.StatusUnauthorized:
		return pipeline.ErrTrnAPIKeyErr(p.Name())
	case http.StatusPaymentRequired:
		return pipeline.ErrTrnQuotaExceedErr(p.Name())
	default:
		return pipeline.NewError(pipeline.ErrTrnTimeout, fmt.Sprintf("DeepL API error %d: %s", resp.StatusCode, string(body)), nil)
	}
}

type deeplRequest struct {
	Text        []string `json:"text"`
	SourceLang  string   `json:"source_lang,omitempty"`
	TargetLang  string   `json:"target_lang"`
	TagHandling string   `json:"tag_handling,omitempty"`
}

type deeplResponse struct {
	Translations []struct {
		Text string `json:"text"`
	} `json:"translations"`
}

var _ Provider = (*DeepLProvider)(nil)
