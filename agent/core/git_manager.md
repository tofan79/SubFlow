# Core: Git Manager

## Commit Convention
```
task(P{n}-{id}): deskripsi singkat     ← task selesai dan verified
wip(P{n}-{id}): deskripsi — RESUME     ← checkpoint sebelum putus
fix(P{n}-{id}): perbaiki bug X         ← bugfix dari reviewer
chore: update agent files              ← update config/memory
```

## Commit Timing
```
WAJIB commit setelah:
  ✓ Task selesai + QA Loop pass
  ✓ Akan rate limit (wip commit)
  ✓ Setiap 30 menit jika masih coding (wip commit)

DILARANG commit jika:
  ✗ go build ./... masih error
  ✗ Ada syntax error
```

## Checkpoint Sebelum Putus
```bash
# 1. Tulis resume marker di kode
# // AGENT_RESUME: P3-02 whisper_local.go — lanjut di func Transcribe() line 87

# 2. Update short_term.json
# resume_hint: "internal/asr/whisper_local.go:87"

# 3. Commit
git add -A
git commit -m "wip(P3-02): whisper subprocess — AGENT_RESUME di whisper_local.go:87"
```

## Recovery Setelah Disconnect
```bash
git log --oneline -5           # lihat posisi terakhir
grep -rn "AGENT_RESUME" .      # temukan resume point
# Baca short_term.json → resume_hint → lanjut dari sana
```

## .gitignore
```
# JANGAN commit
agent/cache/llm_cache.json
agent/cache/file_cache.json
build/
frontend/dist/
frontend/node_modules/
*.exe

# WAJIB commit (penting untuk resume)
agent/memory/long_term.json
agent/memory/short_term.json
agent/tasks/task_list.json
```
