<script setup lang="ts">
import {
  Plus,
  Search,
  Filter,
  CheckCircle2,
  Circle,
  Loader2,
} from "lucide-vue-next";
import { TaskStatus, type Task } from "~/types";

interface Props {
  projectId: string;
}

const { projectId } = defineProps<Props>();

const router = useRouter();
const { tasks, isLoading } = useTasks(projectId);
const statusFilter = ref<"open" | "done">("open");
const searchQuery = ref("");

const openTasks = computed(() =>
  tasks.value.filter((t) => t.status !== TaskStatus.DONE)
);
const doneTasks = computed(() =>
  tasks.value.filter((t) => t.status === TaskStatus.DONE)
);

const displayedTasks = computed(() => {
  const filtered =
    statusFilter.value === "open" ? openTasks.value : doneTasks.value;
  return filtered.filter((t) =>
    t.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const handleTaskClick = (task: Task) => {
  router.push(`/projects/${projectId}/tasks/${task.id}`);
};
</script>

<template>
  <div class="space-y-4">
    <!-- Toolbar -->
    <div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
      <div class="relative w-full sm:w-96">
        <Search
          class="text-muted-foreground/60 absolute top-2.5 left-3 h-4 w-4"
        />
        <UiBaseInput
          v-model="searchQuery"
          placeholder="Search all tasks"
          class-name="pl-9 border-border bg-muted"
        />
      </div>
      <NuxtLink :to="`/projects/${projectId}/tasks/new`">
        <UiBaseButton> <Plus class="mr-2 h-4 w-4" /> New Task </UiBaseButton>
      </NuxtLink>
    </div>

    <!-- GitHub Style List Container -->
    <div
      class="border-border bg-card overflow-hidden rounded-lg border shadow-sm"
    >
      <!-- List Header -->
      <div
        class="border-border bg-muted flex items-center justify-between border-b px-4 py-3"
      >
        <div class="flex items-center gap-4">
          <button
            class="flex items-center gap-1.5 text-sm font-medium transition-colors"
            :class="
              statusFilter === 'open'
                ? 'text-card-foreground'
                : 'text-muted-foreground hover:text-card-foreground'
            "
            @click="statusFilter = 'open'"
          >
            <Circle class="h-4 w-4" />
            {{ openTasks.length }} Open
          </button>
          <button
            class="flex items-center gap-1.5 text-sm font-medium transition-colors"
            :class="
              statusFilter === 'done'
                ? 'text-card-foreground'
                : 'text-muted-foreground hover:text-card-foreground'
            "
            @click="statusFilter = 'done'"
          >
            <CheckCircle2 class="h-4 w-4" />
            {{ doneTasks.length }} Done
          </button>
        </div>

        <div class="text-muted-foreground flex items-center gap-4 text-sm">
          <button class="hover:text-card-foreground flex items-center gap-1">
            Sort <Filter class="h-3 w-3" />
          </button>
        </div>
      </div>

      <!-- List Content -->
      <div v-if="isLoading" class="flex h-64 items-center justify-center">
        <Loader2 class="text-muted-foreground/60 h-8 w-8 animate-spin" />
      </div>
      <div v-else class="divide-border divide-y">
        <TasksIssueItem
          v-for="task in displayedTasks"
          :key="task.id"
          :task="task"
          :project-id="projectId"
          @click="handleTaskClick(task)"
        />

        <div v-if="displayedTasks.length === 0" class="p-12 text-center">
          <div
            class="text-muted mx-auto mb-4 flex h-12 w-12 items-center justify-center"
          >
            <Plus class="h-10 w-10" />
          </div>
          <h3 class="text-card-foreground text-lg font-medium">
            No tasks found
          </h3>
          <p class="text-muted-foreground mt-1 text-sm">
            Try adjusting your filters or search query.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
