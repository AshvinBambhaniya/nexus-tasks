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
    const updatedUser = await useMutation("/api/v1/auth/me", {
      method: "PATCH",
      body: { full_name: fullName.value, email: email.value },
    });
    userStore.setUserData(updatedUser as any);
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
        <h2 class="text-lg font-semibold text-gray-900">Personal Information</h2>
        <p class="text-sm text-gray-500">Update your profile details and how others see you.</p>
      </div>
      
      <UiBaseCard class="overflow-hidden border-gray-100 shadow-sm transition-shadow hover:shadow-md">
        <div class="p-8">
          <div class="flex flex-col items-center gap-8 border-b border-gray-50 pb-8 sm:flex-row">
            <div class="group relative">
              <UiBaseAvatar :fallback="userInitial" class="h-28 w-28 border-4 border-white text-3xl shadow-lg ring-1 ring-gray-100" />
              <button class="absolute right-0 bottom-0 rounded-full border border-gray-200 bg-white p-2 text-gray-500 shadow-sm transition-all hover:bg-gray-50 hover:text-blue-600">
                <Camera class="h-4 w-4" />
              </button>
            </div>
            
            <div class="space-y-2 text-center sm:text-left">
              <h3 class="text-2xl font-bold text-gray-900">
                {{ user?.full_name || "Nexus User" }}
              </h3>
              <div class="flex flex-wrap items-center justify-center gap-4 text-sm text-gray-500 sm:justify-start">
                <span class="flex items-center gap-1.5">
                  <Mail class="h-4 w-4 text-gray-400" />
                  {{ user?.email }}
                </span>
                <span class="flex items-center gap-1.5 font-mono text-xs">
                  <Fingerprint class="h-4 w-4 text-gray-400" />
                  ID: {{ user?.id }}
                </span>
              </div>
              <div class="pt-1">
                <UiBaseBadge variant="secondary" class="bg-green-50 text-[10px] font-bold text-green-700 ring-1 ring-green-100 uppercase">
                  Active Account
                </UiBaseBadge>
              </div>
            </div>
          </div>

          <form @submit.prevent="handleUpdateProfile" class="mt-8 space-y-6">
            <div class="grid gap-6 sm:grid-cols-2">
              <div class="space-y-2">
                <UiBaseLabel for="fullName" class="text-xs font-bold tracking-wide text-gray-500 uppercase">Full Name</UiBaseLabel>
                <div class="relative">
                  <UserIcon class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  <UiBaseInput
                    id="fullName"
                    v-model="fullName"
                    placeholder="John Doe"
                    class="pl-10 transition-all focus:ring-2 focus:ring-blue-100"
                  />
                </div>
              </div>
              <div class="space-y-2">
                <UiBaseLabel for="email" class="text-xs font-bold tracking-wide text-gray-500 uppercase">Email Address</UiBaseLabel>
                <div class="relative">
                  <Mail class="absolute top-1/2 left-3 h-4 w-4 -translate-y-1/2 text-gray-400" />
                  <UiBaseInput
                    id="email"
                    type="email"
                    v-model="email"
                    placeholder="john@example.com"
                    class="pl-10 transition-all focus:ring-2 focus:ring-blue-100"
                  />
                </div>
              </div>
            </div>

            <div class="flex items-center justify-end border-t border-gray-50 pt-6">
              <UiBaseButton 
                type="submit" 
                :disabled="isUpdating"
                class="min-w-[140px] shadow-sm transition-all active:scale-95"
              >
                <Loader2 v-if="isUpdating" class="mr-2 h-4 w-4 animate-spin" />
                <Check v-else-if="isSuccess" class="mr-2 h-4 w-4 text-white" />
                {{ isUpdating ? "Saving..." : isSuccess ? "Changes Saved" : "Save Changes" }}
              </UiBaseButton>
            </div>
          </form>
        </div>
      </UiBaseCard>
    </section>

    <!-- Workspace Section -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-semibold text-gray-900">Current Workspace</h2>
        <p class="text-sm text-gray-500">Settings and information about your active project environment.</p>
      </div>
      
      <UiBaseCard class="border-gray-100 shadow-sm">
        <div class="p-8">
          <div v-if="activeWorkspace" class="grid gap-8 sm:grid-cols-2">
            <div class="space-y-3">
              <UiBaseLabel class="text-xs font-bold tracking-wide text-gray-500 uppercase">Workspace Identity</UiBaseLabel>
              <div class="flex items-center gap-4 rounded-xl border border-gray-100 bg-gray-50/50 p-4 transition-colors hover:bg-gray-50">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg border border-blue-100 bg-white text-blue-600 shadow-sm">
                  <Building class="h-6 w-6" />
                </div>
                <div>
                  <div class="text-sm font-bold text-gray-900">{{ activeWorkspace.name }}</div>
                  <div class="text-xs text-gray-500 italic">Workspace ID: #{{ activeWorkspace.id }}</div>
                </div>
              </div>
            </div>
            
            <div class="space-y-3">
              <UiBaseLabel class="text-xs font-bold tracking-wide text-gray-500 uppercase">Environment Type</UiBaseLabel>
              <div class="flex items-center gap-4 rounded-xl border border-gray-100 bg-gray-50/50 p-4 transition-colors hover:bg-gray-50">
                <div class="flex h-12 w-12 items-center justify-center rounded-lg border border-purple-100 bg-white text-purple-600 shadow-sm">
                  <ExternalLink class="h-6 w-6" />
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-bold text-gray-900 capitalize">{{ activeWorkspace.type.toLowerCase() }}</span>
                    <UiBaseBadge class="bg-blue-100 text-[10px] text-blue-700">Team</UiBaseBadge>
                  </div>
                  <div class="text-xs text-gray-500">Collaborative workspace</div>
                </div>
              </div>
            </div>
          </div>
          
          <div v-else class="rounded-xl border border-dashed border-gray-200 bg-gray-50/50 p-12 text-center">
            <Building class="mx-auto mb-4 h-12 w-12 text-gray-300" />
            <h3 class="font-medium text-gray-900">No Active Workspace</h3>
            <p class="mt-1 text-sm text-gray-500 text-pretty">Select or create a workspace to see details here.</p>
            <UiBaseButton variant="outline" size="sm" class="mt-4" to="/workspaces">Go to Workspaces</UiBaseButton>
          </div>
        </div>
      </UiBaseCard>
    </section>

    <!-- Danger Zone -->
    <section>
      <div class="mb-4">
        <h2 class="text-lg font-semibold text-red-900">Security & Sessions</h2>
        <p class="text-sm text-red-700 opacity-80">Sensitive actions that affect your account access.</p>
      </div>
      
      <div class="overflow-hidden rounded-xl border border-red-100 bg-red-50/20 shadow-sm ring-1 ring-red-100/50">
        <div class="flex flex-col items-center justify-between gap-6 p-8 sm:flex-row">
          <div class="space-y-1 text-center sm:text-left">
            <h3 class="font-bold text-red-900">Sign out of Nexus</h3>
            <p class="text-sm text-red-700">You will be redirected to the login page. All active sessions will remain, but you will need to re-authenticate on this device.</p>
          </div>
          <UiBaseButton
            variant="destructive"
            class="min-w-[140px] bg-red-600 shadow-md shadow-red-200 hover:bg-red-700 active:scale-95"
            @click="logout"
          >
            <LogOut class="mr-2 h-4 w-4" /> Log out
          </UiBaseButton>
        </div>
      </div>
    </section>
  </div>
</template>
