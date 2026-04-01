// Package asr provides automatic speech recognition functionality.
// It supports multiple backends: faster-whisper (local), Groq (cloud), and Deepgram (cloud).
package asr

import (
	"context"
)

// =============================================================================
// ASR OPTIONS
// =============================================================================

// Opts holds options for ASR transcription.
type Opts struct {
	Language    string // ISO 639-1 code, e.g., "en", "ja", or "auto"
	Model       string // e.g., "base", "small", "medium", "large-v3"
	Backend     string // "auto", "cpu", "cuda", "rocm", "coreml", "openvino"
	ComputeType string // "int8", "float16", "float32"
}

// =============================================================================
// ASR SEGMENT
// =============================================================================

// Segment represents a single transcribed segment.
type Segment struct {
	Index      int     // Segment index
	StartMS    int64   // Start time in milliseconds
	EndMS      int64   // End time in milliseconds
	Text       string  // Transcribed text
	Confidence float64 // Confidence score (0.0 - 1.0)
	Language   string  // Detected language if auto-detect
}

// =============================================================================
// ASR PROVIDER INTERFACE
// =============================================================================

// Provider defines the interface for ASR services.
type Provider interface {
	// Transcribe transcribes audio and streams segments.
	// Returns two channels: segments and errors.
	// Segments are sent as they're transcribed.
	// The error channel receives any errors and is closed when done.
	Transcribe(ctx context.Context, audioPath string, opts Opts) (<-chan Segment, <-chan error)

	// EstimateCost estimates the cost in USD for the given duration.
	EstimateCost(durationSeconds float64) float64

	// Name returns the provider name.
	Name() string
}

// =============================================================================
// HARDWARE BACKEND
// =============================================================================

// Backend represents available ASR hardware backends.
type Backend string

const (
	BackendAuto     Backend = "auto"
	BackendCPU      Backend = "cpu"
	BackendCUDA     Backend = "cuda"
	BackendROCm     Backend = "rocm"
	BackendCoreML   Backend = "coreml"
	BackendOpenVINO Backend = "openvino"
)

// AllBackends returns all valid backend values.
func AllBackends() []Backend {
	return []Backend{
		BackendAuto,
		BackendCPU,
		BackendCUDA,
		BackendROCm,
		BackendCoreML,
		BackendOpenVINO,
	}
}

// IsValid checks if the backend is a valid value.
func (b Backend) IsValid() bool {
	switch b {
	case BackendAuto, BackendCPU, BackendCUDA, BackendROCm, BackendCoreML, BackendOpenVINO:
		return true
	default:
		return false
	}
}

// String returns the string representation of the backend.
func (b Backend) String() string {
	return string(b)
}

// =============================================================================
// COMPUTE TYPE
// =============================================================================

// ComputeType represents the precision level for ASR computation.
type ComputeType string

const (
	ComputeInt8    ComputeType = "int8"
	ComputeFloat16 ComputeType = "float16"
	ComputeFloat32 ComputeType = "float32"
)

// ComputeTypeForBackend returns the recommended compute type for a backend.
// CUDA, ROCm, CoreML -> float16
// OpenVINO, CPU -> int8
func ComputeTypeForBackend(backend Backend) ComputeType {
	switch backend {
	case BackendCUDA, BackendROCm, BackendCoreML:
		return ComputeFloat16
	case BackendOpenVINO, BackendCPU:
		return ComputeInt8
	default:
		return ComputeInt8
	}
}

// =============================================================================
// MODEL SIZES
// =============================================================================

// Model represents a Whisper model variant.
type Model string

const (
	ModelBase    Model = "base"
	ModelSmall   Model = "small"
	ModelMedium  Model = "medium"
	ModelLargeV3 Model = "large-v3"
)

// ModelInfo holds information about a Whisper model.
type ModelInfo struct {
	Name        Model
	SizeMB      int
	Quality     string // "medium", "good", "high", "best"
	Bundled     bool   // Included in default installation
	GPUOnly     bool   // Requires GPU for reasonable performance
	DownloadURL string
}

// AllModels returns information about all available models.
func AllModels() []ModelInfo {
	return []ModelInfo{
		{Name: ModelBase, SizeMB: 145, Quality: "medium", Bundled: true, GPUOnly: false},
		{Name: ModelSmall, SizeMB: 465, Quality: "good", Bundled: false, GPUOnly: false},
		{Name: ModelMedium, SizeMB: 1400, Quality: "high", Bundled: false, GPUOnly: false},
		{Name: ModelLargeV3, SizeMB: 2900, Quality: "best", Bundled: false, GPUOnly: true},
	}
}

// =============================================================================
// HARDWARE INFO
// =============================================================================

// HardwareInfo holds detected hardware capabilities.
type HardwareInfo struct {
	Backend     Backend `json:"backend"`
	GPUName     string  `json:"gpuName"`
	CUDAVersion string  `json:"cudaVersion"`
	ROCmVersion string  `json:"rocmVersion"`
	VRAMTotal   int64   `json:"vramTotal"`   // bytes
	VRAMFree    int64   `json:"vramFree"`    // bytes
	ComputeType string  `json:"computeType"` // Recommended compute type
}

// =============================================================================
// PROGRESS EVENT
// =============================================================================

// ProgressEvent represents a progress update during transcription.
type ProgressEvent struct {
	Type    string  `json:"type"`    // "segment", "progress", "done", "error"
	Percent float64 `json:"percent"` // 0.0 - 100.0 for "progress" type
}

// =============================================================================
// PROVIDER NAMES
// =============================================================================

const (
	ProviderWhisperLocal = "whisper-local"
	ProviderGroq         = "groq"
	ProviderDeepgram     = "deepgram"
)

// =============================================================================
// DEFAULT VALUES
// =============================================================================

const (
	DefaultModel       = ModelBase
	DefaultLanguage    = "auto"
	DefaultComputeType = ComputeInt8
)
