# Core: Router — Entry Point Setiap Sesi

## Fungsi
Router adalah langkah pertama setelah baca memory.
Memutuskan ke mana agent harus pergi berdasarkan state saat ini.

---

## Decision Tree

```
BACA short_term.json
        │
        ├── Ada "blocked_by" ? ──────────────────► safety_guard.md
        │
        ├── current_task.status = "in_progress" ► executor.md (RESUME)
        │
        ├── current_task.status = "failed"      ► reviewer.md (DIAGNOSE)
        │       → needs_reviewer = true
        │
        ├── current_task.status = "done"         ► planner.md (NEXT TASK)
        │
        └── current_task.status = "none"         ► planner.md (AMBIL TASK)
```

---

## Multi-Agent Coordination

Sebelum routing, cek siapa agent yang aktif:
```
Baca short_term.json → field "active_agent"
Jika "active_agent" != namaku → WAIT, jangan ambil task
Jika "active_agent" = null   → ambil sesuai role di task_list
```

Contoh: GPT sedang kerjakan P2-04, Gemini mau ambil P6-01.
Gemini cek: P6-01 depends_on P1-03 (done) ✓ → boleh lanjut.
Keduanya bisa jalan paralel karena tidak ada konflik file.

---

## Output Router

Tulis ke short_term.json sebelum serahkan ke modul berikutnya:
```json
{
  "routed_to": "executor | planner | reviewer | safety_guard",
  "route_reason": "penjelasan singkat",
  "active_agent": "claude | gpt | gemini",
  "session_start": "ISO timestamp"
}
```
