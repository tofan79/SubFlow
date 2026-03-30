package subtitle

import (
	"strings"
	"testing"
)

func TestSegmenter_SplitLongLine(t *testing.T) {
	line := "This subtitle line should split cleanly at a space boundary for wrapping"
	parts := SplitLongLine(line, 42)
	if len(parts) < 2 {
		t.Fatalf("expected split line, got %d parts: %#v", len(parts), parts)
	}
	for _, part := range parts {
		if CountChars(part) > 42 {
			t.Fatalf("part exceeds limit: %q", part)
		}
	}
	joined := strings.Join(parts, " ")
	if strings.ReplaceAll(joined, " ", "") != strings.ReplaceAll(line, " ", "") {
		t.Fatalf("wrapped text lost content: %q -> %#v", line, parts)
	}
}

func TestSegmenter_RespectWord(t *testing.T) {
	line := "supercalifragilisticexpialidocious"
	parts := SplitLongLine(line, 10)
	if len(parts) != 1 {
		t.Fatalf("expected no mid-word split, got %#v", parts)
	}
	if parts[0] != line {
		t.Fatalf("word was altered: %#v", parts)
	}
}

func TestSegmenter_CJKCount(t *testing.T) {
	if got := CountChars("你あアA"); got != 7 {
		t.Fatalf("expected CJK width count 7, got %d", got)
	}
}

func TestSegmenter_MaxLines(t *testing.T) {
	cues := []Cue{{
		Index:   1,
		StartMS: 0,
		EndMS:   6000,
		Text:    "Line 1\nLine 2\nLine 3",
	}}
	out := Segment(cues, SegmentConfig{MaxLines: 2})
	if len(out) != 2 {
		t.Fatalf("expected 2 cues, got %d: %#v", len(out), out)
	}
	if out[0].Text != "Line 1\nLine 2" {
		t.Fatalf("first cue mismatch: %#v", out[0])
	}
	if out[1].Text != "Line 3" {
		t.Fatalf("second cue mismatch: %#v", out[1])
	}
}
