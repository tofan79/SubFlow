package rewrite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewEngine(t *testing.T) {
	cfg := Config{
		Provider: ProviderOllama,
		Model:    "llama3.2",
	}
	engine, err := NewEngine(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if engine.ProviderName() != ProviderOllama {
		t.Errorf("expected %s, got %s", ProviderOllama, engine.ProviderName())
	}
}

func TestNewEngine_UnknownProvider(t *testing.T) {
	cfg := Config{Provider: "unknown"}
	_, err := NewEngine(cfg)
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}

func TestBuildSystemPrompt_Natural(t *testing.T) {
	opts := Opts{
		TonePreset:      string(ToneNatural),
		MaxCharsPerLine: 42,
		MaxLines:        2,
		MaxCPS:          17.0,
	}
	prompt := BuildSystemPrompt(opts)

	if !strings.Contains(prompt, "Natural") {
		t.Error("expected Natural tone in prompt")
	}
	if !strings.Contains(prompt, "42") {
		t.Error("expected max chars in prompt")
	}
}

func TestBuildSystemPrompt_Formal(t *testing.T) {
	opts := Opts{TonePreset: string(ToneFormal)}
	prompt := BuildSystemPrompt(opts)
	if !strings.Contains(prompt, "Formal") {
		t.Error("expected Formal tone in prompt")
	}
}

func TestBuildSystemPrompt_Casual(t *testing.T) {
	opts := Opts{TonePreset: string(ToneCasual)}
	prompt := BuildSystemPrompt(opts)
	if !strings.Contains(prompt, "Casual") {
		t.Error("expected Casual tone in prompt")
	}
}

func TestBuildSystemPrompt_Cinematic(t *testing.T) {
	opts := Opts{TonePreset: string(ToneCinematic)}
	prompt := BuildSystemPrompt(opts)
	if !strings.Contains(prompt, "Cinematic") {
		t.Error("expected Cinematic tone in prompt")
	}
}

func TestBuildSystemPrompt_WithGlossary(t *testing.T) {
	opts := Opts{
		Glossary: []GlossaryTerm{
			{SourceTerm: "AI", TargetTerm: "Kecerdasan Buatan"},
		},
	}
	prompt := BuildSystemPrompt(opts)
	if !strings.Contains(prompt, "AI → Kecerdasan Buatan") {
		t.Error("expected glossary in prompt")
	}
}

func TestBuildUserPrompt(t *testing.T) {
	batch := []Input{
		{Source: "Hello", Translated: "Halo", Speaker: "John", Emotion: "happy", Context: "greeting"},
		{Source: "World", Translated: "Dunia"},
	}
	prompt := BuildUserPrompt(batch)

	if !strings.Contains(prompt, "[1]") {
		t.Error("expected line number [1]")
	}
	if !strings.Contains(prompt, "[2]") {
		t.Error("expected line number [2]")
	}
	if !strings.Contains(prompt, "Source: Hello") {
		t.Error("expected source text")
	}
	if !strings.Contains(prompt, "Speaker: John") {
		t.Error("expected speaker")
	}
	if !strings.Contains(prompt, "Emotion: happy") {
		t.Error("expected emotion")
	}
	if !strings.Contains(prompt, "Context: greeting") {
		t.Error("expected context")
	}
}

func TestEstimateTokens(t *testing.T) {
	batch := []Input{
		{Source: "12345678", Translated: "12345678"},
	}
	tokens := EstimateTokens(batch)
	if tokens != 4 {
		t.Errorf("expected 4 tokens (16 chars / 4), got %d", tokens)
	}
}

func TestEstimateTokens_Empty(t *testing.T) {
	batch := []Input{}
	tokens := EstimateTokens(batch)
	if tokens != 0 {
		t.Errorf("expected 0 tokens, got %d", tokens)
	}
}

func TestEngine_RewriteDefaults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"message": map[string]string{"content": "[1] Test"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := Config{
		Provider: ProviderOllama,
		BaseURL:  server.URL,
		Model:    "test",
	}
	engine, _ := NewEngine(cfg)

	result, err := engine.Rewrite(context.Background(), []Input{
		{Source: "Hi", Translated: "Hai"},
	}, Opts{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 result, got %d", len(result))
	}
}

func TestEngine_RewriteEmpty(t *testing.T) {
	cfg := Config{Provider: ProviderOllama}
	engine, _ := NewEngine(cfg)

	result, err := engine.Rewrite(context.Background(), []Input{}, Opts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestEngine_EstimateCost(t *testing.T) {
	cfg := Config{Provider: ProviderOllama}
	engine, _ := NewEngine(cfg)

	cost := engine.EstimateCost(1000)
	if cost != 0 {
		t.Error("expected 0 cost for Ollama")
	}
}
