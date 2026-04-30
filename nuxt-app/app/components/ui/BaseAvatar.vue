<script setup lang="ts">
interface Props {
  src?: string;
  alt?: string;
  fallback: string;
  className?: string;
}

const {
  src = "",
  alt = "",
  fallback,
  className = "",
} = defineProps<Props>();
</script>

<template>
  <div
    :class="
      cn(
        'relative flex h-8 w-8 shrink-0 overflow-hidden rounded-full bg-gray-100',
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
      class="flex h-full w-full items-center justify-center bg-blue-100 text-[10px] font-medium text-blue-600"
    >
      {{ fallback }}
    </div>
  </div>
</template>
