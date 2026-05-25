<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";

definePageMeta({ layout: "dashboard" });

const route = useRoute();
const projectId = computed(() => route.params.projectId as string);

const { project, isLoading } = useProject(projectId.value);

type TabType = "tasks" | "board" | "members" | "settings";
const activeTab = ref<TabType>("tasks");

const tabs: { id: TabType; label: string }[] = [
  { id: "tasks", label: "Tasks" },
  { id: "board", label: "Board" },
  { id: "members", label: "Members" },
  { id: "settings", label: "Settings" },
];
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center justify-center">
    <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
  </div>
  <div
    v-else-if="!project"
    class="text-muted-foreground flex h-full items-center justify-center"
  >
    Project not found.
  </div>
  <div v-else class="flex h-full flex-col space-y-6">
    <!-- Header -->
    <div>
      <div class="mb-1 flex items-center gap-2">
        <h1 class="text-foreground text-2xl font-bold tracking-tight">
          {{ project.name }}
        </h1>
        <span
          v-if="project.is_archived"
          class="rounded-full border border-yellow-200 bg-yellow-100 px-2 py-0.5 text-xs text-yellow-800"
        >
          Archived
        </span>
      </div>
      <p
        v-if="project.description"
        class="text-muted-foreground max-w-2xl text-sm"
      >
        {{ project.description }}
      </p>

      <!-- Tab Navigation -->
      <div class="border-border mt-6 flex gap-6 border-b">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="border-b-2 pb-3 text-sm font-medium transition-colors"
          :class="
            activeTab === tab.id
              ? 'border-primary text-primary'
              : 'text-muted-foreground hover:border-border hover:text-foreground/80 border-transparent'
          "
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="min-h-0 flex-1">
      <div v-if="activeTab === 'tasks'">
        <ProjectTabsTasks :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'board'">
        <ProjectTabsBoard :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'members'">
        <ProjectTabsMembers :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'settings'">
        <ProjectTabsSettings :project="project" />
      </div>
    </div>
  </div>
</template>
