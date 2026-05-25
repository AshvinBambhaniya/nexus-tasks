<script setup lang="ts">
import {
  LogOut,
  Building,
  Loader2,
  Check,
  User as UserIcon,
  Mail,
  Fingerprint,
  Camera,
  ExternalLink,
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
  <div class="space-y-8 pb-10">
    <!-- Profile Section -->
    <section>
      <div class="mb-4">
        <h2 class="text-foreground text-lg font-semibold">
          Personal Information
        </h2>
        <p class="text-muted-foreground text-sm">
          Update your profile details and how others see you.
        </p>
      </div>

      <UiBaseCard
        class="border-border overflow-hidden shadow-sm transition-shadow hover:shadow-md"
      >
        <div class="p-8">
          <div
            class="border-border flex flex-col items-center gap-8 border-b pb-8 sm:flex-row"
          >
            <div class="group relative">
              <UiBaseAvatar
                :fallback="userInitial"
                class="border-card ring-border h-28 w-28 border-4 text-3xl shadow-lg ring-1"
              />
              <button
                class="border-border bg-card text-muted-foreground hover:bg-muted hover:text-primary absolute right-0 bottom-0 rounded-full border p-2 shadow-sm transition-all"
              >
                <Camera class="h-4 w-4" />
              </button>
            </div>

            <div class="space-y-2 text-center sm:text-left">
              <h3 class="text-card-foreground text-2xl font-bold">
                {{ user?.full_name || "Nexus User" }}
              </h3>
              <div
                class="text-muted-foreground flex flex-wrap items-center justify-center gap-4 text-sm sm:justify-start"
              >
                <span class="flex items-center gap-1.5">
                  <Mail class="text-muted-foreground/60 h-4 w-4" />
                  {{ user?.email }}
                </span>
                <span class="flex items-center gap-1.5 font-mono text-xs">
                  <Fingerprint class="text-muted-foreground/60 h-4 w-4" />
                  ID: {{ user?.id }}
                </span>
              </div>
              <div class="pt-1">
                <UiBaseBadge
                  variant="secondary"
                  class="bg-emerald-500/10 text-[10px] font-bold text-emerald-500 uppercase ring-1 ring-emerald-500/20"
                >
                  Active Account
                </UiBaseBadge>
              </div>
            </div>
          </div>

          <form class="mt-8 space-y-6" @submit.prevent="handleUpdateProfile">
            <div class="grid gap-6 sm:grid-cols-2">
              <div class="space-y-2">
                <UiBaseLabel
                  for="fullName"
                  class="text-muted-foreground text-xs font-bold tracking-wide uppercase"
                  >Full Name</UiBaseLabel
                >
                <div class="relative">
                  <UserIcon
                    class="text-muted-foreground/60 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
                  />
                  <UiBaseInput
                    id="fullName"
                    v-model="fullName"
                    placeholder="John Doe"
                    class="focus:ring-primary/10 pl-10 transition-all focus:ring-2"
                  />
                </div>
              </div>
              <div class="space-y-2">
                <UiBaseLabel
                  for="email"
                  class="text-muted-foreground text-xs font-bold tracking-wide uppercase"
                  >Email Address</UiBaseLabel
                >
                <div class="relative">
                  <Mail
                    class="text-muted-foreground/60 absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2"
                  />
                  <UiBaseInput
                    id="email"
                    v-model="email"
                    type="email"
                    placeholder="john@example.com"
                    class="focus:ring-primary/10 pl-10 transition-all focus:ring-2"
                  />
                </div>
              </div>
            </div>

            <div
              class="border-border flex items-center justify-end border-t pt-6"
            >
              <UiBaseButton
                type="submit"
                :disabled="isUpdating"
                class="min-w-[140px] shadow-sm transition-all active:scale-95"
              >
                <Loader2 v-if="isUpdating" class="mr-2 h-4 w-4 animate-spin" />
                <Check
                  v-else-if="isSuccess"
                  class="text-primary-foreground mr-2 h-4 w-4"
                />
                {{
                  isUpdating
                    ? "Saving..."
                    : isSuccess
                      ? "Changes Saved"
                      : "Save Changes"
                }}
              </UiBaseButton>
            </div>
          </form>
        </div>
      </UiBaseCard>
    </section>

    <!-- Workspace Section -->
    <section>
      <div class="mb-4">
        <h2 class="text-foreground text-lg font-semibold">Current Workspace</h2>
        <p class="text-muted-foreground text-sm">
          Settings and information about your active project environment.
        </p>
      </div>

      <UiBaseCard class="border-border shadow-sm">
        <div class="p-8">
          <div v-if="activeWorkspace" class="grid gap-8 sm:grid-cols-2">
            <div class="space-y-3">
              <UiBaseLabel
                class="text-muted-foreground text-xs font-bold tracking-wide uppercase"
                >Workspace Identity</UiBaseLabel
              >
              <div
                class="border-border bg-muted/50 hover:bg-muted flex items-center gap-4 rounded-xl border p-4 transition-colors"
              >
                <div
                  class="border-primary/20 bg-card text-primary flex h-12 w-12 items-center justify-center rounded-lg border shadow-sm"
                >
                  <Building class="h-6 w-6" />
                </div>
                <div>
                  <div class="text-card-foreground text-sm font-bold">
                    {{ activeWorkspace.name }}
                  </div>
                  <div class="text-muted-foreground text-xs italic">
                    Workspace ID: #{{ activeWorkspace.id }}
                  </div>
                </div>
              </div>
            </div>

            <div class="space-y-3">
              <UiBaseLabel
                class="text-muted-foreground text-xs font-bold tracking-wide uppercase"
                >Environment Type</UiBaseLabel
              >
              <div
                class="border-border bg-muted/50 hover:bg-muted flex items-center gap-4 rounded-xl border p-4 transition-colors"
              >
                <div
                  class="bg-card flex h-12 w-12 items-center justify-center rounded-lg border border-purple-500/20 text-purple-500 shadow-sm"
                >
                  <ExternalLink class="h-6 w-6" />
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span
                      class="text-card-foreground text-sm font-bold capitalize"
                      >{{ activeWorkspace.type.toLowerCase() }}</span
                    >
                    <UiBaseBadge class="bg-primary/10 text-primary text-[10px]"
                      >Team</UiBaseBadge
                    >
                  </div>
                  <div class="text-muted-foreground text-xs">
                    Collaborative workspace
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div
            v-else
            class="border-border bg-muted/50 rounded-xl border border-dashed p-12 text-center"
          >
            <Building class="text-muted mx-auto mb-4 h-12 w-12" />
            <h3 class="text-card-foreground font-medium">
              No Active Workspace
            </h3>
            <p class="text-muted-foreground mt-1 text-sm text-pretty">
              Select or create a workspace to see details here.
            </p>
            <UiBaseButton
              variant="outline"
              size="sm"
              class="mt-4"
              to="/workspaces"
              >Go to Workspaces</UiBaseButton
            >
          </div>
        </div>
      </UiBaseCard>
    </section>

    <!-- Danger Zone -->
    <section>
      <div class="mb-4">
        <h2 class="text-destructive text-lg font-semibold">
          Security & Sessions
        </h2>
        <p class="text-destructive/80 text-sm">
          Sensitive actions that affect your account access.
        </p>
      </div>

      <div
        class="border-destructive/20 bg-destructive/10 ring-destructive/20 overflow-hidden rounded-xl border shadow-sm ring-1"
      >
        <div
          class="flex flex-col items-center justify-between gap-6 p-8 sm:flex-row"
        >
          <div class="space-y-1 text-center sm:text-left">
            <h3 class="text-destructive font-bold">Sign out of Nexus</h3>
            <p class="text-destructive/80 text-sm">
              You will be redirected to the login page. All active sessions will
              remain, but you will need to re-authenticate on this device.
            </p>
          </div>
          <UiBaseButton
            variant="destructive"
            class="shadow-destructive/20 min-w-[140px] shadow-md active:scale-95"
            @click="logout"
          >
            <LogOut class="mr-2 h-4 w-4" /> Log out
          </UiBaseButton>
        </div>
      </div>
    </section>
  </div>
</template>
