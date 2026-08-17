<script setup lang="ts">
import {
  Timer,
  Play,
  Square,
  Plus,
  Pencil,
  AlertCircle,
  Search,
} from "lucide-vue-next";
import type { Task, TimeEntry } from "~/types";

const props = defineProps<{
  task: Task;
  entries: TimeEntry[];
  totalLoggedMinutes: number;
  onRefresh: () => void;
}>();

const timerStore = useTimerStore();
const { startTimer, stopTimer } = useTimeTracking();
const { updateTask } = useTasks(props.task.project_id);

const showLogDialog = ref(false);
const showStopDialog = ref(false);
const showEstimateDialog = ref(false);

const stopDescription = ref("");
const isStopping = ref(false);

const estimateHours = ref(0);
const estimateMins = ref(0);
const isSavingEstimate = ref(false);

const searchQuery = ref("");
const activeFilter = ref<"all" | "timer" | "manual">("all");

const isTimerForThisTask = computed(
  () => timerStore.activeTimer?.task_id === props.task.id
);

const estimateMinutes = computed(() => props.task.estimated_minutes || 0);

watch(
  estimateMinutes,
  (newVal) => {
    estimateHours.value = Math.floor(newVal / 60);
    estimateMins.value = newVal % 60;
  },
  { immediate: true }
);

const openEstimateModal = () => {
  estimateHours.value = Math.floor(estimateMinutes.value / 60);
  estimateMins.value = estimateMinutes.value % 60;
  showEstimateDialog.value = true;
};

const handleSaveEstimate = async () => {
  const totalMins = (estimateHours.value || 0) * 60 + (estimateMins.value || 0);
  isSavingEstimate.value = true;
  try {
    await updateTask(props.task.id, {
      estimated_minutes: totalMins || undefined,
    });
    showEstimateDialog.value = false;
    props.onRefresh();
  } catch (err) {
    console.error("Failed to update task estimate", err);
  } finally {
    isSavingEstimate.value = false;
  }
};

const progressPercent = computed(() => {
  if (!estimateMinutes.value) return 0;
  return Math.min(
    Math.round((props.totalLoggedMinutes / estimateMinutes.value) * 100),
    200
  );
});

const progressColor = computed(() => {
  if (progressPercent.value > 100)
    return "bg-red-500 shadow-xs shadow-red-500/50";
  if (progressPercent.value >= 80)
    return "bg-amber-500 shadow-xs shadow-amber-500/50";
  return "bg-emerald-500 shadow-xs shadow-emerald-500/50";
});

const formatMinutes = (minutes: number) => {
  const h = Math.floor(minutes / 60);
  const m = minutes % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
};

const handleStartTimer = async () => {
  try {
    await startTimer(props.task.id);
  } catch (err) {
    console.error("Failed to start timer", err);
  }
};

const handleStopTimer = () => {
  showStopDialog.value = true;
};

const confirmStop = async () => {
  isStopping.value = true;
  try {
    await stopTimer(props.task.id, { description: stopDescription.value });
    showStopDialog.value = false;
    stopDescription.value = "";
    props.onRefresh();
  } catch (err) {
    console.error("Failed to stop timer", err);
  } finally {
    isStopping.value = false;
  }
};

const filteredEntries = computed(() => {
  let list = props.entries;

  if (activeFilter.value === "timer") {
    list = list.filter((e) => !e.is_manual);
  } else if (activeFilter.value === "manual") {
    list = list.filter((e) => e.is_manual);
  }

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase();
    list = list.filter(
      (e) =>
        e.description?.toLowerCase().includes(q) ||
        e.user_full_name?.toLowerCase().includes(q)
    );
  }

  return list;
});
</script>

<template>
  <div class="space-y-6">
    <!-- Header / High-Level Metric -->
    <div
      class="border-border/60 bg-card/60 flex flex-col gap-4 rounded-xl border p-5 shadow-sm"
    >
      <div class="flex items-center justify-between">
        <h3
          class="text-foreground flex items-center gap-2 text-lg font-semibold tracking-tight"
        >
          <Timer class="text-primary h-5 w-5" />
          Time Log
        </h3>
        <span
          v-if="isTimerForThisTask"
          class="flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-3 py-1 text-xs font-medium text-emerald-500"
        >
          <span class="h-2 w-2 animate-ping rounded-full bg-emerald-500" />
          Timer Active
        </span>
      </div>

      <div class="flex flex-col gap-4 md:flex-row md:items-center">
        <!-- Stats -->
        <div class="flex flex-1 items-center gap-6 text-sm">
          <div
            class="group relative cursor-pointer"
            title="Click to edit estimate"
            @click="openEstimateModal"
          >
            <div class="flex items-center gap-2">
              <span
                class="text-muted-foreground text-xs font-medium tracking-wider uppercase"
                >Estimated</span
              >
              <Pencil
                class="text-muted-foreground/40 group-hover:text-primary h-3 w-3 transition-colors"
              />
            </div>
            <span
              class="text-foreground mt-1 block font-mono text-xl font-bold"
            >
              {{ estimateMinutes ? formatMinutes(estimateMinutes) : "Not Set" }}
            </span>
          </div>

          <div class="bg-border/60 h-10 w-px" />

          <div>
            <span
              class="text-muted-foreground text-xs font-medium tracking-wider uppercase"
              >Logged</span
            >
            <span class="text-primary mt-1 block font-mono text-xl font-bold">
              {{ formatMinutes(totalLoggedMinutes) }}
            </span>
          </div>
        </div>

        <!-- Progress Bar (Full Width within flex item) -->
        <div v-if="estimateMinutes" class="flex-1 space-y-2">
          <div
            class="bg-muted/60 border-border/40 h-2.5 w-full overflow-hidden rounded-full border p-0.5"
          >
            <div
              class="h-full rounded-full transition-all duration-500"
              :class="progressColor"
              :style="{ width: `${Math.min(progressPercent, 100)}%` }"
            />
          </div>
          <div
            class="text-muted-foreground flex justify-between text-xs font-medium"
          >
            <span>{{ progressPercent }}% of estimate logged</span>
            <span
              v-if="progressPercent > 100"
              class="flex items-center gap-1 font-semibold text-red-500"
            >
              <AlertCircle class="h-3.5 w-3.5" />
              +{{ formatMinutes(totalLoggedMinutes - estimateMinutes) }} over
              budget
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Toolbar & Search -->
    <div
      class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div class="flex items-center gap-2">
        <button
          v-if="isTimerForThisTask"
          class="flex items-center gap-1.5 rounded-lg border border-red-500/20 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-500 transition-all hover:bg-red-500/20 active:scale-[0.98]"
          @click="handleStopTimer"
        >
          <Square class="h-4 w-4 fill-current" />
          Stop Timer
        </button>
        <button
          v-else
          class="bg-primary text-primary-foreground hover:bg-primary/90 flex items-center gap-1.5 rounded-lg px-4 py-2 text-sm font-semibold transition-all active:scale-[0.98] disabled:opacity-50"
          :disabled="!!timerStore.activeTimer && !isTimerForThisTask"
          @click="handleStartTimer"
        >
          <Play class="h-4 w-4 fill-current" />
          Start Timer
        </button>

        <button
          class="border-border/60 bg-muted/40 hover:bg-muted text-foreground flex items-center gap-1.5 rounded-lg border px-4 py-2 text-sm font-semibold transition-all active:scale-[0.98]"
          @click="showLogDialog = true"
        >
          <Plus class="h-4 w-4" />
          Log Time
        </button>

        <button
          class="text-muted-foreground hover:text-foreground border-border/60 hover:bg-muted flex items-center gap-1.5 rounded-lg border px-4 py-2 text-sm font-medium transition-all active:scale-[0.98]"
          @click="openEstimateModal"
        >
          <Pencil class="h-4 w-4" />
          Set Estimate
        </button>
      </div>

      <div class="flex items-center gap-3">
        <div class="relative">
          <Search
            class="text-muted-foreground absolute top-1/2 left-2.5 h-4 w-4 -translate-y-1/2"
          />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search logs..."
            class="border-border bg-background focus:border-primary focus:ring-primary w-full rounded-md border py-1.5 pr-3 pl-9 text-sm transition-colors outline-none focus:ring-1 sm:w-48"
          />
        </div>
        <div class="border-border bg-muted/30 flex rounded-md border p-1">
          <button
            v-for="filter in ['all', 'timer', 'manual'] as const"
            :key="filter"
            class="rounded px-2.5 py-1 text-xs font-medium capitalize transition-colors"
            :class="
              activeFilter === filter
                ? 'bg-background text-foreground shadow-xs'
                : 'text-muted-foreground hover:text-foreground'
            "
            @click="activeFilter = filter"
          >
            {{ filter }}
          </button>
        </div>
      </div>
    </div>

    <!-- Timeline Feed -->
    <div class="pt-2">
      <TasksTimeLogList :entries="filteredEntries" :on-refresh="onRefresh" />
    </div>

    <!-- Edit Estimate Dialog -->
    <UiBaseDialog
      title="Set Time Estimate"
      description="Define the expected duration for this task."
      :is-open="showEstimateDialog"
      @close="showEstimateDialog = false"
    >
      <form class="space-y-4" @submit.prevent="handleSaveEstimate">
        <div class="space-y-2">
          <UiBaseLabel>Estimated Time</UiBaseLabel>
          <div class="flex items-center gap-3">
            <div class="flex-1">
              <UiBaseInput
                v-model.number="estimateHours"
                type="number"
                :min="0"
                :max="999"
                placeholder="0"
              />
              <span class="text-muted-foreground mt-1 block text-xs"
                >Hours</span
              >
            </div>
            <div class="flex-1">
              <UiBaseInput
                v-model.number="estimateMins"
                type="number"
                :min="0"
                :max="59"
                placeholder="0"
              />
              <span class="text-muted-foreground mt-1 block text-xs"
                >Minutes</span
              >
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <UiBaseButton
            type="button"
            variant="ghost"
            :disabled="isSavingEstimate"
            @click="showEstimateDialog = false"
          >
            Cancel
          </UiBaseButton>
          <UiBaseButton type="submit" :disabled="isSavingEstimate">
            {{ isSavingEstimate ? "Saving..." : "Save Estimate" }}
          </UiBaseButton>
        </div>
      </form>
    </UiBaseDialog>

    <!-- Stop Timer Dialog -->
    <UiBaseDialog
      title="Log Work"
      description="Record what you worked on during this session."
      :is-open="showStopDialog"
      @close="showStopDialog = false"
    >
      <form class="space-y-4" @submit.prevent="confirmStop">
        <div class="text-center">
          <span class="text-primary font-mono text-3xl font-bold">{{
            timerStore.formattedTime
          }}</span>
        </div>
        <div class="space-y-2">
          <UiBaseLabel for="task-stop-desc">Work Description</UiBaseLabel>
          <UiBaseTextArea
            id="task-stop-desc"
            v-model="stopDescription"
            placeholder="What did you work on?"
            :rows="3"
            :disabled="isStopping"
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <UiBaseButton
            type="button"
            variant="ghost"
            :disabled="isStopping"
            @click="showStopDialog = false"
          >
            Cancel
          </UiBaseButton>
          <UiBaseButton type="submit" :disabled="isStopping">
            {{ isStopping ? "Logging..." : "Log Time" }}
          </UiBaseButton>
        </div>
      </form>
    </UiBaseDialog>

    <!-- Log Time Dialog -->
    <TasksLogTimeDialog
      :open="showLogDialog"
      :task-id="task.id"
      @close="showLogDialog = false"
      @logged="
        () => {
          showLogDialog = false;
          onRefresh();
        }
      "
    />
  </div>
</template>
