<script setup lang="ts">
import type { TaskStatus, Task } from "~/types";

interface Props {
  projectId: string;
}

const { projectId } = defineProps<Props>();

const router = useRouter();
const { tasks, isLoading, updateTask } = useTasks(projectId);

const handleTaskMove = async (taskId: string, newStatus: TaskStatus) => {
  try {
    await updateTask(taskId, { status: newStatus });
  } catch (err) {
    console.error("Failed to move task", err);
  }
};

const handleTaskClick = (task: Task) => {
  router.push(`/projects/${projectId}/tasks/${task.id}`);
};
</script>

<template>
  <div v-if="isLoading" class="flex h-64 items-center justify-center">
    <div class="h-8 w-8 animate-spin rounded-full border-b-2 border-gray-900" />
  </div>
  <div v-else class="h-full">
    <TasksBoardView
      :tasks="tasks"
      @task-move="handleTaskMove"
      @task-click="handleTaskClick"
    />
  </div>
</template>
