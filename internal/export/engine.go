// Package export provides subtitle export functionality with QA validation.
package export

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// =============================================================================
// EXPORT FORMAT
// =============================================================================

// Format represents the output subtitle format.
type Format string

const (
	FormatSRT  Format = "srt"
	FormatVTT  Format = "vtt"
	FormatASS  Format = "ass"
	FormatTXT  Format = "txt"
	FormatJSON Format = "json"
)

// ValidFormats is a list of all supported export formats.
var ValidFormats = []Format{FormatSRT, FormatVTT, FormatASS, FormatTXT, FormatJSON}

// IsValidFormat checks if the format is supported.
func IsValidFormat(f Format) bool {
	for _, vf := range ValidFormats {
		if f == vf {
			return true
		}
	}
	return false
}

// Extension returns the file extension for the format.
func (f Format) Extension() string {
	switch f {
	case FormatSRT:
		return ".srt"
	case FormatVTT:
		return ".vtt"
	case FormatASS:
		return ".ass"
	case FormatTXT:
		return ".txt"
	case FormatJSON:
		return ".json"
	default:
		return ".srt"
	}
}

// =============================================================================
// EXPORT LAYER
// =============================================================================

// Layer specifies which text layer to export.
type Layer string

const (
	LayerSource Layer = "source" // Original/ASR text
	LayerL1     Layer = "l1"     // Layer 1 translation
	LayerL2     Layer = "l2"     // Layer 2 rewrite (final)
	LayerDual   Layer = "dual"   // Dual subtitle (source + L2)
)

// =============================================================================
// SUBTITLE CARD (local copy to avoid circular imports)
// =============================================================================

// Card represents a single subtitle entry for export.
type Card struct {
	ID       string // Unique identifier
	Index    int    // Display order (1-based)
	StartMS  int64  // Start time in milliseconds
	EndMS    int64  // End time in milliseconds
	Source   string // Original/ASR text
	L1       string // Layer 1 translation
	L2       string // Layer 2 rewrite (final)
	Speaker  string // Speaker identifier (optional)
	Emotion  string // Detected emotion (optional)
	QAStatus string // "pass", "warn", "error", "pending"
}

// Duration returns the duration in milliseconds.
func (c *Card) Duration() int64 {
	return c.EndMS - c.StartMS
}

// Text returns the appropriate text for the given layer.
func (c *Card) Text(layer Layer) string {
	switch layer {
	case LayerSource:
		return c.Source
	case LayerL1:
		return c.L1
	case LayerL2:
		return c.L2
	default:
		return c.L2
	}
}

// =============================================================================
// QA STATUS (local copy)
// =============================================================================

const (
	QAStatusPass    = "pass"
	QAStatusWarn    = "warn"
	QAStatusError   = "error"
	QAStatusPending = "pending"
)

// =============================================================================
// EXPORT OPTIONS
// =============================================================================

// Options holds export configuration.
type Options struct {
	Format         Format // Output format (srt, vtt, ass, txt, json)
	Layer          Layer  // Which text layer to export
	OutputDir      string // Output directory path
	OutputFilename string // Output filename (without extension, auto-generated if empty)
	SkipQACheck    bool   // If true, skip QA validation before export (not recommended)
	AllowWarnings  bool   // If true, allow export with QA warnings (default: true)
	AllowErrors    bool   // If true, allow export with QA errors (not recommended)
	IncludeEmpty   bool   // If true, include empty cards (default: false)
}

// DefaultOptions returns default export options.
func DefaultOptions() Options {
	return Options{
		Format:        FormatSRT,
		Layer:         LayerL2,
		AllowWarnings: true,
		AllowErrors:   false,
		IncludeEmpty:  false,
	}
}

// =============================================================================
// EXPORT RESULT
// =============================================================================

// Result holds the outcome of an export operation.
type Result struct {
	Success      bool     `json:"success"`
	OutputPath   string   `json:"outputPath"`
	Format       Format   `json:"format"`
	CardsWritten int      `json:"cardsWritten"`
	QAWarnings   int      `json:"qaWarnings"`
	QAErrors     int      `json:"qaErrors"`
	SkippedEmpty int      `json:"skippedEmpty"`
	Error        string   `json:"error,omitempty"`
	ExportedAt   int64    `json:"exportedAt"` // Unix timestamp
	Warnings     []string `json:"warnings,omitempty"`
}

// =============================================================================
// EXPORTER
// =============================================================================

// Exporter handles subtitle export with QA validation.
type Exporter struct {
	projectID string
}

// NewExporter creates a new Exporter.
func NewExporter(projectID string) *Exporter {
	return &Exporter{
		projectID: projectID,
	}
}

// Export exports cards to the specified format.
func (e *Exporter) Export(cards []Card, opts Options) (*Result, error) {
	result := &Result{
		Format:     opts.Format,
		ExportedAt: time.Now().Unix(),
	}

	// Validate format
	if !IsValidFormat(opts.Format) {
		result.Error = fmt.Sprintf("format tidak didukung: %s", opts.Format)
		return result, fmt.Errorf("export.Export: %s", result.Error)
	}

	// Validate output directory
	if opts.OutputDir == "" {
		result.Error = "output directory tidak boleh kosong"
		return result, fmt.Errorf("export.Export: %s", result.Error)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		result.Error = fmt.Sprintf("gagal membuat output directory: %v", err)
		return result, fmt.Errorf("export.Export: %w", err)
	}

	// QA Check
	if !opts.SkipQACheck {
		qaWarnings, qaErrors := e.countQAIssues(cards)
		result.QAWarnings = qaWarnings
		result.QAErrors = qaErrors

		if qaErrors > 0 && !opts.AllowErrors {
			result.Error = fmt.Sprintf("ekspor dibatalkan: %d kartu memiliki QA error", qaErrors)
			result.Warnings = append(result.Warnings, result.Error)
			return result, fmt.Errorf("export.Export: %s", result.Error)
		}

		if qaWarnings > 0 && !opts.AllowWarnings {
			result.Error = fmt.Sprintf("ekspor dibatalkan: %d kartu memiliki QA warning", qaWarnings)
			result.Warnings = append(result.Warnings, result.Error)
			return result, fmt.Errorf("export.Export: %s", result.Error)
		}

		if qaErrors > 0 {
			result.Warnings = append(result.Warnings, fmt.Sprintf("%d kartu memiliki QA error", qaErrors))
		}
	}

	// Filter empty cards if needed
	exportCards := cards
	if !opts.IncludeEmpty {
		exportCards, result.SkippedEmpty = e.filterEmpty(cards, opts.Layer)
	}

	// Generate output path
	outputPath := e.generateOutputPath(opts)
	result.OutputPath = outputPath

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		result.Error = fmt.Sprintf("gagal membuat file: %v", err)
		return result, fmt.Errorf("export.Export: %w", err)
	}
	defer file.Close()

	// Write content based on format
	var writeErr error
	switch opts.Format {
	case FormatSRT:
		if opts.Layer == LayerDual {
			writeErr = e.writeDualSRT(file, exportCards)
		} else {
			writeErr = e.writeSRT(file, exportCards, opts.Layer)
		}
	case FormatVTT:
		if opts.Layer == LayerDual {
			writeErr = e.writeDualVTT(file, exportCards)
		} else {
			writeErr = e.writeVTT(file, exportCards, opts.Layer)
		}
	case FormatASS:
		if opts.Layer == LayerDual {
			writeErr = e.writeDualASS(file, exportCards)
		} else {
			writeErr = e.writeASS(file, exportCards, opts.Layer)
		}
	case FormatTXT:
		writeErr = e.writeTXT(file, exportCards, opts.Layer)
	case FormatJSON:
		writeErr = e.writeJSON(file, exportCards)
	}

	if writeErr != nil {
		result.Error = fmt.Sprintf("gagal menulis file: %v", writeErr)
		return result, fmt.Errorf("export.Export: %w", writeErr)
	}

	result.Success = true
	result.CardsWritten = len(exportCards)

	return result, nil
}

// countQAIssues counts cards with QA warnings and errors.
func (e *Exporter) countQAIssues(cards []Card) (warnings, errors int) {
	for _, card := range cards {
		switch card.QAStatus {
		case QAStatusWarn:
			warnings++
		case QAStatusError:
			errors++
		}
	}
	return
}

// filterEmpty filters out cards with empty text for the specified layer.
func (e *Exporter) filterEmpty(cards []Card, layer Layer) ([]Card, int) {
	var result []Card
	skipped := 0

	for _, card := range cards {
		text := card.Text(layer)
		// For dual layer, check both source and L2
		if layer == LayerDual {
			if strings.TrimSpace(card.Source) == "" && strings.TrimSpace(card.L2) == "" {
				skipped++
				continue
			}
		} else if strings.TrimSpace(text) == "" {
			skipped++
			continue
		}
		result = append(result, card)
	}

	return result, skipped
}

// generateOutputPath creates the output file path.
func (e *Exporter) generateOutputPath(opts Options) string {
	filename := opts.OutputFilename
	if filename == "" {
		filename = fmt.Sprintf("%s_%s_%d", e.projectID, opts.Layer, time.Now().Unix())
	}

	// Ensure no extension in filename
	filename = strings.TrimSuffix(filename, opts.Format.Extension())
	filename = strings.TrimSuffix(filename, ".")

	return filepath.Join(opts.OutputDir, filename+opts.Format.Extension())
}

// =============================================================================
// SRT WRITER
// =============================================================================

func (e *Exporter) writeSRT(w io.Writer, cards []Card, layer Layer) error {
	for i, card := range cards {
		text := card.Text(layer)
		if _, err := fmt.Fprintf(w, "%d\n%s --> %s\n%s\n\n",
			i+1,
			formatSRTTimestamp(card.StartMS),
			formatSRTTimestamp(card.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeDualSRT(w io.Writer, cards []Card) error {
	for i, card := range cards {
		text := card.Source
		if card.L2 != "" {
			text += "\n" + card.L2
		}
		if _, err := fmt.Fprintf(w, "%d\n%s --> %s\n%s\n\n",
			i+1,
			formatSRTTimestamp(card.StartMS),
			formatSRTTimestamp(card.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

// =============================================================================
// VTT WRITER
// =============================================================================

func (e *Exporter) writeVTT(w io.Writer, cards []Card, layer Layer) error {
	if _, err := io.WriteString(w, "WEBVTT\n\n"); err != nil {
		return err
	}

	for _, card := range cards {
		text := card.Text(layer)
		if _, err := fmt.Fprintf(w, "%s --> %s\n%s\n\n",
			formatVTTTimestamp(card.StartMS),
			formatVTTTimestamp(card.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeDualVTT(w io.Writer, cards []Card) error {
	if _, err := io.WriteString(w, "WEBVTT\n\n"); err != nil {
		return err
	}

	for _, card := range cards {
		text := card.Source
		if card.L2 != "" {
			text += "\n" + card.L2
		}
		if _, err := fmt.Fprintf(w, "%s --> %s\n%s\n\n",
			formatVTTTimestamp(card.StartMS),
			formatVTTTimestamp(card.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

// =============================================================================
// ASS WRITER
// =============================================================================

const assHeader = `[Script Info]
Title: SubFlow Export
ScriptType: v4.00+
PlayResX: 1920
PlayResY: 1080
WrapStyle: 0

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: Default,Arial,48,&H00FFFFFF,&H000000FF,&H00000000,&H96000000,0,0,0,0,100,100,0,0,1,2,1,2,20,20,30,1
Style: Source,Arial,40,&H00FFFF00,&H000000FF,&H00000000,&H96000000,0,0,0,0,100,100,0,0,1,2,1,8,20,20,80,1
Style: Translation,Arial,48,&H00FFFFFF,&H000000FF,&H00000000,&H96000000,0,0,0,0,100,100,0,0,1,2,1,2,20,20,30,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
`

func (e *Exporter) writeASS(w io.Writer, cards []Card, layer Layer) error {
	if _, err := io.WriteString(w, assHeader); err != nil {
		return err
	}

	for _, card := range cards {
		text := escapeASSText(card.Text(layer))
		if _, err := fmt.Fprintf(w, "Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n",
			formatASSTimestamp(card.StartMS),
			formatASSTimestamp(card.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeDualASS(w io.Writer, cards []Card) error {
	if _, err := io.WriteString(w, assHeader); err != nil {
		return err
	}

	for _, card := range cards {
		// Write source text (top position, using Source style)
		if card.Source != "" {
			source := escapeASSText(card.Source)
			if _, err := fmt.Fprintf(w, "Dialogue: 0,%s,%s,Source,,0,0,0,,%s\n",
				formatASSTimestamp(card.StartMS),
				formatASSTimestamp(card.EndMS),
				source,
			); err != nil {
				return err
			}
		}

		// Write translation (bottom position, using Translation style)
		if card.L2 != "" {
			translation := escapeASSText(card.L2)
			if _, err := fmt.Fprintf(w, "Dialogue: 1,%s,%s,Translation,,0,0,0,,%s\n",
				formatASSTimestamp(card.StartMS),
				formatASSTimestamp(card.EndMS),
				translation,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

// =============================================================================
// TXT WRITER
// =============================================================================

func (e *Exporter) writeTXT(w io.Writer, cards []Card, layer Layer) error {
	for _, card := range cards {
		text := card.Text(layer)

		// For dual, combine both
		if layer == LayerDual {
			text = card.Source
			if card.L2 != "" {
				text += " | " + card.L2
			}
		}

		if _, err := fmt.Fprintf(w, "[%s] %s\n",
			formatTXTTimestamp(card.StartMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

// =============================================================================
// JSON WRITER
// =============================================================================

func (e *Exporter) writeJSON(w io.Writer, cards []Card) error {
	if _, err := io.WriteString(w, "[\n"); err != nil {
		return err
	}

	for i, card := range cards {
		jsonCard := fmt.Sprintf(`  {
    "index": %d,
    "start_ms": %d,
    "end_ms": %d,
    "source": %q,
    "l1": %q,
    "l2": %q,
    "speaker": %q,
    "emotion": %q,
    "qa_status": %q
  }`,
			card.Index,
			card.StartMS,
			card.EndMS,
			card.Source,
			card.L1,
			card.L2,
			card.Speaker,
			card.Emotion,
			card.QAStatus,
		)

		if i < len(cards)-1 {
			jsonCard += ","
		}
		jsonCard += "\n"

		if _, err := io.WriteString(w, jsonCard); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "]\n"); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// TIMESTAMP FORMATTERS
// =============================================================================

func formatSRTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, ms)
}

func formatVTTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, ms)
}

func formatASSTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	totalCS := int((ms + 5) / 10)
	hours := totalCS / 360000
	totalCS %= 360000
	minutes := totalCS / 6000
	totalCS %= 6000
	seconds := totalCS / 100
	centiseconds := totalCS % 100
	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, centiseconds)
}

func formatTXTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func escapeASSText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\n", "\\N")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	return text
}

// =============================================================================
// BATCH EXPORT
// =============================================================================

// BatchExport exports cards to multiple formats at once.
func (e *Exporter) BatchExport(cards []Card, formats []Format, opts Options) ([]*Result, error) {
	var results []*Result

	for _, format := range formats {
		formatOpts := opts
		formatOpts.Format = format
		result, err := e.Export(cards, formatOpts)
		if err != nil {
			// Continue with other formats even if one fails
			results = append(results, result)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// =============================================================================
// VALIDATION HELPERS
// =============================================================================

// ValidateBeforeExport checks if cards are ready for export.
// Returns (canExport, warnings, errors).
func ValidateBeforeExport(cards []Card, opts Options) (bool, []string, []string) {
	var warnings, errors []string

	qaWarnings := 0
	qaErrors := 0

	for _, card := range cards {
		switch card.QAStatus {
		case QAStatusWarn:
			qaWarnings++
		case QAStatusError:
			qaErrors++
		case QAStatusPending:
			errors = append(errors, fmt.Sprintf("Kartu #%d belum di-QA", card.Index))
		}
	}

	if qaErrors > 0 {
		errors = append(errors, fmt.Sprintf("%d kartu memiliki QA error", qaErrors))
	}

	if qaWarnings > 0 {
		warnings = append(warnings, fmt.Sprintf("%d kartu memiliki QA warning", qaWarnings))
	}

	// Check for empty output directory
	if opts.OutputDir == "" {
		errors = append(errors, "Output directory tidak boleh kosong")
	}

	// Check for valid format
	if !IsValidFormat(opts.Format) {
		errors = append(errors, fmt.Sprintf("Format tidak valid: %s", opts.Format))
	}

	canExport := len(errors) == 0 || (opts.AllowErrors && qaErrors > 0)

	return canExport, warnings, errors
}
