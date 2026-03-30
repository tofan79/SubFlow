# Core: Memory Manager

## Dua Jenis Memory

| | short_term.json | long_term.json |
|---|---|---|
| Reset | Setiap sesi baru (setelah flush) | Tidak pernah |
| Isi | State sesi ini | Keputusan permanen |
| Ditulis | Semua agent | Memory manager saat flush |

---

## Flush Procedure (akhir sesi / sebelum rate limit)

```
1. Baca short_term.json
2. Jika task done → tambah ke long_term.completed_tasks[]
3. Jika ada keputusan arsitektur baru → tambah ke long_term.architecture_decisions[]
4. Update long_term.phase_progress
5. Simpan long_term.json
6. Reset short_term.json (KECUALI jika task masih in_progress)
```

---

## Resume Procedure (sesi baru setelah disconnect)

```
1. Baca short_term.json
2. Cek current_task.status:
   "in_progress" → cari AGENT_RESUME comment → lanjut
   "failed"      → serahkan ke reviewer.md
   "done"        → flush ke long_term → ambil task baru
3. Gabungkan dengan long_term.json untuk full context
```

---

## short_term.json Schema

```json
{
  "session_id": null,
  "session_start": null,
  "active_agent": null,
  "routed_to": null,
  "route_reason": null,
  "waiting_for": null,

  "current_task": {
    "id": null,
    "title": null,
    "assigned_to": null,
    "status": "none",
    "started_at": null,
    "completed_at": null,
    "files_to_create": [],
    "files_to_read": [],
    "acceptance_criteria": [],
    "estimated_lines": 0,
    "build_verified": false,
    "test_verified": false,
    "attempts": 0,
    "error": null,
    "needs_reviewer": false,
    "resume_hint": null,
    "last_checkpoint": null
  },

  "reviewer": null,
  "safety_warnings": [],
  "errors": [],
  "files_modified": []
}
```

---

## long_term.json Schema

```json
{
  "project": "SubFlow",
  "last_updated": null,
  "completed_tasks": [],
  "architecture_decisions": [],
  "known_patterns": [],
  "known_issues": [],
  "phase_progress": {
    "P1": "pending", "P2": "pending", "P3": "pending",
    "P4": "pending", "P5": "pending", "P6": "pending",
    "P7": "pending", "P8": "pending", "P9": "pending", "P10": "pending"
  }
}
```
