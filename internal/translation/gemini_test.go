package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeminiProvider_Name(t *testing.T) {
	p := NewGeminiProvider("test-key", "")
	if p.Name() != ProviderGemini {
		t.Errorf("expected %s, got %s", ProviderGemini, p.Name())
	}
}

func TestGeminiProvider_MaxBatchSize(t *testing.T) {
	p := NewGeminiProvider("test-key", "")
	if p.MaxBatchSize() != 20 {
		t.Errorf("expected 20, got %d", p.MaxBatchSize())
	}
}

func TestGeminiProvider_EstimateCost(t *testing.T) {
	flash := NewGeminiProvider("test-key", "gemini-1.5-flash")
	if flash.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}

	flashCost := flash.EstimateCost(1000)
	if flashCost < 0.00009 || flashCost > 0.00011 {
		t.Errorf("expected ~0.0001 for flash, got %f", flashCost)
	}

	pro := NewGeminiProvider("test-key", "gemini-1.5-pro")
	proCost := pro.EstimateCost(1000)
	if proCost < 0.0019 || proCost > 0.0021 {
		t.Errorf("expected ~0.002 for pro, got %f", proCost)
	}
}

func TestGeminiProvider_DefaultModel(t *testing.T) {
	p := NewGeminiProvider("test-key", "")
	if p.model != "gemini-1.5-flash" {
		t.Errorf("expected default model gemini-1.5-flash, got %s", p.model)
	}
}

func TestGeminiProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}
		if !strings.Contains(r.URL.Path, "/models/gemini-1.5-flash:generateContent") {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if !strings.Contains(r.URL.RawQuery, "key=test-key") {
			t.Error("missing API key in query")
		}

		var req geminiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Contents) != 1 {
			t.Errorf("expected 1 content, got %d", len(req.Contents))
		}

		resp := geminiResponse{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			}{
				{Content: struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				}{Parts: []struct {
					Text string `json:"text"`
				}{{Text: "[1] Halo\n[2] Dunia"}}}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewGeminiProvider("test-key", "gemini-1.5-flash")
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

func TestGeminiProvider_EmptyBatch(t *testing.T) {
	p := NewGeminiProvider("test-key", "")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestGeminiProvider_EmptyAPIKey(t *testing.T) {
	p := NewGeminiProvider("", "")
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestGeminiProvider_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewGeminiProvider("test-key", "")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestGeminiProvider_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewGeminiProvider("bad-key", "")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}

func TestGeminiProvider_GlossaryInPrompt(t *testing.T) {
	var capturedReq geminiRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&capturedReq)
		resp := geminiResponse{
			Candidates: []struct {
				Content struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				} `json:"content"`
			}{
				{Content: struct {
					Parts []struct {
						Text string `json:"text"`
					} `json:"parts"`
				}{Parts: []struct {
					Text string `json:"text"`
				}{{Text: "[1] Test"}}}},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewGeminiProvider("test-key", "")
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

	prompt := capturedReq.Contents[0].Parts[0].Text
	if !strings.Contains(prompt, "Hello → Halo") {
		t.Error("glossary not found in prompt")
	}
}
