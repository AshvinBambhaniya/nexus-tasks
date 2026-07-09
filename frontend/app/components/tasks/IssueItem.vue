<script setup lang="ts">
import {
  MessageSquare,
  Circle,
  CheckCircle2,
  Clock,
  AlertCircle,
  ChevronsUp,
  ChevronUp,
  Minus,
  ChevronDown,
} from "lucide-vue-next";
import { format, isPast, isToday } from "date-fns";
import { cn } from "~/utils/cn";
import { TaskStatus, TaskPriority, type Task } from "~/types";

interface Props {
  task: Task;
  projectId: string;
}

const { task, projectId } = defineProps<Props>();

const emit = defineEmits(["click"]);

const statusIcons = {
  [TaskStatus.BACKLOG]: Clock,
  [TaskStatus.TODO]: Circle,
  [TaskStatus.IN_PROGRESS]: AlertCircle,
  [TaskStatus.DONE]: CheckCircle2,
};

const statusColors = {
  [TaskStatus.BACKLOG]: "text-muted-foreground",
  [TaskStatus.TODO]: "text-blue-500",
  [TaskStatus.IN_PROGRESS]: "text-purple-500",
  [TaskStatus.DONE]: "text-green-500",
};

const statusLabels = {
  [TaskStatus.BACKLOG]: "Backlog",
  [TaskStatus.TODO]: "To Do",
  [TaskStatus.IN_PROGRESS]: "In Progress",
  [TaskStatus.DONE]: "Done",
};

const priorityIcons = {
  [TaskPriority.P0]: ChevronsUp,
  [TaskPriority.P1]: ChevronUp,
  [TaskPriority.P2]: Minus,
  [TaskPriority.P3]: ChevronDown,
};

const priorityColors = {
  [TaskPriority.P0]: "text-red-500",
  [TaskPriority.P1]: "text-orange-500",
  [TaskPriority.P2]: "text-blue-500",
  [TaskPriority.P3]: "text-muted-foreground",
};

const priorityLabels = {
  [TaskPriority.P0]: "Urgent",
  [TaskPriority.P1]: "High",
  [TaskPriority.P2]: "Medium",
  [TaskPriority.P3]: "Low",
};

const statusIcon = computed(() => statusIcons[task.status] || Circle);
const priorityIcon = computed(() => priorityIcons[task.priority] || Minus);
</script>

<template>
  <tr
    class="group hover:bg-muted/50 cursor-pointer transition-colors"
    @click="emit('click', task)"
  >
    <!-- Checkbox -->
    <td class="px-4 py-3" @click.stop>
      <input
        type="checkbox"
        class="border-border text-primary focus:ring-primary/20 bg-background h-3.5 w-3.5 cursor-pointer rounded"
      />
    </td>

    <!-- ID -->
    <td class="text-muted-foreground px-2 py-3 font-mono text-xs">
      TASK-{{ task.number }}
    </td>

    <!-- Title -->
    <td class="px-2 py-3">
      <div class="flex items-center gap-2">
        <NuxtLink
          :to="`/projects/${projectId}/tasks/${task.id}`"
          class="text-foreground hover:text-primary text-sm font-medium transition-colors"
          @click.stop
        >
          {{ task.title }}
        </NuxtLink>
        <div
          v-if="(task.comment_count ?? 0) > 0"
          class="text-muted-foreground bg-muted border-border flex items-center gap-1 rounded-full border px-1.5 text-[10px]"
        >
          <MessageSquare class="h-2.5 w-2.5" />
          {{ task.comment_count }}
        </div>
      </div>
    </td>

    <!-- Assignee -->
    <td class="px-2 py-3">
      <div v-if="task.assignee" class="flex items-center gap-2">
        <UiBaseAvatar
          :fallback="task.assignee.email?.[0]?.toUpperCase() || '?'"
          class-name="h-5 w-5 text-[10px]"
        />
        <span
          class="text-muted-foreground group-hover:text-foreground max-w-[80px] truncate text-xs transition-colors"
        >
          {{ task.assignee.email.split("@")[0] }}
        </span>
      </div>
      <div v-else class="text-muted-foreground/50 text-xs">Unassigned</div>
    </td>

    <!-- Status -->
    <td class="px-2 py-3">
      <div class="flex items-center gap-1.5 text-xs">
        <component
          :is="statusIcon"
          :class="cn('h-3.5 w-3.5', statusColors[task.status])"
        />
        <span
          class="text-muted-foreground group-hover:text-foreground transition-colors"
          >{{ statusLabels[task.status] }}</span
        >
      </div>
    </td>

    <!-- Priority -->
    <td class="px-2 py-3">
      <div class="flex items-center gap-1.5 text-xs">
        <component
          :is="priorityIcon"
          :class="cn('h-3.5 w-3.5', priorityColors[task.priority])"
        />
        <span
          class="text-muted-foreground group-hover:text-foreground transition-colors"
          >{{ priorityLabels[task.priority] }}</span
        >
      </div>
    </td>

    <!-- Due Date -->
    <td class="px-2 py-3 pr-6 text-right">
      <template v-if="task.status === TaskStatus.DONE && task.completed_at">
        <div
          class="flex items-center justify-end gap-1.5 text-xs font-medium text-green-500"
        >
          <CheckCircle2 class="h-3.5 w-3.5" />
          <span>{{ format(new Date(task.completed_at), "MMM d") }}</span>
        </div>
      </template>
      <template v-else-if="task.due_date">
        <div
          :class="
            cn(
              'flex items-center justify-end gap-1.5 text-xs transition-colors',
              isPast(new Date(task.due_date)) &&
                !isToday(new Date(task.due_date)) &&
                task.status !== TaskStatus.DONE
                ? 'font-medium text-red-500'
                : 'text-muted-foreground group-hover:text-foreground'
            )
          "
        >
          <span>{{ format(new Date(task.due_date), "MMM d") }}</span>
        </div>
      </template>
      <span v-else class="text-muted-foreground/50 text-xs">-</span>
    </td>
  </tr>
</template>
