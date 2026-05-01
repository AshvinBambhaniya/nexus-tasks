import type { Task, Comment, TaskPriority, TaskStatus } from "~/types";

export const useTasks = (projectId?: string) => {
  const {
    data: tasks,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Task[]>(
    () => (projectId ? `/api/v1/projects/${projectId}/tasks` : null),
    {
      key: `tasks-list-${projectId}`,
      watch: [() => projectId],
    }
  );

  const createTask = async (task: {
    title: string;
    description?: string;
    priority?: TaskPriority;
    status?: TaskStatus;
    assignee_id?: string;
    due_date?: string;
  }) => {
    if (!projectId) return;
    try {
      await useMutation(`/api/v1/projects/${projectId}/tasks`, {
        method: "POST",
        body: task,
      });
      await refresh();
    } catch (err) {
      console.error("Failed to create task", err);
      throw err;
    }
  };

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

  const deleteTask = async (taskId: string) => {
    try {
      await useMutation(`/api/v1/tasks/${taskId}`, {
        method: "DELETE",
      });
      await refresh();
    } catch (err) {
      console.error("Failed to delete task", err);
      throw err;
    }
  };

  return {
    tasks: computed(() => tasks.value || []),
    isLoading,
    isError: !!error.value,
    createTask,
    updateTask,
    deleteTask,
    refresh,
  };
};

export const useTask = (taskId: string) => {
  const {
    data: task,
    pending: taskLoading,
    error: taskError,
    refresh: refreshTask,
  } = useApi<Task>(`/api/v1/tasks/${taskId}`, {
    key: `task-${taskId}`,
  });

  const {
    data: comments,
    pending: commentsLoading,
    error: commentsError,
    refresh: refreshComments,
  } = useApi<Comment[]>(`/api/v1/tasks/${taskId}/comments`, {
    key: `task-comments-${taskId}`,
  });

  const createComment = async (content: string) => {
    try {
      await useMutation(`/api/v1/tasks/${taskId}/comments`, {
        method: "POST",
        body: { content },
      });
      await refreshComments();
    } catch (err) {
      console.error("Failed to create comment", err);
      throw err;
    }
  };

  const deleteComment = async (commentId: string) => {
    try {
      await useMutation(`/api/v1/comments/${commentId}`, {
        method: "DELETE",
      });
      await refreshComments();
    } catch (err) {
      console.error("Failed to delete comment", err);
      throw err;
    }
  };

  return {
    task,
    comments: computed(() => comments.value || []),
    isLoading: computed(() => taskLoading.value || commentsLoading.value),
    isError: computed(() => !!taskError.value || !!commentsError.value),
    refreshTask,
    refreshComments,
    createComment,
    deleteComment,
  };
};
