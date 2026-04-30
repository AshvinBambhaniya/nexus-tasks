<script setup lang="ts">
import draggable from "vuedraggable";
import type { TaskStatus, Task } from "~/types";

interface Props {
  id: TaskStatus;
  title: string;
  tasks: Task[];
  color?: string;
}

const { id, title, tasks, color = "bg-gray-400" } = defineProps<Props>();

const emit = defineEmits(["task-move", "task-click"]);

// Internal list for draggable
const internalTasks = computed({
  get: () => tasks,
  set: (_val) => {
    // Reordering not supported by backend yet
  },
});

const handleChange = (evt: {
  added?: { element: Task };
  removed?: { element: Task };
  moved?: { element: Task; oldIndex: number; newIndex: number };
}) => {
  if (evt.added) {
    emit("task-move", {
      taskId: evt.added.element.id,
      newStatus: id,
    });
  }
};
</script>

<template>
  <div
    class="flex h-full w-80 min-w-[20rem] flex-col rounded-lg border border-gray-200 bg-gray-50/50"
  >
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-gray-100 p-3">
      <div class="flex items-center gap-2">
        <div :class="cn('h-2.5 w-2.5 rounded-full', color)" />
        <h3 class="text-sm font-semibold text-gray-900">{{ title }}</h3>
        <span
          class="rounded-full bg-gray-200 px-2 py-0.5 text-xs font-medium text-gray-500"
        >
          {{ tasks.length }}
        </span>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-2">
      <draggable
        v-model="internalTasks"
        group="tasks"
        item-key="id"
        class="h-full space-y-2"
        ghost-class="opacity-50"
        drag-class="cursor-grabbing"
        @change="handleChange"
      >
        <template #item="{ element: task }">
          <TasksBoardCard :task="task" @click="emit('task-click', task)" />
        </template>

        <template #footer>
          <div
            v-if="tasks.length === 0"
            class="flex h-24 items-center justify-center rounded-lg border-2 border-dashed border-gray-200 text-xs text-gray-400"
          >
            Empty
          </div>
        </template>
      </draggable>
    </div>
  </div>
</template>
