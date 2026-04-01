import { writable, derived } from 'svelte/store';

export type PipelineStep = 
  | 'idle'
  | 'import'
  | 'asr'
  | 'correct'
  | 'context'
  | 'translate'
  | 'rewrite'
  | 'qa'
  | 'completed'
  | 'error';

export interface PipelineState {
  projectId: string | null;
  currentStep: PipelineStep;
  progress: number;
  status: 'idle' | 'running' | 'paused' | 'completed' | 'error';
  error: string | null;
  stepProgress: Record<PipelineStep, number>;
}

const initialState: PipelineState = {
  projectId: null,
  currentStep: 'idle',
  progress: 0,
  status: 'idle',
  error: null,
  stepProgress: {
    idle: 0,
    import: 0,
    asr: 0,
    correct: 0,
    context: 0,
    translate: 0,
    rewrite: 0,
    qa: 0,
    completed: 0,
    error: 0
  }
};

function createPipelineStore() {
  const { subscribe, set, update } = writable<PipelineState>(initialState);

  return {
    subscribe,
    
    start(projectId: string) {
      update(s => ({
        ...s,
        projectId,
        currentStep: 'import',
        progress: 0,
        status: 'running',
        error: null
      }));
    },

    setStep(step: PipelineStep, progress: number) {
      update(s => ({
        ...s,
        currentStep: step,
        stepProgress: { ...s.stepProgress, [step]: progress }
      }));
    },

    setProgress(progress: number) {
      update(s => ({ ...s, progress }));
    },

    pause() {
      update(s => ({ ...s, status: 'paused' }));
    },

    resume() {
      update(s => ({ ...s, status: 'running' }));
    },

    complete() {
      update(s => ({
        ...s,
        currentStep: 'completed',
        progress: 100,
        status: 'completed'
      }));
    },

    setError(error: string) {
      update(s => ({
        ...s,
        currentStep: 'error',
        status: 'error',
        error
      }));
    },

    reset() {
      set(initialState);
    }
  };
}

export const pipeline = createPipelineStore();

export const isRunning = derived(pipeline, $p => $p.status === 'running');
export const isPaused = derived(pipeline, $p => $p.status === 'paused');
export const hasError = derived(pipeline, $p => $p.status === 'error');
