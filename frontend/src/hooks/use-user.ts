import useSWR from "swr";
import api from "@/lib/api";
import { User } from "@/types";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useUser() {
  const { data, error, isLoading } = useSWR<User>("/api/v1/auth/me", fetcher);

  return {
    user: data,
    isLoading,
    isError: error,
  };
}
