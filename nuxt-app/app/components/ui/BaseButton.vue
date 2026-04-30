<script setup lang="ts">
interface Props {
  variant?: 'primary' | 'ghost' | 'outline' | 'white'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false
})

const variants = {
  primary: 'bg-blue-600 text-white shadow-md shadow-blue-100 hover:bg-blue-700',
  ghost: 'font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50',
  outline: 'border border-slate-200 bg-white text-slate-700 hover:text-slate-900 hover:bg-gray-50',
  white: 'bg-white text-gray-900 hover:bg-gray-100 shadow-md'
}

const sizes = {
  sm: 'h-8 px-3 text-xs',
  md: 'h-10 px-4 py-2',
  lg: 'h-12 px-8 text-base'
}
</script>

<template>
  <button
    :disabled="disabled"
    class="inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-gray-950 disabled:pointer-events-none disabled:opacity-50"
    :class="[variants[variant], sizes[size]]"
  >
    <slot />
  </button>
</template>
