<script setup lang="ts">
interface Props {
  isOpen?: boolean;
}

const { isOpen = false } = defineProps<Props>();

const emit = defineEmits(["close"]);

const { createWorkspace } = useWorkspaces();
const name = ref("");
const isLoading = ref(false);
const error = ref<string | null>(null);

const handleSubmit = async () => {
  if (!name.value.trim()) return;

  isLoading.value = true;
  error.value = null;
  try {
    await createWorkspace(name.value);
    name.value = "";
    emit("close");
  } catch (err: unknown) {
    error.value = getApiErrorMessage(err, "Failed to create workspace");
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <UiBaseDialog
    :is-open="isOpen"
    title="Create Workspace"
    description="Workspaces are where you can organize your teams and projects."
    @close="emit('close')"
  >
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div class="space-y-2">
        <label for="name" class="text-sm font-medium text-gray-700">
          Workspace Name
        </label>
        <UiBaseInput
          id="name"
          v-model="name"
          placeholder="e.g. Acme Corp or Personal"
          required
          :disabled="isLoading"
        />
        <p v-if="error" class="text-sm text-red-600">
          {{ error }}
        </p>
      </div>

      <div class="flex justify-end gap-3 pt-4">
        <UiBaseButton
          type="button"
          variant="ghost"
          :disabled="isLoading"
          @click="emit('close')"
        >
          Cancel
        </UiBaseButton>
        <UiBaseButton type="submit" :disabled="isLoading || !name.trim()">
          {{ isLoading ? "Creating..." : "Create Workspace" }}
        </UiBaseButton>
      </div>
    </form>
  </UiBaseDialog>
</template>
