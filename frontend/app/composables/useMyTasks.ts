import type { Task, TaskWithProject } from "~/types";

export const useMyTasks = () => {
  const {
    data: tasks,
    pending: isLoading,
    error,
    refresh,
  } = useApi<TaskWithProject[]>("/api/v1/tasks/me", {
    key: "my-tasks",
  });

  const updateTask = async (taskId: string, updates: Partial<Task>) => {
    try {
      await useMutation(`/api/v1/tasks/${taskId}`, {
        method: "PATCH",
        body: updates,
      });
      await refresh();
    } catch (err) {
      console.error("Failed to update task", err);
      throw err;
    }
  };

  return {
    tasks: computed(() => tasks.value || []),
    isLoading,
    isError: !!error.value,
    refresh,
    updateTask,
  };
};
