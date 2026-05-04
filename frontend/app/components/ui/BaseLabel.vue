<script setup lang="ts">
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
        'text-sm leading-none font-medium text-gray-900 peer-disabled:cursor-not-allowed peer-disabled:opacity-70',
        className
      )
    "
  >
    <slot />
  </label>
</template>
