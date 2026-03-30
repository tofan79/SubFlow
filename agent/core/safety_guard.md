# Core: Safety Guard — Anti-Halu & Anti-Rusak

## Dipanggil Ketika
- Ada "blocked_by" di short_term.json
- Build gagal 3x berturut-turut dengan error berbeda
- Agent mau hapus > 50 baris kode existing
- Akan tambah dependency baru ke go.mod tanpa konfirmasi Claude

---

## Pemeriksaan Otomatis (sebelum setiap commit)

```bash
# 1. No credential leak
grep -rn "sk-\|Bearer \|api_key\s*=" internal/ frontend/src/

# 2. No subtitle exfiltration
grep -rn "subflow\.io\|subflow\.app" internal/

# 3. No accidental deletion
git diff --stat | grep "^-"
# Jika ada file yang dihapus yang tidak ada di task → STOP
```

---

## Respons Terhadap Anomali

| Severity | Kondisi | Aksi |
|---|---|---|
| Warning | Import package tidak di go.mod | Log, blokir commit sampai go get |
| Block | Build gagal 3x | Stop, set needs_human=true |
| Critical | API key ditemukan di plaintext | Stop, hapus file, alert |
| Critical | > 50 baris existing dihapus tanpa instruksi | Stop, git restore, alert |

---

## Hallucination Symptoms

Tanda agent sedang hallusinasi:
1. Import ke package yang tidak ada
2. Panggil method yang tidak ada di interface
3. Tulis kode tidak berhubungan dengan task
4. Override code yang sudah ada tanpa instruksi

Respons:
```
1. STOP segera
2. Baca ulang task brief di short_term.json
3. Baca ulang interface di types.go
4. Mulai ulang dari awal task (jangan lanjutkan kode yang salah)
5. Jika masih tidak yakin → needs_human: true
```

---

## Kondisi yang SELALU Butuh Manusia

- Akan ubah database schema (migration baru)
- Akan ubah Wails IPC method signature (breaking change)
- Build tetap gagal setelah 3x percobaan
- Keputusan arsitektur bertentangan dengan long_term.json
- Akan tambah package baru ke go.mod

Cara escalate:
```json
{
  "blocked_by": "safety_guard",
  "needs_human": true,
  "reason": "penjelasan spesifik",
  "question": "pertanyaan yang butuh dijawab manusia"
}
```
