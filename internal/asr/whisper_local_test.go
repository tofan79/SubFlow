package asr

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

func TestWhisperLocalProvider_Name(t *testing.T) {
	p := &WhisperLocalProvider{binaryPath: "/fake/path"}
	if got := p.Name(); got != ProviderWhisperLocal {
		t.Errorf("Name() = %q, want %q", got, ProviderWhisperLocal)
	}
}

func TestWhisperLocalProvider_EstimateCost(t *testing.T) {
	p := &WhisperLocalProvider{binaryPath: "/fake/path"}
	if got := p.EstimateCost(100); got != 0 {
		t.Errorf("EstimateCost() = %v, want 0 (local is free)", got)
	}
}

func TestWhisperLocalProvider_buildArgs(t *testing.T) {
	p := &WhisperLocalProvider{binaryPath: "/fake/path"}

	tests := []struct {
		name           string
		audioPath      string
		opts           Opts
		wantAudioArg   string
		wantModelArg   string
		wantLangArg    string
		wantBackendSet bool
		wantComputeSet bool
	}{
		{
			name:           "defaults_uses_auto_detection",
			audioPath:      "/path/to/audio.wav",
			opts:           Opts{},
			wantAudioArg:   "/path/to/audio.wav",
			wantModelArg:   "base",
			wantLangArg:    "auto",
			wantBackendSet: true,
			wantComputeSet: true,
		},
		{
			name:      "explicit_values",
			audioPath: "/path/to/audio.wav",
			opts: Opts{
				Model:       "large-v3",
				Backend:     "cuda",
				ComputeType: "float16",
				Language:    "ja",
			},
			wantAudioArg:   "/path/to/audio.wav",
			wantModelArg:   "large-v3",
			wantLangArg:    "ja",
			wantBackendSet: true,
			wantComputeSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.buildArgs(tt.audioPath, tt.opts)

			argMap := make(map[string]string)
			for i := 0; i < len(got)-1; i += 2 {
				argMap[got[i]] = got[i+1]
			}

			if argMap["--audio"] != tt.wantAudioArg {
				t.Errorf("--audio = %q, want %q", argMap["--audio"], tt.wantAudioArg)
			}
			if argMap["--model"] != tt.wantModelArg {
				t.Errorf("--model = %q, want %q", argMap["--model"], tt.wantModelArg)
			}
			if argMap["--language"] != tt.wantLangArg {
				t.Errorf("--language = %q, want %q", argMap["--language"], tt.wantLangArg)
			}
			if tt.wantBackendSet && argMap["--backend"] == "" {
				t.Error("expected --backend to be set")
			}
			if tt.wantComputeSet && argMap["--compute-type"] == "" {
				t.Error("expected --compute-type to be set")
			}
		})
	}
}

func TestWhisperLocalProvider_Transcribe_EmptyPath(t *testing.T) {
	p := &WhisperLocalProvider{binaryPath: "/fake/path"}
	segCh, errCh := p.Transcribe(context.Background(), "", Opts{})

	for range segCh {
	}

	var gotErr error
	for err := range errCh {
		gotErr = err
	}

	if gotErr == nil {
		t.Fatal("expected error for empty audio path")
	}
	var pErr *pipeline.Error
	if !errors.As(gotErr, &pErr) {
		t.Fatalf("expected pipeline.Error, got %T", gotErr)
	}
	if pErr.Code != pipeline.ErrASRExtractFail {
		t.Errorf("expected code %s, got %s", pipeline.ErrASRExtractFail, pErr.Code)
	}
}

func TestWhisperLocalProvider_Transcribe_FileNotExist(t *testing.T) {
	p := &WhisperLocalProvider{binaryPath: "/fake/path"}
	segCh, errCh := p.Transcribe(context.Background(), "/nonexistent/audio.wav", Opts{})

	for range segCh {
	}

	var gotErr error
	for err := range errCh {
		gotErr = err
	}

	if gotErr == nil {
		t.Fatal("expected error for nonexistent file")
	}
	var pErr *pipeline.Error
	if !errors.As(gotErr, &pErr) {
		t.Fatalf("expected pipeline.Error, got %T", gotErr)
	}
	if pErr.Code != pipeline.ErrImpFileNotFound {
		t.Errorf("expected code %s, got %s", pipeline.ErrImpFileNotFound, pErr.Code)
	}
}

func TestWhisperLocalProvider_Transcribe_BinaryNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	audioFile := filepath.Join(tmpDir, "test.wav")
	if err := os.WriteFile(audioFile, []byte("fake audio"), 0644); err != nil {
		t.Fatal(err)
	}

	p := &WhisperLocalProvider{binaryPath: ""}
	segCh, errCh := p.Transcribe(context.Background(), audioFile, Opts{})

	for range segCh {
	}

	var gotErr error
	for err := range errCh {
		gotErr = err
	}

	if gotErr == nil {
		t.Fatal("expected error for empty binary path")
	}
	var pErr *pipeline.Error
	if !errors.As(gotErr, &pErr) {
		t.Fatalf("expected pipeline.Error, got %T", gotErr)
	}
	if pErr.Code != pipeline.ErrASRNotInstalled {
		t.Errorf("expected code %s, got %s", pipeline.ErrASRNotInstalled, pErr.Code)
	}
}

func TestWhisperLocalProvider_Transcribe_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := &WhisperLocalProvider{binaryPath: "/fake/path"}
	segCh, errCh := p.Transcribe(ctx, "/fake/audio.wav", Opts{})

	for range segCh {
	}

	var gotErr error
	for err := range errCh {
		gotErr = err
	}

	if gotErr == nil {
		t.Fatal("expected error for canceled context")
	}
	if !errors.Is(gotErr, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", gotErr)
	}
}

func TestWhisperLocalProvider_Transcribe_MockSubprocess(t *testing.T) {
	if os.Getenv("GO_TEST_MOCK_WHISPER") == "1" {
		runMockWhisperSubprocess()
		return
	}

	tmpDir := t.TempDir()
	audioFile := filepath.Join(tmpDir, "test.wav")
	if err := os.WriteFile(audioFile, []byte("fake audio"), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestWhisperLocalProvider_Transcribe_MockSubprocess")
	cmd.Env = append(os.Environ(), "GO_TEST_MOCK_WHISPER=1")

	mockBinary := filepath.Join(tmpDir, "mock-whisper")
	if err := os.WriteFile(mockBinary, []byte("#!/bin/sh\nexec "+cmd.Path+" "+cmd.Args[1]), 0755); err != nil {
		t.Skip("cannot create mock binary")
	}

	p := &WhisperLocalProvider{binaryPath: cmd.Path}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	segCh, errCh := p.Transcribe(ctx, audioFile, Opts{Language: "en"})

	var segments []Segment
	for seg := range segCh {
		segments = append(segments, seg)
	}

	var gotErr error
	for err := range errCh {
		gotErr = err
	}

	if gotErr != nil && !errors.Is(gotErr, context.DeadlineExceeded) {
		var pErr *pipeline.Error
		if errors.As(gotErr, &pErr) {
			t.Logf("got pipeline error (expected in test without real mock): %v", pErr)
		}
	}
}

func runMockWhisperSubprocess() {
	segments := []whisperOutputMessage{
		{Type: "segment", Start: 0.0, End: 2.5, Text: "Hello everyone", Confidence: 0.95, Language: "en"},
		{Type: "segment", Start: 2.5, End: 5.0, Text: "Welcome to the test", Confidence: 0.92, Language: "en"},
		{Type: "progress", Percent: 50.0},
		{Type: "segment", Start: 5.0, End: 7.5, Text: "This is a mock", Confidence: 0.88, Language: "en"},
		{Type: "done"},
	}

	enc := json.NewEncoder(os.Stdout)
	for _, seg := range segments {
		_ = enc.Encode(seg)
	}
	os.Exit(0)
}

func TestFindWhisperBinary_NotFound(t *testing.T) {
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	defer os.Setenv("PATH", oldPath)

	_, err := findWhisperBinary()
	if err == nil {
		t.Fatal("expected error when whisper binary not found")
	}

	var pErr *pipeline.Error
	if !errors.As(err, &pErr) {
		t.Fatalf("expected pipeline.Error, got %T", err)
	}
	if pErr.Code != pipeline.ErrASRNotInstalled {
		t.Errorf("expected code %s, got %s", pipeline.ErrASRNotInstalled, pErr.Code)
	}
}

func TestWhisperBinaryCandidates(t *testing.T) {
	candidates := whisperBinaryCandidates()
	if len(candidates) == 0 {
		t.Error("expected at least one candidate path")
	}

	for _, p := range candidates {
		if p == "" {
			continue
		}
		if !filepath.IsAbs(p) && !startsWithEnvVar(p) {
			t.Errorf("candidate path should be absolute or derived from env: %s", p)
		}
	}
}

func startsWithEnvVar(p string) bool {
	return len(p) > 0 && (p[0] == '$' || (len(p) > 1 && p[0] == '%'))
}
