<script setup lang="ts">
import {
  ArrowLeft,
  Loader2,
  Trash2,
  Settings2,
  Clock,
  User as UserIcon,
  CheckCircle2,
  Circle,
  AlertCircle,
  MessageSquare,
  Send,
} from "lucide-vue-next";
import { formatDistanceToNow } from "date-fns";
import VueMarkdown from "vue-markdown-render";
import { cn } from "~/utils/cn";
import { TaskStatus, TaskPriority } from "~/types";
import { useUsersStore } from "~/stores/user";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const router = useRouter();
const projectId = computed(() => route.params.projectId as string);
const taskId = computed(() => route.params.taskId as string);

const userStore = useUsersStore();
const currentUserId = computed(() => userStore.userData?.id);

const {
  task,
  comments,
  isLoading,
  refreshTask,
  refreshComments,
  createComment,
  deleteComment,
} = useTask(projectId.value, taskId.value);

const { updateTask, deleteTask } = useTasks(projectId.value);
const { members } = useProjectMembers(projectId.value);

const commentContent = ref("");
const isSubmittingComment = ref(false);
const isEditingTitle = ref(false);
const titleValue = ref("");

watch(task, (newTask) => {
  if (newTask && !titleValue.value) {
    titleValue.value = newTask.title;
  }
}, { immediate: true });

const handleUpdateStatus = async (status: TaskStatus) => {
  await updateTask(taskId.value, { status });
  await refreshTask();
};

const handleUpdatePriority = async (priority: TaskPriority) => {
  await updateTask(taskId.value, { priority });
  await refreshTask();
};

const handleUpdateAssignee = async (assigneeId: string | null) => {
  await updateTask(taskId.value, { assignee_id: assigneeId || undefined });
  await refreshTask();
};

const handleUpdateDueDate = async (dateStr: string) => {
  await updateTask(taskId.value, { due_date: dateStr || undefined });
  await refreshTask();
};

const handleUpdateTitle = async () => {
  if (!titleValue.value.trim() || titleValue.value === task.value?.title) {
    isEditingTitle.value = false;
    return;
  }
  await updateTask(taskId.value, { title: titleValue.value });
  isEditingTitle.value = false;
  await refreshTask();
};

const handleDelete = async () => {
  if (!confirm("Are you sure you want to delete this task?")) return;
  await deleteTask(taskId.value);
  router.push(`/projects/${projectId.value}`);
};

const handleSubmitComment = async () => {
  if (!commentContent.value.trim()) return;

  isSubmittingComment.value = true;
  try {
    await createComment(commentContent.value);
    commentContent.value = "";
    await refreshComments();
  } catch (err) {
    console.error("Failed to add comment", err);
  } finally {
    isSubmittingComment.value = false;
  }
};

const handleDeleteComment = async (id: string) => {
  if (!confirm("Delete this comment?")) return;
  try {
    await deleteComment(id);
    await refreshComments();
  } catch (err) {
    console.error("Failed to delete comment", err);
  }
};

const statusIcons: Record<TaskStatus, any> = {
  [TaskStatus.BACKLOG]: Clock,
  [TaskStatus.TODO]: Circle,
  [TaskStatus.IN_PROGRESS]: AlertCircle,
  [TaskStatus.DONE]: CheckCircle2,
};

const statusIcon = computed(() => (task.value ? statusIcons[task.value.status] || Circle : Circle));

const authorName = computed(() => {
  if (!task.value?.author) return "Unknown";
  return task.value.author.full_name || task.value.author.email.split("@")[0];
});
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center justify-center">
    <Loader2 class="h-8 w-8 animate-spin text-blue-500" />
  </div>
  <div v-else-if="!task" class="p-8 text-center text-gray-500">
    Task not found
  </div>
  <div v-else class="mx-auto max-w-6xl space-y-6 pb-20">
    <!-- Header -->
    <div class="space-y-4">
      <NuxtLink
        :to="`/projects/${projectId}`"
        class="mb-2 inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900"
      >
        <ArrowLeft class="h-4 w-4" /> Back to project
      </NuxtLink>

      <div class="flex flex-col gap-2">
        <div class="flex items-start justify-between gap-4">
          <div v-if="isEditingTitle" class="flex flex-1 gap-2">
            <UiBaseInput
              v-model="titleValue"
              class-name="h-12 text-2xl font-bold"
              auto-focus
              @keyup.enter="handleUpdateTitle"
              @blur="handleUpdateTitle"
            />
            <UiBaseButton @click="handleUpdateTitle">Save</UiBaseButton>
            <UiBaseButton variant="ghost" @click="isEditingTitle = false">
              Cancel
            </UiBaseButton>
          </div>
          <h1 v-else class="group flex items-center gap-2 text-3xl font-bold text-gray-900">
            {{ task.title }}
            <span class="font-normal text-gray-400">#{{ task.number }}</span>
            <button
              @click="isEditingTitle = true"
              class="rounded p-1 opacity-0 transition-all group-hover:opacity-100 hover:bg-gray-100"
            >
              <Settings2 class="h-4 w-4 text-gray-500" />
            </button>
          </h1>

          <div class="flex shrink-0 gap-2">
            <UiBaseButton
              v-if="task.status === TaskStatus.DONE"
              variant="outline"
              size="sm"
              class="whitespace-nowrap"
              @click="handleUpdateStatus(TaskStatus.TODO)"
            >
              <ArrowLeft class="mr-2 h-4 w-4" />
              Reopen Task
            </UiBaseButton>
            <UiBaseButton
              v-else
              variant="primary"
              size="sm"
              class="whitespace-nowrap !bg-green-600 !hover:bg-green-700 !shadow-none"
              @click="handleUpdateStatus(TaskStatus.DONE)"
            >
              <CheckCircle2 class="mr-2 h-4 w-4" />
              Complete Task
            </UiBaseButton>
            <UiBaseButton
              variant="destructive"
              size="sm"
              class="whitespace-nowrap"
              @click="handleDelete"
            >
              <Trash2 class="mr-2 h-4 w-4" />
              Delete Task
            </UiBaseButton>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3 border-b border-gray-200 pb-6">
          <UiBaseBadge
            :class="
              cn(
                'flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-medium transition-colors',
                task.status === TaskStatus.DONE
                  ? 'bg-green-100 text-green-700 hover:bg-green-200'
                  : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
              )
            "
          >
            <component :is="statusIcon" class="h-4 w-4" />
            {{ task.status.replace('_', ' ') }}
          </UiBaseBadge>
          <span class="flex items-center gap-1.5 text-sm font-medium text-gray-500">
            <UserIcon class="h-4 w-4" />
            <span class="font-semibold text-gray-900">
              {{ authorName }}
            </span>
            opened this task {{ formatDistanceToNow(new Date(task.created_at)) }} ago • {{ comments.length }} comments
          </span>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-4">
      <!-- Main Stream -->
      <div class="space-y-8 lg:col-span-3">
        <!-- Description -->
        <div class="flex gap-4">
          <UiBaseAvatar
            :fallback="authorName[0].toUpperCase()"
            class-name="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
          />
          <div class="flex-1">
            <div class="overflow-hidden rounded-lg border border-gray-200 bg-white shadow-sm">
              <div class="flex items-center justify-between border-b border-gray-200 bg-gray-50/50 px-4 py-2">
                <span class="text-sm font-semibold text-gray-700">
                  Description
                </span>
              </div>
              <VueMarkdown
                v-if="task.description"
                :source="task.description"
                class="prose prose-sm prose-pre:bg-gray-50 prose-pre:border prose-pre:border-gray-100 max-w-none p-4"
              />
              <div v-else class="p-4 text-gray-400 italic">No description provided.</div>
            </div>
          </div>
        </div>

        <!-- Comments List -->
        <div class="relative space-y-8 before:absolute before:top-0 before:bottom-0 before:left-[1.25rem] before:w-0.5 before:bg-gray-100">
          <TasksCommentItem
            v-for="comment in comments"
            :key="comment.id"
            :comment="comment"
            :current-user-id="currentUserId"
            @delete="handleDeleteComment"
          />
        </div>

        <!-- New Comment Box -->
        <div class="flex gap-4 border-t border-gray-200 pt-8">
          <UiBaseAvatar
            :fallback="userStore.userData?.email?.[0].toUpperCase() || '?'"
            class-name="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
          />
          <div class="flex-1">
            <form @submit.prevent="handleSubmitComment" class="space-y-4">
              <UiBaseMarkdownEditor
                v-model="commentContent"
                placeholder="Add a comment..."
                class-name="border-gray-200 shadow-sm"
              />
              <div class="flex justify-end">
                <UiBaseButton
                  type="submit"
                  :disabled="isSubmittingComment || !commentContent.trim()"
                  class-name="bg-blue-600 px-6 text-white hover:bg-blue-700"
                >
                  <Loader2 v-if="isSubmittingComment" class="mr-2 h-4 w-4 animate-spin" />
                  <Send v-else class="mr-2 h-4 w-4" />
                  Comment
                </UiBaseButton>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- Sidebar Controls -->
      <div class="space-y-8 lg:col-span-1">
        <div class="space-y-6">
          <div class="space-y-2 border-b border-gray-100 pb-4">
            <UiBaseLabel class-name="text-xs font-bold tracking-wider text-gray-500 uppercase">
              Due Date
            </UiBaseLabel>
            <UiBaseInput
              type="date"
              :model-value="task.due_date ? task.due_date.split('T')[0] : ''"
              @update:model-value="handleUpdateDueDate"
            />
          </div>

          <div class="space-y-2 border-b border-gray-100 pb-4">
            <UiBaseLabel class-name="text-xs font-bold tracking-wider text-gray-500 uppercase">
              Assignee
            </UiBaseLabel>
            <TasksSelectorsAssigneeSelector
              :model-value="task.assignee_id"
              :members="members"
              @update:model-value="handleUpdateAssignee"
            />
          </div>

          <div class="space-y-2 border-b border-gray-100 pb-4">
            <UiBaseLabel class-name="text-xs font-bold tracking-wider text-gray-500 uppercase">
              Status
            </UiBaseLabel>
            <TasksSelectorsStatusSelector
              :model-value="task.status"
              @update:model-value="handleUpdateStatus"
            />
          </div>

          <div v-if="task.completed_at" class="space-y-2 border-b border-gray-100 pb-4">
            <UiBaseLabel class-name="text-xs font-bold tracking-wider text-gray-500 uppercase">
              Completed On
            </UiBaseLabel>
            <div class="flex items-center gap-2 text-sm text-gray-700">
              <CheckCircle2 class="h-4 w-4 text-green-600" />
              {{ new Date(task.completed_at).toLocaleDateString() }}
            </div>
          </div>

          <div class="space-y-2 border-b border-gray-100 pb-4">
            <UiBaseLabel class-name="text-xs font-bold tracking-wider text-gray-500 uppercase">
              Priority
            </UiBaseLabel>
            <TasksSelectorsPrioritySelector
              :model-value="task.priority"
              @update:model-value="handleUpdatePriority"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
