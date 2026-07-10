<script setup lang="ts">
import {
  LayoutDashboard,
  CheckSquare,
  Inbox,
  Settings,
  LogOut,
  Users,
  Folder,
  ChevronLeft,
  ChevronRight,
  ChevronDown,
  ChevronUp,
  BarChart2,
  Bell,
  Check,
} from "lucide-vue-next";
import { useUsersStore } from "~/stores/user";
import { useUIStore } from "~/stores/ui";
import { useWorkspaceStore } from "~/stores/workspace";
import { useProjectStore } from "~/stores/project";

const route = useRoute();
const { logout } = useAuth();
const userStore = useUsersStore();
const uiStore = useUIStore();
const workspaceStore = useWorkspaceStore();
const projectStore = useProjectStore();

const user = computed(() => userStore.userData);
const isUserLoading = ref(false);

const isCollapsed = computed(() => uiStore.isSidebarCollapsed);

// Data Integration
const { workspaces, activeWorkspace } = useWorkspaces();
const { teams, isLoading: isTeamsLoading } = useTeams();
const { projects, isLoading: isProjectsLoading } = useProjects();

const isWorkspaceDropdownOpen = ref(false);
const isTeamsOpen = ref(true);
const isProjectsOpen = ref(true);

const toggleWorkspaceDropdown = () => {
  if (!isCollapsed.value) {
    isWorkspaceDropdownOpen.value = !isWorkspaceDropdownOpen.value;
  }
};

const selectWorkspace = (id: string) => {
  workspaceStore.setActiveWorkspaceId(id);
  isWorkspaceDropdownOpen.value = false;
};

const selectProject = (id: string) => {
  projectStore.setActiveProjectId(id);
};

const navigation = [
  { name: "Inbox", href: "/inbox", icon: Inbox, badge: "3 new" },
  { name: "My Tasks", href: "/tasks", icon: CheckSquare, badge: "5" },
  { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { name: "Analytics", href: "/analytics", icon: BarChart2 },
  { name: "Notifications", href: "/notifications", icon: Bell },
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
  if (uiStore.isSidebarCollapsed) {
    isWorkspaceDropdownOpen.value = false;
  }
};
</script>

<template>
  <aside
    class="border-border bg-background relative z-40 flex h-full flex-col border-r transition-all duration-300 ease-in-out"
    :class="isCollapsed ? 'w-20' : 'w-72'"
  >
    <!-- User Profile Context Switcher -->
    <div class="border-border relative border-b p-4">
      <div
        class="group hover:bg-muted/50 flex cursor-pointer items-center justify-between rounded-md p-2 transition-colors"
        @click="toggleWorkspaceDropdown"
      >
        <div class="flex items-center gap-3 overflow-hidden">
          <div class="relative">
            <div
              class="bg-muted text-muted-foreground flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-semibold"
            >
              {{ isUserLoading ? "..." : userInitial }}
            </div>
            <!-- Online status dot -->
            <div
              class="border-background absolute right-0 bottom-0 h-2.5 w-2.5 rounded-full border-2 bg-emerald-500"
            />
          </div>
          <div v-if="!isCollapsed" class="min-w-0">
            <p class="text-foreground truncate text-sm font-medium">
              {{ activeWorkspace?.name || "Loading..." }}
            </p>
            <p class="text-muted-foreground truncate text-xs">
              {{ user?.full_name || user?.email?.split("@")[0] }}
            </p>
          </div>
        </div>
        <ChevronDown
          v-if="!isCollapsed"
          class="text-muted-foreground h-4 w-4 opacity-50 group-hover:opacity-100"
        />
      </div>

      <!-- Workspace Dropdown -->
      <div
        v-if="isWorkspaceDropdownOpen && !isCollapsed"
        class="border-border bg-card absolute top-full right-4 left-4 z-50 mt-1 rounded-md border py-1 shadow-lg"
      >
        <div
          class="text-muted-foreground px-2 py-1 text-xs font-semibold tracking-wider uppercase"
        >
          Switch Workspace
        </div>
        <button
          v-for="ws in workspaces"
          :key="ws.id"
          class="text-foreground hover:bg-muted/50 flex w-full items-center justify-between px-3 py-2 text-sm transition-colors"
          @click="selectWorkspace(ws.id)"
        >
          <span class="truncate">{{ ws.name }}</span>
          <Check
            v-if="ws.id === workspaceStore.activeWorkspaceId"
            class="text-primary h-4 w-4"
          />
        </button>
      </div>
    </div>

    <!-- Navigation -->
    <div class="custom-scrollbar flex-1 space-y-6 overflow-y-auto px-3 py-4">
      <nav class="space-y-0.5">
        <NuxtLink
          v-for="item in navigation"
          :key="item.name"
          :to="item.href"
          class="group relative flex items-center justify-between px-3 py-2 text-sm font-medium transition-colors"
          :class="[
            isRouteActive(item.href)
              ? 'bg-muted/30 text-primary border-primary -ml-[2px] border-l-2 pl-[14px]'
              : 'text-muted-foreground hover:bg-muted/30 hover:text-foreground -ml-[2px] border-l-2 border-transparent pl-[14px]',
            isCollapsed
              ? '!ml-0 justify-center rounded-md border-l-0 !pl-0'
              : '',
          ]"
          :title="isCollapsed ? item.name : ''"
        >
          <div class="flex items-center">
            <component
              :is="item.icon"
              class="h-4 w-4 shrink-0 transition-colors"
              :class="[
                isRouteActive(item.href)
                  ? 'text-primary'
                  : 'text-muted-foreground/70 group-hover:text-foreground',
                isCollapsed ? '' : 'mr-3',
              ]"
            />
            <span v-if="!isCollapsed" class="truncate">{{ item.name }}</span>
          </div>
          <span
            v-if="!isCollapsed && item.badge"
            class="bg-muted text-muted-foreground flex h-5 items-center rounded-full px-2 text-[10px] font-semibold"
          >
            {{ item.badge }}
          </span>
        </NuxtLink>
      </nav>

      <!-- TEAMS Tree -->
      <div v-if="!isCollapsed" class="space-y-1">
        <div
          class="group flex cursor-pointer items-center justify-between px-3 py-2"
          @click="isTeamsOpen = !isTeamsOpen"
        >
          <h3
            class="text-muted-foreground text-[10px] font-semibold tracking-widest"
          >
            TEAMS
          </h3>
          <ChevronUp
            v-if="isTeamsOpen"
            class="text-muted-foreground/50 group-hover:text-foreground h-3 w-3"
          />
          <ChevronDown
            v-else
            class="text-muted-foreground/50 group-hover:text-foreground h-3 w-3"
          />
        </div>

        <div v-if="isTeamsOpen" class="space-y-0.5 px-2">
          <div
            v-if="isTeamsLoading"
            class="flex animate-pulse items-center gap-2 px-2 py-1.5"
          >
            <div class="bg-muted h-4 w-4 rounded" />
            <div class="bg-muted h-3 w-20 rounded" />
          </div>
          <div
            v-else-if="teams.length === 0"
            class="text-muted-foreground/50 px-2 py-1.5 text-xs"
          >
            No teams found
          </div>
          <div v-for="team in teams" v-else :key="team.id" class="space-y-0.5">
            <NuxtLink
              :to="`/teams/${team.id}`"
              class="text-muted-foreground hover:bg-muted/30 hover:text-foreground flex cursor-pointer items-center justify-between rounded-md px-2 py-1.5 text-sm transition-colors"
            >
              <div class="flex items-center gap-2">
                <Users class="text-muted-foreground/70 h-4 w-4" />
                <span class="truncate">{{ team.name }}</span>
              </div>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- PROJECTS Tree -->
      <div v-if="!isCollapsed" class="space-y-1">
        <div
          class="group flex cursor-pointer items-center justify-between px-3 py-2"
          @click="isProjectsOpen = !isProjectsOpen"
        >
          <h3
            class="text-muted-foreground text-[10px] font-semibold tracking-widest"
          >
            PROJECTS
          </h3>
          <ChevronUp
            v-if="isProjectsOpen"
            class="text-muted-foreground/50 group-hover:text-foreground h-3 w-3"
          />
          <ChevronDown
            v-else
            class="text-muted-foreground/50 group-hover:text-foreground h-3 w-3"
          />
        </div>

        <div v-if="isProjectsOpen" class="space-y-0.5 px-2">
          <div
            v-if="isProjectsLoading"
            class="flex animate-pulse items-center gap-2 px-2 py-1.5"
          >
            <div class="bg-muted h-4 w-4 rounded" />
            <div class="bg-muted h-3 w-24 rounded" />
          </div>
          <div
            v-else-if="projects.length === 0"
            class="text-muted-foreground/50 px-2 py-1.5 text-xs"
          >
            No projects found
          </div>
          <NuxtLink
            v-for="project in projects"
            v-else
            :key="project.id"
            :to="`/dashboard?projectId=${project.id}`"
            class="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm transition-colors"
            :class="
              projectStore.activeProjectId === project.id
                ? 'bg-muted/30 text-foreground'
                : 'text-muted-foreground hover:bg-muted/30 hover:text-foreground'
            "
            @click="selectProject(project.id)"
          >
            <Folder
              class="h-4 w-4"
              :class="
                projectStore.activeProjectId === project.id
                  ? 'text-primary'
                  : 'text-muted-foreground/70'
              "
            />
            <span class="truncate">{{ project.name }}</span>
          </NuxtLink>
        </div>
      </div>
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
          <NuxtLink
            to="/settings"
            class="text-muted-foreground hover:bg-muted hover:text-foreground rounded-md p-1.5 transition-colors"
            title="Settings"
          >
            <Settings class="h-4 w-4" />
          </NuxtLink>
          <button
            class="text-muted-foreground hover:bg-destructive/10 hover:text-destructive rounded-md p-1.5 transition-colors"
            title="Logout"
            @click="logout"
          >
            <LogOut class="h-4 w-4" />
          </button>
        </div>

        <!-- Toggle Button -->
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

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(150, 150, 150, 0.2);
  border-radius: 4px;
}
.custom-scrollbar:hover::-webkit-scrollbar-thumb {
  background: rgba(150, 150, 150, 0.4);
}
</style>
