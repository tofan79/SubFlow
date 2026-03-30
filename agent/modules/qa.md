# Module: QA Engine

## Prinsip
QA Engine adalah **pure logic** — tidak ada I/O, tidak ada API call.
Input: []SubtitleCard → Output: []QAResult + []SubtitleCard (sudah dikoreksi)
Package `qa` DILARANG import package internal lain.

## 9 QA Checks

| ID | Check | Aturan | Auto-Fix |
|---|---|---|---|
| QA-01 | line_length | Maks 42 char/baris | ✅ Reflow teks |
| QA-02 | line_count | Maks 2 baris/card | ✅ Split ke card baru |
| QA-03 | duration_min | Min 1.0 detik | ✅ Perpanjang end timestamp |
| QA-04 | duration_max | Maks 7.0 detik | ✅ Split + redistribusi |
| QA-05 | cps | Maks 17 char/detik | ✅ Adjust timing |
| QA-06 | overlap | Tidak ada overlap | ✅ Trim card sebelumnya |
| QA-07 | empty_card | Tidak ada card kosong | ✅ Hapus + renumber |
| QA-08 | gap | Min 83ms antar card | ✅ Geser timestamp |
| QA-09 | glossary | Istilah glossary konsisten | ❌ Flag warning saja |

## Interface
```go
type QAValidator struct {
    MaxCharsPerLine int     // default: 42
    MaxLines        int     // default: 2
    MinDurationMS   int64   // default: 1000
    MaxDurationMS   int64   // default: 7000
    MaxCPS          float64 // default: 17.0
    MinGapMS        int64   // default: 83
    Glossary        []GlossaryTerm
}

func (v *QAValidator) Validate(cards []SubtitleCard) []QAResult
func (v *QAValidator) AutoFix(cards []SubtitleCard) ([]SubtitleCard, []QALog)
```

## CPS Calculation
```go
durationSec := float64(card.EndMS - card.StartMS) / 1000.0
charCount := len([]rune(strings.ReplaceAll(card.Text, "\n", "")))
cps := float64(charCount) / durationSec
// Gunakan []rune bukan len() agar CJK dihitung benar
```

## Auto-Fix Loop
```
Jalankan AutoFix → Validate ulang → masih ada error?
→ AutoFix lagi (max 3x total)
→ Masih error setelah 3x → ERR_QA_001, flag manual review
→ QA-09 tidak di-autofix, selalu jadi warning
```

## Test Requirements (wajib semua ada)
```
TestQA_LineLength_Pass, Fail, CJKChar
TestQA_LineCount_Pass, Fail
TestQA_DurationMin_Pass, Fail
TestQA_DurationMax_Pass, Fail
TestQA_CPS_Pass, Fail, ZeroDuration
TestQA_Overlap_Pass, Fail
TestQA_EmptyCard_Detected
TestQA_Gap_Pass, Fail
TestQA_Glossary_Warning
TestQA_AutoFix_AllChecks
TestQA_AutoFix_MaxRetry
```

## QA Report Format (qa_report.json)
```json
{
  "run_at": 1234567890,
  "total_cards": 847,
  "passed": 839,
  "warnings": 5,
  "errors": 3,
  "auto_fixed": 8,
  "results": [
    {
      "cardId": "seg-123",
      "checkId": "QA-05",
      "passed": false,
      "severity": "error",
      "detail": "18.3 CPS (maks 17.0)",
      "autoFixed": true,
      "fixAction": "extend end timestamp +200ms"
    }
  ]
}
```
