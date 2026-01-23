"use client";

import { useState } from "react";
import { useTasks } from "@/hooks/use-tasks";
import { TaskListView } from "@/components/tasks/task-list-view";
import { TaskDialog } from "@/components/tasks/task-dialog";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { Task } from "@/types";

export default function TasksPage() {
  const { tasks, isLoading } = useTasks();
  
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

  return (
    <div className="flex flex-col h-full space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900">All Tasks</h1>
        <Button onClick={handleCreateClick}>
            <Plus className="mr-2 h-4 w-4" /> Create Task
        </Button>
      </div>

      {isLoading ? (
        <div className="flex h-64 items-center justify-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
        </div>
      ) : (
        <div className="flex-1 overflow-hidden">
            <TaskListView tasks={tasks} onTaskClick={handleTaskClick} />
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
