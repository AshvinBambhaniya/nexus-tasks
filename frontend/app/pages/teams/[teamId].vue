<script setup lang="ts">
import { Loader2, Users } from "lucide-vue-next";

definePageMeta({ layout: "dashboard" });

const route = useRoute();
const teamId = computed(() => route.params.teamId as string);

const { team, isLoading } = useTeam(teamId.value);

type TabType = "overview" | "projects" | "members" | "settings";
const activeTab = ref<TabType>("overview");

const tabs = [
  { id: "overview", label: "Overview" },
  { id: "projects", label: "Projects" },
  { id: "members", label: "Members" },
  { id: "settings", label: "Settings" },
];
</script>

<template>
  <div
    v-if="isLoading"
    class="flex h-[calc(100vh-8rem)] items-center justify-center"
  >
    <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
  </div>
  <div
    v-else-if="!team"
    class="text-muted-foreground flex h-[calc(100vh-8rem)] items-center justify-center"
  >
    <div class="flex flex-col items-center gap-4 text-center">
      <div class="bg-muted rounded-full p-4">
        <Users class="text-muted-foreground h-8 w-8" />
      </div>
      <div>
        <h3 class="text-foreground text-lg font-medium">Team not found</h3>
        <p class="text-sm">
          The team you are looking for does not exist or you lack permission.
        </p>
      </div>
    </div>
  </div>
  <div
    v-else
    class="bg-background -mx-8 -mt-6 flex h-[calc(100vh-4rem)] flex-col"
  >
    <!-- Header Area -->
    <div class="border-border bg-background border-b px-8 pt-10 pb-0">
      <div class="mx-auto w-full max-w-[1400px]">
        <h1 class="text-foreground mb-6 text-3xl font-bold tracking-tight">
          {{ team.name }}
        </h1>
        <!-- Navigation Tabs -->
        <nav class="flex gap-6">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="hover:text-foreground relative pb-4 text-sm font-medium transition-colors"
            :class="
              activeTab === tab.id ? 'text-foreground' : 'text-muted-foreground'
            "
            @click="activeTab = tab.id as TabType"
          >
            {{ tab.label }}
            <div
              v-if="activeTab === tab.id"
              class="bg-primary absolute right-0 bottom-[-1px] left-0 h-0.5 rounded-t-full"
            />
          </button>
        </nav>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="bg-background/50 flex-1 overflow-y-auto">
      <div class="mx-auto w-full max-w-[1400px] p-8">
        <div
          v-if="activeTab === 'overview'"
          class="animate-in fade-in duration-300"
        >
          <TeamTabsOverview :team-id="teamId" />
        </div>
        <div
          v-else-if="activeTab === 'projects'"
          class="animate-in fade-in duration-300"
        >
          <!-- Rendering overview for projects for simplicity, we could filter -->
          <TeamTabsOverview :team-id="teamId" :hide-members="true" />
        </div>
        <div
          v-else-if="activeTab === 'members'"
          class="animate-in fade-in max-w-4xl duration-300"
        >
          <TeamTabsMembers :team-id="teamId" />
        </div>
        <div
          v-else-if="activeTab === 'settings'"
          class="animate-in fade-in duration-300"
        >
          <TeamTabsSettings :team="team" />
        </div>
      </div>
    </div>
  </div>
</template>
