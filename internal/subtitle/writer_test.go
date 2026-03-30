package subtitle

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriter_SRT(t *testing.T) {
	cues := []Cue{
		{StartMS: 1200, EndMS: 3500, Text: "Hello everyone"},
		{StartMS: 3600, EndMS: 5100, Text: "Welcome back"},
	}

	var buf bytes.Buffer
	if err := WriteSRT(&buf, cues); err != nil {
		t.Fatalf("WriteSRT error: %v", err)
	}

	want := "1\n00:00:01,200 --> 00:00:03,500\nHello everyone\n\n" +
		"2\n00:00:03,600 --> 00:00:05,100\nWelcome back\n\n"
	if got := buf.String(); got != want {
		t.Fatalf("WriteSRT output mismatch:\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestWriter_VTT(t *testing.T) {
	cues := []Cue{{StartMS: 1200, EndMS: 3500, Text: "Hello everyone"}}

	var buf bytes.Buffer
	if err := WriteVTT(&buf, cues); err != nil {
		t.Fatalf("WriteVTT error: %v", err)
	}

	want := "WEBVTT\n\n00:00:01.200 --> 00:00:03.500\nHello everyone\n\n"
	if got := buf.String(); got != want {
		t.Fatalf("WriteVTT output mismatch:\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestWriter_DualSubtitle(t *testing.T) {
	cues := []DualCue{{
		Cue:         Cue{StartMS: 1000, EndMS: 2500, Text: "Hello"},
		Translation: "Halo",
	}}

	var buf bytes.Buffer
	if err := WriteDualSRT(&buf, cues); err != nil {
		t.Fatalf("WriteDualSRT error: %v", err)
	}

	want := "1\n00:00:01,000 --> 00:00:02,500\nHello\nHalo\n\n"
	if got := buf.String(); got != want {
		t.Fatalf("WriteDualSRT output mismatch:\nwant:\n%q\ngot:\n%q", want, got)
	}
}

func TestWriter_ASSAndTXT(t *testing.T) {
	cues := []Cue{{StartMS: 1234, EndMS: 5678, Text: "Line 1\nLine 2"}}

	var ass bytes.Buffer
	if err := WriteASS(&ass, cues); err != nil {
		t.Fatalf("WriteASS error: %v", err)
	}
	if !strings.Contains(ass.String(), "[Script Info]") || !strings.Contains(ass.String(), "Dialogue: 0,0:00:01.23,0:00:05.68,Default,,0,0,0,,Line 1\\NLine 2") {
		t.Fatalf("WriteASS output unexpected:\n%s", ass.String())
	}

	var txt bytes.Buffer
	if err := WriteTXT(&txt, cues); err != nil {
		t.Fatalf("WriteTXT error: %v", err)
	}
	if got := txt.String(); got != "[00:00:01] Line 1\nLine 2\n" {
		t.Fatalf("WriteTXT output mismatch:\n%q", got)
	}
}
