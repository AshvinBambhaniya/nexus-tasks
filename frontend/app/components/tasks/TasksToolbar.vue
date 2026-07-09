<script setup lang="ts">
import { Search, ChevronDown } from "lucide-vue-next";

defineProps<{
  searchQuery: string;
}>();

defineEmits<{
  (e: "update:searchQuery", value: string): void;
}>();
</script>

<template>
  <div class="mb-2 flex items-center justify-between py-4">
    <div class="relative w-80">
      <div
        class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3"
      >
        <Search class="text-muted-foreground h-4 w-4" />
      </div>
      <input
        type="text"
        :value="searchQuery"
        placeholder="Search or jump to..."
        class="bg-background border-border placeholder:text-muted-foreground focus:ring-primary focus:border-primary text-foreground block w-full rounded-md border py-1.5 pr-3 pl-9 text-sm focus:ring-1 focus:outline-none"
        @input="
          $emit('update:searchQuery', ($event.target as HTMLInputElement).value)
        "
      />
    </div>

    <div class="flex items-center gap-4">
      <button
        class="text-muted-foreground hover:text-foreground hover:bg-muted/50 flex items-center gap-2 rounded-md px-3 py-1.5 text-sm transition-colors"
      >
        Filters
        <ChevronDown class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>
