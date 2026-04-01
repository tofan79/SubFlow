package asr

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const whisperHardwareDetectTimeout = 2 * time.Second

var commandOutputFn = func(name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), whisperHardwareDetectTimeout)
	defer cancel()
	return exec.CommandContext(ctx, name, args...).Output()
}

var dirExistsFn = func(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

var goosFn = func() string { return runtime.GOOS }
var goarchFn = func() string { return runtime.GOARCH }

func DetectBestBackend() (backend Backend) {
	defer func() {
		if r := recover(); r != nil {
			backend = BackendCPU
		}
	}()

	if IsCUDAAvailable() {
		return BackendCUDA
	}
	if IsROCmAvailable() {
		return BackendROCm
	}
	if IsCoreMLAvailable() {
		return BackendCoreML
	}
	if IsOpenVINOAvailable() {
		return BackendOpenVINO
	}
	return BackendCPU
}

func DetectHardware() (info HardwareInfo) {
	defer func() {
		if r := recover(); r != nil {
			info = HardwareInfo{
				Backend:     BackendCPU,
				ComputeType: string(ComputeTypeForBackend(BackendCPU)),
			}
		}
	}()

	backend := DetectBestBackend()
	info.Backend = backend
	info.ComputeType = string(ComputeTypeForBackend(backend))

	switch backend {
	case BackendCUDA:
		if name, ok := detectNVIDIAGPUName(); ok {
			info.GPUName = name
		}
	case BackendROCm:
		if name, ok := detectAMDGPUName(); ok {
			info.GPUName = name
		}
	case BackendCoreML:
		info.GPUName = "Apple Silicon"
	case BackendOpenVINO:
		info.GPUName = "Intel OpenVINO"
	}

	return info
}

func IsCUDAAvailable() bool {
	name, ok := detectNVIDIAGPUName()
	return ok && strings.TrimSpace(name) != ""
}

func IsROCmAvailable() bool {
	name, ok := detectAMDGPUName()
	return ok && strings.TrimSpace(name) != ""
}

func IsCoreMLAvailable() bool {
	return goosFn() == "darwin" && goarchFn() == "arm64"
}

func IsOpenVINOAvailable() bool {
	for _, p := range openVINOPaths() {
		if dirExistsFn(p) {
			return true
		}
	}
	return false
}

func detectNVIDIAGPUName() (string, bool) {
	out, err := commandOutputFn("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
	if err != nil {
		return "", false
	}
	line := firstNonEmptyLine(string(out))
	if line == "" {
		return "", false
	}
	return line, true
}

func detectAMDGPUName() (string, bool) {
	out, err := commandOutputFn("rocm-smi", "--showproductname")
	if err != nil {
		return "", false
	}
	line := firstNonEmptyLine(string(out))
	if line == "" {
		return "", false
	}
	return line, true
}

func firstNonEmptyLine(s string) string {
	for _, ln := range strings.Split(s, "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			return ln
		}
	}
	return ""
}

func openVINOPaths() []string {
	switch goosFn() {
	case "linux":
		return []string{"/opt/intel/openvino"}
	case "windows":
		paths := []string{
			"C:\\Program Files\\Intel\\openvino",
			"C:\\Program Files (x86)\\Intel\\openvino",
			"C:\\Program Files\\Intel\\openvino_2024",
			"C:\\Program Files (x86)\\Intel\\openvino_2024",
		}
		if pf := os.Getenv("ProgramFiles"); pf != "" {
			paths = append(paths, pf+"\\Intel\\openvino")
			paths = append(paths, pf+"\\Intel\\openvino_2024")
		}
		if pf := os.Getenv("ProgramFiles(x86)"); pf != "" {
			paths = append(paths, pf+"\\Intel\\openvino")
			paths = append(paths, pf+"\\Intel\\openvino_2024")
		}
		return paths
	default:
		return nil
	}
}
