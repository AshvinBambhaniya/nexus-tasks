<script setup lang="ts">
import {
  Archive,
  Trash2,
  Check,
  Loader2,
  AlertTriangle,
  RefreshCw,
} from "lucide-vue-next";
import type { Project } from "~/types";
import { useProjects } from "~/composables/useProjects";

interface Props {
  project: Project;
}

const props = defineProps<Props>();

const router = useRouter();
const { refresh: refreshProjects } = useProjects();

const projectName = ref(props.project.name);
const projectDescription = ref(props.project.description || "");
const isUpdating = ref(false);
const isSuccess = ref(false);

// Keep local state in sync with prop if it changes
watch(
  () => props.project,
  (newProject) => {
    projectName.value = newProject.name;
    projectDescription.value = newProject.description || "";
  }
);

const handleUpdateGeneral = async () => {
  isUpdating.value = true;
  isSuccess.value = false;
  try {
    await useMutation(`/api/v2/projects/${props.project.id}`, {
      method: "PATCH",
      body: { name: projectName.value, description: projectDescription.value },
    });
    isSuccess.value = true;
    setTimeout(() => (isSuccess.value = false), 3000);
    // No need to reload, the parent state should update or we can emit
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to update project"));
  } finally {
    isUpdating.value = false;
  }
};

const handleToggleArchive = async () => {
  if (
    !confirm(
      `Are you sure you want to ${props.project.is_archived ? "restore" : "archive"} this project?`
    )
  )
    return;

  isUpdating.value = true;
  try {
    await useMutation(`/api/v2/projects/${props.project.id}`, {
      method: "PATCH",
      body: { is_archived: !props.project.is_archived },
    });
    await refreshProjects();
    // Redirect if archived to avoid staying on an archived project page usually
    if (!props.project.is_archived) {
      router.push("/projects");
    } else {
      window.location.reload();
    }
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to update project status"));
  } finally {
    isUpdating.value = false;
  }
};

const handleDelete = async () => {
  if (
    !confirm(
      "Are you sure you want to permanently delete this project? This action cannot be undone and will remove all tasks and data."
    )
  )
    return;

  isUpdating.value = true;
  try {
    await useMutation(`/api/v2/projects/${props.project.id}`, {
      method: "DELETE",
    });
    await refreshProjects();
    router.push("/projects");
  } catch (err) {
    alert(getApiErrorMessage(err, "Failed to delete project"));
  } finally {
    isUpdating.value = false;
  }
};
</script>

<template>
  <div
    class="animate-in fade-in slide-in-from-bottom-4 max-w-4xl space-y-8 pb-10 duration-300"
  >
    <!-- General Settings -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-semibold text-gray-900">General Settings</h2>
        <p class="text-sm text-gray-500">
          Update project name and description.
        </p>
      </div>

      <UiBaseCard
        class="border-gray-100 shadow-sm transition-shadow hover:shadow-md"
      >
        <form class="space-y-6 p-8" @submit.prevent="handleUpdateGeneral">
          <div class="space-y-2">
            <UiBaseLabel
              for="projectName"
              class="text-xs font-bold tracking-wide text-gray-500 uppercase"
              >Project Name</UiBaseLabel
            >
            <UiBaseInput
              id="projectName"
              v-model="projectName"
              required
              placeholder="E.g., Marketing Campaign 2024"
              class="transition-all focus:ring-2 focus:ring-blue-100"
            />
          </div>

          <div class="space-y-2">
            <UiBaseLabel
              for="projectDesc"
              class="text-xs font-bold tracking-wide text-gray-500 uppercase"
              >Description</UiBaseLabel
            >
            <UiBaseTextArea
              id="projectDesc"
              v-model="projectDescription"
              placeholder="What is this project about?"
              :rows="4"
              class="transition-all focus:ring-2 focus:ring-blue-100"
            />
          </div>

          <div
            class="flex items-center justify-end border-t border-gray-50 pt-6"
          >
            <UiBaseButton
              type="submit"
              :disabled="isUpdating"
              class="min-w-[140px] shadow-sm active:scale-95"
            >
              <Loader2 v-if="isUpdating" class="mr-2 h-4 w-4 animate-spin" />
              <Check v-else-if="isSuccess" class="mr-2 h-4 w-4 text-white" />
              {{
                isUpdating
                  ? "Saving..."
                  : isSuccess
                    ? "Changes Saved"
                    : "Save Changes"
              }}
            </UiBaseButton>
          </div>
        </form>
      </UiBaseCard>
    </section>

    <!-- Archive Section -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-semibold text-gray-900">Project Lifecycle</h2>
        <p class="text-sm text-gray-500">
          Archive or restore this project from your active view.
        </p>
      </div>

      <UiBaseCard class="border-gray-100 shadow-sm">
        <div
          class="flex flex-col items-center justify-between gap-6 p-8 sm:flex-row"
        >
          <div class="space-y-1 text-center sm:text-left">
            <h3 class="font-bold text-gray-900">
              {{ project.is_archived ? "Restore Project" : "Archive Project" }}
            </h3>
            <p class="text-sm text-gray-500">
              {{
                project.is_archived
                  ? "Move this project back to the active list to continue working on it."
                  : "Archived projects are hidden but keep all their data. You can restore them anytime."
              }}
            </p>
          </div>
          <UiBaseButton
            variant="outline"
            :disabled="isUpdating"
            class="min-w-[160px] border-gray-200 transition-all hover:bg-gray-50 active:scale-95"
            @click="handleToggleArchive"
          >
            <RefreshCw v-if="project.is_archived" class="mr-2 h-4 w-4" />
            <Archive v-else class="mr-2 h-4 w-4 text-amber-600" />
            {{ project.is_archived ? "Restore Project" : "Archive Project" }}
          </UiBaseButton>
        </div>
      </UiBaseCard>
    </section>

    <!-- Danger Zone -->
    <section>
      <div class="mb-4 flex items-center gap-2">
        <AlertTriangle class="h-5 w-5 text-red-600" />
        <h2 class="text-lg font-semibold text-red-900">Danger Zone</h2>
      </div>

      <div
        class="overflow-hidden rounded-xl border border-red-100 bg-red-50/20 shadow-sm ring-1 ring-red-100/50"
      >
        <div
          class="flex flex-col items-center justify-between gap-6 p-8 sm:flex-row"
        >
          <div class="space-y-1 text-center sm:text-left">
            <h3 class="font-bold text-red-900">Delete this project</h3>
            <p class="text-sm text-red-700 opacity-80">
              Once you delete a project, there is no going back. Please be
              certain.
            </p>
          </div>
          <UiBaseButton
            variant="destructive"
            :disabled="isUpdating"
            class="min-w-[160px] bg-red-600 shadow-md shadow-red-200 hover:bg-red-700 active:scale-95"
            @click="handleDelete"
          >
            <Trash2 class="mr-2 h-4 w-4" /> Delete Project
          </UiBaseButton>
        </div>
      </div>
    </section>
  </div>
</template>
