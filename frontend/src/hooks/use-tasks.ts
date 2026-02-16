import useSWR, { useSWRConfig } from "swr";
import api from "@/lib/api";
import { Task, TaskStatus, TaskPriority, Comment } from "@/types";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useTask(taskId: number) {
  const {
    data: task,
    error: taskError,
    isLoading: taskLoading,
    mutate: mutateTask,
  } = useSWR<Task>(taskId ? `/api/v1/tasks/${taskId}` : null, fetcher);

  const {
    data: comments,
    error: commentsError,
    isLoading: commentsLoading,
    mutate: mutateComments,
  } = useSWR<Comment[]>(
    taskId ? `/api/v1/tasks/${taskId}/comments` : null,
    fetcher
  );

  const createComment = async (content: string) => {
    try {
      await api.post(`/api/v1/tasks/${taskId}/comments`, { content });
      mutateComments();
    } catch (err) {
      console.error("Failed to create comment", err);
      throw err;
    }
  };

  const deleteComment = async (commentId: number) => {
    try {
      await api.delete(`/api/v1/comments/${commentId}`);
      mutateComments();
    } catch (err) {
      console.error("Failed to delete comment", err);
      throw err;
    }
  };

  return {
    task,
    comments: comments || [],
    isLoading: taskLoading || commentsLoading,
    isError: taskError || commentsError,
    mutateTask,
    mutateComments,
    createComment,
    deleteComment,
  };
}

export function useTasks(projectId?: number) {
  const { mutate } = useSWRConfig();

  const key = projectId ? `/api/v1/projects/${projectId}/tasks` : null;

  const { data, error, isLoading } = useSWR<Task[]>(key, fetcher);

  const createTask = async (task: {
    title: string;
    description?: string;
    priority?: TaskPriority;
    status?: TaskStatus;
    assignee_id?: number;
    due_date?: string;
  }) => {
    if (!projectId) return;
    try {
      await api.post(`/api/v1/projects/${projectId}/tasks`, task);
      mutate(key);
    } catch (err) {
      console.error("Failed to create task", err);
      throw err;
    }
  };

  const updateTask = async (taskId: number, updates: Partial<Task>) => {
    try {
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
