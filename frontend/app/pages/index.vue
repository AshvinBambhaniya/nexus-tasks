<script setup lang="ts">
import {
  Layout,
  Zap,
  ArrowRight,
  Github,
  Code2,
  Terminal,
} from "lucide-vue-next";
import { GITHUB_REPO_URL } from "~/constants";
import { useUsersStore } from "~/stores/user";

const userStore = useUsersStore();
const user = computed(() => userStore.userData);
</script>

<template>
  <div class="selection:bg-primary/30 relative overflow-hidden">
    <!-- Ambient Background Glows -->
    <div
      class="pointer-events-none absolute inset-0 flex justify-center opacity-30 mix-blend-screen"
      aria-hidden="true"
    >
      <div
        class="bg-primary/20 h-[400px] w-[800px] translate-y-[-50%] rounded-full blur-[120px]"
      />
    </div>

    <!-- Hero Section -->
    <section class="relative pt-32 pb-24 lg:pt-48 lg:pb-32">
      <div
        class="relative z-10 container mx-auto px-4 text-center sm:px-6 lg:px-8"
      >
        <div class="mx-auto max-w-4xl">
          <!-- Release Badge -->
          <div
            class="text-primary mb-8 inline-flex items-center rounded-full border border-white/10 bg-white/5 px-3 py-1 font-mono text-xs tracking-wide backdrop-blur-md transition-colors hover:bg-white/10"
          >
            <span class="relative mr-2 flex h-2 w-2">
              <span
                class="bg-primary absolute inline-flex h-full w-full animate-ping rounded-full opacity-75"
              />
              <span
                class="bg-primary relative inline-flex h-2 w-2 rounded-full"
              />
            </span>
            Nexus Tasks v2.0 is Live
          </div>

          <h1
            class="text-foreground mb-8 text-5xl font-extrabold tracking-tight sm:text-7xl"
          >
            Project management for <br class="hidden sm:block" />
            <span
              class="from-primary bg-gradient-to-r to-purple-400 bg-clip-text text-transparent"
              >developers in flow.</span
            >
          </h1>

          <p
            class="text-muted-foreground mx-auto mb-10 max-w-2xl text-lg sm:text-xl"
          >
            Stop fighting complex workflows. Nexus Tasks is a lightning-fast,
            keyboard-first tool built for speed, marrying personal focus with
            team collaboration.
          </p>

          <div
            class="flex flex-col items-center justify-center gap-4 sm:flex-row"
          >
            <template v-if="user">
              <NuxtLink to="/dashboard">
                <UiBaseButton
                  size="lg"
                  class="h-12 px-8 text-base shadow-[0_0_40px_-10px_rgba(94,106,210,0.5)] transition-all hover:shadow-[0_0_60px_-15px_rgba(94,106,210,0.7)]"
                >
                  Launch Workspace <ArrowRight class="ml-2 h-4 w-4" />
                </UiBaseButton>
              </NuxtLink>
            </template>
            <template v-else>
              <NuxtLink to="/register">
                <UiBaseButton
                  size="lg"
                  class="h-12 px-8 text-base shadow-[0_0_40px_-10px_rgba(94,106,210,0.5)] transition-all hover:shadow-[0_0_60px_-15px_rgba(94,106,210,0.7)]"
                >
                  Start Building Free <ArrowRight class="ml-2 h-4 w-4" />
                </UiBaseButton>
              </NuxtLink>
            </template>
            <NuxtLink :to="GITHUB_REPO_URL" target="_blank">
              <UiBaseButton
                variant="outline"
                size="lg"
                class="h-12 border-white/10 bg-transparent px-8 text-base hover:bg-white/5"
              >
                <Github class="mr-2 h-4 w-4" /> Star on GitHub
              </UiBaseButton>
            </NuxtLink>
          </div>

          <!-- Hero UI Preview Placeholder (Abstract) -->
          <div
            class="relative mx-auto mt-20 max-w-5xl rounded-xl border border-white/10 bg-[#0a0a0a]/80 p-2 shadow-2xl backdrop-blur-xl sm:p-4"
          >
            <div
              class="from-primary/10 absolute inset-0 -z-10 bg-gradient-to-b to-transparent blur-2xl"
            />
            <div
              class="overflow-hidden rounded-lg border border-white/5 bg-[#111]"
            >
              <!-- Mockup Header -->
              <div class="flex items-center border-b border-white/5 px-4 py-3">
                <div class="flex gap-1.5">
                  <div class="h-3 w-3 rounded-full bg-white/20" />
                  <div class="h-3 w-3 rounded-full bg-white/20" />
                  <div class="h-3 w-3 rounded-full bg-white/20" />
                </div>
                <div
                  class="mx-auto flex h-6 w-1/3 items-center justify-center rounded bg-white/5 px-2 font-mono text-xs text-white/30"
                >
                  nexus-tasks.app
                </div>
              </div>
              <!-- Mockup Body -->
              <div class="flex h-[300px] sm:h-[400px]">
                <!-- Sidebar -->
                <div
                  class="hidden w-48 space-y-4 border-r border-white/5 p-4 sm:block"
                >
                  <div class="h-4 w-3/4 rounded bg-white/10" />
                  <div class="h-4 w-1/2 rounded bg-white/5" />
                  <div class="h-4 w-5/6 rounded bg-white/5" />
                  <div class="h-4 w-2/3 rounded bg-white/5" />
                </div>
                <!-- Main Content -->
                <div class="flex-1 space-y-4 p-6">
                  <div class="flex justify-between">
                    <div class="h-6 w-1/4 rounded bg-white/10" />
                    <div class="h-6 w-1/6 rounded bg-white/5" />
                  </div>
                  <div
                    class="flex h-12 w-full items-center rounded border border-white/5 bg-white/5 px-4"
                  >
                    <div class="bg-primary/80 mr-3 h-4 w-4 rounded-full" />
                    <div class="h-4 w-1/3 rounded bg-white/20" />
                  </div>
                  <div
                    class="flex h-12 w-full items-center rounded border border-white/5 bg-white/5 px-4"
                  >
                    <div class="mr-3 h-4 w-4 rounded-full bg-yellow-500/80" />
                    <div class="h-4 w-1/2 rounded bg-white/20" />
                  </div>
                  <div
                    class="flex h-12 w-full items-center rounded border border-white/5 bg-white/5 px-4"
                  >
                    <div class="mr-3 h-4 w-4 rounded-full bg-white/20" />
                    <div class="h-4 w-1/4 rounded bg-white/10" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Bento Box Features Grid -->
    <section id="features" class="relative z-10 py-24 sm:py-32">
      <div class="container mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <div class="mx-auto mb-16 max-w-2xl text-center">
          <h2
            class="text-foreground mb-4 text-3xl font-bold tracking-tight sm:text-4xl"
          >
            Designed for the way you code.
          </h2>
          <p class="text-muted-foreground text-lg">
            Everything you need, nothing you don't.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-4 md:grid-cols-3 md:grid-rows-2">
          <!-- Big Feature (Span 2 cols) -->
          <div
            class="group bg-card relative overflow-hidden rounded-2xl border border-white/10 p-8 shadow-lg transition-all hover:border-white/20 md:col-span-2"
          >
            <div
              class="from-primary/10 absolute inset-0 bg-gradient-to-br to-transparent opacity-0 transition-opacity group-hover:opacity-100"
            />
            <div
              class="text-primary mb-4 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-white/5 ring-1 ring-white/10"
            >
              <Terminal class="h-5 w-5" />
            </div>
            <h3 class="text-foreground mb-2 text-xl font-semibold">
              Keyboard First Workflow
            </h3>
            <p class="text-muted-foreground text-sm">
              Command palette driven interface. Hit Cmd+K to navigate, assign,
              or create tasks without touching your mouse.
            </p>
          </div>

          <!-- Small Feature -->
          <div
            class="group bg-card relative overflow-hidden rounded-2xl border border-white/10 p-8 shadow-lg transition-all hover:border-white/20"
          >
            <div
              class="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent opacity-0 transition-opacity group-hover:opacity-100"
            />
            <div
              class="mb-4 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-white/5 text-purple-400 ring-1 ring-white/10"
            >
              <Zap class="h-5 w-5" />
            </div>
            <h3 class="text-foreground mb-2 text-xl font-semibold">
              Real-time Sync
            </h3>
            <p class="text-muted-foreground text-sm">
              Powered by WebSockets. Changes instantly propagate to all your
              devices without refreshing.
            </p>
          </div>

          <!-- Small Feature -->
          <div
            class="group bg-card relative overflow-hidden rounded-2xl border border-white/10 p-8 shadow-lg transition-all hover:border-white/20"
          >
            <div
              class="absolute inset-0 bg-gradient-to-br from-green-500/10 to-transparent opacity-0 transition-opacity group-hover:opacity-100"
            />
            <div
              class="mb-4 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-white/5 text-green-400 ring-1 ring-white/10"
            >
              <Code2 class="h-5 w-5" />
            </div>
            <h3 class="text-foreground mb-2 text-xl font-semibold">
              Markdown Native
            </h3>
            <p class="text-muted-foreground text-sm">
              Full support for GFM markdown in descriptions and comments.
            </p>
          </div>

          <!-- Big Feature (Span 2 cols) -->
          <div
            class="group bg-card relative overflow-hidden rounded-2xl border border-white/10 p-8 shadow-lg transition-all hover:border-white/20 md:col-span-2"
          >
            <div
              class="absolute inset-0 bg-gradient-to-br from-sky-500/10 to-transparent opacity-0 transition-opacity group-hover:opacity-100"
            />
            <div
              class="mb-4 inline-flex h-10 w-10 items-center justify-center rounded-lg bg-white/5 text-sky-400 ring-1 ring-white/10"
            >
              <Layout class="h-5 w-5" />
            </div>
            <h3 class="text-foreground mb-2 text-xl font-semibold">
              Hybrid Workspaces
            </h3>
            <p class="text-muted-foreground text-sm">
              Keep your personal to-dos isolated from team sprints. Switch
              contexts flawlessly from a single unified inbox.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Bottom CTA -->
    <section
      class="relative overflow-hidden border-t border-white/5 py-24 sm:py-32"
    >
      <div
        class="pointer-events-none absolute inset-0 flex justify-center opacity-30 mix-blend-screen"
        aria-hidden="true"
      >
        <div
          class="bg-primary/20 h-[400px] w-[800px] translate-y-[50%] rounded-full blur-[120px]"
        />
      </div>
      <div class="relative z-10 container mx-auto px-4 text-center">
        <h2
          class="text-foreground mb-6 text-3xl font-bold tracking-tight sm:text-5xl"
        >
          Ready to ship faster?
        </h2>
        <p class="text-muted-foreground mx-auto mb-10 max-w-xl text-lg">
          Join the community of developers building better software with less
          friction.
        </p>
        <NuxtLink to="/register">
          <UiBaseButton
            size="lg"
            class="h-14 px-10 text-lg shadow-[0_0_40px_-10px_rgba(94,106,210,0.5)] transition-all hover:shadow-[0_0_60px_-15px_rgba(94,106,210,0.7)]"
          >
            Start Building Free
          </UiBaseButton>
        </NuxtLink>
      </div>
    </section>
  </div>
</template>
