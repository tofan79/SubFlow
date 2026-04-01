package asr

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/subflow/subflow/internal/pipeline"
)

const (
	whisperBinaryName         = "subflow-whisper"
	defaultWhisperTimeout     = 30 * time.Minute
	whisperOutputTypeSegment  = "segment"
	whisperOutputTypeProgress = "progress"
	whisperOutputTypeDone     = "done"
	whisperOutputTypeError    = "error"
)

type WhisperLocalProvider struct {
	binaryPath string
}

func NewWhisperLocalProvider() (*WhisperLocalProvider, error) {
	path, err := findWhisperBinary()
	if err != nil {
		return nil, err
	}
	return &WhisperLocalProvider{binaryPath: path}, nil
}

func NewWhisperLocalProviderWithPath(binaryPath string) *WhisperLocalProvider {
	return &WhisperLocalProvider{binaryPath: binaryPath}
}

func (p *WhisperLocalProvider) Name() string { return ProviderWhisperLocal }

func (p *WhisperLocalProvider) EstimateCost(_ float64) float64 { return 0 }

func (p *WhisperLocalProvider) Transcribe(ctx context.Context, audioPath string, opts Opts) (<-chan Segment, <-chan error) {
	segCh := make(chan Segment, 64)
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

		if strings.TrimSpace(audioPath) == "" {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "audio path is empty", nil)
			return
		}

		if _, err := os.Stat(audioPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				errCh <- pipeline.ErrFileNotFound(audioPath, err)
				return
			}
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, fmt.Sprintf("failed to stat audio file: %s", audioPath), err)
			return
		}

		if strings.TrimSpace(p.binaryPath) == "" {
			errCh <- pipeline.ErrASRNotInstalledErr("")
			return
		}
		if _, err := os.Stat(p.binaryPath); err != nil {
			errCh <- pipeline.ErrASRNotInstalledErr(p.binaryPath)
			return
		}

		execCtx := ctx
		if _, hasDeadline := ctx.Deadline(); !hasDeadline {
			var cancel context.CancelFunc
			execCtx, cancel = context.WithTimeout(ctx, defaultWhisperTimeout)
			defer cancel()
		}

		args := p.buildArgs(audioPath, opts)
		cmd := exec.CommandContext(execCtx, p.binaryPath, args...)
		cmd.Stderr = nil

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to create stdout pipe", err)
			return
		}

		if err := cmd.Start(); err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "failed to start whisper subprocess", err)
			return
		}

		scanner := bufio.NewScanner(stdout)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

		segmentIndex := 0
		for scanner.Scan() {
			if execCtx.Err() != nil {
				errCh <- execCtx.Err()
				_ = cmd.Process.Kill()
				return
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			var msg whisperOutputMessage
			if err := json.Unmarshal([]byte(line), &msg); err != nil {
				continue
			}

			switch msg.Type {
			case whisperOutputTypeSegment:
				seg := Segment{
					Index:      segmentIndex,
					StartMS:    secondsToMS(msg.Start),
					EndMS:      secondsToMS(msg.End),
					Text:       msg.Text,
					Confidence: msg.Confidence,
					Language:   msg.Language,
				}
				segmentIndex++
				select {
				case segCh <- seg:
				case <-execCtx.Done():
					errCh <- execCtx.Err()
					_ = cmd.Process.Kill()
					return
				}
			case whisperOutputTypeError:
				errCh <- pipeline.NewError(msg.Code, msg.Message, nil)
				_ = cmd.Process.Kill()
				return
			case whisperOutputTypeDone, whisperOutputTypeProgress:
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "error reading whisper output", err)
			_ = cmd.Process.Kill()
			return
		}

		if err := cmd.Wait(); err != nil {
			if execCtx.Err() != nil {
				if errors.Is(execCtx.Err(), context.DeadlineExceeded) {
					errCh <- pipeline.ErrASRTimeoutErr(int(defaultWhisperTimeout.Seconds()))
					return
				}
				errCh <- execCtx.Err()
				return
			}
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, fmt.Sprintf("whisper exited with code %d", exitErr.ExitCode()), err)
				return
			}
			errCh <- pipeline.NewError(pipeline.ErrASRExtractFail, "whisper subprocess failed", err)
			return
		}
	}()

	return segCh, errCh
}

func (p *WhisperLocalProvider) buildArgs(audioPath string, opts Opts) []string {
	args := []string{"--audio", audioPath}

	model := strings.TrimSpace(opts.Model)
	if model == "" {
		model = string(DefaultModel)
	}
	args = append(args, "--model", model)

	backend := strings.TrimSpace(opts.Backend)
	if backend == "" || backend == string(BackendAuto) {
		detected := DetectBestBackend()
		backend = string(detected)
	}
	args = append(args, "--backend", backend)

	computeType := strings.TrimSpace(opts.ComputeType)
	if computeType == "" {
		computeType = string(ComputeTypeForBackend(Backend(backend)))
	}
	args = append(args, "--compute-type", computeType)

	language := strings.TrimSpace(opts.Language)
	if language == "" {
		language = DefaultLanguage
	}
	args = append(args, "--language", language)

	return args
}

type whisperOutputMessage struct {
	Type       string  `json:"type"`
	Start      float64 `json:"start,omitempty"`
	End        float64 `json:"end,omitempty"`
	Text       string  `json:"text,omitempty"`
	Confidence float64 `json:"confidence,omitempty"`
	Language   string  `json:"language,omitempty"`
	Percent    float64 `json:"percent,omitempty"`
	Code       string  `json:"code,omitempty"`
	Message    string  `json:"message,omitempty"`
}

func findWhisperBinary() (string, error) {
	if p, err := exec.LookPath(whisperBinaryName); err == nil && p != "" {
		return p, nil
	}

	candidates := whisperBinaryCandidates()
	for _, p := range candidates {
		if p == "" {
			continue
		}
		if isExecutableFile(p) {
			return p, nil
		}
	}

	return "", pipeline.ErrASRNotInstalledErr(whisperBinaryName)
}

func whisperBinaryCandidates() []string {
	var paths []string

	if exePath, err := os.Executable(); err == nil {
		base := filepath.Dir(exePath)
		paths = append(paths,
			filepath.Join(base, "resources", "bin", whisperBinaryName),
			filepath.Join(base, "bin", whisperBinaryName),
			filepath.Join(base, whisperBinaryName),
		)
		if runtime.GOOS == "windows" {
			exe := whisperBinaryName + ".exe"
			paths = append(paths,
				filepath.Join(base, "resources", "bin", exe),
				filepath.Join(base, "bin", exe),
				filepath.Join(base, exe),
			)
		}
	}

	if runtime.GOOS != "windows" {
		paths = append(paths,
			"/usr/local/bin/"+whisperBinaryName,
			"/usr/bin/"+whisperBinaryName,
		)
	} else {
		exe := whisperBinaryName + ".exe"
		if pf := os.Getenv("ProgramFiles"); pf != "" {
			paths = append(paths, filepath.Join(pf, "SubFlow", exe))
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			paths = append(paths, filepath.Join(localAppData, "SubFlow", exe))
		}
	}

	return paths
}
