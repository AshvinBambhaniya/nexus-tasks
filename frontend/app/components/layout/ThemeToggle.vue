<script setup lang="ts">
import { Sun, Moon, Monitor } from "lucide-vue-next";

const colorMode = useColorMode();

const modes = [
  { name: "system", icon: Monitor },
  { name: "light", icon: Sun },
  { name: "dark", icon: Moon },
];

const nextMode = () => {
  const currentIndex = modes.findIndex((m) => m.name === colorMode.preference);
  const nextIndex = (currentIndex + 1) % modes.length;
  colorMode.preference = modes[nextIndex].name;
};

const currentIcon = computed(() => {
  const mode = modes.find((m) => m.name === colorMode.preference);
  return mode ? mode.icon : Monitor;
});
</script>

<template>
  <button
    class="text-muted-foreground hover:bg-muted hover:text-foreground flex h-8 w-8 items-center justify-center rounded-md transition-colors"
    :title="`Switch to next theme (current: ${colorMode.preference})`"
    @click="nextMode"
  >
    <component :is="currentIcon" class="h-4 w-4" />
    <span class="sr-only">Toggle theme</span>
  </button>
</template>
