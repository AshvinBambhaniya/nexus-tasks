<script setup lang="ts">
import { Loader2, MessageSquare, Send } from "lucide-vue-next";
import { TaskStatus, TaskPriority, type Task } from "~/types";
import { useUsersStore } from "~/stores/user";

interface Props {
  isOpen?: boolean;
  task?: Task | null;
  projectId?: string;
}

const { isOpen = false, task = null, projectId = "" } = defineProps<Props>();

const emit = defineEmits(["close"]);

const userStore = useUsersStore();
const currentUserId = computed(() => userStore.userData?.id);

const activeProjectId = computed(() => projectId || task?.project_id || "");

const { createTask, updateTask, deleteTask } = useTasks(activeProjectId.value);
const { members } = useProjectMembers(activeProjectId.value);

// For existing tasks, use useTask hook for comments
const taskId = computed(() => task?.id || "");
const {
  comments,
  isLoading: detailsLoading,
  createComment,
  deleteComment,
} = useTask(taskId.value);

const isSaving = ref(false);
const newComment = ref("");
const isCommenting = ref(false);

const formData = ref({
  title: "",
  description: "",
  status: TaskStatus.TODO,
  priority: TaskPriority.P2,
  assignee_id: null as string | null,
  due_date: "" as string | undefined,
});

watchEffect(() => {
  if (task) {
    formData.value = {
      title: task.title,
      description: task.description || "",
      status: task.status,
      priority: task.priority,
      assignee_id: task.assignee_id || null,
      due_date: task.due_date ? task.due_date.split("T")[0] : "",
    };
  } else {
    formData.value = {
      title: "",
      description: "",
      status: TaskStatus.TODO,
      priority: TaskPriority.P2,
      assignee_id: null,
      due_date: "",
    };
  }
});

const handleSubmit = async () => {
  if (!formData.value.title.trim()) return;
  isSaving.value = true;
  try {
    const payload = {
      ...formData.value,
      assignee_id: formData.value.assignee_id || undefined,
      due_date: formData.value.due_date || undefined,
    };
    if (task) {
      await updateTask(task.id, payload);
    } else {
      await createTask(payload);
    }
    emit("close");
  } catch (err) {
    console.error(err);
  } finally {
    isSaving.value = false;
  }
};

const handleDelete = async () => {
  if (!task) return;
  if (!confirm("Are you sure you want to delete this task?")) return;

  isSaving.value = true;
  try {
    await deleteTask(task.id);
    emit("close");
  } catch (err) {
    console.error(err);
  } finally {
    isSaving.value = false;
  }
};

const handleAddComment = async () => {
  if (!newComment.value.trim()) return;
  isCommenting.value = true;
  try {
    await createComment(newComment.value);
    newComment.value = "";
  } catch (err) {
    console.error(err);
  } finally {
    isCommenting.value = false;
  }
};
</script>

<template>
  <UiBaseDialog
    :is-open="isOpen"
    :title="task ? 'Task Details' : 'Create Task'"
    class-name="max-w-4xl"
    @close="emit('close')"
  >
    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <!-- Main Content: Form -->
      <div class="space-y-6 lg:col-span-2">
        <form class="space-y-4" @submit.prevent="handleSubmit">
          <div class="space-y-2">
            <UiBaseLabel for="title">Title</UiBaseLabel>
            <UiBaseInput
              id="title"
              v-model="formData.title"
              required
              placeholder="e.g. Implement authentication"
              :disabled="isSaving"
            />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="description">Description</UiBaseLabel>
            <UiBaseTextArea
              id="description"
              v-model="formData.description"
              placeholder="Add more details..."
              :rows="5"
              :disabled="isSaving"
            />
          </div>

          <div class="flex justify-between pt-4">
            <UiBaseButton
              v-if="task"
              type="button"
              variant="destructive"
              size="sm"
              :disabled="isSaving"
              @click="handleDelete"
            >
              Delete Task
            </UiBaseButton>
            <div v-else />

            <div class="flex gap-2">
              <UiBaseButton
                type="button"
                variant="ghost"
                size="sm"
                :disabled="isSaving"
                @click="emit('close')"
              >
                Cancel
              </UiBaseButton>
              <UiBaseButton type="submit" size="sm" :disabled="isSaving">
                {{
                  isSaving ? "Saving..." : task ? "Save Changes" : "Create Task"
                }}
              </UiBaseButton>
            </div>
          </div>
        </form>

        <!-- Comments Section -->
        <div v-if="task" class="border-t border-gray-100 pt-8">
          <div class="mb-6 flex items-center gap-2">
            <MessageSquare class="h-5 w-5 text-gray-400" />
            <h3 class="text-lg font-semibold text-gray-900">Activity</h3>
          </div>

          <div class="space-y-6">
            <!-- Comment Form -->
            <div class="flex gap-4">
              <UiBaseAvatar
                :fallback="userStore.userData?.email?.[0].toUpperCase() || '?'"
                class-name="mt-1 h-10 w-10 border border-gray-100 shadow-sm"
              />
              <div class="flex-1 space-y-3">
                <UiBaseTextArea
                  v-model="newComment"
                  placeholder="Write a comment... (Markdown supported)"
                  :rows="3"
                  class-name="bg-gray-50 focus:bg-white transition-colors"
                  :disabled="isCommenting"
                />
                <div class="flex justify-end">
                  <UiBaseButton
                    size="sm"
                    :disabled="isCommenting || !newComment.trim()"
                    @click="handleAddComment"
                  >
                    <template v-if="isCommenting">
                      <Loader2 class="mr-2 h-4 w-4 animate-spin" />
                    </template>
                    <template v-else>
                      <Send class="mr-2 h-4 w-4" />
                    </template>
                    Comment
                  </UiBaseButton>
                </div>
              </div>
            </div>

            <!-- Comment List -->
            <div
              v-if="detailsLoading && comments.length === 0"
              class="flex justify-center py-8"
            >
              <Loader2 class="h-6 w-6 animate-spin text-gray-400" />
            </div>
            <div v-else class="space-y-8">
              <TasksCommentItem
                v-for="comment in comments"
                :key="comment.id"
                :comment="comment"
                :current-user-id="currentUserId"
                @delete="deleteComment"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Sidebar: Meta Information -->
      <div class="space-y-6 lg:border-l lg:border-gray-100 lg:pl-8">
        <div class="space-y-4">
          <div class="space-y-2">
            <UiBaseLabel for="status">Status</UiBaseLabel>
            <TasksSelectorsStatusSelector v-model="formData.status" />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="priority">Priority</UiBaseLabel>
            <TasksSelectorsPrioritySelector v-model="formData.priority" />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="assignee">Assignee</UiBaseLabel>
            <TasksSelectorsAssigneeSelector 
              v-model="formData.assignee_id"
              :members="members" 
            />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="due_date">Due Date</UiBaseLabel>
            <UiBaseInput id="due_date" v-model="formData.due_date" type="date" />
          </div>
        </div>
      </div>
    </div>
  </UiBaseDialog>
</template>
