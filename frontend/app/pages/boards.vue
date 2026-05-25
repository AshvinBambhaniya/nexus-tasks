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
    <div class="border-border flex items-center justify-between border-b pb-6">
      <div>
        <h1 class="text-foreground text-3xl font-bold tracking-tight">
          My Boards
        </h1>
        <p class="text-muted-foreground mt-1">
          Visual kanban of all your active tasks.
        </p>
      </div>
      <div class="flex items-center gap-3">
        <div
          class="border-border bg-muted/50 text-muted-foreground flex items-center gap-2 rounded-lg border px-3 py-1.5 text-xs font-bold tracking-tighter uppercase"
        >
          <Kanban class="h-3.5 w-3.5" />
          Cross-Project View
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="flex flex-1 items-center justify-center">
      <div class="flex flex-col items-center gap-4">
        <Loader2 class="text-primary/20 h-10 w-10 animate-spin" />
        <p class="text-muted-foreground/70 text-sm font-medium">
          Building your boards...
        </p>
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
