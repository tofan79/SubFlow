// Package qa provides pure-logic QA validation for subtitle cards.
// This package has NO I/O operations and NO imports from other internal packages.
package qa

// =============================================================================
// QA CHECK CONSTANTS
// =============================================================================

// CheckID represents a QA check identifier.
type CheckID string

const (
	CheckLineLength  CheckID = "QA-01" // Max 42 chars per line
	CheckLineCount   CheckID = "QA-02" // Max 2 lines per card
	CheckDurationMin CheckID = "QA-03" // Min 1.0 second
	CheckDurationMax CheckID = "QA-04" // Max 7.0 seconds
	CheckCPS         CheckID = "QA-05" // Max 17 chars per second
	CheckOverlap     CheckID = "QA-06" // No overlapping cards
	CheckEmptyCard   CheckID = "QA-07" // No empty cards
	CheckGap         CheckID = "QA-08" // Min 83ms gap between cards
	CheckGlossary    CheckID = "QA-09" // Glossary consistency (warning only)
)

// AllChecks is an ordered list of all QA check IDs.
var AllChecks = []CheckID{
	CheckLineLength,
	CheckLineCount,
	CheckDurationMin,
	CheckDurationMax,
	CheckCPS,
	CheckOverlap,
	CheckEmptyCard,
	CheckGap,
	CheckGlossary,
}

// CheckNames maps check IDs to human-readable names.
var CheckNames = map[CheckID]string{
	CheckLineLength:  "line_length",
	CheckLineCount:   "line_count",
	CheckDurationMin: "duration_min",
	CheckDurationMax: "duration_max",
	CheckCPS:         "cps",
	CheckOverlap:     "overlap",
	CheckEmptyCard:   "empty_card",
	CheckGap:         "gap",
	CheckGlossary:    "glossary",
}

// CheckDescriptions maps check IDs to Indonesian descriptions.
var CheckDescriptions = map[CheckID]string{
	CheckLineLength:  "Maksimum 42 karakter per baris",
	CheckLineCount:   "Maksimum 2 baris per card",
	CheckDurationMin: "Minimum durasi 1.0 detik",
	CheckDurationMax: "Maksimum durasi 7.0 detik",
	CheckCPS:         "Maksimum 17 karakter per detik",
	CheckOverlap:     "Tidak ada overlap antar card",
	CheckEmptyCard:   "Tidak ada card kosong",
	CheckGap:         "Minimum 83ms jarak antar card",
	CheckGlossary:    "Konsistensi istilah glossary",
}

// IsAutoFixable returns true if the check can be auto-fixed.
// QA-09 (glossary) is NOT auto-fixable.
func IsAutoFixable(id CheckID) bool {
	return id != CheckGlossary
}

// =============================================================================
// SEVERITY
// =============================================================================

// Severity represents the severity level of a QA result.
type Severity string

const (
	SeverityPass    Severity = "pass"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// =============================================================================
// SUBTITLE CARD (local copy to avoid importing pipeline package)
// =============================================================================

// SubtitleCard represents a single subtitle entry for QA validation.
// This is a local copy to keep the qa package dependency-free.
type SubtitleCard struct {
	ID       string // Unique identifier
	Index    int    // Display order (1-based)
	StartMS  int64  // Start time in milliseconds
	EndMS    int64  // End time in milliseconds
	Text     string // The subtitle text (L2 or final text)
	Speaker  string // Speaker identifier (optional)
	Emotion  string // Detected emotion (optional)
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
// GLOSSARY TERM (local copy)
// =============================================================================

// GlossaryTerm represents a term that should be translated consistently.
type GlossaryTerm struct {
	SourceTerm    string
	TargetTerm    string
	CaseSensitive bool
}

// =============================================================================
// QA RESULT
// =============================================================================

// Result represents the outcome of a single QA check on a single card.
type Result struct {
	CardID    string   `json:"cardId"`
	CardIndex int      `json:"cardIndex"`
	CheckID   CheckID  `json:"checkId"`
	Passed    bool     `json:"passed"`
	Severity  Severity `json:"severity"`
	Detail    string   `json:"detail"`    // Human-readable explanation
	AutoFixed bool     `json:"autoFixed"` // Was this issue auto-fixed?
	FixAction string   `json:"fixAction"` // Description of fix applied
}

// =============================================================================
// QA LOG
// =============================================================================

// Log represents a single auto-fix action taken.
type Log struct {
	CardID     string  `json:"cardId"`
	CardIndex  int     `json:"cardIndex"`
	CheckID    CheckID `json:"checkId"`
	Action     string  `json:"action"`     // What was done
	OldValue   string  `json:"oldValue"`   // Before fix
	NewValue   string  `json:"newValue"`   // After fix
	LoopNumber int     `json:"loopNumber"` // Which auto-fix iteration (1-3)
}

// =============================================================================
// QA REPORT
// =============================================================================

// Report represents the complete QA report for a project.
type Report struct {
	RunAt      int64    `json:"run_at"`      // Unix timestamp
	TotalCards int      `json:"total_cards"` // Total cards checked
	Passed     int      `json:"passed"`      // Cards with all checks passed
	Warnings   int      `json:"warnings"`    // Cards with warnings (QA-09)
	Errors     int      `json:"errors"`      // Cards with errors
	AutoFixed  int      `json:"auto_fixed"`  // Number of auto-fixes applied
	Results    []Result `json:"results"`     // All individual results
}

// =============================================================================
// QA VALIDATOR CONFIG
// =============================================================================

// Config holds configuration for the QA validator.
type Config struct {
	MaxCharsPerLine int            // Default: 42
	MaxLines        int            // Default: 2
	MinDurationMS   int64          // Default: 1000 (1 second)
	MaxDurationMS   int64          // Default: 7000 (7 seconds)
	MaxCPS          float64        // Default: 17.0
	MinGapMS        int64          // Default: 83
	Glossary        []GlossaryTerm // Terms to check for consistency
}

// DefaultConfig returns the default QA configuration.
func DefaultConfig() Config {
	return Config{
		MaxCharsPerLine: 42,
		MaxLines:        2,
		MinDurationMS:   1000,
		MaxDurationMS:   7000,
		MaxCPS:          17.0,
		MinGapMS:        83,
		Glossary:        nil,
	}
}

// =============================================================================
// QA STATUS CONSTANTS
// =============================================================================

const (
	StatusPass    = "pass"
	StatusWarn    = "warn"
	StatusError   = "error"
	StatusPending = "pending"
)
