<script setup lang="ts">
import {
  CheckSquare,
  Plus,
  Activity,
  Briefcase,
  GitPullRequest,
  MoreHorizontal,
  Flame,
  Zap,
} from "lucide-vue-next";
import { format } from "date-fns";
import { TaskStatus } from "~/types";
import { useUsersStore } from "~/stores/user";

definePageMeta({
  layout: "dashboard",
});

const userStore = useUsersStore();
const user = computed(() => userStore.userData);
const { activeWorkspace } = useWorkspaces();
const { projects, isLoading: projectsLoading } = useProjects();
const { tasks: myTasks, isLoading: tasksLoading } = useMyTasks();

const isTaskDialogOpen = ref(false);

const workspaceProjectIds = computed(
  () => new Set(projects.value.map((p) => p.id))
);
const workspaceTasks = computed(() =>
  myTasks.value.filter((t) => workspaceProjectIds.value.has(t.project_id))
);

const pendingTasksCount = computed(
  () => workspaceTasks.value.filter((t) => t.status !== TaskStatus.DONE).length
);

const recentTasks = computed(() =>
  [...workspaceTasks.value]
    .filter((t) => t.status !== TaskStatus.DONE)
    .sort((a, b) => {
      if (a.priority !== b.priority)
        return a.priority.localeCompare(b.priority);
      return (
        new Date(a.due_date || "2100-01-01").getTime() -
        new Date(b.due_date || "2100-01-01").getTime()
      );
    })
    .slice(0, 7)
);

const getStatusColor = (status: string) => {
  switch (status) {
    case "TODO":
      return "bg-slate-500/10 text-slate-600 dark:text-slate-400 border border-slate-500/20";
    case "IN_PROGRESS":
      return "bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/20";
    case "IN_REVIEW":
      return "bg-purple-500/10 text-purple-600 dark:text-purple-400 border border-purple-500/20";
    case "DONE":
      return "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20";
    default:
      return "bg-muted text-muted-foreground border border-border";
  }
};
const getStatusLabel = (status: string) => status.replace("_", " ");
</script>

<template>
  <div class="mx-auto max-w-[1600px] space-y-6">
    <!-- Header -->
    <div
      class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <h1 class="text-foreground text-2xl font-bold tracking-tight">
          Dashboard
        </h1>
        <p class="text-muted-foreground mt-1 text-sm">
          Overview for
          <span class="text-foreground font-medium">{{
            activeWorkspace?.name || "Workspace"
          }}</span>
        </p>
      </div>
      <div class="flex items-center gap-4">
        <!-- Optional global search input here -->
        <button
          class="bg-primary text-primary-foreground hover:bg-primary/90 inline-flex h-9 items-center justify-center rounded-md px-4 py-2 text-sm font-medium shadow transition-colors disabled:opacity-50"
          :disabled="projects.length === 0"
          @click="isTaskDialogOpen = true"
        >
          <Plus class="mr-2 h-4 w-4" /> New Issue
        </button>
      </div>
    </div>

    <!-- Stats Grid (Bento Top Row) -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <DashboardStatsCard
        title="Active Projects"
        :value="projects.filter((p) => !p.is_archived).length"
        :icon="Activity"
        :is-loading="projectsLoading"
        color="green"
        trend="+3"
      />
      <DashboardStatsCard
        title="Pending Tasks"
        :value="pendingTasksCount"
        :icon="CheckSquare"
        :is-loading="tasksLoading"
        color="orange"
      />
      <DashboardStatsCard
        title="Pull Requests"
        value="8"
        :icon="GitPullRequest"
        :is-loading="false"
        color="blue"
        label="OPEN PRs"
      />
      <DashboardStatsCard
        title="Team Velocity"
        value="94%"
        :icon="Flame"
        :is-loading="false"
        color="purple"
      />
    </div>

    <div class="mt-4 grid grid-cols-1 gap-4 lg:grid-cols-12">
      <!-- Main Column: Recent Projects Bento Boxes -->
      <div class="flex flex-col gap-4 lg:col-span-6 xl:col-span-7">
        <div class="border-border bg-card h-full rounded-xl border p-6">
          <div class="mb-6 flex items-center justify-between">
            <h2 class="text-foreground text-base font-semibold tracking-tight">
              Recent Projects
            </h2>
            <MoreHorizontal class="text-muted-foreground h-4 w-4 opacity-50" />
          </div>

          <div
            v-if="projectsLoading"
            class="grid grid-cols-1 gap-4 sm:grid-cols-2"
          >
            <div
              v-for="i in 4"
              :key="i"
              class="bg-muted border-border h-32 animate-pulse rounded-lg border"
            />
          </div>
          <div
            v-else-if="projects.length === 0"
            class="border-border bg-muted/30 rounded-lg border border-dashed p-8 text-center"
          >
            <div
              class="bg-muted mx-auto flex h-12 w-12 items-center justify-center rounded-full"
            >
              <Briefcase class="text-muted-foreground/70 h-6 w-6" />
            </div>
            <h3 class="text-foreground mt-2 text-sm font-semibold">
              No projects
            </h3>
            <p class="text-muted-foreground mt-1 text-sm">
              Get started by creating a new project.
            </p>
          </div>

          <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <NuxtLink
              v-for="project in projects.slice(0, 6)"
              :key="project.id"
              :to="`/projects/${project.id}`"
              class="block"
            >
              <div
                class="group border-border from-muted/50 hover:border-foreground/20 hover:bg-muted/80 relative flex h-[120px] flex-col justify-between overflow-hidden rounded-lg border bg-gradient-to-br to-transparent p-4 transition-all"
              >
                <div class="flex items-start justify-between">
                  <div>
                    <h3
                      class="text-foreground mb-1 text-sm font-medium tracking-tight"
                    >
                      {{ project.name }}
                    </h3>
                    <p
                      class="text-muted-foreground line-clamp-2 text-xs leading-relaxed opacity-70"
                    >
                      {{ project.description || "No description." }}
                    </p>
                  </div>
                </div>

                <!-- Mock Avatars for project members -->
                <div class="mt-auto flex -space-x-2">
                  <div
                    class="border-background h-6 w-6 rounded-full border bg-slate-300 dark:bg-slate-800"
                  />
                  <div
                    class="border-background h-6 w-6 rounded-full border bg-slate-400 dark:bg-slate-700"
                  />
                  <div
                    class="border-background h-6 w-6 rounded-full border bg-slate-500 dark:bg-slate-600"
                  />
                </div>
              </div>
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Side Column: My Priorities (Dense List) -->
      <div class="flex flex-col gap-4 lg:col-span-6 xl:col-span-5">
        <div
          class="border-border bg-card flex h-full flex-col rounded-xl border p-6"
        >
          <div class="mb-6 flex items-center justify-between">
            <h2 class="text-foreground text-base font-semibold tracking-tight">
              My Priorities
            </h2>
            <div class="flex gap-2">
              <button
                class="border-border bg-muted/50 hover:bg-muted rounded border p-1 transition-colors"
              >
                <CheckSquare class="text-muted-foreground h-3 w-3" />
              </button>
              <button
                class="border-border bg-muted/50 hover:bg-muted rounded border p-1 transition-colors"
              >
                <Settings class="text-muted-foreground h-3 w-3" />
              </button>
            </div>
          </div>

          <!-- Table Header -->
          <div
            class="border-border text-muted-foreground/80 mb-2 grid grid-cols-[1fr_80px_100px_80px] gap-4 border-b pb-2 text-[10px] font-medium tracking-wider uppercase"
          >
            <div class="pl-2">Task</div>
            <div>Assignee</div>
            <div>Status</div>
            <div class="pr-2 text-right">Due Date</div>
          </div>

          <div
            class="custom-scrollbar -mr-1 flex-1 space-y-1 overflow-y-auto pr-1"
          >
            <div v-if="tasksLoading" class="space-y-2">
              <div
                v-for="i in 5"
                :key="i"
                class="bg-muted h-10 animate-pulse rounded-md"
              />
            </div>

            <div
              v-else-if="recentTasks.length === 0"
              class="flex h-32 items-center justify-center text-center"
            >
              <p class="text-muted-foreground text-sm">No pending tasks.</p>
            </div>

            <template v-else>
              <div
                v-for="task in recentTasks"
                :key="task.id"
                class="group hover:bg-muted/50 grid cursor-pointer grid-cols-[1fr_80px_100px_80px] items-center gap-4 rounded-md p-2 transition-colors"
              >
                <!-- Task Title & Priority -->
                <div class="flex min-w-0 items-center gap-3">
                  <span
                    class="flex h-5 w-10 shrink-0 items-center justify-center rounded-sm text-[10px] font-semibold tracking-wider"
                    :class="
                      task.priority === 'P0'
                        ? 'border border-red-500/20 bg-red-500/10 text-red-600 dark:text-red-400'
                        : task.priority === 'P1'
                          ? 'border border-orange-500/20 bg-orange-500/10 text-orange-600 dark:text-orange-400'
                          : task.priority === 'P2'
                            ? 'border border-yellow-500/20 bg-yellow-500/10 text-yellow-600 dark:text-yellow-400'
                            : 'border border-blue-500/20 bg-blue-500/10 text-blue-600 dark:text-blue-400'
                    "
                  >
                    <Zap
                      v-if="task.priority === 'P0'"
                      class="mr-1 h-2.5 w-2.5"
                    />
                    {{ task.priority }}
                  </span>
                  <p
                    class="text-foreground truncate text-sm font-medium tracking-tight"
                  >
                    {{ task.title }}
                  </p>
                </div>

                <!-- Assignee (Mocked with user initials for now) -->
                <div class="flex items-center gap-2">
                  <div
                    class="bg-muted text-muted-foreground border-border flex h-5 w-5 items-center justify-center rounded-full border text-[9px] font-bold"
                  >
                    {{
                      user?.full_name
                        ? user.full_name.substring(0, 2).toUpperCase()
                        : "ME"
                    }}
                  </div>
                  <span
                    class="text-muted-foreground hidden truncate text-xs sm:block"
                  >
                    {{ user?.full_name?.split(" ")[0] || "Me" }}
                  </span>
                </div>

                <!-- Status Badge -->
                <div>
                  <span
                    class="inline-flex items-center justify-center rounded-sm px-2 py-0.5 text-[10px] font-semibold tracking-wider uppercase"
                    :class="getStatusColor(task.status)"
                  >
                    {{ getStatusLabel(task.status) }}
                  </span>
                </div>

                <!-- Due Date -->
                <div class="text-muted-foreground text-right font-mono text-xs">
                  {{
                    task.due_date
                      ? format(new Date(task.due_date), "MM/dd/yy")
                      : "--/--/--"
                  }}
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>

    <!-- Task Dialog Placeholder -->
    <div v-if="isTaskDialogOpen" class="hidden">
      <!-- To be implemented -->
    </div>
  </div>
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
