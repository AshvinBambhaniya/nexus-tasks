import useSWR from "swr";
import api from "@/lib/api";
import { TaskWithWorkspace } from "@/types";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useMyTasks() {
  const { data, error, isLoading, mutate } = useSWR<TaskWithWorkspace[]>("/api/v1/tasks/me", fetcher);

  return {
    tasks: data || [],
    isLoading,
    isError: error,
    refresh: mutate,
  };
}
