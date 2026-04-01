/**
 * Wails bindings wrapper
 * Re-exports all IPC methods and event functions from Wails runtime.
 * 
 * Usage:
 *   import { GetProjects, EventsOn } from '$lib/wails';
 */

// Import Wails runtime functions
// These are injected by Wails at build time
declare global {
  interface Window {
    go: {
      main: {
        App: AppMethods;
      };
    };
    runtime: WailsRuntime;
  }
}

// =============================================================================
// TYPE DEFINITIONS (matching app.go structs)
// =============================================================================

export interface ProjectInfo {
  id: string;
  name: string;
  sourcePath: string;
  state: string;
  createdAt: number;
  updatedAt: number;
  segmentCount: number;
}

export interface PipelineConfig {
  projectId: string;
  sourceLang: string;
  targetLang: string;
  contentMode: string;
  tonePreset: string;
  asrProvider: string;
  asrModel: string;
  translateProvider: string;
  rewriteProvider: string;
}

export interface SettingsData {
  // API Keys
  deeplApiKey: string;
  openaiApiKey: string;
  anthropicApiKey: string;
  geminiApiKey: string;
  groqApiKey: string;
  deepgramApiKey: string;
  xaiApiKey: string;
  qwenApiKey: string;
  ollamaEndpoint: string;

  // ASR Settings
  asrBackend: 'auto' | 'cpu' | 'cuda' | 'rocm' | 'coreml' | 'openvino';
  whisperModel: 'tiny' | 'base' | 'small' | 'medium' | 'large-v3';
  preferredAsr: 'local' | 'groq' | 'deepgram';

  // Translation Settings
  defaultSourceLang: string;
  defaultTargetLang: string;
  defaultTonePreset: 'natural' | 'formal' | 'casual' | 'cinematic';

  // QA Settings
  maxCharsPerLine: number;
  maxLines: number;
  maxCps: number;
  minGapMs: number;

  // Export Settings
  defaultExportFormat: 'srt' | 'vtt' | 'ass' | 'txt';
}

export interface HardwareInfo {
  backend: 'cpu' | 'cuda' | 'rocm' | 'coreml' | 'openvino';
  gpuName: string;
  cudaVersion: string;
  rocmVersion: string;
  vramTotal: number;
  vramFree: number;
  computeType: 'int8' | 'float16' | 'float32';
}

export interface Segment {
  id: string;
  projectId: string;
  index: number;
  startMs: number;
  endMs: number;
  source: string;
  l1: string;
  l2: string;
  speaker: string;
  emotion: string;
  qaStatus: 'pass' | 'warn' | 'error' | 'pending';
}

export interface QAResult {
  cardId: string;
  checkId: string;
  passed: boolean;
  severity: 'error' | 'warning';
  detail: string;
  autoFixed: boolean;
  fixAction: string;
}

export interface QAReport {
  runAt: number;
  totalCards: number;
  passed: number;
  warnings: number;
  errors: number;
  autoFixed: number;
  results: QAResult[];
}

export interface GlossaryTerm {
  id: string;
  sourceTerm: string;
  targetTerm: string;
  caseSensitive: boolean;
  notes: string;
}

export interface ExportOptions {
  projectId: string;
  format: 'srt' | 'vtt' | 'ass' | 'txt';
  outputDir: string;
  layer: 'source' | 'l1' | 'l2';
  dualSubtitle: boolean;
}

export interface CostEstimate {
  asrCost: number;
  translateCost: number;
  rewriteCost: number;
  totalCost: number;
  currency: string;
}

export interface AppStats {
  totalProjects: number;
  totalSegments: number;
  totalCharsProcessed: number;
  totalMinutesAsr: number;
}

// =============================================================================
// APP METHOD SIGNATURES
// =============================================================================

interface AppMethods {
  // Project Management
  GetProjects(): Promise<ProjectInfo[]>;
  GetProject(id: string): Promise<ProjectInfo>;
  CreateProject(sourcePath: string, name: string): Promise<ProjectInfo>;
  DeleteProject(id: string): Promise<void>;

  // Pipeline Control
  StartPipeline(config: PipelineConfig): Promise<void>;
  PausePipeline(projectId: string): Promise<void>;
  ResumePipeline(projectId: string): Promise<void>;
  CancelPipeline(projectId: string): Promise<void>;
  GetPipelineState(projectId: string): Promise<string>;

  // Settings
  GetSettings(): Promise<SettingsData>;
  SetSetting(key: string, value: unknown): Promise<void>;
  SaveSettings(settings: SettingsData): Promise<void>;

  // ASR
  DetectHardware(): Promise<HardwareInfo>;
  GetAvailableModels(): Promise<string[]>;
  DownloadModel(modelName: string): Promise<void>;

  // Segments
  GetSegments(projectId: string): Promise<Segment[]>;
  UpdateSegment(segment: Segment): Promise<void>;
  SplitSegment(segmentId: string, splitAtMs: number, splitAtChar: number): Promise<Segment[]>;
  MergeSegments(segmentId1: string, segmentId2: string): Promise<Segment>;
  DeleteSegment(segmentId: string): Promise<void>;
  RetryL1(segmentId: string): Promise<Segment>;
  RetryL2(segmentId: string): Promise<Segment>;

  // QA
  RunQA(projectId: string): Promise<QAReport>;
  RunQAAutoFix(projectId: string): Promise<QAReport>;

  // Glossary
  GetGlossary(): Promise<GlossaryTerm[]>;
  AddGlossaryTerm(term: GlossaryTerm): Promise<GlossaryTerm>;
  UpdateGlossaryTerm(term: GlossaryTerm): Promise<void>;
  DeleteGlossaryTerm(id: string): Promise<void>;
  ImportGlossary(filePath: string): Promise<number>;
  ExportGlossary(filePath: string): Promise<void>;

  // Export
  Export(options: ExportOptions): Promise<string>;

  // Dialogs
  SelectFile(title: string, filters: string[]): Promise<string>;
  SelectDirectory(title: string): Promise<string>;

  // Cost
  EstimateCost(projectId: string, config: PipelineConfig): Promise<CostEstimate>;

  // Stats
  GetStats(): Promise<AppStats>;
}

// =============================================================================
// WAILS RUNTIME INTERFACE
// =============================================================================

type EventCallback = (...args: unknown[]) => void;

interface WailsRuntime {
  EventsOn(eventName: string, callback: EventCallback): () => void;
  EventsOff(eventName: string, ...additionalEventNames: string[]): void;
  EventsEmit(eventName: string, ...data: unknown[]): void;
  EventsOnce(eventName: string, callback: EventCallback): () => void;
  EventsOnMultiple(eventName: string, callback: EventCallback, maxCallbacks: number): () => void;
  WindowSetTitle(title: string): void;
  WindowFullscreen(): void;
  WindowUnfullscreen(): void;
  WindowIsFullscreen(): Promise<boolean>;
  WindowCenter(): void;
  WindowReload(): void;
  WindowSetSystemDefaultTheme(): void;
  WindowSetLightTheme(): void;
  WindowSetDarkTheme(): void;
  WindowMinimise(): void;
  WindowUnminimise(): void;
  WindowIsMinimised(): Promise<boolean>;
  WindowMaximise(): void;
  WindowUnmaximise(): void;
  WindowIsMaximised(): Promise<boolean>;
  WindowToggleMaximise(): void;
  WindowSetBackgroundColour(R: number, G: number, B: number, A: number): void;
  Quit(): void;
}

// =============================================================================
// APP METHOD EXPORTS
// =============================================================================

// Helper to safely access app methods
function getApp(): AppMethods {
  if (typeof window !== 'undefined' && window.go?.main?.App) {
    return window.go.main.App;
  }
  // Return stub for SSR/build time
  throw new Error('Wails runtime not available');
}

// Project Management
export const GetProjects = () => getApp().GetProjects();
export const GetProject = (id: string) => getApp().GetProject(id);
export const CreateProject = (sourcePath: string, name: string) => 
  getApp().CreateProject(sourcePath, name);
export const DeleteProject = (id: string) => getApp().DeleteProject(id);

// Pipeline Control
export const StartPipeline = (config: PipelineConfig) => getApp().StartPipeline(config);
export const PausePipeline = (projectId: string) => getApp().PausePipeline(projectId);
export const ResumePipeline = (projectId: string) => getApp().ResumePipeline(projectId);
export const CancelPipeline = (projectId: string) => getApp().CancelPipeline(projectId);
export const GetPipelineState = (projectId: string) => getApp().GetPipelineState(projectId);

// Settings
export const GetSettings = () => getApp().GetSettings();
export const SetSetting = (key: string, value: unknown) => getApp().SetSetting(key, value);
export const SaveSettings = (settings: SettingsData) => getApp().SaveSettings(settings);

// ASR
export const DetectHardware = () => getApp().DetectHardware();
export const GetAvailableModels = () => getApp().GetAvailableModels();
export const DownloadModel = (modelName: string) => getApp().DownloadModel(modelName);

// Segments
export const GetSegments = (projectId: string) => getApp().GetSegments(projectId);
export const UpdateSegment = (segment: Segment) => getApp().UpdateSegment(segment);
export const SplitSegment = (segmentId: string, splitAtMs: number, splitAtChar: number) =>
  getApp().SplitSegment(segmentId, splitAtMs, splitAtChar);
export const MergeSegments = (segmentId1: string, segmentId2: string) =>
  getApp().MergeSegments(segmentId1, segmentId2);
export const DeleteSegment = (segmentId: string) => getApp().DeleteSegment(segmentId);
export const RetryL1 = (segmentId: string) => getApp().RetryL1(segmentId);
export const RetryL2 = (segmentId: string) => getApp().RetryL2(segmentId);

// QA
export const RunQA = (projectId: string) => getApp().RunQA(projectId);
export const RunQAAutoFix = (projectId: string) => getApp().RunQAAutoFix(projectId);

// Glossary
export const GetGlossary = () => getApp().GetGlossary();
export const AddGlossaryTerm = (term: GlossaryTerm) => getApp().AddGlossaryTerm(term);
export const UpdateGlossaryTerm = (term: GlossaryTerm) => getApp().UpdateGlossaryTerm(term);
export const DeleteGlossaryTerm = (id: string) => getApp().DeleteGlossaryTerm(id);
export const ImportGlossary = (filePath: string) => getApp().ImportGlossary(filePath);
export const ExportGlossary = (filePath: string) => getApp().ExportGlossary(filePath);

// Export
export const ExportSubtitle = (options: ExportOptions) => getApp().Export(options);

// Dialogs
export const SelectFile = (title: string, filters: string[] = []) => 
  getApp().SelectFile(title, filters);
export const SelectDirectory = (title: string) => getApp().SelectDirectory(title);

// Cost
export const EstimateCost = (projectId: string, config: PipelineConfig) =>
  getApp().EstimateCost(projectId, config);

// Stats
export const GetStats = () => getApp().GetStats();

// =============================================================================
// EVENT EXPORTS
// =============================================================================

// Helper to safely access runtime
function getRuntime(): WailsRuntime {
  if (typeof window !== 'undefined' && window.runtime) {
    return window.runtime;
  }
  // Return stub for SSR/build time
  throw new Error('Wails runtime not available');
}

// Event functions
export const EventsOn = (eventName: string, callback: EventCallback) =>
  getRuntime().EventsOn(eventName, callback);
export const EventsOff = (eventName: string, ...additionalEventNames: string[]) =>
  getRuntime().EventsOff(eventName, ...additionalEventNames);
export const EventsEmit = (eventName: string, ...data: unknown[]) =>
  getRuntime().EventsEmit(eventName, ...data);
export const EventsOnce = (eventName: string, callback: EventCallback) =>
  getRuntime().EventsOnce(eventName, callback);
export const EventsOnMultiple = (eventName: string, callback: EventCallback, maxCallbacks: number) =>
  getRuntime().EventsOnMultiple(eventName, callback, maxCallbacks);

// Window functions
export const WindowSetTitle = (title: string) => getRuntime().WindowSetTitle(title);
export const WindowFullscreen = () => getRuntime().WindowFullscreen();
export const WindowUnfullscreen = () => getRuntime().WindowUnfullscreen();
export const WindowIsFullscreen = () => getRuntime().WindowIsFullscreen();
export const WindowCenter = () => getRuntime().WindowCenter();
export const WindowReload = () => getRuntime().WindowReload();
export const WindowSetDarkTheme = () => getRuntime().WindowSetDarkTheme();
export const WindowMinimise = () => getRuntime().WindowMinimise();
export const WindowUnminimise = () => getRuntime().WindowUnminimise();
export const WindowIsMinimised = () => getRuntime().WindowIsMinimised();
export const WindowMaximise = () => getRuntime().WindowMaximise();
export const WindowUnmaximise = () => getRuntime().WindowUnmaximise();
export const WindowIsMaximised = () => getRuntime().WindowIsMaximised();
export const WindowToggleMaximise = () => getRuntime().WindowToggleMaximise();
export const WindowSetBackgroundColour = (R: number, G: number, B: number, A: number) =>
  getRuntime().WindowSetBackgroundColour(R, G, B, A);
export const Quit = () => getRuntime().Quit();

// =============================================================================
// EVENT NAMES (for type-safe event handling)
// =============================================================================

export const WailsEvents = {
  // Pipeline events
  PIPELINE_STARTED: 'pipeline:started',
  PIPELINE_STEP_CHANGED: 'pipeline:step',
  PIPELINE_PROGRESS: 'pipeline:progress',
  PIPELINE_COMPLETED: 'pipeline:completed',
  PIPELINE_ERROR: 'pipeline:error',
  PIPELINE_PAUSED: 'pipeline:paused',
  PIPELINE_RESUMED: 'pipeline:resumed',

  // ASR events
  ASR_HARDWARE_DETECTED: 'asr:hardware',
  ASR_MODEL_DOWNLOAD_PROGRESS: 'asr:model:progress',
  ASR_MODEL_DOWNLOAD_COMPLETE: 'asr:model:complete',

  // File events
  FILE_DROPPED: 'file:dropped',

  // Segment events
  SEGMENT_UPDATED: 'segment:updated',
  SEGMENT_SPLIT: 'segment:split',
  SEGMENT_MERGED: 'segment:merged',

  // QA events
  QA_STARTED: 'qa:started',
  QA_PROGRESS: 'qa:progress',
  QA_COMPLETED: 'qa:completed',
} as const;

export type WailsEventName = typeof WailsEvents[keyof typeof WailsEvents];
