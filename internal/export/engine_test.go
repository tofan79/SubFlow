package export

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatExtension(t *testing.T) {
	tests := []struct {
		format Format
		ext    string
	}{
		{FormatSRT, ".srt"},
		{FormatVTT, ".vtt"},
		{FormatASS, ".ass"},
		{FormatTXT, ".txt"},
		{FormatJSON, ".json"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			got := tt.format.Extension()
			if got != tt.ext {
				t.Errorf("Extension() = %q, want %q", got, tt.ext)
			}
		})
	}
}

func TestIsValidFormat(t *testing.T) {
	if !IsValidFormat(FormatSRT) {
		t.Error("FormatSRT should be valid")
	}
	if IsValidFormat(Format("invalid")) {
		t.Error("invalid format should not be valid")
	}
}

func TestCardText(t *testing.T) {
	card := Card{
		Source: "Hello",
		L1:     "Halo",
		L2:     "Hai",
	}

	tests := []struct {
		layer Layer
		want  string
	}{
		{LayerSource, "Hello"},
		{LayerL1, "Halo"},
		{LayerL2, "Hai"},
		{Layer("unknown"), "Hai"},
	}

	for _, tt := range tests {
		t.Run(string(tt.layer), func(t *testing.T) {
			got := card.Text(tt.layer)
			if got != tt.want {
				t.Errorf("Text(%s) = %q, want %q", tt.layer, got, tt.want)
			}
		})
	}
}

func TestWriteSRT(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "Hello world"},
		{Index: 2, StartMS: 4000, EndMS: 6500, L2: "Goodbye world"},
	}

	var buf bytes.Buffer
	err := e.writeSRT(&buf, cards, LayerL2)
	if err != nil {
		t.Fatalf("writeSRT failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "1\n00:00:01,000 --> 00:00:03,000\nHello world") {
		t.Errorf("SRT output incorrect: %s", output)
	}
	if !strings.Contains(output, "2\n00:00:04,000 --> 00:00:06,500\nGoodbye world") {
		t.Errorf("SRT output incorrect: %s", output)
	}
}

func TestWriteVTT(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "Hello"},
	}

	var buf bytes.Buffer
	err := e.writeVTT(&buf, cards, LayerL2)
	if err != nil {
		t.Fatalf("writeVTT failed: %v", err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "WEBVTT\n\n") {
		t.Error("VTT should start with WEBVTT header")
	}
	if !strings.Contains(output, "00:00:01.000 --> 00:00:03.000") {
		t.Errorf("VTT timestamp format incorrect: %s", output)
	}
}

func TestWriteASS(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "Hello"},
	}

	var buf bytes.Buffer
	err := e.writeASS(&buf, cards, LayerL2)
	if err != nil {
		t.Fatalf("writeASS failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[Script Info]") {
		t.Error("ASS should contain [Script Info] section")
	}
	if !strings.Contains(output, "[V4+ Styles]") {
		t.Error("ASS should contain [V4+ Styles] section")
	}
	if !strings.Contains(output, "Dialogue: 0,0:00:01.00,0:00:03.00,Default,,0,0,0,,Hello") {
		t.Errorf("ASS dialogue format incorrect: %s", output)
	}
}

func TestWriteDualSRT(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, Source: "Hello", L2: "Halo"},
	}

	var buf bytes.Buffer
	err := e.writeDualSRT(&buf, cards)
	if err != nil {
		t.Fatalf("writeDualSRT failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Hello\nHalo") {
		t.Errorf("Dual SRT should have both source and L2: %s", output)
	}
}

func TestWriteDualASS(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, Source: "Hello", L2: "Halo"},
	}

	var buf bytes.Buffer
	err := e.writeDualASS(&buf, cards)
	if err != nil {
		t.Fatalf("writeDualASS failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Style: Source") {
		t.Error("Dual ASS should have Source style")
	}
	if !strings.Contains(output, "Style: Translation") {
		t.Error("Dual ASS should have Translation style")
	}
	if !strings.Contains(output, ",Source,") {
		t.Error("Dual ASS should have dialogue with Source style")
	}
	if !strings.Contains(output, ",Translation,") {
		t.Error("Dual ASS should have dialogue with Translation style")
	}
}

func TestWriteTXT(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 3661000, EndMS: 3663000, L2: "One hour one minute"},
	}

	var buf bytes.Buffer
	err := e.writeTXT(&buf, cards, LayerL2)
	if err != nil {
		t.Fatalf("writeTXT failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "[01:01:01] One hour one minute") {
		t.Errorf("TXT format incorrect: %s", output)
	}
}

func TestWriteJSON(t *testing.T) {
	e := NewExporter("test-project")
	cards := []Card{
		{
			Index:    1,
			StartMS:  1000,
			EndMS:    3000,
			Source:   "Hello",
			L1:       "Halo",
			L2:       "Hai",
			QAStatus: QAStatusPass,
		},
	}

	var buf bytes.Buffer
	err := e.writeJSON(&buf, cards)
	if err != nil {
		t.Fatalf("writeJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.HasPrefix(output, "[\n") {
		t.Error("JSON should start with array")
	}
	if !strings.HasSuffix(output, "]\n") {
		t.Error("JSON should end with array")
	}
	if !strings.Contains(output, `"source": "Hello"`) {
		t.Errorf("JSON should contain source: %s", output)
	}
	if !strings.Contains(output, `"qa_status": "pass"`) {
		t.Errorf("JSON should contain qa_status: %s", output)
	}
}

func TestExport_Integration(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, Source: "Hello", L1: "Halo", L2: "Hai", QAStatus: QAStatusPass},
		{Index: 2, StartMS: 4000, EndMS: 6000, Source: "World", L1: "Dunia", L2: "Dunia", QAStatus: QAStatusPass},
	}

	opts := Options{
		Format:    FormatSRT,
		Layer:     LayerL2,
		OutputDir: tmpDir,
	}

	result, err := e.Export(cards, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Export should succeed: %s", result.Error)
	}
	if result.CardsWritten != 2 {
		t.Errorf("CardsWritten = %d, want 2", result.CardsWritten)
	}
	if result.Format != FormatSRT {
		t.Errorf("Format = %s, want srt", result.Format)
	}

	content, err := os.ReadFile(result.OutputPath)
	if err != nil {
		t.Fatalf("Cannot read output file: %v", err)
	}

	if !strings.Contains(string(content), "Hai") {
		t.Error("Output should contain L2 text")
	}
}

func TestExport_QAErrors(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "OK", QAStatus: QAStatusPass},
		{Index: 2, StartMS: 4000, EndMS: 6000, L2: "Error", QAStatus: QAStatusError},
	}

	opts := Options{
		Format:      FormatSRT,
		Layer:       LayerL2,
		OutputDir:   tmpDir,
		AllowErrors: false,
	}

	result, err := e.Export(cards, opts)
	if err == nil {
		t.Error("Export should fail with QA errors when AllowErrors=false")
	}
	if result.Success {
		t.Error("Result should not be success")
	}
	if result.QAErrors != 1 {
		t.Errorf("QAErrors = %d, want 1", result.QAErrors)
	}
}

func TestExport_AllowErrors(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "OK", QAStatus: QAStatusPass},
		{Index: 2, StartMS: 4000, EndMS: 6000, L2: "Error", QAStatus: QAStatusError},
	}

	opts := Options{
		Format:      FormatSRT,
		Layer:       LayerL2,
		OutputDir:   tmpDir,
		AllowErrors: true,
	}

	result, err := e.Export(cards, opts)
	if err != nil {
		t.Errorf("Export should succeed with AllowErrors=true: %v", err)
	}
	if !result.Success {
		t.Error("Result should be success")
	}
}

func TestExport_FilterEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "Hello", QAStatus: QAStatusPass},
		{Index: 2, StartMS: 4000, EndMS: 6000, L2: "", QAStatus: QAStatusPass},
		{Index: 3, StartMS: 7000, EndMS: 9000, L2: "   ", QAStatus: QAStatusPass},
		{Index: 4, StartMS: 10000, EndMS: 12000, L2: "World", QAStatus: QAStatusPass},
	}

	opts := Options{
		Format:       FormatSRT,
		Layer:        LayerL2,
		OutputDir:    tmpDir,
		IncludeEmpty: false,
	}

	result, err := e.Export(cards, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if result.CardsWritten != 2 {
		t.Errorf("CardsWritten = %d, want 2 (filtered empty)", result.CardsWritten)
	}
	if result.SkippedEmpty != 2 {
		t.Errorf("SkippedEmpty = %d, want 2", result.SkippedEmpty)
	}
}

func TestExport_DualLayer(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, Source: "Hello", L2: "Halo", QAStatus: QAStatusPass},
	}

	opts := Options{
		Format:    FormatSRT,
		Layer:     LayerDual,
		OutputDir: tmpDir,
	}

	result, err := e.Export(cards, opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	content, err := os.ReadFile(result.OutputPath)
	if err != nil {
		t.Fatalf("Cannot read output: %v", err)
	}

	if !strings.Contains(string(content), "Hello\nHalo") {
		t.Errorf("Dual export should contain both texts: %s", string(content))
	}
}

func TestBatchExport(t *testing.T) {
	tmpDir := t.TempDir()
	e := NewExporter("test-project")
	cards := []Card{
		{Index: 1, StartMS: 1000, EndMS: 3000, L2: "Hello", QAStatus: QAStatusPass},
	}

	opts := Options{
		Layer:     LayerL2,
		OutputDir: tmpDir,
	}

	formats := []Format{FormatSRT, FormatVTT, FormatASS}
	results, err := e.BatchExport(cards, formats, opts)
	if err != nil {
		t.Fatalf("BatchExport failed: %v", err)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for _, r := range results {
		if !r.Success {
			t.Errorf("Export %s failed: %s", r.Format, r.Error)
		}
	}

	files, _ := os.ReadDir(tmpDir)
	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

func TestValidateBeforeExport(t *testing.T) {
	cards := []Card{
		{Index: 1, QAStatus: QAStatusPass},
		{Index: 2, QAStatus: QAStatusWarn},
		{Index: 3, QAStatus: QAStatusError},
		{Index: 4, QAStatus: QAStatusPending},
	}

	opts := Options{
		Format:    FormatSRT,
		OutputDir: "/tmp",
	}

	canExport, warnings, errors := ValidateBeforeExport(cards, opts)

	if canExport {
		t.Error("Should not be able to export with errors")
	}
	if len(warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(warnings))
	}
	if len(errors) != 2 {
		t.Errorf("Expected 2 errors, got %d: %v", len(errors), errors)
	}
}

func TestTimestampFormatters(t *testing.T) {
	ms := int64(3723456)

	srt := formatSRTTimestamp(ms)
	if srt != "01:02:03,456" {
		t.Errorf("SRT timestamp = %s, want 01:02:03,456", srt)
	}

	vtt := formatVTTTimestamp(ms)
	if vtt != "01:02:03.456" {
		t.Errorf("VTT timestamp = %s, want 01:02:03.456", vtt)
	}

	ass := formatASSTimestamp(ms)
	if ass != "1:02:03.46" {
		t.Errorf("ASS timestamp = %s, want 1:02:03.46", ass)
	}

	txt := formatTXTTimestamp(ms)
	if txt != "01:02:03" {
		t.Errorf("TXT timestamp = %s, want 01:02:03", txt)
	}
}

func TestTimestampNegative(t *testing.T) {
	ms := int64(-1000)

	srt := formatSRTTimestamp(ms)
	if srt != "00:00:00,000" {
		t.Errorf("Negative timestamp should be 00:00:00,000, got %s", srt)
	}
}

func TestEscapeASSText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello\nWorld", "Hello\\NWorld"},
		{"Test\\slash", "Test\\\\slash"},
		{"{tag}", "\\{tag\\}"},
		{"Line1\r\nLine2", "Line1\\NLine2"},
	}

	for _, tt := range tests {
		got := escapeASSText(tt.input)
		if got != tt.want {
			t.Errorf("escapeASSText(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestGenerateOutputPath(t *testing.T) {
	e := NewExporter("project123")

	opts := Options{
		Format:         FormatSRT,
		Layer:          LayerL2,
		OutputDir:      "/tmp/exports",
		OutputFilename: "custom_name",
	}

	path := e.generateOutputPath(opts)
	if !strings.HasSuffix(path, ".srt") {
		t.Errorf("Path should end with .srt: %s", path)
	}
	if !strings.Contains(path, "custom_name") {
		t.Errorf("Path should contain custom name: %s", path)
	}
	if filepath.Dir(path) != "/tmp/exports" {
		t.Errorf("Path should be in /tmp/exports: %s", path)
	}
}

func TestExport_InvalidFormat(t *testing.T) {
	e := NewExporter("test")
	opts := Options{
		Format:    Format("invalid"),
		OutputDir: "/tmp",
	}

	_, err := e.Export([]Card{}, opts)
	if err == nil {
		t.Error("Should fail with invalid format")
	}
}

func TestExport_EmptyOutputDir(t *testing.T) {
	e := NewExporter("test")
	opts := Options{
		Format:    FormatSRT,
		OutputDir: "",
	}

	_, err := e.Export([]Card{}, opts)
	if err == nil {
		t.Error("Should fail with empty output dir")
	}
}
