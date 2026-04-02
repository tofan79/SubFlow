import { writable, derived } from 'svelte/store';
import type { Segment } from '$lib/wails';

export type SubtitleLayer = 'source' | 'l1' | 'l2' | 'dual';

export interface EditorState {
  projectId: string | null;
  segments: Segment[];
  selectedSegmentId: string | null;
  currentTimeMs: number;
  isPlaying: boolean;
  playbackRate: number;
  volume: number;
  subtitleLayer: SubtitleLayer;
  zoom: number;
  scrollPositionMs: number;
  isDirty: boolean;
}

const initialState: EditorState = {
  projectId: null,
  segments: [],
  selectedSegmentId: null,
  currentTimeMs: 0,
  isPlaying: false,
  playbackRate: 1.0,
  volume: 1.0,
  subtitleLayer: 'l2',
  zoom: 1.0,
  scrollPositionMs: 0,
  isDirty: false
};

function createEditorStore() {
  const { subscribe, set, update } = writable<EditorState>(initialState);

  return {
    subscribe,

    loadProject(projectId: string, segments: Segment[]) {
      update(s => ({
        ...s,
        projectId,
        segments,
        selectedSegmentId: null,
        currentTimeMs: 0,
        isDirty: false
      }));
    },

    setSegments(segments: Segment[]) {
      update(s => ({ ...s, segments, isDirty: true }));
    },

    updateSegment(segmentId: string, changes: Partial<Segment>) {
      update(s => ({
        ...s,
        segments: s.segments.map(seg =>
          seg.id === segmentId ? { ...seg, ...changes } : seg
        ),
        isDirty: true
      }));
    },

    selectSegment(segmentId: string | null) {
      update(s => ({ ...s, selectedSegmentId: segmentId }));
    },

    setCurrentTime(timeMs: number) {
      update(s => ({ ...s, currentTimeMs: timeMs }));
    },

    setPlaying(isPlaying: boolean) {
      update(s => ({ ...s, isPlaying }));
    },

    togglePlay() {
      update(s => ({ ...s, isPlaying: !s.isPlaying }));
    },

    setPlaybackRate(rate: number) {
      update(s => ({ ...s, playbackRate: Math.max(0.25, Math.min(2.0, rate)) }));
    },

    setVolume(volume: number) {
      update(s => ({ ...s, volume: Math.max(0, Math.min(1, volume)) }));
    },

    setSubtitleLayer(layer: SubtitleLayer) {
      update(s => ({ ...s, subtitleLayer: layer }));
    },

    setZoom(zoom: number) {
      update(s => ({ ...s, zoom: Math.max(0.1, Math.min(10, zoom)) }));
    },

    zoomIn() {
      update(s => ({ ...s, zoom: Math.min(10, s.zoom * 1.2) }));
    },

    zoomOut() {
      update(s => ({ ...s, zoom: Math.max(0.1, s.zoom / 1.2) }));
    },

    setScrollPosition(positionMs: number) {
      update(s => ({ ...s, scrollPositionMs: Math.max(0, positionMs) }));
    },

    addSegment(segment: Segment) {
      update(s => ({
        ...s,
        segments: [...s.segments, segment].sort((a, b) => a.startMs - b.startMs),
        isDirty: true
      }));
    },

    removeSegment(segmentId: string) {
      update(s => ({
        ...s,
        segments: s.segments.filter(seg => seg.id !== segmentId),
        selectedSegmentId: s.selectedSegmentId === segmentId ? null : s.selectedSegmentId,
        isDirty: true
      }));
    },

    markClean() {
      update(s => ({ ...s, isDirty: false }));
    },

    reset() {
      set(initialState);
    }
  };
}

export const editor = createEditorStore();

export const selectedSegment = derived(editor, $e =>
  $e.segments.find(s => s.id === $e.selectedSegmentId) ?? null
);

export const currentSegment = derived(editor, $e =>
  $e.segments.find(s => $e.currentTimeMs >= s.startMs && $e.currentTimeMs <= s.endMs) ?? null
);

export const segmentCount = derived(editor, $e => $e.segments.length);

export const hasUnsavedChanges = derived(editor, $e => $e.isDirty);

export const visibleSegments = derived(editor, $e => {
  const viewportDurationMs = 60000 / $e.zoom;
  const startMs = $e.scrollPositionMs;
  const endMs = startMs + viewportDurationMs;
  return $e.segments.filter(s => s.endMs >= startMs && s.startMs <= endMs);
});
