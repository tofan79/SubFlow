// Package pipeline provides error definitions for the subtitle processing pipeline.
package pipeline

import (
	"fmt"
)

// =============================================================================
// ERROR TYPE
// =============================================================================

// Error represents a pipeline error with code and user message.
type Error struct {
	Code    string // e.g., "ERR_IMP_001"
	Message string // Technical message for logs
	Cause   error  // Underlying error if any
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for errors.Is/errors.As.
func (e *Error) Unwrap() error {
	return e.Cause
}

// UserMessage returns a 3-line user-friendly message in Indonesian.
// Format:
//
//	Line 1: Error code
//	Line 2: What happened
//	Line 3: What user can do
func (e *Error) UserMessage() string {
	msg, ok := userMessages[e.Code]
	if !ok {
		return fmt.Sprintf("[%s]\nTerjadi kesalahan tidak dikenal.\nHubungi developer dengan kode error ini.", e.Code)
	}
	return fmt.Sprintf("[%s]\n%s\n%s", e.Code, msg.what, msg.action)
}

// userMessage holds the user-friendly message components.
type userMessage struct {
	what   string // What happened
	action string // What user can do
}

// =============================================================================
// ERROR CODES
// =============================================================================

// Import errors (ERR_IMP_*)
const (
	ErrImpFileNotFound    = "ERR_IMP_001"
	ErrImpUnsupportedType = "ERR_IMP_002"
	ErrImpFileTooLarge    = "ERR_IMP_003"
	ErrImpCorruptFile     = "ERR_IMP_004"
)

// ASR errors (ERR_ASR_*)
const (
	ErrASRNotInstalled = "ERR_ASR_001"
	ErrASRModelMissing = "ERR_ASR_002"
	ErrASRExtractFail  = "ERR_ASR_003"
	ErrASRTimeout      = "ERR_ASR_004"
	ErrASRNoGPU        = "ERR_ASR_005"
)

// Translation errors (ERR_TRN_*)
const (
	ErrTrnAPIKey      = "ERR_TRN_001"
	ErrTrnRateLimit   = "ERR_TRN_002"
	ErrTrnQuotaExceed = "ERR_TRN_003"
	ErrTrnTimeout     = "ERR_TRN_004"
)

// Rewrite errors (ERR_RWT_*)
const (
	ErrRwtAPIKey    = "ERR_RWT_001"
	ErrRwtRateLimit = "ERR_RWT_002"
	ErrRwtTimeout   = "ERR_RWT_003"
)

// QA errors (ERR_QA_*)
const (
	ErrQAValidation = "ERR_QA_001"
	ErrQAAutofix    = "ERR_QA_002"
)

// Export errors (ERR_EXP_*)
const (
	ErrExpWriteFail   = "ERR_EXP_001"
	ErrExpInvalidPath = "ERR_EXP_002"
)

// Database errors (ERR_DB_*)
const (
	ErrDBConnection = "ERR_DB_001"
	ErrDBMigration  = "ERR_DB_002"
	ErrDBQuery      = "ERR_DB_003"
)

// System errors (ERR_SYS_*)
const (
	ErrSysNotImplemented = "ERR_SYS_001"
)

// =============================================================================
// USER MESSAGES (Indonesian)
// =============================================================================

var userMessages = map[string]userMessage{
	// Import errors
	ErrImpFileNotFound: {
		what:   "File tidak ditemukan.",
		action: "Pastikan file masih ada di lokasi yang dipilih.",
	},
	ErrImpUnsupportedType: {
		what:   "Tipe file tidak didukung.",
		action: "Gunakan file video (MP4/MKV/AVI/MOV) atau subtitle (SRT/VTT/ASS).",
	},
	ErrImpFileTooLarge: {
		what:   "Ukuran file terlalu besar.",
		action: "Maksimum ukuran file adalah 4GB. Coba kompres video terlebih dahulu.",
	},
	ErrImpCorruptFile: {
		what:   "File rusak atau tidak dapat dibaca.",
		action: "Coba buka file di aplikasi lain untuk memastikan file tidak corrupt.",
	},

	// ASR errors
	ErrASRNotInstalled: {
		what:   "Komponen ASR (faster-whisper) tidak ditemukan.",
		action: "Reinstall SubFlow untuk memperbaiki instalasi.",
	},
	ErrASRModelMissing: {
		what:   "Model Whisper tidak ditemukan.",
		action: "Buka Settings dan download model yang diperlukan.",
	},
	ErrASRExtractFail: {
		what:   "Gagal mengekstrak audio dari video.",
		action: "Pastikan video memiliki track audio. Coba convert video ke MP4 terlebih dahulu.",
	},
	ErrASRTimeout: {
		what:   "Proses transkripsi memakan waktu terlalu lama.",
		action: "Coba gunakan model yang lebih kecil atau aktifkan GPU acceleration.",
	},
	ErrASRNoGPU: {
		what:   "GPU yang dipilih tidak tersedia.",
		action: "Pilih backend CPU di Settings atau install driver GPU yang sesuai.",
	},

	// Translation errors
	ErrTrnAPIKey: {
		what:   "API key untuk layanan terjemahan tidak valid.",
		action: "Periksa API key di Settings. Pastikan key masih aktif.",
	},
	ErrTrnRateLimit: {
		what:   "Terlalu banyak request ke layanan terjemahan.",
		action: "Tunggu beberapa menit lalu coba lagi, atau kurangi jumlah segment.",
	},
	ErrTrnQuotaExceed: {
		what:   "Kuota API terjemahan habis.",
		action: "Top up akun API Anda atau gunakan provider lain.",
	},
	ErrTrnTimeout: {
		what:   "Layanan terjemahan tidak merespons.",
		action: "Periksa koneksi internet. Coba lagi dalam beberapa saat.",
	},

	// Rewrite errors
	ErrRwtAPIKey: {
		what:   "API key untuk layanan rewrite tidak valid.",
		action: "Periksa API key di Settings. Pastikan key masih aktif.",
	},
	ErrRwtRateLimit: {
		what:   "Terlalu banyak request ke layanan rewrite.",
		action: "Tunggu beberapa menit lalu coba lagi.",
	},
	ErrRwtTimeout: {
		what:   "Layanan rewrite tidak merespons.",
		action: "Periksa koneksi internet. Coba lagi dalam beberapa saat.",
	},

	// QA errors
	ErrQAValidation: {
		what:   "Subtitle tidak lolos pengecekan kualitas.",
		action: "Periksa tab QA untuk melihat masalah yang ditemukan.",
	},
	ErrQAAutofix: {
		what:   "Auto-fix gagal memperbaiki semua masalah.",
		action: "Edit manual segment yang masih bermasalah di Editor.",
	},

	// Export errors
	ErrExpWriteFail: {
		what:   "Gagal menyimpan file ekspor.",
		action: "Pastikan folder tujuan tidak read-only dan ada ruang disk cukup.",
	},
	ErrExpInvalidPath: {
		what:   "Lokasi ekspor tidak valid.",
		action: "Pilih folder lain yang dapat ditulis.",
	},

	// Database errors
	ErrDBConnection: {
		what:   "Gagal terhubung ke database lokal.",
		action: "Restart SubFlow. Jika masalah berlanjut, hapus folder data dan reinstall.",
	},
	ErrDBMigration: {
		what:   "Gagal memperbarui struktur database.",
		action: "Backup folder data, lalu reinstall SubFlow.",
	},
	ErrDBQuery: {
		what:   "Gagal membaca/menulis data.",
		action: "Restart SubFlow. Jika berlanjut, mungkin database corrupt.",
	},

	// System errors
	ErrSysNotImplemented: {
		what:   "Fitur ini belum tersedia.",
		action: "Fitur sedang dalam pengembangan. Nantikan update berikutnya.",
	},
}

// =============================================================================
// ERROR CONSTRUCTORS
// =============================================================================

// NewError creates a new pipeline error.
func NewError(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Import error constructors
func ErrFileNotFound(path string, cause error) *Error {
	return NewError(ErrImpFileNotFound, fmt.Sprintf("file not found: %s", path), cause)
}

func ErrUnsupportedType(ext string) *Error {
	return NewError(ErrImpUnsupportedType, fmt.Sprintf("unsupported file type: %s", ext), nil)
}

func ErrFileTooLarge(size int64, maxSize int64) *Error {
	return NewError(ErrImpFileTooLarge, fmt.Sprintf("file size %d exceeds max %d", size, maxSize), nil)
}

func ErrCorruptFile(path string, cause error) *Error {
	return NewError(ErrImpCorruptFile, fmt.Sprintf("corrupt file: %s", path), cause)
}

// ASR error constructors
func ErrASRNotInstalledErr(path string) *Error {
	return NewError(ErrASRNotInstalled, fmt.Sprintf("whisper binary not found at: %s", path), nil)
}

func ErrASRModelMissingErr(model string) *Error {
	return NewError(ErrASRModelMissing, fmt.Sprintf("model not found: %s", model), nil)
}

func ErrASRExtractFailErr(cause error) *Error {
	return NewError(ErrASRExtractFail, "audio extraction failed", cause)
}

func ErrASRTimeoutErr(seconds int) *Error {
	return NewError(ErrASRTimeout, fmt.Sprintf("ASR timed out after %d seconds", seconds), nil)
}

func ErrASRNoGPUErr(backend string) *Error {
	return NewError(ErrASRNoGPU, fmt.Sprintf("GPU backend not available: %s", backend), nil)
}

// Translation error constructors
func ErrTrnAPIKeyErr(provider string) *Error {
	return NewError(ErrTrnAPIKey, fmt.Sprintf("invalid API key for %s", provider), nil)
}

func ErrTrnRateLimitErr(provider string, retryAfter int) *Error {
	return NewError(ErrTrnRateLimit, fmt.Sprintf("%s rate limited, retry after %ds", provider, retryAfter), nil)
}

func ErrTrnQuotaExceedErr(provider string) *Error {
	return NewError(ErrTrnQuotaExceed, fmt.Sprintf("%s quota exceeded", provider), nil)
}

func ErrTrnTimeoutErr(provider string, cause error) *Error {
	return NewError(ErrTrnTimeout, fmt.Sprintf("%s request timeout", provider), cause)
}

// Rewrite error constructors
func ErrRwtAPIKeyErr(provider string) *Error {
	return NewError(ErrRwtAPIKey, fmt.Sprintf("invalid API key for %s", provider), nil)
}

func ErrRwtRateLimitErr(provider string) *Error {
	return NewError(ErrRwtRateLimit, fmt.Sprintf("%s rate limited", provider), nil)
}

func ErrRwtTimeoutErr(provider string, cause error) *Error {
	return NewError(ErrRwtTimeout, fmt.Sprintf("%s request timeout", provider), cause)
}

// QA error constructors
func ErrQAValidationErr(failedChecks int) *Error {
	return NewError(ErrQAValidation, fmt.Sprintf("%d QA checks failed", failedChecks), nil)
}

func ErrQAAutofixErr(remaining int, cause error) *Error {
	return NewError(ErrQAAutofix, fmt.Sprintf("autofix incomplete, %d issues remain", remaining), cause)
}

// Export error constructors
func ErrExpWriteFailErr(path string, cause error) *Error {
	return NewError(ErrExpWriteFail, fmt.Sprintf("failed to write: %s", path), cause)
}

func ErrExpInvalidPathErr(path string) *Error {
	return NewError(ErrExpInvalidPath, fmt.Sprintf("invalid export path: %s", path), nil)
}

// Database error constructors
func ErrDBConnectionErr(cause error) *Error {
	return NewError(ErrDBConnection, "database connection failed", cause)
}

func ErrDBMigrationErr(version int, cause error) *Error {
	return NewError(ErrDBMigration, fmt.Sprintf("migration failed at version %d", version), cause)
}

func ErrDBQueryErr(operation string, cause error) *Error {
	return NewError(ErrDBQuery, fmt.Sprintf("database %s failed", operation), cause)
}

// System error constructors
func ErrNotImplemented(feature string) *Error {
	return NewError(ErrSysNotImplemented, fmt.Sprintf("not implemented: %s", feature), nil)
}

// =============================================================================
// ERROR HELPERS
// =============================================================================

// IsRetryable returns true if the error is retryable.
func IsRetryable(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	switch e.Code {
	case ErrTrnRateLimit, ErrTrnTimeout, ErrRwtRateLimit, ErrRwtTimeout, ErrASRTimeout:
		return true
	default:
		return false
	}
}

// GetCode extracts the error code from an error.
// Returns empty string if not a pipeline error.
func GetCode(err error) string {
	e, ok := err.(*Error)
	if !ok {
		return ""
	}
	return e.Code
}
