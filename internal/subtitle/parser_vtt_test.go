package subtitle

import (
	"strings"
	"testing"
)

func TestVTT_ParseBasic(t *testing.T) {
	input := "WEBVTT\n\n00:00:01.200 --> 00:00:03.500 align:middle\nHello everyone\n\n00:00:03.600 --> 00:00:05.100\nWelcome back\n"

	cues, err := ParseVTT(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseVTT() error = %v", err)
	}
	if len(cues) != 2 {
		t.Fatalf("ParseVTT() cues = %d, want 2", len(cues))
	}

	if cues[0].Index != 1 || cues[0].StartMS != 1200 || cues[0].EndMS != 3500 || cues[0].Text != "Hello everyone" {
		t.Fatalf("first cue = %+v, want index 1 start 1200 end 3500 text %q", cues[0], "Hello everyone")
	}
	if cues[1].Index != 2 || cues[1].StartMS != 3600 || cues[1].EndMS != 5100 || cues[1].Text != "Welcome back" {
		t.Fatalf("second cue = %+v, want index 2 start 3600 end 5100 text %q", cues[1], "Welcome back")
	}
}

func TestVTT_ParseWithCueID(t *testing.T) {
	input := "WEBVTT\n\nintro\n00:00:01.000 --> 00:00:02.000\nFirst line\n\nsecond-id\n00:00:02.500 --> 00:00:04.000\nSecond line\n"

	cues, err := ParseVTT(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseVTT() error = %v", err)
	}
	if len(cues) != 2 {
		t.Fatalf("ParseVTT() cues = %d, want 2", len(cues))
	}

	if cues[0].Index != 1 || cues[0].Text != "First line" {
		t.Fatalf("first cue = %+v", cues[0])
	}
	if cues[1].Index != 2 || cues[1].Text != "Second line" {
		t.Fatalf("second cue = %+v", cues[1])
	}
}

func TestVTT_DotTimestamp(t *testing.T) {
	input := "WEBVTT\n\n01:02.345 --> 01:03.678\nOne minute cue\n"

	cues, err := ParseVTT(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ParseVTT() error = %v", err)
	}
	if len(cues) != 1 {
		t.Fatalf("ParseVTT() cues = %d, want 1", len(cues))
	}
	if cues[0].StartMS != 62345 || cues[0].EndMS != 63678 {
		t.Fatalf("cue timing = %+v, want start 62345 end 63678", cues[0])
	}
	if cues[0].Text != "One minute cue" {
		t.Fatalf("cue text = %q, want %q", cues[0].Text, "One minute cue")
	}
}
