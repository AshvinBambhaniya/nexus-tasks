<script setup lang="ts">
// Initialize WebSockets for the dashboard
useSocket();

// Restore active timer on dashboard mount
const { fetchActiveTimer } = useTimeTracking();
onMounted(() => {
  fetchActiveTimer();
});
</script>

<template>
  <div class="bg-background flex h-screen flex-col overflow-hidden">
    <LayoutActiveTimerBar />
    <div class="flex flex-1 overflow-hidden">
      <LayoutSidebar />
      <main class="flex-1 overflow-y-auto p-8">
        <slot />
      </main>
    </div>
    <LayoutCommandPalette />
  </div>
</template>
