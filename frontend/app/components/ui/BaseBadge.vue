<script setup lang="ts">
type BadgeVariant = "default" | "secondary" | "outline" | "destructive";

interface Props {
  variant?: BadgeVariant;
  className?: string;
}

const { variant = "default", className = "" } = defineProps<Props>();

const variants: Record<BadgeVariant, string> = {
  default: "bg-blue-600 text-white hover:bg-blue-700 border-transparent",
  secondary: "bg-gray-100 text-gray-900 hover:bg-gray-200 border-transparent",
  outline: "text-gray-900 border-gray-200",
  destructive: "bg-red-600 text-white hover:bg-red-700 border-transparent",
};
</script>

<template>
  <div
    :class="
      cn(
        'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:ring-2 focus:ring-gray-950 focus:ring-offset-2 focus:outline-none',
        variants[variant],
        className
      )
    "
  >
    <slot />
  </div>
</template>
