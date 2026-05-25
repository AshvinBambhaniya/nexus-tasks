<script setup lang="ts">
import { cn } from "~/utils/cn";

interface Props {
  src?: string;
  alt?: string;
  fallback: string;
  className?: string;
}

const { src = "", alt = "", fallback, className = "" } = defineProps<Props>();
</script>

<template>
  <div
    :class="
      cn(
        'bg-muted relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full',
        className
      )
    "
  >
    <img
      v-if="src"
      class="aspect-square h-full w-full object-cover"
      :src="src"
      :alt="alt"
    />
    <div
      v-else
      class="bg-primary/10 text-primary flex h-full w-full items-center justify-center text-[10px] font-medium"
    >
      {{ fallback }}
    </div>
  </div>
</template>
