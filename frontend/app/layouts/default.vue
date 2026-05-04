<script setup lang="ts"></script>

<template>
  <div class="flex min-h-screen flex-col bg-white selection:bg-blue-100">
    <LayoutTheHeader />
    <main class="flex-1">
      <slot />
    </main>
    <LayoutTheFooter />
  </div>
</template>
