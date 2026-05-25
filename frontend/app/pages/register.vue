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
          Create an account
        </h2>
        <p class="text-muted-foreground mt-2 text-sm">
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
          class="border-destructive/50 bg-destructive/10 text-destructive rounded-md border p-3 text-sm font-medium"
        >
          {{ error }}
        </div>

        <UiBaseButton :disabled="isLoading" class="w-full">
          {{ isLoading ? "Creating Account..." : "Sign Up" }}
        </UiBaseButton>
      </form>

      <div class="text-center text-sm">
        <span class="text-muted-foreground">Already have an account? </span>
        <NuxtLink
          to="/login"
          class="text-primary hover:text-primary/80 font-medium hover:underline"
        >
          Log in
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
