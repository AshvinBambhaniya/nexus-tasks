<script setup lang="ts">
import {
  Loader2,
  Inbox,
  CheckCircle2,
  Calendar,
  Search,
  Circle,
  ChevronRight,
  Filter as FilterIcon,
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
    result = result.filter(t => 
      t.title.toLowerCase().includes(q) || 
      t.description?.toLowerCase().includes(q) ||
      t.project?.name?.toLowerCase().includes(q)
    );
  }

  // 2. Status filter
  if (selectedStatus.value !== "ALL") {
    result = result.filter(t => t.status === selectedStatus.value);
  }

  // 3. Priority filter
  if (selectedPriority.value !== "ALL") {
    result = result.filter(t => t.priority === selectedPriority.value);
  }

  // 4. Timeframe filter
  if (selectedTimeframe.value !== "ALL") {
    result = result.filter(t => {
      if (!t.due_date) return false;
      const date = new Date(t.due_date);
      if (selectedTimeframe.value === "OVERDUE") return isPast(date) && !isToday(date) && t.status !== TaskStatus.DONE;
      if (selectedTimeframe.value === "TODAY") return isToday(date);
      if (selectedTimeframe.value === "UPCOMING") return isAfter(date, today) && !isToday(date);
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
  const newStatus = task.status === TaskStatus.DONE ? TaskStatus.TODO : TaskStatus.DONE;
  try {
    await updateTask(task.id, { status: newStatus });
    await refresh();
  } catch (err) {
    alert("Failed to update task");
  }
};

const getPrioColor = (prio: TaskPriority) => {
  switch (prio) {
    case "P0": return "text-red-600 bg-red-50 border-red-100";
    case "P1": return "text-orange-600 bg-orange-50 border-orange-100";
    case "P2": return "text-blue-600 bg-blue-50 border-blue-100";
    default: return "text-gray-500 bg-gray-50 border-gray-100";
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
  <div class="flex h-full flex-col space-y-6 animate-in fade-in duration-500 max-w-6xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col gap-2 border-b border-gray-100 pb-6">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">My Focus</h1>
      <p class="text-sm text-gray-500">Search and filter across all your assigned tasks.</p>
    </div>

    <!-- Toolbar: Search, Filters, Sort -->
    <div class="flex flex-col gap-4 lg:flex-row lg:items-center">
      <!-- Search -->
      <div class="relative flex-1 group">
        <Search class="absolute top-1/2 left-3.5 h-4.5 w-4.5 -translate-y-1/2 text-gray-400 group-focus-within:text-blue-500 transition-colors" />
        <UiBaseInput
          v-model="searchQuery"
          placeholder="Search by title, project, or description..."
          class="h-11 pl-11 shadow-sm ring-1 ring-gray-200 transition-all focus:ring-2 focus:ring-blue-100"
        />
        <button 
          v-if="searchQuery" 
          @click="searchQuery = ''"
          class="absolute top-1/2 right-3 -translate-y-1/2 rounded-full p-1 text-gray-400 hover:bg-gray-100"
        >
          <X class="h-3.5 w-3.5" />
        </button>
      </div>

      <!-- Filters & Sort Buttons (Mobile/Desktop) -->
      <div class="flex flex-wrap items-center gap-3">
        <!-- Status Select -->
        <select 
          v-model="selectedStatus"
          class="h-11 rounded-xl border border-gray-200 bg-white px-4 text-sm font-medium text-gray-700 shadow-sm transition-all focus:border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-50"
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
          class="h-11 rounded-xl border border-gray-200 bg-white px-4 text-sm font-medium text-gray-700 shadow-sm transition-all focus:border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-50"
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
          class="h-11 rounded-xl border border-gray-200 bg-white px-4 text-sm font-medium text-gray-700 shadow-sm transition-all focus:border-blue-400 focus:outline-none focus:ring-2 focus:ring-blue-50"
        >
          <option value="ALL">Any Time</option>
          <option value="OVERDUE">Overdue</option>
          <option value="TODAY">Today</option>
          <option value="UPCOMING">Upcoming</option>
        </select>

        <!-- Sort Select -->
        <div class="flex items-center gap-2 rounded-xl border border-gray-200 bg-gray-50 px-3">
          <ArrowUpDown class="h-4 w-4 text-gray-400" />
          <select 
            v-model="sortBy"
            class="h-11 border-none bg-transparent text-sm font-bold text-gray-600 focus:outline-none"
          >
            <option value="DUE_DATE">Due Date</option>
            <option value="PRIORITY">Priority</option>
            <option value="CREATED_AT">Newest</option>
          </select>
        </div>

        <button 
          v-if="searchQuery || selectedStatus !== 'ALL' || selectedPriority !== 'ALL' || selectedTimeframe !== 'ALL'"
          @click="resetFilters"
          class="flex h-11 items-center gap-2 rounded-xl px-4 text-sm font-bold text-red-600 hover:bg-red-50"
        >
          <X class="h-4 w-4" />
          Clear
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex flex-1 items-center justify-center">
      <div class="flex flex-col items-center gap-4">
        <Loader2 class="h-10 w-10 animate-spin text-blue-600/20" />
        <p class="text-sm font-medium text-gray-400">Filtering tasks...</p>
      </div>
    </div>

    <!-- Content: Unified List -->
    <div v-else class="flex-1 space-y-4 overflow-y-auto pb-10 pr-2 custom-scrollbar">
      <div v-if="filteredTasks.length === 0" class="flex flex-col items-center justify-center py-24 text-center">
        <div class="mb-6 rounded-full bg-gray-50 p-6 ring-1 ring-gray-100">
          <Inbox class="h-12 w-12 text-gray-300" />
        </div>
        <h3 class="text-xl font-bold text-gray-900">No tasks found</h3>
        <p class="mt-2 max-w-xs text-sm text-gray-500">
          No tasks match your current search or filter criteria. Try adjusting your filters.
        </p>
        <UiBaseButton v-if="searchQuery || selectedStatus !== 'ALL'" variant="outline" class="mt-6" @click="resetFilters">
          Reset all filters
        </UiBaseButton>
      </div>

      <div v-else class="divide-y divide-gray-50 overflow-hidden rounded-2xl border border-gray-100 bg-white shadow-sm ring-1 ring-gray-100/50">
        <div 
          v-for="task in filteredTasks" 
          :key="task.id"
          class="group flex items-center gap-4 px-6 py-4 transition-all hover:bg-blue-50/10"
        >
          <!-- Quick Complete Toggle -->
          <button 
            @click="handleToggleDone(task)"
            class="flex h-6 w-6 shrink-0 items-center justify-center rounded-full border-2 transition-all active:scale-90"
            :class="task.status === TaskStatus.DONE 
              ? 'bg-green-500 border-green-500 text-white' 
              : 'border-gray-200 hover:border-blue-400'"
          >
            <CheckCircle2 v-if="task.status === TaskStatus.DONE" class="h-4 w-4" />
            <Circle v-else class="h-4 w-4 text-transparent group-hover:text-blue-100" />
          </button>

          <!-- Task Main Content -->
          <div class="min-w-0 flex-1 space-y-1">
            <div class="flex flex-wrap items-center gap-3">
              <NuxtLink 
                :to="`/projects/${task.project_id}`"
                class="flex items-center gap-1.5 text-[10px] font-black tracking-widest text-gray-400 uppercase transition-colors hover:text-blue-600"
              >
                <Hash class="h-3 w-3" />
                {{ task.project?.name || 'Project' }}
              </NuxtLink>
              <span class="text-[10px] font-bold text-gray-300">#{{ task.number }}</span>
              <UiBaseBadge 
                variant="outline" 
                :class="[getPrioColor(task.priority), 'h-4.5 px-1.5 text-[9px] font-black uppercase tracking-tighter']"
              >
                {{ task.priority }}
              </UiBaseBadge>
            </div>
            
            <h4 :class="['truncate font-bold text-gray-900 transition-all group-hover:text-blue-700', task.status === TaskStatus.DONE ? 'line-through opacity-50' : '']">
              {{ task.title }}
            </h4>
            
            <p v-if="task.description" class="line-clamp-1 text-xs text-gray-500">
              {{ task.description }}
            </p>
          </div>

          <!-- Metadata & Date -->
          <div class="flex shrink-0 items-center gap-8">
            <div v-if="task.due_date" class="hidden md:flex items-center gap-1.5 text-[11px] font-bold">
              <Calendar class="h-3.5 w-3.5 text-gray-400" />
              <span :class="[
                isPast(new Date(task.due_date)) && !isToday(new Date(task.due_date)) && task.status !== TaskStatus.DONE 
                ? 'text-red-600' 
                : isToday(new Date(task.due_date)) ? 'text-blue-600' : 'text-gray-500'
              ]">
                {{ format(new Date(task.due_date), "MMM d") }}
              </span>
            </div>

            <div class="flex items-center gap-3">
              <UiBaseBadge :class="['hidden sm:inline-flex px-2 py-0.5 text-[9px] font-black uppercase tracking-wider', 
                task.status === TaskStatus.DONE ? 'bg-green-100 text-green-700' : 
                task.status === TaskStatus.IN_PROGRESS ? 'bg-purple-100 text-purple-700' : 'bg-gray-100 text-gray-600']">
                {{ task.status.replace("_", " ") }}
              </UiBaseBadge>
              
              <NuxtLink 
                :to="`/projects/${task.project_id}/tasks`"
                class="flex h-9 w-9 items-center justify-center rounded-xl text-gray-300 transition-all hover:bg-white hover:text-blue-600 hover:shadow-md ring-1 ring-transparent hover:ring-gray-100"
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
  background: #f1f1f1;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #e5e7eb;
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
