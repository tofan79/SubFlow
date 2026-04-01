package asr

import (
	"errors"
	"strings"
	"testing"
)

func TestDetectBestBackend_PriorityCUDA(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		if name == "nvidia-smi" {
			return []byte("NVIDIA RTX 4090\n"), nil
		}
		return nil, errors.New("unexpected command")
	}
	goosFn = func() string { return "linux" }
	goarchFn = func() string { return "amd64" }
	dirExistsFn = func(path string) bool { return false }

	if got := DetectBestBackend(); got != BackendCUDA {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendCUDA)
	}

	hw := DetectHardware()
	if hw.Backend != BackendCUDA {
		t.Fatalf("DetectHardware().Backend=%s, want %s", hw.Backend, BackendCUDA)
	}
	if hw.GPUName != "NVIDIA RTX 4090" {
		t.Fatalf("DetectHardware().GPUName=%q, want %q", hw.GPUName, "NVIDIA RTX 4090")
	}
	if hw.ComputeType != string(ComputeTypeForBackend(BackendCUDA)) {
		t.Fatalf("DetectHardware().ComputeType=%q, want %q", hw.ComputeType, string(ComputeTypeForBackend(BackendCUDA)))
	}
}

func TestDetectBestBackend_PriorityROCm(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		switch name {
		case "nvidia-smi":
			return nil, errors.New("not found")
		case "rocm-smi":
			return []byte("GPU[0]          : Radeon RX 7900 XTX\n"), nil
		default:
			return nil, errors.New("unexpected command")
		}
	}
	goosFn = func() string { return "linux" }
	goarchFn = func() string { return "amd64" }
	dirExistsFn = func(path string) bool { return false }

	if got := DetectBestBackend(); got != BackendROCm {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendROCm)
	}

	hw := DetectHardware()
	if hw.Backend != BackendROCm {
		t.Fatalf("DetectHardware().Backend=%s, want %s", hw.Backend, BackendROCm)
	}
	if !strings.Contains(hw.GPUName, "Radeon") {
		t.Fatalf("DetectHardware().GPUName=%q, want contains %q", hw.GPUName, "Radeon")
	}
}

func TestDetectBestBackend_CoreML(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("not found")
	}
	goosFn = func() string { return "darwin" }
	goarchFn = func() string { return "arm64" }
	dirExistsFn = func(path string) bool { return false }

	if got := DetectBestBackend(); got != BackendCoreML {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendCoreML)
	}

	hw := DetectHardware()
	if hw.Backend != BackendCoreML {
		t.Fatalf("DetectHardware().Backend=%s, want %s", hw.Backend, BackendCoreML)
	}
	if hw.GPUName == "" {
		t.Fatalf("DetectHardware().GPUName should not be empty")
	}
}

func TestDetectBestBackend_OpenVINO(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("not found")
	}
	goosFn = func() string { return "linux" }
	goarchFn = func() string { return "amd64" }
	dirExistsFn = func(path string) bool { return path == "/opt/intel/openvino" }

	if got := DetectBestBackend(); got != BackendOpenVINO {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendOpenVINO)
	}
}

func TestDetectBestBackend_FallbackCPU(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		return nil, errors.New("not found")
	}
	goosFn = func() string { return "linux" }
	goarchFn = func() string { return "amd64" }
	dirExistsFn = func(path string) bool { return false }

	if got := DetectBestBackend(); got != BackendCPU {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendCPU)
	}

	hw := DetectHardware()
	if hw.Backend != BackendCPU {
		t.Fatalf("DetectHardware().Backend=%s, want %s", hw.Backend, BackendCPU)
	}
	if hw.ComputeType != string(ComputeTypeForBackend(BackendCPU)) {
		t.Fatalf("DetectHardware().ComputeType=%q, want %q", hw.ComputeType, string(ComputeTypeForBackend(BackendCPU)))
	}
}

func TestDetectBestBackend_NeverPanics(t *testing.T) {
	oldCmd := commandOutputFn
	oldDir := dirExistsFn
	oldGOOS := goosFn
	oldGOARCH := goarchFn
	defer func() {
		commandOutputFn = oldCmd
		dirExistsFn = oldDir
		goosFn = oldGOOS
		goarchFn = oldGOARCH
	}()

	commandOutputFn = func(name string, args ...string) ([]byte, error) {
		panic("boom")
	}
	goosFn = func() string { return "linux" }
	goarchFn = func() string { return "amd64" }
	dirExistsFn = func(path string) bool { return false }

	if got := DetectBestBackend(); got != BackendCPU {
		t.Fatalf("DetectBestBackend()=%s, want %s", got, BackendCPU)
	}
}
