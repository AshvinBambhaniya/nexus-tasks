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

export const useProject = (id: number) => {
  const {
    data: project,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Project>(`/api/v1/projects/${id}`, {
    key: `project-${id}`,
  });

  return {
    project,
    isLoading,
    isError: !!error.value,
    refresh,
  };
};

export const useProjectMembers = (projectId: number) => {
  const {
    data: members,
    pending: isLoading,
    error,
    refresh,
  } = useApi<ProjectMember[]>(`/api/v1/projects/${projectId}/members`, {
    key: `project-members-${projectId}`,
  });

  const addMember = async (email: string) => {
    try {
      await useMutation(`/api/v1/projects/${projectId}/members`, {
        method: "POST",
        body: { email, role: "MEMBER" },
      });
      await refresh();
    } catch (err) {
      console.error("Failed to add member", err);
      throw err;
    }
  };

  const removeMember = async (userId: number) => {
    try {
      await useMutation(`/api/v1/projects/${projectId}/members/${userId}`, {
        method: "DELETE",
      });
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

export const useProjectTeams = (projectId: number) => {
  const {
    data: teams,
    pending: isLoading,
    error,
    refresh,
  } = useApi<ProjectTeam[]>(`/api/v1/projects/${projectId}/teams`, {
    key: `project-teams-${projectId}`,
  });

  const addTeam = async (teamId: number) => {
    try {
      await useMutation(`/api/v1/projects/${projectId}/teams`, {
        method: "POST",
        body: { team_id: teamId },
      });
      await refresh();
    } catch (err) {
      console.error("Failed to add team", err);
      throw err;
    }
  };

  const removeTeam = async (teamId: number) => {
    try {
      await useMutation(`/api/v1/projects/${projectId}/teams/${teamId}`, {
        method: "DELETE",
      });
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
