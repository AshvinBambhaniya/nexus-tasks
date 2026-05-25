<script setup lang="ts">
import { cn } from "~/utils/cn";

// Define specific types for better autocomplete and strict checking
type ButtonVariant = "primary" | "ghost" | "outline" | "white" | "destructive";
type ButtonSize = "sm" | "md" | "lg";

interface Props {
  variant?: ButtonVariant;
  size?: ButtonSize;
  disabled?: boolean;
  class?: string;
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
  class: className = "",
} = defineProps<Props>();

const variants: Record<ButtonVariant, string> = {
  primary: "bg-primary text-primary-foreground shadow-md hover:opacity-90",
  ghost:
    "font-medium text-muted-foreground hover:text-foreground hover:bg-muted",
  outline: "border border-input bg-transparent text-foreground hover:bg-muted",
  white:
    "bg-white dark:bg-slate-800 text-gray-900 dark:text-gray-100 hover:bg-gray-100 dark:hover:bg-slate-700 shadow-md",
  destructive:
    "bg-destructive text-destructive-foreground shadow-md hover:opacity-90",
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
    :class="
      cn(
        'inline-flex items-center justify-center rounded-md font-semibold transition-colors focus-visible:ring-1 focus-visible:ring-gray-950 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50',
        variants[variant],
        sizes[size],
        className
      )
    "
  >
    <slot />
  </button>
</template>
