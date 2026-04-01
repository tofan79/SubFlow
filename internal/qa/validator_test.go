package qa

import (
	"testing"
)

func TestQA_LineLength_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "This is a short line",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckLineLength)
	if found == nil {
		t.Fatal("expected LineLength result")
	}
	if !found.Passed {
		t.Errorf("expected pass, got fail: %s", found.Detail)
	}
}

func TestQA_LineLength_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "This line is way too long and exceeds the maximum of 42 characters allowed",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckLineLength)
	if found == nil {
		t.Fatal("expected LineLength result")
	}
	if found.Passed {
		t.Error("expected fail, got pass")
	}
	if found.Severity != SeverityError {
		t.Errorf("expected error severity, got %s", found.Severity)
	}
}

func TestQA_LineLength_CJKChar(t *testing.T) {
	v := NewValidator(Config{
		MaxCharsPerLine: 10,
		MaxLines:        2,
		MinDurationMS:   1000,
		MaxDurationMS:   7000,
		MaxCPS:          17.0,
		MinGapMS:        83,
	})

	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "日本語テスト文字", // 8 CJK characters
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckLineLength)
	if found == nil {
		t.Fatal("expected LineLength result")
	}
	if !found.Passed {
		t.Errorf("expected pass for 8 CJK chars with max 10, got fail: %s", found.Detail)
	}

	card.Text = "これは長すぎる日本語テスト" // 13 CJK characters
	results = v.ValidateCard(card, nil)
	found = findResult(results, CheckLineLength)
	if found == nil {
		t.Fatal("expected LineLength result")
	}
	if found.Passed {
		t.Error("expected fail for 13 CJK chars with max 10")
	}
}

func TestQA_LineCount_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "Line one\nLine two",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckLineCount)
	if found == nil {
		t.Fatal("expected LineCount result")
	}
	if !found.Passed {
		t.Errorf("expected pass, got fail: %s", found.Detail)
	}
}

func TestQA_LineCount_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "Line one\nLine two\nLine three",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckLineCount)
	if found == nil {
		t.Fatal("expected LineCount result")
	}
	if found.Passed {
		t.Error("expected fail, got pass")
	}
}

func TestQA_DurationMin_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1500,
		Text:    "Hello",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckDurationMin)
	if found == nil {
		t.Fatal("expected DurationMin result")
	}
	if !found.Passed {
		t.Errorf("expected pass for 1.5s duration, got fail: %s", found.Detail)
	}
}

func TestQA_DurationMin_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   500,
		Text:    "Hello",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckDurationMin)
	if found == nil {
		t.Fatal("expected DurationMin result")
	}
	if found.Passed {
		t.Error("expected fail for 0.5s duration")
	}
}

func TestQA_DurationMax_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   5000,
		Text:    "This is a normal duration subtitle",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckDurationMax)
	if found == nil {
		t.Fatal("expected DurationMax result")
	}
	if !found.Passed {
		t.Errorf("expected pass for 5s duration, got fail: %s", found.Detail)
	}
}

func TestQA_DurationMax_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   10000,
		Text:    "This subtitle is way too long in duration",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckDurationMax)
	if found == nil {
		t.Fatal("expected DurationMax result")
	}
	if found.Passed {
		t.Error("expected fail for 10s duration")
	}
}

func TestQA_CPS_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   3000,
		Text:    "This is about 30 characters!", // 30 chars / 3s = 10 CPS
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckCPS)
	if found == nil {
		t.Fatal("expected CPS result")
	}
	if !found.Passed {
		t.Errorf("expected pass for ~10 CPS, got fail: %s", found.Detail)
	}
}

func TestQA_CPS_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1000,
		Text:    "This text has way too many characters for one second!", // ~54 chars / 1s = 54 CPS
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckCPS)
	if found == nil {
		t.Fatal("expected CPS result")
	}
	if found.Passed {
		t.Error("expected fail for high CPS")
	}
}

func TestQA_CPS_ZeroDuration(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 1000,
		EndMS:   1000,
		Text:    "Hello",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckCPS)
	if found == nil {
		t.Fatal("expected CPS result")
	}
	if found.Passed {
		t.Error("expected fail for zero duration")
	}
	if found.Detail != "Durasi 0 atau negatif" {
		t.Errorf("unexpected detail: %s", found.Detail)
	}
}

func TestQA_Overlap_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	prev := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1000,
		Text:    "First",
	}
	curr := SubtitleCard{
		ID:      "test-2",
		Index:   2,
		StartMS: 1100,
		EndMS:   2000,
		Text:    "Second",
	}

	results := v.ValidateCard(curr, &prev)
	found := findResult(results, CheckOverlap)
	if found == nil {
		t.Fatal("expected Overlap result")
	}
	if !found.Passed {
		t.Errorf("expected pass, got fail: %s", found.Detail)
	}
}

func TestQA_Overlap_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	prev := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1500,
		Text:    "First",
	}
	curr := SubtitleCard{
		ID:      "test-2",
		Index:   2,
		StartMS: 1000,
		EndMS:   2000,
		Text:    "Second",
	}

	results := v.ValidateCard(curr, &prev)
	found := findResult(results, CheckOverlap)
	if found == nil {
		t.Fatal("expected Overlap result")
	}
	if found.Passed {
		t.Error("expected fail for overlapping cards")
	}
}

func TestQA_EmptyCard_Detected(t *testing.T) {
	v := NewValidator(DefaultConfig())
	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "   \n\t  ",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckEmptyCard)
	if found == nil {
		t.Fatal("expected EmptyCard result")
	}
	if found.Passed {
		t.Error("expected fail for empty card")
	}
}

func TestQA_Gap_Pass(t *testing.T) {
	v := NewValidator(DefaultConfig())
	prev := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1000,
		Text:    "First",
	}
	curr := SubtitleCard{
		ID:      "test-2",
		Index:   2,
		StartMS: 1100,
		EndMS:   2000,
		Text:    "Second",
	}

	results := v.ValidateCard(curr, &prev)
	found := findResult(results, CheckGap)
	if found == nil {
		t.Fatal("expected Gap result")
	}
	if !found.Passed {
		t.Errorf("expected pass for 100ms gap, got fail: %s", found.Detail)
	}
}

func TestQA_Gap_Fail(t *testing.T) {
	v := NewValidator(DefaultConfig())
	prev := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   1000,
		Text:    "First",
	}
	curr := SubtitleCard{
		ID:      "test-2",
		Index:   2,
		StartMS: 1050,
		EndMS:   2000,
		Text:    "Second",
	}

	results := v.ValidateCard(curr, &prev)
	found := findResult(results, CheckGap)
	if found == nil {
		t.Fatal("expected Gap result")
	}
	if found.Passed {
		t.Error("expected fail for 50ms gap (min 83ms)")
	}
}

func TestQA_Glossary_Warning(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Glossary = []GlossaryTerm{
		{SourceTerm: "AI", TargetTerm: "Kecerdasan Buatan", CaseSensitive: true},
	}
	v := NewValidator(cfg)

	card := SubtitleCard{
		ID:      "test-1",
		Index:   1,
		StartMS: 0,
		EndMS:   2000,
		Text:    "AI is amazing",
	}

	results := v.ValidateCard(card, nil)
	found := findResult(results, CheckGlossary)
	if found == nil {
		t.Fatal("expected Glossary result")
	}
	if found.Passed {
		t.Error("expected warning for untranslated glossary term")
	}
	if found.Severity != SeverityWarning {
		t.Errorf("expected warning severity, got %s", found.Severity)
	}
}

func findResult(results []Result, checkID CheckID) *Result {
	for _, r := range results {
		if r.CheckID == checkID {
			return &r
		}
	}
	return nil
}
