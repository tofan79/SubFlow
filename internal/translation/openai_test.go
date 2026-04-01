package translation

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

func TestOpenAIProvider_MaxBatchSize(t *testing.T) {
	p := NewOpenAIProvider("test-key", "", "")
	if p.MaxBatchSize() != 20 {
		t.Errorf("expected 20, got %d", p.MaxBatchSize())
	}
}

func TestOpenAIProvider_EstimateCost(t *testing.T) {
	mini := NewOpenAIProvider("test-key", "gpt-4o-mini", "")
	if mini.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}

	miniCost := mini.EstimateCost(1000)
	if miniCost < 0.00039 || miniCost > 0.00041 {
		t.Errorf("expected ~0.0004 for mini, got %f", miniCost)
	}

	gpt4 := NewOpenAIProvider("test-key", "gpt-4", "")
	gpt4Cost := gpt4.EstimateCost(1000)
	if gpt4Cost < 0.0059 || gpt4Cost > 0.0061 {
		t.Errorf("expected ~0.006 for gpt-4, got %f", gpt4Cost)
	}
}

func TestOpenAIProvider_DefaultModel(t *testing.T) {
	p := NewOpenAIProvider("test-key", "", "")
	if p.model != "gpt-4o-mini" {
		t.Errorf("expected default model gpt-4o-mini, got %s", p.model)
	}
}

func TestOpenAIProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}

		var req openaiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(req.Messages))
		}

		resp := openaiResponse{
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

	p := NewOpenAIProvider("test-key", "gpt-4o-mini", server.URL)

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

func TestOpenAIProvider_EmptyBatch(t *testing.T) {
	p := NewOpenAIProvider("test-key", "", "")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestOpenAIProvider_EmptyAPIKey(t *testing.T) {
	p := NewOpenAIProvider("", "", "")
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestOpenAIProvider_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewOpenAIProvider("test-key", "", server.URL)

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestOpenAIProvider_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewOpenAIProvider("bad-key", "", server.URL)

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}

func TestOpenAIProvider_GlossaryInPrompt(t *testing.T) {
	var capturedReq openaiRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := openaiResponse{
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

	p := NewOpenAIProvider("test-key", "", server.URL)
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
	if !contains(systemPrompt, "Hello → Halo") {
		t.Error("glossary not found in system prompt")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
