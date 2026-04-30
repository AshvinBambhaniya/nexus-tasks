<script setup lang="ts">
import { CheckSquare, Github } from "lucide-vue-next";
import { GITHUB_REPO_URL, APP_NAME } from "~/constants";

const year = new Date().getFullYear();
</script>

<template>
  <footer class="border-t border-gray-100 bg-white py-12">
    <div
      class="container mx-auto flex flex-col items-center justify-between gap-6 px-4 sm:px-6 md:flex-row lg:px-8"
    >
      <div class="flex items-center gap-2">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-900 text-white"
        >
          <CheckSquare class="h-4 w-4" />
        </div>
        <div>
          <span class="block text-sm font-bold text-gray-900">
            {{ APP_NAME }}
          </span>
          <span class="text-xs text-gray-500">Built for Builders</span>
        </div>
      </div>
      <p class="text-sm text-gray-500">
        © {{ year }} {{ APP_NAME }}. Distributed under MIT License.
      </p>
      <div class="flex gap-6">
        <NuxtLink
          :to="GITHUB_REPO_URL"
          target="_blank"
          class="flex items-center gap-2 text-sm text-gray-400 transition-colors hover:text-gray-900"
        >
          <Github class="h-5 w-5" />
          <span>Source Code</span>
        </NuxtLink>
      </div>
    </div>
  </footer>
</template>
