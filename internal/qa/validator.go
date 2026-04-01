package qa

import (
	"fmt"
	"strings"
)

type Validator struct {
	config Config
}

func NewValidator(cfg Config) *Validator {
	if cfg.MaxCharsPerLine <= 0 {
		cfg.MaxCharsPerLine = 42
	}
	if cfg.MaxLines <= 0 {
		cfg.MaxLines = 2
	}
	if cfg.MinDurationMS <= 0 {
		cfg.MinDurationMS = 1000
	}
	if cfg.MaxDurationMS <= 0 {
		cfg.MaxDurationMS = 7000
	}
	if cfg.MaxCPS <= 0 {
		cfg.MaxCPS = 17.0
	}
	if cfg.MinGapMS <= 0 {
		cfg.MinGapMS = 83
	}
	return &Validator{config: cfg}
}

func (v *Validator) Config() Config {
	return v.config
}

func (v *Validator) Validate(cards []SubtitleCard) []Result {
	var results []Result

	for i, card := range cards {
		results = append(results, v.checkLineLength(card)...)
		results = append(results, v.checkLineCount(card))
		results = append(results, v.checkDurationMin(card))
		results = append(results, v.checkDurationMax(card))
		results = append(results, v.checkCPS(card))
		results = append(results, v.checkEmptyCard(card))
		results = append(results, v.checkGlossary(card)...)

		if i > 0 {
			results = append(results, v.checkOverlap(cards[i-1], card))
			results = append(results, v.checkGap(cards[i-1], card))
		}
	}

	return results
}

func (v *Validator) ValidateCard(card SubtitleCard, prevCard *SubtitleCard) []Result {
	var results []Result

	results = append(results, v.checkLineLength(card)...)
	results = append(results, v.checkLineCount(card))
	results = append(results, v.checkDurationMin(card))
	results = append(results, v.checkDurationMax(card))
	results = append(results, v.checkCPS(card))
	results = append(results, v.checkEmptyCard(card))
	results = append(results, v.checkGlossary(card)...)

	if prevCard != nil {
		results = append(results, v.checkOverlap(*prevCard, card))
		results = append(results, v.checkGap(*prevCard, card))
	}

	return results
}

func (v *Validator) checkLineLength(card SubtitleCard) []Result {
	var results []Result
	lines := strings.Split(card.Text, "\n")

	for lineNum, line := range lines {
		charCount := len([]rune(line))
		if charCount > v.config.MaxCharsPerLine {
			results = append(results, Result{
				CardID:    card.ID,
				CardIndex: card.Index,
				CheckID:   CheckLineLength,
				Passed:    false,
				Severity:  SeverityError,
				Detail:    fmt.Sprintf("Baris %d: %d karakter (maks %d)", lineNum+1, charCount, v.config.MaxCharsPerLine),
			})
		}
	}

	if len(results) == 0 {
		results = append(results, Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckLineLength,
			Passed:    true,
			Severity:  SeverityPass,
		})
	}

	return results
}

func (v *Validator) checkLineCount(card SubtitleCard) Result {
	lines := strings.Split(card.Text, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines++
		}
	}

	if nonEmptyLines > v.config.MaxLines {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckLineCount,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("%d baris (maks %d)", nonEmptyLines, v.config.MaxLines),
		}
	}

	return Result{
		CardID:    card.ID,
		CardIndex: card.Index,
		CheckID:   CheckLineCount,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkDurationMin(card SubtitleCard) Result {
	duration := card.Duration()

	if duration < v.config.MinDurationMS {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckDurationMin,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("%.2f detik (min %.2f)", float64(duration)/1000.0, float64(v.config.MinDurationMS)/1000.0),
		}
	}

	return Result{
		CardID:    card.ID,
		CardIndex: card.Index,
		CheckID:   CheckDurationMin,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkDurationMax(card SubtitleCard) Result {
	duration := card.Duration()

	if duration > v.config.MaxDurationMS {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckDurationMax,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("%.2f detik (maks %.2f)", float64(duration)/1000.0, float64(v.config.MaxDurationMS)/1000.0),
		}
	}

	return Result{
		CardID:    card.ID,
		CardIndex: card.Index,
		CheckID:   CheckDurationMax,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkCPS(card SubtitleCard) Result {
	duration := card.Duration()

	if duration <= 0 {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckCPS,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    "Durasi 0 atau negatif",
		}
	}

	text := strings.ReplaceAll(card.Text, "\n", "")
	charCount := len([]rune(text))
	durationSec := float64(duration) / 1000.0
	cps := float64(charCount) / durationSec

	if cps > v.config.MaxCPS {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckCPS,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("%.1f CPS (maks %.1f)", cps, v.config.MaxCPS),
		}
	}

	return Result{
		CardID:    card.ID,
		CardIndex: card.Index,
		CheckID:   CheckCPS,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkOverlap(prev, curr SubtitleCard) Result {
	if curr.StartMS < prev.EndMS {
		overlapMS := prev.EndMS - curr.StartMS
		return Result{
			CardID:    curr.ID,
			CardIndex: curr.Index,
			CheckID:   CheckOverlap,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("Overlap %dms dengan card sebelumnya", overlapMS),
		}
	}

	return Result{
		CardID:    curr.ID,
		CardIndex: curr.Index,
		CheckID:   CheckOverlap,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkEmptyCard(card SubtitleCard) Result {
	trimmed := strings.TrimSpace(card.Text)

	if trimmed == "" {
		return Result{
			CardID:    card.ID,
			CardIndex: card.Index,
			CheckID:   CheckEmptyCard,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    "Card kosong",
		}
	}

	return Result{
		CardID:    card.ID,
		CardIndex: card.Index,
		CheckID:   CheckEmptyCard,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkGap(prev, curr SubtitleCard) Result {
	gap := curr.StartMS - prev.EndMS

	if gap > 0 && gap < v.config.MinGapMS {
		return Result{
			CardID:    curr.ID,
			CardIndex: curr.Index,
			CheckID:   CheckGap,
			Passed:    false,
			Severity:  SeverityError,
			Detail:    fmt.Sprintf("Gap %dms (min %dms)", gap, v.config.MinGapMS),
		}
	}

	return Result{
		CardID:    curr.ID,
		CardIndex: curr.Index,
		CheckID:   CheckGap,
		Passed:    true,
		Severity:  SeverityPass,
	}
}

func (v *Validator) checkGlossary(card SubtitleCard) []Result {
	if len(v.config.Glossary) == 0 {
		return nil
	}

	var results []Result
	textLower := strings.ToLower(card.Text)

	for _, term := range v.config.Glossary {
		var found bool
		if term.CaseSensitive {
			found = strings.Contains(card.Text, term.SourceTerm) && !strings.Contains(card.Text, term.TargetTerm)
		} else {
			found = strings.Contains(textLower, strings.ToLower(term.SourceTerm)) &&
				!strings.Contains(textLower, strings.ToLower(term.TargetTerm))
		}

		if found {
			results = append(results, Result{
				CardID:    card.ID,
				CardIndex: card.Index,
				CheckID:   CheckGlossary,
				Passed:    false,
				Severity:  SeverityWarning,
				Detail:    fmt.Sprintf("'%s' harus diterjemahkan '%s'", term.SourceTerm, term.TargetTerm),
			})
		}
	}

	return results
}

func CalculateCPS(text string, durationMS int64) float64 {
	if durationMS <= 0 {
		return 0
	}
	cleanText := strings.ReplaceAll(text, "\n", "")
	charCount := len([]rune(cleanText))
	return float64(charCount) / (float64(durationMS) / 1000.0)
}

func CountChars(text string) int {
	return len([]rune(strings.ReplaceAll(text, "\n", "")))
}

func CountLines(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}
	lines := strings.Split(text, "\n")
	count := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

func MaxLineLength(text string) int {
	lines := strings.Split(text, "\n")
	maxLen := 0
	for _, line := range lines {
		lineLen := len([]rune(line))
		if lineLen > maxLen {
			maxLen = lineLen
		}
	}
	return maxLen
}
