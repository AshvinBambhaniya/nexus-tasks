<script setup lang="ts">
import {
  Timer,
  Plus,
  Play,
  Square,
  AlertCircle,
  Pencil,
} from "lucide-vue-next";
import type { Task } from "~/types";

const props = defineProps<{
  task: Task;
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
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h3
        class="text-foreground flex items-center gap-2 text-sm font-semibold tracking-tight"
      >
        <Timer class="text-primary h-4 w-4" />
        Time Tracking
      </h3>
      <span
        v-if="isTimerForThisTask"
        class="flex items-center gap-1.5 rounded-full border border-emerald-500/20 bg-emerald-500/10 px-2 py-0.5 text-[10px] font-medium text-emerald-500"
      >
        <span class="h-1.5 w-1.5 animate-ping rounded-full bg-emerald-500" />
        Timer Active
      </span>
    </div>

    <div
      class="border-border/60 bg-card/60 space-y-4 rounded-xl border p-4 shadow-xs"
    >
      <!-- Estimate & Logged Metrics -->
      <div class="grid grid-cols-2 gap-3 text-xs">
        <div
          class="bg-muted/30 border-border/40 hover:border-border group relative cursor-pointer rounded-lg border p-2.5 transition-all"
          title="Click to edit estimate"
          @click="openEstimateModal"
        >
          <div class="flex items-center justify-between">
            <span
              class="text-muted-foreground block text-[11px] font-medium tracking-wider uppercase"
              >Estimated</span
            >
            <Pencil
              class="text-muted-foreground/40 group-hover:text-primary h-3 w-3 transition-colors"
            />
          </div>
          <span
            class="text-foreground mt-0.5 block font-mono text-sm font-bold"
          >
            {{ estimateMinutes ? formatMinutes(estimateMinutes) : "Set time" }}
          </span>
        </div>
        <div class="bg-muted/30 border-border/40 rounded-lg border p-2.5">
          <span
            class="text-muted-foreground block text-[11px] font-medium tracking-wider uppercase"
            >Logged</span
          >
          <span class="text-primary mt-0.5 block font-mono text-sm font-bold">
            {{ formatMinutes(totalLoggedMinutes) }}
          </span>
        </div>
      </div>

      <!-- Progress Bar -->
      <div v-if="estimateMinutes" class="space-y-1.5">
        <div
          class="bg-muted/60 border-border/40 h-2 w-full overflow-hidden rounded-full border p-0.5"
        >
          <div
            class="h-full rounded-full transition-all duration-500"
            :class="progressColor"
            :style="{ width: `${Math.min(progressPercent, 100)}%` }"
          />
        </div>
        <div
          class="text-muted-foreground flex justify-between text-[11px] font-medium"
        >
          <span>{{ progressPercent }}% of estimate</span>
          <span
            v-if="progressPercent > 100"
            class="flex items-center gap-1 font-semibold text-red-500"
          >
            <AlertCircle class="h-3 w-3" />
            +{{ formatMinutes(totalLoggedMinutes - estimateMinutes) }} over
          </span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex gap-2 pt-1">
        <button
          v-if="isTimerForThisTask"
          class="flex flex-1 items-center justify-center gap-1.5 rounded-lg border border-red-500/20 bg-red-500/10 px-3 py-2 text-xs font-semibold text-red-500 transition-all hover:bg-red-500/20 active:scale-[0.98]"
          @click="handleStopTimer"
        >
          <Square class="h-3.5 w-3.5 fill-current" />
          Stop Timer
        </button>
        <button
          v-else
          class="bg-primary/10 text-primary hover:bg-primary/20 border-primary/20 flex flex-1 items-center justify-center gap-1.5 rounded-lg border px-3 py-2 text-xs font-semibold transition-all active:scale-[0.98] disabled:opacity-50"
          :disabled="!!timerStore.activeTimer && !isTimerForThisTask"
          @click="handleStartTimer"
        >
          <Play class="h-3.5 w-3.5 fill-current" />
          Start Timer
        </button>
        <button
          class="border-border/60 bg-muted/40 hover:bg-muted text-foreground flex items-center gap-1.5 rounded-lg border px-3 py-2 text-xs font-semibold transition-all active:scale-[0.98]"
          @click="showLogDialog = true"
        >
          <Plus class="h-3.5 w-3.5" />
          Log
        </button>
      </div>
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
          <UiBaseLabel for="sidebar-stop-desc">Work Description</UiBaseLabel>
          <UiBaseTextArea
            id="sidebar-stop-desc"
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
