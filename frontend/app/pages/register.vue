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

const config = useRuntimeConfig();

const handleSSOLogin = () => {
  window.location.href = `${config.public.apiUrl}/api/v2/auth/sso/login`;
};

const error = computed(() => validationError.value || authError.value);
</script>

<template>
  <div
    class="bg-background selection:bg-primary/30 relative flex min-h-screen items-center justify-center overflow-hidden px-4 transition-colors duration-300"
  >
    <!-- Ambient glowing background -->
    <div
      class="pointer-events-none absolute inset-0 flex items-center justify-center mix-blend-multiply dark:mix-blend-screen"
      aria-hidden="true"
    >
      <div
        class="bg-primary/10 h-[600px] w-[600px] rounded-full blur-[120px]"
      />
      <div
        class="absolute h-[400px] w-[400px] -translate-x-1/4 translate-y-1/4 rounded-full bg-sky-500/10 blur-[100px]"
      />
    </div>

    <div class="absolute top-4 right-4 z-50">
      <LayoutThemeToggle />
    </div>

    <!-- Glassmorphic Card -->
    <div
      class="border-border bg-card/60 relative z-10 w-full max-w-[400px] rounded-2xl border p-8 shadow-2xl backdrop-blur-xl transition-colors duration-300"
    >
      <div class="mb-8 flex flex-col items-center text-center">
        <div
          class="bg-muted/50 border-border mb-4 flex h-12 w-12 items-center justify-center rounded-xl border shadow-inner"
        >
          <CheckSquare class="text-primary h-6 w-6" />
        </div>
        <h2 class="text-foreground text-2xl font-bold tracking-tight">
          Create an account
        </h2>
        <p class="text-muted-foreground mt-2 text-sm">
          Start building your workspace
        </p>
      </div>

      <div class="mb-6">
        <UiBaseButton
          variant="outline"
          class="bg-muted/30 hover:bg-muted/50 text-foreground h-11 w-full transition-all"
          @click="handleSSOLogin"
        >
          <CheckSquare class="mr-2 h-4 w-4" /> Continue with SSO
        </UiBaseButton>
      </div>

      <div class="relative mb-6">
        <div class="absolute inset-0 flex items-center">
          <div class="border-border w-full border-t" />
        </div>
        <div class="relative flex justify-center text-xs uppercase">
          <span class="bg-card text-muted-foreground px-2 backdrop-blur-xl"
            >Or sign up with email</span
          >
        </div>
      </div>

      <form class="space-y-4" @submit.prevent="handleSubmit">
        <div class="space-y-4">
          <div class="space-y-2">
            <UiBaseLabel for="fullName" class="text-foreground/80"
              >Full Name</UiBaseLabel
            >
            <UiBaseInput
              id="fullName"
              v-model="fullName"
              placeholder="John Doe"
              type="text"
              required
              :disabled="isLoading"
              class="bg-muted/30 text-foreground placeholder:text-muted-foreground focus-visible:ring-primary h-11"
            />
          </div>
          <div class="space-y-2">
            <UiBaseLabel for="email" class="text-foreground/80"
              >Email Address</UiBaseLabel
            >
            <UiBaseInput
              id="email"
              v-model="email"
              placeholder="name@example.com"
              type="email"
              auto-complete="email"
              required
              :disabled="isLoading"
              class="bg-muted/30 text-foreground placeholder:text-muted-foreground focus-visible:ring-primary h-11"
            />
          </div>
          <div class="space-y-2">
            <UiBaseLabel for="password" class="text-foreground/80"
              >Password</UiBaseLabel
            >
            <UiBaseInput
              id="password"
              v-model="password"
              placeholder="••••••••"
              type="password"
              auto-complete="new-password"
              required
              min-length="8"
              :disabled="isLoading"
              class="bg-muted/30 text-foreground placeholder:text-muted-foreground focus-visible:ring-primary h-11"
            />
          </div>
          <div class="space-y-2">
            <UiBaseLabel for="confirmPassword" class="text-foreground/80"
              >Confirm Password</UiBaseLabel
            >
            <UiBaseInput
              id="confirmPassword"
              v-model="confirmPassword"
              placeholder="••••••••"
              type="password"
              auto-complete="new-password"
              required
              :disabled="isLoading"
              class="bg-muted/30 text-foreground placeholder:text-muted-foreground focus-visible:ring-primary h-11"
            />
          </div>
        </div>

        <div
          v-if="error"
          class="border-destructive/50 bg-destructive/10 text-destructive-foreground rounded-md border p-3 text-sm font-medium"
        >
          {{ error }}
        </div>

        <UiBaseButton
          :disabled="isLoading"
          class="shadow-primary/20 hover:shadow-primary/40 mt-2 h-11 w-full shadow-lg transition-all"
        >
          {{ isLoading ? "Creating Account..." : "Sign Up with Email" }}
        </UiBaseButton>
      </form>

      <div class="mt-6 text-center text-sm">
        <span class="text-muted-foreground">Already have an account? </span>
        <NuxtLink
          to="/login"
          class="text-primary hover:text-primary/80 font-medium transition-colors hover:underline"
        >
          Log in
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
