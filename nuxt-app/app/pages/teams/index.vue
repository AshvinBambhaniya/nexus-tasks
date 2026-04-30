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
        <h1 class="text-2xl font-bold text-gray-900">Teams</h1>
        <p class="mt-1 text-gray-500">
          Organize people into groups for easier management.
        </p>
      </div>
      <UiBaseButton @click="isDialogOpen = true">
        <Plus class="mr-2 h-4 w-4" /> New Team
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
          class="h-40 animate-pulse rounded-xl border border-gray-200 bg-gray-100"
        />
      </div>
      <div
        v-else-if="filteredTeams.length === 0"
        class="flex h-64 flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 text-center"
      >
        <div
          class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100"
        >
          <Users class="h-6 w-6 text-gray-400" />
        </div>
        <h3 class="text-lg font-medium text-gray-900">No teams found</h3>
        <p class="mt-1 max-w-sm text-sm text-gray-500">
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
            class="group flex h-full cursor-pointer flex-col border-gray-200 p-6 transition-all hover:border-blue-200 hover:shadow-md"
          >
            <div class="mb-4 flex items-start justify-between">
              <div
                class="rounded-xl bg-indigo-50 p-2.5 transition-colors group-hover:bg-indigo-100"
              >
                <Users class="h-5 w-5 text-indigo-600" />
              </div>
            </div>

            <div class="flex-1">
              <h3
                class="mb-2 text-lg font-semibold text-gray-900 transition-colors group-hover:text-indigo-700"
              >
                {{ team.name }}
              </h3>
              <p class="line-clamp-2 text-sm leading-relaxed text-gray-500">
                {{ team.description || "No description provided." }}
              </p>
            </div>

            <div
              class="mt-6 flex items-center justify-between border-t border-gray-100 pt-4 text-xs text-gray-500"
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
                class="flex items-center gap-1 font-medium text-indigo-600 opacity-0 transition-opacity group-hover:opacity-100"
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
