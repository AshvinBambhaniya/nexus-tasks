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
const estimateHours = ref(0);
const estimateMins = ref(0);

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
    const totalMins =
      (estimateHours.value || 0) * 60 + (estimateMins.value || 0);
    const payload = {
      ...formData.value,
      assignee_id: formData.value.assignee_id || undefined,
      due_date: formData.value.due_date || undefined,
      estimated_minutes: totalMins > 0 ? totalMins : undefined,
    };
    await createTask(payload);
    router.push(`/projects/${projectId.value}`);
  } catch (err) {
    console.error("Failed to create task", err);
  } finally {
    isSaving.value = false;
  }
};
const isDrafting = ref(false);
const handleMagicDraft = async () => {
  if (!formData.value.title.trim()) {
    alert("Please enter a brief title first to guide the AI.");
    return;
  }

  isDrafting.value = true;
  try {
    const data = await useMutation<{ content: string }>(
      "/api/v2/ai/draft-task",
      {
        method: "POST",
        body: { title: formData.value.title },
      }
    );

    formData.value.description = data.content;
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to generate AI draft"));
  } finally {
    isDrafting.value = false;
  }
};
</script>

<template>
  <div class="mx-auto max-w-4xl space-y-6 pb-20">
    <NuxtLink
      :to="`/projects/${projectId}`"
      class="text-muted-foreground hover:text-foreground inline-flex items-center gap-1 text-sm"
    >
      <ArrowLeft class="h-4 w-4" /> Back to project
    </NuxtLink>

    <div class="flex items-center justify-between">
      <h1 class="text-foreground text-3xl font-bold">Create New Task</h1>
    </div>

    <div class="grid grid-cols-1 gap-8 lg:grid-cols-3">
      <!-- Main Content -->
      <div class="space-y-6 lg:col-span-2">
        <form class="space-y-6" @submit.prevent="handleSubmit">
          <div class="space-y-2">
            <div class="flex items-center justify-between">
              <UiBaseLabel for="title">Title</UiBaseLabel>
              <UiBaseButton
                type="button"
                variant="ghost"
                size="sm"
                class="text-primary hover:bg-primary/10 h-6 px-2 text-[11px] font-bold transition-colors"
                :disabled="isDrafting || !formData.title.trim()"
                @click="handleMagicDraft"
              >
                <Loader2
                  v-if="isDrafting"
                  class="mr-1.5 h-3 w-3 animate-spin"
                />
                <span v-else class="mr-1.5">✨</span>
                {{ isDrafting ? "Drafting..." : "Magic Draft" }}
              </UiBaseButton>
            </div>
            <UiBaseInput
              id="title"
              v-model="formData.title"
              placeholder="What needs to be done?"
              required
              class-name="h-12 text-lg"
              :disabled="isSaving || isDrafting"
            />
          </div>

          <div class="relative space-y-2">
            <UiBaseLabel for="description">Description</UiBaseLabel>
            <div class="relative">
              <UiBaseMarkdownEditor
                v-model="formData.description"
                placeholder="Add more details... (Markdown supported)"
                :disabled="isSaving || isDrafting"
                :class="[isDrafting ? 'pointer-events-none opacity-50' : '']"
              />
              <div
                v-if="isDrafting"
                class="pointer-events-none absolute inset-0 z-10 flex items-center justify-center"
              >
                <div
                  class="bg-background/80 ring-border flex items-center gap-2 rounded-full px-4 py-2 text-sm font-medium shadow-sm ring-1 backdrop-blur-sm"
                >
                  <Loader2 class="text-primary h-4 w-4 animate-spin" />
                  Generating breakdown...
                </div>
              </div>
            </div>
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
              class-name="bg-primary px-8 text-primary-foreground hover:bg-primary/90"
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
      <div class="lg:border-border space-y-6 lg:border-l lg:pl-8">
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

          <div class="space-y-2">
            <UiBaseLabel>Estimated Time</UiBaseLabel>
            <div class="flex items-center gap-2">
              <UiBaseInput
                v-model.number="estimateHours"
                type="number"
                :min="0"
                :max="999"
                placeholder="0"
                class-name="w-20"
                :disabled="isSaving"
              />
              <span class="text-muted-foreground text-xs font-medium">hrs</span>
              <UiBaseInput
                v-model.number="estimateMins"
                type="number"
                :min="0"
                :max="59"
                placeholder="0"
                class-name="w-20"
                :disabled="isSaving"
              />
              <span class="text-muted-foreground text-xs font-medium"
                >mins</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
