<script setup lang="ts">
import { ArrowLeft, Loader2, Send } from "lucide-vue-next";
import { TaskStatus, TaskPriority } from "~/types";

definePageMeta({
  layout: "dashboard",
});

const route = useRoute();
const router = useRouter();
const projectId = computed(() => route.params.projectId as string);

const { createTask } = useTasks(projectId.value);
const { members } = useProjectMembers(projectId.value);

const isSaving = ref(false);
const formData = ref({
  title: "",
  description: "",
  status: TaskStatus.TODO,
  priority: TaskPriority.P2,
  assignee_id: null as string | null,
  due_date: "",
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
    await createTask(payload);
    router.push(`/projects/${projectId.value}`);
  } catch (err) {
    console.error("Failed to create task", err);
  } finally {
    isSaving.value = false;
  }
};
</script>

<template>
  <div class="mx-auto max-w-4xl space-y-6 pb-20">
    <NuxtLink
      :to="`/projects/${projectId}`"
      class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900"
    >
      <ArrowLeft class="h-4 w-4" /> Back to project
    </NuxtLink>

    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-gray-900">Create New Task</h1>
    </div>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <!-- Main Content -->
      <div class="space-y-6 lg:col-span-2">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div class="space-y-2">
            <UiBaseLabel for="title">Title</UiBaseLabel>
            <UiBaseInput
              id="title"
              v-model="formData.title"
              placeholder="What needs to be done?"
              required
              class-name="h-12 text-lg"
              :disabled="isSaving"
            />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="description">Description</UiBaseLabel>
            <UiBaseMarkdownEditor
              v-model="formData.description"
              placeholder="Add more details... (Markdown supported)"
              :disabled="isSaving"
            />
          </div>

          <div class="flex justify-end gap-3 pt-4">
            <UiBaseButton
              type="button"
              variant="ghost"
              :disabled="isSaving"
              @click="router.back()"
            >
              Cancel
            </UiBaseButton>
            <UiBaseButton
              type="submit"
              class-name="bg-blue-600 px-8 text-white hover:bg-blue-700"
              :disabled="isSaving || !formData.title.trim()"
            >
              <Loader2 v-if="isSaving" class="mr-2 h-4 w-4 animate-spin" />
              <Send v-else class="mr-2 h-4 w-4" />
              Create Task
            </UiBaseButton>
          </div>
        </form>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6 lg:border-l lg:border-gray-100 lg:pl-8">
        <div class="space-y-4">
          <div class="space-y-2">
            <UiBaseLabel>Status</UiBaseLabel>
            <TasksSelectorsStatusSelector v-model="formData.status" />
          </div>

          <div class="space-y-2">
            <UiBaseLabel>Priority</UiBaseLabel>
            <TasksSelectorsPrioritySelector v-model="formData.priority" />
          </div>

          <div class="space-y-2">
            <UiBaseLabel>Assignee</UiBaseLabel>
            <TasksSelectorsAssigneeSelector
              v-model="formData.assignee_id"
              :members="members"
            />
          </div>

          <div class="space-y-2">
            <UiBaseLabel for="due_date">Due Date</UiBaseLabel>
            <UiBaseInput
              id="due_date"
              v-model="formData.due_date"
              type="date"
              :disabled="isSaving"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
