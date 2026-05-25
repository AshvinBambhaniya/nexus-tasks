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
        <h1 class="text-foreground text-2xl font-bold">Projects</h1>
        <p class="text-muted-foreground mt-1">
          Manage and organize your team's work.
        </p>
      </div>
      <UiBaseButton @click="isDialogOpen = true">
        <Plus class="mr-2 h-4 w-4" /> New Project
      </UiBaseButton>
    </div>

    <div
      class="border-border flex flex-col items-center justify-between gap-4 border-b pb-4 sm:flex-row"
    >
      <!-- Search -->
      <div class="relative w-full sm:w-72">
        <Search
          class="text-muted-foreground/70 absolute top-2.5 left-3 h-4 w-4"
        />
        <UiBaseInput
          v-model="searchQuery"
          placeholder="Search projects..."
          class-name="pl-9"
        />
      </div>

      <!-- Filter Tabs -->
      <div class="bg-muted flex self-start rounded-lg p-1 sm:self-auto">
        <button
          v-for="status in ['active', 'archived', 'all'] as const"
          :key="status"
          class="rounded-md px-4 py-1.5 text-sm font-medium capitalize transition-all"
          :class="
            filterStatus === status
              ? 'bg-background text-foreground shadow-sm'
              : 'text-muted-foreground hover:text-foreground/80'
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
          class="border-border bg-muted h-48 animate-pulse rounded-xl border"
        />
      </div>
      <div
        v-else-if="filteredProjects.length === 0"
        class="border-border bg-muted/50 flex h-64 flex-col items-center justify-center rounded-lg border-2 border-dashed text-center"
      >
        <div
          class="bg-muted mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
        >
          <Search v-if="searchQuery" class="text-muted-foreground/70 h-6 w-6" />
          <Folder v-else class="text-muted-foreground/70 h-6 w-6" />
        </div>
        <h3 class="text-foreground text-lg font-medium">No projects found</h3>
        <p class="text-muted-foreground mt-1 max-w-sm text-sm">
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
            class="group border-border hover:border-primary/30 flex h-full cursor-pointer flex-col p-6 transition-all hover:shadow-md"
          >
            <div class="mb-4 flex items-start justify-between">
              <div
                class="bg-primary/10 group-hover:bg-primary/20 rounded-xl p-2.5 transition-colors"
              >
                <Folder class="text-primary h-5 w-5" />
              </div>
              <UiBaseBadge v-if="project.is_archived" variant="secondary">
                Archived
              </UiBaseBadge>
            </div>

            <div class="flex-1">
              <h3
                class="text-foreground group-hover:text-primary mb-2 text-lg font-semibold transition-colors"
              >
                {{ project.name }}
              </h3>
              <p
                class="text-muted-foreground line-clamp-2 text-sm leading-relaxed"
              >
                {{ project.description || "No description provided." }}
              </p>
            </div>

            <div
              class="border-border text-muted-foreground mt-6 flex items-center justify-between border-t pt-4 text-xs"
            >
              <div class="flex items-center gap-1">
                <span>
                  Created
                  {{ format(new Date(project.created_at), "MMM d, yyyy") }}
                </span>
              </div>
              <div
                class="text-primary flex items-center gap-1 font-medium opacity-0 transition-opacity group-hover:opacity-100"
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
