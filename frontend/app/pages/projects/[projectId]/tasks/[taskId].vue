<script setup lang="ts">
import {
  ArrowLeft,
  Loader2,
  Trash2,
  Clock,
  CheckCircle2,
  Send,
  Sparkles,
} from "lucide-vue-next";
import { formatDistanceToNow, format } from "date-fns";
import VueMarkdown from "vue-markdown-render";
import type { TaskPriority } from "~/types";
import { TaskStatus } from "~/types";
import { useUsersStore } from "~/stores/user";
import RichTextEditor from "~/components/ui/RichTextEditor.vue";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const router = useRouter();
const projectId = computed(() => route.params.projectId as string);
const taskId = computed(() => route.params.taskId as string);

const userStore = useUsersStore();
const currentUserId = computed(() => userStore.userData?.id);

const formatComment = (htmlContent: string) => {
  if (!currentUserId.value || !htmlContent) return htmlContent || '';
  return htmlContent.replaceAll(`data-id="${currentUserId.value}"`, `data-id="${currentUserId.value}" data-is-me="true"`);
};

const {
  task,
  comments,
  isLoading,
  refreshTask,
  refreshComments,
  createComment,
  deleteComment,
  summarizeComments,
} = useTask(projectId.value, taskId.value);

const { updateTask, deleteTask } = useTasks(projectId.value);
const { members } = useProjectMembers(projectId.value);

const commentContent = ref("");
const isSubmittingComment = ref(false);
const isEditingTitle = ref(false);
const titleValue = ref("");
const titleInput = ref<HTMLInputElement | null>(null);

const aiSummary = ref("");
const isSummarizing = ref(false);

watch(
  task,
  (newTask) => {
    if (newTask && !titleValue.value) {
      titleValue.value = newTask.title;
    }
  },
  { immediate: true }
);

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
    titleValue.value = task.value?.title || "";
    return;
  }
  await updateTask(taskId.value, { title: titleValue.value });
  isEditingTitle.value = false;
  await refreshTask();
};

const enableTitleEdit = () => {
  isEditingTitle.value = true;
  nextTick(() => {
    titleInput.value?.focus();
  });
};

const handleDelete = async () => {
  if (!confirm("Are you sure you want to delete this task?")) return;
  await deleteTask(taskId.value);
  router.push(`/projects/${projectId.value}`);
};

const handleSubmitComment = async (content: string, mentionedUserIds: string[]) => {
  isSubmittingComment.value = true;
  try {
    await createComment(content, mentionedUserIds);
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

const handleSummarize = async () => {
  isSummarizing.value = true;
  try {
    const summary = await summarizeComments();
    if (summary) {
      aiSummary.value = summary;
    }
  } catch (err) {
    console.error("Failed to summarize comments", err);
  } finally {
    isSummarizing.value = false;
  }
};

const authorName = computed(() => {
  if (!task.value?.author) return "Unknown";
  return task.value.author.full_name || task.value.author.email.split("@")[0];
});
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center justify-center">
    <Loader2 class="text-primary h-8 w-8 animate-spin" />
  </div>
  <div v-else-if="!task" class="text-muted-foreground p-8 text-center">
    Task not found
  </div>
  <div v-else class="mx-auto mt-4 flex max-w-6xl flex-col gap-8 pb-20">
    <!-- Header Row (Breadcrumbs & Actions) -->
    <div class="border-border flex items-center justify-between border-b pb-4">
      <div class="flex items-center gap-4">
        <NuxtLink
          :to="`/projects/${projectId}`"
          class="text-muted-foreground hover:text-foreground inline-flex items-center gap-1.5 text-sm font-medium transition-colors"
        >
          <ArrowLeft class="h-4 w-4" /> Back to board
        </NuxtLink>
        <span class="text-muted-foreground/50">/</span>
        <span class="text-muted-foreground font-mono text-sm tracking-tight"
          >TASK-{{ task.number }}</span
        >
      </div>

      <div class="flex items-center gap-2">
        <button
          v-if="task.status === TaskStatus.DONE"
          class="text-muted-foreground hover:text-foreground bg-muted hover:bg-muted/80 flex items-center gap-1.5 rounded-md px-3 py-1.5 text-sm font-medium transition-colors"
          @click="handleUpdateStatus(TaskStatus.TODO)"
        >
          <ArrowLeft class="h-4 w-4" />
          Reopen
        </button>
        <button
          v-else
          class="flex items-center gap-1.5 rounded-md border border-green-600/20 bg-green-600/10 px-3 py-1.5 text-sm font-medium text-green-600 transition-colors hover:bg-green-600/20"
          @click="handleUpdateStatus(TaskStatus.DONE)"
        >
          <CheckCircle2 class="h-4 w-4" />
          Complete
        </button>

        <button
          class="text-muted-foreground hover:text-destructive hover:bg-destructive/10 ml-2 rounded-md p-1.5 transition-colors"
          title="Delete Task"
          @click="handleDelete"
        >
          <Trash2 class="h-4 w-4" />
        </button>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="grid grid-cols-1 gap-12 lg:grid-cols-4">
      <!-- Left Column: Title, Description, Activity -->
      <div class="space-y-10 lg:col-span-3">
        <!-- Title -->
        <div class="group relative">
          <input
            v-if="isEditingTitle"
            ref="titleInput"
            v-model="titleValue"
            class="text-foreground border-primary/50 focus:border-primary w-full border-b bg-transparent pb-1 text-3xl font-bold tracking-tight transition-colors outline-none"
            @keyup.enter="handleUpdateTitle"
            @blur="handleUpdateTitle"
            @keyup.esc="isEditingTitle = false"
          />
          <h1
            v-else
            class="text-foreground cursor-text pr-8 text-3xl leading-tight font-bold tracking-tight"
            @click="enableTitleEdit"
          >
            {{ task.title }}
          </h1>
        </div>

        <!-- Description -->
        <div>
          <VueMarkdown
            v-if="task.description"
            :source="task.description"
            class="prose dark:prose-invert prose-sm prose-pre:bg-muted/50 prose-pre:border prose-pre:border-border max-w-none"
          />
          <div v-else class="text-muted-foreground/60 text-sm italic">
            Add a description...
          </div>
        </div>

        <hr class="border-border" />

        <!-- Activity Stream -->
        <div class="space-y-6">
          <h3 class="text-foreground text-sm font-semibold">Activity</h3>

          <div class="relative space-y-6 pl-4">
            <!-- Continuous Line -->
            <div
              class="bg-border pointer-events-none absolute top-4 bottom-12 left-[27px] w-[1px]"
            />

            <!-- Creation Event -->
            <div class="relative z-10 flex items-start gap-4">
              <div class="bg-background mt-0.5">
                <UiBaseAvatar
                  :fallback="authorName?.[0]?.toUpperCase() || '?'"
                  class-name="h-7 w-7 text-[10px] border border-border bg-muted text-muted-foreground"
                />
              </div>
              <div class="flex-1 pt-1">
                <p class="text-muted-foreground text-sm">
                  <span class="text-foreground font-medium">{{
                    authorName
                  }}</span>
                  created this issue
                  {{ formatDistanceToNow(new Date(task.created_at)) }} ago
                </p>
              </div>
            </div>

            <!-- AI Summary Block -->
            <div
              v-if="comments.length > 0"
              class="relative z-10 my-6 flex items-start gap-4"
            >
              <div class="bg-background mt-0.5">
                <div class="border-border bg-muted flex h-7 w-7 items-center justify-center rounded-full border">
                  <Sparkles class="text-primary h-3.5 w-3.5" />
                </div>
              </div>
              <div class="flex-1">
                <div v-if="aiSummary" class="border-primary/20 bg-primary/5 rounded-lg border p-4 shadow-sm relative overflow-hidden">
                  <div class="absolute inset-0 bg-gradient-to-br from-primary/10 to-transparent pointer-events-none" />
                  <div class="relative z-10">
                    <div class="flex items-center gap-2 mb-2">
                      <Sparkles class="h-4 w-4 text-primary" />
                      <h4 class="text-sm font-semibold text-primary">AI Thread Summary</h4>
                    </div>
                    <VueMarkdown
                      :source="aiSummary"
                      class="prose dark:prose-invert prose-sm max-w-none text-foreground/90 marker:text-primary"
                    />
                  </div>
                </div>
                <button
                  v-else
                  @click="handleSummarize"
                  :disabled="isSummarizing"
                  class="border-border hover:bg-muted/50 bg-background text-muted-foreground hover:text-foreground flex items-center gap-2 rounded-full border px-4 py-1.5 text-xs font-medium transition-colors"
                >
                  <Loader2 v-if="isSummarizing" class="text-primary h-3.5 w-3.5 animate-spin" />
                  <Sparkles v-else class="text-primary h-3.5 w-3.5" />
                  {{ isSummarizing ? "Summarizing thread..." : "✨ Catch me up" }}
                </button>
              </div>
            </div>

            <!-- Comments -->
            <div
              v-for="comment in comments"
              :key="comment.id"
              class="relative z-10 flex items-start gap-4"
            >
              <div class="bg-background mt-0.5">
                <UiBaseAvatar
                  :fallback="
                    (comment.author?.full_name ||
                      comment.author?.email ||
                      '?')?.[0]?.toUpperCase() || '?'
                  "
                  class-name="h-7 w-7 text-[10px] border border-border"
                />
              </div>
              <div class="flex-1">
                <div
                  class="group/comment overflow-hidden rounded-lg border border-border/60 bg-muted/20 transition-colors"
                >
                  <div
                    class="flex items-center justify-between border-b border-border/40 bg-muted/10 px-3 py-2"
                  >
                    <span class="text-muted-foreground text-xs">
                      <span class="text-foreground font-medium">{{
                        comment.author?.full_name || comment.author?.email
                      }}</span>
                      commented
                      {{ formatDistanceToNow(new Date(comment.created_at)) }}
                      ago
                    </span>
                    <button
                      v-if="comment.author_id === currentUserId"
                      class="text-muted-foreground hover:text-destructive opacity-0 transition-opacity group-hover/comment:opacity-100"
                      @click="handleDeleteComment(comment.id)"
                    >
                      <Trash2 class="h-3.5 w-3.5" />
                    </button>
                  </div>
                  <div class="px-3 py-2.5">
                    <div
                      class="prose dark:prose-invert prose-sm max-w-none"
                      v-html="formatComment(comment.content)"
                    ></div>
                  </div>
                </div>
              </div>
            </div>

            <!-- New Comment Box -->
            <div class="relative z-10 flex items-start gap-4 pt-4">
              <div class="bg-background mt-1">
                <UiBaseAvatar
                  :fallback="
                    userStore.userData?.email?.[0]?.toUpperCase() || '?'
                  "
                  class-name="h-7 w-7 text-[10px] border border-border"
                />
              </div>
              <div class="flex-1">
                  <RichTextEditor
                    :project-id="projectId"
                    :is-submitting="isSubmittingComment"
                    @submit="handleSubmitComment"
                  />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Properties Sidebar -->
      <div
        class="lg:border-border/50 space-y-8 lg:col-span-1 lg:border-l lg:pl-8"
      >
        <div class="space-y-6">
          <h3 class="text-foreground text-sm font-semibold">Properties</h3>

          <div class="space-y-4">
            <!-- Status -->
            <div class="space-y-2">
              <span
                class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
              >
                Status
              </span>
              <TasksSelectorsStatusSelector
                :model-value="task.status"
                @update:model-value="handleUpdateStatus"
              />
            </div>

            <!-- Priority -->
            <div class="space-y-2">
              <span
                class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
              >
                Priority
              </span>
              <TasksSelectorsPrioritySelector
                :model-value="task.priority"
                @update:model-value="handleUpdatePriority"
              />
            </div>

            <!-- Assignee -->
            <div class="space-y-2">
              <span
                class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
              >
                Assignee
              </span>
              <TasksSelectorsAssigneeSelector
                :model-value="task.assignee_id"
                :members="members"
                @update:model-value="handleUpdateAssignee"
              />
            </div>

            <!-- Due Date -->
            <div class="space-y-2">
              <span
                class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
              >
                Due Date
              </span>
              <div
                class="border-border bg-background hover:bg-muted focus-within:ring-ring flex items-center rounded-md border px-3 py-2 transition-colors focus-within:ring-2"
              >
                <Clock class="text-muted-foreground mr-2 h-4 w-4 shrink-0" />
                <input
                  type="date"
                  :value="task.due_date ? task.due_date.split('T')[0] : ''"
                  class="text-foreground/80 w-full flex-1 cursor-pointer bg-transparent text-sm font-medium outline-none"
                  @input="
                    (e) =>
                      handleUpdateDueDate((e.target as HTMLInputElement).value)
                  "
                />
              </div>
            </div>
          </div>
        </div>

        <div v-if="task.completed_at" class="border-border/50 border-t pt-6">
          <h3 class="text-foreground mb-3 text-sm font-semibold">Resolution</h3>
          <div
            class="flex items-center gap-2 text-sm font-medium text-green-600"
          >
            <CheckCircle2 class="h-4 w-4" />
            Completed {{ format(new Date(task.completed_at), "MMM d, yyyy") }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
