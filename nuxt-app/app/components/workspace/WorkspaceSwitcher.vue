<script setup lang="ts">
import { ChevronsUpDown, Check, Building2, User } from "lucide-vue-next";
import { WorkspaceType } from "~/types";

const { workspaces, activeWorkspace, isLoading } = useWorkspaces();
const workspaceStore = useWorkspaceStore();

const isOpen = ref(false);
const isDialogOpen = ref(false);
const dropdownRef = ref<HTMLElement | null>(null);

useClickOutside(dropdownRef, () => {
  isOpen.value = false;
});

const handleSelect = (id: number) => {
  workspaceStore.setActiveWorkspaceId(id);
  isOpen.value = false;
};
</script>

<template>
  <div ref="dropdownRef" class="relative w-full">
    <template v-if="isLoading">
      <div class="h-10 w-full animate-pulse rounded bg-gray-200" />
    </template>
    <template v-else>
      <button
        class="flex w-full items-center justify-between rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm font-medium transition-colors hover:bg-gray-50"
        :class="{ 'ring-2 ring-blue-500 ring-offset-1': isOpen }"
        @click="isOpen = !isOpen"
      >
        <span class="flex items-center gap-2 truncate">
          <User
            v-if="activeWorkspace?.type === WorkspaceType.PERSONAL"
            class="h-4 w-4 text-gray-500"
          />
          <Building2 v-else class="h-4 w-4 text-gray-500" />
          <span class="truncate text-gray-900">
            {{ activeWorkspace?.name || "Select Workspace" }}
          </span>
        </span>
        <ChevronsUpDown class="h-4 w-4 shrink-0 text-gray-400" />
      </button>

      <div
        v-if="isOpen"
        class="absolute top-full left-0 z-50 mt-1 w-full overflow-hidden rounded-md border border-gray-200 bg-white shadow-lg"
      >
        <div class="max-h-[300px] overflow-auto p-1">
          <div
            v-if="workspaces.length === 0"
            class="px-2 py-2 text-center text-sm text-gray-500"
          >
            No workspaces found.
          </div>
          <template v-else>
            <button
              v-for="workspace in workspaces"
              :key="workspace.id"
              class="relative flex w-full cursor-pointer items-center rounded-sm px-2 py-2 text-sm outline-none select-none hover:bg-blue-50 hover:text-blue-600"
              :class="
                activeWorkspace?.id === workspace.id
                  ? 'bg-blue-50 text-blue-600'
                  : 'text-gray-700'
              "
              @click="handleSelect(workspace.id)"
            >
              <div class="mr-2 flex h-4 w-4 items-center justify-center">
                <User
                  v-if="workspace.type === WorkspaceType.PERSONAL"
                  class="h-3 w-3"
                />
                <Building2 v-else class="h-3 w-3" />
              </div>
              <span class="flex-1 truncate text-left font-medium">
                {{ workspace.name }}
              </span>
              <Check
                v-if="activeWorkspace?.id === workspace.id"
                class="ml-auto h-4 w-4 text-blue-600"
              />
            </button>
          </template>
        </div>
        <div class="border-t border-gray-100 p-1">
          <button
            class="flex w-full cursor-pointer items-center rounded-sm px-2 py-2 text-sm text-gray-500 hover:bg-gray-100"
            @click="
              isOpen = false;
              isDialogOpen = true;
            "
          >
            <span
              class="mr-2 flex h-4 w-4 items-center justify-center text-lg leading-none"
            >
              +
            </span>
            Create Workspace
          </button>
        </div>
      </div>
    </template>

    <WorkspaceWorkspaceDialog
      :is-open="isDialogOpen"
      @close="isDialogOpen = false"
    />
  </div>
</template>
