<script setup lang="ts">
import draggable from "vuedraggable";
import { Plus, MoreHorizontal } from "lucide-vue-next";
import type { TaskStatus, Task } from "~/types";

interface Props {
  id: TaskStatus;
  title: string;
  tasks: Task[];
}

const { id, title, tasks } = defineProps<Props>();

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
    class="bg-muted/30 hover:border-border/50 flex h-full w-80 min-w-[20rem] flex-col rounded-xl border border-transparent transition-colors"
  >
    <!-- Header -->
    <div class="flex items-center justify-between p-3.5 pb-2">
      <div class="flex items-center gap-2">
        <h3
          class="text-foreground flex items-center gap-2 text-xs font-bold tracking-wider uppercase"
        >
          {{ title }}
        </h3>
        <span class="text-muted-foreground text-xs font-medium">
          {{ tasks.length }}
        </span>
      </div>
      <div class="flex items-center gap-1">
        <button
          class="text-muted-foreground hover:text-foreground hover:bg-muted rounded p-1 transition-colors"
        >
          <Plus class="h-4 w-4" />
        </button>
        <button
          class="text-muted-foreground hover:text-foreground hover:bg-muted rounded p-1 transition-colors"
        >
          <MoreHorizontal class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto px-3 pb-3">
      <draggable
        v-model="internalTasks"
        group="tasks"
        item-key="id"
        class="h-full min-h-[150px] space-y-3 pt-2"
        ghost-class="opacity-40"
        drag-class="cursor-grabbing"
        @change="handleChange"
      >
        <template #item="{ element: task }">
          <TasksBoardCard :task="task" @click="emit('task-click', task)" />
        </template>

        <template #footer>
          <div
            v-if="tasks.length === 0"
            class="text-muted-foreground/50 border-border/50 flex h-24 items-center justify-center rounded-lg border border-dashed text-xs"
          >
            Drop tasks here
          </div>
        </template>
      </draggable>
    </div>
  </div>
</template>
