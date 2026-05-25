<script setup lang="ts">
import {
  MessageSquare,
  Circle,
  CheckCircle2,
  Clock,
  AlertCircle,
  Calendar,
} from "lucide-vue-next";
import { formatDistanceToNow, format, isPast, isToday } from "date-fns";
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
  [TaskStatus.TODO]: "text-blue-500 dark:text-blue-400",
  [TaskStatus.IN_PROGRESS]: "text-purple-500 dark:text-purple-400",
  [TaskStatus.DONE]: "text-green-500 dark:text-green-400",
};

const priorityColors = {
  [TaskPriority.P0]:
    "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-100 dark:border-red-900/50",
  [TaskPriority.P1]:
    "text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-950/30 border-orange-100 dark:border-orange-900/50",
  [TaskPriority.P2]:
    "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-100 dark:border-blue-900/50",
  [TaskPriority.P3]: "text-muted-foreground bg-muted border-border",
};

const icon = computed(() => statusIcons[task.status] || Circle);
</script>

<template>
  <div
    class="group border-border hover:bg-muted/50 flex cursor-pointer items-start gap-3 border-b p-4 transition-colors last:border-0"
    @click="emit('click', task)"
  >
    <div :class="cn('mt-1 shrink-0', statusColors[task.status])">
      <component :is="icon" class="h-5 w-5" />
    </div>

    <div class="min-w-0 flex-1 space-y-1">
      <div class="flex flex-wrap items-center gap-2">
        <NuxtLink
          :to="`/projects/${projectId}/tasks/${task.id}`"
          class="text-foreground hover:text-primary font-bold transition-colors"
          @click.stop
        >
          {{ task.title }}
        </NuxtLink>
        <UiBaseBadge
          variant="outline"
          :class="
            cn(
              'h-5 px-1.5 text-[10px] font-bold uppercase',
              priorityColors[task.priority]
            )
          "
        >
          {{ task.priority }}
        </UiBaseBadge>
      </div>

      <div
        class="text-muted-foreground flex flex-wrap items-center gap-2 text-xs"
      >
        <span>#{{ task.number }}</span>
        <span>
          opened {{ formatDistanceToNow(new Date(task.created_at)) }} ago
        </span>
        <template v-if="task.status === TaskStatus.DONE && task.completed_at">
          <div
            class="flex items-center gap-1 font-medium text-green-600 dark:text-green-400"
          >
            <CheckCircle2 class="h-3 w-3" />
            <span>{{ format(new Date(task.completed_at), "MMM d") }}</span>
          </div>
        </template>
        <template v-else-if="task.due_date">
          <div
            :class="
              cn(
                'flex items-center gap-1',
                isPast(new Date(task.due_date)) &&
                  !isToday(new Date(task.due_date)) &&
                  task.status !== TaskStatus.DONE
                  ? 'font-medium text-red-600 dark:text-red-400'
                  : 'text-muted-foreground'
              )
            "
          >
            <Calendar class="h-3 w-3" />
            <span>{{ format(new Date(task.due_date), "MMM d") }}</span>
          </div>
        </template>
        <div v-if="task.assignee" class="flex items-center gap-1">
          <UiBaseAvatar
            :fallback="task.assignee.email[0].toUpperCase()"
            class-name="h-4 w-4 text-[8px]"
          />
          <span class="hover:text-primary cursor-pointer">
            {{ task.assignee.email.split("@")[0] }}
          </span>
        </div>
      </div>
    </div>

    <div class="text-muted-foreground/70 flex shrink-0 items-center gap-4">
      <div
        v-if="(task.comment_count ?? 0) > 0"
        class="flex items-center gap-1 text-xs"
      >
        <MessageSquare class="h-3.5 w-3.5" />
        <span>{{ task.comment_count }}</span>
      </div>
      <div class="hidden items-center gap-1 text-xs group-hover:flex">
        <MessageSquare class="h-3.5 w-3.5" />
        <span>Details</span>
      </div>
    </div>
  </div>
</template>
