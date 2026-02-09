import useSWR from "swr";
import api from "@/lib/api";
import { TaskWithProject, Task } from "@/types";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useMyTasks() {
  const { data, error, isLoading, mutate } = useSWR<TaskWithProject[]>(
    "/api/v1/tasks/me",
    fetcher
  );

  const updateTask = async (taskId: number, updates: Partial<Task>) => {
    try {
      await api.patch(`/api/v1/tasks/${taskId}`, updates);
      mutate(); // Refresh the list
    } catch (err) {
      console.error("Failed to update task", err);
      throw err;
    }
  };

  return {
    tasks: data || [],
    isLoading,
    isError: error,
    refresh: mutate,
    updateTask,
  };
}
