# Core: QA Loop — Verifikasi Sebelum Done

## 4 Level Verifikasi (wajib semua pass)

```
Level 1 — Build
  go build ./...          → 0 error (Go)
  npm run check           → 0 error (Svelte, jika ada perubahan)
  Timeout: 60 detik

Level 2 — Test
  go test ./[package]/... → semua pass
  Timeout: 30 detik per package

Level 3 — Standards (non-blocking kecuali CRITICAL)
  WARNING : public function tanpa doc comment
  WARNING : fmt.Println tersisa (pakai logger)
  CRITICAL: API key plaintext → BLOCK
  CRITICAL: panic() di luar main.go → BLOCK

Level 4 — Regression
  go test ./... → bandingkan dengan baseline di long_term.json
  Jika ada test yang sebelumnya pass sekarang fail → kirim ke reviewer
```

---

## Failure Response

```
Level 1 fail:
  → Coba fix sendiri (max 2x)
  → Masih fail → set needs_reviewer=true, serahkan ke Claude

Level 2 fail:
  → Coba fix sendiri (max 2x)
  → Masih fail → set needs_reviewer=true

Level 3 CRITICAL:
  → Blokir commit
  → Fix sebelum lanjut

Level 4 regression:
  → Langsung serahkan ke Claude (reviewer)
  → Jangan coba fix sendiri
```

---

## Output QA Loop

```json
{
  "qa_loop": {
    "task_id": "P2-01",
    "levels_passed": [1, 2, 3, 4],
    "flags": [],
    "overall": "pass",
    "approved_for_commit": true
  }
}
```
