import type { Task, Comment, TaskPriority, TaskStatus } from "~/types";

export const useTasks = (projectId?: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: tasks,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Task[]>(
    () =>
      projectId && workspaceStore.activeWorkspaceId
        ? `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks`
        : "",
    {
      key: `tasks-list-${projectId}`,
      watch: [() => projectId, () => workspaceStore.activeWorkspaceId],
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
    if (!projectId || !workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks`,
        {
          method: "POST",
          body: task,
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to create task", err);
      throw err;
    }
  };

  const updateTask = async (taskId: string, updates: Partial<Task>) => {
    if (!projectId || !workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks/${taskId}`,
        {
          method: "PATCH",
          body: updates,
        }
      );
      await refresh();
    } catch (err) {
      console.error("Failed to update task", err);
      throw err;
    }
  };

  const deleteTask = async (taskId: string) => {
    if (!projectId || !workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks/${taskId}`,
        {
          method: "DELETE",
        }
      );
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

export const useMyTasks = () => {
  const {
    data: tasks,
    pending: isLoading,
    error,
    refresh,
  } = useApi<Task[]>("/api/v2/tasks/me", {
    key: "my-tasks-list",
  });

  return {
    tasks: computed(() => tasks.value || []),
    isLoading,
    isError: !!error.value,
    refresh,
  };
};

export const useTask = (projectId: string, taskId: string) => {
  const workspaceStore = useWorkspaceStore();

  const {
    data: task,
    pending: taskLoading,
    error: taskError,
    refresh: refreshTask,
  } = useApi<Task>(
    () =>
      projectId
        ? `/api/v2/workspaces/${workspaceStore.activeWorkspaceId || "0"}/projects/${projectId}/tasks/${taskId}`
        : `/api/v2/tasks/${taskId}`,
    {
      key: `task-${taskId}`,
    }
  );

  const {
    data: comments,
    pending: commentsLoading,
    error: commentsError,
    refresh: refreshComments,
  } = useApi<Comment[]>(
    () =>
      projectId
        ? `/api/v2/workspaces/${workspaceStore.activeWorkspaceId || "0"}/projects/${projectId}/tasks/${taskId}/comments`
        : `/api/v2/tasks/${taskId}/comments`,
    {
      key: `task-comments-${taskId}`,
    }
  );

  const createComment = async (
    content: string,
    mentionedUserIds: string[] = []
  ) => {
    try {
      const endpoint = projectId
        ? `/api/v2/workspaces/${workspaceStore.activeWorkspaceId || "0"}/projects/${projectId}/tasks/${taskId}/comments`
        : `/api/v2/tasks/${taskId}/comments`;

      await useMutation(endpoint, {
        method: "POST",
        body: { content, mentioned_user_ids: mentionedUserIds },
      });
      await refreshComments();
    } catch (err) {
      console.error("Failed to create comment", err);
      throw err;
    }
  };

  const deleteComment = async (commentId: string) => {
    if (!workspaceStore.activeWorkspaceId) return;
    try {
      await useMutation(
        `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks/${taskId}/comments/${commentId}`,
        {
          method: "DELETE",
        }
      );
      await refreshComments();
    } catch (err) {
      console.error("Failed to delete comment", err);
      throw err;
    }
  };

  const summarizeComments = async () => {
    if (!workspaceStore.activeWorkspaceId) return null;
    try {
      const res = await useMutation<{ content: string }>(
        `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${projectId}/tasks/${taskId}/ai/summarize-comments`,
        {
          method: "POST",
        }
      );
      return res?.content || null;
    } catch (err) {
      console.error("Failed to summarize comments", err);
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
    summarizeComments,
  };
};
