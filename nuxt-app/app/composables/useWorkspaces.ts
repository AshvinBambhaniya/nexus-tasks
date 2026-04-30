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
  watchEffect(() => {
    if (
      !isLoading.value &&
      workspaces.value.length > 0 &&
      !workspaceStore.activeWorkspaceId
    ) {
      const firstWorkspace = workspaces.value[0];
      if (firstWorkspace) {
        workspaceStore.setActiveWorkspaceId(firstWorkspace.id);
      }
    }
  });

  const createWorkspace = async (name: string) => {
    try {
      const response = await useMutation<Workspace>("/api/v1/workspaces/", {
        method: "POST",
        body: { name },
      });
      // useMutation already unwraps JSend
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
