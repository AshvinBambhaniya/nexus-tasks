<script setup lang="ts">
import { CheckSquare } from "lucide-vue-next";

definePageMeta({
  layout: false,
});

const email = ref("");
const password = ref("");
const { login, isLoading, error } = useAuth();

const handleSubmit = async () => {
  await login(email.value, password.value);
};
</script>

<template>
  <div class="bg-background flex min-h-screen items-center justify-center px-4">
    <div class="absolute top-4 right-4">
      <LayoutThemeToggle />
    </div>
    <div
      class="border-border bg-card w-full max-w-md space-y-8 rounded-xl border p-8 shadow-sm"
    >
      <div class="flex flex-col items-center text-center">
        <div
          class="bg-primary/10 mb-4 flex h-12 w-12 items-center justify-center rounded-full"
        >
          <CheckSquare class="text-primary h-6 w-6" />
        </div>
        <h2 class="text-card-foreground text-2xl font-bold tracking-tight">
          Welcome back
        </h2>
        <p class="text-muted-foreground mt-2 text-sm">
          Enter your email to sign in to your account
        </p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div class="grid gap-2">
            <UiBaseLabel for="email">Email</UiBaseLabel>
            <UiBaseInput
              id="email"
              v-model="email"
              placeholder="name@example.com"
              type="email"
              auto-complete="email"
              required
              :disabled="isLoading"
            />
          </div>
          <div class="grid gap-2">
            <div class="flex items-center justify-between">
              <UiBaseLabel for="password">Password</UiBaseLabel>
            </div>
            <UiBaseInput
              id="password"
              v-model="password"
              placeholder="••••••••"
              type="password"
              auto-complete="current-password"
              required
              :disabled="isLoading"
            />
          </div>
        </div>

        <div
          v-if="error"
          class="border-destructive/50 bg-destructive/10 text-destructive rounded-md border p-3 text-sm font-medium"
        >
          {{ error }}
        </div>

        <UiBaseButton :disabled="isLoading" class="w-full">
          {{ isLoading ? "Signing In..." : "Sign In" }}
        </UiBaseButton>
      </form>

      <div class="text-center text-sm">
        <span class="text-muted-foreground">Don't have an account? </span>
        <NuxtLink
          to="/register"
          class="text-primary hover:text-primary/80 font-medium hover:underline"
        >
          Sign up
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
