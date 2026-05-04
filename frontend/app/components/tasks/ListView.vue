<script setup lang="ts">
import { format } from "date-fns";
import { TaskStatus, TaskPriority, type Task } from "~/types";

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
  [TaskPriority.P0]: "text-red-600 bg-red-50 border-red-200",
  [TaskPriority.P1]: "text-orange-600 bg-orange-50 border-orange-200",
  [TaskPriority.P2]: "text-blue-600 bg-blue-50 border-blue-100",
  [TaskPriority.P3]: "text-gray-600 bg-gray-50 border-gray-100",
};
</script>

<template>
  <div
    v-if="tasks.length === 0"
    class="flex h-64 items-center justify-center rounded-lg border border-dashed border-gray-200 bg-gray-50 text-gray-500"
  >
    No tasks found. Create one to get started!
  </div>

  <div
    v-else
    class="overflow-hidden rounded-md border border-gray-200 bg-white"
  >
    <div class="overflow-x-auto">
      <table class="w-full text-left text-sm">
        <thead class="bg-gray-50 text-xs font-semibold text-gray-700 uppercase">
          <tr>
            <th class="px-4 py-3">Title</th>
            <th class="px-4 py-3">Status</th>
            <th class="px-4 py-3">Priority</th>
            <th class="px-4 py-3">Due / Completed</th>
            <th class="px-4 py-3 text-right">Created</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr
            v-for="task in tasks"
            :key="task.id"
            class="group cursor-pointer transition-colors hover:bg-gray-50"
            @click="emit('task-click', task)"
          >
            <td class="px-4 py-3 font-semibold text-gray-900">
              <span class="mr-1.5 text-gray-400">#{{ task.number }}</span>
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
            <td class="px-4 py-3 text-gray-500">
              <template
                v-if="task.status === TaskStatus.DONE && task.completed_at"
              >
                <span class="font-medium text-green-600">
                  {{ format(new Date(task.completed_at), "MMM d, yyyy") }}
                </span>
              </template>
              <template v-else-if="task.due_date">
                {{ format(new Date(task.due_date), "MMM d, yyyy") }}
              </template>
              <template v-else>
                <span class="text-gray-300">-</span>
              </template>
            </td>
            <td class="px-4 py-3 text-right text-gray-500">
              {{ new Date(task.created_at).toLocaleDateString() }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
