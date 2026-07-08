<script setup lang="ts">
import { User, Users, Bell, Shield, Palette } from "lucide-vue-next";
import SettingsGeneral from "~/components/settings/SettingsGeneral.vue";
import SettingsMembers from "~/components/settings/SettingsMembers.vue";

definePageMeta({
  layout: "dashboard",
});

const activeTab = ref<"general" | "members">("general");

const tabs = [
  {
    id: "general",
    label: "General",
    icon: User,
    description: "Profile and workspace info",
  },
  {
    id: "members",
    label: "Members",
    icon: Users,
    description: "Manage team access",
  },
] as const;

// Placeholders for future sections
const comingSoonTabs = [
  { label: "Notifications", icon: Bell },
  { label: "Security", icon: Shield },
  { label: "Appearance", icon: Palette },
];
</script>

<template>
  <div class="mx-auto max-w-5xl pb-12">
    <!-- Header -->
    <div class="mb-10">
      <h1 class="text-foreground text-2xl font-bold tracking-tight">
        Settings
      </h1>
      <p class="text-muted-foreground mt-1 text-sm">
        Manage your personal preferences and workspace configuration.
      </p>
    </div>

    <div class="flex flex-col gap-12 md:flex-row">
      <!-- Sidebar Navigation -->
      <aside class="w-full shrink-0 md:w-64">
        <nav class="space-y-6">
          <div>
            <div
              class="text-muted-foreground mb-3 px-3 text-xs font-semibold tracking-wider uppercase"
            >
              Account & Workspace
            </div>
            <div class="space-y-1">
              <button
                v-for="tab in tabs"
                :key="tab.id"
                :class="[
                  'group flex w-full items-center gap-3 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                  activeTab === tab.id
                    ? 'bg-muted text-foreground'
                    : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground',
                ]"
                @click="activeTab = tab.id"
              >
                <component
                  :is="tab.icon"
                  :class="[
                    'h-4 w-4 shrink-0',
                    activeTab === tab.id
                      ? 'text-foreground'
                      : 'text-muted-foreground group-hover:text-foreground',
                  ]"
                />
                {{ tab.label }}
              </button>
            </div>
          </div>

          <div>
            <div
              class="text-muted-foreground mb-3 px-3 text-xs font-semibold tracking-wider uppercase"
            >
              Preferences
            </div>
            <div class="space-y-1">
              <div
                v-for="tab in comingSoonTabs"
                :key="tab.label"
                class="text-muted-foreground flex w-full cursor-not-allowed items-center gap-3 rounded-md px-3 py-2 text-sm font-medium opacity-60"
              >
                <component :is="tab.icon" class="h-4 w-4 shrink-0" />
                <span>{{ tab.label }}</span>
              </div>
            </div>
          </div>
        </nav>
      </aside>

      <!-- Content Area -->
      <main class="min-w-0 flex-1">
        <Transition
          mode="out-in"
          enter-active-class="transition duration-200 ease-out"
          enter-from-class="opacity-0 translate-y-1"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition duration-150 ease-in"
          leave-from-class="opacity-100 translate-y-0"
          leave-to-class="opacity-0 translate-y-1"
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
