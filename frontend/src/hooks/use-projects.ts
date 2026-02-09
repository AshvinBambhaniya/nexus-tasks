import useSWR, { useSWRConfig } from "swr";
import api from "@/lib/api";
import { Project, ProjectMember, ProjectTeam } from "@/types";
import { useWorkspaceStore } from "@/store/workspace-store";
import { create } from "zustand";
import { persist } from "zustand/middleware";

// Add Project Store logic here for selection
interface ProjectState {
  activeProjectId: number | null;
  setActiveProjectId: (id: number) => void;
}

export const useProjectStore = create<ProjectState>()(
  persist(
    (set) => ({
      activeProjectId: null,
      setActiveProjectId: (id) => set({ activeProjectId: id }),
    }),
    {
      name: "project-storage",
    }
  )
);

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useProjectMembers(projectId?: number) {
  const { mutate } = useSWRConfig();

  const key = projectId ? `/api/v1/projects/${projectId}/members` : null;
  const { data, error, isLoading } = useSWR<ProjectMember[]>(key, fetcher);

  const addMember = async (email: string) => {
    if (!projectId) return;
    try {
      await api.post(`/api/v1/projects/${projectId}/members`, {
        email,
        role: "MEMBER",
      });
      mutate(key);
    } catch (err) {
      console.error("Failed to add member", err);
      throw err;
    }
  };

  const removeMember = async (userId: number) => {
    if (!projectId) return;
    try {
      await api.delete(`/api/v1/projects/${projectId}/members/${userId}`);
      mutate(key);
    } catch (err) {
      console.error("Failed to remove member", err);
      throw err;
    }
  };

  return {
    members: data || [],
    isLoading,
    isError: error,
    addMember,
    removeMember,
  };
}

export function useProjectTeams(projectId?: number) {
  const { mutate } = useSWRConfig();

  const key = projectId ? `/api/v1/projects/${projectId}/teams` : null;
  const { data, error, isLoading } = useSWR<ProjectTeam[]>(key, fetcher);

  const addTeam = async (teamId: number) => {
    if (!projectId) return;
    try {
      await api.post(`/api/v1/projects/${projectId}/teams`, {
        team_id: teamId,
      });
      mutate(key);
    } catch (err) {
      console.error("Failed to add team", err);
      throw err;
    }
  };

  const removeTeam = async (teamId: number) => {
    if (!projectId) return;
    try {
      await api.delete(`/api/v1/projects/${projectId}/teams/${teamId}`);
      mutate(key);
    } catch (err) {
      console.error("Failed to remove team", err);
      throw err;
    }
  };

  return {
    teams: data || [],
    isLoading,
    isError: error,
    addTeam,
    removeTeam,
  };
}

export function useProject(id: number) {
  const { data, error, isLoading, mutate } = useSWR<Project>(
    id ? `/api/v1/projects/${id}` : null,
    fetcher
  );

  return {
    project: data,
    isLoading,
    isError: error,
    mutate,
  };
}

export function useProjects() {
  const { mutate } = useSWRConfig();
  const { activeWorkspaceId } = useWorkspaceStore();

  const { data, error, isLoading } = useSWR<Project[]>(
    activeWorkspaceId
      ? `/api/v1/workspaces/${activeWorkspaceId}/projects`
      : null,
    fetcher
  );

  const createProject = async (name: string, description?: string) => {
    if (!activeWorkspaceId) return;
    try {
      await api.post(`/api/v1/workspaces/${activeWorkspaceId}/projects`, {
        name,
        description,
      });
      mutate(`/api/v1/workspaces/${activeWorkspaceId}/projects`);
    } catch (err) {
      console.error("Failed to create project", err);
      throw err;
    }
  };

  const updateProject = async (
    id: number,
    data: { name?: string; description?: string; is_archived?: boolean }
  ) => {
    if (!activeWorkspaceId) return;
    try {
      await api.patch(`/api/v1/projects/${id}`, data);
      mutate(`/api/v1/workspaces/${activeWorkspaceId}/projects`);
    } catch (err) {
      console.error("Failed to update project", err);
      throw err;
    }
  };

  return {
    projects: data || [],
    isLoading,
    isError: error,
    createProject,
    updateProject,
  };
}
