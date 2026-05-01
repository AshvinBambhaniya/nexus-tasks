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
    <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
  </div>
  <div
    v-else-if="projects.length === 0"
    class="flex flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-200 bg-gray-50 p-12 text-center text-gray-500"
  >
    <Folder class="mb-3 h-10 w-10 text-gray-300" />
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
        class="cursor-pointer border-gray-200 p-5 transition-shadow group-hover:border-indigo-200 hover:shadow-md"
      >
        <div class="mb-3 flex items-center justify-between">
          <div
            class="rounded-lg bg-indigo-50 p-2 text-indigo-600 transition-colors group-hover:bg-indigo-100"
          >
            <Folder class="h-5 w-5" />
          </div>
          <ArrowRight
            class="h-4 w-4 transform text-gray-300 transition-colors group-hover:translate-x-1 group-hover:text-indigo-500"
          />
        </div>
        <h3
          class="truncate font-semibold text-gray-900 transition-colors group-hover:text-indigo-700"
        >
          {{ project.name }}
        </h3>
        <p
          v-if="project.description"
          class="mt-1 line-clamp-2 text-sm leading-relaxed text-gray-500"
        >
          {{ project.description }}
        </p>
        <div
          class="mt-4 border-t border-gray-50 pt-4 text-[10px] tracking-wider text-gray-400 uppercase"
        >
          Created {{ format(new Date(project.created_at), "MMM d, yyyy") }}
        </div>
      </UiBaseCard>
    </NuxtLink>
  </div>
</template>
