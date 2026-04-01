package qa

import (
	"fmt"
	"strings"
)

const MaxAutoFixLoops = 3

type AutoFixer struct {
	validator *Validator
}

func NewAutoFixer(cfg Config) *AutoFixer {
	return &AutoFixer{
		validator: NewValidator(cfg),
	}
}

func (f *AutoFixer) Config() Config {
	return f.validator.Config()
}

func (f *AutoFixer) AutoFix(cards []SubtitleCard) ([]SubtitleCard, []Log, error) {
	var allLogs []Log
	result := copyCards(cards)

	for loop := 1; loop <= MaxAutoFixLoops; loop++ {
		results := f.validator.Validate(result)

		hasErrors := false
		for _, r := range results {
			if !r.Passed && r.Severity == SeverityError {
				hasErrors = true
				break
			}
		}

		if !hasErrors {
			break
		}

		var loopLogs []Log
		result, loopLogs = f.applyFixes(result, results, loop)
		allLogs = append(allLogs, loopLogs...)

		if len(loopLogs) == 0 {
			break
		}
	}

	for i := range result {
		result[i] = f.updateQAStatus(result[i], f.validator.ValidateCard(result[i], getPrevCard(result, i)))
	}

	return result, allLogs, nil
}

func (f *AutoFixer) applyFixes(cards []SubtitleCard, results []Result, loop int) ([]SubtitleCard, []Log) {
	var logs []Log
	fixed := copyCards(cards)

	cardErrors := make(map[string][]Result)
	for _, r := range results {
		if !r.Passed && r.Severity == SeverityError && IsAutoFixable(r.CheckID) {
			cardErrors[r.CardID] = append(cardErrors[r.CardID], r)
		}
	}

	var newCards []SubtitleCard
	skipIndices := make(map[int]bool)

	for i, card := range fixed {
		if skipIndices[i] {
			continue
		}

		errs, hasErrors := cardErrors[card.ID]
		if !hasErrors {
			newCards = append(newCards, card)
			continue
		}

		for _, err := range errs {
			var log Log
			card, log = f.fixSingle(card, err, loop, getPrevCard(newCards, len(newCards)))
			if log.Action != "" {
				logs = append(logs, log)
			}
		}

		newCards = append(newCards, card)
	}

	for i := range newCards {
		newCards[i].Index = i + 1
	}

	return newCards, logs
}

func (f *AutoFixer) fixSingle(card SubtitleCard, err Result, loop int, prevCard *SubtitleCard) (SubtitleCard, Log) {
	cfg := f.validator.Config()

	switch err.CheckID {
	case CheckLineLength:
		return f.fixLineLength(card, cfg.MaxCharsPerLine, loop)
	case CheckLineCount:
		return f.fixLineCount(card, cfg.MaxLines, loop)
	case CheckDurationMin:
		return f.fixDurationMin(card, cfg.MinDurationMS, loop)
	case CheckDurationMax:
		return f.fixDurationMax(card, cfg.MaxDurationMS, loop)
	case CheckCPS:
		return f.fixCPS(card, cfg.MaxCPS, loop)
	case CheckOverlap:
		if prevCard != nil {
			return f.fixOverlap(*prevCard, card, loop)
		}
		return card, Log{}
	case CheckEmptyCard:
		return card, Log{
			CardID:     card.ID,
			CardIndex:  card.Index,
			CheckID:    CheckEmptyCard,
			Action:     "mark_for_removal",
			OldValue:   card.Text,
			NewValue:   "",
			LoopNumber: loop,
		}
	case CheckGap:
		if prevCard != nil {
			return f.fixGap(*prevCard, card, cfg.MinGapMS, loop)
		}
		return card, Log{}
	default:
		return card, Log{}
	}
}

func (f *AutoFixer) fixLineLength(card SubtitleCard, maxChars int, loop int) (SubtitleCard, Log) {
	oldText := card.Text
	lines := strings.Split(card.Text, "\n")
	var newLines []string

	for _, line := range lines {
		if len([]rune(line)) <= maxChars {
			newLines = append(newLines, line)
			continue
		}
		wrapped := wrapLine(line, maxChars)
		newLines = append(newLines, wrapped...)
	}

	card.Text = strings.Join(newLines, "\n")

	if card.Text == oldText {
		return card, Log{}
	}

	return card, Log{
		CardID:     card.ID,
		CardIndex:  card.Index,
		CheckID:    CheckLineLength,
		Action:     "reflow_text",
		OldValue:   oldText,
		NewValue:   card.Text,
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixLineCount(card SubtitleCard, maxLines int, loop int) (SubtitleCard, Log) {
	lines := strings.Split(card.Text, "\n")
	var nonEmpty []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmpty = append(nonEmpty, line)
		}
	}

	if len(nonEmpty) <= maxLines {
		return card, Log{}
	}

	oldText := card.Text
	card.Text = strings.Join(nonEmpty[:maxLines], "\n")

	return card, Log{
		CardID:     card.ID,
		CardIndex:  card.Index,
		CheckID:    CheckLineCount,
		Action:     "truncate_lines",
		OldValue:   oldText,
		NewValue:   card.Text,
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixDurationMin(card SubtitleCard, minMS int64, loop int) (SubtitleCard, Log) {
	oldEnd := card.EndMS
	card.EndMS = card.StartMS + minMS

	return card, Log{
		CardID:     card.ID,
		CardIndex:  card.Index,
		CheckID:    CheckDurationMin,
		Action:     "extend_end_timestamp",
		OldValue:   fmt.Sprintf("end=%dms", oldEnd),
		NewValue:   fmt.Sprintf("end=%dms", card.EndMS),
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixDurationMax(card SubtitleCard, maxMS int64, loop int) (SubtitleCard, Log) {
	oldEnd := card.EndMS
	card.EndMS = card.StartMS + maxMS

	return card, Log{
		CardID:     card.ID,
		CardIndex:  card.Index,
		CheckID:    CheckDurationMax,
		Action:     "trim_end_timestamp",
		OldValue:   fmt.Sprintf("end=%dms", oldEnd),
		NewValue:   fmt.Sprintf("end=%dms", card.EndMS),
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixCPS(card SubtitleCard, maxCPS float64, loop int) (SubtitleCard, Log) {
	text := strings.ReplaceAll(card.Text, "\n", "")
	charCount := len([]rune(text))
	requiredDurationSec := float64(charCount) / maxCPS
	requiredDurationMS := int64(requiredDurationSec * 1000)

	if requiredDurationMS <= card.Duration() {
		return card, Log{}
	}

	oldEnd := card.EndMS
	card.EndMS = card.StartMS + requiredDurationMS

	return card, Log{
		CardID:     card.ID,
		CardIndex:  card.Index,
		CheckID:    CheckCPS,
		Action:     "extend_end_timestamp",
		OldValue:   fmt.Sprintf("end=%dms (%.1f CPS)", oldEnd, float64(charCount)/(float64(oldEnd-card.StartMS)/1000.0)),
		NewValue:   fmt.Sprintf("end=%dms (%.1f CPS)", card.EndMS, maxCPS),
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixOverlap(prev, curr SubtitleCard, loop int) (SubtitleCard, Log) {
	if curr.StartMS >= prev.EndMS {
		return curr, Log{}
	}

	oldStart := curr.StartMS
	curr.StartMS = prev.EndMS

	return curr, Log{
		CardID:     curr.ID,
		CardIndex:  curr.Index,
		CheckID:    CheckOverlap,
		Action:     "adjust_start_timestamp",
		OldValue:   fmt.Sprintf("start=%dms", oldStart),
		NewValue:   fmt.Sprintf("start=%dms", curr.StartMS),
		LoopNumber: loop,
	}
}

func (f *AutoFixer) fixGap(prev, curr SubtitleCard, minGapMS int64, loop int) (SubtitleCard, Log) {
	gap := curr.StartMS - prev.EndMS
	if gap >= minGapMS || gap <= 0 {
		return curr, Log{}
	}

	oldStart := curr.StartMS
	curr.StartMS = prev.EndMS + minGapMS

	return curr, Log{
		CardID:     curr.ID,
		CardIndex:  curr.Index,
		CheckID:    CheckGap,
		Action:     "shift_start_timestamp",
		OldValue:   fmt.Sprintf("start=%dms (gap=%dms)", oldStart, gap),
		NewValue:   fmt.Sprintf("start=%dms (gap=%dms)", curr.StartMS, minGapMS),
		LoopNumber: loop,
	}
}

func (f *AutoFixer) updateQAStatus(card SubtitleCard, results []Result) SubtitleCard {
	hasError := false
	hasWarning := false

	for _, r := range results {
		if !r.Passed {
			if r.Severity == SeverityError {
				hasError = true
			} else if r.Severity == SeverityWarning {
				hasWarning = true
			}
		}
	}

	if hasError {
		card.QAStatus = StatusError
	} else if hasWarning {
		card.QAStatus = StatusWarn
	} else {
		card.QAStatus = StatusPass
	}

	return card
}

func wrapLine(line string, maxChars int) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current []string
	currentLen := 0

	for _, word := range words {
		wordLen := len([]rune(word))

		if wordLen > maxChars {
			if len(current) > 0 {
				lines = append(lines, strings.Join(current, " "))
				current = nil
				currentLen = 0
			}
			lines = append(lines, word[:maxChars])
			continue
		}

		newLen := currentLen + wordLen
		if currentLen > 0 {
			newLen++
		}

		if newLen > maxChars {
			lines = append(lines, strings.Join(current, " "))
			current = []string{word}
			currentLen = wordLen
		} else {
			current = append(current, word)
			currentLen = newLen
		}
	}

	if len(current) > 0 {
		lines = append(lines, strings.Join(current, " "))
	}

	return lines
}

func copyCards(cards []SubtitleCard) []SubtitleCard {
	result := make([]SubtitleCard, len(cards))
	copy(result, cards)
	return result
}

func getPrevCard(cards []SubtitleCard, currentIndex int) *SubtitleCard {
	if currentIndex > 0 && currentIndex <= len(cards) {
		return &cards[currentIndex-1]
	}
	return nil
}

func RemoveEmptyCards(cards []SubtitleCard) []SubtitleCard {
	var result []SubtitleCard
	for _, card := range cards {
		if strings.TrimSpace(card.Text) != "" {
			card.Index = len(result) + 1
			result = append(result, card)
		}
	}
	return result
}
