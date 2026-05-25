<script setup lang="ts">
import { Users, Search, Plus, ArrowRight } from "lucide-vue-next";
import { format } from "date-fns";

definePageMeta({ layout: "dashboard" });

const { teams, isLoading } = useTeams();
const searchQuery = ref("");
const isDialogOpen = ref(false);

const filteredTeams = computed(() => {
  return teams.value.filter((team) => {
    return team.name.toLowerCase().includes(searchQuery.value.toLowerCase());
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
        <h1 class="text-foreground text-2xl font-bold">Teams</h1>
        <p class="text-muted-foreground mt-1">
          Organize people into groups for easier management.
        </p>
      </div>
      <UiBaseButton @click="isDialogOpen = true">
        <Plus class="mr-2 h-4 w-4" /> New Team
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
          placeholder="Search teams..."
          class-name="pl-9"
        />
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
          class="border-border bg-muted h-40 animate-pulse rounded-xl border"
        />
      </div>
      <div
        v-else-if="filteredTeams.length === 0"
        class="border-border bg-muted/50 flex h-64 flex-col items-center justify-center rounded-lg border-2 border-dashed text-center"
      >
        <div
          class="bg-muted mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full"
        >
          <Users class="text-muted-foreground/50 h-6 w-6" />
        </div>
        <h3 class="text-foreground text-lg font-medium">No teams found</h3>
        <p class="text-muted-foreground mt-1 max-w-sm text-sm">
          {{
            searchQuery
              ? `No teams matching "${searchQuery}"`
              : "Create a team to group users together."
          }}
        </p>
        <UiBaseButton
          v-if="!searchQuery"
          variant="outline"
          class="mt-4"
          @click="isDialogOpen = true"
        >
          Create Team
        </UiBaseButton>
      </div>
      <div v-else class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
        <NuxtLink
          v-for="team in filteredTeams"
          :key="team.id"
          :to="`/teams/${team.id}`"
          class="block h-full"
        >
          <UiBaseCard
            class="group border-border hover:border-primary/50 flex h-full cursor-pointer flex-col p-6 transition-all hover:shadow-md"
          >
            <div class="mb-4 flex items-start justify-between">
              <div
                class="bg-primary/10 group-hover:bg-primary/20 rounded-xl p-2.5 transition-colors"
              >
                <Users class="text-primary h-5 w-5" />
              </div>
            </div>

            <div class="flex-1">
              <h3
                class="text-foreground group-hover:text-primary mb-2 text-lg font-semibold transition-colors"
              >
                {{ team.name }}
              </h3>
              <p
                class="text-muted-foreground line-clamp-2 text-sm leading-relaxed"
              >
                {{ team.description || "No description provided." }}
              </p>
            </div>

            <div
              class="border-border text-muted-foreground mt-6 flex items-center justify-between border-t pt-4 text-xs"
            >
              <div class="flex items-center gap-1">
                <span>
                  Created
                  {{
                    team.created_at
                      ? format(new Date(team.created_at), "MMM d, yyyy")
                      : "Unknown date"
                  }}
                </span>
              </div>
              <div
                class="text-primary flex items-center gap-1 font-medium opacity-0 transition-opacity group-hover:opacity-100"
              >
                Manage <ArrowRight class="h-3 w-3" />
              </div>
            </div>
          </UiBaseCard>
        </NuxtLink>
      </div>
    </div>

    <TeamDialog :is-open="isDialogOpen" @close="isDialogOpen = false" />
  </div>
</template>
