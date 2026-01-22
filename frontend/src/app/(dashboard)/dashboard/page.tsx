"use client";

import { useWorkspaceStore } from "@/store/workspace-store";

export default function DashboardPage() {
  const { activeWorkspaceId } = useWorkspaceStore();
  
  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
      <p className="mt-2 text-gray-600">
        Welcome to your workspace. Active Workspace ID: {activeWorkspaceId}
      </p>
      <div className="mt-8 rounded-lg border-2 border-dashed border-gray-200 p-12 text-center">
        <p className="text-gray-500">Select a workspace from the sidebar to get started.</p>
      </div>
    </div>
  );
}
