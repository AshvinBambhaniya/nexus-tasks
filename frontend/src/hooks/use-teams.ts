import useSWR, { useSWRConfig } from "swr";
import api from "@/lib/api";
import { Team } from "@/types";
import { useWorkspaceStore } from "@/store/workspace-store";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useTeam(teamId: number) {
  const { data, error, isLoading, mutate } = useSWR<Team>(
    teamId ? `/api/v1/teams/${teamId}` : null,
    fetcher
  );

  return {
    team: data,
    isLoading,
    isError: error,
    mutate
  };
}

export function useTeams() {
  const { mutate } = useSWRConfig();
  const { activeWorkspaceId } = useWorkspaceStore();
  
  const { data, error, isLoading } = useSWR<Team[]>(
    activeWorkspaceId ? `/api/v1/workspaces/${activeWorkspaceId}/teams` : null, 
    fetcher
  );

  const createTeam = async (name: string, description?: string) => {
    if (!activeWorkspaceId) return;
    try {
      await api.post(`/api/v1/workspaces/${activeWorkspaceId}/teams`, { name, description });
      mutate(`/api/v1/workspaces/${activeWorkspaceId}/teams`);
    } catch (err) {
      console.error("Failed to create team", err);
      throw err;
    }
  };

  const updateTeam = async (teamId: number, data: { name?: string; description?: string }) => {
     try {
         await api.patch(`/api/v1/teams/${teamId}`, data);
         mutate(activeWorkspaceId ? `/api/v1/workspaces/${activeWorkspaceId}/teams` : null);
         mutate(`/api/v1/teams/${teamId}`);
     } catch (err) {
         console.error("Failed to update team", err);
         throw err;
     }
  }

  const deleteTeam = async (teamId: number) => {
      try {
          await api.delete(`/api/v1/teams/${teamId}`);
          mutate(activeWorkspaceId ? `/api/v1/workspaces/${activeWorkspaceId}/teams` : null);
      } catch (err) {
          console.error("Failed to delete team", err);
          throw err;
      }
  }

  return {
    teams: data || [],
    isLoading,
    isError: error,
    createTeam,
    updateTeam,
    deleteTeam
  };
}

export function useTeamProjects(teamId: number) {
  const { data, error, isLoading } = useSWR<any[]>(
    teamId ? `/api/v1/teams/${teamId}/projects` : null, 
    fetcher
  );

  return {
    projects: data || [],
    isLoading,
    isError: error,
  };
}

export function useTeamMembers(teamId: number) {
  const { mutate } = useSWRConfig();
  const { data, error, isLoading } = useSWR<any[]>(
    teamId ? `/api/v1/teams/${teamId}/members` : null, 
    fetcher
  );

  const addMember = async (email: string) => {
    try {
      await api.post(`/api/v1/teams/${teamId}/members`, { email, role: "MEMBER" });
      mutate(`/api/v1/teams/${teamId}/members`);
    } catch (err) {
      console.error("Failed to add member", err);
      throw err;
    }
  };

  const removeMember = async (userId: number) => {
    try {
        await api.delete(`/api/v1/teams/${teamId}/members/${userId}`);
        mutate(`/api/v1/teams/${teamId}/members`);
    } catch (err) {
        console.error("Failed to remove member", err);
        throw err;
    }
  }

  return {
    members: data || [],
    isLoading,
    isError: error,
    addMember,
    removeMember
  };
}
