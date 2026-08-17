<script setup lang="ts">
import { Timer, Square, X } from "lucide-vue-next";

const timerStore = useTimerStore();
const { stopTimer, discardTimer } = useTimeTracking();

const showStopDialog = ref(false);
const stopDescription = ref("");
const isProcessing = ref(false);

const handleStop = async () => {
  if (!timerStore.activeTimer) return;
  showStopDialog.value = true;
};

const confirmStop = async () => {
  if (!timerStore.activeTimer) return;
  isProcessing.value = true;
  try {
    await stopTimer(timerStore.activeTimer.task_id, {
      description: stopDescription.value,
    });
    showStopDialog.value = false;
    stopDescription.value = "";
  } catch (err) {
    console.error("Failed to stop timer", err);
  } finally {
    isProcessing.value = false;
  }
};

const handleDiscard = async () => {
  if (!timerStore.activeTimer) return;
  if (!confirm("Discard this timer without logging time?")) return;
  try {
    await discardTimer(timerStore.activeTimer.task_id);
  } catch (err) {
    console.error("Failed to discard timer", err);
  }
};

// Restore ticking on mount (e.g. page refresh)
onMounted(() => {
  if (timerStore.activeTimer) {
    timerStore.startTicking();
  }
});
</script>

<template>
  <div
    v-if="timerStore.activeTimer"
    class="border-border/50 bg-primary/5 flex items-center justify-between border-b px-6 py-2.5"
  >
    <div class="flex items-center gap-3">
      <div class="relative flex items-center">
        <Timer class="text-primary h-4 w-4" />
        <span
          class="bg-primary absolute -top-0.5 -right-0.5 h-2 w-2 animate-pulse rounded-full"
        />
      </div>
      <NuxtLink
        :to="`/projects/${timerStore.activeTimer.task_id}`"
        class="text-foreground text-sm font-medium hover:underline"
      >
        TASK-{{ timerStore.activeTimer.task_number }}:
        {{ timerStore.activeTimer.task_title }}
      </NuxtLink>
    </div>

    <div class="flex items-center gap-4">
      <span class="text-primary font-mono text-lg font-bold tracking-wider">
        {{ timerStore.formattedTime }}
      </span>
      <button
        class="flex items-center gap-1.5 rounded-md bg-red-600/10 px-3 py-1.5 text-sm font-medium text-red-600 transition-colors hover:bg-red-600/20"
        @click="handleStop"
      >
        <Square class="h-3.5 w-3.5" />
        Stop
      </button>
      <button
        class="text-muted-foreground hover:text-foreground rounded-md p-1.5 transition-colors hover:bg-white/5"
        title="Discard timer"
        @click="handleDiscard"
      >
        <X class="h-4 w-4" />
      </button>
    </div>

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
          <UiBaseLabel for="stop-desc">Work Description</UiBaseLabel>
          <UiBaseTextArea
            id="stop-desc"
            v-model="stopDescription"
            placeholder="What did you work on?"
            :rows="3"
            :disabled="isProcessing"
          />
        </div>
        <div class="flex justify-end gap-2 pt-4">
          <UiBaseButton
            type="button"
            variant="ghost"
            :disabled="isProcessing"
            @click="showStopDialog = false"
          >
            Cancel
          </UiBaseButton>
          <UiBaseButton type="submit" :disabled="isProcessing">
            {{ isProcessing ? "Logging..." : "Log Time" }}
          </UiBaseButton>
        </div>
      </form>
    </UiBaseDialog>
  </div>
</template>
