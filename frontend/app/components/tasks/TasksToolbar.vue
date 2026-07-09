<script setup lang="ts">
import { Search, ChevronDown } from "lucide-vue-next";

const props = defineProps<{
  searchQuery: string;
}>();

const emit = defineEmits<{
  (e: 'update:searchQuery', value: string): void;
}>();
</script>

<template>
  <div class="flex items-center justify-between py-4 mb-2">
    <div class="relative w-80">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <Search class="h-4 w-4 text-muted-foreground" />
      </div>
      <input
        type="text"
        :value="searchQuery"
        @input="$emit('update:searchQuery', ($event.target as HTMLInputElement).value)"
        placeholder="Search or jump to..."
        class="block w-full pl-9 pr-3 py-1.5 bg-background border border-border rounded-md text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary text-foreground"
      />
    </div>

    <div class="flex items-center gap-4">
      <button class="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors px-3 py-1.5 rounded-md hover:bg-muted/50">
        Filters
        <ChevronDown class="h-4 w-4" />
      </button>
    </div>
  </div>
</template>
