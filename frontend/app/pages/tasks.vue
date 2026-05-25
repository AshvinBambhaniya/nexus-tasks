<script setup lang="ts">
import { Plus } from "lucide-vue-next";
import type { Task } from "~/types";

definePageMeta({ layout: "dashboard" });

const router = useRouter();
const { tasks, isLoading } = useTasks();

const handleCreateClick = () => {
  // If we don't have a projectId, we might need a project selector or just redirect to projects
  router.push("/projects");
};

const handleTaskClick = (task: Task) => {
  router.push(`/projects/${task.project_id}/tasks/${task.id}`);
};
</script>

<template>
  <div class="flex h-full flex-col space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-foreground text-2xl font-bold tracking-tight">
        All Tasks
      </h1>
      <UiBaseButton @click="handleCreateClick">
        <Plus class="mr-2 h-4 w-4" /> Create Task
      </UiBaseButton>
    </div>

    <div v-if="isLoading" class="flex h-64 items-center justify-center">
      <div
        class="border-primary h-8 w-8 animate-spin rounded-full border-b-2"
      />
    </div>
    <div v-else class="flex-1 overflow-hidden">
      <TasksListView :tasks="tasks" @task-click="handleTaskClick" />
    </div>
  </div>
</template>
