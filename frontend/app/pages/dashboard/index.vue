<script setup lang="ts">
import {
  Folder,
  Users,
  CheckSquare,
  Plus,
  ArrowRight,
  Activity,
  Briefcase,
} from "lucide-vue-next";
import { format } from "date-fns";
import { TaskStatus, type WorkspaceMember } from "~/types";
import { useUsersStore } from "~/stores/user";

definePageMeta({
  layout: "dashboard",
});

const userStore = useUsersStore();
const user = computed(() => userStore.userData);
const { activeWorkspace } = useWorkspaces();
const { projects, isLoading: projectsLoading } = useProjects();
const { tasks: myTasks, isLoading: tasksLoading } = useMyTasks();

// Fetch members for stats
const { data: members } = useApi<WorkspaceMember[]>(
  () =>
    activeWorkspace.value
      ? `/api/v2/workspaces/${activeWorkspace.value.id}/members`
      : "/api/v2/workspaces/0/members",
  {
    key: "workspace-members",
    watch: [activeWorkspace],
  }
);

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
    .slice(0, 5)
);

const getTimeGreeting = () => {
  const hour = new Date().getHours();
  if (hour < 12) return "Good morning";
  if (hour < 18) return "Good afternoon";
  return "Good evening";
};
</script>

<template>
  <div class="space-y-8">
    <!-- Header -->
    <div
      class="flex flex-col gap-1 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <h1 class="text-foreground text-2xl font-bold">
          {{ getTimeGreeting() }},
          {{ user?.full_name?.split(" ")[0] || "there" }}
        </h1>
        <p class="text-muted-foreground">
          Here's what's happening in
          <span class="text-foreground font-medium">
            {{ activeWorkspace?.name }}
          </span>
          today.
        </p>
      </div>
      <div class="mt-4 flex gap-2 sm:mt-0">
        <UiBaseButton
          size="sm"
          class="hidden sm:flex"
          :disabled="projects.length === 0"
          :title="
            projects.length === 0
              ? 'Create a project first'
              : 'Create a new task'
          "
          @click="isTaskDialogOpen = true"
        >
          <Plus class="mr-2 h-4 w-4" /> New Task
        </UiBaseButton>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
      <DashboardStatsCard
        title="Active Projects"
        :value="projects.filter((p) => !p.is_archived).length"
        :icon="Folder"
        :is-loading="projectsLoading"
        color="blue"
      />
      <DashboardStatsCard
        title="Team Members"
        :value="members?.length || 0"
        :icon="Users"
        :is-loading="!members"
        color="indigo"
      />
      <DashboardStatsCard
        title="Pending Tasks"
        :value="pendingTasksCount"
        :icon="CheckSquare"
        :is-loading="tasksLoading"
        color="orange"
      />
      <DashboardStatsCard
        title="Workspace Activity"
        value="Unknown"
        label="Last 7 days"
        :icon="Activity"
        :is-loading="false"
        color="green"
      />
    </div>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <!-- Main Column: Projects -->
      <div class="space-y-6 lg:col-span-2">
        <div class="flex items-center justify-between">
          <h2 class="text-foreground text-lg font-semibold">Recent Projects</h2>
          <NuxtLink
            to="/projects"
            class="text-primary hover:text-primary/80 text-sm font-medium"
          >
            View all
          </NuxtLink>
        </div>

        <div
          v-if="projectsLoading"
          class="grid grid-cols-1 gap-4 sm:grid-cols-2"
        >
          <div
            v-for="i in 4"
            :key="i"
            class="border-border bg-muted h-40 animate-pulse rounded-lg border"
          />
        </div>
        <div
          v-else-if="projects.length === 0"
          class="border-border bg-muted/50 rounded-lg border-2 border-dashed p-8 text-center"
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
          <div class="mt-6">
            <UiBaseButton variant="outline" size="sm">
              Create Project
            </UiBaseButton>
          </div>
        </div>
        <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <NuxtLink
            v-for="project in projects.slice(0, 4)"
            :key="project.id"
            :to="`/projects/${project.id}`"
            class="block h-full"
          >
            <UiBaseCard
              class="group border-border hover:border-primary/50 h-full cursor-pointer p-5 transition-shadow hover:shadow-md"
            >
              <div class="mb-4 flex items-start justify-between">
                <div
                  class="bg-primary/10 group-hover:bg-primary/20 rounded-lg p-2 transition-colors"
                >
                  <Folder class="text-primary h-5 w-5" />
                </div>
                <UiBaseBadge v-if="project.is_archived" variant="secondary">
                  Archived
                </UiBaseBadge>
              </div>
              <h3
                class="text-foreground group-hover:text-primary mb-1 font-semibold transition-colors"
              >
                {{ project.name }}
              </h3>
              <p class="text-muted-foreground line-clamp-2 text-sm">
                {{ project.description || "No description provided." }}
              </p>
              <div
                class="border-border text-muted-foreground mt-4 flex items-center justify-between border-t pt-4 text-xs"
              >
                <span>
                  Updated {{ format(new Date(project.created_at), "MMM d") }}
                </span>
                <ArrowRight
                  class="h-3 w-3 transform opacity-0 transition-opacity group-hover:translate-x-1 group-hover:opacity-100"
                />
              </div>
            </UiBaseCard>
          </NuxtLink>
        </div>
      </div>

      <!-- Side Column: My Tasks -->
      <div class="space-y-6">
        <div class="flex items-center justify-between">
          <h2 class="text-foreground text-lg font-semibold">My Priorities</h2>
          <NuxtLink
            to="/inbox"
            class="text-primary hover:text-primary/80 text-sm font-medium"
          >
            Go to Inbox
          </NuxtLink>
        </div>

        <div class="space-y-3">
          <div v-if="tasksLoading" class="space-y-3">
            <div
              v-for="i in 3"
              :key="i"
              class="bg-muted h-16 animate-pulse rounded-lg"
            />
          </div>
          <div
            v-else-if="recentTasks.length === 0"
            class="border-border bg-muted/50 rounded-lg border p-6 text-center"
          >
            <p class="text-muted-foreground text-sm">
              No pending tasks in this workspace.
            </p>
          </div>
          <template v-else>
            <div
              v-for="task in recentTasks"
              :key="task.id"
              class="group border-border bg-card hover:border-primary/50 flex items-start gap-3 rounded-lg border p-3 transition-all hover:shadow-sm"
            >
              <div
                class="mt-0.5 h-2 w-2 flex-shrink-0 rounded-full"
                :class="
                  task.priority === 'P0'
                    ? 'bg-destructive'
                    : task.priority === 'P1'
                      ? 'bg-orange-500'
                      : 'bg-primary'
                "
              />
              <div class="min-w-0 flex-1">
                <p class="text-foreground truncate text-sm font-medium">
                  {{ task.title }}
                </p>
                <div class="mt-1 flex items-center gap-2">
                  <span
                    class="text-muted-foreground text-[10px] tracking-wider uppercase"
                  >
                    {{ task.project?.name || `Project #${task.project_id}` }}
                  </span>
                  <span
                    v-if="task.due_date"
                    class="text-muted-foreground/70 text-[10px]"
                  >
                    {{ format(new Date(task.due_date), "MMM d") }}
                  </span>
                </div>
              </div>
            </div>
          </template>

          <NuxtLink
            v-if="workspaceTasks.length > 5"
            to="/inbox"
            class="text-muted-foreground hover:text-foreground mt-2 block text-center text-xs"
          >
            + {{ workspaceTasks.length - 5 }} more tasks
          </NuxtLink>
        </div>
      </div>
    </div>

    <!-- Task Dialog Placeholder -->
    <div v-if="isTaskDialogOpen" class="hidden">
      <!-- To be implemented -->
    </div>
  </div>
</template>
