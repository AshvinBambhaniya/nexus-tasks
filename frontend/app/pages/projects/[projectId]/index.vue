<script setup lang="ts">
import { Loader2, Sparkles, X, Copy } from "lucide-vue-next";
import VueMarkdown from "vue-markdown-render";

definePageMeta({ layout: "dashboard" });

const route = useRoute();
const projectId = computed(() => route.params.projectId as string);

const { project, isLoading, generateWeeklyReport } = useProject(
  projectId.value
);

type TabType =
  | "tasks"
  | "board"
  | "analytics"
  | "time-logs"
  | "members"
  | "settings";
const router = useRouter();

const tabs: { id: TabType; label: string }[] = [
  { id: "tasks", label: "Tasks" },
  { id: "board", label: "Board" },
  { id: "analytics", label: "Analytics" },
  { id: "time-logs", label: "Time Logs" },
  { id: "members", label: "Members" },
  { id: "settings", label: "Settings" },
];

const activeTab = computed({
  get: () => {
    const tab = route.query.tab as string;
    return tabs.some((t) => t.id === tab) ? (tab as TabType) : "tasks";
  },
  set: (val: TabType) => {
    router.replace({ query: { ...route.query, tab: val } });
  },
});

const isGeneratingReport = ref(false);
const showReportModal = ref(false);
const generatedReport = ref<string | null>(null);

const handleGenerateReport = async () => {
  isGeneratingReport.value = true;
  showReportModal.value = true;
  generatedReport.value = null;
  try {
    const report = await generateWeeklyReport();
    generatedReport.value = report || "No data available.";
  } catch {
    generatedReport.value = "Failed to generate report.";
  } finally {
    isGeneratingReport.value = false;
  }
};

const copyToClipboard = () => {
  if (generatedReport.value) {
    navigator.clipboard.writeText(generatedReport.value);
  }
};
</script>

<template>
  <div v-if="isLoading" class="flex h-full items-center justify-center">
    <Loader2 class="text-muted-foreground/70 h-8 w-8 animate-spin" />
  </div>
  <div
    v-else-if="!project"
    class="text-muted-foreground flex h-full items-center justify-center"
  >
    Project not found.
  </div>
  <div v-else class="flex h-full flex-col space-y-6">
    <!-- Header -->
    <div>
      <div class="mb-1 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <h1 class="text-foreground text-2xl font-bold tracking-tight">
            {{ project.name }}
          </h1>
          <span
            v-if="project.is_archived"
            class="rounded-full border border-yellow-200 bg-yellow-100 px-2 py-0.5 text-xs text-yellow-800"
          >
            Archived
          </span>
        </div>
        <button
          class="flex items-center gap-2 rounded-md bg-gradient-to-r from-purple-600 to-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-md transition-all hover:from-purple-700 hover:to-indigo-700 hover:shadow-lg focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 focus:outline-none"
          @click="handleGenerateReport"
        >
          <Sparkles class="h-4 w-4" />
          Generate Weekly Report
        </button>
      </div>
      <p
        v-if="project.description"
        class="text-muted-foreground max-w-2xl text-sm"
      >
        {{ project.description }}
      </p>

      <!-- Tab Navigation -->
      <div class="border-border mt-6 flex gap-6 border-b">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="border-b-2 pb-3 text-sm font-medium transition-colors"
          :class="
            activeTab === tab.id
              ? 'border-primary text-primary'
              : 'text-muted-foreground hover:border-border hover:text-foreground/80 border-transparent'
          "
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Tab Content -->
    <div class="min-h-0 flex-1">
      <div v-if="activeTab === 'tasks'">
        <ProjectTabsTasks :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'board'">
        <ProjectTabsBoard :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'analytics'">
        <ProjectTabsAnalytics :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'time-logs'">
        <ProjectTabsTimeLogs :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'members'">
        <ProjectTabsMembers :project-id="projectId" />
      </div>
      <div v-else-if="activeTab === 'settings'">
        <ProjectTabsSettings :project="project" />
      </div>
    </div>

    <!-- Weekly Report Modal -->
    <div
      v-if="showReportModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm transition-all"
      @click.self="showReportModal = false"
    >
      <div
        class="bg-background flex w-full max-w-4xl flex-col overflow-hidden rounded-xl border border-white/10 shadow-2xl"
      >
        <div class="flex items-center justify-between border-b p-4">
          <h2
            class="flex items-center gap-2 text-lg font-semibold text-purple-600 dark:text-purple-400"
          >
            <Sparkles class="h-5 w-5" />
            Weekly Sprint Report
          </h2>
          <div class="flex items-center gap-2">
            <button
              v-if="generatedReport && !isGeneratingReport"
              class="text-muted-foreground hover:bg-muted/50 hover:text-foreground rounded-md p-2 transition-colors"
              title="Copy to clipboard"
              @click="copyToClipboard"
            >
              <Copy class="h-4 w-4" />
            </button>
            <button
              class="text-muted-foreground hover:bg-muted/50 hover:text-foreground rounded-md p-2 transition-colors"
              @click="showReportModal = false"
            >
              <X class="h-4 w-4" />
            </button>
          </div>
        </div>

        <div class="max-h-[75vh] overflow-y-auto p-6">
          <div
            v-if="isGeneratingReport"
            class="flex flex-col items-center justify-center py-12"
          >
            <Loader2 class="mb-6 h-12 w-12 animate-spin text-purple-500" />
            <p class="text-muted-foreground text-lg">
              Synthesizing achievements and comments...
            </p>
          </div>
          <div
            v-else-if="generatedReport"
            class="prose prose-sm dark:prose-invert prose-p:leading-relaxed prose-pre:bg-muted/50 prose-headings:text-purple-600 dark:prose-headings:text-purple-400 max-w-none"
          >
            <VueMarkdown :source="generatedReport" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
