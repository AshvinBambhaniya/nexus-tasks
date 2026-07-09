<script setup lang="ts">
import { ref } from "vue";
import {
  Loader2,
  CheckCircle2,
  Calendar,
  MoreVertical,
  Clock,
} from "lucide-vue-next";
import { format } from "date-fns";
import { TaskStatus } from "~/types";
import { useTask } from "~/composables/useTasks";
import { useUsersStore } from "~/stores/user";
import RichTextEditor from "./ui/RichTextEditor.vue";
import MarkdownViewer from "./ui/MarkdownViewer.vue";

const props = defineProps<{
  taskId: string;
  projectId: string;
}>();

const emit = defineEmits<{
  (e: "toggle-done", task: { id: string; status: TaskStatus }): void;
}>();

const { task, comments, createComment, isLoading, isError } = useTask(
  props.projectId,
  props.taskId
);

const userStore = useUsersStore();
const currentUserId = computed(() => userStore.userData?.id);

const formatComment = (htmlContent: string) => {
  if (!currentUserId.value || !htmlContent) return htmlContent || "";
  return htmlContent.replaceAll(
    `data-id="${currentUserId.value}"`,
    `data-id="${currentUserId.value}" data-is-me="true"`
  );
};

const isSubmitting = ref(false);

const submitComment = async (content: string, mentionedUserIds: string[]) => {
  isSubmitting.value = true;
  try {
    await createComment(content, mentionedUserIds);
  } catch (e) {
    console.error(e);
  } finally {
    isSubmitting.value = false;
  }
};

const getPrioColor = (prio: string | undefined) => {
  if (!prio) return "bg-muted text-muted-foreground";
  switch (prio) {
    case "P0":
      return "text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950/30 border-red-100 dark:border-red-900/50";
    case "P1":
      return "text-orange-600 dark:text-orange-400 bg-orange-50 dark:bg-orange-950/30 border-orange-100 dark:border-orange-900/50";
    case "P2":
      return "text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-950/30 border-blue-100 dark:border-blue-900/50";
    default:
      return "text-muted-foreground bg-muted border-border";
  }
};
</script>

<template>
  <div class="bg-background/30 flex h-full flex-col">
    <div v-if="isLoading" class="flex h-full items-center justify-center">
      <Loader2 class="text-muted-foreground h-8 w-8 animate-spin" />
    </div>
    <div
      v-else-if="isError || !task"
      class="flex h-full flex-col items-center justify-center text-center"
    >
      <p class="text-muted-foreground">Failed to load task details.</p>
    </div>
    <template v-else>
      <!-- Detail Header Toolbar -->
      <div
        class="border-border flex h-16 shrink-0 items-center justify-between border-b px-6"
      >
        <div class="flex items-center gap-3">
          <div
            class="bg-muted text-muted-foreground border-border flex h-6 w-6 items-center justify-center rounded-full border text-[10px] font-bold"
          >
            {{
              task.assignee?.full_name
                ? task.assignee.full_name.substring(0, 2).toUpperCase()
                : "U"
            }}
          </div>
          <span class="text-foreground text-sm font-medium">{{
            task.assignee?.full_name || "Unassigned"
          }}</span>
        </div>

        <div class="flex items-center gap-2">
          <button
            :class="[
              'flex h-8 items-center gap-2 rounded-md border px-3 text-xs font-medium transition-colors',
              task.status === TaskStatus.DONE
                ? 'border-border bg-muted text-foreground'
                : 'border-primary bg-primary/10 text-primary hover:bg-primary/20',
            ]"
            @click="emit('toggle-done', task)"
          >
            <CheckCircle2 class="h-3.5 w-3.5" />
            {{
              task.status === TaskStatus.DONE
                ? "Mark as Undone"
                : "Mark as Done"
            }}
          </button>
          <button
            class="border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground flex h-8 items-center gap-2 rounded-md border px-3 text-xs font-medium transition-colors"
          >
            <Clock class="h-3.5 w-3.5" /> Snooze
          </button>
          <button
            class="border-border bg-card text-muted-foreground hover:bg-muted hover:text-foreground flex h-8 w-8 items-center justify-center rounded-md border transition-colors"
          >
            <MoreVertical class="h-3.5 w-3.5" />
          </button>
        </div>
      </div>

      <!-- Detail Content -->
      <div class="custom-scrollbar flex-1 overflow-y-auto p-6 lg:p-10">
        <div class="mx-auto max-w-3xl space-y-8">
          <!-- Title & Meta -->
          <div>
            <div class="mb-4 flex flex-wrap items-center gap-2">
              <span
                class="bg-muted text-muted-foreground rounded px-2 py-0.5 font-mono text-[10px]"
                >NEX-{{ task.number }}</span
              >
              <span
                :class="[
                  'rounded px-2 py-0.5 text-[10px] font-bold tracking-wider uppercase',
                  getPrioColor(task.priority),
                ]"
              >
                {{ task.priority }}
              </span>
              <span
                class="bg-muted text-muted-foreground rounded px-2 py-0.5 text-[10px] font-bold tracking-wider uppercase"
              >
                {{ task.status?.replace("_", " ") || "UNKNOWN" }}
              </span>
              <span
                v-if="task.due_date"
                class="bg-muted text-muted-foreground flex items-center gap-1 rounded px-2 py-0.5 text-[10px] font-medium"
              >
                <Calendar class="h-3 w-3" />
                {{ format(new Date(task.due_date), "MMM d, yyyy") }}
              </span>
            </div>

            <h1
              class="text-foreground text-2xl font-bold tracking-tight lg:text-3xl"
            >
              {{ task.title }}
            </h1>
          </div>

          <!-- Description -->
          <div
            class="prose prose-sm dark:prose-invert text-muted-foreground max-w-none"
          >
            <MarkdownViewer
              v-if="task.description"
              :content="task.description"
            />
            <p v-else class="italic opacity-50">No description provided.</p>
          </div>

          <hr class="border-border" />

          <!-- Activity Feed (Real Data) -->
          <div>
            <h4
              class="text-muted-foreground mb-6 text-xs font-semibold tracking-widest uppercase"
            >
              Activity
            </h4>

            <div class="space-y-6">
              <!-- Activity Item: Created -->
              <div class="relative flex gap-4">
                <div
                  class="bg-border absolute top-8 bottom-0 left-4 w-px -translate-x-1/2"
                />
                <div
                  class="border-border bg-muted text-muted-foreground z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full border text-[10px] font-bold"
                >
                  +
                </div>
                <div class="flex-1 space-y-1 pt-1.5">
                  <p class="text-foreground text-sm">
                    <span class="font-medium">Task created</span>
                  </p>
                  <p class="text-muted-foreground text-xs">
                    {{
                      format(
                        new Date(task.created_at),
                        "MMM d, yyyy 'at' h:mm a"
                      )
                    }}
                  </p>
                </div>
              </div>

              <!-- Comments -->
              <div
                v-for="(comment, index) in comments"
                :key="comment.id"
                class="relative flex gap-4"
              >
                <div
                  v-if="index !== comments.length - 1"
                  class="bg-border absolute top-8 bottom-0 left-4 w-px -translate-x-1/2"
                />
                <div
                  class="border-border bg-muted text-muted-foreground z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full border text-[10px] font-bold"
                >
                  {{
                    comment.author?.full_name
                      ? comment.author.full_name.substring(0, 2).toUpperCase()
                      : "U"
                  }}
                </div>
                <div class="flex-1 space-y-2 pt-1.5 pb-2">
                  <p class="text-muted-foreground text-xs">
                    <span class="text-foreground font-medium">{{
                      comment.author?.full_name || "User"
                    }}</span>
                    commented •
                    {{ format(new Date(comment.created_at), "MMM d, h:mm a") }}
                  </p>
                  <MarkdownViewer
                    class="border-border bg-card text-foreground rounded-lg border p-3 text-sm transition-colors"
                    :content="formatComment(comment.content)"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Comment Box -->
          <div class="mt-8">
            <RichTextEditor
              :project-id="task.project_id"
              :is-submitting="isSubmitting"
              @submit="submitComment"
            />
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--muted);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: var(--border);
}
</style>
