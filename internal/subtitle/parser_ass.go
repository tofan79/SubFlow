package subtitle

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var assTagPattern = regexp.MustCompile(`\{[^}]*\}`)

func ParseASS(r io.Reader) ([]Cue, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var cues []Cue
	lineNum := 0

	for s.Scan() {
		lineNum++
		line := strings.TrimSuffix(s.Text(), "\r")
		if lineNum == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}

		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "Dialogue:") {
			continue
		}

		cue, err := parseASSDialogueLine(trimmed)
		if err != nil {
			return nil, fmt.Errorf("ass: line %d: %w", lineNum, err)
		}
		cues = append(cues, cue)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("ass: scan failed: %w", err)
	}

	return cues, nil
}

func parseASSDialogueLine(line string) (Cue, error) {
	const prefix = "Dialogue:"

	rest := strings.TrimSpace(strings.TrimPrefix(line, prefix))
	fields := strings.SplitN(rest, ",", 10)
	if len(fields) != 10 {
		return Cue{}, fmt.Errorf("invalid dialogue line %q", line)
	}

	startMS, err := parseASSTimestamp(strings.TrimSpace(fields[1]))
	if err != nil {
		return Cue{}, fmt.Errorf("invalid start timestamp %q: %w", strings.TrimSpace(fields[1]), err)
	}
	endMS, err := parseASSTimestamp(strings.TrimSpace(fields[2]))
	if err != nil {
		return Cue{}, fmt.Errorf("invalid end timestamp %q: %w", strings.TrimSpace(fields[2]), err)
	}
	if endMS <= startMS {
		return Cue{}, fmt.Errorf("end time must be greater than start time")
	}

	text := StripASSTags(fields[9])
	text = strings.NewReplacer("\\N", "\n", "\\n", "\n").Replace(text)

	return Cue{
		Index:   0,
		StartMS: startMS,
		EndMS:   endMS,
		Text:    text,
	}, nil
}

func parseASSTimestamp(ts string) (int64, error) {
	parts := strings.Split(ts, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected H:MM:SS.cc")
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 {
		return 0, fmt.Errorf("invalid hours")
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("invalid minutes")
	}

	secParts := strings.Split(parts[2], ".")
	if len(secParts) != 2 {
		return 0, fmt.Errorf("expected centiseconds")
	}
	seconds, err := strconv.Atoi(secParts[0])
	if err != nil || seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("invalid seconds")
	}
	centiseconds, err := strconv.Atoi(secParts[1])
	if err != nil || centiseconds < 0 || centiseconds > 99 {
		return 0, fmt.Errorf("invalid centiseconds")
	}

	total := int64((((hours*60)+minutes)*60+seconds)*1000 + centiseconds*10)
	return total, nil
}

func StripASSTags(text string) string {
	return assTagPattern.ReplaceAllString(text, "")
}
