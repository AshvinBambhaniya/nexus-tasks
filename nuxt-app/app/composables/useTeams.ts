import type { Team, TeamMember } from "~/types";

export const useTeams = () => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: teams,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Team[]>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams`
        : "/api/v1/workspaces/0/teams",
    {
      key: `teams-list-${workspaceStore.activeWorkspaceId}`,
      watch: [() => workspaceStore.activeWorkspaceId],
    }
  );

  const createTeam = async (name: string, description?: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams`,
        {
          method: "POST",
          body: { name, description },
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to create team", err);
      throw err;
    }
  };

  const updateTeam = async (
    teamId: string,
    data: { name?: string; description?: string }
  ) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}`,
        {
          method: "PATCH",
          body: data,
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to update team", err);
      throw err;
    }
  };

  const deleteTeam = async (teamId: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}`,
        {
          method: "DELETE",
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to delete team", err);
      throw err;
    }
  };

  return {
    teams: computed(() => teams.value || []),
    isLoading,
    isError: !!error.value,
    createTeam,
    updateTeam,
    deleteTeam,
    refresh,
  };
};

export const useTeam = (teamId: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: team,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Team>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}`
        : "/api/v1/workspaces/0/teams/0",
    {
      key: `team-${teamId}`,
    }
  );

  return {
    team,
    isLoading,
    isError: !!error.value,
    refresh,
  };
};

export const useTeamMembers = (teamId: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: members,
    pending: isLoading,
    error,
    refresh,
  } = useApi<TeamMember[]>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}/members`
        : "/api/v1/workspaces/0/teams/0/members",
    {
      key: `team-members-${teamId}`,
    }
  );

  const addMember = async (email: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}/members`,
        {
          method: "POST",
          body: { email, role: "MEMBER" },
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to add member", err);
      throw err;
    }
  };

  const removeMember = async (userId: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/teams/${teamId}/members/${userId}`,
        {
          method: "DELETE",
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to remove member", err);
      throw err;
    }
  };

  return {
    members: computed(() => members.value || []),
    isLoading,
    isError: !!error.value,
    addMember,
    removeMember,
    refresh,
  };
};
