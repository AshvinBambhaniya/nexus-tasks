<script setup lang="ts">
import {
  LogOut,
  Building,
  Loader2,
  Check,
  Mail,
  Camera,
} from "lucide-vue-next";
import { useUsersStore } from "~/stores/user";
import { useWorkspaces } from "~/composables/useWorkspaces";
import { useAuth } from "~/composables/useAuth";
import type { User } from "~/types";

const userStore = useUsersStore();
const { activeWorkspace } = useWorkspaces();
const { logout } = useAuth();

const fullName = ref("");
const email = ref("");
const isUpdating = ref(false);
const isSuccess = ref(false);

const user = computed(() => userStore.userData);

onMounted(() => {
  if (user.value) {
    fullName.value = user.value.full_name || "";
    email.value = user.value.email || "";
  }
});

watch(
  () => user.value,
  (newUser) => {
    if (newUser) {
      fullName.value = newUser.full_name || "";
      email.value = newUser.email || "";
    }
  }
);

const handleUpdateProfile = async () => {
  isUpdating.value = true;
  isSuccess.value = false;
  try {
    const updatedUser = await useMutation<User>("/api/v2/auth/me", {
      method: "PATCH",
      body: { full_name: fullName.value, email: email.value },
    });
    userStore.setUserData(updatedUser);
    isSuccess.value = true;
    setTimeout(() => (isSuccess.value = false), 3000);
  } catch (err) {
    const msg = getApiErrorMessage(err, "Failed to update profile");
    alert(msg);
  } finally {
    isUpdating.value = false;
  }
};

const userInitial = computed(() => {
  return (user.value?.full_name || user.value?.email || "??")
    .substring(0, 2)
    .toUpperCase();
});
</script>

<template>
  <div class="space-y-12">
    <!-- Profile Section -->
    <section
      class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
    >
      <div class="lg:w-1/3">
        <h2 class="text-foreground text-lg font-semibold tracking-tight">
          Personal Information
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Update your profile details and how others see you on the platform.
        </p>
      </div>

      <UiBaseCard class="border-border bg-card lg:w-2/3">
        <div class="p-6">
          <div class="flex flex-col gap-8 sm:flex-row sm:items-center">
            <div class="group relative">
              <UiBaseAvatar
                :fallback="userInitial"
                class="h-24 w-24 border-0 text-3xl shadow-sm"
              />
              <button
                class="border-border bg-background hover:bg-muted text-muted-foreground hover:text-foreground absolute right-0 bottom-0 rounded-full border p-2 shadow-sm transition-colors"
                title="Update avatar"
              >
                <Camera class="h-4 w-4" />
              </button>
            </div>

            <div class="space-y-1">
              <h3 class="text-card-foreground text-xl font-bold">
                {{ user?.full_name || "Nexus User" }}
              </h3>
              <div
                class="text-muted-foreground flex flex-wrap items-center gap-4 text-sm"
              >
                <span class="flex items-center gap-1.5">
                  <Mail class="h-4 w-4" />
                  {{ user?.email }}
                </span>
              </div>
              <div class="flex items-center gap-2 pt-2">
                <UiBaseBadge
                  variant="secondary"
                  class="bg-emerald-500/10 text-[10px] font-bold text-emerald-600 uppercase dark:text-emerald-500"
                >
                  Active Account
                </UiBaseBadge>
                <span class="text-muted-foreground/60 font-mono text-xs">
                  ID: {{ user?.id }}
                </span>
              </div>
            </div>
          </div>

          <form class="mt-8 space-y-6" @submit.prevent="handleUpdateProfile">
            <div class="grid gap-6 sm:grid-cols-2">
              <div class="space-y-2">
                <UiBaseLabel
                  for="fullName"
                  class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
                  >Full Name</UiBaseLabel
                >
                <UiBaseInput
                  id="fullName"
                  v-model="fullName"
                  placeholder="John Doe"
                />
              </div>
              <div class="space-y-2">
                <UiBaseLabel
                  for="email"
                  class="text-muted-foreground text-xs font-semibold tracking-wide uppercase"
                  >Email Address</UiBaseLabel
                >
                <UiBaseInput
                  id="email"
                  v-model="email"
                  type="email"
                  placeholder="john@example.com"
                />
              </div>
            </div>

            <div class="flex items-center justify-end pt-4">
              <UiBaseButton
                type="submit"
                :disabled="isUpdating"
                class="min-w-[120px]"
              >
                <Loader2 v-if="isUpdating" class="mr-2 h-4 w-4 animate-spin" />
                <Check v-else-if="isSuccess" class="mr-2 h-4 w-4" />
                {{ isUpdating ? "Saving..." : isSuccess ? "Saved" : "Save" }}
              </UiBaseButton>
            </div>
          </form>
        </div>
      </UiBaseCard>
    </section>

    <div class="border-border border-t" />

    <!-- Workspace Section -->
    <section
      class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
    >
      <div class="lg:w-1/3">
        <h2 class="text-foreground text-lg font-semibold tracking-tight">
          Current Workspace
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Settings and information about your active project environment.
        </p>
      </div>

      <UiBaseCard class="border-border bg-card lg:w-2/3">
        <div class="p-6">
          <div
            v-if="activeWorkspace"
            class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between"
          >
            <div class="flex items-center gap-4">
              <div
                class="border-border bg-background text-foreground flex h-12 w-12 shrink-0 items-center justify-center rounded-lg border shadow-sm"
              >
                <Building class="h-6 w-6" />
              </div>
              <div>
                <div
                  class="text-card-foreground flex items-center gap-2 text-base font-bold"
                >
                  {{ activeWorkspace.name }}
                  <UiBaseBadge
                    class="bg-primary/10 text-primary hover:bg-primary/20 text-[10px]"
                    >Team</UiBaseBadge
                  >
                </div>
                <div class="text-muted-foreground mt-0.5 text-xs">
                  ID: {{ activeWorkspace.id }}
                </div>
              </div>
            </div>
            <div>
              <UiBaseButton variant="outline" size="sm" to="/workspaces">
                Switch Workspace
              </UiBaseButton>
            </div>
          </div>

          <div
            v-else
            class="flex flex-col items-center justify-center p-8 text-center"
          >
            <Building class="text-muted-foreground/50 mx-auto mb-3 h-10 w-10" />
            <h3 class="text-card-foreground text-sm font-medium">
              No Active Workspace
            </h3>
            <p class="text-muted-foreground mt-1 text-sm text-pretty">
              Select or create a workspace to manage settings.
            </p>
            <UiBaseButton
              variant="outline"
              size="sm"
              class="mt-4"
              to="/workspaces"
            >
              Go to Workspaces
            </UiBaseButton>
          </div>
        </div>
      </UiBaseCard>
    </section>

    <div class="border-border border-t" />

    <!-- Danger Zone -->
    <section
      class="flex flex-col gap-6 lg:flex-row lg:items-start lg:justify-between"
    >
      <div class="lg:w-1/3">
        <h2 class="text-destructive text-lg font-semibold tracking-tight">
          Security & Sessions
        </h2>
        <p class="text-muted-foreground mt-1 text-sm">
          Sensitive actions that affect your account access.
        </p>
      </div>

      <div
        class="border-destructive/30 bg-destructive/5 overflow-hidden rounded-xl border lg:w-2/3"
      >
        <div
          class="flex flex-col gap-4 p-6 sm:flex-row sm:items-center sm:justify-between"
        >
          <div>
            <h3 class="text-destructive font-semibold">Sign out of Nexus</h3>
            <p class="text-muted-foreground mt-1 text-sm">
              You will be redirected to the login page.
            </p>
          </div>
          <UiBaseButton variant="destructive" @click="logout">
            <LogOut class="mr-2 h-4 w-4" /> Log out
          </UiBaseButton>
        </div>
      </div>
    </section>
  </div>
</template>
