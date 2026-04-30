import { defineStore } from "pinia";

export const useProjectStore = defineStore(
  "project",
  () => {
    const activeProjectId = ref<number | null>(null);

    const setActiveProjectId = (id: number | null) => {
      activeProjectId.value = id;
    };

    return { activeProjectId, setActiveProjectId };
  },
  {
    persist: true,
  }
);
