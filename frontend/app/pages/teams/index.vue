<script setup lang="ts">
import { Users, Search, Plus, ArrowRight, MoreHorizontal } from "lucide-vue-next";
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
      <div v-else class="flex flex-col space-y-1 pb-10">
        <NuxtLink
          v-for="team in filteredTeams"
          :key="team.id"
          :to="`/teams/${team.id}`"
          class="group hover:bg-muted/40 border-transparent hover:border-border/60 flex items-center justify-between rounded-lg border px-4 py-3 transition-colors"
        >
          <div class="flex min-w-0 items-center gap-4">
            <div
              class="bg-muted/50 text-muted-foreground group-hover:text-primary group-hover:bg-primary/10 flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-transparent transition-all group-hover:border-primary/20 group-hover:shadow-sm"
            >
              <Users class="h-4 w-4" />
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h3 class="text-foreground truncate text-sm font-semibold tracking-tight">
                  {{ team.name }}
                </h3>
              </div>
              <p class="text-muted-foreground truncate text-xs mt-0.5">
                {{ team.description || "No description provided." }}
              </p>
            </div>
          </div>

          <div class="flex shrink-0 items-center gap-6 pl-4">
            <!-- Mock Avatar Stack for Team Members -->
            <div class="hidden items-center -space-x-2 sm:flex">
              <div class="border-background bg-indigo-500/20 text-indigo-500 flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">T</div>
              <div class="border-background bg-rose-500/20 text-rose-500 flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">E</div>
              <div class="border-background bg-amber-500/20 text-amber-500 flex h-6 w-6 items-center justify-center rounded-full border text-[9px] font-bold shadow-sm">M</div>
            </div>
            <div class="text-muted-foreground/70 hidden text-xs font-medium sm:block w-24 text-right">
              {{ team.created_at ? format(new Date(team.created_at), "MMM d, yy") : "--" }}
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

    <TeamDialog :is-open="isDialogOpen" @close="isDialogOpen = false" />
  </div>
</template>
