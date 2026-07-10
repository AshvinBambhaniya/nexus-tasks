<script setup lang="ts">
import {
  MessageSquare,
  Clock,
  CheckCircle2,
  ChevronsUp,
  ChevronUp,
  Minus,
  ChevronDown,
  MoreHorizontal,
} from "lucide-vue-next";
import { format, isPast, isToday } from "date-fns";
import { cn } from "~/utils/cn";
import { TaskStatus, TaskPriority, type Task } from "~/types";

interface Props {
  task: Task;
}

const { task } = defineProps<Props>();

const emit = defineEmits(["click"]);

const priorityIcons = {
  [TaskPriority.P0]: ChevronsUp,
  [TaskPriority.P1]: ChevronUp,
  [TaskPriority.P2]: Minus,
  [TaskPriority.P3]: ChevronDown,
};

const priorityColors: Record<TaskPriority, string> = {
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

const priorityIcon = computed(() => priorityIcons[task.priority] || Minus);
</script>

<template>
  <div
    class="group cursor-grab active:cursor-grabbing"
    @click="emit('click', task)"
  >
    <div
      class="bg-card border-border/60 hover:border-border relative flex flex-col gap-3 overflow-hidden rounded-xl border p-3.5 shadow-sm transition-all hover:shadow-md"
    >
      <!-- Top Row: ID & Options -->
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-1.5">
          <component
            :is="priorityIcon"
            :class="cn('h-3.5 w-3.5', priorityColors[task.priority])"
          />
          <span class="text-muted-foreground font-mono text-xs"
            >TASK-{{ task.number }}</span
          >
        </div>
        <button
          class="text-muted-foreground hover:text-foreground opacity-0 transition-opacity group-hover:opacity-100"
          @click.stop
        >
          <MoreHorizontal class="h-3.5 w-3.5" />
        </button>
      </div>

      <!-- Title -->
      <h4 class="text-foreground text-sm leading-snug font-medium">
        {{ task.title }}
      </h4>

      <!-- Bottom Row: Meta -->
      <div class="mt-1 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <UiBaseBadge
            variant="outline"
            :class="
              cn(
                'border-border/50 bg-muted/30 h-5 px-1.5 py-0 text-[10px] font-medium',
                priorityColors[task.priority]
              )
            "
          >
            {{ priorityLabels[task.priority] }}
          </UiBaseBadge>

          <div
            v-if="(task.comment_count ?? 0) > 0"
            class="text-muted-foreground flex items-center gap-1 text-[10px] font-medium"
          >
            <MessageSquare class="h-3 w-3" />
            <span class="hidden sm:inline"
              >{{ task.comment_count }} comments</span
            >
            <span class="sm:hidden">{{ task.comment_count }}</span>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <template v-if="task.status === TaskStatus.DONE && task.completed_at">
            <div
              class="flex items-center gap-1 text-[10px] font-medium text-green-500"
            >
              <CheckCircle2 class="h-3.5 w-3.5" />
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
                    ? 'font-medium text-red-500'
                    : 'text-muted-foreground'
                )
              "
            >
              <Clock class="hidden h-3 w-3 sm:block" />
              <span>{{ format(new Date(task.due_date), "MMM d") }}</span>
            </div>
          </template>

          <UiBaseAvatar
            v-if="task.assignee"
            :fallback="task.assignee.email?.[0]?.toUpperCase() || '?'"
            class-name="h-5 w-5 text-[9px] border border-background ring-1 ring-border/50"
          />
        </div>
      </div>
    </div>
  </div>
</template>
