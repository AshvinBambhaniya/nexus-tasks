<script setup lang="ts">
import { Archive, Trash2 } from "lucide-vue-next";
import type { Project } from "~/types";

interface Props {
  project: Project;
}

const { project } = defineProps<Props>();

const router = useRouter();
const { refresh: refreshProjects } = useProjects();
const isUpdating = ref(false);

const handleToggleArchive = async () => {
  isUpdating.value = true;
  try {
    await useMutation(`/api/v1/projects/${project.id}`, {
      method: "PATCH",
      body: { is_archived: !project.is_archived },
    });
    await refreshProjects();
    // In a real app we'd refresh the current project data too
    window.location.reload(); 
  } catch (err) {
    console.error(err);
  } finally {
    isUpdating.value = false;
  }
};

const handleDelete = async () => {
  if (!confirm("Are you sure you want to delete this project? This cannot be undone.")) return;
  
  isUpdating.value = true;
  try {
    await useMutation(`/api/v1/projects/${project.id}`, {
      method: "DELETE",
    });
    await refreshProjects();
    router.push("/projects");
  } catch (err) {
    console.error(err);
  } finally {
    isUpdating.value = false;
  }
};
</script>

<template>
  <div class="max-w-4xl space-y-6">
    <UiBaseCard>
      <div class="border-b border-gray-100 p-6">
        <h3 class="text-lg font-semibold text-gray-900">Archive Project</h3>
        <p class="text-sm text-gray-500">
          Archived projects are hidden from the active list but can still be accessed.
        </p>
      </div>
      <div class="p-6">
        <UiBaseButton
          variant="outline"
          :disabled="isUpdating"
          @click="handleToggleArchive"
        >
          <Archive class="mr-2 h-4 w-4" />
          {{ project.is_archived ? "Restore Project" : "Archive Project" }}
        </UiBaseButton>
      </div>
    </UiBaseCard>

    <div class="rounded-lg border border-red-100 bg-red-50/50 p-6">
      <h3 class="text-lg font-medium text-red-900">Danger Zone</h3>
      <p class="mt-1 text-sm text-red-600">
        Deleting this project will permanently remove all associated tasks and data.
      </p>
      <div class="mt-4">
        <UiBaseButton
          variant="destructive"
          size="sm"
          :disabled="isUpdating"
          @click="handleDelete"
        >
          <Trash2 class="mr-2 h-4 w-4" />
          Delete Project
        </UiBaseButton>
      </div>
    </div>
  </div>
</template>
