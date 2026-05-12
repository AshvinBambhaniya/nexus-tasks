<script setup lang="ts">
import type { Project } from "~/types";

interface Props {
  isOpen?: boolean;
  project?: Project | null;
}

const { isOpen = false, project = null } = defineProps<Props>();

const emit = defineEmits(["close"]);

const name = ref("");
const description = ref("");
const { createProject, refresh: refreshProjects } = useProjects();
// Note: In a full implementation we'd also use updateProject from useProjects
// For now we'll focus on create parity with the dashboard needs
const isLoading = ref(false);

watchEffect(() => {
  if (project) {
    name.value = project.name;
    description.value = project.description || "";
  } else {
    name.value = "";
    description.value = "";
  }
});

const handleSubmit = async () => {
  if (!name.value.trim()) return;

  isLoading.value = true;
  try {
    if (project) {
      await useMutation(`/api/v2/projects/${project.id}`, {
        method: "PATCH",
        body: { name: name.value, description: description.value },
      });
    } else {
      await createProject(name.value, description.value);
    }
    await refreshProjects();
    emit("close");
  } catch (err) {
    console.error(err);
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <UiBaseDialog
    :is-open="isOpen"
    :title="project ? 'Edit Project' : 'Create Project'"
    @close="emit('close')"
  >
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div class="space-y-2">
        <UiBaseLabel for="project-name">Project Name</UiBaseLabel>
        <UiBaseInput
          id="project-name"
          v-model="name"
          placeholder="e.g. Website Redesign"
          required
          :disabled="isLoading"
        />
      </div>
      <div class="space-y-2">
        <UiBaseLabel for="project-desc">Description</UiBaseLabel>
        <UiBaseTextArea
          id="project-desc"
          v-model="description"
          placeholder="Brief description of the project..."
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
          {{
            isLoading
              ? "Saving..."
              : project
                ? "Save Changes"
                : "Create Project"
          }}
        </UiBaseButton>
      </div>
    </form>
  </UiBaseDialog>
</template>
