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
    class="border-border bg-background/80 sticky top-0 z-50 w-full border-b backdrop-blur-md"
  >
    <div
      class="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8"
    >
      <NuxtLink to="/" class="flex items-center gap-2">
        <div
          class="bg-primary text-primary-foreground shadow-primary/20 flex h-8 w-8 items-center justify-center rounded-lg shadow-md"
        >
          <CheckSquare class="h-5 w-5" />
        </div>
        <span class="text-foreground text-xl font-bold tracking-tight">
          {{ APP_NAME }}
        </span>
      </NuxtLink>

      <nav
        class="text-muted-foreground hidden items-center gap-8 text-sm font-medium md:flex"
      >
        <NuxtLink to="/#features" class="hover:text-primary transition-colors">
          Features
        </NuxtLink>
        <NuxtLink
          to="/#tech-stack"
          class="hover:text-primary transition-colors"
        >
          Tech Stack
        </NuxtLink>
        <NuxtLink
          :to="GITHUB_REPO_URL"
          target="_blank"
          class="hover:text-primary transition-colors"
        >
          GitHub
        </NuxtLink>
      </nav>

      <div class="flex items-center gap-2 sm:gap-4">
        <LayoutThemeToggle />
        <template v-if="isLoading">
          <UiBaseButton disabled variant="ghost" size="sm">
            Loading...
          </UiBaseButton>
        </template>
        <template v-else-if="user">
          <NuxtLink to="/dashboard">
            <UiBaseButton class="shadow-primary/10 font-semibold shadow-md">
              Go to Dashboard
            </UiBaseButton>
          </NuxtLink>
        </template>
        <template v-else>
          <NuxtLink to="/login" class="hidden sm:block">
            <UiBaseButton
              variant="ghost"
              class="text-muted-foreground hover:text-foreground font-medium"
            >
              Log in
            </UiBaseButton>
          </NuxtLink>
          <NuxtLink to="/register">
            <UiBaseButton class="shadow-primary/10 font-semibold shadow-md">
              Get Started
            </UiBaseButton>
          </NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>
