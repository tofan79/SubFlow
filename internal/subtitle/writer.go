package subtitle

import (
	"fmt"
	"io"
	"strings"
)

type DualCue struct {
	Cue
	Translation string
}

func WriteSRT(w io.Writer, cues []Cue) error {
	for i, cue := range cues {
		if _, err := fmt.Fprintf(w, "%d\n%s --> %s\n%s\n\n",
			i+1,
			formatSRTTimestamp(cue.StartMS),
			formatSRTTimestamp(cue.EndMS),
			cue.Text,
		); err != nil {
			return err
		}
	}
	return nil
}

func WriteVTT(w io.Writer, cues []Cue) error {
	if _, err := io.WriteString(w, "WEBVTT\n\n"); err != nil {
		return err
	}
	for _, cue := range cues {
		if _, err := fmt.Fprintf(w, "%s --> %s\n%s\n\n",
			formatVTTTimestamp(cue.StartMS),
			formatVTTTimestamp(cue.EndMS),
			cue.Text,
		); err != nil {
			return err
		}
	}
	return nil
}

func WriteASS(w io.Writer, cues []Cue) error {
	if _, err := io.WriteString(w, assHeader); err != nil {
		return err
	}
	for _, cue := range cues {
		if _, err := fmt.Fprintf(w, "Dialogue: 0,%s,%s,Default,,0,0,0,,%s\n",
			formatASSTimestamp(cue.StartMS),
			formatASSTimestamp(cue.EndMS),
			escapeASSText(cue.Text),
		); err != nil {
			return err
		}
	}
	return nil
}

func WriteTXT(w io.Writer, cues []Cue) error {
	for _, cue := range cues {
		if _, err := fmt.Fprintf(w, "[%s] %s\n", formatTXTTimestamp(cue.StartMS), cue.Text); err != nil {
			return err
		}
	}
	return nil
}

func WriteDualSRT(w io.Writer, cues []DualCue) error {
	for i, cue := range cues {
		text := cue.Text
		if cue.Translation != "" {
			text += "\n" + cue.Translation
		}
		if _, err := fmt.Fprintf(w, "%d\n%s --> %s\n%s\n\n",
			i+1,
			formatSRTTimestamp(cue.StartMS),
			formatSRTTimestamp(cue.EndMS),
			text,
		); err != nil {
			return err
		}
	}
	return nil
}

const assHeader = "[Script Info]\n" +
	"ScriptType: v4.00+\n" +
	"PlayResX: 384\n" +
	"PlayResY: 288\n" +
	"\n" +
	"[V4+ Styles]\n" +
	"Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding\n" +
	"Style: Default,Arial,20,&H00FFFFFF,&H000000FF,&H00000000,&H64000000,0,0,0,0,100,100,0,0,1,1,0,2,10,10,10,1\n" +
	"\n" +
	"[Events]\n" +
	"Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\n"

func formatSRTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, ms)
}

func formatVTTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	ms %= 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, ms)
}

func formatASSTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	totalCS := int((ms + 5) / 10)
	hours := totalCS / 360000
	totalCS %= 360000
	minutes := totalCS / 6000
	totalCS %= 6000
	seconds := totalCS / 100
	centiseconds := totalCS % 100
	return fmt.Sprintf("%d:%02d:%02d.%02d", hours, minutes, seconds, centiseconds)
}

func formatTXTTimestamp(ms int64) string {
	if ms < 0 {
		ms = 0
	}
	hours := ms / 3600000
	ms %= 3600000
	minutes := ms / 60000
	ms %= 60000
	seconds := ms / 1000
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func escapeASSText(text string) string {
	text = strings.ReplaceAll(text, "\\", "\\\\")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	text = strings.ReplaceAll(text, "\n", "\\N")
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	return text
}
