import useSWR, { useSWRConfig } from "swr";
import api from "@/lib/api";
import { Workspace } from "@/types";
import { useWorkspaceStore } from "@/store/workspace-store";
import { useEffect } from "react";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useWorkspaces() {
  const { mutate } = useSWRConfig();
  const { data, error, isLoading } = useSWR<Workspace[]>("/api/v1/workspaces", fetcher);
  const { activeWorkspaceId, setActiveWorkspaceId } = useWorkspaceStore();

  const workspaces = data || [];
  const activeWorkspace = workspaces.find((w) => w.id === activeWorkspaceId);

  // Auto-select first workspace if none selected and data loaded
  useEffect(() => {
    if (!isLoading && workspaces.length > 0 && !activeWorkspaceId) {
      setActiveWorkspaceId(workspaces[0].id);
    }
  }, [isLoading, workspaces, activeWorkspaceId, setActiveWorkspaceId]);

  const createWorkspace = async (name: string) => {
    try {
      const response = await api.post("/api/v1/workspaces/", { name });
      const newWorkspace = response.data;
      mutate("/api/v1/workspaces");
      setActiveWorkspaceId(newWorkspace.id);
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
    isError: error,
    createWorkspace,
  };
}
