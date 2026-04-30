<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";

definePageMeta({ layout: "dashboard" });

const route = useRoute();
const projectId = computed(() => parseInt(route.params.projectId as string));

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
    <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
  </div>
  <div
    v-else-if="!project"
    class="flex h-full items-center justify-center text-gray-500"
  >
    Project not found.
  </div>
  <div v-else class="flex h-full flex-col space-y-6">
    <!-- Header -->
    <div>
      <div class="mb-1 flex items-center gap-2">
        <h1 class="text-2xl font-bold tracking-tight text-gray-900">
          {{ project.name }}
        </h1>
        <span
          v-if="project.is_archived"
          class="rounded-full border border-yellow-200 bg-yellow-100 px-2 py-0.5 text-xs text-yellow-800"
        >
          Archived
        </span>
      </div>
      <p v-if="project.description" class="max-w-2xl text-sm text-gray-500">
        {{ project.description }}
      </p>

      <!-- Tab Navigation -->
      <div class="mt-6 flex gap-6 border-b border-gray-200">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="border-b-2 pb-3 text-sm font-medium transition-colors"
          :class="
            activeTab === tab.id
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
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
