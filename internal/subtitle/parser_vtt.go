package subtitle

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParseVTT(r io.Reader) ([]Cue, error) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	const (
		stateHeader = iota
		stateCueStart
		stateCueText
	)

	state := stateHeader
	lineNum := 0
	headerLineNum := 0

	var cues []Cue
	var cur Cue
	var textLines []string

	flushCue := func() {
		cur.Text = strings.Join(textLines, "\n")
		cues = append(cues, cur)
		cur = Cue{}
		textLines = nil
		state = stateCueStart
	}

	for s.Scan() {
		lineNum++
		line := strings.TrimSuffix(s.Text(), "\r")
		if lineNum == 1 {
			line = strings.TrimPrefix(line, "\uFEFF")
		}

		switch state {
		case stateHeader:
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.TrimSpace(line) != "WEBVTT" {
				return nil, fmt.Errorf("vtt: line %d: missing WEBVTT header", lineNum)
			}
			headerLineNum = lineNum
			state = stateCueStart

		case stateCueStart:
			if strings.TrimSpace(line) == "" {
				continue
			}
			if strings.Contains(line, "-->") {
				startMS, endMS, err := parseVTTTimeRange(line)
				if err != nil {
					return nil, fmt.Errorf("vtt: line %d: %w", lineNum, err)
				}
				if endMS <= startMS {
					return nil, fmt.Errorf("vtt: line %d: end time must be greater than start time", lineNum)
				}
				cur = Cue{Index: len(cues) + 1, StartMS: startMS, EndMS: endMS}
				state = stateCueText
				continue
			}
			continue

		case stateCueText:
			if strings.TrimSpace(line) == "" {
				flushCue()
				continue
			}
			textLines = append(textLines, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("vtt: scan failed: %w", err)
	}

	if state == stateHeader {
		return nil, fmt.Errorf("vtt: missing WEBVTT header")
	}
	if state == stateCueText {
		flushCue()
	}
	if state == stateCueStart && len(cues) == 0 && headerLineNum == 0 {
		return nil, fmt.Errorf("vtt: missing WEBVTT header")
	}

	return cues, nil
}

func parseVTTTimeRange(line string) (startMS int64, endMS int64, err error) {
	left, right, ok := strings.Cut(line, "-->")
	if !ok {
		return 0, 0, fmt.Errorf("invalid timestamp range %q (missing -->)", line)
	}
	left = strings.TrimSpace(left)
	right = strings.TrimSpace(right)
	if left == "" || right == "" {
		return 0, 0, fmt.Errorf("invalid timestamp range %q", line)
	}

	if fields := strings.Fields(left); len(fields) > 0 {
		left = fields[0]
	}
	if fields := strings.Fields(right); len(fields) > 0 {
		right = fields[0]
	}

	startMS, err = parseVTTTimestamp(left)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start timestamp %q: %v", left, err)
	}
	endMS, err = parseVTTTimestamp(right)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end timestamp %q: %v", right, err)
	}
	return startMS, endMS, nil
}

func parseVTTTimestamp(ts string) (int64, error) {
	parts := strings.Split(ts, ":")
	if len(parts) != 2 && len(parts) != 3 {
		return 0, fmt.Errorf("expected MM:SS.mmm or HH:MM:SS.mmm")
	}

	var hours int
	var minutesPart string
	var secondsPart string
	var err error

	switch len(parts) {
	case 2:
		minutesPart = parts[0]
		secondsPart = parts[1]
	case 3:
		hours, err = strconv.Atoi(parts[0])
		if err != nil || hours < 0 {
			return 0, fmt.Errorf("invalid hours")
		}
		minutesPart = parts[1]
		secondsPart = parts[2]
	}

	minutes, err := strconv.Atoi(minutesPart)
	if err != nil || minutes < 0 {
		return 0, fmt.Errorf("invalid minutes")
	}

	secParts := strings.Split(secondsPart, ".")
	if len(secParts) != 2 {
		return 0, fmt.Errorf("expected dot milliseconds")
	}
	seconds, err := strconv.Atoi(secParts[0])
	if err != nil || seconds < 0 || seconds > 59 {
		return 0, fmt.Errorf("invalid seconds")
	}
	ms, err := strconv.Atoi(secParts[1])
	if err != nil || ms < 0 || ms > 999 {
		return 0, fmt.Errorf("invalid milliseconds")
	}

	total := int64((((hours*60)+minutes)*60+seconds)*1000 + ms)
	return total, nil
}
