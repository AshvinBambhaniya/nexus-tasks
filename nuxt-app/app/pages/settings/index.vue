<script setup lang="ts">
import { useUsersStore } from "~/stores/user";

definePageMeta({
  layout: "dashboard",
});

const userStore = useUsersStore();
const user = computed(() => userStore.userData);
</script>

<template>
  <div class="max-w-4xl space-y-8">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Settings</h1>
      <p class="text-gray-500">
        Manage your account and workspace preferences.
      </p>
    </div>

    <UiBaseCard class="divide-y divide-gray-100">
      <div class="p-6">
        <h2 class="mb-4 text-lg font-medium text-gray-900">Profile</h2>
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div class="space-y-2">
            <UiBaseLabel for="fullName">Full Name</UiBaseLabel>
            <UiBaseInput
              id="fullName"
              :model-value="user?.full_name || ''"
              disabled
            />
            <p class="text-xs text-gray-400 italic">Managed by system</p>
          </div>
          <div class="space-y-2">
            <UiBaseLabel for="email">Email Address</UiBaseLabel>
            <UiBaseInput id="email" :model-value="user?.email || ''" disabled />
            <p class="text-xs text-gray-400 italic">Managed by system</p>
          </div>
        </div>
      </div>

      <div class="p-6">
        <h2 class="mb-4 text-lg font-medium text-gray-900">Appearance</h2>
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-medium text-gray-700">Theme</p>
            <p class="text-xs text-gray-500">
              Customize how Nexus Tasks looks on your device.
            </p>
          </div>
          <UiBaseBadge variant="secondary">System Default</UiBaseBadge>
        </div>
      </div>
    </UiBaseCard>

    <div class="rounded-lg border border-red-100 bg-red-50/50 p-6">
      <h3 class="text-lg font-medium text-red-900">Danger Zone</h3>
      <p class="mt-1 text-sm text-red-600">
        Deleting your account is permanent and cannot be undone.
      </p>
      <div class="mt-4">
        <UiBaseButton variant="destructive" size="sm">
          Delete Account
        </UiBaseButton>
      </div>
    </div>
  </div>
</template>
