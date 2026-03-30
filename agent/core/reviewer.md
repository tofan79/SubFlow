# Core: Reviewer — Diagnosis & Code Review

## Dipanggil Ketika
1. Executor gagal build/test setelah 2x retry
2. Ada regresi (test sebelumnya pass, sekarang fail)
3. Phase selesai → perlu quality review sebelum lanjut

---

## Mode A: Diagnosis Build Failure

### Kategori Error Go
```
TYPE_MISMATCH    → "cannot use X as type Y"
IMPORT_MISSING   → "cannot find package"
UNDEFINED        → "undefined: FunctionName"
INTERFACE_BREACH → "X does not implement Y"
SYNTAX           → "syntax error: ..."
TEST_FAIL        → "FAIL: TestXxx"
CIRCULAR_IMPORT  → "import cycle not allowed"
```

### Kategori Error Svelte
```
TYPE_ERROR    → TypeScript type mismatch
IMPORT_ERROR  → Cannot find module
WAILS_MISSING → Wails binding not found
PROP_ERROR    → Component prop type mismatch
```

### Output Diagnosis
```json
{
  "reviewer": {
    "mode": "diagnose",
    "reviewed_by": "claude",
    "error_category": "INTERFACE_BREACH",
    "root_cause": "DeepLProvider tidak punya method EstimateCost()",
    "affected_files": ["internal/translation/deepl.go"],
    "fix_instruction": "Tambahkan: func (d *DeepLProvider) EstimateCost(charCount int) float64 { return float64(charCount) * 0.00002 }",
    "fix_complexity": "low",
    "send_to": "gpt"
  }
}
```

---

## Mode B: Quality Review (per Phase)

Go Checklist:
- [ ] Semua public function punya doc comment
- [ ] Tidak ada `return err` telanjang (harus wrap dengan context)
- [ ] Tidak ada hardcoded credential
- [ ] Tidak ada circular import (cek dengan `go build ./...`)
- [ ] Semua goroutine punya cancel path
- [ ] Interface max 5 methods

Svelte Checklist:
- [ ] Semua EventsOn punya EventsOff di onDestroy
- [ ] Semua async punya error handling
- [ ] Tidak ada hardcoded string (pakai konstanta atau i18n)
- [ ] Semua store typed (tidak ada `any`)

---

## Mode C: Regression Analysis

```bash
git log --oneline -10     # lihat commit terakhir
git diff HEAD~1 [file]    # lihat apa yang berubah
go test ./... 2>&1        # run semua test
```

Identifikasi commit yang breaking → buat fix plan → serahkan ke builder.

---

## Eskalasi ke Manusia

Jika setelah analisis masih tidak bisa tentukan fix:
```json
{
  "reviewer": {
    "escalated": true,
    "needs_human": true,
    "question": "Pertanyaan spesifik yang butuh jawaban manusia",
    "context": "konteks lengkap"
  }
}
```
Stop, jangan tebak-tebak.
