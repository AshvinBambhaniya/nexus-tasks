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
    <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
  </div>
  <div
    v-else-if="!team"
    class="text-muted-foreground flex h-full items-center justify-center"
  >
    Team not found.
  </div>
  <div v-else class="flex h-full flex-col space-y-6">
    <!-- Header -->
    <div>
      <div class="mb-1 flex items-center gap-2">
        <h1 class="text-foreground text-2xl font-bold tracking-tight">
          {{ team.name }}
        </h1>
      </div>
      <p
        v-if="team.description"
        class="text-muted-foreground max-w-2xl text-sm"
      >
        {{ team.description }}
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
