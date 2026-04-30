<script setup lang="ts">
interface Props {
  className?: string;
}

const { className = "" } = defineProps<Props>();
</script>

<template>
  <div :class="cn('rounded-lg border bg-white shadow-sm', className)">
    <slot />
  </div>
</template>
