<script setup lang="ts">
import { Archive, Trash2, Loader2, RefreshCw } from "lucide-vue-next";
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
  <div class="mx-auto max-w-4xl space-y-10 pb-10">
    <!-- General Settings -->
    <section class="space-y-4">
      <div>
        <h2 class="text-foreground text-xl font-semibold tracking-tight">
          General Settings
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Manage your project's basic information and identification.
        </p>
      </div>

      <div
        class="border-border bg-card overflow-hidden rounded-lg border shadow-sm"
      >
        <form @submit.prevent="handleUpdateGeneral">
          <div class="space-y-6 p-6">
            <div class="space-y-3">
              <label
                for="projectName"
                class="text-foreground text-sm font-medium"
              >
                Project Name
              </label>
              <input
                id="projectName"
                v-model="projectName"
                required
                placeholder="E.g., Marketing Campaign 2024"
                class="bg-background border-border focus:ring-primary focus:border-primary block h-9 w-full rounded-md border px-3 text-sm transition-all focus:ring-1 focus:outline-none sm:w-[400px]"
              />
            </div>

            <div class="space-y-3">
              <label
                for="projectDesc"
                class="text-foreground text-sm font-medium"
              >
                Description
              </label>
              <textarea
                id="projectDesc"
                v-model="projectDescription"
                placeholder="What is this project about?"
                :rows="3"
                class="bg-background border-border focus:ring-primary focus:border-primary block w-full resize-y rounded-md border px-3 py-2 text-sm transition-all focus:ring-1 focus:outline-none"
              />
            </div>
          </div>

          <div
            class="border-border bg-muted/30 flex items-center justify-between border-t px-6 py-4"
          >
            <p class="text-muted-foreground text-xs">
              Please use 32 characters at maximum for the project name.
            </p>
            <button
              type="submit"
              :disabled="isUpdating"
              class="bg-primary text-primary-foreground hover:bg-primary/90 flex h-9 min-w-[100px] items-center justify-center rounded-md px-4 text-sm font-medium transition-colors disabled:opacity-50"
            >
              <Loader2 v-if="isUpdating" class="h-4 w-4 animate-spin" />
              <span v-else>{{ isSuccess ? "Saved" : "Save" }}</span>
            </button>
          </div>
        </form>
      </div>
    </section>

    <!-- Archive Section -->
    <section class="space-y-4">
      <div>
        <h2 class="text-foreground text-xl font-semibold tracking-tight">
          Archive Project
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Deactivate this project and hide it from active views.
        </p>
      </div>

      <div
        class="border-border bg-card overflow-hidden rounded-lg border shadow-sm"
      >
        <div
          class="flex flex-col justify-between gap-6 p-6 sm:flex-row sm:items-center"
        >
          <div class="max-w-[600px] space-y-1">
            <h3 class="text-foreground text-sm font-medium">
              {{ project.is_archived ? "Restore Project" : "Archive Project" }}
            </h3>
            <p class="text-muted-foreground text-sm">
              {{
                project.is_archived
                  ? "Move this project back to the active list to continue working on it. All data remains intact."
                  : "Archived projects are hidden from the active list but keep all their data. You can restore them anytime."
              }}
            </p>
          </div>
          <button
            :disabled="isUpdating"
            class="bg-secondary text-secondary-foreground hover:bg-secondary/80 border-border flex h-9 items-center justify-center gap-2 rounded-md border px-4 text-sm font-medium whitespace-nowrap transition-colors disabled:opacity-50"
            @click="handleToggleArchive"
          >
            <RefreshCw v-if="project.is_archived" class="h-4 w-4" />
            <Archive v-else class="h-4 w-4" />
            {{ project.is_archived ? "Restore" : "Archive" }}
          </button>
        </div>
      </div>
    </section>

    <!-- Danger Zone -->
    <section class="space-y-4">
      <div>
        <h2 class="text-foreground text-xl font-semibold tracking-tight">
          Danger Zone
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Irreversible and destructive actions.
        </p>
      </div>

      <div
        class="border-destructive/30 bg-card relative overflow-hidden rounded-lg border shadow-sm"
      >
        <div class="bg-destructive/5 pointer-events-none absolute inset-0" />
        <div
          class="relative flex flex-col justify-between gap-6 p-6 sm:flex-row sm:items-center"
        >
          <div class="max-w-[600px] space-y-1">
            <h3 class="text-foreground text-sm font-medium">Delete Project</h3>
            <p class="text-muted-foreground text-sm">
              Once you delete a project, there is no going back. All associated
              tasks, configurations, and data will be permanently removed.
            </p>
          </div>
          <button
            :disabled="isUpdating"
            class="bg-destructive/10 text-destructive hover:bg-destructive hover:text-destructive-foreground border-destructive/30 flex h-9 items-center justify-center gap-2 rounded-md border px-4 text-sm font-medium whitespace-nowrap transition-colors disabled:opacity-50"
            @click="handleDelete"
          >
            <Trash2 class="h-4 w-4" />
            Delete Project
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
