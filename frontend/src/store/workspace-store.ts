import { create } from "zustand";
import { persist } from "zustand/middleware";
import { Workspace } from "@/types";

interface WorkspaceState {
  activeWorkspaceId: number | null;
  setActiveWorkspaceId: (id: number) => void;
}

export const useWorkspaceStore = create<WorkspaceState>()(
  persist(
    (set) => ({
      activeWorkspaceId: null,
      setActiveWorkspaceId: (id) => set({ activeWorkspaceId: id }),
    }),
    {
      name: "workspace-storage",
    }
  )
);
