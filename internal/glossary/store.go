package glossary

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Term struct {
	ID            string
	SourceTerm    string
	TargetTerm    string
	CaseSensitive bool
	Category      string
	Notes         string
}

type Store struct {
	terms map[string]Term
}

func NewStore() *Store {
	return &Store{terms: make(map[string]Term)}
}

func (s *Store) Add(term Term) error {
	if s.terms == nil {
		s.terms = make(map[string]Term)
	}

	if term.ID == "" {
		id, err := newID()
		if err != nil {
			return fmt.Errorf("glossary.Add: generate id: %w", err)
		}
		term.ID = id
	}

	s.terms[term.ID] = term
	return nil
}

func (s *Store) Get(id string) (Term, bool) {
	if s == nil || s.terms == nil {
		return Term{}, false
	}
	term, ok := s.terms[id]
	return term, ok
}

func (s *Store) Delete(id string) bool {
	if s == nil || s.terms == nil {
		return false
	}
	if _, ok := s.terms[id]; !ok {
		return false
	}
	delete(s.terms, id)
	return true
}

func (s *Store) List() []Term {
	if s == nil || len(s.terms) == 0 {
		return nil
	}

	ids := make([]string, 0, len(s.terms))
	for id := range s.terms {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	terms := make([]Term, 0, len(ids))
	for _, id := range ids {
		terms = append(terms, s.terms[id])
	}
	return terms
}

func (s *Store) FindBySource(source string) []Term {
	if s == nil || len(s.terms) == 0 {
		return nil
	}

	terms := make([]Term, 0)
	for _, term := range s.terms {
		if term.CaseSensitive {
			if source == term.SourceTerm {
				terms = append(terms, term)
			}
			continue
		}

		if strings.EqualFold(source, term.SourceTerm) {
			terms = append(terms, term)
		}
	}

	sort.Slice(terms, func(i, j int) bool { return terms[i].ID < terms[j].ID })
	return terms
}

func (s *Store) BuildPromptSection() string {
	terms := s.List()
	if len(terms) == 0 {
		return "Glossary:"
	}

	sort.Slice(terms, func(i, j int) bool {
		if terms[i].SourceTerm == terms[j].SourceTerm {
			return terms[i].ID < terms[j].ID
		}
		return terms[i].SourceTerm < terms[j].SourceTerm
	})

	var b strings.Builder
	b.WriteString("Glossary:")
	for _, term := range terms {
		b.WriteString("\n- ")
		b.WriteString(term.SourceTerm)
		b.WriteString(" → ")
		b.WriteString(term.TargetTerm)
		if term.Notes != "" {
			b.WriteString(" (")
			b.WriteString(term.Notes)
			b.WriteString(")")
		}
	}
	return b.String()
}

func (s *Store) ImportJSON(r io.Reader) error {
	var terms []Term
	if err := json.NewDecoder(r).Decode(&terms); err != nil {
		return fmt.Errorf("glossary.ImportJSON: decode: %w", err)
	}

	for _, term := range terms {
		if err := s.Add(term); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ExportJSON(w io.Writer) error {
	if s == nil {
		return fmt.Errorf("glossary.ExportJSON: nil store")
	}
	terms := s.List()
	if terms == nil {
		terms = []Term{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(terms); err != nil {
		return fmt.Errorf("glossary.ExportJSON: encode: %w", err)
	}
	return nil
}

func newID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}

	var b strings.Builder
	b.Grow(36)
	b.WriteString(hex.EncodeToString(raw[0:4]))
	b.WriteByte('-')
	b.WriteString(hex.EncodeToString(raw[4:6]))
	b.WriteByte('-')
	b.WriteString(hex.EncodeToString(raw[6:8]))
	b.WriteByte('-')
	b.WriteString(hex.EncodeToString(raw[8:10]))
	b.WriteByte('-')
	b.WriteString(hex.EncodeToString(raw[10:16]))
	return b.String(), nil
}
