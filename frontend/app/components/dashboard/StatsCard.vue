<script setup lang="ts">
import type { Component } from "vue";
import { cn } from "~/utils/cn";

type StatsColor = "blue" | "indigo" | "orange" | "green" | "purple";

interface Props {
  title: string;
  value: number | string;
  icon: Component;
  isLoading?: boolean;
  color?: StatsColor;
  label?: string;
  trend?: string;
}

const {
  title,
  value,
  icon,
  isLoading = false,
  color = "blue",
  label = "",
  trend = "",
} = defineProps<Props>();

const textColors: Record<StatsColor, string> = {
  blue: "text-blue-600 dark:text-blue-400 dark:drop-shadow-[0_0_12px_rgba(96,165,250,0.6)]",
  indigo:
    "text-indigo-600 dark:text-indigo-400 dark:drop-shadow-[0_0_12px_rgba(129,140,248,0.6)]",
  orange:
    "text-orange-600 dark:text-orange-400 dark:drop-shadow-[0_0_12px_rgba(251,146,60,0.6)]",
  green:
    "text-emerald-600 dark:text-emerald-400 dark:drop-shadow-[0_0_12px_rgba(52,211,153,0.6)]",
  purple:
    "text-purple-600 dark:text-purple-400 dark:drop-shadow-[0_0_12px_rgba(192,132,252,0.6)]",
};

const iconColors: Record<StatsColor, string> = {
  blue: "text-blue-500",
  indigo: "text-indigo-500",
  orange: "text-orange-500",
  green: "text-emerald-500",
  purple: "text-purple-500",
};
</script>

<template>
  <div
    class="group border-border bg-card hover:border-foreground/10 hover:bg-muted/50 relative overflow-hidden rounded-xl border p-6 transition-all"
  >
    <div class="flex items-center justify-between">
      <h3 class="text-muted-foreground text-sm font-medium">{{ title }}</h3>
      <component
        :is="icon"
        :class="cn('h-4 w-4 opacity-50', iconColors[color])"
      />
    </div>

    <div class="mt-4 flex items-baseline gap-2">
      <template v-if="isLoading">
        <div class="bg-muted h-12 w-24 animate-pulse rounded-md" />
      </template>
      <template v-else>
        <span
          :class="cn('text-5xl font-bold tracking-tighter', textColors[color])"
        >
          {{ value }}
        </span>
      </template>
    </div>

    <div class="mt-6 flex items-center justify-between text-[10px]">
      <span
        class="text-muted-foreground/70 font-semibold tracking-widest uppercase"
        >{{ title }}</span
      >
      <span
        v-if="trend"
        :class="cn('flex items-center gap-1 font-medium', iconColors[color])"
      >
        ↗ {{ trend }}
      </span>
      <span v-else-if="label" class="text-muted-foreground font-medium">
        {{ label }}
      </span>
    </div>
  </div>
</template>
