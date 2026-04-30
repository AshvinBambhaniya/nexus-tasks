<script setup lang="ts">
// Define specific types for better autocomplete and strict checking
type ButtonVariant = "primary" | "ghost" | "outline" | "white" | "destructive";
type ButtonSize = "sm" | "md" | "lg";

interface Props {
  variant?: ButtonVariant;
  size?: ButtonSize;
  disabled?: boolean;
}

/**
 * BEST PRACTICE: Nuxt 4 (Vue 3.4+) Reactive Props Destructure
 * This allows setting defaults directly in the destructure and
 * automatically maintains reactivity.
 */
const {
  variant = "primary",
  size = "md",
  disabled = false,
} = defineProps<Props>();

const variants: Record<ButtonVariant, string> = {
  primary: "bg-blue-600 text-white shadow-md shadow-blue-100 hover:bg-blue-700",
  ghost: "font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50",
  outline:
    "border border-slate-200 bg-white text-slate-700 hover:text-slate-900 hover:bg-gray-50",
  white: "bg-white text-gray-900 hover:bg-gray-100 shadow-md",
  destructive:
    "bg-red-600 text-white shadow-md shadow-red-100 hover:bg-red-700",
};

const sizes: Record<ButtonSize, string> = {
  sm: "h-8 px-3 text-xs",
  md: "h-10 px-4 py-2",
  lg: "h-12 px-8 text-base",
};
</script>

<template>
  <button
    :disabled="disabled"
    class="inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:ring-1 focus-visible:ring-gray-950 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50"
    :class="[variants[variant], sizes[size]]"
  >
    <slot />
  </button>
</template>
