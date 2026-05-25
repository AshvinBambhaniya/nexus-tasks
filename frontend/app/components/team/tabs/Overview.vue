<script setup lang="ts">
import { Folder, ArrowRight, Loader2 } from "lucide-vue-next";
import { format } from "date-fns";

interface Props {
  teamId: string;
}

const { teamId } = defineProps<Props>();

const { team, isLoading } = useTeam(teamId);
const projects = computed(() => team.value?.projects || []);
</script>

<template>
  <div v-if="isLoading" class="flex justify-center p-12">
    <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
  </div>
  <div
    v-else-if="projects.length === 0"
    class="border-border bg-muted text-muted-foreground flex flex-col items-center justify-center rounded-lg border-2 border-dashed p-12 text-center"
  >
    <Folder class="text-muted-foreground/30 mb-3 h-10 w-10" />
    <p>No projects assigned to this team yet.</p>
  </div>
  <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
    <NuxtLink
      v-for="project in projects"
      :key="project.id"
      :to="`/projects/${project.id}`"
      class="group"
    >
      <UiBaseCard
        class="border-border group-hover:border-primary/30 cursor-pointer p-5 transition-shadow hover:shadow-md"
      >
        <div class="mb-3 flex items-center justify-between">
          <div
            class="bg-primary/10 text-primary group-hover:bg-primary/20 rounded-lg p-2 transition-colors"
          >
            <Folder class="h-5 w-5" />
          </div>
          <ArrowRight
            class="text-muted-foreground/30 group-hover:text-primary h-4 w-4 transform transition-colors group-hover:translate-x-1"
          />
        </div>
        <h3
          class="text-foreground group-hover:text-primary truncate font-semibold transition-colors"
        >
          {{ project.name }}
        </h3>
        <p
          v-if="project.description"
          class="text-muted-foreground mt-1 line-clamp-2 text-sm leading-relaxed"
        >
          {{ project.description }}
        </p>
        <div
          class="border-border text-muted-foreground/70 mt-4 border-t pt-4 text-[10px] tracking-wider uppercase"
        >
          Created {{ format(new Date(project.created_at), "MMM d, yyyy") }}
        </div>
      </UiBaseCard>
    </NuxtLink>
  </div>
</template>
