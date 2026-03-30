# Core: Planner — Pilih dan Siapkan Task Berikutnya

## Fungsi
Dipanggil ketika tidak ada task aktif.
Pilih task yang tepat dari task_list.json, siapkan brief untuk executor.

---

## Algoritma Pemilihan Task

```
1. Filter task_list.json:
   - status = "pending"
   - assigned_to = [nama agent ini]
   - semua depends_on sudah "done"

2. Prioritas:
   - Phase terkecil dulu (P1 sebelum P2, dst.)
   - Dalam satu phase: depends_on paling sedikit dulu
   - Estimasi lines terkecil dulu (quick wins)

3. Jika tidak ada task yang siap:
   - Cek apakah ada task blocked karena agent lain belum selesai
   - Tulis waiting_for ke short_term.json
   - Stop sampai dipanggil lagi
```

---

## Pecah Task Besar (> 250 baris estimasi)

```
Task besar → pecah jadi sub-tasks sebelum dikerjakan

Contoh:
  P7-04: "Build Timeline" (400 baris estimasi)
  Pecah jadi:
    P7-04a: Timeline container + card rendering (150 baris)
    P7-04b: Drag handles + timestamp update (150 baris)
    P7-04c: Snap + overlap detection (100 baris)

Tambahkan sub-tasks ke task_list.json
Tandai task asli sebagai "split"
```

---

## Task Brief

Sebelum serahkan ke executor, isi short_term.json:
```json
{
  "current_task": {
    "id": "P2-01",
    "title": "Write subtitle/parser.go",
    "assigned_to": "gpt",
    "status": "planned",
    "files_to_create": ["internal/subtitle/parser.go", "internal/subtitle/parser_test.go"],
    "files_to_read": ["agent/modules/pipeline.md"],
    "acceptance_criteria": [
      "go build ./... sukses",
      "go vet ./... bersih",
      "TestParseSRT_Basic pass",
      "ParseSRT mengembalikan slice Segment dengan StartMS EndMS RawText benar"
    ],
    "estimated_lines": 120,
    "module_ref": "agent/modules/pipeline.md"
  }
}
```
