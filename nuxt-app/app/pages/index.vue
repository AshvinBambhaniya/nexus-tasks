<script setup lang="ts">
import { 
  CheckSquare, 
  Layout, 
  Users, 
  Zap, 
  Shield, 
  ArrowRight, 
  Github, 
  Code2, 
  Terminal, 
  Cpu, 
  GitBranch 
} from 'lucide-vue-next'

const { user, isLoading } = useUser()

const year = new Date().getFullYear()
</script>

<template>
  <div class="flex min-h-screen flex-col bg-white selection:bg-blue-100">
    <!-- Navigation -->
    <header class="sticky top-0 z-50 w-full border-b border-gray-100 bg-white/80 backdrop-blur-md">
      <div class="container mx-auto flex h-16 items-center justify-between px-4 sm:px-6 lg:px-8">
        <div class="flex items-center gap-2">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 text-white shadow-md shadow-blue-200">
            <CheckSquare class="h-5 w-5" />
          </div>
          <span class="text-xl font-bold tracking-tight text-gray-900">
            Nexus Tasks
          </span>
        </div>

        <nav class="hidden items-center gap-8 text-sm font-medium text-gray-600 md:flex">
          <NuxtLink to="#features" class="transition-colors hover:text-blue-600">
            Features
          </NuxtLink>
          <NuxtLink to="#tech-stack" class="transition-colors hover:text-blue-600">
            Tech Stack
          </NuxtLink>
          <NuxtLink 
            to="https://github.com/AshvinBambhaniya/nexus-tasks" 
            target="_blank" 
            class="transition-colors hover:text-blue-600"
          >
            GitHub
          </NuxtLink>
        </nav>

        <div class="flex items-center gap-4">
          <template v-if="isLoading">
            <UiBaseButton disabled variant="ghost" size="sm">
              Loading...
            </UiBaseButton>
          </template>
          <template v-else-if="user">
            <NuxtLink to="/dashboard">
              <UiBaseButton class="font-semibold shadow-md shadow-blue-100">
                Go to Dashboard
              </UiBaseButton>
            </NuxtLink>
          </template>
          <template v-else>
            <NuxtLink to="/login" class="hidden sm:block">
              <UiBaseButton variant="ghost" class="font-medium text-gray-600 hover:text-gray-900">
                Log in
              </UiBaseButton>
            </NuxtLink>
            <NuxtLink to="/register">
              <UiBaseButton class="font-semibold shadow-md shadow-blue-100">
                Get Started
              </UiBaseButton>
            </NuxtLink>
          </template>
        </div>
      </div>
    </header>

    <main class="flex-1">
      <!-- Hero Section -->
      <section class="relative overflow-hidden pt-20 pb-28 lg:pt-32 lg:pb-40">
        <div class="relative z-10 container mx-auto px-4 text-center sm:px-6 lg:px-8">
          <div class="mx-auto max-w-4xl">
            <div class="mb-8 inline-flex items-center rounded-full border border-blue-100 bg-blue-50/50 px-3 py-1 text-sm font-medium text-blue-600 backdrop-blur-sm">
              <span class="mr-2 flex h-2 w-2 animate-pulse rounded-full bg-blue-600"></span>
              v1.0 is Live • Open Source
            </div>
            <h1 class="mb-6 text-4xl font-extrabold tracking-tight text-gray-900 sm:text-6xl">
              The Task Manager for <br class="hidden sm:block" />
              <span class="text-blue-600">Developers in Flow</span>
            </h1>
            <p class="mx-auto mb-10 max-w-2xl text-lg leading-relaxed text-gray-600">
              Stop fighting complex workflows. Nexus Tasks bridges your 
              <strong class="font-semibold text-gray-900">personal focus</strong> 
              with 
              <strong class="font-semibold text-gray-900">team collaboration</strong> 
              in a clean, high-performance interface built for speed.
            </p>

            <div class="mb-16 flex flex-col items-center justify-center gap-4 sm:flex-row">
              <template v-if="user">
                <NuxtLink to="/dashboard">
                  <UiBaseButton size="lg" class="h-12 px-8 text-base">
                    Launch Workspace <ArrowRight class="ml-2 h-4 w-4" />
                  </UiBaseButton>
                </NuxtLink>
              </template>
              <template v-else>
                <NuxtLink to="/register">
                  <UiBaseButton size="lg" class="h-12 px-8 text-base">
                    Start Building Free <ArrowRight class="ml-2 h-4 w-4" />
                  </UiBaseButton>
                </NuxtLink>
              </template>
              <NuxtLink to="https://github.com/AshvinBambhaniya/nexus-tasks" target="_blank">
                <UiBaseButton variant="outline" size="lg" class="h-12 border-slate-200 bg-white px-8 text-base text-slate-700 hover:text-slate-900">
                  <Github class="mr-2 h-4 w-4" /> Star on GitHub
                </UiBaseButton>
              </NuxtLink>
            </div>

            <!-- Tech Stack Pills -->
            <div class="flex flex-wrap justify-center gap-3 opacity-80" id="tech-stack">
              <div class="flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-sm">
                <Code2 class="h-4 w-4 text-gray-400" /> Nuxt 4
              </div>
              <div class="flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-sm">
                <Terminal class="h-4 w-4 text-gray-400" /> Golang
              </div>
              <div class="flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-sm">
                <GitBranch class="h-4 w-4 text-gray-400" /> TypeScript
              </div>
              <div class="flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 shadow-sm">
                <Cpu class="h-4 w-4 text-gray-400" /> Real-time WebSockets
              </div>
            </div>
          </div>
        </div>

        <!-- Decorative Background Elements -->
        <div class="pointer-events-none absolute top-0 left-1/2 -z-10 h-full w-full -translate-x-1/2 overflow-hidden">
          <div class="animate-blob absolute top-[-10%] left-[10%] h-[500px] w-[500px] rounded-full bg-blue-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
          <div class="animate-blob animation-delay-2000 absolute top-[-10%] right-[10%] h-[500px] w-[500px] rounded-full bg-indigo-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
          <div class="animate-blob animation-delay-4000 absolute bottom-[-20%] left-[30%] h-[600px] w-[600px] rounded-full bg-sky-100/50 opacity-60 mix-blend-multiply blur-[100px] filter"></div>
        </div>
      </section>

      <!-- Features Grid -->
      <section id="features" class="border-t border-gray-100 bg-gray-50/50 py-24">
        <div class="container mx-auto px-4 sm:px-6 lg:px-8">
          <div class="mx-auto mb-16 max-w-3xl text-center">
            <h2 class="mb-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
              Designed for the way you code
            </h2>
            <p class="text-lg text-gray-600">
              Nexus Tasks cuts through the enterprise bloat to give you exactly what you need: clarity, speed, and structure.
            </p>
          </div>

          <div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-3">
            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <Layout class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">Hybrid Workspaces</h3>
              <p class="leading-relaxed text-gray-500">Keep your personal side-projects and team sprints in one place. Switch contexts instantly without logging out.</p>
            </div>
            
            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <Code2 class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">Markdown Native</h3>
              <p class="leading-relaxed text-gray-500">Write task descriptions and comments using standard Markdown. Code blocks, lists, and formatting just work.</p>
            </div>

            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <Zap class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">Real-time Sync</h3>
              <p class="leading-relaxed text-gray-500">Collaborate without refreshing. See comments, status changes, and new tasks appear instantly via WebSockets.</p>
            </div>

            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <GitBranch class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">GitHub-style Workflow</h3>
              <p class="leading-relaxed text-gray-500">Manage tasks with a familiar flow: Open, In Progress, Done. Track history and discussions linearly.</p>
            </div>

            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <Users class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">Team Alignment</h3>
              <p class="leading-relaxed text-gray-500">Organize members into Teams. Assign projects to specific groups and manage access with minimal friction.</p>
            </div>

            <div class="group flex flex-col items-start rounded-2xl border border-gray-100 bg-white p-8 text-left shadow-sm transition-all duration-200 hover:border-blue-100 hover:shadow-md">
              <div class="mb-6 flex h-12 w-12 items-center justify-center rounded-xl bg-blue-50 text-blue-600 transition-colors group-hover:bg-blue-600 group-hover:text-white">
                <Shield class="h-6 w-6" />
              </div>
              <h3 class="mb-3 text-xl font-bold text-gray-900">Production Ready</h3>
              <p class="leading-relaxed text-gray-500">Built with a high-performance Golang backend and Type-Safe Nuxt frontend. Secure, scalable, and self-hostable.</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Call to Action -->
      <section class="bg-white py-24">
        <div class="container mx-auto px-4 sm:px-6 lg:px-8">
          <div class="relative overflow-hidden rounded-3xl bg-gray-900 px-6 py-16 text-center shadow-2xl sm:px-12 sm:py-20">
            <div class="relative z-10 mx-auto max-w-2xl">
              <h2 class="mb-6 text-3xl font-bold tracking-tight text-white sm:text-4xl">
                Ready to regain your focus?
              </h2>
              <p class="mb-10 text-lg text-gray-300">
                Join the community of developers building better software with less friction. Open source and free to get started.
              </p>
              <NuxtLink to="/register">
                <UiBaseButton variant="white" size="lg" class="h-14 px-8 text-base font-semibold">
                  Get Started for Free
                </UiBaseButton>
              </NuxtLink>
            </div>

            <!-- Abstract shapes -->
            <svg viewBox="0 0 1024 1024" class="absolute top-1/2 left-1/2 -z-10 h-[64rem] w-[64rem] -translate-x-1/2 [mask-image:radial-gradient(closest-side,white,transparent)]" aria-hidden="true">
              <circle cx="512" cy="512" r="512" fill="url(#gradient)" fill-opacity="0.7" />
              <defs>
                <radialGradient id="gradient">
                  <stop stop-color="#4F46E5" />
                  <stop offset="1" stop-color="#80CAFF" />
                </radialGradient>
              </defs>
            </svg>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-100 bg-white py-12">
      <div class="container mx-auto flex flex-col items-center justify-between gap-6 px-4 sm:px-6 md:flex-row lg:px-8">
        <div class="flex items-center gap-2">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-gray-900 text-white">
            <CheckSquare class="h-4 w-4" />
          </div>
          <div>
            <span class="block text-sm font-bold text-gray-900">Nexus Tasks</span>
            <span class="text-xs text-gray-500">Built for Builders</span>
          </div>
        </div>
        <p class="text-sm text-gray-500">
          © {{ year }} Nexus Tasks. Distributed under MIT License.
        </p>
        <div class="flex gap-6">
          <NuxtLink to="https://github.com/AshvinBambhaniya/nexus-tasks" target="_blank" class="flex items-center gap-2 text-sm text-gray-400 transition-colors hover:text-gray-900">
            <Github class="h-5 w-5" />
            <span>Source Code</span>
          </NuxtLink>
        </div>
      </div>
    </footer>
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
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
</style>
