<script setup lang="ts">
import { TaskStatus, type Task } from "~/types";

interface Props {
  tasks: Task[];
}

const { tasks } = defineProps<Props>();

const emit = defineEmits(["task-move", "task-click"]);

const columns = [
  { id: TaskStatus.BACKLOG, title: "Backlog", color: "bg-gray-400" },
  { id: TaskStatus.TODO, title: "To Do", color: "bg-blue-500" },
  { id: TaskStatus.IN_PROGRESS, title: "In Progress", color: "bg-yellow-500" },
  { id: TaskStatus.DONE, title: "Done", color: "bg-green-500" },
];

const getTasksByStatus = (status: TaskStatus) => {
  return tasks.filter((t) => t.status === status);
};

const handleTaskMove = (event: { taskId: number; newStatus: TaskStatus }) => {
  emit("task-move", event.taskId, event.newStatus);
};
</script>

<template>
  <div class="flex h-full gap-4 overflow-x-auto pb-4">
    <TasksBoardColumn
      v-for="col in columns"
      :id="col.id"
      :key="col.id"
      :title="col.title"
      :color="col.color"
      :tasks="getTasksByStatus(col.id)"
      @task-move="handleTaskMove"
      @task-click="(task) => emit('task-click', task)"
    />
  </div>
</template>
