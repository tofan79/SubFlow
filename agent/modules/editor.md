# Module: Editor UI

## Layout 4 Panel
```
┌──────────────────────────┬────────────────────────────┐
│  VIDEO PREVIEW (60%)     │  SEGMENT LIST (40%)        │
│  [subtitle overlay]      │  # Timestamp  Source  L2   │
│  ◀◀ ▶ ▶▶ 🔊 00:02/10:14  │  ► 3  00:05-00:08  ...    │
├──────────────────────────│  [Translation L1]          │
│  WAVEFORM + TIMELINE     │  Halo semuanya...          │
│  ▓▓▒▒▓▓▒▒▓▓▒▒            │  [AI Rewrite L2]           │
│  [══▓▓▓══][══▓▓══]       │  Halo semuanya, welcome... │
└──────────────────────────┴────────────────────────────┘
```

## Komponen Svelte
```
editor/+page.svelte         ← layout + EditorState store
  VideoPreview.svelte       ← video player + subtitle overlay
  SegmentList.svelte        ← list + QA indicator
    SegmentRow.svelte       ← satu baris segmen
  Timeline.svelte           ← waveform + drag cards
  InlineEditor.svelte       ← edit teks + action buttons
```

## Svelte Store: EditorState
```typescript
interface EditorState {
  segments: Segment[]
  selectedId: string | null
  editingField: 'source' | 'l1' | 'l2' | null
  playbackMs: number
  isPlaying: boolean
  zoom: number
}
```

## Sync Antar Panel
```
Klik segmen di SegmentList → video jump ke timestamp, timeline scroll ke card
Drag card di Timeline      → timestamp diupdate real-time, QA re-run
Edit teks di InlineEditor  → subtitle overlay diupdate, CPS recalculated
```

## QA Indicator Warna
```
🟢 qa-pass    → semua check lulus
🟡 qa-warn    → ada warning (QA-09 glossary)
🔴 qa-error   → ada error harus difix
⚪ qa-pending → belum diproses
```

## Tombol Aksi per Segmen
```
[ Split ]  [ Merge ↑ ]  [ Merge ↓ ]  [ Retry L1 ]  [ Retry L2 ]  [ Hapus ]
```

## Keyboard Shortcuts
```
Space        → Play/Pause
↑ / ↓        → Pindah segmen
Enter        → Aktifkan inline editor
Escape       → Tutup editor, batalkan
Ctrl+S       → Simpan semua
Ctrl+Z       → Undo (max 50)
Ctrl+Y       → Redo
Ctrl+Enter   → Simpan + pindah ke berikutnya
Ctrl+R       → Retry L2 segmen aktif
Ctrl+Scroll  → Zoom timeline
```

## Undo/Redo (50 langkah)
```typescript
interface HistoryEntry {
  type: 'edit' | 'split' | 'merge' | 'delete' | 'timestamp'
  segmentId: string
  before: Partial<Segment>
  after: Partial<Segment>
}
// Stack max 50 entries
// Undo: pop, apply before
// Redo: push ke redo stack
```

## Timeline Drag
```typescript
// Snap ke minimum 83ms dari card tetangga
// Alt + drag = disable snap sementara
// mouseup → commit ke DB via Wails IPC
// mousemove → update store real-time (preview saja)
```

## Dual Subtitle Mode
```
Toggle: [Source] | [L1] | [L2] | [Dual: Source + L2]
Saat Dual: tampilkan source di atas, L2 di bawah
```
