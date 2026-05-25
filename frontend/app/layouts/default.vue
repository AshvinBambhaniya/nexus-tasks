<script setup lang="ts"></script>

<template>
  <div class="bg-background selection:bg-primary/20 flex min-h-screen flex-col">
    <LayoutTheHeader />
    <main class="flex-1">
      <slot />
    </main>
    <LayoutTheFooter />
  </div>
</template>
