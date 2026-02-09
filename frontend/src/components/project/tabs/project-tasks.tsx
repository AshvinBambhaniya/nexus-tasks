"use client";

import { useState } from "react";
import { useTasks } from "@/hooks/use-tasks";
import { TaskListView } from "@/components/tasks/task-list-view";
import { TaskBoardView } from "@/components/tasks/task-board-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { Button } from "@/components/ui/button";
import { Plus, LayoutList, LayoutDashboard, Loader2 } from "lucide-react";
import { Task, TaskStatus } from "@/types";
import { cn } from "@/lib/utils";

interface ProjectTasksProps {
  projectId: number;
}

export function ProjectTasks({ projectId }: ProjectTasksProps) {
  const { tasks, isLoading, updateTask } = useTasks(projectId);
  const [viewMode, setViewMode] = useState<"list" | "board">("list");
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | undefined>(undefined);

  const handleCreateClick = () => {
    setSelectedTask(undefined);
    setIsDialogOpen(true);
  };

  const handleTaskClick = (task: Task) => {
    setSelectedTask(task);
    setIsDialogOpen(true);
  };

  const handleCloseDialog = () => {
    setIsDialogOpen(false);
    setSelectedTask(undefined);
  };

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
      </div>
    );
  }

  return (
    <div className="flex flex-col h-full space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between">
        <div className="flex items-center rounded-lg border border-gray-200 bg-white p-1">
          <button
            onClick={() => setViewMode("list")}
            className={cn(
              "flex items-center gap-2 rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
              viewMode === "list"
                ? "bg-blue-50 text-blue-600"
                : "text-gray-500 hover:bg-gray-50 hover:text-gray-900"
            )}
          >
            <LayoutList className="h-4 w-4" />
            List
          </button>
          <button
            onClick={() => setViewMode("board")}
            className={cn(
              "flex items-center gap-2 rounded-md px-3 py-1.5 text-sm font-medium transition-colors",
              viewMode === "board"
                ? "bg-blue-50 text-blue-600"
                : "text-gray-500 hover:bg-gray-50 hover:text-gray-900"
            )}
          >
            <LayoutDashboard className="h-4 w-4" />
            Board
          </button>
        </div>
        <Button onClick={handleCreateClick}>
          <Plus className="mr-2 h-4 w-4" /> Create Task
        </Button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-hidden">
        {viewMode === "list" ? (
          <TaskListView tasks={tasks} onTaskClick={handleTaskClick} />
        ) : (
          <TaskBoardView
            tasks={tasks}
            onTaskMove={handleTaskMove}
            onTaskClick={handleTaskClick}
          />
        )}
      </div>

      <TaskDialog
        isOpen={isDialogOpen}
        onClose={handleCloseDialog}
        task={selectedTask}
        projectId={projectId}
      />
    </div>
  );
}
