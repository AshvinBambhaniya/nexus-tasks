<script setup lang="ts">
import { CheckSquare, Github } from "lucide-vue-next";
import { GITHUB_REPO_URL, APP_NAME } from "~/constants";

const year = new Date().getFullYear();
</script>

<template>
  <footer class="border-border bg-card border-t py-12">
    <div
      class="container mx-auto flex flex-col items-center justify-between gap-6 px-4 sm:px-6 md:flex-row lg:px-8"
    >
      <div class="flex items-center gap-2">
        <div
          class="bg-foreground text-background flex h-8 w-8 items-center justify-center rounded-lg"
        >
          <CheckSquare class="h-4 w-4" />
        </div>
        <div>
          <span class="text-foreground block text-sm font-bold">
            {{ APP_NAME }}
          </span>
          <span class="text-muted-foreground text-xs">Built for Builders</span>
        </div>
      </div>
      <p class="text-muted-foreground text-sm">
        © {{ year }} {{ APP_NAME }}. Distributed under MIT License.
      </p>
      <div class="flex gap-6">
        <NuxtLink
          :to="GITHUB_REPO_URL"
          target="_blank"
          class="text-muted-foreground/70 hover:text-foreground flex items-center gap-2 text-sm transition-colors"
        >
          <Github class="h-5 w-5" />
          <span>Source Code</span>
        </NuxtLink>
      </div>
    </div>
  </footer>
</template>
