package asr

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

func TestFindFFmpeg_FindsFromPATH(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mock ffmpeg uses POSIX shell")
	}
	toolsDir := t.TempDir()
	ffmpegPath := writeMockTool(t, toolsDir, "ffmpeg")

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", toolsDir+string(os.PathListSeparator)+oldPath)

	got, err := FindFFmpeg()
	if err != nil {
		t.Fatalf("FindFFmpeg returned error: %v", err)
	}
	if filepath.Clean(got) != filepath.Clean(ffmpegPath) {
		t.Fatalf("FindFFmpeg path mismatch: got %q want %q", got, ffmpegPath)
	}
}

func TestExtractAudio_WritesWAVAndPassesArgs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mock ffmpeg uses POSIX shell")
	}
	tmp := t.TempDir()
	toolsDir := filepath.Join(tmp, "tools")
	if err := os.MkdirAll(toolsDir, 0o755); err != nil {
		t.Fatalf("mkdir tools: %v", err)
	}
	_ = writeMockTool(t, toolsDir, "ffmpeg")

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", toolsDir+string(os.PathListSeparator)+oldPath)

	argsFile := filepath.Join(tmp, "ffmpeg_args.txt")
	t.Setenv("SUBFLOW_MOCK_ARGS_FILE", argsFile)

	inPath := filepath.Join(tmp, "input.mp4")
	if err := os.WriteFile(inPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write dummy input: %v", err)
	}
	outPath := filepath.Join(tmp, "out.wav")

	if err := ExtractAudio(context.Background(), inPath, outPath); err != nil {
		t.Fatalf("ExtractAudio error: %v", err)
	}

	b, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(b) < 12 {
		t.Fatalf("output too small: %d", len(b))
	}
	if !bytes.Equal(b[:4], []byte("RIFF")) || !bytes.Equal(b[8:12], []byte("WAVE")) {
		t.Fatalf("output is not a WAV file")
	}

	argsBytes, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("read args file: %v", err)
	}
	gotArgs := splitLines(argsBytes)
	wantArgs := []string{"-i", inPath, "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", "-y", outPath}
	if joinLines(gotArgs) != strings.Join(wantArgs, "\n") {
		t.Fatalf("args mismatch:\n--- got\n%q\n--- want\n%q", gotArgs, wantArgs)
	}
}

func TestExtractAudio_ContextTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mock ffmpeg uses POSIX shell")
	}
	tmp := t.TempDir()
	toolsDir := filepath.Join(tmp, "tools")
	if err := os.MkdirAll(toolsDir, 0o755); err != nil {
		t.Fatalf("mkdir tools: %v", err)
	}
	_ = writeMockTool(t, toolsDir, "ffmpeg")

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", toolsDir+string(os.PathListSeparator)+oldPath)
	t.Setenv("SUBFLOW_MOCK_SLEEP_MS", "200")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	inPath := filepath.Join(tmp, "input.mp4")
	if err := os.WriteFile(inPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write dummy input: %v", err)
	}
	outPath := filepath.Join(tmp, "out.wav")

	err := ExtractAudio(ctx, inPath, outPath)
	if err == nil {
		t.Fatalf("expected error")
	}
	var perr *pipeline.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected pipeline error, got %T", err)
	}
	if perr.Code != pipeline.ErrASRExtractFail {
		t.Fatalf("unexpected error code: %s", perr.Code)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline exceeded")
	}
}

func TestGetAudioDuration_ParsesJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mock ffprobe uses POSIX shell")
	}
	tmp := t.TempDir()
	toolsDir := filepath.Join(tmp, "tools")
	if err := os.MkdirAll(toolsDir, 0o755); err != nil {
		t.Fatalf("mkdir tools: %v", err)
	}
	_ = writeMockTool(t, toolsDir, "ffprobe")

	oldPath := os.Getenv("PATH")
	t.Setenv("PATH", toolsDir+string(os.PathListSeparator)+oldPath)

	argsFile := filepath.Join(tmp, "ffprobe_args.txt")
	t.Setenv("SUBFLOW_MOCK_ARGS_FILE", argsFile)
	t.Setenv("SUBFLOW_MOCK_DURATION", "1.25")

	audioPath := filepath.Join(tmp, "audio.wav")
	if err := os.WriteFile(audioPath, []byte("x"), 0o644); err != nil {
		t.Fatalf("write dummy audio: %v", err)
	}

	secs, err := GetAudioDuration(audioPath)
	if err != nil {
		t.Fatalf("GetAudioDuration error: %v", err)
	}
	if secs != 1.25 {
		t.Fatalf("duration mismatch: got %v want %v", secs, 1.25)
	}

	argsBytes, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatalf("read args file: %v", err)
	}
	gotArgs := splitLines(argsBytes)
	if len(gotArgs) == 0 {
		t.Fatalf("no args recorded")
	}
	if string(gotArgs[len(gotArgs)-1]) != audioPath {
		t.Fatalf("ffprobe last arg mismatch: got %q want %q", string(gotArgs[len(gotArgs)-1]), audioPath)
	}
}

func TestExtractAudio_EmptyPathsAreWrapped(t *testing.T) {
	err := ExtractAudio(context.Background(), "", "")
	if err == nil {
		t.Fatalf("expected error")
	}
	var perr *pipeline.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected pipeline error, got %T", err)
	}
	if perr.Code != pipeline.ErrASRExtractFail {
		t.Fatalf("unexpected error code: %s", perr.Code)
	}
}

func TestGetAudioDuration_EmptyPathWrapped(t *testing.T) {
	_, err := GetAudioDuration("")
	if err == nil {
		t.Fatalf("expected error")
	}
	var perr *pipeline.Error
	if !errors.As(err, &perr) {
		t.Fatalf("expected pipeline error, got %T", err)
	}
	if perr.Code != pipeline.ErrASRExtractFail {
		t.Fatalf("unexpected error code: %s", perr.Code)
	}
}

func writeMockTool(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	script := "#!/bin/sh\n" +
		"set -eu\n" +
		"tool=\"$(basename \"$0\")\"\n" +
		"tool=\"${tool%.*}\"\n" +
		"if [ \"${SUBFLOW_MOCK_SLEEP_MS:-}\" != \"\" ]; then\n" +
		"  ms=\"$SUBFLOW_MOCK_SLEEP_MS\"\n" +
		"  sec=$((ms/1000))\n" +
		"  rem=$((ms%1000))\n" +
		"  if [ \"$rem\" -eq 0 ]; then\n" +
		"    sleep \"$sec\"\n" +
		"  else\n" +
		"    sleep \"$sec.$(printf '%03d' \"$rem\")\"\n" +
		"  fi\n" +
		"fi\n" +
		"if [ \"${SUBFLOW_MOCK_ARGS_FILE:-}\" != \"\" ]; then\n" +
		"  printf '%s\\n' \"$@\" > \"$SUBFLOW_MOCK_ARGS_FILE\"\n" +
		"fi\n" +
		"if [ \"${SUBFLOW_MOCK_FAIL:-0}\" = \"1\" ]; then\n" +
		"  echo mock fail 1>&2\n" +
		"  exit 1\n" +
		"fi\n" +
		"if [ \"$tool\" = \"ffprobe\" ]; then\n" +
		"  d=\"${SUBFLOW_MOCK_DURATION:-1.0}\"\n" +
		"  printf '{\\\"format\\\":{\\\"duration\\\":\\\"%s\\\"}}' \"$d\"\n" +
		"  exit 0\n" +
		"fi\n" +
		"out=\"\"\n" +
		"for out in \"$@\"; do :; done\n" +
		"mkdir -p \"$(dirname \"$out\")\"\n" +
		"printf 'RIFF0000WAVE' > \"$out\"\n" +
		"exit 0\n"

	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatalf("write mock tool: %v", err)
	}
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("chmod mock tool: %v", err)
	}
	return path
}

func splitLines(b []byte) [][]byte {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return nil
	}
	return bytes.Split(b, []byte("\n"))
}

func joinLines(lines [][]byte) string {
	if len(lines) == 0 {
		return ""
	}
	return string(bytes.Join(lines, []byte("\n")))
}
