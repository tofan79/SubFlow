package subtitle

import (
	"math"
	"strings"
	"unicode"
)

type SegmentConfig struct {
	MaxCharsPerLine int
	MaxLines        int
	MaxCPS          float64
	MinDurationMS   int64
	MaxDurationMS   int64
	MinGapMS        int64
}

func DefaultSegmentConfig() SegmentConfig {
	return SegmentConfig{
		MaxCharsPerLine: 42,
		MaxLines:        2,
		MaxCPS:          17.0,
		MinDurationMS:   1000,
		MaxDurationMS:   7000,
		MinGapMS:        83,
	}
}

func Segment(cues []Cue, cfg SegmentConfig) []Cue {
	cfg = normalizeSegmentConfig(cfg)

	var out []Cue
	nextIndex := 1

	for _, cue := range cues {
		segments := segmentCue(cue, cfg)
		for _, seg := range segments {
			seg.Index = nextIndex
			nextIndex++
			out = append(out, seg)
		}
	}

	return out
}

func SplitLongLine(line string, maxChars int) []string {
	if line == "" {
		return []string{""}
	}
	if maxChars <= 0 {
		return []string{line}
	}

	words := splitWords(line)
	if len(words) == 0 {
		return []string{""}
	}

	var out []string
	current := ""

	for _, word := range words {
		if current == "" {
			current = word
			continue
		}

		candidate := current + " " + word
		if CountChars(candidate) <= maxChars {
			current = candidate
			continue
		}

		out = append(out, current)
		current = word
	}

	if current != "" {
		out = append(out, current)
	}

	if len(out) == 0 {
		return []string{""}
	}

	return out
}

func CountChars(s string) int {
	count := 0
	runes := []rune(s)
	for _, r := range runes {
		if isWideCJK(r) {
			count += 2
			continue
		}
		count++
	}
	return count
}

func normalizeSegmentConfig(cfg SegmentConfig) SegmentConfig {
	def := DefaultSegmentConfig()
	if cfg.MaxCharsPerLine <= 0 {
		cfg.MaxCharsPerLine = def.MaxCharsPerLine
	}
	if cfg.MaxLines <= 0 {
		cfg.MaxLines = def.MaxLines
	}
	if cfg.MaxCPS <= 0 {
		cfg.MaxCPS = def.MaxCPS
	}
	if cfg.MinDurationMS <= 0 {
		cfg.MinDurationMS = def.MinDurationMS
	}
	if cfg.MaxDurationMS <= 0 {
		cfg.MaxDurationMS = def.MaxDurationMS
	}
	if cfg.MinGapMS < 0 {
		cfg.MinGapMS = def.MinGapMS
	}
	if cfg.MaxDurationMS < cfg.MinDurationMS {
		cfg.MaxDurationMS = cfg.MinDurationMS
	}
	return cfg
}

func segmentCue(cue Cue, cfg SegmentConfig) []Cue {
	lines := splitCueLines(cue.Text)
	lineLimit := cfg.MaxCharsPerLine
	maxSegmentChars := int(math.Floor(cfg.MaxCPS * float64(cfg.MaxDurationMS) / 1000.0))
	if maxSegmentChars > 0 && maxSegmentChars < lineLimit {
		lineLimit = maxSegmentChars
	}

	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		parts := SplitLongLine(line, lineLimit)
		if len(parts) == 0 {
			parts = []string{""}
		}
		wrapped = append(wrapped, parts...)
	}
	if len(wrapped) == 0 {
		wrapped = []string{""}
	}

	chunks := groupWrappedLines(wrapped, cfg.MaxLines, maxSegmentChars)
	durations := allocateDurations(chunks, cue.EndMS-cue.StartMS, cfg)

	out := make([]Cue, 0, len(chunks))
	start := cue.StartMS
	for i, chunk := range chunks {
		text := strings.Join(chunk, "\n")
		duration := durations[i]
		if duration < cfg.MinDurationMS {
			duration = cfg.MinDurationMS
		}
		if duration > cfg.MaxDurationMS {
			duration = cfg.MaxDurationMS
		}
		end := start + duration
		out = append(out, Cue{
			StartMS: start,
			EndMS:   end,
			Text:    text,
		})
		start = end + cfg.MinGapMS
	}

	return out
}

func splitCueLines(text string) []string {
	if text == "" {
		return []string{""}
	}
	parts := strings.Split(text, "\n")
	if len(parts) == 0 {
		return []string{""}
	}
	return parts
}

func groupWrappedLines(lines []string, maxLines int, maxSegmentChars int) [][]string {
	if maxLines <= 0 {
		maxLines = len(lines)
	}

	var chunks [][]string
	var current []string
	currentChars := 0

	flush := func() {
		if len(current) == 0 {
			return
		}
		chunk := make([]string, len(current))
		copy(chunk, current)
		chunks = append(chunks, chunk)
		current = nil
		currentChars = 0
	}

	for _, line := range lines {
		lineChars := CountChars(line)
		wouldOverflowLines := len(current) >= maxLines
		wouldOverflowChars := maxSegmentChars > 0 && len(current) > 0 && currentChars+lineChars > maxSegmentChars
		if wouldOverflowLines || wouldOverflowChars {
			flush()
		}

		current = append(current, line)
		currentChars += lineChars
	}

	flush()
	if len(chunks) == 0 {
		return [][]string{{""}}
	}
	return chunks
}

func allocateDurations(chunks [][]string, cueDurationMS int64, cfg SegmentConfig) []int64 {
	durations := make([]int64, len(chunks))
	if len(chunks) == 0 {
		return durations
	}

	intrinsic := make([]int64, len(chunks))
	for i, chunk := range chunks {
		width := 0
		for _, line := range chunk {
			width += CountChars(line)
		}
		d := int64(math.Ceil(float64(width) / cfg.MaxCPS * 1000.0))
		if d < cfg.MinDurationMS {
			d = cfg.MinDurationMS
		}
		if d > cfg.MaxDurationMS {
			d = cfg.MaxDurationMS
		}
		intrinsic[i] = d
	}

	copy(durations, intrinsic)
	if cueDurationMS <= 0 {
		return durations
	}

	gapTotal := cfg.MinGapMS * int64(len(chunks)-1)
	available := cueDurationMS - gapTotal
	if available <= 0 {
		return durations
	}

	minTotal := cfg.MinDurationMS * int64(len(chunks))
	if available < minTotal {
		return durations
	}

	currentTotal := int64(0)
	for _, d := range durations {
		currentTotal += d
	}

	if currentTotal < available {
		durations[len(durations)-1] += available - currentTotal
		if durations[len(durations)-1] > cfg.MaxDurationMS {
			durations[len(durations)-1] = cfg.MaxDurationMS
		}
		return durations
	}

	if currentTotal == available {
		return durations
	}

	scaled := make([]int64, len(chunks))
	scaledTotal := int64(0)
	for i, d := range durations {
		s := int64(math.Round(float64(d) * float64(available) / float64(currentTotal)))
		if s < cfg.MinDurationMS {
			s = cfg.MinDurationMS
		}
		if s > cfg.MaxDurationMS {
			s = cfg.MaxDurationMS
		}
		scaled[i] = s
		scaledTotal += s
	}

	delta := available - scaledTotal
	if delta > 0 {
		for i := len(scaled) - 1; i >= 0 && delta > 0; i-- {
			room := cfg.MaxDurationMS - scaled[i]
			if room <= 0 {
				continue
			}
			add := room
			if add > delta {
				add = delta
			}
			scaled[i] += add
			delta -= add
		}
	} else if delta < 0 {
		for i := len(scaled) - 1; i >= 0 && delta < 0; i-- {
			room := scaled[i] - cfg.MinDurationMS
			if room <= 0 {
				continue
			}
			sub := room
			if -sub < delta {
				sub = -delta
			}
			scaled[i] -= sub
			delta += sub
		}
	}

	return scaled
}

func splitWords(line string) []string {
	runes := []rune(line)
	var words []string
	var current []rune

	flush := func() {
		if len(current) == 0 {
			return
		}
		words = append(words, string(current))
		current = nil
	}

	for _, r := range runes {
		if unicode.IsSpace(r) {
			flush()
			continue
		}
		current = append(current, r)
	}
	flush()

	return words
}

func isWideCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || (r >= 0x3040 && r <= 0x30FF)
}
