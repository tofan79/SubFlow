# Module: ASR

## Provider
| Provider | Tipe | Default | Keterangan |
|---|---|---|---|
| faster-whisper | Subprocess lokal | ✓ | Bundle via PyInstaller, tanpa API key |
| Groq Whisper | Cloud API | - | Free tier, sangat cepat |
| Deepgram | Cloud API | - | Terbaik untuk audio berisik |

## Hardware Backend Detection
Urutan prioritas (tercepat ke fallback):
1. NVIDIA CUDA — deteksi: `nvidia-smi --query-gpu=name --format=csv,noheader`
2. AMD ROCm    — deteksi: `rocm-smi --showproductname`
3. Apple CoreML — hanya darwin/arm64 (macOS future)
4. Intel OpenVINO — cek path `/opt/intel/openvino` atau Windows equivalent
5. CPU        — selalu tersedia, gunakan `int8` quantization

## Subprocess Protocol (faster-whisper)
```bash
# Input
subflow-whisper --audio audio.wav --model base --backend cuda --language ja

# Output: JSON stream, satu object per baris
{"type":"segment","start":0.0,"end":2.5,"text":"Hello everyone"}
{"type":"progress","percent":45.2}
{"type":"done","total_segments":847}
{"type":"error","code":"ERR_ASR_003","message":"..."}
```

## Compute Type per Backend
```
CUDA, ROCm, CoreML → float16
OpenVINO, CPU      → int8
```

## Model Sizes
| Model | Size | Quality | Bundled |
|---|---|---|---|
| base | 145MB | Medium | ✓ default |
| small | 465MB | Good | Download opsional |
| medium | 1.4GB | High | Download opsional |
| large-v3 | 2.9GB | Best | GPU only, download opsional |

## FFmpeg Audio Extraction
```go
// Sebelum Whisper, extract audio dari video
exec.Command("ffmpeg", "-i", videoPath,
    "-vn", "-acodec", "pcm_s16le", "-ar", "16000", "-ac", "1",
    "-y", outputWAV)
```
