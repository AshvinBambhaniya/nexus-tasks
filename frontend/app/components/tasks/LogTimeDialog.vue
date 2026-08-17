<script setup lang="ts">
const props = defineProps<{
  open: boolean;
  taskId: string;
}>();

const emit = defineEmits<{
  close: [];
  logged: [];
}>();

const { logManualTime } = useTimeTracking();

const hours = ref(0);
const minutes = ref(30);
const description = ref("");
const date = ref(new Date().toISOString().split("T")[0]);
const isSubmitting = ref(false);

const handleSubmit = async () => {
  const totalMinutes = hours.value * 60 + minutes.value;
  if (totalMinutes < 1) return;
  isSubmitting.value = true;
  try {
    await logManualTime(props.taskId, {
      duration_minutes: totalMinutes,
      description: description.value,
      date: date.value,
    });
    hours.value = 0;
    minutes.value = 30;
    description.value = "";
    emit("logged");
  } catch (err) {
    console.error("Failed to log time", err);
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<template>
  <UiBaseDialog
    title="Log Time Manually"
    description="Record time you've already spent on this task."
    :is-open="open"
    @close="$emit('close')"
  >
    <form class="space-y-4" @submit.prevent="handleSubmit">
      <div class="space-y-2">
        <UiBaseLabel for="log-date">Date</UiBaseLabel>
        <UiBaseInput
          id="log-date"
          v-model="date"
          type="date"
          :disabled="isSubmitting"
        />
      </div>
      <div class="space-y-2">
        <UiBaseLabel>Duration</UiBaseLabel>
        <div class="flex items-center gap-2">
          <UiBaseInput
            v-model.number="hours"
            type="number"
            :min="0"
            :max="99"
            class="w-20"
            :disabled="isSubmitting"
          />
          <span class="text-muted-foreground text-sm">hours</span>
          <UiBaseInput
            v-model.number="minutes"
            type="number"
            :min="0"
            :max="59"
            class="w-20"
            :disabled="isSubmitting"
          />
          <span class="text-muted-foreground text-sm">mins</span>
        </div>
      </div>
      <div class="space-y-2">
        <UiBaseLabel for="manual-log-desc">Work Description</UiBaseLabel>
        <UiBaseTextArea
          id="manual-log-desc"
          v-model="description"
          placeholder="What did you work on?"
          :rows="3"
          :disabled="isSubmitting"
        />
      </div>
      <div class="flex justify-end gap-2 pt-4">
        <UiBaseButton
          type="button"
          variant="ghost"
          :disabled="isSubmitting"
          @click="$emit('close')"
          >Cancel</UiBaseButton
        >
        <UiBaseButton
          type="submit"
          :disabled="isSubmitting || hours * 60 + minutes < 1"
        >
          {{ isSubmitting ? "Logging..." : "Log Time" }}
        </UiBaseButton>
      </div>
    </form>
  </UiBaseDialog>
</template>
