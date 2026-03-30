# Module: Translation & Rewrite

## Dua Layer (keduanya selalu dijalankan)
```
Layer 1 — Translation   → AKURASI (makna tepat)
Layer 2 — AI Rewrite    → NATURALNESS (terdengar manusia)
```

## Layer 1 Providers
| Provider | File | Batch | Concurrency |
|---|---|---|---|
| DeepL | deepl.go | 50 | 10 |
| OpenAI | openai.go | 20 | 5 |
| Anthropic | anthropic.go | 10 | 3 |
| Gemini | gemini.go | 20 | 5 |
| Ollama | ollama.go | 1 | 1 |

Fallback chain (default): DeepL → OpenAI → Anthropic → Ollama

## Layer 2 Providers
OpenAI, OpenRouter, Anthropic, Gemini, Qwen, xAI, Ollama

## Tone Presets
| Preset | Karakteristik |
|---|---|
| natural | Seimbang, percakapan nyata |
| formal | Baku, kalimat lengkap, EYD |
| casual | Santai, boleh gue/lo/nih |
| cinematic | Dramatis, diksi kuat, emosi ditonjolkan |

## System Prompt Layer 1
```
Terjemahkan subtitle {source_lang} berikut ke Bahasa Indonesia.
Tipe konten: {content_mode}
- Pertahankan makna asli secara akurat
- Satu baris input = satu baris output
- Pertahankan formatting marker (<i>, <b>)
{glossary_section}
```

## System Prompt Layer 2
```
Perbaiki terjemahan ini agar natural dalam Bahasa Indonesia.
Teks sumber: {source_text}
Terjemahan mesin: {translated_text}
Konteks: speaker={speaker}, emosi={emotion}, tone={tone}
Gaya: {tone_preset}
Batas: max {max_chars} char/baris, max {max_lines} baris, max {max_cps} CPS
{glossary_section}
Output HANYA teks terjemahan, tanpa penjelasan.
```

## Concurrency Pattern (WAJIB)
```go
sem := make(chan struct{}, maxConcurrent)
// Setiap goroutine: sem <- struct{}{}; defer func() { <-sem }()
```

## Cost Estimation
| Provider | Rate |
|---|---|
| DeepL Free | $0 (up to 500k char/bulan) |
| DeepL Pro | $0.00002/char |
| GPT-4o | $2.50/1M input tokens, $10/1M output |
| Claude Haiku | $0.80/1M input, $4/1M output |
| Groq Whisper | $0.02/menit audio |
| Ollama | $0 (lokal) |
