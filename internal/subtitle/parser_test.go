package subtitle

import (
	"strings"
	"testing"
)

func TestSRT_ParseBasic(t *testing.T) {
	in := "" +
		"1\n" +
		"00:00:01,200 --> 00:00:03,500\n" +
		"Hello everyone\n" +
		"\n" +
		"2\n" +
		"00:00:03,600 --> 00:00:05,100\n" +
		"Welcome back\n" +
		"\n" +
		"3\n" +
		"00:01:00,000 --> 00:01:00,001\n" +
		"End\n"

	cues, err := ParseSRT(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseSRT error: %v", err)
	}
	if len(cues) != 3 {
		t.Fatalf("expected 3 cues, got %d", len(cues))
	}

	if cues[0].Index != 1 || cues[0].StartMS != 1200 || cues[0].EndMS != 3500 || cues[0].Text != "Hello everyone" {
		t.Fatalf("cue[0] mismatch: %+v", cues[0])
	}
	if cues[1].Index != 2 || cues[1].StartMS != 3600 || cues[1].EndMS != 5100 || cues[1].Text != "Welcome back" {
		t.Fatalf("cue[1] mismatch: %+v", cues[1])
	}
	if cues[2].Index != 3 || cues[2].StartMS != 60000 || cues[2].EndMS != 60001 || cues[2].Text != "End" {
		t.Fatalf("cue[2] mismatch: %+v", cues[2])
	}
}

func TestSRT_ParseEmptyLines(t *testing.T) {
	in := "" +
		"\n\n" +
		"1\r\n" +
		"00:00:00,000 --> 00:00:00,010\r\n" +
		"A\r\n" +
		"\r\n" +
		"\r\n" +
		"2\r\n" +
		"00:00:00,020 --> 00:00:00,030\r\n" +
		"B\r\n" +
		"\r\n\r\n"

	cues, err := ParseSRT(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseSRT error: %v", err)
	}
	if len(cues) != 2 {
		t.Fatalf("expected 2 cues, got %d", len(cues))
	}
	if cues[0].Text != "A" || cues[1].Text != "B" {
		t.Fatalf("text mismatch: %+v", cues)
	}
}

func TestSRT_ParseUTF8(t *testing.T) {
	in := "" +
		"1\n" +
		"00:00:10,000 --> 00:00:10,999\n" +
		"你好，世界 👋🌏\n" +
		"\n" +
		"2\n" +
		"00:00:11,000 --> 00:00:11,001\n" +
		"日本語テキスト🙂\n"

	cues, err := ParseSRT(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseSRT error: %v", err)
	}
	if len(cues) != 2 {
		t.Fatalf("expected 2 cues, got %d", len(cues))
	}
	if cues[0].Text != "你好，世界 👋🌏" {
		t.Fatalf("cue[0] text mismatch: %q", cues[0].Text)
	}
	if cues[1].Text != "日本語テキスト🙂" {
		t.Fatalf("cue[1] text mismatch: %q", cues[1].Text)
	}

	if cues[0].StartMS != 10000 || cues[0].EndMS != 10999 {
		t.Fatalf("timestamp mismatch: start=%d end=%d", cues[0].StartMS, cues[0].EndMS)
	}
}

func TestSRT_ParseMultiline(t *testing.T) {
	in := "" +
		"\uFEFF1\n" +
		"00:00:01,000 --> 00:00:02,000\n" +
		"Line 1\n" +
		"Line 2\n" +
		"Line 3\n" +
		"\n"

	cues, err := ParseSRT(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseSRT error: %v", err)
	}
	if len(cues) != 1 {
		t.Fatalf("expected 1 cue, got %d", len(cues))
	}
	if cues[0].Text != "Line 1\nLine 2\nLine 3" {
		t.Fatalf("multiline text mismatch: %q", cues[0].Text)
	}
}
