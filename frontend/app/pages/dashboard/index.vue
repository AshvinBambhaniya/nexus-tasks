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
        <h1 class="text-2xl font-bold text-gray-900">
          {{ getTimeGreeting() }},
          {{ user?.full_name?.split(" ")[0] || "there" }}
        </h1>
        <p class="text-gray-500">
          Here's what's happening in
          <span class="font-medium text-gray-900">
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
          <h2 class="text-lg font-semibold text-gray-900">Recent Projects</h2>
          <NuxtLink
            to="/projects"
            class="text-sm font-medium text-blue-600 hover:text-blue-700"
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
            class="h-40 animate-pulse rounded-lg border border-gray-200 bg-gray-100"
          />
        </div>
        <div
          v-else-if="projects.length === 0"
          class="rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 p-8 text-center"
        >
          <div
            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100"
          >
            <Briefcase class="h-6 w-6 text-gray-400" />
          </div>
          <h3 class="mt-2 text-sm font-semibold text-gray-900">No projects</h3>
          <p class="mt-1 text-sm text-gray-500">
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
              class="group h-full cursor-pointer border-gray-200 p-5 transition-shadow hover:border-blue-200 hover:shadow-md"
            >
              <div class="mb-4 flex items-start justify-between">
                <div
                  class="rounded-lg bg-blue-50 p-2 transition-colors group-hover:bg-blue-100"
                >
                  <Folder class="h-5 w-5 text-blue-600" />
                </div>
                <UiBaseBadge v-if="project.is_archived" variant="secondary">
                  Archived
                </UiBaseBadge>
              </div>
              <h3
                class="mb-1 font-semibold text-gray-900 transition-colors group-hover:text-blue-700"
              >
                {{ project.name }}
              </h3>
              <p class="line-clamp-2 text-sm text-gray-500">
                {{ project.description || "No description provided." }}
              </p>
              <div
                class="mt-4 flex items-center justify-between border-t border-gray-100 pt-4 text-xs text-gray-500"
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
          <h2 class="text-lg font-semibold text-gray-900">My Priorities</h2>
          <NuxtLink
            to="/inbox"
            class="text-sm font-medium text-blue-600 hover:text-blue-700"
          >
            Go to Inbox
          </NuxtLink>
        </div>

        <div class="space-y-3">
          <div v-if="tasksLoading" class="space-y-3">
            <div
              v-for="i in 3"
              :key="i"
              class="h-16 animate-pulse rounded-lg bg-gray-100"
            />
          </div>
          <div
            v-else-if="recentTasks.length === 0"
            class="rounded-lg border border-gray-200 bg-gray-50 p-6 text-center"
          >
            <p class="text-sm text-gray-500">
              No pending tasks in this workspace.
            </p>
          </div>
          <template v-else>
            <div
              v-for="task in recentTasks"
              :key="task.id"
              class="group flex items-start gap-3 rounded-lg border border-gray-200 bg-white p-3 transition-all hover:border-blue-200 hover:shadow-sm"
            >
              <div
                class="mt-0.5 h-2 w-2 flex-shrink-0 rounded-full"
                :class="
                  task.priority === 'P0'
                    ? 'bg-red-500'
                    : task.priority === 'P1'
                      ? 'bg-orange-500'
                      : 'bg-blue-500'
                "
              />
              <div class="min-w-0 flex-1">
                <p class="truncate text-sm font-medium text-gray-900">
                  {{ task.title }}
                </p>
                <div class="mt-1 flex items-center gap-2">
                  <span
                    class="text-[10px] tracking-wider text-gray-500 uppercase"
                  >
                    {{ task.project?.name || `Project #${task.project_id}` }}
                  </span>
                  <span v-if="task.due_date" class="text-[10px] text-gray-400">
                    {{ format(new Date(task.due_date), "MMM d") }}
                  </span>
                </div>
              </div>
            </div>
          </template>

          <NuxtLink
            v-if="workspaceTasks.length > 5"
            to="/inbox"
            class="mt-2 block text-center text-xs text-gray-500 hover:text-gray-900"
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
