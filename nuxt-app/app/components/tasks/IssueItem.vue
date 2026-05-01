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
  [TaskStatus.BACKLOG]: "text-gray-400",
  [TaskStatus.TODO]: "text-blue-500",
  [TaskStatus.IN_PROGRESS]: "text-purple-500",
  [TaskStatus.DONE]: "text-green-500",
};

const priorityColors = {
  [TaskPriority.P0]: "text-red-600 bg-red-50 border-red-100",
  [TaskPriority.P1]: "text-orange-600 bg-orange-50 border-orange-100",
  [TaskPriority.P2]: "text-blue-600 bg-blue-50 border-blue-100",
  [TaskPriority.P3]: "text-gray-500 bg-gray-50 border-gray-100",
};

const icon = computed(() => statusIcons[task.status] || Circle);
</script>

<template>
  <div
    class="group flex cursor-pointer items-start gap-3 border-b border-gray-100 p-4 transition-colors last:border-0 hover:bg-gray-50"
    @click="emit('click', task)"
  >
    <div :class="cn('mt-1 shrink-0', statusColors[task.status])">
      <component :is="icon" class="h-5 w-5" />
    </div>

    <div class="min-w-0 flex-1 space-y-1">
      <div class="flex flex-wrap items-center gap-2">
        <NuxtLink
          :to="`/projects/${projectId}/tasks/${task.id}`"
          class="font-bold text-gray-900 transition-colors hover:text-blue-600"
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

      <div class="flex flex-wrap items-center gap-2 text-xs text-gray-500">
        <span>#{{ task.id }}</span>
        <span>
          opened {{ formatDistanceToNow(new Date(task.created_at)) }} ago
        </span>
        <template v-if="task.status === TaskStatus.DONE && task.completed_at">
          <div class="flex items-center gap-1 font-medium text-green-600">
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
                  ? 'font-medium text-red-600'
                  : 'text-gray-500'
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
          <span class="cursor-pointer hover:text-blue-600">
            {{ task.assignee.email.split("@")[0] }}
          </span>
        </div>
      </div>
    </div>

    <div class="flex shrink-0 items-center gap-4 text-gray-400">
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
