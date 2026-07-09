<script setup lang="ts">
import type { Task } from '~/types';
import TasksTableRow from './TasksTableRow.vue';

defineProps<{
  tasks: Task[];
}>();

defineEmits<{
  (e: 'task-click', task: Task): void;
  (e: 'toggle-status', task: Task): void;
}>();
</script>

<template>
  <div class="rounded-lg border border-border overflow-hidden bg-background">
    <!-- Header Row -->
    <div class="grid grid-cols-[32px_1fr_200px_100px_100px] items-center gap-4 px-4 py-3 border-b border-border bg-muted/20 text-xs font-medium text-muted-foreground">
      <div></div> <!-- Checkbox column -->
      <div>Task Name</div>
      <div>Project Name</div>
      <div>Priority</div>
      <div>Due Date</div>
    </div>

    <!-- Body -->
    <div class="flex flex-col">
      <div v-if="tasks.length === 0" class="px-4 py-16 text-center text-sm text-muted-foreground">
        No tasks found.
      </div>
      <TasksTableRow
        v-for="task in tasks"
        :key="task.id"
        :task="task"
        @click="$emit('task-click', $event)"
        @toggle-status="$emit('toggle-status', $event)"
      />
    </div>
  </div>
</template>
