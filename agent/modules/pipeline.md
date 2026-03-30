# Module: Pipeline

## State Machine
```
imported → audio_extracted → asr_done → corrected → context_done
→ segmented → translated → rewritten → qa_done → exported
```

Aturan wajib:
- Pipeline hanya maju, tidak pernah mundur otomatis
- Setiap step tulis JSON output sebelum update DB
- Jika JSON sudah ada DAN DB state >= step ini → skip (resume safe)

## Package Structure (internal/pipeline/)
```
types.go        — semua shared types + interfaces
errors.go       — semua error codes ERR_*
orchestrator.go — step runner + state machine + resume logic
retry.go        — exponential backoff + fallback chain
```

## Interface Utama
```go
type TranslationProvider interface {
    Translate(ctx context.Context, batch []string, opts TranslationOpts) ([]string, error)
    EstimateCost(charCount int) float64
    MaxBatchSize() int
    Name() string
}

type RewriteProvider interface {
    Rewrite(ctx context.Context, batch []RewriteInput, opts RewriteOpts) ([]string, error)
    EstimateCost(tokenCount int) float64
    Name() string
}

type ASRProvider interface {
    Transcribe(ctx context.Context, audioPath string, opts ASROpts) (<-chan ASRSegment, <-chan error)
    EstimateCost(durationSeconds float64) float64
    Name() string
}
```

## Retry Strategy
```
Attempt 1 → fail → tunggu 1 detik
Attempt 2 → fail → tunggu 2 detik
Attempt 3 → fail → tunggu 4 detik
Semua habis → fallback ke provider berikutnya
Semua provider habis → ERR_TRN_004
```

Retryable: HTTP 429, HTTP 5xx, timeout
Non-retryable: HTTP 4xx, invalid API key

## Wails Events yang Dipancarkan
```
pipeline:progress → {step, current, total, percent}
pipeline:complete → {projectId, segmentCount}
pipeline:error    → {code, step, message, retryable}
cost:estimate     → {layer1, layer2, total}
asr:hardware      → {backend, gpuName, cudaVersion}
file:dropped      → {paths}
```

## State Files per Project
```
{AppData}/SubFlow/projects/{id}/
  raw_asr.json    corrected.json   context.json
  segmented.json  translated.json  rewritten.json
  glossary.json   qa_report.json
```
