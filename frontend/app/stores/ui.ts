import { defineStore } from "pinia";

export const useUIStore = defineStore("ui", {
  state: () => ({
    isSidebarCollapsed: false,
  }),
  actions: {
    toggleSidebar() {
      this.isSidebarCollapsed = !this.isSidebarCollapsed;
    },
    setSidebarCollapsed(value: boolean) {
      this.isSidebarCollapsed = value;
    },
  },
  persist: true,
});
