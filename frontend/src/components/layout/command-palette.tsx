"use client";

import * as React from "react";
import { useRouter } from "next/navigation";
import { Command } from "cmdk";
import {
  Inbox,
  Plus,
  LayoutDashboard,
  Search,
  CheckSquare,
  Building2,
} from "lucide-react";
import { useWorkspaces } from "@/hooks/use-workspaces";
import { useWorkspaceStore } from "@/store/workspace-store";
import { TaskDialog } from "@/components/tasks/task-dialog";

export function CommandPalette() {
  const [open, setOpen] = React.useState(false);
  const [showTaskDialog, setShowTaskDialog] = React.useState(false);
  const router = useRouter();
  const { workspaces } = useWorkspaces();
  const { setActiveWorkspaceId } = useWorkspaceStore();

  // Toggle the menu when ⌘K is pressed
  React.useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        setOpen((open) => !open);
      }
      if (e.key === "Escape") {
        setOpen(false);
      }
    };

    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, []);

  const runCommand = React.useCallback((command: () => void) => {
    setOpen(false);
    command();
  }, []);

  if (!open)
    return (
      <TaskDialog
        isOpen={showTaskDialog}
        onClose={() => setShowTaskDialog(false)}
      />
    );

  return (
    <>
      <div
        className="animate-in fade-in fixed inset-0 z-50 flex items-start justify-center bg-black/50 p-4 pt-[20vh] backdrop-blur-sm duration-200"
        onClick={() => setOpen(false)}
      >
        <Command
          className="animate-in zoom-in-95 w-full max-w-[640px] overflow-hidden rounded-xl border border-gray-200 bg-white shadow-2xl duration-200"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="flex items-center border-b px-3">
            <Search className="mr-2 h-4 w-4 shrink-0 text-gray-400" />
            <Command.Input
              placeholder="Type a command or search..."
              className="flex h-12 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-gray-400 disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
          <Command.List className="max-h-[300px] overflow-y-auto p-2">
            <Command.Empty className="py-6 text-center text-sm text-gray-500">
              No results found.
            </Command.Empty>

            <Command.Group
              heading="Navigation"
              className="px-2 py-1.5 text-xs font-semibold tracking-wider text-gray-500 uppercase"
            >
              <Command.Item
                onSelect={() => runCommand(() => router.push("/inbox"))}
                className="flex cursor-pointer items-center rounded-md px-2 py-2 text-sm text-gray-700 outline-none aria-selected:bg-blue-50 aria-selected:text-blue-600"
              >
                <Inbox className="mr-3 h-4 w-4" />
                Go to Inbox
              </Command.Item>
              <Command.Item
                onSelect={() => runCommand(() => router.push("/tasks"))}
                className="flex cursor-pointer items-center rounded-md px-2 py-2 text-sm text-gray-700 outline-none aria-selected:bg-blue-50 aria-selected:text-blue-600"
              >
                <CheckSquare className="mr-3 h-4 w-4" />
                All Tasks
              </Command.Item>
              <Command.Item
                onSelect={() => runCommand(() => router.push("/dashboard"))}
                className="flex cursor-pointer items-center rounded-md px-2 py-2 text-sm text-gray-700 outline-none aria-selected:bg-blue-50 aria-selected:text-blue-600"
              >
                <LayoutDashboard className="mr-3 h-4 w-4" />
                Dashboard
              </Command.Item>
            </Command.Group>

            <Command.Group
              heading="Actions"
              className="mt-2 px-2 py-1.5 text-xs font-semibold tracking-wider text-gray-500 uppercase"
            >
              <Command.Item
                onSelect={() => runCommand(() => setShowTaskDialog(true))}
                className="flex cursor-pointer items-center rounded-md px-2 py-2 text-sm text-gray-700 outline-none aria-selected:bg-blue-50 aria-selected:text-blue-600"
              >
                <Plus className="mr-3 h-4 w-4" />
                Create New Task
              </Command.Item>
            </Command.Group>

            <Command.Group
              heading="Workspaces"
              className="mt-2 px-2 py-1.5 text-xs font-semibold tracking-wider text-gray-500 uppercase"
            >
              {workspaces.map((workspace) => (
                <Command.Item
                  key={workspace.id}
                  onSelect={() =>
                    runCommand(() => {
                      setActiveWorkspaceId(workspace.id);
                      router.push("/dashboard");
                    })
                  }
                  className="flex cursor-pointer items-center rounded-md px-2 py-2 text-sm text-gray-700 outline-none aria-selected:bg-blue-50 aria-selected:text-blue-600"
                >
                  <Building2 className="mr-3 h-4 w-4" />
                  Switch to {workspace.name}
                </Command.Item>
              ))}
            </Command.Group>
          </Command.List>

          <div className="flex items-center justify-end border-t bg-gray-50 px-3 py-2">
            <p className="text-[10px] text-gray-400">
              Press <kbd className="font-sans text-gray-500">esc</kbd> to close
            </p>
          </div>
        </Command>
      </div>

      <TaskDialog
        isOpen={showTaskDialog}
        onClose={() => setShowTaskDialog(false)}
      />
    </>
  );
}
