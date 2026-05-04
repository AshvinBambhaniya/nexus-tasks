<script setup lang="ts">
// Initialize WebSockets for the dashboard
useSocket();
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-white">
    <LayoutSidebar />
    <main class="flex-1 overflow-y-auto bg-white p-8">
      <slot />
    </main>
    <LayoutCommandPalette />
  </div>
</template>
