package translation

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeepLProvider_Name(t *testing.T) {
	p := NewDeepLProvider("test-key")
	if p.Name() != ProviderDeepL {
		t.Errorf("expected %s, got %s", ProviderDeepL, p.Name())
	}
}

func TestDeepLProvider_MaxBatchSize(t *testing.T) {
	p := NewDeepLProvider("test-key")
	if p.MaxBatchSize() != 50 {
		t.Errorf("expected 50, got %d", p.MaxBatchSize())
	}
}

func TestDeepLProvider_EstimateCost(t *testing.T) {
	p := NewDeepLProvider("test-key")

	if p.EstimateCost(0) != 0 {
		t.Error("expected 0 for 0 chars")
	}

	cost := p.EstimateCost(1000)
	expected := 1000 * 0.00002
	if cost != expected {
		t.Errorf("expected %f, got %f", expected, cost)
	}
}

func TestDeepLProvider_FreeVsPro(t *testing.T) {
	free := NewDeepLProvider("test-key:fx")
	if free.baseURL != deeplFreeURL {
		t.Errorf("expected free URL for :fx key")
	}

	pro := NewDeepLProvider("test-key-pro")
	if pro.baseURL != deeplProURL {
		t.Errorf("expected pro URL for non-:fx key")
	}
}

func TestDeepLProvider_Translate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "DeepL-Auth-Key test-key" {
			t.Error("missing or wrong Authorization header")
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("missing Content-Type header")
		}

		var req deeplRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if len(req.Text) != 2 {
			t.Errorf("expected 2 texts, got %d", len(req.Text))
		}

		resp := deeplResponse{
			Translations: []struct {
				Text string `json:"text"`
			}{
				{Text: "Halo"},
				{Text: "Dunia"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	p := NewDeepLProvider("test-key")
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

func TestDeepLProvider_EmptyBatch(t *testing.T) {
	p := NewDeepLProvider("test-key")
	result, err := p.Translate(context.Background(), []string{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestDeepLProvider_EmptyAPIKey(t *testing.T) {
	p := NewDeepLProvider("")
	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected error for empty API key")
	}
}

func TestDeepLProvider_RateLimitError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	p := NewDeepLProvider("test-key")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected rate limit error")
	}
}

func TestDeepLProvider_AuthError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	p := NewDeepLProvider("bad-key")
	p.baseURL = server.URL

	_, err := p.Translate(context.Background(), []string{"Hello"}, Opts{})
	if err == nil {
		t.Error("expected auth error")
	}
}
