import { defineStore } from "pinia";

export const useWorkspaceStore = defineStore("workspace", {
  state: () => ({
    activeWorkspaceId: null as number | null,
  }),
  actions: {
    setActiveWorkspaceId(id: number) {
      this.activeWorkspaceId = id;
    },
  },
  persist: true,
});
