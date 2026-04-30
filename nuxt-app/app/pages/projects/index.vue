<script setup lang="ts">
import { Folder, Search, Plus, ArrowRight } from "lucide-vue-next";
import { format } from "date-fns";

definePageMeta({ layout: "dashboard" });

const { projects, isLoading } = useProjects();
const searchQuery = ref("");
type FilterStatus = "active" | "archived" | "all";
const filterStatus = ref<FilterStatus>("active");
const isDialogOpen = ref(false);

const filteredProjects = computed(() => {
  return projects.value.filter((project) => {
    const matchesSearch = project.name
      .toLowerCase()
      .includes(searchQuery.value.toLowerCase());

    let matchesStatus = true;
    if (filterStatus.value === "active") {
      matchesStatus = !project.is_archived;
    } else if (filterStatus.value === "archived") {
      matchesStatus = project.is_archived;
    }

    return matchesSearch && matchesStatus;
  });
});
</script>

<template>
  <div class="flex h-full flex-col space-y-8">
    <!-- Header & Controls -->
    <div
      class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div>
        <h1 class="text-2xl font-bold text-gray-900">Projects</h1>
        <p class="mt-1 text-gray-500">
          Manage and organize your team's work.
        </p>
      </div>
      <UiBaseButton @click="isDialogOpen = true">
        <Plus class="mr-2 h-4 w-4" /> New Project
      </UiBaseButton>
    </div>

    <div
      class="flex flex-col items-center justify-between gap-4 border-b border-gray-200 pb-4 sm:flex-row"
    >
      <!-- Search -->
      <div class="relative w-full sm:w-72">
        <Search class="absolute top-2.5 left-3 h-4 w-4 text-gray-400" />
        <UiBaseInput
          v-model="searchQuery"
          placeholder="Search projects..."
          class-name="pl-9"
        />
      </div>

      <!-- Filter Tabs -->
      <div class="flex self-start rounded-lg bg-gray-100 p-1 sm:self-auto">
        <button
          v-for="status in (['active', 'archived', 'all'] as const)"
          :key="status"
          class="rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-all"
          :class="
            filterStatus === status
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-500 hover:text-gray-700'
          "
          @click="filterStatus = status"
        >
          {{ status }}
        </button>
      </div>
    </div>

    <!-- Grid Content -->
    <div class="flex-1 overflow-y-auto">
      <div
        v-if="isLoading"
        class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3"
      >
        <div
          v-for="i in 6"
          :key="i"
          class="h-48 animate-pulse rounded-xl border border-gray-200 bg-gray-100"
        />
      </div>
      <div
        v-else-if="filteredProjects.length === 0"
        class="flex h-64 flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 text-center"
      >
        <div
          class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100"
        >
          <Search v-if="searchQuery" class="h-6 w-6 text-gray-400" />
          <Folder v-else class="h-6 w-6 text-gray-400" />
        </div>
        <h3 class="text-lg font-medium text-gray-900">No projects found</h3>
        <p class="mt-1 max-w-sm text-sm text-gray-500">
          {{
            searchQuery
              ? `No projects matching "${searchQuery}"`
              : filterStatus === "archived"
                ? "No archived projects."
                : "Get started by creating a new project."
          }}
        </p>
        <UiBaseButton
          v-if="!searchQuery && filterStatus !== 'archived'"
          variant="outline"
          class="mt-4"
          @click="isDialogOpen = true"
        >
          Create Project
        </UiBaseButton>
      </div>
      <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="project in filteredProjects"
          :key="project.id"
          :to="`/projects/${project.id}`"
          class="block h-full"
        >
          <UiBaseCard
            class="group flex h-full cursor-pointer flex-col border-gray-200 p-6 transition-all hover:border-blue-200 hover:shadow-md"
          >
            <div class="mb-4 flex items-start justify-between">
              <div
                class="rounded-xl bg-blue-50 p-2.5 transition-colors group-hover:bg-blue-100"
              >
                <Folder class="h-5 w-5 text-blue-600" />
              </div>
              <UiBaseBadge v-if="project.is_archived" variant="secondary">
                Archived
              </UiBaseBadge>
            </div>

            <div class="flex-1">
              <h3
                class="mb-2 text-lg font-semibold text-gray-900 transition-colors group-hover:text-blue-700"
              >
                {{ project.name }}
              </h3>
              <p class="line-clamp-2 text-sm leading-relaxed text-gray-500">
                {{ project.description || "No description provided." }}
              </p>
            </div>

            <div
              class="mt-6 flex items-center justify-between border-t border-gray-100 pt-4 text-xs text-gray-500"
            >
              <div class="flex items-center gap-1">
                <span>
                  Created {{ format(new Date(project.created_at), "MMM d, yyyy") }}
                </span>
              </div>
              <div
                class="flex items-center gap-1 font-medium text-blue-600 opacity-0 transition-opacity group-hover:opacity-100"
              >
                Open <ArrowRight class="h-3 w-3" />
              </div>
            </div>
          </UiBaseCard>
        </NuxtLink>
      </div>
    </div>

    <ProjectDialog :is-open="isDialogOpen" @close="isDialogOpen = false" />
  </div>
</template>
