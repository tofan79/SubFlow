package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds application state and provides IPC methods to frontend.
// All public methods are automatically bound and callable from JavaScript.
type App struct {
	ctx context.Context
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved for runtime calls.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// domReady is called after the frontend DOM has been loaded.
func (a *App) domReady(ctx context.Context) {
	// Emit initial hardware info to frontend
	runtime.EventsEmit(a.ctx, "asr:hardware", map[string]interface{}{
		"backend":     "detecting",
		"gpuName":     "",
		"cudaVersion": "",
	})
}

// shutdown is called when the app is closing.
func (a *App) shutdown(ctx context.Context) {
	// Future: cleanup resources, close DB connections
}

// =============================================================================
// PROJECT MANAGEMENT
// =============================================================================

// ProjectInfo represents a project in the database.
type ProjectInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SourcePath   string `json:"sourcePath"`
	State        string `json:"state"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
	SegmentCount int    `json:"segmentCount"`
}

// GetProjects returns all projects from the database.
func (a *App) GetProjects() ([]ProjectInfo, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetProjects")
}

// GetProject returns a single project by ID.
func (a *App) GetProject(id string) (*ProjectInfo, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetProject")
}

// CreateProject creates a new project from a source file.
func (a *App) CreateProject(sourcePath string, name string) (*ProjectInfo, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: CreateProject")
}

// DeleteProject removes a project and all its data.
func (a *App) DeleteProject(id string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: DeleteProject")
}

// =============================================================================
// PIPELINE CONTROL
// =============================================================================

// PipelineConfig holds configuration for a pipeline run.
type PipelineConfig struct {
	ProjectID         string `json:"projectId"`
	SourceLang        string `json:"sourceLang"`
	TargetLang        string `json:"targetLang"`
	ContentMode       string `json:"contentMode"`
	TonePreset        string `json:"tonePreset"`
	ASRProvider       string `json:"asrProvider"`
	ASRModel          string `json:"asrModel"`
	TranslateProvider string `json:"translateProvider"`
	RewriteProvider   string `json:"rewriteProvider"`
}

// StartPipeline begins the subtitle pipeline for a project.
func (a *App) StartPipeline(config PipelineConfig) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: StartPipeline")
}

// PausePipeline pauses the running pipeline.
func (a *App) PausePipeline(projectID string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: PausePipeline")
}

// ResumePipeline resumes a paused pipeline.
func (a *App) ResumePipeline(projectID string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: ResumePipeline")
}

// CancelPipeline cancels the running pipeline.
func (a *App) CancelPipeline(projectID string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: CancelPipeline")
}

// GetPipelineState returns the current pipeline state for a project.
func (a *App) GetPipelineState(projectID string) (string, error) {
	return "", fmt.Errorf("ERR_NOT_IMPLEMENTED: GetPipelineState")
}

// =============================================================================
// SETTINGS MANAGEMENT
// =============================================================================

// SettingsData holds all application settings.
type SettingsData struct {
	// API Keys (encrypted in storage)
	DeepLAPIKey     string `json:"deeplApiKey"`
	OpenAIAPIKey    string `json:"openaiApiKey"`
	AnthropicAPIKey string `json:"anthropicApiKey"`
	GeminiAPIKey    string `json:"geminiApiKey"`
	GroqAPIKey      string `json:"groqApiKey"`
	DeepgramAPIKey  string `json:"deepgramApiKey"`
	XAIAPIKey       string `json:"xaiApiKey"`
	QwenAPIKey      string `json:"qwenApiKey"`
	OllamaEndpoint  string `json:"ollamaEndpoint"`

	// ASR Settings
	ASRBackend   string `json:"asrBackend"`   // auto, cpu, cuda, rocm, coreml, openvino
	WhisperModel string `json:"whisperModel"` // tiny, base, small, medium, large-v3
	PreferredASR string `json:"preferredAsr"` // local, groq, deepgram

	// Translation Settings
	DefaultSourceLang string `json:"defaultSourceLang"`
	DefaultTargetLang string `json:"defaultTargetLang"`
	DefaultTonePreset string `json:"defaultTonePreset"` // natural, formal, casual, cinematic

	// QA Settings
	MaxCharsPerLine int     `json:"maxCharsPerLine"` // default: 42
	MaxLines        int     `json:"maxLines"`        // default: 2
	MaxCPS          float64 `json:"maxCps"`          // default: 17.0
	MinGapMS        int64   `json:"minGapMs"`        // default: 83

	// Export Settings
	DefaultExportFormat string `json:"defaultExportFormat"` // srt, vtt, ass, txt
}

// GetSettings returns all application settings.
func (a *App) GetSettings() (*SettingsData, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetSettings")
}

// SetSetting updates a single setting.
func (a *App) SetSetting(key string, value interface{}) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: SetSetting")
}

// SaveSettings saves all settings at once.
func (a *App) SaveSettings(settings SettingsData) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: SaveSettings")
}

// =============================================================================
// ASR (AUTOMATIC SPEECH RECOGNITION)
// =============================================================================

// HardwareInfo represents detected hardware capabilities.
type HardwareInfo struct {
	Backend     string `json:"backend"` // cpu, cuda, rocm, coreml, openvino
	GPUName     string `json:"gpuName"`
	CUDAVersion string `json:"cudaVersion"`
	ROCmVersion string `json:"rocmVersion"`
	VRAMTotal   int64  `json:"vramTotal"`   // bytes
	VRAMFree    int64  `json:"vramFree"`    // bytes
	ComputeType string `json:"computeType"` // int8, float16, float32
}

// DetectHardware detects available ASR hardware backends.
func (a *App) DetectHardware() (*HardwareInfo, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: DetectHardware")
}

// GetAvailableModels returns available Whisper models.
func (a *App) GetAvailableModels() ([]string, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetAvailableModels")
}

// DownloadModel downloads a Whisper model.
func (a *App) DownloadModel(modelName string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: DownloadModel")
}

// =============================================================================
// SEGMENT / SUBTITLE MANAGEMENT
// =============================================================================

// Segment represents a single subtitle segment.
type Segment struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Index     int    `json:"index"`
	StartMS   int64  `json:"startMs"`
	EndMS     int64  `json:"endMs"`
	Source    string `json:"source"` // Original/ASR text
	L1        string `json:"l1"`     // Layer 1 translation
	L2        string `json:"l2"`     // Layer 2 rewrite
	Speaker   string `json:"speaker"`
	Emotion   string `json:"emotion"`
	QAStatus  string `json:"qaStatus"` // pass, warn, error, pending
}

// GetSegments returns all segments for a project.
func (a *App) GetSegments(projectID string) ([]Segment, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetSegments")
}

// UpdateSegment updates a single segment.
func (a *App) UpdateSegment(segment Segment) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: UpdateSegment")
}

// SplitSegment splits a segment at the given position.
func (a *App) SplitSegment(segmentID string, splitAtMS int64, splitAtChar int) ([]Segment, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: SplitSegment")
}

// MergeSegments merges two adjacent segments.
func (a *App) MergeSegments(segmentID1 string, segmentID2 string) (*Segment, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: MergeSegments")
}

// DeleteSegment deletes a segment.
func (a *App) DeleteSegment(segmentID string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: DeleteSegment")
}

// RetryL1 re-runs Layer 1 translation for a segment.
func (a *App) RetryL1(segmentID string) (*Segment, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: RetryL1")
}

// RetryL2 re-runs Layer 2 rewrite for a segment.
func (a *App) RetryL2(segmentID string) (*Segment, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: RetryL2")
}

// =============================================================================
// QA ENGINE
// =============================================================================

// QAResult represents a single QA check result.
type QAResult struct {
	CardID    string `json:"cardId"`
	CheckID   string `json:"checkId"`
	Passed    bool   `json:"passed"`
	Severity  string `json:"severity"` // error, warning
	Detail    string `json:"detail"`
	AutoFixed bool   `json:"autoFixed"`
	FixAction string `json:"fixAction"`
}

// QAReport represents a full QA report.
type QAReport struct {
	RunAt      int64      `json:"runAt"`
	TotalCards int        `json:"totalCards"`
	Passed     int        `json:"passed"`
	Warnings   int        `json:"warnings"`
	Errors     int        `json:"errors"`
	AutoFixed  int        `json:"autoFixed"`
	Results    []QAResult `json:"results"`
}

// RunQA runs QA checks on all segments of a project.
func (a *App) RunQA(projectID string) (*QAReport, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: RunQA")
}

// RunQAAutoFix runs QA with auto-fix enabled.
func (a *App) RunQAAutoFix(projectID string) (*QAReport, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: RunQAAutoFix")
}

// =============================================================================
// GLOSSARY MANAGEMENT
// =============================================================================

// GlossaryTerm represents a glossary entry.
type GlossaryTerm struct {
	ID            string `json:"id"`
	SourceTerm    string `json:"sourceTerm"`
	TargetTerm    string `json:"targetTerm"`
	CaseSensitive bool   `json:"caseSensitive"`
	Notes         string `json:"notes"`
}

// GetGlossary returns all glossary terms.
func (a *App) GetGlossary() ([]GlossaryTerm, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetGlossary")
}

// AddGlossaryTerm adds a new glossary term.
func (a *App) AddGlossaryTerm(term GlossaryTerm) (*GlossaryTerm, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: AddGlossaryTerm")
}

// UpdateGlossaryTerm updates an existing glossary term.
func (a *App) UpdateGlossaryTerm(term GlossaryTerm) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: UpdateGlossaryTerm")
}

// DeleteGlossaryTerm deletes a glossary term.
func (a *App) DeleteGlossaryTerm(id string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: DeleteGlossaryTerm")
}

// ImportGlossary imports glossary from JSON file.
func (a *App) ImportGlossary(filePath string) (int, error) {
	return 0, fmt.Errorf("ERR_NOT_IMPLEMENTED: ImportGlossary")
}

// ExportGlossary exports glossary to JSON file.
func (a *App) ExportGlossary(filePath string) error {
	return fmt.Errorf("ERR_NOT_IMPLEMENTED: ExportGlossary")
}

// =============================================================================
// EXPORT
// =============================================================================

// ExportOptions holds export configuration.
type ExportOptions struct {
	ProjectID    string `json:"projectId"`
	Format       string `json:"format"` // srt, vtt, ass, txt
	OutputDir    string `json:"outputDir"`
	Layer        string `json:"layer"`        // source, l1, l2
	DualSubtitle bool   `json:"dualSubtitle"` // source + l2
}

// Export exports subtitles to the specified format.
func (a *App) Export(options ExportOptions) (string, error) {
	return "", fmt.Errorf("ERR_NOT_IMPLEMENTED: Export")
}

// =============================================================================
// FILE DIALOGS
// =============================================================================

// SelectFile opens a file selection dialog.
func (a *App) SelectFile(title string, filters []string) (string, error) {
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{DisplayName: "Supported Files", Pattern: "*.mp4;*.mkv;*.avi;*.mov;*.srt;*.vtt;*.ass"},
			{DisplayName: "Video Files", Pattern: "*.mp4;*.mkv;*.avi;*.mov;*.webm"},
			{DisplayName: "Subtitle Files", Pattern: "*.srt;*.vtt;*.ass;*.ssa"},
			{DisplayName: "All Files", Pattern: "*.*"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("App.SelectFile: %w", err)
	}
	return result, nil
}

// SelectDirectory opens a directory selection dialog.
func (a *App) SelectDirectory(title string) (string, error) {
	result, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return "", fmt.Errorf("App.SelectDirectory: %w", err)
	}
	return result, nil
}

// =============================================================================
// COST ESTIMATION
// =============================================================================

// CostEstimate represents estimated costs for a pipeline run.
type CostEstimate struct {
	ASRCost       float64 `json:"asrCost"`
	TranslateCost float64 `json:"translateCost"`
	RewriteCost   float64 `json:"rewriteCost"`
	TotalCost     float64 `json:"totalCost"`
	Currency      string  `json:"currency"` // always "USD"
}

// EstimateCost estimates the cost for processing a project.
func (a *App) EstimateCost(projectID string, config PipelineConfig) (*CostEstimate, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: EstimateCost")
}

// =============================================================================
// STATS
// =============================================================================

// AppStats represents application-wide statistics.
type AppStats struct {
	TotalProjects       int     `json:"totalProjects"`
	TotalSegments       int     `json:"totalSegments"`
	TotalCharsProcessed int64   `json:"totalCharsProcessed"`
	TotalMinutesASR     float64 `json:"totalMinutesAsr"`
}

// GetStats returns application statistics.
func (a *App) GetStats() (*AppStats, error) {
	return nil, fmt.Errorf("ERR_NOT_IMPLEMENTED: GetStats")
}
