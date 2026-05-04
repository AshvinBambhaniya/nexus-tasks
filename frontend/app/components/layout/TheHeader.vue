<script setup lang="ts">
import { CheckSquare } from "lucide-vue-next";
import { GITHUB_REPO_URL, APP_NAME } from "~/constants";
import { useUsersStore } from "~/stores/user";

const userStore = useUsersStore();
const user = computed(() => userStore.userData);

// Since setUserDataStore doesn't track internal loading,
// we track it locally or just rely on reactive data existence.
const isLoading = ref(false);

onMounted(async () => {
  if (!user.value) {
    isLoading.value = true;
    await setUserDataStore();
    isLoading.value = false;
  }
});
</script>

<template>
  <header
    class="sticky top-0 z-50 w-full border-b border-gray-100 bg-white/80 backdrop-blur-md"
  >
    <div
      class="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8"
    >
      <NuxtLink to="/" class="flex items-center gap-2">
        <div
          class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white shadow-md shadow-blue-200"
        >
          <CheckSquare class="h-5 w-5" />
        </div>
        <span class="text-xl font-bold tracking-tight text-gray-900">
          {{ APP_NAME }}
        </span>
      </NuxtLink>

      <nav
        class="hidden items-center gap-8 text-sm font-medium text-gray-600 md:flex"
      >
        <NuxtLink to="/#features" class="transition-colors hover:text-blue-600">
          Features
        </NuxtLink>
        <NuxtLink
          to="/#tech-stack"
          class="transition-colors hover:text-blue-600"
        >
          Tech Stack
        </NuxtLink>
        <NuxtLink
          :to="GITHUB_REPO_URL"
          target="_blank"
          class="transition-colors hover:text-blue-600"
        >
          GitHub
        </NuxtLink>
      </nav>

      <div class="flex items-center gap-4">
        <template v-if="isLoading">
          <UiBaseButton disabled variant="ghost" size="sm">
            Loading...
          </UiBaseButton>
        </template>
        <template v-else-if="user">
          <NuxtLink to="/dashboard">
            <UiBaseButton class="font-semibold shadow-md shadow-blue-100">
              Go to Dashboard
            </UiBaseButton>
          </NuxtLink>
        </template>
        <template v-else>
          <NuxtLink to="/login" class="hidden sm:block">
            <UiBaseButton
              variant="ghost"
              class="font-medium text-gray-600 hover:text-gray-900"
            >
              Log in
            </UiBaseButton>
          </NuxtLink>
          <NuxtLink to="/register">
            <UiBaseButton class="font-semibold shadow-md shadow-blue-100">
              Get Started
            </UiBaseButton>
          </NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>
