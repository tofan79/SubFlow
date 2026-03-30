# Core: Executor — Tulis Kode

## Fungsi
Mengeksekusi task brief dari planner. Menulis kode yang benar, bersih, dan ter-verify.

---

## State Machine

```
PLANNED → READING → CODING → VERIFYING → COMMITTING → DONE
                        ↑          |
                        └── FIXING ┘ (max 2x)
                                   |
                              (gagal 2x)
                                   ↓
                              FAILED → reviewer.md
```

---

## Step-by-Step

### Step 1: Cek Resume Point
```bash
grep -rn "AGENT_RESUME" internal/ frontend/
```
Jika ketemu → lanjut dari sana, jangan tulis ulang.

### Step 2: Baca File yang Dibutuhkan
- Baca semua `files_to_read` dari task brief
- Baca file yang akan diedit (jangan overwrite tanpa baca dulu)
- Baca `long_term.json` untuk architectural context

### Step 3: Tulis Kode
Urutan yang benar:
1. Types dan interfaces
2. Constructor dan init
3. Main logic
4. Helper functions
5. Test file

Setiap 30 baris atau satu fungsi selesai, tulis checkpoint:
```go
// AGENT_CHECKPOINT: P2-01 parser.go — ParseSRT() selesai, lanjut parseTimestamp()
```

### Step 4: Verifikasi (QA Loop)

**Go:**
```bash
go build ./...           # harus 0 error
go vet ./...             # harus 0 warning  
go test ./[package]/... -v  # semua test pass
```

**Svelte:**
```bash
cd frontend && npm run check  # 0 TypeScript error
```

Jika gagal → coba fix sendiri maksimal 2x.
Jika masih gagal → set `needs_reviewer: true`, serahkan ke reviewer.md.

### Step 5: Commit
```bash
git add -A
git commit -m "task(P2-01): SRT parser dengan timestamp dan segment parsing"
```

Jika rate limit/disconnect sebelum selesai:
```bash
git commit -m "wip(P2-01): SRT parser — AGENT_RESUME di parser.go:87"
```

### Step 6: Update State
```json
{
  "current_task": {
    "status": "done",
    "completed_at": "ISO timestamp",
    "build_verified": true,
    "test_verified": true
  }
}
```
Update task_list.json: set task status = "done".

---

## Go Patterns Wajib

```go
// Error: selalu wrap
return nil, fmt.Errorf("package.FuncName: %w", err)

// Context: selalu propagate
func Do(ctx context.Context, ...) error {
    select {
    case <-ctx.Done(): return ctx.Err()
    default:
    }
}

// Concurrency: semaphore
sem := make(chan struct{}, maxConcurrent)
sem <- struct{}{}; defer func() { <-sem }()
```

## Svelte Patterns Wajib

```typescript
// Selalu cleanup events
onMount(() => EventsOn('event', handler))
onDestroy(() => EventsOff('event'))

// Selalu typed store
const store = writable<TypedState>(initialState)

// Error handling async
async function call() {
  try { await WailsMethod() }
  catch (err) { errorStore.set(String(err)) }
}
```
