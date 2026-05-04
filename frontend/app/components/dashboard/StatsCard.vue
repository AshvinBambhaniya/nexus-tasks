<script setup lang="ts">
import type { Component } from "vue";

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
  blue: "bg-blue-50 text-blue-600",
  indigo: "bg-indigo-50 text-indigo-600",
  orange: "bg-orange-50 text-orange-600",
  green: "bg-green-50 text-green-600",
};
</script>

<template>
  <UiBaseCard class="border-gray-200 p-5 shadow-sm">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-sm font-medium text-gray-500">{{ title }}</p>
        <template v-if="isLoading">
          <div class="mt-1 h-8 w-16 animate-pulse rounded bg-gray-100" />
        </template>
        <template v-else>
          <div class="mt-1 flex items-baseline gap-2">
            <span class="text-2xl font-bold text-gray-900">{{ value }}</span>
            <span v-if="label" class="text-xs font-normal text-gray-500">
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
