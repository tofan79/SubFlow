import { writable, derived } from 'svelte/store';

export interface ProjectInfo {
  id: string;
  name: string;
  sourcePath: string;
  state: 'idle' | 'running' | 'paused' | 'completed' | 'error';
  createdAt: number;
  updatedAt: number;
  segmentCount: number;
}

export interface ProjectState {
  projects: ProjectInfo[];
  currentProjectId: string | null;
  isLoading: boolean;
  error: string | null;
}

const initialState: ProjectState = {
  projects: [],
  currentProjectId: null,
  isLoading: false,
  error: null
};

function createProjectStore() {
  const { subscribe, set, update } = writable<ProjectState>(initialState);

  return {
    subscribe,

    setProjects(projects: ProjectInfo[]) {
      update(s => ({ ...s, projects, isLoading: false, error: null }));
    },

    addProject(project: ProjectInfo) {
      update(s => ({
        ...s,
        projects: [project, ...s.projects],
        currentProjectId: project.id
      }));
    },

    updateProject(project: ProjectInfo) {
      update(s => ({
        ...s,
        projects: s.projects.map(p => p.id === project.id ? project : p)
      }));
    },

    removeProject(id: string) {
      update(s => ({
        ...s,
        projects: s.projects.filter(p => p.id !== id),
        currentProjectId: s.currentProjectId === id ? null : s.currentProjectId
      }));
    },

    selectProject(id: string | null) {
      update(s => ({ ...s, currentProjectId: id }));
    },

    setLoading(isLoading: boolean) {
      update(s => ({ ...s, isLoading }));
    },

    setError(error: string | null) {
      update(s => ({ ...s, error, isLoading: false }));
    },

    reset() {
      set(initialState);
    }
  };
}

export const projectStore = createProjectStore();

export const currentProject = derived(
  projectStore,
  $s => $s.projects.find(p => p.id === $s.currentProjectId) ?? null
);

export const recentProjects = derived(
  projectStore,
  $s => [...$s.projects]
    .sort((a, b) => b.updatedAt - a.updatedAt)
    .slice(0, 5)
);

export const projectCount = derived(
  projectStore,
  $s => $s.projects.length
);
