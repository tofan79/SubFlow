package glossary

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStore_CRUD(t *testing.T) {
	s := NewStore()

	term := Term{SourceTerm: "hello", TargetTerm: "halo", Category: "technical", Notes: "greeting"}
	if err := s.Add(term); err != nil {
		t.Fatalf("Add error: %v", err)
	}

	list := s.List()
	if len(list) != 1 {
		t.Fatalf("expected 1 term, got %d", len(list))
	}
	if list[0].ID == "" {
		t.Fatal("expected generated ID")
	}

	got, ok := s.Get(list[0].ID)
	if !ok {
		t.Fatal("expected Get to find term")
	}
	if got.SourceTerm != "hello" || got.TargetTerm != "halo" || got.Category != "technical" || got.Notes != "greeting" {
		t.Fatalf("Get mismatch: %+v", got)
	}

	if !s.Delete(list[0].ID) {
		t.Fatal("expected Delete to return true")
	}
	if _, ok := s.Get(list[0].ID); ok {
		t.Fatal("expected term to be deleted")
	}
	if s.Delete("missing") {
		t.Fatal("expected Delete missing to return false")
	}
	if len(s.List()) != 0 {
		t.Fatal("expected empty list after delete")
	}
}

func TestStore_FindBySource(t *testing.T) {
	s := NewStore()

	if err := s.Add(Term{ID: "1", SourceTerm: "Apple", TargetTerm: "Apfel"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if err := s.Add(Term{ID: "2", SourceTerm: "Apple", TargetTerm: "りんご", CaseSensitive: true}); err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if err := s.Add(Term{ID: "3", SourceTerm: "Pear", TargetTerm: "Birne"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}

	terms := s.FindBySource("apple")
	if len(terms) != 1 || terms[0].ID != "1" {
		t.Fatalf("expected only case-insensitive match, got %+v", terms)
	}

	terms = s.FindBySource("Apple")
	if len(terms) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(terms))
	}
	if terms[0].ID != "1" || terms[1].ID != "2" {
		t.Fatalf("unexpected order or results: %+v", terms)
	}

	terms = s.FindBySource("APPLE")
	if len(terms) != 1 || terms[0].ID != "1" {
		t.Fatalf("expected only non-case-sensitive match, got %+v", terms)
	}
}

func TestStore_BuildPromptSection(t *testing.T) {
	s := NewStore()
	if err := s.Add(Term{ID: "1", SourceTerm: "hero", TargetTerm: "主人公", Notes: "main character"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if err := s.Add(Term{ID: "2", SourceTerm: "castle", TargetTerm: "城"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}

	got := s.BuildPromptSection()
	want := "Glossary:\n- castle → 城\n- hero → 主人公 (main character)"
	if got != want {
		t.Fatalf("BuildPromptSection mismatch:\nwant: %q\n got: %q", want, got)
	}
}

func TestStore_ImportExportJSON(t *testing.T) {
	orig := NewStore()
	if err := orig.Add(Term{ID: "a", SourceTerm: "alpha", TargetTerm: "アルファ", Category: "technical"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if err := orig.Add(Term{ID: "b", SourceTerm: "beta", TargetTerm: "ベータ", CaseSensitive: true, Notes: "second"}); err != nil {
		t.Fatalf("Add error: %v", err)
	}

	var buf bytes.Buffer
	if err := orig.ExportJSON(&buf); err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}

	var exported []Term
	if err := json.Unmarshal(buf.Bytes(), &exported); err != nil {
		t.Fatalf("unmarshal exported json: %v", err)
	}
	if len(exported) != 2 {
		t.Fatalf("expected 2 exported terms, got %d", len(exported))
	}

	copyStore := NewStore()
	if err := copyStore.ImportJSON(strings.NewReader(buf.String())); err != nil {
		t.Fatalf("ImportJSON error: %v", err)
	}

	terms := copyStore.List()
	if len(terms) != 2 {
		t.Fatalf("expected 2 imported terms, got %d", len(terms))
	}

	if terms[0].ID != "a" || terms[0].SourceTerm != "alpha" || terms[0].TargetTerm != "アルファ" || terms[0].Category != "technical" {
		t.Fatalf("term[0] mismatch: %+v", terms[0])
	}
	if terms[1].ID != "b" || terms[1].SourceTerm != "beta" || terms[1].TargetTerm != "ベータ" || !terms[1].CaseSensitive || terms[1].Notes != "second" {
		t.Fatalf("term[1] mismatch: %+v", terms[1])
	}
}
