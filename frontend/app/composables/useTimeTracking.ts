import type {
  ActiveTimer,
  TaskTimeEntries,
  ProjectTimeAnalytics,
} from "~/types";

export const useTimeTracking = () => {
  const timerStore = useTimerStore();

  const fetchActiveTimer = async () => {
    try {
      const data = await useMutation<ActiveTimer | null>(
        "/api/v2/timer/active",
        {
          method: "GET",
        }
      );
      if (data) {
        timerStore.setActiveTimer(data);
      } else {
        timerStore.clearActiveTimer();
      }
      return data;
    } catch {
      timerStore.clearActiveTimer();
      return null;
    }
  };

  const startTimer = async (taskId: string) => {
    const data = await useMutation<ActiveTimer>(
      `/api/v2/tasks/${taskId}/timer/start`,
      {
        method: "POST",
      }
    );
    if (data) {
      timerStore.setActiveTimer(data);
    }
    return data;
  };

  const stopTimer = async (
    taskId: string,
    body: { description: string; duration_minutes?: number }
  ) => {
    await useMutation(`/api/v2/tasks/${taskId}/timer/stop`, {
      method: "POST",
      body,
    });
    timerStore.clearActiveTimer();
  };

  const discardTimer = async (taskId: string) => {
    await useMutation(`/api/v2/tasks/${taskId}/timer/discard`, {
      method: "POST",
    });
    timerStore.clearActiveTimer();
  };

  const logManualTime = async (
    taskId: string,
    body: { duration_minutes: number; description: string; date?: string }
  ) => {
    return await useMutation(`/api/v2/tasks/${taskId}/time-entries`, {
      method: "POST",
      body,
    });
  };

  const fetchTaskTimeEntries = (taskId: string) => {
    return useApi<TaskTimeEntries>(`/api/v2/tasks/${taskId}/time-entries`, {
      key: `task-time-entries-${taskId}`,
    });
  };

  const deleteTimeEntry = async (entryId: string) => {
    return await useMutation(`/api/v2/time-entries/${entryId}`, {
      method: "DELETE",
    });
  };

  const fetchProjectAnalytics = (workspaceId: string, projectId: string) => {
    return useApi<ProjectTimeAnalytics>(
      `/api/v2/workspaces/${workspaceId}/projects/${projectId}/time-analytics`,
      { key: `project-time-analytics-${projectId}` }
    );
  };

  const fetchProjectTimeEntries = (
    workspaceId: string,
    projectId: string,
    params?: { user_id?: string; start_date?: string; end_date?: string }
  ) => {
    return useApi<import("~/types").ProjectTimeEntry[]>(
      `/api/v2/workspaces/${workspaceId}/projects/${projectId}/time-entries`,
      { key: `project-time-entries-${projectId}`, query: params }
    );
  };

  return {
    fetchActiveTimer,
    startTimer,
    stopTimer,
    discardTimer,
    logManualTime,
    fetchTaskTimeEntries,
    deleteTimeEntry,
    fetchProjectAnalytics,
    fetchProjectTimeEntries,
  };
};
