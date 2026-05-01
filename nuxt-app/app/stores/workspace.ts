import { defineStore } from "pinia";

export const useWorkspaceStore = defineStore("workspace", {
  state: () => ({
    activeWorkspaceId: null as string | null,
  }),
  actions: {
    setActiveWorkspaceId(id: string) {
      this.activeWorkspaceId = id;
    },
  },
  persist: true,
});
