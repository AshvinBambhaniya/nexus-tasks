import { useUsersStore } from "~/stores/user";

export default defineNuxtRouteMiddleware(async (to) => {
  const userStore = useUsersStore();
  const { path } = to;

  // Define public/auth routes
  const isAuthRoute = path === "/login" || path === "/register";

  // If we don't have user data yet, try to fetch it 
  // (this helps during direct entry into a dashboard route)
  if (!userStore.userData && import.meta.client) {
    await setUserDataStore();
  }

  const isAuthenticated = !!userStore.userData;

  // Redirect authenticated users from login/register to dashboard
  if (isAuthenticated && isAuthRoute) {
    return navigateTo("/dashboard");
  }

  // Define protected routes
  const protectedPrefixes = [
    "/dashboard",
    "/inbox",
    "/boards",
    "/tasks",
    "/settings",
    "/projects",
    "/teams",
  ];

  const isProtectedRoute = protectedPrefixes.some((prefix) =>
    path.startsWith(prefix)
  );

  if (!isAuthenticated && isProtectedRoute) {
    // Before kicking them out, double check if they have a cookie
    const token = useCookie("access_token").value;
    if (!token) {
      return navigateTo("/login");
    }
  }
});
