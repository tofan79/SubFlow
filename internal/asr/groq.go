package asr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

const (
	groqDefaultBaseURL     = "https://api.groq.com"
	groqTranscriptionsPath = "/openai/v1/audio/transcriptions"
	groqWhisperModel       = "whisper-large-v3"
)

type GroqProvider struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func NewGroqProvider(apiKey string) *GroqProvider {
	return &GroqProvider{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Minute},
		baseURL:    groqDefaultBaseURL,
	}
}

func (p *GroqProvider) Name() string { return ProviderGroq }

func (p *GroqProvider) EstimateCost(durationSeconds float64) float64 {
	if durationSeconds <= 0 {
		return 0
	}
	return (durationSeconds / 60.0) * 0.02
}

func (p *GroqProvider) Transcribe(ctx context.Context, audioPath string, opts Opts) (<-chan Segment, <-chan error) {
	segCh := make(chan Segment, 16)
	errCh := make(chan error, 1)

	go func() {
		defer close(segCh)
		defer close(errCh)

		if ctx == nil {
			ctx = context.Background()
		}
		if ctx.Err() != nil {
			errCh <- ctx.Err()
			return
		}

		if strings.TrimSpace(p.apiKey) == "" {
			errCh <- pipeline.ErrTrnAPIKeyErr(p.Name())
			return
		}
		if strings.TrimSpace(audioPath) == "" {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "audio path is empty", nil)
			return
		}

		f, err := os.Open(audioPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				errCh <- pipeline.ErrFileNotFound(audioPath, err)
				return
			}
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, fmt.Sprintf("failed to open audio file: %s", audioPath), err)
			return
		}
		defer func() { _ = f.Close() }()

		body, contentType := p.buildMultipartBody(ctx, audioPath, f, opts)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.endpointURL(), body)
		if err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to create groq request", err)
			return
		}
		req.Close = true
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
		req.Header.Set("Content-Type", contentType)

		client := p.httpClient
		if client == nil {
			client = http.DefaultClient
		}
		rt := client.Transport
		if rt == nil {
			rt = http.DefaultTransport
		}
		if tr, ok := rt.(*http.Transport); ok {
			cloned := tr.Clone()
			cloned.DisableKeepAlives = true
			rt = cloned
		}
		clientCopy := *client
		clientCopy.Transport = rt

		resp, err := clientCopy.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				errCh <- ctx.Err()
				return
			}
			if isNetTimeout(err) {
				errCh <- pipeline.NewError(pipeline.ErrASRTimeout, "groq request timeout", err)
				return
			}
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "groq request failed", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- p.httpError(resp)
			return
		}

		var vr groqVerboseResponse
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&vr); err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to decode groq response", err)
			return
		}

		lang := vr.Language
		if lang == "" {
			lang = opts.Language
		}

		if len(vr.Segments) == 0 {
			if strings.TrimSpace(vr.Text) == "" {
				return
			}
			select {
			case segCh <- Segment{Index: 0, StartMS: 0, EndMS: 0, Text: vr.Text, Confidence: 0, Language: lang}:
			case <-ctx.Done():
				errCh <- ctx.Err()
			}
			return
		}

		for _, s := range vr.Segments {
			seg := Segment{
				Index:      s.ID,
				StartMS:    secondsToMS(s.Start),
				EndMS:      secondsToMS(s.End),
				Text:       s.Text,
				Confidence: 0,
				Language:   lang,
			}
			select {
			case segCh <- seg:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}
	}()

	return segCh, errCh
}

func (p *GroqProvider) endpointURL() string {
	base := strings.TrimRight(strings.TrimSpace(p.baseURL), "/")
	if base == "" {
		base = groqDefaultBaseURL
	}
	return base + groqTranscriptionsPath
}

func (p *GroqProvider) buildMultipartBody(ctx context.Context, audioPath string, f *os.File, opts Opts) (io.Reader, string) {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	contentType := w.FormDataContentType()
	pc := &pipeClose{pw: pw}
	done := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			pc.Close(ctx.Err())
		case <-done:
		}
	}()

	go func() {
		defer close(done)
		defer func() {
			if err := w.Close(); err != nil {
				pc.Close(err)
				return
			}
			pc.Close(nil)
		}()

		if ctx.Err() != nil {
			pc.Close(ctx.Err())
			return
		}

		fw, err := w.CreateFormFile("file", filepath.Base(audioPath))
		if err != nil {
			pc.Close(err)
			return
		}
		if _, err := io.Copy(fw, &ctxReader{ctx: ctx, r: f}); err != nil {
			pc.Close(err)
			return
		}

		if err := w.WriteField("model", groqWhisperModel); err != nil {
			pc.Close(err)
			return
		}
		if err := w.WriteField("response_format", "verbose_json"); err != nil {
			pc.Close(err)
			return
		}
		if lang := strings.TrimSpace(opts.Language); lang != "" && lang != "auto" {
			if err := w.WriteField("language", lang); err != nil {
				pc.Close(err)
				return
			}
		}
	}()

	return pr, contentType
}

type pipeClose struct {
	once sync.Once
	pw   *io.PipeWriter
}

func (p *pipeClose) Close(err error) {
	if p == nil || p.pw == nil {
		return
	}
	p.once.Do(func() {
		if err != nil {
			_ = p.pw.CloseWithError(err)
			return
		}
		_ = p.pw.Close()
	})
}

func (p *GroqProvider) httpError(resp *http.Response) error {
	retryAfter := 0
	if ra := strings.TrimSpace(resp.Header.Get("Retry-After")); ra != "" {
		if v, err := strconv.Atoi(ra); err == nil {
			retryAfter = v
		}
	}

	var apiErr groqErrorResponse
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	_ = json.Unmarshal(b, &apiErr)

	msg := strings.TrimSpace(apiErr.Error.Message)
	if msg == "" {
		msg = strings.TrimSpace(string(b))
	}
	if msg == "" {
		msg = resp.Status
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return pipeline.ErrTrnRateLimitErr(p.Name(), retryAfter)
	case http.StatusUnauthorized, http.StatusForbidden:
		return pipeline.ErrTrnAPIKeyErr(p.Name())
	default:
		if resp.StatusCode >= 500 {
			return pipeline.NewError(pipeline.ErrASRTimeout, fmt.Sprintf("%s server error: %s", p.Name(), msg), nil)
		}
		return pipeline.NewError(pipeline.ErrASRExtractFail, fmt.Sprintf("%s API error (%d): %s", p.Name(), resp.StatusCode, msg), nil)
	}
}

type groqVerboseResponse struct {
	Text     string        `json:"text"`
	Language string        `json:"language"`
	Segments []groqSegment `json:"segments"`
}

type groqSegment struct {
	ID    int     `json:"id"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
	Text  string  `json:"text"`
}

type groqErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error"`
}

func secondsToMS(v float64) int64 {
	if v <= 0 {
		return 0
	}
	ms := int64(v*1000.0 + 0.5)
	if ms < 0 {
		return 0
	}
	return ms
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *ctxReader) Read(p []byte) (int, error) {
	if r.ctx.Err() != nil {
		return 0, r.ctx.Err()
	}
	return r.r.Read(p)
}
