"use client";

import { useState } from "react";
import { useTasks } from "@/hooks/use-tasks";
import { TaskListView } from "@/components/tasks/task-list-view";
import { TaskBoardView } from "@/components/tasks/task-board-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { Button } from "@/components/ui/button";
import { LayoutList, Kanban, Plus } from "lucide-react";
import { Task, TaskStatus } from "@/types";

type ViewType = "list" | "board";

export default function TasksPage() {
  const { tasks, isLoading, updateTask } = useTasks();
  const [view, setView] = useState<ViewType>("list");
  
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [selectedTask, setSelectedTask] = useState<Task | undefined>(undefined);

  const handleTaskMove = (taskId: number, newStatus: TaskStatus) => {
    updateTask(taskId, { status: newStatus });
  };

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

  return (
    <div className="flex flex-col h-full space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">Tasks</h1>
        <div className="flex items-center gap-2">
            <div className="flex items-center rounded-lg border border-gray-200 bg-white p-1">
                <button
                    onClick={() => setView("list")}
                    className={`rounded-md p-1.5 transition-colors ${view === "list" ? "bg-gray-100 text-gray-900" : "text-gray-500 hover:text-gray-900"}`}
                >
                    <LayoutList className="h-4 w-4" />
                </button>
                <button
                    onClick={() => setView("board")}
                    className={`rounded-md p-1.5 transition-colors ${view === "board" ? "bg-gray-100 text-gray-900" : "text-gray-500 hover:text-gray-900"}`}
                >
                    <Kanban className="h-4 w-4" />
                </button>
            </div>
            <Button onClick={handleCreateClick}>
                <Plus className="mr-2 h-4 w-4" /> Create Task
            </Button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      ) : (
        <div className="flex-1 overflow-hidden">
            {view === "list" ? (
                <TaskListView tasks={tasks} onTaskClick={handleTaskClick} />
            ) : (
                <TaskBoardView tasks={tasks} onTaskMove={handleTaskMove} onTaskClick={handleTaskClick} />
            )}
        </div>
      )}

      <TaskDialog 
        isOpen={isDialogOpen} 
        onClose={handleCloseDialog} 
        task={selectedTask}
      />
    </div>
  );
}
