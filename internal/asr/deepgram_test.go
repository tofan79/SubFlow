package asr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

func TestDeepgramProviderTranscribe_Success(t *testing.T) {
	audioBytes := []byte("RIFF....WAVEfmt ")
	tmpDir := t.TempDir()
	audioPath := filepath.Join(tmpDir, "audio.wav")
	if err := os.WriteFile(audioPath, audioBytes, 0o600); err != nil {
		t.Fatalf("write temp audio: %v", err)
	}

	var gotAuth string
	var gotQuery map[string]string
	var gotContentType string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/listen" {
			t.Errorf("path = %s, want /v1/listen", r.URL.Path)
		}
		gotAuth = r.Header.Get("Authorization")
		gotContentType = r.Header.Get("Content-Type")
		gotQuery = map[string]string{}
		for k, v := range r.URL.Query() {
			if len(v) > 0 {
				gotQuery[k] = v[0]
			}
		}
		b, _ := io.ReadAll(r.Body)
		gotBody = b

		resp := map[string]any{
			"results": map[string]any{
				"channels": []any{
					map[string]any{
						"alternatives": []any{
							map[string]any{
								"words": []any{
									map[string]any{"word": "hello", "punctuated_word": "Hello", "start": 0.0, "end": 0.5, "confidence": 0.9},
									map[string]any{"word": "world", "punctuated_word": "world.", "start": 0.5, "end": 1.0, "confidence": 0.8},
									map[string]any{"word": "how", "punctuated_word": "How", "start": 1.8, "end": 2.0, "confidence": 0.95},
									map[string]any{"word": "are", "punctuated_word": "are", "start": 2.0, "end": 2.1, "confidence": 0.96},
									map[string]any{"word": "you", "punctuated_word": "you?", "start": 2.1, "end": 2.3, "confidence": 0.97},
								},
							},
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	p := NewDeepgramProvider("testkey")
	p.baseURL = srv.URL

	segCh, errCh := p.Transcribe(context.Background(), audioPath, Opts{Language: "en"})

	var segs []Segment
	for s := range segCh {
		segs = append(segs, s)
	}
	if err := <-errCh; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Token testkey" {
		t.Fatalf("Authorization = %q, want %q", gotAuth, "Token testkey")
	}
	if gotQuery["model"] != "nova-2" || gotQuery["punctuate"] != "true" || gotQuery["utterances"] != "true" || gotQuery["language"] != "en" {
		t.Fatalf("query = %#v, want model=nova-2 punctuate=true utterances=true language=en", gotQuery)
	}
	if gotContentType != "audio/wav" {
		t.Fatalf("Content-Type = %q, want %q", gotContentType, "audio/wav")
	}
	if string(gotBody) != string(audioBytes) {
		t.Fatalf("body mismatch")
	}

	if len(segs) != 2 {
		t.Fatalf("segments = %d, want 2: %#v", len(segs), segs)
	}
	if segs[0].Index != 0 || segs[0].StartMS != 0 || segs[0].EndMS != 1000 || segs[0].Text != "Hello world." || segs[0].Language != "en" {
		t.Fatalf("segment[0] = %#v", segs[0])
	}
	if segs[1].Index != 1 || segs[1].StartMS != 1800 || segs[1].EndMS != 2300 || segs[1].Text != "How are you?" || segs[1].Language != "en" {
		t.Fatalf("segment[1] = %#v", segs[1])
	}
}

func TestDeepgramProviderTranscribe_RateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("rate limited"))
	}))
	defer srv.Close()

	audioPath := filepath.Join(t.TempDir(), "audio.wav")
	_ = os.WriteFile(audioPath, []byte("x"), 0o600)

	p := NewDeepgramProvider("testkey")
	p.baseURL = srv.URL
	_, errCh := p.Transcribe(context.Background(), audioPath, Opts{Language: "en"})

	err := <-errCh
	if err == nil {
		t.Fatalf("expected error")
	}
	if pipeline.GetCode(err) != pipeline.ErrTrnRateLimit {
		t.Fatalf("code = %q, want %q (err=%v)", pipeline.GetCode(err), pipeline.ErrTrnRateLimit, err)
	}
}

func TestDeepgramProviderTranscribe_InvalidAPIKey(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	audioPath := filepath.Join(t.TempDir(), "audio.wav")
	_ = os.WriteFile(audioPath, []byte("x"), 0o600)

	p := NewDeepgramProvider("testkey")
	p.baseURL = srv.URL
	_, errCh := p.Transcribe(context.Background(), audioPath, Opts{Language: "en"})

	err := <-errCh
	if err == nil {
		t.Fatalf("expected error")
	}
	if pipeline.GetCode(err) != pipeline.ErrTrnAPIKey {
		t.Fatalf("code = %q, want %q (err=%v)", pipeline.GetCode(err), pipeline.ErrTrnAPIKey, err)
	}
}

func TestDeepgramProviderTranscribe_ContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(2 * time.Second):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"results":{"channels":[{"alternatives":[{"words":[]}]}]}}`))
			return
		}
	}))
	defer srv.Close()

	audioPath := filepath.Join(t.TempDir(), "audio.wav")
	_ = os.WriteFile(audioPath, []byte("x"), 0o600)

	p := NewDeepgramProvider("testkey")
	p.baseURL = srv.URL

	ctx, cancel := context.WithCancel(context.Background())
	segCh, errCh := p.Transcribe(ctx, audioPath, Opts{Language: "en"})
	cancel()

	// Drain segments (should be none) and read error.
	for range segCh {
	}
	err := <-errCh
	if err == nil {
		t.Fatalf("expected cancellation error")
	}
	if err != context.Canceled {
		t.Fatalf("err = %v, want %v", err, context.Canceled)
	}
}
