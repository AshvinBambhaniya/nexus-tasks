<script setup lang="ts">
import { User, Users, Bell, Shield, Palette } from "lucide-vue-next";
import SettingsGeneral from "~/components/settings/SettingsGeneral.vue";
import SettingsMembers from "~/components/settings/SettingsMembers.vue";

definePageMeta({
  layout: "dashboard",
});

const activeTab = ref<"general" | "members">("general");

const tabs = [
  { id: "general", label: "General", icon: User, description: "Profile and workspace info" },
  { id: "members", label: "Members", icon: Users, description: "Manage team access" },
] as const;

// Placeholders for future sections to make sidebar feel full
const comingSoonTabs = [
  { label: "Notifications", icon: Bell },
  { label: "Security", icon: Shield },
  { label: "Appearance", icon: Palette },
];
</script>

<template>
  <div class="mx-auto max-w-6xl animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div class="mb-8 border-b border-gray-100 pb-6">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">Settings</h1>
      <p class="mt-2 text-gray-500">
        Manage your personal preferences and workspace configuration.
      </p>
    </div>

    <div class="flex flex-col gap-10 lg:flex-row">
      <!-- Sidebar Navigation -->
      <aside class="w-full shrink-0 lg:w-72">
        <nav class="space-y-1">
          <div class="mb-2 px-3 text-xs font-semibold tracking-wider text-gray-400 uppercase">
            Account & Workspace
          </div>
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'group flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-all duration-200',
              activeTab === tab.id
                ? 'bg-blue-600 text-white shadow-md shadow-blue-100'
                : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900',
            ]"
          >
            <component
              :is="tab.icon"
              :class="[
                'h-5 w-5 shrink-0 transition-colors',
                activeTab === tab.id ? 'text-white' : 'text-gray-400 group-hover:text-gray-600',
              ]"
            />
            <div class="text-left">
              <div>{{ tab.label }}</div>
              <div v-if="activeTab !== tab.id" class="text-[11px] font-normal text-gray-400">
                {{ tab.description }}
              </div>
            </div>
          </button>

          <div class="mt-8 mb-2 px-3 text-xs font-semibold tracking-wider text-gray-400 uppercase">
            Preferences
          </div>
          <div
            v-for="tab in comingSoonTabs"
            :key="tab.label"
            class="flex w-full cursor-not-allowed items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium text-gray-400 opacity-60"
          >
            <component :is="tab.icon" class="h-5 w-5 shrink-0" />
            <span>{{ tab.label }}</span>
            <span class="ml-auto text-[10px] font-bold tracking-tighter text-gray-300 uppercase italic">
              Soon
            </span>
          </div>
        </nav>
      </aside>

      <!-- Content Area -->
      <main class="min-w-0 flex-1">
        <Transition
          mode="out-in"
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="transform translate-y-2 opacity-0"
          enter-to-class="transform translate-y-0 opacity-100"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="transform translate-y-0 opacity-100"
          leave-to-class="transform translate-y-2 opacity-0"
        >
          <div :key="activeTab">
            <SettingsGeneral v-if="activeTab === 'general'" />
            <SettingsMembers v-else-if="activeTab === 'members'" />
          </div>
        </Transition>
      </main>
    </div>
  </div>
</template>
