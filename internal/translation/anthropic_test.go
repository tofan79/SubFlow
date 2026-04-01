package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAnthropicProvider_Name(t *testing.T) {
	p := NewAnthropicProvider("test-key", "")
	if p.Name() != ProviderAnthropic {
		t.Errorf("expected %s, got %s", ProviderAnthropic, p.Name())
	}
}

func TestAnthropicProvider_MaxBatchSize(t *testing.T) {
	p := NewAnthropicProvider("test-key", "")
	if p.MaxBatchSize() != 20 {
		t.Errorf("expected 20, got %d", p.MaxBatchSize())
	}
}

func TestAnthropicProvider_EstimateCost(t *testing.T) {
	haiku := NewAnthropicProvider("test-key", "claude-3-haiku-20240307")
	if haiku.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}

	haikuCost := haiku.EstimateCost(1000)
	expectedHaiku := 1000 * 0.0000003
	if haikuCost != expectedHaiku {
		t.Errorf("expected %f for haiku, got %f", expectedHaiku, haikuCost)
	}

	sonnet := NewAnthropicProvider("test-key", "claude-3-sonnet")
	sonnetCost := sonnet.EstimateCost(1000)
	expectedSonnet := 1000 * 0.000004
	if sonnetCost != expectedSonnet {
		t.Errorf("expected %f for sonnet, got %f", expectedSonnet, sonnetCost)
	}

	opus := NewAnthropicProvider("test-key", "claude-3-opus")
	opusCost := opus.EstimateCost(1000)
	expectedOpus := 1000 * 0.00002
	if opusCost != expectedOpus {
		t.Errorf("expected %f for opus, got %f", expectedOpus, opusCost)
	}
}

func TestAnthropicProvider_DefaultModel(t *testing.T) {
	p := NewAnthropicProvider("test-key", "")
	if p.model != "claude-3-haiku-20240307" {
		t.Errorf("expected default model claude-3-haiku-20240307, got %s", p.model)
	}
}

func TestAnthropicProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test-key" {
			t.Error("missing or wrong x-api-key header")
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Error("missing or wrong anthropic-version header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}
		if r.URL.Path != "/v1/messages" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}

		var req anthropicRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.System == "" {
			t.Error("expected system prompt")
		}
		if len(req.Messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(req.Messages))
		}

		resp := anthropicResponse{
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{Type: "text", Text: "[1] Halo\n[2] Dunia"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewAnthropicProvider("test-key", "claude-3-haiku-20240307")
	p.baseURL = server.URL

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

func TestAnthropicProvider_EmptyBatch(t *testing.T) {
	p := NewAnthropicProvider("test-key", "")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestAnthropicProvider_EmptyAPIKey(t *testing.T) {
	p := NewAnthropicProvider("", "")
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestAnthropicProvider_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewAnthropicProvider("test-key", "")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestAnthropicProvider_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewAnthropicProvider("bad-key", "")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}

func TestAnthropicProvider_GlossaryInPrompt(t *testing.T) {
	var capturedReq anthropicRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := anthropicResponse{
			Content: []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			}{
				{Type: "text", Text: "[1] Test"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewAnthropicProvider("test-key", "")
	p.baseURL = server.URL

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

	if !strings.Contains(capturedReq.System, "Hello → Halo") {
		t.Error("glossary not found in system prompt")
	}
}
