<script setup lang="ts">
import { format } from "date-fns";
import { TaskStatus, TaskPriority, type Task } from "~/types";
import { cn } from "~/utils/cn";

interface Props {
  tasks?: Task[];
}

const { tasks = [] } = defineProps<Props>();

const emit = defineEmits(["task-click"]);

const statusVariants: Record<
  TaskStatus,
  "default" | "secondary" | "destructive" | "outline"
> = {
  [TaskStatus.TODO]: "secondary",
  [TaskStatus.IN_PROGRESS]: "default",
  [TaskStatus.DONE]: "outline",
  [TaskStatus.BACKLOG]: "secondary",
};

const priorityColors: Record<TaskPriority, string> = {
  [TaskPriority.P0]: "text-destructive bg-destructive/10 border-destructive/20",
  [TaskPriority.P1]: "text-orange-500 bg-orange-500/10 border-orange-500/20",
  [TaskPriority.P2]: "text-primary bg-primary/10 border-primary/20",
  [TaskPriority.P3]: "text-muted-foreground bg-muted border-border",
};
</script>

<template>
  <div
    v-if="tasks.length === 0"
    class="border-border bg-muted text-muted-foreground flex h-64 items-center justify-center rounded-lg border border-dashed"
  >
    No tasks found. Create one to get started!
  </div>

  <div v-else class="border-border bg-card overflow-hidden rounded-md border">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-sm">
        <thead
          class="bg-muted text-foreground/80 text-xs font-semibold uppercase"
        >
          <tr>
            <th class="px-4 py-3">Title</th>
            <th class="px-4 py-3">Status</th>
            <th class="px-4 py-3">Priority</th>
            <th class="px-4 py-3">Due / Completed</th>
            <th class="px-4 py-3 text-right">Created</th>
          </tr>
        </thead>
        <tbody class="divide-border divide-y">
          <tr
            v-for="task in tasks"
            :key="task.id"
            class="group hover:bg-muted cursor-pointer transition-colors"
            @click="emit('task-click', task)"
          >
            <td class="text-foreground px-4 py-3 font-semibold">
              <span class="text-muted-foreground/70 mr-1.5"
                >#{{ task.number }}</span
              >
              {{ task.title }}
            </td>
            <td class="px-4 py-3">
              <UiBaseBadge :variant="statusVariants[task.status]">
                {{ task.status.replace("_", " ") }}
              </UiBaseBadge>
            </td>
            <td class="px-4 py-3">
              <span
                :class="
                  cn(
                    'inline-flex items-center rounded-md border px-2 py-0.5 text-[10px] font-bold uppercase',
                    priorityColors[task.priority]
                  )
                "
              >
                {{ task.priority }}
              </span>
            </td>
            <td class="text-muted-foreground px-4 py-3">
              <template
                v-if="task.status === TaskStatus.DONE && task.completed_at"
              >
                <span class="font-medium text-green-600 dark:text-green-500">
                  {{ format(new Date(task.completed_at), "MMM d, yyyy") }}
                </span>
              </template>
              <template v-else-if="task.due_date">
                {{ format(new Date(task.due_date), "MMM d, yyyy") }}
              </template>
              <template v-else>
                <span class="text-muted-foreground/30">-</span>
              </template>
            </td>
            <td class="text-muted-foreground px-4 py-3 text-right">
              {{ new Date(task.created_at).toLocaleDateString() }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
