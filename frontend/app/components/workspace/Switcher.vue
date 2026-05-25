<script setup lang="ts">
import { ChevronsUpDown, Check, Building2, User } from "lucide-vue-next";
import { WorkspaceType } from "~/types";

interface Props {
  isCollapsed?: boolean;
}

const { isCollapsed = false } = defineProps<Props>();

const { workspaces, activeWorkspace, isLoading } = useWorkspaces();
const workspaceStore = useWorkspaceStore();

const isOpen = ref(false);
const isDialogOpen = ref(false);
const dropdownRef = ref<HTMLElement | null>(null);

useClickOutside(dropdownRef, () => {
  isOpen.value = false;
});

const handleSelect = (id: string) => {
  workspaceStore.setActiveWorkspaceId(id);
  isOpen.value = false;
};
</script>

<template>
  <div ref="dropdownRef" class="relative w-full">
    <template v-if="isLoading">
      <div class="bg-muted h-10 w-full animate-pulse rounded" />
    </template>
    <template v-else>
      <button
        class="border-border bg-card hover:bg-muted/50 flex w-full items-center justify-between rounded-lg border px-3 py-2 text-sm font-medium transition-colors"
        :class="{
          'ring-ring ring-offset-background ring-2 ring-offset-1': isOpen,
          'justify-center px-2': isCollapsed,
        }"
        :title="isCollapsed ? activeWorkspace?.name : ''"
        @click="isOpen = !isOpen"
      >
        <span
          class="flex items-center gap-2 truncate"
          :class="{ 'justify-center': isCollapsed }"
        >
          <User
            v-if="activeWorkspace?.type === WorkspaceType.PERSONAL"
            class="text-muted-foreground h-4 w-4 shrink-0"
          />
          <Building2 v-else class="text-muted-foreground h-4 w-4 shrink-0" />
          <span v-if="!isCollapsed" class="text-card-foreground truncate">
            {{ activeWorkspace?.name || "Select Workspace" }}
          </span>
        </span>
        <ChevronsUpDown
          v-if="!isCollapsed"
          class="text-muted-foreground/70 h-4 w-4 shrink-0"
        />
      </button>

      <div
        v-if="isOpen"
        class="border-border bg-card absolute top-full left-0 z-50 mt-1 overflow-hidden rounded-md border shadow-lg"
        :class="isCollapsed ? 'w-48' : 'w-full'"
      >
        <div class="max-h-[300px] overflow-auto p-1">
          <div
            v-if="workspaces.length === 0"
            class="text-muted-foreground px-2 py-2 text-center text-sm"
          >
            No workspaces found.
          </div>
          <template v-else>
            <button
              v-for="workspace in workspaces"
              :key="workspace.id"
              class="hover:bg-primary/10 hover:text-primary relative flex w-full cursor-pointer items-center rounded-sm px-2 py-2 text-sm outline-none select-none"
              :class="
                activeWorkspace?.id === workspace.id
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted-foreground'
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
                class="text-primary ml-auto h-4 w-4"
              />
            </button>
          </template>
        </div>
        <div class="border-border/50 border-t p-1">
          <button
            class="text-muted-foreground hover:bg-muted flex w-full cursor-pointer items-center rounded-sm px-2 py-2 text-sm"
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

    <WorkspaceDialog :is-open="isDialogOpen" @close="isDialogOpen = false" />
  </div>
</template>
