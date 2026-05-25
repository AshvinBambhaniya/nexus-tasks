<script setup lang="ts">
import {
  Loader2,
  Inbox,
  CheckCircle2,
  Calendar,
  Search,
  Circle,
  ChevronRight,
  X,
  ArrowUpDown,
  Hash,
} from "lucide-vue-next";
import { format, isPast, isToday, isAfter, startOfToday } from "date-fns";
import { TaskStatus, TaskPriority, type TaskWithProject } from "~/types";

definePageMeta({ layout: "dashboard" });

const { tasks, isLoading, updateTask, refresh } = useMyTasks();

// Search and Filter State
const searchQuery = ref("");
const selectedStatus = ref<TaskStatus | "ALL">("ALL");
const selectedPriority = ref<TaskPriority | "ALL">("ALL");
const selectedTimeframe = ref<"ALL" | "OVERDUE" | "TODAY" | "UPCOMING">("ALL");
const sortBy = ref<"DUE_DATE" | "PRIORITY" | "CREATED_AT">("DUE_DATE");

const filteredTasks = computed(() => {
  if (!tasks.value) return [];

  let result = [...tasks.value];
  const today = startOfToday();

  // 1. Search filter
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.toLowerCase();
    result = result.filter(
      (t) =>
        t.title.toLowerCase().includes(q) ||
        t.description?.toLowerCase().includes(q) ||
        t.project?.name?.toLowerCase().includes(q)
    );
  }

  // 2. Status filter
  if (selectedStatus.value !== "ALL") {
    result = result.filter((t) => t.status === selectedStatus.value);
  }

  // 3. Priority filter
  if (selectedPriority.value !== "ALL") {
    result = result.filter((t) => t.priority === selectedPriority.value);
  }

  // 4. Timeframe filter
  if (selectedTimeframe.value !== "ALL") {
    result = result.filter((t) => {
      if (!t.due_date) return false;
      const date = new Date(t.due_date);
      if (selectedTimeframe.value === "OVERDUE")
        return isPast(date) && !isToday(date) && t.status !== TaskStatus.DONE;
      if (selectedTimeframe.value === "TODAY") return isToday(date);
      if (selectedTimeframe.value === "UPCOMING")
        return isAfter(date, today) && !isToday(date);
      return true;
    });
  }

  // 5. Sorting
  result.sort((a, b) => {
    if (sortBy.value === "DUE_DATE") {
      if (!a.due_date) return 1;
      if (!b.due_date) return -1;
      return new Date(a.due_date).getTime() - new Date(b.due_date).getTime();
    }
    if (sortBy.value === "PRIORITY") {
      const pMap = { P0: 0, P1: 1, P2: 2, P3: 3 };
      return pMap[a.priority] - pMap[b.priority];
    }
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
  });

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

const getPrioColor = (prio: TaskPriority) => {
  switch (prio) {
    case "P0":
      return "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-100 dark:border-red-900/50";
    case "P1":
      return "text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-950/30 border-orange-100 dark:border-orange-900/50";
    case "P2":
      return "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-100 dark:border-blue-900/50";
    default:
      return "text-muted-foreground bg-muted border-border";
  }
};

const resetFilters = () => {
  searchQuery.value = "";
  selectedStatus.value = "ALL";
  selectedPriority.value = "ALL";
  selectedTimeframe.value = "ALL";
};
</script>

<template>
  <div
    class="animate-in fade-in mx-auto flex h-full max-w-6xl flex-col space-y-6 duration-500"
  >
    <!-- Header -->
    <div class="border-border flex flex-col gap-2 border-b pb-6">
      <h1 class="text-foreground text-3xl font-bold tracking-tight">
        My Focus
      </h1>
      <p class="text-muted-foreground text-sm">
        Search and filter across all your assigned tasks.
      </p>
    </div>

    <!-- Toolbar: Search, Filters, Sort -->
    <div class="flex flex-col gap-4 lg:flex-row lg:items-center">
      <!-- Search -->
      <div class="group relative flex-1">
        <Search
          class="text-muted-foreground group-focus-within:text-primary absolute top-1/2 left-3.5 h-4.5 w-4.5 -translate-y-1/2 transition-colors"
        />
        <UiBaseInput
          v-model="searchQuery"
          placeholder="Search by title, project, or description..."
          class="ring-border focus:ring-primary/10 h-11 pl-11 shadow-sm ring-1 transition-all focus:ring-2"
        />
        <button
          v-if="searchQuery"
          class="text-muted-foreground hover:bg-muted absolute top-1/2 right-3 -translate-y-1/2 rounded-full p-1"
          @click="searchQuery = ''"
        >
          <X class="h-3.5 w-3.5" />
        </button>
      </div>

      <!-- Filters & Sort Buttons (Mobile/Desktop) -->
      <div class="flex flex-wrap items-center gap-3">
        <!-- Status Select -->
        <select
          v-model="selectedStatus"
          class="border-border bg-card text-foreground focus:border-primary focus:ring-primary/10 h-11 rounded-xl border px-4 text-sm font-medium shadow-sm transition-all focus:ring-2 focus:outline-none"
        >
          <option value="ALL">All Status</option>
          <option :value="TaskStatus.TODO">Todo</option>
          <option :value="TaskStatus.IN_PROGRESS">In Progress</option>
          <option :value="TaskStatus.BACKLOG">Backlog</option>
          <option :value="TaskStatus.DONE">Done</option>
        </select>

        <!-- Priority Select -->
        <select
          v-model="selectedPriority"
          class="border-border bg-card text-foreground focus:border-primary focus:ring-primary/10 h-11 rounded-xl border px-4 text-sm font-medium shadow-sm transition-all focus:ring-2 focus:outline-none"
        >
          <option value="ALL">All Priority</option>
          <option :value="TaskPriority.P0">P0 - Critical</option>
          <option :value="TaskPriority.P1">P1 - High</option>
          <option :value="TaskPriority.P2">P2 - Medium</option>
          <option :value="TaskPriority.P3">P3 - Low</option>
        </select>

        <!-- Timeframe Select -->
        <select
          v-model="selectedTimeframe"
          class="border-border bg-card text-foreground focus:border-primary focus:ring-primary/10 h-11 rounded-xl border px-4 text-sm font-medium shadow-sm transition-all focus:ring-2 focus:outline-none"
        >
          <option value="ALL">Any Time</option>
          <option value="OVERDUE">Overdue</option>
          <option value="TODAY">Today</option>
          <option value="UPCOMING">Upcoming</option>
        </select>

        <!-- Sort Select -->
        <div
          class="border-border bg-muted flex items-center gap-2 rounded-xl border px-3"
        >
          <ArrowUpDown class="text-muted-foreground h-4 w-4" />
          <select
            v-model="sortBy"
            class="text-muted-foreground h-11 border-none bg-transparent text-sm font-bold focus:outline-none"
          >
            <option value="DUE_DATE">Due Date</option>
            <option value="PRIORITY">Priority</option>
            <option value="CREATED_AT">Newest</option>
          </select>
        </div>

        <button
          v-if="
            searchQuery ||
            selectedStatus !== 'ALL' ||
            selectedPriority !== 'ALL' ||
            selectedTimeframe !== 'ALL'
          "
          class="text-destructive hover:bg-destructive/10 flex h-11 items-center gap-2 rounded-xl px-4 text-sm font-bold"
          @click="resetFilters"
        >
          <X class="h-4 w-4" />
          Clear
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-1 items-center justify-center">
      <div class="flex flex-col items-center gap-4">
        <Loader2 class="text-primary/20 h-10 w-10 animate-spin" />
        <p class="text-muted-foreground text-sm font-medium">
          Filtering tasks...
        </p>
      </div>
    </div>

    <!-- Content: Unified List -->
    <div
      v-else
      class="custom-scrollbar flex-1 space-y-4 overflow-y-auto pr-2 pb-10"
    >
      <div
        v-if="filteredTasks.length === 0"
        class="flex flex-col items-center justify-center py-24 text-center"
      >
        <div class="bg-muted ring-border mb-6 rounded-full p-6 ring-1">
          <Inbox class="text-muted-foreground/50 h-12 w-12" />
        </div>
        <h3 class="text-foreground text-xl font-bold">No tasks found</h3>
        <p class="text-muted-foreground mt-2 max-w-xs text-sm">
          No tasks match your current search or filter criteria. Try adjusting
          your filters.
        </p>
        <UiBaseButton
          v-if="searchQuery || selectedStatus !== 'ALL'"
          variant="outline"
          class="mt-6"
          @click="resetFilters"
        >
          Reset all filters
        </UiBaseButton>
      </div>

      <div
        v-else
        class="divide-border border-border bg-card ring-border divide-y overflow-hidden rounded-2xl border shadow-sm ring-1"
      >
        <div
          v-for="task in filteredTasks"
          :key="task.id"
          class="group hover:bg-primary/5 flex items-center gap-4 px-6 py-4 transition-all"
        >
          <!-- Quick Complete Toggle -->
          <button
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 transition-all active:scale-90"
            :class="
              task.status === TaskStatus.DONE
                ? 'border-green-500 bg-green-500 text-white'
                : 'border-border hover:border-primary'
            "
            @click="handleToggleDone(task)"
          >
            <CheckCircle2
              v-if="task.status === TaskStatus.DONE"
              class="h-4 w-4"
            />
            <Circle
              v-else
              class="group-hover:text-primary/20 h-4 w-4 text-transparent"
            />
          </button>

          <!-- Task Main Content -->
          <div class="min-w-0 flex-1 space-y-1">
            <div class="flex flex-wrap items-center gap-3">
              <NuxtLink
                :to="`/projects/${task.project_id}`"
                class="text-muted-foreground hover:text-primary flex items-center gap-1.5 text-[10px] font-black tracking-widest uppercase transition-colors"
              >
                <Hash class="h-3 w-3" />
                {{ task.project?.name || "Project" }}
              </NuxtLink>
              <span class="text-muted-foreground/50 text-[10px] font-bold"
                >#{{ task.number }}</span
              >
              <UiBaseBadge
                variant="outline"
                :class="[
                  getPrioColor(task.priority),
                  'h-4.5 px-1.5 text-[9px] font-black tracking-tighter uppercase',
                ]"
              >
                {{ task.priority }}
              </UiBaseBadge>
            </div>

            <h4
              :class="[
                'text-foreground group-hover:text-primary truncate font-bold transition-all',
                task.status === TaskStatus.DONE
                  ? 'line-through opacity-50'
                  : '',
              ]"
            >
              {{ task.title }}
            </h4>

            <p
              v-if="task.description"
              class="text-muted-foreground line-clamp-1 text-xs"
            >
              {{ task.description }}
            </p>
          </div>

          <!-- Metadata & Date -->
          <div class="flex shrink-0 items-center gap-8">
            <div
              v-if="task.due_date"
              class="hidden items-center gap-1.5 text-[11px] font-bold md:flex"
            >
              <Calendar class="text-muted-foreground h-3.5 w-3.5" />
              <span
                :class="[
                  isPast(new Date(task.due_date)) &&
                  !isToday(new Date(task.due_date)) &&
                  task.status !== TaskStatus.DONE
                    ? 'text-destructive dark:text-red-400'
                    : isToday(new Date(task.due_date))
                      ? 'text-primary'
                      : 'text-muted-foreground',
                ]"
              >
                {{ format(new Date(task.due_date), "MMM d") }}
              </span>
            </div>

            <div class="flex items-center gap-3">
              <UiBaseBadge
                :class="[
                  'hidden px-2 py-0.5 text-[9px] font-black tracking-wider uppercase sm:inline-flex',
                  task.status === TaskStatus.DONE
                    ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                    : task.status === TaskStatus.IN_PROGRESS
                      ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400'
                      : 'bg-muted text-muted-foreground',
                ]"
              >
                {{ task.status.replace("_", " ") }}
              </UiBaseBadge>

              <NuxtLink
                :to="`/projects/${task.project_id}/tasks`"
                class="text-muted-foreground/50 hover:bg-card hover:text-primary hover:ring-border flex h-9 w-9 items-center justify-center rounded-xl ring-1 ring-transparent transition-all hover:shadow-md"
              >
                <ChevronRight class="h-5 w-5" />
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
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

/* Remove default select styling */
select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%239ca3af' stroke-width='2'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 1rem;
}
</style>
