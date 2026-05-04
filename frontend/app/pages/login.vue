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
  <div class="flex min-h-screen items-center justify-center bg-gray-50 px-4">
    <div
      class="w-full max-w-md space-y-8 rounded-xl border border-gray-200 bg-white p-8 shadow-sm"
    >
      <div class="flex flex-col items-center text-center">
        <div
          class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-blue-100"
        >
          <CheckSquare class="h-6 w-6 text-blue-600" />
        </div>
        <h2 class="text-2xl font-bold tracking-tight text-gray-900">
          Welcome back
        </h2>
        <p class="mt-2 text-sm text-gray-500">
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
          class="rounded-md border border-red-200 bg-red-50 p-3 text-sm font-medium text-red-600"
        >
          {{ error }}
        </div>

        <UiBaseButton :disabled="isLoading" class="w-full">
          {{ isLoading ? "Signing In..." : "Sign In" }}
        </UiBaseButton>
      </form>

      <div class="text-center text-sm">
        <span class="text-gray-500">Don't have an account? </span>
        <NuxtLink
          to="/register"
          class="font-medium text-blue-600 hover:text-blue-500 hover:underline"
        >
          Sign up
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
