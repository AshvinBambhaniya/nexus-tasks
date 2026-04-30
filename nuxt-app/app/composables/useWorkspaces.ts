import type { Workspace } from "~/types";

export const useWorkspaces = () => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: workspacesData,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Workspace[]>("/api/v1/workspaces", {
    key: "workspaces-list",
  });

  const workspaces = computed(() => workspacesData.value || []);

  const activeWorkspace = computed(() =>
    workspaces.value.find((w) => w.id === workspaceStore.activeWorkspaceId)
  );

  // Auto-select first workspace if none selected and data loaded
  watch(
    [isLoading, workspaces],
    ([newIsLoading, newWorkspaces]) => {
      if (
        !newIsLoading &&
        newWorkspaces.length > 0 &&
        !workspaceStore.activeWorkspaceId
      ) {
        workspaceStore.setActiveWorkspaceId(newWorkspaces[0].id);
      }
    },
    { immediate: true }
  );

  const createWorkspace = async (name: string) => {
    try {
      const response = await api$fetch<Workspace>("/api/v1/workspaces/", {
        method: "POST",
        body: { name },
      });
      // api$fetch already unwraps JSend
      const newWorkspace = response as unknown as Workspace;
      await refresh();
      if (newWorkspace && newWorkspace.id) {
        workspaceStore.setActiveWorkspaceId(newWorkspace.id);
      }
      return newWorkspace;
    } catch (err) {
      console.error("Failed to create workspace", err);
      throw err;
    }
  };

  return {
    workspaces,
    activeWorkspace,
    isLoading,
    isError: !!error.value,
    createWorkspace,
    refresh,
  };
};
