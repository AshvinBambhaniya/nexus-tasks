"use client";

import { useState, useRef, useEffect } from "react";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { useWorkspaceStore } from "@/store/workspace-store";
import { cn } from "@/lib/utils";
import { ChevronsUpDown, Check, Building2, User } from "lucide-react";
import { WorkspaceType } from "@/types";
import { WorkspaceDialog } from "./workspace-dialog";

export function WorkspaceSwitcher() {
  const { workspaces, activeWorkspace, isLoading } = useWorkspaces();
  const { setActiveWorkspaceId } = useWorkspaceStore();
  const [open, setOpen] = useState(false);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  if (isLoading) {
    return <div className="h-10 w-full animate-pulse rounded bg-gray-200" />;
  }

  return (
    <div className="relative w-full" ref={dropdownRef}>
      <button
        onClick={() => setOpen(!open)}
        className={cn(
          "flex w-full items-center justify-between rounded-lg border border-gray-200 bg-white px-3 py-2 text-sm font-medium transition-colors hover:bg-gray-50",
          open && "ring-2 ring-blue-500 ring-offset-1"
        )}
      >
        <span className="flex items-center gap-2 truncate">
          {activeWorkspace?.type === WorkspaceType.PERSONAL ? (
            <User className="h-4 w-4 text-gray-500" />
          ) : (
            <Building2 className="h-4 w-4 text-gray-500" />
          )}
          <span className="truncate text-gray-900">
            {activeWorkspace?.name || "Select Workspace"}
          </span>
        </span>
        <ChevronsUpDown className="h-4 w-4 shrink-0 text-gray-400" />
      </button>

      {open && (
        <div className="absolute left-0 top-full z-50 mt-1 w-full overflow-hidden rounded-md border border-gray-200 bg-white shadow-lg">
          <div className="max-h-[300px] overflow-auto p-1">
            {workspaces.length === 0 ? (
              <div className="px-2 py-2 text-sm text-gray-500 text-center">No workspaces found.</div>
            ) : (
              workspaces.map((workspace) => (
                <button
                  key={workspace.id}
                  onClick={() => {
                    setActiveWorkspaceId(workspace.id);
                    setOpen(false);
                  }}
                  className={cn(
                    "relative flex w-full cursor-pointer select-none items-center rounded-sm px-2 py-2 text-sm outline-none hover:bg-blue-50 hover:text-blue-600",
                    activeWorkspace?.id === workspace.id ? "bg-blue-50 text-blue-600" : "text-gray-700"
                  )}
                >
                  <div className="mr-2 flex h-4 w-4 items-center justify-center">
                    {workspace.type === WorkspaceType.PERSONAL ? (
                      <User className="h-3 w-3" />
                    ) : (
                      <Building2 className="h-3 w-3" />
                    )}
                  </div>
                  <span className="flex-1 truncate text-left font-medium">{workspace.name}</span>
                  {activeWorkspace?.id === workspace.id && (
                    <Check className="ml-auto h-4 w-4 text-blue-600" />
                  )}
                </button>
              ))
            )}
          </div>
          <div className="border-t border-gray-100 p-1">
            <button
              onClick={() => {
                setOpen(false);
                setIsDialogOpen(true);
              }}
              className="flex w-full cursor-pointer items-center rounded-sm px-2 py-2 text-sm text-gray-500 hover:bg-gray-100"
            >
              <span className="mr-2 h-4 w-4 flex items-center justify-center text-lg leading-none">+</span>
              Create Workspace
            </button>
          </div>
        </div>
      )}

      <WorkspaceDialog 
        isOpen={isDialogOpen} 
        onClose={() => setIsDialogOpen(false)} 
      />
    </div>
  );
}
