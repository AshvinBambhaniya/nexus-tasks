<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { Loader2 } from "lucide-vue-next";
import type { Task } from "~/types";
import { TaskStatus } from "~/types";
import TasksToolbar from "~/components/tasks/TasksToolbar.vue";
import TasksTable from "~/components/tasks/TasksTable.vue";

definePageMeta({ layout: "dashboard" });

const router = useRouter();
const workspaceStore = useWorkspaceStore();
const { tasks, isLoading } = useMyTasks();

const searchQuery = ref("");

const filteredTasks = computed(() => {
  if (!searchQuery.value.trim()) return tasks.value;
  const q = searchQuery.value.toLowerCase();
  return tasks.value.filter((t) => t.title.toLowerCase().includes(q));
});

const handleTaskClick = (task: Task) => {
  if (task.project_id) {
    router.push(`/projects/${task.project_id}/tasks/${task.id}`);
  }
};

const handleToggleStatus = async (task: Task) => {
  const newStatus =
    task.status === TaskStatus.DONE ? TaskStatus.TODO : TaskStatus.DONE;
  const originalStatus = task.status;

  // Optimistic update
  task.status = newStatus;

  try {
    if (!task.project_id || !workspaceStore.activeWorkspaceId) return;

    await useMutation(
      `/api/v2/workspaces/${workspaceStore.activeWorkspaceId}/projects/${task.project_id}/tasks/${task.id}`,
      {
        method: "PATCH",
        body: { status: newStatus },
      }
    );
  } catch (e) {
    console.error("Failed to update status", e);
    task.status = originalStatus; // Revert
  }
};
</script>

<template>
  <div class="mx-auto flex h-full w-full max-w-6xl flex-col px-8 py-6">
    <h1 class="text-foreground mb-2 text-2xl font-semibold">My Tasks</h1>

    <TasksToolbar v-model:search-query="searchQuery" />

    <div v-if="isLoading" class="flex flex-1 items-center justify-center">
      <Loader2 class="text-muted-foreground h-8 w-8 animate-spin" />
    </div>

    <div v-else class="flex-1 overflow-y-auto pb-10">
      <TasksTable
        :tasks="filteredTasks"
        @task-click="handleTaskClick"
        @toggle-status="handleToggleStatus"
      />
    </div>
  </div>
</template>
