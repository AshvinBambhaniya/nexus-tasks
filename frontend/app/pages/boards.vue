<script setup lang="ts">
import { Kanban, Loader2 } from "lucide-vue-next";
import type { TaskStatus } from "~/types";

definePageMeta({ layout: "dashboard" });

const { tasks, isLoading, updateTask } = useMyTasks();

const handleTaskMove = async (taskId: string, newStatus: TaskStatus) => {
  try {
    await updateTask(taskId, { status: newStatus });
  } catch {
    alert("Failed to move task");
  }
};
</script>

<template>
  <div
    class="animate-in fade-in flex h-full flex-col space-y-6 overflow-hidden duration-500"
  >
    <div
      class="flex items-center justify-between border-b border-gray-100 pb-6"
    >
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">
          My Boards
        </h1>
        <p class="mt-1 text-gray-500">
          Visual kanban of all your active tasks.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <div
          class="flex items-center gap-2 rounded-lg border border-gray-100 bg-gray-50/50 px-3 py-1.5 text-xs font-bold tracking-tighter text-gray-500 uppercase"
        >
          <Kanban class="h-3.5 w-3.5" />
          Cross-Project View
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-1 items-center justify-center">
      <div class="flex flex-col items-center gap-4">
        <Loader2 class="h-10 w-10 animate-spin text-blue-600/20" />
        <p class="text-sm font-medium text-gray-400">Building your boards...</p>
      </div>
    </div>

    <div v-else class="flex-1 overflow-hidden">
      <!-- We use the TasksBoardView which handles columns and drag-drop -->
      <TasksBoardView
        :tasks="tasks"
        class="h-full"
        @task-move="handleTaskMove"
      />
    </div>
  </div>
</template>

<style scoped>
/* Ensure the board container takes full height and allows internal scrolling */
:deep(.board-container) {
  height: 100%;
}
</style>
