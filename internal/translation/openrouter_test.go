package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOpenRouterProvider_Name(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.Name() != ProviderOpenRouter {
		t.Errorf("expected %s, got %s", ProviderOpenRouter, p.Name())
	}
}

func TestOpenRouterProvider_MaxBatchSize(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.MaxBatchSize() != 20 {
		t.Errorf("expected 20, got %d", p.MaxBatchSize())
	}
}

func TestOpenRouterProvider_EstimateCost(t *testing.T) {
	tests := []struct {
		model   string
		chars   int
		minCost float64
		maxCost float64
	}{
		{"openai/gpt-4o-mini", 4000, 0.00014, 0.00016},
		{"openai/gpt-4o", 4000, 0.0049, 0.0051},
		{"anthropic/claude-3-haiku", 4000, 0.00024, 0.00026},
		{"meta-llama/llama-3.1-70b-instruct", 4000, 0.00019, 0.00021},
		{"unknown/model", 4000, 0.0009, 0.0011},
	}

	for _, tc := range tests {
		p := NewOpenRouterProvider("test-key", tc.model, "")
		cost := p.EstimateCost(tc.chars)
		if cost < tc.minCost || cost > tc.maxCost {
			t.Errorf("model %s: expected cost between %f and %f, got %f", tc.model, tc.minCost, tc.maxCost, cost)
		}
	}

	p := NewOpenRouterProvider("test-key", "", "")
	if p.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}
}

func TestOpenRouterProvider_DefaultModel(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.model != "openai/gpt-4o-mini" {
		t.Errorf("expected default model openai/gpt-4o-mini, got %s", p.model)
	}
}

func TestOpenRouterProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}
		if r.Header.Get("HTTP-Referer") == "" {
			t.Error("missing HTTP-Referer header")
		}
		if r.Header.Get("X-Title") == "" {
			t.Error("missing X-Title header")
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}

		var req openRouterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(req.Messages))
		}

		resp := openRouterResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "[1] Halo\n[2] Dunia"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "openai/gpt-4o-mini", server.URL)

	result, err := p.Translate(context.Background(), []string{"Hello", "World"}, Opts{
		SourceLang: "EN",
		TargetLang: "ID",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "Halo" || result[1] != "Dunia" {
		t.Errorf("unexpected translations: %v", result)
	}
}

func TestOpenRouterProvider_EmptyBatch(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestOpenRouterProvider_EmptyAPIKey(t *testing.T) {
	p := NewOpenRouterProvider("", "", "")
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestOpenRouterProvider_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestOpenRouterProvider_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("bad-key", "", server.URL)

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}

func TestOpenRouterProvider_PaymentRequired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected payment required error")
	}
}

func TestOpenRouterProvider_GlossaryInPrompt(t *testing.T) {
	var capturedReq openRouterRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := openRouterResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "[1] Test"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{
		SourceLang: "EN",
		TargetLang: "ID",
		Glossary: []GlossaryTerm{
			{SourceTerm: "Hello", TargetTerm: "Halo"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	systemPrompt := capturedReq.Messages[0].Content
	if !strings.Contains(systemPrompt, "Hello → Halo") {
		t.Error("glossary not found in system prompt")
	}
}

func TestOpenRouterProvider_RequiredHeaders(t *testing.T) {
	var capturedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		resp := openRouterResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "[1] Test"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedHeaders.Get("HTTP-Referer") != "https://subflow.app" {
		t.Errorf("expected HTTP-Referer https://subflow.app, got %s", capturedHeaders.Get("HTTP-Referer"))
	}
	if capturedHeaders.Get("X-Title") != "SubFlow" {
		t.Errorf("expected X-Title SubFlow, got %s", capturedHeaders.Get("X-Title"))
	}
}
