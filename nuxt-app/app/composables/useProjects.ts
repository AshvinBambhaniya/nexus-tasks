import type { Project, ProjectMember, ProjectTeam } from "~/types";

export const useProjects = () => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: projects,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Project[]>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects`
        : "/api/v1/workspaces/0/projects",
    {
      key: `projects-list-${workspaceStore.activeWorkspaceId}`,
      watch: [() => workspaceStore.activeWorkspaceId],
    }
  );

  const createProject = async (name: string, description?: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects`,
        {
          method: "POST",
          body: { name, description },
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to create project", err);
      throw err;
    }
  };

  return {
    projects: computed(() => projects.value || []),
    isLoading,
    isError: !!error.value,
    createProject,
    refresh,
  };
};

export const useProject = (id: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: project,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Project>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${id}`
        : `/api/v1/workspaces/0/projects/${id}`,
    {
      key: `project-${id}`,
    }
  );

  return {
    project,
    isLoading,
    isError: !!error.value,
    refresh,
  };
};

export const useProjectMembers = (projectId: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: members,
    pending: isLoading,
    error,
    refresh,
  } = useApi<ProjectMember[]>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/members`
        : `/api/v1/workspaces/0/projects/${projectId}/members`,
    {
      key: `project-members-${projectId}`,
    }
  );

  const addMember = async (email: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/members`,
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
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/members/${userId}`,
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

export const useProjectTeams = (projectId: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: teams,
    pending: isLoading,
    error,
    refresh,
  } = useApi<ProjectTeam[]>(
    () =>
      workspaceStore.activeWorkspaceId
        ? `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/teams`
        : `/api/v1/workspaces/0/projects/${projectId}/teams`,
    {
      key: `project-teams-${projectId}`,
    }
  );

  const addTeam = async (teamId: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/teams`,
        {
          method: "POST",
          body: { team_id: teamId },
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to add team", err);
      throw err;
    }
  };

  const removeTeam = async (teamId: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v1/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/teams/${teamId}`,
        {
          method: "DELETE",
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to remove team", err);
      throw err;
    }
  };

  return {
    teams: computed(() => teams.value || []),
    isLoading,
    isError: !!error.value,
    addTeam,
    removeTeam,
    refresh,
  };
};
