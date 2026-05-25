<script setup lang="ts">
// Initialize WebSockets for the dashboard
useSocket();
</script>

<template>
  <div class="bg-background flex h-screen overflow-hidden">
    <LayoutSidebar />
    <main class="flex-1 overflow-y-auto p-8">
      <slot />
    </main>
    <LayoutCommandPalette />
  </div>
</template>
