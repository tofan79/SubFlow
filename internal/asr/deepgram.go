package asr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

// DeepgramProvider implements Provider using Deepgram's /v1/listen endpoint.
type DeepgramProvider struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string // allow override for testing
}

// NewDeepgramProvider creates a new Deepgram provider.
func NewDeepgramProvider(apiKey string) *DeepgramProvider {
	return &DeepgramProvider{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Minute,
		},
		baseURL: "https://api.deepgram.com",
	}
}

func (p *DeepgramProvider) Name() string {
	return ProviderDeepgram
}

// EstimateCost estimates cost in USD using Deepgram Nova-2 rate: $0.0043 / minute.
func (p *DeepgramProvider) EstimateCost(durationSeconds float64) float64 {
	if durationSeconds <= 0 {
		return 0
	}
	return (durationSeconds / 60.0) * 0.0043
}

func (p *DeepgramProvider) Transcribe(ctx context.Context, audioPath string, opts Opts) (<-chan Segment, <-chan error) {
	segCh := make(chan Segment, 16)
	errCh := make(chan error, 1)

	go func() {
		defer close(segCh)
		defer close(errCh)

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

		st, err := f.Stat()
		if err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to stat audio file", err)
			return
		}

		endpoint, err := p.listenURL(opts)
		if err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to build deepgram URL", err)
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, f)
		if err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to create deepgram request", err)
			return
		}
		req.Header.Set("Authorization", "Token "+p.apiKey)
		req.Header.Set("Content-Type", guessAudioContentType(audioPath))
		if st.Size() >= 0 {
			req.ContentLength = st.Size()
		}

		resp, err := p.httpClient.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				errCh <- ctx.Err()
				return
			}
			if isNetTimeout(err) {
				errCh <- pipeline.NewError(pipeline.ErrASRTimeout, "deepgram request timeout", err)
				return
			}
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "deepgram request failed", err)
			return
		}
		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			errCh <- deepgramHTTPError(resp, p.Name())
			return
		}

		var dg deepgramListenResponse
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&dg); err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to decode deepgram response", err)
			return
		}

		words := dg.words()
		if len(words) == 0 {
			return
		}

		lang := strings.TrimSpace(opts.Language)
		if lang == "auto" {
			lang = ""
		}

		segments := groupDeepgramWords(words, lang)
		for i := range segments {
			if ctx.Err() != nil {
				errCh <- ctx.Err()
				return
			}
			segments[i].Index = i
			segCh <- segments[i]
		}
	}()

	return segCh, errCh
}

func (p *DeepgramProvider) listenURL(opts Opts) (string, error) {
	base := strings.TrimRight(strings.TrimSpace(p.baseURL), "/")
	if base == "" {
		base = "https://api.deepgram.com"
	}
	u, err := url.Parse(base + "/v1/listen")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("model", "nova-2")
	q.Set("punctuate", "true")
	q.Set("utterances", "true")
	lang := strings.TrimSpace(opts.Language)
	if lang != "" && lang != "auto" {
		q.Set("language", lang)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func guessAudioContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return "application/octet-stream"
	}
	switch ext {
	case ".wav":
		return "audio/wav"
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "audio/webm"
	default:
		if ct := mime.TypeByExtension(ext); ct != "" {
			// mime.TypeByExtension can include charset for some types; Deepgram accepts that.
			return ct
		}
		return "application/octet-stream"
	}
}

func deepgramHTTPError(resp *http.Response, provider string) error {
	retryAfter := 0
	if ra := strings.TrimSpace(resp.Header.Get("Retry-After")); ra != "" {
		if n, err := strconv.Atoi(ra); err == nil {
			retryAfter = n
		}
	}

	// Read a small body for debug context.
	body, _ := readLimited(resp.Body, 16*1024)
	msg := strings.TrimSpace(string(body))
	if msg == "" {
		msg = resp.Status
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return pipeline.ErrTrnAPIKeyErr(provider)
	case http.StatusTooManyRequests:
		return pipeline.ErrTrnRateLimitErr(provider, retryAfter)
	default:
		if resp.StatusCode >= 500 {
			return pipeline.NewError(pipeline.ErrASRTimeout, fmt.Sprintf("%s server error: %s", provider, msg), nil)
		}
		return pipeline.NewError(pipeline.ErrASRExtractFail, fmt.Sprintf("%s API error: %s", provider, msg), nil)
	}
}

func readLimited(r io.Reader, max int64) ([]byte, error) {
	if max <= 0 {
		return nil, nil
	}
	return io.ReadAll(io.LimitReader(r, max))
}

// Deepgram response models (subset).
type deepgramListenResponse struct {
	Results deepgramResults `json:"results"`
}

type deepgramResults struct {
	Channels   []deepgramChannel   `json:"channels"`
	Utterances []deepgramUtterance `json:"utterances"`
}

type deepgramChannel struct {
	Alternatives []deepgramAlternative `json:"alternatives"`
}

type deepgramAlternative struct {
	Words []deepgramWord `json:"words"`
}

type deepgramUtterance struct {
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Transcript string  `json:"transcript"`
	Confidence float64 `json:"confidence"`
}

type deepgramWord struct {
	Word           string  `json:"word"`
	PunctuatedWord string  `json:"punctuated_word"`
	Start          float64 `json:"start"`
	End            float64 `json:"end"`
	Confidence     float64 `json:"confidence"`
}

func (r deepgramListenResponse) words() []deepgramWord {
	if len(r.Results.Channels) == 0 {
		return nil
	}
	ch := r.Results.Channels[0]
	if len(ch.Alternatives) == 0 {
		return nil
	}
	return ch.Alternatives[0].Words
}

func groupDeepgramWords(words []deepgramWord, language string) []Segment {
	const (
		gapSplitSeconds     = 0.80
		maxSegmentSeconds   = 8.0
		maxSegmentWordCount = 14
	)

	var out []Segment

	var (
		curText      strings.Builder
		curStart     float64
		curEnd       float64
		curConfSum   float64
		curConfCount int
		curWordCount int
		prevEnd      float64
		haveCur      bool
	)

	flush := func() {
		if !haveCur {
			return
		}
		text := strings.TrimSpace(curText.String())
		if text != "" {
			conf := 0.0
			if curConfCount > 0 {
				conf = curConfSum / float64(curConfCount)
			}
			out = append(out, Segment{
				StartMS:    int64(math.Round(curStart * 1000)),
				EndMS:      int64(math.Round(curEnd * 1000)),
				Text:       text,
				Confidence: conf,
				Language:   language,
			})
		}
		curText.Reset()
		curStart, curEnd, curConfSum = 0, 0, 0
		curConfCount, curWordCount = 0, 0
		haveCur = false
	}

	for i := range words {
		w := words[i]
		word := strings.TrimSpace(w.PunctuatedWord)
		if word == "" {
			word = strings.TrimSpace(w.Word)
		}
		if word == "" {
			continue
		}

		// Start a new segment.
		if !haveCur {
			haveCur = true
			curStart = w.Start
			curEnd = w.End
			prevEnd = w.End
			curText.WriteString(word)
			curConfSum += w.Confidence
			curConfCount++
			curWordCount++
			if endsSentence(word) {
				flush()
			}
			continue
		}

		gap := w.Start - prevEnd
		wouldBeDuration := w.End - curStart
		split := false
		if gap > gapSplitSeconds {
			split = true
		}
		if curWordCount >= maxSegmentWordCount {
			split = true
		}
		if wouldBeDuration > maxSegmentSeconds {
			split = true
		}

		if split {
			flush()
			// start new segment with this word
			haveCur = true
			curStart = w.Start
			curEnd = w.End
			prevEnd = w.End
			curText.WriteString(word)
			curConfSum += w.Confidence
			curConfCount++
			curWordCount = 1
			if endsSentence(word) {
				flush()
			}
			continue
		}

		curText.WriteByte(' ')
		curText.WriteString(word)
		curEnd = w.End
		prevEnd = w.End
		curConfSum += w.Confidence
		curConfCount++
		curWordCount++
		if endsSentence(word) {
			flush()
		}
	}

	flush()
	return out
}

func endsSentence(word string) bool {
	word = strings.TrimSpace(word)
	if word == "" {
		return false
	}
	last := word[len(word)-1]
	switch last {
	case '.', '?', '!':
		return true
	default:
		return false
	}
}
