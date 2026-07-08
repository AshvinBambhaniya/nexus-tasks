<script setup lang="ts">
import { Loader2, Inbox, Search, Filter } from "lucide-vue-next";
import type { TaskStatus, type TaskWithProject } from "~/types";

definePageMeta({ layout: "dashboard" });

const { tasks, isLoading, updateTask, refresh } = useMyTasks();

// Search and Filter State
const searchQuery = ref("");
const selectedFilter = ref<"ALL" | TaskStatus>("ALL");
const activeTaskId = ref<string | null>(null);

const activeTask = computed(() => {
  if (!activeTaskId.value) return null;
  return tasks.value?.find((t) => t.id === activeTaskId.value) || null;
});

// Watch tasks to auto-select the first one if none is selected
watch(tasks, (newTasks) => {
  if (newTasks && newTasks.length > 0 && !activeTaskId.value) {
    activeTaskId.value = newTasks[0].id;
  }
});

const filteredTasks = computed(() => {
  if (!tasks.value) return [];

  let result = [...tasks.value];

  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(
      (t) =>
        t.title.toLowerCase().includes(q) ||
        t.description?.toLowerCase().includes(q) ||
        t.project?.name?.toLowerCase().includes(q)
    );
  }

  if (selectedFilter.value !== "ALL") {
    result = result.filter((t) => t.status === selectedFilter.value);
  }

  // Sort by newest created first or due date
  result.sort(
    (a, b) =>
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  );

  return result;
});

const handleToggleDone = async (task: TaskWithProject) => {
  const newStatus =
    task.status === TaskStatus.DONE ? TaskStatus.TODO : TaskStatus.DONE;
  try {
    await updateTask(task.id, { status: newStatus });
    await refresh();
  } catch {
    alert("Failed to update task");
  }
};
</script>

<template>
  <div
    class="border-border bg-card animate-in fade-in mx-auto flex h-[calc(100vh-8rem)] max-h-[1000px] w-full max-w-[1600px] overflow-hidden rounded-xl border shadow-sm duration-500"
  >
    <!-- Left Pane: List View -->
    <div
      class="border-border bg-card/50 flex w-[350px] shrink-0 flex-col border-r xl:w-[400px]"
    >
      <!-- Toolbar -->
      <div class="border-border flex flex-col gap-3 border-b p-4">
        <div class="flex items-center justify-between">
          <h2 class="text-foreground text-sm font-semibold tracking-tight">
            Inbox ({{ filteredTasks.length }})
          </h2>
          <button
            class="text-muted-foreground hover:text-foreground transition-colors"
          >
            <Filter class="h-4 w-4" />
          </button>
        </div>

        <div class="relative">
          <Search
            class="text-muted-foreground absolute top-2 left-2.5 h-4 w-4"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search tasks..."
            class="border-border bg-muted/50 text-foreground focus:border-primary focus:ring-primary h-8 w-full rounded-md border pr-3 pl-9 text-sm transition-all focus:ring-1 focus:outline-none"
          />
        </div>

        <div
          class="custom-scrollbar mt-1 flex items-center gap-1 overflow-x-auto pb-1"
        >
          <button
            :class="[
              'rounded-md px-3 py-1 text-xs font-medium whitespace-nowrap transition-colors',
              selectedFilter === 'ALL'
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
            ]"
            @click="selectedFilter = 'ALL'"
          >
            All
          </button>
          <button
            :class="[
              'rounded-md px-3 py-1 text-xs font-medium whitespace-nowrap transition-colors',
              selectedFilter === TaskStatus.TODO
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
            ]"
            @click="selectedFilter = TaskStatus.TODO"
          >
            Todo
          </button>
          <button
            :class="[
              'rounded-md px-3 py-1 text-xs font-medium whitespace-nowrap transition-colors',
              selectedFilter === TaskStatus.IN_PROGRESS
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
            ]"
            @click="selectedFilter = TaskStatus.IN_PROGRESS"
          >
            In Progress
          </button>
          <button
            :class="[
              'rounded-md px-3 py-1 text-xs font-medium whitespace-nowrap transition-colors',
              selectedFilter === TaskStatus.DONE
                ? 'bg-muted text-foreground'
                : 'text-muted-foreground hover:text-foreground hover:bg-muted/50',
            ]"
            @click="selectedFilter = TaskStatus.DONE"
          >
            Done
          </button>
        </div>
      </div>

      <!-- Task List -->
      <div class="custom-scrollbar flex-1 overflow-y-auto">
        <div v-if="isLoading" class="flex items-center justify-center py-10">
          <Loader2 class="text-muted-foreground h-6 w-6 animate-spin" />
        </div>

        <div
          v-else-if="filteredTasks.length === 0"
          class="flex flex-col items-center justify-center px-6 py-20 text-center"
        >
          <Inbox class="text-muted-foreground/50 mb-4 h-10 w-10" />
          <p class="text-foreground text-sm font-medium">
            You're all caught up.
          </p>
          <p class="text-muted-foreground mt-1 text-xs">
            No tasks match your current filters.
          </p>
        </div>

        <div v-else class="flex flex-col">
          <button
            v-for="task in filteredTasks"
            :key="task.id"
            :class="[
              'border-border/50 hover:bg-muted/50 relative flex flex-col items-start gap-1.5 border-b p-4 text-left transition-colors',
              activeTaskId === task.id ? 'bg-muted' : '',
            ]"
            @click="activeTaskId = task.id"
          >
            <!-- Active Indicator -->
            <div
              v-if="activeTaskId === task.id"
              class="bg-primary absolute top-0 left-0 h-full w-[3px]"
            />
            <!-- Unread Indicator (Simulated if not done) -->
            <div
              v-if="task.status !== TaskStatus.DONE"
              class="absolute top-4 right-4 h-2 w-2 rounded-full bg-blue-500 shadow-[0_0_8px_rgba(59,130,246,0.6)]"
            />

            <div class="flex w-full items-center gap-2 pr-6">
              <span
                class="text-muted-foreground/70 text-[10px] font-bold tracking-widest uppercase"
              >
                {{ task.project?.name || "Project" }}
              </span>
              <span class="text-muted-foreground/50 text-[10px]">2m ago</span>
            </div>

            <h4
              :class="[
                'w-full truncate text-sm font-semibold tracking-tight',
                task.status === TaskStatus.DONE
                  ? 'text-muted-foreground line-through'
                  : 'text-foreground',
              ]"
            >
              {{ task.title }}
            </h4>

            <p
              v-if="task.description"
              class="text-muted-foreground line-clamp-1 w-full text-xs"
            >
              {{ task.description }}
            </p>
          </button>
        </div>
      </div>
    </div>

    <!-- Right Pane: Detail View -->
    <div class="bg-background/30 flex flex-1 flex-col">
      <div
        v-if="isLoading && !activeTask"
        class="flex h-full items-center justify-center"
      >
        <Loader2 class="text-muted-foreground h-8 w-8 animate-spin" />
      </div>

      <div
        v-else-if="!activeTask"
        class="flex h-full flex-col items-center justify-center text-center"
      >
        <Inbox class="text-muted-foreground/30 mb-4 h-12 w-12" />
        <h3 class="text-muted-foreground text-lg font-medium">
          Select an item to view details
        </h3>
      </div>

      <TaskDetail
        v-else
        :key="activeTask.id"
        :task-id="activeTask.id"
        :project-id="activeTask.project_id"
        @toggle-done="handleToggleDone"
      />
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
  background: var(--muted);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: var(--border);
}
</style>
