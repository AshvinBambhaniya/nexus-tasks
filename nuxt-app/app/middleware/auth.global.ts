export default defineNuxtRouteMiddleware((to) => {
  const token = useCookie("access_token").value;
  const { path } = to;

  // Redirect authenticated users from login/register to dashboard
  if (token && (path === "/login" || path === "/register")) {
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

  const isProtectedRoute = protectedPrefixes.some((route) =>
    path.startsWith(route)
  );

  if (!token && isProtectedRoute) {
    return navigateTo("/login");
  }
});
