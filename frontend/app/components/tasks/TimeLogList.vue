<script setup lang="ts">
import { Timer, Trash2, Clock, Calendar, PenTool } from "lucide-vue-next";
import { formatDistanceToNow, format } from "date-fns";
import type { TimeEntry } from "~/types";

const props = defineProps<{
  entries: TimeEntry[];
  onRefresh: () => void;
}>();

const { deleteTimeEntry } = useTimeTracking();
const userStore = useUsersStore();

type FilterType = "all" | "timer" | "manual";
const activeFilter = ref<FilterType>("all");

const formatMinutes = (minutes: number | null) => {
  if (!minutes) return "0m";
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
};

const totalLoggedMinutes = computed(() => {
  return props.entries.reduce(
    (sum, entry) => sum + (entry.duration_minutes || 0),
    0
  );
});

const filteredEntries = computed(() => {
  if (activeFilter.value === "timer") {
    return props.entries.filter((e) => !e.is_manual);
  }
  if (activeFilter.value === "manual") {
    return props.entries.filter((e) => e.is_manual);
  }
  return props.entries;
});

const handleDelete = async (entryId: string) => {
  if (!confirm("Are you sure you want to delete this time entry?")) return;
  try {
    await deleteTimeEntry(entryId);
    props.onRefresh();
  } catch (err) {
    console.error("Failed to delete time entry", err);
  }
};
</script>

<template>
  <div class="space-y-4">
    <!-- Header & Summary Bar -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <h4 class="text-foreground text-sm font-semibold tracking-tight">
          Work Log
        </h4>
        <span
          v-if="entries.length > 0"
          class="bg-muted text-muted-foreground rounded-full px-2 py-0.5 font-mono text-xs font-medium"
        >
          {{ entries.length }}
        </span>
      </div>

      <!-- Cumulative Total Badge -->
      <div
        v-if="entries.length > 0"
        class="border-border bg-card/60 flex items-center gap-1.5 rounded-md border px-2.5 py-1 text-xs shadow-xs"
      >
        <Clock class="text-primary h-3.5 w-3.5" />
        <span class="text-muted-foreground font-medium">Total:</span>
        <span class="text-foreground font-mono font-bold">{{
          formatMinutes(totalLoggedMinutes)
        }}</span>
      </div>
    </div>

    <!-- Filter Pills (Visible when entries > 2) -->
    <div v-if="entries.length > 2" class="flex items-center gap-1.5">
      <button
        v-for="filter in ['all', 'timer', 'manual'] as FilterType[]"
        :key="filter"
        class="rounded-md px-2.5 py-1 text-xs font-medium capitalize transition-all"
        :class="
          activeFilter === filter
            ? 'bg-primary/10 text-primary font-semibold'
            : 'text-muted-foreground hover:bg-muted/50 hover:text-foreground'
        "
        @click="activeFilter = filter"
      >
        {{ filter }}
      </button>
    </div>

    <!-- Empty State -->
    <div
      v-if="entries.length === 0"
      class="border-border/60 bg-muted/5 hover:bg-muted/10 flex flex-col items-center justify-center rounded-xl border border-dashed p-6 text-center transition-all"
    >
      <div
        class="border-border/60 bg-background mb-3 flex h-10 w-10 items-center justify-center rounded-full border shadow-xs"
      >
        <Clock class="text-muted-foreground/60 h-5 w-5" />
      </div>
      <p class="text-foreground text-xs font-medium">No work logged yet</p>
      <p class="text-muted-foreground mt-0.5 text-xs">
        Use the live timer or log manual work hours above.
      </p>
    </div>

    <!-- Timeline Work Log Feed -->
    <div v-else class="relative space-y-3 pl-3">
      <!-- Vertical Timeline Connector Line -->
      <div
        v-if="filteredEntries.length > 1"
        class="bg-border/60 pointer-events-none absolute top-4 bottom-4 left-[19px] w-[1px]"
      />

      <div
        v-for="entry in filteredEntries"
        :key="entry.id"
        class="relative z-10 flex items-start gap-3"
      >
        <!-- Timeline Node Icon -->
        <div class="bg-background mt-1 shrink-0">
          <div
            class="flex h-5 w-5 items-center justify-center rounded-full border shadow-xs"
            :class="
              entry.is_manual
                ? 'border-indigo-500/30 bg-indigo-500/10 text-indigo-500'
                : 'border-emerald-500/30 bg-emerald-500/10 text-emerald-500'
            "
          >
            <PenTool v-if="entry.is_manual" class="h-2.5 w-2.5" />
            <Timer v-else class="h-2.5 w-2.5" />
          </div>
        </div>

        <!-- Entry Content Card -->
        <div
          class="border-border/60 bg-card/60 hover:border-border hover:bg-card group min-w-0 flex-1 overflow-hidden rounded-xl border shadow-xs transition-all duration-200"
        >
          <!-- Header Bar -->
          <div
            class="border-border/40 bg-muted/20 flex items-center justify-between border-b px-3.5 py-2"
          >
            <div class="flex min-w-0 items-center gap-2">
              <UiBaseAvatar
                :fallback="(entry.user_full_name || '?')[0]?.toUpperCase()"
                class-name="h-5 w-5 text-[10px] border border-border"
              />
              <span class="text-foreground truncate text-xs font-medium">
                {{ entry.user_full_name || "Team Member" }}
              </span>

              <!-- Log Type Badge -->
              <span
                class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-medium tracking-tight"
                :class="
                  entry.is_manual
                    ? 'bg-indigo-500/10 text-indigo-500 dark:bg-indigo-500/20'
                    : 'bg-emerald-500/10 text-emerald-600 dark:bg-emerald-500/20 dark:text-emerald-400'
                "
              >
                {{ entry.is_manual ? "Manual" : "Timer" }}
              </span>
            </div>

            <!-- Duration Badge & Actions -->
            <div class="flex items-center gap-2.5">
              <span
                class="bg-primary/10 text-primary border-primary/20 rounded-md border px-2 py-0.5 font-mono text-xs font-bold tracking-tight"
              >
                +{{ formatMinutes(entry.duration_minutes) }}
              </span>

              <button
                v-if="entry.user_id === userStore.userData?.id"
                class="text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded-md p-1 opacity-0 transition-all group-hover:opacity-100"
                title="Delete time entry"
                @click="handleDelete(entry.id)"
              >
                <Trash2 class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>

          <!-- Body / Description Area -->
          <div class="px-3.5 py-2.5">
            <p
              v-if="entry.description"
              class="text-foreground/90 text-xs leading-relaxed whitespace-pre-wrap"
            >
              {{ entry.description }}
            </p>
            <p v-else class="text-muted-foreground/50 text-xs italic">
              No work description provided
            </p>

            <div
              class="text-muted-foreground/60 mt-2 flex items-center gap-1 text-[11px]"
            >
              <Calendar class="h-3 w-3" />
              <span
                >{{
                  formatDistanceToNow(
                    new Date(entry.start_time || entry.created_at)
                  )
                }}
                ago</span
              >
              <span
                >({{
                  format(
                    new Date(entry.start_time || entry.created_at),
                    "MMM d, h:mm a"
                  )
                }})</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
