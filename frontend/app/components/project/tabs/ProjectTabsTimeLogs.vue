<script setup lang="ts">
import {
  Download,
  FileDown,
  Search,
  Filter,
  Clock,
  Users,
  Timer,
  PenTool,
} from "lucide-vue-next";
import { format, subDays } from "date-fns";
import type { ProjectTimeEntry } from "~/types";

const props = defineProps<{
  projectId: string;
}>();

const { project } = useProject(props.projectId);
const workspaceId = computed(() => project.value?.workspace_id);

const timeRange = ref<"all" | "7d" | "30d">("all");
const memberFilter = ref<string>("all");
const searchQuery = ref("");

const queryParams = computed(() => {
  const params: Record<string, string> = {};
  if (timeRange.value !== "all") {
    params.end_date = new Date().toISOString();
    params.start_date = subDays(
      new Date(),
      timeRange.value === "7d" ? 7 : 30
    ).toISOString();
  }
  if (memberFilter.value !== "all") {
    params.user_id = memberFilter.value;
  }
  return params;
});

const url = computed(() => {
  if (!workspaceId.value) return null;
  return `/api/v2/workspaces/${workspaceId.value}/projects/${props.projectId}/time-entries`;
});

const { data: entries, pending } = useApi<ProjectTimeEntry[]>(
  () => url.value || "",
  {
    key: `project-timesheet-${props.projectId}`,
    query: queryParams,
    immediate: !!workspaceId.value,
  }
);

const filteredEntries = computed(() => {
  let list = entries.value || [];
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
      (e) =>
        e.description?.toLowerCase().includes(q) ||
        e.user_full_name?.toLowerCase().includes(q) ||
        e.task_title?.toLowerCase().includes(q) ||
        `task-${e.task_number}`.includes(q)
    );
  }
  return list;
});

const formatMinutes = (minutes: number | null) => {
  if (!minutes) return "0m";
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
};

// KPIs
const totalMinutes = computed(() =>
  filteredEntries.value.reduce((sum, e) => sum + (e.duration_minutes || 0), 0)
);
const activeContributorsCount = computed(
  () => new Set(filteredEntries.value.map((e) => e.user_id)).size
);
const totalSessions = computed(() => filteredEntries.value.length);

const allMembers = computed(() => {
  const membersMap = new Map<string, string>();
  (entries.value || []).forEach((e) => {
    if (!membersMap.has(e.user_id)) {
      membersMap.set(e.user_id, e.user_full_name || "Unknown");
    }
  });
  return Array.from(membersMap.entries()).map(([id, name]) => ({ id, name }));
});

const handleExportCSV = () => {
  if (!filteredEntries.value.length) return;
  const headers = [
    "Date",
    "User",
    "Task",
    "Duration (mins)",
    "Type",
    "Description",
  ];
  const rows = filteredEntries.value.map((e) => [
    format(new Date(e.start_time || e.created_at), "yyyy-MM-dd HH:mm"),
    e.user_full_name,
    `TASK-${e.task_number} ${e.task_title}`,
    e.duration_minutes || 0,
    e.is_manual ? "Manual" : "Timer",
    `"${(e.description || "").replace(/"/g, '""')}"`,
  ]);

  const csv = [headers.join(","), ...rows.map((r) => r.join(","))].join("\n");
  const blob = new Blob([csv], { type: "text/csv" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `timesheet-${props.projectId}-${format(new Date(), "yyyy-MM-dd")}.csv`;
  a.click();
};

const handleExportJSON = () => {
  if (!filteredEntries.value.length) return;
  const data = JSON.stringify(filteredEntries.value, null, 2);
  const blob = new Blob([data], { type: "application/json" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = `timesheet-${props.projectId}-${format(new Date(), "yyyy-MM-dd")}.json`;
  a.click();
};
</script>

<template>
  <div class="space-y-6">
    <!-- Toolbar -->
    <div
      class="border-border/60 bg-card/60 flex flex-col gap-4 rounded-xl border p-4 shadow-sm md:flex-row md:items-center md:justify-between"
    >
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <div class="border-border bg-muted/30 flex rounded-md border p-1">
          <button
            v-for="range in [
              { id: 'all', label: 'All Time' },
              { id: '30d', label: 'Last 30 Days' },
              { id: '7d', label: 'Last 7 Days' },
            ] as const"
            :key="range.id"
            class="rounded px-3 py-1.5 text-xs font-medium transition-colors"
            :class="
              timeRange === range.id
                ? 'bg-background text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="timeRange = range.id"
          >
            {{ range.label }}
          </button>
        </div>

        <div class="flex items-center gap-2">
          <Filter class="text-muted-foreground h-4 w-4" />
          <select
            v-model="memberFilter"
            class="border-border bg-background focus:border-primary focus:ring-primary rounded-md border py-1.5 pr-8 pl-3 text-sm transition-colors outline-none focus:ring-1"
          >
            <option value="all">All Members</option>
            <option
              v-for="member in allMembers"
              :key="member.id"
              :value="member.id"
            >
              {{ member.name }}
            </option>
          </select>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <div class="relative">
          <Search
            class="text-muted-foreground absolute top-1/2 left-2.5 h-4 w-4 -translate-y-1/2"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search tasks or logs..."
            class="border-border bg-background focus:border-primary focus:ring-primary w-full rounded-md border py-1.5 pr-3 pl-9 text-sm transition-colors outline-none focus:ring-1 sm:w-64"
          />
        </div>
        <button
          class="border-border/60 hover:bg-muted flex items-center gap-2 rounded-md border px-3 py-1.5 text-sm font-medium transition-colors"
          @click="handleExportCSV"
        >
          <FileDown class="h-4 w-4" />
          CSV
        </button>
        <button
          class="border-border/60 hover:bg-muted flex items-center gap-2 rounded-md border px-3 py-1.5 text-sm font-medium transition-colors"
          @click="handleExportJSON"
        >
          <Download class="h-4 w-4" />
          JSON
        </button>
      </div>
    </div>

    <!-- KPI Cards -->
    <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <div
        class="border-border/60 bg-card/60 flex items-center gap-4 rounded-xl border p-4 shadow-sm"
      >
        <div
          class="border-border bg-muted flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border"
        >
          <Clock class="text-primary h-5 w-5" />
        </div>
        <div>
          <p
            class="text-muted-foreground text-xs font-medium tracking-wider uppercase"
          >
            Total Hours Logged
          </p>
          <p class="text-foreground font-mono text-xl font-bold">
            {{ formatMinutes(totalMinutes) }}
          </p>
        </div>
      </div>

      <div
        class="border-border/60 bg-card/60 flex items-center gap-4 rounded-xl border p-4 shadow-sm"
      >
        <div
          class="border-border bg-muted flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border"
        >
          <Users class="text-primary h-5 w-5" />
        </div>
        <div>
          <p
            class="text-muted-foreground text-xs font-medium tracking-wider uppercase"
          >
            Active Contributors
          </p>
          <p class="text-foreground font-mono text-xl font-bold">
            {{ activeContributorsCount }}
          </p>
        </div>
      </div>

      <div
        class="border-border/60 bg-card/60 flex items-center gap-4 rounded-xl border p-4 shadow-sm"
      >
        <div
          class="border-border bg-muted flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border"
        >
          <Timer class="text-primary h-5 w-5" />
        </div>
        <div>
          <p
            class="text-muted-foreground text-xs font-medium tracking-wider uppercase"
          >
            Total Sessions
          </p>
          <p class="text-foreground font-mono text-xl font-bold">
            {{ totalSessions }}
          </p>
        </div>
      </div>
    </div>

    <!-- Master Timesheet Table -->
    <div
      class="border-border/60 bg-card/60 overflow-hidden rounded-xl border shadow-sm"
    >
      <div v-if="pending" class="flex items-center justify-center p-8">
        <Loader2 class="text-muted-foreground h-6 w-6 animate-spin" />
      </div>
      <div
        v-else-if="filteredEntries.length === 0"
        class="flex flex-col items-center justify-center p-12 text-center"
      >
        <Clock class="text-muted-foreground/50 mb-3 h-8 w-8" />
        <p class="text-foreground text-sm font-medium">No time entries found</p>
        <p class="text-muted-foreground mt-1 text-xs">
          Try adjusting your filters or date range.
        </p>
      </div>
      <table v-else class="w-full text-left text-sm">
        <thead
          class="border-border/60 bg-muted/30 text-muted-foreground border-b text-xs font-medium tracking-wider uppercase"
        >
          <tr>
            <th class="px-4 py-3">Date</th>
            <th class="px-4 py-3">Contributor</th>
            <th class="px-4 py-3">Task</th>
            <th class="px-4 py-3">Duration</th>
            <th class="px-4 py-3">Type</th>
            <th class="px-4 py-3">Description</th>
          </tr>
        </thead>
        <tbody class="divide-border/40 divide-y">
          <tr
            v-for="entry in filteredEntries"
            :key="entry.id"
            class="hover:bg-muted/10 transition-colors"
          >
            <td class="text-muted-foreground px-4 py-3 whitespace-nowrap">
              {{
                format(
                  new Date(entry.start_time || entry.created_at),
                  "MMM d, yyyy"
                )
              }}
              <span class="block text-xs opacity-70">{{
                format(new Date(entry.start_time || entry.created_at), "h:mm a")
              }}</span>
            </td>
            <td class="px-4 py-3 whitespace-nowrap">
              <div class="flex items-center gap-2">
                <UiBaseAvatar
                  :fallback="(entry.user_full_name || '?')[0]?.toUpperCase()"
                  class-name="h-6 w-6 text-[10px] border border-border"
                />
                <span class="text-foreground text-xs font-medium">{{
                  entry.user_full_name || "Unknown"
                }}</span>
              </div>
            </td>
            <td class="px-4 py-3 whitespace-nowrap">
              <NuxtLink
                :to="`/projects/${projectId}/tasks/${entry.task_id}`"
                class="border-border/60 bg-muted/30 hover:bg-muted/50 inline-flex items-center gap-1.5 rounded-md border px-2 py-1 text-xs font-medium transition-colors"
              >
                <span class="text-muted-foreground font-mono"
                  >TASK-{{ entry.task_number }}</span
                >
                <span
                  class="text-foreground max-w-[150px] truncate"
                  :title="entry.task_title"
                  >{{ entry.task_title }}</span
                >
              </NuxtLink>
            </td>
            <td class="px-4 py-3 whitespace-nowrap">
              <span
                class="bg-primary/10 text-primary border-primary/20 rounded-md border px-2 py-0.5 font-mono text-xs font-bold"
              >
                +{{ formatMinutes(entry.duration_minutes) }}
              </span>
            </td>
            <td class="px-4 py-3 whitespace-nowrap">
              <span
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-medium tracking-tight"
                :class="
                  entry.is_manual
                    ? 'bg-indigo-500/10 text-indigo-500'
                    : 'bg-emerald-500/10 text-emerald-600'
                "
              >
                <PenTool v-if="entry.is_manual" class="h-2.5 w-2.5" />
                <Timer v-else class="h-2.5 w-2.5" />
                {{ entry.is_manual ? "Manual" : "Timer" }}
              </span>
            </td>
            <td class="px-4 py-3">
              <p
                class="text-muted-foreground line-clamp-2 text-xs"
                :title="entry.description"
              >
                {{ entry.description || "-" }}
              </p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
