export const useAuth = () => {
  const router = useRouter();
  const workspaceStore = useWorkspaceStore();
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const login = async (email: string, password: string) => {
    isLoading.value = true;
    error.value = null;
    try {
      await api$fetch("/api/v1/auth/login", {
        method: "POST",
        body: { email, password },
      });
      // Clear user data cache if necessary
      refreshCookie("current-user");
      router.push("/dashboard");
    } catch (err: unknown) {
      error.value = getApiErrorMessage(err, "Failed to login");
    } finally {
      isLoading.value = false;
    }
  };

  const register = async (
    email: string,
    password: string,
    fullName?: string
  ) => {
    isLoading.value = true;
    error.value = null;
    try {
      await api$fetch("/api/v1/auth/register", {
        method: "POST",
        body: { email, password, full_name: fullName },
      });
      // Auto-login or redirect
      await login(email, password);
    } catch (err: unknown) {
      error.value = getApiErrorMessage(err, "Failed to register");
      isLoading.value = false;
    }
  };

  const logout = async () => {
    try {
      await api$fetch("/api/v1/auth/logout", { method: "POST" });
    } catch (err) {
      console.error("Logout failed", err);
    }
    // Clear state
    workspaceStore.$reset();
    // Redirect to login
    router.push("/login");
  };

  return {
    login,
    register,
    logout,
    isLoading,
    error,
  };
};
