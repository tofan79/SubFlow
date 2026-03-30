// Package store provides SQLite database operations for SubFlow.
// It handles all persistence including projects, segments, settings, and glossary.
package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Store provides database operations for SubFlow.
type Store struct {
	db *sql.DB
	mu sync.RWMutex
}

// New creates a new Store instance and initializes the database.
func New() (*Store, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, fmt.Errorf("store.New: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("store.New: create dir: %w", err)
	}

	// Open database with WAL mode and foreign keys enabled
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_foreign_keys=ON&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("store.New: open db: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("store.New: ping: %w", err)
	}

	s := &Store{db: db}

	// Run migrations
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, fmt.Errorf("store.New: migrate: %w", err)
	}

	return s, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// getDBPath returns the database file path based on OS.
func getDBPath() (string, error) {
	var baseDir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA not set")
		}
		baseDir = filepath.Join(appData, "SubFlow")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home dir: %w", err)
		}
		baseDir = filepath.Join(home, "Library", "Application Support", "SubFlow")
	default: // linux and others
		// Use XDG_DATA_HOME or fallback to ~/.local/share
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("get home dir: %w", err)
			}
			dataHome = filepath.Join(home, ".local", "share")
		}
		baseDir = filepath.Join(dataHome, "subflow")
	}

	return filepath.Join(baseDir, "subflow.db"), nil
}

// migrate runs all database migrations.
func (s *Store) migrate() error {
	migrations := []string{
		migrationV1,
	}

	// Create migrations table if not exists
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get current version
	var currentVersion int
	row := s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations")
	if err := row.Scan(&currentVersion); err != nil {
		return fmt.Errorf("get current version: %w", err)
	}

	// Apply pending migrations
	for i, migration := range migrations {
		version := i + 1
		if version <= currentVersion {
			continue
		}

		tx, err := s.db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx for v%d: %w", version, err)
		}

		if _, err := tx.Exec(migration); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration v%d: %w", version, err)
		}

		if _, err := tx.Exec("INSERT INTO migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration v%d: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration v%d: %w", version, err)
		}
	}

	return nil
}

// Migration V1: Initial schema
const migrationV1 = `
-- Projects table
CREATE TABLE IF NOT EXISTS projects (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	source_path TEXT NOT NULL,
	state TEXT NOT NULL DEFAULT 'imported',
	created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- Segments table
CREATE TABLE IF NOT EXISTS segments (
	id TEXT PRIMARY KEY,
	project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	idx INTEGER NOT NULL,
	start_ms INTEGER NOT NULL,
	end_ms INTEGER NOT NULL,
	source TEXT NOT NULL DEFAULT '',
	l1 TEXT NOT NULL DEFAULT '',
	l2 TEXT NOT NULL DEFAULT '',
	speaker TEXT NOT NULL DEFAULT '',
	emotion TEXT NOT NULL DEFAULT '',
	qa_status TEXT NOT NULL DEFAULT 'pending',
	created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
CREATE INDEX IF NOT EXISTS idx_segments_project ON segments(project_id);
CREATE INDEX IF NOT EXISTS idx_segments_project_idx ON segments(project_id, idx);

-- Settings table (key-value store)
CREATE TABLE IF NOT EXISTS settings (
	key TEXT PRIMARY KEY,
	value TEXT NOT NULL,
	updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);

-- Glossary table
CREATE TABLE IF NOT EXISTS glossary (
	id TEXT PRIMARY KEY,
	source_term TEXT NOT NULL,
	target_term TEXT NOT NULL,
	case_sensitive INTEGER NOT NULL DEFAULT 0,
	notes TEXT NOT NULL DEFAULT '',
	created_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
CREATE INDEX IF NOT EXISTS idx_glossary_source ON glossary(source_term);

-- Stats table (for tracking usage)
CREATE TABLE IF NOT EXISTS stats (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	total_chars_processed INTEGER NOT NULL DEFAULT 0,
	total_minutes_asr REAL NOT NULL DEFAULT 0.0,
	updated_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
);
INSERT OR IGNORE INTO stats (id) VALUES (1);
`

// =============================================================================
// PROJECT CRUD
// =============================================================================

// Project represents a project record.
type Project struct {
	ID         string
	Name       string
	SourcePath string
	State      string
	CreatedAt  int64
	UpdatedAt  int64
}

// CreateProject creates a new project with a random ID.
func (s *Store) CreateProject(name, sourcePath string) (*Project, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate random ID using SQLite
	var id string
	err := s.db.QueryRow("SELECT lower(hex(randomblob(8)))").Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("store.CreateProject: generate id: %w", err)
	}

	_, err = s.db.Exec(`
		INSERT INTO projects (id, name, source_path, state)
		VALUES (?, ?, ?, 'imported')
	`, id, name, sourcePath)
	if err != nil {
		return nil, fmt.Errorf("store.CreateProject: insert: %w", err)
	}

	return s.getProjectLocked(id)
}

// GetProject returns a project by ID.
func (s *Store) GetProject(id string) (*Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.getProjectLocked(id)
}

func (s *Store) getProjectLocked(id string) (*Project, error) {
	var p Project
	err := s.db.QueryRow(`
		SELECT id, name, source_path, state, created_at, updated_at
		FROM projects WHERE id = ?
	`, id).Scan(&p.ID, &p.Name, &p.SourcePath, &p.State, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("store.GetProject: not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("store.GetProject: %w", err)
	}
	return &p, nil
}

// GetProjects returns all projects ordered by updated_at desc.
func (s *Store) GetProjects() ([]Project, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`
		SELECT id, name, source_path, state, created_at, updated_at
		FROM projects ORDER BY updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("store.GetProjects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.SourcePath, &p.State, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store.GetProjects: scan: %w", err)
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

// UpdateProjectState updates the state of a project.
func (s *Store) UpdateProjectState(id, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec(`
		UPDATE projects SET state = ?, updated_at = strftime('%s', 'now')
		WHERE id = ?
	`, state, id)
	if err != nil {
		return fmt.Errorf("store.UpdateProjectState: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.UpdateProjectState: not found: %s", id)
	}
	return nil
}

// DeleteProject deletes a project and all its segments (cascade).
func (s *Store) DeleteProject(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("store.DeleteProject: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.DeleteProject: not found: %s", id)
	}
	return nil
}

// =============================================================================
// SEGMENT CRUD
// =============================================================================

// Segment represents a subtitle segment record.
type Segment struct {
	ID        string
	ProjectID string
	Index     int
	StartMS   int64
	EndMS     int64
	Source    string
	L1        string
	L2        string
	Speaker   string
	Emotion   string
	QAStatus  string
	CreatedAt int64
	UpdatedAt int64
}

// CreateSegment creates a new segment.
func (s *Store) CreateSegment(seg *Segment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate random ID if not provided
	if seg.ID == "" {
		err := s.db.QueryRow("SELECT lower(hex(randomblob(8)))").Scan(&seg.ID)
		if err != nil {
			return fmt.Errorf("store.CreateSegment: generate id: %w", err)
		}
	}

	_, err := s.db.Exec(`
		INSERT INTO segments (id, project_id, idx, start_ms, end_ms, source, l1, l2, speaker, emotion, qa_status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, seg.ID, seg.ProjectID, seg.Index, seg.StartMS, seg.EndMS, seg.Source, seg.L1, seg.L2, seg.Speaker, seg.Emotion, seg.QAStatus)
	if err != nil {
		return fmt.Errorf("store.CreateSegment: %w", err)
	}
	return nil
}

// GetSegments returns all segments for a project ordered by index.
func (s *Store) GetSegments(projectID string) ([]Segment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`
		SELECT id, project_id, idx, start_ms, end_ms, source, l1, l2, speaker, emotion, qa_status, created_at, updated_at
		FROM segments WHERE project_id = ? ORDER BY idx
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("store.GetSegments: %w", err)
	}
	defer rows.Close()

	var segments []Segment
	for rows.Next() {
		var seg Segment
		if err := rows.Scan(&seg.ID, &seg.ProjectID, &seg.Index, &seg.StartMS, &seg.EndMS,
			&seg.Source, &seg.L1, &seg.L2, &seg.Speaker, &seg.Emotion, &seg.QAStatus,
			&seg.CreatedAt, &seg.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store.GetSegments: scan: %w", err)
		}
		segments = append(segments, seg)
	}
	return segments, rows.Err()
}

// UpdateSegment updates a segment.
func (s *Store) UpdateSegment(seg *Segment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec(`
		UPDATE segments SET
			idx = ?, start_ms = ?, end_ms = ?, source = ?, l1 = ?, l2 = ?,
			speaker = ?, emotion = ?, qa_status = ?, updated_at = strftime('%s', 'now')
		WHERE id = ?
	`, seg.Index, seg.StartMS, seg.EndMS, seg.Source, seg.L1, seg.L2,
		seg.Speaker, seg.Emotion, seg.QAStatus, seg.ID)
	if err != nil {
		return fmt.Errorf("store.UpdateSegment: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.UpdateSegment: not found: %s", seg.ID)
	}
	return nil
}

// DeleteSegment deletes a segment.
func (s *Store) DeleteSegment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM segments WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("store.DeleteSegment: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.DeleteSegment: not found: %s", id)
	}
	return nil
}

// CountSegments returns the number of segments for a project.
func (s *Store) CountSegments(projectID string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM segments WHERE project_id = ?", projectID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("store.CountSegments: %w", err)
	}
	return count, nil
}

// =============================================================================
// SETTINGS CRUD
// =============================================================================

// GetSetting returns a setting value by key.
func (s *Store) GetSetting(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var value string
	err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil // Return empty string for non-existent keys
	}
	if err != nil {
		return "", fmt.Errorf("store.GetSetting: %w", err)
	}
	return value, nil
}

// SetSetting sets a setting value.
func (s *Store) SetSetting(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		INSERT INTO settings (key, value) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = strftime('%s', 'now')
	`, key, value)
	if err != nil {
		return fmt.Errorf("store.SetSetting: %w", err)
	}
	return nil
}

// GetAllSettings returns all settings as a map.
func (s *Store) GetAllSettings() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, fmt.Errorf("store.GetAllSettings: %w", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, fmt.Errorf("store.GetAllSettings: scan: %w", err)
		}
		settings[key] = value
	}
	return settings, rows.Err()
}

// =============================================================================
// GLOSSARY CRUD
// =============================================================================

// GlossaryTerm represents a glossary entry.
type GlossaryTerm struct {
	ID            string
	SourceTerm    string
	TargetTerm    string
	CaseSensitive bool
	Notes         string
	CreatedAt     int64
	UpdatedAt     int64
}

// CreateGlossaryTerm creates a new glossary term.
func (s *Store) CreateGlossaryTerm(term *GlossaryTerm) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate random ID if not provided
	if term.ID == "" {
		err := s.db.QueryRow("SELECT lower(hex(randomblob(8)))").Scan(&term.ID)
		if err != nil {
			return fmt.Errorf("store.CreateGlossaryTerm: generate id: %w", err)
		}
	}

	caseSensitive := 0
	if term.CaseSensitive {
		caseSensitive = 1
	}

	_, err := s.db.Exec(`
		INSERT INTO glossary (id, source_term, target_term, case_sensitive, notes)
		VALUES (?, ?, ?, ?, ?)
	`, term.ID, term.SourceTerm, term.TargetTerm, caseSensitive, term.Notes)
	if err != nil {
		return fmt.Errorf("store.CreateGlossaryTerm: %w", err)
	}
	return nil
}

// GetGlossary returns all glossary terms.
func (s *Store) GetGlossary() ([]GlossaryTerm, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.Query(`
		SELECT id, source_term, target_term, case_sensitive, notes, created_at, updated_at
		FROM glossary ORDER BY source_term
	`)
	if err != nil {
		return nil, fmt.Errorf("store.GetGlossary: %w", err)
	}
	defer rows.Close()

	var terms []GlossaryTerm
	for rows.Next() {
		var t GlossaryTerm
		var caseSensitive int
		if err := rows.Scan(&t.ID, &t.SourceTerm, &t.TargetTerm, &caseSensitive, &t.Notes, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("store.GetGlossary: scan: %w", err)
		}
		t.CaseSensitive = caseSensitive == 1
		terms = append(terms, t)
	}
	return terms, rows.Err()
}

// UpdateGlossaryTerm updates a glossary term.
func (s *Store) UpdateGlossaryTerm(term *GlossaryTerm) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	caseSensitive := 0
	if term.CaseSensitive {
		caseSensitive = 1
	}

	result, err := s.db.Exec(`
		UPDATE glossary SET
			source_term = ?, target_term = ?, case_sensitive = ?, notes = ?,
			updated_at = strftime('%s', 'now')
		WHERE id = ?
	`, term.SourceTerm, term.TargetTerm, caseSensitive, term.Notes, term.ID)
	if err != nil {
		return fmt.Errorf("store.UpdateGlossaryTerm: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.UpdateGlossaryTerm: not found: %s", term.ID)
	}
	return nil
}

// DeleteGlossaryTerm deletes a glossary term.
func (s *Store) DeleteGlossaryTerm(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM glossary WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("store.DeleteGlossaryTerm: %w", err)
	}

	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("store.DeleteGlossaryTerm: not found: %s", id)
	}
	return nil
}

// =============================================================================
// STATS
// =============================================================================

// Stats represents usage statistics.
type Stats struct {
	TotalProjects       int
	TotalSegments       int
	TotalCharsProcessed int64
	TotalMinutesASR     float64
}

// GetStats returns application statistics.
func (s *Store) GetStats() (*Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats Stats

	// Count projects
	err := s.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&stats.TotalProjects)
	if err != nil {
		return nil, fmt.Errorf("store.GetStats: count projects: %w", err)
	}

	// Count segments
	err = s.db.QueryRow("SELECT COUNT(*) FROM segments").Scan(&stats.TotalSegments)
	if err != nil {
		return nil, fmt.Errorf("store.GetStats: count segments: %w", err)
	}

	// Get cumulative stats
	err = s.db.QueryRow("SELECT total_chars_processed, total_minutes_asr FROM stats WHERE id = 1").
		Scan(&stats.TotalCharsProcessed, &stats.TotalMinutesASR)
	if err != nil {
		return nil, fmt.Errorf("store.GetStats: get cumulative: %w", err)
	}

	return &stats, nil
}

// IncrementStats increments the cumulative statistics.
func (s *Store) IncrementStats(chars int64, minutes float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(`
		UPDATE stats SET
			total_chars_processed = total_chars_processed + ?,
			total_minutes_asr = total_minutes_asr + ?,
			updated_at = strftime('%s', 'now')
		WHERE id = 1
	`, chars, minutes)
	if err != nil {
		return fmt.Errorf("store.IncrementStats: %w", err)
	}
	return nil
}
