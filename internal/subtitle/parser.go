package subtitle

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Cue struct {
	Index   int    // 1-based display order
	StartMS int64  // milliseconds
	EndMS   int64  // milliseconds
	Text    string // "\n" as line separator
}

func ParseSRT(r io.Reader) ([]Cue, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	const (
		stateIndex = iota
		stateTime
		stateText
	)

	state := stateIndex
	lineNum := 0
	indexLineNum := 0

	var cues []Cue
	var cur Cue
	var textLines []string

	flushCue := func() {
		cur.Text = strings.Join(textLines, "\n")
		cues = append(cues, cur)
		cur = Cue{}
		textLines = nil
		state = stateIndex
		indexLineNum = 0
	}

	for s.Scan() {
		lineNum++
		line := s.Text()
		line = strings.TrimSuffix(line, "\r")
		if lineNum == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}

		switch state {
		case stateIndex:
			if strings.TrimSpace(line) == "" {
				continue
			}
			idx, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil {
				return nil, fmt.Errorf("srt: line %d: invalid cue index %q", lineNum, line)
			}
			cur = Cue{Index: idx}
			indexLineNum = lineNum
			state = stateTime

		case stateTime:
			if strings.TrimSpace(line) == "" {
				continue
			}
			startMS, endMS, err := parseSRTTimeRange(line)
			if err != nil {
				return nil, fmt.Errorf("srt: line %d: %w", lineNum, err)
			}
			if endMS <= startMS {
				return nil, fmt.Errorf("srt: line %d: end time must be greater than start time", lineNum)
			}
			cur.StartMS = startMS
			cur.EndMS = endMS
			state = stateText

		case stateText:
			if strings.TrimSpace(line) == "" {
				flushCue()
				continue
			}
			textLines = append(textLines, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("srt: scan failed: %w", err)
	}

	if state == stateTime {
		return nil, fmt.Errorf("srt: line %d: missing timestamp after cue index", indexLineNum)
	}
	if state == stateText {
		flushCue()
	}

	return cues, nil
}

func parseSRTTimeRange(line string) (startMS int64, endMS int64, err error) {
	left, right, ok := strings.Cut(line, "-->")
	if !ok {
		return 0, 0, fmt.Errorf("invalid timestamp range %q (missing -->)", line)
	}
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || right == "" {
		return 0, 0, fmt.Errorf("invalid timestamp range %q", line)
	}

	if f := strings.Fields(left); len(f) > 0 {
		left = f[0]
	}
	if f := strings.Fields(right); len(f) > 0 {
		right = f[0]
	}

	startMS, parseErr := parseSRTTimestamp(left)
	if parseErr != nil {
		return 0, 0, fmt.Errorf("invalid start timestamp %q: %v", left, parseErr)
	}
	endMS, parseErr = parseSRTTimestamp(right)
	if parseErr != nil {
		return 0, 0, fmt.Errorf("invalid end timestamp %q: %v", right, parseErr)
	}
	return startMS, endMS, nil
}

func parseSRTTimestamp(ts string) (int64, error) {
	parts := strings.Split(ts, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("expected HH:MM:SS,mmm")
	}

	h, err := strconv.Atoi(parts[0])
	if err != nil || h < 0 {
		return 0, fmt.Errorf("invalid hours")
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("invalid minutes")
	}

	secParts := strings.Split(parts[2], ",")
	if len(secParts) != 2 {
		return 0, fmt.Errorf("expected comma milliseconds")
	}
	s, err := strconv.Atoi(secParts[0])
	if err != nil || s < 0 || s > 59 {
		return 0, fmt.Errorf("invalid seconds")
	}
	ms, err := strconv.Atoi(secParts[1])
	if err != nil || ms < 0 || ms > 999 {
		return 0, fmt.Errorf("invalid milliseconds")
	}

	total := int64(((h*3600 + m*60 + s) * 1000) + ms)
	return total, nil
}
