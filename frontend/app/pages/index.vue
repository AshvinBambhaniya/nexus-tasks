<script setup lang="ts">
import {
  Layout,
  Users,
  Zap,
  Shield,
  ArrowRight,
  Github,
  Code2,
  Terminal,
  Cpu,
  GitBranch,
} from "lucide-vue-next";
import { GITHUB_REPO_URL } from "~/constants";
import { useUsersStore } from "~/stores/user";

const userStore = useUsersStore();
const user = computed(() => userStore.userData);
</script>

<template>
  <div>
    <!-- Hero Section -->
    <section class="relative overflow-hidden pt-20 pb-28 lg:pt-32 lg:pb-40">
      <div
        class="relative z-10 container mx-auto px-4 text-center sm:px-6 lg:px-8"
      >
        <div class="mx-auto max-w-4xl">
          <div
            class="border-primary/10 bg-primary/5 text-primary mb-8 inline-flex items-center rounded-full border px-3 py-1 text-sm font-medium backdrop-blur-sm"
          >
            <span
              class="bg-primary mr-2 flex h-2 w-2 animate-pulse rounded-full"
            />
            v2.0 is Live • Open Source
          </div>
          <h1
            class="text-foreground mb-6 text-4xl font-extrabold tracking-tight sm:text-6xl"
          >
            The Task Manager for <br class="hidden sm:block" />
            <span class="text-primary">Developers in Flow</span>
          </h1>
          <p
            class="text-muted-foreground mx-auto mb-10 max-w-2xl text-lg leading-relaxed"
          >
            Stop fighting complex workflows. Nexus Tasks bridges your
            <strong class="text-foreground font-semibold"
              >personal focus</strong
            >
            with
            <strong class="text-foreground font-semibold"
              >team collaboration</strong
            >
            in a clean, high-performance interface built for speed.
          </p>

          <div
            class="mb-16 flex flex-col items-center justify-center gap-4 sm:flex-row"
          >
            <template v-if="user">
              <NuxtLink to="/dashboard">
                <UiBaseButton
                  size="lg"
                  class="shadow-primary/20 h-12 px-8 text-base shadow-lg"
                >
                  Launch Workspace <ArrowRight class="ml-2 h-4 w-4" />
                </UiBaseButton>
              </NuxtLink>
            </template>
            <template v-else>
              <NuxtLink to="/register">
                <UiBaseButton
                  size="lg"
                  class="shadow-primary/20 h-12 px-8 text-base shadow-lg"
                >
                  Start Building Free <ArrowRight class="ml-2 h-4 w-4" />
                </UiBaseButton>
              </NuxtLink>
            </template>
            <NuxtLink :to="GITHUB_REPO_URL" target="_blank">
              <UiBaseButton
                variant="outline"
                size="lg"
                class="h-12 px-8 text-base"
              >
                <Github class="mr-2 h-4 w-4" /> Star on GitHub
              </UiBaseButton>
            </NuxtLink>
          </div>

          <!-- Tech Stack Pills -->
          <div
            id="tech-stack"
            class="flex flex-wrap justify-center gap-3 opacity-80"
          >
            <div
              class="border-border bg-card text-muted-foreground flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium shadow-sm"
            >
              <Code2 class="text-primary h-4 w-4" /> Nuxt 4
            </div>
            <div
              class="border-border bg-card text-muted-foreground flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium shadow-sm"
            >
              <Terminal class="text-primary h-4 w-4" /> Golang Fiber
            </div>
            <div
              class="border-border bg-card text-muted-foreground flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium shadow-sm"
            >
              <GitBranch class="text-primary h-4 w-4" /> TypeScript
            </div>
            <div
              class="border-border bg-card text-muted-foreground flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium shadow-sm"
            >
              <Cpu class="text-primary h-4 w-4" /> Real-time WebSockets
            </div>
          </div>
        </div>
      </div>

      <!-- Decorative Background Elements -->
      <div
        class="pointer-events-none absolute top-0 left-1/2 -z-10 h-full w-full -translate-x-1/2 overflow-hidden"
      >
        <div
          class="animate-blob bg-primary/10 dark:bg-primary/5 absolute top-[-10%] left-[10%] h-[500px] w-[500px] rounded-full opacity-60 mix-blend-multiply blur-[100px] filter dark:mix-blend-lighten"
        />
        <div
          class="animate-blob animation-delay-2000 absolute top-[-10%] right-[10%] h-[500px] w-[500px] rounded-full bg-indigo-500/10 opacity-60 mix-blend-multiply blur-[100px] filter dark:bg-indigo-500/5 dark:mix-blend-lighten"
        />
        <div
          class="animate-blob animation-delay-4000 absolute bottom-[-20%] left-[30%] h-[600px] w-[600px] rounded-full bg-sky-400/10 opacity-60 mix-blend-multiply blur-[100px] filter dark:bg-sky-400/5 dark:mix-blend-lighten"
        />
      </div>
    </section>

    <!-- Features Grid -->
    <section id="features" class="border-border bg-muted/30 border-t py-24">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div class="mx-auto mb-16 max-w-3xl text-center">
          <h2
            class="text-foreground mb-4 text-3xl font-bold tracking-tight sm:text-4xl"
          >
            Designed for the way you code
          </h2>
          <p class="text-muted-foreground text-lg">
            Nexus Tasks cuts through the enterprise bloat to give you exactly
            what you need: clarity, speed, and structure.
          </p>
        </div>

        <div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <Layout class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              Hybrid Workspaces
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Keep your personal side-projects and team sprints in one place.
              Switch contexts instantly without logging out.
            </p>
          </div>

          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <Code2 class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              Markdown Native
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Write task descriptions and comments using standard Markdown. Code
              blocks, lists, and formatting just work.
            </p>
          </div>

          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <Zap class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              Real-time Sync
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Collaborate without refreshing. See comments, status changes, and
              new tasks appear instantly via WebSockets.
            </p>
          </div>

          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <GitBranch class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              GitHub-style Workflow
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Manage tasks with a familiar flow: Open, In Progress, Done. Track
              history and discussions linearly.
            </p>
          </div>

          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <Users class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              Team Alignment
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Organize members into Teams. Assign projects to specific groups
              and manage access with minimal friction.
            </p>
          </div>

          <div
            class="group border-border bg-card hover:border-primary/30 flex flex-col items-start rounded-2xl border p-8 text-left shadow-sm transition-all duration-200 hover:shadow-md"
          >
            <div
              class="bg-primary/10 text-primary group-hover:bg-primary group-hover:text-primary-foreground mb-6 flex h-12 w-12 items-center justify-center rounded-xl transition-colors"
            >
              <Shield class="h-6 w-6" />
            </div>
            <h3 class="text-foreground mb-3 text-xl font-bold">
              Production Ready
            </h3>
            <p class="text-muted-foreground leading-relaxed">
              Built with a high-performance Golang backend and Type-Safe Nuxt
              frontend. Secure, scalable, and self-hostable.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Call to Action -->
    <section class="bg-background py-24">
      <div class="container mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="bg-foreground relative overflow-hidden rounded-3xl px-6 py-16 text-center shadow-2xl sm:px-12 sm:py-20"
        >
          <div class="relative z-10 mx-auto max-w-2xl">
            <h2
              class="text-background mb-6 text-3xl font-bold tracking-tight sm:text-4xl"
            >
              Ready to regain your focus?
            </h2>
            <p class="text-muted-foreground mb-10 text-lg opacity-80">
              Join the community of developers building better software with
              less friction. Open source and free to get started.
            </p>
            <NuxtLink to="/register">
              <UiBaseButton
                variant="white"
                size="lg"
                class="h-14 px-8 text-base font-semibold"
              >
                Get Started for Free
              </UiBaseButton>
            </NuxtLink>
          </div>

          <!-- Abstract shapes -->
          <svg
            viewBox="0 0 1024 1024"
            class="absolute top-1/2 left-1/2 -z-10 h-[64rem] w-[64rem] -translate-x-1/2 [mask-image:radial-gradient(closest-side,white,transparent)]"
            aria-hidden="true"
          >
            <circle
              cx="512"
              cy="512"
              r="512"
              fill="url(#gradient)"
              fillOpacity="0.7"
            />
            <defs>
              <radialGradient id="gradient">
                <stop stop-color="var(--color-primary)" />
                <stop offset="1" stop-color="#80CAFF" />
              </radialGradient>
            </defs>
          </svg>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.animate-blob {
  animation: blob 7s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
.animation-delay-4000 {
  animation-delay: 4s;
}

@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}
</style>
