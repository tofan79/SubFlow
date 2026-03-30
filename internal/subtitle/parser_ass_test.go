package subtitle

import (
	"strings"
	"testing"
)

func TestASS_ParseDialogue(t *testing.T) {
	in := "[Script Info]\n" +
		"Title: Example\n\n" +
		"[Events]\n" +
		"Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n" +
		"Comment: 0,0:00:00.00,0:00:01.00,Default,,0,0,0,,Ignored\n" +
		"Dialogue: 0,0:00:01.20,0:00:03.50,Default,,0,0,0,,Hello everyone\n" +
		"Dialogue: 1,0:00:04.00,0:00:05.25,Alt,Name,10,20,30,Effect,{\\an8}{\\i1}Hello, world\\NSecond line{\\b1}\n"

	cues, err := ParseASS(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseASS error: %v", err)
	}
	if len(cues) != 2 {
		t.Fatalf("expected 2 cues, got %d", len(cues))
	}

	if cues[0].Index != 0 || cues[0].StartMS != 1200 || cues[0].EndMS != 3500 || cues[0].Text != "Hello everyone" {
		t.Fatalf("cue[0] mismatch: %+v", cues[0])
	}
	if cues[1].Index != 0 || cues[1].StartMS != 4000 || cues[1].EndMS != 5250 || cues[1].Text != "Hello, world\nSecond line" {
		t.Fatalf("cue[1] mismatch: %+v", cues[1])
	}
}

func TestASS_StripTags(t *testing.T) {
	in := "{\\an8}{\\i1}{\\b1}Hello{\\pos(100,200)} world{\\blur2}\\Nnext{\\fnArial}{\\fs24}\\nline"
	want := "Hello world\nnext\nline"

	got := StripASSTags(in)
	got = strings.NewReplacer("\\N", "\n", "\\n", "\n").Replace(got)

	if got != want {
		t.Fatalf("StripASSTags() = %q, want %q", got, want)
	}
}

func TestASS_CentisecondTimestamp(t *testing.T) {
	in := "Dialogue: 0,0:00:01.23,0:00:04.56,Default,,0,0,0,,Timing check\n"

	cues, err := ParseASS(strings.NewReader(in))
	if err != nil {
		t.Fatalf("ParseASS error: %v", err)
	}
	if len(cues) != 1 {
		t.Fatalf("expected 1 cue, got %d", len(cues))
	}
	if cues[0].StartMS != 1230 || cues[0].EndMS != 4560 {
		t.Fatalf("timestamp mismatch: %+v", cues[0])
	}
}
