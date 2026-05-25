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
  ChevronLeft,
  ChevronRight,
} from "lucide-vue-next";
import { useUsersStore } from "~/stores/user";
import { useUIStore } from "~/stores/ui";

const route = useRoute();
const { logout } = useAuth();
const userStore = useUsersStore();
const uiStore = useUIStore();

const user = computed(() => userStore.userData);
const isLoading = ref(false);

const isCollapsed = computed(() => uiStore.isSidebarCollapsed);

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

const toggleSidebar = () => {
  uiStore.toggleSidebar();
};
</script>

<template>
  <aside
    class="border-border bg-muted/30 relative flex h-full flex-col border-r transition-all duration-300 ease-in-out"
    :class="isCollapsed ? 'w-20' : 'w-64'"
  >
    <!-- User Profile (Top) -->
    <div class="border-border border-b p-4">
      <div
        class="flex items-center"
        :class="isCollapsed ? 'justify-center' : 'justify-between'"
      >
        <div class="flex items-center gap-3 overflow-hidden">
          <div
            class="bg-primary text-primary-foreground flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-semibold shadow-sm"
            :title="isCollapsed ? user?.full_name || user?.email : ''"
          >
            {{ isLoading ? "..." : userInitial }}
          </div>
          <div v-if="!isCollapsed" class="min-w-0">
            <p class="text-foreground truncate text-sm font-semibold">
              {{
                isLoading
                  ? "Loading..."
                  : user?.full_name || user?.email?.split("@")[0]
              }}
            </p>
            <p class="text-muted-foreground truncate text-xs">
              {{ isLoading ? "Please wait" : user?.email }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Workspace Switcher -->
    <div class="p-4">
      <WorkspaceSwitcher :is-collapsed="isCollapsed" />
    </div>

    <!-- Navigation -->
    <div class="flex-1 space-y-6 overflow-y-auto px-3 py-2">
      <nav class="space-y-1">
        <NuxtLink
          v-for="item in navigation"
          :key="item.name"
          :to="item.href"
          class="group flex items-center rounded-md px-3 py-2 text-sm font-medium transition-colors"
          :class="[
            isRouteActive(item.href)
              ? 'bg-primary/10 text-primary'
              : 'text-muted-foreground hover:bg-muted hover:text-foreground',
            isCollapsed ? 'justify-center' : '',
          ]"
          :title="isCollapsed ? item.name : ''"
        >
          <component
            :is="item.icon"
            class="h-5 w-5 shrink-0 transition-colors"
            :class="[
              isRouteActive(item.href)
                ? 'text-primary'
                : 'text-muted-foreground/70 group-hover:text-muted-foreground',
              isCollapsed ? '' : 'mr-3',
            ]"
          />
          <span v-if="!isCollapsed" class="truncate">{{ item.name }}</span>
        </NuxtLink>
      </nav>
    </div>

    <!-- Bottom Actions -->
    <div class="border-border border-t p-4">
      <div
        class="flex items-center"
        :class="isCollapsed ? 'flex-col gap-4' : 'justify-between'"
      >
        <div
          class="flex items-center gap-1"
          :class="isCollapsed ? 'flex-col' : ''"
        >
          <LayoutThemeToggle />
          <button
            class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive rounded-md p-1.5 transition-colors"
            title="Logout"
            @click="logout"
          >
            <LogOut class="h-4 w-4" />
          </button>
        </div>

        <!-- Toggle Button (Desktop Only) -->
        <button
          class="text-muted-foreground hover:bg-muted hover:text-foreground hidden rounded-md p-1.5 transition-colors lg:block"
          :title="isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'"
          @click="toggleSidebar"
        >
          <ChevronLeft v-if="!isCollapsed" class="h-4 w-4" />
          <ChevronRight v-else class="h-4 w-4" />
        </button>
      </div>
    </div>
  </aside>
</template>
