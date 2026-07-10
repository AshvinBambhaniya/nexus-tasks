<script setup lang="ts">
import {
  Plus,
  Search,
  Loader2,
  ListFilter,
  ArrowUpDown,
} from "lucide-vue-next";
import { TaskStatus, type Task } from "~/types";

interface Props {
  projectId: string;
}

const { projectId } = defineProps<Props>();

const router = useRouter();
const { tasks, isLoading } = useTasks(projectId);
const statusFilter = ref<"all" | "open" | "done">("all");
const searchQuery = ref("");

const displayedTasks = computed(() => {
  let filtered = tasks.value;
  if (statusFilter.value === "open") {
    filtered = filtered.filter((t) => t.status !== TaskStatus.DONE);
  } else if (statusFilter.value === "done") {
    filtered = filtered.filter((t) => t.status === TaskStatus.DONE);
  }

  return filtered.filter((t) =>
    t.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const handleTaskClick = (task: Task) => {
  router.push(`/projects/${projectId}/tasks/${task.id}`);
};
</script>

<template>
  <div class="flex h-full flex-col space-y-4">
    <!-- Advanced Command Toolbar -->
    <div
      class="bg-background flex flex-col items-center justify-between gap-4 sm:flex-row"
    >
      <div class="flex w-full flex-1 items-center gap-3 sm:w-auto">
        <!-- Search Input -->
        <div class="relative w-full max-w-[320px]">
          <Search
            class="text-muted-foreground absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
          />
          <input
            v-model="searchQuery"
            placeholder="Search issues..."
            class="bg-muted/50 border-border focus:ring-primary focus:border-primary placeholder:text-muted-foreground h-9 w-full rounded-md border pr-4 pl-9 text-sm transition-all focus:ring-1 focus:outline-none"
          />
        </div>

        <!-- Filter Tabs -->
        <div
          class="bg-muted/50 border-border hidden h-9 items-center rounded-md border p-1 sm:flex"
        >
          <button
            class="rounded px-3 py-1 text-xs font-medium transition-colors"
            :class="
              statusFilter === 'all'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="statusFilter = 'all'"
          >
            All
          </button>
          <button
            class="rounded px-3 py-1 text-xs font-medium transition-colors"
            :class="
              statusFilter === 'open'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="statusFilter = 'open'"
          >
            Open
          </button>
          <button
            class="rounded px-3 py-1 text-xs font-medium transition-colors"
            :class="
              statusFilter === 'done'
                ? 'bg-background text-foreground shadow-sm'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="statusFilter = 'done'"
          >
            Done
          </button>
        </div>
      </div>

      <div class="flex w-full items-center gap-2 sm:w-auto">
        <button
          class="border-border bg-card hover:bg-muted/50 text-foreground flex h-9 items-center gap-2 rounded-md border px-3 text-xs font-medium transition-colors"
        >
          <ListFilter class="text-muted-foreground h-3.5 w-3.5" /> Filters
        </button>
        <button
          class="border-border bg-card hover:bg-muted/50 text-foreground flex h-9 items-center gap-2 rounded-md border px-3 text-xs font-medium transition-colors"
        >
          <ArrowUpDown class="text-muted-foreground h-3.5 w-3.5" /> Sort
        </button>
        <NuxtLink :to="`/projects/${projectId}/tasks/new`">
          <button
            class="bg-primary text-primary-foreground hover:bg-primary/90 flex h-9 items-center gap-2 rounded-md px-4 text-xs font-medium shadow-sm transition-colors"
          >
            <Plus class="h-3.5 w-3.5" /> New Issue
          </button>
        </NuxtLink>
      </div>
    </div>

    <!-- Data Table Container -->
    <div
      class="border-border bg-card flex flex-1 flex-col overflow-hidden rounded-lg border shadow-sm"
    >
      <div
        v-if="isLoading"
        class="flex min-h-[400px] flex-1 items-center justify-center"
      >
        <Loader2 class="text-muted-foreground/50 h-8 w-8 animate-spin" />
      </div>

      <div v-else class="flex-1 overflow-auto">
        <table class="text-foreground w-full text-left text-sm">
          <thead class="bg-card border-border sticky top-0 z-10 border-b">
            <tr>
              <th class="w-12 px-4 py-2.5">
                <input
                  type="checkbox"
                  class="border-border text-primary focus:ring-primary/20 bg-muted/50 h-3.5 w-3.5 cursor-pointer rounded"
                />
              </th>
              <th
                class="text-muted-foreground w-[100px] px-2 py-2.5 font-mono text-xs font-semibold"
              >
                ID
              </th>
              <th
                class="text-muted-foreground px-2 py-2.5 text-xs font-semibold"
              >
                Title
              </th>
              <th
                class="text-muted-foreground w-[120px] px-2 py-2.5 text-xs font-semibold"
              >
                Assignee
              </th>
              <th
                class="text-muted-foreground w-[120px] px-2 py-2.5 text-xs font-semibold"
              >
                Status
              </th>
              <th
                class="text-muted-foreground w-[120px] px-2 py-2.5 text-xs font-semibold"
              >
                Priority
              </th>
              <th
                class="text-muted-foreground w-[120px] px-2 py-2.5 pr-6 text-right text-xs font-semibold"
              >
                Due Date
              </th>
            </tr>
          </thead>
          <tbody class="divide-border divide-y">
            <TasksIssueItem
              v-for="task in displayedTasks"
              :key="task.id"
              :task="task"
              :project-id="projectId"
              @click="handleTaskClick(task)"
            />

            <!-- Empty State -->
            <tr v-if="displayedTasks.length === 0">
              <td colspan="7" class="py-16 text-center">
                <div
                  class="bg-muted/50 border-border mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full border"
                >
                  <Search class="text-muted-foreground h-5 w-5" />
                </div>
                <h3 class="text-foreground text-sm font-medium">
                  No tasks found
                </h3>
                <p class="text-muted-foreground mt-1 text-xs">
                  Try adjusting your filters or search query.
                </p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
