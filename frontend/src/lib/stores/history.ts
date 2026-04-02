import { writable, derived, get } from 'svelte/store';
import type { Segment } from '$lib/wails';

const MAX_HISTORY = 50;

interface HistoryEntry {
  timestamp: number;
  action: string;
  segments: Segment[];
  selectedSegmentId: string | null;
}

interface HistoryState {
  entries: HistoryEntry[];
  currentIndex: number;
}

const initialState: HistoryState = {
  entries: [],
  currentIndex: -1
};

function createHistoryStore() {
  const { subscribe, set, update } = writable<HistoryState>(initialState);

  return {
    subscribe,

    push(action: string, segments: Segment[], selectedSegmentId: string | null) {
      update(s => {
        const newEntries = s.entries.slice(0, s.currentIndex + 1);

        newEntries.push({
          timestamp: Date.now(),
          action,
          segments: JSON.parse(JSON.stringify(segments)),
          selectedSegmentId
        });

        if (newEntries.length > MAX_HISTORY) {
          newEntries.shift();
        }

        return {
          entries: newEntries,
          currentIndex: newEntries.length - 1
        };
      });
    },

    undo(): HistoryEntry | null {
      const state = get({ subscribe });
      if (state.currentIndex <= 0) return null;

      let result: HistoryEntry | null = null;
      update(s => {
        if (s.currentIndex > 0) {
          result = s.entries[s.currentIndex - 1];
          return { ...s, currentIndex: s.currentIndex - 1 };
        }
        return s;
      });
      return result;
    },

    redo(): HistoryEntry | null {
      const state = get({ subscribe });
      if (state.currentIndex >= state.entries.length - 1) return null;

      let result: HistoryEntry | null = null;
      update(s => {
        if (s.currentIndex < s.entries.length - 1) {
          result = s.entries[s.currentIndex + 1];
          return { ...s, currentIndex: s.currentIndex + 1 };
        }
        return s;
      });
      return result;
    },

    clear() {
      set(initialState);
    },

    init(segments: Segment[], selectedSegmentId: string | null) {
      set({
        entries: [{
          timestamp: Date.now(),
          action: 'init',
          segments: JSON.parse(JSON.stringify(segments)),
          selectedSegmentId
        }],
        currentIndex: 0
      });
    }
  };
}

export const history = createHistoryStore();

export const canUndo = derived(history, $h => $h.currentIndex > 0);
export const canRedo = derived(history, $h => $h.currentIndex < $h.entries.length - 1);
export const historyLength = derived(history, $h => $h.entries.length);
export const currentHistoryIndex = derived(history, $h => $h.currentIndex);
