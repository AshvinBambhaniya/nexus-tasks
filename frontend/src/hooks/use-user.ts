import useSWR from "swr";
import api from "@/lib/api";
import { User } from "@/types";

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export function useUser() {
  const { data, error, isLoading, mutate } = useSWR<User>(
    "/api/v1/auth/me",
    fetcher
  );

  const updateUser = async (updates: Partial<User>) => {
    try {
      const response = await api.patch("/api/v1/auth/me", updates);
      mutate(response.data);
      return response.data;
    } catch (err) {
      console.error("Failed to update user", err);
      throw err;
    }
  };

  return {
    user: data,
    isLoading,
    isError: error,
    updateUser,
  };
}
