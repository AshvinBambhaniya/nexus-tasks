"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { useSWRConfig } from "swr";

export function useAuth() {
  const router = useRouter();
  const { mutate } = useSWRConfig();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const login = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      const formData = new FormData();
      formData.append("username", email);
      formData.append("password", password);

      await api.post("/api/v1/auth/login", formData, {
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
      });

      // Token is now set in HttpOnly cookie by backend

      router.push("/dashboard");
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.detail || "Failed to login");
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (
    email: string,
    password: string,
    fullName?: string
  ) => {
    setIsLoading(true);
    setError(null);
    try {
      await api.post("/api/v1/auth/register", {
        email,
        password,
        full_name: fullName,
      });
      // Auto-login after register
      await login(email, password);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (err: any) {
      setError(err.response?.data?.detail || "Failed to register");
      setIsLoading(false);
    }
  };

  const logout = async () => {
    try {
      await api.post("/api/v1/auth/logout");
    } catch (error) {
      console.error("Logout failed", error);
    }

    localStorage.removeItem("workspace-storage");
    mutate(() => true, undefined, { revalidate: false });
    router.push("/login");
  };

  return { login, register, logout, isLoading, error };
}
