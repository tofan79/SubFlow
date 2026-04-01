import { writable, derived } from 'svelte/store';

export interface SettingsData {
  deeplApiKey: string;
  openaiApiKey: string;
  anthropicApiKey: string;
  geminiApiKey: string;
  groqApiKey: string;
  deepgramApiKey: string;
  xaiApiKey: string;
  qwenApiKey: string;
  openrouterApiKey: string;
  ollamaEndpoint: string;
  asrBackend: 'auto' | 'cpu' | 'cuda' | 'rocm' | 'coreml' | 'openvino';
  whisperModel: 'tiny' | 'base' | 'small' | 'medium' | 'large-v3';
  preferredAsr: 'local' | 'groq' | 'deepgram';
  defaultSourceLang: string;
  defaultTargetLang: string;
  defaultTonePreset: 'natural' | 'formal' | 'casual' | 'cinematic';
  maxCharsPerLine: number;
  maxLines: number;
  maxCps: number;
  minGapMs: number;
  defaultExportFormat: 'srt' | 'vtt' | 'ass' | 'txt';
}

export interface HardwareInfo {
  backend: string;
  gpuName: string;
  cudaVersion: string;
  rocmVersion: string;
  vramTotal: number;
  vramFree: number;
  computeType: string;
}

export interface SettingsState {
  settings: SettingsData;
  hardware: HardwareInfo | null;
  isLoading: boolean;
  isSaving: boolean;
  isDirty: boolean;
  error: string | null;
}

const defaultSettings: SettingsData = {
  deeplApiKey: '',
  openaiApiKey: '',
  anthropicApiKey: '',
  geminiApiKey: '',
  groqApiKey: '',
  deepgramApiKey: '',
  xaiApiKey: '',
  qwenApiKey: '',
  openrouterApiKey: '',
  ollamaEndpoint: 'http://localhost:11434',
  asrBackend: 'auto',
  whisperModel: 'medium',
  preferredAsr: 'local',
  defaultSourceLang: 'en',
  defaultTargetLang: 'id',
  defaultTonePreset: 'natural',
  maxCharsPerLine: 42,
  maxLines: 2,
  maxCps: 17.0,
  minGapMs: 83,
  defaultExportFormat: 'srt'
};

const initialState: SettingsState = {
  settings: defaultSettings,
  hardware: null,
  isLoading: false,
  isSaving: false,
  isDirty: false,
  error: null
};

function createSettingsStore() {
  const { subscribe, set, update } = writable<SettingsState>(initialState);

  return {
    subscribe,

    setSettings(settings: SettingsData) {
      update(s => ({
        ...s,
        settings,
        isLoading: false,
        isDirty: false,
        error: null
      }));
    },

    updateSetting<K extends keyof SettingsData>(key: K, value: SettingsData[K]) {
      update(s => ({
        ...s,
        settings: { ...s.settings, [key]: value },
        isDirty: true
      }));
    },

    setHardware(hardware: HardwareInfo) {
      update(s => ({ ...s, hardware }));
    },

    setLoading(isLoading: boolean) {
      update(s => ({ ...s, isLoading }));
    },

    setSaving(isSaving: boolean) {
      update(s => ({ ...s, isSaving }));
    },

    markSaved() {
      update(s => ({ ...s, isDirty: false, isSaving: false }));
    },

    setError(error: string | null) {
      update(s => ({ ...s, error, isLoading: false, isSaving: false }));
    },

    reset() {
      set(initialState);
    }
  };
}

export const settingsStore = createSettingsStore();

export const hasApiKeys = derived(
  settingsStore,
  $s => {
    const keys = $s.settings;
    return !!(
      keys.deeplApiKey ||
      keys.openaiApiKey ||
      keys.anthropicApiKey ||
      keys.geminiApiKey ||
      keys.groqApiKey ||
      keys.deepgramApiKey ||
      keys.openrouterApiKey
    );
  }
);

export const hardwareDisplay = derived(
  settingsStore,
  $s => {
    if (!$s.hardware) return 'Mendeteksi...';
    const h = $s.hardware;
    if (h.backend === 'cpu') return 'CPU';
    if (h.gpuName) return `${h.gpuName} (${h.backend.toUpperCase()})`;
    return h.backend.toUpperCase();
  }
);
