import type { ActiveTimer } from "~/types";

export const useTimerStore = defineStore(
  "timer",
  () => {
    const activeTimer = ref<ActiveTimer | null>(null);
    const elapsedSeconds = ref(0);
    let intervalId: ReturnType<typeof setInterval> | null = null;

    function setActiveTimer(timer: ActiveTimer) {
      activeTimer.value = timer;
      startTicking();
    }

    function clearActiveTimer() {
      activeTimer.value = null;
      stopTicking();
    }

    function startTicking() {
      stopTicking();
      if (!activeTimer.value) return;
      intervalId = setInterval(() => {
        if (activeTimer.value) {
          const start = new Date(activeTimer.value.start_time).getTime();
          elapsedSeconds.value = Math.floor((Date.now() - start) / 1000);
        }
      }, 1000);
    }

    function stopTicking() {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = null;
      }
      elapsedSeconds.value = 0;
    }

    const formattedTime = computed(() => {
      const total = elapsedSeconds.value;
      const h = Math.floor(total / 3600);
      const m = Math.floor((total % 3600) / 60);
      const s = total % 60;
      return `${String(h).padStart(2, "0")}:${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
    });

    return {
      activeTimer,
      elapsedSeconds,
      formattedTime,
      setActiveTimer,
      clearActiveTimer,
      startTicking,
      stopTicking,
    };
  },
  { persist: true }
);
