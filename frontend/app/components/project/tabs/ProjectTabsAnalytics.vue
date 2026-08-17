<script setup lang="ts">
import { Clock, TrendingUp, AlertTriangle, Users } from "lucide-vue-next";

const props = defineProps<{ projectId: string }>();

const workspaceStore = useWorkspaceStore();
const { fetchProjectAnalytics } = useTimeTracking();

const { data: analytics, pending: isLoading } = fetchProjectAnalytics(
  workspaceStore.activeWorkspaceId || "",
  props.projectId
);

const formatMinutes = (minutes: number) => {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
};

const maxLoggedByTask = computed(() => {
  if (!analytics.value?.by_task?.length) return 1;
  return Math.max(
    ...analytics.value.by_task.map((t) =>
      Math.max(t.logged_minutes, t.estimated_minutes || 0)
    ),
    1
  );
});

const maxLoggedByMember = computed(() => {
  if (!analytics.value?.by_member?.length) return 1;
  return Math.max(...analytics.value.by_member.map((m) => m.logged_minutes), 1);
});
</script>

<template>
  <div
    v-if="isLoading"
    class="text-muted-foreground flex items-center justify-center py-20"
  >
    Loading analytics...
  </div>
  <div
    v-else-if="
      !analytics || (!analytics.by_task?.length && !analytics.by_member?.length)
    "
    class="text-muted-foreground flex flex-col items-center justify-center py-20"
  >
    <Clock class="mb-4 h-12 w-12 opacity-30" />
    <p class="text-lg font-medium">No time data yet</p>
    <p class="text-sm">Start tracking time on tasks to see analytics here.</p>
  </div>
  <div v-else class="space-y-8">
    <!-- KPI Cards -->
    <div class="grid grid-cols-2 gap-4 lg:grid-cols-4">
      <div class="bg-card border-border rounded-xl border p-5">
        <div
          class="text-muted-foreground mb-1 flex items-center gap-2 text-xs font-medium tracking-wide uppercase"
        >
          <Clock class="h-3.5 w-3.5" /> Estimated
        </div>
        <p class="text-foreground text-2xl font-bold">
          {{ formatMinutes(analytics.total_estimated_minutes) }}
        </p>
      </div>
      <div class="bg-card border-border rounded-xl border p-5">
        <div
          class="text-muted-foreground mb-1 flex items-center gap-2 text-xs font-medium tracking-wide uppercase"
        >
          <TrendingUp class="h-3.5 w-3.5" /> Logged
        </div>
        <p class="text-foreground text-2xl font-bold">
          {{ formatMinutes(analytics.total_logged_minutes) }}
        </p>
      </div>
      <div class="bg-card border-border rounded-xl border p-5">
        <div
          class="text-muted-foreground mb-1 flex items-center gap-2 text-xs font-medium tracking-wide uppercase"
        >
          Accuracy
        </div>
        <p class="text-foreground text-2xl font-bold">
          {{ analytics.estimate_accuracy_percent }}%
        </p>
      </div>
      <div class="bg-card border-border rounded-xl border p-5">
        <div
          class="text-muted-foreground mb-1 flex items-center gap-2 text-xs font-medium tracking-wide uppercase"
        >
          <AlertTriangle class="h-3.5 w-3.5" /> Over Budget
        </div>
        <p class="text-foreground text-2xl font-bold">
          {{ analytics.over_budget_task_count }}
        </p>
      </div>
    </div>

    <!-- Task Breakdown -->
    <div class="bg-card border-border rounded-xl border p-6">
      <h3 class="text-foreground mb-4 text-sm font-semibold">
        Estimated vs Logged (by Task)
      </h3>
      <div class="space-y-4">
        <div
          v-for="task in analytics.by_task"
          :key="task.task_id"
          class="space-y-1"
        >
          <div class="flex items-center justify-between text-sm">
            <span class="text-foreground font-medium"
              >TASK-{{ task.task_number }}: {{ task.task_title }}</span
            >
            <span
              v-if="task.is_over_budget"
              class="text-xs font-medium text-red-500"
            >
              +{{
                formatMinutes(
                  task.logged_minutes - (task.estimated_minutes || 0)
                )
              }}
              over
            </span>
          </div>
          <div class="space-y-1">
            <div class="flex items-center gap-2">
              <div class="bg-muted h-3 flex-1 overflow-hidden rounded-full">
                <div
                  class="bg-primary/40 h-full rounded-full"
                  :style="{
                    width: `${((task.estimated_minutes || 0) / maxLoggedByTask) * 100}%`,
                  }"
                />
              </div>
              <span class="text-muted-foreground w-16 text-right text-xs"
                >{{ formatMinutes(task.estimated_minutes || 0) }} est</span
              >
            </div>
            <div class="flex items-center gap-2">
              <div class="bg-muted h-3 flex-1 overflow-hidden rounded-full">
                <div
                  class="h-full rounded-full"
                  :class="task.is_over_budget ? 'bg-red-500' : 'bg-primary'"
                  :style="{
                    width: `${(task.logged_minutes / maxLoggedByTask) * 100}%`,
                  }"
                />
              </div>
              <span class="text-muted-foreground w-16 text-right text-xs"
                >{{ formatMinutes(task.logged_minutes) }} log</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Team Breakdown -->
    <div class="bg-card border-border rounded-xl border p-6">
      <h3
        class="text-foreground mb-4 flex items-center gap-2 text-sm font-semibold"
      >
        <Users class="h-4 w-4" /> Hours by Team Member
      </h3>
      <div class="space-y-3">
        <div
          v-for="member in analytics.by_member"
          :key="member.user_id"
          class="flex items-center gap-3"
        >
          <UiBaseAvatar
            :fallback="(member.full_name || '?')[0]?.toUpperCase()"
            class-name="h-7 w-7 text-[10px] border border-border"
          />
          <div class="min-w-0 flex-1">
            <div class="flex items-center justify-between text-sm">
              <span class="text-foreground font-medium">{{
                member.full_name
              }}</span>
              <span class="text-primary text-sm font-semibold">{{
                formatMinutes(member.logged_minutes)
              }}</span>
            </div>
            <div class="bg-muted mt-1 h-2 w-full overflow-hidden rounded-full">
              <div
                class="bg-primary h-full rounded-full transition-all"
                :style="{
                  width: `${(member.logged_minutes / maxLoggedByMember) * 100}%`,
                }"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
