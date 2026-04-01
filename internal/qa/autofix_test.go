package qa

import (
	"testing"
)

func TestQA_AutoFix_AllChecks(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   500,
			Text:    "This line is definitely way too long and will need to be wrapped properly for display",
		},
		{
			ID:      "card-2",
			Index:   2,
			StartMS: 400,
			EndMS:   1500,
			Text:    "Overlapping card",
		},
		{
			ID:      "card-3",
			Index:   3,
			StartMS: 1510,
			EndMS:   2000,
			Text:    "Small gap card",
		},
	}

	fixed, logs, err := fixer.AutoFix(cards)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(logs) == 0 {
		t.Error("expected some auto-fix logs")
	}

	if len(fixed) != 3 {
		t.Errorf("expected 3 cards, got %d", len(fixed))
	}

	for _, log := range logs {
		t.Logf("Fix: card=%s check=%s action=%s", log.CardID, log.CheckID, log.Action)
	}
}

func TestQA_AutoFix_MaxRetry(t *testing.T) {
	cfg := Config{
		MaxCharsPerLine: 5,
		MaxLines:        1,
		MinDurationMS:   1000,
		MaxDurationMS:   7000,
		MaxCPS:          17.0,
		MinGapMS:        83,
	}
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   2000,
			Text:    "This is an impossible text that cannot be fixed with max 5 chars per line",
		},
	}

	_, logs, err := fixer.AutoFix(cards)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loopCounts := make(map[int]int)
	for _, log := range logs {
		loopCounts[log.LoopNumber]++
	}

	for loop := range loopCounts {
		if loop > MaxAutoFixLoops {
			t.Errorf("loop %d exceeds max %d", loop, MaxAutoFixLoops)
		}
	}

	t.Logf("Loops used: %v", loopCounts)
}

func TestQA_AutoFix_LineLength(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   3000,
			Text:    "This line has more than forty-two characters and should be wrapped",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	foundFix := false
	for _, log := range logs {
		if log.CheckID == CheckLineLength {
			foundFix = true
			break
		}
	}
	if !foundFix {
		t.Error("expected LineLength fix")
	}

	maxLen := MaxLineLength(fixed[0].Text)
	if maxLen > cfg.MaxCharsPerLine {
		t.Errorf("line still too long after fix: %d chars", maxLen)
	}
}

func TestQA_AutoFix_DurationMin(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   500,
			Text:    "Short",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	foundFix := false
	for _, log := range logs {
		if log.CheckID == CheckDurationMin {
			foundFix = true
			break
		}
	}
	if !foundFix {
		t.Error("expected DurationMin fix")
	}

	if fixed[0].Duration() < cfg.MinDurationMS {
		t.Errorf("duration still too short: %dms", fixed[0].Duration())
	}
}

func TestQA_AutoFix_CPS(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   1000,
			Text:    "This has way too many characters for just one second of display time!",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	foundFix := false
	for _, log := range logs {
		if log.CheckID == CheckCPS {
			foundFix = true
			break
		}
	}
	if !foundFix {
		t.Error("expected CPS fix")
	}

	cps := CalculateCPS(fixed[0].Text, fixed[0].Duration())
	if cps > cfg.MaxCPS {
		t.Errorf("CPS still too high: %.1f", cps)
	}
}

func TestQA_AutoFix_Overlap(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   2000,
			Text:    "First card",
		},
		{
			ID:      "card-2",
			Index:   2,
			StartMS: 1500,
			EndMS:   3000,
			Text:    "Overlapping",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	foundFix := false
	for _, log := range logs {
		if log.CheckID == CheckOverlap {
			foundFix = true
			break
		}
	}
	if !foundFix {
		t.Error("expected Overlap fix")
	}

	if fixed[1].StartMS < fixed[0].EndMS {
		t.Errorf("cards still overlap: card1 ends at %d, card2 starts at %d", fixed[0].EndMS, fixed[1].StartMS)
	}
}

func TestQA_AutoFix_Gap(t *testing.T) {
	cfg := DefaultConfig()
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   1000,
			Text:    "First",
		},
		{
			ID:      "card-2",
			Index:   2,
			StartMS: 1050,
			EndMS:   2000,
			Text:    "Second",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	foundFix := false
	for _, log := range logs {
		if log.CheckID == CheckGap {
			foundFix = true
			break
		}
	}
	if !foundFix {
		t.Error("expected Gap fix")
	}

	gap := fixed[1].StartMS - fixed[0].EndMS
	if gap < cfg.MinGapMS {
		t.Errorf("gap still too small: %dms", gap)
	}
}

func TestQA_AutoFix_EmptyCard_Removal(t *testing.T) {
	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   1000,
			Text:    "Valid card",
		},
		{
			ID:      "card-2",
			Index:   2,
			StartMS: 1100,
			EndMS:   2000,
			Text:    "   ",
		},
		{
			ID:      "card-3",
			Index:   3,
			StartMS: 2100,
			EndMS:   3000,
			Text:    "Another valid",
		},
	}

	cleaned := RemoveEmptyCards(cards)

	if len(cleaned) != 2 {
		t.Errorf("expected 2 cards after removal, got %d", len(cleaned))
	}

	if cleaned[0].Index != 1 || cleaned[1].Index != 2 {
		t.Error("indices not renumbered correctly")
	}
}

func TestQA_AutoFix_GlossaryNotFixed(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Glossary = []GlossaryTerm{
		{SourceTerm: "AI", TargetTerm: "Kecerdasan Buatan", CaseSensitive: true},
	}
	fixer := NewAutoFixer(cfg)

	cards := []SubtitleCard{
		{
			ID:      "card-1",
			Index:   1,
			StartMS: 0,
			EndMS:   2000,
			Text:    "AI is amazing",
		},
	}

	fixed, logs, _ := fixer.AutoFix(cards)

	for _, log := range logs {
		if log.CheckID == CheckGlossary {
			t.Error("glossary should NOT be auto-fixed")
		}
	}

	if fixed[0].Text != "AI is amazing" {
		t.Error("text should remain unchanged for glossary issues")
	}

	if fixed[0].QAStatus != StatusWarn {
		t.Errorf("expected warn status for glossary issue, got %s", fixed[0].QAStatus)
	}
}

func TestQA_CalculateCPS(t *testing.T) {
	tests := []struct {
		text       string
		durationMS int64
		expected   float64
	}{
		{"Hello", 1000, 5.0},
		{"Hello World", 2000, 5.5},
		{"日本語", 1000, 3.0},
		{"Line1\nLine2", 2000, 5.0},
		{"", 1000, 0.0},
		{"Text", 0, 0.0},
	}

	for _, tc := range tests {
		result := CalculateCPS(tc.text, tc.durationMS)
		if !floatEquals(result, tc.expected, 0.01) {
			t.Errorf("CalculateCPS(%q, %d) = %.2f, expected %.2f", tc.text, tc.durationMS, result, tc.expected)
		}
	}
}

func TestQA_CountChars(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"Hello", 5},
		{"Hello\nWorld", 10},
		{"日本語", 3},
		{"", 0},
	}

	for _, tc := range tests {
		result := CountChars(tc.text)
		if result != tc.expected {
			t.Errorf("CountChars(%q) = %d, expected %d", tc.text, result, tc.expected)
		}
	}
}

func TestQA_CountLines(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"Hello", 1},
		{"Hello\nWorld", 2},
		{"", 0},
		{"  \n  ", 0},
		{"Line1\n\nLine2", 2},
	}

	for _, tc := range tests {
		result := CountLines(tc.text)
		if result != tc.expected {
			t.Errorf("CountLines(%q) = %d, expected %d", tc.text, result, tc.expected)
		}
	}
}

func floatEquals(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}
