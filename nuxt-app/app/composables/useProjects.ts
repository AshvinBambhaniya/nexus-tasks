import type { Project } from "~/types";

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
        : "/api/v1/workspaces/0/projects", // Fallback to avoid null in reactive url
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
