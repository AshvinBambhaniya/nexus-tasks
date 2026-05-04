import { defineStore } from "pinia";

export const useProjectStore = defineStore(
  "project",
  () => {
    const activeProjectId = ref<string | null>(null);

    const setActiveProjectId = (id: string | null) => {
      activeProjectId.value = id;
    };

    return { activeProjectId, setActiveProjectId };
  },
  {
    persist: true,
  }
);
