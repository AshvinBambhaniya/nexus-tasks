<script setup lang="ts">
import { cn } from "~/utils/cn";

interface Props {
  for?: string;
  className?: string;
}

const { for: htmlFor = "", className = "" } = defineProps<Props>();
</script>

<template>
  <label
    :for="htmlFor"
    :class="
      cn(
        'text-foreground text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
        className
      )
    "
  >
    <slot />
  </label>
</template>
