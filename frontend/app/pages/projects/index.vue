<script setup lang="ts">
import { Folder, Search, Plus, ArrowRight, MoreHorizontal } from "lucide-vue-next";
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
      <div v-else class="flex flex-col space-y-1 pb-10">
        <NuxtLink
          v-for="project in filteredProjects"
          :key="project.id"
          :to="`/projects/${project.id}`"
          class="group hover:bg-muted/40 border-transparent hover:border-border/60 flex items-center justify-between rounded-lg border px-4 py-3 transition-colors"
        >
          <div class="flex min-w-0 items-center gap-4">
            <div
              class="bg-muted/50 text-muted-foreground group-hover:text-primary group-hover:bg-primary/10 flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-transparent transition-all group-hover:border-primary/20 group-hover:shadow-sm"
            >
              <Folder class="h-4 w-4" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-foreground truncate text-sm font-semibold tracking-tight">
                  {{ project.name }}
                </h3>
                <span
                  class="bg-emerald-500/10 text-emerald-500 border-emerald-500/20 inline-flex items-center justify-center rounded-sm border px-1.5 py-0.5 text-[9px] font-bold tracking-wider uppercase shadow-sm"
                  :class="project.is_archived ? 'bg-muted text-muted-foreground border-border' : ''"
                >
                  {{ project.is_archived ? "Archived" : "Active" }}
                </span>
              </div>
              <p class="text-muted-foreground truncate text-xs mt-0.5">
                {{ project.description || "No description provided." }}
              </p>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-6 pl-4">
            <!-- Mock Avatar Stack -->
            <div class="hidden items-center -space-x-2 sm:flex">
              <div class="border-background bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300 flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">A</div>
              <div class="border-background bg-slate-300 dark:bg-slate-600 text-slate-700 dark:text-slate-300 flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">M</div>
              <div class="border-background bg-slate-400 dark:bg-slate-500 text-white flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">J</div>
            </div>
            <div class="text-muted-foreground/70 hidden text-xs font-medium sm:block w-24 text-right">
              {{ format(new Date(project.created_at), "MMM d, yy") }}
            </div>
            <div class="flex h-8 w-8 items-center justify-center rounded-md hover:bg-muted/80 transition-colors">
              <MoreHorizontal
                class="text-muted-foreground/50 hover:text-foreground h-4 w-4 opacity-0 transition-opacity group-hover:opacity-100"
              />
            </div>
          </div>
        </NuxtLink>
      </div>
    </div>

    <ProjectDialog :is-open="isDialogOpen" @close="isDialogOpen = false" />
  </div>
</template>
