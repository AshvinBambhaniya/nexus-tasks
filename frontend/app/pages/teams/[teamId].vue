<script setup lang="ts">
import { Loader2 } from "lucide-vue-next";

definePageMeta({ layout: "dashboard" });

const route = useRoute();
const teamId = computed(() => route.params.teamId as string);

const { team, isLoading } = useTeam(teamId.value);

type TabType = "projects" | "members" | "settings";
const activeTab = ref<TabType>("projects");

const tabs: { id: TabType; label: string }[] = [
  { id: "projects", label: "Overview" },
  { id: "members", label: "Members" },
  { id: "settings", label: "Settings" },
];
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center justify-center">
    <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
  </div>
  <div
    v-else-if="!team"
    class="flex h-full items-center justify-center text-gray-500"
  >
    Team not found.
  </div>
  <div v-else class="flex h-full flex-col space-y-6">
    <!-- Header -->
    <div>
      <div class="mb-1 flex items-center gap-2">
        <h1 class="text-2xl font-bold tracking-tight text-gray-900">
          {{ team.name }}
        </h1>
      </div>
      <p v-if="team.description" class="max-w-2xl text-sm text-gray-500">
        {{ team.description }}
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
      <div v-if="activeTab === 'projects'">
        <TeamTabsOverview :team-id="teamId" />
      </div>
      <div v-else-if="activeTab === 'members'">
        <TeamTabsMembers :team-id="teamId" />
      </div>
      <div v-else-if="activeTab === 'settings'">
        <TeamTabsSettings :team="team" />
      </div>
    </div>
  </div>
</template>
