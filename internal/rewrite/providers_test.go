package rewrite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProvider_Name(t *testing.T) {
	p := NewOpenAIProvider("test-key", "", "")
	if p.Name() != ProviderOpenAI {
		t.Errorf("expected %s, got %s", ProviderOpenAI, p.Name())
	}
}

func TestOpenAIProvider_DefaultModel(t *testing.T) {
	p := NewOpenAIProvider("test-key", "", "")
	if p.model != "gpt-4o-mini" {
		t.Errorf("expected default model gpt-4o-mini, got %s", p.model)
	}
}

func TestOpenAIProvider_EstimateCost(t *testing.T) {
	p := NewOpenAIProvider("test-key", "gpt-4o-mini", "")
	cost := p.EstimateCost(1000)
	if cost < 0.0001 || cost > 0.0002 {
		t.Errorf("unexpected cost for gpt-4o-mini: %f", cost)
	}
}

func TestOpenAIProvider_Rewrite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or wrong Authorization header")
		}
		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": "[1] Halo\n[2] Dunia"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenAIProvider("test-key", "gpt-4o-mini", server.URL)
	result, err := p.Rewrite(context.Background(), []Input{
		{Source: "Hello", Translated: "Halo"},
		{Source: "World", Translated: "Dunia"},
	}, Opts{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
}

func TestOpenAIProvider_EmptyAPIKey(t *testing.T) {
	p := NewOpenAIProvider("", "", "")
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestAnthropicProvider_Name(t *testing.T) {
	p := NewAnthropicProvider("test-key", "")
	if p.Name() != ProviderAnthropic {
		t.Errorf("expected %s, got %s", ProviderAnthropic, p.Name())
	}
}

func TestAnthropicProvider_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Error("missing x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Error("missing anthropic-version header")
		}
		resp := map[string]any{
			"content": []map[string]string{{"text": "[1] Test"}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewAnthropicProvider("test-key", "")
	p.baseURL = server.URL
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGeminiProvider_Name(t *testing.T) {
	p := NewGeminiProvider("test-key", "")
	if p.Name() != ProviderGemini {
		t.Errorf("expected %s, got %s", ProviderGemini, p.Name())
	}
}

func TestGeminiProvider_APIKeyInURL(t *testing.T) {
	var requestedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedURL = r.URL.String()
		resp := map[string]any{
			"candidates": []map[string]any{
				{"content": map[string]any{"parts": []map[string]string{{"text": "[1] Test"}}}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewGeminiProvider("my-api-key", "gemini-1.5-flash")
	p.baseURL = server.URL
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if requestedURL == "" || len(requestedURL) < 10 {
		t.Error("expected URL with API key")
	}
}

func TestQwenProvider_Name(t *testing.T) {
	p := NewQwenProvider("test-key", "", "")
	if p.Name() != ProviderQwen {
		t.Errorf("expected %s, got %s", ProviderQwen, p.Name())
	}
}

func TestXAIProvider_Name(t *testing.T) {
	p := NewXAIProvider("test-key", "")
	if p.Name() != ProviderXAI {
		t.Errorf("expected %s, got %s", ProviderXAI, p.Name())
	}
}

func TestOllamaProvider_Name(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.Name() != ProviderOllama {
		t.Errorf("expected %s, got %s", ProviderOllama, p.Name())
	}
}

func TestOllamaProvider_ZeroCost(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.EstimateCost(1000) != 0 {
		t.Error("expected 0 cost for Ollama")
	}
	if p.EstimateCost(0) != 0 {
		t.Error("expected 0 cost for 0 tokens")
	}
}

func TestOllamaProvider_NoAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("Ollama should not have Authorization header")
		}
		resp := map[string]any{
			"message": map[string]string{"content": "[1] Test"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOllamaProvider_Rewrite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		resp := map[string]any{
			"message": map[string]string{"content": "[1] Halo\n[2] Dunia"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")
	result, err := p.Rewrite(context.Background(), []Input{
		{Source: "Hello", Translated: "Halo"},
		{Source: "World", Translated: "Dunia"},
	}, Opts{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "Halo" || result[1] != "Dunia" {
		t.Errorf("unexpected results: %v", result)
	}
}

func TestParseRewriteResponse(t *testing.T) {
	content := "[1] First line\n[2] Second line\n\n[3] Third line"
	result, err := parseRewriteResponse(content, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 results, got %d", len(result))
	}
	if result[0] != "First line" {
		t.Errorf("unexpected first: %s", result[0])
	}
	if result[1] != "Second line" {
		t.Errorf("unexpected second: %s", result[1])
	}
	if result[2] != "Third line" {
		t.Errorf("unexpected third: %s", result[2])
	}
}

func TestParseRewriteResponse_OutOfOrder(t *testing.T) {
	content := "[3] Third\n[1] First\n[2] Second"
	result, err := parseRewriteResponse(content, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0] != "First" || result[1] != "Second" || result[2] != "Third" {
		t.Errorf("unexpected results: %v", result)
	}
}

func TestHandleHTTPError_RateLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "test")
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestHandleHTTPError_Auth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewOpenAIProvider("bad-key", "", server.URL)
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}

func TestOpenRouterProvider_Name(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.Name() != ProviderOpenRouter {
		t.Errorf("expected %s, got %s", ProviderOpenRouter, p.Name())
	}
}

func TestOpenRouterProvider_DefaultModel(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.model != "openai/gpt-4o-mini" {
		t.Errorf("expected default model openai/gpt-4o-mini, got %s", p.model)
	}
}

func TestOpenRouterProvider_EstimateCost(t *testing.T) {
	tests := []struct {
		model   string
		tokens  int
		minCost float64
		maxCost float64
	}{
		{"openai/gpt-4o-mini", 1000, 0.0001, 0.0002},
		{"openai/gpt-4o", 1000, 0.004, 0.006},
		{"anthropic/claude-3-haiku", 1000, 0.0002, 0.0003},
		{"anthropic/claude-3-sonnet", 1000, 0.002, 0.004},
		{"anthropic/claude-3-opus", 1000, 0.01, 0.02},
		{"google/gemini-1.5-flash", 1000, 0.00005, 0.0001},
		{"google/gemini-1.5-pro", 1000, 0.001, 0.002},
		{"meta-llama/llama-3-70b", 1000, 0.0001, 0.0003},
		{"mistralai/mistral-7b", 1000, 0.0001, 0.0003},
		{"unknown/model", 1000, 0.0005, 0.002},
	}

	for _, tc := range tests {
		p := NewOpenRouterProvider("test-key", tc.model, "")
		cost := p.EstimateCost(tc.tokens)
		if cost < tc.minCost || cost > tc.maxCost {
			t.Errorf("model %s: unexpected cost %f, expected between %f and %f", tc.model, cost, tc.minCost, tc.maxCost)
		}
	}
}

func TestOpenRouterProvider_ZeroCost(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	if p.EstimateCost(0) != 0 {
		t.Error("expected 0 cost for 0 tokens")
	}
	if p.EstimateCost(-100) != 0 {
		t.Error("expected 0 cost for negative tokens")
	}
}

func TestOpenRouterProvider_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("HTTP-Referer") != "https://subflow.app" {
			t.Error("missing or wrong HTTP-Referer header")
		}
		if r.Header.Get("X-Title") != "SubFlow" {
			t.Error("missing or wrong X-Title header")
		}
		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": "[1] Test"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOpenRouterProvider_Rewrite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"choices": []map[string]any{
				{"message": map[string]string{"content": "[1] Halo\n[2] Dunia"}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "openai/gpt-4o-mini", server.URL)
	result, err := p.Rewrite(context.Background(), []Input{
		{Source: "Hello", Translated: "Halo"},
		{Source: "World", Translated: "Dunia"},
	}, Opts{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %d", len(result))
	}
	if result[0] != "Halo" || result[1] != "Dunia" {
		t.Errorf("unexpected results: %v", result)
	}
}

func TestOpenRouterProvider_EmptyAPIKey(t *testing.T) {
	p := NewOpenRouterProvider("", "", "")
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestOpenRouterProvider_EmptyBatch(t *testing.T) {
	p := NewOpenRouterProvider("test-key", "", "")
	result, err := p.Rewrite(context.Background(), []Input{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d items", len(result))
	}
}

func TestOpenRouterProvider_PaymentRequired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusPaymentRequired)
	}))
	defer server.Close()

	p := NewOpenRouterProvider("test-key", "", server.URL)
	_, err := p.Rewrite(context.Background(), []Input{{Source: "Hi", Translated: "Hai"}}, Opts{})
	if err == nil {
		t.Error("expected payment required error")
	}
}
