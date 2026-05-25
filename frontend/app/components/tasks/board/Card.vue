<script setup lang="ts">
import { MessageSquare, Clock, CheckCircle2 } from "lucide-vue-next";
import { format, isPast, isToday } from "date-fns";
import { cn } from "~/utils/cn";
import { TaskStatus, TaskPriority, type Task } from "~/types";

interface Props {
  task: Task;
}

const { task } = defineProps<Props>();

const emit = defineEmits(["click"]);

const priorityColors: Record<TaskPriority, string> = {
  [TaskPriority.P0]:
    "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-100 dark:border-red-900/50",
  [TaskPriority.P1]:
    "text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-950/30 border-orange-100 dark:border-orange-900/50",
  [TaskPriority.P2]:
    "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-100 dark:border-blue-900/50",
  [TaskPriority.P3]: "text-muted-foreground bg-muted border-border",
};
</script>

<template>
  <div class="cursor-grab active:cursor-grabbing" @click="emit('click', task)">
    <UiBaseCard class="border-border transition-all hover:shadow-md">
      <div class="space-y-3 p-3">
        <div class="flex items-start justify-between gap-2">
          <span
            class="text-foreground line-clamp-2 text-sm leading-tight font-medium"
          >
            {{ task.title }}
          </span>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <UiBaseBadge
              variant="outline"
              :class="
                cn(
                  'h-5 border px-1.5 py-0 text-[10px] font-bold uppercase',
                  priorityColors[task.priority]
                )
              "
            >
              {{ task.priority }}
            </UiBaseBadge>
            <template
              v-if="task.status === TaskStatus.DONE && task.completed_at"
            >
              <div
                class="flex items-center gap-1 text-[10px] font-medium text-green-600 dark:text-green-400"
              >
                <CheckCircle2 class="h-3 w-3" />
                <span>{{ format(new Date(task.completed_at), "MMM d") }}</span>
              </div>
            </template>
            <template v-else-if="task.due_date">
              <div
                :class="
                  cn(
                    'flex items-center gap-1 text-[10px]',
                    isPast(new Date(task.due_date)) &&
                      !isToday(new Date(task.due_date))
                      ? 'font-medium text-red-600 dark:text-red-400'
                      : 'text-muted-foreground'
                  )
                "
              >
                <Clock class="h-3 w-3" />
                <span>{{ format(new Date(task.due_date), "MMM d") }}</span>
              </div>
            </template>
          </div>

          <div class="flex items-center gap-2">
            <div
              v-if="(task.comment_count ?? 0) > 0"
              class="text-muted-foreground flex items-center gap-1 text-xs"
            >
              <MessageSquare class="h-3 w-3" />
              <span>{{ task.comment_count }}</span>
            </div>
            <UiBaseAvatar
              v-if="task.assignee"
              :fallback="task.assignee.email[0].toUpperCase()"
              class-name="h-5 w-5 border border-border text-[9px]"
            />
          </div>
        </div>
      </div>
    </UiBaseCard>
  </div>
</template>
