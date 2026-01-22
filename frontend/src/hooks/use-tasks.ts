import useSWR, { useSWRConfig } from "swr";
import api from "@/lib/api";
import { Task, TaskStatus, TaskPriority } from "@/types";
import { useWorkspaceStore } from "@/store/workspace-store";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useTasks() {
  const { activeWorkspaceId } = useWorkspaceStore();
  const { mutate } = useSWRConfig();
  
  const key = activeWorkspaceId ? `/api/v1/workspaces/${activeWorkspaceId}/tasks` : null;
  
  const { data, error, isLoading } = useSWR<Task[]>(key, fetcher);

  const createTask = async (task: { title: string; description?: string; priority?: TaskPriority; status?: TaskStatus }) => {
    if (!activeWorkspaceId) return;
    try {
      await api.post(`/api/v1/workspaces/${activeWorkspaceId}/tasks`, task);
      mutate(key);
    } catch (err) {
      console.error("Failed to create task", err);
      throw err;
    }
  };

  const updateTask = async (taskId: number, updates: Partial<Task>) => {
    try {
      // Optimistic update could be added here
      await api.patch(`/api/v1/tasks/${taskId}`, updates);
      mutate(key);
    } catch (err) {
      console.error("Failed to update task", err);
      throw err;
    }
  };

  const deleteTask = async (taskId: number) => {
    try {
      await api.delete(`/api/v1/tasks/${taskId}`);
      mutate(key);
    } catch (err) {
      console.error("Failed to delete task", err);
      throw err;
    }
  };

  return {
    tasks: data || [],
    isLoading,
    isError: error,
    createTask,
    updateTask,
    deleteTask,
  };
}
