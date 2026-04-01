package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOllamaProvider_Name(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.Name() != ProviderOllama {
		t.Errorf("expected %s, got %s", ProviderOllama, p.Name())
	}
}

func TestOllamaProvider_MaxBatchSize(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.MaxBatchSize() != 10 {
		t.Errorf("expected 10, got %d", p.MaxBatchSize())
	}
}

func TestOllamaProvider_EstimateCost(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}
	if p.EstimateCost(1000) != 0 {
		t.Error("expected 0 for any char count (local model)")
	}
	if p.EstimateCost(1000000) != 0 {
		t.Error("expected 0 for any char count (local model)")
	}
}

func TestOllamaProvider_DefaultModel(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.model != "llama3.2" {
		t.Errorf("expected default model llama3.2, got %s", p.model)
	}
}

func TestOllamaProvider_DefaultBaseURL(t *testing.T) {
	p := NewOllamaProvider("", "")
	if p.baseURL != ollamaDefaultURL {
		t.Errorf("expected default URL %s, got %s", ollamaDefaultURL, p.baseURL)
	}
}

func TestOllamaProvider_CustomBaseURL(t *testing.T) {
	p := NewOllamaProvider("http://custom:11434/", "")
	if p.baseURL != "http://custom:11434" {
		t.Errorf("expected trimmed URL http://custom:11434, got %s", p.baseURL)
	}
}

func TestOllamaProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}
		if r.URL.Path != "/api/chat" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}

		var req ollamaRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Model != "llama3.2" {
			t.Errorf("expected model llama3.2, got %s", req.Model)
		}
		if req.Stream != false {
			t.Error("expected stream=false")
		}
		if len(req.Messages) != 1 {
			t.Errorf("expected 1 message, got %d", len(req.Messages))
		}

		resp := ollamaResponse{
			Message: struct {
				Content string `json:"content"`
			}{Content: "[1] Halo\n[2] Dunia"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")

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

func TestOllamaProvider_EmptyBatch(t *testing.T) {
	p := NewOllamaProvider("", "")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestOllamaProvider_ModelNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "model not found"}`))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "nonexistent-model")

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected model not found error")
	}
	if !strings.Contains(err.Error(), "model not found") {
		t.Errorf("expected 'model not found' in error, got: %v", err)
	}
}

func TestOllamaProvider_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal error"}`))
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected server error")
	}
}

func TestOllamaProvider_GlossaryInPrompt(t *testing.T) {
	var capturedReq ollamaRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := ollamaResponse{
			Message: struct {
				Content string `json:"content"`
			}{Content: "[1] Test"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")

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

	prompt := capturedReq.Messages[0].Content
	if !strings.Contains(prompt, "Hello → Halo") {
		t.Error("glossary not found in prompt")
	}
}

func TestOllamaProvider_NoAPIKeyRequired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Error("Ollama should not have Authorization header")
		}
		resp := ollamaResponse{
			Message: struct {
				Content string `json:"content"`
			}{Content: "[1] Test"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewOllamaProvider(server.URL, "llama3.2")

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
