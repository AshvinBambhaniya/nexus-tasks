<script setup lang="ts">
import {
  LayoutDashboard,
  CheckSquare,
  Inbox,
  Settings,
  LogOut,
  Users,
  Kanban,
  Folder,
} from "lucide-vue-next";
import { useUsersStore } from "~/stores/user";

const route = useRoute();
const { logout } = useAuth();
const userStore = useUsersStore();

const user = computed(() => userStore.userData);
// We can track global loading in the store if needed,
// for now we'll assume false after initial layout fetch
const isLoading = ref(false);

const navigation = [
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { name: "My Focus", href: "/inbox", icon: Inbox },
  { name: "Projects", href: "/projects", icon: Folder },
  { name: "Teams", href: "/teams", icon: Users },
  { name: "Boards", href: "/boards", icon: Kanban },
  { name: "All Tasks", href: "/tasks", icon: CheckSquare },
  { name: "Settings", href: "/settings", icon: Settings },
];

const userInitial = computed(() => {
  if (user.value?.full_name) {
    return user.value.full_name
      .split(" ")
      .map((n) => n[0])
      .join("")
      .substring(0, 2)
      .toUpperCase();
  }
  return user.value?.email?.substring(0, 2).toUpperCase() || "??";
});

const isRouteActive = (href: string) => {
  return (
    route.path === href ||
    (href !== "/dashboard" && route.path.startsWith(href))
  );
};
</script>

<template>
  <div class="flex h-full w-64 flex-col border-r border-gray-200 bg-gray-50/50">
    <div class="p-4">
      <WorkspaceSwitcher />
    </div>

    <div class="flex-1 space-y-6 overflow-y-auto px-3 py-2">
      <nav class="space-y-1">
        <NuxtLink
          v-for="item in navigation"
          :key="item.name"
          :to="item.href"
          class="group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors"
          :class="
            isRouteActive(item.href)
              ? 'bg-blue-50 text-blue-600'
              : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
          "
        >
          <component
            :is="item.icon"
            class="mr-3 h-5 w-5 flex-shrink-0 transition-colors"
            :class="
              isRouteActive(item.href)
                ? 'text-blue-600'
                : 'text-gray-400 group-hover:text-gray-500'
            "
          />
          {{ item.name }}
        </NuxtLink>
      </nav>
    </div>

    <div class="border-t border-gray-200 p-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div
            class="flex h-8 w-8 items-center justify-center rounded-full bg-blue-600 text-xs font-semibold text-white"
          >
            {{ isLoading ? "..." : userInitial }}
          </div>
          <div class="text-sm">
            <p class="font-medium text-gray-700">
              {{
                isLoading
                  ? "Loading..."
                  : user?.full_name || user?.email?.split("@")[0]
              }}
            </p>
            <p class="text-xs text-gray-500">
              {{ isLoading ? "Please wait" : user?.email }}
            </p>
          </div>
        </div>
        <button
          class="rounded-md p-1.5 text-gray-400 transition-colors hover:bg-red-50 hover:text-red-600"
          title="Logout"
          @click="logout"
        >
          <LogOut class="h-4 w-4" />
        </button>
      </div>
    </div>
  </div>
</template>
