<script setup lang="ts">
import { computed } from "vue";
import { Square, CheckSquare, Zap, Flame, AlertCircle } from "lucide-vue-next";
import type { Task } from "~/types";
import { TaskStatus } from "~/types";

const props = defineProps<{
  task: Task;
}>();

defineEmits<{
  (e: "click" | "toggle-status", task: Task): void;
}>();

const isCompleted = computed(() => props.task.status === TaskStatus.DONE);

const priorityBadge = computed(() => {
  switch (props.task.priority) {
    case "P0":
      return { class: "text-red-500 bg-red-500/10", icon: Flame };
    case "P1":
      return { class: "text-yellow-500 bg-yellow-500/10", icon: Zap };
    case "P2":
      return { class: "text-orange-500 bg-orange-500/10", icon: Zap };
    case "P3":
      return { class: "text-blue-500 bg-blue-500/10", icon: AlertCircle };
    default:
      return { class: "text-muted-foreground bg-muted/50", icon: AlertCircle };
  }
});

const formatDate = (dateString?: string) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  const mm = String(date.getMonth() + 1).padStart(2, "0");
  const dd = String(date.getDate()).padStart(2, "0");
  const yy = String(date.getFullYear()).slice(-2);
  return `${mm}/${dd}/${yy}`;
};
</script>

<template>
  <div
    class="border-border/50 hover:bg-muted/30 group grid cursor-pointer grid-cols-[32px_1fr_200px_100px_100px] items-center gap-4 border-b px-4 py-2 transition-colors last:border-b-0"
    @click="$emit('click', task)"
  >
    <!-- Checkbox -->
    <div class="flex items-center justify-center">
      <button
        class="text-muted-foreground/60 hover:text-foreground transition-colors focus:outline-none"
        @click.stop="$emit('toggle-status', task)"
      >
        <CheckSquare v-if="isCompleted" class="text-muted-foreground h-4 w-4" />
        <Square v-else class="h-4 w-4" />
      </button>
    </div>

    <!-- Task Name -->
    <div
      :class="[
        'truncate text-sm transition-all duration-200',
        isCompleted
          ? 'text-muted-foreground line-through opacity-60'
          : 'text-foreground font-medium',
      ]"
    >
      {{ task.title }}
    </div>

    <!-- Project Name -->
    <div>
      <span
        v-if="task.project"
        class="inline-flex items-center rounded bg-blue-500/10 px-2 py-0.5 text-xs font-medium text-blue-500/80"
      >
        {{ task.project.name }}
      </span>
      <span v-else class="text-muted-foreground/50 text-xs">--</span>
    </div>

    <!-- Priority -->
    <div class="flex items-center">
      <span
        class="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-[10px] font-bold"
        :class="priorityBadge.class"
      >
        <component :is="priorityBadge.icon" class="h-3 w-3" />
        {{ task.priority }}
      </span>
    </div>

    <!-- Due Date -->
    <div class="text-muted-foreground text-sm">
      {{ formatDate(task.due_date) }}
    </div>
  </div>
</template>
