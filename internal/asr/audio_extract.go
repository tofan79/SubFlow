package asr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

const (
	defaultExtractTimeout = 2 * time.Minute
	defaultProbeTimeout   = 15 * time.Second
)

func ExtractAudio(ctx context.Context, videoPath, outputPath string) error {
	if strings.TrimSpace(videoPath) == "" {
		return wrapASRExtractErr(fmt.Errorf("video path is empty"))
	}
	if strings.TrimSpace(outputPath) == "" {
		return wrapASRExtractErr(fmt.Errorf("output path is empty"))
	}

	ffmpegPath, err := FindFFmpeg()
	if err != nil {
		return wrapASRExtractErr(err)
	}

	if dir := filepath.Dir(outputPath); dir != "." && dir != "" {
		if mkErr := os.MkdirAll(dir, 0o755); mkErr != nil {
			return wrapASRExtractErr(fmt.Errorf("create output dir: %w", mkErr))
		}
	}

	extractCtx := ctx
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		extractCtx, cancel = context.WithTimeout(ctx, defaultExtractTimeout)
		defer cancel()
	}

	args := []string{"-i", videoPath, "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1", "-y", outputPath}
	cmd := exec.CommandContext(extractCtx, ffmpegPath, args...)
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined

	if runErr := cmd.Run(); runErr != nil {
		if errors.Is(extractCtx.Err(), context.DeadlineExceeded) || errors.Is(extractCtx.Err(), context.Canceled) {
			return wrapASRExtractErr(extractCtx.Err())
		}
		return wrapASRExtractErr(fmt.Errorf("ffmpeg failed: %w: %s", runErr, strings.TrimSpace(combined.String())))
	}

	if _, statErr := os.Stat(outputPath); statErr != nil {
		return wrapASRExtractErr(fmt.Errorf("output not created: %w", statErr))
	}

	return nil
}

func GetAudioDuration(audioPath string) (float64, error) {
	if strings.TrimSpace(audioPath) == "" {
		return 0, wrapASRExtractErr(fmt.Errorf("audio path is empty"))
	}

	ffprobePath, err := FindFFprobe()
	if err != nil {
		return 0, wrapASRExtractErr(err)
	}

	probeCtx, cancel := context.WithTimeout(context.Background(), defaultProbeTimeout)
	defer cancel()

	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		audioPath,
	}
	cmd := exec.CommandContext(probeCtx, ffprobePath, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if runErr := cmd.Run(); runErr != nil {
		if errors.Is(probeCtx.Err(), context.DeadlineExceeded) || errors.Is(probeCtx.Err(), context.Canceled) {
			return 0, wrapASRExtractErr(probeCtx.Err())
		}
		return 0, wrapASRExtractErr(fmt.Errorf("ffprobe failed: %w: %s", runErr, strings.TrimSpace(stderr.String())))
	}

	var parsed struct {
		Format struct {
			Duration string `json:"duration"`
		} `json:"format"`
	}
	if unmarshalErr := json.Unmarshal(out.Bytes(), &parsed); unmarshalErr != nil {
		return 0, wrapASRExtractErr(fmt.Errorf("parse ffprobe output: %w", unmarshalErr))
	}
	if strings.TrimSpace(parsed.Format.Duration) == "" {
		return 0, wrapASRExtractErr(fmt.Errorf("ffprobe duration missing"))
	}
	secs, parseErr := strconv.ParseFloat(parsed.Format.Duration, 64)
	if parseErr != nil {
		return 0, wrapASRExtractErr(fmt.Errorf("parse duration: %w", parseErr))
	}
	if secs < 0 {
		return 0, wrapASRExtractErr(fmt.Errorf("negative duration"))
	}

	return secs, nil
}

func FindFFmpeg() (string, error) {
	p, err := findBinary("ffmpeg")
	if err != nil {
		return "", wrapASRExtractErr(err)
	}
	return p, nil
}

func FindFFprobe() (string, error) {
	p, err := findBinary("ffprobe")
	if err != nil {
		return "", wrapASRExtractErr(err)
	}
	return p, nil
}

func findBinary(name string) (string, error) {
	if p, err := exec.LookPath(name); err == nil && p != "" {
		return p, nil
	}

	candidates := commonBinaryCandidates(name)
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if isExecutableFile(p) {
			return p, nil
		}
	}

	return "", fmt.Errorf("%s not found", name)
}

func isExecutableFile(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	if st.IsDir() {
		return false
	}
	if runtime.GOOS == "windows" {
		return true
	}
	return st.Mode()&0o111 != 0
}

func wrapASRExtractErr(err error) error {
	if err == nil {
		return nil
	}
	var perr *pipeline.Error
	if errors.As(err, &perr) && perr.Code == pipeline.ErrASRExtractFail {
		return err
	}
	return pipeline.ErrASRExtractFailErr(err)
}

func commonBinaryCandidates(name string) []string {
	var paths []string

	if runtime.GOOS != "windows" {
		paths = append(paths,
			"/usr/bin/"+name,
			"/usr/local/bin/"+name,
			"/bin/"+name,
		)
	} else {
		exe := name + ".exe"
		paths = append(paths,
			filepath.Join(os.Getenv("ProgramFiles"), "FFmpeg", "bin", exe),
			filepath.Join(os.Getenv("ProgramFiles(x86)"), "FFmpeg", "bin", exe),
		)
	}

	if exePath, err := os.Executable(); err == nil {
		base := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(base, "resources", "bin", name),
			filepath.Join(base, "bin", name),
			filepath.Join(base, name),
		)
		if runtime.GOOS == "windows" {
			exe := name + ".exe"
			paths = append(paths,
				filepath.Join(base, "resources", "bin", exe),
				filepath.Join(base, "bin", exe),
				filepath.Join(base, exe),
			)
		}
	}

	return paths
}
