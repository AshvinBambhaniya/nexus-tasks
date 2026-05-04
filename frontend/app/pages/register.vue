<script setup lang="ts">
import { CheckSquare } from "lucide-vue-next";

definePageMeta({
  layout: false,
});

const fullName = ref("");
const email = ref("");
const password = ref("");
const confirmPassword = ref("");
const validationError = ref<string | null>(null);

const { register, isLoading, error: authError } = useAuth();

const handleSubmit = async () => {
  validationError.value = null;

  if (password.value !== confirmPassword.value) {
    validationError.value = "Passwords do not match";
    return;
  }

  await register(email.value, password.value, fullName.value);
};

const error = computed(() => validationError.value || authError.value);
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
          Create an account
        </h2>
        <p class="mt-2 text-sm text-gray-500">
          Enter your details below to get started
        </p>
      </div>

      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div class="grid gap-2">
            <UiBaseLabel for="fullName">Full Name</UiBaseLabel>
            <UiBaseInput
              id="fullName"
              v-model="fullName"
              placeholder="John Doe"
              type="text"
              required
              :disabled="isLoading"
            />
          </div>
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
            <UiBaseLabel for="password">Password</UiBaseLabel>
            <UiBaseInput
              id="password"
              v-model="password"
              placeholder="••••••••"
              type="password"
              auto-complete="new-password"
              required
              min-length="8"
              :disabled="isLoading"
            />
          </div>
          <div class="grid gap-2">
            <UiBaseLabel for="confirmPassword">Confirm Password</UiBaseLabel>
            <UiBaseInput
              id="confirmPassword"
              v-model="confirmPassword"
              placeholder="••••••••"
              type="password"
              auto-complete="new-password"
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
          {{ isLoading ? "Creating Account..." : "Sign Up" }}
        </UiBaseButton>
      </form>

      <div class="text-center text-sm">
        <span class="text-gray-500">Already have an account? </span>
        <NuxtLink
          to="/login"
          class="font-medium text-blue-600 hover:text-blue-500 hover:underline"
        >
          Log in
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
