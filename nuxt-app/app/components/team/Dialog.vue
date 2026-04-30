<script setup lang="ts">
interface Props {
  isOpen?: boolean;
}

const { isOpen = false } = defineProps<Props>();

const emit = defineEmits(["close"]);

const name = ref("");
const description = ref("");
const { createTeam } = useTeams();
const isLoading = ref(false);

const handleSubmit = async () => {
  if (!name.value.trim()) return;

  isLoading.value = true;
  try {
    await createTeam(name.value, description.value);
    name.value = "";
    description.value = "";
    emit("close");
  } catch (err) {
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <UiBaseDialog title="Create Team" :is-open="isOpen" @close="emit('close')">
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div class="space-y-2">
        <UiBaseLabel for="team-name">Team Name</UiBaseLabel>
        <UiBaseInput
          id="team-name"
          v-model="name"
          placeholder="e.g. Engineering"
          required
          :disabled="isLoading"
        />
      </div>
      <div class="space-y-2">
        <UiBaseLabel for="team-desc">Description</UiBaseLabel>
        <UiBaseTextArea
          id="team-desc"
          v-model="description"
          placeholder="Brief description of the team..."
          :disabled="isLoading"
        />
      </div>
      <div class="flex justify-end gap-2 pt-4">
        <UiBaseButton
          type="button"
          variant="ghost"
          :disabled="isLoading"
          @click="emit('close')"
        >
          Cancel
        </UiBaseButton>
        <UiBaseButton type="submit" :disabled="isLoading || !name.trim()">
          {{ isLoading ? "Creating..." : "Create Team" }}
        </UiBaseButton>
      </div>
    </form>
  </UiBaseDialog>
</template>
