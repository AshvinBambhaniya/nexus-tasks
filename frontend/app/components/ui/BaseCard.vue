<script setup lang="ts">
import { cn } from "~/utils/cn";

interface Props {
  className?: string;
}

const { className = "" } = defineProps<Props>();
</script>

<template>
  <div
    :class="
      cn(
        'border-border bg-card text-card-foreground rounded-lg border shadow-sm',
        className
      )
    "
  >
    <slot />
  </div>
</template>
