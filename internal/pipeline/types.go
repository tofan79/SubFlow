// Package pipeline provides types and interfaces for the subtitle processing pipeline.
package pipeline

import (
	"context"
)

// =============================================================================
// PIPELINE STATES
// =============================================================================

// State represents the current state of a pipeline.
type State string

const (
	StateImported       State = "imported"
	StateAudioExtracted State = "audio_extracted"
	StateASRDone        State = "asr_done"
	StateCorrected      State = "corrected"
	StateContextDone    State = "context_done"
	StateSegmented      State = "segmented"
	StateTranslated     State = "translated"
	StateRewritten      State = "rewritten"
	StateQADone         State = "qa_done"
	StateExported       State = "exported"
)

// StateOrder defines the valid order of pipeline states.
var StateOrder = []State{
	StateImported,
	StateAudioExtracted,
	StateASRDone,
	StateCorrected,
	StateContextDone,
	StateSegmented,
	StateTranslated,
	StateRewritten,
	StateQADone,
	StateExported,
}

// IsValidTransition checks if transitioning from `from` to `to` is valid.
func IsValidTransition(from, to State) bool {
	fromIdx := -1
	toIdx := -1
	for i, s := range StateOrder {
		if s == from {
			fromIdx = i
		}
		if s == to {
			toIdx = i
		}
	}
	// Can only move forward in the pipeline
	return fromIdx >= 0 && toIdx >= 0 && toIdx == fromIdx+1
}

// =============================================================================
// TRANSLATION PROVIDER INTERFACE
// =============================================================================

// TranslationOpts holds options for translation.
type TranslationOpts struct {
	SourceLang  string            // e.g., "en", "ja"
	TargetLang  string            // e.g., "id"
	ContentMode string            // e.g., "movie", "documentary", "anime"
	Glossary    []GlossaryTerm    // Terms to preserve
	Context     map[string]string // Additional context
}

// GlossaryTerm represents a term that should be translated consistently.
type GlossaryTerm struct {
	SourceTerm    string
	TargetTerm    string
	CaseSensitive bool
}

// TranslationProvider defines the interface for translation services.
type TranslationProvider interface {
	// Translate translates a batch of strings.
	// Input and output slices have the same length.
	Translate(ctx context.Context, batch []string, opts TranslationOpts) ([]string, error)

	// EstimateCost estimates the cost in USD for translating the given character count.
	EstimateCost(charCount int) float64

	// MaxBatchSize returns the maximum number of strings per batch.
	MaxBatchSize() int

	// Name returns the provider name (e.g., "deepl", "openai").
	Name() string
}

// =============================================================================
// REWRITE PROVIDER INTERFACE
// =============================================================================

// RewriteInput holds input for a single rewrite operation.
type RewriteInput struct {
	Source     string // Original text
	Translated string // Layer 1 translation
	Speaker    string // Speaker name/ID if known
	Emotion    string // Detected emotion if available
	Context    string // Scene context
}

// RewriteOpts holds options for rewrite.
type RewriteOpts struct {
	TonePreset      string         // "natural", "formal", "casual", "cinematic"
	MaxCharsPerLine int            // Default: 42
	MaxLines        int            // Default: 2
	MaxCPS          float64        // Default: 17.0
	Glossary        []GlossaryTerm // Terms to preserve
}

// RewriteProvider defines the interface for AI rewrite services (Layer 2).
type RewriteProvider interface {
	// Rewrite rewrites a batch of translations for naturalness.
	Rewrite(ctx context.Context, batch []RewriteInput, opts RewriteOpts) ([]string, error)

	// EstimateCost estimates the cost in USD for the given token count.
	EstimateCost(tokenCount int) float64

	// Name returns the provider name.
	Name() string
}

// =============================================================================
// ASR PROVIDER INTERFACE
// =============================================================================

// ASROpts holds options for ASR.
type ASROpts struct {
	Language    string // ISO 639-1 code, e.g., "en", "ja", or "auto"
	Model       string // e.g., "base", "small", "medium", "large-v3"
	Backend     string // "auto", "cpu", "cuda", "rocm", "coreml", "openvino"
	ComputeType string // "int8", "float16", "float32"
}

// ASRSegment represents a single transcribed segment.
type ASRSegment struct {
	Index      int     // Segment index
	StartMS    int64   // Start time in milliseconds
	EndMS      int64   // End time in milliseconds
	Text       string  // Transcribed text
	Confidence float64 // Confidence score (0.0 - 1.0)
	Language   string  // Detected language if auto-detect
}

// ASRProvider defines the interface for ASR services.
type ASRProvider interface {
	// Transcribe transcribes audio and streams segments.
	// Returns two channels: segments and errors.
	// Segments are sent as they're transcribed.
	// The error channel receives any errors and is closed when done.
	Transcribe(ctx context.Context, audioPath string, opts ASROpts) (<-chan ASRSegment, <-chan error)

	// EstimateCost estimates the cost in USD for the given duration.
	EstimateCost(durationSeconds float64) float64

	// Name returns the provider name.
	Name() string
}

// =============================================================================
// SUBTITLE CARD
// =============================================================================

// SubtitleCard represents a single subtitle entry for processing.
type SubtitleCard struct {
	ID       string // Unique identifier
	Index    int    // Display order
	StartMS  int64  // Start time in milliseconds
	EndMS    int64  // End time in milliseconds
	Source   string // Original/ASR text
	L1       string // Layer 1 translation
	L2       string // Layer 2 rewrite (final)
	Speaker  string // Speaker identifier
	Emotion  string // Detected emotion
	QAStatus string // "pass", "warn", "error", "pending"
}

// Duration returns the duration in milliseconds.
func (c *SubtitleCard) Duration() int64 {
	return c.EndMS - c.StartMS
}

// DurationSeconds returns the duration in seconds.
func (c *SubtitleCard) DurationSeconds() float64 {
	return float64(c.EndMS-c.StartMS) / 1000.0
}

// =============================================================================
// PIPELINE PROGRESS
// =============================================================================

// Progress represents the current progress of a pipeline step.
type Progress struct {
	ProjectID string  `json:"projectId"`
	Step      string  `json:"step"`    // Current step name
	Current   int     `json:"current"` // Current item
	Total     int     `json:"total"`   // Total items
	Percent   float64 `json:"percent"` // 0.0 - 100.0
	Message   string  `json:"message"` // Human-readable status
}

// =============================================================================
// COST ESTIMATE
// =============================================================================

// CostEstimate represents estimated costs for pipeline operations.
type CostEstimate struct {
	ASR       float64 `json:"asr"`
	Translate float64 `json:"translate"`
	Rewrite   float64 `json:"rewrite"`
	Total     float64 `json:"total"`
	Currency  string  `json:"currency"` // Always "USD"
}

// =============================================================================
// HARDWARE INFO
// =============================================================================

// HardwareBackend represents available ASR hardware backends.
type HardwareBackend string

const (
	BackendAuto     HardwareBackend = "auto"
	BackendCPU      HardwareBackend = "cpu"
	BackendCUDA     HardwareBackend = "cuda"
	BackendROCm     HardwareBackend = "rocm"
	BackendCoreML   HardwareBackend = "coreml"
	BackendOpenVINO HardwareBackend = "openvino"
)

// HardwareInfo holds detected hardware capabilities.
type HardwareInfo struct {
	Backend     HardwareBackend `json:"backend"`
	GPUName     string          `json:"gpuName"`
	CUDAVersion string          `json:"cudaVersion"`
	ROCmVersion string          `json:"rocmVersion"`
	VRAMTotal   int64           `json:"vramTotal"`   // bytes
	VRAMFree    int64           `json:"vramFree"`    // bytes
	ComputeType string          `json:"computeType"` // Recommended compute type
}

// =============================================================================
// TONE PRESETS
// =============================================================================

// TonePreset represents a predefined tone for Layer 2 rewrite.
type TonePreset string

const (
	ToneNatural   TonePreset = "natural"   // Balanced, real conversation
	ToneFormal    TonePreset = "formal"    // Proper, complete sentences, EYD
	ToneCasual    TonePreset = "casual"    // Relaxed, may use gue/lo/nih
	ToneCinematic TonePreset = "cinematic" // Dramatic, strong diction, emotional
)

// ToneDescriptions provides human-readable descriptions for tone presets.
var ToneDescriptions = map[TonePreset]string{
	ToneNatural:   "Seimbang, percakapan nyata",
	ToneFormal:    "Baku, kalimat lengkap, sesuai EYD",
	ToneCasual:    "Santai, boleh gue/lo/nih",
	ToneCinematic: "Dramatis, diksi kuat, emosi ditonjolkan",
}

// =============================================================================
// CONTENT MODES
// =============================================================================

// ContentMode represents the type of content being processed.
type ContentMode string

const (
	ContentMovie       ContentMode = "movie"
	ContentTVSeries    ContentMode = "tv_series"
	ContentDocumentary ContentMode = "documentary"
	ContentAnime       ContentMode = "anime"
	ContentNews        ContentMode = "news"
	ContentEducation   ContentMode = "education"
	ContentYouTube     ContentMode = "youtube"
)
