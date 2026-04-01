package asr

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

func TestGroqProvider_NameAndCost(t *testing.T) {
	p := NewGroqProvider("k")
	if p.Name() != "groq" {
		t.Fatalf("Name() = %q, want %q", p.Name(), "groq")
	}
	if got := p.EstimateCost(0); got != 0 {
		t.Fatalf("EstimateCost(0) = %v, want 0", got)
	}
	got := p.EstimateCost(60)
	want := 0.02
	if got != want {
		t.Fatalf("EstimateCost(60) = %v, want %v", got, want)
	}
}

func TestGroqProvider_Transcribe_SendsSegments(t *testing.T) {
	audioPath := writeTempFile(t, "audio.wav", "RIFFxxxxWAVE")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/openai/v1/audio/transcriptions" {
			t.Errorf("path = %s, want /openai/v1/audio/transcriptions", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-key")
		}

		ct := r.Header.Get("Content-Type")
		mt, params, err := mime.ParseMediaType(ct)
		if err != nil {
			t.Errorf("parse content-type: %v", err)
		}
		if mt != "multipart/form-data" {
			t.Errorf("media type = %q, want multipart/form-data", mt)
		}
		mr := multipart.NewReader(r.Body, params["boundary"])

		fields := map[string]string{}
		sawFile := false
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Errorf("next part: %v", err)
				break
			}
			name := part.FormName()
			if name == "file" {
				sawFile = true
				if part.FileName() != filepath.Base(audioPath) {
					t.Errorf("filename = %q, want %q", part.FileName(), filepath.Base(audioPath))
				}
				b, _ := io.ReadAll(part)
				if len(b) == 0 {
					t.Errorf("file content empty")
				}
				continue
			}
			b, _ := io.ReadAll(part)
			fields[name] = string(b)
		}
		if !sawFile {
			t.Errorf("missing file part")
		}
		if fields["model"] != "whisper-large-v3" {
			t.Errorf("model = %q, want %q", fields["model"], "whisper-large-v3")
		}
		if fields["response_format"] != "verbose_json" {
			t.Errorf("response_format = %q, want %q", fields["response_format"], "verbose_json")
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"language": "en",
			"segments": []map[string]any{
				{"id": 0, "start": 0.0, "end": 1.23, "text": "hello"},
				{"id": 1, "start": 1.23, "end": 2.0, "text": "world"},
			},
		})
	}))
	defer ts.Close()

	p := NewGroqProvider("test-key")
	p.baseURL = ts.URL
	p.httpClient = ts.Client()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	segCh, errCh := p.Transcribe(ctx, audioPath, Opts{Language: "auto"})
	segs := drainSegments(segCh)
	if err := drainErr(errCh); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(segs) != 2 {
		t.Fatalf("segments = %d", len(segs))
	}
	if segs[0].Text != "hello" || segs[0].StartMS != 0 || segs[0].EndMS != 1230 {
		t.Fatalf("seg0 = %#v", segs[0])
	}
	if segs[1].Text != "world" || segs[1].StartMS != 1230 || segs[1].EndMS != 2000 {
		t.Fatalf("seg1 = %#v", segs[1])
	}
	if segs[0].Language != "en" {
		t.Fatalf("language = %q", segs[0].Language)
	}
}

func TestGroqProvider_Transcribe_RateLimited(t *testing.T) {
	audioPath := writeTempFile(t, "audio.wav", "data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "7")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
	}))
	defer ts.Close()

	p := NewGroqProvider("k")
	p.baseURL = ts.URL
	p.httpClient = ts.Client()

	segCh, errCh := p.Transcribe(context.Background(), audioPath, Opts{})
	_ = drainSegments(segCh)
	err := drainErr(errCh)
	if err == nil {
		t.Fatalf("expected error")
	}
	if pipeline.GetCode(err) != pipeline.ErrTrnRateLimit {
		t.Fatalf("code = %q, want %q, err=%v", pipeline.GetCode(err), pipeline.ErrTrnRateLimit, err)
	}
}

func TestGroqProvider_Transcribe_InvalidAPIKey(t *testing.T) {
	audioPath := writeTempFile(t, "audio.wav", "data")

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":{"message":"bad key"}}`))
	}))
	defer ts.Close()

	p := NewGroqProvider("k")
	p.baseURL = ts.URL
	p.httpClient = ts.Client()

	segCh, errCh := p.Transcribe(context.Background(), audioPath, Opts{})
	_ = drainSegments(segCh)
	err := drainErr(errCh)
	if err == nil {
		t.Fatalf("expected error")
	}
	if pipeline.GetCode(err) != pipeline.ErrTrnAPIKey {
		t.Fatalf("code = %q, want %q, err=%v", pipeline.GetCode(err), pipeline.ErrTrnAPIKey, err)
	}
}

func TestGroqProvider_Transcribe_ContextCanceled(t *testing.T) {
	audioPath := writeTempFile(t, "audio.wav", "data")
	started := make(chan struct{})
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		select {
		case <-started:
		default:
			close(started)
		}
		<-r.Context().Done()
		if r.Body != nil {
			_ = r.Body.Close()
		}
		return nil, r.Context().Err()
	})}

	p := NewGroqProvider("k")
	p.httpClient = client

	ctx, cancel := context.WithCancel(context.Background())
	segCh, errCh := p.Transcribe(ctx, audioPath, Opts{})
	<-started
	cancel()

	_ = drainSegments(segCh)
	err := drainErr(errCh)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func writeTempFile(t *testing.T, name string, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	return p
}

func drainSegments(ch <-chan Segment) []Segment {
	var out []Segment
	for s := range ch {
		out = append(out, s)
	}
	return out
}

func drainErr(ch <-chan error) error {
	var last error
	for err := range ch {
		if err != nil {
			last = err
		}
	}
	return last
}
