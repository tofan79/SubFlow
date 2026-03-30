# Module: Subtitle Parser & Segmenter

## Package Rules
Package `subtitle` adalah pure logic — DILARANG import package internal lain.

## Format yang Didukung
| Format | Parser | Writer |
|---|---|---|
| SRT | parser.go | writer.go |
| VTT | parser_vtt.go | writer.go |
| ASS/SSA | parser_ass.go | writer.go |
| TXT | - | writer.go |

## SubtitleCard Type
```go
type SubtitleCard struct {
    ID      string  // UUID
    Index   int     // 1-based
    StartMS int64   // milliseconds
    EndMS   int64   // milliseconds
    Text    string  // "\n" sebagai line separator
}
```

## SRT Parser
```
Format SRT:
  1
  00:00:01,200 --> 00:00:03,500
  Hello everyone

  2
  00:00:03,600 --> 00:00:05,100
  Welcome back to my channel

Timestamp format: HH:MM:SS,mmm
```

## VTT Parser
```
Format VTT:
  WEBVTT

  00:00:01.200 --> 00:00:03.500
  Hello everyone

Timestamp format: HH:MM:SS.mmm (titik bukan koma)
```

## ASS Parser
```
Format ASS — ambil hanya baris Dialogue:
  Dialogue: 0,0:00:01.20,0:00:03.50,Default,,0,0,0,,Hello everyone

Parse: layer, start, end, style, name, marginL, marginR, marginV, effect, text
Strip formatting tags: {\an8}, {\i1}, {\b1}, dll
```

## Segmentation Rules (HARD LIMITS)
```
Max chars per baris : 42
Max baris per card  : 2
Min durasi          : 1.0 detik (1000ms)
Max durasi          : 7.0 detik (7000ms)
Max CPS             : 17 char/detik
Min gap             : 83ms (2 frame @ 24fps)
Line break          : Di batas klausa, bukan tengah kata
```

## Writer Output
```go
// SRT: nomor, timestamp HH:MM:SS,mmm --> HH:MM:SS,mmm, teks, baris kosong
// VTT: header WEBVTT, timestamp dengan titik
// ASS: header minimal, Dialogue lines
// TXT: timestamp [HH:MM:SS] teks
// Dual Subtitle: source di atas, terjemahan di bawah, satu card
```

## Test Requirements
```
TestSRT_ParseBasic          — 3 segmen normal
TestSRT_ParseEmptyLines     — toleransi baris kosong tambahan  
TestSRT_ParseUTF8           — karakter CJK dan emoji
TestVTT_ParseBasic
TestASS_ParseDialogue
TestASS_StripTags           — pastikan tag formatting dihapus
TestWriter_SRT
TestWriter_DualSubtitle     — pastikan source + terjemahan dalam satu card
TestSegmenter_SplitLongLine — baris > 42 char harus dipecah
TestSegmenter_RespectWord   — pecah di spasi, bukan tengah kata
```
