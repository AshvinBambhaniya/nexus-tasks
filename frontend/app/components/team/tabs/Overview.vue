<script setup lang="ts">
import { Folder, Loader2, MoreHorizontal } from "lucide-vue-next";

interface Props {
  teamId: string;
  hideMembers?: boolean;
}

const { teamId, hideMembers = false } = defineProps<Props>();
const { team, isLoading } = useTeam(teamId);
const projects = computed(() => team.value?.projects || []);
const { members } = useTeamMembers(teamId);
</script>

<template>
  <div class="flex flex-col gap-12 pb-12 lg:flex-row">
    <!-- Active Projects Section (60%) -->
    <div class="flex-1">
      <div class="mb-6 flex items-center justify-between">
        <h2 class="text-foreground text-xl font-semibold tracking-tight">
          Projects
        </h2>
      </div>

      <div v-if="isLoading" class="flex justify-center p-12">
        <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
      </div>
      <div
        v-else-if="projects.length === 0"
        class="border-border bg-card text-muted-foreground flex flex-col items-center justify-center rounded-xl border border-dashed p-16 text-center shadow-sm"
      >
        <div class="bg-muted mb-4 rounded-full p-4">
          <Folder class="text-muted-foreground/50 h-8 w-8" />
        </div>
        <h3 class="text-foreground text-lg font-medium">No active projects</h3>
        <p class="text-muted-foreground mt-1 max-w-sm text-sm">
          This team hasn't been assigned any projects yet.
        </p>
      </div>
      <div v-else class="grid gap-4 sm:grid-cols-2">
        <NuxtLink
          v-for="project in projects"
          :key="project.id"
          :to="`/projects/${project.id}`"
          class="group block h-full"
        >
          <div
            class="border-border bg-card hover:border-primary/50 relative flex h-full flex-col overflow-hidden rounded-xl border shadow-sm transition-all duration-200"
          >
            <div class="p-5">
              <h3
                class="text-foreground mb-1.5 text-base font-semibold tracking-tight"
              >
                {{ project.name }}
              </h3>
              <p
                v-if="project.description"
                class="text-muted-foreground line-clamp-2 text-sm leading-relaxed"
              >
                {{ project.description }}
              </p>
              <p v-else class="text-muted-foreground/50 text-sm italic">
                No description
              </p>
            </div>

            <div class="mt-auto px-5 pt-4 pb-5">
              <div class="flex items-center justify-between">
                <div
                  class="text-muted-foreground flex items-center gap-1.5 text-xs font-medium"
                >
                  <Folder class="h-3.5 w-3.5" />
                  GitHub
                </div>
                <UiBaseBadge
                  variant="secondary"
                  class="rounded-full px-2 py-0.5 text-[10px] font-medium"
                >
                  Staging
                </UiBaseBadge>
              </div>
            </div>
          </div>
        </NuxtLink>
      </div>
    </div>

    <!-- Team Members Section (40%) -->
    <div v-if="!hideMembers" class="w-full shrink-0 lg:w-[350px]">
      <div class="mb-6 flex items-center justify-between">
        <h2 class="text-foreground text-xl font-semibold tracking-tight">
          Team Members
        </h2>
      </div>

      <div class="space-y-4">
        <div
          v-for="member in members.slice(0, 5)"
          :key="member.user_id"
          class="group flex items-center justify-between"
        >
          <div class="flex items-center gap-3">
            <div
              class="bg-primary/10 text-primary flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-xs font-bold uppercase"
            >
              {{ (member.email || "?")[0] }}
            </div>
            <div>
              <div
                class="text-foreground mb-1 text-sm leading-none font-medium"
              >
                {{ member.email.split("@")[0] }}
              </div>
              <div class="text-muted-foreground text-xs capitalize">
                {{ (member.role || "Member").toLowerCase() }}
              </div>
            </div>
          </div>

          <div class="flex items-center gap-3">
            <div
              class="text-muted-foreground flex items-center gap-1.5 text-xs"
            >
              <div
                class="h-1.5 w-1.5 rounded-full"
                :class="
                  member.role === 'ADMIN'
                    ? 'bg-emerald-500'
                    : 'bg-muted-foreground'
                "
              />
              {{ member.role === "ADMIN" ? "Online" : "Offline" }}
            </div>
            <button
              class="text-muted-foreground hover:text-foreground opacity-0 transition-opacity group-hover:opacity-100"
            >
              <MoreHorizontal class="h-4 w-4" />
            </button>
          </div>
        </div>

        <button
          v-if="members.length > 5"
          class="text-muted-foreground hover:text-foreground mt-2 inline-block text-sm font-medium transition-colors"
        >
          View all {{ members.length }} members
        </button>
      </div>
    </div>
  </div>
</template>
