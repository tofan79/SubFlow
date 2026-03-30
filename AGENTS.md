# SubFlow — Multi-Agent Build System

> Entry point untuk OpenCode + Oh My Open Agent.
> Tiga agent bekerja paralel dengan tugas berbeda.
> Baca file ini dulu sebelum apapun.

---

## Startup Sequence (WAJIB, setiap sesi)

```
1. Baca agent/rules.txt
2. Baca agent/context.txt
3. Baca agent/memory/long_term.json
4. Baca agent/memory/short_term.json
5. Baca agent/tasks/task_list.json
6. Tentukan siapa kamu (Claude/GPT/Gemini) → baca agent/agents/[namamu].md
7. Jalankan agent/core/router.md
```

---

## Tiga Agent, Tiga Peran

| Agent | File | Tanggung Jawab |
|---|---|---|
| **Claude** | `agent/agents/claude.md` | Arsitek + Reviewer — PRD logic, code review, keputusan arsitektur, QA diagnosis |
| **GPT** | `agent/agents/gpt.md` | Builder — semua Go backend, Wails IPC, SQLite, pipeline, ASR |
| **Gemini** | `agent/agents/gemini.md` | Frontend + Integration — Svelte UI, Flowbite, testing, packaging |

---

## Stack

```
Desktop Shell : Wails v2
Backend       : Go latest
Frontend      : Svelte + SvelteKit + Tailwind + Flowbite Svelte latest
Database      : SQLite (go-sqlite3, CGO)
ASR Lokal     : faster-whisper subprocess (CPU/CUDA/ROCm/CoreML/OpenVINO)
OS Target     : Windows 10+ 64-bit, Linux
```

---

## Anti-Halu Checklist (cek sebelum mulai coding)

- [ ] Sudah baca `short_term.json` → tahu posisi sekarang
- [ ] Sudah baca `long_term.json` → tahu keputusan arsitektur
- [ ] Task yang mau dikerjakan ada di `task_list.json` dengan status `pending`
- [ ] Semua dependencies task sudah `done`
- [ ] Tidak akan import package yang belum ada di `go.mod`
- [ ] Tidak akan menulis placeholder `// TODO`
