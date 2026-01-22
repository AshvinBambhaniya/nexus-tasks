"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import api from "@/lib/api";
import { useSWRConfig } from "swr";
import Cookies from "js-cookie";

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

      const res = await api.post("/api/v1/auth/login", formData, {
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
      });

      // Store token in cookie (expires in 7 days)
      Cookies.set("token", res.data.access_token, { expires: 7, secure: true, sameSite: "strict" });
      
      router.push("/dashboard");
    } catch (err: any) {
      setError(err.response?.data?.detail || "Failed to login");
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (email: string, password: string) => {
    setIsLoading(true);
    setError(null);
    try {
      await api.post("/api/v1/auth/register", { email, password });
      // Auto-login after register
      await login(email, password);
    } catch (err: any) {
      setError(err.response?.data?.detail || "Failed to register");
      setIsLoading(false);
    }
  };

  const logout = () => {
    Cookies.remove("token");
    localStorage.removeItem("workspace-storage");
    mutate(() => true, undefined, { revalidate: false });
    router.push("/login");
  };

  return { login, register, logout, isLoading, error };
}
