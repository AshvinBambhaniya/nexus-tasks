<script setup lang="ts">
import { cn } from "~/utils/cn";

type BadgeVariant = "default" | "secondary" | "outline" | "destructive";

interface Props {
  variant?: BadgeVariant;
  className?: string;
}

const { variant = "default", className = "" } = defineProps<Props>();

const variants: Record<BadgeVariant, string> = {
  default:
    "bg-primary text-primary-foreground hover:opacity-90 border-transparent",
  secondary:
    "bg-secondary text-secondary-foreground hover:opacity-90 border-transparent",
  outline: "text-foreground border-border",
  destructive:
    "bg-destructive text-destructive-foreground hover:opacity-90 border-transparent",
};
</script>

<template>
  <div
    :class="
      cn(
        'focus:ring-ring inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:ring-2 focus:ring-offset-2 focus:outline-none',
        variants[variant],
        className
      )
    "
  >
    <slot />
  </div>
</template>
