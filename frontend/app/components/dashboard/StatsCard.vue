<script setup lang="ts">
import type { Component } from "vue";
import { cn } from "~/utils/cn";

type StatsColor = "blue" | "indigo" | "orange" | "green";

interface Props {
  title: string;
  value: number | string;
  icon: Component;
  isLoading?: boolean;
  color?: StatsColor;
  label?: string;
}

const {
  title,
  value,
  icon,
  isLoading = false,
  color = "blue",
  label = "",
} = defineProps<Props>();

const colors: Record<StatsColor, string> = {
  blue: "bg-primary/10 text-primary",
  indigo: "bg-indigo-500/10 text-indigo-500",
  orange: "bg-orange-500/10 text-orange-500",
  green: "bg-emerald-500/10 text-emerald-500",
};
</script>

<template>
  <UiBaseCard class="border-border p-5 shadow-sm">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-muted-foreground text-sm font-medium">{{ title }}</p>
        <template v-if="isLoading">
          <div class="bg-muted mt-1 h-8 w-16 animate-pulse rounded" />
        </template>
        <template v-else>
          <div class="mt-1 flex items-baseline gap-2">
            <span class="text-card-foreground text-2xl font-bold">{{
              value
            }}</span>
            <span
              v-if="label"
              class="text-muted-foreground text-xs font-normal"
            >
              {{ label }}
            </span>
          </div>
        </template>
      </div>
      <div :class="cn('rounded-lg p-2', colors[color])">
        <!-- Render the icon component using the built-in component tag -->
        <component :is="icon" class="h-5 w-5" />
      </div>
    </div>
  </UiBaseCard>
</template>
